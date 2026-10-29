package server

import (
	"encoding/json"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

// contentSecurityPolicy is sent with every embedded asset. No inline script
// or style; fonts come from the same origin only.
const contentSecurityPolicy = "default-src 'none'; script-src 'self'; style-src 'self'; font-src 'self'; connect-src 'self'; img-src 'self' data:; base-uri 'none'; form-action 'none'; frame-ancestors 'none'"

// assetTypes lists what may be served from each embedded prefix.
var assetTypes = map[string]map[string]string{
	"/js/":    {".js": "text/javascript; charset=utf-8"},
	"/fonts/": {".woff2": "font/woff2", ".txt": "text/plain; charset=utf-8", "NOTICE": "text/plain; charset=utf-8"},
}

// serveAsset serves the embedded files that hold no data. It reports whether
// the path was an asset path (so the caller stops), including a refusal.
func (s *Server) serveAsset(w http.ResponseWriter, r *http.Request) bool {
	p := r.URL.Path
	for prefix, types := range assetTypes {
		if !strings.HasPrefix(p, prefix) {
			continue
		}
		name := strings.TrimPrefix(p, prefix)
		// one flat directory per prefix: no separators, no dots-only names,
		// no encoded or backslash forms (r.URL.Path is already decoded)
		if name == "" || strings.ContainsAny(name, "/\\\x00") || strings.HasPrefix(name, ".") || path.Clean(name) != name {
			http.NotFound(w, r)
			return true
		}
		ctype, ok := types[path.Ext(name)]
		if !ok {
			ctype, ok = types[name]
		}
		if !ok {
			http.NotFound(w, r)
			return true
		}
		s.asset(w, "web"+prefix+name, ctype)
		return true
	}
	return false
}

// uiConfig is the body of GET /ui/config. It never carries the token.
func (s *Server) uiConfig(w http.ResponseWriter) {
	root, err := filepath.Abs(s.Root)
	if err != nil {
		root = s.Root
	}
	b, _ := json.Marshal(map[string]string{
		"root":     root,
		"identity": osUser(),
	})
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(append(b, '\n'))
}

func osUser() string {
	if u, err := user.Current(); err == nil && u.Username != "" {
		return u.Username
	}
	for _, k := range []string{"USER", "USERNAME", "LOGNAME"} {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}
