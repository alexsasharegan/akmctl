package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/api"
)

// script makes the CLI think it runs on a terminal that types answers.
func script(t *testing.T, tty bool, answers string) {
	t.Helper()
	oldIn, oldTTY := stdinReader, interactive
	stdinReader = strings.NewReader(answers)
	interactive = func(io.Writer) bool { return tty }
	t.Cleanup(func() { stdinReader, interactive = oldIn, oldTTY })
}

func tree(t *testing.T) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "basic")
	filepath.WalkDir(repo, func(p string, d os.DirEntry, err error) error {
		rel, _ := filepath.Rel(repo, p)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(root, rel), 0o755)
		}
		b, _ := os.ReadFile(p)
		return os.WriteFile(filepath.Join(root, rel), b, 0o644)
	})
	os.MkdirAll(filepath.Join(root, ".git"), 0o755)
	return root
}

func cfg(root string, args ...string) (int, string, string) {
	var out, errb bytes.Buffer
	code := Run(append(args, "--repo", root), &out, &errb, root)
	return code, out.String(), errb.String()
}

func has(root, rel string) bool {
	_, err := os.Lstat(filepath.Join(root, filepath.FromSlash(rel)))
	return err == nil
}

func TestPrompter(t *testing.T) {
	var out bytes.Buffer
	p := newPrompter(strings.NewReader("maybe\nY\n\nno\n"), &out)
	if !p.yesNo("Go?", false) { // "maybe" is re-asked, then Y
		t.Error("Y is yes")
	}
	if p.yesNo("Again?", false) { // empty takes the default
		t.Error("empty should take the default no")
	}
	if p.yesNo("Third?", true) { // "no"
		t.Error("no")
	}
	if !p.yesNo("End?", true) { // end of input takes the default
		t.Error("EOF default")
	}
	if !strings.Contains(out.String(), "Go? [y/N] ") || !strings.Contains(out.String(), "Third? [Y/n] ") || !strings.Contains(out.String(), "Please answer y or n.") {
		t.Errorf("%q", out.String())
	}
	out.Reset()
	p = newPrompter(strings.NewReader("a/b\n.x\nknowledge\n"), &out)
	if got := p.docsDir(); got != "knowledge" {
		t.Error(got)
	}
	if !strings.Contains(out.String(), "Where should the records go? (docs): ") || strings.Count(out.String(), "Where should") != 3 {
		t.Errorf("%q", out.String())
	}
	if got := newPrompter(strings.NewReader("\n"), &out).docsDir(); got != "docs" {
		t.Error(got)
	}
	if got := newPrompter(strings.NewReader(""), &out).docsDir(); got != "docs" {
		t.Error(got)
	}
}

func TestAskInit(t *testing.T) {
	var out bytes.Buffer
	req := api.Request{}
	askInit(newPrompter(strings.NewReader("notes\ny\nn\n"), &out), req, false)
	// not in a repository and git init declined: the commit skill is not asked
	if req["docs_dir"] != "notes" || req["wikilinks"] != true || req["git_init"] != false {
		t.Errorf("%v", req)
	}
	if _, ok := req["git_commits"]; ok || strings.Contains(out.String(), "akm-commit") {
		t.Errorf("asked about commits: %q", out.String())
	}
	// the order: folder, wikilinks, git init, then the commit skill
	req = api.Request{}
	out.Reset()
	askInit(newPrompter(strings.NewReader("\n\n\ny\n"), &out), req, false)
	order := []string{"Where should the records go? (docs): ", "Use [[wikilinks]] between records? [y/N] ", "Run git init? [Y/n] ", "Install the akm-commit skill (agents commit after each step)? [y/N] "}
	last := -1
	for _, q := range order {
		i := strings.Index(out.String(), q)
		if i <= last {
			t.Errorf("prompt %q out of order or missing in %q", q, out.String())
		}
		last = i
	}
	if req["docs_dir"] != "docs" || req["wikilinks"] != false || req["git_init"] != true || req["git_commits"] != true {
		t.Errorf("%v", req)
	}
	// flags answer first; inside a repository git init is not asked
	req = api.Request{"docs_dir": "d", "wikilinks": false, "git_commits": false}
	out.Reset()
	askInit(newPrompter(strings.NewReader(""), &out), req, true)
	if out.Len() != 0 || len(req) != 3 {
		t.Errorf("asked %q, req %v", out.String(), req)
	}
	// inside a repository the commit skill is asked, git init is not
	req = api.Request{"docs_dir": "d", "wikilinks": false}
	out.Reset()
	askInit(newPrompter(strings.NewReader("y\n"), &out), req, true)
	if strings.Contains(out.String(), "git init") || req["git_commits"] != true {
		t.Errorf("%q %v", out.String(), req)
	}
}

