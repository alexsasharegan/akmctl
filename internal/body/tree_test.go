package body

import (
	"bytes"
	"errors"
	"math/rand"
	"strings"
	"testing"
)

const sample = "# Title\n\nintro\n\n## Goal\n\nDo it.\n\nDone when:\n\n- [ ] one\n- [x] two\n\nafter\n- [ ] not in the run\n\n## Next\n\nStep.\n\n## Notes\n\nfree\n\n## Log\n\n- 2026-10-01: filed\n- 2026-10-02: claimed\n"

func mustSplit(t *testing.T, s string) *Tree {
	t.Helper()
	tr, err := Split([]byte(s))
	if err != nil {
		t.Fatal(err)
	}
	return tr
}

func TestSplitJoin(t *testing.T) {
	tr := mustSplit(t, sample)
	if string(tr.Join()) != sample {
		t.Fatal("Join differs")
	}
	if got := string(tr.Preamble()); got != "# Title\n\nintro\n\n" {
		t.Errorf("preamble %q", got)
	}
	var names []string
	for _, s := range tr.Sections() {
		names = append(names, s.Name)
	}
	if strings.Join(names, ",") != "Goal,Next,Notes,Log" {
		t.Errorf("sections %v", names)
	}
	if s, ok := tr.Section("LOG"); !ok || s.Line != 25 {
		t.Errorf("Log section = %+v %v", s, ok)
	}
	if got := string(tr.Text(mustSection(t, tr, "next"))); got != "\nStep.\n\n" {
		t.Errorf("next text %q", got)
	}
	// a heading in a fence is not a section; closing hashes are not part of the name
	tr = mustSplit(t, "## A ##\n```\n## no\n```\n## C#\n")
	if s := tr.Sections(); len(s) != 2 || s[0].Name != "A" || s[1].Name != "C#" {
		t.Errorf("sections %+v", s)
	}
	if _, err := Split([]byte("## A\n\n<div>\n")); !errors.Is(err, ErrUnsupported) {
		t.Errorf("Split of an HTML line = %v, want ErrUnsupported", err)
	}
}

func mustSection(t *testing.T, tr *Tree, name string) Section {
	t.Helper()
	s, ok := tr.Section(name)
	if !ok {
		t.Fatalf("no section %s", name)
	}
	return s
}

func TestChecklist(t *testing.T) {
	tr := mustSplit(t, sample)
	items, err := tr.Checklist("goal")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[0].Text != "one" || items[0].Done || !items[1].Done || items[1].Pos != 2 || items[0].Line != 11 {
		t.Fatalf("items = %+v", items)
	}
	if items, _ := tr.Checklist("next"); len(items) != 0 {
		t.Errorf("next has a checklist: %+v", items)
	}
	if _, err := tr.Checklist("missing"); !errors.Is(err, ErrNoSection) {
		t.Errorf("missing section: %v", err)
	}
}

func TestSetChecked(t *testing.T) {
	for _, eol := range []string{"\n", "\r\n"} {
		src := strings.ReplaceAll(sample, "\n", eol)
		tr := mustSplit(t, src)
		out, err := tr.SetChecked("goal", 1, true, "")
		if err != nil {
			t.Fatal(err)
		}
		if want := strings.Replace(src, "- [ ] one", "- [x] one", 1); string(out) != want {
			t.Errorf("tick (%q):\n%q", eol, out)
		}
		out, _ = tr.SetChecked("goal", 2, false, "")
		if want := strings.Replace(src, "- [x] two", "- [ ] two", 1); string(out) != want {
			t.Errorf("untick: %q", out)
		}
		out, _ = tr.SetChecked("goal", 2, true, "") // already ticked
		if string(out) != src {
			t.Errorf("a no-op changed the body")
		}
	}
	tr := mustSplit(t, sample)
	if _, err := tr.SetChecked("goal", 3, true, ""); !errors.Is(err, ErrNoItem) {
		t.Errorf("item 3 (outside the run): %v", err)
	}
	if _, err := tr.SetChecked("goal", 1, true, "0000"); !errors.Is(err, ErrStale) {
		t.Errorf("stale etag: %v", err)
	}
	if _, err := tr.SetChecked("goal", 1, true, tr.Etag(mustSection(t, tr, "goal"))); err != nil {
		t.Errorf("good etag: %v", err)
	}
}

