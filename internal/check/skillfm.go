package check

import (
	"regexp"
	"strings"
)

// blockHeader matches a top-level key whose value is a YAML block scalar: `key: |`, `key: >-`, `key: |2+`.
var blockHeader = regexp.MustCompile(`^([A-Za-z0-9_-]+):[ \t]+([|>])([+-]?)([1-9]?)([+-]?)[ \t]*(#.*)?$`)

// readSkillFrontmatter prepares a SKILL.md for rule S2. Records use the strict frontmatter reader, which
// rejects block scalars (spec, "Frontmatter subset"); a skill's frontmatter is written by tools that fold
// long descriptions with `>` or `|`, so S2 reads those itself. Each top-level block scalar is replaced by
// a quoted placeholder and comment lines, which keeps every other line number, and its real value is
// returned by key. Everything else is left for frontmatter.Parse to accept or reject as it does for records.
func readSkillFrontmatter(src []byte) ([]byte, map[string]string) {
	text := string(src)
	eol := "\n"
	if strings.Contains(text, "\r\n") {
		eol = "\r\n"
	}
	lines := strings.Split(text, "\n")
	for i := range lines {
		lines[i] = strings.TrimSuffix(lines[i], "\r")
	}
	if len(lines) == 0 || lines[0] != "---" {
		return src, nil
	}
	end := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			end = i
			break
		}
	}
	if end < 0 {
		return src, nil
	}
	vals := map[string]string{}
	for i := 1; i < end; i++ {
		m := blockHeader.FindStringSubmatch(lines[i])
		if m == nil {
			continue
		}
		j := i + 1
		for j < end && (strings.TrimSpace(lines[j]) == "" || lines[j][0] == ' ') {
			j++
		}
		body := lines[i+1 : j]
		if _, dup := vals[m[1]]; !dup {
			vals[m[1]] = blockValue(m[2] == ">", m[3]+m[5], m[4], body)
		}
		lines[i] = m[1] + `: "x"`
		for k := i + 1; k < j; k++ {
			lines[k] = "#"
		}
		i = j - 1
	}
	return []byte(strings.Join(lines, eol)), vals
}

// blockValue is the value of a block scalar: style folded or literal, chomp "", "-" or "+", an optional
// explicit indentation digit, and the lines under the header (blank lines included).
func blockValue(folded bool, chomp, indent string, body []string) string {
	// Trailing blank lines belong to the chomping, not to the content.
	n := len(body)
	for n > 0 && strings.TrimSpace(body[n-1]) == "" {
		n--
	}
	trail := len(body) - n
	body = body[:n]
	ind := 0
	if indent != "" {
		ind = int(indent[0] - '0')
	} else {
		for _, l := range body {
			if strings.TrimSpace(l) != "" {
				ind = len(l) - len(strings.TrimLeft(l, " "))
				break
			}
		}
	}
	rows := make([]string, len(body))
	for i, l := range body {
		if len(l) >= ind {
			rows[i] = l[ind:]
		} else {
			rows[i] = strings.TrimLeft(l, " ")
		}
	}
	var out string
	if !folded {
		out = strings.Join(rows, "\n")
	} else {
		var b strings.Builder
		for i, r := range rows {
			if r == "" {
				b.WriteString("\n")
			} else if i > 0 {
				prev := rows[i-1]
				switch {
				case prev == "":
					// the blank line already wrote the break
				case strings.HasPrefix(r, " ") || strings.HasPrefix(r, "\t") || strings.HasPrefix(prev, " ") || strings.HasPrefix(prev, "\t"):
					b.WriteString("\n")
				default:
					b.WriteString(" ")
				}
			}
			b.WriteString(r)
		}
		out = b.String()
	}
	switch {
	case out == "":
		return ""
	case chomp == "-":
		return out
	case chomp == "+":
		return out + "\n" + strings.Repeat("\n", trail)
	}
	return out + "\n"
}
