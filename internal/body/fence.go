package body

import "strings"

// fenceOpen reports whether line opens a code fence at column 0, and with which
// character and run length. A backtick fence may not have a backtick in its info
// string. This is the only place the rule lives; the scanner and FenceTracker use it.
func fenceOpen(line string) (ch byte, n int, ok bool) {
	if n := runLen(line, '`'); n >= 3 && !strings.Contains(line[n:], "`") {
		return '`', n, true
	}
	if n := runLen(line, '~'); n >= 3 {
		return '~', n, true
	}
	return 0, 0, false
}

// fenceClose reports whether line closes a fence opened with ch and n: the same
// character, at least as long, up to three spaces of indent and nothing after it.
func fenceClose(line string, ch byte, n int) bool {
	c := strings.TrimLeft(line, " ")
	m := runLen(c, ch)
	return len(line)-len(c) <= 3 && m >= n && strings.TrimSpace(c[m:]) == ""
}

// FenceTracker follows code fences one line at a time, with the scanner's rule: a fence
// is closed by the same character, at least as long as the opener. A fence still
// open at the end of the text runs to the end. The zero value is outside a fence.
type FenceTracker struct {
	ch byte
	n  int
}

// In reports whether line, the next line of the text without its line ending,
// belongs to a fence: its opening line, a line inside, or its closing line.
func (f *FenceTracker) In(line string) bool {
	if f.ch != 0 {
		if fenceClose(line, f.ch, f.n) {
			f.ch, f.n = 0, 0
		}
		return true
	}
	if ch, n, ok := fenceOpen(line); ok {
		f.ch, f.n = ch, n
		return true
	}
	return false
}
