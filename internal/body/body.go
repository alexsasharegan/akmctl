// Package body finds the blocks of a record's Markdown body (ADRs 261031-gam0,
// 261031-cpbg and 261031-nvka): headings, fences, task items and list items, with
// their byte offsets, so that an edit can change the bytes it means to and no
// others. It reads a narrow subset of Markdown. A construct outside the subset is
// an error that names it and the line (the tool fails closed); it is never
// guessed at. A goldmark adapter in the tests is the oracle that the subset is
// read the way CommonMark reads it. On top of the blocks, tree.go splits a body into
// sections and edits them by splice (ADR 261031-gam0).
package body

import (
	"errors"
	"fmt"
)

// Kind is the kind of a block the scanner claims.
type Kind int

const (
	// Heading is an ATX heading at column 0 (one to six hashes, then a space).
	Heading Kind = iota + 1
	// Fence is a fenced code block, from its opening line to its closing line, or
	// to the last line when it is not closed.
	Fence
	// Task is a task list item, "- [ ] text" or "- [x] text", with at most three
	// leading spaces. It spans its first line.
	Task
	// ListItem is a bullet item at column 0 that is not a task. It spans its first
	// line.
	ListItem
	// Quote is a run of consecutive blockquote lines ("> text", with at most three
	// leading spaces). The scanner does not read inside it: it is opaque text.
	Quote
)

func (k Kind) String() string {
	switch k {
	case Heading:
		return "heading"
	case Fence:
		return "fence"
	case Task:
		return "task"
	case ListItem:
		return "list item"
	case Quote:
		return "quote"
	}
	return fmt.Sprintf("kind(%d)", int(k))
}

// Block is one block the scanner claims. Start is the offset of the first byte of
// its first line; End is the offset just after the last byte of its last line,
// before that line's ending ("\n" or "\r\n"), so body[Start:End] is the block's
// text without a trailing newline. Line and EndLine are 1-based.
type Block struct {
	Kind    Kind
	Level   int  // Heading: 1 to 6
	Done    bool // Task: the box is ticked
	Start   int
	End     int
	Line    int
	EndLine int
}

// Layout lists the blocks of a body in source order. A line that is inside a
// fence belongs to the Fence block and to no other. Refusals lists the local
// constructs a tolerant scan stepped over (see LocateTolerant); a strict Locate
// returns an error at the first one and leaves it empty.
type Layout struct {
	Blocks   []Block
	Refusals []Refusal
}

// Refusal is a local construct the tolerant scan did not read: it is treated as
// plain text and recorded, so an edit can tell whether it lies in the section it
// means to change.
type Refusal struct {
	Construct string
	Line      int // 1-based
	Start     int // offset of the line
}

// Locator finds the blocks of a body. The scanner implements it; the tests also
// implement it with goldmark, and one battery runs against both.
type Locator interface {
	Locate(body []byte) (Layout, error)
}

// ErrUnsupported is the error a Locator returns for a construct outside the
// supported subset; test for it with errors.Is.
var ErrUnsupported = errors.New("unsupported Markdown construct")

// UnsupportedError says which construct was met and on which line (1-based).
// Structural is true when the construct could change where a heading, fence or
// checklist is (ADR 261101-tr8e), so no edit of the body is safe; a local
// construct only affects its own lines.
type UnsupportedError struct {
	Construct  string
	Line       int
	Structural bool
}

// localConstructs are the refusals confined to their own lines: they cannot hide
// or invent a heading or a fence. Every other refusal is structural, so a
// construct that is not listed here is treated as the more cautious kind.
var localConstructs = map[string]bool{
	"table delimiter row":                    true,
	"empty list item":                        true,
	"ordered-list task":                      true,
	"task item indented four or more spaces": true,
	"tab after a list marker":                true,
	"task box with unusual contents":         true,
	"task box not followed by a space":       true,
	"task item without text":                 true,
}

// structuralConstructs are the refusals that could change where a heading, fence
// or checklist is. TestEveryRefusalIsClassified checks that this and
// localConstructs together name every construct the scanner refuses.
var structuralConstructs = map[string]bool{
	"control character":                true,
	"bare carriage return line ending": true,
	"tab-indented line":                true,
	"indented fence":                   true,
	"indented heading":                 true,
	"HTML line":                        true,
	"spaced thematic break":            true,
	"setext heading underline or thematic break after text":   true,
	"setext heading underline":                                true,
	"indented setext underline, thematic break or empty item": true,
	"indented code in a list item":                            true,
	"block inside a list item":                                true,
	"empty heading":                                           true,
}

// IsStructural reports whether refusing construct c is structural.
func IsStructural(c string) bool { return !localConstructs[c] }

func (e *UnsupportedError) Error() string {
	return fmt.Sprintf("line %d: unsupported Markdown construct: %s", e.Line, e.Construct)
}

func (e *UnsupportedError) Is(target error) bool { return target == ErrUnsupported }
