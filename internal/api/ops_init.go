package api

import (
	"os"
	"path/filepath"

	"github.com/alexsasharegan/akmctl/internal/check"
	"github.com/alexsasharegan/akmctl/internal/initrepo"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// opInit creates the layout of a profile in the directory dir (default: the
// base directory itself), relative to repo.Root. It runs before any repository
// exists, so repo only supplies the base directory. It never overwrites a file.
// With git_init it also runs `git init` when the directory is not inside a Git
// working tree; a failure of that is reported, not fatal.
func opInit(repo *record.Repo, q Request) (Obj, *Error) {
	p, ok := initrepo.Get(q.str("profile"))
	if !ok {
		return nil, invalid("profile", "profile: unknown value %q", q.str("profile"))
	}
	dir := "."
	if d, given := q["dir"].(string); given {
		dir = d
		if d != "." {
			if err := record.CheckRelPath(d); err != nil {
				return nil, invalid("dir", "dir %q %v", d, err)
			}
		}
	}
	if repo == nil || repo.Root == "" {
		return nil, &Error{Code: "unreadable", Message: "init needs a base directory"}
	}
	target := filepath.Join(repo.Root, filepath.FromSlash(dir))
	if st, err := os.Stat(target); err == nil && !st.IsDir() {
		return nil, refused("not_a_directory", "%s exists and is not a directory", dir)
	}
	// git first: git_commits only works inside a work tree
	if err := os.MkdirAll(target, 0o755); err != nil {
		return nil, unreadable(err)
	}
	gitState, gitErr := "none", ""
	if q.boolean("git_init") {
		if record.InGitRepo(target) {
			gitState = "existing"
		} else if gerr := initrepo.GitInit(target); gerr != nil {
			gitState, gitErr = "failed", gerr.Error()
		} else {
			gitState = "initialised"
		}
	}
	commits := q.boolean("git_commits")
	forcedOff := commits && !record.InGitRepo(target)
	if forcedOff {
		commits = false
	}
	res, err := initrepo.InitWith(target, p, initrepo.Options{GitCommits: commits})
	if err != nil {
		return nil, unreadable(err)
	}
	out := O("profile", p.Name, "dir", dir, "layout", check.LayoutVersion, "written", nonNil(res.Written), "skipped", nonNil(res.Skipped),
		"git_commits", commits, "git_commits_forced_off", forcedOff, "git", gitState)
	if gitErr != "" {
		out = append(out, KV{"git_error", gitErr})
	}
	return out, nil
}

func nonNil(l []string) []string {
	if l == nil {
		return []string{}
	}
	return l
}
