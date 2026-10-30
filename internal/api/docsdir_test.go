package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/record"
)

// call loads the repository at root and runs one operation against it.
func callOK(t *testing.T, root, op string, req Request) Obj {
	t.Helper()
	repo, err := record.Load(root)
	if err != nil {
		t.Fatalf("%s: load: %v", op, err)
	}
	resp, aerr := Run(repo, Info{"akmctl", "test"}, op, req)
	if aerr != nil {
		t.Fatalf("%s %v: %v", op, req, aerr)
	}
	return resp
}

func findingPaths(resp Obj) []string {
	var out []string
	f, _ := resp.Get("findings")
	for _, x := range f.([]Obj) {
		p, _ := x.Get("file")
		out = append(out, p.(string))
	}
	return out
}

// A repository made with a docs dir other than docs works end to end: every
// operation follows the manifest to its folder.
func TestCustomDocsDirEndToEnd(t *testing.T) {
	base := t.TempDir()
	if err := os.MkdirAll(filepath.Join(base, ".git"), 0o755); err != nil { // a work tree, for config and .gitignore
		t.Fatal(err)
	}
	initResp, aerr := Run(&record.Repo{Root: base}, Info{"akmctl", "test"}, "init", Request{"docs_dir": "knowledge", "wikilinks": true})
	if aerr != nil {
		t.Fatal(aerr)
	}
	if d, _ := initResp.Get("docs_dir"); d != "knowledge" {
		t.Errorf("docs_dir %v", d)
	}
	if exists(base, "docs") {
		t.Fatal("a docs folder exists")
	}
	if b, _ := os.ReadFile(filepath.Join(base, ".gitignore")); string(b) != "# the advisory lock akmctl takes while it writes\nknowledge/.akm.lock\n" {
		t.Errorf(".gitignore %q", b)
	}

	c := callOK(t, base, "check", Request{})
	if e, _ := c.Get("errors"); e != 0 {
		t.Fatalf("check after init: %v", c)
	}
	if w, _ := c.Get("warnings"); w != 0 {
		t.Fatalf("check after init: %v", c)
	}
	v := callOK(t, base, "version", Request{})
	repo, _ := v.Get("repository")
	if wd, _ := repo.(Obj).Get("work_dir"); wd != "knowledge/work" {
		t.Errorf("version work_dir %v", wd)
	}

	// new, set, claim, section, close, index
	resp := callOK(t, base, "new", Request{"kind": "work", "slug": "first", "title": "First"})
	rec, _ := resp.Get("record")
	id, _ := rec.(Obj).Get("id")
	if p, _ := rec.(Obj).Get("path"); !strings.HasPrefix(p.(string), "knowledge/work/") {
		t.Errorf("path %v", p)
	}
	callOK(t, base, "section", Request{"id": id, "name": "Goal", "text": "Try it.\n\n- [ ] it works"})
	callOK(t, base, "set", Request{"id": id, "changes": map[string]any{"status": "ready"}})
	callOK(t, base, "claim", Request{"id": id, "as": "agent"})
	if !exists(base, "knowledge/.akm.lock") {
		t.Error("the lock is not in the docs dir")
	}
	if exists(base, "docs") {
		t.Error("a write created docs/")
	}
	callOK(t, base, "section", Request{"id": id, "name": "Next", "text": "Close it."})
	callOK(t, base, "section", Request{"id": id, "name": "Outcome", "text": "It works."})
	callOK(t, base, "close", Request{"id": id, "status": "done"})
	if !exists(base, "knowledge/work/done") {
		t.Error("closed item not under knowledge/work/done")
	}
	callOK(t, base, "new", Request{"kind": "adr", "slug": "use-it"})
	callOK(t, base, "index", Request{})
	if !exists(base, "knowledge/adr/README.md") {
		t.Error("no index")
	}
	mr := callOK(t, base, "migrate", Request{"mode": "report"})
	if st, _ := mr.Get("state"); st != "current" {
		t.Errorf("migrate state %v", st)
	}
	if exists(base, "docs") {
		t.Error("an operation created docs/")
	}

	// config on and off
	callOK(t, base, "config", Request{"key": "git_commits", "value": true})
	if b, _ := os.ReadFile(filepath.Join(base, "knowledge", ".layout")); !strings.Contains(string(b), "git_commits: true\n") {
		t.Errorf("manifest after config on: %q", b)
	}
	if !exists(base, ".agents/skills/akm-commit/SKILL.md") {
		t.Error("skill not installed")
	}
	if c := callOK(t, base, "check", Request{}); true {
		if e, _ := c.Get("errors"); e != 0 {
			t.Errorf("check after config on: %v", c)
		}
	}
	callOK(t, base, "config", Request{"key": "git_commits", "value": false, "confirm": true})
	if b, _ := os.ReadFile(filepath.Join(base, "knowledge", ".layout")); !strings.Contains(string(b), "git_commits: false\n") {
		t.Errorf("manifest after config off: %q", b)
	}
	if exists(base, ".agents/skills/akm-commit") {
		t.Error("skill not removed")
	}

	// a repository found from a subdirectory
	for _, sub := range []string{"knowledge/work", "knowledge/adr", "."} {
		got, err := record.FindRoot(filepath.Join(base, filepath.FromSlash(sub)))
		if err != nil {
			t.Fatal(err)
		}
		if a, b := evalSymlinks(got), evalSymlinks(base); a != b {
			t.Errorf("from %s: root %s, want %s", sub, a, b)
		}
	}
}

