// Package frontmatter reads and edits YAML-lite, the frontmatter format of AKM
// records: a strict subset of YAML 1.2 core. Parsing keeps byte offsets so that
// an edit changes only the bytes of the entry it touches.
package frontmatter

import "strconv"

// Kind is the type of a Value.
type Kind uint8

const (
	Null Kind = iota
	String
	Int
	Bool
	List
)

// Value is a YAML-lite scalar or a list of scalars.
type Value struct {
	Kind  Kind
	Str   string
	Int   int64
	Bool  bool
	Items []Value // List only; every item is a scalar
}

func NewString(s string) Value { return Value{Kind: String, Str: s} }
func NewInt(n int64) Value     { return Value{Kind: Int, Int: n} }
func NewBool(b bool) Value     { return Value{Kind: Bool, Bool: b} }
func NewNull() Value           { return Value{} }
func NewList(v ...Value) Value { return Value{Kind: List, Items: v} }
func NewStrings(s ...string) Value {
	items := make([]Value, len(s))
	for i, x := range s {
		items[i] = NewString(x)
	}
	return Value{Kind: List, Items: items}
}

// Equal reports whether two values hold the same data.
func (v Value) Equal(o Value) bool {
	if v.Kind != o.Kind {
		return false
	}
	switch v.Kind {
	case String:
		return v.Str == o.Str
	case Int:
		return v.Int == o.Int
	case Bool:
		return v.Bool == o.Bool
	case List:
		if len(v.Items) != len(o.Items) {
			return false
		}
		for i := range v.Items {
			if !v.Items[i].Equal(o.Items[i]) {
				return false
			}
		}
	}
	return true
}

// Strings returns the string items of a list, or the string itself for a
// String value, and nil for anything else. Non-string items are skipped.
func (v Value) Strings() []string {
	switch v.Kind {
	case String:
		return []string{v.Str}
	case List:
		var out []string
		for _, it := range v.Items {
			if it.Kind == String {
				out = append(out, it.Str)
			}
		}
		return out
	}
	return nil
}

func (v Value) String() string {
	switch v.Kind {
	case String:
		return v.Str
	case Int:
		return strconv.FormatInt(v.Int, 10)
	case Bool:
		return strconv.FormatBool(v.Bool)
	case List:
		s := "["
		for i, it := range v.Items {
			if i > 0 {
				s += ", "
			}
			s += it.String()
		}
		return s + "]"
	}
	return "null"
}
