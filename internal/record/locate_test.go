package record

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func put(t *testing.T, root, rel, text string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLocate(t *testing.T) {
	cases := []struct {
		name  string
		files []string
		want  string
		found bool
		err   string
	}{
		{"none", []string{"README.md"}, "", false, ""},
		{"docs", []string{"docs/.layout"}, "docs", true, ""},
		{"other folder", []string{"knowledge/.layout"}, "knowledge", true, ""},
		{"docs wins over another", []string{"docs/.layout", "knowledge/.layout"}, "docs", true, ""},
		{"hidden and node_modules skipped", []string{".obsidian/.layout", "node_modules/x/.layout", "node_modules/.layout"}, "", false, ""},
		{"nested is not scanned", []string{"a/b/.layout"}, "", false, ""},
		{"two", []string{"notes/.layout", "wiki/.layout"}, "", false, "notes/.layout and wiki/.layout"},
		{"folder named .layout is not a manifest", []string{"x/.layout/y"}, "", false, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := t.TempDir()
			for _, f := range c.files {
				put(t, root, f, "layout: 0.3.0\n")
			}
			dir, found, err := Locate(root)
			if c.err != "" {
				if err == nil || !strings.Contains(err.Error(), c.err) {
					t.Fatalf("err = %v, want %q", err, c.err)
				}
				if _, lerr := Load(root); lerr == nil {
					t.Error("Load accepted two manifests")
				}
				return
			}
			if err != nil || dir != c.want || found != c.found {
				t.Fatalf("Locate = %q %v %v", dir, found, err)
			}
		})
	}
}

func TestLoadFollowsTheManifestFolder(t *testing.T) {
	root := t.TempDir()
	put(t, root, "knowledge/.layout", "layout: 0.3.0\nwork_dir: knowledge/work\nadr_dir: knowledge/adr\nroadmap: knowledge/ROADMAP.md\n")
	r, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if r.MetaDir() != "knowledge" || r.ManifestPath() != "knowledge/.layout" || r.DenyPath() != "knowledge/.akm-deny" || r.MigrateLogPath() != "knowledge/.akm-migrate-log" {
		t.Errorf("paths: %s %s %s", r.ManifestPath(), r.DenyPath(), r.MigrateLogPath())
	}
	if r.WorkDir() != "knowledge/work" {
		t.Errorf("work dir %s", r.WorkDir())
	}
	if len(r.Problems) != 0 {
		t.Errorf("problems %v", r.Problems)
	}
	// no manifest: the defaults
	empty, _ := Load(t.TempDir())
	if empty.MetaDir() != "docs" || empty.ManifestPath() != "docs/.layout" {
		t.Errorf("default %s", empty.ManifestPath())
	}
}

func TestFindRootFromSubdirectory(t *testing.T) {
	root := t.TempDir()
	put(t, root, "knowledge/.layout", "layout: 0.3.0\n")
	put(t, root, "knowledge/work/x.md", "x")
	put(t, root, "src/deep/y.go", "x")
	want, _ := filepath.EvalSymlinks(root)
	for _, from := range []string{".", "knowledge/work", "src/deep"} {
		got, err := FindRoot(filepath.Join(root, filepath.FromSlash(from)))
		if err != nil {
			t.Fatalf("%s: %v", from, err)
		}
		if g, _ := filepath.EvalSymlinks(got); g != want {
			t.Errorf("%s: root %s, want %s", from, g, want)
		}
	}
	// the nearest ancestor wins: a nested AGENTS.md marks its own root
	put(t, root, "src/AGENTS.md", "x")
	got, _ := FindRoot(filepath.Join(root, "src", "deep"))
	if filepath.Base(got) != "src" {
		t.Errorf("nearest root = %s", got)
	}
	// two manifests still mark the root; Load reports them
	two := t.TempDir()
	put(t, two, "a/.layout", "layout: 0.3.0\n")
	put(t, two, "b/.layout", "layout: 0.3.0\n")
	put(t, two, "a/sub/z.md", "x")
	got, err := FindRoot(filepath.Join(two, "a", "sub"))
	if err != nil || filepath.Base(got) != filepath.Base(two) {
		t.Fatalf("ambiguous root: %s %v", got, err)
	}
}