func evalSymlinks(p string) string {
	if q, err := filepath.EvalSymlinks(p); err == nil {
		return q
	}
	return p
}

// Findings name the real path of the manifest and the deny list.
func TestCustomDocsDirFindingPaths(t *testing.T) {
	base := t.TempDir()
	if _, aerr := Run(&record.Repo{Root: base}, Info{"akmctl", "test"}, "init", Request{"docs_dir": "knowledge"}); aerr != nil {
		t.Fatal(aerr)
	}
	// git_commits on outside a work tree: S4 names the manifest
	m := filepath.Join(base, "knowledge", ".layout")
	b, _ := os.ReadFile(m)
	os.WriteFile(m, append(b, "git_commits: true\n"...), 0o644)
	// a deny-list hit and a bad deny pattern
	os.WriteFile(filepath.Join(base, "knowledge", ".akm-deny"), []byte("forbidden phrase\nre:(\n"), 0o644)
	a, _ := os.ReadFile(filepath.Join(base, "AGENTS.md"))
	os.WriteFile(filepath.Join(base, "AGENTS.md"), append(a, "\nthe forbidden phrase\n"...), 0o644)
	c := callOK(t, base, "check", Request{})
	paths := strings.Join(findingPaths(c), " ")
	if !strings.Contains(paths, "knowledge/.layout") || !strings.Contains(paths, "knowledge/.akm-deny") {
		t.Errorf("finding paths %q in %v", paths, c)
	}
	if strings.Contains(paths, "docs/") {
		t.Errorf("a finding names docs/: %q", paths)
	}
	f, _ := c.Get("findings")
	var denied bool
	for _, x := range f.([]Obj) {
		if r, _ := x.Get("rule"); r == "SF2" {
			if msg, _ := x.Get("message"); strings.Contains(msg.(string), "knowledge/.akm-deny") {
				denied = true
			}
		}
	}
	if !denied {
		t.Errorf("no SF2 finding naming knowledge/.akm-deny: %v", f)
	}
}

// Two manifests make the repository unreadable, naming both.
func TestTwoManifestsAreAnError(t *testing.T) {
	base := t.TempDir()
	for _, d := range []string{"notes", "wiki"} {
		if _, aerr := Run(&record.Repo{Root: base}, Info{"akmctl", "test"}, "init", Request{"docs_dir": d, "dir": d + "-repo"}); aerr != nil {
			t.Fatal(aerr)
		}
	}
	// move them side by side under one root
	root := t.TempDir()
	for _, d := range []string{"notes", "wiki"} {
		if err := os.Rename(filepath.Join(base, d+"-repo", d), filepath.Join(root, d)); err != nil {
			t.Fatal(err)
		}
	}
	_, err := record.Load(root)
	if err == nil || !strings.Contains(err.Error(), "notes/.layout") || !strings.Contains(err.Error(), "wiki/.layout") {
		t.Fatalf("Load: %v", err)
	}
	// init in that root with a third folder or docs is refused too
	_, aerr := Run(&record.Repo{Root: root}, Info{"akmctl", "test"}, "init", Request{})
	if aerr == nil {
		t.Error("init over two manifests succeeded")
	}
}
