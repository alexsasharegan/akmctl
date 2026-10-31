package api

import (
	"os/exec"
	"strings"
	"testing"
)

func gitIn(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root, "-c", "user.name=t", "-c", "user.email=t@example.com", "-c", "commit.gpgsign=false"}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

// commitsRepo copies a fixture repository, makes it a git repository and commits
// the starting state.
func commitsRepo(t *testing.T, fixture string) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	root := copyRepo(t, fixture, "")
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "start")
	return root
}

// commitWith commits everything with the message of a commit object.
func commitWith(t *testing.T, root string, resp Obj, edit func(string) string) {
	t.Helper()
	v, ok := resp.Get("commit")
	if !ok {
		t.Fatalf("no commit member: %v", resp)
	}
	m, _ := v.(Obj).Get("message")
	msg := m.(string)
	if edit != nil {
		msg = edit(msg)
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", msg)
}

func commitsReport(t *testing.T, root string) (clean bool, problems []string, checked, skipped int) {
	t.Helper()
	resp := mustOK(t, root, "commits", Request{})
	c, _ := resp.Get("clean")
	ck, _ := resp.Get("checked")
	sk, _ := resp.Get("skipped")
	cs, _ := resp.Get("commits")
	for _, x := range cs.([]Obj) {
		ps, _ := x.Get("problems")
		for _, p := range ps.([]Obj) {
			r, _ := p.Get("rule")
			problems = append(problems, r.(string))
		}
	}
	return c.(bool), problems, ck.(int), sk.(int)
}

func TestCommitsAcceptsTheToolsOwnMessages(t *testing.T) {
	root := commitsRepo(t, "claims")
	commitWith(t, root, mustOK(t, root, "claim", Request{"id": "261022-e5f6", "as": "agent-1", "as_of": "2026-10-28T09:30"}), nil)
	commitWith(t, root, mustOK(t, root, "release", Request{"id": "261022-e5f6", "as_of": "2026-10-28"}), nil)
	commitWith(t, root, mustOK(t, root, "set", Request{"id": "261022-e5f6", "changes": map[string]any{"status": "ready"}, "as_of": "2026-10-28"}), nil)
	commitWith(t, root, mustOK(t, root, "new", Request{"kind": "work", "slug": "frob", "as_of": "2026-10-28"}), nil)
	commitWith(t, root, mustOK(t, root, "close", Request{"id": "261022-e5f6", "status": "done", "outcome": "Done.", "as_of": "2026-10-28"}), nil)
	clean, problems, checked, skipped := commitsReport(t, root)
	if !clean || len(problems) != 0 || checked != 5 || skipped != 1 {
		t.Errorf("clean=%v problems=%v checked=%d skipped=%d", clean, problems, checked, skipped)
	}
}

func TestCommitsAcceptsDrop(t *testing.T) {
	root := commitsRepo(t, "closable")
	commitWith(t, root, mustOK(t, root, "close", Request{"id": "261022-e5f6", "status": "dropped", "outcome": "No.", "as_of": "2026-10-28"}), nil)
	if clean, problems, checked, _ := commitsReport(t, root); !clean || checked != 1 {
		t.Errorf("clean=%v problems=%v checked=%d", clean, problems, checked)
	}
}

func TestCommitsFlagsWhatDoesNotAgree(t *testing.T) {
	cases := []struct {
		name string
		edit func(string) string
		want string
	}{
		{"trailer disagrees with the subject", func(m string) string { return strings.Replace(m, "AKM-Step: close", "AKM-Step: drop", 1) }, "CM2"},
		{"missing item trailer", func(m string) string { return strings.Replace(m, "\nAKM-Item: 261022-e5f6", "", 1) }, "CM2"},
		{"subject is not the format", func(m string) string { return strings.Replace(m, "(close): done", ": done", 1) }, "CM1"},
		{"step outside the closed set", func(m string) string { return strings.Replace(m, "(close)", "(finish)", 1) }, "CM1"},
		{"the diff is for another item", func(m string) string { return strings.ReplaceAll(m, "261022-e5f6", "261020-a2b3") }, "CM3"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := commitsRepo(t, "closable")
			commitWith(t, root, mustOK(t, root, "close", Request{"id": "261022-e5f6", "status": "done", "outcome": "Done.", "as_of": "2026-10-28"}), c.edit)
			clean, problems, _, _ := commitsReport(t, root)
			if clean || !has(problems, c.want) {
				t.Errorf("clean=%v problems=%v, want %s", clean, problems, c.want)
			}
		})
	}
}

func TestCommitsFlagsAStepThatDidNotHappen(t *testing.T) {
	root := commitsRepo(t, "closable")
	// a close subject on a commit that only edits the item
	resp := mustOK(t, root, "set", Request{"id": "261022-e5f6", "changes": map[string]any{"priority": 4}, "as_of": "2026-10-28"})
	_ = resp
	msg := "261022-e5f6(close): done\n\nAKM-Step: close\nAKM-Item: 261022-e5f6"
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", msg)
	clean, problems, _, _ := commitsReport(t, root)
	if clean || !has(problems, "CM3") {
		t.Errorf("clean=%v problems=%v", clean, problems)
	}
}

func TestCommitsSkipsOrdinaryCommits(t *testing.T) {
	root := commitsRepo(t, "basic")
	if clean, _, checked, skipped := commitsReport(t, root); !clean || checked != 0 || skipped != 1 {
		t.Errorf("clean=%v checked=%d skipped=%d", clean, checked, skipped)
	}
}

func TestCommitsRefusesWhatIsNotGitOrNotARange(t *testing.T) {
	root := copyRepo(t, "basic", "")
	wantCode(t, root, "commits", Request{}, "refused", "not_git")
	g := commitsRepo(t, "basic")
	for _, r := range []string{"--output=x", "-n5", "a b", "HEAD;ls"} {
		wantCode(t, g, "commits", Request{"range": r}, "invalid_request", "")
	}
	wantCode(t, g, "commits", Request{"range": "nope..HEAD"}, "invalid_request", "")
	wantCode(t, g, "commits", Request{"limit": 501}, "invalid_request", "")
}

func has(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}
