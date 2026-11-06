package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// agents prints what init wrote for the same settings, and writes nothing.
func TestAgentsEqualsInit(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"defaults", []string{"init", "--yes"}},
		{"docs dir and wikilinks", []string{"init", "--yes", "--docs-dir", "knowledge", "--wikilinks"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := filepath.Join(t.TempDir(), "proj")
			if code, _, errb := cfg(root, tc.args...); code != 0 {
				t.Fatalf("init: %d %s", code, errb)
			}
			want, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
			if err != nil {
				t.Fatal(err)
			}
			before := agentsTree(t, root)
			code, out, errb := cfg(root, "agents")
			if code != 0 || errb != "" {
				t.Fatalf("agents: %d %s", code, errb)
			}
			if out != string(want) {
				t.Errorf("agents differs from the AGENTS.md init wrote:\n%s", out)
			}
			if after := agentsTree(t, root); after != before {
				t.Errorf("agents wrote something:\nbefore %s\nafter  %s", before, after)
			}
		})
	}
}

// A repository that edited its AGENTS.md still gets the shipped text.
func TestAgentsIgnoresEdits(t *testing.T) {
	root := filepath.Join(t.TempDir(), "proj")
	if code, _, errb := cfg(root, "init", "--yes"); code != 0 {
		t.Fatal(errb)
	}
	os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("mine\n"), 0o644)
	code, out, _ := cfg(root, "agents")
	if code != 0 || strings.Contains(out, "mine") || !strings.Contains(out, "# proj") {
		t.Errorf("code %d, out %q", code, out)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "AGENTS.md")); string(b) != "mine\n" {
		t.Error("AGENTS.md was changed")
	}
}

func TestAgentsOutsideARepository(t *testing.T) {
	dir := t.TempDir()
	if code, _, errb := cfg(dir, "agents"); code != 2 || errb == "" {
		t.Errorf("code %d, stderr %q", code, errb)
	}
	if code, _, _ := cfg(dir, "agents", "extra"); code != 2 {
		t.Error("an argument is refused")
	}
}

// agentsTree lists every path and size under root.
func agentsTree(t *testing.T, root string) string {
	t.Helper()
	var b strings.Builder
	filepath.Walk(root, func(p string, fi os.FileInfo, err error) error {
		if err == nil {
			b.WriteString(p + ":" + fi.ModTime().String() + "\n")
		}
		return nil
	})
	return b.String()
}
