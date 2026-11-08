package conformance

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/api"
	"github.com/alexsasharegan/akmctl/internal/testutil"
)

type gitCommit struct {
	Message string `json:"message"`
	Date    string `json:"date"`
	// Files maps a path to its new text; null deletes the file.
	Files map[string]*string `json:"files"`
}

type conformanceCase struct {
	Name     string            `json:"name"`
	Repo     string            `json:"repo"`
	Request  api.Request       `json:"request"`
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
	// Git is a history the runner replays into the copy (ADR 261108-15dw).
	Git []gitCommit `json:"git"`
	// Bindings restricts the case to some bindings (api, http, cli).
	Bindings []string `json:"bindings"`
}

const (
	casesDir = "../../testdata/api/cases"
	reposDir = "../../testdata/api/repos"
)

// update makes TestConformance rewrite the expected response (and the after
// files) of a failing case from the tool's own output, in the source fixtures
// and in the copy, for a person to review in the diff. Matchers in the old
// expected response stay where the new output still fits them. Only the api
// binding runs, so run the test again without the flag afterwards:
//
//	go test ./internal/conformance -run TestConformance -update
var update = flag.Bool("update", false, "rewrite the expected responses of the conformance fixtures from the tool's output")

func TestConformance(t *testing.T) {
	files, _ := filepath.Glob(filepath.Join(casesDir, "*.json"))
	if len(files) == 0 {
		t.Fatal("no fixture cases")
	}
	impl := api.Operations()
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		var c conformanceCase
		if err := json.Unmarshal(raw, &c); err != nil {
			t.Fatal(f, err)
		}
		// the operation is the file name's prefix (set-15.json runs set): the
		// response of a refusal that never reached the operation names "http"
		op, _, _ := strings.Cut(strings.TrimSuffix(filepath.Base(f), ".json"), "-")
		if !slices.Contains(impl, op) {
			t.Logf("skip %s: operation %s not implemented yet", filepath.Base(f), op)
			continue
		}
		for _, bd := range bindings {
			if len(c.Bindings) > 0 && !slices.Contains(c.Bindings, bd.name) {
				continue
			}
			if *update && bd.name != firstBinding(&c) {
				continue
			}
			t.Run(strings.TrimSuffix(filepath.Base(f), ".json")+"/"+bd.name, func(t *testing.T) {
				runCase(t, filepath.Base(f), &c, op, bd)
			})
		}
	}
}

// snapshot is the files under root without the version control directory (git
// may still write there, for example its automatic maintenance) and the lock file.
func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	m := map[string]string{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if rel = filepath.ToSlash(rel); rel == "docs/.akm.lock" {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		m[rel] = string(b)
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot %s: %v", root, err)
	}
	return m
}

