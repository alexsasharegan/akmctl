package record

import (
	"errors"
	mdbody "github.com/alexsasharegan/akmctl/internal/body"
	"strings"
	"testing"
)

func TestEditSection(t *testing.T) {
	body := "\n# T\n\n## Goal\nOld goal.\n\n## Log\n- one\n"
	cases := []struct {
		name, section, text string
		app                 bool
		eol, want           string
	}{
		{"replace keeps the separator", "goal", "New goal.", false, "\n", "\n# T\n\n## Goal\nNew goal.\n\n## Log\n- one\n"},
		{"replace is case-insensitive and trims newlines", "GOAL", "\nA\n\nB\n\n", false, "\n", "\n# T\n\n## Goal\nA\n\nB\n\n## Log\n- one\n"},
		{"empty text clears the section", "goal", "", false, "\n", "\n# T\n\n## Goal\n\n## Log\n- one\n"},
		{"append after the last line", "log", "- two", true, "\n", "\n# T\n\n## Goal\nOld goal.\n\n## Log\n- one\n- two\n"},
		{"append adds the bullet to a bare Log line", "log", "two", true, "\n", "\n# T\n\n## Goal\nOld goal.\n\n## Log\n- one\n- two\n"},
		{"append keeps a numbered item", "log", "2. two", true, "\n", "\n# T\n\n## Goal\nOld goal.\n\n## Log\n- one\n2. two\n"},
		{"append bullets only the first line", "Log", "two\n  more", true, "\n", "\n# T\n\n## Goal\nOld goal.\n\n## Log\n- one\n- two\n  more\n"},
		{"append to another section is verbatim", "goal", "more", true, "\n", "\n# T\n\n## Goal\nOld goal.\nmore\n\n## Log\n- one\n"},
		{"append nothing changes nothing", "log", "\n", true, "\n", body},
		{"create before Log", "Notes", "n", false, "\n", "\n# T\n\n## Goal\nOld goal.\n\n## Notes\n\nn\n\n## Log\n- one\n"},
		{"lowercase name is capitalised", "outcome", "done", false, "\n", "\n# T\n\n## Goal\nOld goal.\n\n## Outcome\n\ndone\n\n## Log\n- one\n"},
		{"crlf", "goal", "a\nb", false, "\r\n", "\r\n# T\r\n\r\n## Goal\r\na\r\nb\r\n\r\n## Log\r\n- one\r\n"},
	}
	for _, c := range cases {
		in := body
		if c.eol == "\r\n" {
			in = "\r\n# T\r\n\r\n## Goal\r\nOld goal.\r\n\r\n## Log\r\n- one\r\n"
		}
		if got := set(t, in, c.section, c.text, c.app, c.eol); got != c.want {
			t.Errorf("%s:\ngot  %q\nwant %q", c.name, got, c.want)
		}
	}
	// no Log: the section goes to the end; the last section can be replaced and
	// a missing final newline is added
	if got := set(t, "# T\n\n## Goal\nx", "Notes", "n", false, "\n"); got != "# T\n\n## Goal\nx\n\n## Notes\n\nn\n" {
		t.Errorf("end: %q", got)
	}
	if got := set(t, "# T\n\n## Goal\nx", "goal", "y", false, "\n"); got != "# T\n\n## Goal\ny\n" {
		t.Errorf("last: %q", got)
	}
	// a heading inside a code fence is not a heading
	fenced := "## A\n```\n## Notes\n```\n"
	if got := set(t, fenced, "Notes", "n", false, "\n"); got != fenced+"\n## Notes\n\nn\n" {
		t.Errorf("fence: %q", got)
	}
	// sibling sections are untouched
	if got := set(t, "## A\na\n\n## B\nb\n\n## C\nc\n", "b", "B2", false, "\n"); got != "## A\na\n\n## B\nB2\n\n## C\nc\n" {
		t.Errorf("sibling: %q", got)
	}
}

