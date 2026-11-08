package conformance

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

// Matchers (ADR 261108-15dw). A string in an expected response may hold:
//
//	$any            the whole string: any value of any type, but present
//	$id(YYMMDD)     an id of that date; $id(YYMMDD,b) is a second, distinct one
//	$hash           a git object name, 7 to 40 hex digits
//
// The first occurrence of a marker binds the actual text and every other
// occurrence of the same marker means that text, in the response and in the
// paths and files a case expects afterwards.
var reMarker = regexp.MustCompile(`\$id\(\d{6}(?:,\w+)?\)|\$hash`)

const idAlphabet = `[0-9a-hjkmnp-tv-z]{4}`

// binds maps a marker to the text it stands for in this case.
type binds map[string]string

func hasMarker(s string) bool { return s == "$any" || reMarker.MatchString(s) }

// matchString reports whether actual fits the expected string, binding the
// markers it holds. Without a marker the strings must be equal.
func matchString(exp, act string, b binds) bool {
	if !reMarker.MatchString(exp) {
		return exp == act
	}
	var pat strings.Builder
	pat.WriteString("^")
	var keys []string
	last := 0
	for _, loc := range reMarker.FindAllStringIndex(exp, -1) {
		pat.WriteString(regexp.QuoteMeta(exp[last:loc[0]]))
		m := exp[loc[0]:loc[1]]
		if v, ok := b[m]; ok {
			pat.WriteString(regexp.QuoteMeta(v))
		} else {
			if m == "$hash" {
				pat.WriteString(`([0-9a-f]{7,40})`)
			} else {
				pat.WriteString(`(` + m[len("$id("):len("$id(")+6] + `-` + idAlphabet + `)`)
			}
			keys = append(keys, m)
		}
		last = loc[1]
	}
	pat.WriteString(regexp.QuoteMeta(exp[last:]) + "$")
	// a marker used twice in one string must agree: the regexp captures each use
	// separately, so check the captures after
	re := regexp.MustCompile(pat.String())
	got := re.FindStringSubmatch(act)
	if got == nil {
		return false
	}
	seen := map[string]string{}
	for i, k := range keys {
		if v, ok := seen[k]; ok && v != got[i+1] {
			return false
		}
		seen[k] = got[i+1]
	}
	for k, v := range seen {
		b[k] = v
	}
	return true
}

// subst replaces the bound markers in s; an unbound marker stays as it is.
func subst(s string, b binds) string {
	return reMarker.ReplaceAllStringFunc(s, func(m string) string {
		if v, ok := b[m]; ok {
			return v
		}
		return m
	})
}

// match compares an expected value with an actual one, both as decoded from
// JSON, and returns the first difference with its path.
func match(exp, act any, path string, b binds) error {
	if s, ok := exp.(string); ok {
		if s == "$any" {
			if act == nil {
				return fmt.Errorf("%s: $any, but the member is null", path)
			}
			return nil
		}
		a, ok := act.(string)
		if !ok || !matchString(s, a, b) {
			return fmt.Errorf("%s: got %v, want %q", path, show(act), s)
		}
		return nil
	}
	switch e := exp.(type) {
	case map[string]any:
		a, ok := act.(map[string]any)
		if !ok {
			return fmt.Errorf("%s: got %v, want an object", path, show(act))
		}
		var keys []string
		for k := range e {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v, present := a[k]
			if !present {
				return fmt.Errorf("%s.%s: missing", path, k)
			}
			if err := match(e[k], v, path+"."+k, b); err != nil {
				return err
			}
		}
		for k := range a {
			if _, ok := e[k]; !ok {
				return fmt.Errorf("%s.%s: unexpected member %v", path, k, show(a[k]))
			}
		}
		return nil
	case []any:
		a, ok := act.([]any)
		if !ok || len(a) != len(e) {
			return fmt.Errorf("%s: got %v, want %d elements", path, show(act), len(e))
		}
		for i := range e {
			if err := match(e[i], a[i], fmt.Sprintf("%s[%d]", path, i), b); err != nil {
				return err
			}
		}
		return nil
	}
	if !reflect.DeepEqual(exp, act) {
		return fmt.Errorf("%s: got %v, want %v", path, show(act), show(exp))
	}
	return nil
}

func show(v any) string {
	s := fmt.Sprintf("%v", v)
	if len(s) > 120 {
		s = s[:120] + "..."
	}
	return s
}
