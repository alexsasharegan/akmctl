package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const repo = "../../testdata/api/repos/basic"

func run(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, errb bytes.Buffer
	code := Run(append(args, "--repo", repo), &out, &errb, ".")
	return code, out.String(), errb.String()
}

func TestJSONOutputIsOneLine(t *testing.T) {
	code, out, _ := run(t, "ls", "--format", "json", "--status", "ready", "--fields", "title")
	if code != 0 || strings.Count(out, "\n") != 1 {
		t.Fatalf("code %d, output %q", code, out)
	}
	var got struct {
		OK    bool
		Count int
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil || !got.OK || got.Count != 1 {
		t.Fatalf("%v %+v", err, got)
	}
}

func TestFlagsMapToParameters(t *testing.T) {
	code, out, _ := run(t, "ls", "--format=json", "--count-only", "--no-next", "--type", "work", "--limit=2")
	if code != 0 {
		t.Fatal(out)
	}
	for _, want := range []string{`"count_only":true`, `"next":false`, `"limit":2`, `"total":6`} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
	}
}

func TestPositionals(t *testing.T) {
	_, out, _ := run(t, "search", "JSON", "output", "--format", "ids")
	if strings.TrimSpace(out) != "261020-a2b3" {
		t.Fatalf("search words: %q", out)
	}
	code, out, _ := run(t, "find", "261022-e5f6")
	if code != 0 || strings.TrimSpace(out) != "docs/work/261022-e5f6-document-the-exit-codes.md" {
		t.Fatalf("%d %q", code, out)
	}
}

func TestErrorsExitTwo(t *testing.T) {
	for _, args := range [][]string{{"show", "261099-zzzz"}, {"ls", "--status", "finished"}, {"ls", "--bogus"}, {"show"}, {"nonsense"}, {"ls", "extra"}} {
		code, _, _ := run(t, args...)
		if code != 2 {
			t.Errorf("%v exits %d", args, code)
		}
	}
	code, out, _ := run(t, "show", "261099-zzzz", "--format", "json")
	if code != 2 || !strings.Contains(out, `"code":"not_found"`) {
		t.Fatalf("%d %s", code, out)
	}
}

func TestNoRepository(t *testing.T) {
	var out, errb bytes.Buffer
	code := Run([]string{"ls", "--format", "json"}, &out, &errb, t.TempDir())
	if code != 2 || !strings.Contains(out.String(), `"code":"unreadable"`) {
		t.Fatalf("%d %s", code, out.String())
	}
}

func TestWriteCommands(t *testing.T) {
	root := t.TempDir()
	copyDir(t, "../../testdata/api/repos/claims", root)
	do := func(args ...string) (int, string) {
		var out, errb bytes.Buffer
		code := Run(append(args, "--repo", root, "--format", "json"), &out, &errb, ".")
		return code, out.String()
	}
	code, out := do("new", "work", "add-a-thing", "--title", "Add a thing", "--field", "priority=2", "--field", "blocked_by=a person,another", "--as-of", "2026-10-28")
	if code != 0 || !strings.Contains(out, `"priority":2`) || !strings.Contains(out, `"blocked_by":["a person","another"]`) {
		t.Fatalf("new: %d %s", code, out)
	}
	code, out = do("set", "261022-e5f6", "priority=5", "summary=Short.", "parent=null", "--as-of", "2026-10-28")
	if code != 0 || !strings.Contains(out, `"changed":["parent","priority","summary"]`) {
		t.Fatalf("set: %d %s", code, out)
	}
	code, out = do("claim", "261022-e5f6", "--as", "agent-1", "--as-of", "2026-10-28T09:30")
	if code != 0 || !strings.Contains(out, `"claimed_by":"agent-1"`) {
		t.Fatalf("claim: %d %s", code, out)
	}
	code, out = do("release", "261022-e5f6", "--as-of", "2026-10-28")
	if code != 0 {
		t.Fatalf("release: %d %s", code, out)
	}
	code, out = do("close", "261022-e5f6", "done")
	if code != 2 || !strings.Contains(out, `"reason":"empty_outcome"`) {
		t.Fatalf("close: %d %s", code, out)
	}
	code, _ = do("check")
	if code == 2 {
		t.Fatal("check could not run")
	}
}

func copyDir(t *testing.T, src, dst string) {
	t.Helper()
	filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		rel, _ := filepath.Rel(src, p)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		b, _ := os.ReadFile(p)
		return os.WriteFile(filepath.Join(dst, rel), b, 0o644)
	})
}

