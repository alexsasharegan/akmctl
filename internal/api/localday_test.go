package api

import (
	"strings"
	"testing"
	"time"
)

// An evening in Arizona (MST, UTC-7) is already the next day in UTC: new
// records carry the local calendar day, while claimed_at stays in UTC.
func TestNewAndClaimUseTheLocalDay(t *testing.T) {
	mst := time.FixedZone("MST", -7*3600)
	old := timeNow
	timeNow = func() time.Time { return time.Date(2026, 10, 29, 18, 30, 0, 0, mst) } // 2026-10-30T01:30Z
	t.Cleanup(func() { timeNow = old })

	root := copyRepo(t, "basic", "extension.claims: 0.3.0\n")
	resp := mustOK(t, root, "new", Request{"kind": "work", "slug": "evening-note"})
	rec, _ := resp.Get("record")
	id, _ := rec.(Obj).Get("id")
	if !strings.HasPrefix(id.(string), "261029-") {
		t.Fatalf("id %v does not carry the local day", id)
	}
	path, _ := resp.Get("path")
	text := read(t, root, path.(string))
	for _, want := range []string{"created: 2026-10-29\n", "updated: 2026-10-29\n", "- 2026-10-29: filed."} {
		if !strings.Contains(text, want) {
			t.Errorf("missing %q in\n%s", want, text)
		}
	}
	// as_of still wins over the clock
	resp = mustOK(t, root, "new", Request{"kind": "work", "slug": "dated-note", "as_of": "2026-11-02"})
	rec, _ = resp.Get("record")
	if id, _ = rec.(Obj).Get("id"); !strings.HasPrefix(id.(string), "261102-") {
		t.Fatalf("as_of ignored: %v", id)
	}
	// claim: the log line uses the local day, claimed_at the UTC timestamp
	mustOK(t, root, "claim", Request{"id": e5f6, "as": "agent-1"})
	text = read(t, root, "docs/work/261022-e5f6-document-the-exit-codes.md")
	for _, want := range []string{"claimed_at: 2026-10-30T01:30\n", "- 2026-10-29 18:30-07:00: claimed (agent-1)."} {
		if !strings.Contains(text, want) {
			t.Errorf("missing %q in\n%s", want, text)
		}
	}
}

// A Log line carries the writer's local time and offset; at records when the
// thing happened (ADR 261101-gw7b).
func TestLogStampAndAt(t *testing.T) {
	mst := time.FixedZone("MST", -7*3600)
	old := timeNow
	timeNow = func() time.Time { return time.Date(2026, 11, 1, 16, 2, 30, 0, mst) }
	t.Cleanup(func() { timeNow = old })

	root := tickRepo(t, "\n")
	mustOK(t, root, "log", Request{"id": e5f6ID, "text": "plain"})
	mustOK(t, root, "log", Request{"id": e5f6ID, "text": "same day", "at": "2026-11-01 11:05"})
	mustOK(t, root, "log", Request{"id": e5f6ID, "text": "other zone", "at": "2026-11-01T11:05+01:00"})
	mustOK(t, root, "log", Request{"id": e5f6ID, "text": "earlier day", "at": "2026-10-31T23:50"})
	got := read(t, root, e5f6Path)
	for _, want := range []string{
		"- 2026-11-01 16:02-07:00: plain\n",
		"- 2026-11-01 16:02-07:00: same day (at 11:05-07:00)\n",
		"- 2026-11-01 16:02-07:00: other zone (at 11:05+01:00)\n",
		"- 2026-11-01 16:02-07:00: earlier day (at 2026-10-31 23:50-07:00)\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in\n%s", want, got)
		}
	}
	for _, at := range []string{"2026-11-01 16:03", "2026-11-02 00:00", "2026-11-01T23:30+00:00"} {
		_, err := call(t, root, "log", Request{"id": e5f6ID, "text": "x", "at": at})
		if err == nil || err.Code != "invalid_request" {
			t.Errorf("at %s after the write: %v", at, err)
		}
	}
	for _, at := range []string{"yesterday", "2026-13-01 10:00", "2026-11-01"} {
		_, err := call(t, root, "log", Request{"id": e5f6ID, "text": "x", "at": at})
		if err == nil || err.Code != "invalid_request" {
			t.Errorf("bad at %s: %v", at, err)
		}
	}
}
