package initrepo_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
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

func profile(t *testing.T, name string) initrepo.Profile {
	t.Helper()
	p, ok := initrepo.Get(name)
	if !ok {
		t.Fatal("no profile", name)
	}
	return p
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

func TestProfilesCheckClean(t *testing.T) {
	for _, name := range initrepo.Profiles {
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "fresh", "proj") // neither level exists
			res, err := initrepo.Init(dir, profile(t, name))
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
			if name == "notes" {
				if repo.WorkDir() != "notes/work" || repo.ADRDir() != "notes/adr" || repo.RoadmapPath() != "notes/ROADMAP.md" {
					t.Errorf("notes folders: %s %s %s", repo.WorkDir(), repo.ADRDir(), repo.RoadmapPath())
				}
				if _, err := os.Stat(filepath.Join(dir, ".obsidian")); err == nil {
					t.Error("init must not write .obsidian")
				}
				if v := repo.Layout["reference_style"]; v.Str != "wikilink" {
					t.Errorf("reference_style = %q", v.Str)
				}
			} else if repo.WorkDir() != "docs/work" {
				t.Errorf("code work dir %s", repo.WorkDir())
			}
			files := snapshot(t, dir)
			if files["CLAUDE.md"] != "@AGENTS.md\n" {
				t.Errorf("CLAUDE.md = %q", files["CLAUDE.md"])
			}
			if !strings.Contains(files[".gitignore"], "docs/.akm.lock") {
				t.Errorf(".gitignore = %q", files[".gitignore"])
			}
			if !strings.HasPrefix(files["AGENTS.md"], "# proj\n") {
				t.Errorf("AGENTS.md title: %q", files["AGENTS.md"][:20])
			}
			if strings.Contains(files["AGENTS.md"], "{{") || strings.Contains(files["README.md"], "{{") {
				t.Error("an unreplaced marker")
			}
			if name == "code" && strings.Contains(files[".agents/skills/work-items/SKILL.md"], "notes/") {
				t.Error("code skill mentions notes/")
			}
			if name == "notes" && strings.Contains(files[".agents/skills/work-items/SKILL.md"], "docs/work") {
				t.Error("notes skill still points at docs/work")
			}
		})
	}
}

func TestNoOverwriteAndIdempotent(t *testing.T) {
	for _, name := range initrepo.Profiles {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			mine := "# mine\n\nkeep me\n"
			if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte(mine), 0o644); err != nil {
				t.Fatal(err)
			}
			res, err := initrepo.Init(dir, profile(t, name))
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
			res2, err := initrepo.Init(dir, profile(t, name))
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
	if _, err := initrepo.Init(f, profile(t, "code")); err == nil {
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
	files, _ := initrepo.Files(profile(t, "code"), "x")
	for _, f := range files {
		if f.Path == record.ManifestPath && !strings.HasPrefix(string(f.Data), "layout: "+check.LayoutVersion+"\n") {
			t.Errorf("manifest: %q", f.Data)
		}
	}
}

// The embedded skills are a copy of this repository's own.
func TestSkillsMatchRepository(t *testing.T) {
	src := snapshot(t, "../../.agents/skills")
	got := snapshot(t, "templates/skills")
	if !reflect.DeepEqual(src, got) {
		t.Error("internal/initrepo/templates/skills differs from .agents/skills: copy the skills over")
	}
}

// A cold agent can use the written files to list, claim, work and close an item.
func TestWorkLifecycle(t *testing.T) {
	for _, name := range initrepo.Profiles {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			if _, err := initrepo.Init(dir, profile(t, name)); err != nil {
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
		{"no profile", api.Request{}, "invalid_request"},
		{"unknown profile", api.Request{"profile": "wiki"}, "invalid_request"},
		{"dir climbs out", api.Request{"profile": "code", "dir": "../x"}, "invalid_request"},
		{"dir is absolute", api.Request{"profile": "code", "dir": "/tmp/x"}, "invalid_request"},
		{"unknown parameter", api.Request{"profile": "code", "force": true}, "invalid_request"},
	} {
		if _, aerr := api.Run(repo, info, "init", c.req); aerr == nil || aerr.Code != c.code {
			t.Errorf("%s: %v", c.name, aerr)
		}
	}
	resp, aerr := api.Run(repo, info, "init", api.Request{"profile": "code", "dir": "a/b"})
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
	if _, aerr := api.Run(repo, info, "init", api.Request{"profile": "code", "dir": "f"}); aerr == nil || aerr.Code != "refused" {
		t.Errorf("dir is a file: %v", aerr)
	}
}
