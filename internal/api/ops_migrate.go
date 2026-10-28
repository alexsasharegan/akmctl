package api

import (
	"errors"
	"time"

	"github.com/alexsasharegan/akmctl/internal/check"
	"github.com/alexsasharegan/akmctl/internal/frontmatter"
	"github.com/alexsasharegan/akmctl/internal/fsx"
	"github.com/alexsasharegan/akmctl/internal/migrate"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// repoLayout is the layout version recorded in the manifest, "0.0.0" when none.
func repoLayout(repo *record.Repo) string {
	if v, ok := repo.Layout["layout"]; ok && v.Kind == frontmatter.String && v.Str != "" {
		return v.Str
	}
	return "0.0.0"
}

func cmpVersion(a, b string) int {
	var x, y [3]int
	for i, s := range []string{a, b} {
		var p [3]int
		n := 0
		for _, part := range splitDots(s) {
			if n < 3 {
				p[n] = atoi(part)
				n++
			}
		}
		if i == 0 {
			x = p
		} else {
			y = p
		}
	}
	for i := 0; i < 3; i++ {
		if x[i] != y[i] {
			if x[i] < y[i] {
				return -1
			}
			return 1
		}
	}
	return 0
}

func splitDots(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r == '.' {
			out = append(out, cur)
			cur = ""
		} else {
			cur += string(r)
		}
	}
	return append(out, cur)
}

func atoi(s string) int {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			break
		}
		n = n*10 + int(r-'0')
	}
	return n
}

// classify gives a finding its class for a target version: auto when an auto
// step fixes its rule, manual when a manual step does or when it is an error,
// info for a warning.
func classify(f check.Finding, steps []migrate.Step) string {
	class := ""
	for _, st := range steps {
		for _, r := range st.Fixes {
			if r == f.Rule {
				if st.Class == "auto" {
					return "auto"
				}
				class = st.Class
			}
		}
	}
	if class != "" {
		return class
	}
	if f.Severity == "E" {
		return "manual"
	}
	return "info"
}

func findingObj(f check.Finding, class any) Obj {
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
	return O("rule", f.Rule, "severity", f.Severity, "file", file, "line", line, "id", id, "message", f.Message, "class", class)
}

func stepObj(r migrate.Result) Obj {
	files := r.Files
	if files == nil {
		files = []string{}
	}
	return O("id", r.Step.ID, "kind", r.Step.Kind, "class", r.Step.Class, "files", files, "status", r.Status)
}

func opMigrate(repo *record.Repo, q Request) (Obj, *Error) {
	mode := q.str("mode")
	if mode == "" {
		mode = "report"
	}
	to := q.str("target")
	if to == "" {
		to = check.LayoutVersion
	}
	from := repoLayout(repo)
	chain, found := migrate.Chain(from, to)
	steps := migrate.Steps(chain)
	state := "current"
	switch {
	case from == to:
	case found:
		state = "behind"
	case cmpVersion(from, to) > 0:
		state = "ahead"
	default:
		state = "unsupported"
	}
	// the target's rules: this release has one rule set, so only its own version can be checked
	var findings []check.Finding
	if to == check.LayoutVersion {
		findings = check.Run(repo, check.Options{}).Findings
	}
	asObjs := func() []Obj {
		out := make([]Obj, 0, len(findings))
		for _, f := range findings {
			out = append(out, findingObj(f, classify(f, steps)))
		}
		return out
	}
	finish := func(results []migrate.Result, state string, extra ...KV) Obj {
		so := make([]Obj, 0, len(results))
		pending := false
		for _, r := range results {
			so = append(so, stepObj(r))
			if r.Status == "pending" && r.Step.Class != "info" {
				pending = true
			}
		}
		o := O("mode", mode, "from", from, "to", to, "state", state, "steps", so, "findings", asObjs(), "clean", !pending && state == "current")
		return append(o, extra...)
	}
	switch mode {
	case "report":
		res := make([]migrate.Result, 0, len(steps))
		for _, s := range steps {
			res = append(res, migrate.Result{Step: s, Files: []string{}, Status: "pending"})
		}
		return finish(res, state), nil
	case "plan":
		return finish(migrate.Plan(repo, steps), state), nil
	case "handoff":
		res := migrate.Plan(repo, steps)
		var msgs []string
		for _, f := range findings {
			if f.Severity == "E" {
				msgs = append(msgs, f.Rule+" "+f.File+": "+f.Message)
			}
		}
		return finish(res, state, KV{"handoff", migrate.Handoff(from, to, res, msgs)}), nil
	}
	// apply_safe
	if !q.boolean("acknowledge_backup") {
		return nil, refused("no_backup_ack", "apply_safe needs acknowledge_backup: true (a backup tag is made, and data loss is the maintainer's risk)")
	}
	if state == "current" {
		return finish(nil, state), nil
	}
	if state != "behind" {
		return nil, refused("no_path", "no migration path from layout %s to %s in this release", from, to)
	}
	tag, branch, err := migrate.Preflight(repo.Root, from, to)
	if err != nil {
		var rf *migrate.Refusal
		if errors.As(err, &rf) {
			return nil, refused(rf.Reason, "%s", rf.Message)
		}
		return nil, unreadable(err)
	}
	unlock, lerr := fsx.LockRepo(repo.Root)
	if lerr != nil {
		return nil, unreadable(lerr)
	}
	defer unlock()
	if err := migrate.Backup(repo.Root, tag, branch); err != nil {
		return nil, unreadable(err)
	}
	if err := repo.Reload(); err != nil {
		return nil, unreadable(err)
	}
	results, err := migrate.Apply(repo, steps)
	if err != nil {
		return nil, &Error{"unreadable", err.Error() + " (the backup tag " + tag + " holds the state before the migration)", O("tag", tag, "branch", branch)}
	}
	if err := migrate.AppendLog(repo.Root, from, to, results, clock(q)); err != nil {
		return nil, unreadable(err)
	}
	manual := false
	for _, r := range results {
		if r.Status == "pending" && r.Step.Class == "manual" {
			manual = true
		}
	}
	newState := "behind"
	if !manual {
		if err := migrate.SetLayout(repo.Root, to); err != nil {
			return nil, unreadable(err)
		}
		newState = "current"
	}
	if err := repo.Reload(); err != nil {
		return nil, unreadable(err)
	}
	findings = check.Run(repo, check.Options{}).Findings
	return finish(results, newState), nil
}

var _ = time.Now
