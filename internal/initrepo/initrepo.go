// Package initrepo creates the AKM layout in a new repository: the files of a
// profile, written without ever overwriting one that exists.
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

// Profiles lists the profile names.
var Profiles = []string{"code", "notes"}

// Profile is the shape of a new repository.
type Profile struct {
	Name                 string
	WorkDir, ADRDir, Map string // Map is the roadmap file
	WikiLinks            bool
}

// Get returns the profile called name.
func Get(name string) (Profile, bool) {
	switch name {
	case "code":
		return Profile{Name: "code", WorkDir: record.DefaultWorkDir, ADRDir: record.DefaultADRDir, Map: record.DefaultRoadmap}, true
	case "notes":
		return Profile{Name: "notes", WorkDir: "notes/work", ADRDir: "notes/adr", Map: "notes/ROADMAP.md", WikiLinks: true}, true
	}
	return Profile{}, false
}

// File is one file of a profile, with a path relative to the repository root.
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

// Files returns the files of profile p for a project called name. A name that
// would need escaping in Markdown is replaced by "project".
func Files(p Profile, name string) ([]File, error) {
	if !plainName.MatchString(name) {
		name = "project"
	}
	tmpl := func(file string) (string, error) {
		b, err := templates.ReadFile("templates/" + file)
		return string(b), err
	}
	var commands, arch, about, extra string
	if p.Name == "code" {
		commands = strings.Join([]string{
			"- Build: `<build command>`",
			"- Test: `<test command>`",
			"- Lint and format: `<lint command>`",
			"- Run: `<run command>`",
			"- Records check: `akmctl check`; `akmctl index` regenerates `" + p.ADRDir + "/README.md`.",
			"- Format records: `npx prettier@<pinned version> --check AGENTS.md \"docs/**/*.md\"` (pin the version you use).",
			"- Full check before a commit: the test command, `akmctl check` and the format command above.",
		}, "\n")
		arch = "- `docs/architecture.md`: the map of the parts and the invariants (write it once there is more than one part).\n"
		about = "Decisions, direction and work are kept as Markdown records under `docs/`, in the AKM layout; start with [AGENTS.md](AGENTS.md)."
	} else {
		commands = strings.Join([]string{
			"- Records check: `akmctl check`; `akmctl index` regenerates `" + p.ADRDir + "/README.md`.",
			"- Format records: `npx prettier@<pinned version> --check AGENTS.md \"notes/**/*.md\"` (pin the version you use).",
			"- Full check before a commit: `akmctl check` and the format command above.",
		}, "\n")
		extra = "- `docs/.layout`: the manifest. `reference_style: wikilink` makes the tool write Obsidian links; keep Obsidian's own settings out of the records.\n"
		about = "The records (work, decisions and the roadmap) live in the AKM layout under `notes/`; start with [AGENTS.md](AGENTS.md). Open this folder as an Obsidian vault."
	}
	fill := strings.NewReplacer("{{NAME}}", name, "{{COMMANDS}}", commands, "{{ARCH}}\n", arch, "{{EXTRA}}\n", extra,
		"{{ABOUT}}", about, "{{WORK}}", p.WorkDir, "{{ADR}}", p.ADRDir, "{{ROADMAP}}", p.Map)

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

	gitignore := "# the advisory lock akmctl takes while it writes\ndocs/.akm.lock\n"
	if p.WikiLinks {
		gitignore += "\n# Obsidian state that differs per machine\n.obsidian/workspace.json\n.obsidian/workspaces.json\n.trash/\n"
	}
	add(".gitignore", gitignore)
	add(".prettierrc.json", "{ \"embeddedLanguageFormatting\": \"off\" }\n")
	add(".prettierignore", p.ADRDir+"/README.md\n")

	manifest := "layout: " + check.LayoutVersion + "\n"
	for _, ext := range []string{"adr-index", "work-items", "claims", "safety"} {
		manifest += "extension." + ext + ": " + check.LayoutVersion + "\n"
	}
	if p.WorkDir != record.DefaultWorkDir {
		manifest += "work_dir: " + p.WorkDir + "\n"
	}
	if p.ADRDir != record.DefaultADRDir {
		manifest += "adr_dir: " + p.ADRDir + "\n"
	}
	if p.Map != record.DefaultRoadmap {
		manifest += "roadmap: " + p.Map + "\n"
	}
	if p.WikiLinks {
		manifest += "reference_style: wikilink\n"
	}
	add(record.ManifestPath, manifest)
	deny, _ := tmpl("akm-deny.tmpl")
	add("docs/.akm-deny", deny)

	roadmap, _ := tmpl("ROADMAP.md.tmpl")
	add(p.Map, roadmap)
	add(p.ADRDir+"/README.md", index.Generate(&record.Repo{}))
	add(p.WorkDir+"/.gitkeep", "")
	add(p.WorkDir+"/done/.gitkeep", "")

	paths := strings.NewReplacer("docs/work", p.WorkDir, "docs/adr", p.ADRDir)
	err = fs.WalkDir(templates, "templates/skills", func(fp string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := templates.ReadFile(fp)
		if err != nil {
			return err
		}
		add(".agents/"+strings.TrimPrefix(fp, "templates/"), paths.Replace(string(b)))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

// Init writes the files of the profile under target, creating the directory
// when it is absent. A file that exists, whatever it holds, is left alone and
// listed as skipped.
func Init(target string, p Profile) (Result, error) {
	var res Result
	if st, err := os.Stat(target); err == nil && !st.IsDir() {
		return res, fmt.Errorf("%s exists and is not a directory", target)
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return res, err
	}
	files, err := Files(p, filepath.Base(abs))
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
