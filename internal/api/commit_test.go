package api

import (
	"regexp"
	"strings"
	"testing"
)

var reCommitMessage = regexp.MustCompile(`^(\d{6}-[0-9a-hjkmnp-tv-z]{4})\(([a-z]+)\): (.+)\n\nAKM-Step: ([a-z]+)\nAKM-Item: (\d{6}-[0-9a-hjkmnp-tv-z]{4})$`)

// checkCommit asserts the commit member is well formed and consistent: the
// subject and trailers agree with step and item (ADR 261030-tbkh, Confirmation).
func checkCommit(t *testing.T, resp Obj, step string) Obj {
	t.Helper()
	v, ok := resp.Get("commit")
	if !ok {
		t.Fatalf("no commit member: %v", resp)
	}
	c := v.(Obj)
	gs, _ := c.Get("step")
	gi, _ := c.Get("item")
	gd, _ := c.Get("detail")
	gm, _ := c.Get("message")
	m := reCommitMessage.FindStringSubmatch(gm.(string))
	if m == nil {
		t.Fatalf("message does not match the format: %q", gm)
	}
	if gs != step || m[2] != step || m[4] != step || m[1] != gi || m[5] != gi || m[3] != gd {
		t.Errorf("step/item/detail disagree with the message: %v %q", c, gm)
	}
	if !commitSteps[step] {
		t.Errorf("step %q is not in the closed set", step)
	}
	if ps, _ := c.Get("paths"); len(ps.([]any)) == 0 {
		t.Errorf("no paths")
	}
	return c
}

func TestNewReturnsCommit(t *testing.T) {
	root := copyRepo(t, "basic", "")
	c := checkCommit(t, mustOK(t, root, "new", Request{"kind": "work", "slug": "frob", "title": "Frob the widget"}), "file")
	if d, _ := c.Get("detail"); d != "Frob the widget" {
		t.Errorf("detail = %v", d)
	}
	checkCommit(t, mustOK(t, root, "new", Request{"kind": "adr", "slug": "use-frob"}), "file")
}

func TestSetStatusChangeReturnsCommit(t *testing.T) {
	root := copyRepo(t, "basic", "")
	// 261026-m1n2 is raw
	checkCommit(t, mustOK(t, root, "set", Request{"id": "261026-m1n2", "changes": map[string]any{"status": "ready"}}), "ready")
	checkCommit(t, mustOK(t, root, "set", Request{"id": "261026-m1n2", "changes": map[string]any{"status": "active"}}), "start")
	r := mustOK(t, root, "set", Request{"id": "261026-m1n2", "changes": map[string]any{"status": "blocked", "blocked_by": []any{"legal review"}}})
	if d, _ := checkCommit(t, r, "block").Get("detail"); d != "blocked on legal review" {
		t.Errorf("detail = %v", d)
	}
	checkCommit(t, mustOK(t, root, "set", Request{"id": "261026-m1n2", "changes": map[string]any{"status": "ready"}}), "unblock")
}

func TestSetWithoutStatusChangeReturnsNoCommit(t *testing.T) {
	root := copyRepo(t, "basic", "")
	for _, ch := range []map[string]any{{"priority": 3}, {"status": "ready"}} {
		resp := mustOK(t, root, "set", Request{"id": "261022-e5f6", "changes": ch})
		if _, ok := resp.Get("commit"); ok {
			t.Errorf("commit for %v", ch)
		}
	}
}

func TestCloseAndDropReturnCommit(t *testing.T) {
	root := copyRepo(t, "closable", "")
	c := checkCommit(t, mustOK(t, root, "close", Request{"id": "261020-a2b3", "status": "done", "outcome": "Done.", "as_of": "2026-10-28"}), "close")
	ps, _ := c.Get("paths")
	if len(ps.([]any)) != 2 || !strings.Contains(ps.([]any)[1].(string), "/done/") {
		t.Errorf("paths = %v", ps)
	}
	checkCommit(t, mustOK(t, root, "close", Request{"id": "261022-e5f6", "status": "dropped", "outcome": "No.", "as_of": "2026-10-28"}), "drop")
}
