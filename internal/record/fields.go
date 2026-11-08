package record

import (
	"strings"

	mdbody "github.com/alexsasharegan/akmctl/internal/body"
	"github.com/alexsasharegan/akmctl/internal/frontmatter"
)

var (
	workFields = []string{"id", "title", "summary", "status", "kind", "parent", "blocked_by", "priority", "cycle", "keep", "created", "updated", "closed", "authored_by", "decisions"}
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
	var fence mdbody.FenceTracker
	flush := func() {
		if cur != "" {
			text[cur] = strings.Trim(strings.Join(buf, "\n"), "\n")
		}
		buf = nil
	}
	for _, l := range strings.Split(body, "\n") {
		l = strings.TrimRight(l, "\r")
		if !fence.In(l) && strings.HasPrefix(l, "## ") {
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
	var fence mdbody.FenceTracker
	for _, l := range strings.Split(body, "\n") {
		l = strings.TrimRight(l, "\r")
		if !fence.In(l) && strings.HasPrefix(l, "# ") {
			return strings.TrimSpace(l[2:])
		}
	}
	return ""
}

// RetitleHeading rewrites the first "# " heading of a body outside fenced code
// blocks from old to new, and reports whether it did. A body whose first heading
// is not old is left alone: its author made it differ, and the tool does not
// choose between them (rule WK12 warns).
func RetitleHeading(body, old, new string) (string, bool) {
	var fence mdbody.FenceTracker
	lines := strings.Split(body, "\n")
	for i, l := range lines {
		t := strings.TrimRight(l, "\r")
		if fence.In(t) || !strings.HasPrefix(t, "# ") {
			continue
		}
		if strings.TrimSpace(t[2:]) != old {
			return body, false
		}
		lines[i] = "# " + new + l[len(t):]
		return strings.Join(lines, "\n"), true
	}
	return body, false
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
// returns the new body; the section is created when it is missing. eol is the
// line ending for a body that has none of its own. The error is an
// *mdbody.UnsupportedError when the body has a construct the tool cannot edit
// around (ADR 261101-tr8e).
func AppendToSection(body, name, line, eol string) (string, error) {
	tr, err := mdbody.Split([]byte(body))
	if err != nil {
		return "", err
	}
	if _, ok := tr.Section(name); !ok {
		// the section goes at the end, its first line straight under the heading
		b := strings.TrimRight(body, "\r\n")
		return b + eol + eol + "## " + strings.ToUpper(name[:1]) + name[1:] + eol + line + eol, nil
	}
	return EditSection(body, name, line, true, eol)
}

// listItem reports whether text starts with a Markdown list marker: "- ", "* ",
// "+ " or a number followed by "." or ")" and a space.
func listItem(text string) bool {
	if len(text) > 1 && strings.ContainsRune("-*+", rune(text[0])) && text[1] == ' ' {
		return true
	}
	i := 0
	for i < len(text) && text[i] >= '0' && text[i] <= '9' {
		i++
	}
	return i > 0 && i+1 < len(text) && (text[i] == '.' || text[i] == ')') && text[i+1] == ' '
}

// EditSection replaces the text of the "## name" section of a body, or appends
// lines to it, through the section tree, and returns the new body. The heading is
// matched without regard to case, outside code fences. A missing section is
// created: before the "## Log" section when there is one, otherwise at the end.
// Bytes outside the section are never changed. A body with a structural
// construct, or a local one in the section, is refused with an
// *mdbody.UnsupportedError and nothing is written.
func EditSection(body, name, text string, appendMode bool, eol string) (string, error) {
	tr, err := mdbody.Split([]byte(body))
	if err != nil {
		return "", err
	}
	tr.WithEOL(eol)
	name = strings.TrimSpace(name)
	text = strings.ReplaceAll(text, "\r\n", "\n")
	if appendMode {
		text = strings.TrimRight(text, "\n")
		if strings.EqualFold(name, "log") && text != "" && !listItem(text) {
			// a bare line would join the previous bullet
			text = "- " + text
		}
		out, err := tr.AppendText(name, text)
		return string(out), err
	}
	if _, ok := tr.Section(name); !ok {
		out, err := tr.AddSection(name, text)
		return string(out), err
	}
	sec, _ := tr.Section(name)
	out, err := tr.ReplaceSection(name, text, "")
	if err == nil && sec.End == len(body) && len(out) > 0 && out[len(out)-1] != '\n' {
		// the last section ends the file with a line ending
		out = append(out, tr.WithEOL(eol).LineEnding()...)
	}
	return string(out), err
}
