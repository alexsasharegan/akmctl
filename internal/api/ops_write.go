package api

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	mdbody "github.com/alexsasharegan/akmctl/internal/body"
	"github.com/alexsasharegan/akmctl/internal/depgraph"
	"github.com/alexsasharegan/akmctl/internal/frontmatter"
	"github.com/alexsasharegan/akmctl/internal/fsx"
	"github.com/alexsasharegan/akmctl/internal/record"
)

const idAlphabet = "0123456789abcdefghjkmnpqrstvwxyz"

var reSlug = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// timeNow is the system clock; tests replace it.
var timeNow = time.Now

// clock is the time of a write: as_of when given, otherwise now in the local
// zone, so the dates the tool writes (ids, created, updated, log lines) are the
// local calendar day. Timestamps (claimed_at) are written in UTC.
func clock(q Request) time.Time {
	if s := q.str("as_of"); s != "" {
		layout := "2006-01-02"
		if len(s) > 10 {
			layout = "2006-01-02T15:04"
		}
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return timeNow()
}

func today(t time.Time) string { return t.Format("2006-01-02") }

// logStamp is the stamp of a Log line (ADR 261101-gw7b): the writer's local
// time with its offset, "2006-01-02 15:04-07:00". A date-only as_of gives a
// date-only stamp, because the tool does not invent a time of day; an as_of
// with a time is read as UTC and so carries +00:00.
func logStamp(q Request, t time.Time) string {
	if len(q.str("as_of")) == 10 {
		return today(t)
	}
	return t.Format("2006-01-02 15:04-07:00")
}

// logAt parses the `at` of a log call, in the zone of the write when it has no
// offset, and renders the suffix " (at ...)": the time alone on the stamp's
// day, with the date otherwise. It refuses a time after the write.
func logAt(q Request, now time.Time) (string, *Error) {
	s := q.str("at")
	if s == "" {
		return "", nil
	}
	layout := "2006-01-02 15:04"
	if s[10] == 'T' {
		layout = "2006-01-02T15:04"
	}
	var at time.Time
	var err error
	if len(s) > 16 {
		at, err = time.Parse(layout+"-07:00", s)
	} else {
		at, err = time.ParseInLocation(layout, s, now.Location())
	}
	if err != nil {
		return "", invalid("at", "at is not a valid time: %s", s)
	}
	dayOnly := len(q.str("as_of")) == 10 // the clock has no time of day to compare
	if (!dayOnly && at.After(now)) || (dayOnly && today(at) > today(now)) {
		return "", invalid("at", "at is after the time of the write")
	}
	if today(at) == today(now) {
		return " (at " + at.Format("15:04-07:00") + ")", nil
	}
	return " (at " + at.Format("2006-01-02 15:04-07:00") + ")", nil
}

func refused(reason, format string, a ...any) *Error {
	e := invalid("", format, a...)
	e.Code = "refused"
	e.Details = O("reason", reason)
	return e
}

func refusedRule(rule, format string, a ...any) *Error {
	e := refused("rule", format, a...)
	e.Details = O("reason", "rule", "rule", rule)
	return e
}

// unsupportedBody is the refusal for a body write the section tree will not make:
// the body holds a construct akmctl cannot edit around, and nothing is written.
func unsupportedBody(e *edit, err error) *Error {
	var u *mdbody.UnsupportedError
	if !errors.As(err, &u) {
		return unreadable(err)
	}
	// the line is the file's, as rule B1 reports it
	line := u.Line + strings.Count(string(e.rec.Text[:len(e.rec.Text)-len(e.rec.Body)]), "\n")
	re := refused("unsupported_markdown", "%s has %s at line %d, which akmctl cannot edit around; edit the file by hand", e.rec.ID, u.Construct, line)
	re.Details = O("reason", "unsupported_markdown", "construct", u.Construct, "line", line, "structural", u.Structural)
	return re
}

func unreadable(err error) *Error { return &Error{Code: "unreadable", Message: err.Error()} }

// edit is a locked, compare-and-swap rewrite of one record file.
type edit struct {
	repo *record.Repo
	rec  *record.Record
	doc  *frontmatter.Doc
	text string // body after the frontmatter, editable
	eol  string
	old  string // file text when the lock was taken
	dry  bool   // dry_run: nothing is written and the lock is not taken
	diff Obj    // the would-be change, set by commit when dry
}

// openEdit locks the repository, re-reads the record from disk and checks the
// etag. The caller must call the returned release function.
func openEdit(repo *record.Repo, q Request) (*edit, func(), *Error) {
	dry := q.boolean("dry_run")
	unlock := func() {}
	var rec *record.Record
	var err *Error
	if dry {
		// A dry run takes no lock, so it can read while a writer replaces the file.
		// Where a replace is not atomic for readers (Windows) the record can be
		// missing for a moment; read again a few times before saying not_found.
		for i := 0; ; i++ {
			if lerr := repo.Reload(); lerr != nil {
				return nil, nil, unreadable(lerr)
			}
			if rec, err = lookup(repo, q.str("id")); err == nil || err.Code != "not_found" || i == 40 {
				break
			}
			time.Sleep(5 * time.Millisecond)
		}
		if err != nil {
			return nil, nil, err
		}
	} else {
		var lerr error
		if unlock, lerr = fsx.LockRepo(repo.Root, repo.MetaDir()); lerr != nil {
			return nil, nil, unreadable(lerr)
		}
		if lerr := repo.Reload(); lerr != nil {
			unlock()
			return nil, nil, unreadable(lerr)
		}
		if rec, err = lookup(repo, q.str("id")); err != nil {
			unlock()
			return nil, nil, err
		}
	}
	if want := q.str("if_match"); want != "" && want != rec.Etag {
		unlock()
		return nil, nil, &Error{"stale", "the record changed since it was read", O("etag", rec.Etag)}
	}
	doc, perr := frontmatter.Parse(rec.Text)
	if perr != nil {
		unlock()
		return nil, nil, unreadable(perr)
	}
	eol := "\n"
	if strings.Contains(string(rec.Text), "\r\n") {
		eol = "\r\n"
	}
	return &edit{repo: repo, rec: rec, doc: doc, text: string(doc.Body()), eol: eol, old: string(rec.Text), dry: dry}, unlock, nil
}

func (e *edit) bytes() []byte {
	b := e.doc.Bytes()
	head := b[:len(b)-len(e.doc.Body())]
	return append(append([]byte(nil), head...), e.text...)
}

// commit writes the edited file (to dest when the record moves) and returns the
// record read back from disk.
func (e *edit) commit(dest string) (*record.Record, *Error) {
	data := e.bytes()
	if e.dry {
		return e.preview(data, dest)
	}
	full := filepath.Join(e.repo.Root, filepath.FromSlash(e.rec.Path))
	if string(data) != e.old {
		if err := fsx.WriteAtomic(full, data); err != nil {
			return nil, unreadable(err)
		}
	}
	if dest != "" && dest != e.rec.Path {
		to := filepath.Join(e.repo.Root, filepath.FromSlash(dest))
		if err := os.MkdirAll(filepath.Dir(to), 0o755); err == nil {
			err = fsx.Rename(full, to)
			if err != nil {
				fsx.WriteAtomic(full, []byte(e.old))
				return nil, unreadable(err)
			}
		} else {
			fsx.WriteAtomic(full, []byte(e.old))
			return nil, unreadable(err)
		}
	}
	if err := e.repo.Reload(); err != nil {
		return nil, unreadable(err)
	}
	rec, aerr := lookup(e.repo, e.rec.ID)
	return rec, aerr
}

// preview is commit for a dry run: the record as it would be after the write
// (with the current etag) and the diff, with nothing written.
func (e *edit) preview(data []byte, dest string) (*record.Record, *Error) {
	rel := e.rec.Path
	if dest != "" {
		rel = dest
	}
	rec, err := record.FromText(rel, e.rec.Type, e.rec.ID, data)
	if err != nil {
		return nil, unreadable(err)
	}
	rec.Etag = e.rec.Etag
	e.diff = Diff(e.rec.Path, rel, e.old, string(data))
	return rec, nil
}

// finish adds the diff to the response of a dry run.
func (e *edit) finish(out Obj) Obj {
	if !e.dry {
		return out
	}
	if e.diff == nil {
		e.diff = Diff(e.rec.Path, e.rec.Path, e.old, e.old)
	}
	return append(out, KV{"diff", e.diff})
}

func refString(repo *record.Repo, id string) string {
	if v, ok := repo.Layout["reference_style"]; ok && v.Str == "wikilink" {
		if recs := repo.ByID(id); len(recs) == 1 {
			base := strings.TrimSuffix(path.Base(recs[0].Path), ".md")
			return "[[" + base + "|" + id + "]]"
		}
	}
	return id
}

// ---- new ----

func opNew(repo *record.Repo, q Request) (Obj, *Error) {
	kind := q.str("kind")
	slug := q.str("slug")
	if !reSlug.MatchString(slug) || len(slug) > 60 {
		return nil, invalid("slug", "slug must be lowercase words joined by hyphens, at most 60 characters")
	}
	now := clock(q)
	fields, _ := q["fields"].(map[string]any)
	title := q.str("title")
	if title == "" {
		title = strings.ReplaceAll(slug, "-", " ")
	}
	if strings.ContainsAny(title, "\r\n") {
		return nil, invalid("title", "title must be one line")
	}
	doc, body, aerr := newRecordText(repo, kind, title, fields, now)
	if aerr != nil {
		return nil, aerr
	}
	unlock, err := fsx.LockRepo(repo.Root, repo.MetaDir())
	if err != nil {
		return nil, unreadable(err)
	}
	defer unlock()
	if err := repo.Reload(); err != nil {
		return nil, unreadable(err)
	}
	id, err := mintID(repo, now)
	if err != nil {
		return nil, unreadable(err)
	}
	var similar []Obj
	if kind == record.Work {
		similar = similarItems(repo, slug, title)
	}
	rel := repo.WorkDir() + "/" + id + "-" + slug + ".md"
	if kind == record.ADR {
		rel = repo.ADRDir() + "/" + id + "-" + slug + ".md"
	}
	text := strings.Replace(doc, "{{id}}", id, 1) + body
	if err := fsx.WriteAtomic(filepath.Join(repo.Root, filepath.FromSlash(rel)), []byte(text)); err != nil {
		return nil, unreadable(err)
	}
	if err := repo.Reload(); err != nil {
		return nil, unreadable(err)
	}
	rec, aerr := lookup(repo, id)
	if aerr != nil {
		return nil, aerr
	}
	out := O("record", recordObj(repo, rec, nil, true), "path", rec.Path)
	if len(similar) > 0 {
		out = append(out, KV{"similar", similar})
	}
	out = append(out, KV{"commit", commitObj("file", id, title, rec.Path)})
	return out, nil
}

func mintID(repo *record.Repo, now time.Time) (string, error) {
	for i := 0; i < 50; i++ {
		var b [4]byte
		if _, err := rand.Read(b[:]); err != nil {
			return "", err
		}
		suffix := make([]byte, 4)
		for j := range b {
			suffix[j] = idAlphabet[b[j]&31]
		}
		id := now.Format("060102") + "-" + string(suffix)
		if len(repo.ByID(id)) == 0 {
			return id, nil
		}
	}
	return "", os.ErrExist
}

// newRecordText builds the frontmatter (with {{id}} as the id placeholder) and
// the body of a new record.
func newRecordText(repo *record.Repo, kind, title string, f map[string]any, now time.Time) (string, string, *Error) {
	allowed := map[string]bool{"summary": true}
	if kind == record.Work {
		for _, k := range []string{"kind", "parent", "priority", "blocked_by", "decisions"} {
			allowed[k] = true
		}
	} else {
		allowed["scope"] = true
	}
	for k := range f {
		if !allowed[k] {
			return "", "", invalid("fields", "fields: %q cannot be set when creating a %s", k, kind)
		}
	}
	fm, _ := frontmatter.Parse([]byte("---\nid: x\n---\n"))
	set := func(k string, v frontmatter.Value) *Error {
		if err := fm.Set(k, v); err != nil {
			return invalid("fields", "fields.%s: %v", k, err)
		}
		return nil
	}
	set("id", frontmatter.NewString("{{id}}"))
	set("title", frontmatter.NewString(title))
	summary, hasSummary := f["summary"]
	if hasSummary {
		s, ok := summary.(string)
		if !ok || len([]rune(s)) > 160 || strings.ContainsAny(s, "\r\n") {
			return "", "", invalid("fields", "fields.summary must be one line of at most 160 characters")
		}
		set("summary", frontmatter.NewString(s))
	} else if kind == record.ADR {
		s := title
		if r := []rune(s); len(r) > 160 {
			s = string(r[:160])
		}
		set("summary", frontmatter.NewString(s))
	}
	day := today(now)
	if kind == record.Work {
		set("status", frontmatter.NewString("raw"))
		var parentID string
		var blockers []frontmatter.Value
		for _, k := range []string{"kind", "parent", "priority", "blocked_by", "decisions"} {
			v, ok := f[k]
			if !ok {
				continue
			}
			fv, err := fieldValue(repo, record.Work, k, v)
			if err != nil {
				return "", "", err
			}
			if k == "parent" && fv.Kind == frontmatter.String {
				if rerr := parentEdgeRule(repo, "", "raw", record.RefID(fv.Str)); rerr != nil {
					return "", "", rerr
				}
				parentID = record.RefID(fv.Str)
			}
			if k == "blocked_by" && fv.Kind == frontmatter.List {
				blockers = fv.Items
			}
			if err := set(k, fv); err != nil {
				return "", "", err
			}
		}
		if rerr := newItemCycle(repo, parentID, blockers); rerr != nil {
			return "", "", rerr
		}
		set("created", frontmatter.NewString(day))
		set("updated", frontmatter.NewString(day))
		return string(fm.Bytes()), "\n# " + title + "\n\n## Goal\n\n## Log\n\n- " + day + ": filed.\n", nil
	}
	set("status", frontmatter.NewString("proposed"))
	set("date", frontmatter.NewString(day))
	if v, ok := f["scope"]; ok {
		fv, err := fieldValue(repo, record.ADR, "scope", v)
		if err != nil {
			return "", "", err
		}
		if err := set("scope", fv); err != nil {
			return "", "", err
		}
	}
	return string(fm.Bytes()), "\n# " + title + "\n\n## Context\n\n## Considered Options\n\n## Decision\n\n## Consequences\n", nil
}

// fieldValue converts a JSON value to a frontmatter value for a field of a
// record type, checking its type and enum.
func fieldValue(repo *record.Repo, typ, key string, v any) (frontmatter.Value, *Error) {
	bad := func(format string, a ...any) (frontmatter.Value, *Error) {
		return frontmatter.Value{}, invalid("fields", "%s: "+format, append([]any{key}, a...)...)
	}
	if v == nil {
		return frontmatter.NewNull(), nil
	}
	switch key {
	case "title", "summary", "authored_by", "claimed_by":
		s, ok := v.(string)
		if !ok || strings.ContainsAny(s, "\r\n") || (key != "authored_by" && strings.TrimSpace(s) == "") {
			return bad("must be one line of text")
		}
		if key == "summary" && len([]rune(s)) > 160 {
			return bad("at most 160 characters")
		}
		return frontmatter.NewString(s), nil
	case "kind":
		s, _ := v.(string)
		if !contains([]string{"feature", "bug", "chore", "spike"}, s) {
			return bad("must be feature, bug, chore or spike")
		}
		return frontmatter.NewString(s), nil
	case "status":
		s, _ := v.(string)
		set := workStatus
		if typ == record.ADR {
			set = []string{"proposed", "accepted", "rejected", "deprecated", "superseded"}
		}
		if !contains(set, s) {
			return bad("must be one of %s", strings.Join(set, ", "))
		}
		return frontmatter.NewString(s), nil
	case "cycle":
		s, ok := v.(string)
		if !ok || !reCycle.MatchString(s) {
			return bad("must be a label of 1 to 64 characters: letters, digits, '.', '_' or '-', starting with a letter or digit")
		}
		return frontmatter.NewString(s), nil
	case "keep":
		b, ok := v.(bool)
		if !ok {
			return bad("must be true or false")
		}
		return frontmatter.NewBool(b), nil
	case "priority":
		n, ok := asInt(v)
		if !ok {
			return bad("must be an integer")
		}
		return frontmatter.NewInt(int64(n)), nil
	case "parent", "supersedes", "superseded_by":
		s, ok := v.(string)
		if !ok || !record.ValidID(s) {
			return bad("must be a record id")
		}
		t := repo.ByID(s)
		want := record.Work
		if key != "parent" {
			want = record.ADR
		}
		if len(t) == 0 || t[0].Type != want {
			return frontmatter.Value{}, &Error{Code: "not_found", Message: key + " " + s + " does not resolve to a " + want + " record"}
		}
		return frontmatter.NewString(refString(repo, s)), nil
	case "blocked_by", "decisions", "scope":
		l, ok := asStrings(v)
		if !ok {
			return bad("must be a list of strings")
		}
		items := make([]frontmatter.Value, 0, len(l))
		for _, x := range l {
			if strings.ContainsAny(x, "\r\n") || strings.TrimSpace(x) == "" {
				return bad("entries must be one non-empty line")
			}
			if key == "decisions" {
				if !record.ValidID(x) {
					return bad("entries must be ADR ids")
				}
				t := repo.ByID(x)
				if len(t) == 0 || t[0].Type != record.ADR {
					return frontmatter.Value{}, &Error{Code: "not_found", Message: "decisions entry " + x + " does not resolve to an ADR"}
				}
				x = refString(repo, x)
			} else if key == "blocked_by" && record.ValidID(x) {
				x = refString(repo, x)
			}
			items = append(items, frontmatter.NewString(x))
		}
		return frontmatter.NewList(items...), nil
	}
	return bad("is not a field that can be set")
}

// blockedByRule refuses a blocked_by entry that would close a dependency cycle
// (WK18): the item would wait for target while target already waits, directly or
// through its children and its own blockers, for the item.
func blockedByRule(repo *record.Repo, id string, entries []frontmatter.Value, had []frontmatter.Value) *Error {
	known := map[string]bool{}
	for _, it := range had {
		known[record.RefID(it.Str)] = true
	}
	var g *depgraph.Graph
	for _, it := range entries {
		target := record.RefID(it.Str)
		if it.Kind != frontmatter.String || !record.ValidID(target) || known[target] {
			continue
		}
		if t := repo.ByID(target); len(t) == 0 || t[0].Type != record.Work {
			continue
		}
		if target == id {
			return refusedRule("WK18", "%s cannot be blocked by itself", id)
		}
		if g == nil {
			g = depgraph.Build(repo)
		}
		if path := g.Path(target, id); path != nil {
			return refusedRule("WK18", "blocked_by %s would make a dependency cycle: %s -> %s", target, id, strings.Join(path, " -> "))
		}
	}
	return nil
}

// newItemCycle refuses (WK18) a new item whose parent already waits, directly or
// through its own blockers and children, for something the new item would wait
// for: the parent waits for the new item (its child), which would wait for the
// blocker, which waits for the parent.
func newItemCycle(repo *record.Repo, parentID string, blockers []frontmatter.Value) *Error {
	if parentID == "" {
		return nil
	}
	var g *depgraph.Graph
	for _, it := range blockers {
		target := record.RefID(it.Str)
		if it.Kind != frontmatter.String || !record.ValidID(target) {
			continue
		}
		if t := repo.ByID(target); len(t) == 0 || t[0].Type != record.Work {
			continue
		}
		if target == parentID {
			return refusedRule("WK18", "blocked_by %s would make a dependency cycle: new item -> %s -> new item", target, parentID)
		}
		if g == nil {
			g = depgraph.Build(repo)
		}
		if path := g.Path(target, parentID); path != nil {
			return refusedRule("WK18", "blocked_by %s would make a dependency cycle: new item -> %s -> new item", target, strings.Join(path, " -> "))
		}
	}
	return nil
}

// parentEdgeRule refuses a new parent edge from the item id (status st; id is
// empty for an item being filed) to parentID when the edge would close a loop
// (WK4) or put an open item under a closed parent (WK8). The parent chain is
// walked upwards: with a single parent per item, the edge closes a loop exactly
// when the item is on the new parent's chain. A loop already in the chain that
// does not pass through the item is not this edge's doing and ends the walk.
func parentEdgeRule(repo *record.Repo, id, st, parentID string) *Error {
	if id != "" {
		chain := []string{id}
		seen := map[string]bool{}
		for cur := parentID; cur != "" && !seen[cur]; {
			seen[cur] = true
			chain = append(chain, cur)
			if cur == id {
				return refusedRule("WK4", "parent %s would make a loop: %s", parentID, strings.Join(chain, " -> "))
			}
			t := repo.ByID(cur)
			if len(t) == 0 || t[0].Type != record.Work {
				break
			}
			cur = record.RefID(t[0].Str("parent"))
		}
	}
	if id != "" {
		g := depgraph.Build(repo)
		if old := record.RefID(repo.ByID(id)[0].Str("parent")); old != "" {
			g = g.Without(old, id)
		}
		if path := g.Path(id, parentID); path != nil {
			return refusedRule("WK18", "parent %s would make a dependency cycle: %s -> %s", parentID, strings.Join(path, " -> "), id)
		}
	}
	if t := repo.ByID(parentID); len(t) > 0 && (t[0].Status() == "done" || t[0].Status() == "dropped") && st != "done" && st != "dropped" {
		who := "a new item"
		if id != "" {
			who = id
		}
		return refusedRule("WK8", "%s is %s, so %s cannot be its child while open", parentID, t[0].Status(), who)
	}
	return nil
}

// ---- set ----

var setFields = map[string][]string{
	record.Work: {"title", "summary", "status", "kind", "parent", "blocked_by", "priority", "cycle", "keep", "authored_by", "decisions"},
	record.ADR:  {"title", "summary", "status", "scope", "supersedes", "superseded_by", "authored_by"},
}

func opSet(repo *record.Repo, q Request) (Obj, *Error) {
	changes, _ := q["changes"].(map[string]any)
	if len(changes) == 0 {
		return nil, invalid("changes", "changes must name at least one field")
	}
	e, unlock, err := openEdit(repo, q)
	if err != nil {
		return nil, err
	}
	defer unlock()
	typ := e.rec.Type
	keys := make([]string, 0, len(changes))
	for k := range changes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	reason, cerr0 := correctionReason(q)
	if cerr0 != nil {
		return nil, cerr0
	}
	closed := typ == record.Work && repo.IsDone(e.rec.Path)
	if closed {
		for _, k := range keys {
			switch {
			case k == "keep" || k == "status" || !contains(setFields[typ], k): // status is WK9's to refuse; unknown fields are invalid below
			case reason == "":
				return nil, closedRefusal(e.rec.ID, k)
			case k != "title" && k != "summary":
				return nil, &Error{"refused", fmt.Sprintf("%s is closed; %s is never changed on a closed item, even with a correction (only title and summary are)", e.rec.ID, k), O("reason", "closed_field", "field", k)}
			}
		}
	}
	oldStatus := e.rec.Status()
	var changed []string
	for _, k := range keys {
		if !contains(setFields[typ], k) {
			return nil, invalid("changes", "changes: %q cannot be set on a %s (id, created, updated, closed, date and claim fields are managed by the tool)", k, typ)
		}
		v := changes[k]
		if v == nil && (k == "title" || k == "status" || (k == "summary" && typ == record.ADR)) {
			return nil, invalid("changes", "changes: %s cannot be removed", k)
		}
		if k == "status" && typ == record.Work {
			if s, _ := v.(string); s == "done" || s == "dropped" {
				return nil, invalid("changes", "changes: close the item with close to make it %s", s)
			}
		}
		fv, ferr := fieldValue(repo, typ, k, v)
		if ferr != nil {
			return nil, ferr
		}
		if k == "parent" && fv.Kind == frontmatter.String && record.RefID(fv.Str) == e.rec.ID {
			return nil, refusedRule("WK4", "an item cannot be its own parent")
		}
		cur, has := e.rec.Fields[k]
		if k == "blocked_by" && typ == record.Work && fv.Kind == frontmatter.List && !(has && sameValue(k, cur, fv)) {
			if rerr := blockedByRule(repo, e.rec.ID, fv.Items, cur.Items); rerr != nil {
				return nil, rerr
			}
		}
		if k == "parent" && fv.Kind == frontmatter.String && !(has && sameValue(k, cur, fv)) {
			if rerr := parentEdgeRule(repo, e.rec.ID, e.rec.Status(), record.RefID(fv.Str)); rerr != nil {
				return nil, rerr
			}
		}
		if fv.Kind == frontmatter.Null {
			if has {
				if _, derr := e.doc.Delete(k); derr != nil {
					return nil, unreadable(derr)
				}
				changed = append(changed, k)
			}
			continue
		}
		if has && sameValue(k, cur, fv) {
			continue
		}
		if serr := e.doc.Set(k, fv); serr != nil {
			return nil, invalid("changes", "changes.%s: %v", k, serr)
		}
		changed = append(changed, k)
	}
	if len(changed) == 0 {
		return e.finish(O("record", recordObj(repo, e.rec, nil, true), "changed", []string{})), nil
	}
	// status rules
	newStatus := oldStatus
	if v, ok := e.doc.Get("status"); ok {
		newStatus = v.Str
	}
	if typ == record.Work {
		by, _ := e.doc.Get("blocked_by")
		if newStatus == "blocked" && len(by.Items) == 0 {
			return nil, refusedRule("WK5", "status blocked needs a non-empty blocked_by")
		}
		finished := repo.IsDone(e.rec.Path)
		if finished && newStatus != "done" && newStatus != "dropped" {
			return nil, refusedRule("WK9", "a closed item stays in %s/ with status done or dropped", repo.DoneDir())
		}
	} else if newStatus == "superseded" {
		if v, ok := e.doc.Get("superseded_by"); !ok || v.Kind != frontmatter.String {
			return nil, refusedRule("A3", "a superseded ADR needs superseded_by")
		}
	}
	if closed && reason != "" {
		var fixed []string
		for _, k := range changed {
			if k != "keep" {
				fixed = append(fixed, k)
			}
		}
		if len(fixed) > 0 {
			b, lerr := logInto(e, q, []byte(e.text), fmt.Sprintf("corrected %s: %s", strings.Join(fixed, " and "), reason))
			if lerr != nil {
				return nil, lerr
			}
			e.text = string(b)
		}
	}
	now := clock(q)
	if typ == record.Work {
		if serr := e.doc.Set("updated", frontmatter.NewString(today(now))); serr != nil {
			return nil, unreadable(serr)
		}
	} else if newStatus != oldStatus {
		if serr := e.doc.Set("date", frontmatter.NewString(today(now))); serr != nil {
			return nil, unreadable(serr)
		}
	}
	rec, cerr := e.commit("")
	if cerr != nil {
		return nil, cerr
	}
	out := O("record", recordObj(repo, rec, nil, true), "changed", changed)
	if typ == record.Work {
		if step := statusStep(oldStatus, newStatus); step != "" {
			detail := map[string]string{"ready": "ready", "start": "started", "unblock": "unblocked"}[step]
			if step == "block" {
				detail = "blocked on " + firstBlocker(e.doc)
			}
			out = append(out, KV{"commit", commitObj(step, rec.ID, detail, rec.Path)})
		}
	}
	return e.finish(out), nil
}

// firstBlocker is the first blocked_by entry of the document, as written.
func firstBlocker(doc *frontmatter.Doc) string {
	if v, ok := doc.Get("blocked_by"); ok && len(v.Items) > 0 {
		return record.RefID(v.Items[0].Str)
	}
	return "its blocker"
}

// sameValue compares a stored value with a new one, reading references as ids.
func sameValue(key string, cur, nv frontmatter.Value) bool {
	return jsonEqual(record.JSONValue(key, cur), record.JSONValue(key, nv))
}

func jsonEqual(a, b any) bool {
	x, _ := Marshal(a)
	y, _ := Marshal(b)
	return string(x) == string(y)
}

// ---- close ----

func opClose(repo *record.Repo, q Request) (Obj, *Error) {
	e, unlock, err := openEdit(repo, q)
	if err != nil {
		return nil, err
	}
	defer unlock()
	if e.rec.Type != record.Work {
		return nil, invalid("id", "only work items can be closed; set an ADR's status instead")
	}
	if st := e.rec.Status(); st == "done" || st == "dropped" {
		return nil, refused("already_closed", "%s is already %s", e.rec.ID, st)
	}
	status := q.str("status")
	if o := q.str("outcome"); strings.TrimSpace(o) != "" {
		var terr error
		if e.text, terr = record.EditSection(e.text, "Outcome", o, false, e.eol); terr != nil {
			return nil, unsupportedBody(e, terr)
		}
	}
	_, text := record.Sections(e.text)
	if strings.TrimSpace(text["outcome"]) == "" {
		return nil, refused("empty_outcome", "%s has no ## Outcome; write it before closing", e.rec.ID)
	}
	if status == "done" || status == "dropped" {
		for _, ch := range repo.Records {
			if ch.Type == record.Work && record.RefID(ch.Str("parent")) == e.rec.ID && ch.Status() != "done" && ch.Status() != "dropped" {
				return nil, &Error{"refused", "child " + ch.ID + " is not done or dropped", O("reason", "open_children", "rule", "WK8")}
			}
		}
	}
	if status == "done" {
		if rerr := uncheckedGoal(e); rerr != nil {
			return nil, rerr
		}
	}
	if serr := e.doc.Set("status", frontmatter.NewString(status)); serr != nil {
		return nil, unreadable(serr)
	}
	for _, k := range []string{"claimed_by", "claimed_at"} {
		e.doc.Delete(k)
	}
	day := today(clock(q))
	if serr := e.doc.Set("updated", frontmatter.NewString(day)); serr != nil {
		return nil, unreadable(serr)
	}
	if serr := e.doc.SetAfter("closed", "updated", frontmatter.NewString(day)); serr != nil {
		return nil, unreadable(serr)
	}
	from := e.rec.Path
	to := repo.DoneDir() + "/" + path.Base(from)
	deps := dependents(repo, e.rec.ID)
	// The move breaks relative Markdown links in the file and to it: plan the
	// rewrite now, from the files as they are, and write it after the move.
	whole := string(e.bytes())
	headLen := len(whole) - len(e.text)
	moved, selfN := moveSelfLinks(repo, from, to, whole)
	if selfN > 0 {
		e.text = moved[headLen:] // the frontmatter is not scanned, so the head is unchanged
	}
	inbound := planMovedInbound(repo, from, to)
	rec, cerr := e.commit(to)
	if cerr != nil {
		return nil, cerr
	}
	var failed []Obj
	if !e.dry {
		failed = applyLinkEdits(repo, inbound)
		if len(inbound) > 0 {
			repo.Reload()
		}
	}
	out := O("record", recordObj(repo, rec, nil, true), "moved_from", from, "moved_to", rec.Path)
	if len(deps) > 0 {
		out = append(out, KV{"dependents", deps})
	}
	if rep := linkReport(rec.Path, selfN, inbound); len(rep) > 0 {
		out = append(out, KV{"links_rewritten", rep})
	}
	if len(failed) > 0 {
		out = append(out, KV{"link_errors", failed})
	}
	step, detail := "close", "done"
	if status == "dropped" {
		step, detail = "drop", "dropped"
	}
	paths := []string{from, rec.Path}
	for _, ed := range inbound {
		paths = append(paths, ed.path)
	}
	out = append(out, KV{"commit", commitObj(step, rec.ID, detail, paths...)})
	return e.finish(out), nil
}

// dependents lists the open work items whose blocked_by names id, in id order:
// the items that were waiting on the one being closed. The tool reports them
// and edits none (section 4.12 of the API spec).
func dependents(repo *record.Repo, id string) []Obj {
	var out []Obj
	for _, r := range repo.Records {
		if r.Type != record.Work || r.ID == id || !similarOpen[r.Status()] {
			continue
		}
		for _, it := range r.Fields["blocked_by"].Items {
			if it.Kind == frontmatter.String && record.RefID(it.Str) == id {
				out = append(out, O("id", r.ID, "status", r.Status(), "title", r.Str("title"), "path", r.Path))
				break
			}
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, _ := out[i].Get("id")
		b, _ := out[j].Get("id")
		return a.(string) < b.(string)
	})
	return out
}

// ---- claim and release ----

func claimsOn(repo *record.Repo) *Error {
	if !repo.Extension("claims") {
		return refused("extension_off", "the claims extension is not enabled in %s", repo.ManifestPath())
	}
	return nil
}

func claimStaleHours(repo *record.Repo) int {
	if v, ok := repo.Layout["claim_stale_hours"]; ok && v.Kind == frontmatter.Int && v.Int > 0 {
		return int(v.Int)
	}
	return 24
}

func claimStale(repo *record.Repo, claimedAt string, now time.Time) bool {
	hours := claimStaleHours(repo)
	t, err := time.Parse("2006-01-02T15:04", claimedAt)
	if err != nil {
		return true
	}
	return now.Sub(t) > time.Duration(hours)*time.Hour
}

func opClaim(repo *record.Repo, q Request) (Obj, *Error) {
	if err := claimsOn(repo); err != nil {
		return nil, err
	}
	e, unlock, err := openEdit(repo, q)
	if err != nil {
		return nil, err
	}
	defer unlock()
	if e.rec.Type != record.Work {
		return nil, invalid("id", "only work items can be claimed")
	}
	st := e.rec.Status()
	if st != "ready" && st != "active" {
		return nil, refused("status", "%s is %s; only ready or active items can be claimed", e.rec.ID, st)
	}
	// An active item needs a ## Next (WK6). Unlike set, claim can ask for it first:
	// the item is still ready and section writes it.
	if _, text := record.Sections(e.rec.Body); strings.TrimSpace(text["next"]) == "" {
		return nil, refusedRule("WK6", "%s has an empty ## Next; write it with section, then claim", e.rec.ID)
	}
	now := clock(q)
	who := q.str("as")
	prev := e.rec.Str("claimed_by")
	var previous any
	logLine := "claimed (" + who + ")."
	if prev != "" {
		previous = prev
		if prev != who && !claimStale(repo, e.rec.Str("claimed_at"), now) && !q.boolean("takeover") {
			return nil, refused("claimed", "%s is claimed by %s", e.rec.ID, prev)
		}
		if prev != who {
			logLine = "took over from " + prev + " (" + who + ")."
		} else if st == "active" {
			// the same holder already holds an active item: the claim is
			// refreshed (claimed_at) but not logged a second time
			logLine = ""
		}
	}
	set := func(k string, v frontmatter.Value) *Error {
		if serr := e.doc.Set(k, v); serr != nil {
			return unreadable(serr)
		}
		return nil
	}
	if st == "ready" {
		if serr := set("status", frontmatter.NewString("active")); serr != nil {
			return nil, serr
		}
	}
	for _, kv := range []struct{ k, v string }{{"claimed_by", who}, {"claimed_at", now.UTC().Format("2006-01-02T15:04")}, {"updated", today(now)}} {
		if serr := set(kv.k, frontmatter.NewString(kv.v)); serr != nil {
			return nil, serr
		}
	}
	if logLine != "" {
		var terr error
		if e.text, terr = record.AppendToSection(e.text, "log", "- "+logStamp(q, now)+": "+logLine, e.eol); terr != nil {
			return nil, unsupportedBody(e, terr)
		}
	}
	rec, cerr := e.commit("")
	if cerr != nil {
		return nil, cerr
	}
	return e.finish(O("record", recordObj(repo, rec, nil, true), "previous_holder", previous, "commit", commitObj("claim", rec.ID, "claimed by "+who, rec.Path))), nil
}

func opRelease(repo *record.Repo, q Request) (Obj, *Error) {
	if err := claimsOn(repo); err != nil {
		return nil, err
	}
	e, unlock, err := openEdit(repo, q)
	if err != nil {
		return nil, err
	}
	defer unlock()
	if e.rec.Type != record.Work {
		return nil, invalid("id", "only work items can be released")
	}
	prev := e.rec.Str("claimed_by")
	if prev == "" {
		return nil, refused("not_claimed", "%s has no claim", e.rec.ID)
	}
	now := clock(q)
	e.doc.Delete("claimed_by")
	e.doc.Delete("claimed_at")
	if serr := e.doc.Set("updated", frontmatter.NewString(today(now))); serr != nil {
		return nil, unreadable(serr)
	}
	var terr error
	if e.text, terr = record.AppendToSection(e.text, "log", "- "+logStamp(q, now)+": released ("+prev+").", e.eol); terr != nil {
		return nil, unsupportedBody(e, terr)
	}
	rec, cerr := e.commit("")
	if cerr != nil {
		return nil, cerr
	}
	return e.finish(O("record", recordObj(repo, rec, nil, true), "commit", commitObj("release", rec.ID, "released", rec.Path))), nil
}

// ---- section ----

func opSection(repo *record.Repo, q Request) (Obj, *Error) {
	name := strings.TrimSpace(q.str("name"))
	if name == "" || strings.ContainsAny(name, "\r\n") || strings.HasPrefix(name, "#") {
		return nil, invalid("name", "name must be the text of a ## heading: one line, no hashes")
	}
	mode := q.str("mode")
	if mode == "" {
		mode = "replace"
	}
	if mode == "replace" && strings.EqualFold(name, "log") {
		return nil, invalid("mode", "the Log section is append-only: use mode append")
	}
	e, unlock, err := openEdit(repo, q)
	if err != nil {
		return nil, err
	}
	defer unlock()
	reason, cerr0 := correctionReason(q)
	if cerr0 != nil {
		return nil, cerr0
	}
	closed := repo.IsDone(e.rec.Path)
	if closed {
		switch {
		case reason == "":
			return nil, closedRefusal(e.rec.ID, "the "+strings.ToLower(name)+" section")
		case strings.EqualFold(name, "log"):
			return nil, &Error{"refused", fmt.Sprintf("%s is closed; its Log is written by the tool, not edited, even with a correction", e.rec.ID), O("reason", "closed_field", "field", "log")}
		}
	}
	text, terr := record.EditSection(e.text, name, q.str("text"), mode == "append", e.eol)
	if terr != nil {
		return nil, unsupportedBody(e, terr)
	}
	if closed && text != e.text {
		b, lerr := logInto(e, q, []byte(text), fmt.Sprintf("corrected %s: %s", strings.ToLower(name), reason))
		if lerr != nil {
			return nil, lerr
		}
		text = string(b)
	}
	changed := text != e.text
	rec := e.rec
	if changed {
		e.text = text
		if e.rec.Type == record.Work {
			if serr := e.doc.Set("updated", frontmatter.NewString(today(clock(q)))); serr != nil {
				return nil, unreadable(serr)
			}
		}
		var cerr *Error
		if rec, cerr = e.commit(""); cerr != nil {
			return nil, cerr
		}
	}
	_, secs := record.Sections(rec.Body)
	return e.finish(O("record", recordObj(repo, rec, nil, true), "section", O("name", strings.ToLower(name), "text", secs[strings.ToLower(name)]), "changed", changed)), nil
}

// ---- tick and log ----

// sectionOut is the section as the write operations return it: name in lower
// case, the text as `show` returns it, and the etag a client sends back as
// section_etag.
func sectionOut(body, name string) Obj {
	_, text := record.Sections(body)
	key := strings.ToLower(name)
	etag := ""
	if tr, err := mdbody.Split([]byte(body)); err == nil {
		if sec, ok := tr.Section(name); ok {
			etag = tr.Etag(sec)
		}
	}
	return O("name", key, "text", text[key], "etag", etag)
}

// staleSection is the answer to a section_etag that no longer matches.
func staleSection(tr *mdbody.Tree, sec mdbody.Section) *Error {
	return &Error{"stale", "the " + strings.ToLower(sec.Name) + " section changed since it was read", O("section_etag", tr.Etag(sec))}
}

// uncheckedGoal refuses `close done` while the Goal holds an empty box: the
// criteria are not met, or the box is stale and the Goal should say so. A body
// the tree cannot read is not checked (the check needs the Goal's lines).
func uncheckedGoal(e *edit) *Error {
	tr, err := mdbody.Split([]byte(e.text))
	if err != nil {
		return nil
	}
	items, err := tr.Unchecked("goal")
	if err != nil || len(items) == 0 {
		return nil
	}
	var rows []any
	var msg []string
	for _, it := range items {
		row := O("line", it.Line, "text", it.Text)
		if it.Pos > 0 {
			row = append(row, KV{"position", it.Pos})
		}
		rows = append(rows, row)
		msg = append(msg, fmt.Sprintf("line %d %q", it.Line, it.Text))
	}
	return &Error{"refused", fmt.Sprintf("%s has unchecked boxes in its Goal: %s; tick each with evidence (tick --evidence), remove one whose requirement is gone (remove), change the Goal with section, or close it dropped", e.rec.ID, strings.Join(msg, ", ")), O("reason", "unchecked_goal", "items", rows)}
}

func opTick(repo *record.Repo, q Request) (Obj, *Error) {
	name := strings.TrimSpace(q.str("section"))
	if name == "" {
		name = "goal"
	}
	done := true
	if v, ok := q["done"].(bool); ok {
		done = v
	}
	pos, _ := asInt(q["item"])
	e, unlock, err := openEdit(repo, q)
	if err != nil {
		return nil, err
	}
	defer unlock()
	if repo.IsDone(e.rec.Path) {
		return nil, refused("closed", "%s is closed; its sections are not edited", e.rec.ID)
	}
	tr, terr := mdbody.Split([]byte(e.text))
	if terr != nil {
		return nil, unsupportedBody(e, terr)
	}
	sec, ok := tr.Section(name)
	if !ok {
		return nil, &Error{Code: "not_found", Message: "no section " + strings.ToLower(name) + " in " + e.rec.ID}
	}
	if want := q.str("section_etag"); want != "" && want != tr.Etag(sec) {
		return nil, staleSection(tr, sec)
	}
	out, terr := tr.SetChecked(name, pos, done, "")
	switch {
	case errors.Is(terr, mdbody.ErrNoItem):
		return nil, &Error{Code: "not_found", Message: fmt.Sprintf("%s has no item %d in its %s checklist", e.rec.ID, pos, strings.ToLower(name))}
	case terr != nil:
		return nil, unsupportedBody(e, terr)
	}
	if ev := strings.Join(strings.Fields(q.str("evidence")), " "); q.str("evidence") != "" {
		if ev == "" || strings.ContainsAny(q.str("evidence"), "\r\n") {
			return nil, invalid("evidence", "evidence must be one non-empty line")
		}
		items, _ := tr.Checklist(name)
		verb := "ticked"
		if !done {
			verb = "unticked"
		}
		line := fmt.Sprintf("%s \"%s\": %s", verb, items[pos-1].Text, ev)
		var lerr *Error
		if out, lerr = logInto(e, q, out, line); lerr != nil {
			return nil, lerr
		}
	}
	return finishBodyEdit(repo, q, e, string(out), name, func(rec *record.Record) Obj {
		items, _ := mustSplit(rec.Body).Checklist(name)
		it := items[pos-1]
		return O("position", pos, "done", it.Done, "text", it.Text)
	}, "item")
}

// correctionReason reads the correction parameter: one non-empty line, or ""
// when the call has none (ADR 261101-dxcn).
func correctionReason(q Request) (string, *Error) {
	raw := q.str("correction")
	if raw == "" {
		return "", nil
	}
	r := strings.Join(strings.Fields(raw), " ")
	if r == "" || strings.ContainsAny(raw, "\r\n") {
		return "", invalid("correction", "correction must be one non-empty line")
	}
	return r, nil
}

// closedRefusal is the answer to an edit of a closed work item that carries no
// correction.
func closedRefusal(id, what string) *Error {
	return &Error{"refused", fmt.Sprintf("%s is closed; %s is not edited without a correction: pass correction, a one-line reason, which is written to its Log", id, what), O("reason", "closed")}
}

// logInto appends line to the Log of body (creating the section when it has
// none) with the stamp of this write.
func logInto(e *edit, q Request, body []byte, line string) ([]byte, *Error) {
	stamp := logStamp(q, clock(q))
	otr, serr := mdbody.Split(body)
	if serr != nil {
		return nil, unsupportedBody(e, serr)
	}
	if _, has := otr.Section("log"); has {
		b, aerr := otr.WithEOL(e.eol).AppendLog(stamp, line)
		if aerr != nil {
			return nil, unsupportedBody(e, aerr)
		}
		return b, nil
	}
	s, aerr := record.AppendToSection(string(body), "log", "- "+stamp+": "+line, e.eol)
	if aerr != nil {
		return nil, unsupportedBody(e, aerr)
	}
	return []byte(s), nil
}

// opRemove deletes one unchecked box of a checklist and logs it. The record
// loses that line: the tool does not track git, so confirm says the caller
// accepts that.
func opRemove(repo *record.Repo, q Request) (Obj, *Error) {
	name := strings.TrimSpace(q.str("section"))
	if name == "" {
		name = "goal"
	}
	reason := strings.Join(strings.Fields(q.str("reason")), " ")
	if reason == "" || strings.ContainsAny(q.str("reason"), "\r\n") {
		return nil, invalid("reason", "reason must be one non-empty line")
	}
	pos, _ := asInt(q["item"])
	e, unlock, err := openEdit(repo, q)
	if err != nil {
		return nil, err
	}
	defer unlock()
	if repo.IsDone(e.rec.Path) {
		return nil, refused("closed", "%s is closed; its sections are not edited", e.rec.ID)
	}
	tr, terr := mdbody.Split([]byte(e.text))
	if terr != nil {
		return nil, unsupportedBody(e, terr)
	}
	sec, ok := tr.Section(name)
	if !ok {
		return nil, &Error{Code: "not_found", Message: "no section " + strings.ToLower(name) + " in " + e.rec.ID}
	}
	if want := q.str("section_etag"); want != "" && want != tr.Etag(sec) {
		return nil, staleSection(tr, sec)
	}
	out, it, terr := tr.RemoveItem(name, pos, "")
	switch {
	case errors.Is(terr, mdbody.ErrNoItem):
		return nil, &Error{Code: "not_found", Message: fmt.Sprintf("%s has no item %d in its %s checklist", e.rec.ID, pos, strings.ToLower(name))}
	case errors.Is(terr, mdbody.ErrChecked):
		return nil, &Error{"refused", fmt.Sprintf("item %d (%q) is ticked, which records that the work was done; it is not removed. Say more with log, or change the Goal with section", pos, it.Text), O("reason", "ticked_box", "item", pos, "text", it.Text)}
	case terr != nil:
		return nil, unsupportedBody(e, terr)
	}
	if confirm, _ := q["confirm"].(bool); !confirm && !q.boolean("dry_run") {
		return nil, &Error{"refused", fmt.Sprintf("remove deletes item %d (%q) from %s and logs it. This loses history: the tool does not track git, so the line is recoverable only from a commit. Commit first if you want it back, then pass confirm to accept the loss", pos, it.Text, e.rec.ID), O("reason", "confirm_required", "item", pos, "text", it.Text)}
	}
	var lerr *Error
	if out, lerr = logInto(e, q, out, fmt.Sprintf("removed \"%s\": %s", it.Text, reason)); lerr != nil {
		return nil, lerr
	}
	return finishBodyEdit(repo, q, e, string(out), name, func(rec *record.Record) Obj {
		items, _ := mustSplit(rec.Body).Checklist(name)
		rows := make([]Obj, 0, len(items))
		for _, x := range items {
			rows = append(rows, O("position", x.Pos, "done", x.Done, "text", x.Text))
		}
		return O("position", pos, "text", it.Text, "items", rows)
	}, "removed")
}

func opLog(repo *record.Repo, q Request) (Obj, *Error) {
	text := strings.Join(strings.Fields(q.str("text")), " ")
	if text == "" || strings.ContainsAny(q.str("text"), "\r\n") {
		return nil, invalid("text", "text must be one non-empty line")
	}
	e, unlock, err := openEdit(repo, q)
	if err != nil {
		return nil, err
	}
	defer unlock()
	if repo.IsDone(e.rec.Path) {
		return nil, refused("closed", "%s is closed; its sections are not edited", e.rec.ID)
	}
	tr, terr := mdbody.Split([]byte(e.text))
	if terr != nil {
		return nil, unsupportedBody(e, terr)
	}
	now := clock(q)
	suffix, aerr0 := logAt(q, now)
	if aerr0 != nil {
		return nil, aerr0
	}
	text += suffix
	entry := "- " + logStamp(q, now) + ": " + text
	var out string
	if sec, ok := tr.Section("log"); ok {
		if want := q.str("section_etag"); want != "" && want != tr.Etag(sec) {
			return nil, staleSection(tr, sec)
		}
		b, aerr := tr.WithEOL(e.eol).AppendLog(logStamp(q, now), text)
		if aerr != nil {
			return nil, unsupportedBody(e, aerr)
		}
		out = string(b)
	} else {
		if q.str("section_etag") != "" {
			return nil, &Error{Code: "not_found", Message: "no section log in " + e.rec.ID}
		}
		var aerr error
		if out, aerr = record.AppendToSection(e.text, "log", entry, e.eol); aerr != nil {
			return nil, unsupportedBody(e, aerr)
		}
	}
	return finishBodyEdit(repo, q, e, out, "log", func(*record.Record) Obj { return O("text", entry) }, "entry")
}

func mustSplit(body string) *mdbody.Tree {
	tr, _ := mdbody.Split([]byte(body))
	return tr
}

// finishBodyEdit writes a new body for tick and log: it sets updated when the
// text changed, commits once and answers with the record, the section with its
// etag, `changed` and the operation's own member.
func finishBodyEdit(repo *record.Repo, q Request, e *edit, text, secName string, member func(*record.Record) Obj, memberName string) (Obj, *Error) {
	changed := text != e.text
	rec := e.rec
	if changed {
		e.text = text
		if e.rec.Type == record.Work {
			if serr := e.doc.Set("updated", frontmatter.NewString(today(clock(q)))); serr != nil {
				return nil, unreadable(serr)
			}
		}
		var cerr *Error
		if rec, cerr = e.commit(""); cerr != nil {
			return nil, cerr
		}
	}
	return e.finish(O("record", recordObj(repo, rec, nil, true), "section", sectionOut(rec.Body, secName), memberName, member(rec), "changed", changed)), nil
}
