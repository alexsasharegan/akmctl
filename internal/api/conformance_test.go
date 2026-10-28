package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/record"
)

type conformanceCase struct {
	Name     string            `json:"name"`
	Repo     string            `json:"repo"`
	Request  Request           `json:"request"`
	Expected json.RawMessage   `json:"expected"`
	After    map[string]string `json:"after"`
	Removed  []string          `json:"removed"`
}

// TestConformance runs every fixture case of the API spec whose operation is
// implemented and compares the response as JSON values.
func TestConformance(t *testing.T) {
	files, _ := filepath.Glob("../../testdata/api/cases/*.json")
	if len(files) == 0 {
		t.Fatal("no fixture cases")
	}
	ran := 0
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		var c conformanceCase
		if err := json.Unmarshal(raw, &c); err != nil {
			t.Fatal(f, err)
		}
		var exp struct{ Operation string }
		json.Unmarshal(c.Expected, &exp)
		if _, ok := specs[exp.Operation]; !ok {
			t.Logf("skip %s: operation %s not implemented yet", filepath.Base(f), exp.Operation)
			continue
		}
		t.Run(filepath.Base(f), func(t *testing.T) {
			// every case runs on a copy, so write cases cannot touch the fixtures
			root := t.TempDir()
			copyTree(t, filepath.Join("../../testdata/api/repos", c.Repo), root)
			before := snapshot(t, root)
			repo, err := record.Load(root)
			if err != nil {
				t.Fatal(err)
			}
			got, _ := Run(repo, Info{"akmctl", "test"}, exp.Operation, c.Request)
			gb, err := Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			var g, w any
			json.Unmarshal(gb, &g)
			json.Unmarshal(c.Expected, &w)
			if !reflect.DeepEqual(g, w) {
				wb, _ := json.MarshalIndent(w, "", " ")
				gbi, _ := json.MarshalIndent(g, "", " ")
				t.Fatalf("%s\ngot:\n%s\nwant:\n%s", c.Name, gbi, wb)
			}
			now := snapshot(t, root)
			delete(now, "docs/.akm.lock")
			delete(before, "docs/.akm.lock")
			if c.After == nil && c.Removed == nil {
				if !reflect.DeepEqual(before, now) {
					t.Fatalf("%s: files changed but the case expects none", c.Name)
				}
				return
			}
			for p, want := range c.After {
				if now[p] != want {
					t.Errorf("%s: %s is\n%s\nwant\n%s", c.Name, p, now[p], want)
				}
			}
			for _, p := range c.Removed {
				if _, ok := now[p]; ok {
					t.Errorf("%s: %s should be gone", c.Name, p)
				}
			}
			for p := range now {
				if _, ok := before[p]; !ok {
					if _, expected := c.After[p]; !expected {
						t.Errorf("%s: unexpected new file %s", c.Name, p)
					}
				}
			}
		})
		ran++
	}
	t.Logf("ran %d of %d cases", ran, len(files))
}

func copyTree(t *testing.T, src, dst string) {
	t.Helper()
	filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		rel, _ := filepath.Rel(src, p)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		b, _ := os.ReadFile(p)
		return os.WriteFile(filepath.Join(dst, rel), b, 0o644)
	})
}

// snapshot maps every file under root to its text.
func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	m := map[string]string{}
	filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			rel, _ := filepath.Rel(root, p)
			b, _ := os.ReadFile(p)
			m[filepath.ToSlash(rel)] = string(b)
		}
		return nil
	})
	return m
}
