package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alexsasharegan/akmctl/internal/record"
)

// repoWith copies the clean fixture repository to a temp dir and applies edits
// (path -> function on the file text; nil deletes the file).
func repoWith(t *testing.T, edits map[string]func(string) string) *record.Repo {
	t.Helper()
	src := "../../testdata/api/repos/basic"
	dst := t.TempDir()
	filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		rel, _ := filepath.Rel(src, p)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		b, _ := os.ReadFile(p)
		return os.WriteFile(filepath.Join(dst, rel), b, 0o644)
	})
	for p, f := range edits {
		full := filepath.Join(dst, p)
		if f == nil {
			os.Remove(full)
			continue
		}
		b, _ := os.ReadFile(full)
		os.MkdirAll(filepath.Dir(full), 0o755)
		if err := os.WriteFile(full, []byte(f(string(b))), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	r, err := record.Load(dst)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func rules(res Result) map[string]int {
	m := map[string]int{}
	for _, f := range res.Findings {
		m[f.Rule]++
	}
	return m
}

func sub(old, new string) func(string) string {
	return func(s string) string {
		if !strings.Contains(s, old) {
			panic("test setup: " + old + " not found")
		}
		return strings.Replace(s, old, new, 1)
	}
}

var asOf = time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC)

func TestCleanRepository(t *testing.T) {
	res := Run(repoWith(t, nil), Options{AsOf: asOf})
	if len(res.Findings) != 0 || res.State != "current" {
		t.Fatalf("%+v", res)
	}
}

const (
	a2b3 = "docs/work/261020-a2b3-add-json-output-to-ls.md"
	e5f6 = "docs/work/261022-e5f6-document-the-exit-codes.md"
	g7h8 = "docs/work/261024-g7h8-publish-the-first-release.md"
	j9k0 = "docs/work/done/261010-j9k0-choose-a-licence.md"
	p3q4 = "docs/adr/261012-p3q4-store-records-as-markdown-files.md"
	r5s6 = "docs/adr/261018-r5s6-use-a-local-server-for-the-ui.md"
	rm   = "docs/ROADMAP.md"
)

func TestRules(t *testing.T) {
	cases := []struct {
		name  string
		rule  string
		sev   string
		edits map[string]func(string) string
	}{
		{"WK1 unparsable", "WK1", "E", map[string]func(string) string{e5f6: sub("status: ready", "status: ready\nstatus: raw")}},
		{"WK1 wrong id field", "WK1", "E", map[string]func(string) string{e5f6: sub("id: 261022-e5f6", "id: 261022-zzzz")}},
		{"WK2 created differs", "WK2", "W", map[string]func(string) string{e5f6: sub("created: 2026-10-22", "created: 2026-10-23")}},
		{"WK3 bad status", "WK3", "E", map[string]func(string) string{e5f6: sub("status: ready", "status: waiting")}},
		{"WK3 bad date", "WK3", "E", map[string]func(string) string{e5f6: sub("updated: 2026-10-22", "updated: 2026-13-45")}},
		{"WK4 unresolved parent", "WK4", "E", map[string]func(string) string{e5f6: sub("parent: 261015-c4d5", "parent: 261099-aaaa")}},
		{"WK4 self parent", "WK4", "E", map[string]func(string) string{e5f6: sub("parent: 261015-c4d5", "parent: 261022-e5f6")}},
		{"WK4 cycle", "WK4", "E", map[string]func(string) string{"docs/work/261015-c4d5-machine-readable-listings.md": sub("status: active", "status: active\nparent: 261020-a2b3")}},
		{"WK5 blocked without blocked_by", "WK5", "E", map[string]func(string) string{g7h8: sub("blocked_by: [261020-a2b3, \"legal review\"]\n", "")}},
		{"WK5 unresolved blocker", "WK5", "E", map[string]func(string) string{g7h8: sub("261020-a2b3", "261099-aaaa")}},
		{"WK6 active without next", "WK6", "E", map[string]func(string) string{a2b3: sub("## Next\nWrite the encoder test.\n\nSaved: 2026-10-27 18:00\n\n", "")}},
		{"WK7 no goal", "WK7", "E", map[string]func(string) string{e5f6: sub("## Goal", "## Aim")}},
		{"WK7 empty outcome", "WK7", "E", map[string]func(string) string{j9k0: sub("MIT, in `LICENSE`.", "")}},
		{"WK8 done parent with open child", "WK8", "E", map[string]func(string) string{"docs/work/261015-c4d5-machine-readable-listings.md": sub("status: active", "status: done")}},
		{"WK9 finished item outside done", "WK9", "E", map[string]func(string) string{e5f6: sub("status: ready", "status: dropped")}},
		{"WK10 stale", "WK10", "W", map[string]func(string) string{a2b3: sub("updated: 2026-10-27", "updated: 2026-08-01")}},
		{"WK11 bad kind", "WK11", "E", map[string]func(string) string{"docs/work/261015-c4d5-machine-readable-listings.md": sub("kind: feature", "kind: epic")}},
		{"WK11 priority text", "WK11", "E", map[string]func(string) string{a2b3: sub("priority: 2", "priority: high")}},
		{"WK11 decision not an ADR", "WK11", "E", map[string]func(string) string{a2b3: sub("priority: 2", "priority: 2\ndecisions: [261024-g7h8]")}},
		{"WK11 long summary", "WK11", "E", map[string]func(string) string{a2b3: sub("summary: Print ls results as JSON.", "summary: "+strings.Repeat("x", 161))}},
		{"WK12 heading differs", "WK12", "W", map[string]func(string) string{e5f6: sub("# Document the exit codes", "# Something else")}},
		{"I1 duplicate id", "I1", "E", map[string]func(string) string{"docs/work/261022-e5f6-copy.md": func(string) string {
			return "---\nid: 261022-e5f6\ntitle: Copy\nstatus: raw\ncreated: 2026-10-22\nupdated: 2026-10-22\n---\n# Copy\n"
		}}},
		{"A1 missing summary", "A1", "E", map[string]func(string) string{p3q4: sub("summary: Records are Markdown files with frontmatter; read this before changing the format.\n", "")}},
		{"A2 bad status", "A2", "E", map[string]func(string) string{r5s6: sub("status: proposed", "status: maybe")}},
		{"A3 superseded without target", "A3", "E", map[string]func(string) string{p3q4: sub("status: accepted", "status: superseded")}},
		{"A4 stale index", "A4", "E", map[string]func(string) string{"docs/adr/README.md": sub("Proposed", "Pending")}},
		{"A4 missing index", "A4", "E", map[string]func(string) string{"docs/adr/README.md": nil}},
		{"RM1 unresolved", "RM1", "E", map[string]func(string) string{rm: sub("## Later\n", "## Later\n- [261099-bbbb] Ghost\n")}},
		{"RM2 unknown section", "RM2", "E", map[string]func(string) string{rm: sub("## Later", "## Someday")}},
		{"RM2 duplicate", "RM2", "E", map[string]func(string) string{rm: sub("## Later\n", "## Later\n- [261015-c4d5] Again\n")}},
		{"RM3 done on the roadmap", "RM3", "E", map[string]func(string) string{rm: sub("## Later\n", "## Later\n- [261010-j9k0] Licence\n")}},
		{"RM4 top-level not on roadmap", "RM4", "W", map[string]func(string) string{rm: sub("- [261015-c4d5] Machine-readable listings\n", "")}},
		{"G1 missing heading", "G1", "E", map[string]func(string) string{"AGENTS.md": sub("## Rules", "## Conduct")}},
		{"G1 missing file", "G1", "E", map[string]func(string) string{"AGENTS.md": nil}},
		{"G2 long", "G2", "W", map[string]func(string) string{"AGENTS.md": func(s string) string { return s + strings.Repeat("- x\n", 80) }}},
		{"G3 missing script", "G3", "W", map[string]func(string) string{"AGENTS.md": sub("`akmctl check`", "`scripts/check.sh`")}},
		{"L1 broken link", "L1", "E", map[string]func(string) string{"AGENTS.md": func(s string) string { return s + "\nSee [the guide](docs/guide.md).\n" }}},
		{"L2 unresolved wikilink", "L2", "W", map[string]func(string) string{"AGENTS.md": func(s string) string { return s + "\nSee [[Nowhere]].\n" }}},
		{"S1 behind", "S1", "W", map[string]func(string) string{"docs/.layout": sub("layout: 1.0.0", "layout: 0.9.0")}},
		{"S1 missing manifest", "S1", "W", map[string]func(string) string{"docs/.layout": nil}},
		{"S2 skill name", "S2", "E", map[string]func(string) string{".agents/skills/demo/SKILL.md": func(string) string { return "---\nname: other\ndescription: x\n---\n" }}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := Run(repoWith(t, c.edits), Options{AsOf: asOf})
			var sev string
			for _, f := range res.Findings {
				if f.Rule == c.rule {
					sev = f.Severity
				}
			}
			if sev != c.sev {
				t.Fatalf("want %s %s, got %v", c.sev, c.rule, res.Findings)
			}
		})
	}
}

