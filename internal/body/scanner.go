package body

import (
	"regexp"
	"strings"
)

// Scanner is the Locator that reads the supported subset line by line. By
// default it returns an error at the first construct it does not read. With
// Tolerant set it steps over a local construct (see localConstructs), treating
// the line as plain text and recording it in Layout.Refusals, and still returns an
// error for a structural one.
type Scanner struct{ Tolerant bool }

// Locate finds the blocks of body with the strict Scanner.
func Locate(body []byte) (Layout, error) { return Scanner{}.Locate(body) }

// LocateTolerant finds the blocks of body, stepping over local constructs.
func LocateTolerant(body []byte) (Layout, error) { return Scanner{Tolerant: true}.Locate(body) }

var (
	headingRe     = regexp.MustCompile(`^(#{1,6})([ \t]|$)`)
	indHeadingRe  = regexp.MustCompile(`^ {1,3}#{1,6}([ \t]|$)`)
	taskRe        = regexp.MustCompile(`^( {0,3})[-*+] {1,4}\[([ xX])\] +[ \t]*\S`)
	glueTaskRe    = regexp.MustCompile(`^ {0,3}[-*+] {1,4}\[[ xX]\][^ ]`)
	codeInItemRe  = regexp.MustCompile(`^ {0,3}[-*+] {5,}\S`)
	tabMarkerRe   = regexp.MustCompile(`^ {0,3}[-*+] *\t`)
	emptyTaskRe   = regexp.MustCompile(`^ {0,3}[-*+] {1,4}\[[ xX]\]\s*$`)
	orderedTask   = regexp.MustCompile(`^\s*\d+[.)]([ \t]+(\d+[.)]|[-*+]|#{1,6}|>))*[ \t]+\[`)
	ordTaskRe     = regexp.MustCompile(`^( {0,3})(\d{1,9})([.)]) {1,4}\[([ xX])\] +[ \t]*\S`)
	ordItemRe     = regexp.MustCompile(`^( {0,3})(\d{1,9})([.)]) +\S`)
	deepTaskRe    = regexp.MustCompile(`^ {4,}[-*+] \[[ xX]\]`)
	bulletRe      = regexp.MustCompile(`^[-*+] +\S`)
	htmlRe        = regexp.MustCompile(`^ {0,3}<[a-zA-Z/!?]`)
	htmlCommentRe = regexp.MustCompile(`^ {0,3}<!--.*-->`)
	htmlInlineRe  = regexp.MustCompile(`^ {0,3}</?(?:a|b|i|u|em|strong|code|kbd|mark|span|sub|sup|small|del|ins)(?:[ \t][^<>]*)?>[ \t]*[^ \t<]`)
	quoteRe       = regexp.MustCompile(`^ {0,3}>`)
	tableRowRe    = regexp.MustCompile(`^[ 	|:-]*$`)
	spacedThemRe  = regexp.MustCompile(`^ {0,3}(-([ \t]*-){2,}|\*([ \t]*\*){2,}|_([ \t]*_){2,})[ \t]*$`)
	compactThemRe = regexp.MustCompile(`^ {0,3}(-{3,}|\*{3,}|_{3,})[ \t]*$`)
	emptyItemRe   = regexp.MustCompile(`^ {0,3}[-*+][ \t]*$`)
	dashRunRe     = regexp.MustCompile(`^ {0,3}:?-+:?[ \t]*$`)
	itemBlockRe   = regexp.MustCompile("^(#{1,6}([ \\t]|$)|>|<|[-*+]([ \\t]|$)|\\d+[.)]([ \\t]|$)|(-{3,}|\\*{3,}|_{3,}) *$)")
	indRuleRe     = regexp.MustCompile(`^ +(-+|=+|\*+|_+)[ \t]*$`)
	leadTabRe     = regexp.MustCompile(`^ *\t`)
	boxLikeRe     = regexp.MustCompile(`^ {0,3}[-*+] {1,4}\[([^\]]{0,2})\]`)
	setextEqRe    = regexp.MustCompile(`^ {0,3}=+[ \t]*$`)
	indFenceRe    = regexp.MustCompile("^ +(```|~~~)")
)