func TestInitAsksOnATerminalOnly(t *testing.T) {
	// a terminal: the answers decide
	script(t, true, "knowledge\ny\nn\n")
	base := t.TempDir()
	var out, errb bytes.Buffer
	if code := Run([]string{"init", "proj"}, &out, &errb, base); code != 0 {
		t.Fatal(code, errb.String())
	}
	for _, q := range []string{"Where should the records go? (docs): ", "Use [[wikilinks]] between records? [y/N] ", "Run git init? [Y/n] "} {
		if !strings.Contains(errb.String(), q) {
			t.Errorf("missing prompt %q in %q", q, errb.String())
		}
	}
	if !has(base, "proj/knowledge/work") || has(base, "proj/docs") {
		t.Error("the folder answer was not used")
	}
	if b, _ := os.ReadFile(filepath.Join(base, "proj", "knowledge", ".layout")); !strings.Contains(string(b), "reference_style: wikilink") {
		t.Error("the wikilinks answer was not used")
	}
	if has(base, "proj/.git") || has(base, "proj/.agents/skills/akm-commit") || has(base, "proj/.gitignore") {
		t.Error("git, the skill or .gitignore appeared")
	}
	// not a terminal: nothing is asked, even with input waiting; the defaults apply
	script(t, false, "notes\ny\ny\n")
	out.Reset()
	errb.Reset()
	if code := Run([]string{"init", "other"}, &out, &errb, base); code != 0 || errb.Len() != 0 {
		t.Errorf("bare init without a terminal: %d %q", code, errb.String())
	}
	if !has(base, "other/docs/.layout") || has(base, "other/.git") || has(base, "other/.agents/skills/akm-commit") {
		t.Error("defaults were not used")
	}
	// a terminal with --yes or --no-input: not asked either
	for _, flag := range []string{"--yes", "--no-input"} {
		script(t, true, "notes\n")
		errb.Reset()
		dir := "d" + flag[2:3]
		if code := Run([]string{"init", flag, dir}, &out, &errb, base); code != 0 {
			t.Errorf("%s: %d %q", flag, code, errb.String())
		}
		if strings.Contains(errb.String(), "Where should") || !has(base, dir+"/docs/.layout") {
			t.Errorf("%s asked or ignored the defaults", flag)
		}
	}
	// flags skip their questions: only the folder is asked
	script(t, true, "wiki\n")
	errb.Reset()
	if code := Run([]string{"init", "--no-wikilinks", "--no-git-commits", "--no-git-init", "flagged"}, &out, &errb, base); code != 0 {
		t.Fatal(code, errb.String())
	}
	if strings.Contains(errb.String(), "akm-commit") || strings.Contains(errb.String(), "git init?") || strings.Contains(errb.String(), "wikilinks") {
		t.Errorf("asked what the flags answered: %q", errb.String())
	}
	if !has(base, "flagged/wiki/.layout") {
		t.Error("the folder answer was not used")
	}
}

func TestInitGitCommitsFlag(t *testing.T) {
	script(t, false, "")
	base := t.TempDir()
	os.MkdirAll(filepath.Join(base, ".git"), 0o755)
	var out, errb bytes.Buffer
	if code := Run([]string{"init", "--git-commits", "--yes", "p"}, &out, &errb, base); code != 0 {
		t.Fatal(code, errb.String())
	}
	if !has(base, "p/.agents/skills/akm-commit/SKILL.md") {
		t.Error("skill not installed")
	}
}

