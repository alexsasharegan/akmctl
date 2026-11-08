package cli

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/testutil"
)

var updateGolden = flag.Bool("update", false, "rewrite the golden files of the text view from the tool's output")

// A new id is the date of --as-of and four random characters; the golden files
// keep the date and write the characters as xxxx.
var reNewID = regexp.MustCompile(`261028-[0-9a-hjkmnp-tv-z]{4}`)

// TestTextView pins the text a person or an agent reads: the table view of each
// command, the commit message block and the progress counts of a parent. The
// JSON view is tested elsewhere (work item 261108-qv1w, finding F15). To change
// a golden file on purpose run
//
//	go test ./internal/cli -run TestTextView -update
//
// and review the diff.
func TestTextView(t *testing.T) {
	cases := []struct {
		name  string
		repo  string
		write bool // the case runs on a copy it may change
		args  []string
	}{
		{"ls", "basic", false, []string{"ls"}},
		{"ls-tree", "basic", false, []string{"ls", "--tree"}},
		{"ls-ready-fields", "basic", false, []string{"ls", "--status", "ready,active", "--fields", "title,priority"}},
		{"ls-count-only", "basic", false, []string{"ls", "--type", "work", "--count-only"}},
		{"show-parent", "basic", false, []string{"show", "261015-c4d5"}},
		{"show-child-section", "basic", false, []string{"show", "261022-e5f6", "--section", "Goal"}},
		{"board", "basic", false, []string{"board"}},
		{"search", "basic", false, []string{"search", "JSON"}},
		{"check-broken", "broken", false, []string{"check", "--as-of", "2026-10-28"}},
		{"close-dry-run", "basic", false, []string{"close", "261020-a2b3", "done", "--outcome", "Shipped.", "--as-of", "2026-10-28", "--dry-run"}},
		{"close-dropped", "closable", false, []string{"close", "261020-a2b3", "dropped", "--outcome", "Not needed.", "--as-of", "2026-10-28", "--dry-run"}},
		{"set-dry-run", "basic", false, []string{"set", "261022-e5f6", "priority=1", "--as-of", "2026-10-28", "--dry-run"}},
		{"claim-dry-run", "claims", false, []string{"claim", "261022-e5f6", "--as", "agent-1", "--as-of", "2026-10-28T09:30", "--dry-run"}},
		{"split", "basic", true, []string{"split", "261015-c4d5", "parse-flags", "write-docs", "--kind", "chore", "--as-of", "2026-10-28"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), c.repo)
			testutil.CopyTree(t, filepath.Join("../../testdata/api/repos", c.repo), dir)
			var out, errb bytes.Buffer
			code := Run(append(append([]string(nil), c.args...), "--repo", dir), &out, &errb, ".")
			got := "exit " + string(rune('0'+code)) + "\n" + out.String()
			if errb.Len() > 0 {
				got += "--- stderr\n" + errb.String()
			}
			got = strings.ReplaceAll(got, dir, "<repo>")
			got = reNewID.ReplaceAllString(got, "261028-xxxx")
			golden := filepath.Join("testdata", "golden", c.name+".txt")
			if *updateGolden {
				os.MkdirAll(filepath.Dir(golden), 0o755)
				if err := os.WriteFile(golden, []byte(got), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("%v (run with -update to write it)", err)
			}
			if got != string(want) {
				t.Errorf("text view of %v changed:\n--- got\n%s--- want\n%s", c.args, got, want)
			}
		})
	}
}
