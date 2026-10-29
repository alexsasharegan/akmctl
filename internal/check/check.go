// Package check runs the rules of the AKM layout spec over a repository.
package check

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/alexsasharegan/akmctl/internal/frontmatter"
	"github.com/alexsasharegan/akmctl/internal/index"
	"github.com/alexsasharegan/akmctl/internal/migrate"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// LayoutVersion is the layout version this tool checks against.
const LayoutVersion = "1.0.0"

// Finding is one rule violation.
type Finding struct {
	Rule     string
	Severity string // "E" or "W"
	File     string // "" when not tied to a file
	Line     int    // 0 when not tied to a line
	ID       string
	Message  string
}

// Options configure a run.
type Options struct {
	Rules []string  // rule ids or prefixes; empty means all
	AsOf  time.Time // the clock for staleness rules
}

// Result is the outcome of a run.
type Result struct {
	Findings []Finding
	State    string // reconcile state of the manifest
	Steps    []string
}

var (
	reDate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	reTS   = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$`)
)

var workStatuses = []string{"raw", "ready", "active", "blocked", "done", "dropped"}
var adrStatuses = []string{"proposed", "accepted", "rejected", "deprecated", "superseded"}
var kinds = []string{"feature", "bug", "chore", "spike"}

type checker struct {
	repo *record.Repo
	opt  Options
	out  []Finding
	// settings
	staleDays, claimStaleHours, warnLines, maxLines int
}

func (c *checker) add(rule, sev, file string, line int, id, format string, a ...any) {
	if len(c.opt.Rules) > 0 {
		ok := false
		for _, r := range c.opt.Rules {
			if strings.HasPrefix(rule, r) {
				ok = true
			}
		}
		if !ok {
			return
		}
	}
	c.out = append(c.out, Finding{rule, sev, file, line, id, fmt.Sprintf(format, a...)})
}

// Run checks the repository.
func Run(repo *record.Repo, opt Options) Result {
	if opt.AsOf.IsZero() {
		opt.AsOf = time.Now().UTC()
	}
	c := &checker{repo: repo, opt: opt, staleDays: 30, claimStaleHours: 24, warnLines: 80, maxLines: 150}
	c.settings()
	c.problems()
	c.work()
	c.adrs()
	c.ids()
	c.roadmap()
	c.agents()
	c.links()
	c.skills()
	if repo.Extension("claims") {
		c.claims()
	}
	if repo.Extension("safety") {
		c.safety()
	}
	state, steps := c.manifest()
	sort.SliceStable(c.out, func(i, j int) bool {
		a, b := c.out[i], c.out[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		return a.Rule < b.Rule
	})
	return Result{c.out, state, steps}
}

func (c *checker) settings() {
	get := func(k string, p *int) {
		if v, ok := c.repo.Layout[k]; ok && v.Kind == frontmatter.Int && v.Int > 0 {
			*p = int(v.Int)
		}
	}
	get("stale_days", &c.staleDays)
	get("claim_stale_hours", &c.claimStaleHours)
	get("agents_md_warn_lines", &c.warnLines)
	get("agents_md_max_lines", &c.maxLines)
}

func in(l []string, s string) bool {
	for _, x := range l {
		if x == s {
			return true
		}
	}
	return false
}

func validDate(s string) bool {
	if !reDate.MatchString(s) {
		return false
	}
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

// ref resolves a reference (bare id or wikilink) to a record.
func (c *checker) ref(s string) *record.Record {
	id := record.RefID(s)
	if !record.ValidID(id) {
		return nil
	}
	if recs := c.repo.ByID(id); len(recs) > 0 {
		return recs[0]
	}
	return nil
}

func (c *checker) problems() {
	for _, p := range c.repo.Problems {
		switch {
		case p.Rule != "":
			c.add(p.Rule, "E", p.Path, p.Line, "", "manifest: %s", p.Msg)
		case c.repo.InWork(p.Path):
			c.add("WK1", "E", p.Path, p.Line, "", "%s", p.Msg)
		case c.repo.InADR(p.Path):
			c.add("A1", "E", p.Path, p.Line, "", "%s", p.Msg)
		case p.Path == record.ManifestPath:
			c.add("S1", "W", p.Path, p.Line, "", "manifest: %s", p.Msg)
		}
	}
}

func (c *checker) line(r *record.Record, k string) int { return r.Lines[k] }

func (c *checker) work() {
	for _, r := range c.repo.Records {
		if r.Type != record.Work {
			continue
		}
		c.workRecord(r)
	}
}

func (c *checker) workRecord(r *record.Record) {
	st := r.Status()
	done := c.repo.IsDone(r.Path)
	// WK1: id field
	if v, ok := r.Fields["id"]; !ok || v.Kind != frontmatter.String || v.Str != r.ID {
		c.add("WK1", "E", r.Path, c.line(r, "id"), r.ID, "id field must equal the filename's id %s", r.ID)
	}
	// WK2
	if cr := r.Str("created"); cr != "" && validDate(cr) && cr != r.IDDate() {
		c.add("WK2", "W", r.Path, c.line(r, "created"), r.ID, "created %s differs from the date in the id (%s)", cr, r.IDDate())
	}
	// WK3
	if st == "" {
		c.add("WK3", "E", r.Path, 0, r.ID, "status is missing")
	} else if !in(workStatuses, st) {
		c.add("WK3", "E", r.Path, c.line(r, "status"), r.ID, "status %q is not one of %s", st, strings.Join(workStatuses, ", "))
	}
	for _, k := range []string{"title", "created", "updated"} {
		if r.Str(k) == "" {
			c.add("WK3", "E", r.Path, c.line(r, k), r.ID, "%s is missing", k)
		}
	}
	for _, k := range []string{"created", "updated"} {
		if v := r.Str(k); v != "" && !validDate(v) {
			c.add("WK3", "E", r.Path, c.line(r, k), r.ID, "%s %q is not a date (YYYY-MM-DD)", k, v)
		}
	}
	// WK4
	if p := r.Str("parent"); p != "" {
		pr := c.ref(p)
		switch {
		case pr == nil || pr.Type != record.Work:
			c.add("WK4", "E", r.Path, c.line(r, "parent"), r.ID, "parent %s does not resolve to an existing item", record.RefID(p))
		case pr.ID == r.ID:
			c.add("WK4", "E", r.Path, c.line(r, "parent"), r.ID, "an item cannot be its own parent")
		case c.cycle(r):
			c.add("WK4", "E", r.Path, c.line(r, "parent"), r.ID, "parent chain forms a cycle")
		}
	}
	// WK5
	by := r.Fields["blocked_by"]
	if st == "blocked" && len(by.Items) == 0 {
		c.add("WK5", "E", r.Path, c.line(r, "blocked_by"), r.ID, "a blocked item needs a non-empty blocked_by")
	}
	for _, it := range by.Items {
		if it.Kind == frontmatter.String {
			if id := record.RefID(it.Str); record.ValidID(id) && c.ref(id) == nil {
				c.add("WK5", "E", r.Path, c.line(r, "blocked_by"), r.ID, "blocked_by %s does not resolve to a record", id)
			}
		}
	}
	// WK6, WK7
	names, text := record.Sections(r.Body)
	_ = names
	if (st == "active" || st == "blocked") && strings.TrimSpace(text["next"]) == "" {
		c.add("WK6", "E", r.Path, 0, r.ID, "status %s requires a non-empty ## Next", st)
	}
	if (st == "ready" || st == "active" || st == "blocked" || st == "done") && !hasSection(names, "goal") {
		c.add("WK7", "E", r.Path, 0, r.ID, "status %s requires a ## Goal", st)
	}
	if (st == "done" || st == "dropped") && strings.TrimSpace(text["outcome"]) == "" {
		c.add("WK7", "E", r.Path, 0, r.ID, "status %s requires a non-empty ## Outcome", st)
	}
	// WK8
	if st == "done" {
		for _, ch := range c.repo.Records {
			if ch.Type == record.Work && record.RefID(ch.Str("parent")) == r.ID && ch.Status() != "done" && ch.Status() != "dropped" {
				c.add("WK8", "E", r.Path, c.line(r, "status"), r.ID, "done parent has open child %s", ch.ID)
			}
		}
	}
	// WK9
	if in(workStatuses, st) {
		finished := st == "done" || st == "dropped"
		switch {
		case done && !finished:
			c.add("WK9", "E", r.Path, c.line(r, "status"), r.ID, "an item in %s/ must be done or dropped, not %s", c.repo.DoneDir(), st)
		case !done && finished:
			c.add("WK9", "E", r.Path, c.line(r, "status"), r.ID, "a %s item belongs in %s/", st, c.repo.DoneDir())
		}
	}
	// WK10
	if st == "active" {
		if u := r.Str("updated"); validDate(u) {
			t, _ := time.Parse("2006-01-02", u)
			if days := int(c.opt.AsOf.Sub(t).Hours() / 24); days > c.staleDays {
				c.add("WK10", "W", r.Path, c.line(r, "updated"), r.ID, "active item not updated for %d days", days)
			}
		}
	}
	// WK11
	if k := r.Str("kind"); k != "" && !in(kinds, k) {
		c.add("WK11", "E", r.Path, c.line(r, "kind"), r.ID, "kind %q is not one of %s", k, strings.Join(kinds, ", "))
	}
	if v, ok := r.Fields["priority"]; ok && v.Kind != frontmatter.Int && v.Kind != frontmatter.Null {
		c.add("WK11", "E", r.Path, c.line(r, "priority"), r.ID, "priority must be an integer")
	}
	for _, k := range []string{"blocked_by", "decisions"} {
		if v, ok := r.Fields[k]; ok && v.Kind != frontmatter.List && v.Kind != frontmatter.Null {
			c.add("WK11", "E", r.Path, c.line(r, k), r.ID, "%s must be a list", k)
		}
	}
	for _, it := range r.Fields["decisions"].Items {
		if d := c.ref(it.Str); it.Kind != frontmatter.String || d == nil || d.Type != record.ADR {
			c.add("WK11", "E", r.Path, c.line(r, "decisions"), r.ID, "decisions entry %s does not resolve to an ADR", it.String())
		}
	}
	if n := len([]rune(r.Str("summary"))); n > 160 {
		c.add("WK11", "E", r.Path, c.line(r, "summary"), r.ID, "summary is %d characters; the limit is 160", n)
	}
	// WK12
	if h := record.FirstHeading(r.Body); h != "" && r.Str("title") != "" && h != r.Str("title") {
		c.add("WK12", "W", r.Path, 0, r.ID, "title differs from the first heading %q", h)
	}
}

func hasSection(names []string, n string) bool { return in(names, n) }

func (c *checker) cycle(r *record.Record) bool {
	seen := map[string]bool{r.ID: true}
	cur := r
	for {
		p := c.ref(cur.Str("parent"))
		if p == nil || p.Type != record.Work {
			return false
		}
		if seen[p.ID] {
			return true
		}
		seen[p.ID] = true
		cur = p
	}
}

func (c *checker) adrs() {
	for _, r := range c.repo.Records {
		if r.Type != record.ADR {
			continue
		}
		if v, ok := r.Fields["id"]; !ok || v.Kind != frontmatter.String || v.Str != r.ID {
			c.add("A1", "E", r.Path, c.line(r, "id"), r.ID, "id field must equal the filename's id %s", r.ID)
		}
		for _, k := range []string{"title", "summary", "status", "date"} {
			if r.Str(k) == "" {
				c.add("A1", "E", r.Path, c.line(r, k), r.ID, "%s is missing", k)
			}
		}
		if n := len([]rune(r.Str("summary"))); n > 160 {
			c.add("A1", "E", r.Path, c.line(r, "summary"), r.ID, "summary is %d characters; the limit is 160", n)
		}
		if st := r.Status(); st != "" && !in(adrStatuses, st) {
			c.add("A2", "E", r.Path, c.line(r, "status"), r.ID, "status %q is not one of %s", st, strings.Join(adrStatuses, ", "))
		}
		if d := r.Str("date"); d != "" && !validDate(d) {
			c.add("A2", "E", r.Path, c.line(r, "date"), r.ID, "date %q is not a date (YYYY-MM-DD)", d)
		}
		c.supersession(r)
	}
	// A4
	want := index.Generate(c.repo)
	indexPath := c.repo.ADRIndexPath()
	got, err := os.ReadFile(filepath.Join(c.repo.Root, filepath.FromSlash(indexPath)))
	hasADR := false
	for _, r := range c.repo.Records {
		if r.Type == record.ADR {
			hasADR = true
		}
	}
	switch {
	case err != nil && hasADR:
		c.add("A4", "E", indexPath, 0, "", "the ADR index is missing; run akmctl index")
	case err == nil && string(got) != want:
		c.add("A4", "E", indexPath, firstDiffLine(string(got), want), "", "the ADR index differs from what akmctl index generates; run akmctl index")
	}
}

func firstDiffLine(a, b string) int {
	al, bl := strings.Split(a, "\n"), strings.Split(b, "\n")
	for i := 0; i < len(al) && i < len(bl); i++ {
		if al[i] != bl[i] {
			return i + 1
		}
	}
	if len(al) < len(bl) {
		return len(al)
	}
	return len(bl)
}

func (c *checker) supersession(r *record.Record) {
	sup, by := r.Str("supersedes"), r.Str("superseded_by")
	if r.Status() == "superseded" && by == "" {
		c.add("A3", "E", r.Path, c.line(r, "status"), r.ID, "a superseded ADR needs superseded_by")
	}
	if by != "" {
		t := c.ref(by)
		switch {
		case t == nil || t.Type != record.ADR:
			c.add("A3", "E", r.Path, c.line(r, "superseded_by"), r.ID, "superseded_by %s does not resolve to an ADR", record.RefID(by))
		case t.ID == r.ID:
			c.add("A3", "E", r.Path, c.line(r, "superseded_by"), r.ID, "an ADR cannot supersede itself")
		default:
			if s := t.Str("supersedes"); s != "" && record.RefID(s) != r.ID {
				c.add("A3", "E", r.Path, c.line(r, "superseded_by"), r.ID, "%s supersedes %s, not this ADR", t.ID, record.RefID(s))
			}
		}
	}
	if sup != "" {
		t := c.ref(sup)
		if t == nil || t.Type != record.ADR {
			c.add("A3", "E", r.Path, c.line(r, "supersedes"), r.ID, "supersedes %s does not resolve to an ADR", record.RefID(sup))
		} else if t.ID == r.ID {
			c.add("A3", "E", r.Path, c.line(r, "supersedes"), r.ID, "an ADR cannot supersede itself")
		}
	}
}

// ids is rule I1.
func (c *checker) ids() {
	byID := map[string][]*record.Record{}
	for _, r := range c.repo.Records {
		byID[r.ID] = append(byID[r.ID], r)
	}
	for id, rs := range byID {
		for _, r := range rs[1:] {
			c.add("I1", "E", r.Path, 0, id, "id %s is also used by %s", id, rs[0].Path)
		}
	}
}

var roadmapEntry = regexp.MustCompile(`^- \[([^\]]*)\]`)

func (c *checker) roadmap() {
	path := c.repo.RoadmapPath()
	b, err := os.ReadFile(filepath.Join(c.repo.Root, filepath.FromSlash(path)))
	if err != nil {
		return
	}
	section := ""
	seen := map[string]bool{}
	for i, l := range strings.Split(string(b), "\n") {
		n := i + 1
		l = strings.TrimRight(l, "\r")
		if strings.HasPrefix(l, "## ") {
			section = strings.TrimSpace(l[3:])
			if !in([]string{"Now", "Next", "Later", "Dropped"}, section) {
				c.add("RM2", "E", path, n, "", "unknown roadmap section %q (use Now, Next, Later, Dropped)", section)
			}
			continue
		}
		m := roadmapEntry.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		id := m[1]
		r := c.ref(id)
		if r == nil || r.Type != record.Work {
			c.add("RM1", "E", path, n, "", "roadmap entry %s does not resolve to a work item", id)
			continue
		}
		if seen[id] {
			c.add("RM2", "E", path, n, id, "%s appears more than once", id)
		}
		seen[id] = true
		st := r.Status()
		switch {
		case section != "Dropped" && (st == "done" || st == "dropped"):
			c.add("RM3", "E", path, n, id, "%s is %s but listed under %s", id, st, section)
		case section == "Dropped" && st != "dropped":
			c.add("RM3", "E", path, n, id, "%s is listed under Dropped but its status is %s", id, st)
		}
	}
	for _, r := range c.repo.Records {
		if r.Type == record.Work && (r.Status() == "ready" || r.Status() == "active") && r.Str("parent") == "" && !seen[r.ID] {
			c.add("RM4", "W", r.Path, 0, r.ID, "top-level %s item is not on the roadmap", r.Status())
		}
	}
}

var requiredHeadings = []string{"Commands", "Where things live", "Rules", "Do not, without asking"}

func (c *checker) agents() {
	b, err := os.ReadFile(filepath.Join(c.repo.Root, "AGENTS.md"))
	if err != nil {
		c.add("G1", "E", "AGENTS.md", 0, "", "AGENTS.md is missing")
		return
	}
	lines := strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	have := map[string]bool{}
	fence := false
	for _, l := range lines {
		if strings.HasPrefix(l, "```") {
			fence = !fence
		}
		if !fence && strings.HasPrefix(l, "## ") {
			have[strings.TrimSpace(l[3:])] = true
		}
	}
	for _, h := range requiredHeadings {
		if !have[h] {
			c.add("G1", "E", "AGENTS.md", 0, "", "missing the heading \"## %s\"", h)
		}
	}
	switch n := len(lines); {
	case n > c.maxLines:
		c.add("G2", "E", "AGENTS.md", 0, "", "AGENTS.md has %d lines; the limit is %d", n, c.maxLines)
	case n > c.warnLines:
		c.add("G2", "W", "AGENTS.md", 0, "", "AGENTS.md has %d lines; over %d is long", n, c.warnLines)
	}
	c.commandPaths(lines)
}

