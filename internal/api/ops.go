package api

import (
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alexsasharegan/akmctl/internal/check"
	"github.com/alexsasharegan/akmctl/internal/frontmatter"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// Info describes the running tool for the version operation.
type Info struct {
	Name, Version string
}

// Run executes a read operation against repo and returns the response object,
// ok or not, and the error if it failed. The response always has the envelope.
func Run(repo *record.Repo, info Info, op string, req Request) (Obj, *Error) {
	if req == nil {
		req = Request{}
	}
	echo := echoRequest(req)
	if err := validate(op, req); err != nil {
		return failure(op, echo, err), err
	}
	var body Obj
	var err *Error
	if writesRecords[op] {
		if e := manifestPathError(repo); e != nil {
			return failure(op, echo, e), e
		}
	}
	switch op {
	case "version":
		body = opVersion(info, repo)
	case "ls":
		body, err = opLs(repo, req)
	case "search":
		body, err = opSearch(repo, req)
	case "show":
		body, err = opShow(repo, req)
	case "board":
		body, err = opBoard(repo, req)
	case "find":
		body, err = opFind(repo, req)
	case "when":
		body, err = opWhen(req)
	case "check":
		body, err = opCheck(repo, req)
	case "index":
		body, err = opIndex(repo, req)
	case "migrate":
		body, err = opMigrate(repo, req)
	case "init":
		body, err = opInit(repo, req)
	case "config":
		body, err = opConfig(repo, req)
	case "new":
		body, err = opNew(repo, req)
	case "set":
		body, err = opSet(repo, req)
	case "section":
		body, err = opSection(repo, req)
	case "close":
		body, err = opClose(repo, req)
	case "claim":
		body, err = opClaim(repo, req)
	case "release":
		body, err = opRelease(repo, req)
	}
	if err != nil {
		return failure(op, echo, err), err
	}
	out := O("api_version", Version, "operation", op, "request", echo, "ok", true)
	return append(out, body...), nil
}

// writesRecords lists the operations that write under the manifest's folders.
var writesRecords = map[string]bool{"new": true, "set": true, "section": true, "close": true, "claim": true, "release": true, "index": true}

// manifestPathError refuses a write while the manifest names a bad folder
// (rule S3): the defaults would be used, and a write would land in the wrong place.
func manifestPathError(repo *record.Repo) *Error {
	for _, p := range repo.Problems {
		if p.Rule == "S3" {
			return refusedRule("S3", "%s: %s", p.Path, p.Msg)
		}
	}
	return nil
}

func echoRequest(r Request) Obj {
	keys := make([]string, 0, len(r))
	for k := range r {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	o := Obj{}
	for _, k := range keys {
		o = append(o, KV{k, r[k]})
	}
	return o
}

// Failure builds the response of an operation that could not run.
func Failure(op string, req Request, e *Error) Obj { return failure(op, echoRequest(req), e) }

func failure(op string, echo Obj, e *Error) Obj {
	eo := O("code", e.Code, "message", e.Message)
	if len(e.Details) > 0 {
		eo = append(eo, KV{"details", e.Details})
	}
	return O("api_version", Version, "operation", op, "request", echo, "ok", false, "error", eo)
}

func opVersion(info Info, repo *record.Repo) Obj {
	return O("tool", O("name", info.Name, "version", info.Version), "api_versions", []string{Version},
		"layout_versions", []string{check.LayoutVersion}, "extensions", []string{"claims"}, "repository", repositoryObj(repo))
}

// repositoryObj describes the repository the call ran against; null when the
// tool runs outside any repository (the CLI's version command).
func repositoryObj(repo *record.Repo) any {
	if repo == nil || repo.Root == "" {
		return nil
	}
	layout := "0.0.0"
	if v, ok := repo.Layout["layout"]; ok && v.Kind == frontmatter.String {
		layout = v.Str
	}
	exts := []string{}
	for k := range repo.Layout {
		if n, ok := strings.CutPrefix(k, "extension."); ok {
			exts = append(exts, n)
		}
	}
	sort.Strings(exts)
	style := "id"
	if v, ok := repo.Layout["reference_style"]; ok && v.Kind == frontmatter.String && v.Str == "wikilink" {
		style = "wikilink"
	}
	return O("name", filepath.Base(repo.Root), "layout", layout, "extensions", exts, "reference_style", style, "claim_stale_hours", claimStaleHours(repo),
		"work_dir", repo.WorkDir(), "adr_dir", repo.ADRDir(), "git_commits", repo.GitCommits())
}

// fieldsObj renders the frontmatter of r for the chosen field list. A nil list
// means the type's default set; "*" means all fields.
func fieldsObj(repo *record.Repo, r *record.Record, names []string, full bool) Obj {
	var keys []string
	switch {
	case full || (len(names) == 1 && names[0] == "*"):
		keys = append(keys, record.TypeFields(r.Type, repo.Extension("claims"))...)
		for _, k := range r.Keys {
			if !contains(keys, k) {
				keys = append(keys, k)
			}
		}
	case names != nil:
		keys = names
	default:
		keys = record.HeaderFields(r.Type)
	}
	o := Obj{}
	for _, k := range keys {
		o = append(o, KV{k, fieldJSON(r, k)})
	}
	return o
}

func fieldJSON(r *record.Record, k string) any {
	v, ok := r.Fields[k]
	if !ok {
		return nil
	}
	return record.JSONValue(k, v)
}

func contains(l []string, s string) bool {
	for _, x := range l {
		if x == s {
			return true
		}
	}
	return false
}

func recordObj(repo *record.Repo, r *record.Record, names []string, full bool) Obj {
	return O("id", r.ID, "type", r.Type, "path", r.Path, "id_date", r.IDDate(), "etag", r.Etag, "fields", fieldsObj(repo, r, names, full))
}

// ---- ls ----

func opLs(repo *record.Repo, q Request) (Obj, *Error) {
	var recs []*record.Record
	for _, r := range repo.Records {
		if matchesFilters(r, q) {
			recs = append(recs, r)
		}
	}
	sortKeys := q.strs("sort")
	if sortKeys == nil {
		sortKeys = []string{"id"}
	}
	if err := checkSortKeys(sortKeys); err != nil {
		return nil, err
	}
	sortRecords(repo, recs, sortKeys)
	total := len(recs)
	page := paginate(recs, q.integer("offset", 0), q.integer("limit", -1))
	if q.boolean("count_only") {
		return O("total", total, "count", 0, "records", []Obj{}), nil
	}
	names := q.strs("fields")
	out := make([]Obj, 0, len(page))
	for _, r := range page {
		o := recordObj(repo, r, names, false)
		if q.boolean("next") && r.Type == record.Work {
			if l, ok := r.NextLine(); ok && l != "" {
				o = append(o, KV{"next_line", l})
			} else {
				o = append(o, KV{"next_line", nil})
			}
		}
		out = append(out, o)
	}
	return O("total", total, "count", len(out), "records", out), nil
}

func matchesFilters(r *record.Record, q Request) bool {
	if t := q.str("type"); t != "" && r.Type != t {
		return false
	}
	if st := q.strs("status"); st != nil && !contains(st, r.Status()) {
		return false
	}
	if _, ok := q["parent"]; ok {
		if r.Type != record.Work || record.RefID(r.Str("parent")) != q.str("parent") {
			return false
		}
	}
	if _, ok := q["blocked"]; ok {
		if r.Type != record.Work {
			return false
		}
		blocked := len(r.Fields["blocked_by"].Items) > 0
		if blocked != q.boolean("blocked") {
			return false
		}
	}
	if c := q.str("cycle"); c != "" {
		if r.Type != record.Work || r.Str("cycle") != c {
			return false
		}
	}
	if sc, ok := q["scope"]; ok {
		if r.Type != record.ADR || !scopeMatch(r.Fields["scope"].Strings(), sc.(string)) {
			return false
		}
	}
	if since := q.str("since"); since != "" {
		d := r.Str("updated")
		if r.Type == record.ADR {
			d = r.Str("date")
		}
		if d == "" || d < since {
			return false
		}
	}
	return true
}

func scopeMatch(entries []string, value string) bool {
	for _, e := range entries {
		if e == value || (strings.HasSuffix(e, "/") && strings.HasPrefix(value, e)) {
			return true
		}
	}
	return false
}

var sortable = map[string]bool{"id": true, "title": true, "status": true, "kind": true, "priority": true, "created": true, "updated": true, "date": true, "roadmap": true}

func checkSortKeys(keys []string) *Error {
	for _, k := range keys {
		if !sortable[strings.TrimPrefix(k, "-")] {
			return invalid("sort", "sort: unknown key %q", k)
		}
	}
	return nil
}

type sortVal struct {
	present bool
	num     int64
	str     string
	isNum   bool
}

func valueOf(repo *record.Repo, r *record.Record, key string) sortVal {
	switch key {
	case "id":
		return sortVal{present: true, str: r.ID}
	case "roadmap":
		if n := repo.RoadmapRank(r.ID); n >= 0 {
			return sortVal{present: true, isNum: true, num: int64(n)}
		}
		return sortVal{}
	}
	v, ok := r.Fields[key]
	switch {
	case !ok || v.Kind == frontmatter.Null:
		return sortVal{}
	case v.Kind == frontmatter.Int:
		return sortVal{present: true, isNum: true, num: v.Int}
	case v.Kind == frontmatter.String:
		return sortVal{present: true, str: v.Str}
	}
	return sortVal{}
}

func sortRecords(repo *record.Repo, recs []*record.Record, keys []string) {
	sort.SliceStable(recs, func(i, j int) bool {
		for _, k := range keys {
			desc := strings.HasPrefix(k, "-")
			a, b := valueOf(repo, recs[i], strings.TrimPrefix(k, "-")), valueOf(repo, recs[j], strings.TrimPrefix(k, "-"))
			if a.present != b.present {
				return a.present // absent sorts last in both directions
			}
			if !a.present {
				continue
			}
			c := 0
			if a.isNum && b.isNum {
				switch {
				case a.num < b.num:
					c = -1
				case a.num > b.num:
					c = 1
				}
			} else {
				c = strings.Compare(a.str, b.str)
			}
			if c != 0 {
				return (c < 0) != desc
			}
		}
		return false
	})
}

func paginate(recs []*record.Record, offset, limit int) []*record.Record {
	if offset > len(recs) {
		offset = len(recs)
	}
	recs = recs[offset:]
	if limit >= 0 && limit < len(recs) {
		recs = recs[:limit]
	}
	return recs
}

// ---- search ----

type hit struct {
	rec     *record.Record
	matches []Obj
}

func opSearch(repo *record.Repo, q Request) (Obj, *Error) {
	words := strings.Fields(strings.ToLower(q.str("text")))
	if len(words) == 0 {
		return nil, invalid("text", "text must not be empty")
	}
	var hits []hit
	for _, r := range repo.Records {
		if t := q.str("type"); t != "" && r.Type != t {
			continue
		}
		if st := q.strs("status"); st != nil && !contains(st, r.Status()) {
			continue
		}
		if m, ok := searchRecord(r, words, q.boolean("body")); ok {
			hits = append(hits, hit{r, m})
		}
	}
	sort.SliceStable(hits, func(i, j int) bool {
		if len(hits[i].matches) != len(hits[j].matches) {
			return len(hits[i].matches) > len(hits[j].matches)
		}
		return hits[i].rec.ID < hits[j].rec.ID
	})
	total := len(hits)
	off := q.integer("offset", 0)
	if off > len(hits) {
		off = len(hits)
	}
	hits = hits[off:]
	if lim := q.integer("limit", -1); lim >= 0 && lim < len(hits) {
		hits = hits[:lim]
	}
	names := q.strs("fields")
	out := make([]Obj, 0, len(hits))
	for _, h := range hits {
		o := recordObj(repo, h.rec, names, false)
		out = append(out, append(o, KV{"matches", h.matches}))
	}
	return O("total", total, "count", len(out), "records", out), nil
}

func searchRecord(r *record.Record, words []string, body bool) ([]Obj, bool) {
	lines := strings.Split(string(r.Text), "\n")
	var ms []Obj
	covered := make([]bool, len(words))
	check := func(field string, lineNo int, text string) {
		low := strings.ToLower(text)
		any := false
		for i, w := range words {
			if strings.Contains(low, w) {
				covered[i], any = true, true
			}
		}
		if any {
			sn := strings.TrimRight(text, "\r")
			if len(sn) > 160 {
				sn = sn[:160]
			}
			ms = append(ms, O("field", field, "line", lineNo, "snippet", sn))
		}
	}
	for _, f := range []string{"title", "summary", "scope"} {
		v, ok := r.Fields[f]
		if !ok {
			continue
		}
		n := r.Lines[f]
		check(f, n, valueText(v, lines, n))
	}
	if body {
		start := len(lines)
		for i, l := range lines {
			if i > 0 && strings.TrimRight(l, "\r") == "---" {
				start = i + 1
				break
			}
		}
		for i := start; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) != "" {
				check("body", i+1, lines[i])
			}
		}
	}
	for _, c := range covered {
		if !c {
			return nil, false
		}
	}
	return ms, true
}

