package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	cwd := t.TempDir()
	var out, errb bytes.Buffer
	do := func(args ...string) int {
		out.Reset()
		errb.Reset()
		return Run(args, &out, &errb, cwd)
	}
	if code := do("init", "proj"); code != 0 {
		t.Fatal(code, errb.String())
	}
	if !strings.Contains(out.String(), "wrote    docs/.layout") {
		t.Errorf("output %q", out.String())
	}
	if _, err := os.Stat(filepath.Join(cwd, "proj", "AGENTS.md")); err != nil {
		t.Error(err)
	}
	// a second run changes nothing and still exits 0
	if code := do("init", "proj", "--format", "json"); code != 0 || !strings.Contains(out.String(), `"written":[]`) {
		t.Errorf("second run: %d %s", code, out.String())
	}
	// the working directory itself, with a docs dir and wikilinks
	if code := do("init", "--docs-dir", "knowledge", "--wikilinks"); code != 0 {
		t.Fatal(code, errb.String())
	}
	if _, err := os.Stat(filepath.Join(cwd, "knowledge", "work", "done")); err != nil {
		t.Error(err)
	}
	if b, _ := os.ReadFile(filepath.Join(cwd, "knowledge", ".layout")); !strings.Contains(string(b), "reference_style: wikilink\n") || !strings.Contains(string(b), "work_dir: knowledge/work\n") {
		t.Errorf("manifest %q", b)
	}
	if _, err := os.Stat(filepath.Join(cwd, "docs")); err == nil {
		t.Error("a docs folder appeared")
	}
	// an absolute directory; the retired --profile and bad folder names
	abs := filepath.Join(t.TempDir(), "abs")
	if code := do("init", abs); code != 0 {
		t.Fatal(code, errb.String())
	}
	if _, err := os.Stat(filepath.Join(abs, "docs", ".layout")); err != nil {
		t.Error(err)
	}
	if _, err := os.Stat(filepath.Join(abs, ".gitignore")); err == nil {
		t.Error(".gitignore written outside a work tree")
	}
	if code := do("init", "--profile", "code", abs); code != 2 {
		t.Errorf("--profile: %d", code)
	}
	for _, bad := range []string{"a/b", "..", ".", ".hidden", ""} {
		if code := do("init", "--docs-dir", bad, filepath.Join(t.TempDir(), "z")); code != 2 {
			t.Errorf("--docs-dir %q: %d", bad, code)
		}
	}
}
