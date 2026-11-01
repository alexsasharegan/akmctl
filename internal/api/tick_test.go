package api

import (
	mdbody "github.com/alexsasharegan/akmctl/internal/body"
	"github.com/alexsasharegan/akmctl/internal/record"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// tickRepo is the basic repository with a three-box Goal on e5f6.
func tickRepo(t *testing.T, eol string) string {
	root := copyRepo(t, "basic", "")
	p := filepath.Join(root, e5f6Path)
	b, _ := os.ReadFile(p)
	s := strings.Replace(string(b), "- [ ] It is done.", "- [ ] one\n- [x] two\n- [ ] three", 1)
	if eol == "\r\n" {
		s = strings.ReplaceAll(s, "\n", "\r\n")
	}
	os.WriteFile(p, []byte(s), 0o644)
	return root
}

func sectionEtag(t *testing.T, root, name string) string {
	t.Helper()
	repo, err := record.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	rec, aerr := lookup(repo, e5f6ID)
	if aerr != nil {
		t.Fatal(aerr)
	}
	tr, serr := mdbody.Split([]byte(rec.Body))
	if serr != nil {
		t.Fatal(serr)
	}
	sec, ok := tr.Section(name)
	if !ok {
		t.Fatalf("no section %s", name)
	}
	return tr.Etag(sec)
}

// differing counts the bytes at which two texts of one length differ.
func differing(a, b string) int {
	if len(a) != len(b) {
		return -1
	}
	n := 0
	for i := range a {
		if a[i] != b[i] {
			n++
		}
	}
	return n
}

func TestTickChangesOneByte(t *testing.T) {
	for _, eol := range []string{"\n", "\r\n"} {
		root := tickRepo(t, eol)
		before := read(t, root, e5f6Path)
		resp := mustOK(t, root, "tick", Request{"id": e5f6ID, "item": 3, "as_of": "2026-10-29"})
		after := read(t, root, e5f6Path)
		// one box, and updated (same length: a date for a date)
		want := strings.Replace(before, "updated: 2026-10-22", "updated: 2026-10-29", 1)
		if d := differing(want, after); d != 1 {
			t.Errorf("%q: %d bytes differ beyond updated:\n%s", eol, d, after)
		}
		if !strings.Contains(after, "[x] three") {
			t.Errorf("%q: not ticked", eol)
		}
		if ch, _ := resp.Get("changed"); ch != true {
			t.Errorf("changed = %v", ch)
		}
		it, _ := resp.Get("item")
		if d, _ := it.(Obj).Get("done"); d != true {
			t.Errorf("item = %v", it)
		}
		// undo, and a no-op writes nothing and leaves updated alone
		mustOK(t, root, "tick", Request{"id": e5f6ID, "item": 3, "done": false, "as_of": "2026-10-29"})
		resp = mustOK(t, root, "tick", Request{"id": e5f6ID, "item": 3, "done": false, "as_of": "2026-11-05"})
		if ch, _ := resp.Get("changed"); ch != false || strings.Contains(read(t, root, e5f6Path), "2026-11-05") {
			t.Errorf("a tick that changes nothing wrote the file")
		}
	}
}

func TestLogAddsOneLine(t *testing.T) {
	for _, eol := range []string{"\n", "\r\n"} {
		root := tickRepo(t, eol)
		mustOK(t, root, "log", Request{"id": e5f6ID, "text": "first", "as_of": "2026-10-29"})
		mustOK(t, root, "log", Request{"id": e5f6ID, "text": "second  line", "as_of": "2026-10-30"})
		got := read(t, root, e5f6Path)
		want := "## Log" + eol + eol + "- 2026-10-29: first" + eol + "- 2026-10-30: second line" + eol
		if !strings.HasSuffix(got, want) && !strings.Contains(got, "- 2026-10-29: first"+eol+"- 2026-10-30: second line"+eol) {
			t.Errorf("%q: log lines:\n%q", eol, got)
		}
		if !strings.Contains(got, "updated: 2026-10-30") {
			t.Errorf("updated not set")
		}
	}
}

func TestTickAndLogRefusals(t *testing.T) {
	root := tickRepo(t, "\n")
	code := func(op string, req Request) string {
		_, err := call(t, root, op, req)
		if err == nil {
			return "ok"
		}
		return err.Code
	}
	if c := code("tick", Request{"id": e5f6ID, "item": 9}); c != "not_found" {
		t.Errorf("missing item: %s", c)
	}
	if c := code("tick", Request{"id": e5f6ID, "item": 1, "section": "nonesuch"}); c != "not_found" {
		t.Errorf("missing section: %s", c)
	}
	if c := code("tick", Request{"id": "261010-j9k0", "item": 1}); c != "refused" {
		t.Errorf("closed item: %s", c)
	}
	if c := code("log", Request{"id": "261010-j9k0", "text": "x"}); c != "refused" {
		t.Errorf("closed item, log: %s", c)
	}
	if c := code("log", Request{"id": e5f6ID, "text": "a\nb"}); c != "invalid_request" {
		t.Errorf("multi-line log: %s", c)
	}
	if c := code("tick", Request{"id": e5f6ID, "item": 1, "if_match": "0000000000000000"}); c != "stale" {
		t.Errorf("stale file etag: %s", c)
	}
	if c := code("tick", Request{"id": e5f6ID, "item": 0}); c != "invalid_request" {
		t.Errorf("item 0: %s", c)
	}
}

func TestSectionEtag(t *testing.T) {
	root := tickRepo(t, "\n")
	mustOK(t, root, "log", Request{"id": e5f6ID, "text": "start"}) // gives the file a Log, so a new section lands before it
	et := sectionEtag(t, root, "goal")
	// a change elsewhere in the file does not stale the Goal's etag
	mustOK(t, root, "section", Request{"id": e5f6ID, "name": "Notes", "text": "n"})
	mustOK(t, root, "tick", Request{"id": e5f6ID, "item": 1, "section_etag": et})
	// but a change to the Goal does, and the answer carries the current etag
	_, err := call(t, root, "tick", Request{"id": e5f6ID, "item": 3, "section_etag": et})
	if err == nil || err.Code != "stale" {
		t.Fatalf("%v", err)
	}
	if got, _ := err.Details.Get("section_etag"); got != sectionEtag(t, root, "goal") {
		t.Errorf("stale answered with %v", got)
	}
	// the Log etag is checked by log
	le := sectionEtag(t, root, "log")
	mustOK(t, root, "log", Request{"id": e5f6ID, "text": "a", "section_etag": le})
	if _, err := call(t, root, "log", Request{"id": e5f6ID, "text": "b", "section_etag": le}); err == nil || err.Code != "stale" {
		t.Errorf("log with a stale etag: %v", err)
	}
}

func TestTickRefusesUnreadableBody(t *testing.T) {
	root := tickRepo(t, "\n")
	p := filepath.Join(root, e5f6Path)
	b, _ := os.ReadFile(p)
	// a local construct in Notes does not stop a tick in Goal; in Goal it does
	os.WriteFile(p, append(b, []byte("\n## Notes\n\n*\tordered\n")...), 0o644)
	mustOK(t, root, "tick", Request{"id": e5f6ID, "item": 1})
	_, err := call(t, root, "tick", Request{"id": e5f6ID, "item": 1, "section": "notes"})
	if err == nil || err.Code != "refused" {
		t.Fatalf("%v", err)
	}
	if r, _ := err.Details.Get("reason"); r != "unsupported_markdown" {
		t.Errorf("details = %v", err.Details)
	}
	// a structural construct stops every write
	os.WriteFile(p, append(b, []byte("\n<div>\n")...), 0o644)
	if _, err := call(t, root, "tick", Request{"id": e5f6ID, "item": 1}); err == nil || err.Code != "refused" {
		t.Fatalf("structural: %v", err)
	}
	if _, err := call(t, root, "log", Request{"id": e5f6ID, "text": "x"}); err == nil || err.Code != "refused" {
		t.Fatalf("structural log: %v", err)
	}
}
