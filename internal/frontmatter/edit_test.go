package frontmatter

import "testing"

func edited(t *testing.T, in string, f func(*Doc) error) string {
	t.Helper()
	d, err := Parse([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	if err := f(d); err != nil {
		t.Fatal(err)
	}
	return string(d.Bytes())
}

func TestSetKeepsEverythingElse(t *testing.T) {
	in := "---\n# keep\nid:   261027-m3qz   # why\nstatus: \"active\"\ntags:\n  - a\n  - b\nlist: [x,   y]\n---\nbody\n"
	cases := []struct {
		name string
		key  string
		v    Value
		want string
	}{
		{"scalar keeps comment", "id", NewString("261028-aaaa"), "---\n# keep\nid:   261028-aaaa   # why\nstatus: \"active\"\ntags:\n  - a\n  - b\nlist: [x,   y]\n---\nbody\n"},
		{"same value is a no-op", "status", NewString("active"), in},
		{"quote changes only the value", "status", NewString("done"), "---\n# keep\nid:   261027-m3qz   # why\nstatus: done\ntags:\n  - a\n  - b\nlist: [x,   y]\n---\nbody\n"},
		{"block list stays block", "tags", NewStrings("a", "b", "c"), "---\n# keep\nid:   261027-m3qz   # why\nstatus: \"active\"\ntags:\n  - a\n  - b\n  - c\nlist: [x,   y]\n---\nbody\n"},
		{"block list to scalar", "tags", NewString("z"), "---\n# keep\nid:   261027-m3qz   # why\nstatus: \"active\"\ntags: z\nlist: [x,   y]\n---\nbody\n"},
		{"inline list", "list", NewStrings("p", "q r", "s,t"), "---\n# keep\nid:   261027-m3qz   # why\nstatus: \"active\"\ntags:\n  - a\n  - b\nlist: [p, q r, \"s,t\"]\n---\nbody\n"},
		{"new key at the end", "priority", NewInt(2), "---\n# keep\nid:   261027-m3qz   # why\nstatus: \"active\"\ntags:\n  - a\n  - b\nlist: [x,   y]\npriority: 2\n---\nbody\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := edited(t, in, func(d *Doc) error { return d.Set(c.key, c.v) })
			if got != c.want {
				t.Fatalf("got\n%q\nwant\n%q", got, c.want)
			}
		})
	}
}

func TestSetQuoting(t *testing.T) {
	in := "---\na: 1\n---\n"
	for v, want := range map[string]string{
		"plain":       "a: plain",
		"2026-10-28":  "a: 2026-10-28",
		"12:30":       "a: 12:30",
		"yes":         "a: yes",
		"02":          "a: \"02\"",
		"1.2":         "a: \"1.2\"",
		"true":        "a: \"true\"",
		"True":        "a: \"True\"",
		"null":        "a: \"null\"",
		"x: y":        "a: \"x: y\"",
		"a # b":       "a: \"a # b\"",
		"[x]":         "a: \"[x]\"",
		" lead":       "a: \" lead\"",
		"say \"hi\"":  "a: say \"hi\"",
		"-":           "a: \"-\"",
		"back\\slash": "a: back\\slash",
		"":            "a: \"\"",
	} {
		d, _ := Parse([]byte(in))
		if err := d.Set("a", NewString(v)); err != nil {
			t.Fatalf("%q: %v", v, err)
		}
		got := string(d.Bytes())[4 : len(d.Bytes())-5]
		if got != want+"\n" && got != want {
			t.Errorf("%q written as %q, want %q", v, got, want)
		}
		back, _ := d.Get("a")
		if back.Kind != String || back.Str != v {
			t.Errorf("%q read back as %v", v, back)
		}
	}
	d, _ := Parse([]byte(in))
	if err := d.Set("a", NewString("two\nlines")); err == nil {
		t.Error("line break accepted")
	}
}

func TestDeleteRename(t *testing.T) {
	in := "---\na: 1\ntags:\n- x\n- y\nb: 2\n---\n"
	got := edited(t, in, func(d *Doc) error { _, err := d.Delete("tags"); return err })
	if got != "---\na: 1\nb: 2\n---\n" {
		t.Fatalf("delete: %q", got)
	}
	got = edited(t, in, func(d *Doc) error { return d.Rename("b", "c") })
	if got != "---\na: 1\ntags:\n- x\n- y\nc: 2\n---\n" {
		t.Fatalf("rename: %q", got)
	}
	d, _ := Parse([]byte(in))
	if d.Rename("a", "b") == nil {
		t.Fatal("rename onto an existing key")
	}
}

func TestEmptyAndCRLF(t *testing.T) {
	got := edited(t, "---\r\nowner:\r\n---\r\n", func(d *Doc) error { return d.Set("owner", NewString("me")) })
	if got != "---\r\nowner: me\r\n---\r\n" {
		t.Fatalf("%q", got)
	}
	got = edited(t, "---\r\na: 1\r\n---\r\n", func(d *Doc) error { return d.Set("b", NewBool(true)) })
	if got != "---\r\na: 1\r\nb: true\r\n---\r\n" {
		t.Fatalf("%q", got)
	}
}

func TestBareAppend(t *testing.T) {
	d, _ := ParseBare([]byte("layout: 1.0.0"))
	if err := d.Set("stale_days", NewInt(30)); err != nil {
		t.Fatal(err)
	}
	if string(d.Bytes()) != "layout: 1.0.0\nstale_days: 30\n" {
		t.Fatalf("%q", d.Bytes())
	}
	if ok, _ := d.Delete("stale_days"); !ok || string(d.Bytes()) != "layout: 1.0.0" {
		t.Fatalf("%q", d.Bytes())
	}
}
