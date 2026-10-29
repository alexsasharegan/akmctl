package migrate

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/frontmatter"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// target is one Markdown file a field step reads.
type target struct {
	Path string // repository-relative, with /
	Text []byte
}

// targets lists the files a step works on. With dir set, every Markdown file in
// that folder (the subfolders only with recursive; not README.md, not dot
// files), whatever its name: this reaches notes that have no id yet. Without
// dir, the loaded records, narrowed by scope.
func targets(repo *record.Repo, s Step) ([]target, []string, error) {
	var out []target
	var notes []string
	if s.Dir == "" {
		for _, r := range repo.Records {
			if s.Scope != "" && s.Scope != "all" && s.Scope != r.Type {
				continue
			}
			out = append(out, target{r.Path, r.Text})
		}
		return out, notes, nil
	}
	if !safeRel(s.Dir) {
		return nil, nil, fmt.Errorf("dir %q must be a repository-relative path", s.Dir)
	}
	dir := path.Clean(s.Dir)
	abs := filepath.Join(repo.Root, filepath.FromSlash(dir))
	if st, err := os.Stat(abs); err != nil || !st.IsDir() {
		return nil, []string{"directory " + dir + " does not exist; nothing to do"}, nil
	}
	var paths []string
	filepath.WalkDir(abs, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if p != abs && (!s.Recursive || strings.HasPrefix(d.Name(), ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		n := d.Name()
		if d.Type().IsRegular() && !strings.HasPrefix(n, ".") && strings.EqualFold(filepath.Ext(n), ".md") && !strings.EqualFold(n, "README.md") {
			paths = append(paths, p)
		}
		return nil
	})
	sort.Strings(paths)
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, nil, err
		}
		rel, _ := filepath.Rel(repo.Root, p)
		out = append(out, target{filepath.ToSlash(rel), b})
	}
	return out, notes, nil
}

// parseTarget reads the frontmatter of a target; ok is false for a file to
// leave alone. A file with no frontmatter block is left alone silently, one
// whose block cannot be read is reported.
func parseTarget(t target) (doc *frontmatter.Doc, note string, ok bool) {
	doc, err := frontmatter.Parse(t.Text)
	if err == nil {
		return doc, "", true
	}
	var fe *frontmatter.Error
	if errors.As(err, &fe) && strings.HasPrefix(fe.Msg, "frontmatter must start") {
		return nil, "", false
	}
	return nil, t.Path + ": skipped, frontmatter cannot be read (" + err.Error() + ")", false
}

func statusAllowed(doc *frontmatter.Doc, statuses []string) bool {
	if len(statuses) == 0 {
		return true
	}
	v, _ := doc.Get("status")
	for _, st := range statuses {
		if v.Kind == frontmatter.String && v.Str == st {
			return true
		}
	}
	return false
}

// fieldStep runs rename_field, set_default, remove_field, copy_field and
// title_from_h1 over the targets of the step.
func fieldStep(repo *record.Repo, s Step, write func(string, []byte) error) ([]string, error) {
	ts, notes, err := targets(repo, s)
	if err != nil {
		return notes, err
	}
	for _, t := range ts {
		doc, note, ok := parseTarget(t)
		if note != "" {
			notes = append(notes, note)
		}
		if !ok {
			continue
		}
		before := string(doc.Bytes())
		switch s.Kind {
		case "rename_field":
			_, hasFrom := doc.Get(s.From)
			_, hasTo := doc.Get(s.To)
			if !hasFrom || hasTo {
				continue
			}
			err = doc.Rename(s.From, s.To)
		case "copy_field":
			v, hasFrom := doc.Get(s.From)
			_, hasTo := doc.Get(s.To)
			if !hasFrom || hasTo {
				continue
			}
			err = doc.Set(s.To, v)
		case "remove_field":
			_, err = doc.Delete(s.Field)
		case "set_default":
			if _, has := doc.Get(s.Field); has {
				continue
			}
			var v frontmatter.Value
			if v, err = toValue(s.Value); err == nil {
				err = doc.Set(s.Field, v)
			}
		case "title_from_h1":
			if _, has := doc.Get("title"); has {
				continue
			}
			h := record.FirstHeading(string(doc.Body()))
			if h == "" {
				notes = append(notes, t.Path+": no title and no # heading; left alone")
				continue
			}
			err = doc.Set("title", frontmatter.NewString(h))
		}
		if err != nil {
			return notes, fmt.Errorf("%s: %w", t.Path, err)
		}
		if string(doc.Bytes()) == before {
			continue
		}
		if err := write(t.Path, doc.Bytes()); err != nil {
			return notes, err
		}
	}
	return notes, nil
}

