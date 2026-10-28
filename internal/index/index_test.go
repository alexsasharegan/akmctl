package index

import (
	"os"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/record"
)

func TestMatchesFixture(t *testing.T) {
	repo, err := record.Load("../../testdata/api/repos/basic")
	if err != nil {
		t.Fatal(err)
	}
	want, _ := os.ReadFile("../../testdata/api/repos/basic/docs/adr/README.md")
	if got := Generate(repo); got != string(want) {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
}

func TestCell(t *testing.T) {
	if got := cell("a  b\t| c"); got != `a b \| c` {
		t.Fatal(got)
	}
}