// Locate implements Locator.
func (sc Scanner) Locate(body []byte) (Layout, error) {
	var lay Layout
	for i, c := range body {
		if (c < 0x20 && c != '\t' && c != '\n' && c != '\r') || c == 0x7f {
			return Layout{}, &UnsupportedError{Construct: "control character", Line: 1 + strings.Count(string(body[:i]), "\n"), Structural: true}
		}
		if c == '\r' && (i+1 == len(body) || body[i+1] != '\n') {
			return Layout{}, &UnsupportedError{Construct: "bare carriage return line ending", Line: 1 + strings.Count(string(body[:i]), "\n"), Structural: true}
		}
	}
	var (
		fch       byte // fence character while inside a fence, else 0
		flen      int
		fStart    Block
		prevBlank = true
		prev      string // the previous line, "" after a blank one
		ord       ordState
		lineNo    int
	)
	off := 0
	for off <= len(body) {
		if off == len(body) && off > 0 && body[off-1] == '\n' {
			break
		}
		nl := indexByte(body[off:], '\n')
		end := len(body)
		next := len(body) + 1
		if nl >= 0 {
			end = off + nl
			next = end + 1
		}
		raw := string(body[off:end])
		t := strings.TrimSuffix(raw, "\r")
		lineEnd := off + len(t)
		lineNo++

		fail := func(c string) (Layout, error) {
			return Layout{}, &UnsupportedError{Construct: c, Line: lineNo, Structural: IsStructural(c)}
		}
		// refuse deals with a construct the scanner does not read: it stops with
		// an error, or, in a tolerant scan of a local construct, records the line
		// and goes on. It reports whether the scan must stop.
		var refusal error
		refuse := func(c string) bool {
			if sc.Tolerant && !IsStructural(c) {
				lay.Refusals = append(lay.Refusals, Refusal{Construct: c, Line: lineNo, Start: off})
				return false
			}
			_, refusal = fail(c)
			return true
		}

		if fch != 0 {
			if fenceClose(t, fch, flen) {
				fStart.End, fStart.EndLine = lineEnd, lineNo
				lay.Blocks = append(lay.Blocks, fStart)
				fch, flen = 0, 0
			} else {
				fStart.End, fStart.EndLine = lineEnd, lineNo
			}
			prevBlank, ord = false, ordState{}
			off = next
			continue
		}

		blank := strings.TrimSpace(t) == ""
		switch {
		case blank:
			prevBlank, prev, ord = true, "", ordState{}
			off = next
			continue
		case leadTabRe.MatchString(t):
			// Tried as plain text in work item 261101-4j86 and kept: the oracle reads a
			// space-and-tab indent before a task marker ("* a\n \t* [X]") as a task of
			// the list, which no line rule here matches (saved fuzz finding 672f74408e54df82).
			return fail("tab-indented line")
		}

		// a fence opens at column 0
		if ch, n, ok := fenceOpen(t); ok {
			fch, flen = ch, n
			fStart = Block{Kind: Fence, Start: off, End: lineEnd, Line: lineNo, EndLine: lineNo}
			prevBlank, ord = false, ordState{}
			off = next
			continue
		}

		var ordTask []string
		skip := false // a local construct was stepped over: the line is plain text
		switch {
		case indFenceRe.MatchString(t):
			return fail("indented fence")
		case htmlRe.MatchString(t) && !plainHTML(t):
			return fail("HTML line")
		case quoteRe.MatchString(t):
			// a quote is opaque: its lines are not read, and a run of them is one block
			if n := len(lay.Blocks); n > 0 && lay.Blocks[n-1].Kind == Quote && lay.Blocks[n-1].EndLine == lineNo-1 {
				lay.Blocks[n-1].End, lay.Blocks[n-1].EndLine = lineEnd, lineNo
			} else {
				lay.Blocks = append(lay.Blocks, Block{Kind: Quote, Start: off, End: lineEnd, Line: lineNo, EndLine: lineNo})
			}
			skip = true
		case isDelimiterLike(t):
			if !tableDelimiter(prev, t) {
				return fail("table delimiter row")
			}
		case spacedThemRe.MatchString(t) && !compactThemRe.MatchString(t):
			return fail("spaced thematic break")
		case (compactThemRe.MatchString(t) || dashRunRe.MatchString(t)) && !prevBlank:
			return fail("setext heading underline or thematic break after text")
		case indRuleRe.MatchString(t):
			return fail("indented setext underline, thematic break or empty item")
		case emptyItemRe.MatchString(t):
			if refuse("empty list item") {
				return Layout{}, refusal
			}
			skip = true
		case setextEqRe.MatchString(t) && !prevBlank:
			return fail("setext heading underline")
		case indHeadingRe.MatchString(t):
			return fail("indented heading")
		case orderedTask.MatchString(t):
			// "1. [ ] x" is a task like "- [ ] x" (work item 261101-wder). It is claimed
			// when it is a list item by CommonMark's rules (see ordState.listItem); any
			// other shape, such as nested markers or a box without text, stays a local
			// refusal.
			if m := ordTaskRe.FindStringSubmatch(t); m != nil && ord.listItem(m[1], m[2], m[3], prevBlank) {
				ordTask = m
			} else if refuse("ordered-list task") {
				return Layout{}, refusal
			}
			skip = ordTask == nil
		case deepTaskRe.MatchString(t):
			if refuse("task item indented four or more spaces") {
				return Layout{}, refusal
			}
			skip = true
		case codeInItemRe.MatchString(t):
			return fail("indented code in a list item")
		case tabMarkerRe.MatchString(t):
			if refuse("tab after a list marker") {
				return Layout{}, refusal
			}
			skip = true
		case blockAfterMarker(t):
			return fail("block inside a list item")
		case unusualBox(t):
			if refuse("task box with unusual contents") {
				return Layout{}, refusal
			}
			skip = true
		case glueTaskRe.MatchString(t):
			if refuse("task box not followed by a space") {
				return Layout{}, refusal
			}
			skip = true
		case emptyTaskRe.MatchString(t):
			if refuse("task item without text") {
				return Layout{}, refusal
			}
			skip = true
		}

		if !skip {
			switch {
			case headingRe.MatchString(t):
				m := headingRe.FindStringSubmatch(t)
				if strings.Trim(strings.TrimSpace(t[len(m[1]):]), "#") == "" {
					return fail("empty heading")
				}
				lay.Blocks = append(lay.Blocks, Block{Kind: Heading, Level: len(m[1]), Start: off, End: lineEnd, Line: lineNo, EndLine: lineNo})
			case ordTask != nil:
				lay.Blocks = append(lay.Blocks, Block{Kind: Task, Done: ordTask[4] != " ", Start: off, End: lineEnd, Line: lineNo, EndLine: lineNo})
			case taskRe.MatchString(t):
				m := taskRe.FindStringSubmatch(t)
				lay.Blocks = append(lay.Blocks, Block{Kind: Task, Done: m[2] != " ", Start: off, End: lineEnd, Line: lineNo, EndLine: lineNo})
			case bulletRe.MatchString(t):
				lay.Blocks = append(lay.Blocks, Block{Kind: ListItem, Start: off, End: lineEnd, Line: lineNo, EndLine: lineNo})
			}
		}
		ord = ord.next(t, prevBlank)
		prevBlank, prev = false, t
		off = next
	}
	if fch != 0 {
		lay.Blocks = append(lay.Blocks, fStart)
	}
	return lay, nil
}

