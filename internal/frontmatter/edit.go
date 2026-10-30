package frontmatter

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Get returns the value of key. A key with no value reads as Null and present.
func (d *Doc) Get(key string) (Value, bool) {
	if e := d.entry(key); e != nil {
		return e.Value, true
	}
	return Value{}, false
}

// Line returns the 1-based line of key, or 0.
func (d *Doc) Line(key string) int {
	if e := d.entry(key); e != nil {
		return e.Line
	}
	return 0
}

// Bytes returns the whole file, edits included.
func (d *Doc) Bytes() []byte { return d.src }

// Body returns the bytes after the closing delimiter line (everything after the
// block for a bare document: nothing).
func (d *Doc) Body() []byte { return d.src[d.bodyStart:] }

func (d *Doc) entry(key string) *Entry {
	for i := range d.Entries {
		if d.Entries[i].Key == key {
			return &d.Entries[i]
		}
	}
	return nil
}

// Set changes the value of key, or appends the key at the end of the block. Only
// the bytes of that entry change. Setting a value equal to the current one
// writes nothing. A trailing comment on the key line is kept.
func (d *Doc) Set(key string, v Value) error {
	if !validKey(key) {
		return errors.New("invalid key " + strconv.Quote(key))
	}
	e := d.entry(key)
	if e != nil && e.Value.Equal(v) {
		return nil
	}
	flat, err := formatValue(v, d.eol, e)
	if err != nil {
		return err
	}
	if e == nil {
		return d.reparse(d.insertAt(d.blockEnd, key+": "+flat.inline+d.eol))
	}
	var edits []edit
	switch {
	case e.block && v.Kind == List && len(v.Items) > 0:
		edits = append(edits, edit{e.keyEnd + 1, e.end, flat.blockTail})
	case e.block:
		// one edit: the value goes on the key line (before any comment) and the items go
		edits = append(edits, edit{e.colonEnd, e.end, " " + flat.inline + string(d.src[e.colonEnd:d.keyLineEnd(e)])})
	case e.valStart == e.valEnd:
		edits = append(edits, edit{e.colonEnd, e.colonEnd, " " + flat.inline})
	default:
		edits = append(edits, edit{e.valStart, e.valEnd, flat.inline})
	}
	return d.reparse(apply(d.src, edits))
}

// SetAfter is Set, except that a key that is not present yet is inserted on the
// line after the key named after (or at the end of the block when after is
// absent too).
func (d *Doc) SetAfter(key, after string, v Value) error {
	a := d.entry(after)
	if d.entry(key) != nil || a == nil {
		return d.Set(key, v)
	}
	if !validKey(key) {
		return errors.New("invalid key " + strconv.Quote(key))
	}
	flat, err := formatValue(v, d.eol, nil)
	if err != nil {
		return err
	}
	return d.reparse(apply(d.src, []edit{{a.end, a.end, d.eol + key + ": " + flat.inline}}))
}

// Delete removes the entry for key. It reports whether the key existed.
func (d *Doc) Delete(key string) (bool, error) {
	e := d.entry(key)
	if e == nil {
		return false, nil
	}
	end := e.end
	switch {
	case strings.HasPrefix(string(d.src[end:]), "\r\n"):
		end += 2
	case strings.HasPrefix(string(d.src[end:]), "\n"):
		end++
	}
	start := e.start
	if end == len(d.src) && start > 0 {
		// last line of a bare file with no trailing newline: take the previous newline instead
		start -= len(d.eol)
	}
	return true, d.reparse(apply(d.src, []edit{{start, end, ""}}))
}

// Rename changes a key, keeping its value and position.
func (d *Doc) Rename(from, to string) error {
	e := d.entry(from)
	if e == nil {
		return errors.New("no key " + from)
	}
	if !validKey(to) {
		return errors.New("invalid key " + strconv.Quote(to))
	}
	if d.entry(to) != nil {
		return errors.New("key " + to + " already exists")
	}
	return d.reparse(apply(d.src, []edit{{e.keyStart, e.keyEnd, to}}))
}

type edit struct {
	start, end int
	text       string
}

