package api

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/alexsasharegan/akmctl/internal/frontmatter"
	"github.com/alexsasharegan/akmctl/internal/fsx"
	"github.com/alexsasharegan/akmctl/internal/initrepo"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// settingDefs lists the settings the config operation knows.
var settingDefs = []struct {
	key, desc string
	paths     []string
}{
	{"git_commits", "Install the akm-commit skill, which tells an agent to commit after each logical step. Off: the skill is not installed.", []string{initrepo.SkillDir(initrepo.CommitSkill)}},
}

func settingsObj(repo *record.Repo, only string) []Obj {
	out := []Obj{}
	for _, d := range settingDefs {
		if only != "" && d.key != only {
			continue
		}
		var v any
		switch d.key {
		case "git_commits":
			v = repo.GitCommits()
		}
		out = append(out, O("key", d.key, "value", v, "description", d.desc, "paths", d.paths))
	}
	return out
}

// opConfig lists the settings (no key), reads one (key) or changes one (key and
// value). Changing git_commits writes the manifest key and installs or removes
// the optional skill akm-commit.
func opConfig(repo *record.Repo, q Request) (Obj, *Error) {
	_, hasKey := q["key"]
	_, hasVal := q["value"]
	if hasVal && !hasKey {
		return nil, invalid("key", "value needs a key")
	}
	if !hasVal {
		for _, p := range []string{"dry_run", "confirm", "confirm_edited"} {
			if _, given := q[p]; given {
				return nil, invalid(p, "%s needs a value to set", p)
			}
		}
		return O("settings", settingsObj(repo, q.str("key"))), nil
	}
	return setGitCommits(repo, q, q["value"].(bool))
}

// commitPlan is what turning git_commits on or off would do.
type commitPlan struct {
	write, remove, modify []string
	files                 []initrepo.File // contents of write
	conflicts             []string
	edits                 []Obj // diffs from the shipped copy to the installed files (removal)
}

func (p *commitPlan) obj() Obj {
	return O("write", nonNil(p.write), "remove", nonNil(p.remove), "modify", nonNil(p.modify))
}

func (p *commitPlan) editsObj() []Obj {
	if p.edits == nil {
		return []Obj{}
	}
	return p.edits
}

func setGitCommits(repo *record.Repo, q Request, on bool) (Obj, *Error) {
	if _, err := os.Stat(filepath.Join(repo.Root, filepath.FromSlash(record.ManifestPath))); err != nil {
		return nil, refused("no_manifest", "there is no %s: run akmctl init or akmctl migrate first", record.ManifestPath)
	}
	if on && !record.InGitRepo(repo.Root) {
		return nil, refused("not_a_git_repo", "this repository is not inside a git work tree: run git init first; git commits need one")
	}
	dry := q.boolean("dry_run")
	if !dry {
		unlock, lerr := fsx.LockRepo(repo.Root)
		if lerr != nil {
			return nil, unreadable(lerr)
		}
		defer unlock()
	}
	if lerr := repo.Reload(); lerr != nil {
		return nil, unreadable(lerr)
	}
	plan, err := planGitCommits(repo, on)
	if err != nil {
		return nil, err
	}
	if len(plan.conflicts) > 0 {
		return nil, &Error{Code: "refused", Message: "not overwriting what is not the tool's own: " + joinPaths(plan.conflicts),
			Details: O("reason", "conflict", "paths", plan.conflicts)}
	}
	needsConfirm := len(plan.remove) > 0
	edited := len(plan.edits) > 0
	if !dry && needsConfirm {
		if !q.boolean("confirm") || (edited && !q.boolean("confirm_edited")) {
			msg := "removing " + initrepo.SkillDir(initrepo.CommitSkill) + " needs confirm: true"
			if edited && q.boolean("confirm") {
				msg = "the installed skill differs from the shipped copy (see edits); removing it needs confirm_edited: true"
			}
			return nil, &Error{Code: "refused", Message: msg,
				Details: O("reason", "needs_confirmation", "edited", edited, "plan", plan.obj(), "edits", plan.editsObj())}
		}
	}
	changed := len(plan.write) > 0 || len(plan.remove) > 0 || len(plan.modify) > 0
	if !dry && changed {
		if err := applyGitCommits(repo, plan, on); err != nil {
			return nil, err
		}
		if lerr := repo.Reload(); lerr != nil {
			return nil, unreadable(lerr)
		}
	}
	return O("key", "git_commits", "value", on, "changed", changed, "applied", changed && !dry,
		"needs_confirmation", needsConfirm, "plan", plan.obj(), "edits", plan.editsObj(), "settings", settingsObj(repo, "")), nil
}

