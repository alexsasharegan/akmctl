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

// TestMoveClosedReportsInboundLinks is work item 261108-3t7t: moving closed
// records leaves links to them broken, which the step reports as move_file does.
// A link between two records that move together still resolves, so it is not.
func TestMoveClosedReportsInboundLinks(t *testing.T) {
	files := map[string]string{
		"docs/.layout":            "layout: 0.1.0\n",
		"docs/notes/one.md":       "---\nstatus: done\n---\n# One\n\nsee [two](two.md) and [live](live.md)\n",
		"docs/notes/two.md":       "---\nstatus: dropped\n---\n# Two\n",
		"docs/notes/live.md":      "---\nstatus: active\n---\n# Live\n\nsee [one](one.md) and [two](./two.md#top)\n",
		"docs/a.md":               "# A\n\n[one](notes/one.md) and [live](notes/live.md)\n\nPath [[docs/notes/two]], bare [[two]].\n",
		"docs/notes/sub/inner.md": "# Inner\n\n[one](../one.md)\n",
	}
	repo := filesRepo(t, files)
	s := Step{ID: "close", Kind: "move_closed", Class: "auto", Dir: "docs/notes"}
	want := []string{
		"inbound link left broken: docs/a.md:3 points at docs/notes/one.md, moved to docs/notes/done/one.md",
		"inbound wikilink left broken: docs/a.md:5 points at docs/notes/two.md, moved to docs/notes/done/two.md",
		"inbound link left broken: docs/notes/live.md:6 points at docs/notes/one.md, moved to docs/notes/done/one.md",
		"inbound link left broken: docs/notes/live.md:6 points at docs/notes/two.md, moved to docs/notes/done/two.md",
		"inbound link left broken: docs/notes/sub/inner.md:3 points at docs/notes/one.md, moved to docs/notes/done/one.md",
	}
	check := func(label string, notes []string) {
		t.Helper()
		if strings.Join(notes, "\n") != strings.Join(want, "\n") {
			t.Fatalf("%s notes:\n%s\nwant:\n%s", label, strings.Join(notes, "\n"), strings.Join(want, "\n"))
		}
	}
	check("plan", PlanEnv(repo, []Step{s}, Env{})[0].Notes)
	if !exists(repo, "docs/notes/one.md") {
		t.Fatal("plan moved a file")
	}
	check("apply", apply(t, repo, s).Notes)
	if exists(repo, "docs/notes/one.md") || !exists(repo, "docs/notes/done/one.md") {
		t.Fatal("the records were not moved")
	}
}

// TestMoveClosedBoundsInboundNotes: past the bound the notes end in one line
// that counts the rest.
func TestMoveClosedBoundsInboundNotes(t *testing.T) {
	files := map[string]string{"docs/.layout": "layout: 0.1.0\n", "docs/notes/one.md": "---\nstatus: done\n---\n# One\n"}
	var links strings.Builder
	for i := 0; i < maxInboundNotes+7; i++ {
		links.WriteString("[one](notes/one.md)\n\n")
	}
	files["docs/many.md"] = "# Many\n\n" + links.String()
	repo := filesRepo(t, files)
	notes := PlanEnv(repo, []Step{{ID: "close", Kind: "move_closed", Class: "auto", Dir: "docs/notes"}}, Env{})[0].Notes
	if len(notes) != maxInboundNotes+1 || notes[maxInboundNotes] != "... and 7 more inbound links left broken" {
		t.Fatalf("%d notes, last %q", len(notes), notes[len(notes)-1])
	}
}
