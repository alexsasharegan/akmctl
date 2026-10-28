package server

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/api"
)

const token = "0123456789abcdef0123456789abcdef"

func setup(t *testing.T) (*httptest.Server, string) {
	t.Helper()
	root := t.TempDir()
	src := "../../testdata/api/repos/claims"
	filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		rel, _ := filepath.Rel(src, p)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(root, rel), 0o755)
		}
		b, _ := os.ReadFile(p)
		return os.WriteFile(filepath.Join(root, rel), b, 0o644)
	})
	s := &Server{Root: root, Token: token, Info: api.Info{Name: "akmctl", Version: "test"}}
	ts := httptest.NewServer(s.Handler())
	s.Hosts = []string{strings.TrimPrefix(ts.URL, "http://")}
	t.Cleanup(ts.Close)
	return ts, root
}

type reply struct {
	status int
	h      http.Header
	body   map[string]any
}

func do(t *testing.T, ts *httptest.Server, method, path string, body string, hdr map[string]string) reply {
	t.Helper()
	var rd io.Reader
	if body != "" {
		rd = strings.NewReader(body)
	}
	req, _ := http.NewRequest(method, ts.URL+path, rd)
	if _, ok := hdr["Authorization"]; !ok {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range hdr {
		if v == "-" {
			req.Header.Del(k)
		} else if k == "Host" {
			req.Host = v
		} else {
			req.Header.Set(k, v)
		}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	r := reply{res.StatusCode, res.Header, nil}
	json.Unmarshal(b, &r.body)
	return r
}

func errCode(r reply) string {
	e, _ := r.body["error"].(map[string]any)
	c, _ := e["code"].(string)
	return c
}

func reason(r reply) string {
	e, _ := r.body["error"].(map[string]any)
	d, _ := e["details"].(map[string]any)
	s, _ := d["reason"].(string)
	return s
}

func TestAuth(t *testing.T) {
	ts, _ := setup(t)
	if r := do(t, ts, "GET", "/api/v1/version", "", map[string]string{"Authorization": "-"}); r.status != 401 || errCode(r) != "unauthorized" {
		t.Fatalf("no token: %d %v", r.status, r.body)
	}
	if r := do(t, ts, "GET", "/api/v1/version", "", map[string]string{"Authorization": "Bearer nope"}); r.status != 401 {
		t.Fatalf("wrong token: %d", r.status)
	}
	if r := do(t, ts, "GET", "/api/v1/version", "", map[string]string{"Authorization": "Basic " + token}); r.status != 401 {
		t.Fatalf("wrong scheme: %d", r.status)
	}
	if r := do(t, ts, "GET", "/api/v1/version", "", nil); r.status != 200 || r.body["ok"] != true {
		t.Fatalf("with token: %d %v", r.status, r.body)
	}
	// every route needs the token, reads included
	for _, p := range []string{"/api/v1/records", "/api/v1/board", "/api/v1/check", "/api/v1/records/261022-e5f6", "/api/v1/nowhere"} {
		if r := do(t, ts, "GET", p, "", map[string]string{"Authorization": "-"}); r.status != 401 {
			t.Errorf("%s without a token: %d", p, r.status)
		}
	}
}

func TestHostAndOrigin(t *testing.T) {
	ts, _ := setup(t)
	if r := do(t, ts, "GET", "/api/v1/version", "", map[string]string{"Host": "evil.example:80"}); r.status != 403 {
		t.Fatalf("bad host: %d", r.status)
	}
	if r := do(t, ts, "GET", "/api/v1/version", "", map[string]string{"Origin": "http://evil.example"}); r.status != 403 {
		t.Fatalf("cross origin: %d", r.status)
	}
	if r := do(t, ts, "GET", "/api/v1/version", "", map[string]string{"Origin": "https://" + strings.TrimPrefix(ts.URL, "http://")}); r.status != 403 {
		t.Fatalf("https origin: %d", r.status)
	}
	if r := do(t, ts, "GET", "/api/v1/version", "", map[string]string{"Origin": ts.URL}); r.status != 200 {
		t.Fatalf("same origin: %d", r.status)
	}
	r := do(t, ts, "OPTIONS", "/api/v1/records", "", map[string]string{"Authorization": "-", "Origin": "http://evil.example"})
	if r.status == 200 || r.h.Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("preflight answered: %d %v", r.status, r.h)
	}
}

func TestPageNeedsNoTokenAndHoldsNoData(t *testing.T) {
	ts, _ := setup(t)
	for _, p := range []string{"/", "/app.js", "/app.css"} {
		res, err := http.Get(ts.URL + p)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode != 200 || strings.Contains(string(b), token) || strings.Contains(string(b), "261022-e5f6") {
			t.Fatalf("%s: %d", p, res.StatusCode)
		}
		if !strings.Contains(res.Header.Get("Content-Security-Policy"), "default-src 'none'") || res.Header.Get("Cache-Control") != "no-store" {
			t.Fatalf("%s headers: %v", p, res.Header)
		}
	}
}

func TestReads(t *testing.T) {
	ts, _ := setup(t)
	r := do(t, ts, "GET", "/api/v1/records?type=work&status=ready,active&sort=roadmap,-priority&fields=title&limit=2&next=true", "", nil)
	if r.status != 200 || r.body["count"] != float64(2) || r.body["total"] != float64(3) {
		t.Fatalf("%d %v", r.status, r.body)
	}
	req, _ := r.body["request"].(map[string]any)
	if req["next"] != true || req["limit"] != float64(2) {
		t.Fatalf("request echo %v", req)
	}
	if r := do(t, ts, "GET", "/api/v1/records/261020-a2b3?section=next", "", nil); r.status != 200 || r.h.Get("ETag") == "" {
		t.Fatalf("show: %d %v", r.status, r.h)
	}
	if r := do(t, ts, "GET", "/api/v1/records/261099-zzzz", "", nil); r.status != 404 || errCode(r) != "not_found" {
		t.Fatalf("missing: %d", r.status)
	}
	if r := do(t, ts, "GET", "/api/v1/records?bogus=1", "", nil); r.status != 400 {
		t.Fatalf("unknown param: %d", r.status)
	}
	if r := do(t, ts, "GET", "/api/v1/find/261022-e5f6", "", nil); r.status != 200 {
		t.Fatalf("find: %d", r.status)
	}
	if r := do(t, ts, "GET", "/api/v1/when/261022-e5f6?as_of=2026-10-28", "", nil); r.status != 200 || r.body["age_days"] != float64(6) {
		t.Fatalf("when: %d %v", r.status, r.body)
	}
	if r := do(t, ts, "GET", "/api/v1/check?as_of=2026-10-28", "", nil); r.status != 200 || r.body["clean"] == nil {
		t.Fatalf("check: %d", r.status)
	}
	if r := do(t, ts, "GET", "/api/v1/board", "", nil); r.status != 200 {
		t.Fatalf("board: %d", r.status)
	}
	if r := do(t, ts, "GET", "/api/v1/nowhere", "", nil); r.status != 404 {
		t.Fatalf("route: %d", r.status)
	}
	if r := do(t, ts, "DELETE", "/api/v1/records", "", nil); r.status != 404 {
		t.Fatalf("method: %d", r.status)
	}
}

func etagOf(t *testing.T, ts *httptest.Server, id string) string {
	t.Helper()
	r := do(t, ts, "GET", "/api/v1/records/"+id, "", nil)
	return strings.Trim(r.h.Get("ETag"), `"`)
}

func TestWritesNeedIfMatch(t *testing.T) {
	ts, root := setup(t)
	path := "/api/v1/records/261022-e5f6"
	r := do(t, ts, "PATCH", path, `{"changes":{"priority":3}}`, nil)
	if r.status != 428 || errCode(r) != "refused" || reason(r) != "if_match_required" {
		t.Fatalf("no If-Match: %d %v", r.status, r.body)
	}
	b, _ := os.ReadFile(filepath.Join(root, "docs/work/261022-e5f6-document-the-exit-codes.md"))
	if strings.Contains(string(b), "priority") {
		t.Fatal("the write went through")
	}
	et := etagOf(t, ts, "261022-e5f6")
	r = do(t, ts, "PATCH", path, `{"changes":{"priority":3}}`, map[string]string{"If-Match": `"` + et + `"`})
	if r.status != 200 || r.h.Get("ETag") == "" || r.h.Get("ETag") == `"`+et+`"` {
		t.Fatalf("with If-Match: %d %v", r.status, r.body)
	}
	r = do(t, ts, "PATCH", path, `{"changes":{"priority":4}}`, map[string]string{"If-Match": `"` + et + `"`})
	if r.status != 412 || errCode(r) != "stale" {
		t.Fatalf("stale: %d %v", r.status, r.body)
	}
	// the etag may also come in the body
	et = etagOf(t, ts, "261022-e5f6")
	r = do(t, ts, "PATCH", path, `{"changes":{"priority":5},"if_match":"`+et+`"}`, nil)
	if r.status != 200 {
		t.Fatalf("body if_match: %d %v", r.status, r.body)
	}
}

func TestWriteRoutes(t *testing.T) {
	ts, _ := setup(t)
	r := do(t, ts, "POST", "/api/v1/records", `{"kind":"work","slug":"from-http","title":"From HTTP","as_of":"2026-10-28"}`, nil)
	if r.status != 201 || r.body["ok"] != true {
		t.Fatalf("new: %d %v", r.status, r.body)
	}
	rec := r.body["record"].(map[string]any)
	id := rec["id"].(string)
	// claim, release
	et := strings.Trim(r.h.Get("ETag"), `"`)
	r = do(t, ts, "PATCH", "/api/v1/records/"+id, `{"changes":{"status":"ready"}}`, map[string]string{"If-Match": et})
	if r.status != 200 {
		t.Fatal(r.body)
	}
	et = strings.Trim(r.h.Get("ETag"), `"`)
	r = do(t, ts, "POST", "/api/v1/records/"+id+"/claim", `{"as":"http-agent","as_of":"2026-10-28T10:00"}`, map[string]string{"If-Match": et})
	if r.status != 200 {
		t.Fatalf("claim: %d %v", r.status, r.body)
	}
	et = strings.Trim(r.h.Get("ETag"), `"`)
	r = do(t, ts, "POST", "/api/v1/records/"+id+"/close", `{"status":"done"}`, map[string]string{"If-Match": et})
	if r.status != 409 || reason(r) != "empty_outcome" {
		t.Fatalf("close: %d %v", r.status, r.body)
	}
	r = do(t, ts, "DELETE", "/api/v1/records/"+id+"/claim", "", map[string]string{"If-Match": et})
	if r.status != 200 {
		t.Fatalf("release: %d %v", r.status, r.body)
	}
	if r := do(t, ts, "DELETE", "/api/v1/records/"+id+"/claim", "", nil); r.status != 428 {
		t.Fatalf("release without If-Match: %d", r.status)
	}
	if r := do(t, ts, "POST", "/api/v1/migrate", `{"mode":"report"}`, nil); r.status != 200 || r.body["state"] != "current" {
		t.Fatalf("migrate: %d %v", r.status, r.body)
	}
	if r := do(t, ts, "POST", "/api/v1/migrate", `{"mode":"apply_safe"}`, nil); r.status != 409 || reason(r) != "no_backup_ack" {
		t.Fatalf("migrate apply: %d %v", r.status, r.body)
	}
	if r := do(t, ts, "POST", "/api/v1/index", `{"dry_run":true}`, nil); r.status != 200 {
		t.Fatalf("index: %d %v", r.status, r.body)
	}
}

func TestBadWrites(t *testing.T) {
	ts, _ := setup(t)
	if r := do(t, ts, "POST", "/api/v1/records", `{"kind":"work","slug":"x"}`, map[string]string{"Content-Type": "text/plain"}); r.status != 400 {
		t.Fatalf("content type: %d", r.status)
	}
	if r := do(t, ts, "POST", "/api/v1/records", `not json`, nil); r.status != 400 {
		t.Fatalf("bad json: %d", r.status)
	}
	if r := do(t, ts, "PATCH", "/api/v1/records/261022-e5f6", `{"id":"261020-a2b3","changes":{"priority":1}}`, map[string]string{"If-Match": "0000000000000000"}); r.status != 400 {
		t.Fatalf("id mismatch: %d", r.status)
	}
}

func TestListensOnLoopbackOnly(t *testing.T) {
	l, err := Listen(0)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	host, _, _ := net.SplitHostPort(l.Addr().String())
	if host != "127.0.0.1" {
		t.Fatalf("bound to %s", host)
	}
}

func TestServeStopsOnCancel(t *testing.T) {
	l, _ := Listen(0)
	s := &Server{Root: t.TempDir(), Token: token}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- s.Serve(ctx, l) }()
	cancel()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestToken(t *testing.T) {
	a, _ := NewToken()
	b, _ := NewToken()
	if len(a) != 32 || a == b {
		t.Fatalf("%q %q", a, b)
	}
}

// A dry run writes nothing, so it needs no If-Match; a real write still does.
func TestDryRunNeedsNoIfMatch(t *testing.T) {
	ts, root := setup(t)
	file := filepath.Join(root, "docs/work/261022-e5f6-document-the-exit-codes.md")
	before, _ := os.ReadFile(file)
	et := etagOf(t, ts, "261022-e5f6")
	routes := []struct{ method, path, body string }{
		{"PATCH", "/api/v1/records/261022-e5f6", `{"changes":{"priority":3},"dry_run":true}`},
		{"PUT", "/api/v1/records/261022-e5f6/sections/Goal", `{"text":"New.","dry_run":true}`},
		{"POST", "/api/v1/records/261022-e5f6/close", `{"status":"dropped","outcome":"No.","dry_run":true}`},
		{"POST", "/api/v1/records/261022-e5f6/claim", `{"as":"a","dry_run":true}`},
		{"DELETE", "/api/v1/records/261020-a2b3/claim", `{"dry_run":true}`},
	}
	for _, rt := range routes {
		r := do(t, ts, rt.method, rt.path, rt.body, nil)
		if r.status != 200 {
			t.Fatalf("%s %s dry run: %d %v", rt.method, rt.path, r.status, r.body)
		}
		d, _ := r.body["diff"].(map[string]any)
		if hs, _ := d["hunks"].([]any); len(hs) == 0 {
			t.Fatalf("%s %s: no hunks: %v", rt.method, rt.path, r.body)
		}
		if got := r.h.Get("ETag"); got != `"`+et+`"` && rt.method != "DELETE" {
			t.Fatalf("%s: ETag %s, want the current %s", rt.path, got, et)
		}
		// without dry_run, or with dry_run false, the same request needs If-Match
		real := strings.Replace(strings.Replace(rt.body, `,"dry_run":true`, "", 1), `{"dry_run":true}`, "", 1)
		if r := do(t, ts, rt.method, rt.path, real, nil); r.status != 428 || reason(r) != "if_match_required" {
			t.Fatalf("%s %s real write without If-Match: %d %v", rt.method, rt.path, r.status, r.body)
		}
		off := strings.Replace(rt.body, `"dry_run":true`, `"dry_run":false`, 1)
		if r := do(t, ts, rt.method, rt.path, off, nil); r.status != 428 {
			t.Fatalf("%s %s dry_run false without If-Match: %d", rt.method, rt.path, r.status)
		}
	}
	after, _ := os.ReadFile(file)
	if string(before) != string(after) {
		t.Fatal("a dry run or a refused write changed the file")
	}
	// a stale If-Match still fails a dry run
	r := do(t, ts, "PATCH", "/api/v1/records/261022-e5f6", `{"changes":{"priority":3},"dry_run":true}`, map[string]string{"If-Match": `"0000000000000000"`})
	if r.status != 412 {
		t.Fatalf("stale dry run: %d", r.status)
	}
}

func TestSectionRoute(t *testing.T) {
	ts, root := setup(t)
	path := "/api/v1/records/261022-e5f6/sections/Goal"
	et := etagOf(t, ts, "261022-e5f6")
	if r := do(t, ts, "PUT", path, `{"text":"New goal."}`, nil); r.status != 428 {
		t.Fatalf("no If-Match: %d", r.status)
	}
	if r := do(t, ts, "PUT", path, `{"text":"x","name":"Other"}`, map[string]string{"If-Match": `"` + et + `"`}); r.status != 400 {
		t.Fatalf("name mismatch: %d", r.status)
	}
	if r := do(t, ts, "PUT", path, `{"text":"New goal."}`, map[string]string{"Content-Type": "text/plain", "If-Match": `"` + et + `"`}); r.status != 400 {
		t.Fatalf("content type: %d", r.status)
	}
	r := do(t, ts, "PUT", path, `{"text":"New goal.","as_of":"2026-10-28"}`, map[string]string{"If-Match": `"` + et + `"`})
	if r.status != 200 || r.body["changed"] != true {
		t.Fatalf("put: %d %v", r.status, r.body)
	}
	if sec, _ := r.body["section"].(map[string]any); sec["name"] != "goal" || sec["text"] != "New goal." {
		t.Fatalf("section: %v", r.body["section"])
	}
	if _, ok := r.body["diff"]; ok {
		t.Fatal("diff on a real write")
	}
	b, _ := os.ReadFile(filepath.Join(root, "docs/work/261022-e5f6-document-the-exit-codes.md"))
	if !strings.Contains(string(b), "## Goal\nNew goal.\n") {
		t.Fatal(string(b))
	}
	if r := do(t, ts, "PUT", path, `{"text":"again"}`, map[string]string{"If-Match": `"` + et + `"`}); r.status != 412 {
		t.Fatalf("stale: %d", r.status)
	}
	if r := do(t, ts, "POST", path, `{"text":"x"}`, nil); r.status != 404 {
		t.Fatalf("wrong method: %d", r.status)
	}
	if r := do(t, ts, "PUT", "/api/v1/records/261010-j9k0/sections/Goal", `{"text":"x","dry_run":true}`, nil); r.status != 409 || reason(r) != "closed" {
		t.Fatalf("closed: %d %v", r.status, r.body)
	}
	// the version route reports the repository
	v := do(t, ts, "GET", "/api/v1/version", "", nil)
	repo, _ := v.body["repository"].(map[string]any)
	if repo["layout"] != "1.0.0" || repo["reference_style"] != "id" {
		t.Fatalf("version: %v", v.body)
	}
}
