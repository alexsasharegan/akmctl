package api

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/check"
	"github.com/alexsasharegan/akmctl/internal/record"
)

func copyRepo(t *testing.T, name string, extraLayout string) string {
	t.Helper()
	src := filepath.Join("../../testdata/api/repos", name)
	dst := t.TempDir()
	filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		rel, _ := filepath.Rel(src, p)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		b, _ := os.ReadFile(p)
		return os.WriteFile(filepath.Join(dst, rel), b, 0o644)
	})
	if extraLayout != "" {
		f, _ := os.OpenFile(filepath.Join(dst, "docs/.layout"), os.O_APPEND|os.O_WRONLY, 0o644)
		f.WriteString(extraLayout)
		f.Close()
	}
	return dst
}

func call(t *testing.T, root, op string, req Request) (Obj, *Error) {
	t.Helper()
	repo, err := record.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	return Run(repo, Info{"akmctl", "test"}, op, req)
}

func mustOK(t *testing.T, root, op string, req Request) Obj {
	t.Helper()
	resp, err := call(t, root, op, req)
	if err != nil {
		t.Fatalf("%s %v: %v", op, req, err)
	}
	return resp
}

func wantCode(t *testing.T, root, op string, req Request, code, reason string) {
	t.Helper()
	_, err := call(t, root, op, req)
	if err == nil || err.Code != code {
		t.Fatalf("%s %v: got %v, want %s", op, req, err, code)
	}
	if reason != "" {
		if r, _ := err.Details.Get("reason"); r != reason {
			t.Fatalf("%s: reason %v, want %s", op, r, reason)
		}
	}
}

func read(t *testing.T, root, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func clean(t *testing.T, root string) {
	t.Helper()
	repo, _ := record.Load(root)
	if res := check.Run(repo, check.Options{}); len(res.Findings) != 0 {
		t.Fatalf("check: %+v", res.Findings)
	}
}

const (
	e5f6 = "261022-e5f6"
	a2b3 = "261020-a2b3"
	c4d5 = "261015-c4d5"
)

func TestNewWork(t *testing.T) {
	root := copyRepo(t, "basic", "")
	resp := mustOK(t, root, "new", Request{"kind": "work", "slug": "add-a-thing", "title": "Add a thing", "as_of": "2026-10-28T09:30",
		"fields": map[string]any{"summary": "One line.", "parent": c4d5, "priority": 1, "kind": "feature"}})
	p, _ := resp.Get("path")
	if !regexp.MustCompile(`^docs/work/261028-[0-9a-hjkmnp-tv-z]{4}-add-a-thing\.md$`).MatchString(p.(string)) {
		t.Fatalf("path %v", p)
	}
	text := read(t, root, p.(string))
	for _, want := range []string{"title: Add a thing\n", "summary: One line.\n", "status: raw\n", "parent: " + c4d5 + "\n", "priority: 1\n", "created: 2026-10-28\n", "updated: 2026-10-28\n", "# Add a thing\n", "## Log\n\n- 2026-10-28: filed.\n"} {
		if !strings.Contains(text, want) {
			t.Errorf("missing %q in\n%s", want, text)
		}
	}
	clean(t, root)
}

func TestNewADRNeedsIndexRun(t *testing.T) {
	root := copyRepo(t, "basic", "")
	mustOK(t, root, "new", Request{"kind": "adr", "slug": "use-sqlite", "title": "Use SQLite", "as_of": "2026-10-28", "fields": map[string]any{"scope": []any{"src/store/"}}})
	repo, _ := record.Load(root)
	res := check.Run(repo, check.Options{})
	if len(res.Findings) != 1 || res.Findings[0].Rule != "A4" {
		t.Fatalf("%+v", res.Findings)
	}
	mustOK(t, root, "index", Request{})
	clean(t, root)
}

func TestNewRejects(t *testing.T) {
	root := copyRepo(t, "basic", "")
	wantCode(t, root, "new", Request{"kind": "work", "slug": "Bad Slug"}, "invalid_request", "")
	wantCode(t, root, "new", Request{"kind": "work", "slug": "x", "fields": map[string]any{"status": "done"}}, "invalid_request", "")
	wantCode(t, root, "new", Request{"kind": "work", "slug": "x", "fields": map[string]any{"parent": "261099-aaaa"}}, "not_found", "")
	wantCode(t, root, "new", Request{"kind": "work", "slug": "x", "fields": map[string]any{"priority": "high"}}, "invalid_request", "")
	wantCode(t, root, "new", Request{"kind": "work", "slug": "x", "title": "two\nlines"}, "invalid_request", "")
}

func TestNewIDsAreUnique(t *testing.T) {
	root := copyRepo(t, "basic", "")
	seen := map[string]bool{}
	for i := 0; i < 40; i++ {
		resp := mustOK(t, root, "new", Request{"kind": "work", "slug": "item", "as_of": "2026-10-28"})
		rec, _ := resp.Get("record")
		id, _ := rec.(Obj).Get("id")
		if seen[id.(string)] {
			t.Fatalf("duplicate id %v", id)
		}
		seen[id.(string)] = true
	}
}

func TestSetChangesOnlyWhatItMust(t *testing.T) {
	root := copyRepo(t, "basic", "")
	before := read(t, root, "docs/work/261022-e5f6-document-the-exit-codes.md")
	resp := mustOK(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"priority": 4, "summary": "Document them."}, "as_of": "2026-10-28"})
	ch, _ := resp.Get("changed")
	if strings.Join(ch.([]string), ",") != "priority,summary" {
		t.Fatalf("changed %v", ch)
	}
	after := read(t, root, "docs/work/261022-e5f6-document-the-exit-codes.md")
	want := strings.Replace(before, "updated: 2026-10-22\n", "updated: 2026-10-28\n", 1)
	want = strings.Replace(want, "parent: 261015-c4d5\n", "parent: 261015-c4d5\npriority: 4\nsummary: Document them.\n", 1)
	// new keys are appended at the end of the block, after updated
	want = strings.Replace(want, "parent: 261015-c4d5\npriority: 4\nsummary: Document them.\n", "parent: 261015-c4d5\n", 1)
	want = strings.Replace(want, "updated: 2026-10-28\n---", "updated: 2026-10-28\npriority: 4\nsummary: Document them.\n---", 1)
	if after != want {
		t.Fatalf("got\n%s\nwant\n%s", after, want)
	}
	// setting the same values again writes nothing
	mustOK(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"priority": 4}, "as_of": "2030-01-01"})
	if read(t, root, "docs/work/261022-e5f6-document-the-exit-codes.md") != after {
		t.Fatal("a no-op set changed the file")
	}
}

