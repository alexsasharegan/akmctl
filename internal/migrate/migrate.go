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
// set_default, remove_field, copy_field, title_from_h1, goal_scaffold,
// rename_heading, move_file, move_closed, add_file, mint_ids) or "manual".
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
	Dir         string   `json:"dir,omitempty"`        // the directory, repository-relative (mint_ids; the field steps when set)
	DateField   string   `json:"date_field,omitempty"` // mint_ids: the field the id date comes from, default created
	Recursive   bool     `json:"recursive,omitempty"`  // also the subdirectories
	Statuses    []string `json:"statuses,omitempty"`   // goal_scaffold: only records with one of these statuses
	Rename      *bool    `json:"rename,omitempty"`     // mint_ids: rename the files (default true)
	Seed        string   `json:"seed,omitempty"`       // mint_ids: makes the id suffixes repeatable
	SlugFrom    string   `json:"slug_from,omitempty"`  // mint_ids: name (default) or title
	MaxSlug     int      `json:"max_slug,omitempty"`   // mint_ids: the longest slug, 1 to 100 (default 60)
}

// Rename is one file move a step planned or made, repository-relative.
type Rename struct {
	From string `json:"from"`
	To   string `json:"to"`
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
func Chain(from, to string) ([]Migration, bool) { return ChainIn(Registry, from, to) }

// ChainIn is Chain over the given migrations (the embedded ones plus any a
// caller supplies).
func ChainIn(registry []Migration, from, to string) ([]Migration, bool) {
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
		for _, m := range registry {
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
	Step    Step
	Files   []string
	Notes   []string // what the step decided or could not do (mint_ids: old name, new name, id, references)
	Renames []Rename // mint_ids: the files renamed, or with rename false the renames left to the caller
	Status  string   // pending, applied or skipped
}

// Plan lists what each step would change, without writing.
func Plan(repo *record.Repo, steps []Step) []Result { return PlanEnv(repo, steps, Env{}) }

// PlanEnv is Plan with the random source and the clock given. The steps are
// planned in order on a scratch copy of the repository's Markdown files, each
// applied there after it is planned, so a step that depends on an earlier one
// (a date field renamed before ids are minted from it, a legacy id renamed
// before ids are minted) is planned as it will run. Nothing in the repository
// is written. If the scratch copy cannot be made, each step is planned alone
// against the working tree.
func PlanEnv(repo *record.Repo, steps []Step, env Env) []Result {
	out := make([]Result, 0, len(steps))
	sc, cleanup := scratchRepo(repo)
	if cleanup != nil {
		defer cleanup()
	}
	for _, s := range steps {
		files, notes, renames := []string{}, []string(nil), []Rename(nil)
		if s.Class == "auto" {
			switch {
			case sc == nil || s.Kind == "move_file" || s.Kind == "add_file":
				// these touch files that are not Markdown: planned against the real tree
				files, notes, renames, _ = run(repo, s, true, env)
			default:
				files, notes, renames, _ = run(sc, s, false, env)
				sc.Reload()
			}
		}
		out = append(out, Result{Step: s, Files: files, Notes: notes, Renames: renames, Status: "pending"})
	}
	return out
}

// scratchRepo copies the manifest and every Markdown file of the repository to a
// temporary directory and loads it. It returns nil, nil when that fails.
func scratchRepo(repo *record.Repo) (*record.Repo, func()) {
	dir, err := os.MkdirTemp("", "akmctl-plan-")
	if err != nil {
		return nil, nil
	}
	cleanup := func() { os.RemoveAll(dir) }
	cp := func(rel string) error {
		b, err := os.ReadFile(filepath.Join(repo.Root, filepath.FromSlash(rel)))
		if err != nil {
			return err
		}
		dst := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		return os.WriteFile(dst, b, 0o644)
	}
	if err := cp(repo.ManifestPath()); err != nil && !os.IsNotExist(err) {
		cleanup()
		return nil, nil
	}
	werr := filepath.WalkDir(repo.Root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if p != repo.Root && (d.Name() == ".git" || d.Name() == "node_modules" || strings.HasPrefix(d.Name(), ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type().IsRegular() && strings.EqualFold(filepath.Ext(d.Name()), ".md") {
			rel, _ := filepath.Rel(repo.Root, p)
			return cp(filepath.ToSlash(rel))
		}
		return nil
	})
	if werr != nil {
		cleanup()
		return nil, nil
	}
	sc, err := record.Load(dir)
	if err != nil {
		cleanup()
		return nil, nil
	}
	return sc, cleanup
}

// Apply runs the auto steps in order. Each step is idempotent: when there is
// nothing to change it is "skipped". Manual and info steps stay "pending".
func Apply(repo *record.Repo, steps []Step) ([]Result, error) {
	return ApplyEnv(repo, steps, Env{})
}

// ApplyEnv is Apply with the random source and the clock given.
func ApplyEnv(repo *record.Repo, steps []Step, env Env) ([]Result, error) {
	var out []Result
	for _, s := range steps {
		if s.Class != "auto" {
			out = append(out, Result{Step: s, Files: []string{}, Status: "pending"})
			continue
		}
		files, notes, renames, err := run(repo, s, false, env)
		if err != nil {
			return out, fmt.Errorf("step %s: %w", s.ID, err)
		}
		st := "skipped"
		if len(files) > 0 {
			st = "applied"
		}
		out = append(out, Result{Step: s, Files: files, Notes: notes, Renames: renames, Status: st})
		if err := repo.Reload(); err != nil {
			return out, err
		}
	}
	return out, nil
}

// run applies (or, with dry, only lists) one step and returns the files it
// changes, as repository-relative paths, and the notes it has to report.
func run(repo *record.Repo, s Step, dry bool, env Env) ([]string, []string, []Rename, error) {
	if s.Kind == "mint_ids" {
		return mintIDs(repo, s, dry, env)
	}
	files, notes, err := runBasic(repo, s, dry)
	return files, notes, nil, err
}

func runBasic(repo *record.Repo, s Step, dry bool) ([]string, []string, error) {
	files := []string{}
	var notes []string
	write := func(rel string, data []byte) error {
		files = append(files, rel)
		if dry {
			return nil
		}
		return fsx.WriteAtomic(filepath.Join(repo.Root, filepath.FromSlash(rel)), data)
	}
	switch s.Kind {
	case "rename_field", "set_default", "remove_field", "copy_field", "title_from_h1":
		var err error
		notes, err = fieldStep(repo, s, write)
		if err != nil {
			return files, notes, err
		}
	case "move_closed":
		var err error
		files, notes, err = moveClosed(repo, s, dry)
		if err != nil {
			return files, notes, err
		}
	case "goal_scaffold":
		var err error
		notes, err = goalScaffold(repo, s, write)
		if err != nil {
			return files, notes, err
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
				return files, notes, err
			}
		}
	case "move_file":
		from, to := filepath.Join(repo.Root, filepath.FromSlash(s.From)), filepath.Join(repo.Root, filepath.FromSlash(s.To))
		if !safeRel(s.From) || !safeRel(s.To) {
			return files, notes, fmt.Errorf("path escapes the repository")
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
				return files, notes, err
			}
			if err := os.Rename(from, to); err != nil {
				return files, notes, err
			}
		}
	case "add_file":
		if !safeRel(s.Path) {
			return files, notes, fmt.Errorf("path escapes the repository")
		}
		if _, err := os.Stat(filepath.Join(repo.Root, filepath.FromSlash(s.Path))); err == nil {
			break // never overwrite
		}
		if err := write(s.Path, []byte(s.Content)); err != nil {
			return files, notes, err
		}
	default:
		return files, notes, fmt.Errorf("unknown step kind %q", s.Kind)
	}
	sort.Strings(files)
	return files, notes, nil
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

var versionRe = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

var stepKinds = map[string]bool{"rename_field": true, "set_default": true, "remove_field": true, "copy_field": true, "title_from_h1": true, "goal_scaffold": true,
	"rename_heading": true, "move_file": true, "add_file": true, "mint_ids": true, "move_closed": true, "manual": true}

// Validate checks a migration that did not come from the embedded set: versions,
// step ids, kinds, classes and the parameters each kind needs.
func Validate(m Migration) error {
	if !versionRe.MatchString(m.From) || !versionRe.MatchString(m.To) {
		return fmt.Errorf("from and to must be versions like 0.2.0")
	}
	if len(m.Steps) == 0 {
		return fmt.Errorf("a migration needs at least one step")
	}
	seen := map[string]bool{}
	for i, s := range m.Steps {
		at := fmt.Sprintf("step %d", i+1)
		if s.ID == "" || seen[s.ID] {
			return fmt.Errorf("%s: id is missing or used twice", at)
		}
		seen[s.ID] = true
		at = "step " + s.ID
		if !stepKinds[s.Kind] {
			return fmt.Errorf("%s: unknown kind %q", at, s.Kind)
		}
		if s.Class != "auto" && s.Class != "manual" && s.Class != "info" {
			return fmt.Errorf("%s: class must be auto, manual or info", at)
		}
		need := map[string][]string{}
		switch s.Kind {
		case "rename_field", "copy_field":
			need["from"], need["to"] = []string{s.From}, []string{s.To}
		case "set_default":
			need["field"] = []string{s.Field}
		case "remove_field":
			need["field"] = []string{s.Field}
		case "rename_heading":
			need["from"], need["to"] = []string{s.From}, []string{s.To}
		case "move_file":
			need["from"], need["to"] = []string{s.From}, []string{s.To}
		case "add_file":
			need["path"] = []string{s.Path}
		case "mint_ids", "move_closed":
			need["dir"] = []string{s.Dir}
			if s.SlugFrom != "" && s.SlugFrom != "name" && s.SlugFrom != "title" {
				return fmt.Errorf("%s: slug_from must be name or title", at)
			}
			if s.MaxSlug < 0 || s.MaxSlug > 100 {
				return fmt.Errorf("%s: max_slug must be between 1 and 100", at)
			}
		}
		for k, v := range need {
			if v[0] == "" {
				return fmt.Errorf("%s: %s is required for %s", at, k, s.Kind)
			}
		}
		for _, p := range []string{s.Dir, s.Path} {
			if p != "" && !safeRel(p) {
				return fmt.Errorf("%s: %q must stay inside the repository", at, p)
			}
		}
		if s.Kind == "move_file" && (!safeRel(s.From) || !safeRel(s.To)) {
			return fmt.Errorf("%s: paths must stay inside the repository", at)
		}
	}
	return nil
}
