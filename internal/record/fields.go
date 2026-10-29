package record

import (
	"strings"

	"github.com/alexsasharegan/akmctl/internal/frontmatter"
)

var (
	workFields = []string{"id", "title", "summary", "status", "kind", "parent", "blocked_by", "priority", "created", "updated", "authored_by", "decisions"}
	claimField = []string{"claimed_by", "claimed_at"}
	adrFields  = []string{"id", "title", "summary", "status", "date", "scope", "supersedes", "superseded_by", "authored_by"}
	workHeader = []string{"id", "title", "summary", "status", "kind", "parent", "priority", "updated"}
	adrHeader  = []string{"id", "title", "summary", "status", "date", "scope"}
)

// TypeFields lists the frontmatter fields of a record type in spec order.
func TypeFields(typ string, claims bool) []string {
	if typ == ADR {
		return adrFields
	}
	if claims {
		return append(append([]string(nil), workFields...), claimField...)
	}
	return workFields
}

// HeaderFields is the default field set of ls, search and board.
func HeaderFields(typ string) []string {
	if typ == ADR {
		return adrHeader
	}
	return workHeader
}

// refFields hold ids or wikilinks to records and are returned as bare ids.
var refFields = map[string]bool{"parent": true, "blocked_by": true, "decisions": true, "supersedes": true, "superseded_by": true}

// JSONValue converts a field value for output: references become bare ids.
func JSONValue(key string, v frontmatter.Value) any {
	switch v.Kind {
	case frontmatter.String:
		if refFields[key] {
			return RefID(v.Str)
		}
		return v.Str
	case frontmatter.Int:
		return v.Int
	case frontmatter.Bool:
		return v.Bool
	case frontmatter.List:
		out := make([]any, 0, len(v.Items))
		for _, it := range v.Items {
			out = append(out, JSONValue(key, it))
		}
		return out
	}
	return nil
}

// Sections splits a body into its "## " sections: names (lowercase) in order and
// their text. Fenced code blocks are skipped when looking for headings.
func Sections(body string) (names []string, text map[string]string) {
	text = map[string]string{}
	var cur string
	var buf []string
	fence := false
	flush := func() {
		if cur != "" {
			text[cur] = strings.Trim(strings.Join(buf, "\n"), "\n")
		}
		buf = nil
	}
	for _, l := range strings.Split(body, "\n") {
		l = strings.TrimRight(l, "\r")
		if strings.HasPrefix(l, "```") || strings.HasPrefix(l, "~~~") {
			fence = !fence
		}
		if !fence && strings.HasPrefix(l, "## ") {
			flush()
			cur = strings.ToLower(strings.TrimSpace(l[3:]))
			names = append(names, cur)
			continue
		}
		buf = append(buf, l)
	}
	flush()
	return names, text
}

// FirstHeading is the text of the first "# " heading of a body outside fenced
// code blocks, or "".
func FirstHeading(body string) string {
	fence := false
	for _, l := range strings.Split(body, "\n") {
		l = strings.TrimRight(l, "\r")
		if strings.HasPrefix(l, "```") || strings.HasPrefix(l, "~~~") {
			fence = !fence
		}
		if !fence && strings.HasPrefix(l, "# ") {
			return strings.TrimSpace(l[2:])
		}
	}
	return ""
}

// NextLine is the first line of the Next section, or "".
func (r *Record) NextLine() (string, bool) {
	_, text := Sections(r.Body)
	t, ok := text["next"]
	if !ok {
		return "", false
	}
	for _, l := range strings.Split(t, "\n") {
		if strings.TrimSpace(l) != "" {
			return l, true
		}
	}
	return "", true
}

