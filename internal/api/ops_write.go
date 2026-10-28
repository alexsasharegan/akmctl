package api

import (
	"crypto/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/alexsasharegan/akmctl/internal/frontmatter"
	"github.com/alexsasharegan/akmctl/internal/fsx"
	"github.com/alexsasharegan/akmctl/internal/record"
)

const idAlphabet = "0123456789abcdefghjkmnpqrstvwxyz"

var reSlug = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// clock is the time of a write: as_of when given, otherwise now (UTC).
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
	return time.Now().UTC()
}

func today(t time.Time) string { return t.Format("2006-01-02") }

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
	rec, err := lookup(repo, q.str("id"))
	if err != nil {
		return nil, nil, err
	}
	dry := q.boolean("dry_run")
	unlock := func() {}
	if !dry {
		var lerr error
		if unlock, lerr = fsx.LockRepo(repo.Root); lerr != nil {
			return nil, nil, unreadable(lerr)
		}
	}
	ok := false
	defer func() {
		if !ok {
			unlock()
		}
	}()
	if lerr := repo.Reload(); lerr != nil {
		return nil, nil, unreadable(lerr)
	}
	rec, err = lookup(repo, q.str("id"))
	if err != nil {
		return nil, nil, err
	}
	if want := q.str("if_match"); want != "" && want != rec.Etag {
		return nil, nil, &Error{"stale", "the record changed since it was read", O("etag", rec.Etag)}
	}
	doc, perr := frontmatter.Parse(rec.Text)
	if perr != nil {
		return nil, nil, unreadable(perr)
	}
	eol := "\n"
	if strings.Contains(string(rec.Text), "\r\n") {
		eol = "\r\n"
	}
	ok = true
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
			err = os.Rename(full, to)
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
	unlock, err := fsx.LockRepo(repo.Root)
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
	rel := "docs/work/" + id + "-" + slug + ".md"
	if kind == record.ADR {
		rel = "docs/adr/" + id + "-" + slug + ".md"
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
	return O("record", recordObj(repo, rec, nil, true), "path", rec.Path), nil
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
		for _, k := range []string{"kind", "parent", "priority", "blocked_by", "decisions"} {
			v, ok := f[k]
			if !ok {
				continue
			}
			fv, err := fieldValue(repo, record.Work, k, v)
			if err != nil {
				return "", "", err
			}
			if err := set(k, fv); err != nil {
				return "", "", err
			}
		}
		set("created", frontmatter.NewString(day))
		set("updated", frontmatter.NewString(day))
		return string(fm.Bytes()), "\n# " + title + "\n\n## Goal\n\n## Log\n- " + day + ": filed.\n", nil
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

// ---- set ----

var setFields = map[string][]string{
	record.Work: {"title", "summary", "status", "kind", "parent", "blocked_by", "priority", "authored_by", "decisions"},
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
	oldStatus := e.rec.Status()
	var changed []string
	for _, k := range keys {
		if !contains(setFields[typ], k) {
			return nil, invalid("changes", "changes: %q cannot be set on a %s (id, created, updated, date and claim fields are managed by the tool)", k, typ)
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
		finished := strings.HasPrefix(e.rec.Path, "docs/work/done/")
		if finished && newStatus != "done" && newStatus != "dropped" {
			return nil, refusedRule("WK9", "a closed item stays in docs/work/done/ with status done or dropped")
		}
	} else if newStatus == "superseded" {
		if v, ok := e.doc.Get("superseded_by"); !ok || v.Kind != frontmatter.String {
			return nil, refusedRule("A3", "a superseded ADR needs superseded_by")
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
	return e.finish(O("record", recordObj(repo, rec, nil, true), "changed", changed)), nil
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
		e.text = record.SetSection(e.text, "Outcome", o, false, e.eol)
	}
	_, text := record.Sections(e.text)
	if strings.TrimSpace(text["outcome"]) == "" {
		return nil, refused("empty_outcome", "%s has no ## Outcome; write it before closing", e.rec.ID)
	}
	if status == "done" {
		for _, ch := range repo.Records {
			if ch.Type == record.Work && record.RefID(ch.Str("parent")) == e.rec.ID && ch.Status() != "done" && ch.Status() != "dropped" {
				return nil, &Error{"refused", "child " + ch.ID + " is not done or dropped", O("reason", "open_children", "rule", "WK8")}
			}
		}
	}
	if serr := e.doc.Set("status", frontmatter.NewString(status)); serr != nil {
		return nil, unreadable(serr)
	}
	for _, k := range []string{"claimed_by", "claimed_at"} {
		e.doc.Delete(k)
	}
	if serr := e.doc.Set("updated", frontmatter.NewString(today(clock(q)))); serr != nil {
		return nil, unreadable(serr)
	}
	from := e.rec.Path
	to := "docs/work/done/" + path.Base(from)
	rec, cerr := e.commit(to)
	if cerr != nil {
		return nil, cerr
	}
	return e.finish(O("record", recordObj(repo, rec, nil, true), "moved_from", from, "moved_to", rec.Path)), nil
}

// ---- claim and release ----

func claimsOn(repo *record.Repo) *Error {
	if !repo.Extension("claims") {
		return refused("extension_off", "the claims extension is not enabled in docs/.layout")
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
	for _, kv := range []struct{ k, v string }{{"claimed_by", who}, {"claimed_at", now.Format("2006-01-02T15:04")}, {"updated", today(now)}} {
		if serr := set(kv.k, frontmatter.NewString(kv.v)); serr != nil {
			return nil, serr
		}
	}
	e.text = record.AppendToSection(e.text, "log", "- "+today(now)+": "+logLine, e.eol)
	rec, cerr := e.commit("")
	if cerr != nil {
		return nil, cerr
	}
	return e.finish(O("record", recordObj(repo, rec, nil, true), "previous_holder", previous)), nil
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
	e.text = record.AppendToSection(e.text, "log", "- "+today(now)+": released ("+prev+").", e.eol)
	rec, cerr := e.commit("")
	if cerr != nil {
		return nil, cerr
	}
	return e.finish(O("record", recordObj(repo, rec, nil, true))), nil
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
	if strings.HasPrefix(e.rec.Path, "docs/work/done/") {
		return nil, refused("closed", "%s is closed; its sections are not edited", e.rec.ID)
	}
	text := record.SetSection(e.text, name, q.str("text"), mode == "append", e.eol)
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
