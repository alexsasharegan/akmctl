package api

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/migrate"
)

// fakeRegistry is a migration from a made-up layout 0.1.0, where records used
// created_at and a "Result" section and had no .gitignore.
func fakeRegistry(t *testing.T) {
	t.Helper()
	old := migrate.Registry
	migrate.Registry = []migrate.Migration{{From: "0.1.0", To: "0.3.0", Steps: []migrate.Step{
		{ID: "rename-created", Kind: "rename_field", Class: "auto", Description: "created_at becomes created", Scope: "work", From: "created_at", To: "created", Fixes: []string{"WK3"}},
		{ID: "rename-result", Kind: "rename_heading", Class: "auto", Description: "## Result becomes ## Outcome", Scope: "work", From: "Result", To: "Outcome", Fixes: []string{"WK7"}},
		{ID: "default-kind", Kind: "set_default", Class: "auto", Description: "kind defaults to chore", Scope: "work", Field: "kind", Value: "chore"},
		{ID: "add-ignore", Kind: "add_file", Class: "auto", Description: "docs/.gitignore lists the lock", Path: "docs/.gitignore", Content: ".akm.lock\n"},
		{ID: "review-adrs", Kind: "manual", Class: "manual", Description: "check every ADR summary reads as one sentence", Fixes: []string{"A1"}},
	}}}
	t.Cleanup(func() { migrate.Registry = old })
}

func gitRepo(t *testing.T, root string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"add", "-A"}, {"-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "-m", "start"}} {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s", args, out)
		}
	}
}

// oldRepo is the basic fixture rewritten as layout 0.1.0.
func oldRepo(t *testing.T) string {
	root := copyRepo(t, "basic", "")
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	must(os.WriteFile(filepath.Join(root, "docs/.layout"), []byte("layout: 0.1.0\nextension.adr-index: 0.3.0\n"), 0o644))
	filepath.WalkDir(filepath.Join(root, "docs/work"), func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			b, _ := os.ReadFile(p)
			s := strings.Replace(string(b), "\ncreated: ", "\ncreated_at: ", 1)
			s = strings.Replace(s, "## Outcome", "## Result", 1)
			s = strings.Replace(s, "\nkind: feature", "", 1)
			os.WriteFile(p, []byte(s), 0o644)
		}
		return nil
	})
	return root
}

func TestMigrateReportAndPlan(t *testing.T) {
	fakeRegistry(t)
	root := oldRepo(t)
	resp := mustOK(t, root, "migrate", Request{})
	if st, _ := resp.Get("state"); st != "behind" {
		t.Fatalf("state %v", st)
	}
	if c, _ := resp.Get("clean"); c != false {
		t.Fatal("an old repo is not clean")
	}
	resp = mustOK(t, root, "migrate", Request{"mode": "plan"})
	steps, _ := resp.Get("steps")
	var created Obj
	for _, s := range steps.([]Obj) {
		if id, _ := s.Get("id"); id == "rename-created" {
			created = s
		}
	}
	files, _ := created.Get("files")
	if len(files.([]string)) != 6 {
		t.Fatalf("rename-created would touch %v", files)
	}
	// planning wrote nothing
	if b, _ := os.ReadFile(filepath.Join(root, "docs/work/261022-e5f6-document-the-exit-codes.md")); !strings.Contains(string(b), "created_at:") {
		t.Fatal("plan changed a file")
	}
	// classified findings: WK3 is fixed by an auto step
	resp = mustOK(t, root, "check", Request{"target": "0.3.0", "as_of": "2026-10-28"})
	fs, _ := resp.Get("findings")
	classes := map[string]any{}
	for _, f := range fs.([]Obj) {
		r, _ := f.Get("rule")
		c, _ := f.Get("class")
		classes[r.(string)] = c
	}
	if classes["WK3"] != "auto" {
		t.Fatalf("classes %v", classes)
	}
	rec, _ := resp.Get("reconcile")
	if s, _ := rec.(Obj).Get("state"); s != "behind" {
		t.Fatalf("reconcile %v", rec)
	}
}

