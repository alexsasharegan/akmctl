package body

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// The section tree (ADR 261031-gam0). A body is a preamble, the text before the
// first "## " heading, and an ordered list of sections. The tree is a view over
// the bytes it was split from: Join returns them unchanged, and an edit builds
// the new bytes by splicing only the span it means to change, so unknown
// sections and hand formatting are never rewritten.

// Errors the edits return.
var (
	ErrNoSection  = errors.New("no such section")
	ErrNoItem     = errors.New("no such item")
	ErrStale      = errors.New("the section changed since it was read")
	ErrAppendOnly = errors.New("the Log section is append-only")
	ErrChecked    = errors.New("the item is ticked")
)

// Section is one "## " section: its heading line and the text under it, up to
// the next "## " heading.
type Section struct {
	Name  string // the heading text, trimmed
	Line  int    // 1-based line of the heading
	Start int    // offset of the heading line
	Body  int    // offset of the first byte after the heading line
	End   int    // offset of the next section's heading, or the end of the body
}

// Key is the name compared case-insensitively, as sections are addressed.
func (s Section) Key() string { return strings.ToLower(s.Name) }

// Tree is a body split into a preamble and sections.
type Tree struct {
	src      []byte
	sections []Section
	refusals []Refusal // local constructs the scan stepped over
	eol      string    // the line ending for a body with no line break (default "\n")
}

// WithEOL sets the line ending the tree writes when the body has no line break
// of its own (an empty body, or one line); a body that has one keeps it.
func (t *Tree) WithEOL(eol string) *Tree {
	t.eol = eol
	return t
}

// Split reads body (the text after the frontmatter) with a tolerant Locate and
// splits it at its "## " headings. A structural construct outside the supported
// subset is an error (ErrUnsupported, with the construct and the line), so a body
// whose sections the scanner cannot be sure of is never split by guesswork. A
// local construct is stepped over: the tree splits, and an edit is refused only
// when the local construct lies in the section it would change (ADR 261101-tr8e).
func Split(body []byte) (*Tree, error) {
	lay, err := LocateTolerant(body)
	if err != nil {
		return nil, err
	}
	return build(body, lay), nil
}

// build splits src at the level-2 headings of lay.
func build(src []byte, lay Layout) *Tree {
	t := &Tree{src: src, refusals: lay.Refusals}
	for _, b := range lay.Blocks {
		if b.Kind != Heading || b.Level != 2 {
			continue
		}
		body := b.End
		if body < len(src) && src[body] == '\r' {
			body++
		}
		if body < len(src) && src[body] == '\n' {
			body++
		}
		t.sections = append(t.sections, Section{Name: headingName(string(src[b.Start:b.End])), Line: b.Line, Start: b.Start, Body: body, End: len(src)})
	}
	for i := 0; i+1 < len(t.sections); i++ {
		t.sections[i].End = t.sections[i+1].Start
	}
	return t
}

// headingName is the text of an ATX heading line: no marker, no closing hashes
// (a run of hashes after a space), trimmed.
func headingName(line string) string {
	s := strings.TrimSpace(strings.TrimLeft(line, "#"))
	if i := strings.LastIndexAny(s, " \t"); i >= 0 && strings.Trim(s[i+1:], "#") == "" {
		s = strings.TrimSpace(s[:i])
	} else if strings.Trim(s, "#") == "" {
		s = ""
	}
	return s
}

// Refusals lists the local constructs the split stepped over, in order.
func (t *Tree) Refusals() []Refusal { return append([]Refusal(nil), t.refusals...) }

// refusedIn is the error for an edit of section s when a local construct lies
// inside it, or nil.
func (t *Tree) refusedIn(s Section) error {
	for _, r := range t.refusals {
		if r.Start >= s.Start && r.Start < s.End {
			return &UnsupportedError{Construct: r.Construct, Line: r.Line}
		}
	}
	return nil
}

// Join returns the bytes the tree was split from.
func (t *Tree) Join() []byte { return append([]byte(nil), t.src...) }

// Preamble is the text before the first section.
func (t *Tree) Preamble() []byte {
	end := len(t.src)
	if len(t.sections) > 0 {
		end = t.sections[0].Start
	}
	return t.src[:end]
}

// Sections lists the sections in order.
func (t *Tree) Sections() []Section { return append([]Section(nil), t.sections...) }

// Section finds the first section named name, without regard to case.
func (t *Tree) Section(name string) (Section, bool) {
	for _, s := range t.sections {
		if strings.EqualFold(s.Name, strings.TrimSpace(name)) {
			return s, true
		}
	}
	return Section{}, false
}

// Text is the text under the heading of s, without the heading line.
func (t *Tree) Text(s Section) []byte { return t.src[s.Body:s.End] }

