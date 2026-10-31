package api

import "testing"

func TestCloseReportsDependents(t *testing.T) {
	root := copyRepo(t, "closable", "")
	// 261024-g7h8 is blocked_by 261020-a2b3
	resp := mustOK(t, root, "close", Request{"id": "261020-a2b3", "status": "dropped", "as_of": "2026-10-28"})
	v, ok := resp.Get("dependents")
	if !ok || len(v.([]Obj)) != 1 {
		t.Fatalf("dependents = %v", v)
	}
	if id, _ := v.([]Obj)[0].Get("id"); id != "261024-g7h8" {
		t.Errorf("id = %v", id)
	}
	// it reports and edits nothing: the dependent is as it was
	after := mustOK(t, root, "show", Request{"id": "261024-g7h8"})
	rec, _ := after.Get("record")
	f, _ := rec.(Obj).Get("fields")
	if st, _ := f.(Obj).Get("status"); st != "blocked" {
		t.Errorf("dependent status = %v", st)
	}
}

func TestCloseWithNoDependentsHasNoMember(t *testing.T) {
	root := copyRepo(t, "closable", "")
	resp := mustOK(t, root, "close", Request{"id": "261022-e5f6", "status": "dropped", "outcome": "Not needed.", "as_of": "2026-10-28"})
	if _, ok := resp.Get("dependents"); ok {
		t.Errorf("dependents present: %v", resp)
	}
}

func TestCloseIgnoresClosedDependents(t *testing.T) {
	root := copyRepo(t, "closable", "")
	mustOK(t, root, "section", Request{"id": "261024-g7h8", "name": "outcome", "text": "Moot."})
	mustOK(t, root, "close", Request{"id": "261024-g7h8", "status": "dropped", "as_of": "2026-10-28"})
	resp := mustOK(t, root, "close", Request{"id": "261020-a2b3", "status": "dropped", "as_of": "2026-10-28"})
	if _, ok := resp.Get("dependents"); ok {
		t.Errorf("a dropped item counted as a dependent: %v", resp)
	}
}
