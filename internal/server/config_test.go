package server

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigRoutes(t *testing.T) {
	ts, root := setup(t)
	r := do(t, ts, "GET", "/api/v1/config", "", nil)
	if r.status != 200 || r.body["operation"] != "config" || r.body["settings"] == nil {
		t.Fatalf("list: %d %v", r.status, r.body)
	}
	r = do(t, ts, "GET", "/api/v1/config?key=git_commits", "", nil)
	if r.status != 200 || len(r.body["settings"].([]any)) != 1 {
		t.Fatalf("get: %d %v", r.status, r.body)
	}
	// no work tree: refused, 409
	r = do(t, ts, "POST", "/api/v1/config", `{"key":"git_commits","value":true}`, nil)
	if r.status != 409 || errCode(r) != "refused" {
		t.Fatalf("outside a work tree: %d %v", r.status, r.body)
	}
	os.MkdirAll(filepath.Join(root, ".git"), 0o755)
	r = do(t, ts, "POST", "/api/v1/config", `{"key":"git_commits","value":true,"dry_run":true}`, nil)
	if r.status != 200 || r.body["applied"] != false {
		t.Fatalf("dry run: %d %v", r.status, r.body)
	}
	r = do(t, ts, "POST", "/api/v1/config", `{"key":"git_commits","value":true}`, nil)
	if r.status != 200 || r.body["applied"] != true {
		t.Fatalf("on: %d %v", r.status, r.body)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents", "skills", "akm-commit", "SKILL.md")); err != nil {
		t.Error(err)
	}
	r = do(t, ts, "POST", "/api/v1/config", `{"key":"git_commits","value":false}`, nil)
	if r.status != 409 {
		t.Fatalf("off without confirm: %d %v", r.status, r.body)
	}
	r = do(t, ts, "POST", "/api/v1/config", `{"key":"git_commits","value":false,"confirm":true}`, nil)
	if r.status != 200 {
		t.Fatalf("off: %d %v", r.status, r.body)
	}
	r = do(t, ts, "POST", "/api/v1/config", `{"key":"color","value":false}`, nil)
	if r.status != 400 {
		t.Fatalf("unknown key: %d", r.status)
	}
}

func TestSkillsRoutes(t *testing.T) {
	ts, root := setup(t)
	r := do(t, ts, "GET", "/api/v1/skills", "", nil)
	if r.status != 200 || r.body["operation"] != "skills" || r.body["clean"] != false {
		t.Fatalf("status: %d %v", r.status, r.body)
	}
	r = do(t, ts, "GET", "/api/v1/skills?mode=install", "", nil)
	if r.status != 400 || errCode(r) != "invalid_request" {
		t.Fatalf("a GET must not write: %d %v", r.status, r.body)
	}
	r = do(t, ts, "POST", "/api/v1/skills", `{"mode":"install","confirm":true}`, nil)
	if r.status != 200 || r.body["applied"] != true {
		t.Fatalf("install: %d %v", r.status, r.body)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents/skills/akm/.akm-skill.json")); err != nil {
		t.Fatal("marker not written")
	}
}
