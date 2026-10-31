package api

import (
	"strings"
	"testing"
)

const doneItem = "261010-j9k0"

// newItem files a work item and returns its id.
func newItem(t *testing.T, root, slug, parent string) string {
	t.Helper()
	req := Request{"kind": "work", "slug": slug}
	if parent != "" {
		req["fields"] = map[string]any{"parent": parent}
	}
	resp := mustOK(t, root, "new", req)
	rec, _ := resp.Get("record")
	id, _ := rec.(Obj).Get("id")
	return id.(string)
}

func wantRule(t *testing.T, err *Error, rule string, in ...string) {
	t.Helper()
	if err == nil || err.Code != "refused" {
		t.Fatalf("want refused %s, got %v", rule, err)
	}
	if got, _ := err.Details.Get("rule"); got != rule {
		t.Fatalf("rule = %v, want %s (%v)", got, rule, err)
	}
	for _, s := range in {
		if !strings.Contains(err.Message, s) {
			t.Errorf("message %q lacks %q", err.Message, s)
		}
	}
}

func TestSetRefusesParentLoopsAtAnyDepth(t *testing.T) {
	root := copyRepo(t, "basic", "")
	a := newItem(t, root, "tier-one", "")
	b := newItem(t, root, "tier-two", a)
	c := newItem(t, root, "leaf", b)
	d := newItem(t, root, "deeper", c)
	for _, tc := range []struct{ id, parent string }{{a, b}, {a, c}, {a, d}, {b, d}} {
		_, err := call(t, root, "set", Request{"id": tc.id, "changes": map[string]any{"parent": tc.parent}})
		wantRule(t, err, "WK4", tc.id, tc.parent)
	}
	clean(t, root)
	// a legal move still works
	mustOK(t, root, "set", Request{"id": d, "changes": map[string]any{"parent": a}})
	clean(t, root)
}

func TestSetKeepsAnExistingLoopEditable(t *testing.T) {
	// a violation that is already present is not a reason to refuse an unrelated change
	root := copyRepo(t, "basic", "")
	mustOK(t, root, "set", Request{"id": "261015-c4d5", "changes": map[string]any{"priority": 1}})
}

func TestOpenItemCannotJoinAClosedParent(t *testing.T) {
	root := copyRepo(t, "basic", "")
	open := newItem(t, root, "open-one", "")
	_, err := call(t, root, "set", Request{"id": open, "changes": map[string]any{"parent": doneItem}})
	wantRule(t, err, "WK8", open, doneItem)
	_, err = call(t, root, "new", Request{"kind": "work", "slug": "late-child", "fields": map[string]any{"parent": doneItem}})
	wantRule(t, err, "WK8", doneItem)
	clean(t, root)
}

func TestClosedItemMayJoinAClosedParent(t *testing.T) {
	root := copyRepo(t, "basic", "")
	x := newItem(t, root, "finished", "")
	mustOK(t, root, "section", Request{"id": x, "name": "Goal", "text": "Done when it is done."})
	mustOK(t, root, "section", Request{"id": x, "name": "Outcome", "text": "Shipped."})
	mustOK(t, root, "close", Request{"id": x, "status": "done"})
	mustOK(t, root, "set", Request{"id": x, "changes": map[string]any{"parent": doneItem}})
	clean(t, root)
}