func TestLogEntriesAndAppend(t *testing.T) {
	tr := mustSplit(t, sample)
	es, err := tr.LogEntries()
	if err != nil || len(es) != 2 || es[1].Date != "2026-10-02" || es[1].Text != "claimed" {
		t.Fatalf("entries %+v %v", es, err)
	}
	out, err := tr.AppendLog("2026-10-03", "built\nit")
	if err != nil || string(out) != sample+"- 2026-10-03: built it\n" {
		t.Errorf("append: %q %v", out, err)
	}
	// before a following section, after trailing blanks, in CRLF
	src := "## Log\r\n\r\n- 2026-10-01: a\r\n\r\n## Outcome\r\n\r\nx\r\n"
	out, _ = mustSplit(t, src).AppendLog("2026-10-02", "b")
	if want := "## Log\r\n\r\n- 2026-10-01: a\r\n- 2026-10-02: b\r\n\r\n## Outcome\r\n\r\nx\r\n"; string(out) != want {
		t.Errorf("crlf append: %q", out)
	}
	// no final newline; empty section; heading last with no ending
	for in, want := range map[string]string{
		"## Log\n- 2026-10-01: a": "## Log\n- 2026-10-01: a\n- 2026-10-02: b",
		"## Log\n\n## Next\nx\n":  "## Log\n\n- 2026-10-02: b\n\n## Next\nx\n",
		"## Log":                  "## Log\n\n- 2026-10-02: b",
		"## Log\n":                "## Log\n\n- 2026-10-02: b\n",
	} {
		out, _ := mustSplit(t, in).AppendLog("2026-10-02", "b")
		if string(out) != want {
			t.Errorf("append to %q = %q, want %q", in, out, want)
		}
	}
	if _, err := mustSplit(t, "## Goal\n").AppendLog("2026-10-02", "b"); !errors.Is(err, ErrNoSection) {
		t.Errorf("no Log: %v", err)
	}
}

func TestReplaceSection(t *testing.T) {
	tr := mustSplit(t, sample)
	out, err := tr.ReplaceSection("Next", "New step.\nSecond.", "")
	if err != nil {
		t.Fatal(err)
	}
	if want := strings.Replace(sample, "## Next\n\nStep.\n\n", "## Next\n\nNew step.\nSecond.\n\n", 1); string(out) != want {
		t.Errorf("replace:\n%q", out)
	}
	src := strings.ReplaceAll(sample, "\n", "\r\n")
	out, _ = mustSplit(t, src).ReplaceSection("next", "A\nB\n", "")
	if want := strings.Replace(src, "Step.", "A\r\nB", 1); string(out) != want {
		t.Errorf("crlf replace: %q", out)
	}
	for in, want := range map[string]string{
		"## Next\n\n## Log\n":    "## Next\n\nx\n\n## Log\n",
		"## Next\n":              "## Next\n\nx\n",
		"## Next":                "## Next\n\nx\n",
		"## Next\ny":             "## Next\nx",
		"## Next\n\ny\n\n## L\n": "## Next\n\nx\n\n## L\n",
	} {
		out, err := mustSplit(t, in).ReplaceSection("next", "x", "")
		if err != nil || string(out) != want {
			t.Errorf("replace in %q = %q (%v), want %q", in, out, err, want)
		}
	}
	if _, err := tr.ReplaceSection("log", "x", ""); !errors.Is(err, ErrAppendOnly) {
		t.Errorf("replace Log: %v", err)
	}
	if _, err := tr.ReplaceSection("next", "x", "0000"); !errors.Is(err, ErrStale) {
		t.Errorf("stale: %v", err)
	}
	if _, err := tr.ReplaceSection("nope", "x", ""); !errors.Is(err, ErrNoSection) {
		t.Errorf("no section: %v", err)
	}
}

