package api

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/alexsasharegan/akmctl/internal/fsx"
	"github.com/alexsasharegan/akmctl/internal/initrepo"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// opSkills reports, installs and updates the managed skills (ADR 261030-h6t7,
// ADR 261030-wsqb). A write mode returns the plan and writes nothing unless
// confirm is true and dry_run is not.

// skillInfo is what is known about one managed skill after reading the disk.
type skillInfo struct {
	skill    initrepo.Skill
	dir      string // repository-relative folder
	state    string // missing, current, outdated, edited, unmarked, conflict
	files    []Obj  // {path, state}
	stale    []string
	hasMark  bool
	unmarked bool // edited, but only because there is no marker to compare with
}

// planAction is one change a write mode would make.
type planAction struct {
	skill   string
	kind    string // write, replace, marker, remove
	path    string
	data    []byte
	old     []byte
	blocked bool // needs confirm_edited
}

func opSkills(repo *record.Repo, q Request, info Info) (Obj, *Error) {
	mode := q.str("mode")
	if mode == "" {
		mode = "status"
	}
	tool := info.Version
	if tool == "" {
		tool = "dev"
	}
	write := mode != "status"
	if write {
		if _, err := os.Stat(filepath.Join(repo.Root, filepath.FromSlash(repo.ManifestPath()))); err != nil {
			return nil, refused("no_manifest", "there is no %s: run akmctl init or akmctl migrate first", repo.ManifestPath())
		}
	} else {
		for _, p := range []string{"dry_run", "confirm", "confirm_edited"} {
			if _, given := q[p]; given {
				return nil, invalid(p, "%s needs mode install or update", p)
			}
		}
	}
	dry := q.boolean("dry_run")
	if write && !dry && q.boolean("confirm") {
		unlock, lerr := fsx.LockRepo(repo.Root, repo.MetaDir())
		if lerr != nil {
			return nil, unreadable(lerr)
		}
		defer unlock()
	}
	if lerr := repo.Reload(); lerr != nil {
		return nil, unreadable(lerr)
	}
	managed, merr := initrepo.Managed(repo.WorkDir(), repo.ADRDir(), repo.GitCommits())
	if merr != nil {
		return nil, &Error{Code: "internal", Message: merr.Error()}
	}
	abs := func(rel string) string { return filepath.Join(repo.Root, filepath.FromSlash(rel)) }
	parentConflict := false
	for _, parent := range []string{".agents", ".agents/skills"} {
		if st, err := os.Lstat(abs(parent)); err == nil && !st.IsDir() {
			parentConflict = true
		}
	}
	var infos []*skillInfo
	for _, sk := range managed {
		si, err := inspectSkill(repo.Root, sk, parentConflict)
		if err != nil {
			return nil, unreadable(err)
		}
		infos = append(infos, si)
	}
	renamed, rerr := inspectRenamed(repo.Root, managed)
	if rerr != nil {
		return nil, unreadable(rerr)
	}

	reports := make([]Obj, 0, len(infos)+len(renamed))
	clean := true
	for _, si := range infos {
		reports = append(reports, si.obj())
		if si.state != "current" {
			clean = false
		}
	}
	for _, r := range renamed {
		reports = append(reports, r.obj())
		clean = false
	}
	if !write {
		return O("mode", mode, "clean", clean, "skills", reports), nil
	}

	confirmEdited := q.boolean("confirm_edited")
	var actions []planAction
	var skipped []Obj
	skip := func(name, reason string) { skipped = append(skipped, O("skill", name, "reason", reason)) }
	for _, si := range infos {
		repoFiles := si.skill.RepoFiles(tool)
		marker := repoFiles[0]
		for _, f := range repoFiles {
			if strings.HasSuffix(f.Path, "/"+initrepo.MarkerName) {
				marker = f
			}
		}
		switch si.state {
		case "current":
		case "conflict":
			skip(si.skill.Name, "conflict")
		case "missing":
			for _, f := range repoFiles {
				actions = append(actions, planAction{skill: si.skill.Name, kind: "write", path: f.Path, data: f.Data})
			}
		case "unmarked":
			for _, f := range si.skill.Files {
				p := si.dir + "/" + f.Path
				if _, err := os.Lstat(abs(p)); errors.Is(err, fs.ErrNotExist) {
					actions = append(actions, planAction{skill: si.skill.Name, kind: "write", path: p, data: f.Data})
				}
			}
			actions = append(actions, planAction{skill: si.skill.Name, kind: "marker", path: marker.Path, data: marker.Data})
		case "outdated", "edited":
			if mode == "install" {
				skip(si.skill.Name, si.state)
				continue
			}
			blocked := si.state == "edited" && !confirmEdited
			for _, f := range si.skill.Files {
				p := si.dir + "/" + f.Path
				cur, err := os.ReadFile(abs(p))
				switch {
				case errors.Is(err, fs.ErrNotExist):
					actions = append(actions, planAction{skill: si.skill.Name, kind: "write", path: p, data: f.Data, blocked: blocked})
				case err != nil:
					return nil, unreadable(err)
				case !bytes.Equal(cur, f.Data):
					actions = append(actions, planAction{skill: si.skill.Name, kind: "replace", path: p, data: f.Data, old: cur, blocked: blocked})
				}
			}
			actions = append(actions, planAction{skill: si.skill.Name, kind: "marker", path: marker.Path, data: marker.Data, blocked: blocked})
		}
	}
	for _, r := range renamed {
		if mode == "install" {
			skip(r.name, "renamed: run update")
			continue
		}
		actions = append(actions, planAction{skill: r.name, kind: "remove", path: r.dir, blocked: !r.unedited && !confirmEdited})
	}

	planObjs := make([]Obj, 0, len(actions))
	for _, a := range actions {
		o := O("skill", a.skill, "action", a.kind, "path", a.path, "needs_confirm_edited", a.blocked)
		if a.kind == "replace" {
			o = append(o, KV{"diff", Diff(a.path, a.path, string(a.old), string(a.data))})
		}
		planObjs = append(planObjs, o)
	}
	apply := q.boolean("confirm") && !dry
	var written, replaced, removed []string
	if apply {
		blockedSkills := map[string]bool{}
		for _, a := range actions {
			if a.blocked {
				blockedSkills[a.skill] = true
			}
		}
		for _, a := range actions {
			if a.blocked {
				continue
			}
			switch a.kind {
			case "write", "marker":
				if err := fsx.WriteAtomic(abs(a.path), a.data); err != nil {
					return nil, unreadable(err)
				}
				written = append(written, a.path)
			case "replace":
				if err := fsx.WriteAtomic(abs(a.path), a.data); err != nil {
					return nil, unreadable(err)
				}
				replaced = append(replaced, a.path)
			case "remove":
				st, err := os.Lstat(abs(a.path))
				if err != nil || !st.IsDir() || !isRenamedOld(filepath.Base(a.path)) {
					return nil, refused("conflict", "%s is not a skill folder the tool may remove", a.path)
				}
				if err := os.RemoveAll(abs(a.path)); err != nil {
					return nil, unreadable(err)
				}
				removed = append(removed, a.path)
			}
		}
		for name := range blockedSkills {
			skip(name, "edited: needs confirm_edited")
		}
	}
	sort.Strings(written)
	sort.Strings(replaced)
	sort.Strings(removed)
	pending := false
	for _, a := range actions {
		if !a.blocked {
			pending = true
		}
	}
	out := O("mode", mode, "clean", clean, "skills", reports, "plan", planObjs, "applied", apply && (len(written)+len(replaced)+len(removed)) > 0,
		"needs_confirmation", pending && !apply, "written", nonNil(written), "replaced", nonNil(replaced), "removed", nonNil(removed), "skipped", nonNilObj(skipped))
	return out, nil
}

