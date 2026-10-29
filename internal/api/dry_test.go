package api

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alexsasharegan/akmctl/internal/record"
)

const e5f6File = "docs/work/261022-e5f6-document-the-exit-codes.md"

// state records every file under root with its mtime, size and text.
type fileState struct {
	mtime time.Time
	text  string
}

func tree(t *testing.T, root string) map[string]fileState {
	t.Helper()
	m := map[string]fileState{}
	filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			st, _ := d.Info()
			b, _ := os.ReadFile(p)
			rel, _ := filepath.Rel(root, p)
			if filepath.Base(rel) == ".akm.lock" {
				return nil
			}
			m[filepath.ToSlash(rel)] = fileState{st.ModTime(), string(b)}
		}
		return nil
	})
	return m
}

func sameTree(a, b map[string]fileState) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if w, ok := b[k]; !ok || !v.mtime.Equal(w.mtime) || v.text != w.text {
			return false
		}
	}
	return true
}

func putOutcome(t *testing.T, root string) {
	t.Helper()
	p := filepath.Join(root, "docs/work/261020-a2b3-add-json-output-to-ls.md")
	b, _ := os.ReadFile(p)
	os.WriteFile(p, append(b, []byte("\n## Outcome\nShipped.\n")...), 0o644)
}

// A dry run changes nothing on disk: not a file, not a modification time, and
// no lock file either.
func TestDryRunNeverTouchesTheRepository(t *testing.T) {
	root := copyRepo(t, "claims", "")
	putOutcome(t, root)
	// age every file, so that a rewrite would show in the mtime
	old := time.Now().Add(-48 * time.Hour)
	filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			os.Chtimes(p, old, old)
		}
		return nil
	})
	before := tree(t, root)
	for _, c := range []struct {
		op  string
		req Request
	}{
		{"set", Request{"id": e5f6, "changes": map[string]any{"priority": 9}}},
		{"set", Request{"id": e5f6, "changes": map[string]any{"parent": nil}}},
		{"section", Request{"id": e5f6, "name": "Goal", "text": "x"}},
		{"section", Request{"id": a2b3, "name": "Log", "text": "- x", "mode": "append"}},
		{"close", Request{"id": a2b3, "status": "done"}},
		{"close", Request{"id": e5f6, "status": "dropped", "outcome": "Not needed."}},
		{"claim", Request{"id": e5f6, "as": "agent-1"}},
		{"release", Request{"id": a2b3}},
	} {
		c.req["dry_run"] = true
		resp := mustOK(t, root, c.op, c.req)
		d, ok := resp.Get("diff")
		if !ok {
			t.Fatalf("%s: no diff in a dry run", c.op)
		}
		if hs, _ := d.(Obj).Get("hunks"); len(hs.([]Obj)) == 0 {
			t.Fatalf("%s %v: the diff is empty", c.op, c.req)
		}
		if !sameTree(before, tree(t, root)) {
			t.Fatalf("%s %v changed the repository", c.op, c.req)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "docs/.akm.lock")); err == nil {
		t.Fatal("a dry run created the lock file")
	}
}

// A write without dry_run has no diff member.
func TestNoDiffWithoutDryRun(t *testing.T) {
	root := copyRepo(t, "basic", "")
	resp := mustOK(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"priority": 9}})
	if _, ok := resp.Get("diff"); ok {
		t.Fatal("diff on a real write")
	}
}

// The record of a dry run is the would-be state with the current etag, and the
// next real write with the same request produces exactly that state.
func TestDryRunPreviewMatchesTheWrite(t *testing.T) {
	root := copyRepo(t, "basic", "")
	putOutcome(t, root)
	req := Request{"id": a2b3, "status": "done", "outcome": "Shipped in 0.2.", "as_of": "2026-10-28"}
	repo, _ := record.Load(root)
	cur := repo.ByID(a2b3)[0]

	dreq := Request{"dry_run": true}
	for k, v := range req {
		dreq[k] = v
	}
	prev := mustOK(t, root, "close", dreq)
	prec, _ := prev.Get("record")
	pr := prec.(Obj)
	if et, _ := pr.Get("etag"); et != cur.Etag {
		t.Fatalf("dry-run etag %v, want the current %s", et, cur.Etag)
	}
	if p, _ := pr.Get("path"); p != "docs/work/done/261020-a2b3-add-json-output-to-ls.md" {
		t.Fatalf("dry-run path %v", p)
	}
	d, _ := prev.Get("diff")
	if np, _ := d.(Obj).Get("new_path"); np != "docs/work/done/261020-a2b3-add-json-output-to-ls.md" {
		t.Fatalf("new_path %v", np)
	}

	real := mustOK(t, root, "close", req)
	rrec, _ := real.Get("record")
	rr := rrec.(Obj)
	pf, _ := pr.Get("fields")
	rf, _ := rr.Get("fields")
	if !jsonEqual(pf, rf) {
		t.Fatalf("preview fields differ from the written record:\n%v\n%v", pf, rf)
	}
}

