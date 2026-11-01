package api

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestCloseRewritesLinks: closing moves the item into done/ and repairs the
// relative Markdown links the move breaks, in the item and to it (work item
// 261031-wjs3).
func TestCloseRewritesLinks(t *testing.T) {
	root := copyRepo(t, "basic", "")
	item := "docs/work/261020-a2b3-add-json-output-to-ls.md"
	write := func(rel, text string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(filepath.Join(root, rel)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, rel), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	b, _ := os.ReadFile(filepath.Join(root, item))
	adrs, _ := filepath.Glob(filepath.Join(root, "docs/adr/*.md"))
	if len(adrs) == 0 {
		t.Fatal("fixture has no ADR to link to")
	}
	adr := filepath.Base(adrs[0])
	// the item links out to a sibling, an ADR, itself, a missing file and an anchor
	write(item, string(b)+
		"\nSee [sib](261015-c4d5-machine-readable-listings.md#top), [adr](../adr/"+adr+"), [me](./261020-a2b3-add-json-output-to-ls.md), [gone](missing.md), [up](#x).\r\n"+
		"```\n[code](261015-c4d5-machine-readable-listings.md)\n```\n\n## Outcome\nShipped.\n")
	// other files link to it
	write("docs/work/other.md", "---\nid: x\n---\n[a](261020-a2b3-add-json-output-to-ls.md#s) `[b](261020-a2b3-add-json-output-to-ls.md)`\n[c](<./261020-a2b3-add-json-output-to-ls.md> \"t\")\n```\n[d](261020-a2b3-add-json-output-to-ls.md)\n```\n")
	write("docs/notes/idea.md", "[e](../work/261020-a2b3-add-json-output-to-ls.md)\n")
	write("docs/ROADMAP.md", "# Roadmap\n\n[f](work/261020-a2b3-add-json-output-to-ls.md)\n")

	// a dry run reports the plan and writes nothing
	dry := mustOK(t, root, "close", Request{"id": a2b3, "status": "done", "dry_run": true})
	if _, ok := dry.Get("links_rewritten"); !ok {
		t.Fatalf("dry run lists no links: %v", dry)
	}
	if s := read(t, root, "docs/notes/idea.md"); s != "[e](../work/261020-a2b3-add-json-output-to-ls.md)\n" {
		t.Fatalf("dry run wrote: %q", s)
	}

	resp := mustOK(t, root, "close", Request{"id": a2b3, "status": "done"})
	moved := "docs/work/done/261020-a2b3-add-json-output-to-ls.md"
	own := read(t, root, moved)
	for _, want := range []string{
		"[sib](../261015-c4d5-machine-readable-listings.md#top)",
		"[adr](../../adr/" + adr + ")",
		"[me](./261020-a2b3-add-json-output-to-ls.md)",
		"[gone](missing.md)",
		"[up](#x).\r\n",
		"```\n[code](261015-c4d5-machine-readable-listings.md)\n```",
	} {
		if !strings.Contains(own, want) {
			t.Errorf("moved item lacks %q:\n%s", want, own)
		}
	}
	for rel, want := range map[string]string{
		"docs/work/other.md": "---\nid: x\n---\n[a](done/261020-a2b3-add-json-output-to-ls.md#s) `[b](261020-a2b3-add-json-output-to-ls.md)`\n[c](<./done/261020-a2b3-add-json-output-to-ls.md> \"t\")\n```\n[d](261020-a2b3-add-json-output-to-ls.md)\n```\n",
		"docs/notes/idea.md": "[e](../work/done/261020-a2b3-add-json-output-to-ls.md)\n",
		"docs/ROADMAP.md":    "# Roadmap\n\n[f](work/done/261020-a2b3-add-json-output-to-ls.md)\n",
	} {
		if got := read(t, root, rel); got != want {
			t.Errorf("%s:\n got  %q\n want %q", rel, got, want)
		}
	}
	if v, _ := resp.Get("link_errors"); v != nil {
		t.Errorf("link_errors: %v", v)
	}
	c, _ := resp.Get("commit")
	paths, _ := c.(Obj).Get("paths")
	got := map[string]bool{}
	for _, p := range paths.([]any) {
		got[p.(string)] = true
	}
	for _, p := range []string{item, moved, "docs/work/other.md", "docs/notes/idea.md", "docs/ROADMAP.md"} {
		if !got[p] {
			t.Errorf("commit.paths lacks %s: %v", p, paths)
		}
	}
	// nothing is left broken
	ck, _ := Marshal(mustOK(t, root, "check", Request{}))
	if r := string(ck); strings.Count(r, `"L1"`) != 1 || !strings.Contains(r, "link missing.md does not resolve") { // only the link that was already broken
		t.Errorf("check reports broken links after close:\n%s", r)
	}
}

// TestCloseReportsLinkWriteFailure: a file whose links cannot be rewritten does
// not undo the move or stop the other files; it is reported in link_errors.
func TestCloseReportsLinkWriteFailure(t *testing.T) {
	if os.Geteuid() == 0 || runtime.GOOS == "windows" {
		t.Skip("needs a folder that refuses writes")
	}
	root := copyRepo(t, "basic", "")
	item := "docs/work/261020-a2b3-add-json-output-to-ls.md"
	p := filepath.Join(root, item)
	b, _ := os.ReadFile(p)
	os.WriteFile(p, append(b, []byte("\n## Outcome\nShipped.\n")...), 0o644)
	os.MkdirAll(filepath.Join(root, "docs/locked"), 0o755)
	os.WriteFile(filepath.Join(root, "docs/locked/a.md"), []byte("[x](../work/261020-a2b3-add-json-output-to-ls.md)\n"), 0o644)
	os.WriteFile(filepath.Join(root, "docs/b.md"), []byte("[y](work/261020-a2b3-add-json-output-to-ls.md)\n"), 0o644)
	os.Chmod(filepath.Join(root, "docs/locked"), 0o555)
	t.Cleanup(func() { os.Chmod(filepath.Join(root, "docs/locked"), 0o755) })

	resp := mustOK(t, root, "close", Request{"id": a2b3, "status": "done"})
	if _, err := os.Stat(filepath.Join(root, "docs/work/done/261020-a2b3-add-json-output-to-ls.md")); err != nil {
		t.Fatal("the move did not happen")
	}
	ev, _ := resp.Get("link_errors")
	if errs, _ := ev.([]Obj); len(errs) != 1 {
		t.Fatalf("link_errors = %v", ev)
	} else if p, _ := errs[0].Get("path"); p != "docs/locked/a.md" {
		t.Fatalf("link_errors path = %v", p)
	}
	if got := read(t, root, "docs/b.md"); got != "[y](work/done/261020-a2b3-add-json-output-to-ls.md)\n" {
		t.Fatalf("the other file was not rewritten: %q", got)
	}
}
