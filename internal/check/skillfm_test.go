package check

import (
	"testing"

	"go.yaml.in/yaml/v3"
)

// S2 accepts block scalars in a skill's frontmatter, and reads their value as YAML does.
func TestS2BlockScalars(t *testing.T) {
	skill := func(fm string) map[string]func(string) string {
		return map[string]func(string) string{".agents/skills/demo/SKILL.md": func(string) string { return "---\n" + fm + "---\nbody\n" }}
	}
	cases := []struct {
		name string
		fm   string
		ok   bool
	}{
		{"plain", "name: demo\ndescription: x\n", true},
		{"folded description", "name: demo\ndescription: >\n  Use this when a long\n  line is folded.\n", true},
		{"literal description", "name: demo\ndescription: |\n  Line one.\n  Line two.\n", true},
		{"strip and keep indicators", "name: demo\ndescription: >-\n  a\nother: |+\n  b\n\n", true},
		{"block before name", "description: >\n  x\nname: demo\n", true},
		{"folded name", "name: >\n  demo\ndescription: x\n", true},
		{"comment after header", "name: demo\ndescription: > # why\n  x\n", true},
		{"empty block description", "name: demo\ndescription: >\n\nother: x\n", false},
		{"block name differs", "name: |\n  other\ndescription: x\n", false},
		{"no description", "name: demo\nother: |\n  x\n", false},
		{"still rejects an anchor", "name: demo\ndescription: x\nnote: &a v\n", false},
		{"still rejects a duplicate key", "name: demo\ndescription: >\n  x\ndescription: y\n", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := Run(repoWith(t, skill(c.fm)), Options{AsOf: asOf})
			var s2 bool
			for _, f := range res.Findings {
				if f.Rule == "S2" {
					s2 = true
				}
			}
			if s2 == c.ok {
				t.Fatalf("ok=%v but S2 reported=%v: %v", c.ok, s2, res.Findings)
			}
		})
	}
}

// The value of a block scalar matches a YAML parser's.
func TestBlockValueMatchesYAML(t *testing.T) {
	bodies := []string{
		"  a\n  b\n", "  a\n\n  b\n", "  a\n\n\n  b\n", "  a\n   more\n  b\n", "  a\n  b\n\n\n", "  one\n", "\n  a\n",
	}
	for _, header := range []string{"|", ">", "|-", ">-", "|+", ">+", "|2", ">2-"} {
		for _, body := range bodies {
			src := "k: " + header + "\n" + body
			var want map[string]string
			if err := yaml.Unmarshal([]byte(src), &want); err != nil {
				continue
			}
			full := "---\n" + src + "---\n"
			_, got := readSkillFrontmatter([]byte(full))
			if got["k"] != want["k"] {
				t.Errorf("%q: got %q, want %q", src, got["k"], want["k"])
			}
		}
	}
}
