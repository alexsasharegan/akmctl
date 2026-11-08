package api

import (
	"reflect"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/testutil"
)

// historyRepo is commitsRepo with git_commits on, which history needs.
func historyRepo(t *testing.T, fixture string) string {
	t.Helper()
	testutil.RequireTool(t, "git")
	root := copyRepo(t, fixture, "git_commits: true\n")
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "start")
	return root
}

type move struct{ item, step, from, to string }

func historyMoves(t *testing.T, root string, req Request) (resp Obj, got []move) {
	t.Helper()
	resp = mustOK(t, root, "history", req)
	ms, _ := resp.Get("moves")
	for _, m := range ms.([]Obj) {
		g := func(k string) string { v, _ := m.Get(k); return v.(string) }
		got = append(got, move{g("item"), g("step"), g("from"), g("to")})
	}
	return resp, got
}

func TestHistoryListsStepCommitsOldestFirstWithStatuses(t *testing.T) {
	root := historyRepo(t, "claims")
	commitWith(t, root, mustOK(t, root, "claim", Request{"id": "261022-e5f6", "as": "agent-1", "as_of": "2026-10-28T09:30"}), nil)
	commitWith(t, root, mustOK(t, root, "release", Request{"id": "261022-e5f6", "as_of": "2026-10-28"}), nil)
	commitWith(t, root, mustOK(t, root, "new", Request{"kind": "work", "slug": "frob", "as_of": "2026-10-28"}), nil)
	mustOK(t, root, "tick", Request{"id": "261022-e5f6", "item": 1, "as_of": "2026-10-28"})
	commitWith(t, root, mustOK(t, root, "close", Request{"id": "261022-e5f6", "status": "done", "outcome": "Done.", "as_of": "2026-10-28"}), func(m string) string { return m + "\nAssisted-by: test-agent" })
	resp, got := historyMoves(t, root, Request{})
	var steps []string
	for _, m := range got {
		steps = append(steps, m.step)
	}
	if want := []string{"claim", "release", "file", "close"}; !reflect.DeepEqual(steps, want) {
		t.Fatalf("steps %v, want %v", steps, want)
	}
	if got[0].from != "ready" || got[0].to != "active" {
		t.Errorf("claim moved %q to %q", got[0].from, got[0].to)
	}
	if got[2].from != "" || got[2].to != "raw" {
		t.Errorf("file moved %q to %q, want \"\" to raw", got[2].from, got[2].to)
	}
	if got[3].to != "done" || got[3].from != "active" {
		t.Errorf("close moved %q to %q, want active to done (a rename into done/)", got[3].from, got[3].to)
	}
	ms, _ := resp.Get("moves")
	tr, _ := ms.([]Obj)[3].Get("trailers")
	if !reflect.DeepEqual(tr, []string{"Assisted-by: test-agent"}) {
		t.Errorf("trailers %v", tr)
	}
	if a, _ := resp.Get("available"); a != true {
		t.Errorf("available %v", a)
	}
	if h, _ := resp.Get("head"); len(h.(string)) != 40 && len(h.(string)) != 64 {
		t.Errorf("head %q", h)
	}
}

func TestHistoryShowsThePathOfOneItem(t *testing.T) {
	root := historyRepo(t, "claims")
	id := "261022-e5f6"
	commitWith(t, root, mustOK(t, root, "claim", Request{"id": id, "as": "a", "as_of": "2026-10-28T09:30"}), nil)
	commitWith(t, root, mustOK(t, root, "release", Request{"id": id, "as_of": "2026-10-28"}), nil)
	commitWith(t, root, mustOK(t, root, "set", Request{"id": id, "changes": map[string]any{"status": "ready"}, "as_of": "2026-10-28"}), nil)
	commitWith(t, root, mustOK(t, root, "new", Request{"kind": "work", "slug": "other", "as_of": "2026-10-28"}), nil)
	_, got := historyMoves(t, root, Request{"item": id})
	if len(got) != 3 {
		t.Fatalf("moves of %s: %v", id, got)
	}
	for i := 1; i < len(got); i++ {
		if got[i].from != got[i-1].to {
			t.Errorf("step %d starts at %q but the one before ended at %q", i, got[i].from, got[i-1].to)
		}
	}
}

