// Package fsx holds the file operations every write goes through: the
// repository lock and the atomic write.
package fsx

import (
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Rename is os.Rename. On Windows a rename over a file that a reader has open
// fails with "access denied" until the reader closes it, so it is tried again
// for up to a second; a reader holds a file for the length of one read.
func Rename(from, to string) error {
	err := os.Rename(from, to)
	for i := 0; err != nil && runtime.GOOS == "windows" && i < 100; i++ {
		time.Sleep(10 * time.Millisecond)
		err = os.Rename(from, to)
	}
	return err
}

// WriteAtomic writes through a temporary file in the same directory and renames
// it over the target, so a crash leaves the old file or the new one.
func WriteAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".akm-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}
	if st, err := os.Stat(path); err == nil {
		os.Chmod(name, st.Mode().Perm())
	} else {
		os.Chmod(name, 0o644)
	}
	if err := Rename(name, path); err != nil {
		os.Remove(name)
		return err
	}
	return nil
}
