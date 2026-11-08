package conformance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/api"
	"github.com/alexsasharegan/akmctl/internal/cli"
	"github.com/alexsasharegan/akmctl/internal/record"
	"github.com/alexsasharegan/akmctl/internal/server"
)

// A binding runs one request on a repository directory and returns the response
// body as JSON. injected lists the members of the echoed request that the
// binding added itself (an If-Match turned into if_match, a path segment): they
// are removed from the echo before it is compared. skip says why the binding
// cannot run this operation.
type binding struct {
	name string
	run  func(t *testing.T, root, op string, req map[string]any, noIfMatch bool) (body []byte, injected []string, skip string)
}

var bindings = []binding{{"api", runAPI}, {"http", runHTTP}, {"cli", runCLI}}

var info = api.Info{Name: "akmctl", Version: "test"}

func runAPI(t *testing.T, root, op string, req map[string]any, _ bool) ([]byte, []string, string) {
	t.Helper()
	repo, err := record.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := api.Run(repo, info, op, req)
	b, err := api.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	return b, nil, ""
}

// ---- HTTP -------------------------------------------------------------

const token = "0123456789abcdef0123456789abcdef"

// needsIfMatch are the operations the HTTP binding refuses without an If-Match
// when they change a record (spec section 6).
var needsIfMatch = map[string]bool{"set": true, "section": true, "tick": true, "log": true, "close": true, "claim": true, "release": true}

// route is the HTTP request for an operation: method, path, whether the other
// parameters go in the query (reads) or the body, and the parameters the path
// carries.
type route struct {
	method, path string
	query        bool
	inPath       []string
}

func routeOf(op string, req map[string]any) (route, string) {
	id, _ := req["id"].(string)
	switch op {
	case "version", "search", "board", "check", "commits", "history":
		return route{"GET", "/" + op, true, nil}, ""
	case "ls":
		return route{"GET", "/records", true, nil}, ""
	case "show", "when":
		p := "/when/" + id
		if op == "show" {
			p = "/records/" + id
		}
		return route{"GET", p, true, []string{"id"}}, ""
	case "find":
		if id != "" {
			return route{"GET", "/find/" + id, true, []string{"id"}}, ""
		}
		return route{"GET", "/find", true, nil}, ""
	case "new":
		return route{"POST", "/records", false, nil}, ""
	case "index", "bases", "init", "migrate":
		return route{"POST", "/" + op, false, nil}, ""
	case "config":
		if _, ok := req["value"]; ok {
			return route{"POST", "/config", false, nil}, ""
		}
		return route{"GET", "/config", true, nil}, ""
	case "skills":
		if m, _ := req["mode"].(string); m == "install" || m == "update" {
			return route{"POST", "/skills", false, nil}, ""
		}
		return route{"GET", "/skills", true, nil}, ""
	case "set":
		return route{"PATCH", "/records/" + id, false, []string{"id"}}, ""
	case "section":
		name, _ := req["name"].(string)
		return route{"PUT", "/records/" + id + "/sections/" + url.PathEscape(name), false, []string{"id", "name"}}, ""
	case "tick":
		sec, _ := req["section"].(string)
		if sec == "" {
			sec = "goal"
		}
		n, _ := req["item"].(float64)
		return route{"PUT", fmt.Sprintf("/records/%s/sections/%s/items/%d", id, url.PathEscape(sec), int(n)), false, []string{"id"}}, ""
	case "log", "close":
		return route{"POST", "/records/" + id + "/" + op, false, []string{"id"}}, ""
	case "claim":
		return route{"POST", "/records/" + id + "/claim", false, []string{"id"}}, ""
	case "release":
		return route{"DELETE", "/records/" + id + "/claim", false, []string{"id"}}, ""
	}
	return route{}, "the HTTP binding has no route for " + op
}

func queryOf(req map[string]any, skip []string) (url.Values, string) {
	q := url.Values{}
	for k, v := range req {
		if contains(skip, k) {
			continue
		}
		switch x := v.(type) {
		case string:
			q.Set(k, x)
		case bool:
			q.Set(k, strconv.FormatBool(x))
		case float64:
			q.Set(k, strconv.Itoa(int(x)))
		case []any:
			var parts []string
			for _, e := range x {
				s, ok := e.(string)
				if !ok {
					return nil, "a list that is not strings cannot go in a query"
				}
				parts = append(parts, s)
			}
			q.Set(k, strings.Join(parts, ","))
		default:
			return nil, "a " + fmt.Sprintf("%T", v) + " cannot go in a query"
		}
	}
	return q, ""
}