func TestLinksThatResolve(t *testing.T) {
	res := Run(repoWith(t, map[string]func(string) string{
		"AGENTS.md": func(s string) string {
			return s + "\n[roadmap](docs/ROADMAP.md) [adr](<docs/adr/README.md>) [site](https://example.com) [top](#commands) `[x](nope.md)` [[ROADMAP]] [[ROADMAP|alias]] [[ROADMAP#now]]\n```\n[y](nope.md) [[Nowhere]]\n```\n"
		},
	}), Options{AsOf: asOf})
	if len(res.Findings) != 0 {
		t.Fatalf("%+v", res.Findings)
	}
}

func TestClaims(t *testing.T) {
	r := repoWith(t, map[string]func(string) string{"docs/.layout": func(s string) string { return s + "extension.claims: 1.0.0\n" }})
	res := Run(r, Options{AsOf: asOf})
	if rules(res)["CL1"] != 2 {
		t.Fatalf("%v", res.Findings)
	}
	r = repoWith(t, map[string]func(string) string{
		"docs/.layout": func(s string) string { return s + "extension.claims: 1.0.0\n" },
		a2b3:           sub("status: active", "status: active\nclaimed_by: agent-1\nclaimed_at: 2026-10-27T10:00"),
		"docs/work/261015-c4d5-machine-readable-listings.md": sub("status: active", "status: active\nclaimed_by: agent-1\nclaimed_at: 2026-10-27T10:00"),
	})
	if got := rules(Run(r, Options{AsOf: asOf})); got["CL2"] != 2 || got["CL1"] != 0 {
		t.Fatalf("%v", got)
	}
}

func TestRuleFilter(t *testing.T) {
	r := repoWith(t, map[string]func(string) string{e5f6: sub("status: ready", "status: waiting"), rm: sub("## Later\n", "## Later\n- [261099-bbbb] Ghost\n")})
	got := rules(Run(r, Options{AsOf: asOf, Rules: []string{"RM"}}))
	if len(got) != 1 || got["RM1"] != 1 {
		t.Fatalf("%v", got)
	}
}
