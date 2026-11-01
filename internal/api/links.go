package api

import (
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/alexsasharegan/akmctl/internal/fsx"
	"github.com/alexsasharegan/akmctl/internal/mdlinks"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// linkEdit is a file whose Markdown links a move rewrites.
type linkEdit struct {
	path string
	text string
	n    int // links changed
}

// moveSelfLinks rebases the relative links of a file that moves from one path to
// another (given relative to the repository root): a link to another file keeps
// pointing at it, and a link to the file itself follows it. A link that does not
// resolve to a file is left as written.
func moveSelfLinks(repo *record.Repo, from, to, text string) (string, int) {
	return mdlinks.Rewrite(text, func(target string) (string, bool) {
		resolved := path.Join(path.Dir(from), target)
		if resolved != from {
			if _, err := os.Stat(filepath.Join(repo.Root, filepath.FromSlash(resolved))); err != nil {
				return "", false
			}
		} else {
			resolved = to
		}
		return mdlinks.Rel(path.Dir(to), resolved), true
	})
}

// planMovedInbound finds the links in the repository's other Markdown files that
// point at from and rewrites them to point at to, reading each file from disk.
// moved is the file that is being moved; its own links are not read here.
func planMovedInbound(repo *record.Repo, from, to string) []linkEdit {
	var out []linkEdit
	for _, f := range mdlinks.Files(repo) {
		if f == from || f == to {
			continue
		}
		b, err := os.ReadFile(filepath.Join(repo.Root, filepath.FromSlash(f)))
		if err != nil {
			continue
		}
		dir := path.Dir(f)
		text, n := mdlinks.Rewrite(string(b), func(target string) (string, bool) {
			if path.Join(dir, target) != from {
				return "", false
			}
			return mdlinks.Rel(dir, to), true
		})
		if n > 0 {
			out = append(out, linkEdit{f, text, n})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].path < out[j].path })
	return out
}

// applyLinkEdits writes the edits, each atomically, and returns the files it
// could not write. A failed file does not stop the others; the caller reports it.
func applyLinkEdits(repo *record.Repo, edits []linkEdit) []Obj {
	var failed []Obj
	for _, e := range edits {
		if err := fsx.WriteAtomic(filepath.Join(repo.Root, filepath.FromSlash(e.path)), []byte(e.text)); err != nil {
			failed = append(failed, O("path", e.path, "error", err.Error()))
		}
	}
	return failed
}

// linkReport lists the files whose links were (or, on a dry run, would be)
// rewritten, in path order, with the number of links changed in each.
func linkReport(self string, selfN int, edits []linkEdit) []Obj {
	var out []Obj
	if selfN > 0 {
		out = append(out, O("path", self, "count", selfN))
	}
	for _, e := range edits {
		out = append(out, O("path", e.path, "count", e.n))
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, _ := out[i].Get("path")
		b, _ := out[j].Get("path")
		return a.(string) < b.(string)
	})
	return out
}
