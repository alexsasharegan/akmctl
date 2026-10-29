package server

import (
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The embedded client: every file the page loads is listed here, so a new
// asset is a deliberate change that is served and tested.
var clientModules = []string{"adrs.js", "anim.js", "board.js", "boot.js", "capture.js", "check.js", "detail.js", "flow.js", "live.js", "logic.js", "moves.js", "review.js", "rows.js", "store.js", "tables.js", "ui.js", "util.js", "wizard.js"}

func TestClientAssetList(t *testing.T) {
	entries, err := fs.ReadDir(web, "web/js")
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, e := range entries {
		got = append(got, e.Name())
	}
	sort.Strings(got)
	want := append([]string(nil), clientModules...)
	sort.Strings(want)
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("web/js holds %v; the test lists %v: update clientModules and the page", got, want)
	}
}

func TestClientModulesAreServed(t *testing.T) {
	ts, _ := setup(t)
	for _, name := range append([]string{"/app.js", "/app.css", "/"}, prefixed("/js/", clientModules)...) {
		res, b := get(t, ts, name, nil)
		if res.StatusCode != 200 || len(b) == 0 {
			t.Fatalf("%s: %d", name, res.StatusCode)
		}
		want := "text/javascript"
		switch {
		case name == "/":
			want = "text/html"
		case name == "/app.css":
			want = "text/css"
		}
		if !strings.HasPrefix(res.Header.Get("Content-Type"), want) {
			t.Fatalf("%s: type %q", name, res.Header.Get("Content-Type"))
		}
		csp := res.Header.Get("Content-Security-Policy")
		if !strings.Contains(csp, "script-src 'self'") || strings.Contains(csp, "unsafe") {
			t.Fatalf("%s: csp %q", name, csp)
		}
		if strings.Contains(string(b), token) {
			t.Fatalf("%s holds the token", name)
		}
	}
}

func prefixed(p string, names []string) []string {
	var out []string
	for _, n := range names {
		out = append(out, p+n)
	}
	return out
}