// valueText returns the searchable text of a field: the line in the file.
func valueText(v frontmatter.Value, lines []string, n int) string {
	if n >= 1 && n <= len(lines) {
		return lines[n-1]
	}
	return v.String()
}

// ---- show, find, when, board ----

func lookup(repo *record.Repo, id string) (*record.Record, *Error) {
	recs := repo.ByID(id)
	switch len(recs) {
	case 0:
		return nil, &Error{Code: "not_found", Message: "no record with id " + id}
	case 1:
		return recs[0], nil
	}
	paths := make([]string, len(recs))
	for i, r := range recs {
		paths[i] = r.Path
	}
	return nil, &Error{"ambiguous", "id " + id + " names more than one record", O("paths", paths)}
}

func opShow(repo *record.Repo, q Request) (Obj, *Error) {
	r, err := lookup(repo, q.str("id"))
	if err != nil {
		return nil, err
	}
	names, text := record.Sections(r.Body)
	if names == nil {
		names = []string{}
	}
	out := O("record", recordObj(repo, r, nil, true), "sections", names)
	if sec, ok := q["section"]; ok {
		name := strings.ToLower(strings.TrimSpace(sec.(string)))
		t, found := text[name]
		if !found {
			return nil, &Error{Code: "not_found", Message: "no section " + name + " in " + r.ID}
		}
		return append(out, KV{"section", O("name", name, "text", t)}), nil
	}
	return append(out, KV{"body", r.Body}), nil
}

