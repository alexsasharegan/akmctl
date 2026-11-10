package api

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alexsasharegan/akmctl/internal/check"
	"github.com/alexsasharegan/akmctl/internal/fsx"
	"github.com/alexsasharegan/akmctl/internal/index"
	"github.com/alexsasharegan/akmctl/internal/migrate"
	"github.com/alexsasharegan/akmctl/internal/record"
)

func opCheck(repo *record.Repo, q Request) (Obj, *Error) {
	if t := q.str("target"); t != "" && t != check.LayoutVersion {
		// any layout the embedded migrations reach from the repository's own, or the
		// repository's own: the findings are the tool's, classified by the chain's steps
		if from := repoLayout(repo); t != from {
			if _, ok := migrate.Chain(from, t); !ok {
				reach := []string{check.LayoutVersion}
				for _, v := range migrate.Reachable(from) {
					if v != check.LayoutVersion {
						reach = append(reach, v)
					}
				}
				return nil, invalid("target", "target: layout version %s is not reachable from this repository's layout %s; supported targets: %s", t, from, strings.Join(reach, ", "))
			}
		}
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
	var steps []migrate.Step
	if t := q.str("target"); t != "" {
		chain, _ := migrate.Chain(repoLayout(repo), t)
		steps = migrate.Steps(chain)
	}
	errs, warns := 0, 0
	findings := make([]Obj, 0, len(res.Findings))
	for _, f := range res.Findings {
		if f.Severity == "E" {
			errs++
		} else {
			warns++
		}
		var class any
		if q.str("target") != "" {
			class = classify(f, steps)
		}
		findings = append(findings, findingObj(f, class))
	}
	rec := O("state", res.State)
	if len(res.Steps) > 0 {
		rec = append(rec, KV{"steps", res.Steps})
	}
	return O("clean", errs == 0, "errors", errs, "warnings", warns, "count", len(findings), "findings", findings, "reconcile", rec), nil
}

func opIndex(repo *record.Repo, q Request) (Obj, *Error) {
	want := index.Generate(repo)
	full := filepath.Join(repo.Root, filepath.FromSlash(repo.ADRIndexPath()))
	cur, err := os.ReadFile(full)
	changed := err != nil || string(cur) != want
	written := false
	if changed && !q.boolean("dry_run") {
		unlock, lerr := fsx.LockRepo(repo.Root, repo.MetaDir())
		if lerr != nil {
			return nil, unreadable(lerr)
		}
		defer unlock()
		if err := fsx.WriteAtomic(full, []byte(want)); err != nil {
			return nil, &Error{Code: "unreadable", Message: "cannot write " + repo.ADRIndexPath() + ": " + err.Error()}
		}
		written = true
	}
	out := O("path", repo.ADRIndexPath(), "changed", changed, "written", written)
	if q.boolean("include_content") {
		out = append(out, KV{"content", want})
	}
	return out, nil
}
