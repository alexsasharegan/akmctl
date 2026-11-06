package api

import (
	"strconv"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/record"
)

// history lists the AKM step commits of a range as moves, oldest first (ADR
// 261106-60xv). It shares the git plumbing of commits and judges nothing: a
// malformed step commit is still a move, and commits is where it is reported.

const (
	historyDefaultLimit = 100
	historyMaxLimit     = 500
)

func opHistory(repo *record.Repo, q Request) (Obj, *Error) {
	rng := q.str("range")
	if rng != "" && (strings.HasPrefix(rng, "-") || !reRevRange.MatchString(rng)) {
		return nil, invalid("range", "range must be a git revision range such as origin/main..HEAD")
	}
	limit := q.integer("limit", historyDefaultLimit)
	if limit > historyMaxLimit {
		return nil, invalid("limit", "limit is at most %d", historyMaxLimit)
	}
	item := q.str("item")
	unavailable := func(reason string) (Obj, *Error) {
		return O("available", false, "reason", reason, "head", "", "moves", []Obj{}, "truncated", false), nil
	}
	if _, err := runGit(repo.Root, "rev-parse", "--is-inside-work-tree"); err != nil {
		return unavailable("not_git")
	}
	if !repo.GitCommits() {
		return unavailable("git_commits_off")
	}
	head, err := runGit(repo.Root, "rev-parse", "--verify", "-q", "HEAD")
	head = strings.TrimSpace(head)
	if err != nil {
		// no commits yet
		return O("available", true, "head", "", "moves", []Obj{}, "truncated", false), nil
	}
	// the subject of a step commit starts with <id>(; the filter is a first cut and the
	// subject is checked again below, so a body line that happens to match is dropped
	id := `[0-9]{6}-[0-9a-hjkmnp-tv-z]{4}`
	if item != "" {
		id = item
	}
	args := []string{"log", "--no-merges", "--extended-regexp", "--grep=^" + id + "\\(", "-n", strconv.Itoa(limit + 1), "--format=%H%x00%aI%x00%an%x00%B%x1e"}
	if rng != "" {
		args = append(args, rng, "--")
	} else {
		args = append(args, "HEAD", "--")
	}
	raw, err := runGit(repo.Root, args...)
	if err != nil {
		return nil, invalid("range", "git cannot read the range %q", rng)
	}
	var found []Obj // newest first, as git prints them
	for _, chunk := range strings.Split(raw, "\x1e") {
		chunk = strings.TrimLeft(chunk, "\n")
		if strings.TrimSpace(chunk) == "" {
			continue
		}
		parts := strings.SplitN(chunk, "\x00", 4)
		if len(parts) != 4 {
			continue
		}
		sha, when, author, msg := parts[0], parts[1], parts[2], strings.TrimRight(parts[3], "\n")
		lines := strings.Split(msg, "\n")
		m := reCommitSubject.FindStringSubmatch(lines[0])
		if m == nil || !commitSteps[m[2]] || (item != "" && m[1] != item) {
			continue
		}
		from, to := moveStatus(repo, sha, m[1])
		found = append(found, O("sha", sha, "time", when, "author", author, "item", m[1], "step", m[2], "from", from, "to", to, "trailers", messageTrailers(lines)))
	}
	truncated := len(found) > limit
	if truncated {
		found = found[:limit]
	}
	moves := make([]Obj, 0, len(found))
	for i := len(found) - 1; i >= 0; i-- { // oldest first
		moves = append(moves, found[i])
	}
	return O("available", true, "head", head, "moves", moves, "truncated", truncated), nil
}

// messageTrailers returns the trailers of the last paragraph of a message as
// "Key: value" strings, in order, without the AKM- ones.
func messageTrailers(lines []string) []string {
	out := []string{}
	start := len(lines)
	for i := len(lines) - 1; i >= 1; i-- {
		if strings.TrimSpace(lines[i]) == "" {
			start = i + 1
			break
		}
	}
	for _, l := range lines[start:] {
		k, _, ok := strings.Cut(l, ": ")
		if !ok || k == "" || strings.ContainsAny(k, " \t") || strings.HasPrefix(l, "AKM-Step:") || strings.HasPrefix(l, "AKM-Item:") {
			continue
		}
		out = append(out, strings.TrimSpace(l))
	}
	return out
}

// moveStatus is the status of an item's file before and after a commit: the
// file at the parent (or the old name of a moved file) and at the commit. A
// file the commit adds has no "from"; one that cannot be read has no "to".
func moveStatus(repo *record.Repo, sha, item string) (from, to string) {
	raw, err := runGit(repo.Root, "diff-tree", "-r", "-M", "--root", "--no-commit-id", "--name-status", "--relative", "-z", sha)
	if err != nil {
		return "", ""
	}
	isRecord := func(p string) bool {
		base := p[strings.LastIndex(p, "/")+1:]
		return strings.HasPrefix(base, item+"-") && strings.HasSuffix(base, ".md") && (repo.InWork(p) || repo.InADR(p))
	}
	for _, c := range parseNameStatus(raw) {
		if !isRecord(c.path) && !(c.from != "" && isRecord(c.from)) {
			continue
		}
		if c.status[0] != 'D' {
			to, _, _ = recordState(repo, sha, c.path)
		}
		if c.status[0] != 'A' {
			old := c.path
			if c.from != "" {
				old = c.from
			}
			from, _, _ = recordState(repo, sha+"^", old)
		}
		return from, to
	}
	return "", ""
}
