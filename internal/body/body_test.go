package body

import (
	"errors"
	"strings"
	"testing"
)

func TestLocateBlocks(t *testing.T) {
	src := "intro\n\n## Goal\r\n\n- [ ] one\n- [X] two\n* plain\n  - [x] nested\n\n```go\n## not\n```\n\n### Sub ###\n"
	lay, err := Locate([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	type want struct {
		kind  Kind
		line  int
		text  string
		level int
		done  bool
	}
	wants := []want{
		{Heading, 3, "## Goal", 2, false},
		{Task, 5, "- [ ] one", 0, false},
		{Task, 6, "- [X] two", 0, true},
		{ListItem, 7, "* plain", 0, false},
		{Task, 8, "  - [x] nested", 0, true},
		{Fence, 10, "```go\n## not\n```", 0, false},
		{Heading, 14, "### Sub ###", 3, false},
	}
	if len(lay.Blocks) != len(wants) {
		t.Fatalf("%d blocks, want %d: %+v", len(lay.Blocks), len(wants), lay.Blocks)
	}
	for i, w := range wants {
		b := lay.Blocks[i]
		got := src[b.Start:b.End]
		if b.Kind != w.kind || b.Line != w.line || got != w.text || b.Level != w.level || b.Done != w.done {
			t.Errorf("block %d = %v line %d %q level %d done %t; want %v line %d %q level %d done %t",
				i, b.Kind, b.Line, got, b.Level, b.Done, w.kind, w.line, w.text, w.level, w.done)
		}
	}
	if f := lay.Blocks[5]; f.EndLine != 12 {
		t.Errorf("fence ends at line %d, want 12", f.EndLine)
	}
}

func TestFenceTracking(t *testing.T) {
	for name, src := range map[string]string{
		"tilde does not close backticks": "```\n~~~\n## in\n```\n## out\n",
		"four holds three":               "````\n```\n## in\n```\n````\n## out\n",
		"info string on open":            "```go\n## in\n```\n## out\n",
		"longer close":                   "```\n## in\n`````\n## out\n",
	} {
		lay, err := Locate([]byte(src))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		var heads []string
		for _, b := range lay.Blocks {
			if b.Kind == Heading {
				heads = append(heads, src[b.Start:b.End])
			}
		}
		if len(heads) != 1 || heads[0] != "## out" {
			t.Errorf("%s: headings %q, want only %q", name, heads, "## out")
		}
	}
	// an unclosed fence runs to the end of the input
	lay, err := Locate([]byte("```\ntext\n## hidden\n"))
	if err != nil || len(lay.Blocks) != 1 || lay.Blocks[0].Kind != Fence || lay.Blocks[0].EndLine != 3 {
		t.Errorf("unclosed fence: %+v, %v", lay.Blocks, err)
	}
}

func TestOffsetsAcrossLineEndings(t *testing.T) {
	for _, nl := range []string{"\n", "\r\n"} {
		src := "## A" + nl + nl + "- [ ] x" + nl + "- [x] y"
		lay, err := Locate([]byte(src))
		if err != nil || len(lay.Blocks) != 3 {
			t.Fatalf("%q: %+v, %v", nl, lay.Blocks, err)
		}
		for _, b := range lay.Blocks {
			if strings.ContainsAny(src[b.Start:b.End], "\r\n") {
				t.Errorf("%q: block %q includes a line ending", nl, src[b.Start:b.End])
			}
		}
		if got := src[lay.Blocks[2].Start:lay.Blocks[2].End]; got != "- [x] y" {
			t.Errorf("%q: last block %q", nl, got)
		}
	}
}

func TestRefusals(t *testing.T) {
	for _, c := range []struct {
		construct, src string
		line           int
	}{
		{"indented fence", "## A\n\n  ```\n  x\n  ```\n", 3},
		{"HTML line", "## A\n\n<div>\n## not\n</div>\n", 3},
		{"blockquote", "> ## quoted\n", 1},
		{"tab-indented line", "text\n\tcode\n", 2},
		{"tab-indented line", "text\n \t* [X]\n", 2},
		{"table delimiter row", "a | b\n---|---\n", 2},
		{"table delimiter row", "| a | b | c |\n|---|---|\n", 2},
		{"table delimiter row", "* |\n|-\n", 2},
		{"indented setext underline, thematic break or empty item", "- a\n    -\n", 2},
		{"empty list item", "- a\n\n-\n", 3},
		{"bare carriage return line ending", "## A\r\r- x\n", 1},
		{"block inside a list item", "- ## h\n", 1},
		{"block inside a list item", "- ```\n", 1},
		{"block inside a list item", "+ ** *\n", 1},
		{"block inside a list item", "text\n  * # [X]\n", 2},
		{"setext heading underline or thematic break after text", "Title\n-\n", 2},
		{"setext heading underline or thematic break after text", "- [X] 0\n-:\n", 2},
		{"ordered-list task", "1. [ ] one\n", 1},
		{"ordered-list task", "0) 0) [X]\n", 1},
		{"task item indented four or more spaces", "- a\n    - [ ] deep\n", 2},
		{"indented code in a list item", "-      x\n", 1},
		{"control character", "## A\n\n* \f\n", 3},
		{"tab after a list marker", "*\t0\n", 1},
		{"tab after a list marker", "* \t0\n", 1},
		{"task box with unusual contents", "* [\t]\n", 1},
		{"task box with unusual contents", "- []  x\n", 1},
		{"task box not followed by a space", "* [X]0\n", 1},
		{"task box not followed by a space", "* [X]\t0\n", 1},
		{"task item without text", "- [ ]\n", 1},
		{"indented heading", "text\n\n  ## indented\n", 3},
		{"empty heading", "##\n", 1},
		{"setext heading underline", "Title\n===\n", 2},
		{"setext heading underline or thematic break after text", "Title\n---\n", 2},
		{"spaced thematic break", "* * *\n", 1},
	} {
		_, err := Locate([]byte(c.src))
		var u *UnsupportedError
		if !errors.Is(err, ErrUnsupported) || !errors.As(err, &u) {
			t.Errorf("%s: err = %v, want ErrUnsupported", c.construct, err)
			continue
		}
		if u.Construct != c.construct || u.Line != c.line {
			t.Errorf("%s: refused %q at line %d, want line %d", c.construct, u.Construct, u.Line, c.line)
		}
	}
}

func TestPlainThematicBreakAfterBlank(t *testing.T) {
	lay, err := Locate([]byte("## A\n\n---\n\n## B\n"))
	if err != nil || len(lay.Blocks) != 2 {
		t.Fatalf("%+v, %v", lay.Blocks, err)
	}
}

func TestTable(t *testing.T) {
	src := "## A\n\n| a | b |\n|---|:-:|\n| - [ ] x | ## y |\n| 1 | 2 |\n\n## B\n"
	lay, err := Locate([]byte(src))
	if err != nil || len(lay.Blocks) != 2 {
		t.Fatalf("%+v, %v", lay.Blocks, err)
	}
}

func TestEmpty(t *testing.T) {
	for _, s := range []string{"", "\n", "\n\n", "text"} {
		lay, err := Locate([]byte(s))
		if err != nil || len(lay.Blocks) != 0 {
			t.Errorf("%q: %+v, %v", s, lay.Blocks, err)
		}
	}
}