// plainHTML reports an HTML line that cannot open an HTML block, so it cannot
// swallow a heading and is read as a line of text (work item 261101-v65x):
//   - a comment that ends on its own line (type 2 ends at the line holding "-->");
//   - a line that starts with an inline-only tag (not one of the block names that
//     open a type 6 block, nor script, style or pre) and has text after the tag.
//     A tag alone on its line is a type 7 block, which runs to the next blank
//     line, so it stays refused.
//
// The oracle battery checks both against goldmark.
func plainHTML(t string) bool {
	return htmlCommentRe.MatchString(t) || htmlInlineRe.MatchString(t)
}

// ordState remembers whether the line before was an ordered list item, so that the
// next one can be its sibling.
type ordState struct {
	ok     bool
	indent int
	delim  string
}

// listItem reports whether an ordered line is a list item and not text of the
// paragraph before it. An ordered item may interrupt a paragraph only when it
// starts at 1; otherwise it must follow a blank line, or be the next item of the
// list the line before began (same delimiter and indent).
func (o ordState) listItem(indent, num, delim string, prevBlank bool) bool {
	return prevBlank || num == "1" || (o.ok && o.delim == delim && o.indent == len(indent))
}

// next is the state after line t, which was read with prevBlank before it.
func (o ordState) next(t string, prevBlank bool) ordState {
	m := ordItemRe.FindStringSubmatch(t)
	if m == nil || !o.listItem(m[1], m[2], m[3], prevBlank) {
		return ordState{}
	}
	return ordState{ok: true, indent: len(m[1]), delim: m[3]}
}

