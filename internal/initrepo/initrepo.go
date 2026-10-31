// Package initrepo creates the AKM layout in a new repository from a few
// options, writing without ever overwriting a file that exists.
package initrepo

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/check"
	"github.com/alexsasharegan/akmctl/internal/index"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// The templates are plain files; {{NAME}}-style markers are replaced by strings.NewReplacer.
// templates/skills is a byte-for-byte copy of .agents/skills of this repository
// (a test keeps the two equal).
//
//go:embed all:templates
var templates embed.FS

// DefaultDocsDir is the folder for the records when init is given none.
const DefaultDocsDir = "docs"

// ValidateDocsDir reports why name cannot be the docs folder: it must be one
// path segment that is safe as a folder name on macOS, Linux and Windows and
// as a line of an ignore file.
func ValidateDocsDir(name string) error {
	switch {
	case name == "":
		return errors.New("is empty")
	case name == "." || name == "..":
		return errors.New("must be a folder name, not \".\" or \"..\"")
	case strings.ContainsAny(name, "/\\"):
		return errors.New("must be a single folder name, with no path separator")
	case strings.HasPrefix(name, "."):
		return errors.New("must not start with a dot (a hidden folder is not found again)")
	case strings.HasPrefix(name, "#") || strings.HasPrefix(name, "!"):
		return errors.New("must not start with # or !")
	case name == "node_modules":
		return errors.New("is reserved")
	case name != strings.TrimSpace(name) || strings.HasSuffix(name, "."):
		return errors.New("must not start or end with a space, or end with a dot")
	}
	for _, c := range name {
		if c < 0x20 || c == 0x7f || strings.ContainsRune(`:*?"<>|[]`, c) {
			return fmt.Errorf("has a character (%q) that is not safe in a folder name", c)
		}
	}
	return nil
}

// File is one file of a new repository, with a path relative to the repository root.
type File struct {
	Path string
	Data []byte
}

// Result lists what a run wrote and what it left alone, as paths relative to
// the target directory, sorted.
type Result struct {
	Written, Skipped []string
}

var plainName = regexp.MustCompile(`^(?:[A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9 ._-]*[A-Za-z0-9])$`)

// Options are the choices of init.
type Options struct {
	DocsDir    string // the single folder for the records and the manifest; "" means docs
	WikiLinks  bool   // write reference_style: wikilink
	GitCommits bool   // write git_commits: true and install the optional skill akm-commit
	Gitignore  bool   // write .gitignore (the target is inside a work tree)
	Obsidian   bool   // the target has an .obsidian folder: ignore its per-machine files
	Tool       string // the tool version recorded in each skill's marker; "" means dev
}

// Files returns the files of a new repository called name with default options.
func Files(name string) ([]File, error) { return FilesWith(name, Options{}) }

