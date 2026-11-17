package api

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/alexsasharegan/akmctl/internal/fsx"
	"github.com/alexsasharegan/akmctl/internal/mdlinks"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// heldRepos are the repositories whose lock an apply holds for a whole
// sequence. The key is the *record.Repo that apply loaded for itself, which no
// other request has, so only the writes of that sequence skip the lock.
var heldRepos sync.Map

// lockRepo takes the repository lock for one write, or nothing when the write
// belongs to an apply that already holds it (the lock is not re-entrant).
func lockRepo(repo *record.Repo) (func(), error) {
	if _, ok := heldRepos.Load(repo); ok {
		return func() {}, nil
	}
	return fsx.LockRepo(repo.Root, repo.MetaDir())
}

// opApply makes an ordered list of writes as one step (ADR 261116-2qyf): it
// holds the repository lock from the first write to the last, and when a write
// fails it puts every file back as it was before the first one, from the bytes
// it read under the lock. The response is the response of each write, or the
// error of the one that failed.
func opApply(repo *record.Repo, info Info, q Request) (Obj, *Error) {
	writes, werr := writeList(q)
	if werr != nil {
		return nil, werr
	}
	for i, w := range writes {
		if _, ok := w["dry_run"]; ok {
			return nil, atWrite(i, invalid("writes", "writes[%d]: dry_run does not belong in apply; use preview to see a sequence before it is made", i))
		}
	}
	// A Repo of its own: the caller's may be shared, and the lock is skipped by
	// pointer.
	own, lerr := record.Load(repo.Root)
	if lerr != nil {
		return nil, unreadable(lerr)
	}
	unlock, lerr := fsx.LockRepo(own.Root, own.MetaDir())
	if lerr != nil {
		return nil, unreadable(lerr)
	}
	defer unlock()
	heldRepos.Store(own, struct{}{})
	defer heldRepos.Delete(own)
	defer repo.Reload()
	if lerr := own.Reload(); lerr != nil {
		return nil, unreadable(lerr)
	}
	before, lerr := snapshotFiles(own)
	if lerr != nil {
		return nil, unreadable(lerr)
	}
	fail := func(i int, e *Error) (Obj, *Error) {
		d := Obj{{"write", i}}
		if failed := before.restore(own); len(failed) > 0 {
			d = append(d, KV{"rolled_back", false}, KV{"rollback_errors", failed})
		} else {
			d = append(d, KV{"rolled_back", true})
		}
		return nil, &Error{Code: e.Code, Message: e.Message, Details: append(d, e.Details...)}
	}

	etags := map[string]string{} // the etag of each record after its last write
	newID := ""
	out := make([]Obj, 0, len(writes))
	for i, w := range writes {
		op, req := nextWrite(w, q, newID, etags)
		resp, err := Run(own, info, op, req)
		if err != nil {
			return fail(i, err)
		}
		if le, ok := resp.Get("link_errors"); ok {
			// close moved the file but could not rewrite every link to it
			return fail(i, &Error{Code: "unreadable", Message: fmt.Sprintf("%s could not rewrite the links to the moved file in %d file(s)", op, len(le.([]Obj))), Details: O("link_errors", le)})
		}
		if rec, ok := resp.Get("record"); ok {
			id, _ := rec.(Obj).Get("id")
			if op == "new" {
				newID, _ = id.(string)
			}
			if et, ok := rec.(Obj).Get("etag"); ok {
				etags[id.(string)], _ = et.(string)
			}
		}
		out = append(out, resp)
	}
	return O("writes", out), nil
}

// fileSnapshot is the content of the Markdown files a write can change, by path
// relative to the repository root, as they were when it was taken.
type fileSnapshot map[string][]byte

// writablePaths lists the files the writes of a sequence can change or make:
// the records, and the Markdown files whose links a close rewrites.
func writablePaths(repo *record.Repo) []string {
	seen := map[string]bool{}
	var out []string
	add := func(p string) {
		if !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	for _, p := range mdlinks.Files(repo) {
		add(p)
	}
	for _, r := range repo.Records {
		add(r.Path)
	}
	sort.Strings(out)
	return out
}

func snapshotFiles(repo *record.Repo) (fileSnapshot, error) {
	s := fileSnapshot{}
	for _, p := range writablePaths(repo) {
		b, err := os.ReadFile(filepath.Join(repo.Root, filepath.FromSlash(p)))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		s[p] = b
	}
	return s, nil
}

// restore puts the files back as the snapshot has them: a file that changed or
// went away is written again, and a file made since is removed. It returns the
// files it could not put back.
func (s fileSnapshot) restore(repo *record.Repo) []Obj {
	var failed []Obj
	repo.Reload() // to see the files made since; on an error the old list still names the folders
	for _, p := range writablePaths(repo) {
		if _, ok := s[p]; ok {
			continue
		}
		if err := os.Remove(filepath.Join(repo.Root, filepath.FromSlash(p))); err != nil && !os.IsNotExist(err) {
			failed = append(failed, O("path", p, "error", err.Error()))
		}
	}
	paths := make([]string, 0, len(s))
	for p := range s {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		full := filepath.Join(repo.Root, filepath.FromSlash(p))
		if cur, err := os.ReadFile(full); err == nil && bytes.Equal(cur, s[p]) {
			continue
		}
		if err := fsx.WriteAtomic(full, s[p]); err != nil {
			failed = append(failed, O("path", p, "error", err.Error()))
		}
	}
	repo.Reload()
	return failed
}
