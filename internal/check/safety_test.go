package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The secrets below are built from pieces so that no scanner reading this file
// finds a real-looking credential. They are the public documentation examples
// or obvious fakes.
var (
	awsKey  = "AKIA" + "IOSFODNN7EXAMPLE"
	ghToken = "ghp_" + strings.Repeat("a1B2c3D4e5", 4)
	pemHead = "-----BEGIN RSA " + "PRIVATE KEY-----"
	jwt     = "eyJ" + "hbGciOiJIUzI1NiJ9" + ".eyJ" + "zdWIiOiIxMjM0NTY3ODkwIn0" + ".dBjftJeZ4CVPmB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	random  = "q8Zk3VxT9mWn2LpR7sYb4HdJ6cFg1AeU5oXi"
)

func withLine(line string) map[string]func(string) string {
	return map[string]func(string) string{
		"AGENTS.md":    func(s string) string { return s + "\n" + line + "\n" },
		"docs/.layout": func(s string) string { return s + "extension.safety: 0.2.0\n" },
	}
}

func safetyRules(t *testing.T, edits map[string]func(string) string) []Finding {
	t.Helper()
	var out []Finding
	for _, f := range Run(repoWith(t, edits), Options{AsOf: asOf}).Findings {
		if strings.HasPrefix(f.Rule, "SF") {
			out = append(out, f)
		}
	}
	return out
}

func TestSafetyOffByDefault(t *testing.T) {
	e := withLine("key " + awsKey)
	delete(e, "docs/.layout")
	if got := safetyRules(t, e); len(got) != 0 {
		t.Fatalf("%+v", got)
	}
}

func TestSF1Patterns(t *testing.T) {
	for name, line := range map[string]string{
		"aws":   "Use " + awsKey + " for the bucket.",
		"gh":    "token " + ghToken,
		"pem":   pemHead,
		"jwt":   "Bearer " + jwt,
		"slack": "xox" + "b-1234567890-abcdefghij",
		"pass":  "password: Tr0ub4dor&3xyz",
		"env":   "API_KEY=9f8e7d6c5b4a3f2e1d",
	} {
		t.Run(name, func(t *testing.T) {
			got := safetyRules(t, withLine(line))
			if len(got) != 1 || got[0].Rule != "SF1" || got[0].Severity != "E" || got[0].File != "AGENTS.md" || got[0].Line == 0 {
				t.Fatalf("%+v", got)
			}
			// the finding never prints the matched text
			for _, secret := range []string{awsKey, ghToken, jwt, "Tr0ub4dor", "9f8e7d6c5b4a3f2e1d", "1234567890-abcdefghij"} {
				if strings.Contains(got[0].Message, secret) {
					t.Fatalf("message leaks the secret: %s", got[0].Message)
				}
			}
		})
	}
}

func TestSF1IgnoresPlaceholders(t *testing.T) {
	for _, line := range []string{
		"password: <your password>",
		"password: changeme",
		"password = ${DB_PASSWORD}",
		"api_key: your_api_key_here",
		"The password must be at least 12 characters.",
		"secret: example-value-123",
		"token: xxxxxxxxxxxx",
		"commit 3f786850e387550fdab836ed7e6dc881de23001b fixed it",
		"See docs/work/261022-e5f6-document-the-exit-codes.md for the plan.",
		"id: 261022-e5f6 and parent 261015-c4d5",
	} {
		if got := safetyRules(t, withLine(line)); len(got) != 0 {
			t.Errorf("%q: %+v", line, got)
		}
	}
}

func TestSF3Entropy(t *testing.T) {
	got := safetyRules(t, withLine("The value is "+random+" in the config."))
	if len(got) != 1 || got[0].Rule != "SF3" || got[0].Severity != "W" {
		t.Fatalf("%+v", got)
	}
	if strings.Contains(got[0].Message, random) {
		t.Fatal("leaks")
	}
}

func TestSF2DenyList(t *testing.T) {
	e := withLine("We call the customer Zorblax Industries about it.")
	e["docs/.akm-deny"] = func(string) string { return "# private names\nZorblax\nre:(?i)project +nightingale\n" }
	got := safetyRules(t, e)
	if len(got) != 1 || got[0].Rule != "SF2" || strings.Contains(got[0].Message, "Zorblax") {
		t.Fatalf("%+v", got)
	}
	e = withLine("About Project  Nightingale.")
	e["docs/.akm-deny"] = func(string) string { return "re:(?i)project +nightingale\n" }
	if got := safetyRules(t, e); len(got) != 1 || got[0].Rule != "SF2" {
		t.Fatalf("regexp: %+v", got)
	}
	e["docs/.akm-deny"] = func(string) string { return "re:([\n" }
	if got := safetyRules(t, e); len(got) != 1 || got[0].File != "docs/.akm-deny" {
		t.Fatalf("bad regexp: %+v", got)
	}
}

func TestSF4AllowMarkers(t *testing.T) {
	got := safetyRules(t, withLine("The docs example is "+awsKey+" <!-- akm-allow: SF1 public documentation example -->"))
	if len(got) != 1 || got[0].Rule != "SF4" || got[0].Severity != "W" || !strings.Contains(got[0].Message, "public documentation example") {
		t.Fatalf("%+v", got)
	}
	// no reason: the marker is an error and suppresses nothing
	got = safetyRules(t, withLine("The docs example is "+awsKey+" <!-- akm-allow: SF1 -->"))
	rules := map[string]string{}
	for _, f := range got {
		rules[f.Rule] = f.Severity
	}
	if rules["SF4"] != "E" || rules["SF1"] != "E" {
		t.Fatalf("%+v", got)
	}
}

func TestSafetyScansRecordsAndRoadmap(t *testing.T) {
	e := map[string]func(string) string{
		"docs/.layout":    func(s string) string { return s + "extension.safety: 0.2.0\n" },
		"docs/ROADMAP.md": func(s string) string { return s + "\nnote " + awsKey + "\n" },
		a2b3:              func(s string) string { return s + "\nsaved " + pemHead + "\n" },
	}
	got := safetyRules(t, e)
	files := map[string]bool{}
	for _, f := range got {
		files[f.File] = true
	}
	if !files["docs/ROADMAP.md"] || !files[a2b3] || len(got) != 2 {
		t.Fatalf("%+v", got)
	}
}

// TestSafetyCorpus measures false positives on ordinary Markdown: set
// AKM_CORPUS to a directory of notes. It prints the findings and never fails.
func TestSafetyCorpus(t *testing.T) {
	dir := os.Getenv("AKM_CORPUS")
	if dir == "" {
		t.Skip("AKM_CORPUS not set")
	}
	c := &checker{repo: repoWith(t, nil)}
	files, lines := 0, 0
	filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".md") {
			return nil
		}
		b, _ := os.ReadFile(p)
		files++
		for i, l := range strings.Split(string(b), "\n") {
			lines++
			before := len(c.out)
			c.safetyLine(filepath.Base(p), i+1, l, nil)
			for _, f := range c.out[before:] {
				t.Logf("%s:%d %s %s", f.File, f.Line, f.Rule, f.Message)
			}
		}
		return nil
	})
	t.Logf("%d files, %d lines, %d findings", files, lines, len(c.out))
}
