package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const e5f6ID = "261022-e5f6"
const e5f6Path = "docs/work/261022-e5f6-document-the-exit-codes.md"

func TestCloseWritesClosedForDoneAndDropped(t *testing.T) {
	for _, status := range []string{"done", "dropped"} {
		root := copyRepo(t, "basic", "extension.claims: 0.3.0\n")
		mustOK(t, root, "claim", Request{"id": a2b3, "as": "agent-1", "as_of": "2026-10-28T09:30"})
		p := filepath.Join(root, "docs/work/261020-a2b3-add-json-output-to-ls.md")
		b, _ := os.ReadFile(p)
		os.WriteFile(p, append(b, []byte("\n## Outcome\nShipped.\n")...), 0o644)
		resp := mustOK(t, root, "close", Request{"id": a2b3, "status": status, "as_of": "2026-10-29"})
		to, _ := resp.Get("moved_to")
		text := read(t, root, to.(string))
		if !strings.Contains(text, "updated: 2026-10-29\nclosed: 2026-10-29\n") {
			t.Fatalf("%s: closed is not on the line after updated:\n%s", status, text)
		}
		if strings.Count(text, "closed:") != 1 {
			t.Fatalf("%s: %s", status, text)
		}
	}
}

func TestClosedIsOnlyWrittenByClose(t *testing.T) {
	root := copyRepo(t, "basic", "")
	_, err := call(t, root, "set", Request{"id": e5f6ID, "changes": map[string]any{"closed": "2026-10-28"}})
	if err == nil || err.Code != "invalid_request" || !strings.Contains(err.Message, "closed") || !strings.Contains(err.Message, "managed by the tool") {
		t.Fatalf("%v", err)
	}
	// not even to remove it
	_, err = call(t, root, "set", Request{"id": e5f6ID, "changes": map[string]any{"closed": nil}})
	if err == nil || err.Code != "invalid_request" {
		t.Fatalf("%v", err)
	}
	// other operations leave it alone
	mustOK(t, root, "set", Request{"id": e5f6ID, "changes": map[string]any{"priority": 2}})
	mustOK(t, root, "section", Request{"id": e5f6ID, "name": "Goal", "text": "Again."})
	if strings.Contains(read(t, root, e5f6Path), "closed") {
		t.Fatal("closed written by set or section")
	}
}

func TestSetCycleAndKeep(t *testing.T) {
	root := copyRepo(t, "basic", "")
	resp := mustOK(t, root, "set", Request{"id": e5f6ID, "changes": map[string]any{"cycle": "2026-PI3-S4", "keep": true}, "as_of": "2026-10-28"})
	if ch, _ := resp.Get("changed"); strings.Join(ch.([]string), ",") != "cycle,keep" {
		t.Fatalf("%v", ch)
	}
	text := read(t, root, e5f6Path)
	if !strings.Contains(text, "cycle: 2026-PI3-S4\n") || !strings.Contains(text, "keep: true\n") {
		t.Fatal(text)
	}
	clean(t, root)
	shown := mustOK(t, root, "show", Request{"id": e5f6ID})
	rec, _ := shown.Get("record")
	f, _ := rec.(Obj).Get("fields")
	if v, _ := f.(Obj).Get("cycle"); v != "2026-PI3-S4" {
		t.Fatalf("cycle %v", v)
	}
	if v, _ := f.(Obj).Get("keep"); v != true {
		t.Fatalf("keep %v", v)
	}
	// keep: false is stored as false; null removes both
	mustOK(t, root, "set", Request{"id": e5f6ID, "changes": map[string]any{"keep": false}})
	if !strings.Contains(read(t, root, e5f6Path), "keep: false\n") {
		t.Fatal(read(t, root, e5f6Path))
	}
	mustOK(t, root, "set", Request{"id": e5f6ID, "changes": map[string]any{"cycle": nil, "keep": nil}})
	if text := read(t, root, e5f6Path); strings.Contains(text, "cycle") || strings.Contains(text, "keep") {
		t.Fatal(text)
	}
}

func TestSetCycleKeepRejects(t *testing.T) {
	root := copyRepo(t, "basic", "")
	for _, ch := range []map[string]any{
		{"cycle": "has space"}, {"cycle": ""}, {"cycle": "-lead"}, {"cycle": strings.Repeat("a", 65)}, {"cycle": 5},
		{"keep": "yes"}, {"keep": 1},
	} {
		if _, err := call(t, root, "set", Request{"id": e5f6ID, "changes": ch}); err == nil || err.Code != "invalid_request" {
			t.Fatalf("%v: %v", ch, err)
		}
	}
	mustOK(t, root, "set", Request{"id": e5f6ID, "changes": map[string]any{"cycle": strings.Repeat("a", 64)}})
}

func TestLsFiltersByCycle(t *testing.T) {
	root := copyRepo(t, "basic", "")
	mustOK(t, root, "set", Request{"id": e5f6ID, "changes": map[string]any{"cycle": "2026-Q4"}})
	resp := mustOK(t, root, "ls", Request{"cycle": "2026-Q4", "fields": []any{"cycle", "keep"}})
	if n, _ := resp.Get("total"); n != 1 {
		t.Fatalf("%v", n)
	}
	if _, err := call(t, root, "ls", Request{"cycle": "has space"}); err == nil || err.Code != "invalid_request" {
		t.Fatalf("%v", err)
	}
}
