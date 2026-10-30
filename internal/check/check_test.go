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
		{"WK13 done item without closed", "WK13", "W", map[string]func(string) string{j9k0: sub("closed: 2026-10-12\n", "")}},
		{"WK14 open item with closed", "WK14", "E", map[string]func(string) string{e5f6: sub("updated: 2026-10-22", "updated: 2026-10-22\nclosed: 2026-10-22")}},
		{"WK15 closed not a date", "WK15", "E", map[string]func(string) string{j9k0: sub("closed: 2026-10-12", "closed: 2026-13-45")}},
		{"WK15 closed before created", "WK15", "E", map[string]func(string) string{j9k0: sub("closed: 2026-10-12", "closed: 2026-10-09")}},
		{"WK16 cycle with a space", "WK16", "E", map[string]func(string) string{a2b3: sub("priority: 2", "priority: 2\ncycle: has space")}},
		{"WK16 cycle too long", "WK16", "E", map[string]func(string) string{a2b3: sub("priority: 2", "priority: 2\ncycle: "+strings.Repeat("a", 65))}},
		{"WK16 cycle starts with a dash", "WK16", "E", map[string]func(string) string{a2b3: sub("priority: 2", "priority: 2\ncycle: \"-x\"")}},
		{"WK17 keep not a boolean", "WK17", "E", map[string]func(string) string{a2b3: sub("priority: 2", "priority: 2\nkeep: maybe")}},
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
		{"S1 behind", "S1", "W", map[string]func(string) string{"docs/.layout": sub("layout: 0.3.0", "layout: 0.1.0")}},
		{"S1 additive", "S1", "W", map[string]func(string) string{"docs/.layout": sub("layout: 0.3.0", "layout: 0.2.0")}},
		{"S1 missing manifest", "S1", "W", map[string]func(string) string{"docs/.layout": nil}},
		{"S3 bad flag", "S3", "E", map[string]func(string) string{"docs/.layout": func(s string) string { return s + "git_commits: maybe\n" }}},
		{"S4 on without the skill", "S4", "W", map[string]func(string) string{"docs/.layout": func(s string) string { return s + "git_commits: true\n" }, ".git/HEAD": func(string) string { return "ref: refs/heads/main\n" }}},
		{"S4 on outside a work tree", "S4", "W", map[string]func(string) string{"docs/.layout": func(s string) string { return s + "git_commits: true\n" }, ".agents/skills/akm-commit/SKILL.md": func(string) string { return "---\nname: akm-commit\ndescription: x\n---\n" }}},
		{"S4 off with the skill", "S4", "W", map[string]func(string) string{".agents/skills/akm-commit/SKILL.md": func(string) string { return "---\nname: akm-commit\ndescription: x\n---\n" }}},
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

func TestCommitSkillAgrees(t *testing.T) {
	skill := func(string) string { return "---\nname: akm-commit\ndescription: x\n---\n" }
	on := func(s string) string { return s + "git_commits: true\n" }
	for name, edits := range map[string]map[string]func(string) string{
		"absent and absent": {},
		"false and absent":  {"docs/.layout": func(s string) string { return s + "git_commits: false\n" }},
		"true and present":  {"docs/.layout": on, ".agents/skills/akm-commit/SKILL.md": skill, ".git/HEAD": func(string) string { return "x\n" }},
	} {
		if got := rules(Run(repoWith(t, edits), Options{AsOf: asOf})); got["S4"] != 0 || got["S3"] != 0 {
			t.Errorf("%s: %v", name, got)
		}
	}
}

