// Package cli is the command-line binding of the API: flags in, JSON or text out.
package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/api"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// Version of the tool. Releases set it at link time:
// -ldflags "-X github.com/alexsasharegan/akmctl/internal/cli.Version=0.1.0".
var Version = "0.0.0-dev"

// positional names the parameter that takes the positional arguments.
var positional = map[string][]string{"search": {"text"}, "show": {"id"}, "find": {"id"}, "when": {"id"}, "new": {"kind", "slug"}, "close": {"id", "status"}, "claim": {"id"}, "release": {"id"}, "init": {"dir"}, "set": {"id"}, "section": {"id", "name"}}

// Run executes one command line and returns the exit code.
func Run(args []string, stdout, stderr io.Writer, cwd string) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		usage(stderr)
		if len(args) == 0 {
			return 2
		}
		return 0
	}
	op := args[0]
	if op == "--version" {
		op = "version"
	}
	if op == "ui" {
		return runUI(args[1:], stderr, cwd)
	}
	known := api.Params(op)
	if known == nil || !isOp(op) {
		fmt.Fprintf(stderr, "akmctl: unknown command %q (try: akmctl help)\n", op)
		return 2
	}
	format, repoDir, outFile := "table", cwd, ""
	req := api.Request{}
	var pos []string
	yes, noInput := false, false
	rest := args[1:]
	for i := 0; i < len(rest); i++ {
		a := rest[i]
		if !strings.HasPrefix(a, "--") || a == "--" {
			pos = append(pos, a)
			continue
		}
		name, val, hasVal := strings.Cut(a[2:], "=")
		value := func() (string, bool) {
			if hasVal {
				return val, true
			}
			if i+1 < len(rest) {
				i++
				return rest[i], true
			}
			return "", false
		}
		switch name {
		case "field":
			v, _ := value()
			k, raw, ok := strings.Cut(v, "=")
			if !ok || k == "" {
				fmt.Fprintln(stderr, "akmctl: --field needs key=value")
				return 2
			}
			m, _ := req["fields"].(map[string]any)
			if m == nil {
				m = map[string]any{}
			}
			m[k] = parseFieldValue(k, raw)
			req["fields"] = m
			continue
		case "out":
			v, _ := value()
			outFile = v
			continue
		case "yes", "no-input":
			if op != "init" && !(op == "config" && name == "yes") {
				fmt.Fprintf(stderr, "akmctl: --%s does not belong to %s\n", name, op)
				return 2
			}
			if name == "yes" {
				yes = true
			} else {
				noInput = true
			}
			continue
		case "migration":
			v, ok := value()
			if !ok || op != "migrate" {
				fmt.Fprintln(stderr, "akmctl: --migration <file> belongs to migrate")
				return 2
			}
			ms, err := readMigrationFile(v)
			if err != nil {
				fmt.Fprintln(stderr, "akmctl:", err)
				return 2
			}
			prev, _ := req["migrations"].([]any)
			req["migrations"] = append(prev, ms...)
			continue
		case "format":
			v, ok := value()
			if !ok || (v != "json" && v != "table" && v != "ids") {
				fmt.Fprintln(stderr, "akmctl: --format needs json, table or ids")
				return 2
			}
			format = v
			continue
		case "repo":
			v, ok := value()
			if !ok {
				fmt.Fprintln(stderr, "akmctl: --repo needs a directory")
				return 2
			}
			repoDir = v
			continue
		case "api-version":
			if v, ok := value(); !ok || !strings.HasPrefix(api.Version, v) {
				fmt.Fprintf(stderr, "akmctl: API version %q not supported (this tool speaks %s)\n", v, api.Version)
				return 2
			}
			continue
		}
		neg := false
		pname := strings.ReplaceAll(name, "-", "_")
		kind, ok := known[pname]
		if !ok && strings.HasPrefix(name, "no-") {
			pname, neg = strings.ReplaceAll(name[3:], "-", "_"), true
			kind, ok = known[pname]
		}
		if !ok {
			// pass it on: the API rejects it as an unknown parameter, in the JSON form
			req[strings.ReplaceAll(name, "-", "_")] = true
			continue
		}
		switch kind {
		case "bool":
			b := !neg
			if hasVal {
				b = val == "true" || val == "1"
			}
			req[pname] = b
		case "int":
			v, ok := value()
			n, err := strconv.Atoi(v)
			if !ok || err != nil {
				req[pname] = v // the API reports the wrong type
			} else {
				req[pname] = n
			}
		case "strings":
			v, _ := value()
			prev, _ := req[pname].([]string)
			req[pname] = append(prev, splitList(v)...)
		default:
			v, _ := value()
			if pname == "mode" {
				v = strings.ReplaceAll(v, "-", "_") // apply-safe on the command line
			}
			if pname == "text" && v == "-" {
				b, _ := io.ReadAll(os.Stdin) // --text - reads standard input
				v = string(b)
			}
			req[pname] = v
		}
	}
	names := positional[op]
	switch {
	case op == "config":
		// positional words are read by runConfig
	case op == "search":
		if len(pos) > 0 {
			req["text"] = strings.Join(pos, " ")
		}
	case op == "set":
		if len(pos) > 0 {
			req["id"] = pos[0]
			changes := map[string]any{}
			for _, kv := range pos[1:] {
				k, raw, ok := strings.Cut(kv, "=")
				if !ok || k == "" {
					fmt.Fprintf(stderr, "akmctl set: %q is not key=value\n", kv)
					return 2
				}
				changes[k] = parseFieldValue(k, raw)
			}
			if len(changes) > 0 {
				req["changes"] = changes
			}
		}
	default:
		if len(pos) > len(names) {
			fmt.Fprintf(stderr, "akmctl %s: unexpected argument %q\n", op, pos[len(names)])
			return 2
		}
		for i, v := range pos {
			req[names[i]] = v
		}
	}
	info := api.Info{Name: "akmctl", Version: Version}
	if op == "init" {
		// there is no repository yet: the base is the working directory (or
		// --repo), and a dir that is absolute or climbs out of it becomes the base
		base := repoDir
		if d, ok := req["dir"].(string); ok && !relativeDir(d) {
			base = d
			if !filepath.IsAbs(d) {
				base = filepath.Join(repoDir, d)
			}
			base, _ = filepath.Abs(base)
			delete(req, "dir")
		} else if ok {
			req["dir"] = filepath.ToSlash(filepath.Clean(d))
		}
		if !initQuestions(req, base, yes, noInput, stdout, stderr) {
			return 2
		}
		resp, aerr := api.Run(&record.Repo{Root: base}, info, op, req)
		print(stdout, stderr, format, op, resp, aerr == nil)
		if aerr != nil {
			return 2
		}
		return 0
	}
	root, err := record.FindRoot(repoDir)
	if err == nil {
		if flagGiven(args, "--repo") {
			root, err = repoDir, nil
		}
	}
	var repo *record.Repo
	if err == nil {
		repo, err = record.Load(root)
	}
	if err != nil {
		if op == "version" {
			repo, err = &record.Repo{}, nil
		} else {
			resp := api.Failure(op, req, &api.Error{Code: "unreadable", Message: err.Error()})
			print(stdout, stderr, format, op, resp, false)
			return 2
		}
	}
	if op == "config" {
		return runConfig(repo, info, pos, req, yes, format, stdout, stderr)
	}
	resp, aerr := api.Run(repo, info, op, req)
	if h, ok := resp.Get("handoff"); ok && outFile != "" {
		if err := os.WriteFile(outFile, []byte(fmt.Sprint(h)), 0o644); err != nil {
			fmt.Fprintln(stderr, "akmctl:", err)
			return 2
		}
		fmt.Fprintln(stderr, "akmctl: handoff written to", outFile)
	}
	print(stdout, stderr, format, op, resp, aerr == nil)
	if aerr != nil {
		return 2
	}
	if clean, ok := resp.Get("clean"); ok && clean == false {
		return 1
	}
	return 0
}

