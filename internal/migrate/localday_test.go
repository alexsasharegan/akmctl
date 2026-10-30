package migrate

import (
	"strings"
	"testing"
	"time"
)

// A note with no usable date gets the clock's calendar day in the clock's own
// zone: 18:30 in Arizona is already the next day in UTC, and the id must not
// read the next day.
func TestMintIDsUsesTheLocalDayOfTheClock(t *testing.T) {
	repo := filesRepo(t, map[string]string{"inbox/Evening note.md": "---\ntitle: Evening note\n---\n# Evening note\n"})
	mst := time.FixedZone("MST", -7*3600)
	env := Env{Now: time.Date(2026, 10, 29, 18, 30, 0, 0, mst), Rand: groupReader(10)}
	res := PlanEnv(repo, []Step{{ID: "m", Kind: "mint_ids", Class: "auto", Dir: "inbox"}}, env)[0]
	if len(res.Renames) != 1 || !strings.HasPrefix(res.Renames[0].To, "inbox/261029-") {
		t.Fatalf("%+v", res.Renames)
	}
}
