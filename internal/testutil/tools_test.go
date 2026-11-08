package testutil

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

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

func TestCopyTreeAndSnapshot(t *testing.T) {
	src, dst := t.TempDir(), filepath.Join(t.TempDir(), "copy")
	os.MkdirAll(filepath.Join(src, "a", "b"), 0o755)
	os.WriteFile(filepath.Join(src, "a", "b", "x.md"), []byte("x\n"), 0o644)
	os.WriteFile(filepath.Join(src, "y.md"), []byte("y\n"), 0o644)
	CopyTree(t, src, dst)
	want := map[string]string{"a/b/x.md": "x\n", "y.md": "y\n"}
	if got := Snapshot(t, dst); !reflect.DeepEqual(got, want) {
		t.Fatalf("%v, want %v", got, want)
	}
}