// readMigrationFile reads a migration file: one migration object, or a list of
// them. The API checks the content.
func readMigrationFile(p string) ([]any, error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, fmt.Errorf("%s: not JSON: %v", p, err)
	}
	switch x := v.(type) {
	case map[string]any:
		return []any{x}, nil
	case []any:
		return x, nil
	}
	return nil, fmt.Errorf("%s: a migration file holds an object or a list of objects", p)
}

// relativeDir reports whether d stays inside the directory it is relative to.
func relativeDir(d string) bool {
	c := filepath.ToSlash(filepath.Clean(d))
	return !filepath.IsAbs(d) && c != ".." && !strings.HasPrefix(c, "../") && !strings.Contains(d, ":")
}

func flagGiven(args []string, name string) bool {
	for _, a := range args {
		if a == name || strings.HasPrefix(a, name+"=") {
			return true
		}
	}
	return false
}

func isOp(op string) bool {
	for _, o := range api.Operations() {
		if o == op {
			return true
		}
	}
	return false
}

func splitList(v string) []string {
	var out []string
	for _, x := range strings.Split(v, ",") {
		if x = strings.TrimSpace(x); x != "" {
			out = append(out, x)
		}
	}
	return out
}

func print(stdout, stderr io.Writer, format, op string, resp api.Obj, ok bool) {
	if format == "json" {
		b, err := api.Marshal(resp)
		if err != nil {
			fmt.Fprintln(stderr, "akmctl:", err)
			return
		}
		fmt.Fprintf(stdout, "%s\n", b)
		return
	}
	if !ok {
		e, _ := resp.Get("error")
		eo := e.(api.Obj)
		code, _ := eo.Get("code")
		msg, _ := eo.Get("message")
		fmt.Fprintf(stderr, "akmctl: %v: %v\n", code, msg)
		return
	}
	if format == "ids" {
		for _, r := range recordsOf(resp) {
			id, _ := r.Get("id")
			fmt.Fprintln(stdout, id)
		}
		return
	}
	text(stdout, op, resp)
}

