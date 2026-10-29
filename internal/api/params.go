package api

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/record"
)

// Version is the API version the tool implements.
const Version = "1.0"

// Error is a failed operation: a code from the closed set of the spec, a
// message, and optional details.
type Error struct {
	Code    string
	Message string
	Details Obj
}

func (e *Error) Error() string { return e.Code + ": " + e.Message }

func invalid(param, format string, a ...any) *Error {
	return &Error{"invalid_request", fmt.Sprintf(format, a...), O("parameter", param)}
}

// Request holds the parameters of one call, as given.
type Request map[string]any

// Parameter kinds.
const (
	kString  = "string"
	kBool    = "bool"
	kInt     = "int"
	kStrings = "strings"
	kObject  = "object"
	kObjects = "objects" // a list of objects
)

type param struct {
	kind     string
	pattern  *regexp.Regexp
	enum     []string
	min      int
	required bool
	empty    bool // an empty string is a value
}

var (
	reDate    = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	reDateTS  = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}(T\d{2}:\d{2})?$`)
	reID      = regexp.MustCompile(`^\d{6}-[0-9a-hjkmnp-tv-z]{4}$`)
	reVersion = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	reEtag    = regexp.MustCompile(`^[0-9a-f]{16}$`)
)

var workStatus = []string{"raw", "ready", "active", "blocked", "done", "dropped"}
var anyStatus = append(append([]string(nil), workStatus...), "proposed", "accepted", "rejected", "deprecated", "superseded")

func s(extra ...func(*param)) param {
	p := param{kind: kString}
	for _, f := range extra {
		f(&p)
	}
	return p
}
func pat(re *regexp.Regexp) func(*param) { return func(p *param) { p.pattern = re } }
func req(p *param)                       { p.required = true }
func empty(p *param)                     { p.empty = true }

// specs lists the parameters of each read operation.
var specs = map[string]map[string]param{
	"version": {},
	"ls": {
		"type": {kind: kString, enum: []string{"work", "adr"}}, "status": {kind: kStrings, enum: anyStatus},
		"parent": s(pat(reID)), "blocked": {kind: kBool}, "scope": s(), "since": s(pat(reDate)),
		"sort": {kind: kStrings}, "fields": {kind: kStrings}, "limit": {kind: kInt, min: 1}, "offset": {kind: kInt, min: 0},
		"count_only": {kind: kBool}, "next": {kind: kBool},
	},
	"search": {
		"text": s(req), "body": {kind: kBool}, "type": {kind: kString, enum: []string{"work", "adr"}}, "status": {kind: kStrings, enum: anyStatus},
		"fields": {kind: kStrings}, "limit": {kind: kInt, min: 1}, "offset": {kind: kInt, min: 0},
	},
	"show":    {"id": s(req, pat(reID)), "section": s()},
	"board":   {"statuses": {kind: kStrings, enum: workStatus}},
	"find":    {"id": s(req, pat(reID))},
	"when":    {"id": s(req, pat(reID)), "as_of": s(pat(reDate))},
	"check":   {"target": s(pat(reVersion)), "rules": {kind: kStrings}, "as_of": s(pat(reDateTS))},
	"new":     {"kind": {kind: kString, enum: []string{"work", "adr"}, required: true}, "slug": s(req), "title": s(), "fields": {kind: kObject}, "as_of": s(pat(reDateTS))},
	"set":     {"id": s(req, pat(reID)), "changes": {kind: kObject, required: true}, "if_match": s(pat(reEtag)), "dry_run": {kind: kBool}, "as_of": s(pat(reDateTS))},
	"section": {"id": s(req, pat(reID)), "name": s(req), "text": s(req, empty), "mode": {kind: kString, enum: []string{"replace", "append"}}, "if_match": s(pat(reEtag)), "dry_run": {kind: kBool}, "as_of": s(pat(reDateTS))},
	"close":   {"id": s(req, pat(reID)), "status": {kind: kString, enum: []string{"done", "dropped"}, required: true}, "outcome": s(), "if_match": s(pat(reEtag)), "dry_run": {kind: kBool}, "as_of": s(pat(reDateTS))},
	"claim":   {"id": s(req, pat(reID)), "as": s(req), "if_match": s(pat(reEtag)), "takeover": {kind: kBool}, "dry_run": {kind: kBool}, "as_of": s(pat(reDateTS))},
	"release": {"id": s(req, pat(reID)), "if_match": s(pat(reEtag)), "dry_run": {kind: kBool}, "as_of": s(pat(reDateTS))},
	"migrate": {"mode": {kind: kString, enum: []string{"report", "plan", "apply_safe", "handoff"}}, "target": s(pat(reVersion)), "acknowledge_backup": {kind: kBool}, "migrations": {kind: kObjects}, "as_of": s(pat(reDateTS))},
	"index":   {"dry_run": {kind: kBool}, "include_content": {kind: kBool}},
}

// Operations lists the implemented operation names.
func Operations() []string {
	var out []string
	for k := range specs {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Params returns the parameter names and kinds of an operation, for bindings.
func Params(op string) map[string]string {
	out := map[string]string{}
	for k, p := range specs[op] {
		out[k] = p.kind
	}
	return out
}

// validate checks a request against the operation's parameters.
func validate(op string, r Request) *Error {
	sp, ok := specs[op]
	if !ok {
		return &Error{"invalid_request", "unknown operation " + op, O("parameter", "operation")}
	}
	for k := range r {
		if _, ok := sp[k]; !ok {
			return invalid(k, "unknown parameter %q", k)
		}
	}
	for name, p := range sp {
		v, given := r[name]
		if !given {
			if p.required {
				return invalid(name, "missing required parameter %q", name)
			}
			continue
		}
		switch p.kind {
		case kString:
			str, ok := v.(string)
			if !ok {
				return invalid(name, "%s must be a string", name)
			}
			if err := checkString(name, str, p); err != nil {
				return err
			}
		case kBool:
			if _, ok := v.(bool); !ok {
				return invalid(name, "%s must be true or false", name)
			}
		case kInt:
			n, ok := asInt(v)
			if !ok {
				return invalid(name, "%s must be an integer", name)
			}
			if n < p.min {
				return invalid(name, "%s must be at least %d", name, p.min)
			}
		case kObject:
			if _, ok := v.(map[string]any); !ok {
				return invalid(name, "%s must be an object", name)
			}
		case kObjects:
			l, ok := v.([]any)
			if !ok {
				if ms, isMaps := v.([]map[string]any); isMaps {
					for _, m := range ms {
						l = append(l, m)
					}
					ok = true
				}
			}
			if !ok {
				return invalid(name, "%s must be a list of objects", name)
			}
			for _, x := range l {
				if _, ok := x.(map[string]any); !ok {
					return invalid(name, "%s must be a list of objects", name)
				}
			}
		case kStrings:
			l, ok := asStrings(v)
			if !ok {
				return invalid(name, "%s must be a list of strings", name)
			}
			for _, x := range l {
				if err := checkString(name, x, p); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func checkString(name, v string, p param) *Error {
	if p.enum != nil {
		for _, e := range p.enum {
			if e == v {
				return nil
			}
		}
		return invalid(name, "%s: unknown value %q", name, v)
	}
	if p.pattern != nil && !p.pattern.MatchString(v) {
		return invalid(name, "%s: %q has the wrong form", name, v)
	}
	if p.required && !p.empty && strings.TrimSpace(v) == "" {
		return invalid(name, "%s must not be empty", name)
	}
	return nil
}

func asInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case float64:
		if n == float64(int(n)) {
			return int(n), true
		}
	}
	return 0, false
}

func asStrings(v any) ([]string, bool) {
	switch l := v.(type) {
	case []string:
		return l, true
	case []any:
		out := make([]string, len(l))
		for i, x := range l {
			s, ok := x.(string)
			if !ok {
				return nil, false
			}
			out[i] = s
		}
		return out, true
	}
	return nil, false
}

func (r Request) str(k string) string   { s, _ := r[k].(string); return s }
func (r Request) boolean(k string) bool { b, _ := r[k].(bool); return b }
func (r Request) strs(k string) []string {
	l, _ := asStrings(r[k])
	return l
}
func (r Request) integer(k string, def int) int {
	if v, ok := r[k]; ok {
		n, _ := asInt(v)
		return n
	}
	return def
}

var _ = record.Work
