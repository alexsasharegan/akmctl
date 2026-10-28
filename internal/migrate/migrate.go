// Package migrate reconciles a repository with a layout version: it plans the
// steps from the manifest's version to a target, applies the safe ones, and
// writes the handoff for the rest. The steps are data: each release of the tool
// embeds the migrations it supports.
package migrate

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/frontmatter"
	"github.com/alexsasharegan/akmctl/internal/fsx"
	"github.com/alexsasharegan/akmctl/internal/record"
)

//go:embed migrations
var embedded embed.FS

// Step is one declarative step. Kind is one of the primitives (rename_field,
// set_default, rename_heading, move_file, add_file) or "manual".
type Step struct {
	ID          string   `json:"id"`
	Kind        string   `json:"kind"`
	Class       string   `json:"class"` // auto, manual or info
	Description string   `json:"description"`
	Fixes       []string `json:"fixes,omitempty"` // rule ids this step resolves, for classifying findings
	Scope       string   `json:"scope,omitempty"` // work, adr or all
	Field       string   `json:"field,omitempty"`
	From        string   `json:"from,omitempty"`
	To          string   `json:"to,omitempty"`
	Value       any      `json:"value,omitempty"`
	Path        string   `json:"path,omitempty"`
	Content     string   `json:"content,omitempty"`
}

// Migration is the steps from one layout version to the next.
type Migration struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Steps []Step `json:"steps"`
}

// Registry holds the migrations of this release; tests replace it.
var Registry = loadEmbedded()

func loadEmbedded() []Migration {
	var out []Migration
	ents, _ := embedded.ReadDir("migrations")
	for _, e := range ents {
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		b, _ := embedded.ReadFile("migrations/" + e.Name())
		var m Migration
		if err := json.Unmarshal(b, &m); err != nil {
			panic("embedded migration " + e.Name() + ": " + err.Error())
		}
		out = append(out, m)
	}
	return out
}

// Chain returns the migrations that lead from one version to another, or false
// when there is no path. From equal to To is an empty chain.
func Chain(from, to string) ([]Migration, bool) {
	if from == to {
		return nil, true
	}
	type node struct {
		v    string
		path []Migration
	}
	seen := map[string]bool{from: true}
	queue := []node{{from, nil}}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		for _, m := range Registry {
			if m.From != n.v || seen[m.To] {
				continue
			}
			p := append(append([]Migration(nil), n.path...), m)
			if m.To == to {
				return p, true
			}
			seen[m.To] = true
			queue = append(queue, node{m.To, p})
		}
	}
	return nil, false
}

// Steps flattens a chain.
func Steps(ms []Migration) []Step {
	var out []Step
	for _, m := range ms {
		out = append(out, m.Steps...)
	}
	return out
}

// Result is the outcome of one step.
type Result struct {
	Step   Step
	Files  []string
	Status string // pending, applied or skipped
}

// Plan lists what each step would change, without writing.
func Plan(repo *record.Repo, steps []Step) []Result {
	out := make([]Result, 0, len(steps))
	for _, s := range steps {
		files := []string{}
		if s.Class == "auto" {
			files, _ = run(repo, s, true)
		}
		out = append(out, Result{s, files, "pending"})
	}
	return out
}

// Apply runs the auto steps in order. Each step is idempotent: when there is
// nothing to change it is "skipped". Manual and info steps stay "pending".
func Apply(repo *record.Repo, steps []Step) ([]Result, error) {
	var out []Result
	for _, s := range steps {
		if s.Class != "auto" {
			out = append(out, Result{s, []string{}, "pending"})
			continue
		}
		files, err := run(repo, s, false)
		if err != nil {
			return out, fmt.Errorf("step %s: %w", s.ID, err)
		}
		st := "skipped"
		if len(files) > 0 {
			st = "applied"
		}
		out = append(out, Result{s, files, st})
		if err := repo.Reload(); err != nil {
			return out, err
		}
	}
	return out, nil
}

