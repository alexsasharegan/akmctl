package api

import (
	"strings"
	"testing"
)

// texts is the content of every file under root, without the lock file and
// without modification times: a file put back by a rollback is written again.
func texts(t *testing.T, root string) map[string]string {
	t.Helper()
	m := map[string]string{}
	for k, v := range tree(t, root) {
		m[k] = v.text
	}
	return m
}

func sameTexts(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if w, ok := b[k]; !ok || v != w {
			return false
		}
	}
	return true
}

// A sequence is made in order under one lock: a later write to a record needs no
// etag of its own, and each entry is the response of that write.
func TestApplyMakesTheWritesInOrder(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "claims", "")
	rec := mustOK(t, root, "show", Request{"id": e5f6})
	r, _ := rec.Get("record")
	etag, _ := r.(Obj).Get("etag")

	resp := mustOK(t, root, "apply", Request{"as_of": "2026-11-16T10:00", "writes": []any{
		map[string]any{"op": "section", "id": e5f6, "name": "Goal", "text": "First.", "if_match": etag},
		map[string]any{"op": "section", "id": e5f6, "name": "Goal", "text": "Second."},
		map[string]any{"op": "claim", "id": e5f6, "as": "agent-1"},
	}})
	ws, _ := resp.Get("writes")
	if len(ws.([]Obj)) != 3 {
		t.Fatalf("writes %v", ws)
	}
	if op, _ := previewEntry(t, resp, 2).Get("operation"); op != "claim" {
		t.Errorf("entry 3 is %v", op)
	}
	if _, ok := previewEntry(t, resp, 2).Get("commit"); !ok {
		t.Error("the claim's entry has no commit")
	}
	got := mustOK(t, root, "show", Request{"id": e5f6})
	gp, _ := got.Get("record")
	path, _ := gp.(Obj).Get("path")
	if s := read(t, root, path.(string)); !strings.Contains(s, "Second.") || strings.Contains(s, "First.") {
		t.Errorf("file after the sequence: %s", s)
	}
	gr, _ := got.Get("record")
	f, _ := gr.(Obj).Get("fields")
	if by, _ := f.(Obj).Get("claimed_by"); by != "agent-1" {
		t.Errorf("claimed_by %v", by)
	}
}

// A write that fails undoes the earlier ones: the new ADR is gone, the closed
// item is back in its folder and the links to it are as they were.
func TestApplyRollsBackWhenAWriteFails(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "closelinks", "")
	before := texts(t, root)

	_, err := call(t, root, "apply", Request{"as_of": "2026-11-16T10:00", "writes": []any{
		map[string]any{"op": "new", "kind": "adr", "slug": "use-json", "title": "Use JSON"},
		map[string]any{"op": "set", "id": a2b3, "changes": map[string]any{"decisions": []any{"$new"}}},
		map[string]any{"op": "close", "id": a2b3, "status": "done", "outcome": "Shipped."},
		map[string]any{"op": "set", "id": "261001-zzzz", "changes": map[string]any{"priority": 1}},
	}})
	if err == nil || err.Code != "not_found" {
		t.Fatalf("want not_found, got %v", err)
	}
	if w, _ := err.Details.Get("write"); w != 3 {
		t.Errorf("details.write %v", w)
	}
	if rb, _ := err.Details.Get("rolled_back"); rb != true {
		t.Errorf("details.rolled_back %v (%v)", rb, err.Details)
	}
	after := texts(t, root)
	if !sameTexts(before, after) {
		for k, v := range after {
			if before[k] != v {
				t.Errorf("%s differs after the rollback", k)
			}
		}
		for k := range before {
			if _, ok := after[k]; !ok {
				t.Errorf("%s is missing after the rollback", k)
			}
		}
		t.Fatal("the repository is not as it was")
	}
	// the same sequence without the bad write goes through, so the first three did write
	resp := mustOK(t, root, "apply", Request{"as_of": "2026-11-16T10:00", "writes": []any{
		map[string]any{"op": "new", "kind": "adr", "slug": "use-json", "title": "Use JSON"},
		map[string]any{"op": "set", "id": a2b3, "changes": map[string]any{"decisions": []any{"$new"}}},
		map[string]any{"op": "close", "id": a2b3, "status": "done", "outcome": "Shipped."},
	}})
	if _, ok := previewEntry(t, resp, 2).Get("links_rewritten"); !ok {
		t.Errorf("the close rewrote no links, so the rollback above proved less than it should: %v", previewEntry(t, resp, 2))
	}
	if sameTexts(before, texts(t, root)) {
		t.Fatal("the sequence wrote nothing")
	}
}

// A stale first write, a dry_run and an unknown op are refused with the position
// of the write, and nothing is written.
func TestApplyRefusals(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "claims", "")
	before := texts(t, root)
	_, err := call(t, root, "apply", Request{"writes": []any{
		map[string]any{"op": "section", "id": e5f6, "name": "Goal", "text": "x"},
		map[string]any{"op": "section", "id": e5f6, "name": "Next", "text": "y", "if_match": "0000000000000000"},
	}})
	// the etag of a later write to the same record is chained, so only a first write can be stale
	if err != nil {
		t.Fatalf("a chained etag was checked as given: %v", err)
	}
	root = copyRepo(t, "claims", "")
	_, err = call(t, root, "apply", Request{"writes": []any{
		map[string]any{"op": "section", "id": e5f6, "name": "Goal", "text": "x", "if_match": "0000000000000000"},
	}})
	if err == nil || err.Code != "stale" {
		t.Fatalf("want stale, got %v", err)
	}
	_, err = call(t, root, "apply", Request{"writes": []any{
		map[string]any{"op": "section", "id": e5f6, "name": "Goal", "text": "x", "dry_run": true},
	}})
	if err == nil || err.Code != "invalid_request" {
		t.Fatalf("want invalid_request for dry_run, got %v", err)
	}
	_, err = call(t, root, "apply", Request{"writes": []any{map[string]any{"op": "index"}}})
	if err == nil || err.Code != "invalid_request" {
		t.Fatalf("want invalid_request for the op, got %v", err)
	}
	if !sameTexts(before, texts(t, root)) {
		t.Fatal("a refused sequence wrote")
	}
}
