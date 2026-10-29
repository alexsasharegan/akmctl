package check

import (
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	mdLink   = regexp.MustCompile(`\[[^\]]*\]\(([^)\n]*)\)`)
	wikiLink = regexp.MustCompile(`\[\[([^\]\n|#]*)(?:#[^\]\n|]*)?(?:\|[^\]\n]*)?\]\]`)
	codeSpan = regexp.MustCompile("`[^`\n]*`")
)

// links is rules L1 (relative Markdown links) and L2 (wikilinks) over AGENTS.md
// and every Markdown file under docs/, the work folder and the ADR folder (they
// differ from docs/ when the manifest moves them), and the roadmap file.
func (c *checker) links() {
	var files []string
	if _, err := os.Stat(filepath.Join(c.repo.Root, "AGENTS.md")); err == nil {
		files = append(files, "AGENTS.md")
	}
	seen := map[string]bool{"AGENTS.md": true}
	for _, dir := range c.linkRoots() {
		filepath.WalkDir(filepath.Join(c.repo.Root, filepath.FromSlash(dir)), func(p string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
				rel, _ := filepath.Rel(c.repo.Root, p)
				rel = filepath.ToSlash(rel)
				if !seen[rel] {
					seen[rel] = true
					files = append(files, rel)
				}
			}
			return nil
		})
	}
	if rm := c.repo.RoadmapPath(); !seen[rm] {
		if _, err := os.Stat(filepath.Join(c.repo.Root, filepath.FromSlash(rm))); err == nil {
			files = append(files, rm)
		}
	}
	var byName map[string][]string // wikilink index, built on first use
	for _, f := range files {
		b, err := os.ReadFile(filepath.Join(c.repo.Root, filepath.FromSlash(f)))
		if err != nil {
			continue
		}
		for _, ln := range bodyLines(string(b)) {
			text := codeSpan.ReplaceAllString(ln.text, "")
			for _, m := range mdLink.FindAllStringSubmatch(text, -1) {
				target, skip := linkTarget(m[1])
				if skip {
					continue
				}
				p := path.Join(path.Dir(f), target)
				if _, err := os.Stat(filepath.Join(c.repo.Root, filepath.FromSlash(p))); err != nil {
					c.add("L1", "E", f, ln.n, "", "link %s does not resolve to a file", m[1])
				}
			}
			for _, m := range wikiLink.FindAllStringSubmatch(text, -1) {
				if byName == nil {
					byName = c.markdownIndex()
				}
				name := strings.TrimSpace(m[1])
				if name == "" {
					continue
				}
				switch n := len(resolveWiki(byName, name)); {
				case n == 0:
					c.add("L2", "W", f, ln.n, "", "wikilink [[%s]] does not resolve to a file", name)
				case n > 1:
					c.add("L2", "W", f, ln.n, "", "wikilink [[%s]] matches %d files", name, n)
				}
			}
		}
	}
}

type numbered struct {
	n    int
	text string
}

// bodyLines returns the lines of a Markdown file that can hold links: not the
// frontmatter and not fenced code blocks. Line numbers are those of the file.
func bodyLines(s string) []numbered {
	lines := strings.Split(s, "\n")
	start := 0
	if len(lines) > 0 && strings.TrimRight(lines[0], "\r") == "---" {
		for i := 1; i < len(lines); i++ {
			if strings.TrimRight(lines[i], "\r") == "---" {
				start = i + 1
				break
			}
		}
	}
	var out []numbered
	fence := false
	for i := start; i < len(lines); i++ {
		l := strings.TrimRight(lines[i], "\r")
		if strings.HasPrefix(strings.TrimLeft(l, " "), "```") || strings.HasPrefix(strings.TrimLeft(l, " "), "~~~") {
			fence = !fence
			continue
		}
		if !fence {
			out = append(out, numbered{i + 1, l})
		}
	}
	return out
}

// linkTarget cleans a Markdown link destination. skip is true for external
// links, anchors and anything that is not a relative file path.
func linkTarget(raw string) (string, bool) {
	t := strings.TrimSpace(raw)
	if strings.HasPrefix(t, "<") {
		if i := strings.Index(t, ">"); i > 0 {
			t = t[1:i]
		}
	} else if i := strings.IndexAny(t, " \t"); i > 0 {
		t = t[:i] // a title after the destination
	}
	if t == "" || strings.HasPrefix(t, "#") || strings.HasPrefix(t, "/") || strings.Contains(t, "://") || strings.HasPrefix(t, "mailto:") {
		return "", true
	}
	if i := strings.IndexAny(t, "#?"); i >= 0 {
		t = t[:i]
	}
	if u, err := url.PathUnescape(t); err == nil {
		t = u
	}
	return t, t == ""
}

// markdownIndex maps a file name without .md to the paths that have it.
func (c *checker) markdownIndex() map[string][]string {
	idx := map[string][]string{}
	filepath.WalkDir(c.repo.Root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && (d.Name() == ".git" || d.Name() == "node_modules") {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			rel, _ := filepath.Rel(c.repo.Root, p)
			idx[strings.TrimSuffix(d.Name(), ".md")] = append(idx[strings.TrimSuffix(d.Name(), ".md")], filepath.ToSlash(rel))
		}
		return nil
	})
	return idx
}

// resolveWiki finds the files a wikilink can mean: by file name, or by path
// when the name has a slash.
func resolveWiki(idx map[string][]string, name string) []string {
	name = strings.TrimSuffix(name, ".md")
	if !strings.Contains(name, "/") {
		return idx[name]
	}
	var out []string
	for _, paths := range idx {
		for _, p := range paths {
			if strings.TrimSuffix(p, ".md") == name || strings.HasSuffix(strings.TrimSuffix(p, ".md"), "/"+name) {
				out = append(out, p)
			}
		}
	}
	return out
}

// linkRoots are the folders whose Markdown the link rules read: docs/ and the
// work and ADR folders when they are outside it.
func (c *checker) linkRoots() []string {
	roots := []string{"docs", c.repo.WorkDir(), c.repo.ADRDir()}
	var out []string
	for _, r := range roots {
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
