package api

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/migrate"
)

// mintRegistry is a migration from layout 0.1.0 whose only step mints ids for
// the notes directory of the fixture.
func mintRegistry(t *testing.T) {
	t.Helper()
	old := migrate.Registry
	migrate.Registry = []migrate.Migration{{From: "0.1.0", To: "0.2.0", Steps: []migrate.Step{
		{ID: "mint-notes", Kind: "mint_ids", Class: "auto", Description: "give every note an id", Dir: "docs/notes"},
	}}}
	t.Cleanup(func() { migrate.Registry = old; migrateRand = nil })
}

// seed makes the random source repeatable: the same draws as the unit test.
func seed(t *testing.T) {
	t.Helper()
	var b []byte
	for _, g := range []byte{0, 0, 1, 2, 2, 3, 3, 4} {
		b = append(b, g, g, g, g)
	}
	migrateRand = bytes.NewReader(b)
	t.Cleanup(func() { migrateRand = nil })
}

func mintFixtureRepo(t *testing.T) string {
	t.Helper()
	src := "../../testdata/migrate/mint_ids"
	dst := t.TempDir()
	err := filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		rel, _ := filepath.Rel(src, p)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		b, _ := os.ReadFile(p)
		return os.WriteFile(filepath.Join(dst, rel), b, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	return dst
}

func stepOf(t *testing.T, resp Obj, id string) Obj {
	t.Helper()
	steps, _ := resp.Get("steps")
	for _, s := range steps.([]Obj) {
		if v, _ := s.Get("id"); v == id {
			return s
		}
	}
	t.Fatalf("no step %s", id)
	return nil
}

// snap reads every file under root except .git, which git itself may change at
// any time (background maintenance).
func snap(t *testing.T, root string) map[string]string {
	t.Helper()
	m := map[string]string{}
	filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			b, _ := os.ReadFile(p)
			rel, _ := filepath.Rel(root, p)
			m[filepath.ToSlash(rel)] = string(b)
		}
		return nil
	})
	return m
}

func sameSnap(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if w, ok := b[k]; !ok || v != w {
			return false
		}
	}
	return true
}

func gitOut(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %s", args, out)
	}
	return strings.TrimSpace(string(out))
}