// Etag identifies the bytes of s, its heading line included. An edit that
// carries an etag is refused with ErrStale when the section has changed.
func (t *Tree) Etag(s Section) string {
	h := sha256.Sum256(t.src[s.Start:s.End])
	return hex.EncodeToString(h[:8])
}

// lineSpan is one line of the body: its bytes without the line ending, and the
// ending ("\n", "\r\n" or "" at the end of a body with no final newline).
type lineSpan struct {
	start, end int // the text, without the ending
	eol        string
}

func (t *Tree) lines(from, to int) []lineSpan {
	var out []lineSpan
	for off := from; off < to; {
		nl := bytes.IndexByte(t.src[off:to], '\n')
		if nl < 0 {
			out = append(out, lineSpan{off, to, ""})
			break
		}
		end := off + nl
		eol := "\n"
		if end > off && t.src[end-1] == '\r' {
			end--
			eol = "\r\n"
		}
		out = append(out, lineSpan{off, end, eol})
		off += nl + 1
	}
	return out
}

// fileEOL is the line ending the body uses, from its first line break.
func (t *Tree) fileEOL() string {
	if i := bytes.IndexByte(t.src, '\n'); i > 0 && t.src[i-1] == '\r' {
		return "\r\n"
	} else if i < 0 && t.eol != "" {
		return t.eol
	}
	return "\n"
}

// LineEnding is the line ending the tree writes: the body's own, or the one set
// by WithEOL when the body has no line break.
func (t *Tree) LineEnding() string { return t.fileEOL() }

// ---- Goal: the checklist ----

// Item is one checklist item.
type Item struct {
	Pos   int  // 1-based position in the checklist
	Done  bool // the box is ticked
	Text  string
	Line  int // 1-based line in the body
	Start int // offset of the line
	End   int // offset after the line's text, before its line ending
	box   int // offset of the character between the brackets
}

// Checklist is the first contiguous run of task lines ("- [ ] text", "1. [x] text",
// bullet and ordered alike, one per line) in the section, or nil when it has none. A task line after a
// blank line or other text is not part of the run; it stays in the prose.
func (t *Tree) Checklist(name string) ([]Item, error) {
	s, ok := t.Section(name)
	if !ok {
		return nil, ErrNoSection
	}
	lay, err := LocateTolerant(t.src[s.Body:s.End])
	if err != nil {
		return nil, err
	}
	var out []Item
	prevLine := 0
	for _, b := range lay.Blocks {
		if b.Kind != Task {
			if len(out) > 0 {
				break // another block ends the run
			}
			continue
		}
		if len(out) > 0 && b.Line != prevLine+1 {
			break
		}
		start, end := s.Body+b.Start, s.Body+b.End
		line := string(t.src[start:end])
		open := strings.IndexByte(line, '[')
		out = append(out, Item{
			Pos: len(out) + 1, Done: b.Done, Text: strings.TrimSpace(line[open+3:]),
			Line: s.Line + b.Line, Start: start, End: end, box: start + open + 1,
		})
		prevLine = b.Line
	}
	return out, nil
}

// Unchecked lists every task of section name whose box is empty, in order, not
// only the first run that Checklist reads. Pos is the item's position in the
// checklist (the number tick takes), or 0 for a task outside that run, which tick
// cannot address.
func (t *Tree) Unchecked(name string) ([]Item, error) {
	s, ok := t.Section(name)
	if !ok {
		return nil, ErrNoSection
	}
	run, err := t.Checklist(name)
	if err != nil {
		return nil, err
	}
	pos := map[int]int{}
	for _, it := range run {
		pos[it.Start] = it.Pos
	}
	lay, err := LocateTolerant(t.src[s.Body:s.End])
	if err != nil {
		return nil, err
	}
	var out []Item
	for _, b := range lay.Blocks {
		if b.Kind != Task || b.Done {
			continue
		}
		start, end := s.Body+b.Start, s.Body+b.End
		line := string(t.src[start:end])
		open := strings.IndexByte(line, '[')
		out = append(out, Item{Pos: pos[start], Text: strings.TrimSpace(line[open+3:]), Line: s.Line + b.Line, Start: start, End: end})
	}
	return out, nil
}

// SetChecked ticks or unticks item pos (1-based) of the checklist in section
// name and returns the new body: only the character between the brackets
// changes. etag, when not empty, must be the Etag of the section as it was read.
func (t *Tree) SetChecked(name string, pos int, done bool, etag string) ([]byte, error) {
	s, ok := t.Section(name)
	if !ok {
		return nil, ErrNoSection
	}
	if err := t.refusedIn(s); err != nil {
		return nil, err
	}
	if etag != "" && etag != t.Etag(s) {
		return nil, ErrStale
	}
	items, err := t.Checklist(name)
	if err != nil {
		return nil, err
	}
	if pos < 1 || pos > len(items) {
		return nil, fmt.Errorf("%w: item %d of %d", ErrNoItem, pos, len(items))
	}
	it := items[pos-1]
	if it.Done == done {
		return t.Join(), nil
	}
	c := byte(' ')
	if done {
		c = 'x'
	}
	out := t.Join()
	out[it.box] = c
	return out, nil
}

