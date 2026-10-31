package cli

import (
	"fmt"
	"io"

	"github.com/alexsasharegan/akmctl/internal/api"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// runSkills reports the managed skills (status) or installs and updates them.
// A write is planned first and shown; it goes on with --yes, or after a yes on
// a terminal. Edited skills need a second yes.
func runSkills(repo *record.Repo, info api.Info, pos []string, req api.Request, yes bool, format string, stdout, stderr io.Writer) int {
	if len(pos) > 1 {
		fmt.Fprintln(stderr, "akmctl skills: usage: akmctl skills [status|install|update] [--yes] [--dry-run]")
		return 2
	}
	if len(pos) == 1 {
		req["mode"] = pos[0]
	}
	run := func(r api.Request) (api.Obj, bool) {
		resp, aerr := api.Run(repo, info, "skills", r)
		return resp, aerr == nil
	}
	mode, _ := req["mode"].(string)
	if mode == "" || mode == "status" || req["dry_run"] == true {
		resp, ok := run(req)
		print(stdout, stderr, format, "skills", resp, ok)
		if !ok {
			return 2
		}
		return 0
	}
	pre := api.Request{}
	for k, v := range req {
		pre[k] = v
	}
	pre["dry_run"] = true
	resp, ok := run(pre)
	if !ok {
		print(stdout, stderr, format, "skills", resp, false)
		return 2
	}
	plan, _ := resp.Get("plan")
	actions, _ := plan.([]api.Obj)
	if len(actions) == 0 {
		print(stdout, stderr, format, "skills", resp, true)
		return 0
	}
	edited := false
	for _, a := range actions {
		if b, _ := a.Get("needs_confirm_edited"); b == true {
			edited = true
		}
	}
	skillsText(stderr, resp)
	req["confirm"] = true
	if !yes {
		if !interactive(stdout) {
			fmt.Fprintln(stderr, "akmctl skills: a write needs a terminal to ask on, or --yes")
			return 2
		}
		p := newPrompter(stdinReader, stderr)
		if !p.yesNo("Apply these changes?", false) {
			fmt.Fprintln(stderr, "nothing changed")
			return 1
		}
		if edited {
			req["confirm_edited"] = p.yesNo("Some skills were edited or have no marker. Overwrite or remove them anyway?", false)
		}
	}
	resp, ok = run(req)
	print(stdout, stderr, format, "skills", resp, ok)
	if !ok {
		return 2
	}
	return 0
}

// skillsText prints a skills response for people.
func skillsText(w io.Writer, resp api.Obj) {
	ss, _ := resp.Get("skills")
	for _, s := range ss.([]api.Obj) {
		n, _ := s.Get("name")
		st, _ := s.Get("state")
		pth, _ := s.Get("path")
		fmt.Fprintf(w, "%-10v %-24v %v\n", st, n, pth)
		if note, ok := s.Get("note"); ok {
			fmt.Fprintf(w, "           %v\n", note)
		}
	}
	plan, ok := resp.Get("plan")
	if !ok {
		return
	}
	for _, a := range plan.([]api.Obj) {
		k, _ := a.Get("action")
		pth, _ := a.Get("path")
		flag := ""
		if b, _ := a.Get("needs_confirm_edited"); b == true {
			flag = "  (edited: needs a second yes)"
		}
		fmt.Fprintf(w, "  %-8v %v%s\n", k, pth, flag)
		if d, has := a.Get("diff"); has {
			writeDiff(w, d.(api.Obj))
		}
	}
	if sk, ok := resp.Get("skipped"); ok {
		for _, o := range sk.([]api.Obj) {
			n, _ := o.Get("skill")
			r, _ := o.Get("reason")
			fmt.Fprintf(w, "  skipped  %v (%v)\n", n, r)
		}
	}
	switch ap, _ := resp.Get("applied"); {
	case ap == true:
		fmt.Fprintln(w, "applied")
	case len(plan.([]api.Obj)) > 0:
		fmt.Fprintln(w, "nothing written (a plan only)")
	}
}
