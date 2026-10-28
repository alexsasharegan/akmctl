package frontmatter

import (
	"fmt"
	"strconv"
	"testing"

	"go.yaml.in/yaml/v3"
)

// oracle decodes a block with a YAML 1.2 core parser and reduces it to the same
// shape as Doc: date-shaped values are read as text, as the spec requires.
func oracle(src string) (map[string]Value, error) {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(src), &root); err != nil {
		return nil, err
	}
	out := map[string]Value{}
	if root.Kind == 0 {
		return out, nil
	}
	m := root.Content[0]
	if m.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("not a mapping")
	}
	for i := 0; i < len(m.Content); i += 2 {
		if _, dup := out[m.Content[i].Value]; dup {
			return nil, fmt.Errorf("duplicate key %s", m.Content[i].Value)
		}
		v, err := oracleValue(m.Content[i+1])
		if err != nil {
			return nil, err
		}
		out[m.Content[i].Value] = v
	}
	return out, nil
}

func oracleValue(n *yaml.Node) (Value, error) {
	switch n.Kind {
	case yaml.SequenceNode:
		v := Value{Kind: List}
		for _, c := range n.Content {
			it, err := oracleValue(c)
			if err != nil {
				return v, err
			}
			v.Items = append(v.Items, it)
		}
		return v, nil
	case yaml.ScalarNode:
		switch n.ShortTag() {
		case "!!str", "!!timestamp":
			return NewString(n.Value), nil
		case "!!int":
			i, err := strconv.ParseInt(n.Value, 10, 64)
			return NewInt(i), err
		case "!!bool":
			return NewBool(n.Value == "true"), nil
		case "!!null":
			return NewNull(), nil
		}
		return Value{}, fmt.Errorf("oracle read %q as %s", n.Value, n.ShortTag())
	}
	return Value{}, fmt.Errorf("unsupported node kind %d", n.Kind)
}

func sameAsOracle(t *testing.T, name string, block string) {
	t.Helper()
	d, err := Parse([]byte("---\n" + block + "\n---\n"))
	if err != nil {
		t.Errorf("%s: parse: %v", name, err)
		return
	}
	want, err := oracle(block)
	if err != nil {
		t.Errorf("%s: oracle: %v", name, err)
		return
	}
	if len(want) != len(d.Entries) {
		t.Errorf("%s: %d keys, oracle has %d", name, len(d.Entries), len(want))
	}
	for _, e := range d.Entries {
		w, ok := want[e.Key]
		if !ok || !e.Value.Equal(w) {
			t.Errorf("%s: %s = %v, oracle %v", name, e.Key, e.Value, w)
		}
	}
}

// Every accepted fixture means the same to a YAML 1.2 core parser.
func TestOracleFixtures(t *testing.T) {
	for _, f := range loadFixtures(t) {
		if f.Accept {
			sameAsOracle(t, f.Name, f.Input)
		}
	}
}

// Every rejected fixture is either invalid YAML or one of the documented stricter-than-YAML cases.
func TestOracleRejects(t *testing.T) {
	stricter := map[string]bool{ // valid YAML that YAML-lite rejects on purpose
		"leading zero": true, "plus sign": true, "float": true, "capital True": true, "capital Null": true,
		"hex": true, "exponent": true, "underscore": true, "inf": true, "quoted backslash escape": true,
		"anchor": true, "alias": true, "tag": true, "block scalar": true, "flow mapping": true,
		"nested list": true, "unquoted wikilink": true, "true key": true, "multi-line quoted": true,
		"uneven block list": true, "indented line": true,
	}
	for _, f := range loadFixtures(t) {
		if f.Accept || stricter[f.Name] {
			continue
		}
		if _, err := oracle(f.Input); err == nil {
			t.Errorf("%s: invalid in YAML-lite and accepted by the oracle, but not listed as stricter", f.Name)
		}
	}
}

// After an edit, the oracle and the tool still agree.
func TestOracleAfterEdit(t *testing.T) {
	d, _ := Parse([]byte("---\nid: 261027-m3qz\ntags:\n  - a\n  - b\nn: 2\n---\n"))
	for _, e := range []struct {
		k string
		v Value
	}{{"id", NewString("261028-zzzz")}, {"tags", NewStrings("x", "y z", "1", "true")}, {"n", NewNull()}, {"new", NewString("2026-10-28")}, {"s", NewString("it's \"x\"")}} {
		if err := d.Set(e.k, e.v); err != nil {
			t.Fatal(err)
		}
	}
	block := string(d.Bytes())
	block = block[4 : len(block)-4]
	sameAsOracle(t, "after edit", block)
}