func readWeb(t *testing.T, name string) string {
	t.Helper()
	b, err := web.ReadFile("web/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestPageHoldsNoInlineCode(t *testing.T) {
	html := readWeb(t, "index.html")
	for _, re := range []string{`(?s)<script(?:\s[^>]*)?>[^<\s]`, `<style`, `\sstyle\s*=`, `\son[a-z]+\s*=`, `javascript:`} {
		if regexp.MustCompile(re).MatchString(html) {
			t.Fatalf("index.html matches %q", re)
		}
	}
	// every script and stylesheet it names is served
	ts, _ := setup(t)
	refs := regexp.MustCompile(`(?:src|href)="(/[^"]*)"`).FindAllStringSubmatch(html, -1)
	if len(refs) < 3 {
		t.Fatalf("page names %d assets", len(refs))
	}
	for _, m := range refs {
		if res, _ := get(t, ts, m[1], nil); res.StatusCode != 200 {
			t.Fatalf("%s: %d", m[1], res.StatusCode)
		}
	}
	if !strings.Contains(html, `type="module" src="/app.js"`) || !strings.Contains(html, `src="/js/boot.js"`) {
		t.Fatal("page lacks the module entry or the theme boot script")
	}
}

func TestFontsInCSSAreServed(t *testing.T) {
	css := readWeb(t, "app.css")
	ts, _ := setup(t)
	urls := regexp.MustCompile(`url\("([^"]+)"\)`).FindAllStringSubmatch(css, -1)
	if len(urls) != 4 {
		t.Fatalf("app.css names %d font files, want 4", len(urls))
	}
	for _, m := range urls {
		if !strings.HasPrefix(m[1], "/fonts/") {
			t.Fatalf("font url %q is not under /fonts/", m[1])
		}
		if res, _ := get(t, ts, m[1], nil); res.StatusCode != 200 || res.Header.Get("Content-Type") != "font/woff2" {
			t.Fatalf("%s: %d %q", m[1], res.StatusCode, res.Header.Get("Content-Type"))
		}
	}
	// local() comes first in every @font-face, so an installed Monaspace wins
	for _, face := range regexp.MustCompile(`@font-face\{[^}]*\}`).FindAllString(css, -1) {
		if strings.Index(face, "local(") < 0 || strings.Index(face, "local(") > strings.Index(face, "url(") {
			t.Fatalf("local() is not first in %q", face)
		}
	}
	if strings.Contains(css, "@import") || regexp.MustCompile(`url\(["']?(https?:)?//`).MatchString(css) {
		t.Fatal("app.css reaches out to another origin")
	}
}

func TestCSSMotionOnlyWhenAllowed(t *testing.T) {
	css := readWeb(t, "app.css")
	const open = "@media (prefers-reduced-motion:no-preference){"
	if !strings.Contains(css, open) {
		t.Fatal("no reduced-motion block")
	}
	// every such block is taken out; what is left must hold no motion
	rest := css
	for {
		i := strings.Index(rest, open)
		if i < 0 {
			break
		}
		depth, j := 1, i+len(open)
		for ; j < len(rest) && depth > 0; j++ {
			switch rest[j] {
			case '{':
				depth++
			case '}':
				depth--
			}
		}
		rest = rest[:i] + rest[j:]
	}
	if regexp.MustCompile(`(animation|transition)\s*:|(^|[;{\s])scroll-behavior\s*:`).MatchString(rest) {
		t.Fatal("motion is declared outside the prefers-reduced-motion:no-preference block")
	}
	if strings.Contains(css, "style=") {
		t.Fatal("css holds style=")
	}
}

func clientSources(t *testing.T) map[string]string {
	t.Helper()
	out := map[string]string{"app.js": readWeb(t, "app.js")}
	for _, n := range clientModules {
		out["js/"+n] = readWeb(t, "js/"+n)
	}
	return out
}

func TestClientSourceRules(t *testing.T) {
	banned := []string{"innerHTML", "outerHTML", "insertAdjacentHTML", "document.write", "createContextualFragment", "srcdoc", "eval(", "new Function", "cssText", `setAttribute("style"`, "setAttribute('style'", ".style =", "javascript:", "XMLHttpRequest", "importScripts", "WebSocket", "EventSource"}
	for name, src := range clientSources(t) {
		for _, b := range banned {
			if strings.Contains(src, b) {
				t.Errorf("%s uses %q", name, b)
			}
		}
		if name != "js/store.js" && strings.Contains(src, "fetch(") {
			t.Errorf("%s calls fetch; every request goes through store.js", name)
		}
		if strings.Contains(src, "http://") || strings.Contains(src, "https://") {
			t.Errorf("%s names an absolute URL", name)
		}
	}
	store := clientSources(t)["js/store.js"]
	if strings.Contains(store, "localStorage") || !strings.Contains(store, `"Bearer "`) {
		t.Error("store.js must keep the token in sessionStorage and send it as a Bearer header")
	}
	if !strings.Contains(store, "history.replaceState") && !strings.Contains(store, "hist.replaceState") {
		t.Error("store.js must strip the token from the address bar")
	}
}

func TestClientImportsResolve(t *testing.T) {
	re := regexp.MustCompile(`(?m)^\s*(?:import|export)\b[^;]*?from\s+"([^"]+)"`)
	for name, src := range clientSources(t) {
		for _, m := range re.FindAllStringSubmatch(src, -1) {
			target := path.Join(path.Dir(name), m[1])
			if !strings.HasPrefix(m[1], "./") {
				t.Errorf("%s imports %q: only relative imports are served", name, m[1])
			}
			if _, err := web.ReadFile("web/" + target); err != nil {
				t.Errorf("%s imports %q, which is not embedded", name, m[1])
			}
		}
	}
}

// The pure modules run under Node when it is installed (CI has it). They are
// copied next to a package.json that marks them as modules.
func TestPureModulesUnderNode(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not installed")
	}
	dir := t.TempDir()
	for _, n := range []string{"util.js", "logic.js", "store.js", "live.js", "moves.js", "review.js", "capture.js"} {
		if err := os.WriteFile(filepath.Join(dir, n), []byte(readWeb(t, "js/"+n)), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"type":"module"}`), 0o644)
	test, err := os.ReadFile("testdata/logic.test.mjs")
	if err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(dir, "logic.test.mjs"), test, 0o644)
	cmd := exec.Command(node, "logic.test.mjs")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("node: %v\n%s", err, out)
	}
}
