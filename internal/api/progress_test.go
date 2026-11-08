package api

import (
	"reflect"
	"testing"
)

func progressOfResp(t *testing.T, resp Obj) Obj {
	t.Helper()
	rec, _ := resp.Get("record")
	ro, ok := rec.(Obj).Get("progress")
	if !ok {
		return nil
	}
	return ro.(Obj)
}

func num(o Obj, k string) int { v, _ := o.Get(k); return v.(int) }

func TestShowProgressCountsChildren(t *testing.T) {
	root := copyRepo(t, "closable", "")
	// 261015-c4d5 has two open children; drop one and finish the other
	mustOK(t, root, "close", Request{"id": "261022-e5f6", "status": "dropped", "outcome": "No.", "as_of": "2026-10-28"})
	mustOK(t, root, "close", Request{"id": "261020-a2b3", "status": "done", "outcome": "Done.", "as_of": "2026-10-28"})
	ro := progressOfResp(t, mustOK(t, root, "show", Request{"id": "261015-c4d5"}))
	if ro == nil {
		t.Fatal("no progress on a parent")
	}
	got := map[string]int{}
	for _, kv := range ro {
		got[kv.K] = kv.V.(int)
	}
	want := map[string]int{"children_done": 1, "children_dropped": 1, "children_total": 2, "descendants_done": 1, "descendants_dropped": 1, "descendants_total": 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("progress = %v, want %v", got, want)
	}
	// a leaf has none
	if ro := progressOfResp(t, mustOK(t, root, "show", Request{"id": "261024-g7h8"})); ro != nil {
		t.Errorf("leaf progress = %v", ro)
	}
}

func TestProgressDescendantsCountGrandchildren(t *testing.T) {
	root := copyRepo(t, "basic", "")
	// make 261026-m1n2 a child of 261020-a2b3, which is a child of 261015-c4d5
	mustOK(t, root, "set", Request{"id": "261026-m1n2", "changes": map[string]any{"parent": "261020-a2b3"}})
	ro := progressOfResp(t, mustOK(t, root, "show", Request{"id": "261015-c4d5"}))
	if num(ro, "children_total") != 2 || num(ro, "descendants_total") != 3 {
		t.Errorf("progress = %v", ro)
	}
}

func TestLsTreeOrdersChildrenUnderParent(t *testing.T) {
	root := copyRepo(t, "basic", "")
	resp := mustOK(t, root, "ls", Request{"type": "work", "tree": true})
	v, _ := resp.Get("records")
	var ids []string
	depths := map[string]int{}
	for _, r := range v.([]Obj) {
		id, _ := r.Get("id")
		d, _ := r.Get("depth")
		ids = append(ids, id.(string))
		depths[id.(string)] = d.(int)
	}
	pos := func(id string) int {
		for i, x := range ids {
			if x == id {
				return i
			}
		}
		return -1
	}
	if !(pos("261015-c4d5") < pos("261020-a2b3") && pos("261020-a2b3") < pos("261022-e5f6") && pos("261022-e5f6") < pos("261024-g7h8")) {
		t.Errorf("order = %v", ids)
	}
	if depths["261015-c4d5"] != 0 || depths["261020-a2b3"] != 1 {
		t.Errorf("depths = %v", depths)
	}
	// a child whose parent the filter hides is a root
	resp = mustOK(t, root, "ls", Request{"type": "work", "tree": true, "status": []any{"ready"}})
	v, _ = resp.Get("records")
	for _, r := range v.([]Obj) {
		if d, _ := r.Get("depth"); d.(int) != 0 {
			id, _ := r.Get("id")
			t.Errorf("%v depth %v with its parent filtered out", id, d)
		}
	}
	// without tree there is no depth
	resp = mustOK(t, root, "ls", Request{"type": "work"})
	v, _ = resp.Get("records")
	if _, ok := v.([]Obj)[0].Get("depth"); ok {
		t.Error("depth without tree")
	}
}

func TestCloseNamesParentAndOpenSiblings(t *testing.T) {
	root := copyRepo(t, "closable", "")
	r := mustOK(t, root, "close", Request{"id": "261020-a2b3", "status": "done", "outcome": "Done.", "as_of": "2026-10-28"})
	if p, _ := r.Get("parent"); p != "261015-c4d5" {
		t.Errorf("parent = %v", p)
	}
	if s, _ := r.Get("open_siblings"); !reflect.DeepEqual(s, []string{"261022-e5f6"}) {
		t.Errorf("open_siblings = %v", s)
	}
	if c, _ := r.Get("parent_closable"); c != false {
		t.Errorf("parent_closable = %v", c)
	}
	// the last child: the parent is now closable, also in a dry run
	r = mustOK(t, root, "close", Request{"id": "261022-e5f6", "status": "dropped", "outcome": "No.", "as_of": "2026-10-28", "dry_run": true})
	if s, _ := r.Get("open_siblings"); !reflect.DeepEqual(s, []string{}) {
		t.Errorf("open_siblings = %#v", s)
	}
	if c, _ := r.Get("parent_closable"); c != true {
		t.Errorf("parent_closable = %v", c)
	}
	// an item with no parent reports none
	r = mustOK(t, root, "close", Request{"id": "261024-g7h8", "status": "dropped", "outcome": "No.", "as_of": "2026-10-28"})
	for _, k := range []string{"parent", "open_siblings", "parent_closable"} {
		if _, ok := r.Get(k); ok {
			t.Errorf("%s on an item with no parent", k)
		}
	}
}
