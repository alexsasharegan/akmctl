package frontmatter

import (
	"strings"
	"testing"
)

// FuzzAgainstOracle: whatever Parse accepts, a YAML 1.2 core parser must read the
// same way (the conformance rule of the spec). Run with
// go test -fuzz=FuzzAgainstOracle -fuzztime=30s ./internal/frontmatter
func FuzzAgainstOracle(f *testing.F) {
	for _, fx := range loadFixturesF(f) {
		f.Add(fx)
	}
	f.Fuzz(func(t *testing.T, block string) {
		if strings.Contains(block, "---") || strings.Contains(block, "\r") {
			return
		}
		d, err := Parse([]byte("---\n" + block + "\n---\n"))
		if err != nil {
			return
		}
		want, err := oracle(block)
		if err != nil {
			t.Fatalf("tool accepts %q, oracle fails: %v", block, err)
		}
		if len(want) != len(d.Entries) {
			t.Fatalf("%q: %d keys, oracle %d", block, len(d.Entries), len(want))
		}
		for _, e := range d.Entries {
			if w, ok := want[e.Key]; !ok || !e.Value.Equal(w) {
				t.Fatalf("%q: %s = %v, oracle %v", block, e.Key, e.Value, w)
			}
		}
		// an edit that changes nothing must change no byte
		for _, e := range d.Entries {
			before := string(d.Bytes())
			if err := d.Set(e.Key, e.Value); err != nil {
				t.Fatalf("%q: no-op set of %s: %v", block, e.Key, err)
			}
			if string(d.Bytes()) != before {
				t.Fatalf("%q: no-op set changed bytes", block)
			}
		}
	})
}

func loadFixturesF(f *testing.F) []string {
	var out []string
	for _, fx := range loadFixturesAny(f) {
		out = append(out, fx.Input)
	}
	return out
}
