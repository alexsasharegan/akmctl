package migrate

import (
	"path"
	"strings"
	"testing"
)

func TestMintIDsSlugFromTitle(t *testing.T) {
	files := map[string]string{
		"adr/019.md": "---\ntitle: Use a local server for the workbench UI\ncreated: 2026-03-04\n---\n# Heading\n",
		"adr/20.md":  "---\ncreated: 2026-03-04\n---\n# Store records as plain Markdown files\n",
		"adr/21.md":  "---\ntitle: A very long decision title that keeps going well past the limit\ncreated: 2026-03-04\n---\n",
		"adr/22.md":  "---\ncreated: 2026-03-04\n---\n",
	}
	run := func(from string, max int) []string {
		repo := filesRepo(t, files)
		s := Step{ID: "m", Kind: "mint_ids", Class: "auto", Dir: "adr", Seed: "slugs", SlugFrom: from, MaxSlug: max}
		res, err := ApplyEnv(repo, []Step{s}, Env{})
		if err != nil {
			t.Fatal(err)
		}
		var slugs []string
		for _, f := range res[0].Files {
			b := path.Base(f)
			slugs = append(slugs, strings.TrimSuffix(b[12:], ".md"))
		}
		return slugs
	}
	// default: the numeric file name is the slug
	if got := strings.Join(run("", 0), ","); !strings.Contains(got, "019,") || !strings.Contains(got, ",20") {
		t.Fatalf("name slugs %s", got)
	}
	// title: numeric names give way to the title (or the first heading)
	got := run("title", 0)
	want := []string{ // in path order of the files minted from 019, 20, 21 and 22 as sorted by new name
		"a-very-long-decision-title-that-keeps-going-well-past-the", // the 60-character default limit ends on a word
		"use-a-local-server-for-the-workbench-ui",
		"store-records-as-plain-markdown-files",
		"22", // no title at all: the name is the fallback
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("title slugs\n got %v\nwant %v", got, want)
	}
	// a length limit ends on a whole word
	got = run("title", 24)
	for _, g := range got {
		if len(g) > 24 || strings.HasSuffix(g, "-") {
			t.Errorf("slug %q breaks the limit", g)
		}
	}
	if !strings.Contains(strings.Join(got, ","), "use-a-local-server-for") || !strings.Contains(strings.Join(got, ","), "a-very-long-decision") {
		t.Fatalf("limited slugs %v", got)
	}
	// the limit applies to the name too
	if g := run("", 2); len(g[0]) > 2 {
		t.Fatalf("name slugs %v", g)
	}
}

func TestMintIDsSlugOptionsAreValidated(t *testing.T) {
	step := func(s Step) error {
		s.ID, s.Kind, s.Class, s.Dir = "m", "mint_ids", "auto", "adr"
		return Validate(Migration{From: "0.9.0", To: "1.0.0", Steps: []Step{s}})
	}
	if err := step(Step{SlugFrom: "title", MaxSlug: 40}); err != nil {
		t.Fatal(err)
	}
	if step(Step{SlugFrom: "body"}) == nil || step(Step{MaxSlug: 500}) == nil || step(Step{MaxSlug: -1}) == nil {
		t.Fatal("bad options accepted")
	}
}