func nonNilObj(l []Obj) []Obj {
	if l == nil {
		return []Obj{}
	}
	return l
}

func isRenamedOld(name string) bool {
	for _, r := range initrepo.Renames {
		if r.Old == name {
			return true
		}
	}
	return false
}

func (si *skillInfo) obj() Obj {
	o := O("name", si.skill.Name, "path", si.dir, "state", si.state, "marker", si.hasMark, "files", si.files)
	if si.unmarked {
		o = append(o, KV{"note", "no marker: the files differ from the shipped copy, so the tool cannot tell an old copy from an edited one"})
	}
	if len(si.stale) > 0 {
		o = append(o, KV{"stale", si.stale})
	}
	return o
}

// inspectSkill reads the disk and decides the state of one managed skill.
func inspectSkill(root string, sk initrepo.Skill, parentConflict bool) (*skillInfo, error) {
	abs := func(rel string) string { return filepath.Join(root, filepath.FromSlash(rel)) }
	si := &skillInfo{skill: sk, dir: initrepo.SkillDir(sk.Name), files: []Obj{}}
	if parentConflict {
		si.state = "conflict"
		return si, nil
	}
	st, err := os.Lstat(abs(si.dir))
	switch {
	case errors.Is(err, fs.ErrNotExist):
		si.state = "missing"
		for _, f := range sk.Files {
			si.files = append(si.files, O("path", f.Path, "state", "missing"))
		}
		return si, nil
	case err != nil:
		return nil, err
	case !st.IsDir():
		si.state = "conflict"
		return si, nil
	}
	disk := map[string][]byte{}
	conflict := false
	for _, f := range sk.Files {
		fi, err := os.Lstat(abs(si.dir + "/" + f.Path))
		switch {
		case errors.Is(err, fs.ErrNotExist):
		case err != nil:
			return nil, err
		case !fi.Mode().IsRegular():
			conflict = true
		default:
			b, rerr := os.ReadFile(abs(si.dir + "/" + f.Path))
			if rerr != nil {
				return nil, rerr
			}
			disk[f.Path] = b
		}
	}
	if conflict {
		si.state = "conflict"
		return si, nil
	}
	var recorded map[string]string
	if b, rerr := os.ReadFile(abs(si.dir + "/" + initrepo.MarkerName)); rerr == nil {
		if _, files, ok := initrepo.ParseMarker(b); ok {
			recorded, si.hasMark = files, true
		}
	}
	shipped := map[string]bool{}
	for _, f := range sk.Files {
		shipped[f.Path] = true
	}
	allEqual, existingEqual := true, true
	for _, f := range sk.Files {
		d, ok := disk[f.Path]
		if !ok {
			allEqual = false
			continue
		}
		if !bytes.Equal(d, f.Data) {
			allEqual, existingEqual = false, false
		}
	}
	// the stale files: recorded, not shipped any more
	for p := range recorded {
		if !shipped[p] {
			si.stale = append(si.stale, p)
		}
	}
	sort.Strings(si.stale)
	switch {
	case allEqual && si.hasMark:
		si.state = "current"
	case existingEqual && !si.hasMark:
		si.state = "unmarked"
	case !si.hasMark:
		si.state, si.unmarked = "edited", true
	default:
		si.state = "outdated"
		for p, h := range recorded {
			if !shipped[p] {
				continue
			}
			if d, ok := disk[p]; !ok || initrepo.Hash(d) != h {
				si.state = "edited"
			}
		}
	}
	for _, f := range sk.Files {
		fs := "current"
		d, ok := disk[f.Path]
		switch {
		case !ok:
			fs = "missing"
		case bytes.Equal(d, f.Data):
		case si.hasMark && recorded[f.Path] == initrepo.Hash(d):
			fs = "outdated"
		default:
			fs = "edited"
		}
		si.files = append(si.files, O("path", f.Path, "state", fs))
	}
	return si, nil
}