// AppendToSection adds a line at the end of the "## name" section of a body and
// returns the new body; the section is created at the end when it is missing.
// eol is the line ending to use.
func AppendToSection(body, name, line, eol string) string {
	lines := strings.Split(body, "\n")
	fence := false
	start, end := -1, len(lines)
	for i, l := range lines {
		t := strings.TrimRight(l, "\r")
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			fence = !fence
		}
		if fence || !strings.HasPrefix(t, "## ") {
			continue
		}
		if start >= 0 {
			end = i
			break
		}
		if strings.EqualFold(strings.TrimSpace(t[3:]), name) {
			start = i
		}
	}
	if start < 0 {
		b := strings.TrimRight(body, "\r\n")
		heading := strings.ToUpper(name[:1]) + name[1:]
		return b + eol + eol + "## " + heading + eol + line + eol
	}
	last := start
	for i := start + 1; i < end; i++ {
		if strings.TrimSpace(lines[i]) != "" {
			last = i
		}
	}
	out := append([]string(nil), lines[:last+1]...)
	cr := ""
	if eol == "\r\n" {
		cr = "\r"
	}
	out = append(out, line+cr)
	out = append(out, lines[last+1:]...)
	return strings.Join(out, "\n")
}

// SetSection replaces the text of the "## name" section of a body, or appends
// lines to it, and returns the new body. The heading is matched without regard
// to case, outside code fences. A missing section is created: before the "## Log"
// section when there is one, otherwise at the end. Blank lines around the
// section text are kept, and eol is the line ending to write.
func SetSection(body, name, text string, appendMode bool, eol string) string {
	cr := ""
	if eol == "\r\n" {
		cr = "\r"
	}
	name = strings.TrimSpace(name)
	text = strings.ReplaceAll(text, "\r\n", "\n")
	if appendMode {
		text = strings.TrimRight(text, "\n")
	} else {
		text = strings.Trim(text, "\n")
	}
	var tl []string
	if text != "" {
		for _, l := range strings.Split(text, "\n") {
			tl = append(tl, l+cr)
		}
	}
	if appendMode && tl == nil {
		return body
	}
	var lines []string
	endsNL := strings.HasSuffix(body, "\n")
	if body != "" {
		lines = strings.Split(strings.TrimSuffix(body, "\n"), "\n")
	}
	blank := func(l string) bool { return strings.TrimSpace(l) == "" }
	start, end, logAt := -1, len(lines), -1
	fence := false
	for i, l := range lines {
		t := strings.TrimRight(l, "\r")
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			fence = !fence
		}
		if fence || !strings.HasPrefix(t, "## ") {
			continue
		}
		h := strings.TrimSpace(t[3:])
		if start >= 0 {
			end = i
			break
		}
		if strings.EqualFold(h, name) {
			start = i
		} else if strings.EqualFold(h, "log") && logAt < 0 {
			logAt = i
		}
	}
	var out []string
	if start < 0 {
		heading := name
		if heading == strings.ToLower(heading) {
			r := []rune(heading)
			heading = strings.ToUpper(string(r[0])) + string(r[1:])
		}
		block := []string{"## " + heading + cr}
		if tl != nil {
			block = append(block, cr) // a blank line after the heading, as Prettier writes it
		}
		block = append(block, tl...)
		if logAt >= 0 {
			out = append(out, lines[:logAt]...)
			if logAt > 0 && !blank(lines[logAt-1]) {
				out = append(out, cr)
			}
			out = append(out, block...)
			out = append(out, cr)
			out = append(out, lines[logAt:]...)
		} else {
			out = append(out, lines...)
			if len(lines) > 0 && !blank(lines[len(lines)-1]) {
				out = append(out, cr)
			}
			out = append(out, block...)
			endsNL = true
		}
		return strings.Join(out, "\n") + nl(endsNL)
	}
	content := lines[start+1 : end]
	first, last := -1, -1
	for i, l := range content {
		if !blank(l) {
			if first < 0 {
				first = i
			}
			last = i
		}
	}
	var lead, core, trail []string
	if first < 0 {
		lead = content
	} else {
		lead, core, trail = content[:first], content[first:last+1], content[last+1:]
	}
	if appendMode {
		core = append(append([]string(nil), core...), tl...)
	} else {
		core = tl
	}
	if len(trail) == 0 && end < len(lines) {
		trail = []string{cr}
	}
	out = append(out, lines[:start+1]...)
	out = append(out, lead...)
	out = append(out, core...)
	out = append(out, trail...)
	out = append(out, lines[end:]...)
	if end == len(lines) {
		endsNL = true
	}
	return strings.Join(out, "\n") + nl(endsNL)
}

func nl(b bool) string {
	if b {
		return "\n"
	}
	return ""
}
