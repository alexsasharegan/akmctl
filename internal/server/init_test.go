package server

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitRoute(t *testing.T) {
	ts, root := setup(t)
	r := do(t, ts, "POST", "/api/v1/init", `{"profile":"notes","dir":"vault"}`, nil)
	if r.status != 200 || r.body["ok"] != true || r.body["operation"] != "init" {
		t.Fatalf("init: %d %v", r.status, r.body)
	}
	if _, err := os.Stat(filepath.Join(root, "vault", "notes", "work", "done")); err != nil {
		t.Error(err)
	}
	r = do(t, ts, "POST", "/api/v1/init", `{"profile":"notes","dir":"../out"}`, nil)
	if r.status != 400 || errCode(r) != "invalid_request" {
		t.Fatalf("dir outside the repository: %d %v", r.status, r.body)
	}
}
