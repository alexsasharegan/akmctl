package api

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	oTasks = "Notes Vault/9 pkm/tasks"
	oADR   = "Notes Vault/9 pkm/adr"
)

// legacyRepo is a notes repository with no layout history: a manifest that
// names the folders but has no layout key, and notes without ids.
func legacyRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	basic, _ := os.ReadFile("../../testdata/api/repos/basic/AGENTS.md")
	files := map[string]string{
		"AGENTS.md":                              string(basic),
		"docs/.layout":                           "work_dir: " + oTasks + "\nadr_dir: " + oADR + "\n",
		oTasks + "/2026-10-20 Write the docs.md": "---\ntype: task\nstatus: ready\ncreated_at: 2026-10-20\nupdated_at: 2026-10-21\n---\n# Write the docs\n\nWhy it matters.\n\n## Log\n\n- 2026-10-20: filed.\n",
		oTasks + "/2026-10-22 Second.md":         "---\ntype: task\nstatus: raw\ncreated_at: 2026-10-22\nupdated_at: 2026-10-22\n---\n# Second\n",
		oADR + "/ADR-001 Use SQLite.md":          "---\nid: 1\ntype: adr\nstatus: accepted\ndate: 2026-10-01\n---\n# Use SQLite\n",
	}
	for p, text := range files {
		full := filepath.Join(root, filepath.FromSlash(p))
		os.MkdirAll(filepath.Dir(full), 0o755)
		if err := os.WriteFile(full, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func oneOff() map[string]any {
	steps := []any{
		map[string]any{"id": "created", "kind": "rename_field", "class": "auto", "description": "created_at becomes created", "dir": oTasks, "from": "created_at", "to": "created"},
		map[string]any{"id": "updated", "kind": "rename_field", "class": "auto", "description": "updated_at becomes updated", "dir": oTasks, "from": "updated_at", "to": "updated"},
		map[string]any{"id": "title", "kind": "title_from_h1", "class": "auto", "description": "title from the heading", "dir": oTasks},
		map[string]any{"id": "goal", "kind": "goal_scaffold", "class": "auto", "description": "add a Goal", "dir": oTasks, "statuses": []any{"ready", "active"}},
		map[string]any{"id": "adr-number", "kind": "rename_field", "class": "auto", "description": "legacy id", "dir": oADR, "from": "id", "to": "adr_number"},
		map[string]any{"id": "ids-tasks", "kind": "mint_ids", "class": "auto", "description": "ids", "dir": oTasks, "rename": false, "seed": "t"},
		map[string]any{"id": "ids-adrs", "kind": "mint_ids", "class": "auto", "description": "ids", "dir": oADR, "date_field": "date", "rename": false, "seed": "a"},
	}
	return map[string]any{"from": "0.0.0", "to": "0.2.0", "steps": steps}
}

func stepsByID(resp Obj) map[string]Obj {
	st, _ := resp.Get("steps")
	m := map[string]Obj{}
	for _, s := range st.([]Obj) {
		id, _ := s.Get("id")
		m[id.(string)] = s
	}
	return m
}

func renamesOf(s Obj) []string {
	r, ok := s.Get("renames")
	if !ok {
		return nil
	}
	var out []string
	for _, x := range r.([]Obj) {
		f, _ := x.Get("from")
		to, _ := x.Get("to")
		out = append(out, f.(string)+" -> "+to.(string))
	}
	return out
}

func TestOneOffMigrationPlan(t *testing.T) {
	root := legacyRepo(t)
	// without the parameter there is no path
	resp := mustOK(t, root, "migrate", Request{})
	if st, _ := resp.Get("state"); st != "unsupported" {
		t.Fatalf("state %v", st)
	}
	before := tree(t, root)
	req := Request{"mode": "plan", "migrations": []any{oneOff()}, "as_of": "2026-10-28"}
	resp = mustOK(t, root, "migrate", req)
	if st, _ := resp.Get("state"); st != "behind" {
		t.Fatalf("state %v", st)
	}
	if f, _ := resp.Get("from"); f != "0.0.0" {
		t.Fatalf("from %v", f)
	}
	if !sameTree(before, tree(t, root)) {
		t.Fatal("plan wrote")
	}
	steps := stepsByID(resp)
	rn := renamesOf(steps["ids-tasks"])
	if len(rn) != 2 || !strings.HasPrefix(rn[0], oTasks+"/2026-10-20 Write the docs.md -> "+oTasks+"/261020-") || !strings.HasSuffix(rn[0], "-write-the-docs.md") {
		t.Fatalf("renames %v", rn)
	}
	// the ADR's numeric id was renamed first, so it gets a minted id from its date
	ra := renamesOf(steps["ids-adrs"])
	if len(ra) != 1 || !strings.HasPrefix(ra[0], oADR+"/ADR-001 Use SQLite.md -> "+oADR+"/261001-") {
		t.Fatalf("adr renames %v", ra)
	}
	// the plan is the same on every call, and the apply agrees with it
	again := stepsByID(mustOK(t, root, "migrate", req))
	if strings.Join(renamesOf(again["ids-tasks"]), "|") != strings.Join(rn, "|") {
		t.Fatal("the plan changed between calls")
	}
	// the goal step reports a note per file
	if n, _ := steps["goal"].Get("notes"); len(n.([]string)) != 1 {
		t.Fatalf("goal notes %v", n)
	}
	// the plan classifies each step; a step with a plan lists its files
	if f, _ := steps["created"].Get("files"); len(f.([]string)) != 2 {
		t.Fatalf("created files %v", f)
	}
}

func TestOneOffMigrationApply(t *testing.T) {
	root := legacyRepo(t)
	gitRepo(t, root)
	req := Request{"mode": "apply_safe", "migrations": []any{oneOff()}, "as_of": "2026-10-28T10:00"}
	wantCode(t, root, "migrate", req, "refused", "no_backup_ack")
	plan := stepsByID(mustOK(t, root, "migrate", Request{"mode": "plan", "migrations": []any{oneOff()}, "as_of": "2026-10-28"}))
	req["acknowledge_backup"] = true
	resp := mustOK(t, root, "migrate", req)
	if st, _ := resp.Get("state"); st != "current" {
		t.Fatalf("state %v", st)
	}
	applied := stepsByID(resp)
	for _, id := range []string{"ids-tasks", "ids-adrs"} {
		if strings.Join(renamesOf(applied[id]), "|") != strings.Join(renamesOf(plan[id]), "|") {
			t.Errorf("%s: apply %v, plan %v", id, renamesOf(applied[id]), renamesOf(plan[id]))
		}
	}
	// ids written, files not renamed, text converted
	text := read(t, root, oTasks+"/2026-10-20 Write the docs.md")
	for _, want := range []string{"id: \"261020-", "title: Write the docs\n", "created: 2026-10-20\n", "updated: 2026-10-21\n", "# Write the docs\n\n## Goal\n\nWhy it matters.\n\n## Log\n"} {
		if !strings.Contains(text, want) {
			t.Errorf("missing %q in\n%s", want, text)
		}
	}
	if !strings.Contains(read(t, root, oADR+"/ADR-001 Use SQLite.md"), "adr_number: 1\n") {
		t.Error("adr_number missing")
	}
	// manifest keeps the folders and gains the layout; the log has one line per applied step
	man := read(t, root, "docs/.layout")
	if !strings.Contains(man, "layout: 0.2.0") || !strings.Contains(man, "work_dir: "+oTasks) {
		t.Fatalf("manifest %q", man)
	}
	if n := strings.Count(read(t, root, "docs/.migrate-log"), "\n"); n != 7 {
		t.Fatalf("log lines %d", n)
	}
	// the tag and branch exist, and nothing was committed
	tags, _ := exec.Command("git", "-C", root, "tag", "--list").Output()
	if strings.TrimSpace(string(tags)) != "pre-migrate-0.0.0-to-0.2.0" {
		t.Fatalf("tags %q", tags)
	}
	if out, _ := exec.Command("git", "-C", root, "status", "--porcelain").Output(); len(out) == 0 {
		t.Fatal("expected uncommitted changes")
	}
	if n, _ := exec.Command("git", "-C", root, "rev-list", "--count", "HEAD").Output(); strings.TrimSpace(string(n)) != "1" {
		t.Fatalf("commits %s", n)
	}
	// the migration is over: the repository is current and a new call does nothing
	resp = mustOK(t, root, "migrate", Request{"mode": "apply_safe", "migrations": []any{oneOff()}, "acknowledge_backup": true})
	if st, _ := resp.Get("state"); st != "current" {
		t.Fatalf("state %v", st)
	}
}

func TestOneOffMigrationKeepsTheSafeguards(t *testing.T) {
	root := legacyRepo(t)
	req := Request{"mode": "apply_safe", "migrations": []any{oneOff()}, "acknowledge_backup": true}
	wantCode(t, root, "migrate", req, "refused", "not_git")
	gitRepo(t, root)
	os.WriteFile(filepath.Join(root, "dirty.txt"), []byte("x"), 0o644)
	wantCode(t, root, "migrate", req, "refused", "dirty_tree")
	os.Remove(filepath.Join(root, "dirty.txt"))
	exec.Command("git", "-C", root, "branch", "migrate-0.0.0-to-0.2.0").Run()
	wantCode(t, root, "migrate", req, "refused", "branch_exists")
}

func TestOneOffMigrationIsValidated(t *testing.T) {
	root := legacyRepo(t)
	bad := func(mutate func(m map[string]any)) Request {
		m := oneOff()
		mutate(m)
		return Request{"mode": "plan", "migrations": []any{m}}
	}
	step := func(m map[string]any) map[string]any { return m["steps"].([]any)[0].(map[string]any) }
	for name, req := range map[string]Request{
		"unknown kind":      bad(func(m map[string]any) { step(m)["kind"] = "run_script" }),
		"unknown member":    bad(func(m map[string]any) { step(m)["command"] = "rm -rf /" }),
		"bad class":         bad(func(m map[string]any) { step(m)["class"] = "root" }),
		"dir escapes":       bad(func(m map[string]any) { step(m)["dir"] = "../elsewhere" }),
		"missing from":      bad(func(m map[string]any) { delete(step(m), "from") }),
		"bad version":       bad(func(m map[string]any) { m["to"] = "one" }),
		"no steps":          bad(func(m map[string]any) { m["steps"] = []any{} }),
		"duplicate step id": bad(func(m map[string]any) { m["steps"] = append(m["steps"].([]any), step(m)) }),
		"not objects":       {"migrations": []any{"x"}},
		"not a list":        {"migrations": "x"},
	} {
		_, err := call(t, root, "migrate", req)
		if err == nil || err.Code != "invalid_request" {
			t.Errorf("%s: got %v", name, err)
		}
	}
}
