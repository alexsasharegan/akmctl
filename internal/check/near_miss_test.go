package check

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

type edits = map[string]func(string) string

func put(text string) func(string) string { return func(string) string { return text } }

func app(text string) func(string) string { return func(s string) string { return s + text } }

// padTo makes a file exactly n lines long.
func padTo(n int) func(string) string {
	return func(s string) string {
		s = strings.TrimRight(s, "\n") + "\n"
		for strings.Count(s, "\n") < n {
			s += "- x\n"
		}
		return s
	}
}

const (
	c4d5 = "docs/work/261015-c4d5-machine-readable-listings.md"
	m1n2 = "docs/work/261026-m1n2-explore-a-plugin-system.md"
)

const claimsOn = "extension.claims: 0.3.0\n"

// nearMisses are, for each rule, a repository that sits next to the case that
// raises it in TestRules and must not raise it (work item 261108-qv1w, finding
// F18). Each is a boundary or a look-alike, so a rule that grows too eager
// fails here.
var nearMisses = []struct {
	name  string
	rule  string
	edits edits
}{
	{"WK1 quoted id", "WK1", edits{e5f6: sub("id: 261022-e5f6", "id: \"261022-e5f6\"")}},
	{"WK2 quoted created", "WK2", edits{e5f6: sub("created: 2026-10-22", "created: \"2026-10-22\"")}},
	{"WK3 leap-day date", "WK3", edits{e5f6: sub("updated: 2026-10-22", "updated: 2028-02-29")}},
	{"WK4 parent is another open item", "WK4", edits{e5f6: sub("parent: 261015-c4d5", "parent: 261024-g7h8")}},
	{"WK5 blocked by free text only", "WK5", edits{g7h8: sub("blocked_by: [261020-a2b3, \"legal review\"]", "blocked_by: [\"legal review\"]")}},
	{"WK6 ready item without a Next", "WK6", nil},
	{"WK7 raw item without a Goal", "WK7", edits{m1n2: sub("## Notes", "## Aim")}},
	{"WK8 open parent with open children", "WK8", edits{c4d5: sub("status: active", "status: ready")}},
	{"WK9 dropped item in done", "WK9", edits{j9k0: sub("status: done", "status: dropped")}},
	{"WK10 exactly 30 days", "WK10", edits{a2b3: sub("updated: 2026-10-27", "updated: 2026-09-28")}},
	{"WK11 priority 0 and a 160-character summary", "WK11", edits{a2b3: sub("summary: Print ls results as JSON.", "summary: "+strings.Repeat("x", 160)), e5f6: sub("status: ready", "status: ready\npriority: 0")}},
	{"WK12 quoted title", "WK12", edits{e5f6: sub("title: Document the exit codes", "title: \"Document the exit codes\"")}},
	{"WK13 done item with closed", "WK13", nil},
	{"WK14 done item with closed", "WK14", nil},
	{"WK15 closed on the day created", "WK15", edits{j9k0: sub("closed: 2026-10-12", "closed: 2026-10-10")}},
	{"WK16 cycle of 64 characters", "WK16", edits{a2b3: sub("priority: 2", "priority: 2\ncycle: "+strings.Repeat("a", 64))}},
	{"WK17 keep false", "WK17", edits{a2b3: sub("priority: 2", "priority: 2\nkeep: false")}},
	{"WK18 one-way block", "WK18", edits{e5f6: sub("parent: 261015-c4d5", "parent: 261015-c4d5\nblocked_by: [261020-a2b3]")}},
	{"WK19 blocked_by a sibling", "WK19", edits{e5f6: sub("parent: 261015-c4d5", "parent: 261015-c4d5\nblocked_by: [261020-a2b3]")}},
	{"B1 indented text inside a fence", "B1", edits{e5f6: app("\n```\n    <div>\n\ttab\n```\n")}},
	{"I1 neighbouring id", "I1", edits{"docs/work/261022-e5f7-copy.md": put("---\nid: 261022-e5f7\ntitle: Copy\nstatus: raw\ncreated: 2026-10-22\nupdated: 2026-10-22\n---\n# Copy\n")}},
	{"A1 160-character summary", "A1", edits{p3q4: sub("summary: Records are Markdown files with frontmatter; read this before changing the format.", "summary: "+strings.Repeat("x", 160))}},
	{"A2 deprecated is a status", "A2", edits{r5s6: sub("status: proposed", "status: deprecated")}},
	{"A3 superseded with both links", "A3", edits{r5s6: sub("status: proposed", "status: superseded\nsuperseded_by: 261012-p3q4"), p3q4: sub("status: accepted", "status: accepted\nsupersedes: 261018-r5s6")}},
	{"A4 body edit leaves the index", "A4", edits{p3q4: app("\nMore context.\n")}},
	{"RM1 entry resolves", "RM1", edits{rm: sub("## Later\n", "## Later\n- [261026-m1n2] Explore\n")}},
	{"RM2 dropped section", "RM2", edits{j9k0: sub("status: done", "status: dropped"), rm: sub("## Dropped", "## Dropped\n- [261010-j9k0] Licence")}},
	{"RM3 dropped item under Dropped", "RM3", edits{j9k0: sub("status: done", "status: dropped"), rm: sub("## Dropped", "## Dropped\n- [261010-j9k0] Licence")}},
	{"RM4 raw top-level item is not required", "RM4", nil},
	{"G1 trailing spaces on a heading", "G1", edits{"AGENTS.md": sub("## Rules", "## Rules  ")}},
	{"G2 exactly 80 lines", "G2", edits{"AGENTS.md": padTo(80)}},
	{"G3 script that exists", "G3", edits{"AGENTS.md": sub("`akmctl check`", "`scripts/check.sh`"), "scripts/check.sh": put("#!/bin/sh\n")}},
	{"L1 link in a code span", "L1", edits{"AGENTS.md": app("\nUse `[x](nope.md)` here.\n")}},
	{"L2 wikilink in a code span", "L2", edits{"AGENTS.md": app("\nUse `[[Nowhere]]` here.\n")}},
	{"S1 current manifest with an extra key", "S1", edits{"docs/.layout": app("stale_days: 45\n")}},
	{"S2 name equals the folder", "S2", edits{".agents/skills/demo/SKILL.md": put("---\nname: demo\ndescription: x\n---\n")}},
	{"S3 git_commits false", "S3", edits{"docs/.layout": app("git_commits: false\n")}},
	{"S4 on with the skill in a work tree", "S4", edits{"docs/.layout": app("git_commits: true\n"), ".agents/skills/akm-commit/SKILL.md": put("---\nname: akm-commit\ndescription: x\n---\n"), ".git/HEAD": put("ref: refs/heads/main\n")}},
	{"CL1 active item claimed", "CL1", edits{"docs/.layout": app(claimsOn), a2b3: sub("status: active", "status: active\nclaimed_by: agent-1\nclaimed_at: 2026-10-27T10:00"), c4d5: sub("status: active", "status: active\nclaimed_by: agent-2\nclaimed_at: 2026-10-27T10:00")}},
	{"CL2 two holders, one item each", "CL2", edits{"docs/.layout": app(claimsOn), a2b3: sub("status: active", "status: active\nclaimed_by: agent-1\nclaimed_at: 2026-10-27T10:00"), c4d5: sub("status: active", "status: active\nclaimed_by: agent-2\nclaimed_at: 2026-10-27T10:00")}},
	{"SF2 deny list that matches nothing", "SF2", edits{"docs/.layout": app("extension.safety: 0.3.0\n"), "docs/.akm-deny": put("Zorblax\n")}},
}

