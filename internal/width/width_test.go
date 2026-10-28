package width

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
)

type cell struct {
	ID      int    `json:"id"`
	Cls     string `json:"cls"`
	Hex     string `json:"hex"`
	W       int    `json:"w"`
	Changed bool   `json:"changed"`
}

// TestPrettier compares String with the widths Prettier 3.9.10 chose when it
// padded one-cell tables (cells it rewrote are not scored). Some classes are
// known to differ; the allowed misses are listed so that a new one fails.
func TestPrettier(t *testing.T) {
	raw, err := os.ReadFile("testdata/prettier-3.9.10.json")
	if err != nil {
		t.Fatal(err)
	}
	var cells []cell
	if err := json.Unmarshal(raw, &cells); err != nil {
		t.Fatal(err)
	}
	// ids of cells known to differ from Prettier (see the evaluation note); keep this list short
	known := map[int]bool{148: true} // a ZWJ pair that is not a standard emoji sequence; Prettier does not join it, and the tool has no list of standard sequences
	scored, wrong := 0, 0
	for _, c := range cells {
		if c.Changed || c.Cls == "invalid-utf8" { // records are valid UTF-8; the parser rejects anything else
			continue
		}
		b, _ := hex.DecodeString(c.Hex)
		s := string(b)
		scored++
		if got := String(s); got != c.W {
			wrong++
			if !known[c.ID] {
				t.Errorf("cell %d (%s) %+q: width %d, Prettier %d", c.ID, c.Cls, shorten(s), got, c.W)
			}
		}
	}
	t.Logf("%d scored, %d differ from Prettier", scored, wrong)
}

func shorten(s string) string {
	if r := []rune(s); len(r) > 24 {
		return string(r[:24]) + "..."
	}
	return s
}

func TestBasics(t *testing.T) {
	for s, w := range map[string]int{"": 0, "abc": 3, "漢字": 4, "café": 4, "é": 1, "😀": 2, "✔️": 2, "á": 1, "👨‍👩‍👧‍👦": 2, "🇺🇸": 2, "a\tb": 2} {
		if got := String(s); got != w {
			t.Errorf("String(%+q) = %d, want %d", s, got, w)
		}
	}
}
