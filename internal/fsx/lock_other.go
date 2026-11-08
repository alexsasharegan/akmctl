//go:build windows

package fsx

import (
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// lockWait is how long LockRepo waits for another writer. A writer holds the
// lock for the length of one file write, so a wait this long means the lock is
// stuck, not busy.
const lockWait = 30 * time.Second

// LockRepo takes the lock by creating <metaDir>/.akm.lock exclusively and retrying;
// flock is not available here. It waits up to lockWait, trying again after a short,
// random pause so that writers that started together do not keep colliding. The
// release function removes the file and tries again while Windows refuses (a file
// that something else has open cannot be deleted yet): a lock left behind would
// block every later writer (work item 261029-q8z3).
func LockRepo(root, metaDir string) (func(), error) {
	dir := filepath.Join(root, filepath.FromSlash(metaDir))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, ".akm.lock")
	deadline := time.Now().Add(lockWait)
	for {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			f.Close()
			return func() { removeLock(path) }, nil
		}
		if time.Now().After(deadline) {
			return nil, err
		}
		time.Sleep(time.Duration(2+rand.Intn(14)) * time.Millisecond)
	}
}

// removeLock deletes the lock file, trying for up to ten seconds.
func removeLock(path string) {
	for i := 0; i < 1000; i++ {
		err := os.Remove(path)
		if err == nil || os.IsNotExist(err) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}