func TestConfigCommand(t *testing.T) {
	script(t, false, "")
	root := tree(t)
	code, out, _ := cfg(root, "config")
	if code != 0 || !strings.Contains(out, "git_commits = false") {
		t.Fatalf("list: %d %q", code, out)
	}
	// dry run: the plan, nothing written
	code, out, _ = cfg(root, "config", "git-commits", "on", "--dry-run")
	if code != 0 || !strings.Contains(out, "would be turned on") || has(root, ".agents") {
		t.Fatalf("dry run: %d %q", code, out)
	}
	if code, out, _ = cfg(root, "config", "git-commits", "on"); code != 0 || !has(root, ".agents/skills/akm-commit/SKILL.md") {
		t.Fatalf("on: %d %q", code, out)
	}
	if code, out, _ = cfg(root, "config", "git-commits", "on"); code != 0 || !strings.Contains(out, "nothing to do") {
		t.Fatalf("on again: %d %q", code, out)
	}
	if code, _, _ = cfg(root, "config", "git-commits", "maybe"); code != 2 {
		t.Error("bad word")
	}
	// off without a terminal and without --yes: refused, nothing removed
	code, _, errb := cfg(root, "config", "git-commits", "off")
	if code != 2 || !strings.Contains(errb, "--yes") || !has(root, ".agents/skills/akm-commit") {
		t.Fatalf("off, no terminal: %d %q", code, errb)
	}
	// off on a terminal, answer no
	script(t, true, "n\n")
	code, _, errb = cfg(root, "config", "git-commits", "off")
	if code != 1 || !has(root, ".agents/skills/akm-commit") || !strings.Contains(errb, "remove .agents/skills/akm-commit/SKILL.md") {
		t.Fatalf("off, no: %d %q", code, errb)
	}
	// off with edits: two questions
	skill := filepath.Join(root, ".agents", "skills", "akm-commit", "SKILL.md")
	b, _ := os.ReadFile(skill)
	os.WriteFile(skill, append(b, "mine\n"...), 0o644)
	script(t, true, "y\nn\n")
	code, _, errb = cfg(root, "config", "git-commits", "off")
	if code != 1 || !has(root, ".agents/skills/akm-commit") || !strings.Contains(errb, "+mine") || !strings.Contains(errb, "edited") {
		t.Fatalf("off, edited, second no: %d %q", code, errb)
	}
	script(t, true, "y\ny\n")
	if code, _, errb = cfg(root, "config", "git-commits", "off"); code != 0 || has(root, ".agents/skills/akm-commit") {
		t.Fatalf("off, edited, yes: %d %q", code, errb)
	}
	// --yes skips the questions
	script(t, false, "")
	cfg(root, "config", "git-commits", "on")
	if code, _, errb = cfg(root, "config", "git-commits", "off", "--yes"); code != 0 || has(root, ".agents/skills/akm-commit") {
		t.Fatalf("off --yes: %d %q", code, errb)
	}
	// no work tree: on is refused
	os.RemoveAll(filepath.Join(root, ".git"))
	code, _, errb = cfg(root, "config", "git-commits", "on")
	if code != 2 || !strings.Contains(errb, "git work tree") {
		t.Fatalf("no work tree: %d %q", code, errb)
	}
	// a conflict is reported with the path
	os.MkdirAll(filepath.Join(root, ".git"), 0o755)
	os.MkdirAll(filepath.Join(root, ".agents", "skills", "akm-commit"), 0o755)
	os.WriteFile(skill, []byte("x"), 0o644)
	code, out, _ = cfg(root, "config", "git-commits", "on", "--format", "json")
	if code != 2 || !strings.Contains(out, `"reason":"conflict"`) || !strings.Contains(out, "SKILL.md") {
		t.Fatalf("conflict: %d %q", code, out)
	}
}
