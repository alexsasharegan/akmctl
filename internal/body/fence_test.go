package body

import "testing"

func TestFenceRule(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []bool
	}{
		{"backticks", []string{"a", "```go", "x", "```", "b"}, []bool{false, true, true, true, false}},
		{"tilde does not close backticks", []string{"```", "~~~", "## h", "```", "## i"}, []bool{true, true, true, true, false}},
		{"backtick does not close tildes", []string{"~~~", "```", "## h", "~~~", "## i"}, []bool{true, true, true, true, false}},
		{"shorter run does not close", []string{"````", "```", "## h", "````", "## i"}, []bool{true, true, true, true, false}},
		{"closer with text does not close", []string{"```", "``` x", "## h", "```"}, []bool{true, true, true, true}},
		{"backtick in info string is not a fence", []string{"``` a`b", "## h"}, []bool{false, false}},
		{"unclosed runs to the end", []string{"```", "## h", "## i"}, []bool{true, true, true}},
		{"indented opener is not claimed", []string{"  ```", "## h"}, []bool{false, false}},
	}
	for _, c := range cases {
		var f FenceTracker
		for i, l := range c.in {
			if got := f.In(l); got != c.want[i] {
				t.Errorf("%s: line %d %q: In = %v, want %v", c.name, i, l, got, c.want[i])
			}
		}
	}
}