// renamedInfo is a skill folder under a name the tool has retired.
type renamedInfo struct {
	name, dir, newName string
	unedited           bool
}

func (r renamedInfo) obj() Obj {
	return O("name", r.name, "path", r.dir, "state", "renamed", "replaced_by", r.newName, "unedited", r.unedited)
}

// inspectRenamed finds folders with a retired name. A folder is unedited when
// its marker is valid and every file in it is a recorded one with its recorded
// hash (the marker aside); a folder with no marker is never unedited.
func inspectRenamed(root string, managed []initrepo.Skill) ([]renamedInfo, error) {
	var out []renamedInfo
	for _, r := range initrepo.Renames {
		dir := initrepo.SkillDir(r.Old)
		full := filepath.Join(root, filepath.FromSlash(dir))
		st, err := os.Lstat(full)
		if errors.Is(err, fs.ErrNotExist) || errors.Is(err, syscall.ENOTDIR) {
			continue // a parent that is a file is reported as a conflict on the managed skills
		}
		if err != nil {
			return nil, err
		}
		if !st.IsDir() {
			continue
		}
		ri := renamedInfo{name: r.Old, dir: dir, newName: r.New}
		if b, rerr := os.ReadFile(filepath.Join(full, initrepo.MarkerName)); rerr == nil {
			if _, files, ok := initrepo.ParseMarker(b); ok {
				ri.unedited = true
				seen := 0
				werr := filepath.WalkDir(full, func(p string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if d.IsDir() || d.Name() == initrepo.MarkerName {
						return nil
					}
					rel, _ := filepath.Rel(full, p)
					want, known := files[filepath.ToSlash(rel)]
					data, rerr := os.ReadFile(p)
					if !d.Type().IsRegular() || rerr != nil || !known || initrepo.Hash(data) != want {
						ri.unedited = false
					}
					seen++
					return nil
				})
				if werr != nil {
					return nil, werr
				}
				if seen != len(files) {
					ri.unedited = false
				}
			}
		}
		out = append(out, ri)
	}
	return out, nil
}
