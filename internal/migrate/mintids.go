package migrate

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alexsasharegan/akmctl/internal/frontmatter"
	"github.com/alexsasharegan/akmctl/internal/fsx"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// IDAlphabet is the alphabet of the four-character id suffix.
const IDAlphabet = "0123456789abcdefghjkmnpqrstvwxyz"

// Env is what a run needs from outside: the random source for id suffixes (nil
// reads the system source) and the clock for a record with no usable date (the
// zero time reads the system clock). Tests set both.
type Env struct {
	Rand io.Reader
	Now  time.Time
}

func (e Env) rand() io.Reader {
	if e.Rand != nil {
		return e.Rand
	}
	return rand.Reader
}

func (e Env) now() time.Time {
	if !e.Now.IsZero() {
		return e.Now.UTC()
	}
	return time.Now().UTC()
}

var (
	datePrefixRe = regexp.MustCompile(`^(?:\d{4}-\d{2}-\d{2}|\d{8})(?:[-_ .]+|$)`)
	nameIDRe     = regexp.MustCompile(`^(\d{6}-[0-9a-hjkmnp-tv-z]{4})-.+`)
	idValueRe    = regexp.MustCompile(`^\d{6}-[0-9a-hjkmnp-tv-z]{4}$`)
	slugJunkRe   = regexp.MustCompile(`[^a-z0-9]+`)
)

// slugify makes lowercase words joined by hyphens, at most 60 characters. Text
// outside ASCII letters and digits separates words and is dropped.
func slugify(s string) string {
	s = strings.Trim(slugJunkRe.ReplaceAllString(strings.ToLower(s), "-"), "-")
	if len(s) > 60 {
		s = strings.TrimRight(s[:60], "-")
	}
	return s
}

// slugFor takes the slug from the file name (extension and date prefix
// removed), then from the title, then "untitled".
func slugFor(name, title string) string {
	base := strings.TrimSuffix(name, path.Ext(name))
	base = datePrefixRe.ReplaceAllString(base, "")
	if s := slugify(base); s != "" {
		return s
	}
	if s := slugify(title); s != "" {
		return s
	}
	return "untitled"
}

// idDate gives the YYMMDD part for a date field value: a date, or an ISO
// timestamp (read as UTC when it carries an offset).
func idDate(v string) (string, bool) {
	v = strings.TrimSpace(v)
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02T15:04", "2006-01-02 15:04:05", "2006-01-02 15:04"} {
		if t, err := time.Parse(layout, v); err == nil {
			if layout == time.RFC3339Nano {
				t = t.UTC()
			}
			return t.Format("060102"), true
		}
	}
	if len(v) >= 10 {
		if t, err := time.Parse("2006-01-02", v[:10]); err == nil && (len(v) == 10 || strings.ContainsAny(v[10:11], "T ")) {
			return t.Format("060102"), true
		}
	}
	return "", false
}

type mintCtx struct {
	repo *record.Repo
	root string
	used map[string]bool
	env  Env
	seed string // when set, the suffix comes from a hash of seed and file, not the random source
}

// mint draws an id for date. With a seed, the suffix is the first four
// characters of sha256(seed, name, try): the same seed and name give the same
// id on every run, and the try counter redraws on a collision.
func (m *mintCtx) mint(date, name string) (string, error) {
	for i := 0; i < 200; i++ {
		var b [4]byte
		if m.seed != "" {
			sum := sha256.Sum256([]byte(m.seed + "\x00" + name + "\x00" + strconv.Itoa(i)))
			copy(b[:], sum[:4])
		} else if _, err := io.ReadFull(m.env.rand(), b[:]); err != nil {
			return "", err
		}
		suffix := make([]byte, 4)
		for j := range b {
			suffix[j] = IDAlphabet[b[j]&31]
		}
		id := date + "-" + string(suffix)
		if !m.used[id] {
			m.used[id] = true
			return id, nil
		}
	}
	return "", errors.New("no free id after 200 tries")
}

// usedIDs collects the ids already present in the three record directories and
// in dir: file name prefixes and `id` fields of every Markdown file.
func (m *mintCtx) collectUsed(dir string, recursive bool) {
	scan := func(abs string, rec bool) {
		filepath.WalkDir(abs, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if p != abs && !rec {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.EqualFold(filepath.Ext(d.Name()), ".md") {
				return nil
			}
			if mm := nameIDRe.FindStringSubmatch(d.Name()); mm != nil {
				m.used[mm[1]] = true
			}
			if b, err := os.ReadFile(p); err == nil {
				if doc, err := frontmatter.Parse(b); err == nil {
					if v, ok := doc.Get("id"); ok && v.Kind == frontmatter.String && idValueRe.MatchString(v.Str) {
						m.used[v.Str] = true
					}
				}
			}
			return nil
		})
	}
	for _, d := range []string{m.repo.ADRDir(), m.repo.WorkDir(), m.repo.DoneDir()} {
		scan(filepath.Join(m.root, filepath.FromSlash(d)), false)
	}
	scan(filepath.Join(m.root, filepath.FromSlash(dir)), recursive)
}

