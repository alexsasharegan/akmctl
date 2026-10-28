package frontmatter

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

type fixture struct {
	Name   string          `json:"name"`
	Input  string          `json:"input"`
	Accept bool            `json:"accept"`
	Value  json.RawMessage `json:"value,omitempty"` // for accepted single-key inputs: the expected value of the key
	Key    string          `json:"key,omitempty"`
}

func loadFixtures(t *testing.T) []fixture { return loadFixturesAny(t) }

func loadFixturesAny(t testing.TB) []fixture {
	t.Helper()
	raw, err := os.ReadFile("../../testdata/yaml-lite.json")
	if err != nil {
		t.Fatal(err)
	}
	var fs []fixture
	if err := json.Unmarshal(raw, &fs); err != nil {
		t.Fatal(err)
	}
	return fs
}

func wrap(in string) []byte { return []byte("---\n" + in + "\n---\nbody\n") }

func toJSON(v Value) any {
	switch v.Kind {
	case String:
		return v.Str
	case Int:
		return float64(v.Int)
	case Bool:
		return v.Bool
	case List:
		out := []any{}
		for _, it := range v.Items {
			out = append(out, toJSON(it))
		}
		return out
	}
	return nil
}

func TestFixtures(t *testing.T) {
	for _, f := range loadFixtures(t) {
		t.Run(f.Name, func(t *testing.T) {
			d, err := Parse(wrap(f.Input))
			if !f.Accept {
				if err == nil {
					t.Fatalf("accepted %q; want an error", f.Input)
				}
				if _, ok := err.(*Error); !ok {
					t.Fatalf("error is %T, want *Error", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("rejected %q: %v", f.Input, err)
			}
			if f.Key != "" {
				v, ok := d.Get(f.Key)
				if !ok {
					t.Fatalf("key %s missing", f.Key)
				}
				var want any
				if err := json.Unmarshal(f.Value, &want); err != nil {
					t.Fatal(err)
				}
				gb, _ := json.Marshal(toJSON(v))
				wb, _ := json.Marshal(want)
				if string(gb) != string(wb) {
					t.Fatalf("%s = %s, want %s", f.Key, gb, wb)
				}
			}
		})
	}
}

func TestErrorLineNumbers(t *testing.T) {
	_, err := Parse([]byte("---\na: 1\nb: 02\n---\n"))
	e, ok := err.(*Error)
	if !ok || e.Line != 3 {
		t.Fatalf("got %v, want an error on line 3", err)
	}
}

func TestStructure(t *testing.T) {
	bad := map[string]string{
		"no opening":     "a: 1\n---\n",
		"unterminated":   "---\na: 1\n",
		"bom":            "\xef\xbb\xbf---\na: 1\n---\n",
		"tab":            "---\na:\t1\n---\n",
		"bad escape":     "---\na: \"x\\ny\"\n---\n",
		"unclosed quote": "---\na: \"x\n---\n",
		"invalid utf8":   "---\na: \xff\n---\n",
		"orphan item":    "---\n- a\n---\n",
		"item after val": "---\na: 1\n- b\n---\n",
		"bool key":       "---\nTrue: 1\n---\n",
		"empty item":     "---\na:\n  - \n---\n",
		"text after":     "---\na: \"x\" y\n---\n",
	}
	for name, in := range bad {
		if _, err := Parse([]byte(in)); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
	d, err := Parse([]byte("---\r\na: 1\r\nb: [x, y]\r\n---\r\nbody\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	if string(d.Body()) != "body\r\n" {
		t.Fatalf("body = %q", d.Body())
	}
	if v, _ := d.Get("b"); len(v.Items) != 2 {
		t.Fatalf("b = %v", v)
	}
	if !strings.Contains(string(d.Bytes()), "\r\n") {
		t.Fatal("CRLF lost")
	}
}

func TestBare(t *testing.T) {
	d, err := ParseBare([]byte("layout: 1.0.0\nstale_days: 30\n"))
	if err != nil {
		t.Fatal(err)
	}
	if v, _ := d.Get("layout"); v.Str != "1.0.0" {
		t.Fatalf("layout = %v", v)
	}
	if v, _ := d.Get("stale_days"); v.Int != 30 {
		t.Fatalf("stale_days = %v", v)
	}
}