// run applies (or, with dry, only lists) one step and returns the files it
// changes, as repository-relative paths.
func run(repo *record.Repo, s Step, dry bool) ([]string, error) {
	files := []string{}
	write := func(rel string, data []byte) error {
		files = append(files, rel)
		if dry {
			return nil
		}
		return fsx.WriteAtomic(filepath.Join(repo.Root, filepath.FromSlash(rel)), data)
	}
	switch s.Kind {
	case "rename_field", "set_default":
		for _, r := range repo.Records {
			if s.Scope != "" && s.Scope != "all" && s.Scope != r.Type {
				continue
			}
			doc, err := frontmatter.Parse(r.Text)
			if err != nil {
				continue
			}
			if s.Kind == "rename_field" {
				_, hasFrom := doc.Get(s.From)
				_, hasTo := doc.Get(s.To)
				if !hasFrom || hasTo {
					continue
				}
				if err := doc.Rename(s.From, s.To); err != nil {
					return files, err
				}
			} else {
				if _, ok := doc.Get(s.Field); ok {
					continue
				}
				v, err := toValue(s.Value)
				if err != nil {
					return files, err
				}
				if err := doc.Set(s.Field, v); err != nil {
					return files, err
				}
			}
			if err := write(r.Path, doc.Bytes()); err != nil {
				return files, err
			}
		}
	case "rename_heading":
		re := regexp.MustCompile(`(?m)^## ` + regexp.QuoteMeta(s.From) + `[ \t]*(\r?)$`)
		for _, r := range repo.Records {
			if s.Scope != "" && s.Scope != "all" && s.Scope != r.Type {
				continue
			}
			doc, err := frontmatter.Parse(r.Text)
			if err != nil {
				continue
			}
			body := string(doc.Body())
			nb := replaceOutsideFences(re, body, "## "+s.To)
			if nb == body {
				continue
			}
			b := doc.Bytes()
			head := b[:len(b)-len(doc.Body())]
			if err := write(r.Path, append(append([]byte(nil), head...), nb...)); err != nil {
				return files, err
			}
		}
	case "move_file":
		from, to := filepath.Join(repo.Root, filepath.FromSlash(s.From)), filepath.Join(repo.Root, filepath.FromSlash(s.To))
		if !safeRel(s.From) || !safeRel(s.To) {
			return files, fmt.Errorf("path escapes the repository")
		}
		if _, err := os.Stat(from); err != nil {
			break
		}
		if _, err := os.Stat(to); err == nil {
			break // never overwrite
		}
		files = append(files, s.To)
		if !dry {
			if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
				return files, err
			}
			if err := os.Rename(from, to); err != nil {
				return files, err
			}
		}
	case "add_file":
		if !safeRel(s.Path) {
			return files, fmt.Errorf("path escapes the repository")
		}
		if _, err := os.Stat(filepath.Join(repo.Root, filepath.FromSlash(s.Path))); err == nil {
			break // never overwrite
		}
		if err := write(s.Path, []byte(s.Content)); err != nil {
			return files, err
		}
	default:
		return files, fmt.Errorf("unknown step kind %q", s.Kind)
	}
	sort.Strings(files)
	return files, nil
}

func safeRel(p string) bool {
	c := path.Clean(p)
	return p != "" && !path.IsAbs(c) && c != ".." && !strings.HasPrefix(c, "../")
}

func replaceOutsideFences(re *regexp.Regexp, body, repl string) string {
	lines := strings.Split(body, "\n")
	fence := false
	for i, l := range lines {
		t := strings.TrimRight(l, "\r")
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			fence = !fence
			continue
		}
		if !fence && re.MatchString(l) {
			cr := ""
			if strings.HasSuffix(l, "\r") {
				cr = "\r"
			}
			lines[i] = repl + cr
		}
	}
	return strings.Join(lines, "\n")
}

func toValue(v any) (frontmatter.Value, error) {
	switch x := v.(type) {
	case nil:
		return frontmatter.NewNull(), nil
	case string:
		return frontmatter.NewString(x), nil
	case bool:
		return frontmatter.NewBool(x), nil
	case float64:
		return frontmatter.NewInt(int64(x)), nil
	case int:
		return frontmatter.NewInt(int64(x)), nil
	}
	return frontmatter.Value{}, fmt.Errorf("unsupported default value %v", v)
}
