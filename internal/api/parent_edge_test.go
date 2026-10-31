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

func block(t *testing.T, root, id string, by ...string) *Error {
	t.Helper()
	list := make([]any, len(by))
	for i, b := range by {
		list[i] = b
	}
	_, err := call(t, root, "set", Request{"id": id, "changes": map[string]any{"blocked_by": list}})
	return err
}

func TestSetRefusesBlockedByCycles(t *testing.T) {
	root := copyRepo(t, "basic", "")
	a := newItem(t, root, "item-a", "")
	b := newItem(t, root, "item-b", "")
	c := newItem(t, root, "item-c", "")
	if err := block(t, root, a, b); err != nil {
		t.Fatal(err)
	}
	if err := block(t, root, b, c); err != nil {
		t.Fatal(err)
	}
	wantRule(t, block(t, root, a, a), "WK18", a)                 // itself
	wantRule(t, block(t, root, b, a), "WK18", a, b)              // mutual
	wantRule(t, block(t, root, c, a), "WK18", a, b, c)           // longer loop, named in order
	if err := block(t, root, a, b, "legal review"); err != nil { // free text and an entry already there
		t.Fatal(err)
	}
	clean(t, root)
}

func TestBlockedByCyclesThroughTheParentTree(t *testing.T) {
	root := copyRepo(t, "basic", "")
	p := newItem(t, root, "project", "")
	c := newItem(t, root, "child", p)
	g := newItem(t, root, "grandchild", c)
	other := newItem(t, root, "other", "")
	// a parent waits for its children, so a child cannot wait for the parent, at any depth
	wantRule(t, block(t, root, c, p), "WK18", p, c)
	wantRule(t, block(t, root, g, p), "WK18", p, g)
	// through a third item that waits for the project
	if err := block(t, root, other, p); err != nil {
		t.Fatal(err)
	}
	wantRule(t, block(t, root, g, other), "WK18", other, g)
	// a parent listing its own child is redundant, not circular
	if err := block(t, root, p, c); err != nil {
		t.Fatal(err)
	}
	clean(t, root)
}

func TestParentEdgeCannotCompleteADependencyCycle(t *testing.T) {
	root := copyRepo(t, "basic", "")
	p := newItem(t, root, "project", "")
	c := newItem(t, root, "stray", "")
	if err := block(t, root, c, p); err != nil { // c waits for p: fine while unrelated
		t.Fatal(err)
	}
	_, err := call(t, root, "set", Request{"id": c, "changes": map[string]any{"parent": p}})
	wantRule(t, err, "WK18", p, c)
	// moving a child away from the parent it waits on is allowed
	q := newItem(t, root, "elsewhere", "")
	mustOK(t, root, "set", Request{"id": c, "changes": map[string]any{"parent": q}})
	clean(t, root)
}
