package api

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/record"
)

const skillFile = ".agents/skills/akm-commit/SKILL.md"

// cfgRepo is a copy of the basic fixture repository; tree makes it a work tree.
func cfgRepo(t *testing.T, tree bool) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "basic")
	copyTree(t, "../../testdata/api/repos/basic", root)
	if tree {
		if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func runCfg(t *testing.T, root string, req Request) (Obj, *Error) {
	t.Helper()
	repo, err := record.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	return Run(repo, Info{"akmctl", "test"}, "config", req)
}

func cfgOK(t *testing.T, root string, req Request) Obj {
	t.Helper()
	resp, e := runCfg(t, root, req)
	if e != nil {
		t.Fatalf("%v: %v", req, e)
	}
	return resp
}

func cfgRefused(t *testing.T, root string, req Request, reason string) *Error {
	t.Helper()
	_, e := runCfg(t, root, req)
	if e == nil || e.Code != "refused" {
		t.Fatalf("%v: want refused, got %v", req, e)
	}
	if r, _ := e.Details.Get("reason"); r != reason {
		t.Fatalf("%v: reason %v, want %s", req, r, reason)
	}
	return e
}

func layoutText(t *testing.T, root string) string {
	t.Helper()
	b, _ := os.ReadFile(filepath.Join(root, "docs", ".layout"))
	return string(b)
}

func exists(root, rel string) bool {
	_, err := os.Lstat(filepath.Join(root, filepath.FromSlash(rel)))
	return err == nil
}

func planOf(resp Obj, k string) []string {
	p, _ := resp.Get("plan")
	v, _ := p.(Obj).Get(k)
	return v.([]string)
}

func TestConfigList(t *testing.T) {
	root := cfgRepo(t, false)
	resp := cfgOK(t, root, Request{})
	ss, _ := resp.Get("settings")
	l := ss.([]Obj)
	if len(l) != 1 {
		t.Fatalf("%v", l)
	}
	if k, _ := l[0].Get("key"); k != "git_commits" {
		t.Error(k)
	}
	if v, _ := l[0].Get("value"); v != false {
		t.Error(v)
	}
	if _, e := runCfg(t, root, Request{"dry_run": true}); e == nil || e.Code != "invalid_request" {
		t.Errorf("dry_run without a value: %v", e)
	}
}

func TestConfigOnOffIdempotent(t *testing.T) {
	root := cfgRepo(t, true)
	on := Request{"key": "git_commits", "value": true}
	resp := cfgOK(t, root, on)
	if ch, _ := resp.Get("changed"); ch != true {
		t.Fatal("not changed")
	}
	if !exists(root, skillFile) || !strings.Contains(layoutText(t, root), "git_commits: true\n") {
		t.Fatal("skill or key missing")
	}
	// on again: a no-op, nothing in the plan
	before := layoutText(t, root)
	resp = cfgOK(t, root, on)
	if ch, _ := resp.Get("changed"); ch != false {
		t.Error("second on changed something")
	}
	if len(planOf(resp, "write"))+len(planOf(resp, "modify"))+len(planOf(resp, "remove")) != 0 || layoutText(t, root) != before {
		t.Error("second on is not a no-op")
	}
	// dry-run off: the plan, nothing removed, no confirm needed for a dry run
	resp = cfgOK(t, root, Request{"key": "git_commits", "value": false, "dry_run": true})
	if got := planOf(resp, "remove"); len(got) == 0 || got[0] != skillFile {
		t.Errorf("remove plan %v", got)
	}
	if nc, _ := resp.Get("needs_confirmation"); nc != true {
		t.Error("needs_confirmation")
	}
	if !exists(root, skillFile) || layoutText(t, root) != before {
		t.Fatal("dry run wrote")
	}
	// off without confirm: refused with the plan
	e := cfgRefused(t, root, Request{"key": "git_commits", "value": false}, "needs_confirmation")
	if _, ok := e.Details.Get("plan"); !ok {
		t.Error("no plan in the refusal")
	}
	if !exists(root, skillFile) {
		t.Fatal("removed without confirm")
	}
	// off with confirm
	cfgOK(t, root, Request{"key": "git_commits", "value": false, "confirm": true})
	if exists(root, ".agents/skills/akm-commit") || !strings.Contains(layoutText(t, root), "git_commits: false\n") {
		t.Fatal("not removed or key not false")
	}
	if !exists(root, ".agents/skills") {
		t.Error("removed more than the skill directory")
	}
	// off again: a no-op
	resp = cfgOK(t, root, Request{"key": "git_commits", "value": false})
	if ch, _ := resp.Get("changed"); ch != false {
		t.Error("second off changed something")
	}
}

func TestConfigEditedNeedsSecondConfirm(t *testing.T) {
	root := cfgRepo(t, true)
	cfgOK(t, root, Request{"key": "git_commits", "value": true})
	skill := filepath.Join(root, filepath.FromSlash(skillFile))
	b, _ := os.ReadFile(skill)
	os.WriteFile(skill, append(b, "\nMy own line.\n"...), 0o644)
	extra := filepath.Join(root, ".agents", "skills", "akm-commit", "notes.md")
	os.WriteFile(extra, []byte("mine\n"), 0o644)
	off := Request{"key": "git_commits", "value": false}
	e := cfgRefused(t, root, off, "needs_confirmation")
	if ed, _ := e.Details.Get("edited"); ed != true {
		t.Error("edited not reported")
	}
	edits, _ := e.Details.Get("edits")
	if n := len(edits.([]Obj)); n != 2 {
		t.Fatalf("%d diffs, want 2", n)
	}
	d, _ := edits.([]Obj)[0].Get("hunks")
	if len(d.([]Obj)) == 0 {
		t.Error("empty diff")
	}
	// confirm alone is not enough
	off["confirm"] = true
	cfgRefused(t, root, off, "needs_confirmation")
	if !exists(root, skillFile) {
		t.Fatal("removed on one confirm")
	}
	off["confirm_edited"] = true
	cfgOK(t, root, off)
	if exists(root, ".agents/skills/akm-commit") {
		t.Fatal("not removed")
	}
}

func TestConfigConflicts(t *testing.T) {
	// a file where the directory goes
	root := cfgRepo(t, true)
	os.MkdirAll(filepath.Join(root, ".agents", "skills"), 0o755)
	os.WriteFile(filepath.Join(root, ".agents", "skills", "akm-commit"), []byte("x"), 0o644)
	before := layoutText(t, root)
	e := cfgRefused(t, root, Request{"key": "git_commits", "value": true}, "conflict")
	if p, _ := e.Details.Get("paths"); len(p.([]string)) != 1 || layoutText(t, root) != before {
		t.Errorf("paths %v", p)
	}
	// an edited file of the skill, with the key off
	root = cfgRepo(t, true)
	os.MkdirAll(filepath.Join(root, ".agents", "skills", "akm-commit"), 0o755)
	os.WriteFile(filepath.Join(root, filepath.FromSlash(skillFile)), []byte("my own skill\n"), 0o644)
	e = cfgRefused(t, root, Request{"key": "git_commits", "value": true}, "conflict")
	if p, _ := e.Details.Get("paths"); p.([]string)[0] != skillFile {
		t.Errorf("paths %v", p)
	}
	if got, _ := os.ReadFile(filepath.Join(root, filepath.FromSlash(skillFile))); string(got) != "my own skill\n" {
		t.Error("overwritten")
	}
	if strings.Contains(layoutText(t, root), "git_commits") {
		t.Error("manifest changed on a conflict")
	}
	// a link in place of the directory is never followed, on or off
	if runtime.GOOS != "windows" {
		root = cfgRepo(t, true)
		outside := t.TempDir()
		os.WriteFile(filepath.Join(outside, "keep"), []byte("k"), 0o644)
		os.MkdirAll(filepath.Join(root, ".agents", "skills"), 0o755)
		if err := os.Symlink(outside, filepath.Join(root, ".agents", "skills", "akm-commit")); err != nil {
			t.Skip(err)
		}
		cfgRefused(t, root, Request{"key": "git_commits", "value": true}, "conflict")
		cfgRefused(t, root, Request{"key": "git_commits", "value": false, "confirm": true, "confirm_edited": true}, "conflict")
		if !exists(outside, "keep") {
			t.Fatal("followed the link")
		}
	}
}

func TestConfigNeedsGitWorkTree(t *testing.T) {
	root := cfgRepo(t, false)
	cfgRefused(t, root, Request{"key": "git_commits", "value": true}, "not_a_git_repo")
	if exists(root, ".agents") || strings.Contains(layoutText(t, root), "git_commits") {
		t.Fatal("wrote outside a work tree")
	}
	// off always works
	cfgOK(t, root, Request{"key": "git_commits", "value": false})
}

func TestConfigNeedsManifest(t *testing.T) {
	root := cfgRepo(t, true)
	os.Remove(filepath.Join(root, "docs", ".layout"))
	cfgRefused(t, root, Request{"key": "git_commits", "value": true}, "no_manifest")
}

func TestConfigRepairsBadValue(t *testing.T) {
	root := cfgRepo(t, true)
	os.WriteFile(filepath.Join(root, "docs", ".layout"), []byte(layoutText(t, root)+"git_commits: maybe\n"), 0o644)
	cfgOK(t, root, Request{"key": "git_commits", "value": true})
	if !strings.HasSuffix(layoutText(t, root), "git_commits: true\n") {
		t.Errorf("%q", layoutText(t, root))
	}
}

func TestInitGitCommitsAndGitInit(t *testing.T) {
	run := func(t *testing.T, base string, req Request) Obj {
		t.Helper()
		resp, e := Run(&record.Repo{Root: base}, Info{"akmctl", "test"}, "init", req)
		if e != nil {
			t.Fatal(e)
		}
		return resp
	}
	// outside a work tree: forced off, reported
	base := t.TempDir()
	resp := run(t, base, Request{"profile": "code", "git_commits": true})
	if f, _ := resp.Get("git_commits_forced_off"); f != true {
		t.Error("not forced off")
	}
	if c, _ := resp.Get("git_commits"); c != false || exists(base, ".agents/skills/akm-commit") || strings.Contains(layoutText(t, base), "git_commits") {
		t.Error("installed outside a work tree")
	}
	// inside one: installed
	base = t.TempDir()
	os.MkdirAll(filepath.Join(base, ".git"), 0o755)
	resp = run(t, base, Request{"profile": "code", "git_commits": true, "git_init": true})
	if g, _ := resp.Get("git"); g != "existing" {
		t.Errorf("git %v", g)
	}
	if !exists(base, skillFile) || !strings.Contains(layoutText(t, base), "git_commits: true\n") {
		t.Error("skill or key missing")
	}
	// git init in a new directory, then the skill goes in
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("no git on PATH")
	}
	base = t.TempDir()
	resp = run(t, base, Request{"profile": "code", "dir": "proj", "git_commits": true, "git_init": true})
	if g, _ := resp.Get("git"); g != "initialised" {
		t.Errorf("git %v", resp)
	}
	if !exists(base, "proj/.git") || !exists(base, "proj/"+skillFile) {
		t.Error("git init or skill missing")
	}
	// the written repository checks clean
	repo, _ := record.Load(filepath.Join(base, "proj"))
	if !repo.GitCommits() {
		t.Error("key not read")
	}
}

func TestInitGitInitFailureIsNotFatal(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("no git on PATH")
	}
	t.Setenv("PATH", t.TempDir()) // git is gone
	resp, e := Run(&record.Repo{Root: t.TempDir()}, Info{"akmctl", "test"}, "init", Request{"profile": "code", "git_init": true, "git_commits": true})
	if e != nil {
		t.Fatal(e)
	}
	if g, _ := resp.Get("git"); g != "failed" {
		t.Errorf("git %v", g)
	}
	if _, ok := resp.Get("git_error"); !ok {
		t.Error("no git_error")
	}
	if f, _ := resp.Get("git_commits_forced_off"); f != true {
		t.Error("not forced off after a failed git init")
	}
}
