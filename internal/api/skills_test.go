package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/initrepo"
	"github.com/alexsasharegan/akmctl/internal/record"
)

func skillState(t *testing.T, resp Obj, name string) string {
	t.Helper()
	l, _ := resp.Get("skills")
	for _, s := range l.([]Obj) {
		if n, _ := s.Get("name"); n == name {
			st, _ := s.Get("state")
			return st.(string)
		}
	}
	t.Fatalf("no skill %s in %v", name, resp)
	return ""
}

func skillPath(root, rel string) string { return filepath.Join(root, filepath.FromSlash(rel)) }

func TestSkillsInstallStatusUpdate(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "basic", "")
	resp := mustOK(t, root, "skills", Request{})
	if n, _ := resp.Get("notes"); n == nil || !strings.Contains(n.([]string)[0], "git-commits on") {
		t.Fatalf("no note about akm-commit: %v", resp)
	}
	if c, _ := resp.Get("clean"); c != false || skillState(t, resp, "akm") != "missing" {
		t.Fatalf("a bare repo has no skills: %v", resp)
	}
	// a write without confirm is a plan
	resp = mustOK(t, root, "skills", Request{"mode": "install"})
	if nc, _ := resp.Get("needs_confirmation"); nc != true {
		t.Fatalf("want a plan: %v", resp)
	}
	if _, err := os.Stat(skillPath(root, ".agents/skills/akm/SKILL.md")); err == nil {
		t.Fatal("a plan wrote a file")
	}
	resp = mustOK(t, root, "skills", Request{"mode": "install", "confirm": true})
	if ap, _ := resp.Get("applied"); ap != true {
		t.Fatalf("not applied: %v", resp)
	}
	for _, p := range []string{".agents/skills/akm/.akm-skill.json", ".agents/skills/akm-track-work/SKILL.md", ".agents/skills/akm-record-decision/.akm-skill.json"} {
		if _, err := os.Stat(skillPath(root, p)); err != nil {
			t.Errorf("%s not written", p)
		}
	}
	if _, err := os.Stat(skillPath(root, ".agents/skills/akm-commit")); err == nil {
		t.Error("the optional skill came with install")
	}
	resp = mustOK(t, root, "skills", Request{"mode": "status"})
	if c, _ := resp.Get("clean"); c != true {
		t.Fatalf("not clean after install: %v", resp)
	}
	// an edit is reported, and update leaves it unless confirm_edited
	skill := skillPath(root, ".agents/skills/akm-track-work/SKILL.md")
	orig, _ := os.ReadFile(skill)
	os.WriteFile(skill, append(orig, "\nmy note\n"...), 0o644)
	if st := skillState(t, mustOK(t, root, "skills", Request{}), "akm-track-work"); st != "edited" {
		t.Fatalf("state %s", st)
	}
	resp = mustOK(t, root, "skills", Request{"mode": "update", "confirm": true})
	if b, _ := os.ReadFile(skill); !strings.Contains(string(b), "my note") {
		t.Fatal("update overwrote an edit")
	}
	if sk, _ := resp.Get("skipped"); len(sk.([]Obj)) != 1 {
		t.Fatalf("skipped %v", sk)
	}
	mustOK(t, root, "skills", Request{"mode": "update", "confirm": true, "confirm_edited": true})
	if b, _ := os.ReadFile(skill); string(b) != string(orig) {
		t.Fatal("confirm_edited did not restore the shipped text")
	}
}

func TestSkillsOutdatedUpdatesWithoutAsking(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "basic", "")
	mustOK(t, root, "skills", Request{"mode": "install", "confirm": true})
	// an older release: the file differs and the marker records the old hash
	name := "akm-record-decision"
	managed, _ := initrepo.Managed("docs/work", "docs/adr", false)
	var sk initrepo.Skill
	for _, m := range managed {
		if m.Name == name {
			sk = m
		}
	}
	old := []byte("an older text\n")
	rel := sk.Files[0].Path
	os.WriteFile(skillPath(root, initrepo.SkillDir(name)+"/"+rel), old, 0o644)
	older := sk
	older.Files = append([]initrepo.File(nil), sk.Files...)
	older.Files[0].Data = old
	os.WriteFile(skillPath(root, initrepo.SkillDir(name)+"/"+initrepo.MarkerName), older.Marker("0.0.1"), 0o644)
	if st := skillState(t, mustOK(t, root, "skills", Request{}), name); st != "outdated" {
		t.Fatalf("state %s", st)
	}
	mustOK(t, root, "skills", Request{"mode": "install", "confirm": true})
	if b, _ := os.ReadFile(skillPath(root, initrepo.SkillDir(name)+"/"+rel)); string(b) != string(old) {
		t.Fatal("install changed an outdated skill")
	}
	mustOK(t, root, "skills", Request{"mode": "update", "confirm": true})
	if st := skillState(t, mustOK(t, root, "skills", Request{}), name); st != "current" {
		t.Fatalf("after update %s", st)
	}
}