func recordsOf(resp api.Obj) []api.Obj {
	if v, ok := resp.Get("records"); ok {
		if rs, ok := v.([]api.Obj); ok {
			return rs
		}
	}
	if v, ok := resp.Get("columns"); ok {
		var out []api.Obj
		for _, c := range v.([]api.Obj) {
			cv, _ := c.Get("records")
			out = append(out, cv.([]api.Obj)...)
		}
		return out
	}
	return nil
}

func field(r api.Obj, k string) string {
	f, _ := r.Get("fields")
	if fo, ok := f.(api.Obj); ok {
		if v, ok := fo.Get(k); ok && v != nil {
			return fmt.Sprint(v)
		}
	}
	return ""
}

func row(w io.Writer, r api.Obj) {
	id, _ := r.Get("id")
	date, _ := r.Get("id_date")
	status := field(r, "status")
	fmt.Fprintf(w, "%v\t%v\t%s\t%s\n", id, date, status, field(r, "title"))
}

func text(w io.Writer, op string, resp api.Obj) {
	switch op {
	case "ls", "search":
		for _, r := range recordsOf(resp) {
			row(w, r)
		}
		total, _ := resp.Get("total")
		count, _ := resp.Get("count")
		if total != count {
			fmt.Fprintf(w, "(%v of %v)\n", count, total)
		}
	case "board":
		v, _ := resp.Get("columns")
		for _, c := range v.([]api.Obj) {
			st, _ := c.Get("status")
			n, _ := c.Get("count")
			fmt.Fprintf(w, "== %v (%v)\n", st, n)
			rs, _ := c.Get("records")
			for _, r := range rs.([]api.Obj) {
				row(w, r)
			}
		}
	case "show":
		v, _ := resp.Get("record")
		r := v.(api.Obj)
		f, _ := r.Get("fields")
		for _, kv := range f.(api.Obj) {
			if kv.V != nil {
				fmt.Fprintf(w, "%s: %v\n", kv.K, kv.V)
			}
		}
		if b, ok := resp.Get("body"); ok {
			fmt.Fprintf(w, "\n%s", b)
		}
		if s, ok := resp.Get("section"); ok {
			t, _ := s.(api.Obj).Get("text")
			fmt.Fprintf(w, "\n%v\n", t)
		}
	case "migrate":
		st, _ := resp.Get("state")
		from, _ := resp.Get("from")
		to, _ := resp.Get("to")
		fmt.Fprintf(w, "layout %v -> %v: %v\n", from, to, st)
		ss, _ := resp.Get("steps")
		for _, x := range ss.([]api.Obj) {
			id, _ := x.Get("id")
			cl, _ := x.Get("class")
			status, _ := x.Get("status")
			files, _ := x.Get("files")
			fmt.Fprintf(w, "  %v [%v] %v %v\n", id, cl, status, files)
		}
		if h, ok := resp.Get("handoff"); ok {
			fmt.Fprintf(w, "\n%v", h)
		}
	case "check":
		fs, _ := resp.Get("findings")
		for _, f := range fs.([]api.Obj) {
			file, _ := f.Get("file")
			line, _ := f.Get("line")
			rule, _ := f.Get("rule")
			sev, _ := f.Get("severity")
			msg, _ := f.Get("message")
			loc := ""
			if file != nil {
				loc = fmt.Sprint(file)
				if line != nil {
					loc += fmt.Sprintf(":%v", line)
				}
				loc += ": "
			}
			fmt.Fprintf(w, "%s[%v] %v %v\n", loc, sev, rule, msg)
		}
		e, _ := resp.Get("errors")
		wn, _ := resp.Get("warnings")
		fmt.Fprintf(w, "%v errors, %v warnings\n", e, wn)
	case "init":
		dir, _ := resp.Get("dir")
		wr, _ := resp.Get("written")
		sk, _ := resp.Get("skipped")
		docs, _ := resp.Get("docs_dir")
		fmt.Fprintf(w, "initialised %v (records in %v/)\n", dir, docs)
		for _, p := range wr.([]string) {
			fmt.Fprintln(w, "  wrote   ", p)
		}
		for _, p := range sk.([]string) {
			fmt.Fprintln(w, "  skipped ", p, "(exists)")
		}
		switch g, _ := resp.Get("git"); g {
		case "initialised":
			fmt.Fprintln(w, "  git init done")
		case "failed":
			ge, _ := resp.Get("git_error")
			fmt.Fprintln(w, "  git init failed:", ge)
		}
		if f, _ := resp.Get("git_commits_forced_off"); f == true {
			fmt.Fprintln(w, "not a git repository: git commits stay off (run git init, then: akmctl config git-commits on)")
		}
	case "config":
		configText(w, resp)
	case "index":
		p, _ := resp.Get("path")
		ch, _ := resp.Get("changed")
		wr, _ := resp.Get("written")
		fmt.Fprintf(w, "%v: changed=%v written=%v\n", p, ch, wr)
	case "new", "set", "section", "close", "claim", "release":
		v, _ := resp.Get("record")
		r := v.(api.Obj)
		row(w, r)
		if d, ok := resp.Get("diff"); ok {
			writeDiff(w, d.(api.Obj))
		}
		if p, ok := resp.Get("path"); ok {
			fmt.Fprintln(w, p)
		}
		if sv, ok := resp.Get("similar"); ok {
			fmt.Fprintln(w, "similar open items (add to one instead of filing, if it fits):")
			for _, o := range sv.([]api.Obj) {
				id, _ := o.Get("id")
				st, _ := o.Get("status")
				ti, _ := o.Get("title")
				fmt.Fprintf(w, "  %v\t%v\t%v\n", id, st, ti)
			}
		}
		if p, ok := resp.Get("moved_to"); ok {
			fmt.Fprintln(w, "moved to", p)
		}
		if cm, ok := resp.Get("commit"); ok {
			msg, _ := cm.(api.Obj).Get("message")
			ps, _ := cm.(api.Obj).Get("paths")
			fmt.Fprintln(w, "commit message (use as given; stage these paths and the code):")
			for _, l := range strings.Split(msg.(string), "\n") {
				fmt.Fprintln(w, "  "+l)
			}
			for _, p := range ps.([]any) {
				fmt.Fprintln(w, "  stage:", p)
			}
		}
		if dv, ok := resp.Get("dependents"); ok {
			fmt.Fprintln(w, "open items that were waiting on this one (unblock or re-point each):")
			for _, o := range dv.([]api.Obj) {
				id, _ := o.Get("id")
				st, _ := o.Get("status")
				ti, _ := o.Get("title")
				fmt.Fprintf(w, "  %v\t%v\t%v\n", id, st, ti)
			}
		}
		if c, ok := resp.Get("changed"); ok {
			fmt.Fprintln(w, "changed:", c)
		}
	case "find":
		p, _ := resp.Get("path")
		fmt.Fprintln(w, p)
	case "when":
		d, _ := resp.Get("date")
		wd, _ := resp.Get("weekday")
		age, _ := resp.Get("age_days")
		fmt.Fprintf(w, "%v %v, %v days ago\n", d, wd, age)
	case "version":
		t, _ := resp.Get("tool")
		n, _ := t.(api.Obj).Get("name")
		v, _ := t.(api.Obj).Get("version")
		fmt.Fprintf(w, "%v %v (API %v)\n", n, v, api.Version)
	}
}

