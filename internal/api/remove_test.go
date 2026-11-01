package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// remove deletes one unchecked box and logs it (work item 261101-zq7c).
func TestRemoveBox(t *testing.T) {
	for _, eol := range []string{"\n", "\r\n"} {
		root := tickRepo(t, eol) // boxes: one (empty), two (ticked), three (empty)
		before := read(t, root, e5f6Path)
		// no confirm: refused, names the box and the loss, writes nothing
		_, err := call(t, root, "remove", Request{"id": e5f6ID, "item": 3, "reason": "dropped from scope"})
		if err == nil || err.Code != "refused" {
			t.Fatalf("%q: no confirm: %v", eol, err)
		}
		if r, _ := err.Details.Get("reason"); r != "confirm_required" {
			t.Errorf("%q: reason = %v", eol, r)
		}
		if !strings.Contains(err.Message, "loses history") || !strings.Contains(err.Message, `"three"`) || !strings.Contains(err.Message, "confirm") {
			t.Errorf("message = %s", err.Message)
		}
		if read(t, root, e5f6Path) != before {
			t.Fatalf("%q: a refused remove wrote", eol)
		}
		// dry_run needs no confirm and writes nothing
		resp := mustOK(t, root, "remove", Request{"id": e5f6ID, "item": 3, "reason": "dropped from scope", "dry_run": true})
		if _, ok := resp.Get("diff"); !ok || read(t, root, e5f6Path) != before {
			t.Errorf("%q: dry_run: diff %v, file changed %v", eol, ok, read(t, root, e5f6Path) != before)
		}
		// confirmed: the line goes, later boxes move up, the Log says why
		resp = mustOK(t, root, "remove", Request{"id": e5f6ID, "item": 1, "reason": "dropped  from scope", "confirm": true, "as_of": "2026-10-29"})
		got := read(t, root, e5f6Path)
		if strings.Contains(got, "- [ ] one") || !strings.Contains(got, "- [x] two"+eol+"- [ ] three"+eol) {
			t.Errorf("%q: box not removed cleanly:\n%q", eol, got)
		}
		if !strings.Contains(got, `- 2026-10-29: removed "one": dropped from scope`+eol) && !strings.HasSuffix(got, `removed "one": dropped from scope`) {
			t.Errorf("%q: no Log line:\n%q", eol, got)
		}
		if strings.ContainsRune(strings.ReplaceAll(got, "\r\n", ""), '\r') || (eol == "\r\n" && strings.Contains(strings.ReplaceAll(got, "\r\n", ""), "\n")) {
			t.Errorf("%q: mixed line endings", eol)
		}
		rm, _ := resp.Get("removed")
		items, _ := rm.(Obj).Get("items")
		rows, _ := items.([]Obj)
		if len(rows) != 2 {
			t.Fatalf("%q: items = %v", eol, items)
		}
		if p, _ := rows[1].Get("position"); p != 2 {
			t.Errorf("%q: renumbered position = %v", eol, p)
		}
		if tx, _ := rows[1].Get("text"); tx != "three" {
			t.Errorf("%q: text = %v", eol, tx)
		}
	}
}

func TestRemoveRefusals(t *testing.T) {
	root := tickRepo(t, "\n")
	code := func(req Request) string {
		req["id"] = e5f6ID
		req["confirm"] = true
		if req["reason"] == nil {
			req["reason"] = "r"
		}
		_, err := call(t, root, "remove", req)
		if err == nil {
			return "ok"
		}
		if r, _ := err.Details.Get("reason"); r != nil {
			return err.Code + "/" + r.(string)
		}
		return err.Code
	}
	if c := code(Request{"item": 2}); c != "refused/ticked_box" {
		t.Errorf("ticked box: %s", c)
	}
	if c := code(Request{"item": 9}); c != "not_found" {
		t.Errorf("missing item: %s", c)
	}
	if c := code(Request{"item": 1, "section": "nonesuch"}); c != "not_found" {
		t.Errorf("missing section: %s", c)
	}
	if c := code(Request{"item": 1, "reason": "a\nb"}); c != "invalid_request" {
		t.Errorf("multi-line reason: %s", c)
	}
	if c := code(Request{"item": 1, "reason": "  "}); c != "invalid_request" {
		t.Errorf("blank reason: %s", c)
	}
	if c := code(Request{"item": 1, "section_etag": "0000000000000000"}); c != "stale" {
		t.Errorf("stale section: %s", c)
	}
	if _, err := call(t, root, "remove", Request{"id": e5f6ID, "item": 1, "confirm": true}); err == nil || err.Code != "invalid_request" {
		t.Errorf("no reason: %v", err)
	}
	// a closed item is not edited
	if _, err := call(t, root, "remove", Request{"id": "261010-j9k0", "item": 1, "reason": "r", "confirm": true}); err == nil || err.Code != "refused" {
		t.Errorf("closed item: %v", err)
	}
	// a structural construct stops the write
	p := filepath.Join(root, e5f6Path)
	b, _ := os.ReadFile(p)
	os.WriteFile(p, append(b, []byte("\n<div>\n")...), 0o644)
	if c := code(Request{"item": 1}); c != "refused/unsupported_markdown" {
		t.Errorf("structural: %s", c)
	}
}

// With every box removed and the Log kept, the Goal closes: removing is the way out.
func TestRemoveThenClose(t *testing.T) {
	root := tickRepo(t, "\n")
	mustOK(t, root, "section", Request{"id": e5f6ID, "name": "Outcome", "text": "Done."})
	_, err := call(t, root, "close", Request{"id": e5f6ID, "status": "done"})
	if err == nil || !strings.Contains(err.Message, "remove") {
		t.Fatalf("close message does not name remove: %v", err)
	}
	mustOK(t, root, "remove", Request{"id": e5f6ID, "item": 3, "reason": "gone", "confirm": true})
	mustOK(t, root, "remove", Request{"id": e5f6ID, "item": 1, "reason": "gone", "confirm": true})
	mustOK(t, root, "close", Request{"id": e5f6ID, "status": "done"})
}

func TestRemoveLastLineWithoutEnding(t *testing.T) {
	root := tickRepo(t, "\n")
	p := filepath.Join(root, e5f6Path)
	b, _ := os.ReadFile(p)
	s := strings.Replace(string(b), "## Goal", "## Goal", 1)
	// make the Goal's last box the last line of the file
	i := strings.Index(s, "- [ ] three")
	os.WriteFile(p, []byte(s[:i]+"- [ ] three"), 0o644)
	mustOK(t, root, "remove", Request{"id": e5f6ID, "item": 3, "reason": "r", "confirm": true})
	got := read(t, root, e5f6Path)
	if strings.Contains(got, "three\"") && !strings.Contains(got, `removed "three": r`) {
		t.Errorf("%q", got)
	}
	if strings.Contains(got, "- [ ] three") {
		t.Errorf("box kept:\n%s", got)
	}
}
