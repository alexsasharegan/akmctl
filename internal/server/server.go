// Package server is the HTTP binding of the API (loopback only) and the
// workbench page that uses it. See the API spec, section 6, for the rules.
package server

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/alexsasharegan/akmctl/internal/api"
	"github.com/alexsasharegan/akmctl/internal/record"
)

//go:embed web/*
var web embed.FS

// Server serves one repository.
type Server struct {
	Root  string
	Token string
	Info  api.Info
	// Hosts are the Host header values accepted (host:port).
	Hosts []string
}

// NewToken returns 128 bits of randomness as hex.
func NewToken() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

// Listen opens the loopback listener. Port 0 picks a free port. There is no
// option to bind another interface.
func Listen(port int) (net.Listener, error) {
	return net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)))
}

// Serve runs the server on l until ctx is cancelled.
func (s *Server) Serve(ctx context.Context, l net.Listener) error {
	_, port, _ := net.SplitHostPort(l.Addr().String())
	s.Hosts = []string{"127.0.0.1:" + port, "localhost:" + port, "[::1]:" + port}
	srv := &http.Server{
		Handler:           s.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    32 << 10,
	}
	go func() {
		<-ctx.Done()
		sh, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		srv.Shutdown(sh)
	}()
	err := srv.Serve(l)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// Handler returns the HTTP handler.
func (s *Server) Handler() http.Handler {
	return http.HandlerFunc(s.serve)
}

var statusOf = map[string]int{
	"invalid_request": http.StatusBadRequest, "unauthorized": http.StatusUnauthorized, "not_found": http.StatusNotFound,
	"ambiguous": http.StatusConflict, "stale": http.StatusPreconditionFailed, "refused": http.StatusConflict,
	"unreadable": http.StatusInternalServerError, "internal": http.StatusInternalServerError,
}

func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("Cache-Control", "no-store")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Referrer-Policy", "no-referrer")
	if !s.hostOK(r.Host) {
		s.fail(w, "http", http.StatusForbidden, &api.Error{Code: "unauthorized", Message: "unexpected Host header"})
		return
	}
	if o := r.Header.Get("Origin"); o != "" && !s.originOK(o) {
		s.fail(w, "http", http.StatusForbidden, &api.Error{Code: "unauthorized", Message: "cross-origin request refused"})
		return
	}
	// the page and its assets hold no data, so they need no token
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		switch r.URL.Path {
		case "/", "/index.html":
			s.asset(w, "web/index.html", "text/html; charset=utf-8")
			return
		case "/app.js":
			s.asset(w, "web/app.js", "text/javascript; charset=utf-8")
			return
		case "/app.css":
			s.asset(w, "web/app.css", "text/css; charset=utf-8")
			return
		}
	}
	if !s.tokenOK(r.Header.Get("Authorization")) {
		w.Header().Set("WWW-Authenticate", "Bearer")
		s.fail(w, "http", http.StatusUnauthorized, &api.Error{Code: "unauthorized", Message: "missing or wrong token"})
		return
	}
	const base = "/api/v1/"
	if !strings.HasPrefix(r.URL.Path, base) {
		s.fail(w, "http", http.StatusNotFound, &api.Error{Code: "not_found", Message: "no such route"})
		return
	}
	op, req, ok, ferr := s.route(r, strings.Split(strings.Trim(strings.TrimPrefix(r.URL.Path, base), "/"), "/"))
	if ferr != nil {
		s.fail(w, "http", statusFor(ferr), ferr)
		return
	}
	if !ok {
		s.fail(w, "http", http.StatusNotFound, &api.Error{Code: "not_found", Message: "no such route"})
		return
	}
	s.run(w, op, req)
}

func (s *Server) hostOK(host string) bool {
	for _, h := range s.Hosts {
		if strings.EqualFold(h, host) {
			return true
		}
	}
	return false
}

func (s *Server) originOK(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil || u.Scheme != "http" {
		return false
	}
	return s.hostOK(u.Host)
}

func (s *Server) tokenOK(auth string) bool {
	const p = "Bearer "
	if !strings.HasPrefix(auth, p) || s.Token == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(auth[len(p):]), []byte(s.Token)) == 1
}

func (s *Server) asset(w http.ResponseWriter, name, ctype string) {
	b, err := web.ReadFile(name)
	if err != nil {
		http.NotFound(w, nil)
		return
	}
	h := w.Header()
	h.Set("Content-Type", ctype)
	h.Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; connect-src 'self'; img-src 'self' data:; base-uri 'none'; form-action 'none'; frame-ancestors 'none'")
	w.Write(b)
}

func statusFor(e *api.Error) int {
	if e.Code == "refused" {
		if r, _ := e.Details.Get("reason"); r == "if_match_required" {
			return http.StatusPreconditionRequired
		}
	}
	if st, ok := statusOf[e.Code]; ok {
		return st
	}
	return http.StatusInternalServerError
}

func (s *Server) fail(w http.ResponseWriter, op string, status int, e *api.Error) {
	s.write(w, status, api.Failure(op, api.Request{}, e))
}

