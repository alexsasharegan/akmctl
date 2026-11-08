package api

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/frontmatter"
	"github.com/alexsasharegan/akmctl/internal/fsx"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// bases writes an Obsidian Bases file that mirrors the workbench: table, cards and
// Kanban views of the work items and a table of the ADRs. It is written once and
// never overwritten, like init, so the maintainer can edit it freely in Obsidian.
// The format is Obsidian's, and the subset used here (filters, formulas,
// properties, and the table, cards and kanban views) is the one Obsidian 1.14
// documents; akmctl cannot open Obsidian, so the maintainer checks the views render.

const (
	basesName          = "akm.base"
	basesDefaultStale  = 30
	basesRecentlyShown = 14 // days a closed item stays in "Recently closed"
)

func opBases(repo *record.Repo, q Request) (Obj, *Error) {
	rel := q.str("path")
	if rel == "" {
		rel = repo.MetaDir() + "/" + basesName
	}
	if err := basesPath(rel); err != nil {
		return nil, invalid("path", "%s", err.Error())
	}
	vault, found := vaultRoot(repo.Root)
	work, perr := vaultFolder(repo.Root, vault, repo.WorkDir())
	if perr != nil {
		return nil, unreadable(perr)
	}
	adr, perr := vaultFolder(repo.Root, vault, repo.ADRDir())
	if perr != nil {
		return nil, unreadable(perr)
	}
	stale := basesDefaultStale
	if v, ok := repo.Layout["stale_days"]; ok && v.Kind == frontmatter.Int && v.Int > 0 {
		stale = int(v.Int)
	}
	content := BasesFile(work, adr, stale)

	full := filepath.Join(repo.Root, filepath.FromSlash(rel))
	if _, err := os.Lstat(full); err == nil {
		return nil, refused("exists", "%s already exists; akmctl never overwrites it. Delete it or pass another path to write a new one", rel)
	}
	written := false
	if !q.boolean("dry_run") {
		unlock, lerr := fsx.LockRepo(repo.Root, repo.MetaDir())
		if lerr != nil {
			return nil, unreadable(lerr)
		}
		defer unlock()
		if _, err := os.Lstat(full); err == nil { // made while the lock was awaited
			return nil, refused("exists", "%s already exists; akmctl never overwrites it. Delete it or pass another path to write a new one", rel)
		}
		if err := fsx.WriteAtomic(full, []byte(content)); err != nil {
			return nil, &Error{Code: "unreadable", Message: "cannot write " + rel + ": " + err.Error()}
		}
		written = true
	}
	vrel, _ := filepath.Rel(repo.Root, vault)
	out := O("path", rel, "written", written, "vault_root", filepath.ToSlash(vrel), "vault_found", found, "work_folder", work, "adr_folder", adr)
	if q.boolean("include_content") {
		out = append(out, KV{"content", content})
	}
	return out, nil
}

// basesPath checks a path for the file: relative to the repository root, with "/",
// inside it, and ending in .base.
func basesPath(p string) error {
	switch {
	case strings.ContainsAny(p, "\\\x00\r\n") || path.IsAbs(p) || strings.HasPrefix(p, "/"):
		return errString("path is relative to the repository root, with / as the separator")
	case path.Clean(p) != p || p == ".." || strings.HasPrefix(p, "../") || strings.Contains(p, "/../"):
		return errString("path must stay inside the repository, with no . or .. segment")
	case !strings.HasSuffix(p, ".base") || len(p) <= len(".base"):
		return errString("path must name a file ending in .base")
	}
	return nil
}

type errString string

func (e errString) Error() string { return string(e) }

// vaultRoot finds the nearest folder at or above the repository root that holds an
// .obsidian folder. With none, the repository root stands in for the vault.
func vaultRoot(repoRoot string) (string, bool) {
	dir := repoRoot
	for {
		if st, err := os.Stat(filepath.Join(dir, ".obsidian")); err == nil && st.IsDir() {
			return dir, true
		}
		up := filepath.Dir(dir)
		if up == dir {
			return repoRoot, false
		}
		dir = up
	}
}

