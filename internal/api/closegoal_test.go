package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// close done is refused while the Goal has an empty box (work item 261031-zek2).
func TestCloseRefusesUncheckedGoal(t *testing.T) {
	root := tickRepo(t, "\n") // boxes: one (empty), two (ticked), three (empty)
	mustOK(t, root, "section", Request{"id": e5f6ID, "name": "Outcome", "text": "Done."})
	_, err := call(t, root, "close", Request{"id": e5f6ID, "status": "done"})
	if err == nil || err.Code != "refused" {
		t.Fatalf("close with empty boxes: %v", err)
	}
	if r, _ := err.Details.Get("reason"); r != "unchecked_goal" {
		t.Errorf("reason = %v", r)
	}
	items, _ := err.Details.Get("items")
	rows, _ := items.([]any)
	if len(rows) != 2 {
		t.Fatalf("items = %v", items)
	}
	first := rows[0].(Obj)
	if p, _ := first.Get("position"); p != 1 {
		t.Errorf("first position = %v", p)
	}
	if tx, _ := first.Get("text"); tx != "one" {
		t.Errorf("first text = %v", tx)
	}
	if !strings.Contains(err.Message, `"three"`) {
		t.Errorf("message does not name each box: %s", err.Message)
	}
	// nothing was written
	if _, serr := os.Stat(filepath.Join(root, e5f6Path)); serr != nil {
		t.Errorf("file moved: %v", serr)
	}
	// dropped is not affected
	if _, err := call(t, root, "close", Request{"id": e5f6ID, "status": "dropped", "dry_run": true}); err != nil {
		t.Errorf("dropped: %v", err)
	}
	// ticking each box with evidence lets done through
	mustOK(t, root, "tick", Request{"id": e5f6ID, "item": 1, "evidence": "TestOne"})
	mustOK(t, root, "tick", Request{"id": e5f6ID, "item": 3, "evidence": "TestThree"})
	mustOK(t, root, "close", Request{"id": e5f6ID, "status": "done"})
}

// A box outside the first run is named without a position, because tick cannot
// address it; a Goal with no boxes closes.
func TestCloseGoalEdges(t *testing.T) {
	root := copyRepo(t, "basic", "")
	p := filepath.Join(root, e5f6Path)
	b, _ := os.ReadFile(p)
	s := strings.Replace(string(b), "- [ ] It is done.", "- [x] a\n\nprose\n\n- [ ] stray", 1)
	os.WriteFile(p, []byte(s), 0o644)
	mustOK(t, root, "section", Request{"id": e5f6ID, "name": "Outcome", "text": "Done."})
	_, err := call(t, root, "close", Request{"id": e5f6ID, "status": "done"})
	if err == nil {
		t.Fatal("stray box did not refuse")
	}
	items, _ := err.Details.Get("items")
	row := items.([]any)[0].(Obj)
	if _, has := row.Get("position"); has {
		t.Errorf("a box outside the run has a position: %v", row)
	}
	mustOK(t, root, "section", Request{"id": e5f6ID, "name": "Goal", "text": "No boxes here."})
	mustOK(t, root, "close", Request{"id": e5f6ID, "status": "done"})
}

// tick --evidence writes the Log line in the same write.
func TestTickEvidence(t *testing.T) {
	for _, eol := range []string{"\n", "\r\n"} {
		root := tickRepo(t, eol)
		mustOK(t, root, "tick", Request{"id": e5f6ID, "item": 1, "evidence": "TestOne  passes", "as_of": "2026-10-29"})
		got := read(t, root, e5f6Path)
		if !strings.Contains(got, "- [x] one"+eol) || !strings.Contains(got, "- 2026-10-29: ticked \"one\": TestOne passes"+eol) {
			t.Errorf("%q: evidence not written:\n%q", eol, got)
		}
		// a second tick appends after the first line; undo says so
		mustOK(t, root, "tick", Request{"id": e5f6ID, "item": 1, "done": false, "evidence": "regressed", "as_of": "2026-10-30"})
		got = read(t, root, e5f6Path)
		if !strings.Contains(got, "- [ ] one") || !strings.Contains(got, "ticked \"one\": TestOne passes"+eol+"- 2026-10-30: unticked \"one\": regressed") {
			t.Errorf("%q: undo evidence:\n%q", eol, got)
		}
		// a box already in the asked state still records the evidence
		mustOK(t, root, "tick", Request{"id": e5f6ID, "item": 2, "evidence": "was already done", "as_of": "2026-10-30"})
		if !strings.Contains(read(t, root, e5f6Path), "ticked \"two\": was already done") {
			t.Errorf("%q: no-op tick lost its evidence", eol)
		}
	}
	root := tickRepo(t, "\n")
	for _, ev := range []string{"a\nb", "  "} {
		if _, err := call(t, root, "tick", Request{"id": e5f6ID, "item": 1, "evidence": ev}); err == nil || err.Code != "invalid_request" {
			t.Errorf("evidence %q: %v", ev, err)
		}
	}
	if !strings.Contains(read(t, root, e5f6Path), "- [ ] one") {
		t.Error("a refused evidence ticked the box")
	}
}