func TestMigrateMintIDsModes(t *testing.T) {
	mintRegistry(t)
	root := mintFixtureRepo(t)
	gitRepo(t, root)
	before := snap(t, root)

	// report: the step is listed, nothing is computed or written
	resp := mustOK(t, root, "migrate", Request{})
	files, _ := stepOf(t, resp, "mint-notes").Get("files")
	if st, _ := resp.Get("state"); st != "behind" || len(files.([]string)) != 0 {
		t.Fatalf("report: %v %v", st, files)
	}

	// plan: old name, new name and id, nothing written
	seed(t)
	resp = mustOK(t, root, "migrate", Request{"mode": "plan", "as_of": "2026-10-28T10:00"})
	step := stepOf(t, resp, "mint-notes")
	files, _ = step.Get("files")
	notes, _ := step.Get("notes")
	if len(files.([]string)) != 6 || files.([]string)[0] != "docs/notes/260904-2222-beta.md" {
		t.Fatalf("plan files %v", files)
	}
	if !strings.Contains(strings.Join(notes.([]string), "\n"), "docs/notes/gamma.md -> docs/notes/261028-3333-gamma.md (id 261028-3333, date today") {
		t.Fatalf("plan notes %v", notes)
	}
	if st, _ := step.Get("status"); st != "pending" {
		t.Fatalf("plan status %v", st)
	}
	if !sameSnap(before, snap(t, root)) {
		t.Fatal("plan changed the tree")
	}

	// handoff carries the notes
	seed(t)
	resp = mustOK(t, root, "migrate", Request{"mode": "handoff", "as_of": "2026-10-28T10:00"})
	if h, _ := resp.Get("handoff"); !strings.Contains(h.(string), "reference not rewritten: docs/overview.md line 3") {
		t.Fatalf("handoff:\n%v", h)
	}
	if !sameSnap(before, snap(t, root)) {
		t.Fatal("handoff changed the tree")
	}

	// apply_safe: refused without the acknowledgement, applied with it, never committed
	wantCode(t, root, "migrate", Request{"mode": "apply_safe"}, "refused", "no_backup_ack")
	seed(t)
	resp = mustOK(t, root, "migrate", Request{"mode": "apply_safe", "acknowledge_backup": true, "as_of": "2026-10-28T10:00"})
	step = stepOf(t, resp, "mint-notes")
	if st, _ := step.Get("status"); st != "applied" {
		t.Fatalf("apply status %v", st)
	}
	if st, _ := resp.Get("state"); st != "current" {
		t.Fatalf("state %v", st)
	}
	for _, p := range []string{"docs/notes/261022-2222-alpha-note.md", "docs/notes/261028-4444-no-frontmatter.md", "docs/notes/261001-0000-needs-slugging-co.md"} {
		if _, err := os.Stat(filepath.Join(root, p)); err != nil {
			t.Errorf("missing %s", p)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "docs/notes/alpha-note.md")); err == nil {
		t.Error("old name still there")
	}
	if b, _ := os.ReadFile(filepath.Join(root, "docs/notes/261022-2222-alpha-note.md")); !strings.HasPrefix(string(b), "---\nid: \"261022-2222\"\ntitle: Alpha note\n") {
		t.Fatalf("alpha:\n%s", b)
	}
	logb, _ := os.ReadFile(filepath.Join(root, migrate.LogPath))
	if !strings.Contains(string(logb), `"step":"mint-notes"`) || !strings.Contains(string(logb), "alpha-note.md -\\u003e docs/notes/261022-2222-alpha-note.md") {
		t.Fatalf("log:\n%s", logb)
	}
	if n := gitOut(t, root, "log", "--oneline"); strings.Count(n, "\n") != 0 {
		t.Fatalf("the tool committed:\n%s", n)
	}
	if gitOut(t, root, "branch", "--show-current") != "migrate-0.1.0-to-0.2.0" || gitOut(t, root, "tag", "--list") != "pre-migrate-0.1.0-to-0.2.0" {
		t.Fatal("no branch or tag")
	}
	// the manifest is on 0.2.0 now, so another run has nothing to do
	resp = mustOK(t, root, "migrate", Request{"mode": "apply_safe", "acknowledge_backup": true})
	if c, _ := resp.Get("clean"); c != true {
		t.Fatalf("second call: %v", resp)
	}
}

// A second run, from a committed result with the manifest put back, changes
// nothing but the manifest.
func TestMigrateMintIDsTwice(t *testing.T) {
	mintRegistry(t)
	root := mintFixtureRepo(t)
	gitRepo(t, root)
	seed(t)
	mustOK(t, root, "migrate", Request{"mode": "apply_safe", "acknowledge_backup": true, "as_of": "2026-10-28T10:00"})
	gitOut(t, root, "add", "-A")
	gitOut(t, root, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "-m", "migrated")
	// back to the old manifest on a fresh branch name, so the safeguards pass again
	gitOut(t, root, "tag", "-d", "pre-migrate-0.1.0-to-0.2.0")
	gitOut(t, root, "branch", "-m", "first-run")
	if err := os.WriteFile(filepath.Join(root, "docs/.layout"), []byte("layout: 0.1.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitOut(t, root, "add", "docs/.layout")
	gitOut(t, root, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "-m", "manifest back")
	before := snap(t, root)

	migrateRand = bytes.NewReader(bytes.Repeat([]byte{9}, 64))
	resp := mustOK(t, root, "migrate", Request{"mode": "plan"})
	if f, _ := stepOf(t, resp, "mint-notes").Get("files"); len(f.([]string)) != 0 {
		t.Fatalf("plan after the run: %v", f)
	}
	resp = mustOK(t, root, "migrate", Request{"mode": "apply_safe", "acknowledge_backup": true})
	if st, _ := stepOf(t, resp, "mint-notes").Get("status"); st != "skipped" {
		t.Fatalf("second run status %v", st)
	}
	after := snap(t, root)
	for p, v := range after {
		switch p {
		case "docs/.layout", migrate.LogPath:
			continue
		}
		if w, ok := before[p]; !ok || w != v {
			t.Errorf("%s changed on the second run", p)
		}
	}
	if len(after) != len(before) {
		t.Errorf("file count %d, was %d", len(after), len(before))
	}
	if out := gitOut(t, root, "status", "--porcelain"); strings.Contains(out, "docs/notes") {
		t.Fatalf("notes changed:\n%s", out)
	}
}
