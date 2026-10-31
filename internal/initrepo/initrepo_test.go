package initrepo_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/api"
	"github.com/alexsasharegan/akmctl/internal/check"
	"github.com/alexsasharegan/akmctl/internal/initrepo"
	"github.com/alexsasharegan/akmctl/internal/record"
)

func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	m := map[string]string{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			rel, _ := filepath.Rel(root, p)
			b, _ := os.ReadFile(p)
			m[filepath.ToSlash(rel)] = string(b)
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// variants are the option sets the tests walk through: the defaults and a
// custom folder with wikilinks.
var variants = []struct {
	name string
	opt  initrepo.Options
}{
	{"default", initrepo.Options{Gitignore: true}},
	{"knowledge", initrepo.Options{DocsDir: "knowledge", WikiLinks: true, Gitignore: true, Obsidian: true}},
}

func run(t *testing.T, root, op string, req api.Request) api.Obj {
	t.Helper()
	repo, err := record.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	resp, aerr := api.Run(repo, api.Info{Name: "akmctl", Version: "test"}, op, req)
	if aerr != nil {
		t.Fatalf("%s %v: %v", op, req, aerr)
	}
	return resp
}

func TestOptionsCheckClean(t *testing.T) {
	for _, v := range variants {
		t.Run(v.name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "fresh", "proj") // neither level exists
			res, err := initrepo.InitWith(dir, v.opt)
			if err != nil {
				t.Fatal(err)
			}
			if len(res.Written) == 0 || len(res.Skipped) != 0 {
				t.Fatalf("written %v skipped %v", res.Written, res.Skipped)
			}
			resp := run(t, dir, "check", api.Request{})
			for _, k := range []string{"errors", "warnings"} {
				if n, _ := resp.Get(k); n != 0 {
					f, _ := resp.Get("findings")
					t.Errorf("%s = %v: %v", k, n, f)
				}
			}
			repo, _ := record.Load(dir)
			docs := v.opt.DocsDir
			if docs == "" {
				docs = "docs"
			}
			if repo.MetaDir() != docs || repo.WorkDir() != docs+"/work" || repo.ADRDir() != docs+"/adr" || repo.RoadmapPath() != docs+"/ROADMAP.md" {
				t.Errorf("folders: %s %s %s %s", repo.MetaDir(), repo.WorkDir(), repo.ADRDir(), repo.RoadmapPath())
			}
			if _, err := os.Stat(filepath.Join(dir, ".obsidian")); err == nil {
				t.Error("init must not write .obsidian")
			}
			if wiki := repo.Layout["reference_style"].Str == "wikilink"; wiki != v.opt.WikiLinks {
				t.Errorf("reference_style wikilink = %v", wiki)
			}
			files := snapshot(t, dir)
			if files["CLAUDE.md"] != "@AGENTS.md\n" {
				t.Errorf("CLAUDE.md = %q", files["CLAUDE.md"])
			}
			if !strings.Contains(files[".gitignore"], docs+"/.akm.lock\n") {
				t.Errorf(".gitignore = %q", files[".gitignore"])
			}
			if got := strings.Contains(files[".gitignore"], ".obsidian/workspace.json"); got != v.opt.Obsidian {
				t.Errorf("obsidian lines = %v in %q", got, files[".gitignore"])
			}
			ag := files["AGENTS.md"]
			if !strings.HasPrefix(ag, "# proj\n") {
				t.Errorf("AGENTS.md title: %q", ag[:20])
			}
			if strings.Contains(ag, "{{") || strings.Contains(files["README.md"], "{{") {
				t.Error("an unreplaced marker")
			}
			for _, bad := range []string{"- Build:", "- Test:", "- Lint", "- Run:"} {
				if strings.Contains(ag, bad) {
					t.Errorf("AGENTS.md has a %q line", bad)
				}
			}
			if !strings.Contains(ag, "--check AGENTS.md \""+docs+"/**/*.md\"") || !strings.Contains(ag, "`"+docs+"/.akm-deny`") {
				t.Errorf("AGENTS.md does not use the docs dir %s", docs)
			}
			if docs != "docs" {
				if _, err := os.Stat(filepath.Join(dir, "docs")); err == nil {
					t.Error("a docs folder was created")
				}
				if strings.Contains(files[".agents/skills/akm-track-work/SKILL.md"], "docs/work") {
					t.Error("the skill still points at docs/work")
				}
			}
		})
	}
}

