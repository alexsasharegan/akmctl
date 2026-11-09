package api

import (
	"strings"
	"testing"

	"github.com/alexsasharegan/akmctl/internal/migrate"
)

// twoStepRegistry is a chain 0.1.0 -> 0.2.0 -> 0.3.0: the first step fixes WK3
// (created_at), the second WK7 (the Result heading).
func twoStepRegistry(t *testing.T) {
	t.Helper()
	old := migrate.Registry
	migrate.Registry = []migrate.Migration{
		{From: "0.1.0", To: "0.2.0", Steps: []migrate.Step{
			{ID: "rename-created", Kind: "rename_field", Class: "auto", Scope: "work", From: "created_at", To: "created", Fixes: []string{"WK3"}},
		}},
		{From: "0.2.0", To: "0.3.0", Steps: []migrate.Step{
			{ID: "rename-result", Kind: "rename_heading", Class: "auto", Scope: "work", From: "Result", To: "Outcome", Fixes: []string{"WK7"}},
		}},
	}
	t.Cleanup(func() { migrate.Registry = old })
}

func checkClasses(t *testing.T, root, target string) map[string]any {
	t.Helper()
	resp := mustOK(t, root, "check", Request{"target": target, "as_of": "2026-10-28"})
	fs, _ := resp.Get("findings")
	classes := map[string]any{}
	for _, f := range fs.([]Obj) {
		r, _ := f.Get("rule")
		c, _ := f.Get("class")
		classes[r.(string)] = c
	}
	return classes
}

func TestCheckTargetAcceptsAnyLayoutTheChainReaches(t *testing.T) {
	twoStepRegistry(t)
	root := oldRepo(t)

	// the middle layout: only the first step applies, so WK3 is auto and WK7 is not
	mid := checkClasses(t, root, "0.2.0")
	if mid["WK3"] != "auto" {
		t.Fatalf("target 0.2.0: WK3 is %v, want auto (classes %v)", mid["WK3"], mid)
	}
	if c, ok := mid["WK7"]; !ok || c == "auto" {
		t.Fatalf("target 0.2.0: WK7 is %v, want a finding the first step does not fix (classes %v)", c, mid)
	}

	// the end of the chain: both steps apply
	end := checkClasses(t, root, "0.3.0")
	if end["WK3"] != "auto" || end["WK7"] != "auto" {
		t.Fatalf("target 0.3.0: %v", end)
	}

	// the repository's own layout: an empty chain, nothing is auto
	own := checkClasses(t, root, "0.1.0")
	if own["WK3"] == "auto" || own["WK7"] == "auto" {
		t.Fatalf("target 0.1.0: %v", own)
	}
}

func TestCheckTargetRefusesALayoutTheChainDoesNotReach(t *testing.T) {
	twoStepRegistry(t)
	root := oldRepo(t)
	for _, target := range []string{"0.9.0", "0.0.5"} {
		_, err := call(t, root, "check", Request{"target": target})
		if err == nil || err.Code != "invalid_request" {
			t.Fatalf("target %s: got %v, want invalid_request", target, err)
		}
		for _, want := range []string{target, "0.1.0", "0.2.0", "0.3.0"} {
			if !strings.Contains(err.Message, want) {
				t.Errorf("target %s: message %q lacks %s", target, err.Message, want)
			}
		}
	}
}
