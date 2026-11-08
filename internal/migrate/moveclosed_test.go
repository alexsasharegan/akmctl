package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alexsasharegan/akmctl/internal/record"
	"github.com/alexsasharegan/akmctl/internal/testutil"
)

// vaultCopy loads a copy of the vault fixture with the named task files given
// a closed status by hand (what a migrated note can carry, and what `akmctl
// set` cannot write).
func vaultCopy(t *testing.T, closed map[string]string) *record.Repo {
	t.Helper()
	src := "../../testdata/api/repos/vault"
	dst := t.TempDir()
	testutil.CopyTree(t, src, dst)
	for name, st := range closed {
		var found bool
		filepath.WalkDir(dst, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() || filepath.Base(p) != name {
				return err
			}
			b, rerr := os.ReadFile(p)
			if rerr != nil {
				t.Fatal(rerr)
			}
			lines := strings.Split(string(b), "\n")
			for i, l := range lines {
				if strings.HasPrefix(l, "status: ") {
					lines[i] = "status: " + st
					break
				}
			}
			if werr := os.WriteFile(p, []byte(strings.Join(lines, "\n")), 0o644); werr != nil {
				t.Fatal(werr)
			}
			found = true
			return nil
		})
		if !found {
			t.Fatalf("no file %s in the vault fixture", name)
		}
	}
	r, err := record.Load(dst)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestMoveClosedOnTheVaultFixture(t *testing.T) {
	const dir = "Notes/tasks"
	repo := vaultCopy(t, map[string]string{
		"261022-e5f6-document-the-exit-codes.md": "done",
		"261026-m1n2-explore-a-plugin-system.md": "dropped",
	})
	s := Step{ID: "close", Kind: "move_closed", Class: "auto", Dir: dir}
	want := []string{dir + "/done/261022-e5f6-document-the-exit-codes.md", dir + "/done/261026-m1n2-explore-a-plugin-system.md"}
	plan := PlanEnv(repo, []Step{s}, Env{})[0]
	if strings.Join(plan.Files, "|") != strings.Join(want, "|") {
		t.Fatalf("plan %v", plan.Files)
	}
	if !exists(repo, dir+"/261022-e5f6-document-the-exit-codes.md") {
		t.Fatal("plan moved a file")
	}
	res := apply(t, repo, s)
	if res.Status != "applied" || strings.Join(res.Files, "|") != strings.Join(want, "|") {
		t.Fatalf("%+v", res)
	}
	for _, w := range want {
		if !exists(repo, w) {
			t.Errorf("missing %s", w)
		}
	}
	if exists(repo, dir+"/261022-e5f6-document-the-exit-codes.md") || !exists(repo, dir+"/261020-a2b3-add-json-output-to-ls.md") {
		t.Fatal("wrong files moved")
	}
	if !exists(repo, dir+"/done/261010-j9k0-choose-a-licence.md") {
		t.Fatal("record already in done/ was moved")
	}
	if again := apply(t, repo, s); again.Status != "skipped" {
		t.Fatalf("second run %+v", again)
	}
}

func TestMoveClosedOptionsAndLogSize(t *testing.T) {
	files := map[string]string{}
	for i := 0; i < 100; i++ {
		files[fmt.Sprintf("notes/note-%03d.md", i)] = "---\nstatus: done\n---\n# N\n"
	}
	files["docs/.layout"] = "layout: 0.1.0\n"
	files["notes/open.md"] = "---\nstatus: active\n---\n# Open\n"
	files["notes/README.md"] = "---\nstatus: done\n---\nindex\n"
	repo := filesRepo(t, files)
	s := Step{ID: "close", Kind: "move_closed", Class: "auto", Dir: "notes", To: "archive/closed"}
	res, err := ApplyEnv(repo, []Step{s}, Env{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res[0].Files) != 100 || !exists(repo, "notes/open.md") || !exists(repo, "notes/README.md") || !exists(repo, "archive/closed/note-042.md") {
		t.Fatalf("%d files moved", len(res[0].Files))
	}
	if err := AppendLog(repo.Root, repo.MigrateLogPath(), "0.1.0", "0.3.0", res, time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	log := readRel(t, repo, repo.MigrateLogPath())
	if strings.Count(log, "\n") != 1 || !strings.Contains(log, `"count":100`) || len(log) > 8000 {
		t.Fatalf("log of %d bytes, %d lines", len(log), strings.Count(log, "\n"))
	}
	bad := Migration{From: "0.1.0", To: "0.3.0", Steps: []Step{{ID: "a", Kind: "move_closed", Class: "auto"}}}
	if Validate(bad) == nil {
		t.Fatal("dir should be required")
	}
}
