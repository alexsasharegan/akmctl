package api

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/record"
)

// opFind turns an id into a record's path, or a note name (the target of a
// wikilink) into the path of the one Markdown file that has it.
func opFind(repo *record.Repo, q Request) (Obj, *Error) {
	id, name := q.str("id"), q.str("name")
	switch {
	case id != "" && name != "":
		return nil, invalid("name", "give id or name, not both")
	case id == "" && name == "":
		return nil, invalid("id", "missing required parameter \"id\" (or give name)")
	case id != "":
		r, err := lookup(repo, id)
		if err != nil {
			return nil, err
		}
		return O("id", r.ID, "type", r.Type, "path", r.Path), nil
	}
	p, err := findNote(repo, name)
	if err != nil {
		return nil, err
	}
	out := O("name", name, "path", p)
	for _, r := range repo.Records {
		if r.Path == p {
			out = append(out, KV{"id", r.ID}, KV{"type", r.Type})
			break
		}
	}
	return out, nil
}

// skipDir is a folder a note name is never looked up in: hidden folders (.git,
// .obsidian, .trash, .agents) and dependencies.
func skipDir(name string) bool { return strings.HasPrefix(name, ".") || name == "node_modules" }

// findNote finds the Markdown file for a note name the way a wikilink names
// it: the file name without .md, in any folder, with the folders it gives as a
// trailing path ("Projects/Vault Views"). The match ignores case. Links are not
// followed. No match is not_found; several are ambiguous.
func findNote(repo *record.Repo, name string) (string, *Error) {
	want := strings.ToLower(strings.TrimSuffix(strings.Trim(strings.TrimSpace(name), "/"), ".md"))
	if want == "" || strings.Contains(want, "\\") || strings.Contains(want, "..") {
		return "", invalid("name", "name is not a note name")
	}
	var hits []string
	werr := filepath.WalkDir(repo.Root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if d != nil && d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			if p != repo.Root && skipDir(d.Name()) {
				return fs.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() || !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}
		rel, rerr := filepath.Rel(repo.Root, p)
		if rerr != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		stem := strings.ToLower(strings.TrimSuffix(rel, filepath.Ext(rel)))
		if stem == want || strings.HasSuffix(stem, "/"+want) {
			hits = append(hits, rel)
		}
		return nil
	})
	if werr != nil {
		return "", unreadable(werr)
	}
	sort.Strings(hits)
	switch len(hits) {
	case 0:
		return "", &Error{Code: "not_found", Message: "no note named " + name}
	case 1:
		return hits[0], nil
	}
	return "", &Error{"ambiguous", "name " + name + " names more than one note", O("paths", hits)}
}