// TestDeviations: a typed section that does not have the typed shape keeps the
// rest as opaque text, and the typed part still edits.
func TestDeviations(t *testing.T) {
	src := "## Goal\n\n- [ ] one\n- [ ] two\n\nprose after the list\n\n- [ ] later\n\n## Log\n\n- 2026-10-01: line\n  continued here\n- not a dated line\n- 2026-10-02: next\n"
	tr := mustSplit(t, src)
	items, _ := tr.Checklist("goal")
	if len(items) != 2 {
		t.Fatalf("run = %+v", items)
	}
	out, err := tr.SetChecked("goal", 2, true, "")
	if err != nil || string(out) != strings.Replace(src, "- [ ] two", "- [x] two", 1) {
		t.Errorf("tick in a Goal with prose after the list: %q %v", out, err)
	}
	es, _ := tr.LogEntries()
	if len(es) != 2 || es[0].Text != "line" || es[1].Text != "next" {
		t.Errorf("entries = %+v", es)
	}
	out, _ = tr.AppendLog("2026-10-03", "z")
	if string(out) != src+"- 2026-10-03: z\n" {
		t.Errorf("append after opaque lines: %q", out)
	}
}

// TestSplitJoinRecords: the lossless round trip holds for every record of this
// repository, and no record is refused.
func TestSplitJoinRecords(t *testing.T) {
	bodies := repoBodies(t)
	for p, b := range bodies {
		tr, err := Split(b)
		if err != nil {
			t.Errorf("%s: %v", p, err)
			continue
		}
		if !bytes.Equal(tr.Join(), b) {
			t.Errorf("%s: Join differs", p)
		}
		var pre = tr.Preamble()
		sum := len(pre)
		for _, s := range tr.Sections() {
			sum += s.End - s.Start
		}
		if sum != len(b) {
			t.Errorf("%s: preamble and sections cover %d of %d bytes", p, sum, len(b))
		}
	}
}

// TestTreeAgainstOracle: for every generated document the scanner accepts, the
// sections found are the level-2 headings the oracle finds, and each typed edit
// changes nothing outside the span it means to change.
func TestTreeAgainstOracle(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	checked, edited, refused := 0, 0, 0
	for i := 0; i < 20000; i++ {
		src := generate(r)
		tr, err := Split(src)
		if err != nil {
			continue
		}
		orc, oerr := goldmarkLocator{}.Locate(src)
		if oerr != nil {
			t.Fatal(oerr)
		}
		got, want := names(tr), names(build(src, orc))
		if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
			t.Fatalf("%q: sections %q, oracle %q", src, got, want)
		}
		checked++
		for _, name := range []string{"goal", "next", "log"} {
			s, ok := tr.Section(name)
			if !ok {
				continue
			}
			if items, err := tr.Checklist(name); err == nil && len(items) > 0 {
				k := r.Intn(len(items)) + 1
				out, err := tr.SetChecked(name, k, !items[k-1].Done, "")
				switch {
				case errors.Is(err, ErrUnsupported):
					if tr.refusedIn(s) == nil {
						t.Fatalf("%q: tick refused with no local construct in %s: %v", src, name, err)
					}
					refused++
				case err != nil:
					t.Fatalf("%q: %v", src, err)
				default:
					if d := diffCount(string(src), out); d != 1 {
						t.Fatalf("%q: ticking changed %d bytes", src, d)
					}
					edited++
				}
			}
			var out []byte
			var err error
			switch name {
			case "log":
				out, err = tr.AppendLog("2026-11-01", "e")
			default:
				out, err = tr.ReplaceSection(name, "replaced", "")
			}
			if errors.Is(err, ErrUnsupported) {
				if tr.refusedIn(s) == nil {
					t.Fatalf("%q: %s edit refused with no local construct in it: %v", src, name, err)
				}
				refused++
				continue
			}
			if err != nil {
				t.Fatalf("%q: %s: %v", src, name, err)
			}
			if !bytes.HasPrefix(out, src[:s.Start]) || !bytes.HasSuffix(out, src[s.End:]) {
				t.Fatalf("%q: editing %s changed text outside the section:\n%q", src, name, out)
			}
			if _, err := Split(out); err != nil {
				t.Fatalf("%q: edited body does not split: %v\n%q", src, err, out)
			}
			edited++
		}
	}
	t.Logf("%d documents split, %d edits checked, %d refused for a local construct in their section", checked, edited, refused)
	if refused < 100 {
		t.Errorf("only %d edits were refused; the local-construct path is barely exercised", refused)
	}
	if checked < 2000 || edited < 2000 {
		t.Errorf("too few documents or edits (%d, %d)", checked, edited)
	}
}

