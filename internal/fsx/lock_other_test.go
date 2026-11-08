//go:build windows

package fsx

import (
	"sync"
	"testing"
)

// TestLockRepoManyWriters takes and releases the lock from many goroutines at once, as
// concurrent writers do; the old fixed-pause loop could give up on a busy lock and
// a failed remove left it behind (work item 261029-q8z3).
func TestLockRepoManyWriters(t *testing.T) {
	root := t.TempDir()
	var wg sync.WaitGroup
	inside := 0
	var mu sync.Mutex
	for g := 0; g < 16; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				unlock, err := LockRepo(root, "docs")
				if err != nil {
					t.Error(err)
					return
				}
				mu.Lock()
				inside++
				if inside != 1 {
					t.Error("two holders of the lock")
				}
				inside--
				mu.Unlock()
				unlock()
			}
		}()
	}
	wg.Wait()
}
