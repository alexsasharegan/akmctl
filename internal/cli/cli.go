// Package cli is the command-line binding of the API: flags in, JSON or text out.
package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/api"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// Version of the tool.
const Version = "0.0.0-dev"

// positional names the parameter that takes the positional arguments.
var positional = map[string]string{"search": "text", "show": "id", "find": "id", "when": "id"}

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
	known := api.Params(op)
	if known == nil || !isOp(op) {
		fmt.Fprintf(stderr, "akmctl: unknown command %q (try: akmctl help)\n", op)
		return 2
	}
	format, repoDir := "table", cwd
	req := api.Request{}
	var pos []string
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
			req[pname] = v
		}
	}
	if pn := positional[op]; pn != "" && len(pos) > 0 {
		req[pn] = strings.Join(pos, " ")
	} else if len(pos) > 0 {
		fmt.Fprintf(stderr, "akmctl %s: unexpected argument %q\n", op, pos[0])
		return 2
	}
	info := api.Info{Name: "akmctl", Version: Version}
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
	resp, aerr := api.Run(repo, info, op, req)
	print(stdout, stderr, format, op, resp, aerr == nil)
	if aerr != nil {
		return 2
	}
	return 0
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
	case "find":
		p, _ := resp.Get("path")
		fmt.Fprintln(w, p)
	case "when":
		d, _ := resp.Get("date")
		wd, _ := resp.Get("weekday")
		age, _ := resp.Get("age_days")
		fmt.Fprintf(w, "%v (UTC) %v, %v days ago\n", d, wd, age)
	case "version":
		t, _ := resp.Get("tool")
		n, _ := t.(api.Obj).Get("name")
		v, _ := t.(api.Obj).Get("version")
		fmt.Fprintf(w, "%v %v (API %v)\n", n, v, api.Version)
	}
}

func usage(w io.Writer) {
	ops := api.Operations()
	sort.Strings(ops)
	fmt.Fprintf(w, "usage: akmctl <command> [flags]\n\ncommands: %s\n\nglobal flags: --repo <dir>  --format json|table|ids\n", strings.Join(ops, ", "))
}

var _ = os.Stdout