// coveredElsewhere are the rules whose near misses live in a test of their own.
var coveredElsewhere = map[string]string{
	"SF1": "TestSF1IgnoresPlaceholders",
	"SF3": "TestSF3SkipsRecordNamesAndLinkTargets",
	"SF4": "TestSF4AllowMarkers (a marker with a reason is a warning, not an error)",
}

func TestNearMisses(t *testing.T) {
	for _, c := range nearMisses {
		t.Run(c.name, func(t *testing.T) {
			res := Run(repoWith(t, c.edits), Options{AsOf: asOf})
			for _, f := range res.Findings {
				if f.Rule == c.rule {
					t.Fatalf("%s must not be raised, got %+v", c.rule, f)
				}
			}
		})
	}
}

// TestEveryRuleHasANearMiss fails when a rule the tool raises has no near-miss
// case: add the rule to nearMisses (or to coveredElsewhere with the test).
func TestEveryRuleHasANearMiss(t *testing.T) {
	have := map[string]bool{}
	for _, c := range nearMisses {
		have[c.rule] = true
	}
	for r := range coveredElsewhere {
		have[r] = true
	}
	var missing []string
	for _, r := range ruleIDs() {
		if !have[r] {
			missing = append(missing, r)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Fatalf("no near-miss case for %v", missing)
	}
}

// ruleIDs reads the rule ids from the table in the specification, the one list
// of what the tool is meant to check.
func ruleIDs() []string {
	b, err := os.ReadFile("../../docs/spec/agent-ready-repo-spec.md")
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(`(?m)^\| ([A-Z]{1,2}[0-9]{1,2}) +\| (?:E|W|W/E) +\|`)
	var ids []string
	for _, m := range re.FindAllStringSubmatch(string(b), -1) {
		ids = append(ids, m[1])
	}
	if len(ids) < 40 {
		panic("the rule table of the specification was not found")
	}
	return ids
}
