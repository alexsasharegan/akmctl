package api

import (
	"os"
	"path/filepath"
	"time"

	"github.com/alexsasharegan/akmctl/internal/check"
	"github.com/alexsasharegan/akmctl/internal/index"
	"github.com/alexsasharegan/akmctl/internal/record"
)

func opCheck(repo *record.Repo, q Request) (Obj, *Error) {
	if t := q.str("target"); t != "" && t != check.LayoutVersion {
		return nil, invalid("target", "target: layout version %s is not supported (this tool checks %s)", t, check.LayoutVersion)
	}
	opt := check.Options{Rules: q.strs("rules")}
	if s := q.str("as_of"); s != "" {
		layout := "2006-01-02"
		if len(s) > 10 {
			layout = "2006-01-02T15:04"
		}
		t, err := time.Parse(layout, s)
		if err != nil {
			return nil, invalid("as_of", "as_of is not a valid date or timestamp")
		}
		opt.AsOf = t
	}
	res := check.Run(repo, opt)
	errs, warns := 0, 0
	findings := make([]Obj, 0, len(res.Findings))
	for _, f := range res.Findings {
		if f.Severity == "E" {
			errs++
		} else {
			warns++
		}
		var file, id, line any
		if f.File != "" {
			file = f.File
		}
		if f.ID != "" {
			id = f.ID
		}
		if f.Line > 0 {
			line = f.Line
		}
		findings = append(findings, O("rule", f.Rule, "severity", f.Severity, "file", file, "line", line, "id", id, "message", f.Message, "class", nil))
	}
	rec := O("state", res.State)
	if len(res.Steps) > 0 {
		rec = append(rec, KV{"steps", res.Steps})
	}
	return O("clean", errs == 0, "errors", errs, "warnings", warns, "count", len(findings), "findings", findings, "reconcile", rec), nil
}

func opIndex(repo *record.Repo, q Request) (Obj, *Error) {
	want := index.Generate(repo)
	full := filepath.Join(repo.Root, filepath.FromSlash(index.Path))
	cur, err := os.ReadFile(full)
	changed := err != nil || string(cur) != want
	written := false
	if changed && !q.boolean("dry_run") {
		if err := writeAtomic(full, []byte(want)); err != nil {
			return nil, &Error{Code: "unreadable", Message: "cannot write " + index.Path + ": " + err.Error()}
		}
		written = true
	}
	out := O("path", index.Path, "changed", changed, "written", written)
	if q.boolean("include_content") {
		out = append(out, KV{"content", want})
	}
	return out, nil
}

// writeAtomic writes through a temporary file in the same directory and renames
// it over the target, so a crash leaves the old file or the new one.
func writeAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".akm-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}
	if st, err := os.Stat(path); err == nil {
		os.Chmod(name, st.Mode().Perm())
	} else {
		os.Chmod(name, 0o644)
	}
	if err := os.Rename(name, path); err != nil {
		os.Remove(name)
		return err
	}
	return nil
}