func TestGitignoreOnlyWhenAsked(t *testing.T) {
	dir := t.TempDir()
	if _, err := initrepo.InitWith(dir, initrepo.Options{}); err != nil {
		t.Fatal(err)
	}
	if _, ok := snapshot(t, dir)[".gitignore"]; ok {
		t.Error(".gitignore written without a work tree")
	}
}

func TestValidateDocsDir(t *testing.T) {
	for _, ok := range []string{"docs", "knowledge", "My Notes", "notes-2", "a.b"} {
		if err := initrepo.ValidateDocsDir(ok); err != nil {
			t.Errorf("%q: %v", ok, err)
		}
	}
	for _, bad := range []string{"", ".", "..", "a/b", "a\\b", ".hidden", "x:y", "a*", "a?", "a|b", "a\x00b", "node_modules", " a", "a ", "a.", "#a", "!a", "[a]"} {
		if err := initrepo.ValidateDocsDir(bad); err == nil {
			t.Errorf("%q accepted", bad)
		}
	}
}

func TestNoOverwriteAndIdempotent(t *testing.T) {
	for _, v := range variants {
		t.Run(v.name, func(t *testing.T) {
			dir := t.TempDir()
			mine := "# mine\n\nkeep me\n"
			if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte(mine), 0o644); err != nil {
				t.Fatal(err)
			}
			res, err := initrepo.InitWith(dir, v.opt)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(res.Skipped, []string{"AGENTS.md"}) {
				t.Errorf("skipped %v", res.Skipped)
			}
			after := snapshot(t, dir)
			if after["AGENTS.md"] != mine {
				t.Error("AGENTS.md was overwritten")
			}
			res2, err := initrepo.InitWith(dir, v.opt)
			if err != nil {
				t.Fatal(err)
			}
			if len(res2.Written) != 0 || len(res2.Skipped) != len(res.Written)+1 {
				t.Errorf("second run: written %v, skipped %d", res2.Written, len(res2.Skipped))
			}
			if !reflect.DeepEqual(snapshot(t, dir), after) {
				t.Error("a second run changed files")
			}
		})
	}
}

func TestTargetIsAFile(t *testing.T) {
	f := filepath.Join(t.TempDir(), "file")
	os.WriteFile(f, []byte("x"), 0o644)
	if _, err := initrepo.Init(f); err == nil {
		t.Error("init into a file succeeded")
	}
}