// TestFenceRule: every function that looks for a heading follows the same fence
// rule: a fence closes only on the same character and at least as many of it
// (work item 261031-2j14).
func TestFenceRule(t *testing.T) {
	mixed := "## A\n```\n~~~\n## Notes\n```\n## B\nb\n"  // ~~~ inside backticks does not close it
	long := "## A\n````\n```\n## Notes\n````\n## B\nb\n" // ``` inside ```` does not close it
	tilde := "## A\n~~~\n```\n## Notes\n~~~\n## B\nb\n"  // ``` inside ~~~ does not close it
	for name, body := range map[string]string{"mixed": mixed, "long": long, "tilde": tilde} {
		names, _ := Sections(body)
		if len(names) != 2 || names[0] != "a" || names[1] != "b" {
			t.Errorf("%s: Sections names = %v, want [a b]", name, names)
		}
		if got := set(t, body, "Notes", "n", false, "\n"); got != body+"\n## Notes\n\nn\n" {
			t.Errorf("%s: SetSection saw a heading inside the fence:\n%s", name, got)
		}
		if got := mustAppend(t, body, "Notes", "n"); got != strings.TrimRight(body, "\n")+"\n\n## Notes\nn\n" {
			t.Errorf("%s: AppendToSection saw a heading inside the fence:\n%s", name, got)
		}
		if got := FirstHeading("```\n~~~\n# in\n```\n# out\n"); got != "out" {
			t.Errorf("FirstHeading = %q, want out", got)
		}
	}
}

// set is EditSection for tests that expect no refusal.
func set(t *testing.T, body, name, text string, appendMode bool, eol string) string {
	t.Helper()
	got, err := EditSection(body, name, text, appendMode, eol)
	if err != nil {
		t.Fatalf("EditSection(%q, %q): %v", body, name, err)
	}
	return got
}

// TestEditSectionEdges: creation in an empty body, in CRLF, with no Log, and the
// refusal rule (ADR 261101-tr8e).
func TestEditSectionEdges(t *testing.T) {
	eq := func(name, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("%s:\ngot  %q\nwant %q", name, got, want)
		}
	}
	eq("empty body", set(t, "", "goal", "g", false, "\n"), "## Goal\n\ng\n")
	eq("empty body, CRLF", set(t, "", "goal", "g", false, "\r\n"), "## Goal\r\n\r\ng\r\n")
	eq("create before Log, CRLF", set(t, "## Goal\r\nx\r\n\r\n## Log\r\n- a\r\n", "Notes", "n\nm", false, "\r\n"), "## Goal\r\nx\r\n\r\n## Notes\r\n\r\nn\r\nm\r\n\r\n## Log\r\n- a\r\n")
	eq("create with no Log", set(t, "# T\n\ntext", "Notes", "n", false, "\n"), "# T\n\ntext\n\n## Notes\n\nn\n")
	eq("create, no text", set(t, "# T\n", "Notes", "", false, "\n"), "# T\n\n## Notes\n")
	eq("append creates", set(t, "## Goal\nx\n", "log", "- a", true, "\n"), "## Goal\nx\n\n## Log\n\n- a\n")
	// a heading is capitalised only when the name is all lowercase
	eq("mixed case kept", set(t, "## A\n", "ToDo", "t", false, "\n"), "## A\n\n## ToDo\n\nt\n")

	// refusals: structural anywhere, local only in the edited section
	structural := "## Goal\nx\n\n<div>\n\n## Notes\nn\n"
	for _, name := range []string{"goal", "notes", "other"} {
		_, err := EditSection(structural, name, "y", false, "\n")
		var u *mdbody.UnsupportedError
		if !errors.As(err, &u) || !u.Structural || u.Line != 4 {
			t.Errorf("structural, %s: err = %v", name, err)
		}
	}
	local := "## Goal\n- [ ] a\n\n## Notes\n*\tb\n"
	if _, err := EditSection(local, "goal", "y", false, "\n"); err != nil {
		t.Errorf("a local construct in another section refused the write: %v", err)
	}
	_, err := EditSection(local, "notes", "y", false, "\n")
	var u *mdbody.UnsupportedError
	if !errors.As(err, &u) || u.Structural || u.Line != 5 {
		t.Errorf("local construct in the edited section: err = %v", err)
	}
	if _, err := EditSection(local, "notes", "y", true, "\n"); err == nil {
		t.Errorf("append to a section with a local construct was not refused")
	}
}

func mustAppend(t *testing.T, body, name, line string) string {
	t.Helper()
	got, err := AppendToSection(body, name, line, "\n")
	if err != nil {
		t.Fatal(err)
	}
	return got
}