func TestDryRunStillChecksTheEtagAndRules(t *testing.T) {
	root := copyRepo(t, "basic", "")
	wantCode(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"priority": 1}, "if_match": "0000000000000000", "dry_run": true}, "stale", "")
	wantCode(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"status": "blocked"}, "dry_run": true}, "refused", "rule")
	wantCode(t, root, "close", Request{"id": a2b3, "status": "done", "dry_run": true}, "refused", "empty_outcome")
	wantCode(t, root, "claim", Request{"id": e5f6, "as": "x", "dry_run": true}, "refused", "extension_off")
	// an unchanged set has an empty diff, not a missing one
	resp := mustOK(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"parent": c4d5}, "dry_run": true})
	d, _ := resp.Get("diff")
	if hs, _ := d.(Obj).Get("hunks"); len(hs.([]Obj)) != 0 {
		t.Fatal("expected no hunks")
	}
	if ch, _ := resp.Get("changed"); len(ch.([]string)) != 0 {
		t.Fatal("expected no changes")
	}
}

func TestSection(t *testing.T) {
	root := copyRepo(t, "basic", "")
	// replace, with updated set
	mustOK(t, root, "section", Request{"id": e5f6, "name": "goal", "text": "New goal.", "as_of": "2026-10-28"})
	text := read(t, root, e5f6File)
	if !strings.Contains(text, "## Goal\nNew goal.\n") || !strings.Contains(text, "updated: 2026-10-28\n") {
		t.Fatal(text)
	}
	// a replace that changes nothing writes nothing and leaves updated alone
	resp := mustOK(t, root, "section", Request{"id": e5f6, "name": "goal", "text": "New goal.", "as_of": "2026-11-02"})
	if ch, _ := resp.Get("changed"); ch != false {
		t.Fatal("changed should be false")
	}
	if read(t, root, e5f6File) != text {
		t.Fatal("a no-op replace rewrote the file")
	}
	// missing section is created before Log
	mustOK(t, root, "section", Request{"id": a2b3, "name": "Notes", "text": "n"})
	a := read(t, root, "docs/work/261020-a2b3-add-json-output-to-ls.md")
	if !strings.Contains(a, "## Notes\n\nn\n\n## Log\n") {
		t.Fatal(a)
	}
	// and at the end when there is no Log
	mustOK(t, root, "section", Request{"id": e5f6, "name": "Notes", "text": "n"})
	if !strings.HasSuffix(read(t, root, e5f6File), "\n\n## Notes\n\nn\n") {
		t.Fatal(read(t, root, e5f6File))
	}
	// the response carries the section as show gives it
	resp = mustOK(t, root, "section", Request{"id": a2b3, "name": "Log", "text": "- two", "mode": "append"})
	s, _ := resp.Get("section")
	if n, _ := s.(Obj).Get("name"); n != "log" {
		t.Fatal(n)
	}
	if tx, _ := s.(Obj).Get("text"); tx != "- 2026-10-20: started.\n- two" {
		t.Fatalf("%q", tx)
	}
	// errors
	wantCode(t, root, "section", Request{"id": e5f6, "name": "  ", "text": "x"}, "invalid_request", "")
	wantCode(t, root, "section", Request{"id": e5f6, "name": "## Goal", "text": "x"}, "invalid_request", "")
	wantCode(t, root, "section", Request{"id": e5f6, "name": "LOG", "text": "x"}, "invalid_request", "")
	wantCode(t, root, "section", Request{"id": e5f6, "name": "Goal", "text": "x", "mode": "prepend"}, "invalid_request", "")
	wantCode(t, root, "section", Request{"id": "261099-zzzz", "name": "Goal", "text": "x"}, "not_found", "")
	wantCode(t, root, "section", Request{"id": e5f6, "name": "Goal", "text": "x", "if_match": "0000000000000000"}, "stale", "")
	wantCode(t, root, "section", Request{"id": "261010-j9k0", "name": "Goal", "text": "x"}, "refused", "closed")
	// an empty text is a value: it clears the section
	mustOK(t, root, "section", Request{"id": e5f6, "name": "Goal", "text": ""})
	if !strings.Contains(read(t, root, e5f6File), "## Goal\n\n## Notes\n") {
		t.Fatal(read(t, root, e5f6File))
	}
	// the other sections and the frontmatter are untouched; check still passes
	clean(t, root)
}

func TestSectionKeepsCRLF(t *testing.T) {
	root := copyRepo(t, "basic", "")
	p := filepath.Join(root, e5f6File)
	b, _ := os.ReadFile(p)
	os.WriteFile(p, []byte(strings.ReplaceAll(string(b), "\n", "\r\n")), 0o644)
	mustOK(t, root, "section", Request{"id": e5f6, "name": "Goal", "text": "a\nb", "as_of": "2026-10-28"})
	text := read(t, root, e5f6File)
	if strings.Contains(strings.ReplaceAll(text, "\r\n", ""), "\n") || !strings.Contains(text, "## Goal\r\na\r\nb\r\n") {
		t.Fatalf("%q", text)
	}
}