func TestMigrateFlags(t *testing.T) {
	root := t.TempDir()
	copyDir(t, "../../testdata/api/repos/basic", root)
	var out, errb bytes.Buffer
	code := Run([]string{"migrate", "--mode", "handoff", "--target", "0.3.0", "--out", filepath.Join(root, "handoff.md"), "--repo", root, "--format", "json"}, &out, &errb, ".")
	if code != 0 || !strings.Contains(out.String(), `"state":"current"`) {
		t.Fatalf("%d %s %s", code, out.String(), errb.String())
	}
	if b, err := os.ReadFile(filepath.Join(root, "handoff.md")); err != nil || !strings.Contains(string(b), "Finish the migration") {
		t.Fatalf("handoff file: %v %s", err, b)
	}
	out.Reset()
	code = Run([]string{"migrate", "--mode", "apply-safe", "--repo", root, "--format", "json"}, &out, &errb, ".")
	if code != 2 || !strings.Contains(out.String(), `"reason":"no_backup_ack"`) {
		t.Fatalf("%d %s", code, out.String())
	}
}

// The CLI reads the repository directly, so these run against the fixture
// without writing: dry runs only.
func TestSectionDryRunPrintsADiff(t *testing.T) {
	before, _ := os.ReadFile(filepath.Join(repo, "docs/work/261022-e5f6-document-the-exit-codes.md"))
	code, out, errs := run(t, "section", "261022-e5f6", "goal", "--text", "New goal.", "--dry-run")
	if code != 0 {
		t.Fatalf("%d %s", code, errs)
	}
	for _, want := range []string{"--- docs/work/261022-e5f6-document-the-exit-codes.md\n", "+++ ", "@@ ", "-Do the thing.\n", "+New goal.\n"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in\n%s", want, out)
		}
	}
	after, _ := os.ReadFile(filepath.Join(repo, "docs/work/261022-e5f6-document-the-exit-codes.md"))
	if string(before) != string(after) {
		t.Fatal("dry run wrote")
	}
	code, out, _ = run(t, "close", "261020-a2b3", "done", "--outcome", "Shipped.", "--dry-run", "--format", "json")
	if code != 0 || !strings.Contains(out, `"dry_run":true`) || !strings.Contains(out, `"outcome":"Shipped."`) || !strings.Contains(out, `"new_path":"docs/work/done/`) {
		t.Fatalf("%d %s", code, out)
	}
	if code, _, _ := run(t, "section", "261022-e5f6", "--text", "x", "--dry-run"); code != 2 {
		t.Fatalf("missing name should be invalid_request: %d", code)
	}
}

func TestMigrationFileFlag(t *testing.T) {
	root := t.TempDir()
	for p, text := range map[string]string{
		"docs/.layout":             "work_dir: Vault Notes/tasks\n",
		"Vault Notes/tasks/a b.md": "---\ncreated_at: 2026-10-20\n---\n# A\n",
	} {
		full := filepath.Join(root, filepath.FromSlash(p))
		os.MkdirAll(filepath.Dir(full), 0o755)
		os.WriteFile(full, []byte(text), 0o644)
	}
	one := `{"from":"0.0.0","to":"0.3.0","steps":[{"id":"c","kind":"rename_field","class":"auto","description":"d","dir":"Vault Notes/tasks","from":"created_at","to":"created"}]}`
	file := filepath.Join(t.TempDir(), "one-off.json")
	os.WriteFile(file, []byte(one), 0o644)
	var out, errb bytes.Buffer
	code := Run([]string{"migrate", "--mode", "plan", "--migration", file, "--repo", root, "--format", "json"}, &out, &errb, ".")
	if code != 1 || !strings.Contains(out.String(), `"state":"behind"`) || !strings.Contains(out.String(), `"files":["Vault Notes/tasks/a b.md"]`) {
		t.Fatalf("code %d out %s err %s", code, out.String(), errb.String())
	}
	// a list of migrations is accepted too
	os.WriteFile(file, []byte("["+one+"]"), 0o644)
	out.Reset()
	if code := Run([]string{"migrate", "--migration", file, "--repo", root, "--format", "json"}, &out, &errb, "."); code != 1 || !strings.Contains(out.String(), `"state":"behind"`) {
		t.Fatalf("list: code %d out %s", code, out.String())
	}
	// errors: a missing file, bad JSON, a bad migration, another command
	os.WriteFile(file, []byte("{nope"), 0o644)
	for _, args := range [][]string{
		{"migrate", "--migration", filepath.Join(t.TempDir(), "missing.json")},
		{"migrate", "--migration", file},
		{"ls", "--migration", file},
		{"migrate", "--migration"},
	} {
		out.Reset()
		errb.Reset()
		if code := Run(append(args, "--repo", root), &out, &errb, "."); code != 2 {
			t.Errorf("%v: code %d", args, code)
		}
	}
	os.WriteFile(file, []byte(`{"from":"0.0.0","to":"0.3.0","steps":[{"id":"x","kind":"shell","class":"auto"}]}`), 0o644)
	out.Reset()
	if code := Run([]string{"migrate", "--migration", file, "--repo", root, "--format", "json"}, &out, &errb, "."); code != 2 || !strings.Contains(out.String(), "invalid_request") {
		t.Fatalf("bad migration: code %d out %s", code, out.String())
	}
}
