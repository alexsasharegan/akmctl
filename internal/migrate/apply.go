package migrate

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alexsasharegan/akmctl/internal/frontmatter"
	"github.com/alexsasharegan/akmctl/internal/fsx"
)

// Safeguard errors, with the reason code the API reports.
type Refusal struct {
	Reason  string
	Message string
}

func (r *Refusal) Error() string { return r.Message }

func git(root string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// Preflight checks the safeguards before anything is written: a git work tree
// with no uncommitted changes. It returns the tag and branch to create.
func Preflight(root, from, to string) (tag, branch string, err error) {
	if _, e := git(root, "rev-parse", "--is-inside-work-tree"); e != nil {
		return "", "", &Refusal{"not_git", "the repository is not a git work tree, so no backup tag can be made"}
	}
	if out, e := git(root, "status", "--porcelain"); e != nil || out != "" {
		return "", "", &Refusal{"dirty_tree", "the git tree has uncommitted changes; commit or stash them first"}
	}
	tag = "pre-migrate-" + from + "-to-" + to
	branch = "migrate-" + from + "-to-" + to
	if out, _ := git(root, "tag", "--list", tag); out != "" {
		return "", "", &Refusal{"tag_exists", "the backup tag " + tag + " already exists"}
	}
	if out, _ := git(root, "branch", "--list", branch); out != "" {
		return "", "", &Refusal{"branch_exists", "the branch " + branch + " already exists"}
	}
	return tag, branch, nil
}

// Backup creates the tag and switches to the new branch.
func Backup(root, tag, branch string) error {
	if out, err := git(root, "tag", tag); err != nil {
		return errors.New("git tag: " + out)
	}
	if out, err := git(root, "switch", "-c", branch); err != nil {
		return errors.New("git switch: " + out)
	}
	return nil
}

type logLine struct {
	Time  string   `json:"time"`
	From  string   `json:"from"`
	To    string   `json:"to"`
	Step  string   `json:"step"`
	Count int      `json:"count"`
	Files []string `json:"files"`
	Notes []string `json:"notes,omitempty"`
}

// AppendLog adds one line per applied step to the migration log.
func AppendLog(root, logPath, from, to string, results []Result, now time.Time) error {
	f, err := os.OpenFile(filepath.Join(root, filepath.FromSlash(logPath)), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, r := range results {
		if r.Status != "applied" {
			continue
		}
		b, _ := json.Marshal(logLine{now.UTC().Format("2006-01-02T15:04:05Z"), from, to, r.Step.ID, len(r.Files), r.Files, r.Notes})
		if _, err := f.Write(append(b, '\n')); err != nil {
			return err
		}
	}
	return nil
}

// SetLayout writes the layout version into the manifest, keeping everything else.
func SetLayout(root, manifestPath, version string) error {
	p := filepath.Join(root, filepath.FromSlash(manifestPath))
	b, err := os.ReadFile(p)
	if err != nil {
		b = nil
	}
	doc, err := frontmatter.ParseBare(b)
	if err != nil {
		return err
	}
	if err := doc.Set("layout", frontmatter.NewString(version)); err != nil {
		return err
	}
	return fsx.WriteAtomic(p, doc.Bytes())
}
