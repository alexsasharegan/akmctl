package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// render prints a diff in the unified form, for comparing in tests.
func render(d Obj) string {
	b, _ := Marshal(d)
	var m struct {
		Path    string `json:"path"`
		NewPath string `json:"new_path"`
		Hunks   []struct {
			Header string `json:"header"`
			Lines  []struct {
				Op   string `json:"op"`
				Text string `json:"text"`
			} `json:"lines"`
		} `json:"hunks"`
	}
	json.Unmarshal(b, &m)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s > %s\n", m.Path, m.NewPath)
	for _, h := range m.Hunks {
		sb.WriteString(h.Header + "\n")
		for _, l := range h.Lines {
			sb.WriteString(map[string]string{"ctx": " ", "add": "+", "rem": "-"}[l.Op] + l.Text + "\n")
		}
	}
	return sb.String()
}

func nums(n int) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&sb, "l%d\n", i)
	}
	return sb.String()
}

func TestDiffNoChange(t *testing.T) {
	d := Diff("a.md", "a.md", "x\ny\n", "x\ny\n")
	if got := render(d); got != "a.md > a.md\n" {
		t.Fatal(got)
	}
	hs, _ := d.Get("hunks")
	if hs == nil || len(hs.([]Obj)) != 0 {
		t.Fatalf("hunks must be an empty array, got %#v", hs)
	}
	b, _ := Marshal(d)
	if !strings.Contains(string(b), `"hunks":[]`) {
		t.Fatal(string(b))
	}
}

func TestDiffRename(t *testing.T) {
	got := render(Diff("docs/work/a.md", "docs/work/done/a.md", "status: active\nbody\n", "status: done\nbody\n"))
	want := "docs/work/a.md > docs/work/done/a.md\n@@ -1,2 +1,2 @@\n-status: active\n+status: done\n body\n"
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
}

func TestDiffContextAndHunks(t *testing.T) {
	old := nums(30)
	// two changes 20 lines apart: two hunks with three lines of context
	nw := strings.Replace(strings.Replace(old, "l5\n", "L5\n", 1), "l25\n", "L25\n", 1)
	got := render(Diff("f", "f", old, nw))
	want := "f > f\n@@ -2,7 +2,7 @@\n l2\n l3\n l4\n-l5\n+L5\n l6\n l7\n l8\n@@ -22,7 +22,7 @@\n l22\n l23\n l24\n-l25\n+L25\n l26\n l27\n l28\n"
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
	// changes six lines apart share context and make one hunk; seven apart do not
	near := strings.Replace(strings.Replace(old, "l5\n", "L5\n", 1), "l13\n", "L13\n", 1)
	if n := strings.Count(render(Diff("f", "f", old, near)), "@@ -"); n != 2 {
		t.Fatalf("7 unchanged lines between changes: %d hunks", n)
	}
	near = strings.Replace(strings.Replace(old, "l5\n", "L5\n", 1), "l11\n", "L11\n", 1)
	if n := strings.Count(render(Diff("f", "f", old, near)), "@@ -"); n != 1 {
		t.Fatalf("5 unchanged lines between changes: %d hunks", n)
	}
}

func TestDiffInsertAndDelete(t *testing.T) {
	got := render(Diff("f", "f", "a\nb\nc\n", "a\nb\nX\nY\nc\n"))
	if want := "f > f\n@@ -1,3 +1,5 @@\n a\n b\n+X\n+Y\n c\n"; got != want {
		t.Fatalf("got\n%s", got)
	}
	got = render(Diff("f", "f", "a\nb\nX\nY\nc\n", "a\nb\nc\n"))
	if want := "f > f\n@@ -1,5 +1,3 @@\n a\n b\n-X\n-Y\n c\n"; got != want {
		t.Fatalf("got\n%s", got)
	}
	// from and to an empty file
	if got := render(Diff("f", "f", "", "a\n")); got != "f > f\n@@ -0,0 +1 @@\n+a\n" {
		t.Fatal(got)
	}
	if got := render(Diff("f", "f", "a\n", "")); got != "f > f\n@@ -1 +0,0 @@\n-a\n" {
		t.Fatal(got)
	}
}

func TestDiffTrailingNewline(t *testing.T) {
	// no phantom empty line from the final newline
	if got := render(Diff("f", "f", "a\nb\n", "a\nb\nc\n")); got != "f > f\n@@ -1,2 +1,3 @@\n a\n b\n+c\n" {
		t.Fatal(got)
	}
	// a missing final newline is a change of the last line
	got := render(Diff("f", "f", "a\nb", "a\nb\n"))
	if want := "f > f\n@@ -1,2 +1,2 @@\n a\n-b\n+b\n"; got != want {
		t.Fatalf("got\n%s", got)
	}
	if got := render(Diff("f", "f", "a\nb\n", "a\nb")); got != "f > f\n@@ -1,2 +1,2 @@\n a\n-b\n+b\n" {
		t.Fatal(got)
	}
	if got := render(Diff("f", "f", "a\nb", "a\nb")); got != "f > f\n" {
		t.Fatal(got)
	}
}

func TestDiffCRLF(t *testing.T) {
	// a CR stays in the text of its line, so a line-ending change is visible
	old, nw := "a\r\nb\r\nc\r\n", "a\r\nB\r\nc\r\n"
	d := Diff("f", "f", old, nw)
	b, _ := Marshal(d)
	if !strings.Contains(string(b), `"text":"a\r"`) || !strings.Contains(string(b), `"text":"B\r"`) {
		t.Fatal(string(b))
	}
	if got := render(Diff("f", "f", "a\r\nb\r\n", "a\nb\n")); !strings.Contains(got, "-a\r\n") || !strings.Contains(got, "+a\n") {
		t.Fatalf("line-ending change not shown:\n%q", got)
	}
	if hs, _ := Diff("f", "f", old, old).Get("hunks"); len(hs.([]Obj)) != 0 {
		t.Fatal("identical CRLF text must have no hunks")
	}
}

func TestDiffHugeChangeFallsBack(t *testing.T) {
	var a, b strings.Builder
	for i := 0; i < 3000; i++ {
		fmt.Fprintf(&a, "a%d\n", i)
		fmt.Fprintf(&b, "b%d\n", i)
	}
	got := render(Diff("f", "f", a.String(), b.String()))
	if !strings.Contains(got, "@@ -1,3000 +1,3000 @@") || strings.Count(got, "\n-a") != 3000 || strings.Count(got, "\n+b") != 3000 {
		t.Fatal(got[:200])
	}
}
