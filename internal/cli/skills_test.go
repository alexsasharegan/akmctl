package cli

import (
	"strings"
	"testing"
)

func TestSkillsCommand(t *testing.T) {
	script(t, false, "")
	root := tree(t)
	if code, out, _ := cfg(root, "skills"); code != 0 || !strings.Contains(out, "missing") {
		t.Fatalf("status: %d %s", code, out)
	}
	// a write with no terminal and no --yes is refused
	if code, _, errb := cfg(root, "skills", "install"); code != 2 || !strings.Contains(errb, "--yes") {
		t.Fatalf("no terminal: %d %s", code, errb)
	}
	if has(root, ".agents/skills/akm/SKILL.md") {
		t.Fatal("refused write wrote")
	}
	if code, _, _ := cfg(root, "skills", "install", "--dry-run"); code != 0 || has(root, ".agents/skills/akm/SKILL.md") {
		t.Fatal("dry run wrote")
	}
	if code, _, errb := cfg(root, "skills", "install", "--yes"); code != 0 || !has(root, ".agents/skills/akm/.akm-skill.json") {
		t.Fatalf("install: %d %s", code, errb)
	}
	if code, out, _ := cfg(root, "skills", "update", "--yes"); code != 0 || strings.Contains(out, "replace") {
		t.Fatalf("update of a current repo: %d %s", code, out)
	}
	if code, _, _ := cfg(root, "skills", "bogus"); code != 2 {
		t.Fatal("unknown mode accepted")
	}
	// a terminal: answer no
	script(t, true, "n\n")
	root2 := tree(t)
	if code, _, _ := cfg(root2, "skills", "install"); code != 1 || has(root2, ".agents/skills/akm/SKILL.md") {
		t.Fatal("a no still wrote")
	}
}