func runLen(t string, c byte) int {
	n := 0
	for n < len(t) && t[n] == c {
		n++
	}
	return n
}

func indexByte(b []byte, c byte) int {
	for i, x := range b {
		if x == c {
			return i
		}
	}
	return -1
}

// tableDelimiter reports whether delim is the delimiter row of a pipe-led table
// whose header row is prev: both start with a pipe and have the same number of
// cells, and every delimiter cell is dashes with optional colons. The table's body
// rows need no care: they cannot start a claimed block, and a line that does
// (heading, fence, bullet) ends the table.
func tableDelimiter(prev, delim string) bool {
	h, d := strings.TrimSpace(prev), strings.TrimSpace(delim)
	if !strings.HasPrefix(h, "|") || !strings.HasPrefix(d, "|") {
		return false
	}
	hc, dc := cells(h), cells(d)
	if len(hc) != len(dc) {
		return false
	}
	for _, c := range dc {
		c = strings.Trim(strings.TrimSpace(c), ":")
		if c == "" || strings.Trim(c, "-") != "" {
			return false
		}
	}
	return true
}

// cells splits a pipe-led row into its cells, ignoring escaped pipes' effect only
// by refusing: a backslash before a pipe makes the count unsure, so it counts as
// a mismatch.
func cells(row string) []string {
	if strings.Contains(row, `\|`) {
		return nil
	}
	row = strings.TrimPrefix(row, "|")
	row = strings.TrimSuffix(row, "|")
	return strings.Split(row, "|")
}

// blockInItem reports whether the text after a bullet marker starts a block of
// its own (a heading, quote, HTML, fence, nested list or thematic break) instead
// of a paragraph. The oracle locates an item by its first paragraph, so the
// Scanner declines these rather than claim a line it cannot cross-check.
func blockInItem(rest string) bool {
	return itemBlockRe.MatchString(rest) || spacedThemRe.MatchString(rest) || compactThemRe.MatchString(rest) || runLen(rest, '`') >= 3 || runLen(rest, '~') >= 3
}

// unusualBox reports a bullet line that starts with a bracket pair that is almost
// a task box ("[]", "[\t]", "[xx]") but is not "[ ]", "[x]" or "[X]".
func unusualBox(t string) bool {
	m := boxLikeRe.FindStringSubmatch(t)
	return m != nil && m[1] != " " && m[1] != "x" && m[1] != "X"
}

var markerRe = regexp.MustCompile(`^ {0,3}[-*+] +(\S.*)$`)

// blockAfterMarker reports a bullet line, at any indent up to three spaces, whose
// text starts a block of its own.
func blockAfterMarker(t string) bool {
	m := markerRe.FindStringSubmatch(t)
	return m != nil && blockInItem(m[1])
}

// isDelimiterLike reports a line made only of pipes, dashes, colons and blanks
// with at least one pipe and one dash: a table delimiter row, or something that
// could be mistaken for one.
func isDelimiterLike(t string) bool {
	return tableRowRe.MatchString(t) && strings.Contains(t, "|") && strings.Contains(t, "-")
}
