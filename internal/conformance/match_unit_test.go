package conformance

import "testing"

func TestMatchers(t *testing.T) {
	b := binds{}
	if !matchString("docs/work/$id(261028)-x.md", "docs/work/261028-ab3d-x.md", b) || b["$id(261028)"] != "261028-ab3d" {
		t.Fatalf("bind: %v", b)
	}
	// a bound marker means its text
	if !matchString("$id(261028)", "261028-ab3d", b) || matchString("$id(261028)", "261028-zzzz", b) {
		t.Fatal("a bound id must repeat")
	}
	// a second, distinct id
	if !matchString("$id(261028,b)", "261028-ffff", b) || b["$id(261028,b)"] != "261028-ffff" {
		t.Fatal("second id")
	}
	// wrong date, wrong alphabet (the letters i, l, o, u are not used)
	for _, bad := range []string{"261029-ab3d", "261028-abid", "261028-ab3", "x261028-ab3d"} {
		if matchString("$id(261028)", bad, binds{}) {
			t.Errorf("%q matched", bad)
		}
	}
	if !matchString("$hash", "0123abc", binds{}) || matchString("$hash", "xyz", binds{}) || matchString("$hash", "012345", binds{}) {
		t.Fatal("hash")
	}
	hb := binds{}
	if matchString("$hash $hash", "0123abc 4567def", hb) {
		t.Fatal("one marker twice in a string must agree")
	}
	if !matchString("$hash..$hash", "0123abc..0123abc", hb) {
		t.Fatal("the same hash twice")
	}
	if err := match(map[string]any{"a": "$any", "b": []any{"$id(261028)"}}, map[string]any{"a": 1.0, "b": []any{"261028-ab3d"}}, "$", binds{}); err != nil {
		t.Fatal(err)
	}
	if err := match(map[string]any{"a": "$any"}, map[string]any{"a": nil}, "$", binds{}); err == nil {
		t.Fatal("$any matched null")
	}
	if err := match(map[string]any{"a": 1.0}, map[string]any{"a": 1.0, "extra": true}, "$", binds{}); err == nil {
		t.Fatal("an unexpected member passed")
	}
}
