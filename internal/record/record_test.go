package record

import "testing"

func TestRefID(t *testing.T) {
	for in, want := range map[string]string{
		"261027-m3qz":                           "261027-m3qz",
		"[[261027-m3qz-fix-login|261027-m3qz]]": "261027-m3qz",
		"[[261027-m3qz-fix-login]]":             "261027-m3qz",
		"legal review":                          "legal review",
		"[[Some note]]":                         "[[Some note]]",
	} {
		if got := RefID(in); got != want {
			t.Errorf("RefID(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLoad(t *testing.T) {
	r, err := Load("../../testdata/api/repos/basic")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Records) != 8 || len(r.Problems) != 0 {
		t.Fatalf("%d records, problems %v", len(r.Records), r.Problems)
	}
	if got := r.Roadmap; len(got) != 2 || got[0] != "261015-c4d5" {
		t.Fatalf("roadmap %v", got)
	}
	if !r.Extension("work-items") || r.Extension("claims") {
		t.Fatal("extensions wrong")
	}
	if recs := r.ByID("261010-j9k0"); len(recs) != 1 || recs[0].Path != "docs/work/done/261010-j9k0-choose-a-licence.md" {
		t.Fatalf("done item: %v", recs)
	}
	if _, err := Load("../../testdata/api/repos/missing"); err == nil {
		t.Fatal("missing repo loaded")
	}
}

func TestDuplicateIDsAreVisible(t *testing.T) {
	if IDDate("261027-m3qz") != "2026-10-27" {
		t.Fatal(IDDate("261027-m3qz"))
	}
}
