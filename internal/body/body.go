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
// fence belongs to the Fence block and to no other.
type Layout struct {
	Blocks []Block
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
type UnsupportedError struct {
	Construct string
	Line      int
}

func (e *UnsupportedError) Error() string {
	return fmt.Sprintf("line %d: unsupported Markdown construct: %s", e.Line, e.Construct)
}

func (e *UnsupportedError) Is(target error) bool { return target == ErrUnsupported }
