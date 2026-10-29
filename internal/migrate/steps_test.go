package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alexsasharegan/akmctl/internal/record"
)

const (
	tasks = "Vault Notes/9 pkm/tasks"
	adrs  = "Vault Notes/9 pkm/adr"
)

// filesRepo writes files (path to text) under a temp dir and loads it.
func filesRepo(t *testing.T, files map[string]string) *record.Repo {
	t.Helper()
	root := t.TempDir()
	for p, text := range files {
		full := filepath.Join(root, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	r, err := record.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func readRel(t *testing.T, repo *record.Repo, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(repo.Root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func exists(repo *record.Repo, rel string) bool {
	_, err := os.Stat(filepath.Join(repo.Root, filepath.FromSlash(rel)))
	return err == nil
}

func apply(t *testing.T, repo *record.Repo, s Step) Result {
	t.Helper()
	res, err := ApplyEnv(repo, []Step{s}, Env{Now: time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	return res[0]
}

func legacyFiles() map[string]string {
	return map[string]string{
		tasks + "/2026-10-20 Write the docs.md":   "---\ntype: task\nstatus: ready\ncreated_at: 2026-10-20\nupdated_at: 2026-10-21\nproject: \"[[Docs]]\"\n---\n# Write the docs\n\nIntro paragraph.\n\n## Log\n\n- 2026-10-20: filed.\n",
		tasks + "/2026-10-22 Already has both.md": "---\ntype: task\nstatus: raw\ncreated_at: 2026-10-22\ncreated: 2026-10-23\n---\n# Already has both\n",
		tasks + "/Plain note.md":                  "# No frontmatter\n\nText.\n",
		tasks + "/README.md":                      "---\ncreated_at: 2026-01-01\n---\nindex\n",
	}
}

func TestRenameFieldWithDir(t *testing.T) {
	repo := filesRepo(t, legacyFiles())
	s := Step{ID: "r", Kind: "rename_field", Class: "auto", Dir: tasks, From: "created_at", To: "created"}
	plan := PlanEnv(repo, []Step{s}, Env{})[0]
	if len(plan.Files) != 1 || plan.Files[0] != tasks+"/2026-10-20 Write the docs.md" {
		t.Fatalf("plan %v", plan.Files)
	}
	if strings.Contains(readRel(t, repo, plan.Files[0]), "\ncreated:") {
		t.Fatal("plan wrote")
	}
	res := apply(t, repo, s)
	if res.Status != "applied" || len(res.Files) != 1 {
		t.Fatalf("%+v", res)
	}
	text := readRel(t, repo, tasks+"/2026-10-20 Write the docs.md")
	if !strings.Contains(text, "created: 2026-10-20\n") || strings.Contains(text, "created_at") || !strings.Contains(text, "updated_at: 2026-10-21") {
		t.Fatal(text)
	}
	// the one that already has `created` is left alone; README and the no-frontmatter note too
	if !strings.Contains(readRel(t, repo, tasks+"/2026-10-22 Already has both.md"), "created_at: 2026-10-22") {
		t.Fatal("an existing target was overwritten")
	}
	if readRel(t, repo, tasks+"/Plain note.md") != legacyFiles()[tasks+"/Plain note.md"] {
		t.Fatal("a note without frontmatter changed")
	}
	if again := apply(t, repo, s); again.Status != "skipped" {
		t.Fatalf("second run: %+v", again)
	}
}

func TestRemoveAndCopyAndDefault(t *testing.T) {
	repo := filesRepo(t, legacyFiles())
	apply(t, repo, Step{ID: "c", Kind: "copy_field", Class: "auto", Dir: tasks, From: "project", To: "parent"})
	text := readRel(t, repo, tasks+"/2026-10-20 Write the docs.md")
	if !strings.Contains(text, "project: \"[[Docs]]\"\n") || !strings.Contains(text, "parent: \"[[Docs]]\"\n") {
		t.Fatal(text)
	}
	if res := apply(t, repo, Step{ID: "c", Kind: "copy_field", Class: "auto", Dir: tasks, From: "project", To: "parent"}); res.Status != "skipped" {
		t.Fatalf("copy is not idempotent: %+v", res)
	}
	// copy never overwrites an existing target
	apply(t, repo, Step{ID: "c2", Kind: "copy_field", Class: "auto", Dir: tasks, From: "status", To: "parent"})
	if !strings.Contains(readRel(t, repo, tasks+"/2026-10-20 Write the docs.md"), "parent: \"[[Docs]]\"") {
		t.Fatal("copy overwrote")
	}
	res := apply(t, repo, Step{ID: "x", Kind: "remove_field", Class: "auto", Dir: tasks, Field: "type"})
	if res.Status != "applied" || len(res.Files) != 2 {
		t.Fatalf("%+v", res)
	}
	if strings.Contains(readRel(t, repo, tasks+"/2026-10-20 Write the docs.md"), "type:") {
		t.Fatal("type still there")
	}
	if res := apply(t, repo, Step{ID: "x", Kind: "remove_field", Class: "auto", Dir: tasks, Field: "type"}); res.Status != "skipped" {
		t.Fatalf("remove is not idempotent: %+v", res)
	}
	apply(t, repo, Step{ID: "d", Kind: "set_default", Class: "auto", Dir: tasks, Field: "kind", Value: "chore"})
	if !strings.Contains(readRel(t, repo, tasks+"/2026-10-20 Write the docs.md"), "kind: chore") {
		t.Fatal("default not set in dir mode")
	}
}

func TestTitleFromH1(t *testing.T) {
	files := legacyFiles()
	files[tasks+"/Has title.md"] = "---\ntype: task\ntitle: Kept\n---\n# Other heading\n"
	files[tasks+"/No heading.md"] = "---\ntype: task\n---\nJust text.\n"
	files[tasks+"/Colon.md"] = "---\ntype: task\n---\n```\n# not a heading\n```\n\n# Fix: the thing # 2\n"
	repo := filesRepo(t, files)
	s := Step{ID: "t", Kind: "title_from_h1", Class: "auto", Dir: tasks}
	res := apply(t, repo, s)
	if res.Status != "applied" {
		t.Fatalf("%+v", res)
	}
	if got := readRel(t, repo, tasks+"/2026-10-20 Write the docs.md"); !strings.Contains(got, "title: Write the docs\n") {
		t.Fatal(got)
	}
	if got := readRel(t, repo, tasks+"/Has title.md"); !strings.Contains(got, "title: Kept\n") || strings.Contains(got, "Other heading\n---") {
		t.Fatal(got)
	}
	if got := readRel(t, repo, tasks+"/Colon.md"); !strings.Contains(got, `title: "Fix: the thing # 2"`) {
		t.Fatalf("a title that needs quoting: %s", got)
	}
	if got := readRel(t, repo, tasks+"/No heading.md"); strings.Contains(got, "title") {
		t.Fatal(got)
	}
	foundNote := false
	for _, n := range res.Notes {
		if strings.Contains(n, "No heading.md") {
			foundNote = true
		}
	}
	if !foundNote {
		t.Fatalf("notes %v", res.Notes)
	}
	if again := apply(t, repo, s); again.Status != "skipped" {
		t.Fatalf("second run: %+v", again.Files)
	}
}

func TestGoalScaffold(t *testing.T) {
	files := map[string]string{
		tasks + "/a.md": "---\ntype: task\nstatus: ready\n---\n# A\n\nWhy this matters.\n\nSecond paragraph.\n\n## Done when\n\n- [ ] one\n- [ ] two\n\n## Log\n\n- 2026-10-20: filed.\n",
		tasks + "/b.md": "---\ntype: task\nstatus: active\n---\n# B\n\n## Next\n\nGo.\n",
		tasks + "/c.md": "---\ntype: task\nstatus: ready\n---\n# C\n\n## Goal\n\nAlready there.\n",
		tasks + "/d.md": "---\ntype: task\nstatus: raw\n---\n# D\n\nText of a raw note.\n",
		tasks + "/e.md": "---\ntype: task\nstatus: ready\n---\n# E\n\nIntro only.\n",
		tasks + "/f.md": "---\ntype: task\nstatus: done\n---\n# F\n\n## Done when\n\n- [x] it\n\n## Outcome\n\nDone.\n",
		tasks + "/g.md": "---\ntype: task\nstatus: ready\n---\r\n# G\r\n\r\nCRLF intro.\r\n\r\n## Log\r\n\r\n- x\r\n",
	}
	repo := filesRepo(t, files)
	s := Step{ID: "g", Kind: "goal_scaffold", Class: "auto", Dir: tasks, Statuses: []string{"ready", "active", "blocked", "done"}}
	res := apply(t, repo, s)
	if res.Status != "applied" {
		t.Fatalf("%+v", res)
	}
	want := map[string]string{
		"a.md": "---\ntype: task\nstatus: ready\n---\n# A\n\n## Goal\n\nWhy this matters.\n\nSecond paragraph.\n\nDone when:\n\n- [ ] one\n- [ ] two\n\n## Log\n\n- 2026-10-20: filed.\n",
		"b.md": "---\ntype: task\nstatus: active\n---\n# B\n\n## Goal\n\n## Next\n\nGo.\n",
		"c.md": files[tasks+"/c.md"],
		"d.md": files[tasks+"/d.md"], // raw is outside statuses
		"e.md": "---\ntype: task\nstatus: ready\n---\n# E\n\n## Goal\n\nIntro only.\n",
		"f.md": "---\ntype: task\nstatus: done\n---\n# F\n\n## Goal\n\nDone when:\n\n- [x] it\n\n## Outcome\n\nDone.\n",
	}
	for name, w := range want {
		if got := readRel(t, repo, tasks+"/"+name); got != w {
			t.Errorf("%s:\n got %q\nwant %q", name, got, w)
		}
	}
	if got := readRel(t, repo, tasks+"/g.md"); !strings.Contains(got, "# G\r\n\r\n## Goal\r\n\r\nCRLF intro.\r\n\r\n## Log\r\n") {
		t.Errorf("crlf: %q", got)
	}
	// one note per changed file, naming it
	if len(res.Notes) != 5 {
		t.Fatalf("notes %v", res.Notes)
	}
	for _, n := range res.Notes {
		if !strings.Contains(n, "review") {
			t.Errorf("note %q does not ask for a review", n)
		}
	}
	if again := apply(t, repo, s); again.Status != "skipped" || len(again.Notes) != 0 {
		t.Fatalf("second run: %+v", again)
	}
}

func TestGoalScaffoldWithoutHeading(t *testing.T) {
	out, how, changed := scaffoldGoal("\nSome text with no heading.\n\n## Log\n\n- x\n")
	if !changed || out != "## Goal\n\nSome text with no heading.\n\n## Log\n\n- x\n" || how != "introduction moved under it" {
		t.Fatalf("%q %s", out, how)
	}
}

func TestStepsOnRecordsWithoutDir(t *testing.T) {
	repo := filesRepo(t, map[string]string{
		"docs/work/261020-a2b3-x.md": "---\nid: \"261020-a2b3\"\ntitle: X\nstatus: raw\ncreated_at: 2026-10-20\nupdated: 2026-10-20\n---\n# X\n",
	})
	res := apply(t, repo, Step{ID: "r", Kind: "rename_field", Class: "auto", Scope: "work", From: "created_at", To: "created"})
	if res.Status != "applied" || res.Files[0] != "docs/work/261020-a2b3-x.md" {
		t.Fatalf("%+v", res)
	}
}

func TestDirMustStayInside(t *testing.T) {
	repo := filesRepo(t, legacyFiles())
	for _, d := range []string{"../out", "/abs"} {
		_, err := ApplyEnv(repo, []Step{{ID: "r", Kind: "remove_field", Class: "auto", Dir: d, Field: "x"}}, Env{})
		if err == nil {
			t.Errorf("dir %q accepted", d)
		}
	}
	// a missing directory is a note, not an error
	res := apply(t, repo, Step{ID: "r", Kind: "remove_field", Class: "auto", Dir: "nowhere", Field: "x"})
	if res.Status != "skipped" || len(res.Notes) != 1 {
		t.Fatalf("%+v", res)
	}
}

// ---- mint_ids: rename false, seed, and the order with rename_field ----

func mintFiles() map[string]string {
	return map[string]string{
		tasks + "/2026-10-20 Write the docs.md": "---\ntype: task\ncreated: 2026-10-20\n---\n# Write the docs\n",
		tasks + "/Plain.md":                     "---\ntype: task\ncreated: 2026-09-01\n---\n# Plain\n",
		tasks + "/notes without a date.md":      "---\ntype: task\n---\n# n\n",
	}
}

func TestMintIDsNoRenameMatchesRename(t *testing.T) {
	no, yes := false, true
	base := Step{ID: "m", Kind: "mint_ids", Class: "auto", Dir: tasks, Seed: "vault-1"}
	env := Env{Now: time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)}
	a := filesRepo(t, mintFiles())
	b := filesRepo(t, mintFiles())

	nr := base
	nr.Rename = &no
	plan := PlanEnv(a, []Step{nr}, env)[0]
	if len(plan.Renames) != 3 {
		t.Fatalf("renames %v", plan.Renames)
	}
	// plan wrote nothing; the apply agrees with the plan
	if !exists(a, tasks+"/Plain.md") || strings.Contains(readRel(t, a, tasks+"/Plain.md"), "id:") {
		t.Fatal("plan wrote")
	}
	res, err := ApplyEnv(a, []Step{nr}, env)
	if err != nil {
		t.Fatal(err)
	}
	if len(res[0].Renames) != 3 {
		t.Fatalf("applied renames %v", res[0].Renames)
	}
	for i, r := range plan.Renames {
		if res[0].Renames[i] != r {
			t.Errorf("plan %v, apply %v", r, res[0].Renames[i])
		}
	}
	// no file moved; the id is written
	for _, r := range res[0].Renames {
		if !exists(a, r.From) || exists(a, r.To) {
			t.Errorf("%s moved or %s exists", r.From, r.To)
		}
		if id := filepath.Base(r.To)[:11]; !strings.Contains(readRel(t, a, r.From), "id: \""+id+"\"\n") {
			t.Errorf("%s has no id %s", r.From, id)
		}
		if !strings.Contains(strings.Join(res[0].Notes, "\n"), r.From+" -> "+r.To) || !strings.Contains(strings.Join(res[0].Notes, "\n"), "not renamed") {
			t.Errorf("notes %v", res[0].Notes)
		}
	}
	// a second run changes nothing
	if again, _ := ApplyEnv(a, []Step{nr}, env); again[0].Status != "skipped" || len(again[0].Renames) != 0 {
		t.Fatalf("second run: %+v", again[0])
	}

	// the names are exactly what rename: true produces for the same seed
	yr := base
	yr.Rename = &yes
	got, err := ApplyEnv(b, []Step{yr}, env)
	if err != nil {
		t.Fatal(err)
	}
	if len(got[0].Renames) != 3 {
		t.Fatalf("%v", got[0].Renames)
	}
	for i, r := range got[0].Renames {
		if r != res[0].Renames[i] {
			t.Errorf("rename: true gave %v, rename: false planned %v", r, res[0].Renames[i])
		}
		if !exists(b, r.To) || exists(b, r.From) {
			t.Errorf("%v not moved", r)
		}
	}
	// the default is to rename
	c := filesRepo(t, mintFiles())
	d := base
	if r, _ := ApplyEnv(c, []Step{d}, env); !exists(c, r[0].Renames[0].To) {
		t.Fatal("default must rename")
	}
	// a different seed gives different ids
	e := filesRepo(t, mintFiles())
	d2 := nr
	d2.Seed = "vault-2"
	other := PlanEnv(e, []Step{d2}, env)[0]
	if other.Renames[0] == plan.Renames[0] {
		t.Fatal("the seed made no difference")
	}
}

func TestMintIDsSeedRedrawsOnCollision(t *testing.T) {
	no := false
	s := Step{ID: "m", Kind: "mint_ids", Class: "auto", Dir: tasks, Seed: "s", Rename: &no}
	repo := filesRepo(t, mintFiles())
	m := &mintCtx{repo: repo, root: repo.Root, used: map[string]bool{}, seed: "s"}
	first, _ := m.mint("261020", "x.md")
	m2 := &mintCtx{repo: repo, root: repo.Root, used: map[string]bool{first: true}, seed: "s"}
	second, _ := m2.mint("261020", "x.md")
	if first == second || !strings.HasPrefix(second, "261020-") {
		t.Fatalf("%s %s", first, second)
	}
	_ = s
}

func TestADRsWithNumericIdsRenameFirst(t *testing.T) {
	files := map[string]string{
		adrs + "/ADR-001 Use SQLite.md": "---\nid: 1\ntype: adr\nstatus: accepted\ndate: 2026-10-01\n---\n# Use SQLite\n",
		adrs + "/ADR-002 Use Go.md":     "---\nid: 2\ntype: adr\nstatus: accepted\ndate: 2026-10-02\n---\n# Use Go\n",
	}
	no := false
	rename := Step{ID: "n", Kind: "rename_field", Class: "auto", Dir: adrs, From: "id", To: "adr_number"}
	mint := Step{ID: "m", Kind: "mint_ids", Class: "auto", Dir: adrs, DateField: "date", Seed: "adr", Rename: &no}
	// minted without the rename first: the numeric id counts as an id, so nothing is minted
	repo := filesRepo(t, files)
	if r, _ := ApplyEnv(repo, []Step{mint}, Env{}); r[0].Status != "skipped" {
		t.Fatalf("mint_ids must not touch a file that has an id field: %+v", r[0])
	}
	// in the documented order, both get ids from `date`
	repo = filesRepo(t, files)
	res, err := ApplyEnv(repo, []Step{rename, mint}, Env{})
	if err != nil {
		t.Fatal(err)
	}
	if res[0].Status != "applied" || res[1].Status != "applied" || len(res[1].Renames) != 2 {
		t.Fatalf("%+v", res)
	}
	text := readRel(t, repo, adrs+"/ADR-001 Use SQLite.md")
	if !strings.HasPrefix(text, "---\nid: \"261001-") || !strings.Contains(text, "adr_number: 1\n") || strings.Contains(text, "\nid: 1") {
		t.Fatal(text)
	}
	if !strings.HasPrefix(res[1].Renames[0].To, adrs+"/261001-") || !strings.HasSuffix(res[1].Renames[0].To, "-adr-001-use-sqlite.md") {
		t.Fatalf("%v", res[1].Renames)
	}
}

func TestMintIDsUsesTheManifestFolders(t *testing.T) {
	// an id already used in the configured work folder is never drawn again
	files := map[string]string{
		"docs/.layout": "work_dir: Vault Notes/tasks\nadr_dir: Vault Notes/adr\n",
		"Vault Notes/tasks/261020-aaaa-existing.md": "---\nid: \"261020-aaaa\"\n---\n# e\n",
		"Vault Notes/inbox/2026-10-20 New.md":       "---\ncreated: 2026-10-20\n---\n# New\n",
	}
	repo := filesRepo(t, files)
	// groupReader(10,...) draws "aaaa" (alphabet index 10 is "a"), then "bbbb"
	env := Env{Rand: groupReader(10, 11)}
	res := PlanEnv(repo, []Step{{ID: "m", Kind: "mint_ids", Class: "auto", Dir: "Vault Notes/inbox"}}, env)[0]
	if len(res.Renames) != 1 || res.Renames[0].To != "Vault Notes/inbox/261020-bbbb-new.md" {
		t.Fatalf("%+v", res.Renames)
	}
}
