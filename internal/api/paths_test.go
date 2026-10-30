package api

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/check"
	"github.com/alexsasharegan/akmctl/internal/record"
)

const (
	vWork = "Deep Thought/9 pkm/tasks"
	vADR  = "Deep Thought/9 pkm/adr"
)

func TestVersionReportsFolders(t *testing.T) {
	root := copyRepo(t, "vault", "")
	resp := mustOK(t, root, "version", Request{})
	repo, _ := resp.Get("repository")
	r := repo.(Obj)
	if w, _ := r.Get("work_dir"); w != vWork {
		t.Fatalf("work_dir %v", w)
	}
	if a, _ := r.Get("adr_dir"); a != vADR {
		t.Fatalf("adr_dir %v", a)
	}
}

func TestWritesFollowTheManifestFolders(t *testing.T) {
	root := copyRepo(t, "vault", "")
	clean(t, root)
	// new work and new adr land in their folders
	resp := mustOK(t, root, "new", Request{"kind": "work", "slug": "add-a-thing", "as_of": "2026-10-28"})
	p, _ := resp.Get("path")
	if !regexp.MustCompile(`^Deep Thought/9 pkm/tasks/261028-[0-9a-hjkmnp-tv-z]{4}-add-a-thing\.md$`).MatchString(p.(string)) {
		t.Fatalf("work path %v", p)
	}
	resp = mustOK(t, root, "new", Request{"kind": "adr", "slug": "use-sqlite", "as_of": "2026-10-28"})
	p, _ = resp.Get("path")
	if !strings.HasPrefix(p.(string), vADR+"/261028-") {
		t.Fatalf("adr path %v", p)
	}
	// the index goes beside the ADRs
	resp = mustOK(t, root, "index", Request{})
	if ip, _ := resp.Get("path"); ip != vADR+"/README.md" {
		t.Fatalf("index path %v", ip)
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "adr")); err == nil {
		t.Fatal("the default folder was created")
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "work")); err == nil {
		t.Fatal("the default work folder was created")
	}
	// set and section edit in place; close moves into done/ of the work folder
	mustOK(t, root, "set", Request{"id": a2b3, "changes": map[string]any{"priority": 5}})
	mustOK(t, root, "section", Request{"id": a2b3, "name": "Outcome", "text": "Shipped."})
	resp = mustOK(t, root, "close", Request{"id": a2b3, "status": "done", "as_of": "2026-10-28"})
	to, _ := resp.Get("moved_to")
	if to != vWork+"/done/261020-a2b3-add-json-output-to-ls.md" {
		t.Fatalf("moved to %v", to)
	}
	// a closed record's sections are not edited; a done one cannot be reopened outside done/
	wantCode(t, root, "section", Request{"id": a2b3, "name": "Log", "text": "x", "mode": "append"}, "refused", "closed")
	wantCode(t, root, "set", Request{"id": a2b3, "changes": map[string]any{"status": "ready"}}, "refused", "rule")
	repo, _ := record.Load(root)
	for _, f := range check.Run(repo, check.Options{}).Findings {
		if f.Severity == "E" {
			t.Errorf("finding %+v", f)
		}
	}
}

func TestClaimsFollowTheManifestFolders(t *testing.T) {
	root := copyRepo(t, "vault", "extension.claims: 0.3.0\n")
	resp := mustOK(t, root, "claim", Request{"id": e5f6, "as": "agent-1", "as_of": "2026-10-28T09:30"})
	rec, _ := resp.Get("record")
	if p, _ := rec.(Obj).Get("path"); p != vWork+"/261022-e5f6-document-the-exit-codes.md" {
		t.Fatalf("path %v", p)
	}
	mustOK(t, root, "release", Request{"id": e5f6})
}

func TestBadManifestPathRefusesWrites(t *testing.T) {
	root := copyRepo(t, "basic", "work_dir: ../outside\n")
	for _, c := range []struct {
		op  string
		req Request
	}{
		{"new", Request{"kind": "work", "slug": "x"}},
		{"set", Request{"id": e5f6, "changes": map[string]any{"priority": 1}}},
		{"close", Request{"id": a2b3, "status": "done"}},
		{"index", Request{}},
	} {
		wantCode(t, root, c.op, c.req, "refused", "rule")
	}
	if _, err := os.Stat(filepath.Join(root, "..", "outside")); err == nil {
		t.Fatal("wrote outside the repository")
	}
	// reads still work and check reports S3
	mustOK(t, root, "ls", Request{})
	resp := mustOK(t, root, "check", Request{"rules": []string{"S3"}})
	if n, _ := resp.Get("errors"); n != 1 {
		t.Fatalf("errors %v", n)
	}
}

func TestCloseCreatesOutcomeWithBlankLine(t *testing.T) {
	root := copyRepo(t, "basic", "")
	resp := mustOK(t, root, "close", Request{"id": a2b3, "status": "done", "outcome": "Shipped.", "as_of": "2026-10-28"})
	to, _ := resp.Get("moved_to")
	if text := read(t, root, to.(string)); !strings.Contains(text, "## Outcome\n\nShipped.\n") {
		t.Fatalf("no blank line after ## Outcome:\n%s", text)
	}
}
