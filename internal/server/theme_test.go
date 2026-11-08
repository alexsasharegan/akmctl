package server

import (
	"encoding/json"
	"io/fs"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The colour themes (work item 261108-hp9k): web/themes/themes.json lists them,
// web/themes/<id>.css holds one token block per variant, and these tests keep a
// new theme from shipping with a token missing.

type themeVariant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Mode string `json:"mode"`
}

type themeEntry struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Variants []themeVariant `json:"variants"`
}

type themeRegistry struct {
	Default string       `json:"default"`
	Themes  []themeEntry `json:"themes"`
}

type cssBlock struct {
	selector string
	tokens   map[string]bool
	scheme   string
}

var (
	cssComment = regexp.MustCompile(`(?s)/\*.*?\*/`)
	cssRule    = regexp.MustCompile(`([^{}]+)\{([^{}]*)\}`)
	cssToken   = regexp.MustCompile(`(--[a-z0-9-]+)\s*:`)
	cssScheme  = regexp.MustCompile(`color-scheme\s*:\s*(light|dark)`)
)

func loadRegistry(t *testing.T) themeRegistry {
	t.Helper()
	var reg themeRegistry
	if err := json.Unmarshal([]byte(readWeb(t, "themes/themes.json")), &reg); err != nil {
		t.Fatal(err)
	}
	return reg
}

func parseBlocks(css string) []cssBlock {
	var out []cssBlock
	for _, m := range cssRule.FindAllStringSubmatch(cssComment.ReplaceAllString(css, ""), -1) {
		b := cssBlock{selector: strings.Join(strings.Fields(m[1]), " "), tokens: map[string]bool{}}
		for _, d := range cssToken.FindAllStringSubmatch(m[2], -1) {
			b.tokens[d[1]] = true
		}
		if s := cssScheme.FindStringSubmatch(m[2]); s != nil {
			b.scheme = s[1]
		}
		out = append(out, b)
	}
	return out
}

// selectorOf is the selector a variant's block must carry: the first variant of
// a mode is the theme's default for it and has no data-variant.
func selectorOf(th themeEntry, i int) string {
	v := th.Variants[i]
	sel := `:root[data-theme="` + th.ID + `"][data-mode="` + v.Mode + `"]`
	for j := 0; j < i; j++ {
		if th.Variants[j].Mode == v.Mode {
			return sel + `[data-variant="` + v.ID + `"]`
		}
	}
	return sel
}

func tokenNames(m map[string]bool) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func TestThemeRegistryAndTokens(t *testing.T) {
	reg := loadRegistry(t)
	if len(reg.Themes) == 0 {
		t.Fatal("themes.json lists no theme")
	}
	// the files in web/themes are the themes in themes.json, one for one
	entries, err := fs.ReadDir(web, "web/themes")
	if err != nil {
		t.Fatal(err)
	}
	files := map[string]bool{}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".css") {
			files[strings.TrimSuffix(e.Name(), ".css")] = true
		}
	}
	var ref map[string]bool // the tokens every variant must define: the default theme's light block
	seen := map[string]bool{}
	defaultFound := false
	for _, th := range reg.Themes {
		if !regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`).MatchString(th.ID) || th.Name == "" || seen[th.ID] {
			t.Fatalf("theme %q: bad, unnamed or repeated id", th.ID)
		}
		seen[th.ID] = true
		if th.ID == reg.Default {
			defaultFound = true
		}
		if !files[th.ID] {
			t.Fatalf("themes.json lists %q but web/themes/%s.css is missing", th.ID, th.ID)
		}
		delete(files, th.ID)
		modes := map[string]int{}
		ids := map[string]bool{}
		for i, v := range th.Variants {
			if (v.Mode != "light" && v.Mode != "dark") || v.Name == "" || !regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`).MatchString(v.ID) || ids[v.Mode+"/"+v.ID] {
				t.Fatalf("theme %s variant %d (%+v): bad mode, name or id", th.ID, i, v)
			}
			ids[v.Mode+"/"+v.ID] = true
			modes[v.Mode]++
		}
		if modes["light"] == 0 || modes["dark"] == 0 {
			t.Fatalf("theme %s needs a light and a dark variant, has %v", th.ID, modes)
		}
		// each variant has one block with the right selector, scheme and every token
		blocks := parseBlocks(readWeb(t, "themes/"+th.ID+".css"))
		if len(blocks) != len(th.Variants) {
			t.Fatalf("theme %s: %d variants in themes.json but %d blocks in %s.css", th.ID, len(th.Variants), len(blocks), th.ID)
		}
		for i, v := range th.Variants {
			want := selectorOf(th, i)
			var got *cssBlock
			for j := range blocks {
				// the default theme's default light block also answers for a bare :root (the page before any theme is chosen)
				if blocks[j].selector == want || blocks[j].selector == ":root, "+want {
					got = &blocks[j]
				}
			}
			if got == nil {
				t.Fatalf("theme %s variant %s: no block with selector %s", th.ID, v.ID, want)
			}
			if got.scheme != v.Mode {
				t.Fatalf("theme %s variant %s: color-scheme is %q, want %q", th.ID, v.ID, got.scheme, v.Mode)
			}
			if ref == nil {
				ref = got.tokens
				if len(ref) < 30 {
					t.Fatalf("the reference block defines only %d tokens", len(ref))
				}
			}
			var missing, extra []string
			for k := range ref {
				if !got.tokens[k] {
					missing = append(missing, k)
				}
			}
			for k := range got.tokens {
				if !ref[k] {
					extra = append(extra, k)
				}
			}
			sort.Strings(missing)
			sort.Strings(extra)
			if len(missing) > 0 || len(extra) > 0 {
				t.Fatalf("theme %s variant %s: missing tokens %v, extra tokens %v", th.ID, v.ID, missing, extra)
			}
		}
	}
	if !defaultFound {
		t.Fatalf("default theme %q is not in themes.json", reg.Default)
	}
	if len(files) > 0 {
		t.Fatalf("web/themes holds CSS that themes.json does not list: %v", files)
	}
}