func TestSetRemovesAndRefuses(t *testing.T) {
	root := copyRepo(t, "basic", "")
	mustOK(t, root, "set", Request{"id": a2b3, "changes": map[string]any{"priority": nil}})
	if strings.Contains(read(t, root, "docs/work/261020-a2b3-add-json-output-to-ls.md"), "priority") {
		t.Fatal("priority not removed")
	}
	wantCode(t, root, "set", Request{"id": a2b3, "changes": map[string]any{"title": nil}}, "invalid_request", "")
	wantCode(t, root, "set", Request{"id": a2b3, "changes": map[string]any{"id": "261099-aaaa"}}, "invalid_request", "")
	wantCode(t, root, "set", Request{"id": a2b3, "changes": map[string]any{"created": "2020-01-01"}}, "invalid_request", "")
	wantCode(t, root, "set", Request{"id": a2b3, "changes": map[string]any{"status": "done"}}, "invalid_request", "")
	wantCode(t, root, "set", Request{"id": a2b3, "changes": map[string]any{"status": "nope"}}, "invalid_request", "")
	wantCode(t, root, "set", Request{"id": a2b3, "changes": map[string]any{"claimed_by": "x"}}, "invalid_request", "")
	wantCode(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"status": "blocked"}}, "refused", "rule")
	wantCode(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"parent": e5f6}}, "refused", "rule")
	wantCode(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"parent": "261099-aaaa"}}, "not_found", "")
	wantCode(t, root, "set", Request{"id": "261099-zzzz", "changes": map[string]any{"priority": 1}}, "not_found", "")
	wantCode(t, root, "set", Request{"id": "261010-j9k0", "changes": map[string]any{"status": "ready"}}, "refused", "rule")
	mustOK(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"status": "blocked", "blocked_by": []any{"a person"}}})
}

func TestSetADRStatusDate(t *testing.T) {
	root := copyRepo(t, "basic", "")
	mustOK(t, root, "set", Request{"id": "261018-r5s6", "changes": map[string]any{"status": "accepted"}, "as_of": "2026-10-29"})
	if !strings.Contains(read(t, root, "docs/adr/261018-r5s6-use-a-local-server-for-the-ui.md"), "date: 2026-10-29\n") {
		t.Fatal("date not refreshed")
	}
	wantCode(t, root, "set", Request{"id": "261018-r5s6", "changes": map[string]any{"status": "superseded"}}, "refused", "rule")
}

