package frontmatter

import "testing"

func TestSetAfter(t *testing.T) {
	src := "---\nupdated: 2026-10-28\nclaimed_by: a\n---\n"
	got := edited(t, src, func(d *Doc) error { return d.SetAfter("closed", "updated", NewString("2026-10-29")) })
	if got != "---\nupdated: 2026-10-28\nclosed: 2026-10-29\nclaimed_by: a\n---\n" {
		t.Fatalf("%q", got)
	}
	// an existing key is changed in place; a missing anchor appends
	got = edited(t, got, func(d *Doc) error { return d.SetAfter("closed", "claimed_by", NewString("2026-11-01")) })
	if got != "---\nupdated: 2026-10-28\nclosed: 2026-11-01\nclaimed_by: a\n---\n" {
		t.Fatalf("%q", got)
	}
	got = edited(t, src, func(d *Doc) error { return d.SetAfter("closed", "nope", NewString("x")) })
	if got != "---\nupdated: 2026-10-28\nclaimed_by: a\nclosed: x\n---\n" {
		t.Fatalf("%q", got)
	}
	got = edited(t, "---\r\nupdated: 1\r\n---\r\n", func(d *Doc) error { return d.SetAfter("closed", "updated", NewString("b")) })
	if got != "---\r\nupdated: 1\r\nclosed: b\r\n---\r\n" {
		t.Fatalf("%q", got)
	}
}
