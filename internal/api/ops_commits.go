package api

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/frontmatter"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// commits checks the history against the commit object (ADR 261030-tbkh):
// every commit that claims to be an AKM step has the subject and trailers the
// tool would have returned, and its diff does what the step says. It is
// read-only and the only operation besides migrate that runs git.

const (
	commitsDefaultLimit = 50
	commitsMaxLimit     = 500
)

// reRevRange is a conservative revision range: no leading dash, no spaces and
// no shell or option characters, so a request cannot pass an option to git.
var reRevRange = regexp.MustCompile(`^[A-Za-z0-9_./~^@{}-]+(\.\.\.?[A-Za-z0-9_./~^@{}-]+)?$`)

var reCommitSubject = regexp.MustCompile(`^(\d{6}-[0-9a-hjkmnp-tv-z]{4})\(([a-z]+)\): (.+)$`)

type gitChange struct {
	status string // A, M, D, R, ...
	path   string // the new path
	from   string // the old path of a rename
}

func opCommits(repo *record.Repo, q Request) (Obj, *Error) {
	rng := q.str("range")
	if rng != "" && (strings.HasPrefix(rng, "-") || !reRevRange.MatchString(rng)) {
		return nil, invalid("range", "range must be a git revision range such as origin/main..HEAD")
	}
	limit := q.integer("limit", commitsDefaultLimit)
	if limit > commitsMaxLimit {
		return nil, invalid("limit", "limit is at most %d", commitsMaxLimit)
	}
	if _, err := runGit(repo.Root, "rev-parse", "--is-inside-work-tree"); err != nil {
		return nil, refused("not_git", "the repository is not a git work tree")
	}
	args := []string{"log", "--no-merges", "-n", strconv.Itoa(limit), "--format=%H%x00%B%x1e"}
	if rng != "" {
		args = append(args, rng, "--")
	} else {
		args = append(args, "HEAD", "--")
	}
	raw, err := runGit(repo.Root, args...)
	if err != nil {
		if rng == "" {
			// no commits yet
			return O("checked", 0, "skipped", 0, "clean", true, "commits", []Obj{}), nil
		}
		return nil, invalid("range", "git cannot read the range %q", rng)
	}
	checked, skipped := 0, 0
	clean := true
	out := []Obj{}
	for _, chunk := range strings.Split(raw, "\x1e") {
		chunk = strings.TrimLeft(chunk, "\n")
		if strings.TrimSpace(chunk) == "" {
			continue
		}
		parts := strings.SplitN(chunk, "\x00", 2)
		if len(parts) != 2 {
			continue
		}
		sha, msg := parts[0], strings.TrimRight(parts[1], "\n")
		c, ok := verifyCommit(repo, sha, msg)
		if !ok {
			skipped++
			continue
		}
		checked++
		if probs, _ := c.Get("problems"); len(probs.([]Obj)) > 0 {
			clean = false
		}
		out = append(out, c)
	}
	return O("checked", checked, "skipped", skipped, "clean", clean, "commits", out), nil
}