func replayGit(t *testing.T, root string, commits []gitCommit) {
	t.Helper()
	testutil.RequireTool(t, "git")
	git := func(env []string, args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), append([]string{
			"GIT_CONFIG_GLOBAL=" + os.DevNull, "GIT_CONFIG_SYSTEM=" + os.DevNull, "GIT_CONFIG_NOSYSTEM=1",
			"GIT_AUTHOR_NAME=Fixture Author", "GIT_AUTHOR_EMAIL=author@example.com",
			"GIT_COMMITTER_NAME=Fixture Author", "GIT_COMMITTER_EMAIL=author@example.com",
		}, env...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	git(nil, "init", "-q", "-b", "main")
	git(nil, "config", "commit.gpgsign", "false")
	git(nil, "config", "core.autocrlf", "false")
	git(nil, "config", "gc.auto", "0")
	git(nil, "config", "maintenance.auto", "false")
	for _, c := range commits {
		var paths []string
		for p, text := range c.Files {
			full := filepath.Join(root, filepath.FromSlash(p))
			if text == nil {
				os.Remove(full)
			} else {
				if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(full, []byte(*text), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			paths = append(paths, p)
		}
		slices.Sort(paths)
		git(nil, append([]string{"add", "-A", "--"}, paths...)...)
		git([]string{"GIT_AUTHOR_DATE=" + c.Date, "GIT_COMMITTER_DATE=" + c.Date},
			"commit", "-q", "--allow-empty", "-m", c.Message)
	}
}

func runCase(t *testing.T, file string, c *conformanceCase, op string, bd binding) {
	// every case runs on a copy, so write cases cannot touch the fixtures
	// the directory is named like the fixture repository: version reports it
	root := filepath.Join(t.TempDir(), c.Repo)
	testutil.CopyTree(t, filepath.Join(reposDir, c.Repo), root)
	if c.GitWorkTree {
		if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if len(c.Git) > 0 {
		replayGit(t, root, c.Git)
	}
	before := snapshot(t, root)
	req := map[string]any{}
	for k, v := range c.Request {
		req[k] = v
	}
	gb, injected, skip := bd.run(t, root, op, req, noIfMatch(c.Expected))
	if skip != "" {
		t.Skip(skip)
	}
	var g, w any
	if err := json.Unmarshal(gb, &g); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, gb)
	}
	json.Unmarshal(c.Expected, &w)
	if gm, ok := g.(map[string]any); ok {
		if rq, ok := gm["request"].(map[string]any); ok {
			for _, k := range injected {
				delete(rq, k)
			}
		}
	}
	now := snapshot(t, root)
	b := binds{}
	problem := compare(op, g, w, before, c, b)
	if *update {
		updateCase(t, file, c, op, gb, g, w, before, now, problem, binds{})
		return
	}
	if problem != nil {
		t.Fatalf("%s: %v", c.Name, problem)
	}
	checkFiles(t, c, b, before, now)
}

// compare checks the response against the expected one: the diff by applying it
// (equal diffs can differ in where a run of lines sits, spec section 9), the
// implementation members of version not at all, the rest with matchers.
func compare(op string, g, w any, before map[string]string, c *conformanceCase, b binds) error {
	gm, gok := g.(map[string]any)
	wm, wok := w.(map[string]any)
	if !gok || !wok {
		return match(w, g, "", b)
	}
	gc, wc := copyMap(gm), copyMap(wm)
	if op == "version" {
		// tool, api_versions, layout_versions and extensions describe the
		// implementation (spec section 9); only the rest is compared
		for _, k := range []string{"tool", "api_versions", "layout_versions", "extensions"} {
			delete(gc, k)
			delete(wc, k)
		}
	}
	gd, gHas := gm["diff"].(map[string]any)
	wd, wHas := wm["diff"].(map[string]any)
	if gHas && wHas {
		delete(gc, "diff")
		delete(wc, "diff")
		np, _ := gd["new_path"].(string)
		if gd["path"] != wd["path"] || !matchString(fmt.Sprint(wd["new_path"]), np, b) {
			return fmt.Errorf("diff paths %v -> %v, want %v -> %v", gd["path"], gd["new_path"], wd["path"], wd["new_path"])
		}
		want, ok := lookupBound(c.WouldBe, np, b)
		if !ok {
			return fmt.Errorf("the case has a diff but no would_be text for %s", np)
		}
		old, _ := gd["path"].(string)
		got, err := applyDiff(before[old], gd)
		if err != nil || !matchString(want, got, b) {
			return fmt.Errorf("the diff does not produce the would-be file (%v):\n%s\nwant:\n%s", err, got, want)
		}
	}
	return match(wc, gc, "", b)
}

func copyMap(m map[string]any) map[string]any {
	o := make(map[string]any, len(m))
	for k, v := range m {
		o[k] = v
	}
	return o
}

// lookupBound finds the entry of m whose key, with the bound markers replaced,
// is path (or matches it when it still holds a marker).
func lookupBound(m map[string]string, path string, b binds) (string, bool) {
	for k, v := range m {
		if k == path || matchString(k, path, b) {
			return v, true
		}
	}
	return "", false
}

// checkFiles compares the files after the run with the case's after, after_paths
// and removed members, and fails on a file the case does not name.
func checkFiles(t *testing.T, c *conformanceCase, b binds, before, now map[string]string) {
	t.Helper()
	if c.After == nil && c.Removed == nil && c.AfterPaths == nil {
		if !reflect.DeepEqual(before, now) {
			t.Fatalf("%s: files changed but the case expects none", c.Name)
		}
		return
	}
	named := map[string]bool{}
	find := func(pat string) (string, bool) {
		if _, ok := now[pat]; ok && !hasMarker(pat) {
			return pat, true
		}
		for p := range now {
			if !named[p] && matchString(pat, p, b) {
				return p, true
			}
		}
		return "", false
	}
	paths := make([]string, 0, len(c.After))
	for p := range c.After {
		paths = append(paths, p)
	}
	slices.Sort(paths)
	for _, p := range paths {
		got, ok := find(p)
		if !ok {
			t.Errorf("%s: %s should exist", c.Name, p)
			continue
		}
		named[got] = true
		if !matchString(c.After[p], now[got], b) {
			t.Errorf("%s: %s is\n%s\nwant\n%s", c.Name, got, now[got], c.After[p])
		}
	}
	for _, p := range c.AfterPaths {
		if strings.HasSuffix(p, "/") {
			// a folder: it holds at least one new file, and every new file in it
			// is named (a shipped skill's file list is not part of the contract)
			found := false
			for q := range now {
				if _, old := before[q]; !old && strings.HasPrefix(q, subst(p, b)) {
					named[q], found = true, true
				}
			}
			if !found {
				t.Errorf("%s: %s should hold a new file", c.Name, p)
			}
		} else if got, ok := find(p); !ok {
			t.Errorf("%s: %s should exist", c.Name, p)
		} else {
			named[got] = true
		}
	}
	for _, p := range c.Removed {
		if _, ok := now[subst(p, b)]; ok {
			t.Errorf("%s: %s should be gone", c.Name, p)
		}
	}
	for p := range now {
		if _, ok := before[p]; !ok && !named[p] {
			t.Errorf("%s: unexpected new file %s", c.Name, p)
		}
	}
}

// ---- update mode ------------------------------------------------------

var fixtureSources = []string{"../../testdata/api/cases", "../../docs/spec/akm-api/.fixtures/cases"}

func updateCase(t *testing.T, file string, c *conformanceCase, op string, gb []byte, g, w any, before, now map[string]string, problem error, b binds) {
	t.Helper()
	rewrite := map[string]any{}
	if problem != nil && op != "version" {
		rewrite["expected"] = keepRequest(rebuild(gb, w, b), c.Expected)
		if gd, ok := g.(map[string]any)["diff"].(map[string]any); ok {
			np, _ := gd["new_path"].(string)
			old, _ := gd["path"].(string)
			applied, err := applyDiff(before[old], gd)
			if err != nil {
				t.Fatalf("the diff does not apply (%v)", err)
			}
			wb := map[string]string{}
			for k, v := range c.WouldBe {
				wb[k] = v
			}
			wb[np] = applied
			rewrite["would_be"] = wb
		}
	}
	if len(c.After) > 0 {
		changed, na := false, map[string]string{}
		for p, want := range c.After {
			if hasMarker(p) || hasMarker(want) {
				na[p] = want // a case with matchers is kept as written
				continue
			}
			na[p] = now[p]
			changed = changed || now[p] != want
		}
		if changed {
			rewrite["after"] = na
		}
	}
	if len(rewrite) == 0 {
		return
	}
	rewriteCase(t, file, rewrite)
	t.Logf("%s: rewrote %d member(s) of %s", c.Name, len(rewrite), file)
}

// rebuild returns the actual response bytes with every string that fits a
// matcher of the old expected value replaced by that matcher.
func rebuild(act json.RawMessage, exp any, b binds) json.RawMessage {
	dec := json.NewDecoder(bytes.NewReader(act))
	var walk func(exp any) []byte
	walk = func(exp any) []byte {
		tok, err := dec.Token()
		if err != nil {
			return nil
		}
		switch x := tok.(type) {
		case json.Delim:
			var out bytes.Buffer
			out.WriteByte(byte(x))
			for i := 0; dec.More(); i++ {
				if i > 0 {
					out.WriteByte(',')
				}
				var child any
				if x == '{' {
					kt, _ := dec.Token()
					k := kt.(string)
					kb, _ := json.Marshal(k)
					out.Write(kb)
					out.WriteByte(':')
					if m, ok := exp.(map[string]any); ok {
						child = m[k]
					}
				} else if l, ok := exp.([]any); ok && i < len(l) {
					child = l[i]
				}
				out.Write(walk(child))
			}
			dec.Token() // the closing delimiter
			if x == '{' {
				out.WriteByte('}')
			} else {
				out.WriteByte(']')
			}
			return out.Bytes()
		case string:
			if s, ok := exp.(string); ok && hasMarker(s) && (s == "$any" || matchString(s, x, b)) {
				kb, _ := marshalFixture(s)
				return kb
			}
			kb, _ := marshalFixture(x)
			return kb
		default:
			kb, _ := marshalFixture(tok)
			return kb
		}
	}
	return walk(exp)
}

// spliceMember replaces the value of a top-level member of a JSON object file
// with value, indented to sit at depth one, and keeps every other byte.
func spliceMember(raw []byte, key string, value any) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	if _, err := dec.Token(); err != nil {
		return nil, err
	}
	last := int64(0)
	indent := func(value any) ([]byte, error) {
		vb, err := marshalFixture(value)
		if err != nil {
			return nil, err
		}
		var ind bytes.Buffer
		if err := json.Indent(&ind, vb, "  ", "  "); err != nil {
			return nil, err
		}
		return ind.Bytes(), nil
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
		ind, err := indent(value)
		if err != nil {
			return nil, err
		}
		out := append([]byte(nil), raw[:start]...)
		out = append(out, ": "...)
		out = append(out, ind...)
		return append(out, raw[end:]...), nil
	}
	ind, err := indent(value) // a member the file does not have yet goes last
	if err != nil {
		return nil, err
	}
	kb, _ := json.Marshal(key)
	out := append([]byte(nil), raw[:last]...)
	out = append(out, ",\n  "...)
	out = append(out, kb...)
	out = append(out, ": "...)
	out = append(out, ind...)
	return append(out, raw[last:]...), nil
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

// noIfMatch reports that the case expects the refusal of a write that carries
// no If-Match, so the HTTP runner must not add one.
func noIfMatch(expected json.RawMessage) bool {
	var e struct {
		Error struct{ Details struct{ Reason string } }
	}
	json.Unmarshal(expected, &e)
	return e.Error.Details.Reason == "if_match_required"
}

// firstBinding is the binding -update runs a case with: api, unless the case
// is restricted to others.
func firstBinding(c *conformanceCase) string {
	if len(c.Bindings) > 0 {
		return c.Bindings[0]
	}
	return "api"
}