// app.css holds structure and non-colour tokens only: no colour literal, and
// every var(--x) it uses is defined by app.css or by the themes.
func TestAppCSSHasNoColourOfItsOwn(t *testing.T) {
	css := cssComment.ReplaceAllString(readWeb(t, "app.css"), "")
	for _, re := range []string{`#[0-9a-fA-F]{3,8}\b`, `\b(?:rgb|rgba|hsl|hsla|hwb|lab|lch|oklab|oklch)\((?:[^v)]|v[^a])`, `prefers-color-scheme`, `data-theme`} {
		if m := regexp.MustCompile(re).FindString(css); m != "" {
			t.Fatalf("app.css has %q (%s): colours belong in web/themes", m, re)
		}
	}
	defined := map[string]bool{}
	for _, d := range cssToken.FindAllStringSubmatch(css, -1) {
		defined[d[1]] = true
	}
	reg := loadRegistry(t)
	for _, b := range parseBlocks(readWeb(t, "themes/"+reg.Default+".css")) {
		for k := range b.tokens {
			defined[k] = true
		}
	}
	defined["--toast-ms"] = true // set from script on each toast (ui.js)
	var undefined []string
	for _, m := range regexp.MustCompile(`var\((--[a-z0-9-]+)`).FindAllStringSubmatch(css, -1) {
		if !defined[m[1]] {
			undefined = append(undefined, m[1])
		}
	}
	if len(undefined) > 0 {
		sort.Strings(undefined)
		t.Fatalf("app.css uses undefined tokens %v", undefined)
	}
}

func TestThemesAreServed(t *testing.T) {
	ts, _ := setup(t)
	reg := loadRegistry(t)
	res, b := get(t, ts, "/themes.css", nil)
	if res.StatusCode != 200 || !strings.HasPrefix(res.Header.Get("Content-Type"), "text/css") {
		t.Fatalf("/themes.css: %d %q", res.StatusCode, res.Header.Get("Content-Type"))
	}
	if strings.Contains(string(b), token) {
		t.Fatal("/themes.css holds the token")
	}
	got := parseBlocks(string(b))
	for _, th := range reg.Themes {
		for i := range th.Variants {
			want := selectorOf(th, i)
			ok := false
			for _, g := range got {
				ok = ok || g.selector == want || g.selector == ":root, "+want
			}
			if !ok {
				t.Fatalf("/themes.css lacks %s", want)
			}
		}
	}
	res, b = get(t, ts, "/themes/themes.json", nil)
	if res.StatusCode != 200 || !strings.HasPrefix(res.Header.Get("Content-Type"), "application/json") || !json.Valid(b) {
		t.Fatalf("/themes/themes.json: %d %q", res.StatusCode, res.Header.Get("Content-Type"))
	}
	// the stylesheet is part of the page
	if !strings.Contains(readWeb(t, "index.html"), `href="/themes.css"`) {
		t.Fatal("index.html does not load /themes.css")
	}
}