type mintItem struct {
	oldRel, newRel string
	id             string
	text           []byte // the new file text
	renamed        bool
	planned        string // with rename false, the name the caller should give the file
	note           string
}

// mintIDs plans (and unless dry applies) the mint_ids step. files are the new
// repository-relative paths; notes say what the tool decided or could not do.
func mintIDs(repo *record.Repo, s Step, dry bool, env Env) (files, notes []string, renames []Rename, err error) {
	root := repo.Root
	files, notes = []string{}, []string{}
	if s.Dir == "" || !safeRel(s.Dir) {
		return files, notes, nil, fmt.Errorf("mint_ids: dir %q must be a repository-relative path", s.Dir)
	}
	doRename := s.Rename == nil || *s.Rename
	dir := path.Clean(s.Dir)
	absDir := filepath.Join(root, filepath.FromSlash(dir))
	if st, e := os.Stat(absDir); e != nil || !st.IsDir() {
		return files, []string{"directory " + dir + " does not exist; nothing to do"}, nil, nil
	}
	field := s.DateField
	if field == "" {
		field = "created"
	}
	m := &mintCtx{repo: repo, root: root, used: map[string]bool{}, env: env, seed: s.Seed}
	m.collectUsed(dir, s.Recursive)

	// the candidates, in path order so a seeded run is repeatable
	var cands []string
	filepath.WalkDir(absDir, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return nil
		}
		if d.IsDir() {
			if p != absDir && (!s.Recursive || strings.HasPrefix(d.Name(), ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		n := d.Name()
		if !d.Type().IsRegular() || strings.HasPrefix(n, ".") || !strings.EqualFold(filepath.Ext(n), ".md") || strings.EqualFold(n, "README.md") {
			return nil
		}
		cands = append(cands, p)
		return nil
	})
	sort.Strings(cands)

	var items []mintItem
	for _, p := range cands {
		rel, _ := filepath.Rel(root, p)
		rel = filepath.ToSlash(rel)
		b, e := os.ReadFile(p)
		if e != nil {
			return files, notes, nil, e
		}
		it, note, skip := planOne(m, rel, b, field, doRename)
		if note != "" {
			notes = append(notes, note)
		}
		if skip {
			continue
		}
		items = append(items, *it)
	}
	if len(items) == 0 {
		return files, notes, nil, nil
	}

	refs := findReferences(root, items)
	for _, it := range items {
		files = append(files, it.newRel)
		notes = append(notes, it.note)
		switch {
		case it.renamed:
			renames = append(renames, Rename{it.oldRel, it.newRel})
		case it.planned != "":
			renames = append(renames, Rename{it.oldRel, it.planned})
		}
	}
	notes = append(notes, refs...)
	if dry {
		sort.Strings(files)
		return files, notes, renames, nil
	}
	for _, it := range items {
		from, to := filepath.Join(root, filepath.FromSlash(it.oldRel)), filepath.Join(root, filepath.FromSlash(it.newRel))
		if it.renamed {
			if _, e := os.Stat(to); e == nil {
				return files, notes, renames, fmt.Errorf("mint_ids: %s already exists", it.newRel)
			}
			if e := fsx.Rename(from, to); e != nil {
				return files, notes, renames, e
			}
		}
		if e := fsx.WriteAtomic(to, it.text); e != nil {
			return files, notes, renames, e
		}
	}
	sort.Strings(files)
	return files, notes, renames, nil
}

// planOne decides one file. skip is true when the file stays as it is; note is
// then a report line or empty.
func planOne(m *mintCtx, rel string, b []byte, field string, doRename bool) (it *mintItem, note string, skip bool) {
	name := path.Base(rel)
	dirRel := path.Dir(rel)
	doc, err := frontmatter.Parse(b)
	hasFM := err == nil
	if err != nil {
		var fe *frontmatter.Error
		if !errors.As(err, &fe) || !strings.HasPrefix(fe.Msg, "frontmatter must start") {
			return nil, rel + ": skipped, frontmatter cannot be read (" + err.Error() + ")", true
		}
	}
	title := ""
	if hasFM {
		if _, ok := doc.Get("id"); ok {
			return nil, "", true
		}
		if v, ok := doc.Get("title"); ok && v.Kind == frontmatter.String {
			title = v.Str
		}
	}
	if title == "" {
		for _, l := range strings.Split(string(b), "\n") {
			if strings.HasPrefix(l, "# ") {
				title = strings.TrimSpace(strings.TrimPrefix(l, "# "))
				break
			}
		}
	}
	eol := "\n"
	if i := strings.IndexByte(string(b), '\n'); i > 0 && b[i-1] == '\r' {
		eol = "\r\n"
	}

	var extra []string
	// a name that already carries an id keeps it: only the field is written
	if mm := nameIDRe.FindStringSubmatch(name); mm != nil {
		id := mm[1]
		text, e := withID(b, doc, hasFM, id, eol)
		if e != nil {
			return nil, rel + ": skipped, " + e.Error(), true
		}
		note = fmt.Sprintf("%s: name already starts with id %s; wrote the id field, no rename", rel, id)
		if !hasFM {
			note += "; the file had no frontmatter, created it"
		}
		return &mintItem{oldRel: rel, newRel: rel, id: id, text: text, note: note}, "", false
	}

	date, src := "", ""
	if hasFM {
		if v, ok := doc.Get(field); ok && v.Kind == frontmatter.String {
			if d, ok := idDate(v.Str); ok {
				date, src = d, "from "+field
			}
		}
	}
	if date == "" {
		date = m.env.now().Format("060102")
		src = "today (" + field + " absent or not a date)"
	}
	id, e := m.mint(date, rel)
	if e != nil {
		return nil, rel + ": skipped, " + e.Error(), true
	}
	newName := id + "-" + slugFor(name, title) + ".md"
	newRel := newName
	if dirRel != "." {
		newRel = dirRel + "/" + newName
	}
	text, e := withID(b, doc, hasFM, id, eol)
	if e != nil {
		return nil, rel + ": skipped, " + e.Error(), true
	}
	if !hasFM {
		extra = append(extra, "the file had no frontmatter, created it")
	}
	note = fmt.Sprintf("%s -> %s (id %s, date %s)", rel, newRel, id, src)
	if !doRename {
		// the id is written; the caller renames the file (Obsidian rewrites links)
		extra = append(extra, "not renamed: rename the file as listed in renames")
		if len(extra) > 0 {
			note += "; " + strings.Join(extra, "; ")
		}
		return &mintItem{oldRel: rel, newRel: rel, planned: newRel, id: id, text: text, note: note}, "", false
	}
	if len(extra) > 0 {
		note += "; " + strings.Join(extra, "; ")
	}
	return &mintItem{oldRel: rel, newRel: newRel, id: id, text: text, renamed: true, note: note}, "", false
}

// withID puts `id: "<id>"` first in the frontmatter, creating the block when
// the file has none.
func withID(b []byte, doc *frontmatter.Doc, hasFM bool, id, eol string) ([]byte, error) {
	line := `id: "` + id + `"` + eol
	var out []byte
	if hasFM {
		first := strings.IndexByte(string(b), '\n') + 1
		out = append(append(append([]byte(nil), b[:first]...), line...), b[first:]...)
	} else {
		out = append([]byte("---"+eol+line+"---"+eol+eol), b...)
	}
	chk, err := frontmatter.Parse(out)
	if err != nil {
		return nil, err
	}
	if v, ok := chk.Get("id"); !ok || v.Str != id {
		return nil, errors.New("the id field did not read back")
	}
	return out, nil
}

// findReferences reports Markdown files that mention an old file name, as a
// wikilink or as a path or link ending in the name. Nothing is rewritten.
func findReferences(root string, items []mintItem) []string {
	type target struct {
		old, base string
		re        *regexp.Regexp
		newRel    string
	}
	var ts []target
	moved := map[string]string{}
	for _, it := range items {
		if !it.renamed {
			continue
		}
		moved[it.oldRel] = it.newRel
		base := strings.TrimSuffix(path.Base(it.oldRel), path.Ext(it.oldRel))
		q := regexp.QuoteMeta(base)
		ts = append(ts, target{it.oldRel, base, regexp.MustCompile(`\[\[(?:[^\]|#\n]*/)?` + q + `(?:\.md)?(?:[\]|#])|(?:^|[^A-Za-z0-9_.\-])` + q + `\.md\b`), it.newRel})
	}
	if len(ts) == 0 {
		return nil
	}
	hits := map[string][]string{} // key: referrer + "\x00" + old
	var order []string
	filepath.WalkDir(root, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return nil
		}
		if d.IsDir() {
			if p != root && (strings.HasPrefix(d.Name(), ".") || d.Name() == "node_modules") {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() || !strings.EqualFold(filepath.Ext(d.Name()), ".md") {
			return nil
		}
		b, e := os.ReadFile(p)
		if e != nil {
			return nil
		}
		rel, _ := filepath.Rel(root, p)
		rel = filepath.ToSlash(rel)
		shown := rel
		if n, ok := moved[rel]; ok {
			shown = n
		}
		for ln, l := range strings.Split(string(b), "\n") {
			for _, t := range ts {
				if t.old == rel || !t.re.MatchString(l) {
					continue
				}
				k := shown + "\x00" + t.old
				if _, ok := hits[k]; !ok {
					order = append(order, k)
				}
				hits[k] = append(hits[k], fmt.Sprint(ln+1))
			}
		}
		return nil
	})
	sort.Strings(order)
	var out []string
	for _, k := range order {
		parts := strings.SplitN(k, "\x00", 2)
		out = append(out, fmt.Sprintf("reference not rewritten: %s line %s names %s (now %s)", parts[0], strings.Join(hits[k], ", "), path.Base(parts[1]), path.Base(moved[parts[1]])))
	}
	return out
}