var scriptPath = regexp.MustCompile("`((?:\\./)?scripts/[A-Za-z0-9_./-]+)`")

// commandPaths is the part of G3 that can be checked without running anything:
// scripts/... paths named in the Commands section must exist.
func (c *checker) commandPaths(lines []string) {
	in := false
	for i, l := range lines {
		if strings.HasPrefix(l, "## ") {
			in = strings.TrimSpace(l[3:]) == "Commands"
			continue
		}
		if !in {
			continue
		}
		for _, m := range scriptPath.FindAllStringSubmatch(l, -1) {
			if _, err := os.Stat(filepath.Join(c.repo.Root, filepath.FromSlash(strings.TrimPrefix(m[1], "./")))); err != nil {
				c.add("G3", "W", "AGENTS.md", i+1, "", "command path %s does not exist", m[1])
			}
		}
	}
}

func (c *checker) claims() {
	holders := map[string][]*record.Record{}
	for _, r := range c.repo.Records {
		if r.Type != record.Work || r.Status() != "active" {
			continue
		}
		by := r.Str("claimed_by")
		if by == "" {
			c.add("CL1", "W", r.Path, 0, r.ID, "active item has no claimed_by")
			continue
		}
		if at := r.Str("claimed_at"); at != "" && !reTS.MatchString(at) {
			c.add("CL1", "W", r.Path, c.line(r, "claimed_at"), r.ID, "claimed_at %q is not a timestamp (YYYY-MM-DDTHH:MM)", at)
		}
		holders[by] = append(holders[by], r)
	}
	for by, rs := range holders {
		if len(rs) > 1 {
			for _, r := range rs {
				c.add("CL2", "W", r.Path, c.line(r, "claimed_by"), r.ID, "%s holds %d active items", by, len(rs))
			}
		}
	}
}