// RemoveItem deletes item pos (1-based) of the checklist in section name, line
// ending included, and returns the new body with the removed item. Later items
// move up by one. A ticked item is refused with ErrChecked: it is the record
// that the work was done. etag, when not empty, must be the Etag of the section
// as it was read.
func (t *Tree) RemoveItem(name string, pos int, etag string) ([]byte, Item, error) {
	s, ok := t.Section(name)
	if !ok {
		return nil, Item{}, ErrNoSection
	}
	if err := t.refusedIn(s); err != nil {
		return nil, Item{}, err
	}
	if etag != "" && etag != t.Etag(s) {
		return nil, Item{}, ErrStale
	}
	items, err := t.Checklist(name)
	if err != nil {
		return nil, Item{}, err
	}
	if pos < 1 || pos > len(items) {
		return nil, Item{}, fmt.Errorf("%w: item %d of %d", ErrNoItem, pos, len(items))
	}
	it := items[pos-1]
	if it.Done {
		return nil, it, ErrChecked
	}
	from, to := it.Start, it.End
	switch {
	case to < len(t.src) && t.src[to] == '\r' && to+1 < len(t.src) && t.src[to+1] == '\n':
		to += 2
	case to < len(t.src) && t.src[to] == '\n':
		to++
	case from > 0 && t.src[from-1] == '\n': // the last line has no ending: take the one before it
		from--
		if from > 0 && t.src[from-1] == '\r' {
			from--
		}
	}
	src := t.Join()
	return append(src[:from:from], src[to:]...), it, nil
}

// ---- Log: dated entries ----

// Entry is one Log entry, "- YYYY-MM-DD: text".
type Entry struct {
	Date   string
	Time   string // "HH:MM", empty for a date-only line
	Offset string // "+HH:MM" or "-HH:MM", empty for a date-only line
	Text   string
	Line   int // 1-based line in the body
	Start  int
	End    int
}

var logLine = regexp.MustCompile(`^- (\d{4}-\d{2}-\d{2})(?: (\d{2}:\d{2})([+-]\d{2}:\d{2}))?: ?(.*)$`)

// LogEntries lists the lines of the Log section that have the entry form, in
// order. Any other line (a continuation, a note) is opaque text and is not
// listed.
func (t *Tree) LogEntries() ([]Entry, error) {
	s, ok := t.Section("log")
	if !ok {
		return nil, ErrNoSection
	}
	var out []Entry
	for i, l := range t.lines(s.Body, s.End) {
		if m := logLine.FindStringSubmatch(string(t.src[l.start:l.end])); m != nil {
			out = append(out, Entry{Date: m[1], Time: m[2], Offset: m[3], Text: m[4], Line: s.Line + 1 + i, Start: l.start, End: l.end})
		}
	}
	return out, nil
}

// AppendLog adds "- stamp: text" (the stamp is a date, or a date with a time and
// offset, ADR 261101-gw7b) after the last non-blank line of the Log section
// and returns the new body. The line ends as the line before it does.
func (t *Tree) AppendLog(stamp, text string) ([]byte, error) {
	s, ok := t.Section("log")
	if !ok {
		return nil, ErrNoSection
	}
	if err := t.refusedIn(s); err != nil {
		return nil, err
	}
	text = strings.Join(strings.Fields(text), " ")
	entry := "- " + stamp + ": " + text
	return t.appendLine(s, entry), nil
}

// appendLine inserts a line after the last non-blank line of s (after its
// heading when it has none).
func (t *Tree) appendLine(s Section, line string) []byte {
	eol := t.fileEOL()
	at := s.Body
	var last lineSpan
	found := false
	for _, l := range t.lines(s.Body, s.End) {
		if strings.TrimSpace(string(t.src[l.start:l.end])) != "" {
			last, found = l, true
		}
	}
	var ins string
	if found {
		if last.eol != "" {
			eol = last.eol
		}
		at = last.end + len(last.eol)
		if last.eol == "" { // no final newline: break the line first, and add none after
			ins = eol + line
		} else {
			ins = line + eol
		}
	} else {
		// an empty section: a blank line after the heading, then the line
		if s.Body > s.Start && t.src[s.Body-1] == '\n' {
			ins = eol + line + eol
		} else { // the heading is the last line and has no ending
			ins = eol + eol + line
		}
	}
	out := make([]byte, 0, len(t.src)+len(ins))
	out = append(out, t.src[:at]...)
	out = append(out, ins...)
	out = append(out, t.src[at:]...)
	return out
}

