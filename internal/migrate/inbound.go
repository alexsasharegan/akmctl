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
	return inboundLinksOfMoves(repo, map[string]string{from: to})
}

// maxInboundNotes bounds the inbound-link lines of one move_closed step: a
// vault can hold thousands of links to the records it moves. The notes beyond
// it are counted in a last line (work item 261108-3t7t).
const maxInboundNotes = 50

// inboundLinksOfMoves is inboundLinks for several moves (from path to new path)
// read in one pass over the files. A file that moves too is read at its new
// place: its relative link to a file that moves along with it still resolves,
// so it is not reported.
func inboundLinksOfMoves(repo *record.Repo, moves map[string]string) []string {
	var notes []string
	for _, f := range mdlinks.Files(repo) {
		b, err := os.ReadFile(filepath.Join(repo.Root, filepath.FromSlash(f)))
		if err != nil {
			continue
		}
		at := f // where the file will be after the moves
		if to, ok := moves[f]; ok {
			at = to
		}
		dir, newDir := path.Dir(f), path.Dir(at)
		for _, l := range mdlinks.Scan(string(b)) {
			from := path.Join(dir, l.Target)
			if to, ok := moves[from]; ok && path.Join(newDir, l.Target) != to {
				notes = append(notes, fmt.Sprintf("inbound link left broken: %s:%d points at %s, moved to %s", f, l.Line, from, to))
			}
		}
		for _, w := range mdlinks.Wikilinks(string(b)) {
			for from, to := range moves {
				if pathQualified(w.Target, from) {
					notes = append(notes, fmt.Sprintf("inbound wikilink left broken: %s:%d points at %s, moved to %s", f, w.Line, from, to))
				}
			}
		}
	}
	return notes
}

// boundNotes keeps the first max notes and sums up the rest in one line.
func boundNotes(notes []string, max int) []string {
	if len(notes) <= max {
		return notes
	}
	out := append([]string(nil), notes[:max]...)
	return append(out, fmt.Sprintf("... and %d more inbound links left broken", len(notes)-max))
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