func TestSkillsRetiredNameIsRemoved(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "basic", "")
	mustOK(t, root, "skills", Request{"mode": "install", "confirm": true})
	// the old folder with no marker (0.3.x wrote none) counts as edited
	oldDir := skillPath(root, ".agents/skills/akm-work-items")
	os.MkdirAll(oldDir, 0o755)
	os.WriteFile(filepath.Join(oldDir, "SKILL.md"), []byte("old\n"), 0o644)
	if st := skillState(t, mustOK(t, root, "skills", Request{}), "akm-work-items"); st != "renamed" {
		t.Fatalf("state %s", st)
	}
	mustOK(t, root, "skills", Request{"mode": "update", "confirm": true})
	if _, err := os.Stat(oldDir); err != nil {
		t.Fatal("an unmarked old folder was removed without confirm_edited")
	}
	mustOK(t, root, "skills", Request{"mode": "update", "confirm": true, "confirm_edited": true})
	if _, err := os.Stat(oldDir); err == nil {
		t.Fatal("the retired folder is still there")
	}
	// a marked, unedited old folder goes with a plain confirm
	os.MkdirAll(oldDir, 0o755)
	data := []byte("old\n")
	os.WriteFile(filepath.Join(oldDir, "SKILL.md"), data, 0o644)
	m := initrepo.Skill{Name: "akm-work-items", Files: []initrepo.File{{Path: "SKILL.md", Data: data}}}.Marker("0.3.0")
	os.WriteFile(filepath.Join(oldDir, initrepo.MarkerName), m, 0o644)
	mustOK(t, root, "skills", Request{"mode": "update", "confirm": true})
	if _, err := os.Stat(oldDir); err == nil {
		t.Fatal("an unedited retired folder was kept")
	}
}

func TestSkillsRefusals(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "basic", "")
	wantCode(t, root, "skills", Request{"confirm": true}, "invalid_request", "")
	wantCode(t, root, "skills", Request{"mode": "remove"}, "invalid_request", "")
	// dry_run never writes, even with confirm
	mustOK(t, root, "skills", Request{"mode": "install", "dry_run": true, "confirm": true})
	if _, err := os.Stat(skillPath(root, ".agents")); err == nil {
		t.Fatal("a dry run wrote")
	}
	// a file where the folder should be is a conflict, not overwritten
	os.WriteFile(skillPath(root, ".agents"), []byte("x"), 0o644)
	resp := mustOK(t, root, "skills", Request{"mode": "install", "confirm": true})
	if skillState(t, resp, "akm") != "conflict" {
		t.Fatalf("%v", resp)
	}
	if b, _ := os.ReadFile(skillPath(root, ".agents")); string(b) != "x" {
		t.Fatal("the conflicting file changed")
	}
}

func TestSkillsNewerThanTheBinary(t *testing.T) {
	t.Parallel()
	root := copyRepo(t, "basic", "")
	mustOK(t, root, "skills", Request{"mode": "install", "confirm": true})
	run := func(version string, req Request) Obj {
		t.Helper()
		repo, err := record.Load(root)
		if err != nil {
			t.Fatal(err)
		}
		resp, aerr := Run(repo, Info{"akmctl", version}, "skills", req)
		if aerr != nil {
			t.Fatal(aerr)
		}
		return resp
	}
	// a skill written by a newer tool than the running one
	name := "akm-track-work"
	managed, _ := initrepo.Managed("docs/work", "docs/adr", false)
	for _, m := range managed {
		if m.Name == name {
			os.WriteFile(skillPath(root, initrepo.SkillDir(name)+"/"+initrepo.MarkerName), m.Marker("0.4.0"), 0o644)
		}
	}
	resp := run("0.3.4-0.20261108200140-e4537d4a1b5d", Request{})
	notes, _ := resp.Get("notes")
	if n, _ := notes.([]string); len(n) == 0 || !strings.Contains(n[0], "install akmctl 0.4.0 or newer") {
		t.Fatalf("no note naming the version: %v", resp)
	}
	if c, _ := resp.Get("clean"); c != false {
		t.Fatalf("clean with skills from a newer tool: %v", resp)
	}
	sl, _ := resp.Get("skills")
	for _, s := range sl.([]Obj) {
		got, _ := s.Get("requires_tool")
		if n, _ := s.Get("name"); (n == name) != (got == "0.4.0") {
			t.Errorf("%v: requires_tool %v", n, got)
		}
	}
	// update leaves such a skill alone
	os.WriteFile(skillPath(root, initrepo.SkillDir(name)+"/SKILL.md"), []byte("newer text\n"), 0o644)
	resp = run("0.3.4", Request{"mode": "update", "confirm": true})
	if b, _ := os.ReadFile(skillPath(root, initrepo.SkillDir(name)+"/SKILL.md")); string(b) != "newer text\n" {
		t.Fatal("update downgraded a skill")
	}
	if sk, _ := resp.Get("skipped"); len(sk.([]Obj)) != 1 {
		t.Fatalf("skipped %v", sk)
	}
	// the same tool, a newer one, a development stamp and an unknown version raise nothing
	for _, v := range []string{"0.4.0", "0.4.1", "1.0.0", "0.0.0-dev+0123456", "test", ""} {
		if n, _ := run(v, Request{}).Get("notes"); n != nil && strings.Contains(strings.Join(n.([]string), ""), "newer than this") {
			t.Errorf("version %q: noted a newer skill", v)
		}
	}
}

func TestToolNewer(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want bool
	}{
		{"0.4.0", "0.3.4", true},
		{"0.3.4", "0.4.0", false},
		{"0.3.4", "0.3.4", false},
		{"v0.3.4", "0.3.4-0.20261108200140-e4537d4a1b5d", true},
		{"0.3.4-0.20261108200141-aaaaaaaaaaaa", "0.3.4-0.20261108200140-e4537d4a1b5d", true},
		{"0.3.4-0.20261108200140-e4537d4a1b5d", "0.3.4-0.20261108200141-aaaaaaaaaaaa", false},
		{"0.3.5-0.20261108200140-e4537d4a1b5d", "0.3.4", true},
		{"1.0.0", "0.0.0-dev+abc", false},
		{"", "0.3.4", false},
		{"0.3.4", "", false},
		{"nonsense", "0.3.4", false},
	} {
		if got := toolNewer(c.a, c.b); got != c.want {
			t.Errorf("toolNewer(%q, %q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}
