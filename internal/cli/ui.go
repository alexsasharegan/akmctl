package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/api"
	"github.com/alexsasharegan/akmctl/internal/record"
	"github.com/alexsasharegan/akmctl/internal/server"
)

// runUI starts the HTTP binding and the workbench page on the loopback
// interface and blocks until interrupted.
func runUI(args []string, stderr io.Writer, cwd string) int {
	port, open, repoDir := 0, false, cwd
	for i := 0; i < len(args); i++ {
		name, val, hasVal := strings.Cut(strings.TrimPrefix(args[i], "--"), "=")
		next := func() string {
			if hasVal {
				return val
			}
			if i+1 < len(args) {
				i++
				return args[i]
			}
			return ""
		}
		switch {
		case !strings.HasPrefix(args[i], "--"):
			fmt.Fprintf(stderr, "akmctl ui: unexpected argument %q\n", args[i])
			return 2
		case name == "port":
			n, err := strconv.Atoi(next())
			if err != nil || n < 0 || n > 65535 {
				fmt.Fprintln(stderr, "akmctl ui: --port needs a number from 0 to 65535")
				return 2
			}
			port = n
		case name == "open":
			open = true
		case name == "repo":
			repoDir = next()
		default:
			fmt.Fprintf(stderr, "akmctl ui: unknown flag --%s\n", name)
			return 2
		}
	}
	root := repoDir
	if !flagGiven(args, "--repo") {
		var err error
		if root, err = record.FindRoot(repoDir); err != nil {
			fmt.Fprintln(stderr, "akmctl ui:", err)
			return 2
		}
	}
	if _, err := record.Load(root); err != nil {
		fmt.Fprintln(stderr, "akmctl ui:", err)
		return 2
	}
	token, err := server.NewToken()
	if err != nil {
		fmt.Fprintln(stderr, "akmctl ui:", err)
		return 2
	}
	l, err := server.Listen(port)
	if err != nil {
		fmt.Fprintln(stderr, "akmctl ui:", err)
		return 2
	}
	url := "http://" + l.Addr().String() + "/#token=" + token
	fmt.Fprintf(stderr, "akmctl ui: serving %s on the loopback interface only\nopen: %s\nthe token is shown once and kept in memory; Ctrl-C stops the server\n", root, url)
	if open {
		openBrowser(url)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	srv := &server.Server{Root: root, Token: token, Info: api.Info{Name: "akmctl", Version: Version}}
	if err := srv.Serve(ctx, l); err != nil {
		fmt.Fprintln(stderr, "akmctl ui:", err)
		return 2
	}
	return 0
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Start()
}
