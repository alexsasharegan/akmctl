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
	res, err := initrepo.Init(target, p)
	if err != nil {
		return nil, unreadable(err)
	}
	return O("profile", p.Name, "dir", dir, "layout", check.LayoutVersion, "written", nonNil(res.Written), "skipped", nonNil(res.Skipped)), nil
}

func nonNil(l []string) []string {
	if l == nil {
		return []string{}
	}
	return l
}