// Versions come from the tool's constants, not from the templates.
func TestVersionsFromConstants(t *testing.T) {
	err := fs.WalkDir(os.DirFS("templates"), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, _ := os.ReadFile(filepath.Join("templates", p))
		if strings.Contains(string(b), check.LayoutVersion) {
			t.Errorf("template %s names the layout version %s", p, check.LayoutVersion)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	files, _ := initrepo.Files("x")
	for _, f := range files {
		if f.Path == "docs/.layout" && !strings.HasPrefix(string(f.Data), "layout: "+check.LayoutVersion+"\n") {
			t.Errorf("manifest: %q", f.Data)
		}
	}
}

// The embedded skills are a copy of this repository's own.
func TestSkillsMatchRepository(t *testing.T) {
	src := snapshot(t, "../../.agents/skills")
	// This repository has the optional skill on (git_commits); it must match its shipped copy.
	opt := map[string]string{}
	for k, v := range src {
		if rest, ok := strings.CutPrefix(k, "akm-commit/"); ok {
			opt[rest] = v
			delete(src, k)
		}
	}
	if len(opt) > 0 && !reflect.DeepEqual(opt, snapshot(t, "templates/optional-skills/akm-commit")) {
		t.Error("internal/initrepo/templates/optional-skills/akm-commit differs from .agents/skills/akm-commit: copy it over")
	}
	got := snapshot(t, "templates/skills")
	if !reflect.DeepEqual(src, got) {
		t.Error("internal/initrepo/templates/skills differs from .agents/skills: copy the skills over")
	}
}

// The always-installed skills are the three akm- ones, in the drift test and in what init writes.
func TestAlwaysInstalledSkills(t *testing.T) {
	want := []string{"akm", "akm-record-decision", "akm-track-work"}
	for _, dir := range []string{"../../.agents/skills", "templates/skills"} {
		ents, err := os.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for _, e := range ents {
			// this repository may have the optional skill installed
			if e.IsDir() && e.Name() != "akm-commit" {
				got = append(got, e.Name())
			}
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s holds %v, want %v", dir, got, want)
		}
	}
	files, _ := initrepo.Files("x")
	have := map[string]bool{}
	for _, f := range files {
		have[f.Path] = true
	}
	for _, n := range want {
		if p := ".agents/skills/" + n + "/SKILL.md"; !have[p] {
			t.Errorf("init does not write %s", p)
		}
	}
	if have[".agents/skills/akm-commit/SKILL.md"] {
		t.Error("the optional skill akm-commit is installed by default")
	}
}

// Every API operation has a reference file in the akm skill.
func TestEveryOperationHasReference(t *testing.T) {
	ops := api.Operations()
	if len(ops) == 0 {
		t.Fatal("no operations")
	}
	for _, op := range ops {
		for _, dir := range []string{"../../.agents/skills/akm/references", "templates/skills/akm/references"} {
			b, err := os.ReadFile(filepath.Join(dir, op+".md"))
			if err != nil {
				t.Errorf("operation %s has no reference: %v", op, err)
				continue
			}
			if n := strings.Count(string(b), "\n"); n > 40 {
				t.Errorf("%s/%s.md has %d lines, want at most 40", dir, op, n)
			}
		}
	}
}

// The optional skill is shipped, checks the manifest first and never stages everything.
func TestOptionalCommitSkill(t *testing.T) {
	b, err := os.ReadFile("templates/optional-skills/akm-commit/SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, want := range []string{"name: akm-commit", "git_commits", "STOP", "git add -A"} {
		if !strings.Contains(s, want) {
			t.Errorf("akm-commit lacks %q", want)
		}
	}
}

// A cold agent can use the written files to list, claim, work and close an item.
func TestWorkLifecycle(t *testing.T) {
	for _, v := range variants {
		t.Run(v.name, func(t *testing.T) {
			dir := t.TempDir()
			if _, err := initrepo.InitWith(dir, v.opt); err != nil {
				t.Fatal(err)
			}
			ls := run(t, dir, "ls", api.Request{"status": []string{"ready", "active", "blocked"}})
			if n, _ := ls.Get("total"); n != 0 {
				t.Fatalf("a new repository lists %v items", n)
			}
			resp := run(t, dir, "new", api.Request{"kind": "work", "slug": "first-task", "title": "First task"})
			rec, _ := resp.Get("record")
			id, _ := rec.(api.Obj).Get("id")
			run(t, dir, "section", api.Request{"id": id, "name": "Goal", "text": "Try it.\n\n- [ ] it works"})
			run(t, dir, "set", api.Request{"id": id, "changes": map[string]any{"status": "ready"}})
			run(t, dir, "claim", api.Request{"id": id, "as": "agent"})
			run(t, dir, "section", api.Request{"id": id, "name": "Next", "text": "Close it."})
			run(t, dir, "section", api.Request{"id": id, "name": "Outcome", "text": "It works."})
			run(t, dir, "close", api.Request{"id": id, "status": "done"})
			run(t, dir, "new", api.Request{"kind": "adr", "slug": "use-it"})
			run(t, dir, "index", api.Request{})
			c := run(t, dir, "check", api.Request{})
			if e, _ := c.Get("errors"); e != 0 {
				f, _ := c.Get("findings")
				t.Errorf("check after the lifecycle: %v", f)
			}
		})
	}
}

func TestAPIRequests(t *testing.T) {
	base := t.TempDir()
	repo := &record.Repo{Root: base}
	info := api.Info{Name: "akmctl", Version: "test"}
	for _, c := range []struct {
		name string
		req  api.Request
		code string
	}{
		{"profile is gone", api.Request{"profile": "code"}, "invalid_request"},
		{"docs_dir empty", api.Request{"docs_dir": ""}, "invalid_request"},
		{"docs_dir dot", api.Request{"docs_dir": "."}, "invalid_request"},
		{"docs_dir dotdot", api.Request{"docs_dir": ".."}, "invalid_request"},
		{"docs_dir nested", api.Request{"docs_dir": "a/b"}, "invalid_request"},
		{"docs_dir hidden", api.Request{"docs_dir": ".notes"}, "invalid_request"},
		{"docs_dir backslash", api.Request{"docs_dir": "a\\b"}, "invalid_request"},
		{"docs_dir colon", api.Request{"docs_dir": "a:b"}, "invalid_request"},
		{"docs_dir not a string", api.Request{"docs_dir": 3}, "invalid_request"},
		{"dir climbs out", api.Request{"dir": "../x"}, "invalid_request"},
		{"dir is absolute", api.Request{"dir": "/tmp/x"}, "invalid_request"},
		{"unknown parameter", api.Request{"force": true}, "invalid_request"},
	} {
		_, aerr := api.Run(repo, info, "init", c.req)
		if aerr == nil || aerr.Code != c.code {
			t.Errorf("%s: %v", c.name, aerr)
			continue
		}
		if p := strings.TrimPrefix(c.name, "docs_dir"); p != c.name {
			if d, _ := aerr.Details.Get("parameter"); d != "docs_dir" {
				t.Errorf("%s: details.parameter = %v", c.name, d)
			}
		}
	}
	// a bare init works
	resp, aerr := api.Run(repo, info, "init", api.Request{"dir": "a/b"})
	if aerr != nil {
		t.Fatal(aerr)
	}
	if w, _ := resp.Get("written"); len(w.([]string)) == 0 {
		t.Error("nothing written")
	}
	if _, err := os.Stat(filepath.Join(base, "a", "b", "docs", ".layout")); err != nil {
		t.Error(err)
	}
	os.WriteFile(filepath.Join(base, "f"), []byte("x"), 0o644)
	if _, aerr := api.Run(repo, info, "init", api.Request{"dir": "f"}); aerr == nil || aerr.Code != "refused" {
		t.Errorf("dir is a file: %v", aerr)
	}
}

// A second init with another docs dir would make two manifests: it is refused.
func TestSecondManifestRefused(t *testing.T) {
	base := t.TempDir()
	repo := &record.Repo{Root: base}
	info := api.Info{Name: "akmctl", Version: "test"}
	if _, aerr := api.Run(repo, info, "init", api.Request{"docs_dir": "knowledge"}); aerr != nil {
		t.Fatal(aerr)
	}
	_, aerr := api.Run(repo, info, "init", api.Request{})
	if aerr == nil || aerr.Code != "refused" {
		t.Fatalf("second manifest: %v", aerr)
	}
	if r, _ := aerr.Details.Get("reason"); r != "manifest_exists" {
		t.Errorf("reason = %v", r)
	}
	// the same dir again is idempotent
	if _, aerr := api.Run(repo, info, "init", api.Request{"docs_dir": "knowledge"}); aerr != nil {
		t.Errorf("repeat: %v", aerr)
	}
}

// Skill names follow the grammar of ADR 261030-wsqb: the router is akm, every
// other skill is akm-<verb>[-<object>] in lower-case words.
func TestSkillNamesFollowTheGrammar(t *testing.T) {
	re := regexp.MustCompile(`^akm(-[a-z]+)+$`)
	for _, root := range []string{"../../.agents/skills", "templates/skills", "templates/optional-skills"} {
		entries, err := os.ReadDir(root)
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range entries {
			if e.IsDir() && e.Name() != "akm" && !re.MatchString(e.Name()) {
				t.Errorf("%s/%s is not akm or akm-<verb>[-<object>]", root, e.Name())
			}
		}
	}
}