// ---- Replacing a section ----

// ReplaceSection replaces the text of section name and returns the new body. The
// blank lines before and after the old text stay; the new text is written with
// the body's line endings. A section with no text gets a blank line after its
// heading. The Log section cannot be replaced (ErrAppendOnly). etag, when not
// empty, must be the Etag of the section as it was read.
func (t *Tree) ReplaceSection(name, text, etag string) ([]byte, error) {
	s, ok := t.Section(name)
	if !ok {
		return nil, ErrNoSection
	}
	if strings.EqualFold(s.Name, "log") {
		return nil, ErrAppendOnly
	}
	if err := t.refusedIn(s); err != nil {
		return nil, err
	}
	if etag != "" && etag != t.Etag(s) {
		return nil, ErrStale
	}
	eol := t.fileEOL()
	text = strings.Trim(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	var repl string
	if text != "" {
		repl = strings.ReplaceAll(text, "\n", eol)
	}
	ls := t.lines(s.Body, s.End)
	first, last := -1, -1
	for i, l := range ls {
		if strings.TrimSpace(string(t.src[l.start:l.end])) != "" {
			if first < 0 {
				first = i
			}
			last = i
		}
	}
	var out []byte
	switch {
	case first < 0: // no text yet
		out = append(out, t.src[:s.Body]...)
		if !bytes.HasSuffix(t.src[:s.Body], []byte("\n")) {
			out = append(out, eol...)
		}
		if repl != "" {
			out = append(out, eol...)
			out = append(out, repl...)
			out = append(out, eol...)
		}
		out = append(out, t.src[s.Body:]...)
	default:
		a := ls[first].start
		b := ls[last].end
		if repl == "" {
			// Clearing removes the lines whole. When the heading is followed by a
			// blank line (Prettier's style) the blank lines after the text go too,
			// so one blank line stays between the heading and the next section; a
			// section that ends the file keeps only its heading.
			if first > 0 {
				b = s.End
				if s.End >= len(t.src) {
					a = s.Body
				}
			} else {
				b += len(ls[last].eol)
			}
		}
		out = append(out, t.src[:a]...)
		out = append(out, repl...)
		out = append(out, t.src[b:]...)
	}
	return out, nil
}

// ---- Adding a section and appending text ----

// AddSection creates the section name with text and returns the new body. It goes
// before the first "## Log" section, otherwise at the end. The heading is written
// as given, or capitalised when the name is all lowercase; a section with text has
// a blank line after its heading, as Prettier writes it, and one blank line
// separates it from its neighbours. It does not check that the name is free:
// callers look first.
func (t *Tree) AddSection(name, text string) ([]byte, error) {
	eol := t.fileEOL()
	name = strings.TrimSpace(name)
	if name == strings.ToLower(name) {
		r := []rune(name)
		name = strings.ToUpper(string(r[:1])) + string(r[1:])
	}
	block := "## " + name + eol
	if text = strings.Trim(strings.ReplaceAll(text, "\r\n", "\n"), "\n"); text != "" {
		block += eol + strings.ReplaceAll(text, "\n", eol) + eol
	}
	if lg, ok := t.Section("log"); ok {
		pre := t.src[:lg.Start]
		out := append([]byte(nil), pre...)
		if ls := t.lines(0, lg.Start); len(ls) > 0 && strings.TrimSpace(string(t.src[ls[len(ls)-1].start:ls[len(ls)-1].end])) != "" {
			out = append(out, eol...)
		}
		out = append(out, block...)
		out = append(out, eol...)
		return append(out, t.src[lg.Start:]...), nil
	}
	out := append([]byte(nil), t.src...)
	if len(out) > 0 {
		if !bytes.HasSuffix(out, []byte("\n")) {
			out = append(out, eol...)
		}
		ls := t.lines(0, len(t.src))
		if last := ls[len(ls)-1]; strings.TrimSpace(string(t.src[last.start:last.end])) != "" {
			out = append(out, eol...)
		}
	}
	return append(out, block...), nil
}

// AppendText adds text after the last non-blank line of section name and returns
// the new body; a missing section is created with it (AddSection). The text is
// written as given, one line ending per line. Nothing is written for empty text.
func (t *Tree) AppendText(name, text string) ([]byte, error) {
	text = strings.TrimRight(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	if text == "" {
		return t.Join(), nil
	}
	s, ok := t.Section(name)
	if !ok {
		return t.AddSection(name, text)
	}
	if err := t.refusedIn(s); err != nil {
		return nil, err
	}
	return t.appendLine(s, strings.ReplaceAll(text, "\n", t.fileEOL())), nil
}
