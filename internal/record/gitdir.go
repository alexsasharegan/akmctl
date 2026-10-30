package record

import (
	"os"
	"path/filepath"
)

// InGitRepo reports whether dir, or the nearest of its ancestors that exists,
// is inside a Git working tree: a .git directory or file on the way up.
func InGitRepo(dir string) bool {
	d, err := filepath.Abs(dir)
	if err != nil {
		return false
	}
	for {
		if _, err := os.Lstat(filepath.Join(d, ".git")); err == nil {
			return true
		}
		parent := filepath.Dir(d)
		if parent == d {
			return false
		}
		d = parent
	}
}