func opFind(repo *record.Repo, q Request) (Obj, *Error) {
	r, err := lookup(repo, q.str("id"))
	if err != nil {
		return nil, err
	}
	return O("id", r.ID, "type", r.Type, "path", r.Path), nil
}

func opWhen(q Request) (Obj, *Error) {
	id := q.str("id")
	d, err := time.Parse("2006-01-02", record.IDDate(id))
	if err != nil {
		return nil, invalid("id", "id %q does not hold a valid date", id)
	}
	y, m, dd := time.Now().Date()
	asOf := time.Date(y, m, dd, 0, 0, 0, 0, time.UTC) // the local day
	if s := q.str("as_of"); s != "" {
		asOf, err = time.Parse("2006-01-02", s)
		if err != nil {
			return nil, invalid("as_of", "as_of is not a valid date")
		}
	}
	return O("id", id, "date", d.Format("2006-01-02"), "weekday", d.Weekday().String(), "age_days", int(asOf.Sub(d).Hours()/24)), nil
}

func opBoard(repo *record.Repo, q Request) (Obj, *Error) {
	statuses := q.strs("statuses")
	if statuses == nil {
		statuses = []string{"active", "blocked", "ready"}
	}
	totals := map[string]int{}
	for _, r := range repo.Records {
		if r.Type == record.Work {
			totals[r.Status()]++
		}
	}
	cols := make([]Obj, 0, len(statuses))
	for _, st := range statuses {
		var recs []*record.Record
		for _, r := range repo.Records {
			if r.Type == record.Work && r.Status() == st {
				recs = append(recs, r)
			}
		}
		sortRecords(repo, recs, []string{"roadmap", "-priority", "id"})
		rs := make([]Obj, 0, len(recs))
		for _, r := range recs {
			rs = append(rs, recordObj(repo, r, append(append([]string(nil), record.HeaderFields(record.Work)...), "blocked_by", "cycle", "keep"), false))
		}
		cols = append(cols, O("status", st, "count", len(rs), "records", rs))
	}
	tot := Obj{}
	for _, st := range workStatus {
		tot = append(tot, KV{st, totals[st]})
	}
	return O("columns", cols, "totals", tot), nil
}