func joinPaths(l []string) string {
	s := ""
	for i, p := range l {
		if i > 0 {
			s += ", "
		}
		s += p
	}
	return s
}

// planGitCommits works out the changes from what is on disk. It writes nothing.
func planGitCommits(repo *record.Repo, on bool) (*commitPlan, *Error) {
	root := repo.Root
	skillDir := initrepo.SkillDir(initrepo.CommitSkill)
	abs := func(rel string) string { return filepath.Join(root, filepath.FromSlash(rel)) }
	p := &commitPlan{}
	shipped, serr := initrepo.OptionalSkillFiles(initrepo.CommitSkill, repo.WorkDir(), repo.ADRDir())
	if serr != nil {
		return nil, &Error{Code: "internal", Message: serr.Error()}
	}
	// the key: present with the wanted value, or (off) absent
	if v, ok := repo.Layout["git_commits"]; !ok && on || ok && (v.Kind != frontmatter.Bool || v.Bool != on) {
		p.modify = []string{record.ManifestPath}
	}
	// a parent that is a file or a link is not ours to write through
	for _, parent := range []string{".agents", ".agents/skills"} {
		if st, err := os.Lstat(abs(parent)); err == nil && !st.IsDir() {
			p.conflicts = append(p.conflicts, parent)
		}
	}
	if len(p.conflicts) > 0 {
		return p, nil
	}
	st, err := os.Lstat(abs(skillDir))
	exists := err == nil
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, unreadable(err)
	}
	if exists && !st.IsDir() {
		p.conflicts = []string{skillDir}
		return p, nil
	}
	if on {
		for _, f := range shipped {
			cur, err := os.Lstat(abs(f.Path))
			switch {
			case errors.Is(err, fs.ErrNotExist):
				p.write, p.files = append(p.write, f.Path), append(p.files, f)
			case err != nil || !cur.Mode().IsRegular():
				p.conflicts = append(p.conflicts, f.Path)
			default:
				b, rerr := os.ReadFile(abs(f.Path))
				if rerr != nil {
					return nil, unreadable(rerr)
				}
				if !bytes.Equal(b, f.Data) {
					p.conflicts = append(p.conflicts, f.Path)
				}
			}
		}
		return p, nil
	}
	if !exists {
		return p, nil
	}
	want := map[string]string{}
	for _, f := range shipped {
		want[f.Path] = string(f.Data)
	}
	werr := filepath.WalkDir(abs(skillDir), func(fp string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, fp)
		rel = filepath.ToSlash(rel)
		switch {
		case d.IsDir():
			return nil
		case !d.Type().IsRegular():
			p.conflicts = append(p.conflicts, rel)
			return nil
		}
		b, err := os.ReadFile(fp)
		if err != nil {
			return err
		}
		p.remove = append(p.remove, rel)
		if w, known := want[rel]; !known || w != string(b) {
			p.edits = append(p.edits, Diff(rel, rel, w, string(b)))
		}
		return nil
	})
	if werr != nil {
		return nil, unreadable(werr)
	}
	sort.Strings(p.remove)
	sort.Strings(p.conflicts)
	return p, nil
}

// applyGitCommits does the plan. For on: the skill files, then the key. For
// off: the skill directory, then the key. It runs under the repository lock.
func applyGitCommits(repo *record.Repo, p *commitPlan, on bool) *Error {
	root := repo.Root
	if on {
		for _, f := range p.files {
			if err := fsx.WriteAtomic(filepath.Join(root, filepath.FromSlash(f.Path)), f.Data); err != nil {
				return unreadable(err)
			}
		}
	} else if len(p.remove) > 0 {
		dir := filepath.Join(root, filepath.FromSlash(initrepo.SkillDir(initrepo.CommitSkill)))
		if st, err := os.Lstat(dir); err != nil || !st.IsDir() {
			return refused("conflict", "%s is no longer a directory", initrepo.SkillDir(initrepo.CommitSkill))
		}
		if err := os.RemoveAll(dir); err != nil {
			return unreadable(err)
		}
	}
	if len(p.modify) > 0 {
		full := filepath.Join(root, filepath.FromSlash(record.ManifestPath))
		b, err := os.ReadFile(full)
		if err != nil {
			return unreadable(err)
		}
		doc, err := frontmatter.ParseBare(b)
		if err != nil {
			return unreadable(err)
		}
		if err := doc.Set("git_commits", frontmatter.NewBool(on)); err != nil {
			return unreadable(err)
		}
		if err := fsx.WriteAtomic(full, doc.Bytes()); err != nil {
			return unreadable(err)
		}
	}
	return nil
}
