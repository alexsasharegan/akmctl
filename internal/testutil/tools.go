// Package testutil holds helpers shared by the tests of several packages.
package testutil

import (
	"os"
	"os/exec"
	"testing"
)

// RequireTool returns the path of an external program a test needs. When it is
// missing the test is skipped, unless AKM_REQUIRE_TOOLS=1 (CI sets it), when the
// test fails: a runner without git or node must not pass with those tests
// silently skipped (work item 261108-z9n8).
func RequireTool(t testing.TB, name string) string {
	t.Helper()
	p, err := exec.LookPath(name)
	if err == nil {
		return p
	}
	if os.Getenv("AKM_REQUIRE_TOOLS") == "1" {
		t.Fatalf("%s is required (AKM_REQUIRE_TOOLS=1) but not on PATH: %v", name, err)
	}
	t.Skipf("%s is not installed", name)
	return ""
}
