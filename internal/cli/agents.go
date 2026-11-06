package cli

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/alexsasharegan/akmctl/internal/initrepo"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// runAgents prints the AGENTS.md that init would write for this repository: the
// shipped template filled with the records folder, the wikilink style and the
// project name the repository already has. It reads the manifest and writes
// nothing, so a repository whose AGENTS.md predates a template change can
// compare the two and keep the lines it judges.
func runAgents(args []string, stdout, stderr io.Writer, cwd string) int {
	repoDir, given := cwd, false
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--repo" && i+1 < len(args):
			i++
			repoDir, given = args[i], true
		case len(args[i]) > 7 && args[i][:7] == "--repo=":
			repoDir, given = args[i][7:], true
		default:
			fmt.Fprintf(stderr, "akmctl agents: unexpected argument %q (usage: akmctl agents [--repo <dir>])\n", args[i])
			return 2
		}
	}
	root, err := record.FindRoot(repoDir)
	if given {
		root, err = repoDir, nil
		if e := insideRepository("agents", repoDir); e != nil {
			err = e
		}
	}
	if err == nil {
		if _, found, lerr := record.Locate(root); lerr != nil {
			err = lerr
		} else if !found {
			err = fmt.Errorf("%s has no manifest (run akmctl init first)", root)
		}
	}
	var repo *record.Repo
	if err == nil {
		repo, err = record.Load(root)
	}
	if err != nil {
		fmt.Fprintln(stderr, "akmctl agents:", err)
		return 2
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		abs = root
	}
	opt := initrepo.Options{DocsDir: repo.MetaDir(), Tool: Version}
	if v, ok := repo.Layout["reference_style"]; ok && v.Str == "wikilink" {
		opt.WikiLinks = true
	}
	text, err := initrepo.AgentsFile(filepath.Base(abs), opt)
	if err != nil {
		fmt.Fprintln(stderr, "akmctl agents:", err)
		return 2
	}
	fmt.Fprint(stdout, text)
	return 0
}
