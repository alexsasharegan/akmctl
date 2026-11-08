package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSetTitleRenamesTheHeading is work item 261108-kdz1: set title= left the
// "# heading" as it was, so check warned WK12 and set said nothing.
func TestSetTitleRenamesTheHeading(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "basic", "")
	mustOK(t, root, "set", Request{"id": e5f6ID, "changes": map[string]any{"title": "Explain the exit codes"}, "as_of": "2026-10-28"})
	text := read(t, root, e5f6Path)
	if !strings.Contains(text, "title: Explain the exit codes\n") || !strings.Contains(text, "\n# Explain the exit codes\n") || strings.Contains(text, "# Document the exit codes") {
		t.Fatal(text)
	}
	clean(t, root) // no WK12
	// a heading that already differed from the old title is the author's, and stays
	mustOK(t, root, "set", Request{"id": e5f6ID, "changes": map[string]any{"title": "Third name"}})
	text = read(t, root, e5f6Path)
	if !strings.Contains(text, "\n# Third name\n") {
		t.Fatal(text)
	}
	custom := strings.Replace(text, "\n# Third name\n", "\n# A heading of my own\n", 1)
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(e5f6Path)), []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}
	mustOK(t, root, "set", Request{"id": e5f6ID, "changes": map[string]any{"title": "Fourth name"}})
	if got := read(t, root, e5f6Path); !strings.Contains(got, "\n# A heading of my own\n") || !strings.Contains(got, "title: Fourth name\n") {
		t.Fatal(got)
	}
}
