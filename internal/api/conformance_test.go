package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/record"
)

type conformanceCase struct {
	Name     string            `json:"name"`
	Repo     string            `json:"repo"`
	Request  Request           `json:"request"`
	Expected json.RawMessage   `json:"expected"`
	After    map[string]string `json:"after"`
	// AfterPaths are files that must exist afterwards, content not compared
	// (init: the wording of a template is not part of the contract).
	AfterPaths []string          `json:"after_paths"`
	Removed    []string          `json:"removed"`
	WouldBe    map[string]string `json:"would_be"`
	// GitWorkTree makes the copy a work tree of the version control system
	// (an empty .git directory).
	GitWorkTree bool `json:"git_work_tree"`
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
			// the directory is named like the fixture repository: version reports it
			root := filepath.Join(t.TempDir(), c.Repo)
			copyTree(t, filepath.Join("../../testdata/api/repos", c.Repo), root)
			if c.GitWorkTree {
				if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
					t.Fatal(err)
				}
			}
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
			if exp.Operation == "version" {
				// tool, api_versions, layout_versions and extensions describe the
				// implementation (spec section 9); only the rest is compared
				for _, m := range []any{g, w} {
					for _, k := range []string{"tool", "api_versions", "layout_versions", "extensions"} {
						delete(m.(map[string]any), k)
					}
				}
			}
			if gd, ok := g.(map[string]any)["diff"]; ok {
				// equal diffs can differ in where a run of lines sits, so a diff is
				// checked by applying it (spec section 9)
				wd, _ := w.(map[string]any)["diff"].(map[string]any)
				delete(g.(map[string]any), "diff")
				delete(w.(map[string]any), "diff")
				gdm, _ := gd.(map[string]any)
				if gdm["path"] != wd["path"] || gdm["new_path"] != wd["new_path"] {
					t.Fatalf("diff paths %v -> %v, want %v -> %v", gdm["path"], gdm["new_path"], wd["path"], wd["new_path"])
				}
				if want, ok := c.WouldBe[gdm["new_path"].(string)]; !ok {
					t.Fatal("the case has a diff but no would_be text")
				} else if got, err := applyDiff(snapshot(t, root)[gdm["path"].(string)], gdm); err != nil || got != want {
					t.Fatalf("the diff does not produce the would-be file (%v):\n%s\nwant:\n%s", err, got, want)
				}
			}
			if !reflect.DeepEqual(g, w) {
				wb, _ := json.MarshalIndent(w, "", " ")
				gbi, _ := json.MarshalIndent(g, "", " ")
				t.Fatalf("%s\ngot:\n%s\nwant:\n%s", c.Name, gbi, wb)
			}
			now := snapshot(t, root)
			delete(now, "docs/.akm.lock")
			delete(before, "docs/.akm.lock")
			if c.After == nil && c.Removed == nil && c.AfterPaths == nil {
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
			for _, p := range c.AfterPaths {
				if _, ok := now[p]; !ok {
					t.Errorf("%s: %s should exist", c.Name, p)
				}
			}
			for _, p := range c.Removed {
				if _, ok := now[p]; ok {
					t.Errorf("%s: %s should be gone", c.Name, p)
				}
			}
			for p := range now {
				if _, ok := before[p]; !ok {
					if _, expected := c.After[p]; !expected && !contains(c.AfterPaths, p) {
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

// applyDiff applies the hunks of a diff member to the old file text, checking
// that every context and removed line is what the file has and that the line
// numbers and headers agree. It returns the new text.
func applyDiff(old string, d map[string]any) (string, error) {
	lines := strings.Split(strings.TrimSuffix(old, "\n"), "\n")
	if old == "" {
		lines = nil
	}
	var out []string
	pos := 0 // lines of old consumed
	hunks, _ := d["hunks"].([]any)
	for _, h := range hunks {
		hm := h.(map[string]any)
		var oldN, newN, oldStart, newStart int
		for _, l := range hm["lines"].([]any) {
			lm := l.(map[string]any)
			op, text := lm["op"].(string), lm["text"].(string)
			if op != "add" {
				n := int(lm["old"].(float64))
				if oldN == 0 {
					oldStart = n
					for pos < n-1 {
						out = append(out, lines[pos])
						pos++
					}
				}
				if n != pos+1 || lines[pos] != text {
					return "", fmt.Errorf("old line %d is not %q", n, text)
				}
				pos++
				oldN++
			} else if lm["old"] != nil {
				return "", fmt.Errorf("an added line has an old number")
			}
			if op != "rem" {
				n := int(lm["new"].(float64))
				if newN == 0 {
					newStart = n
				}
				if n != len(out)+1 {
					return "", fmt.Errorf("new line number %d, expected %d", n, len(out)+1)
				}
				out = append(out, text)
				newN++
			} else if lm["new"] != nil {
				return "", fmt.Errorf("a removed line has a new number")
			}
		}
		span := func(s, n int) string {
			if n == 1 {
				return fmt.Sprint(s)
			}
			return fmt.Sprintf("%d,%d", s, n)
		}
		if want := "@@ -" + span(oldStart, oldN) + " +" + span(newStart, newN) + " @@"; hm["header"] != want {
			return "", fmt.Errorf("header %v, want %s", hm["header"], want)
		}
	}
	out = append(out, lines[pos:]...)
	return strings.Join(out, "\n") + "\n", nil
}
