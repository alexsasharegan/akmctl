package migrate

import (
	"strings"
	"testing"
)

// TestMoveFileReportsInboundLinks is work item 261029-m0vm: a move_file step
// leaves links to the old path broken and says where, in the plan and on apply.
func TestMoveFileReportsInboundLinks(t *testing.T) {
	files := map[string]string{
		"docs/.layout":     "layout: 0.1.0\n",
		"docs/old/note.md": "---\nstatus: active\n---\n# Note\n",
		"docs/a.md":        "# A\n\nline\nsee [the note](old/note.md#top) here\n",
		"docs/sub/b.md":    "# B\n\n[up](../old/note.md) and [other](../a.md) and `[code](../old/note.md)`\n",
		"docs/c.md":        "# C\n\n[elsewhere](old/other.md) and [web](https://example.com/old/note.md)\n",
		"docs/wiki.md":     "# W\n\nPath [[old/note]], full [[docs/old/note.md#h|alias]], bare [[note]], other [[old/other]].\n\n| a |\n| - |\n| [[docs/old/note\\|x]] |\n",
		"docs/fenced.md":   "# F\n\n```\n[x](old/note.md)\n```\n",
		"AGENTS.md":        "# Agents\n\n[note](docs/old/note.md)\n",
	}
	repo := filesRepo(t, files)
	s := Step{ID: "mv", Kind: "move_file", Class: "auto", From: "docs/old/note.md", To: "docs/new/note.md"}
	want := []string{
		"inbound link left broken: AGENTS.md:3",
		"inbound link left broken: docs/a.md:4",
		"inbound link left broken: docs/sub/b.md:3",
		"inbound wikilink left broken: docs/wiki.md:3",
		"inbound wikilink left broken: docs/wiki.md:3",
		"inbound wikilink left broken: docs/wiki.md:7",
	}
	check := func(label string, notes []string) {
		t.Helper()
		if len(notes) != len(want) {
			t.Fatalf("%s: %d notes: %q", label, len(notes), notes)
		}
		for i, w := range want {
			if !strings.HasPrefix(notes[i], w+" points at docs/old/note.md, moved to docs/new/note.md") {
				t.Fatalf("%s: note %d = %q, want %s", label, i, notes[i], w)
			}
		}
	}
	plan := PlanEnv(repo, []Step{s}, Env{})
	check("plan", plan[0].Notes)
	if !exists(repo, "docs/old/note.md") {
		t.Fatal("plan moved the file")
	}
	res, err := Apply(repo, []Step{s})
	if err != nil {
		t.Fatal(err)
	}
	check("apply", res[0].Notes)
	if !exists(repo, "docs/new/note.md") || exists(repo, "docs/old/note.md") {
		t.Fatal("file not moved")
	}
	if got := readRel(t, repo, "docs/a.md"); got != files["docs/a.md"] {
		t.Fatalf("an inbound link was rewritten: %q", got)
	}
	// A second run finds nothing to move and nothing to report.
	res, err = Apply(repo, []Step{s})
	if err != nil || len(res[0].Notes) != 0 || res[0].Status != "skipped" {
		t.Fatalf("second run: %v %+v", err, res[0])
	}
}
