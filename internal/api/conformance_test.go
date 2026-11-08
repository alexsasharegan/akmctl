package api

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/record"
	"github.com/alexsasharegan/akmctl/internal/testutil"
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

// update makes TestConformance rewrite the expected response (and the after
// files) of a failing case from the tool's own output, in the source fixtures
// and in the copy, for a person to review in the diff:
//
//	go test ./internal/api -run TestConformance -update
//
// A case whose only difference is a diff member is not rewritten (its would_be
// text is by hand), and neither are the implementation members of version.
var update = flag.Bool("update", false, "rewrite the expected responses of the conformance fixtures from the tool's output")

// fixtureSources are the case directories -update writes: the copy the tests
// read and the source it comes from (scripts/sync-fixtures.sh keeps them equal).
var fixtureSources = []string{"../../testdata/api/cases", "../../docs/spec/akm-api/.fixtures/cases"}

// spliceMember replaces the value of a top-level member of a JSON object file
// with value, indented to sit at depth one, and keeps every other byte.
func spliceMember(raw []byte, key string, value any) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	if _, err := dec.Token(); err != nil {
		return nil, err
	}
	last := int64(0)
	add := func(value any) ([]byte, error) {
		vb, err := marshalFixture(value)
		if err != nil {
			return nil, err
		}
		var ind bytes.Buffer
		if err := json.Indent(&ind, vb, "  ", "  "); err != nil {
			return nil, err
		}
		kb, _ := json.Marshal(key)
		out := append([]byte(nil), raw[:last]...)
		out = append(out, ",\n  "...)
		out = append(out, kb...)
		out = append(out, ": "...)
		out = append(out, ind.Bytes()...)
		return append(out, raw[last:]...), nil
	}
	for dec.More() {
		k, err := dec.Token()
		if err != nil {
			return nil, err
		}
		start := dec.InputOffset()
		var skip json.RawMessage
		if err := dec.Decode(&skip); err != nil {
			return nil, err
		}
		end := dec.InputOffset()
		last = end
		if k != key {
			continue
		}
		vb, err := marshalFixture(value)
		if err != nil {
			return nil, err
		}
		var ind bytes.Buffer
		if err := json.Indent(&ind, vb, "  ", "  "); err != nil {
			return nil, err
		}
		out := append([]byte(nil), raw[:start]...)
		out = append(out, ": "...)
		out = append(out, ind.Bytes()...)
		return append(out, raw[end:]...), nil
	}
	return add(value) // a member the file does not have yet goes last
}

// marshalFixture writes a value without escaping <, > and & (the fixtures keep them).
func marshalFixture(v any) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(b.Bytes(), []byte("\n")), nil
}

// keepRequest returns the response got with the request member of the old
// expected response: the echo of a request is a map, so the tool writes its
// keys sorted, and a case keeps them in the order a person wrote them.
func keepRequest(got, old json.RawMessage) json.RawMessage {
	member := func(raw json.RawMessage, key string) json.RawMessage {
		var m map[string]json.RawMessage
		if json.Unmarshal(raw, &m) != nil {
			return nil
		}
		return m[key]
	}
	oldReq := member(old, "request")
	if oldReq == nil {
		return got
	}
	dec := json.NewDecoder(bytes.NewReader(got))
	dec.Token()
	var out bytes.Buffer
	out.WriteByte('{')
	for first := true; dec.More(); first = false {
		k, _ := dec.Token()
		var v json.RawMessage
		dec.Decode(&v)
		if k == "request" {
			v = oldReq
		}
		if !first {
			out.WriteByte(',')
		}
		kb, _ := json.Marshal(k)
		out.Write(kb)
		out.WriteByte(':')
		out.Write(v)
	}
	out.WriteByte('}')
	return out.Bytes()
}

// rewriteCase writes new members into a case file in every fixture directory.
func rewriteCase(t *testing.T, name string, members map[string]any) {
	t.Helper()
	for _, dir := range fixtureSources {
		p := filepath.Join(dir, name)
		raw, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		for _, k := range []string{"expected", "after", "would_be"} {
			v, ok := members[k]
			if !ok {
				continue
			}
			if raw, err = spliceMember(raw, k, v); err != nil {
				t.Fatalf("%s: %v", p, err)
			}
		}
		if err := os.WriteFile(p, raw, 0o644); err != nil {
			t.Fatal(err)
		}
	}
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
			testutil.CopyTree(t, filepath.Join("../../testdata/api/repos", c.Repo), root)
			if c.GitWorkTree {
				if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			before := testutil.Snapshot(t, root)
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
			var rewrite = map[string]any{}
			if *update {
				// a case with a diff or a version response is compared in parts, so
				// it is not rewritten from the whole response
				gm, wm := g.(map[string]any), w.(map[string]any)
				if gd, hasDiff := gm["diff"].(map[string]any); hasDiff {
					// the diff is checked by applying it, so compare the rest and
					// rewrite the would-be text from the applied diff
					wd, _ := wm["diff"].(map[string]any)
					np, _ := gd["new_path"].(string)
					applied, aerr := applyDiff(testutil.Snapshot(t, root)[gd["path"].(string)], gd)
					if aerr != nil {
						t.Fatalf("the diff does not apply (%v)", aerr)
					}
					rest := func(m map[string]any) map[string]any {
						o := map[string]any{}
						for k, v := range m {
							if k != "diff" {
								o[k] = v
							}
						}
						return o
					}
					if wd == nil || gd["path"] != wd["path"] || np != wd["new_path"] || c.WouldBe[np] != applied || !reflect.DeepEqual(rest(gm), rest(wm)) {
						rewrite["expected"] = keepRequest(gb, c.Expected)
						rewrite["would_be"] = map[string]string{np: applied}
					}
				} else if exp.Operation != "version" && !reflect.DeepEqual(g, w) {
					rewrite["expected"] = keepRequest(gb, c.Expected)
				}
			}
			if exp.Operation == "version" {
				// tool, api_versions, layout_versions and extensions describe the
				// implementation (spec section 9); only the rest is compared
				for _, m := range []any{g, w} {
					for _, k := range []string{"tool", "api_versions", "layout_versions", "extensions"} {
						delete(m.(map[string]any), k)
					}
				}
			}
			if gd, ok := g.(map[string]any)["diff"]; ok && rewrite["expected"] == nil {
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
				} else if got, err := applyDiff(testutil.Snapshot(t, root)[gdm["path"].(string)], gdm); err != nil || got != want {
					t.Fatalf("the diff does not produce the would-be file (%v):\n%s\nwant:\n%s", err, got, want)
				}
			}
			if !reflect.DeepEqual(g, w) && rewrite["expected"] == nil {
				wb, _ := json.MarshalIndent(w, "", " ")
				gbi, _ := json.MarshalIndent(g, "", " ")
				t.Fatalf("%s\ngot:\n%s\nwant:\n%s", c.Name, gbi, wb)
			}
			now := testutil.Snapshot(t, root)
			delete(now, "docs/.akm.lock")
			delete(before, "docs/.akm.lock")
			if *update {
				// after files that differ are rewritten too
				changed, newAfter := false, map[string]string{}
				for p, want := range c.After {
					newAfter[p] = now[p]
					if now[p] != want {
						changed = true
					}
				}
				if changed {
					rewrite["after"] = newAfter
				}
				if len(rewrite) > 0 {
					rewriteCase(t, filepath.Base(f), rewrite)
					t.Logf("%s: rewrote %d member(s) of %s", c.Name, len(rewrite), filepath.Base(f))
					return
				}
			}
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
