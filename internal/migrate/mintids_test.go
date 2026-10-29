package migrate

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alexsasharegan/akmctl/internal/record"
)

const mintFixture = "../../testdata/migrate/mint_ids"

// groupReader gives the random source four equal bytes per draw, so a draw of
// n yields the suffix of the alphabet's n-th character four times.
func groupReader(groups ...byte) *bytes.Reader {
	var b []byte
	for _, g := range groups {
		b = append(b, g, g, g, g)
	}
	return bytes.NewReader(b)
}

func copyTree(t *testing.T, src, dst string) {
	t.Helper()
	err := filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dst, rel), b, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func mintRepo(t *testing.T) *record.Repo {
	t.Helper()
	root := t.TempDir()
	copyTree(t, mintFixture, root)
	repo, err := record.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

func snapshotTree(t *testing.T, root string) map[string]string {
	t.Helper()
	m := map[string]string{}
	filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			b, _ := os.ReadFile(p)
			rel, _ := filepath.Rel(root, p)
			m[filepath.ToSlash(rel)] = string(b)
		}
		return nil
	})
	return m
}

var mintStep = Step{ID: "mint-notes", Kind: "mint_ids", Class: "auto", Dir: "docs/notes"}

var mintEnv = func() Env {
	// draws: needs-slugging 0000; alpha 0000 (taken by a record in done/), 1111 (taken by an ADR), 2222;
	// beta 2222 (another date, free); gamma 3333; no-frontmatter 3333 (taken in this batch), 4444
	return Env{Rand: groupReader(0, 0, 1, 2, 2, 3, 3, 4), Now: time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)}
}

