package migrate

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/mdlinks"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// inboundLinks reports the relative Markdown links, in the repository's other
// files, that point at from, which a move_file step is about to move to to. The
// step does not repair them (close does; work item 261029-m0vm), so each is a
// note: file, line, and where the file goes. They are read from the files as they
// are before the move, in path order.
func inboundLinks(repo *record.Repo, from, to string) []string {
	var notes []string
	for _, f := range mdlinks.Files(repo) {
		if f == from || f == to {
			continue
		}
		b, err := os.ReadFile(filepath.Join(repo.Root, filepath.FromSlash(f)))
		if err != nil {
			continue
		}
		dir := path.Dir(f)
		for _, l := range mdlinks.Scan(string(b)) {
			if path.Join(dir, l.Target) == from {
				notes = append(notes, fmt.Sprintf("inbound link left broken: %s:%d points at %s, moved to %s", f, l.Line, from, to))
			}
		}
		for _, w := range mdlinks.Wikilinks(string(b)) {
			if pathQualified(w.Target, from) {
				notes = append(notes, fmt.Sprintf("inbound wikilink left broken: %s:%d points at %s, moved to %s", f, w.Line, from, to))
			}
		}
	}
	return notes
}

// pathQualified reports whether a wikilink target names the file from by a path
// ("old/note" or "docs/old/note"). A bare name ("note") is resolved by name, so
// a move does not break it and it is not reported.
func pathQualified(target, from string) bool {
	if !strings.Contains(target, "/") {
		return false
	}
	name := strings.TrimSuffix(from, ".md")
	target = strings.TrimPrefix(target, "/")
	return name == target || strings.HasSuffix(name, "/"+target)
}
