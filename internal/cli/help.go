package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/api"
)

// examples holds one usage line per command. A command with no entry fails
// TestEveryCommandHasHelp, so a new command cannot ship without one.
var examples = map[string]string{
	"board":   "akmctl board",
	"check":   "akmctl check",
	"claim":   "akmctl claim 261101-m3pv --as claude-sonnet-5-5",
	"close":   `akmctl close 261101-m3pv done --outcome "What was done and why."`,
	"commits": "akmctl commits",
	"history": "akmctl history --range origin/main..HEAD",
	"config":  "akmctl config git-commits on --dry-run",
	"find":    "akmctl find 261101-m3pv",
	"index":   "akmctl index",
	"init":    "akmctl init --yes",
	"log":     `akmctl log 261101-m3pv "Found the flag parser in cli.go."`,
	"ls":      "akmctl ls --status ready,active --fields title",
	"migrate": "akmctl migrate --mode report",
	"new":     "akmctl new work fix-the-parser",
	"release": "akmctl release 261101-m3pv",
	"remove":  "akmctl remove 261101-m3pv 2",
	"search":  "akmctl search toast duration",
	"section": `akmctl section 261101-m3pv Next --text "The next step."`,
	"set":     "akmctl set 261101-m3pv priority=2",
	"show":    "akmctl show 261101-m3pv",
	"skills":  "akmctl skills status",
	"tick":    `akmctl tick 261101-m3pv 1 --evidence "go test passes"`,
	"ui":      "akmctl ui --open",
	"version": "akmctl version",
	"when":    "akmctl when 261101-m3pv",
}

// commandNames lists every command help covers: the operations and ui.
func commandNames() []string {
	names := append(api.Operations(), "ui")
	sort.Strings(names)
	return names
}

// commandUsage prints the flags and an example for one command.
func commandUsage(w io.Writer, op string) {
	if op == "ui" {
		fmt.Fprintf(w, "usage: akmctl ui [flags]\n\nflags:\n  --port <n>  port to listen on (default: any free port)\n  --open      open the workbench in the browser\n\nexample:\n  %s\n", examples[op])
		return
	}
	pos := positional[op]
	skip := map[string]bool{}
	var head strings.Builder
	for _, p := range pos {
		skip[p] = true
		head.WriteString(" [" + p + "]")
	}
	fmt.Fprintf(w, "usage: akmctl %s%s [flags]\n", op, head.String())
	params := api.Params(op)
	var names []string
	for k := range params {
		if !skip[k] {
			names = append(names, k)
		}
	}
	sort.Strings(names)
	fmt.Fprintf(w, "\nflags:\n")
	for _, k := range names {
		flag := "--" + strings.ReplaceAll(k, "_", "-")
		switch params[k] {
		case "bool":
			fmt.Fprintf(w, "  %s\n", flag)
		case "strings":
			fmt.Fprintf(w, "  %s <a,b>\n", flag)
		case "int":
			fmt.Fprintf(w, "  %s <n>\n", flag)
		default:
			fmt.Fprintf(w, "  %s <text>\n", flag)
		}
	}
	fmt.Fprintf(w, "  --repo <dir>  --format json|table|ids  (every command)\n")
	fmt.Fprintf(w, "\nexample:\n  %s\n", examples[op])
}

// wantsHelp reports whether a command's own arguments ask for its help.
func wantsHelp(args []string) bool {
	for _, a := range args {
		if a == "--" {
			return false
		}
		if a == "-h" || a == "--help" {
			return true
		}
	}
	return false
}
