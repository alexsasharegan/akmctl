package api

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// The API fixtures have one source, docs/spec/akm-api/.fixtures, and a copy the tests read,
// testdata/api and testdata/yaml-lite.json (work item 261029-x4ws). This test fails when
// the copy differs from the source in any file or byte; scripts/sync-fixtures.sh fixes it.
func TestFixtureCopiesAgree(t *testing.T) {
	src := "../../docs/spec/akm-api/.fixtures"
	pairs := [][2]string{
		{filepath.Join(src, "cases"), "../../testdata/api/cases"},
		{filepath.Join(src, "repos"), "../../testdata/api/repos"},
		{filepath.Join(src, "yaml-lite.json"), "../../testdata/yaml-lite.json"},
	}
	var bad []string
	for _, p := range pairs {
		a, b := listFiles(t, p[0]), listFiles(t, p[1])
		for rel, ab := range a {
			bb, ok := b[rel]
			switch {
			case !ok:
				bad = append(bad, "missing from testdata: "+filepath.Join(p[1], rel))
			case !bytes.Equal(ab, bb):
				bad = append(bad, "differs: "+filepath.Join(p[1], rel))
			}
		}
		for rel := range b {
			if _, ok := a[rel]; !ok {
				bad = append(bad, "not in the source: "+filepath.Join(p[1], rel))
			}
		}
	}
	if len(bad) > 0 {
		sort.Strings(bad)
		if len(bad) > 20 {
			bad = append(bad[:20], "...")
		}
		t.Fatalf("testdata is not a copy of docs/spec/akm-api/.fixtures; edit the source and run scripts/sync-fixtures.sh:\n%s", join(bad))
	}
}

// listFiles reads a file, or every file under a directory, keyed by the path below it.
func listFiles(t *testing.T, root string) map[string][]byte {
	t.Helper()
	out := map[string][]byte{}
	info, err := os.Stat(root)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		b, err := os.ReadFile(root)
		if err != nil {
			t.Fatal(err)
		}
		out["."] = b
		return out
	}
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, rerr := os.ReadFile(p)
		if rerr != nil {
			return rerr
		}
		rel, _ := filepath.Rel(root, p)
		out[filepath.ToSlash(rel)] = b
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func join(l []string) string {
	s := ""
	for _, x := range l {
		s += "  " + x + "\n"
	}
	return s
}
