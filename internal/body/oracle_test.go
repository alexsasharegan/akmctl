package body

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

// goldmarkLocator is the oracle: it reads a body with a CommonMark and GFM parser
// and reports the same kinds of block as the Scanner, so one battery can compare
// them. It lives in tests only (ADR 261031-cpbg); the binary does not link it.
//
// It reports what the Scanner claims and nothing else: ATX headings that are
// children of the document, task items, and bullet items that start at column 0
// in a top-level list. It reports no Fence blocks (goldmark gives a fenced block
// only its content lines); the battery checks fences by what they hide.
type goldmarkLocator struct{}

func (goldmarkLocator) Locate(src []byte) (Layout, error) {
	doc := goldmark.New(goldmark.WithExtensions(extension.GFM)).Parser().Parse(text.NewReader(src))
	var lay Layout
	add := func(b Block) { lay.Blocks = append(lay.Blocks, b) }
	lineAt := func(off int) (line, start, end int) {
		line = bytes.Count(src[:off], []byte("\n")) + 1
		start = off
		for start > 0 && src[start-1] != '\n' {
			start--
		}
		end = len(src)
		if i := bytes.IndexByte(src[start:], '\n'); i >= 0 {
			end = start + i
		}
		for end > start && src[end-1] == '\r' {
			end--
		}
		return
	}
	firstLine := func(n ast.Node) (off int, ok bool) {
		for c := n; c != nil; c = c.FirstChild() {
			if c.Type() == ast.TypeBlock && c.Lines().Len() > 0 {
				return c.Lines().At(0).Start, true
			}
			if c.FirstChild() == nil {
				break
			}
		}
		return 0, false
	}
	ast.Walk(doc, func(n ast.Node, enter bool) (ast.WalkStatus, error) {
		if !enter {
			return ast.WalkContinue, nil
		}
		switch v := n.(type) {
		case *ast.Heading:
			if v.Parent() == doc && v.Lines().Len() > 0 {
				ln, s, e := lineAt(v.Lines().At(0).Start)
				if src[s] == '#' { // ATX; a setext heading starts with its text
					add(Block{Kind: Heading, Level: v.Level, Start: s, End: e, Line: ln, EndLine: ln})
				}
			}
		case *east.TaskCheckBox:
			// the scanner treats a blockquote as opaque, so a box inside one is not a task it claims
			for a := v.Parent(); a != nil; a = a.Parent() {
				if _, quoted := a.(*ast.Blockquote); quoted {
					return ast.WalkContinue, nil
				}
			}
			if off, ok := firstLine(v.Parent()); ok {
				ln, s, e := lineAt(off)
				add(Block{Kind: Task, Done: v.IsChecked, Start: s, End: e, Line: ln, EndLine: ln})
			}
		case *ast.ListItem:
			list, _ := v.Parent().(*ast.List)
			if list == nil || list.Parent() != doc || list.IsOrdered() {
				break
			}
			if fc := v.FirstChild(); fc != nil && fc.FirstChild() != nil {
				if _, isTask := fc.FirstChild().(*east.TaskCheckBox); isTask {
					break
				}
			}
			if off, ok := firstLine(v); ok {
				ln, s, e := lineAt(off)
				if s < len(src) && strings.ContainsRune("-*+", rune(src[s])) { // the marker is at column 0
					add(Block{Kind: ListItem, Start: s, End: e, Line: ln, EndLine: ln})
				}
			}
		}
		return ast.WalkContinue, nil
	})
	sort.SliceStable(lay.Blocks, func(i, j int) bool { return lay.Blocks[i].Start < lay.Blocks[j].Start })
	return lay, nil
}

// project reduces a Layout to what both sides report: headings, tasks and list
// items, with their level or state and their first line, one per line of text.
func project(l Layout) []string {
	var out []string
	for _, b := range l.Blocks {
		switch b.Kind {
		case Heading:
			out = append(out, fmt.Sprintf("heading%d@%d[%d:%d]", b.Level, b.Line, b.Start, b.End))
		case Task:
			out = append(out, fmt.Sprintf("task(done=%t)@%d[%d:%d]", b.Done, b.Line, b.Start, b.End))
		case ListItem:
			out = append(out, fmt.Sprintf("item@%d[%d:%d]", b.Line, b.Start, b.End))
		}
	}
	return out
}

// compare runs the Scanner and the oracle on src. refused reports that the Scanner
// declined the input (always allowed); otherwise a non-empty diff means the Scanner
// was confidently wrong.
func compare(src []byte) (refused bool, diff string) {
	got, err := Locate(src)
	if err != nil {
		return true, ""
	}
	want, _ := goldmarkLocator{}.Locate(src)
	g, w := project(got), project(want)
	if strings.Join(g, " ") != strings.Join(w, " ") {
		return false, fmt.Sprintf("scanner %v\noracle  %v", g, w)
	}
	if d := fenceLeak(got); d != "" {
		return false, d
	}
	return false, ""
}

// fenceLeak checks the fence blocks by what they hide: no other block may start
// inside one.
func fenceLeak(l Layout) string {
	for _, f := range l.Blocks {
		if f.Kind != Fence {
			continue
		}
		for _, b := range l.Blocks {
			if b.Kind != Fence && b.Line >= f.Line && b.Line <= f.EndLine {
				return fmt.Sprintf("%s at line %d is inside the fence at lines %d to %d", b.Kind, b.Line, f.Line, f.EndLine)
			}
		}
	}
	return ""
}
