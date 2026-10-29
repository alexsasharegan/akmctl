package record

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckRelPath(t *testing.T) {
	for _, ok := range []string{"docs/work", "Deep Thought/9 pkm/tasks", "tasks", "a/b/c.md", "é ü/x"} {
		if err := CheckRelPath(ok); err != nil {
			t.Errorf("%q: %v", ok, err)
		}
	}
	for _, bad := range []string{"", "/abs/path", "../out", "a/../b", "a/./b", "a//b", "a/", "C:/x", "a\\b", "a/b\n", "./a"} {
		if err := CheckRelPath(bad); err == nil {
			t.Errorf("%q was accepted", bad)
		}
	}
}

func TestDefaultsAndAccessors(t *testing.T) {
	r, err := Load("../../testdata/api/repos/basic")
	if err != nil {
		t.Fatal(err)
	}
	if r.WorkDir() != "docs/work" || r.DoneDir() != "docs/work/done" || r.ADRDir() != "docs/adr" || r.RoadmapPath() != "docs/ROADMAP.md" || r.ADRIndexPath() != "docs/adr/README.md" {
		t.Fatalf("defaults: %s %s %s %s", r.WorkDir(), r.DoneDir(), r.ADRDir(), r.RoadmapPath())
	}
	var empty Repo
	if empty.WorkDir() != "docs/work" || empty.RoadmapPath() != "docs/ROADMAP.md" {
		t.Fatal("a zero Repo must give the defaults")
	}
}

func TestLoadCustomFolders(t *testing.T) {
	r, err := Load("../../testdata/api/repos/vault")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Problems) != 0 {
		t.Fatalf("problems %v", r.Problems)
	}
	if r.WorkDir() != "Deep Thought/9 pkm/tasks" || r.ADRDir() != "Deep Thought/9 pkm/adr" || r.RoadmapPath() != "Deep Thought/9 pkm/ROADMAP.md" {
		t.Fatalf("%s %s %s", r.WorkDir(), r.ADRDir(), r.RoadmapPath())
	}
	if len(r.Records) != 8 {
		t.Fatalf("%d records", len(r.Records))
	}
	rec := r.ByID("261010-j9k0")[0]
	if rec.Path != "Deep Thought/9 pkm/tasks/done/261010-j9k0-choose-a-licence.md" || !r.IsDone(rec.Path) {
		t.Fatalf("done item at %s", rec.Path)
	}
	if len(r.Roadmap) != 2 || r.Roadmap[0] != "261015-c4d5" {
		t.Fatalf("roadmap %v", r.Roadmap)
	}
	if !r.InWork(rec.Path) || r.InADR(rec.Path) {
		t.Fatal("InWork / InADR")
	}
}

func layoutRepo(t *testing.T, layout string) *Repo {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", ".layout"), []byte(layout), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestBadManifestPathsAreProblemsAndDefaultsHold(t *testing.T) {
	for _, c := range []struct{ name, layout, key, msg string }{
		{"parent dir", "work_dir: ../elsewhere\n", "work_dir", `".."`},
		{"absolute", "adr_dir: /etc/adr\n", "adr_dir", "absolute"},
		{"backslash", "work_dir: a\\b\n", "work_dir", "forward slashes"},
		{"drive", "roadmap: C:/r.md\n", "roadmap", "forward slashes"},
		{"not a string", "work_dir: 5\n", "work_dir", "path string"},
		{"roadmap not markdown", "roadmap: notes/plan.txt\n", "roadmap", ".md"},
		{"same folder", "work_dir: x\nadr_dir: x\n", "adr_dir", "inside one another"},
		{"nested", "work_dir: x\nadr_dir: x/adr\n", "adr_dir", "inside one another"},
		{"roadmap inside work", "work_dir: w\nroadmap: w/RM.md\n", "roadmap", "inside work_dir"},
	} {
		r := layoutRepo(t, c.layout)
		var got *Problem
		for i := range r.Problems {
			if r.Problems[i].Rule == "S3" {
				got = &r.Problems[i]
			}
		}
		if got == nil || got.Path != ManifestPath || got.Line == 0 || !strings.Contains(got.Msg, c.msg) || !strings.Contains(got.Msg, c.key) {
			t.Errorf("%s: problem %+v (all %v)", c.name, got, r.Problems)
			continue
		}
		if c.key == "roadmap" && r.RoadmapPath() == DefaultRoadmap {
			continue // the other keys of this case are valid and stay
		}
		if r.WorkDir() != DefaultWorkDir || r.ADRDir() != DefaultADRDir || r.RoadmapPath() != DefaultRoadmap {
			t.Errorf("%s: defaults not restored: %s %s %s", c.name, r.WorkDir(), r.ADRDir(), r.RoadmapPath())
		}
	}
	r := layoutRepo(t, "work_dir: tasks\n")
	if len(r.Problems) != 0 || r.WorkDir() != "tasks" || r.DoneDir() != "tasks/done" || r.ADRDir() != DefaultADRDir {
		t.Fatalf("one key alone: %v %s", r.Problems, r.WorkDir())
	}
}
