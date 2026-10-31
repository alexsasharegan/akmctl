package api

import (
	"sort"
	"strings"
	"unicode"

	"github.com/alexsasharegan/akmctl/internal/record"
)

// similarMax is the most open items `new work` reports as similar.
const similarMax = 5

var similarOpen = map[string]bool{"raw": true, "ready": true, "active": true, "blocked": true}

var stopWords = map[string]bool{
	"a": true, "an": true, "and": true, "are": true, "as": true, "at": true, "be": true, "by": true,
	"for": true, "from": true, "in": true, "is": true, "it": true, "of": true, "on": true, "or": true,
	"the": true, "to": true, "with": true,
}

// words splits text into its lower-case words: runs of letters and digits,
// without the stop words, each once, in order of first appearance.
func words(text string) []string {
	var out []string
	seen := map[string]bool{}
	for _, w := range strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}) {
		if !stopWords[w] && !seen[w] {
			seen[w] = true
			out = append(out, w)
		}
	}
	return out
}

// similarItems returns up to similarMax open work items that share words with
// the new item's slug and title. The rule is in section 4.10 of the API spec.
func similarItems(repo *record.Repo, slug, title string) []Obj {
	mine := words(strings.ReplaceAll(slug, "-", " ") + " " + title)
	if len(mine) == 0 {
		return nil
	}
	need := 2
	if len(mine) < need {
		need = len(mine)
	}
	type hit struct {
		r      *record.Record
		shared []string
	}
	var hits []hit
	for _, r := range repo.Records {
		if r.Type != record.Work || !similarOpen[r.Status()] {
			continue
		}
		theirs := map[string]bool{}
		for _, w := range words(r.Str("title") + " " + r.Str("summary")) {
			theirs[w] = true
		}
		var shared []string
		for _, w := range mine {
			if theirs[w] {
				shared = append(shared, w)
			}
		}
		if len(shared) >= need {
			hits = append(hits, hit{r, shared})
		}
	}
	sort.SliceStable(hits, func(i, j int) bool {
		if len(hits[i].shared) != len(hits[j].shared) {
			return len(hits[i].shared) > len(hits[j].shared)
		}
		return hits[i].r.ID < hits[j].r.ID
	})
	if len(hits) > similarMax {
		hits = hits[:similarMax]
	}
	out := make([]Obj, 0, len(hits))
	for _, h := range hits {
		sh := make([]any, len(h.shared))
		for i, w := range h.shared {
			sh[i] = w
		}
		out = append(out, O("id", h.r.ID, "status", h.r.Status(), "title", h.r.Str("title"), "path", h.r.Path, "words", sh))
	}
	return out
}