func TestCloseWithOutcome(t *testing.T) {
	root := copyRepo(t, "basic", "")
	resp := mustOK(t, root, "close", Request{"id": a2b3, "status": "done", "outcome": "Shipped in 0.2.", "as_of": "2026-10-28"})
	to, _ := resp.Get("moved_to")
	text := read(t, root, to.(string))
	if !strings.Contains(text, "## Outcome\n\nShipped in 0.2.\n\n## Log\n") || !strings.Contains(text, "status: done\n") {
		t.Fatal(text)
	}
	if _, err := os.Stat(filepath.Join(root, "docs/work/261020-a2b3-add-json-output-to-ls.md")); err == nil {
		t.Fatal("old file still there")
	}
	// a whitespace outcome does not count; the file's own Outcome is used when the parameter is absent
	root = copyRepo(t, "basic", "")
	wantCode(t, root, "close", Request{"id": a2b3, "status": "done", "outcome": "  \n"}, "refused", "empty_outcome")
	putOutcome(t, root)
	mustOK(t, root, "close", Request{"id": a2b3, "status": "done"})
	// the parameter replaces an existing Outcome
	root = copyRepo(t, "closable", "")
	resp = mustOK(t, root, "close", Request{"id": a2b3, "status": "dropped", "outcome": "Replaced."})
	to, _ = resp.Get("moved_to")
	if text := read(t, root, to.(string)); !strings.Contains(text, "## Outcome\nReplaced.\n") || strings.Contains(text, "Shipped") {
		t.Fatal(text)
	}
	// a refusal after the outcome was prepared writes nothing
	root = copyRepo(t, "basic", "")
	before := tree(t, root)
	wantCode(t, root, "close", Request{"id": c4d5, "status": "done", "outcome": "x"}, "refused", "open_children")
	if !sameTree(before, tree(t, root)) {
		t.Fatal("a refused close wrote something")
	}
}

func TestVersionRepository(t *testing.T) {
	root := copyRepo(t, "claims", "claim_stale_hours: 6\nreference_style: wikilink\n")
	resp := mustOK(t, root, "version", Request{})
	r, _ := resp.Get("repository")
	ro, ok := r.(Obj)
	if !ok {
		t.Fatalf("repository: %#v", r)
	}
	b, _ := Marshal(ro)
	want := `{"name":"` + filepath.Base(root) + `","layout":"1.0.0","extensions":["adr-index","claims","work-items"],"reference_style":"wikilink","claim_stale_hours":6,"work_dir":"docs/work","adr_dir":"docs/adr"}`
	if string(b) != want {
		t.Fatalf("got  %s\nwant %s", b, want)
	}
	// no repository (the CLI's version outside one): null
	resp, _ = Run(&record.Repo{}, Info{"akmctl", "t"}, "version", Request{})
	if r, ok := resp.Get("repository"); !ok || r != nil {
		t.Fatalf("repository should be null, got %#v", r)
	}
}

// Writers and dry-run readers at the same time: every dry run sees a whole file
// (never a partial one), a dry run never fails with unreadable, and the lock
// file is not needed to read. Run with -race.
func TestWritesAndDryRunsRace(t *testing.T) {
	root := copyRepo(t, "basic", "")
	var wg sync.WaitGroup
	errs := make(chan string, 200)
	for w := 0; w < 4; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 15; i++ {
				r, _ := record.Load(root)
				_, err := Run(r, Info{}, "set", Request{"id": e5f6, "changes": map[string]any{"priority": w*100 + i}})
				if err != nil {
					errs <- "set: " + err.Error()
				}
				_, err = Run(r, Info{}, "section", Request{"id": e5f6, "name": "Log", "text": "- w", "mode": "append"})
				if err != nil {
					errs <- "section: " + err.Error()
				}
			}
		}()
	}
	for d := 0; d < 4; d++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 30; i++ {
				r, _ := record.Load(root)
				resp, err := Run(r, Info{}, "set", Request{"id": e5f6, "changes": map[string]any{"priority": 7777}, "dry_run": true})
				if err != nil {
					errs <- "dry set: " + err.Error()
					continue
				}
				if _, ok := resp.Get("diff"); !ok {
					errs <- "dry set: no diff"
				}
				_, err = Run(r, Info{}, "section", Request{"id": e5f6, "name": "Goal", "text": "dry", "dry_run": true})
				if err != nil {
					errs <- "dry section: " + err.Error()
				}
			}
		}()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		t.Error(e)
	}
	text := read(t, root, e5f6File)
	if strings.Contains(text, "7777") || strings.Contains(text, "dry") {
		t.Fatal("a dry run reached the file")
	}
	if n := strings.Count(text, "- w\n"); n != 60 {
		t.Fatalf("%d log lines, want 60", n)
	}
	clean(t, root)
	ents, _ := os.ReadDir(filepath.Join(root, "docs/work"))
	for _, e := range ents {
		if strings.Contains(e.Name(), ".tmp") {
			t.Fatalf("temp file left: %s", e.Name())
		}
	}
}
