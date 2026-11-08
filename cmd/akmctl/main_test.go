package main

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/testutil"
)

// TestBuiltBinary builds the command and runs it, so main and the usage text
// are exercised as a person meets them (work item 261108-qv1w, finding F19).
func TestBuiltBinary(t *testing.T) {
	goTool := testutil.RequireTool(t, "go")
	exe := filepath.Join(t.TempDir(), "akmctl")
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}
	if out, err := exec.Command(goTool, "build", "-o", exe, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	run := func(args ...string) (string, int) {
		cmd := exec.Command(exe, args...)
		cmd.Dir = t.TempDir()
		out, err := cmd.CombinedOutput()
		code := 0
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else if err != nil {
			t.Fatal(err)
		}
		return string(out), code
	}
	if out, code := run("version"); code != 0 || !strings.HasPrefix(out, "akmctl ") || !strings.Contains(out, "(API 1.0)") {
		t.Errorf("version: %d %q", code, out)
	}
	if out, code := run("help"); code != 0 || !strings.Contains(out, "usage:") {
		t.Errorf("help: %d %q", code, out)
	}
	if out, code := run("nonsense"); code != 2 || out == "" {
		t.Errorf("an unknown command: %d %q", code, out)
	}
}
