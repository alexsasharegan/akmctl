package testutil

import "testing"

func TestRequireToolFindsAProgram(t *testing.T) {
	t.Setenv("AKM_REQUIRE_TOOLS", "1")
	if RequireTool(t, "go") == "" {
		t.Fatal("no path for go")
	}
}

func TestRequireToolSkipsWithoutTheSwitch(t *testing.T) {
	t.Setenv("AKM_REQUIRE_TOOLS", "")
	t.Run("missing", func(t *testing.T) {
		RequireTool(t, "akm-no-such-tool")
		t.Fatal("not skipped")
	})
}
