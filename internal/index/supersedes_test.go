package index_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/check"
	"github.com/alexsasharegan/akmctl/internal/index"
	"github.com/alexsasharegan/akmctl/internal/record"
	"github.com/alexsasharegan/akmctl/internal/testutil"
)

// A row says which relation each id is (work item 261028-1k0e), and rule A4
// accepts the index of a repository with a superseded record.
func TestSupersedesCellNamesTheRelation(t *testing.T) {
	root := t.TempDir()
	testutil.CopyTree(t, "../../testdata/api/repos/basic", root)
	adr := filepath.Join(root, "docs/adr")
	old, new := "", ""
	es, _ := os.ReadDir(adr)
	for _, e := range es {
		if e.Name() == "README.md" {
			continue
		}
		p := filepath.Join(adr, e.Name())
		b, _ := os.ReadFile(p)
		switch s := string(b); {
		case strings.Contains(s, "status: accepted") && old == "":
			old = p
		case strings.Contains(s, "status: proposed") && new == "":
			new = p
		}
	}
	if old == "" || new == "" {
		t.Fatal("the fixture needs an accepted and a proposed ADR")
	}
	oldID, newID := filepath.Base(old)[:11], filepath.Base(new)[:11]
	edit := func(p, from, to string) {
		b, _ := os.ReadFile(p)
		if !strings.Contains(string(b), from) {
			t.Fatalf("%s has no %q", p, from)
		}
		os.WriteFile(p, []byte(strings.Replace(string(b), from, to, 1)), 0o644)
	}
	edit(old, "status: accepted", "status: superseded\nsuperseded_by: "+newID)
	edit(new, "status: proposed", "status: accepted\nsupersedes: "+oldID)
	repo, err := record.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	got := index.Generate(repo)
	for _, want := range []string{"| superseded by " + newID + " ", "| supersedes " + oldID + " "} {
		if !strings.Contains(got, want) {
			t.Fatalf("no %q in\n%s", want, got)
		}
	}
	os.WriteFile(filepath.Join(adr, "README.md"), []byte(got), 0o644)
	repo, _ = record.Load(root)
	for _, f := range check.Run(repo, check.Options{}).Findings {
		if f.Rule == "A4" {
			t.Fatalf("A4 on a generated index: %+v", f)
		}
	}
	if out := os.Getenv("AKM_INDEX_OUT"); out != "" { // for a manual prettier --check
		os.WriteFile(out, []byte(got), 0o644)
	}
}
