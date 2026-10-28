// Package record loads the records of an AKM repository: work items, ADRs and
// the roadmap, from the files under docs/.
package record

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/frontmatter"
)

const (
	Work = "work"
	ADR  = "adr"
)

var (
	idRe   = regexp.MustCompile(`^\d{6}-[0-9a-hjkmnp-tv-z]{4}$`)
	fileRe = regexp.MustCompile(`^(\d{6}-[0-9a-hjkmnp-tv-z]{4})-.+\.md$`)
	wikiRe = regexp.MustCompile(`^\[\[(\d{6}-[0-9a-hjkmnp-tv-z]{4})[^\]|]*(\|[^\]]*)?\]\]$`)
)

// ValidID reports whether s matches the identifier pattern.
func ValidID(s string) bool { return idRe.MatchString(s) }

// RefID reduces a reference to its bare id: an id stays as it is, a quoted
// wikilink gives the leading id of its target, and any other text (free text in
// blocked_by) is returned unchanged.
func RefID(s string) string {
	if m := wikiRe.FindStringSubmatch(s); m != nil {
		return m[1]
	}
	return s
}

// Problem is a file the loader could not use.
type Problem struct {
	Path string
	Line int
	Msg  string
}

// Record is one ADR or work item.
type Record struct {
	ID     string
	Type   string
	Path   string // relative to the repository root, with /
	Keys   []string
	Fields map[string]frontmatter.Value
	Lines  map[string]int
	Body   string
	Etag   string
	Text   []byte
}

// IDDate is the UTC date part of the id as YYYY-MM-DD.
func (r *Record) IDDate() string { return IDDate(r.ID) }

// IDDate converts the YYMMDD prefix of an id to YYYY-MM-DD.
func IDDate(id string) string { return "20" + id[0:2] + "-" + id[2:4] + "-" + id[4:6] }

// Str returns a string field, or "".
func (r *Record) Str(k string) string {
	v := r.Fields[k]
	if v.Kind == frontmatter.String {
		return v.Str
	}
	return ""
}

// Status returns the status field.
func (r *Record) Status() string { return r.Str("status") }

// Repo is a loaded repository.
type Repo struct {
	Root     string
	Records  []*Record // sorted by id, then path
	Roadmap  []string  // ids under Now, Next and Later in file order
	Dropped  []string
	Layout   map[string]frontmatter.Value
	Problems []Problem
}

// FindRoot returns the nearest ancestor of dir that has docs/.layout or AGENTS.md.
func FindRoot(dir string) (string, error) {
	d, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	for {
		for _, m := range []string{"docs/.layout", "AGENTS.md"} {
			if _, err := os.Stat(filepath.Join(d, filepath.FromSlash(m))); err == nil {
				return d, nil
			}
		}
		p := filepath.Dir(d)
		if p == d {
			return "", errors.New("no repository root (docs/.layout or AGENTS.md) above " + dir)
		}
		d = p
	}
}

// FromText builds a record from the text of its file: rel is the path relative
// to the repository root, typ is Work or ADR and id is the leading id of the
// file name.
func FromText(rel, typ, id string, text []byte) (*Record, error) {
	doc, err := frontmatter.Parse(text)
	if err != nil {
		return nil, err
	}
	rec := &Record{ID: id, Type: typ, Path: rel, Fields: map[string]frontmatter.Value{}, Lines: map[string]int{}, Body: string(doc.Body()), Text: text}
	sum := sha256.Sum256(text)
	rec.Etag = hex.EncodeToString(sum[:])[:16]
	for _, en := range doc.Entries {
		rec.Keys = append(rec.Keys, en.Key)
		rec.Fields[en.Key] = en.Value
		rec.Lines[en.Key] = en.Line
	}
	return rec, nil
}

// Load reads every record of the repository at root. Files that cannot be used
// are listed in Problems; the rest load.
func Load(root string) (*Repo, error) {
	if st, err := os.Stat(root); err != nil || !st.IsDir() {
		return nil, errors.New("cannot read repository at " + root)
	}
	r := &Repo{Root: root, Layout: map[string]frontmatter.Value{}}
	for _, d := range []struct{ dir, typ string }{{"docs/work", Work}, {"docs/work/done", Work}, {"docs/adr", ADR}} {
		ents, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(d.dir)))
		if err != nil {
			continue
		}
		for _, e := range ents {
			m := fileRe.FindStringSubmatch(e.Name())
			if e.IsDir() || m == nil {
				continue
			}
			rel := d.dir + "/" + e.Name()
			text, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
			if err != nil {
				r.Problems = append(r.Problems, Problem{rel, 0, err.Error()})
				continue
			}
			rec, err := FromText(rel, d.typ, m[1], text)
			if err != nil {
				p := Problem{rel, 0, err.Error()}
				var fe *frontmatter.Error
				if errors.As(err, &fe) {
					p.Line, p.Msg = fe.Line, fe.Msg
				}
				r.Problems = append(r.Problems, p)
				continue
			}
			r.Records = append(r.Records, rec)
		}
	}
	sort.SliceStable(r.Records, func(i, j int) bool {
		if r.Records[i].ID != r.Records[j].ID {
			return r.Records[i].ID < r.Records[j].ID
		}
		return r.Records[i].Path < r.Records[j].Path
	})
	r.loadRoadmap()
	r.loadLayout()
	return r, nil
}

func (r *Repo) loadLayout() {
	b, err := os.ReadFile(filepath.Join(r.Root, "docs", ".layout"))
	if err != nil {
		return
	}
	doc, err := frontmatter.ParseBare(b)
	if err != nil {
		r.Problems = append(r.Problems, Problem{"docs/.layout", 0, err.Error()})
		return
	}
	for _, e := range doc.Entries {
		r.Layout[e.Key] = e.Value
	}
}

// Extension reports whether the manifest enables an extension.
func (r *Repo) Extension(name string) bool {
	_, ok := r.Layout["extension."+name]
	return ok
}

var roadmapLine = regexp.MustCompile(`^- \[(\d{6}-[0-9a-hjkmnp-tv-z]{4})\]`)

func (r *Repo) loadRoadmap() {
	b, err := os.ReadFile(filepath.Join(r.Root, "docs", "ROADMAP.md"))
	if err != nil {
		return
	}
	section := ""
	for _, l := range strings.Split(string(b), "\n") {
		l = strings.TrimRight(l, "\r")
		if strings.HasPrefix(l, "## ") {
			section = strings.TrimSpace(l[3:])
			continue
		}
		if m := roadmapLine.FindStringSubmatch(l); m != nil {
			switch section {
			case "Now", "Next", "Later":
				r.Roadmap = append(r.Roadmap, m[1])
			case "Dropped":
				r.Dropped = append(r.Dropped, m[1])
			}
		}
	}
}

// RoadmapRank returns the position of id on the roadmap (Now, Next, Later), or -1.
func (r *Repo) RoadmapRank(id string) int {
	for i, x := range r.Roadmap {
		if x == id {
			return i
		}
	}
	return -1
}

// ByID returns every record with the id (more than one means a collision).
func (r *Repo) ByID(id string) []*Record {
	var out []*Record
	for _, rec := range r.Records {
		if rec.ID == id {
			out = append(out, rec)
		}
	}
	return out
}

// Reload reads the repository again from disk into r.
func (r *Repo) Reload() error {
	n, err := Load(r.Root)
	if err != nil {
		return err
	}
	*r = *n
	return nil
}