func TestMintIDsApply(t *testing.T) {
	repo := mintRepo(t)
	root := repo.Root
	before := snapshotTree(t, root)
	res, err := ApplyEnv(repo, []Step{mintStep}, mintEnv())
	if err != nil {
		t.Fatal(err)
	}
	if res[0].Status != "applied" {
		t.Fatalf("status %s", res[0].Status)
	}
	wantFiles := []string{
		"docs/notes/260904-2222-beta.md",
		"docs/notes/261001-0000-needs-slugging-co.md",
		"docs/notes/261005-9999-name-carries-id.md",
		"docs/notes/261022-2222-alpha-note.md",
		"docs/notes/261028-3333-gamma.md",
		"docs/notes/261028-4444-no-frontmatter.md",
	}
	if strings.Join(res[0].Files, "\n") != strings.Join(wantFiles, "\n") {
		t.Fatalf("files:\n%s", strings.Join(res[0].Files, "\n"))
	}
	after := snapshotTree(t, root)
	// the old names are gone by rename, never by deletion: the text moved
	for old, id := range map[string]string{
		"docs/notes/alpha-note.md":                      "docs/notes/261022-2222-alpha-note.md",
		"docs/notes/beta.md":                            "docs/notes/260904-2222-beta.md",
		"docs/notes/gamma.md":                           "docs/notes/261028-3333-gamma.md",
		"docs/notes/2026-10-01 Needs Slugging & Co!.md": "docs/notes/261001-0000-needs-slugging-co.md",
		"docs/notes/no-frontmatter.md":                  "docs/notes/261028-4444-no-frontmatter.md",
		"docs/notes/261005-9999-name-carries-id.md":     "docs/notes/261005-9999-name-carries-id.md",
	} {
		b, ok := after[id]
		if !ok {
			t.Errorf("%s missing", id)
			continue
		}
		if old != id {
			if _, still := after[old]; still {
				t.Errorf("%s still exists", old)
			}
		}
		orig := before[old]
		if i := strings.Index(orig, "\n---\n"); strings.HasPrefix(orig, "---\n") && i >= 0 {
			orig = orig[i:]
		}
		if !strings.HasSuffix(b, orig) {
			t.Errorf("%s lost its body", id)
		}
	}
	// the id is the first field, quoted; the rest of the file is untouched
	if got, want := after["docs/notes/261022-2222-alpha-note.md"], "---\nid: \"261022-2222\"\n"+before["docs/notes/alpha-note.md"][4:]; got != want {
		t.Errorf("alpha:\n%s\nwant:\n%s", got, want)
	}
	// a timestamp with an offset is read in UTC
	if !strings.HasPrefix(after["docs/notes/260904-2222-beta.md"], "---\nid: \"260904-2222\"\n") {
		t.Errorf("beta:\n%s", after["docs/notes/260904-2222-beta.md"])
	}
	// no frontmatter: a block is created
	if got, want := after["docs/notes/261028-4444-no-frontmatter.md"], "---\nid: \"261028-4444\"\n---\n\n"+before["docs/notes/no-frontmatter.md"]; got != want {
		t.Errorf("no-frontmatter:\n%s", got)
	}
	// a name that carries an id keeps it
	if got, want := after["docs/notes/261005-9999-name-carries-id.md"], "---\nid: \"261005-9999\"\n"+before["docs/notes/261005-9999-name-carries-id.md"][4:]; got != want {
		t.Errorf("name-id:\n%s", got)
	}
	// untouched: minted, bare id, README, hidden, subdirectory, unterminated, other records, the referrer
	for _, p := range []string{
		"docs/notes/261003-7777-already-minted.md", "docs/notes/minted-bare.md", "docs/notes/README.md",
		"docs/notes/.hidden.md", "docs/notes/sub/inner.md", "docs/notes/unterminated.md", "docs/overview.md",
		"docs/adr/261022-1111-also-taken.md", "docs/work/done/261022-0000-holds-a-drawn-suffix.md", "docs/.layout",
	} {
		if after[p] != before[p] {
			t.Errorf("%s changed", p)
		}
	}
	if len(after) != len(before) {
		t.Errorf("file count %d, was %d", len(after), len(before))
	}
	// notes: date source, creation of frontmatter, the skipped file, the references
	notes := strings.Join(res[0].Notes, "\n")
	for _, want := range []string{
		"docs/notes/gamma.md -> docs/notes/261028-3333-gamma.md (id 261028-3333, date today (created absent or not a date))",
		"docs/notes/alpha-note.md -> docs/notes/261022-2222-alpha-note.md (id 261022-2222, date from created)",
		"docs/notes/no-frontmatter.md -> docs/notes/261028-4444-no-frontmatter.md",
		"the file had no frontmatter, created it",
		"docs/notes/261005-9999-name-carries-id.md: name already starts with id 261005-9999; wrote the id field, no rename",
		"docs/notes/unterminated.md: skipped, frontmatter cannot be read",
		"reference not rewritten: docs/overview.md line 3 names alpha-note.md (now 261022-2222-alpha-note.md)",
		"reference not rewritten: docs/overview.md line 3 names beta.md (now 260904-2222-beta.md)",
		"reference not rewritten: docs/overview.md line 4 names gamma.md (now 261028-3333-gamma.md)",
	} {
		if !strings.Contains(notes, want) {
			t.Errorf("notes lack %q:\n%s", want, notes)
		}
	}
	if strings.Contains(notes, "line 5") || strings.Contains(notes, "inner") {
		t.Errorf("a near miss was reported:\n%s", notes)
	}

	// the result loads: ids match file names and are unique
	re, err := record.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, r := range re.Records {
		if seen[r.ID] {
			t.Errorf("duplicate id %s", r.ID)
		}
		seen[r.ID] = true
	}

	// a second run changes nothing
	snap := snapshotTree(t, root)
	res2, err := ApplyEnv(re, []Step{mintStep}, Env{Rand: groupReader(9, 9, 9, 9, 9, 9, 9, 9)})
	if err != nil {
		t.Fatal(err)
	}
	if res2[0].Status != "skipped" || len(res2[0].Files) != 0 {
		t.Fatalf("second run: %s %v", res2[0].Status, res2[0].Files)
	}
	if !equalMaps(snap, snapshotTree(t, root)) {
		t.Fatal("the second run changed files")
	}
}

func equalMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if w, ok := b[k]; !ok || v != w {
			return false
		}
	}
	return true
}

func TestMintIDsPlanWritesNothing(t *testing.T) {
	repo := mintRepo(t)
	before := snapshotTree(t, repo.Root)
	res := PlanEnv(repo, []Step{mintStep}, mintEnv())
	if len(res[0].Files) != 6 || res[0].Status != "pending" {
		t.Fatalf("plan %v %s", res[0].Files, res[0].Status)
	}
	if !equalMaps(before, snapshotTree(t, repo.Root)) {
		t.Fatal("plan wrote")
	}
	// the plan and the run agree when the source agrees
	res2, err := ApplyEnv(repo, []Step{mintStep}, mintEnv())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(res[0].Files, ",") != strings.Join(res2[0].Files, ",") {
		t.Fatalf("plan %v, run %v", res[0].Files, res2[0].Files)
	}
}

