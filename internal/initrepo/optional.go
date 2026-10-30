package initrepo

import (
	"io/fs"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// CommitSkill is the optional skill that git_commits switches.
const CommitSkill = "akm-commit"

// OptionalSkills lists the skills the tool may install or remove after init:
// each name starts with "akm-" and the shipped copy is in templates/optional-skills.
var OptionalSkills = []string{CommitSkill}

// SkillDir is the repository-relative directory of a skill.
func SkillDir(name string) string { return ".agents/skills/" + name }

// IsOptionalSkill reports whether name is a known optional skill.
func IsOptionalSkill(name string) bool {
	for _, s := range OptionalSkills {
		if s == name && strings.HasPrefix(s, "akm-") {
			return true
		}
	}
	return false
}

// OptionalSkillFiles returns the shipped files of an optional skill with
// repository-relative paths (.agents/skills/<name>/...), sorted. The folders of
// the repository replace docs/work and docs/adr in the text, as for the skills
// that init always writes.
func OptionalSkillFiles(name, workDir, adrDir string) ([]File, error) {
	if !IsOptionalSkill(name) {
		return nil, &fs.PathError{Op: "optional skill", Path: name, Err: fs.ErrNotExist}
	}
	paths := strings.NewReplacer("docs/work", workDir, "docs/adr", adrDir)
	root := "templates/optional-skills/" + name
	var files []File
	err := fs.WalkDir(templates, root, func(fp string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := templates.ReadFile(fp)
		if err != nil {
			return err
		}
		files = append(files, File{SkillDir(name) + "/" + strings.TrimPrefix(fp, root+"/"), []byte(paths.Replace(string(b)))})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

// GitInit runs `git init` in dir, which must exist. The environment loses the
// variables that would point git at another repository.
func GitInit(dir string) error {
	cmd := exec.Command("git", "init", "--quiet")
	cmd.Dir = dir
	for _, e := range os.Environ() {
		k, _, _ := strings.Cut(e, "=")
		if !strings.HasPrefix(k, "GIT_DIR") && !strings.HasPrefix(k, "GIT_WORK_TREE") && !strings.HasPrefix(k, "GIT_INDEX_FILE") {
			cmd.Env = append(cmd.Env, e)
		}
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return &gitError{msg}
	}
	return nil
}

type gitError struct{ msg string }

func (e *gitError) Error() string { return "git init: " + e.msg }
