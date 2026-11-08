package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"
)

// the AKM and file properties a view may name; anything else is a mistake
var basesKnown = map[string]bool{
	"status": true, "parent": true, "updated": true, "closed": true, "decisions": true, "blocked_by": true,
	"cycle": true, "priority": true, "kind": true, "keep": true, "claimed_by": true, "date": true, "scope": true,
	"file.name": true, "formula.age_days": true, "formula.stale": true,
}

func parseBases(t *testing.T, text string) map[string]any {
	t.Helper()
	var doc map[string]any
	if err := yaml.Unmarshal([]byte(text), &doc); err != nil {
		t.Fatalf("the Obsidian Bases file does not parse as YAML: %v\n%s", err, text)
	}
	return doc
}

func TestBasesFileParsesAndNamesOnlyAKMFields(t *testing.T) {
	doc := parseBases(t, BasesFile("docs/work", "docs/adr", 30))
	views, _ := doc["views"].([]any)
	if len(views) != 10 {
		t.Fatalf("got %d views, want 10", len(views))
	}
	kinds := map[string]int{}
	for _, v := range views {
		view := v.(map[string]any)
		kinds[view["type"].(string)]++
		for _, p := range view["order"].([]any) {
			if !basesKnown[p.(string)] {
				t.Errorf("%v orders by %q, which is not an AKM field", view["name"], p)
			}
		}
		for _, s := range asList(view["sort"]) {
			if p := s.(map[string]any)["property"].(string); !basesKnown[p] {
				t.Errorf("%v sorts by %q", view["name"], p)
			}
		}
		if g, ok := view["groupBy"].(map[string]any); ok && !basesKnown[g["property"].(string)] {
			t.Errorf("%v groups by %q", view["name"], g["property"])
		}
		if view["type"] == "kanban" {
			if got := asList(view["groupOrder"]); len(got) != 4 || got[0] != "raw" || got[3] != "blocked" {
				t.Errorf("the Kanban group order is %v, want the open statuses in order", got)
			}
		}
	}
	if kinds["table"] != 6 || kinds["cards"] != 3 || kinds["kanban"] != 1 {
		t.Errorf("view types %v", kinds)
	}
	for k := range doc["formulas"].(map[string]any) {
		if _, ok := doc["properties"].(map[string]any)["formula."+k]; !ok {
			t.Errorf("formula %s has no property entry", k)
		}
	}
}

func asList(v any) []any {
	l, _ := v.([]any)
	return l
}

func TestBasesFileUsesTheManifestsFoldersAndStaleDays(t *testing.T) {
	text := BasesFile("notes/it's/work", "notes/adr", 45)
	parseBases(t, text)
	for _, want := range []string{`'file.inFolder("notes/it''s/work")'`, `'file.inFolder("notes/adr")'`, ".days > 45,"} {
		if !strings.Contains(text, want) {
			t.Errorf("the file lacks %s", want)
		}
	}
	if strings.Contains(text, "docs/work") {
		t.Error("the default folder leaked into a file for other folders")
	}
}

func TestBasesWritesOnceAndNeverOverwrites(t *testing.T) {
	root := copyRepo(t, "basic", "")
	resp := mustOK(t, root, "bases", Request{})
	if p, _ := resp.Get("path"); p != "docs/akm.base" {
		t.Fatalf("path %v", p)
	}
	if w, _ := resp.Get("written"); w != true {
		t.Fatal("not written")
	}
	full := filepath.Join(root, "docs", "akm.base")
	first, err := os.ReadFile(full)
	if err != nil {
		t.Fatal(err)
	}
	parseBases(t, string(first))
	// the maintainer edits it; a second call refuses and leaves the edit alone
	os.WriteFile(full, []byte("mine\n"), 0o644)
	wantCode(t, root, "bases", Request{}, "refused", "exists")
	if b, _ := os.ReadFile(full); string(b) != "mine\n" {
		t.Errorf("the file was overwritten: %q", b)
	}
	wantCode(t, root, "bases", Request{"dry_run": true}, "refused", "exists")
}

func TestBasesDryRunAndContentWriteNothing(t *testing.T) {
	root := copyRepo(t, "basic", "")
	resp := mustOK(t, root, "bases", Request{"dry_run": true, "include_content": true})
	if w, _ := resp.Get("written"); w != false {
		t.Error("a dry run wrote")
	}
	c, _ := resp.Get("content")
	parseBases(t, c.(string))
	if _, err := os.Stat(filepath.Join(root, "docs", "akm.base")); err == nil {
		t.Error("the file exists after a dry run")
	}
}

func TestBasesFindsTheVaultRoot(t *testing.T) {
	// the repository is a folder inside a vault: folder filters start at the vault
	vault := t.TempDir()
	os.MkdirAll(filepath.Join(vault, ".obsidian"), 0o755)
	inner := filepath.Join(vault, "projects", "repo")
	src := copyRepo(t, "basic", "")
	if err := os.MkdirAll(filepath.Dir(inner), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(src, inner); err != nil {
		t.Skipf("cannot move the fixture: %v", err)
	}
	resp := mustOK(t, inner, "bases", Request{"include_content": true})
	if f, _ := resp.Get("vault_found"); f != true {
		t.Error("the vault was not found")
	}
	if v, _ := resp.Get("vault_root"); v != "../.." {
		t.Errorf("vault_root %v", v)
	}
	if w, _ := resp.Get("work_folder"); w != "projects/repo/docs/work" {
		t.Errorf("work_folder %v", w)
	}
	c, _ := resp.Get("content")
	if !strings.Contains(c.(string), `file.inFolder("projects/repo/docs/work")`) {
		t.Error("the filter does not start at the vault root")
	}
}

func TestBasesRefusesBadPaths(t *testing.T) {
	root := copyRepo(t, "basic", "")
	for _, p := range []string{"/abs/x.base", "../x.base", "docs/../x.base", "docs/x.md", ".base", `docs\x.base`, "docs//x.base"} {
		wantCode(t, root, "bases", Request{"path": p}, "invalid_request", "")
	}
	mustOK(t, root, "bases", Request{"path": "views/mine.base"})
	if _, err := os.Stat(filepath.Join(root, "views", "mine.base")); err != nil {
		t.Errorf("a path of the caller's was not written: %v", err)
	}
}