func TestClaims(t *testing.T) {
	r := repoWith(t, map[string]func(string) string{"docs/.layout": func(s string) string { return s + "extension.claims: 0.3.0\n" }})
	res := Run(r, Options{AsOf: asOf})
	if rules(res)["CL1"] != 2 {
		t.Fatalf("%v", res.Findings)
	}
	r = repoWith(t, map[string]func(string) string{
		"docs/.layout": func(s string) string { return s + "extension.claims: 0.3.0\n" },
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

func vaultRepo(t *testing.T, edits map[string]func(string) string) *record.Repo {
	t.Helper()
	src := "../../testdata/api/repos/vault"
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
		full := filepath.Join(dst, filepath.FromSlash(p))
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

func TestCustomFolders(t *testing.T) {
	const w = "Deep Thought/9 pkm/tasks/"
	if res := Run(vaultRepo(t, nil), Options{AsOf: asOf}); len(res.Findings) != 0 || res.State != "current" {
		t.Fatalf("%+v", res)
	}
	cases := []struct {
		name, rule string
		edits      map[string]func(string) string
		msg        string
	}{
		{"WK9 names the done folder", "WK9", map[string]func(string) string{w + "261022-e5f6-document-the-exit-codes.md": sub("status: ready", "status: dropped")}, "belongs in Deep Thought/9 pkm/tasks/done/"},
		{"WK1 for an unreadable task", "WK1", map[string]func(string) string{w + "261022-e5f6-document-the-exit-codes.md": sub("status: ready", "status: ready\nstatus: raw")}, ""},
		{"A1 for an unreadable ADR", "A1", map[string]func(string) string{"Deep Thought/9 pkm/adr/261012-p3q4-store-records-as-markdown-files.md": sub("status: accepted", "status: accepted\nstatus: proposed")}, ""},
		{"A4 index in the ADR folder", "A4", map[string]func(string) string{"Deep Thought/9 pkm/adr/README.md": func(s string) string { return s + "x\n" }}, ""},
		{"RM1 in the roadmap file", "RM1", map[string]func(string) string{"Deep Thought/9 pkm/ROADMAP.md": sub("## Later\n", "## Later\n- [261099-bbbb] Ghost\n")}, ""},
		{"L1 reaches the work folder", "L1", map[string]func(string) string{w + "261022-e5f6-document-the-exit-codes.md": func(s string) string { return s + "\n[gone](missing.md)\n" }}, ""},
		{"S3 for a bad folder", "S3", map[string]func(string) string{"docs/.layout": sub("work_dir: Deep Thought/9 pkm/tasks", "work_dir: ../tasks")}, `".."`},
	}
	for _, c := range cases {
		res := Run(vaultRepo(t, c.edits), Options{AsOf: asOf})
		found := false
		for _, f := range res.Findings {
			if f.Rule == c.rule && strings.Contains(f.Message, c.msg) {
				found = true
				if f.File == "" || strings.HasPrefix(f.File, "docs/work") || strings.HasPrefix(f.File, "docs/adr") {
					t.Errorf("%s: file %q", c.name, f.File)
				}
			}
		}
		if !found {
			t.Errorf("%s: no %s finding in %+v", c.name, c.rule, res.Findings)
		}
	}
}

func TestBadManifestPathIsAnError(t *testing.T) {
	res := Run(repoWith(t, map[string]func(string) string{"docs/.layout": func(s string) string { return s + "adr_dir: /abs\n" }}), Options{AsOf: asOf})
	for _, f := range res.Findings {
		if f.Rule == "S3" && f.Severity == "E" && f.File == "docs/.layout" && f.Line > 0 {
			return
		}
	}
	t.Fatalf("no S3 error: %+v", res.Findings)
}

// Obsidian escapes the alias pipe inside a table cell as \|, so [[name\|alias]] names the note "name".
func TestL2EscapedPipeInTableCell(t *testing.T) {
	edit := map[string]func(string) string{"AGENTS.md": func(s string) string {
		return s + "\n| a | b |\n| - | - |\n| [[AGENTS\\|the agents file]] | x |\n"
	}}
	for _, f := range Run(repoWith(t, edit), Options{AsOf: asOf}).Findings {
		if f.Rule == "L2" {
			t.Fatalf("%+v", f)
		}
	}
}
