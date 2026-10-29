package record

import "testing"

func TestSetSection(t *testing.T) {
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
		if got := SetSection(in, c.section, c.text, c.app, c.eol); got != c.want {
			t.Errorf("%s:\ngot  %q\nwant %q", c.name, got, c.want)
		}
	}
	// no Log: the section goes to the end; the last section can be replaced and
	// a missing final newline is added
	if got := SetSection("# T\n\n## Goal\nx", "Notes", "n", false, "\n"); got != "# T\n\n## Goal\nx\n\n## Notes\n\nn\n" {
		t.Errorf("end: %q", got)
	}
	if got := SetSection("# T\n\n## Goal\nx", "goal", "y", false, "\n"); got != "# T\n\n## Goal\ny\n" {
		t.Errorf("last: %q", got)
	}
	// a heading inside a code fence is not a heading
	fenced := "## A\n```\n## Notes\n```\n"
	if got := SetSection(fenced, "Notes", "n", false, "\n"); got != fenced+"\n## Notes\n\nn\n" {
		t.Errorf("fence: %q", got)
	}
	// sibling sections are untouched
	if got := SetSection("## A\na\n\n## B\nb\n\n## C\nc\n", "b", "B2", false, "\n"); got != "## A\na\n\n## B\nB2\n\n## C\nc\n" {
		t.Errorf("sibling: %q", got)
	}
}