func TestMintIDsDateField(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "docs/notes"), 0o755)
	os.WriteFile(filepath.Join(root, "docs/notes/a.md"), []byte("---\ntitle: A\nwritten: 2025-03-04\ncreated: 2026-01-01\n---\n"), 0o644)
	os.WriteFile(filepath.Join(root, "docs/notes/b.md"), []byte("---\ntitle: B\nwritten: soon\n---\n"), 0o644)
	repo, _ := record.Load(root)
	s := mintStep
	s.DateField = "written"
	res, err := ApplyEnv(repo, []Step{s}, Env{Rand: groupReader(5, 6), Now: time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	want := "docs/notes/250304-5555-a.md,docs/notes/261028-6666-b.md"
	if strings.Join(res[0].Files, ",") != want {
		t.Fatalf("files %v", res[0].Files)
	}
	if !strings.Contains(strings.Join(res[0].Notes, "\n"), "today (written absent or not a date)") {
		t.Fatalf("notes %v", res[0].Notes)
	}
}

func TestMintIDsRecursiveAndTitleSlug(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "docs/notes/sub"), 0o755)
	os.WriteFile(filepath.Join(root, "docs/notes/sub/日本語.md"), []byte("---\ntitle: Kanji Only Name\ncreated: 2026-02-03\n---\n"), 0o644)
	os.WriteFile(filepath.Join(root, "docs/notes/"+strings.Repeat("long-word-", 12)+".md"), []byte("---\ncreated: 2026-02-03\n---\n"), 0o644)
	repo, _ := record.Load(root)
	s := mintStep
	s.Recursive = true
	res, err := ApplyEnv(repo, []Step{s}, Env{Rand: groupReader(1, 2)})
	if err != nil {
		t.Fatal(err)
	}
	if len(res[0].Files) != 2 {
		t.Fatalf("files %v", res[0].Files)
	}
	got := strings.Join(res[0].Files, ",")
	if !strings.Contains(got, "docs/notes/sub/260203-2222-kanji-only-name.md") && !strings.Contains(got, "docs/notes/sub/260203-1111-kanji-only-name.md") {
		t.Fatalf("title slug missing: %s", got)
	}
	for _, f := range res[0].Files {
		slug := strings.TrimSuffix(filepath.Base(f)[12:], ".md")
		if len(slug) > 60 || strings.HasSuffix(slug, "-") {
			t.Errorf("slug %q", slug)
		}
	}
}

func TestMintIDsRefusesBadDir(t *testing.T) {
	repo := mintRepo(t)
	for _, d := range []string{"", "../x", "/abs"} {
		s := mintStep
		s.Dir = d
		if _, err := ApplyEnv(repo, []Step{s}, Env{}); err == nil {
			t.Errorf("dir %q accepted", d)
		}
	}
	s := mintStep
	s.Dir = "docs/missing"
	res, err := ApplyEnv(repo, []Step{s}, Env{})
	if err != nil || res[0].Status != "skipped" {
		t.Fatalf("missing dir: %v %v", res, err)
	}
}

func TestSlugFor(t *testing.T) {
	for in, want := range map[string]string{
		"2026-10-01 Needs Slugging & Co!.md": "needs-slugging-co",
		"20261001_report.md":                 "report",
		"2026-10-01.md":                      "t",
		"Hello_World.md":                     "hello-world",
		"!!!.md":                             "t",
	} {
		if got := slugFor(in, "T"); got != want {
			t.Errorf("slugFor(%q) = %q, want %q", in, got, want)
		}
	}
	if slugFor("!!!.md", "") != "untitled" {
		t.Error("untitled")
	}
}

func TestIDDate(t *testing.T) {
	for in, want := range map[string]string{
		"2026-10-22":                "261022",
		"2026-10-22T08:00":          "261022",
		"2026-10-22T08:00:00Z":      "261022",
		"2026-09-03T23:30:00-04:00": "260904",
		"2026-10-22 08:00":          "261022",
	} {
		if got, ok := idDate(in); !ok || got != want {
			t.Errorf("idDate(%q) = %q %v", in, got, ok)
		}
	}
	for _, in := range []string{"", "soon", "2026-13-01", "2026-10-22x", "22 Oct 2026"} {
		if _, ok := idDate(in); ok {
			t.Errorf("idDate(%q) accepted", in)
		}
	}
}
