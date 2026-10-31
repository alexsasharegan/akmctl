package api

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/record"
)

// Hidden files follow ADR 261030-1wj5: .layout is the one recognisable
// configuration file, and every other file the tool keeps for itself is .akm*.
// The rest are files other tools define (git, Prettier) that init writes.
func TestHiddenFilesFollowTheNamespaceRule(t *testing.T) {
	base := t.TempDir()
	if _, aerr := Run(&record.Repo{Root: base}, Info{"akmctl", "test"}, "init", Request{"dir": "p"}); aerr != nil {
		t.Fatal(aerr)
	}
	root := filepath.Join(base, "p")
	mustOK(t, root, "new", Request{"kind": "work", "slug": "first"})
	mustOK(t, root, "index", Request{})
	other := map[string]bool{".gitkeep": true, ".gitignore": true, ".prettierignore": true, ".prettierrc.json": true}
	filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		n := d.Name()
		if err != nil || d.IsDir() || !strings.HasPrefix(n, ".") {
			return nil
		}
		if n != ".layout" && !strings.HasPrefix(n, ".akm") && !other[n] {
			t.Errorf("%s is a hidden file that is neither .layout, .akm* nor a known file of another tool", p)
		}
		return nil
	})
}
