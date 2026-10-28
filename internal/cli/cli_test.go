package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

const repo = "../../testdata/api/repos/basic"

func run(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, errb bytes.Buffer
	code := Run(append(args, "--repo", repo), &out, &errb, ".")
	return code, out.String(), errb.String()
}

func TestJSONOutputIsOneLine(t *testing.T) {
	code, out, _ := run(t, "ls", "--format", "json", "--status", "ready", "--fields", "title")
	if code != 0 || strings.Count(out, "\n") != 1 {
		t.Fatalf("code %d, output %q", code, out)
	}
	var got struct {
		OK    bool
		Count int
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil || !got.OK || got.Count != 1 {
		t.Fatalf("%v %+v", err, got)
	}
}

func TestFlagsMapToParameters(t *testing.T) {
	code, out, _ := run(t, "ls", "--format=json", "--count-only", "--no-next", "--type", "work", "--limit=2")
	if code != 0 {
		t.Fatal(out)
	}
	for _, want := range []string{`"count_only":true`, `"next":false`, `"limit":2`, `"total":6`} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
	}
}

func TestPositionals(t *testing.T) {
	_, out, _ := run(t, "search", "JSON", "output", "--format", "ids")
	if strings.TrimSpace(out) != "261020-a2b3" {
		t.Fatalf("search words: %q", out)
	}
	code, out, _ := run(t, "find", "261022-e5f6")
	if code != 0 || strings.TrimSpace(out) != "docs/work/261022-e5f6-document-the-exit-codes.md" {
		t.Fatalf("%d %q", code, out)
	}
}

func TestErrorsExitTwo(t *testing.T) {
	for _, args := range [][]string{{"show", "261099-zzzz"}, {"ls", "--status", "finished"}, {"ls", "--bogus"}, {"show"}, {"nonsense"}, {"ls", "extra"}} {
		code, _, _ := run(t, args...)
		if code != 2 {
			t.Errorf("%v exits %d", args, code)
		}
	}
	code, out, _ := run(t, "show", "261099-zzzz", "--format", "json")
	if code != 2 || !strings.Contains(out, `"code":"not_found"`) {
		t.Fatalf("%d %s", code, out)
	}
}

func TestNoRepository(t *testing.T) {
	var out, errb bytes.Buffer
	code := Run([]string{"ls", "--format", "json"}, &out, &errb, t.TempDir())
	if code != 2 || !strings.Contains(out.String(), `"code":"unreadable"`) {
		t.Fatalf("%d %s", code, out.String())
	}
}