func names(tr *Tree) []string {
	var out []string
	for _, s := range tr.Sections() {
		out = append(out, s.Key())
	}
	return out
}

func diffCount(a string, b []byte) int {
	if len(a) != len(b) {
		return -1
	}
	n := 0
	for i := range b {
		if a[i] != b[i] {
			n++
		}
	}
	return n
}

// FuzzTree: whatever Split accepts joins back to the same bytes, and the edits
// change nothing outside their section. Run with
// go test -fuzz=FuzzTree -fuzztime=30s ./internal/body
func FuzzTree(f *testing.F) {
	f.Add(sample)
	for _, b := range repoBodies(f) {
		f.Add(string(b))
	}
	f.Fuzz(func(t *testing.T, s string) {
		tr, err := Split([]byte(s))
		if err != nil {
			return
		}
		if string(tr.Join()) != s {
			t.Fatalf("%q: Join differs", s)
		}
		for _, sec := range tr.Sections() {
			sec, _ = tr.Section(sec.Name) // edits address the first section of a name
			var out []byte
			var err error
			if sec.Key() == "log" {
				out, err = tr.AppendLog("2026-11-01", "e")
			} else {
				out, err = tr.ReplaceSection(sec.Name, "x", "")
			}
			if errors.Is(err, ErrUnsupported) {
				if tr.refusedIn(sec) == nil {
					t.Fatalf("%q: %s edit refused with no local construct in it: %v", s, sec.Name, err)
				}
				continue
			}
			if err != nil {
				t.Fatalf("%q: %s: %v", s, sec.Name, err)
			}
			if !bytes.HasPrefix(out, []byte(s[:sec.Start])) || !bytes.HasSuffix(out, []byte(s[sec.End:])) {
				t.Fatalf("%q: editing %s changed text outside it: %q", s, sec.Name, out)
			}
		}
	})
}

// TestLocalConstructScope: a local construct blocks only the edits of its own
// section. A blockquote is not a refusal at all.
func TestLocalConstructScope(t *testing.T) {
	src := "## Goal\n\n- [ ] one\n\n## Notes\n\n> a callout\n\n- [ ]\n\n## Log\n\n- 2026-10-01: a\n"
	tr := mustSplit(t, src)
	if r := tr.Refusals(); len(r) != 1 || r[0].Construct != "task item without text" || r[0].Line != 9 {
		t.Fatalf("refusals = %+v", r)
	}
	if out, err := tr.SetChecked("goal", 1, true, ""); err != nil || !strings.Contains(string(out), "- [x] one") {
		t.Errorf("tick beside a Notes section with a local construct: %v", err)
	}
	if _, err := tr.AppendLog("2026-10-02", "b"); err != nil {
		t.Errorf("log append: %v", err)
	}
	_, err := tr.ReplaceSection("notes", "x", "")
	var u *UnsupportedError
	if !errors.As(err, &u) || u.Structural || u.Line != 9 || !errors.Is(err, ErrUnsupported) {
		t.Errorf("replacing the Notes section = %v, want a local refusal on line 9", err)
	}
	// a blockquote inside the section being edited is fine
	tr = mustSplit(t, "## Goal\n\n> note\n\n- [ ] one\n")
	if out, err := tr.SetChecked("goal", 1, true, ""); err != nil || !strings.Contains(string(out), "- [x] one") {
		t.Errorf("tick in a Goal with a blockquote: %v", err)
	}
}

