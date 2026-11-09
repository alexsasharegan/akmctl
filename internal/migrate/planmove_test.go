package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A move_file before a mint_ids on the same directory: the plan must apply the
// move on its scratch copy, so the moved file is not counted as an id to mint.
func TestPlanSimulatesMoveFileBeforeMintIDs(t *testing.T) {
	repo := filesRepo(t, map[string]string{
		"docs/.layout":                "layout: 0.1.0\n",
		"docs/notes/Alpha.md":         "---\ntitle: Alpha\n---\n# Alpha\n",
		"docs/notes/Candidates.md":    "---\ntitle: Candidates\n---\n# Candidates\n",
		"docs/notes/raw-export.txt":   "plain text\n",
		"docs/notes/already/Other.md": "---\ntitle: Other\n---\n# Other\n",
	})
	steps := []Step{
		{ID: "move-md", Kind: "move_file", Class: "auto", From: "docs/notes/Candidates.md", To: "docs/candidates/Candidates.md"},
		{ID: "move-txt", Kind: "move_file", Class: "auto", From: "docs/notes/raw-export.txt", To: "docs/raw/raw-export.txt"},
		{ID: "add", Kind: "add_file", Class: "auto", Path: "docs/notes/Alpha.md", Content: "never overwritten\n"},
		{ID: "mint", Kind: "mint_ids", Class: "auto", Dir: "docs/notes"},
	}
	env := mintEnv()
	plan := PlanEnv(repo, steps, env)
	if len(plan[0].Files) != 1 || len(plan[1].Files) != 1 {
		t.Fatalf("moves planned %v, %v", plan[0].Files, plan[1].Files)
	}
	if len(plan[2].Files) != 0 {
		t.Fatalf("add_file over an existing file planned %v", plan[2].Files)
	}
	if exists(repo, "docs/candidates/Candidates.md") || !exists(repo, "docs/notes/raw-export.txt") {
		t.Fatal("plan moved a file")
	}
	got, err := ApplyEnv(repo, steps, mintEnv())
	if err != nil {
		t.Fatal(err)
	}
	for i := range steps {
		if strings.Join(plan[i].Files, "|") != strings.Join(got[i].Files, "|") {
			t.Errorf("step %s: plan %v, apply %v", steps[i].ID, plan[i].Files, got[i].Files)
		}
		if len(plan[i].Renames) != len(got[i].Renames) {
			t.Errorf("step %s: plan %d renames, apply %d", steps[i].ID, len(plan[i].Renames), len(got[i].Renames))
		}
	}
	for _, f := range plan[3].Files {
		if strings.Contains(f, "ndidates") {
			t.Errorf("the moved file was counted: %s", f)
		}
	}
	if len(plan[3].Renames) != 1 {
		t.Fatalf("mint planned %d renames, want 1 (Alpha; the directory is not recursive)", len(plan[3].Renames))
	}
}

// On the vault fixture, with a note moved out of the tasks folder before ids are
// minted there, the plan and apply agree on the files and the renames.
func TestPlanAgreesWithApplyOnTheVaultFixture(t *testing.T) {
	repo := vaultCopy(t, nil)
	for name, text := range map[string]string{
		"Notes/tasks/Candidates.md": "---\ntitle: Candidates\n---\n# Candidates\n",
		"Notes/tasks/Loose note.md": "---\ntitle: Loose note\n---\n# Loose note\n",
	} {
		full := filepath.Join(repo.Root, filepath.FromSlash(name))
		if err := os.WriteFile(full, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := repo.Reload(); err != nil {
		t.Fatal(err)
	}
	steps := []Step{
		{ID: "move", Kind: "move_file", Class: "auto", From: "Notes/tasks/Candidates.md", To: "Notes/Candidates.md"},
		{ID: "mint", Kind: "mint_ids", Class: "auto", Dir: "Notes/tasks"},
	}
	plan := PlanEnv(repo, steps, mintEnv())
	got, err := ApplyEnv(repo, steps, mintEnv())
	if err != nil {
		t.Fatal(err)
	}
	if len(got[1].Renames) != 1 {
		t.Fatalf("apply renamed %d, want 1", len(got[1].Renames))
	}
	for i := range steps {
		if strings.Join(plan[i].Files, "|") != strings.Join(got[i].Files, "|") || len(plan[i].Renames) != len(got[i].Renames) {
			t.Errorf("step %s: plan %v (%d renames), apply %v (%d renames)", steps[i].ID, plan[i].Files, len(plan[i].Renames), got[i].Files, len(got[i].Renames))
		}
	}
}
