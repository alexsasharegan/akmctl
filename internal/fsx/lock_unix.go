//go:build !windows

package fsx

import (
	"os"
	"path/filepath"
	"syscall"
)

// LockRepo takes the exclusive advisory lock <metaDir>/.akm.lock for one write and
// returns the function that releases it.
func LockRepo(root, metaDir string) (func(), error) {
	dir := filepath.Join(root, filepath.FromSlash(metaDir))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(filepath.Join(dir, ".akm.lock"), os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		return nil, err
	}
	return func() {
		syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		f.Close()
	}, nil
}
