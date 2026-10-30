package api

import (
	"os"
	"path/filepath"

	"github.com/alexsasharegan/akmctl/internal/check"
	"github.com/alexsasharegan/akmctl/internal/initrepo"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// opInit creates the layout in the directory dir (default: the base directory
// itself), relative to repo.Root. It runs before any repository exists, so
// repo only supplies the base directory. It never overwrites a file. With
// git_init it also runs the version-control init when the directory is not
// inside a work tree; a failure of that is reported, not fatal.
func opInit(repo *record.Repo, q Request) (Obj, *Error) {
	docsDir := initrepo.DefaultDocsDir
	if d, given := q["docs_dir"].(string); given {
		if err := initrepo.ValidateDocsDir(d); err != nil {
			return nil, invalid("docs_dir", "docs_dir %q %v", d, err)
		}
		docsDir = d
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
	// a second manifest would make discovery ambiguous
	if have, found, err := record.Locate(target); err != nil {
		return nil, &Error{Code: "unreadable", Message: err.Error()}
	} else if found && have != docsDir {
		return nil, refused("manifest_exists", "%s/%s already exists; init would make a second manifest: use docs_dir %q or remove it", have, record.ManifestName, have)
	}
	// the work tree first: git_commits and .gitignore need one
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
	inTree := record.InGitRepo(target)
	commits := q.boolean("git_commits")
	forcedOff := commits && !inTree
	if forcedOff {
		commits = false
	}
	wiki := q.boolean("wikilinks")
	obsidian := false
	if st, err := os.Stat(filepath.Join(target, ".obsidian")); err == nil && st.IsDir() {
		obsidian = true
	}
	res, err := initrepo.InitWith(target, initrepo.Options{DocsDir: docsDir, WikiLinks: wiki, GitCommits: commits, Gitignore: inTree, Obsidian: obsidian})
	if err != nil {
		return nil, unreadable(err)
	}
	out := O("docs_dir", docsDir, "wikilinks", wiki, "dir", dir, "layout", check.LayoutVersion, "written", nonNil(res.Written), "skipped", nonNil(res.Skipped),
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
