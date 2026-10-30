package cli

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/api"
	"github.com/alexsasharegan/akmctl/internal/initrepo"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// initQuestions completes the init request. On a terminal (and without --yes or
// --no-input) it asks for what the flags left open: where the records go,
// wikilinks, git init (outside a repository) and the commit skill. Otherwise it
// asks nothing and the defaults apply. It reports false after printing why the
// request cannot go on.
func initQuestions(req api.Request, base string, yes, noInput bool, stdout, stderr io.Writer) bool {
	target := base
	if d, ok := req["dir"].(string); ok {
		target = filepath.Join(base, filepath.FromSlash(d))
	}
	if ask := !yes && !noInput && interactive(stdout); ask {
		askInit(newPrompter(stdinReader, stderr), req, record.InGitRepo(target))
	}
	return true
}

// askInit asks for what req leaves open; the answers are added to req.
func askInit(p *prompter, req api.Request, inRepo bool) {
	if _, ok := req["docs_dir"]; !ok {
		req["docs_dir"] = p.docsDir()
	}
	if _, ok := req["wikilinks"]; !ok {
		req["wikilinks"] = p.yesNo("Use [[wikilinks]] between records?", false)
	}
	if _, ok := req["git_init"]; !ok && !inRepo {
		req["git_init"] = p.yesNo("Run git init?", true)
	}
	willBeRepo := inRepo || req["git_init"] == true
	if _, ok := req["git_commits"]; !ok && willBeRepo {
		req["git_commits"] = p.yesNo("Install the akm-commit skill (agents commit after each step)?", false)
	}
}

// runConfig lists the settings or changes one. A removal is shown first and
// confirmed on a terminal, or with --yes; without either it is refused.
func runConfig(repo *record.Repo, info api.Info, pos []string, req api.Request, yes bool, format string, stdout, stderr io.Writer) int {
	fail := func(resp api.Obj) int {
		print(stdout, stderr, format, "config", resp, false)
		return 2
	}
	if len(pos) > 0 {
		if pos[0] != "git-commits" && pos[0] != "git_commits" || len(pos) > 2 {
			fmt.Fprintln(stderr, "akmctl config: usage: akmctl config [git-commits on|off] [--yes] [--dry-run]")
			return 2
		}
		req["key"] = "git_commits"
		if len(pos) == 2 {
			switch strings.ToLower(pos[1]) {
			case "on", "true":
				req["value"] = true
			case "off", "false":
				req["value"] = false
			default:
				fmt.Fprintf(stderr, "akmctl config: %q is not on or off\n", pos[1])
				return 2
			}
		}
	}
	if _, set := req["value"]; !set || req["dry_run"] == true {
		resp, aerr := api.Run(repo, info, "config", req)
		print(stdout, stderr, format, "config", resp, aerr == nil)
		if aerr != nil {
			return 2
		}
		return 0
	}
	// the plan first: it carries the conflicts and refusals, and what a removal
	// would delete
	pre := api.Request{}
	for k, v := range req {
		pre[k] = v
	}
	pre["dry_run"] = true
	resp, aerr := api.Run(repo, info, "config", pre)
	if aerr != nil {
		return fail(resp)
	}
	if nc, _ := resp.Get("needs_confirmation"); nc == true {
		plan, _ := resp.Get("plan")
		rm, _ := plan.(api.Obj).Get("remove")
		edits, _ := resp.Get("edits")
		edited := len(edits.([]api.Obj)) > 0
		fmt.Fprintf(stderr, "This removes %s (%d files):\n", initrepo.SkillDir(initrepo.CommitSkill), len(rm.([]string)))
		for _, f := range rm.([]string) {
			fmt.Fprintln(stderr, "  remove", f)
		}
		req["confirm"] = true
		if edited {
			fmt.Fprintln(stderr, "These files differ from the copy the tool ships (shipped text to installed text):")
			for _, d := range edits.([]api.Obj) {
				writeDiff(stderr, d)
			}
			req["confirm_edited"] = true
		}
		if !yes {
			if !interactive(stdout) {
				fmt.Fprintln(stderr, "akmctl config: removal needs a terminal to ask on, or --yes")
				return 2
			}
			p := newPrompter(stdinReader, stderr)
			if !p.yesNo("Remove them and turn git_commits off?", false) ||
				(edited && !p.yesNo("The installed files were edited. Remove them anyway?", false)) {
				fmt.Fprintln(stderr, "nothing changed")
				return 1
			}
		}
	}
	resp, aerr = api.Run(repo, info, "config", req)
	print(stdout, stderr, format, "config", resp, aerr == nil)
	if aerr != nil {
		return 2
	}
	return 0
}

// configText prints a config response for people.
func configText(w io.Writer, resp api.Obj) {
	if _, set := resp.Get("plan"); set {
		v, _ := resp.Get("value")
		ch, _ := resp.Get("changed")
		ap, _ := resp.Get("applied")
		state := map[any]string{true: "on", false: "off"}[v]
		switch {
		case ch != true:
			fmt.Fprintf(w, "git_commits is already %s: nothing to do\n", state)
		case ap != true:
			fmt.Fprintf(w, "git_commits would be turned %s:\n", state)
		default:
			fmt.Fprintf(w, "git_commits turned %s:\n", state)
		}
		plan, _ := resp.Get("plan")
		for _, k := range []string{"write", "remove", "modify"} {
			l, _ := plan.(api.Obj).Get(k)
			for _, p := range l.([]string) {
				fmt.Fprintf(w, "  %-7s %s\n", k, p)
			}
		}
		if nc, _ := resp.Get("needs_confirmation"); nc == true && ap != true {
			fmt.Fprintln(w, "  (removing needs confirmation)")
		}
		return
	}
	ss, _ := resp.Get("settings")
	for _, s := range ss.([]api.Obj) {
		k, _ := s.Get("key")
		v, _ := s.Get("value")
		d, _ := s.Get("description")
		fmt.Fprintf(w, "%v = %v\n    %v\n", k, v, d)
	}
}