func TestAddSection(t *testing.T) {
	cases := []struct{ name, in, sec, text, want string }{
		{"before Log", "## Goal\ng\n\n## Log\n- a\n", "Notes", "n", "## Goal\ng\n\n## Notes\n\nn\n\n## Log\n- a\n"},
		{"before Log, no blank line above", "## Goal\ng\n## Log\n- a\n", "Notes", "n", "## Goal\ng\n\n## Notes\n\nn\n\n## Log\n- a\n"},
		{"no Log", "# T\n\n## Goal\ng\n", "outcome", "done", "# T\n\n## Goal\ng\n\n## Outcome\n\ndone\n"},
		{"no Log, no final newline", "## Goal\ng", "Notes", "n", "## Goal\ng\n\n## Notes\n\nn\n"},
		{"empty body", "", "goal", "g", "## Goal\n\ng\n"},
		{"CRLF", "## Goal\r\ng\r\n\r\n## Log\r\n- a\r\n", "Notes", "a\nb", "## Goal\r\ng\r\n\r\n## Notes\r\n\r\na\r\nb\r\n\r\n## Log\r\n- a\r\n"},
		{"no text", "## Goal\ng\n", "Notes", "", "## Goal\ng\n\n## Notes\n"},
	}
	for _, c := range cases {
		tr, err := Split([]byte(c.in))
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		got, err := tr.AddSection(c.sec, c.text)
		if err != nil || string(got) != c.want {
			t.Errorf("%s:\ngot  %q (%v)\nwant %q", c.name, got, err, c.want)
		}
		if _, err := Split(got); err != nil {
			t.Errorf("%s: the result does not split: %v", c.name, err)
		}
	}
}

func TestLogEntriesBothForms(t *testing.T) {
	src := "## Log\n- 2026-10-01: old\n- 2026-11-01 14:30-07:00: new (at 11:05-07:00)\n- 2026-11-01 21:30+00:00: utc\n- 2026-11-01 14:30: no offset, opaque\n- 2026-11-01T14:30-07:00: T form, opaque\n"
	es, err := mustSplit(t, src).LogEntries()
	if err != nil || len(es) != 3 {
		t.Fatalf("entries %+v %v", es, err)
	}
	want := []Entry{
		{Date: "2026-10-01", Text: "old"},
		{Date: "2026-11-01", Time: "14:30", Offset: "-07:00", Text: "new (at 11:05-07:00)"},
		{Date: "2026-11-01", Time: "21:30", Offset: "+00:00", Text: "utc"},
	}
	for i, w := range want {
		g := es[i]
		if g.Date != w.Date || g.Time != w.Time || g.Offset != w.Offset || g.Text != w.Text {
			t.Errorf("entry %d = %+v, want %+v", i, g, w)
		}
	}
	out, _ := mustSplit(t, "## Log\n- 2026-10-01: a\n").AppendLog("2026-11-01 14:30-07:00", "b")
	if string(out) != "## Log\n- 2026-10-01: a\n- 2026-11-01 14:30-07:00: b\n" {
		t.Errorf("append: %q", out)
	}
}

// Ordered and bullet tasks share one count: position n stays a stable index
// (work item 261101-wder).
func TestOrderedChecklist(t *testing.T) {
	for _, eol := range []string{"\n", "\r\n"} {
		src := strings.ReplaceAll("## Goal\n\n- [ ] a\n1. [ ] b\n2. [x] c\n- [ ] d\n\n## Log\n", "\n", eol)
		tr := mustSplit(t, src)
		items, err := tr.Checklist("goal")
		if err != nil || len(items) != 4 || items[1].Text != "b" || items[2].Pos != 3 || !items[2].Done {
			t.Fatalf("%q: items %+v %v", eol, items, err)
		}
		out, err := tr.SetChecked("goal", 2, true, "")
		want := strings.Replace(src, "1. [ ] b", "1. [x] b", 1)
		if err != nil || string(out) != want {
			t.Errorf("%q: tick 2 = %q %v", eol, out, err)
		}
		out, err = tr.SetChecked("goal", 3, false, "")
		if want = strings.Replace(src, "2. [x] c", "2. [ ] c", 1); err != nil || string(out) != want {
			t.Errorf("%q: untick 3 = %q %v", eol, out, err)
		}
	}
}