func TestStaleEtag(t *testing.T) {
	root := copyRepo(t, "basic", "")
	repo, _ := record.Load(root)
	etag := repo.ByID(e5f6)[0].Etag
	mustOK(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"priority": 1}, "if_match": etag})
	wantCode(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"priority": 2}, "if_match": etag}, "stale", "")
	if !strings.Contains(read(t, root, "docs/work/261022-e5f6-document-the-exit-codes.md"), "priority: 1\n") {
		t.Fatal("stale write went through")
	}
}

// Two writers with the same etag: exactly one wins.
func TestCompareAndSwapRace(t *testing.T) {
	root := copyRepo(t, "basic", "")
	repo, _ := record.Load(root)
	etag := repo.ByID(e5f6)[0].Etag
	var wg sync.WaitGroup
	results := make([]*Error, 8)
	for i := range results {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, _ := record.Load(root)
			_, results[i] = Run(r, Info{}, "set", Request{"id": e5f6, "changes": map[string]any{"priority": i + 10}, "if_match": etag})
		}()
	}
	wg.Wait()
	won := 0
	for _, e := range results {
		if e == nil {
			won++
		} else if e.Code != "stale" {
			t.Fatalf("unexpected error %v", e)
		}
	}
	if won != 1 {
		t.Fatalf("%d writers won", won)
	}
}

func TestClose(t *testing.T) {
	root := copyRepo(t, "basic", "")
	wantCode(t, root, "close", Request{"id": a2b3, "status": "done"}, "refused", "empty_outcome")
	// give the item an outcome
	p := filepath.Join(root, "docs/work/261020-a2b3-add-json-output-to-ls.md")
	b, _ := os.ReadFile(p)
	os.WriteFile(p, append(b, []byte("\n## Outcome\nShipped.\n")...), 0o644)
	resp := mustOK(t, root, "close", Request{"id": a2b3, "status": "done", "as_of": "2026-10-28"})
	to, _ := resp.Get("moved_to")
	if to != "docs/work/done/261020-a2b3-add-json-output-to-ls.md" {
		t.Fatalf("moved to %v", to)
	}
	if _, err := os.Stat(p); err == nil {
		t.Fatal("old file still there")
	}
	text := read(t, root, to.(string))
	if !strings.Contains(text, "status: done\n") || !strings.Contains(text, "updated: 2026-10-28\n") {
		t.Fatal(text)
	}
	wantCode(t, root, "close", Request{"id": a2b3, "status": "dropped"}, "refused", "already_closed")
}

func TestCloseParentWithOpenChildren(t *testing.T) {
	root := copyRepo(t, "basic", "")
	p := filepath.Join(root, "docs/work/261015-c4d5-machine-readable-listings.md")
	b, _ := os.ReadFile(p)
	os.WriteFile(p, append(b, []byte("\n## Outcome\nDone.\n")...), 0o644)
	wantCode(t, root, "close", Request{"id": c4d5, "status": "done"}, "refused", "open_children")
	mustOK(t, root, "close", Request{"id": c4d5, "status": "dropped"})
}

func TestClaimAndRelease(t *testing.T) {
	root := copyRepo(t, "basic", "")
	wantCode(t, root, "claim", Request{"id": e5f6, "as": "agent-1"}, "refused", "extension_off")
	root = copyRepo(t, "basic", "extension.claims: 0.3.0\n")
	resp := mustOK(t, root, "claim", Request{"id": e5f6, "as": "agent-1", "as_of": "2026-10-28T09:30"})
	if pv, _ := resp.Get("previous_holder"); pv != nil {
		t.Fatalf("previous %v", pv)
	}
	text := read(t, root, "docs/work/261022-e5f6-document-the-exit-codes.md")
	for _, want := range []string{"status: active\n", "claimed_by: agent-1\n", "claimed_at: 2026-10-28T09:30\n", "## Log\n- 2026-10-28: claimed (agent-1).\n"} {
		if !strings.Contains(text, want) {
			t.Errorf("missing %q in\n%s", want, text)
		}
	}
	// another holder, not stale
	wantCode(t, root, "claim", Request{"id": e5f6, "as": "agent-2", "as_of": "2026-10-28T12:00"}, "refused", "claimed")
	// the same holder refreshes
	mustOK(t, root, "claim", Request{"id": e5f6, "as": "agent-1", "as_of": "2026-10-28T12:00"})
	text = read(t, root, "docs/work/261022-e5f6-document-the-exit-codes.md")
	if n := strings.Count(text, "claimed (agent-1)."); n != 1 || !strings.Contains(text, "claimed_at: 2026-10-28T12:00\n") {
		t.Fatalf("same-holder re-claim must refresh claimed_at without a second log line:\n%s", text)
	}
	// stale: takeover without the flag
	resp = mustOK(t, root, "claim", Request{"id": e5f6, "as": "agent-2", "as_of": "2026-10-30T12:00"})
	if pv, _ := resp.Get("previous_holder"); pv != "agent-1" {
		t.Fatalf("previous %v", pv)
	}
	if !strings.Contains(read(t, root, "docs/work/261022-e5f6-document-the-exit-codes.md"), "took over from agent-1 (agent-2).") {
		t.Fatal("no takeover line")
	}
	// forced takeover of a fresh claim
	mustOK(t, root, "claim", Request{"id": e5f6, "as": "agent-3", "takeover": true, "as_of": "2026-10-30T12:05"})
	mustOK(t, root, "release", Request{"id": e5f6, "as_of": "2026-10-30"})
	text = read(t, root, "docs/work/261022-e5f6-document-the-exit-codes.md")
	if strings.Contains(text, "claimed_by") || !strings.Contains(text, "released (agent-3).") || !strings.Contains(text, "status: active\n") {
		t.Fatal(text)
	}
	wantCode(t, root, "release", Request{"id": e5f6}, "refused", "not_claimed")
	wantCode(t, root, "claim", Request{"id": "261026-m1n2", "as": "x"}, "refused", "status")
}