// vaultFolder is a folder of the repository as the vault sees it: relative to the
// vault root, with "/" as the separator.
func vaultFolder(repoRoot, vault, repoRel string) (string, error) {
	rel, err := filepath.Rel(vault, filepath.Join(repoRoot, filepath.FromSlash(repoRel)))
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}

// yamlQuote is a YAML single-quoted scalar.
func yamlQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }

// BasesFile is the file's text. work and adr are folders relative to the vault root
// and stale is the manifest's stale_days.
func BasesFile(work, adr string, stale int) string {
	inWork := "        - " + yamlQuote(`file.inFolder("`+work+`")`) + "\n"
	open := "        - status != \"done\"\n        - status != \"dropped\"\n"
	var b strings.Builder
	w := b.WriteString
	w("# Obsidian Bases file written by `akmctl bases` for this repository. It mirrors the\n")
	w("# akmctl workbench. akmctl wrote it once and never overwrites it: edit it freely.\n")
	w("# Written for Obsidian 1.14 (the Kanban view needs 1.14; table and cards 1.9). akmctl cannot open Obsidian to test it.\n")
	w("formulas:\n")
	w("  age_days: if(updated, (today() - updated).days, \"\")\n")
	w("  stale: if(updated, (today() - updated).days > " + strconv.Itoa(stale) + ", true)\n")
	w("properties:\n  formula.age_days:\n    displayName: days since update\n  formula.stale:\n    displayName: stale\n")
	w("views:\n")
	view := func(kind, name, filters string, rest string) {
		w("  - type: " + kind + "\n    name: " + name + "\n    filters:\n      and:\n" + inWork + filters + rest)
	}
	order := func(props ...string) string {
		s := "    order:\n"
		for _, p := range props {
			s += "      - " + p + "\n"
		}
		return s
	}
	sortBy := func(prop, dir string) string {
		return "      - property: " + prop + "\n        direction: " + dir + "\n"
	}
	group := func(prop string) string { return "    groupBy:\n      property: " + prop + "\n      direction: ASC\n" }

	view("table", "Open by parent", open, order("parent", "file.name", "status", "formula.age_days")+"    sort:\n"+sortBy("parent", "ASC")+sortBy("status", "ASC"))
	view("table", "Ready and active", "        - status == \"ready\" || status == \"active\"\n", order("file.name", "status", "parent", "updated"))
	view("table", "Blocked", "        - status == \"blocked\"\n", order("file.name", "blocked_by", "parent", "updated"))
	view("table", "Stale", open+"        - formula.stale\n", order("file.name", "status", "parent", "formula.age_days"))
	view("table", "Recently closed", "        - status == \"done\" || status == \"dropped\"\n        - closed\n        - (today() - closed).days <= "+strconv.Itoa(basesRecentlyShown)+"\n", order("file.name", "status", "parent", "decisions", "closed"))
	view("cards", "Cards by status", open, group("status")+order("file.name", "parent", "formula.age_days")+"    sort:\n"+sortBy("updated", "DESC"))
	view("cards", "Cards by parent", open, group("parent")+order("file.name", "status", "formula.age_days")+"    sort:\n"+sortBy("status", "ASC"))
	view("cards", "Cards by cycle", open, group("cycle")+order("file.name", "status", "parent", "priority")+"    sort:\n"+sortBy("priority", "DESC"))
	view("kanban", "Kanban by status", open, group("status")+"    groupOrder:\n      - raw\n      - ready\n      - active\n      - blocked\n"+order("file.name", "parent", "formula.age_days"))
	w("  - type: table\n    name: ADRs\n    filters:\n      and:\n")
	w("        - " + yamlQuote(`file.inFolder("`+adr+`")`) + "\n")
	w("        - file.name != \"README.md\"\n")
	w(order("file.name", "status", "date", "scope") + "    sort:\n" + sortBy("date", "DESC"))
	return b.String()
}