// FilesWith returns the files of a new repository for a project called name. A
// name that would need escaping in Markdown is replaced by "project". The
// docs dir must be valid (ValidateDocsDir).
func FilesWith(name string, opt Options) ([]File, error) {
	if opt.DocsDir == "" {
		opt.DocsDir = DefaultDocsDir
	}
	if opt.Tool == "" {
		opt.Tool = "dev"
	}
	if err := ValidateDocsDir(opt.DocsDir); err != nil {
		return nil, fmt.Errorf("docs dir %q %v", opt.DocsDir, err)
	}
	dir := opt.DocsDir
	workDir, adrDir, roadmapPath := dir+"/work", dir+"/adr", dir+"/ROADMAP.md"
	if !plainName.MatchString(name) {
		name = "project"
	}
	tmpl := func(file string) (string, error) {
		b, err := templates.ReadFile("templates/" + file)
		return string(b), err
	}
	commands := strings.Join([]string{
		"- Records check: `akmctl check`; `akmctl index` regenerates `" + adrDir + "/README.md`.",
		"- Format records: `npx prettier@<pinned version> --check AGENTS.md \"" + dir + "/**/*.md\"` (pin the version you use).",
		"- Full check before a commit: `akmctl check` and the format command above.",
	}, "\n")
	extra := "- `" + dir + "/.layout`: the manifest."
	if opt.WikiLinks {
		extra += " `reference_style: wikilink` makes the tool write wikilinks."
	}
	extra += "\n"
	about := "The records (work, decisions and the roadmap) live in the AKM layout under `" + dir + "/`; start with [AGENTS.md](AGENTS.md)."
	fill := strings.NewReplacer("{{NAME}}", name, "{{COMMANDS}}", commands, "{{EXTRA}}\n", extra,
		"{{ABOUT}}", about, "{{WORK}}", workDir, "{{ADR}}", adrDir, "{{ROADMAP}}", roadmapPath, "{{DOCS}}", dir)

	var files []File
	add := func(rel, text string) { files = append(files, File{rel, []byte(text)}) }

	agents, err := tmpl("AGENTS.md.tmpl")
	if err != nil {
		return nil, err
	}
	add("AGENTS.md", fill.Replace(agents))
	add("CLAUDE.md", "@AGENTS.md\n")
	readme, _ := tmpl("README.md.tmpl")
	add("README.md", fill.Replace(readme))

	if opt.Gitignore {
		ignore := "# the advisory lock akmctl takes while it writes\n" + dir + "/.akm.lock\n"
		if opt.Obsidian {
			ignore += "\n# Obsidian state that differs per machine\n.obsidian/workspace.json\n.obsidian/workspaces.json\n.trash/\n"
		}
		add(".gitignore", ignore)
	}
	add(".prettierrc.json", "{ \"embeddedLanguageFormatting\": \"off\" }\n")
	add(".prettierignore", adrDir+"/README.md\n")

	manifest := "layout: " + check.LayoutVersion + "\n"
	for _, ext := range []string{"adr-index", "work-items", "claims", "safety"} {
		manifest += "extension." + ext + ": " + check.LayoutVersion + "\n"
	}
	if workDir != record.DefaultWorkDir {
		manifest += "work_dir: " + workDir + "\n"
	}
	if adrDir != record.DefaultADRDir {
		manifest += "adr_dir: " + adrDir + "\n"
	}
	if roadmapPath != record.DefaultRoadmap {
		manifest += "roadmap: " + roadmapPath + "\n"
	}
	if opt.WikiLinks {
		manifest += "reference_style: wikilink\n"
	}
	if opt.GitCommits {
		manifest += "git_commits: true\n"
	}
	add(dir+"/"+record.ManifestName, manifest)
	deny, _ := tmpl("akm-deny.tmpl")
	add(dir+"/.akm-deny", deny)

	roadmap, _ := tmpl("ROADMAP.md.tmpl")
	add(roadmapPath, roadmap)
	add(adrDir+"/README.md", index.Generate(&record.Repo{}))
	add(workDir+"/.gitkeep", "")
	add(workDir+"/done/.gitkeep", "")

	skills, err := Managed(workDir, adrDir, opt.GitCommits)
	if err != nil {
		return nil, err
	}
	for _, sk := range skills {
		files = append(files, sk.RepoFiles(opt.Tool)...)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

// Init writes the files of a new repository under target, creating the directory
// when it is absent. A file that exists, whatever it holds, is left alone and
// listed as skipped.
func Init(target string) (Result, error) { return InitWith(target, Options{}) }

// InitWith is Init with options.
func InitWith(target string, opt Options) (Result, error) {
	var res Result
	if st, err := os.Stat(target); err == nil && !st.IsDir() {
		return res, fmt.Errorf("%s exists and is not a directory", target)
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return res, err
	}
	files, err := FilesWith(filepath.Base(abs), opt)
	if err != nil {
		return res, err
	}
	for _, f := range files {
		full := filepath.Join(target, filepath.FromSlash(f.Path))
		if _, err := os.Lstat(full); err == nil {
			res.Skipped = append(res.Skipped, f.Path)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return res, err
		}
		out, err := os.OpenFile(full, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if errors.Is(err, fs.ErrExist) {
			res.Skipped = append(res.Skipped, f.Path)
			continue
		}
		if err != nil {
			return res, err
		}
		_, werr := out.Write(f.Data)
		if cerr := out.Close(); werr == nil {
			werr = cerr
		}
		if werr != nil {
			return res, werr
		}
		res.Written = append(res.Written, f.Path)
	}
	return res, nil
}
