package api

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/record"
)

// The matrix of work item 261031-z0n8, cell by cell. "corr" is no with
// correction true; "append" is yes with append true.
var editMatrix = []struct {
	name                                string
	section                             bool
	raw, ready, active, blocked, closed string
}{
	{"title", false, "yes", "yes", "yes", "yes", "corr"},
	{"summary", false, "yes", "yes", "yes", "yes", "corr"},
	{"kind", false, "yes", "yes", "warn", "warn", "no"},
	{"priority", false, "yes", "yes", "warn", "yes", "no"},
	{"parent", false, "yes", "yes", "warn", "warn", "no"},
	{"blocked_by", false, "yes", "yes", "yes", "yes", "no"},
	{"cycle", false, "yes", "yes", "yes", "yes", "no"},
	{"keep", false, "yes", "yes", "yes", "yes", "yes"},
	{"decisions", false, "yes", "yes", "yes", "yes", "no"},
	{"authored_by", false, "yes", "yes", "yes", "yes", "no"},
	{"status", false, "yes", "yes", "yes", "yes", "no"},
	{"project", false, "no", "no", "no", "no", "no"},
	{"id", false, "no", "no", "no", "no", "no"},
	{"created", false, "no", "no", "no", "no", "no"},
	{"updated", false, "no", "no", "no", "no", "no"},
	{"closed", false, "no", "no", "no", "no", "no"},
	{"claimed_by", false, "no", "no", "no", "no", "no"},
	{"claimed_at", false, "no", "no", "no", "no", "no"},
	{"goal", true, "yes", "yes", "warn", "yes", "corr"},
	{"next", true, "yes", "yes", "yes", "yes", "no"},
	{"context", true, "yes", "yes", "yes", "yes", "corr"},
	{"outcome", true, "yes", "yes", "yes", "yes", "corr"},
	{"log", true, "append", "append", "append", "append", "no"},
}

func cellOf(s editState) string {
	switch {
	case s.State == editNo && s.Correction:
		return "corr"
	case s.State == editYes && s.Append:
		return "append"
	}
	return s.State
}

func TestEditStatesFollowTheMatrix(t *testing.T) {
	t.Parallel()
	seen := map[string]bool{}
	for _, row := range editMatrix {
		seen[row.name] = true
		want := map[string]string{"raw": row.raw, "ready": row.ready, "active": row.active, "blocked": row.blocked, "done": row.closed, "dropped": row.closed}
		for _, status := range workStatus {
			closed := status == "done" || status == "dropped"
			var s editState
			if row.section {
				s = sectionEditState(row.name, status, closed)
			} else {
				s = fieldEditState(row.name, status, closed)
			}
			if got := cellOf(s); got != want[status] {
				t.Errorf("%s in %s: got %s, want %s", row.name, status, got, want[status])
			}
			if (s.State == editYes) != (s.Reason == "") {
				t.Errorf("%s in %s: state %s with reason %q (a reason goes with warn and no only)", row.name, status, s.State, s.Reason)
			}
		}
	}
	// Every field the tool knows for a work item has a row.
	for _, f := range record.TypeFields(record.Work, true) {
		if !seen[f] {
			t.Errorf("field %s has no row in the matrix", f)
		}
	}
	if len(workStatus) != 6 {
		t.Fatalf("the matrix covers six statuses, the tool has %v", workStatus)
	}
}

// What edit says agrees with what set and section do on a closed item.
func TestEditStatesAgreeWithTheWritesOnAClosedItem(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "basic", "")
	p := filepath.Join(root, "docs/work/261020-a2b3-add-json-output-to-ls.md")
	b, _ := os.ReadFile(p)
	os.WriteFile(p, append(b, []byte("\n## Outcome\nShipped.\n")...), 0o644)
	mustOK(t, root, "close", Request{"id": a2b3, "status": "done", "as_of": "2026-10-29"})
	resp := mustOK(t, root, "show", Request{"id": a2b3})
	ev, ok := resp.Get("edit")
	if !ok {
		t.Fatal("show has no edit member")
	}
	fields, _ := ev.(Obj).Get("fields")
	values := map[string]any{"title": "New title", "summary": "New summary.", "kind": "bug", "priority": 3, "parent": nil, "blocked_by": nil, "cycle": "2026-W45", "keep": true, "decisions": nil, "authored_by": "someone"}
	for _, kv := range fields.(Obj) {
		v, settable := values[kv.K]
		if !settable {
			continue
		}
		st, _ := kv.V.(Obj).Get("state")
		_, corr := kv.V.(Obj).Get("correction")
		_, plain := call(t, root, "set", Request{"id": a2b3, "changes": map[string]any{kv.K: v}, "dry_run": true})
		_, withReason := call(t, root, "set", Request{"id": a2b3, "changes": map[string]any{kv.K: v}, "correction": "typo", "dry_run": true})
		if (st == "no") != (plain != nil) {
			t.Errorf("%s: state %v but a plain set answered %v", kv.K, st, plain)
		}
		if st == "no" && corr != (withReason == nil) {
			t.Errorf("%s: correction %v but a set with a correction answered %v", kv.K, corr, withReason)
		}
	}
	sections, _ := ev.(Obj).Get("sections")
	for _, name := range []string{"goal", "next", "outcome", "log"} {
		if _, ok := sections.(Obj).Get(name); !ok {
			t.Errorf("edit.sections has no %s", name)
		}
	}
}

func TestShowOfAnADRHasNoEdit(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "basic", "")
	repo, err := record.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range repo.Records {
		if r.Type == record.ADR {
			resp := mustOK(t, root, "show", Request{"id": r.ID})
			if _, ok := resp.Get("edit"); ok {
				t.Fatalf("show of ADR %s carries edit", r.ID)
			}
			return
		}
	}
	t.Skip("the fixture has no ADR")
}