func apply(src []byte, edits []edit) []byte {
	// edits are given in any order but never overlap; apply from the highest offset down
	out := append([]byte(nil), src...)
	for len(edits) > 0 {
		hi := 0
		for i, e := range edits {
			if e.start > edits[hi].start {
				hi = i
			}
		}
		e := edits[hi]
		edits = append(edits[:hi:hi], edits[hi+1:]...)
		out = append(out[:e.start:e.start], append([]byte(e.text), out[e.end:]...)...)
	}
	return out
}

func (d *Doc) insertAt(off int, text string) []byte {
	src := d.src
	if !d.delim && off > 0 && src[off-1] != '\n' {
		text = d.eol + text
	}
	return apply(src, []edit{{off, off, text}})
}

func (d *Doc) keyLineEnd(e *Entry) int {
	i := e.start
	for i < len(d.src) && d.src[i] != '\n' && d.src[i] != '\r' {
		i++
	}
	return i
}

func (d *Doc) reparse(src []byte) error {
	var n *Doc
	var err error
	if d.delim {
		n, err = Parse(src)
	} else {
		n, err = ParseBare(src)
	}
	if err != nil {
		return errors.New("internal: edit produced an invalid block: " + err.Error())
	}
	*d = *n
	return nil
}

func validKey(k string) bool {
	if k == "" || !isKeyStart(k[0]) {
		return false
	}
	for i := 1; i < len(k); i++ {
		if !isKeyChar(k[i]) {
			return false
		}
	}
	switch strings.ToLower(k) {
	case "true", "false", "null":
		return false
	}
	return true
}

type formatted struct {
	inline    string // text after "key: "
	blockTail string // for an existing block list: the text after "key:" (items on following lines)
	block     []string
}

func formatValue(v Value, eol string, e *Entry) (formatted, error) {
	var f formatted
	if v.Kind != List {
		s, err := formatScalar(v, false)
		f.inline = s
		return f, err
	}
	parts := make([]string, len(v.Items))
	for i, it := range v.Items {
		if it.Kind == List {
			return f, errors.New("nested lists are not supported")
		}
		s, err := formatScalar(it, true)
		if err != nil {
			return f, err
		}
		parts[i] = s
	}
	f.inline = "[" + strings.Join(parts, ", ") + "]"
	if e != nil && e.block {
		var b strings.Builder
		pad := strings.Repeat(" ", e.indent)
		for _, p := range parts {
			// block items are plain scalars, which may contain commas; reuse the flow form, which is also valid here
			b.WriteString(eol + pad + "- " + p)
		}
		f.blockTail = b.String()
		f.block = parts
	}
	return f, nil
}

func formatScalar(v Value, flow bool) (string, error) {
	switch v.Kind {
	case Null:
		return "null", nil
	case Bool:
		return strconv.FormatBool(v.Bool), nil
	case Int:
		return strconv.FormatInt(v.Int, 10), nil
	case String:
		if canPlain(v.Str, flow) {
			return v.Str, nil
		}
		if !utf8.ValidString(v.Str) {
			return "", errors.New("string is not valid UTF-8")
		}
		var b strings.Builder
		b.WriteByte('"')
		for i := 0; i < len(v.Str); i++ {
			c := v.Str[i]
			switch {
			case c == '"' || c == '\\':
				b.WriteByte('\\')
				b.WriteByte(c)
			case c < 0x20 || c == 0x7f:
				return "", errors.New("string holds a control character or line break; YAML-lite cannot write it")
			default:
				b.WriteByte(c)
			}
		}
		b.WriteByte('"')
		return b.String(), nil
	}
	return "", errors.New("cannot format value")
}

// canPlain reports whether s reads back as exactly the string s when written unquoted.
func canPlain(s string, flow bool) bool {
	if s == "" || s != strings.TrimSpace(s) || strings.ContainsAny(s, "\t\r\n") {
		return false
	}
	if s[0] == '"' || s[0] == '\'' {
		return false
	}
	v, end, msg := scalarAt(s, 0, flow)
	return msg == "" && end == len(s) && v.Kind == String && v.Str == s
}
