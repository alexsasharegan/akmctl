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