func TestMigrateApplySafe(t *testing.T) {
	fakeRegistry(t)
	root := oldRepo(t)
	gitRepo(t, root)
	wantCode(t, root, "migrate", Request{"mode": "apply_safe"}, "refused", "no_backup_ack")
	resp := mustOK(t, root, "migrate", Request{"mode": "apply_safe", "acknowledge_backup": true, "as_of": "2026-10-28T10:00"})
	st, _ := resp.Get("state")
	// the manual step keeps the manifest on the old version
	if st != "behind" {
		t.Fatalf("state %v", st)
	}
	steps, _ := resp.Get("steps")
	got := map[string]any{}
	for _, s := range steps.([]Obj) {
		id, _ := s.Get("id")
		status, _ := s.Get("status")
		got[id.(string)] = status
	}
	if got["rename-created"] != "applied" || got["add-ignore"] != "applied" || got["review-adrs"] != "pending" {
		t.Fatalf("steps %v", got)
	}
	text, _ := os.ReadFile(filepath.Join(root, "docs/work/261022-e5f6-document-the-exit-codes.md"))
	if strings.Contains(string(text), "created_at") || !strings.Contains(string(text), "\ncreated: 2026-10-22\n") || !strings.Contains(string(text), "kind: chore") {
		t.Fatalf("e5f6:\n%s", text)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "docs/.gitignore")); string(b) != ".akm.lock\n" {
		t.Fatalf("gitignore %q", b)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "docs/work/done/261010-j9k0-choose-a-licence.md")); !strings.Contains(string(b), "## Outcome") {
		t.Fatal("heading not renamed")
	}
	logb, _ := os.ReadFile(filepath.Join(root, "docs/.migrate-log"))
	if !strings.Contains(string(logb), `"step":"rename-created"`) || strings.Count(string(logb), "\n") != 4 {
		t.Fatalf("log:\n%s", logb)
	}
	// a backup tag and a branch exist, and the manifest still says 0.1.0
	out, _ := exec.Command("git", "-C", root, "tag", "--list").Output()
	br, _ := exec.Command("git", "-C", root, "branch", "--show-current").Output()
	if strings.TrimSpace(string(out)) != "pre-migrate-0.1.0-to-0.3.0" || strings.TrimSpace(string(br)) != "migrate-0.1.0-to-0.3.0" {
		t.Fatalf("tag %q branch %q", out, br)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "docs/.layout")); !strings.Contains(string(b), "layout: 0.1.0") {
		t.Fatalf("manifest %s", b)
	}
	// the tool did not commit: the tree is dirty, so a second apply is refused
	wantCode(t, root, "migrate", Request{"mode": "apply_safe", "acknowledge_backup": true}, "refused", "dirty_tree")
}

func TestMigrateIdempotentAndFinishes(t *testing.T) {
	fakeRegistry(t)
	// without the manual step the manifest moves to 0.3.0
	migrate.Registry[0].Steps = migrate.Registry[0].Steps[:4]
	root := oldRepo(t)
	gitRepo(t, root)
	resp := mustOK(t, root, "migrate", Request{"mode": "apply_safe", "acknowledge_backup": true})
	if st, _ := resp.Get("state"); st != "current" {
		t.Fatalf("state %v", st)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "docs/.layout")); !strings.Contains(string(b), "layout: 0.3.0") || !strings.Contains(string(b), "extension.adr-index") {
		t.Fatalf("manifest %s", b)
	}
	// running the steps again changes nothing
	repo, _ := loadRepo(t, root)
	res, err := migrate.Apply(repo, migrate.Registry[0].Steps)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range res {
		if r.Status != "skipped" {
			t.Errorf("%s ran again: %v", r.Step.ID, r.Files)
		}
	}
	mustOK(t, root, "migrate", Request{"mode": "apply_safe", "acknowledge_backup": true}) // current: nothing to do
}

func TestMigrateRefusals(t *testing.T) {
	fakeRegistry(t)
	root := oldRepo(t) // not a git repository
	wantCode(t, root, "migrate", Request{"mode": "apply_safe", "acknowledge_backup": true}, "refused", "not_git")
	gitRepo(t, root)
	exec.Command("git", "-C", root, "tag", "pre-migrate-0.1.0-to-0.3.0").Run()
	wantCode(t, root, "migrate", Request{"mode": "apply_safe", "acknowledge_backup": true}, "refused", "tag_exists")
	root2 := copyRepo(t, "basic", "")
	os.WriteFile(filepath.Join(root2, "docs/.layout"), []byte("layout: 0.0.5\n"), 0o644)
	resp := mustOK(t, root2, "migrate", Request{})
	if st, _ := resp.Get("state"); st != "unsupported" {
		t.Fatalf("state %v", st)
	}
	wantCode(t, root2, "migrate", Request{"mode": "apply_safe", "acknowledge_backup": true}, "refused", "no_path")
	os.WriteFile(filepath.Join(root2, "docs/.layout"), []byte("layout: 2.0.0\n"), 0o644)
	resp = mustOK(t, root2, "migrate", Request{})
	if st, _ := resp.Get("state"); st != "ahead" {
		t.Fatalf("state %v", st)
	}
}

func TestMigrateHandoff(t *testing.T) {
	fakeRegistry(t)
	root := oldRepo(t)
	resp := mustOK(t, root, "migrate", Request{"mode": "handoff"})
	h, _ := resp.Get("handoff")
	for _, want := range []string{"from layout 0.1.0 to 0.3.0", "review-adrs", "akmctl check --target 0.3.0", "Never delete a file", "pre-migrate-0.1.0-to-0.3.0"} {
		if !strings.Contains(h.(string), want) {
			t.Errorf("handoff lacks %q", want)
		}
	}
}

func TestCheckReportsBehindWithSteps(t *testing.T) {
	fakeRegistry(t)
	root := oldRepo(t)
	resp := mustOK(t, root, "check", Request{"rules": []any{"S1"}})
	rec, _ := resp.Get("reconcile")
	steps, _ := rec.(Obj).Get("steps")
	if s, _ := rec.(Obj).Get("state"); s != "behind" || len(steps.([]string)) != 5 {
		t.Fatalf("reconcile %v", rec)
	}
}