func (s *Server) write(w http.ResponseWriter, status int, body api.Obj) {
	b, err := api.Marshal(body)
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if rec, ok := body.Get("record"); ok {
		if ro, ok := rec.(api.Obj); ok {
			if et, ok := ro.Get("etag"); ok {
				w.Header().Set("ETag", `"`+fmt.Sprint(et)+`"`)
			}
		}
	}
	w.WriteHeader(status)
	w.Write(append(b, '\n'))
}

func (s *Server) run(w http.ResponseWriter, op string, req api.Request) {
	repo, err := record.Load(s.Root)
	if err != nil {
		s.fail(w, op, http.StatusInternalServerError, &api.Error{Code: "unreadable", Message: err.Error()})
		return
	}
	resp, aerr := api.Run(repo, s.Info, op, req)
	if aerr != nil {
		s.write(w, statusFor(aerr), resp)
		return
	}
	status := http.StatusOK
	if op == "new" {
		status = http.StatusCreated
	}
	s.write(w, status, resp)
}

// writeOps need an If-Match when they change an existing record.
var needsIfMatch = map[string]bool{"set": true, "close": true, "claim": true, "release": true}

// route maps a method and path to an operation and its request.
func (s *Server) route(r *http.Request, seg []string) (string, api.Request, bool, *api.Error) {
	m := r.Method
	var op string
	var id string
	switch {
	case len(seg) == 1 && seg[0] == "version" && m == "GET":
		op = "version"
	case len(seg) == 1 && seg[0] == "records" && m == "GET":
		op = "ls"
	case len(seg) == 1 && seg[0] == "records" && m == "POST":
		op = "new"
	case len(seg) == 1 && seg[0] == "search" && m == "GET":
		op = "search"
	case len(seg) == 1 && seg[0] == "board" && m == "GET":
		op = "board"
	case len(seg) == 1 && seg[0] == "check" && m == "GET":
		op = "check"
	case len(seg) == 1 && seg[0] == "index" && m == "POST":
		op = "index"
	case len(seg) == 2 && seg[0] == "records" && m == "GET":
		op, id = "show", seg[1]
	case len(seg) == 2 && seg[0] == "records" && m == "PATCH":
		op, id = "set", seg[1]
	case len(seg) == 2 && seg[0] == "find" && m == "GET":
		op, id = "find", seg[1]
	case len(seg) == 2 && seg[0] == "when" && m == "GET":
		op, id = "when", seg[1]
	case len(seg) == 3 && seg[0] == "records" && seg[2] == "close" && m == "POST":
		op, id = "close", seg[1]
	case len(seg) == 3 && seg[0] == "records" && seg[2] == "claim" && m == "POST":
		op, id = "claim", seg[1]
	case len(seg) == 3 && seg[0] == "records" && seg[2] == "claim" && m == "DELETE":
		op, id = "release", seg[1]
	default:
		return "", nil, false, nil
	}
	req := api.Request{}
	if m == "GET" {
		q := r.URL.Query()
		kinds := api.Params(op)
		for k, vs := range q {
			req[k] = queryValue(kinds[k], vs)
		}
	} else if m != "DELETE" || r.ContentLength > 0 {
		if m != "DELETE" && !strings.HasPrefix(strings.ToLower(r.Header.Get("Content-Type")), "application/json") && r.ContentLength != 0 {
			return "", nil, false, &api.Error{Code: "invalid_request", Message: "write requests need Content-Type: application/json"}
		}
		body, err := io.ReadAll(http.MaxBytesReader(nil, r.Body, 1<<20))
		if err != nil {
			return "", nil, false, &api.Error{Code: "invalid_request", Message: "cannot read the request body"}
		}
		if len(strings.TrimSpace(string(body))) > 0 {
			if err := json.Unmarshal(body, &req); err != nil {
				return "", nil, false, &api.Error{Code: "invalid_request", Message: "the body is not a JSON object"}
			}
		}
	}
	if id != "" {
		if got, ok := req["id"]; ok && got != id {
			return "", nil, false, &api.Error{Code: "invalid_request", Message: "id in the body differs from the id in the path", Details: api.O("parameter", "id")}
		}
		req["id"] = id
	}
	if needsIfMatch[op] {
		if _, ok := req["if_match"]; !ok {
			h := strings.Trim(strings.TrimSpace(r.Header.Get("If-Match")), `"`)
			if h == "" {
				return "", nil, false, &api.Error{Code: "refused", Message: "a write needs If-Match with the record's etag (read the record first)", Details: api.O("reason", "if_match_required")}
			}
			req["if_match"] = h
		}
	}
	return op, req, true, nil
}

func queryValue(kind string, vs []string) any {
	v := vs[len(vs)-1]
	switch kind {
	case "bool":
		return v == "true" || v == "1"
	case "int":
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
		return v
	case "strings":
		var out []string
		for _, x := range vs {
			for _, p := range strings.Split(x, ",") {
				if p = strings.TrimSpace(p); p != "" {
					out = append(out, p)
				}
			}
		}
		return out
	}
	return v
}
