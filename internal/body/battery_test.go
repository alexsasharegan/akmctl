package body

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/frontmatter"
)

// hostile is a table of inputs on which a careless reader goes wrong. The
// battery asserts that the Scanner refuses each one or agrees with the oracle.
var hostile = []string{
	"## A\n\n```\n## not\n```\n\n## B\n",
	"## A\n\n~~~\n## not\n~~~\n\n## B\n",
	"## A\n\n```\n~~~\n## not\n```\n## B\n",
	"## A\n\n````\n```\n## not\n```\n````\n\n## B\n",
	"## A\n\n```\ntext\n## B\n",
	"## A\n\n  ```\n## not\n  ```\n\n## B\n",
	"## A\n\n<div>\n## not\n</div>\n\n## B\n",
	"Sub\n---\n\n## B\n",
	"## A\n\n    ## not a heading\n\n## B\n",
	"## A\n\n> ## not\n\n## B\n",
	"## A\n\n- ## not\n\n## B\n",
	"## A ##\n\n## B\n",
	"##A\n\n## B\n",
	"## A\n\n- [ ] a\n  - [x] b\n- [ ] c\n",
	"## A\n\n```\n- [ ] not\n```\n- [ ] yes\n",
	"## A\n\n- [ ] a\ncontinued\n- [x] b\n",
	"## A\n\n* [ ] a\n+ [X] b\n",
	"## A\n\n| a |\n|---|\n| - [ ] x |\n",
	"## A\r\n\r\n- [ ] x\r\n- [x] y\r\n\r\n## B\r\n",
	"## A\n\n- [ ] x",
	"## A\n\n1. [ ] one\n2. [x] two\n",
	"## A\n\n1. one\n2. two\n\n- bullet\n",
	"text\n- bullet right after text\n",
	"## A\n\n- item\n\n  - nested item\n",
	"- [ ] a\n- [ ] b\n\n---\n\n## C\n",
	"## Log\n\n- 2026-10-31: one\n- 2026-10-31: two\n",
	"# T\n\n## Goal\n\ntext\n\n- [ ] done when\n\n## Log\n\n- 2026-10-31: x\n",
	"\n\n\n",
	"- \n",
	"-\n",
	"* * *\n",
	"1) paren\n",
	"## \n",
}

func TestHostileCases(t *testing.T) {
	for _, s := range hostile {
		if refused, diff := compare([]byte(s)); !refused && diff != "" {
			t.Errorf("%q\n%s", s, diff)
		}
	}
}

// repoBodies returns the body of every record under docs/work, docs/work/done and
// docs/adr, keyed by path.
func repoBodies(t testing.TB) map[string][]byte {
	t.Helper()
	root := filepath.Join("..", "..", "docs")
	out := map[string][]byte{}
	for _, dir := range []string{"work", filepath.Join("work", "done"), "adr"} {
		files, _ := filepath.Glob(filepath.Join(root, dir, "*.md"))
		for _, f := range files {
			if filepath.Base(f) == "README.md" {
				continue
			}
			src, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			d, err := frontmatter.Parse(src)
			if err != nil {
				t.Fatalf("%s: %v", f, err)
			}
			out[f] = d.Body()
		}
	}
	return out
}

// TestRecords: none of this repository's records may be refused, and the Scanner
// must agree with the oracle on every one.
func TestRecords(t *testing.T) {
	bodies := repoBodies(t)
	if len(bodies) < 100 {
		t.Fatalf("only %d records found; is the working directory wrong?", len(bodies))
	}
	for path, b := range bodies {
		_, err := Locate(b)
		if err != nil {
			t.Errorf("%s: refused: %v", path, err)
			continue
		}
		if _, diff := compare(b); diff != "" {
			t.Errorf("%s:\n%s", path, diff)
		}
	}
}

// fragments are the lines a generated document is made of. The first group is
// ordinary record text; the second is hostile.
var (
	ordinary = []string{"## Goal", "## Log", "## Next", "text", "more text", "", "", "- [ ] item", "- [x] done",
		"- plain", "- 2026-10-31: log", "```", "~~~", "# Title", "### Sub", "* [ ] star"}
	nasty = []string{"##x", "## A ##", "  ## indented", "    ## code", "- [ ] item", "1. [ ] ordered", "  - [ ] nested",
		"```go", "  ```", "````", "`````", "~~~~", "<div>", "</div>", "<!-- c -->", "<br>", "> quote", "> ## in quote",
		"> - [ ] q", "| a | b |", "|---|---|", "---", "===", "Setext", "\ttab", "[ref]: http://x", "* * *", "-", "##", "- [ ]"}
)

func generate(r *rand.Rand) []byte {
	var ls []string
	for i, n := 0, 3+r.Intn(10); i < n; i++ {
		if r.Intn(10) < 7 {
			ls = append(ls, ordinary[r.Intn(len(ordinary))])
		} else {
			ls = append(ls, nasty[r.Intn(len(nasty))])
		}
	}
	sep := "\n"
	if r.Intn(6) == 0 {
		sep = "\r\n"
	}
	s := strings.Join(ls, sep)
	if r.Intn(2) == 0 {
		s += sep
	}
	return []byte(s)
}

// TestGenerated runs a fixed-seed stream of generated documents through the
// battery. It fails on any case where the Scanner is confidently wrong, and on a
// run so cautious that it proves nothing.
func TestGenerated(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	const n = 30000
	refused, wrong := 0, 0
	for i := 0; i < n; i++ {
		src := generate(r)
		ref, diff := compare(src)
		switch {
		case ref:
			refused++
		case diff != "":
			wrong++
			if wrong <= 5 {
				t.Errorf("%q\n%s", src, diff)
			}
		}
	}
	t.Logf("generated %d documents: %d refused, %d compared, %d wrong", n, refused, n-refused, wrong)
	if n-refused < n/5 {
		t.Errorf("only %d of %d documents were compared; the generator or the refusals are too aggressive", n-refused, n)
	}
}

// FuzzAgainstOracle: whatever the Scanner accepts, the oracle must read the same
// way. Run with go test -fuzz=FuzzAgainstOracle -fuzztime=30s ./internal/body
func FuzzAgainstOracle(f *testing.F) {
	for _, s := range hostile {
		f.Add(s)
	}
	for _, b := range repoBodies(f) {
		f.Add(string(b))
	}
	f.Fuzz(func(t *testing.T, s string) {
		if refused, diff := compare([]byte(s)); !refused && diff != "" {
			t.Fatalf("%q\n%s", s, diff)
		}
	})
}
