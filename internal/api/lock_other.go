//go:build windows

package api

import (
	"os"
	"path/filepath"
	"time"
)

// lockRepo takes the lock by creating docs/.akm.lock exclusively and retrying;
// flock is not available here.
func lockRepo(root string) (func(), error) {
	dir := filepath.Join(root, "docs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, ".akm.lock")
	for i := 0; ; i++ {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			f.Close()
			return func() { os.Remove(path) }, nil
		}
		if i > 200 {
			return nil, err
		}
		time.Sleep(25 * time.Millisecond)
	}
}
