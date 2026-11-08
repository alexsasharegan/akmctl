package api

import "testing"

// TestScopeMatch is the path-prefix rule of the ls filter scope (spec section
// 4.2): an entry ending in / is a prefix, any other entry matches only itself.
func TestScopeMatch(t *testing.T) {
	for _, c := range []struct {
		entries []string
		value   string
		want    bool
	}{
		{[]string{"src/store/"}, "src/store/open.go", true},
		{[]string{"src/store/"}, "src/store/", true},
		{[]string{"src/store/"}, "src/storefront/x.go", false},
		{[]string{"src/store"}, "src/store/open.go", false},
		{[]string{"cli"}, "cli", true},
		{[]string{"cli"}, "client", false},
		{[]string{"docs/", "cli"}, "cli", true},
		{nil, "cli", false},
	} {
		if got := scopeMatch(c.entries, c.value); got != c.want {
			t.Errorf("scopeMatch(%v, %q) = %v, want %v", c.entries, c.value, got, c.want)
		}
	}
}