// writeDiff prints the diff of a dry run in unified form.
func writeDiff(w io.Writer, d api.Obj) {
	p, _ := d.Get("path")
	np, _ := d.Get("new_path")
	fmt.Fprintf(w, "--- %v\n+++ %v\n", p, np)
	hs, _ := d.Get("hunks")
	for _, h := range hs.([]api.Obj) {
		hd, _ := h.Get("header")
		fmt.Fprintln(w, hd)
		ls, _ := h.Get("lines")
		for _, l := range ls.([]api.Obj) {
			op, _ := l.Get("op")
			t, _ := l.Get("text")
			mark := map[any]string{"ctx": " ", "add": "+", "rem": "-"}[op]
			fmt.Fprintf(w, "%s%v\n", mark, t)
		}
	}
}

func usage(w io.Writer) {
	ops := append(api.Operations(), "ui")
	sort.Strings(ops)
	fmt.Fprintf(w, "usage: akmctl <command> [flags]\n\ncommands: %s\n\nglobal flags: --repo <dir>  --format json|table|ids\nmigrate flags: --mode report|plan|apply-safe|handoff  --target <version>  --acknowledge-backup  --migration <file>  --out <file>\ninit: akmctl init [dir] [--docs-dir <name>] [--wikilinks|--no-wikilinks] [--git-commits|--no-git-commits] [--git-init|--no-git-init] [--yes|--no-input]  (writes missing files, never overwrites; asks on a terminal)\nconfig: akmctl config  |  akmctl config git-commits on|off [--yes] [--dry-run]\nui flags: --port <n> (default: any free port)  --open\n", strings.Join(ops, ", "))
}

var _ = os.Stdout

// parseFieldValue reads the value of key=value on the command line: null for
// "null" or nothing, a list for the list fields, an integer for priority.
func parseFieldValue(key, raw string) any {
	if raw == "" || raw == "null" {
		return nil
	}
	switch key {
	case "blocked_by", "decisions", "scope":
		var out []any
		for _, x := range splitList(raw) {
			out = append(out, x)
		}
		if out == nil {
			out = []any{}
		}
		return out
	case "priority":
		if n, err := strconv.Atoi(raw); err == nil {
			return n
		}
	case "keep":
		switch raw {
		case "true":
			return true
		case "false":
			return false
		}
	}
	return raw
}