func TestHistoryReportsOnlyTheItemOfTheSubject(t *testing.T) {
	root := historyRepo(t, "claims")
	// the commit also changes another item, but it only claims e5f6
	mustOK(t, root, "set", Request{"id": "261024-g7h8", "changes": map[string]any{"priority": 4}, "as_of": "2026-10-28"})
	commitWith(t, root, mustOK(t, root, "claim", Request{"id": "261022-e5f6", "as": "a", "as_of": "2026-10-28T09:30"}), nil)
	_, got := historyMoves(t, root, Request{})
	if len(got) != 1 || got[0].item != "261022-e5f6" {
		t.Errorf("moves %v, want one for 261022-e5f6", got)
	}
}

func TestHistoryRangeLimitAndHead(t *testing.T) {
	root := historyRepo(t, "claims")
	commitWith(t, root, mustOK(t, root, "claim", Request{"id": "261022-e5f6", "as": "a", "as_of": "2026-10-28T09:30"}), nil)
	resp, _ := historyMoves(t, root, Request{})
	head, _ := resp.Get("head")
	// nothing new since head
	resp, got := historyMoves(t, root, Request{"range": head.(string) + "..HEAD"})
	if a, _ := resp.Get("available"); a != true || len(got) != 0 {
		t.Errorf("available %v moves %v, want an empty answer", a, got)
	}
	commitWith(t, root, mustOK(t, root, "release", Request{"id": "261022-e5f6", "as_of": "2026-10-28"}), nil)
	commitWith(t, root, mustOK(t, root, "new", Request{"kind": "work", "slug": "more", "as_of": "2026-10-28"}), nil)
	_, got = historyMoves(t, root, Request{"range": head.(string) + "..HEAD"})
	if len(got) != 2 || got[0].step != "release" || got[1].step != "file" {
		t.Errorf("since head: %v", got)
	}
	// a limit keeps the newest and says so
	resp, got = historyMoves(t, root, Request{"limit": 1})
	if tr, _ := resp.Get("truncated"); tr != true || len(got) != 1 || got[0].step != "file" {
		t.Errorf("limit 1: truncated %v moves %v", tr, got)
	}
}

func TestHistoryIgnoresCommitsThatAreNotSteps(t *testing.T) {
	root := historyRepo(t, "claims")
	mustOK(t, root, "set", Request{"id": "261022-e5f6", "changes": map[string]any{"priority": 4}, "as_of": "2026-10-28"})
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "an ordinary commit")
	resp, got := historyMoves(t, root, Request{})
	if a, _ := resp.Get("available"); a != true || len(got) != 0 {
		t.Errorf("available %v moves %v, want none", a, got)
	}
}

func TestHistoryWithoutGitOrWithTheSettingOff(t *testing.T) {
	cases := []struct {
		name, reason string
		repo         func(t *testing.T) string
	}{
		{"not a work tree", "not_git", func(t *testing.T) string { return copyRepo(t, "claims", "git_commits: true\n") }},
		{"git_commits off", "git_commits_off", func(t *testing.T) string { return commitsRepo(t, "claims") }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp, got := historyMoves(t, c.repo(t), Request{})
			if a, _ := resp.Get("available"); a != false {
				t.Errorf("available %v", a)
			}
			if r, _ := resp.Get("reason"); r != c.reason {
				t.Errorf("reason %v, want %s", r, c.reason)
			}
			if h, _ := resp.Get("head"); h != "" || len(got) != 0 {
				t.Errorf("head %v moves %v, want none", h, got)
			}
		})
	}
}

func TestHistoryRejectsABadRangeAndItem(t *testing.T) {
	root := historyRepo(t, "claims")
	wantCode(t, root, "history", Request{"range": "--output=x"}, "invalid_request", "")
	wantCode(t, root, "history", Request{"range": "nosuchrev..HEAD"}, "invalid_request", "")
	wantCode(t, root, "history", Request{"item": "nope"}, "invalid_request", "")
	wantCode(t, root, "history", Request{"limit": 501}, "invalid_request", "")
	if !strings.Contains(strings.Join(Operations(), " "), "history") {
		t.Error("history is not listed in Operations()")
	}
}
