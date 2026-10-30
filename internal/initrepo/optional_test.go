package initrepo_test

import (
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/initrepo"
)

func TestOptionalSkill(t *testing.T) {
	if !initrepo.IsOptionalSkill("akm-commit") || initrepo.IsOptionalSkill("work-items") || initrepo.IsOptionalSkill("../x") {
		t.Fatal("IsOptionalSkill")
	}
	files, err := initrepo.OptionalSkillFiles("akm-commit", "docs/work", "docs/adr")
	if err != nil || len(files) == 0 {
		t.Fatal(files, err)
	}
	for _, f := range files {
		if !strings.HasPrefix(f.Path, ".agents/skills/akm-commit/") {
			t.Errorf("path %s", f.Path)
		}
	}
	if _, err := initrepo.OptionalSkillFiles("work-items", "docs/work", "docs/adr"); err == nil {
		t.Error("a skill that is not optional was served")
	}
}

func TestFilesWithGitCommits(t *testing.T) {
	p := profile(t, "code")
	plain, _ := initrepo.Files(p, "x")
	with, err := initrepo.FilesWith(p, "x", initrepo.Options{GitCommits: true})
	if err != nil || len(with) <= len(plain) {
		t.Fatal(len(with), len(plain), err)
	}
	for _, f := range with {
		if f.Path == "docs/.layout" && !strings.HasSuffix(string(f.Data), "git_commits: true\n") {
			t.Errorf("manifest %q", f.Data)
		}
	}
}
