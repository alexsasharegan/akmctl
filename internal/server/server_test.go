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

func get(t *testing.T, ts *httptest.Server, path string, hdr map[string]string) (*http.Response, []byte) {
	t.Helper()
	req, _ := http.NewRequest("GET", ts.URL+path, nil)
	for k, v := range hdr {
		req.Header.Set(k, v)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return res, b
}

func TestEmbeddedAssetsNeedNoToken(t *testing.T) {
	ts, _ := setup(t)
	cases := map[string]string{
		"/js/util.js":                        "text/javascript",
		"/fonts/monaspace-neon-var.woff2":    "font/woff2",
		"/fonts/monaspace-radon-var.woff2":   "font/woff2",
		"/fonts/monaspace-argon-var.woff2":   "font/woff2",
		"/fonts/monaspace-krypton-var.woff2": "font/woff2",
		"/fonts/monaspace-neon-LICENSE.txt":  "text/plain",
		"/fonts/NOTICE":                      "text/plain",
	}
	for p, ct := range cases {
		res, b := get(t, ts, p, nil)
		if res.StatusCode != 200 || !strings.HasPrefix(res.Header.Get("Content-Type"), ct) || len(b) == 0 {
			t.Fatalf("%s: %d %q", p, res.StatusCode, res.Header.Get("Content-Type"))
		}
		if strings.Contains(string(b), token) {
			t.Fatalf("%s holds the token", p)
		}
		csp := res.Header.Get("Content-Security-Policy")
		if !strings.Contains(csp, "font-src 'self'") || strings.Contains(csp, "unsafe-inline") || strings.Contains(csp, "unsafe-eval") {
			t.Fatalf("%s csp: %q", p, csp)
		}
		if res.Header.Get("Access-Control-Allow-Origin") != "" || res.Header.Get("Cache-Control") != "no-store" {
			t.Fatalf("%s headers: %v", p, res.Header)
		}
	}
	// the page itself carries the same policy
	res, _ := get(t, ts, "/", nil)
	if !strings.Contains(res.Header.Get("Content-Security-Policy"), "font-src 'self'") {
		t.Fatal("page csp lacks font-src")
	}
}

func TestAssetRefusals(t *testing.T) {
	ts, _ := setup(t)
	for _, p := range []string{
		"/js/../../x", "/js/%2e%2e/%2e%2e/x", "/js/..%2fapp.js", "/js/%2E%2E%2Fserver.go", "/fonts/..%5cNOTICE",
		"/js/", "/fonts/", "/js/nope.js", "/fonts/nope.woff2", "/js/util.js/", "/js/sub/util.js", "/fonts/%00",
		"/fonts/monaspace-neon-var.woff2.gz", "/js/util.css", "/web/app.js", "/js/.hidden.js",
	} {
		res, _ := get(t, ts, p, nil)
		if res.StatusCode == 200 {
			t.Fatalf("%s served", p)
		}
	}
	// a name outside the prefixes is not an asset: it needs the token and is then not found
	if res, _ := get(t, ts, "/other.js", nil); res.StatusCode != 401 {
		t.Fatalf("/other.js: %d", res.StatusCode)
	}
	// assets are GET/HEAD only
	req, _ := http.NewRequest("POST", ts.URL+"/js/util.js", nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != 401 {
		t.Fatalf("POST asset: %d", res.StatusCode)
	}
}

func TestAssetsKeepHostAndOriginRules(t *testing.T) {
	ts, _ := setup(t)
	for _, p := range []string{"/js/util.js", "/fonts/monaspace-neon-var.woff2", "/ui/config"} {
		if res, _ := get(t, ts, p, map[string]string{"Origin": "http://evil.example", "Authorization": "Bearer " + token}); res.StatusCode != 403 {
			t.Fatalf("%s foreign origin: %d", p, res.StatusCode)
		}
		req, _ := http.NewRequest("GET", ts.URL+p, nil)
		req.Host = "evil.example"
		req.Header.Set("Authorization", "Bearer "+token)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != 403 {
			t.Fatalf("%s foreign host: %d", p, res.StatusCode)
		}
	}
	// no preflight answer and no CORS headers
	req, _ := http.NewRequest("OPTIONS", ts.URL+"/ui/config", nil)
	req.Header.Set("Access-Control-Request-Method", "GET")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode == 200 || res.StatusCode == 204 || res.Header.Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("preflight: %d %v", res.StatusCode, res.Header)
	}
}

func TestUIConfig(t *testing.T) {
	ts, root := setup(t)
	if res, _ := get(t, ts, "/ui/config", nil); res.StatusCode != 401 {
		t.Fatalf("no token: %d", res.StatusCode)
	}
	if res, _ := get(t, ts, "/ui/config", map[string]string{"Authorization": "Bearer wrong"}); res.StatusCode != 401 {
		t.Fatalf("wrong token: %d", res.StatusCode)
	}
	res, b := get(t, ts, "/ui/config", map[string]string{"Authorization": "Bearer " + token})
	if res.StatusCode != 200 || !strings.HasPrefix(res.Header.Get("Content-Type"), "application/json") {
		t.Fatalf("config: %d %q", res.StatusCode, res.Header.Get("Content-Type"))
	}
	if strings.Contains(string(b), token) {
		t.Fatal("config carries the token")
	}
	var cfg map[string]string
	if err := json.Unmarshal(b, &cfg); err != nil {
		t.Fatal(err)
	}
	if !filepath.IsAbs(cfg["root"]) || cfg["root"] != root {
		t.Fatalf("root %q, want %q", cfg["root"], root)
	}
	if cfg["identity"] == "" {
		t.Fatal("no identity")
	}
	if v, ok := cfg["obsidian_vault"]; !ok || v != "" {
		t.Fatalf("obsidian_vault %q %v", v, ok)
	}
	// other methods are not the config route
	if r := do(t, ts, "POST", "/ui/config", `{}`, nil); r.status != 404 {
		t.Fatalf("POST config: %d", r.status)
	}
}

func TestUIConfigObsidianVault(t *testing.T) {
	root := t.TempDir()
	s := &Server{Root: root, Token: token, ObsidianVault: "Deep Thought", Info: api.Info{Name: "akmctl", Version: "test"}}
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()
	s.Hosts = []string{strings.TrimPrefix(ts.URL, "http://")}
	_, b := get(t, ts, "/ui/config", map[string]string{"Authorization": "Bearer " + token})
	var cfg map[string]string
	json.Unmarshal(b, &cfg)
	if cfg["obsidian_vault"] != "Deep Thought" {
		t.Fatalf("%v", cfg)
	}
}