// A claimed item that is blocked and then unblocked to active (set, then
// claim by the same identity) keeps one 'claimed' line in its log.
func TestUnblockClaimedItemLogsOneClaim(t *testing.T) {
	root := copyRepo(t, "basic", "extension.claims: 0.3.0\n")
	mustOK(t, root, "claim", Request{"id": e5f6, "as": "agent-1", "as_of": "2026-10-28T09:30"})
	mustOK(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"status": "blocked", "blocked_by": []any{"waiting on review"}}, "as_of": "2026-10-28"})
	mustOK(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"status": "active", "blocked_by": nil}, "as_of": "2026-10-29"})
	mustOK(t, root, "claim", Request{"id": e5f6, "as": "agent-1", "as_of": "2026-10-29T08:00"})
	text := read(t, root, "docs/work/261022-e5f6-document-the-exit-codes.md")
	if n := strings.Count(text, "claimed (agent-1)."); n != 1 {
		t.Fatalf("want one claimed line, got %d:\n%s", n, text)
	}
	// another identity still has to ask for a takeover
	wantCode(t, root, "claim", Request{"id": e5f6, "as": "agent-2", "as_of": "2026-10-29T08:05"}, "refused", "claimed")
}

func TestCloseClearsClaim(t *testing.T) {
	root := copyRepo(t, "basic", "extension.claims: 0.3.0\n")
	mustOK(t, root, "claim", Request{"id": a2b3, "as": "agent-1", "as_of": "2026-10-28T09:30"})
	p := filepath.Join(root, "docs/work/261020-a2b3-add-json-output-to-ls.md")
	b, _ := os.ReadFile(p)
	os.WriteFile(p, append(b, []byte("\n## Outcome\nShipped.\n")...), 0o644)
	resp := mustOK(t, root, "close", Request{"id": a2b3, "status": "done"})
	to, _ := resp.Get("moved_to")
	if strings.Contains(read(t, root, to.(string)), "claimed_") {
		t.Fatal("claim not cleared")
	}
}

func TestWikilinkReferenceStyle(t *testing.T) {
	root := copyRepo(t, "basic", "reference_style: wikilink\n")
	mustOK(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"parent": a2b3}})
	text := read(t, root, "docs/work/261022-e5f6-document-the-exit-codes.md")
	if !strings.Contains(text, `parent: "[[261020-a2b3-add-json-output-to-ls|261020-a2b3]]"`) {
		t.Fatal(text)
	}
	resp := mustOK(t, root, "show", Request{"id": e5f6})
	rec, _ := resp.Get("record")
	f, _ := rec.(Obj).Get("fields")
	if p, _ := f.(Obj).Get("parent"); p != a2b3 {
		t.Fatalf("parent %v", p)
	}
	clean(t, root)
}

func TestNoTempFilesLeft(t *testing.T) {
	root := copyRepo(t, "basic", "")
	mustOK(t, root, "set", Request{"id": e5f6, "changes": map[string]any{"priority": 9}})
	ents, _ := os.ReadDir(filepath.Join(root, "docs/work"))
	for _, e := range ents {
		if strings.HasSuffix(e.Name(), ".tmp") {
			t.Fatalf("left %s", e.Name())
		}
	}
}

func loadRepo(t *testing.T, root string) (*record.Repo, error) {
	t.Helper()
	return record.Load(root)
}
