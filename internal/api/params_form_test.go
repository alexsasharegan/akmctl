package api

import (
	"strings"
	"testing"
)

// Every pattern-checked parameter says what a valid value looks like, so the
// error does not send the caller to a reference page (work item 261101-y3ch).
func TestEveryPatternParameterDescribesItsForm(t *testing.T) {
	for op, params := range specs {
		for name, p := range params {
			if p.pattern != nil && strings.TrimSpace(p.form) == "" {
				t.Errorf("%s.%s has a pattern and no form description", op, name)
			}
		}
	}
}

func TestPatternErrorNamesTheExpectedForm(t *testing.T) {
	root := commitsRepo(t, "claims")
	cases := []struct{ op, param, value, want string }{
		{"log", "at", "09:46", "YYYY-MM-DD HH:MM"},
		{"show", "id", "nope", "YYMMDD-xxxx"},
		{"when", "as_of", "yesterday", "YYYY-MM-DD"},
	}
	for _, c := range cases {
		req := Request{c.param: c.value}
		switch c.op {
		case "log":
			req["id"], req["text"] = "261022-e5f6", "x"
		case "when":
			req["id"] = "261022-e5f6"
		}
		_, err := call(t, root, c.op, req)
		if err == nil {
			t.Errorf("%s %s=%q: accepted", c.op, c.param, c.value)
			continue
		}
		if !strings.Contains(err.Message, "has the wrong form; expected ") || !strings.Contains(err.Message, c.want) {
			t.Errorf("%s %s=%q: %q does not name the form %q", c.op, c.param, c.value, err.Message, c.want)
		}
	}
}
