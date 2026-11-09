package api

import (
	"fmt"
	"strings"
	"testing"
)

func previewEntry(t *testing.T, resp Obj, n int) Obj {
	t.Helper()
	ws, _ := resp.Get("writes")
	l, ok := ws.([]Obj)
	if !ok || n >= len(l) {
		t.Fatalf("writes %v has no entry %d", ws, n)
	}
	return l[n]
}

func diffText(t *testing.T, e Obj) string {
	t.Helper()
	d, ok := e.Get("diff")
	if !ok {
		t.Fatalf("no diff in %v", e)
	}
	return fmt.Sprint(d)
}

// Each diff is as it applies after the earlier writes, and the repository is not
// touched.
func TestPreviewShowsEachWriteAfterTheEarlierOnes(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "claims", "")
	before := tree(t, root)
	rec := mustOK(t, root, "show", Request{"id": e5f6})
	r, _ := rec.Get("record")
	etag, _ := r.(Obj).Get("etag")

	resp := mustOK(t, root, "preview", Request{"as_of": "2026-11-09T10:00", "writes": []any{
		map[string]any{"op": "section", "id": e5f6, "name": "Goal", "text": "First.", "if_match": etag},
		map[string]any{"op": "section", "id": e5f6, "name": "Goal", "text": "Second."},
		map[string]any{"op": "claim", "id": e5f6, "as": "agent-1"},
		map[string]any{"op": "set", "id": e5f6, "changes": map[string]any{"priority": 2}},
	}})
	if ws, _ := resp.Get("writes"); len(ws.([]Obj)) != 4 {
		t.Fatalf("writes %v", ws)
	}
	// write 1 replaces the original text; write 2 replaces write 1's
	if d := diffText(t, previewEntry(t, resp, 0)); !strings.Contains(d, "Do the thing.") || !strings.Contains(d, "First.") {
		t.Errorf("write 1: %s", d)
	}
	d2 := diffText(t, previewEntry(t, resp, 1))
	if !strings.Contains(d2, "First.") || !strings.Contains(d2, "Second.") || strings.Contains(d2, "Do the thing.") {
		t.Errorf("write 2 is not against the file after write 1: %s", d2)
	}
	// the claim sees the section text of write 2 as context and sets the status
	if d := diffText(t, previewEntry(t, resp, 2)); !strings.Contains(d, "claimed_by") || !strings.Contains(d, "active") {
		t.Errorf("write 3: %s", d)
	}
	// a set after the claim is against the claimed file
	if d := diffText(t, previewEntry(t, resp, 3)); !strings.Contains(d, "priority") {
		t.Errorf("write 4: %s", d)
	}
	if !sameTree(before, tree(t, root)) {
		t.Fatal("preview changed the repository")
	}
}

func TestPreviewFailureNamesTheWrite(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "claims", "")
	// a stale first write keeps its code and says where it failed
	_, err := call(t, root, "preview", Request{"writes": []any{
		map[string]any{"op": "section", "id": e5f6, "name": "Goal", "text": "x", "if_match": "0000000000000000"},
	}})
	if err == nil || err.Code != "stale" {
		t.Fatalf("got %v, want stale", err)
	}
	if w, _ := err.Details.Get("write"); w != 0 {
		t.Fatalf("details %v", err.Details)
	}
	// a later write that fails is named by its position
	_, err = call(t, root, "preview", Request{"writes": []any{
		map[string]any{"op": "section", "id": e5f6, "name": "Goal", "text": "x"},
		map[string]any{"op": "set", "id": e5f6, "changes": map[string]any{"id": "261022-zzzz"}},
	}})
	if err == nil {
		t.Fatal("setting the id was accepted")
	}
	if w, _ := err.Details.Get("write"); w != 1 {
		t.Fatalf("details %v", err.Details)
	}
	// an unknown op and an empty list are refused before anything is copied
	_, err = call(t, root, "preview", Request{"writes": []any{map[string]any{"op": "index"}}})
	if err == nil || err.Code != "invalid_request" {
		t.Fatalf("got %v", err)
	}
	_, err = call(t, root, "preview", Request{"writes": []any{}})
	if err == nil || err.Code != "invalid_request" {
		t.Fatalf("got %v", err)
	}
}

// A new ADR and a set that names it by "$new": the set sees the id the new write
// made on the scratch copy, and nothing is made for real.
func TestPreviewNewRecordAndItsPlaceholder(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "claims", "")
	before := tree(t, root)
	resp := mustOK(t, root, "preview", Request{"as_of": "2026-11-09", "writes": []any{
		map[string]any{"op": "new", "kind": "adr", "slug": "use-a-scratch-copy", "title": "Use a scratch copy"},
		map[string]any{"op": "set", "id": e5f6, "changes": map[string]any{"decisions": []any{"$new"}}},
	}})
	first := previewEntry(t, resp, 0)
	if _, ok := first.Get("note"); !ok {
		t.Error("the new write has no note about its provisional id")
	}
	if d := diffText(t, first); !strings.Contains(d, "Use a scratch copy") {
		t.Errorf("new: %s", d)
	}
	rec, _ := first.Get("record")
	id, _ := rec.(Obj).Get("id")
	if d := diffText(t, previewEntry(t, resp, 1)); strings.Contains(d, "$new") || !strings.Contains(d, id.(string)) {
		t.Errorf("set did not get the provisional id %v: %s", id, d)
	}
	if !sameTree(before, tree(t, root)) {
		t.Fatal("preview changed the repository")
	}
}
