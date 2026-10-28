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
	Name     string          `json:"name"`
	Repo     string          `json:"repo"`
	Request  Request         `json:"request"`
	Expected json.RawMessage `json:"expected"`
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
			repo, err := record.Load(filepath.Join("../../testdata/api/repos", c.Repo))
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
		})
		ran++
	}
	t.Logf("ran %d of %d cases", ran, len(files))
}