func runHTTP(t *testing.T, root, op string, req map[string]any, noIfMatch bool) ([]byte, []string, string) {
	t.Helper()
	rt, why := routeOf(op, req)
	if why != "" {
		return nil, nil, why
	}
	s := &server.Server{Root: root, Token: token, Info: info}
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()
	s.Hosts = []string{strings.TrimPrefix(ts.URL, "http://")}
	do := func(method, path string, body []byte, hdr map[string]string) (int, []byte) {
		rq, err := http.NewRequest(method, ts.URL+"/api/v1"+path, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		rq.Header.Set("Authorization", "Bearer "+token)
		if len(body) > 0 {
			rq.Header.Set("Content-Type", "application/json")
		}
		for k, v := range hdr {
			rq.Header.Set(k, v)
		}
		resp, err := http.DefaultClient.Do(rq)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, b
	}
	var injected []string
	hdr := map[string]string{}
	rest := map[string]any{}
	for k, v := range req {
		if !contains(rt.inPath, k) {
			rest[k] = v
		}
	}
	if op == "tick" {
		if _, ok := req["section"]; !ok {
			injected = append(injected, "section")
		}
	}
	// a write that changes a record carries the record's current etag, as a
	// client of the HTTP binding must, unless the case brings its own
	if needsIfMatch[op] && req["dry_run"] != true && !noIfMatch {
		_, own := req["if_match"]
		_, sec := req["section_etag"]
		if !own && !sec {
			st, b := do("GET", "/records/"+fmt.Sprint(req["id"]), nil, nil)
			var rec struct {
				Record struct{ Etag string } `json:"record"`
			}
			if json.Unmarshal(b, &rec); st != 200 || rec.Record.Etag == "" {
				return nil, nil, "cannot read the etag of the record first"
			}
			hdr["If-Match"] = `"` + rec.Record.Etag + `"`
			injected = append(injected, "if_match")
		}
	}
	var status int
	var out []byte
	if rt.query {
		q, why := queryOf(rest, nil)
		if why != "" {
			return nil, nil, why
		}
		p := rt.path
		if len(q) > 0 {
			p += "?" + q.Encode()
		}
		status, out = do(rt.method, p, nil, hdr)
	} else {
		var body []byte
		if len(rest) > 0 || rt.method != "DELETE" {
			body, _ = json.Marshal(rest)
		}
		status, out = do(rt.method, rt.path, body, hdr)
	}
	// the HTTP status must agree with the envelope (spec section 2)
	var env struct {
		OK    bool `json:"ok"`
		Error struct {
			Code    string
			Details map[string]any
		} `json:"error"`
	}
	if err := json.Unmarshal(out, &env); err != nil {
		t.Fatalf("not JSON (status %d): %s", status, out)
	}
	if want := statusWant(op, env.OK, env.Error.Code, env.Error.Details); status != want {
		t.Errorf("HTTP status %d, want %d for %s", status, want, string(out[:min(len(out), 200)]))
	}
	return bytes.TrimSpace(out), injected, ""
}

// statusWant is the table of spec section 2.
func statusWant(op string, ok bool, code string, details map[string]any) int {
	if ok {
		if op == "new" {
			return 201
		}
		return 200
	}
	if r, _ := details["reason"].(string); code == "refused" && r == "if_match_required" {
		return 428
	}
	switch code {
	case "invalid_request":
		return 400
	case "unauthorized":
		return 401
	case "not_found":
		return 404
	case "ambiguous", "refused":
		return 409
	case "stale":
		return 412
	}
	return 500
}

// ---- CLI --------------------------------------------------------------

// cliReads are the operations the CLI prints as JSON with --format json and that
// change nothing, so one run answers them. A write goes through the CLI's own
// prompts and flags, which have their tests.
var cliReads = map[string]bool{"version": true, "ls": true, "search": true, "show": true, "board": true, "find": true, "when": true, "check": true, "commits": true, "history": true, "config": true, "skills": true}

// cliPositional are the parameters the CLI takes as words (internal/cli).
var cliPositional = map[string][]string{"search": {"text"}, "show": {"id"}, "find": {"id"}, "when": {"id"}}

func runCLI(t *testing.T, root, op string, req map[string]any, _ bool) ([]byte, []string, string) {
	t.Helper()
	if !cliReads[op] {
		return nil, nil, "the CLI is run for reads only"
	}
	if op == "config" {
		if _, ok := req["value"]; ok {
			return nil, nil, "the CLI is run for reads only"
		}
	}
	if m, _ := req["mode"].(string); op == "skills" && m != "" && m != "status" {
		return nil, nil, "the CLI is run for reads only"
	}
	var args []string
	args = append(args, op)
	words := map[string]bool{}
	for _, name := range cliPositional[op] {
		if v, ok := req[name].(string); ok {
			args = append(args, v)
			words[name] = true
		}
	}
	var keys []string
	for k := range req {
		if !words[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		flag := "--" + strings.ReplaceAll(k, "_", "-")
		switch x := req[k].(type) {
		case string:
			args = append(args, flag, x)
		case bool:
			if x {
				args = append(args, flag)
			} else {
				args = append(args, "--no-"+strings.ReplaceAll(k, "_", "-"))
			}
		case float64:
			args = append(args, flag, strconv.Itoa(int(x)))
		case []any:
			var parts []string
			for _, e := range x {
				s, ok := e.(string)
				if !ok {
					return nil, nil, "a list that is not strings has no flag form"
				}
				parts = append(parts, s)
			}
			args = append(args, flag, strings.Join(parts, ","))
		default:
			return nil, nil, "a " + fmt.Sprintf("%T", x) + " has no flag form"
		}
	}
	args = append(args, "--repo", root, "--format", "json")
	var out, errb bytes.Buffer
	old := cli.Version
	cli.Version = info.Version
	code := cli.Run(args, &out, &errb, root)
	cli.Version = old
	if out.Len() == 0 {
		t.Fatalf("akmctl %v printed nothing (exit %d): %s", args, code, errb.String())
	}
	return bytes.TrimSpace(out.Bytes()), nil, ""
}

func contains(l []string, s string) bool {
	for _, x := range l {
		if x == s {
			return true
		}
	}
	return false
}