// goalScaffold inserts "## Goal" after the H1 of each record that has none. The
// text between the H1 and the next "## " heading moves under it; a "## Done
// when" section becomes a "Done when:" lead-in with its list, inside Goal. Every
// changed file gets a note, because the text needs a human's review.
func goalScaffold(repo *record.Repo, s Step, write func(string, []byte) error) ([]string, error) {
	ts, notes, err := targets(repo, s)
	if err != nil {
		return notes, err
	}
	for _, t := range ts {
		doc, note, ok := parseTarget(t)
		if note != "" {
			notes = append(notes, note)
		}
		if !ok || !statusAllowed(doc, s.Statuses) {
			continue
		}
		body := string(doc.Body())
		out, how, changed := scaffoldGoal(body)
		if !changed {
			continue
		}
		b := doc.Bytes()
		head := b[:len(b)-len(doc.Body())]
		if err := write(t.Path, append(append([]byte(nil), head...), out...)); err != nil {
			return notes, err
		}
		notes = append(notes, t.Path+": added ## Goal ("+how+"); review the text")
	}
	return notes, nil
}

// scaffoldGoal returns the new body, a short description of what it did and
// whether it changed anything.
func scaffoldGoal(body string) (string, string, bool) {
	crlf := strings.Contains(body, "\r\n")
	if crlf {
		body = strings.ReplaceAll(body, "\r\n", "\n")
	}
	lines := strings.Split(body, "\n")
	// headings outside fences
	type head struct {
		i     int
		level int // 1 or 2
		name  string
	}
	var hs []head
	fence := false
	for i, l := range lines {
		if strings.HasPrefix(l, "```") || strings.HasPrefix(l, "~~~") {
			fence = !fence
			continue
		}
		if fence {
			continue
		}
		switch {
		case strings.HasPrefix(l, "## "):
			hs = append(hs, head{i, 2, strings.ToLower(strings.TrimSpace(strings.TrimRight(strings.TrimSpace(l[3:]), ":")))})
		case strings.HasPrefix(l, "# "):
			hs = append(hs, head{i, 1, ""})
		}
	}
	h1 := -1
	for _, h := range hs {
		if h.level == 2 && h.name == "goal" {
			return body, "", false
		}
		if h.level == 1 && h1 < 0 {
			h1 = h.i
		}
	}
	// the end of the intro: the first "## " heading after the H1
	first := len(lines)
	for _, h := range hs {
		if h.level == 2 && h.i > h1 {
			first = h.i
			break
		}
	}
	// the Done when section: from its heading to the next heading of any level
	dw0, dw1 := -1, -1
	for k, h := range hs {
		if h.level == 2 && h.name == "done when" && h.i > h1 {
			dw0, dw1 = h.i, len(lines)
			if k+1 < len(hs) {
				dw1 = hs[k+1].i
			}
			break
		}
	}
	trim := func(l []string) []string {
		for len(l) > 0 && strings.TrimSpace(l[0]) == "" {
			l = l[1:]
		}
		for len(l) > 0 && strings.TrimSpace(l[len(l)-1]) == "" {
			l = l[:len(l)-1]
		}
		return l
	}
	intro := trim(lines[h1+1 : first]) // h1 is -1 without a heading: the text before the first ## heading
	var done []string
	if dw0 >= 0 {
		done = trim(lines[dw0+1 : dw1])
	}
	var head1 []string
	if h1 >= 0 {
		head1 = lines[:h1+1]
	}
	rest := lines[first:]
	if dw0 >= 0 { // drop the Done when section from the rest
		var kept []string
		for k := first; k < len(lines); k++ {
			if k >= dw0 && k < dw1 {
				continue
			}
			kept = append(kept, lines[k])
		}
		rest = kept
	}
	var out []string
	out = append(out, head1...)
	if h1 >= 0 {
		out = append(out, "")
	}
	out = append(out, "## Goal", "")
	how := "no introductory text"
	if len(intro) > 0 {
		out = append(out, intro...)
		out = append(out, "")
		how = "introduction moved under it"
	}
	if len(done) > 0 {
		out = append(out, "Done when:", "")
		out = append(out, done...)
		out = append(out, "")
		if len(intro) > 0 {
			how += ", Done when kept as a list"
		} else {
			how = "Done when kept as a list"
		}
	}
	if len(rest) > 0 && strings.TrimSpace(strings.Join(rest, "")) != "" {
		out = append(out, trimLeadingBlank(rest)...)
	}
	res := strings.Join(out, "\n")
	if !strings.HasSuffix(res, "\n") {
		res += "\n"
	}
	if crlf {
		res = strings.ReplaceAll(res, "\n", "\r\n")
	}
	return res, how, true
}

func trimLeadingBlank(l []string) []string {
	for len(l) > 0 && strings.TrimSpace(l[0]) == "" {
		l = l[1:]
	}
	return l
}