func runGit(root string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

// verifyCommit returns the report of one commit, or false when the commit makes
// no claim to be an AKM step (no AKM- trailer and no step subject).
func verifyCommit(repo *record.Repo, sha, msg string) (Obj, bool) {
	lines := strings.Split(msg, "\n")
	subject := lines[0]
	var stepTr, itemTr []string
	for _, l := range lines[1:] {
		if v, ok := strings.CutPrefix(l, "AKM-Step: "); ok {
			stepTr = append(stepTr, strings.TrimSpace(v))
		} else if v, ok := strings.CutPrefix(l, "AKM-Item: "); ok {
			itemTr = append(itemTr, strings.TrimSpace(v))
		} else if strings.HasPrefix(l, "AKM-") {
			stepTr = append(stepTr, l) // an unknown AKM- trailer still counts as a claim
		}
	}
	m := reCommitSubject.FindStringSubmatch(subject)
	if m == nil && len(stepTr) == 0 && len(itemTr) == 0 {
		return nil, false
	}
	var problems []Obj
	add := func(rule, format string, a ...any) {
		problems = append(problems, O("rule", rule, "message", fmt.Sprintf(format, a...)))
	}
	step, item := "", ""
	switch {
	case m == nil:
		add("CM1", "the subject %q is not <id>(<step>): <detail>", subject)
	case !commitSteps[m[2]]:
		add("CM1", "step %q is not one of the AKM steps", m[2])
	default:
		item, step = m[1], m[2]
	}
	switch {
	case len(stepTr) != 1 || len(itemTr) != 1:
		add("CM2", "a step commit has exactly one AKM-Step and one AKM-Item trailer (found %d and %d)", len(stepTr), len(itemTr))
	default:
		if step != "" && stepTr[0] != step {
			add("CM2", "AKM-Step is %q but the subject says %q", stepTr[0], step)
		}
		if item != "" && itemTr[0] != item {
			add("CM2", "AKM-Item is %q but the subject says %q", itemTr[0], item)
		}
	}
	if step != "" && item != "" {
		for _, p := range diffProblems(repo, sha, step, item) {
			add("CM3", "%s", p)
		}
	}
	if problems == nil {
		problems = []Obj{}
	}
	return O("sha", sha, "subject", subject, "step", step, "item", item, "ok", len(problems) == 0, "problems", problems), true
}

// diffProblems checks that what the commit changed is what the step says.
func diffProblems(repo *record.Repo, sha, step, item string) []string {
	raw, err := runGit(repo.Root, "diff-tree", "-r", "-M", "--root", "--no-commit-id", "--name-status", "--relative", "-z", sha)
	if err != nil {
		return []string{"git cannot read the diff of " + sha}
	}
	changes := parseNameStatus(raw)
	isRecord := func(p string) bool {
		base := p[strings.LastIndex(p, "/")+1:]
		return strings.HasPrefix(base, item+"-") && strings.HasSuffix(base, ".md") &&
			(repo.InWork(p) || repo.InADR(p))
	}
	var touched []gitChange
	for _, c := range changes {
		if isRecord(c.path) || (c.from != "" && isRecord(c.from)) {
			touched = append(touched, c)
		}
	}
	if len(touched) == 0 {
		return []string{"no file of " + item + " changed"}
	}
	// the item's file after the commit: prefer the one that exists there
	var cur *gitChange
	for i := range touched {
		if touched[i].status[0] != 'D' {
			cur = &touched[i]
		}
	}
	state := func(rev, path string) (status string, claimed bool, ok bool) {
		out, err := runGit(repo.Root, "show", rev+":./"+path)
		if err != nil {
			return "", false, false
		}
		doc, perr := frontmatter.Parse([]byte(out))
		if perr != nil {
			return "", false, false
		}
		v, _ := doc.Get("status")
		cb, has := doc.Get("claimed_by")
		return v.Str, has && cb.Kind == frontmatter.String && cb.Str != "", true
	}
	switch step {
	case "file":
		if cur == nil || cur.status[0] != 'A' {
			return []string{"file: the commit does not add a file for " + item}
		}
		return nil
	case "close", "drop":
		want := map[string]string{"close": "done", "drop": "dropped"}[step]
		if cur == nil || !repo.IsDone(cur.path) || cur.status[0] == 'M' {
			return []string{step + ": the commit does not move the file of " + item + " into " + repo.DoneDir() + "/"}
		}
		st, _, ok := state(sha, cur.path)
		if !ok || st != want {
			return []string{step + ": the file's status is " + quote(st) + ", want " + want}
		}
		return nil
	}
	// the other steps edit the item's file in place
	if cur == nil || cur.status[0] != 'M' {
		return []string{step + ": the commit does not modify the file of " + item}
	}
	after, aClaim, aok := state(sha, cur.path)
	before, bClaim, bok := state(sha+"^", cur.path)
	if !aok || !bok {
		return []string{step + ": cannot read the file of " + item + " before and after the commit"}
	}
	var bad string
	switch step {
	case "claim":
		if after != "active" || !aClaim {
			bad = "after the commit the status is " + quote(after) + " and claimed_by is " + present(aClaim) + "; a claim leaves it active and claimed"
		}
	case "release":
		if !bClaim || aClaim {
			bad = "a release removes an existing claimed_by (before: " + present(bClaim) + ", after: " + present(aClaim) + ")"
		}
	case "ready":
		if after != "ready" || before == "ready" {
			bad = "status went from " + quote(before) + " to " + quote(after) + ", want a change to ready"
		}
	case "start":
		if after != "active" || before == "active" {
			bad = "status went from " + quote(before) + " to " + quote(after) + ", want a change to active"
		}
	case "block":
		if after != "blocked" || before == "blocked" {
			bad = "status went from " + quote(before) + " to " + quote(after) + ", want a change to blocked"
		}
	case "unblock":
		if before != "blocked" || after == "blocked" {
			bad = "status went from " + quote(before) + " to " + quote(after) + ", want a change from blocked"
		}
	}
	if bad != "" {
		return []string{step + ": " + bad}
	}
	return nil
}

func quote(s string) string { return "\"" + s + "\"" }

func present(b bool) string {
	if b {
		return "present"
	}
	return "absent"
}

// parseNameStatus reads `git diff-tree --name-status -z`: a status, then one
// path (two for a rename or copy), each ended by NUL.
func parseNameStatus(raw string) []gitChange {
	f := strings.Split(strings.TrimRight(raw, "\x00"), "\x00")
	var out []gitChange
	for i := 0; i < len(f); {
		st := f[i]
		if st == "" {
			i++
			continue
		}
		if (st[0] == 'R' || st[0] == 'C') && i+2 < len(f) {
			out = append(out, gitChange{status: st, from: f[i+1], path: f[i+2]})
			i += 3
			continue
		}
		if i+1 >= len(f) {
			break
		}
		out = append(out, gitChange{status: st, path: f[i+1]})
		i += 2
	}
	return out
}
