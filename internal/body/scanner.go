package body

import (
	"regexp"
	"strings"
)

// Scanner is the Locator that reads the supported subset line by line.
type Scanner struct{}

// Locate finds the blocks of body with the Scanner.
func Locate(body []byte) (Layout, error) { return Scanner{}.Locate(body) }

var (
	headingRe     = regexp.MustCompile(`^(#{1,6})([ \t]|$)`)
	indHeadingRe  = regexp.MustCompile(`^ {1,3}#{1,6}([ \t]|$)`)
	taskRe        = regexp.MustCompile(`^( {0,3})[-*+] {1,4}\[([ xX])\] +[ \t]*\S`)
	glueTaskRe    = regexp.MustCompile(`^ {0,3}[-*+] {1,4}\[[ xX]\][^ ]`)
	codeInItemRe  = regexp.MustCompile(`^ {0,3}[-*+] {5,}\S`)
	tabMarkerRe   = regexp.MustCompile(`^ {0,3}[-*+] *\t`)
	emptyTaskRe   = regexp.MustCompile(`^ {0,3}[-*+] {1,4}\[[ xX]\]\s*$`)
	orderedTask   = regexp.MustCompile(`^\s*\d+[.)]([ \t]+(\d+[.)]|[-*+]|#{1,6}|>))*[ \t]+\[`)
	deepTaskRe    = regexp.MustCompile(`^ {4,}[-*+] \[[ xX]\]`)
	bulletRe      = regexp.MustCompile(`^[-*+] +\S`)
	htmlRe        = regexp.MustCompile(`^ {0,3}<[a-zA-Z/!?]`)
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
func (Scanner) Locate(body []byte) (Layout, error) {
	var lay Layout
	for i, c := range body {
		if (c < 0x20 && c != '\t' && c != '\n' && c != '\r') || c == 0x7f {
			return Layout{}, &UnsupportedError{Construct: "control character", Line: 1 + strings.Count(string(body[:i]), "\n")}
		}
		if c == '\r' && (i+1 == len(body) || body[i+1] != '\n') {
			return Layout{}, &UnsupportedError{Construct: "bare carriage return line ending", Line: 1 + strings.Count(string(body[:i]), "\n")}
		}
	}
	var (
		fch       byte // fence character while inside a fence, else 0
		flen      int
		fStart    Block
		prevBlank = true
		prev      string // the previous line, "" after a blank one
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
			return Layout{}, &UnsupportedError{Construct: c, Line: lineNo}
		}

		if fch != 0 {
			if fenceClose(t, fch, flen) {
				fStart.End, fStart.EndLine = lineEnd, lineNo
				lay.Blocks = append(lay.Blocks, fStart)
				fch, flen = 0, 0
			} else {
				fStart.End, fStart.EndLine = lineEnd, lineNo
			}
			prevBlank = false
			off = next
			continue
		}

		blank := strings.TrimSpace(t) == ""
		switch {
		case blank:
			prevBlank, prev = true, ""
			off = next
			continue
		case leadTabRe.MatchString(t):
			return fail("tab-indented line")
		}

		// a fence opens at column 0
		if ch, n, ok := fenceOpen(t); ok {
			fch, flen = ch, n
			fStart = Block{Kind: Fence, Start: off, End: lineEnd, Line: lineNo, EndLine: lineNo}
			prevBlank = false
			off = next
			continue
		}

		switch {
		case indFenceRe.MatchString(t):
			return fail("indented fence")
		case htmlRe.MatchString(t):
			return fail("HTML line")
		case quoteRe.MatchString(t):
			return fail("blockquote")
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
			return fail("empty list item")
		case setextEqRe.MatchString(t) && !prevBlank:
			return fail("setext heading underline")
		case indHeadingRe.MatchString(t):
			return fail("indented heading")
		case orderedTask.MatchString(t):
			return fail("ordered-list task")
		case deepTaskRe.MatchString(t):
			return fail("task item indented four or more spaces")
		case codeInItemRe.MatchString(t):
			return fail("indented code in a list item")
		case tabMarkerRe.MatchString(t):
			return fail("tab after a list marker")
		case blockAfterMarker(t):
			return fail("block inside a list item")
		case unusualBox(t):
			return fail("task box with unusual contents")
		case glueTaskRe.MatchString(t):
			return fail("task box not followed by a space")
		case emptyTaskRe.MatchString(t):
			return fail("task item without text")
		}

		switch {
		case headingRe.MatchString(t):
			m := headingRe.FindStringSubmatch(t)
			if strings.Trim(strings.TrimSpace(t[len(m[1]):]), "#") == "" {
				return fail("empty heading")
			}
			lay.Blocks = append(lay.Blocks, Block{Kind: Heading, Level: len(m[1]), Start: off, End: lineEnd, Line: lineNo, EndLine: lineNo})
		case taskRe.MatchString(t):
			m := taskRe.FindStringSubmatch(t)
			lay.Blocks = append(lay.Blocks, Block{Kind: Task, Done: m[2] != " ", Start: off, End: lineEnd, Line: lineNo, EndLine: lineNo})
		case bulletRe.MatchString(t):
			lay.Blocks = append(lay.Blocks, Block{Kind: ListItem, Start: off, End: lineEnd, Line: lineNo, EndLine: lineNo})
		}
		prevBlank, prev = false, t
		off = next
	}
	if fch != 0 {
		lay.Blocks = append(lay.Blocks, fStart)
	}
	return lay, nil
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
