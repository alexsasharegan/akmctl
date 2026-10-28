package frontmatter

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Error is a parse error with the 1-based line number in the file (the opening
// delimiter is line 1).
type Error struct {
	Line int
	Msg  string
}

func (e *Error) Error() string { return fmt.Sprintf("line %d: %s", e.Line, e.Msg) }

// Entry is one key of the block with its value and byte spans.
type Entry struct {
	Key   string
	Value Value
	Line  int // line of the key

	start, end       int // whole entry, first byte of the key line to the end of the last line (no eol)
	keyStart, keyEnd int
	colonEnd         int // offset just after the colon
	valStart, valEnd int // value text on the key line; both colonEnd when empty or a block list
	block            bool
	indent           int // indentation of block list items
}

// Doc is a parsed frontmatter block. Edits go through Set, Delete and Rename.
type Doc struct {
	Entries []Entry

	src       []byte
	eol       string
	delim     bool // the block is wrapped in --- lines
	blockEnd  int  // offset where an appended entry goes
	bodyStart int  // offset of the first byte after the closing delimiter line
}

var bom = []byte{0xEF, 0xBB, 0xBF}

type line struct {
	off  int
	text string
	n    int
	eol  string
}

func splitLines(src []byte, from, firstN int) []line {
	var out []line
	n := firstN
	for off := from; off < len(src); n++ {
		i := bytes.IndexByte(src[off:], '\n')
		var l line
		if i < 0 {
			l = line{off, string(src[off:]), n, ""}
			off = len(src)
		} else {
			t := src[off : off+i]
			e := "\n"
			if len(t) > 0 && t[len(t)-1] == '\r' {
				t = t[:len(t)-1]
				e = "\r\n"
			}
			l = line{off, string(t), n, e}
			off += i + 1
		}
		out = append(out, l)
	}
	return out
}

// Parse reads a record's frontmatter: the file starts with a --- line, and the
// block ends at the next line that is exactly ---.
func Parse(src []byte) (*Doc, error) {
	if bytes.HasPrefix(src, bom) {
		return nil, &Error{1, "byte-order mark"}
	}
	ls := splitLines(src, 0, 1)
	if len(ls) == 0 || ls[0].text != "---" {
		return nil, &Error{1, "frontmatter must start with a --- line"}
	}
	d := &Doc{src: src, eol: ls[0].eol, delim: true}
	if d.eol == "" {
		return nil, &Error{1, "unterminated frontmatter"}
	}
	closeAt := -1
	for i := 1; i < len(ls); i++ {
		if ls[i].text == "---" {
			closeAt = i
			break
		}
	}
	if closeAt < 0 {
		return nil, &Error{len(ls), "unterminated frontmatter: no closing --- line"}
	}
	d.blockEnd = ls[closeAt].off
	d.bodyStart = ls[closeAt].off + len(ls[closeAt].text) + len(ls[closeAt].eol)
	if err := d.parseBlock(ls[1:closeAt]); err != nil {
		return nil, err
	}
	return d, nil
}

// ParseBare reads a block with no --- delimiters, such as the manifest
// docs/.layout.
func ParseBare(src []byte) (*Doc, error) {
	if bytes.HasPrefix(src, bom) {
		return nil, &Error{1, "byte-order mark"}
	}
	ls := splitLines(src, 0, 1)
	d := &Doc{src: src, eol: "\n", blockEnd: len(src), bodyStart: len(src)}
	if len(ls) > 0 && ls[0].eol != "" {
		d.eol = ls[0].eol
	}
	if err := d.parseBlock(ls); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Doc) parseBlock(ls []line) error {
	seen := map[string]bool{}
	var cur *Entry
	for _, l := range ls {
		if !utf8.ValidString(l.text) {
			return &Error{l.n, "invalid UTF-8"}
		}
		if strings.Contains(l.text, "\t") {
			return &Error{l.n, "tab character"}
		}
		for _, r := range l.text {
			if !printable(r) {
				return &Error{l.n, fmt.Sprintf("control or non-character code point U+%04X", r)}
			}
		}
		t := strings.TrimLeft(l.text, " ")
		if t == "" || t[0] == '#' {
			continue
		}
		ind := len(l.text) - len(t)
		if t == "-" || strings.HasPrefix(t, "- ") {
			if cur == nil || (cur.Value.Kind != Null && !cur.block) || (cur.Value.Kind == Null && cur.valEnd != cur.colonEnd) {
				return &Error{l.n, "list item without an empty key above it"}
			}
			if !cur.block {
				cur.block = true
				cur.indent = ind
				cur.Value = Value{Kind: List}
			} else if ind != cur.indent {
				return &Error{l.n, "list items must be indented by the same number of spaces"}
			}
			rest := t[1:]
			sp := len(rest) - len(strings.TrimLeft(rest, " "))
			v, end, msg := scalarAt(l.text, ind+1+sp, false)
			if msg != "" {
				return &Error{l.n, msg}
			}
			if msg = trailer(l.text, end); msg != "" {
				return &Error{l.n, msg}
			}
			cur.Value.Items = append(cur.Value.Items, v)
			cur.end = l.off + len(l.text)
			continue
		}
		if ind > 0 {
			return &Error{l.n, "indented line"}
		}
		e, err := parseEntry(l)
		if err != nil {
			return err
		}
		if seen[e.Key] {
			return &Error{l.n, "duplicate key " + e.Key}
		}
		seen[e.Key] = true
		d.Entries = append(d.Entries, e)
		cur = &d.Entries[len(d.Entries)-1]
	}
	return nil
}

func isKeyStart(c byte) bool { return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') }
func isKeyChar(c byte) bool {
	return isKeyStart(c) || (c >= '0' && c <= '9') || c == '.' || c == '-'
}

func parseEntry(l line) (Entry, error) {
	s := l.text
	if s == "" || !isKeyStart(s[0]) {
		return Entry{}, &Error{l.n, "expected a key"}
	}
	i := 1
	for i < len(s) && isKeyChar(s[i]) {
		i++
	}
	if i >= len(s) || s[i] != ':' || (i+1 < len(s) && s[i+1] != ' ') {
		return Entry{}, &Error{l.n, "expected \"key: value\""}
	}
	key := s[:i]
	if len(key) > 255 {
		return Entry{}, &Error{l.n, "key longer than 255 characters"}
	}
	switch strings.ToLower(key) {
	case "true", "false", "null":
		return Entry{}, &Error{l.n, "key " + key + " is not a string key in YAML 1.2"}
	}
	e := Entry{Key: key, Line: l.n, start: l.off, end: l.off + len(s), keyStart: l.off, keyEnd: l.off + i, colonEnd: l.off + i + 1}
	vs := i + 1
	for vs < len(s) && s[vs] == ' ' {
		vs++
	}
	if vs == len(s) || s[vs] == '#' {
		e.valStart, e.valEnd = e.colonEnd, e.colonEnd
		return e, nil
	}
	var v Value
	var end int
	var msg string
	if s[vs] == '[' {
		v, end, msg = listAt(s, vs)
	} else {
		v, end, msg = scalarAt(s, vs, false)
	}
	if msg != "" {
		return Entry{}, &Error{l.n, msg}
	}
	if msg = trailer(s, end); msg != "" {
		return Entry{}, &Error{l.n, msg}
	}
	e.Value = v
	e.valStart, e.valEnd = l.off+vs, l.off+end
	return e, nil
}

// trailer checks that only spaces and a comment follow a value ending at end.
func trailer(s string, end int) string {
	rest := s[end:]
	t := strings.TrimLeft(rest, " ")
	if t == "" {
		return ""
	}
	if t[0] == '#' && len(rest) > len(t) {
		return ""
	}
	return "unexpected text after the value"
}

// listAt parses an inline list starting at s[i] == '['.
func listAt(s string, i int) (Value, int, string) {
	v := Value{Kind: List}
	i++
	for {
		for i < len(s) && s[i] == ' ' {
			i++
		}
		if i >= len(s) {
			return v, 0, "unterminated or multi-line list"
		}
		switch s[i] {
		case ']':
			return v, i + 1, ""
		case ',':
			return v, 0, "empty list item"
		case '[', '{':
			return v, 0, "nested list or mapping"
		}
		it, end, msg := scalarAt(s, i, true)
		if msg != "" {
			return v, 0, msg
		}
		v.Items = append(v.Items, it)
		i = end
		for i < len(s) && s[i] == ' ' {
			i++
		}
		if i >= len(s) {
			return v, 0, "unterminated or multi-line list"
		}
		if s[i] == ',' {
			i++
		} else if s[i] != ']' {
			return v, 0, "expected , or ] in a list"
		}
	}
}

// scalarAt parses one scalar starting at s[i]; flow is true inside an inline
// list. It returns the value and the offset just after it.
func scalarAt(s string, i int, flow bool) (Value, int, string) {
	if i >= len(s) {
		return Value{}, 0, "empty list item"
	}
	switch s[i] {
	case '"':
		var b strings.Builder
		for j := i + 1; j < len(s); j++ {
			switch s[j] {
			case '\\':
				if j+1 >= len(s) || (s[j+1] != '"' && s[j+1] != '\\') {
					return Value{}, 0, "unsupported escape (only \\\" and \\\\)"
				}
				b.WriteByte(s[j+1])
				j++
			case '"':
				return NewString(b.String()), j + 1, ""
			default:
				b.WriteByte(s[j])
			}
		}
		return Value{}, 0, "unterminated or multi-line quoted string"
	case '\'':
		var b strings.Builder
		for j := i + 1; j < len(s); j++ {
			if s[j] == '\'' {
				if j+1 < len(s) && s[j+1] == '\'' {
					b.WriteByte('\'')
					j++
					continue
				}
				return NewString(b.String()), j + 1, ""
			}
			b.WriteByte(s[j])
		}
		return Value{}, 0, "unterminated or multi-line quoted string"
	}
	// plain
	c := s[i]
	if strings.IndexByte(",[]{}#&*!|>%@`", c) >= 0 {
		switch c {
		case '&', '*':
			return Value{}, 0, "anchors and aliases are not supported"
		case '!':
			return Value{}, 0, "tags are not supported"
		case '|', '>':
			return Value{}, 0, "block scalars are not supported"
		case '[', '{':
			return Value{}, 0, "flow collections are not supported here"
		}
		return Value{}, 0, "plain value cannot start with " + string(c)
	}
	if flow && c == ':' {
		return Value{}, 0, "plain value in a list cannot start with \":\""
	}
	if (c == '-' || c == '?' || c == ':') && (i+1 == len(s) || s[i+1] == ' ' || (flow && (s[i+1] == ',' || s[i+1] == ']'))) {
		return Value{}, 0, "plain value cannot start with \"" + string(c) + " \""
	}
	j := i
	for j < len(s) {
		if s[j] == '#' && j > i && s[j-1] == ' ' {
			break
		}
		if flow && (s[j] == ',' || s[j] == ']') {
			break
		}
		if flow && (s[j] == '[' || s[j] == '{' || s[j] == '}' || s[j] == '?') {
			return Value{}, 0, "plain value in a list cannot contain " + string(s[j])
		}
		j++
	}
	text := strings.TrimRight(s[i:j], " ")
	end := i + len(text)
	if strings.Contains(text, ": ") || strings.HasSuffix(text, ":") {
		return Value{}, 0, "\": \" in a plain value (quote it)"
	}
	v, msg := typePlain(text)
	return v, end, msg
}

// typePlain gives an unquoted scalar its YAML 1.2 core type, or an error for a
// form YAML-lite does not support.
func typePlain(s string) (Value, string) {
	switch s {
	case "null", "~":
		return NewNull(), ""
	case "true":
		return NewBool(true), ""
	case "false":
		return NewBool(false), ""
	case "True", "TRUE", "False", "FALSE", "Null", "NULL":
		return Value{}, "capitalised " + s + " is not supported (write it lower case or quote it)"
	}
	if isCanonicalInt(s) {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return Value{}, "integer out of range"
		}
		return NewInt(n), ""
	}
	if numberLike(s) {
		return Value{}, "number form not supported (only canonical decimal integers; quote it for text)"
	}
	return NewString(s), ""
}

func isCanonicalInt(s string) bool {
	t := strings.TrimPrefix(s, "-")
	if t == "" || (len(t) > 1 && t[0] == '0') {
		return false
	}
	for i := 0; i < len(t); i++ {
		if t[i] < '0' || t[i] > '9' {
			return false
		}
	}
	return true
}

// numberLike reports whether s is a number in YAML 1.2 core or 1.1 other than a
// canonical integer: signs, leading zeros, underscores, floats, hex, octal,
// exponents, .inf and .nan.
func numberLike(s string) bool {
	if strings.Contains(s, "_") {
		if u := strings.ReplaceAll(s, "_", ""); u != s && numberLike(u) {
			return true
		}
	}
	t := strings.TrimLeft(s, "+-")
	if len(s)-len(t) > 1 {
		return false
	}
	switch strings.ToLower(t) {
	case ".inf", ".nan":
		return true
	}
	if len(t) > 2 && t[0] == '0' && strings.IndexByte("xXoObB", t[1]) >= 0 {
		return true // 0x1F, 0o17, 0b101 in any case: a number in some YAML versions
	}
	i, digits := 0, 0
	for i < len(t) && (t[i] >= '0' && t[i] <= '9' || t[i] == '_') {
		if t[i] != '_' {
			digits++
		}
		i++
	}
	frac := 0
	hasDot := false
	if i < len(t) && t[i] == '.' {
		hasDot = true
		i++
		for i < len(t) && (t[i] >= '0' && t[i] <= '9' || t[i] == '_') {
			frac++
			i++
		}
	}
	if digits == 0 && frac == 0 {
		return false
	}
	if digits == 0 && !hasDot {
		return false
	}
	if i < len(t) && (t[i] == 'e' || t[i] == 'E') {
		i++
		if i < len(t) && (t[i] == '+' || t[i] == '-') {
			i++
		}
		e := 0
		for i < len(t) && t[i] >= '0' && t[i] <= '9' {
			e++
			i++
		}
		if e == 0 {
			return false
		}
	}
	return i == len(t)
}

// printable is YAML's c-printable without tab and line breaks, which are handled
// separately: a block holding anything else is invalid YAML.
func printable(r rune) bool {
	switch {
	case r >= 0x20 && r <= 0x7E, r == 0x85:
		return true
	case r >= 0xA0 && r <= 0xD7FF, r >= 0xE000 && r <= 0xFFFD, r >= 0x10000 && r <= 0x10FFFF:
		return true
	}
	return false
}
