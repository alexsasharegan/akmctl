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
	if code := do("init", "--profile", "code", "proj"); code != 0 {
		t.Fatal(code, errb.String())
	}
	if !strings.Contains(out.String(), "wrote    docs/.layout") {
		t.Errorf("output %q", out.String())
	}
	if _, err := os.Stat(filepath.Join(cwd, "proj", "AGENTS.md")); err != nil {
		t.Error(err)
	}
	// a second run changes nothing and still exits 0
	if code := do("init", "--profile", "code", "proj", "--format", "json"); code != 0 || !strings.Contains(out.String(), `"written":[]`) {
		t.Errorf("second run: %d %s", code, out.String())
	}
	// the working directory itself
	if code := do("init", "--profile", "notes"); code != 0 {
		t.Fatal(code, errb.String())
	}
	if _, err := os.Stat(filepath.Join(cwd, "notes", "work", "done")); err != nil {
		t.Error(err)
	}
	// an absolute directory, and a missing or wrong profile
	abs := filepath.Join(t.TempDir(), "abs")
	if code := do("init", "--profile=code", abs); code != 0 {
		t.Fatal(code, errb.String())
	}
	if _, err := os.Stat(filepath.Join(abs, "docs", ".layout")); err != nil {
		t.Error(err)
	}
	if code := do("init", abs); code != 2 {
		t.Errorf("no profile: %d", code)
	}
	if code := do("init", "--profile", "wiki", abs); code != 2 {
		t.Errorf("bad profile: %d", code)
	}
}
