package api

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// crlfRepo copies a fixture repository and turns every Markdown file into CRLF,
// as a Windows checkout with core.autocrlf=true does.
func crlfRepo(t *testing.T) string {
	t.Helper()
	root := copyRepo(t, "basic", "")
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".md") {
			return err
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		b = bytes.ReplaceAll(bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n")), []byte("\n"), []byte("\r\n"))
		return os.WriteFile(p, b, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	return root
}

// mixedEndings lists the Markdown files under root that hold both CRLF and bare LF.
func mixedEndings(t *testing.T, root string) []string {
	t.Helper()
	var out []string
	filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".md") {
			return err
		}
		b, _ := os.ReadFile(p)
		crlf := bytes.Count(b, []byte("\r\n"))
		lf := bytes.Count(b, []byte("\n")) - crlf
		if crlf > 0 && lf > 0 {
			rel, _ := filepath.Rel(root, p)
			out = append(out, filepath.ToSlash(rel))
		}
		return nil
	})
	return out
}

// TestWritesOnACRLFCheckout simulates a Windows checkout on any platform (work
// item 261028-jzsx; the real filesystem behavior is the PowerShell harness,
// scripts/windows-smoke.ps1): every write keeps the file's CRLF and the links a
// move rewrites use forward slashes.
func TestWritesOnACRLFCheckout(t *testing.T) {
	t.Parallel()
	root := crlfRepo(t)
	item := "docs/work/261020-a2b3-add-json-output-to-ls.md"

	mustOK(t, root, "set", Request{"id": a2b3, "changes": map[string]any{"priority": 1}})
	mustOK(t, root, "section", Request{"id": a2b3, "name": "Next", "text": "Next line one.\nNext line two."})
	mustOK(t, root, "log", Request{"id": a2b3, "text": "a line from the CRLF test"})
	if got := mixedEndings(t, root); len(got) != 0 {
		t.Fatalf("set, section and log mixed line endings in %v:\n%q", got, read(t, root, item))
	}
	if s := read(t, root, item); !strings.Contains(s, "Next line one.\r\nNext line two.\r\n") {
		t.Errorf("a multi-line section write did not take the file's CRLF:\n%q", s)
	}

	// A link that the move must rewrite, in a file that is CRLF.
	other := "docs/work/261022-e5f6-document-the-exit-codes.md"
	if _, err := os.Stat(filepath.Join(root, other)); err != nil {
		t.Fatalf("fixture has no %s", other)
	}
	b := read(t, root, other)
	if err := os.WriteFile(filepath.Join(root, other), []byte(b+"See [a](261020-a2b3-add-json-output-to-ls.md).\r\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mustOK(t, root, "section", Request{"id": a2b3, "name": "Outcome", "text": "Done on CRLF."})
	mustOK(t, root, "close", Request{"id": a2b3, "status": "done"})

	if got := mixedEndings(t, root); len(got) != 0 {
		t.Fatalf("close mixed line endings in %v", got)
	}
	s := read(t, root, other)
	if !strings.Contains(s, "[a](done/261020-a2b3-add-json-output-to-ls.md).\r\n") {
		t.Errorf("the inbound link was not rewritten with a forward slash and CRLF:\n%q", s)
	}
	if strings.Contains(s, `done\`) {
		t.Errorf("a rewritten link holds a backslash:\n%q", s)
	}
}
