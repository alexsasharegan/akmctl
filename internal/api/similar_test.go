package api

import (
	"reflect"
	"testing"
)

func similarIDs(t *testing.T, resp Obj) []string {
	t.Helper()
	v, ok := resp.Get("similar")
	if !ok {
		return nil
	}
	var ids []string
	for _, o := range v.([]Obj) {
		id, _ := o.Get("id")
		ids = append(ids, id.(string))
	}
	return ids
}

func TestNewWorkReportsSimilarOpenItems(t *testing.T) {
	root := copyRepo(t, "basic", "")
	// "Add JSON output to ls" is active; "Document the exit codes" is ready.
	resp := mustOK(t, root, "new", Request{"kind": "work", "slug": "json-output-for-board"})
	ids := similarIDs(t, resp)
	if len(ids) != 1 {
		t.Fatalf("want one similar item, got %v", ids)
	}
	v, _ := resp.Get("similar")
	o := v.([]Obj)[0]
	if w, _ := o.Get("words"); !reflect.DeepEqual(w, []any{"json", "output"}) {
		t.Errorf("words = %v", w)
	}
	if st, _ := o.Get("status"); st != "active" {
		t.Errorf("status = %v", st)
	}
}

func TestNewWorkNoMatchHasNoSimilar(t *testing.T) {
	root := copyRepo(t, "basic", "")
	resp := mustOK(t, root, "new", Request{"kind": "work", "slug": "rotate-the-signing-keys"})
	if _, ok := resp.Get("similar"); ok {
		t.Errorf("similar present for no match: %v", resp)
	}
}

func TestNewSimilarOneWordIsNotEnough(t *testing.T) {
	root := copyRepo(t, "basic", "")
	// shares only "output" with the active item
	resp := mustOK(t, root, "new", Request{"kind": "work", "slug": "colour-output"})
	if ids := similarIDs(t, resp); ids != nil {
		t.Errorf("one shared word matched: %v", ids)
	}
}

func TestNewSimilarSingleWordItemNeedsThatWord(t *testing.T) {
	root := copyRepo(t, "basic", "")
	// fewer than two words: every word must be shared ("a" is a stop word)
	resp := mustOK(t, root, "new", Request{"kind": "work", "slug": "a-plugin"})
	if ids := similarIDs(t, resp); len(ids) != 1 {
		t.Errorf("want the plugin item, got %v", ids)
	}
}

func TestNewSimilarSkipsClosedAndAdr(t *testing.T) {
	root := copyRepo(t, "basic", "")
	first := mustOK(t, root, "new", Request{"kind": "work", "slug": "audit-the-flux-capacitor"})
	id, _ := first.Get("record")
	rid, _ := id.(Obj).Get("id")
	mustOK(t, root, "section", Request{"id": rid, "name": "outcome", "text": "Not needed."})
	mustOK(t, root, "close", Request{"id": rid, "status": "dropped"})
	resp := mustOK(t, root, "new", Request{"kind": "work", "slug": "audit-flux-capacitor"})
	if ids := similarIDs(t, resp); ids != nil {
		t.Errorf("a dropped item matched: %v", ids)
	}
	adr := mustOK(t, root, "new", Request{"kind": "adr", "slug": "json-output-format"})
	if _, ok := adr.Get("similar"); ok {
		t.Errorf("new adr returned similar")
	}
}

func TestNewSimilarIsCappedAndOrdered(t *testing.T) {
	root := copyRepo(t, "basic", "")
	for i := 0; i < 7; i++ {
		mustOK(t, root, "new", Request{"kind": "work", "slug": "frobnicate-widget", "title": "Frobnicate widget"})
	}
	resp := mustOK(t, root, "new", Request{"kind": "work", "slug": "frobnicate-widget"})
	ids := similarIDs(t, resp)
	if len(ids) != 5 {
		t.Fatalf("want 5, got %d", len(ids))
	}
	for i := 1; i < len(ids); i++ {
		if ids[i-1] >= ids[i] {
			t.Errorf("not ordered by id: %v", ids)
		}
	}
}
