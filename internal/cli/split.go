package cli

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/api"
	"github.com/alexsasharegan/akmctl/internal/record"
)

var reSplitSlug = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// runSplit is `akmctl split <id> <slug>...`: one `new` per slug, each child a raw work item
// under <id>. It is a command-line convenience over the `new` operation, not an operation of
// its own (work item 261031-ph53). It writes only the new children; the parent is left as it
// is, and the command prints what the parent still needs.
func runSplit(args []string, stdout, stderr io.Writer, cwd string) int {
	repoDir, given, format := cwd, false, "table"
	fields := map[string]any{}
	asOf := ""
	var pos []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if !strings.HasPrefix(a, "--") {
			pos = append(pos, a)
			continue
		}
		name, val, hasVal := strings.Cut(a[2:], "=")
		value := func() (string, bool) {
			if hasVal {
				return val, true
			}
			if i+1 < len(args) {
				i++
				return args[i], true
			}
			return "", false
		}
		v, ok := value()
		if !ok {
			fmt.Fprintf(stderr, "akmctl split: --%s needs a value\n", name)
			return 2
		}
		switch name {
		case "repo":
			repoDir, given = v, true
		case "format":
			if v != "json" && v != "table" && v != "ids" {
				fmt.Fprintln(stderr, "akmctl: --format needs json, table or ids")
				return 2
			}
			format = v
		case "kind", "priority":
			fields[name] = parseFieldValue(name, v)
		case "as-of":
			asOf = v
		default:
			fmt.Fprintf(stderr, "akmctl split: unknown flag --%s (usage: akmctl split <id> <slug>... [--kind <kind>] [--priority <n>])\n", name)
			return 2
		}
	}
	if len(pos) < 2 {
		fmt.Fprintln(stderr, "akmctl split: needs the parent's id and at least one slug (usage: akmctl split <id> <slug>...)")
		return 2
	}
	parentID, slugs := pos[0], pos[1:]
	seen := map[string]bool{}
	for _, s := range slugs {
		if !reSplitSlug.MatchString(s) || len(s) > 60 {
			fmt.Fprintf(stderr, "akmctl split: %q is not a slug (lowercase words joined by hyphens, at most 60 characters); nothing was created\n", s)
			return 2
		}
		if seen[s] {
			fmt.Fprintf(stderr, "akmctl split: the slug %q is given twice; nothing was created\n", s)
			return 2
		}
		seen[s] = true
	}
	root, err := record.FindRoot(repoDir)
	if given {
		root, err = repoDir, nil
		if e := insideRepository("new", repoDir); e != nil {
			err = e
		}
	}
	var repo *record.Repo
	if err == nil {
		repo, err = record.Load(root)
	}
	if err != nil {
		fmt.Fprintln(stderr, "akmctl split:", err)
		return 2
	}
	info := api.Info{Name: "akmctl", Version: Version}
	show, aerr := api.Run(repo, info, "show", api.Request{"id": parentID})
	if aerr != nil {
		print(stdout, stderr, "table", "show", show, false)
		return 2
	}
	rv, _ := show.Get("record")
	rec := rv.(api.Obj)
	if t, _ := rec.Get("type"); t != record.Work {
		fmt.Fprintf(stderr, "akmctl split: %s is not a work item; nothing was created\n", parentID)
		return 2
	}
	pf, _ := rec.Get("fields")
	get := func(k string) any { v, _ := pf.(api.Obj).Get(k); return v }
	if st := fmt.Sprint(get("status")); st == "done" || st == "dropped" {
		fmt.Fprintf(stderr, "akmctl split: %s is %s; nothing was created\n", parentID, st)
		return 2
	}
	// children sort with the parent's work unless a priority is given
	if _, ok := fields["priority"]; !ok {
		if p := get("priority"); p != nil {
			fields["priority"] = p
		}
	}
	fields["parent"] = parentID
	var children []api.Obj
	for n, slug := range slugs {
		req := api.Request{"kind": record.Work, "slug": slug, "fields": copyFields(fields)}
		if asOf != "" {
			req["as_of"] = asOf
		}
		resp, aerr := api.Run(repo, info, "new", req)
		if aerr != nil {
			print(stdout, stderr, "table", "new", resp, false)
			if n > 0 {
				fmt.Fprintf(stderr, "akmctl split: stopped after %d of %d children (they are filed; commit them, then split the rest)\n", n, len(slugs))
			}
			return 2
		}
		children = append(children, resp)
	}
	left := parentLeft(parentID, get("status"), get("claimed_by"))
	if format == "json" {
		cs := make([]any, len(children))
		for i, c := range children {
			cs[i] = c
		}
		out := api.Obj{{K: "operation", V: "split"}, {K: "ok", V: true}, {K: "parent", V: parentID}, {K: "children", V: cs}, {K: "parent_left", V: left}}
		b, merr := api.Marshal(out)
		if merr != nil {
			fmt.Fprintln(stderr, "akmctl:", merr)
			return 2
		}
		fmt.Fprintf(stdout, "%s\n", b)
		return 0
	}
	for _, c := range children {
		if format == "ids" {
			rv, _ := c.Get("record")
			id, _ := rv.(api.Obj).Get("id")
			fmt.Fprintln(stdout, id)
			continue
		}
		text(stdout, "new", c)
	}
	if format == "table" {
		fmt.Fprintf(stdout, "left for the parent %s (split changed none of it):\n", parentID)
		for i, l := range left {
			fmt.Fprintf(stdout, "  %s. %s\n", strconv.Itoa(i+1), l)
		}
	}
	return 0
}

func copyFields(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// parentLeft lists what a split leaves for the parent. It edits none of it.
func parentLeft(id string, status, claimedBy any) []string {
	out := []string{
		"write its ## Goal as a list of the children (plain bullets with links, in dependency order) and move each old box to the child that owns it: akmctl section " + id + " Goal",
		"write a ## Next that points at the first child: akmctl section " + id + " Next",
	}
	if claimedBy != nil {
		out = append(out, "release the claim; work the children, not the parent: akmctl release "+id)
	}
	if fmt.Sprint(status) != "ready" {
		out = append(out, "set it ready (a parent waits for its children; an unclaimed active item makes the next agent ask): akmctl set "+id+" status=ready")
	}
	return out
}
