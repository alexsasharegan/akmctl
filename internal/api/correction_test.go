package api

import (
	"strings"
	"testing"
)

const closedPath = "docs/work/done/261010-j9k0-choose-a-licence.md"

// A closed work item is read-only except through a logged correction
// (ADR 261101-dxcn, work item 261101-3vwz).
func TestClosedItemCorrection(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "basic", "")
	id := "261010-j9k0"
	closedReason := func(op string, req Request) string {
		_, err := call(t, root, op, req)
		if err == nil {
			return "ok"
		}
		r, _ := err.Details.Get("reason")
		s, _ := r.(string)
		return err.Code + "/" + s
	}
	// without a correction: refused for every field but keep
	for _, f := range []string{"summary", "title", "priority", "kind", "parent", "cycle", "blocked_by", "decisions"} {
		if c := closedReason("set", Request{"id": id, "changes": map[string]any{f: "x"}}); c != "refused/closed" && c != "invalid_request/" {
			t.Errorf("set %s without a correction: %s", f, c)
		}
	}
	if c := closedReason("set", Request{"id": id, "changes": map[string]any{"summary": "x"}}); c != "refused/closed" {
		t.Errorf("summary: %s", c)
	}
	if c := closedReason("section", Request{"id": id, "name": "Outcome", "text": "x"}); c != "refused/closed" {
		t.Errorf("section: %s", c)
	}
	if c := closedReason("log", Request{"id": id, "text": "x"}); c != "refused/closed" {
		t.Errorf("log: %s", c)
	}
	// keep needs no correction and writes no Log line
	mustOK(t, root, "set", Request{"id": id, "changes": map[string]any{"keep": true}, "as_of": "2026-10-28"})
	if got := read(t, root, closedPath); !strings.Contains(got, "keep: true") || strings.Contains(got, "## Log") {
		t.Errorf("keep on a closed item:\n%s", got)
	}
	// with a correction: title and summary only, one Log line naming them
	if c := closedReason("set", Request{"id": id, "changes": map[string]any{"priority": 1}, "correction": "r"}); c != "refused/closed_field" {
		t.Errorf("priority with a correction: %s", c)
	}
	if c := closedReason("set", Request{"id": id, "changes": map[string]any{"summary": "x"}, "correction": "a\nb"}); c != "invalid_request/" {
		t.Errorf("multi-line correction: %s", c)
	}
	mustOK(t, root, "set", Request{"id": id, "changes": map[string]any{"summary": "Chose MIT.", "title": "Pick a license"}, "correction": "wording", "as_of": "2026-10-28"})
	got := read(t, root, closedPath)
	if !strings.Contains(got, "summary: Chose MIT.") || !strings.Contains(got, "- 2026-10-28: corrected summary and title: wording") {
		t.Errorf("correction not applied or not logged:\n%s", got)
	}
	// the same value again changes nothing and logs nothing
	mustOK(t, root, "set", Request{"id": id, "changes": map[string]any{"summary": "Chose MIT."}, "correction": "again", "as_of": "2026-10-29"})
	if strings.Contains(read(t, root, closedPath), "again") {
		t.Error("a no-op correction wrote a Log line")
	}
	// a section: any but the Log; the line names it; a no-op writes nothing
	if c := closedReason("section", Request{"id": id, "name": "Log", "text": "- x", "mode": "append", "correction": "r"}); c != "refused/closed_field" {
		t.Errorf("log with a correction: %s", c)
	}
	mustOK(t, root, "section", Request{"id": id, "name": "Outcome", "text": "MIT, in `LICENSE` and the header.", "correction": "header missed", "as_of": "2026-10-28"})
	got = read(t, root, closedPath)
	if !strings.Contains(got, "header.") || !strings.Contains(got, "- 2026-10-28: corrected outcome: header missed") {
		t.Errorf("section correction:\n%s", got)
	}
	n := strings.Count(got, "corrected")
	mustOK(t, root, "section", Request{"id": id, "name": "Outcome", "text": "MIT, in `LICENSE` and the header.", "correction": "again"})
	if strings.Count(read(t, root, closedPath), "corrected") != n {
		t.Error("a no-op section correction wrote a Log line")
	}
	// the item stays closed and valid
	if c := closedReason("set", Request{"id": id, "changes": map[string]any{"status": "ready"}, "correction": "r"}); c != "refused/rule" {
		t.Errorf("reopen: %s", c)
	}
	// not clean(): a corrected title no longer matches the heading (WK12, a warning, as for an open item)
}

// On an open item a correction is not needed, and when given it does nothing.
func TestCorrectionOnAnOpenItemIsIgnored(t *testing.T) {
	root := tickRepo(t, "\n")
	mustOK(t, root, "set", Request{"id": e5f6ID, "changes": map[string]any{"summary": "s"}, "correction": "why"})
	if strings.Contains(read(t, root, e5f6Path), "corrected") {
		t.Error("an open item got a correction line")
	}
}