func (c *checker) skills() {
	dir := filepath.Join(c.repo.Root, ".agents", "skills")
	ents, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range ents {
		if !e.IsDir() {
			continue
		}
		rel := ".agents/skills/" + e.Name() + "/SKILL.md"
		b, err := os.ReadFile(filepath.Join(c.repo.Root, filepath.FromSlash(rel)))
		if err != nil {
			c.add("S2", "E", ".agents/skills/"+e.Name(), 0, "", "skill has no SKILL.md")
			continue
		}
		doc, err := frontmatter.Parse(b)
		if err != nil {
			c.add("S2", "E", rel, 0, "", "SKILL.md frontmatter: %v", err)
			continue
		}
		if v, _ := doc.Get("name"); v.Str != e.Name() {
			c.add("S2", "E", rel, doc.Line("name"), "", "name %q must equal the directory %q", v.Str, e.Name())
		}
		if v, _ := doc.Get("description"); v.Kind != frontmatter.String || strings.TrimSpace(v.Str) == "" {
			c.add("S2", "E", rel, doc.Line("description"), "", "description is missing")
		}
	}
}

// manifest is rule S1: it reconciles docs/.layout with the layout version the
// tool supports.
func (c *checker) manifest() (string, []string) {
	v, ok := c.repo.Layout["layout"]
	if !ok {
		if _, err := os.Stat(filepath.Join(c.repo.Root, filepath.FromSlash(record.ManifestPath))); err != nil {
			c.add("S1", "W", record.ManifestPath, 0, "", "the manifest %s is missing", record.ManifestPath)
		} else {
			c.add("S1", "W", record.ManifestPath, 0, "", "the manifest has no layout key")
		}
		return "absent", nil
	}
	repoV := v.Str
	switch cmp := compareVersions(repoV, LayoutVersion); {
	case cmp == 0:
		return "current", nil
	case cmp > 0:
		c.add("S1", "W", record.ManifestPath, c.manifestLine("layout"), "", "the repository is on layout %s, newer than this tool (%s): upgrade the tool", repoV, LayoutVersion)
		return "ahead", nil
	default:
		if chain, ok := migrate.Chain(repoV, LayoutVersion); ok {
			var steps []string
			for _, st := range migrate.Steps(chain) {
				steps = append(steps, st.ID+": "+st.Description)
			}
			c.add("S1", "W", record.ManifestPath, c.manifestLine("layout"), "", "the repository is on layout %s; run akmctl migrate to reach %s (%d steps)", repoV, LayoutVersion, len(steps))
			return "behind", steps
		}
		c.add("S1", "W", record.ManifestPath, c.manifestLine("layout"), "", "the repository is on layout %s; this tool supports %s and has no migration path from it", repoV, LayoutVersion)
		return "unsupported", nil
	}
}

func (c *checker) manifestLine(key string) int {
	b, err := os.ReadFile(filepath.Join(c.repo.Root, filepath.FromSlash(record.ManifestPath)))
	if err != nil {
		return 0
	}
	if d, err := frontmatter.ParseBare(b); err == nil {
		return d.Line(key)
	}
	return 0
}

func compareVersions(a, b string) int {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < 3; i++ {
		var x, y int
		if i < len(pa) {
			fmt.Sscanf(pa[i], "%d", &x)
		}
		if i < len(pb) {
			fmt.Sscanf(pb[i], "%d", &y)
		}
		if x != y {
			if x < y {
				return -1
			}
			return 1
		}
	}
	return 0
}
