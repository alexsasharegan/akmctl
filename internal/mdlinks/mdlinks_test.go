package mdlinks

import (
	"path"
	"testing"
)

func rewriteTo(from, to string) func(string) (string, bool) {
	return func(t string) (string, bool) {
		if path.Clean(t) == from {
			return to, true
		}
		return "", false
	}
}

func TestScanSkips(t *testing.T) {
	text := "---\nsee: [x](front.md)\n---\n" +
		"[a](a.md) `[b](b.md)` [c](https://x/c.md) [d](#top) [e](/abs.md)\n" +
		"```\n[f](f.md)\n```\n" +
		"~~~\n```\n[g](g.md)\n~~~\n" +
		"[h](<h one.md#sec> \"t\") [i](i%20two.md?x=1)\r\n[j](./j.md 'title')\n"
	var got []string
	for _, l := range Scan(text) {
		got = append(got, l.Target)
	}
	want := []string{"a.md", "h one.md", "i two.md", "./j.md"}
	if len(got) != len(want) {
		t.Fatalf("targets = %q, want %q", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("targets = %q, want %q", got, want)
		}
	}
}

func TestRewriteKeepsTheRest(t *testing.T) {
	text := "[a](old.md#sec) [b](<old.md> \"t\") [c](./old.md) `[d](old.md)` [e](other.md)\r\n```\n[f](old.md)\n```\n"
	got, n := Rewrite(text, rewriteTo("old.md", "done/old.md"))
	want := "[a](done/old.md#sec) [b](<done/old.md> \"t\") [c](./done/old.md) `[d](old.md)` [e](other.md)\r\n```\n[f](old.md)\n```\n"
	if n != 3 || got != want {
		t.Errorf("n = %d\ngot  %q\nwant %q", n, got, want)
	}
	if got, n := Rewrite("[a](x%20y.md)", rewriteTo("x y.md", "d/x y.md")); n != 1 || got != "[a](d/x%20y.md)" {
		t.Errorf("escaped: %q %d", got, n)
	}
	if got, n := Rewrite("[a](x.md)", rewriteTo("x.md", "d/a b.md")); n != 1 || got != "[a](d/a%20b.md)" {
		t.Errorf("space: %q %d", got, n)
	}
}

func TestRel(t *testing.T) {
	for _, c := range []struct{ from, to, want string }{
		{"docs/work", "docs/work/done/x.md", "done/x.md"},
		{"docs/work/done", "docs/adr/a.md", "../../adr/a.md"},
		{"docs/work/done", "docs/work/x.md", "../x.md"},
		{"docs/work/done", "docs/work/done/x.md", "x.md"},
		{".", "docs/work/x.md", "docs/work/x.md"},
		{"docs", "AGENTS.md", "../AGENTS.md"},
	} {
		if got := Rel(c.from, c.to); got != c.want {
			t.Errorf("Rel(%q, %q) = %q, want %q", c.from, c.to, got, c.want)
		}
	}
}
