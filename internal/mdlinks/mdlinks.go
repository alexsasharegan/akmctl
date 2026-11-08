// Package mdlinks finds the relative Markdown links of a repository's files and
// rewrites their destinations. Rule L1 of check and the link repair of close
// (work item 261031-wjs3) read links through the same scanner, so what one
// accepts the other repairs.
package mdlinks

import (
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/body"
	"github.com/alexsasharegan/akmctl/internal/record"
)

var (
	inline   = regexp.MustCompile(`\[[^\]]*\]\(([^)\n]*)\)`)
	codeSpan = regexp.MustCompile("`[^`\n]*`")
)

// Link is one inline link whose destination is a relative file path.
type Link struct {
	Line   int    // 1-based line of the link
	Raw    string // the destination as written between the parentheses
	Target string // the file path in it: no title, fragment, query or escapes
	// Start and End delimit Target as written, in bytes of the scanned text; a
	// fragment, a query, a title and angle brackets lie outside them.
	Start, End int
}

// Scan lists the relative file links of a Markdown file. It skips the
// frontmatter, fenced code blocks and code spans, and links that are external,
// absolute or only an anchor.
func Scan(text string) []Link {
	var out []Link
	off := 0
	var fence body.FenceTracker
	start := frontmatterEnd(text)
	line := 0
	for off < len(text) {
		line++
		nl := strings.IndexByte(text[off:], '\n')
		end := len(text)
		if nl >= 0 {
			end = off + nl
		}
		l := strings.TrimRight(text[off:end], "\r")
		lineOff := off
		off = end + 1
		if lineOff < start {
			continue
		}
		if fence.In(strings.TrimLeft(l, " ")) {
			continue
		}
		masked := codeSpan.ReplaceAllStringFunc(l, func(s string) string { return strings.Repeat(" ", len(s)) })
		for _, m := range inline.FindAllStringSubmatchIndex(masked, -1) {
			raw := masked[m[2]:m[3]]
			s, e, target, ok := destination(raw)
			if !ok {
				continue
			}
			out = append(out, Link{Line: line, Raw: raw, Target: target, Start: lineOff + m[2] + s, End: lineOff + m[2] + e})
		}
	}
	return out
}

// frontmatterEnd is the offset just after the closing "---" line of a leading
// frontmatter block, or 0 when the text has none.
func frontmatterEnd(text string) int {
	first, _, _ := strings.Cut(text, "\n")
	if strings.TrimRight(first, "\r") != "---" {
		return 0
	}
	off := len(first) + 1
	for off < len(text) {
		nl := strings.IndexByte(text[off:], '\n')
		end := len(text)
		if nl >= 0 {
			end = off + nl
		}
		if strings.TrimRight(text[off:end], "\r") == "---" {
			return end + 1
		}
		off = end + 1
	}
	return 0
}

// destination picks the file path out of a link destination (raw, as written
// between the parentheses). s and e delimit it within raw; target is it
// unescaped. ok is false for external, absolute and anchor-only destinations.
func destination(raw string) (s, e int, target string, ok bool) {
	i := 0
	for i < len(raw) && (raw[i] == ' ' || raw[i] == '\t') {
		i++
	}
	s = i
	if i < len(raw) && raw[i] == '<' {
		s = i + 1
		e = strings.IndexByte(raw[s:], '>')
		if e < 0 {
			return 0, 0, "", false
		}
		e += s
	} else {
		e = len(raw)
		if j := strings.IndexAny(raw[s:], " \t"); j >= 0 {
			e = s + j // a title follows
		}
	}
	t := raw[s:e]
	if t == "" || strings.HasPrefix(t, "#") || strings.HasPrefix(t, "/") || strings.Contains(t, "://") || strings.HasPrefix(t, "mailto:") {
		return 0, 0, "", false
	}
	if j := strings.IndexAny(t, "#?"); j >= 0 {
		e = s + j
		t = t[:j]
	}
	if u, err := url.PathUnescape(t); err == nil {
		t = u
	}
	return s, e, t, t != ""
}

