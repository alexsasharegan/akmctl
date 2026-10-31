package fsx

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteAtomicSyncsBeforeAndAfterTheRename(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "rec.md")
	var events []string
	oldF, oldD := syncFile, syncDir
	defer func() { syncFile, syncDir = oldF, oldD }()
	syncFile = func(f *os.File) error {
		// the temp file still has its temp name and the target is not there yet
		if _, err := os.Stat(target); err == nil {
			t.Error("target exists before the file sync")
		}
		events = append(events, "file")
		return nil
	}
	syncDir = func(d string) error {
		if b, err := os.ReadFile(target); err != nil || string(b) != "new" {
			t.Errorf("folder sync before the rename: %v %q", err, b)
		}
		events = append(events, "dir:"+filepath.Base(d))
		return nil
	}
	if err := WriteAtomic(target, []byte("new")); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(events, ","); got != "file,dir:"+filepath.Base(dir) {
		t.Errorf("order = %s", got)
	}
}

func TestWriteAtomicFileSyncFailureLeavesTheOldFile(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "rec.md")
	os.WriteFile(target, []byte("old"), 0o644)
	oldF := syncFile
	defer func() { syncFile = oldF }()
	syncFile = func(*os.File) error { return errors.New("disk full") }
	if err := WriteAtomic(target, []byte("new")); err == nil {
		t.Fatal("want an error")
	}
	if b, _ := os.ReadFile(target); string(b) != "old" {
		t.Errorf("target = %q", b)
	}
	left, _ := filepath.Glob(filepath.Join(dir, ".akm-*.tmp"))
	if len(left) != 0 {
		t.Errorf("temp files left: %v", left)
	}
}

func TestWriteAtomicFolderSyncFailureIsNotAnError(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "rec.md")
	oldD := syncDir
	defer func() { syncDir = oldD }()
	syncDir = func(string) error { return errors.New("not supported") }
	if err := WriteAtomic(target, []byte("new")); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(target); string(b) != "new" {
		t.Errorf("target = %q", b)
	}
}
