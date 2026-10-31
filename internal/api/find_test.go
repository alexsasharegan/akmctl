package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindByNote(t *testing.T) {
	root := copyRepo(t, "basic", "")
	write := func(rel, text string) {
		full := filepath.Join(root, filepath.FromSlash(rel))
		os.MkdirAll(filepath.Dir(full), 0o755)
		os.WriteFile(full, []byte(text), 0o644)
	}
	write("Projects/Vault Views.md", "# Vault Views\n")
	write("Projects/Other/Plan.md", "# Plan\n")
	write("Archive/Plan.md", "# Plan\n")
	write(".obsidian/Hidden.md", "x")
	write("node_modules/pkg/Dep.md", "x")

	resp := mustOK(t, root, "find", Request{"name": "vault views"})
	if p, _ := resp.Get("path"); p != "Projects/Vault Views.md" {
		t.Fatalf("path %v", p)
	}
	if _, has := resp.Get("id"); has {
		t.Error("a note that is not a record has an id")
	}
	// a record by its file name stem carries its id and type
	resp = mustOK(t, root, "find", Request{"name": "261022-e5f6-document-the-exit-codes"})
	if id, _ := resp.Get("id"); id != "261022-e5f6" {
		t.Fatalf("record: %v", resp)
	}
	// a folder in the name picks one of two notes with the same file name
	resp = mustOK(t, root, "find", Request{"name": "Archive/Plan"})
	if p, _ := resp.Get("path"); p != "Archive/Plan.md" {
		t.Fatalf("path %v", p)
	}
	wantCode(t, root, "find", Request{"name": "Plan"}, "ambiguous", "")
	wantCode(t, root, "find", Request{"name": "Hidden"}, "not_found", "")
	wantCode(t, root, "find", Request{"name": "Dep"}, "not_found", "")
	wantCode(t, root, "find", Request{"name": "../outside"}, "invalid_request", "")
	wantCode(t, root, "find", Request{}, "invalid_request", "")
	wantCode(t, root, "find", Request{"id": "261022-e5f6", "name": "x"}, "invalid_request", "")
}