// Rewrite maps the target of every link in text through fn and returns the new
// text and the number of links changed. fn gets the target as Scan reports it
// and returns the replacement path, or ok false to leave the link alone. Only the
// path is replaced; a fragment, query, title and angle brackets stay.
func Rewrite(text string, fn func(target string) (string, bool)) (string, int) {
	links := Scan(text)
	var b strings.Builder
	last, n := 0, 0
	for _, l := range links {
		nt, ok := fn(l.Target)
		if !ok {
			continue
		}
		if strings.HasPrefix(l.Target, "./") && !strings.HasPrefix(nt, "../") && !strings.HasPrefix(nt, "./") {
			nt = "./" + nt
		}
		if nt == l.Target {
			continue
		}
		old := text[l.Start:l.End]
		b.WriteString(text[last:l.Start])
		b.WriteString(encode(nt, strings.Contains(old, "%") || strings.ContainsAny(nt, " ()")))
		last = l.End
		n++
	}
	b.WriteString(text[last:])
	return b.String(), n
}

// encode writes a path as a link destination: escaped per segment when the
// original was escaped or the path holds a character that would end the link.
func encode(p string, escape bool) string {
	if !escape {
		return p
	}
	segs := strings.Split(p, "/")
	for i, s := range segs {
		s = url.PathEscape(s)
		s = strings.NewReplacer("(", "%28", ")", "%29").Replace(s)
		segs[i] = s
	}
	return strings.Join(segs, "/")
}

// Rel is the path of to as seen from the folder from, both relative to the
// repository root, with forward slashes.
func Rel(from, to string) string {
	f := splitPath(from)
	t := splitPath(to)
	i := 0
	for i < len(f) && i < len(t)-1 && f[i] == t[i] {
		i++
	}
	parts := make([]string, 0, len(f)-i+len(t)-i)
	for range f[i:] {
		parts = append(parts, "..")
	}
	parts = append(parts, t[i:]...)
	return strings.Join(parts, "/")
}

func splitPath(p string) []string {
	p = path.Clean(p)
	if p == "." || p == "" {
		return nil
	}
	return strings.Split(p, "/")
}

// Files lists the Markdown files whose links are read, relative to the
// repository root: AGENTS.md, every Markdown file under the manifest, work and
// ADR folders, and the roadmap. A folder whose name starts with a dot is skipped, as
// Obsidian skips it (work item 261029-z4gk).
func Files(repo *record.Repo) []string {
	var files []string
	seen := map[string]bool{}
	add := func(rel string) {
		if !seen[rel] {
			seen[rel] = true
			files = append(files, rel)
		}
	}
	if _, err := os.Stat(filepath.Join(repo.Root, "AGENTS.md")); err == nil {
		add("AGENTS.md")
	}
	for _, dir := range roots(repo) {
		top := filepath.Join(repo.Root, filepath.FromSlash(dir))
		filepath.WalkDir(top, func(p string, d fs.DirEntry, err error) error {
			if err == nil && d.IsDir() && p != top && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir // a hidden folder (.trash, .fixtures) is not part of the notes
			}
			if err == nil && !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
				rel, _ := filepath.Rel(repo.Root, p)
				add(filepath.ToSlash(rel))
			}
			return nil
		})
	}
	if rm := repo.RoadmapPath(); !seen[rm] {
		if _, err := os.Stat(filepath.Join(repo.Root, filepath.FromSlash(rm))); err == nil {
			add(rm)
		}
	}
	return files
}

// roots are the folders whose Markdown is read: the manifest folder (docs/ by
// default) and the work and ADR folders when they are outside it.
func roots(repo *record.Repo) []string {
	var out []string
	for _, r := range []string{repo.MetaDir(), repo.WorkDir(), repo.ADRDir()} {
		covered := false
		for _, o := range out {
			if r == o || strings.HasPrefix(r, o+"/") {
				covered = true
			}
		}
		if !covered {
			out = append(out, r)
		}
	}
	return out
}
