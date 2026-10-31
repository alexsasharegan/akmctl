package initrepo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/fs"
	"sort"
	"strings"
)

// The managed skills and their markers (ADR 261030-h6t7).

// MarkerName is the file in each managed skill folder that records what the tool
// wrote there (ADR 261030-1wj5 names it).
const MarkerName = ".akm-skill.json"

// Skill is one managed skill as the binary ships it. Its files have paths
// relative to the skill folder.
type Skill struct {
	Name  string
	Files []File
}

// Rename says that a skill folder was renamed; the installer retires the old one.
type Rename struct{ Old, New string }

// Renames lists the skill renames the tool knows about (ADR 261030-wsqb).
var Renames = []Rename{{Old: "akm-work-items", New: "akm-track-work"}}

// Managed returns the skills the tool manages in a repository: the router, the
// skills `init` always writes, and the optional akm-commit when gitCommits is
// true. The folders of the repository replace docs/work and docs/adr in the text.
func Managed(workDir, adrDir string, gitCommits bool) ([]Skill, error) {
	paths := strings.NewReplacer("docs/work", workDir, "docs/adr", adrDir)
	entries, err := fs.ReadDir(templates, "templates/skills")
	if err != nil {
		return nil, err
	}
	var out []Skill
	read := func(root, name string) error {
		sk := Skill{Name: name}
		werr := fs.WalkDir(templates, root, func(fp string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			b, err := templates.ReadFile(fp)
			if err != nil {
				return err
			}
			sk.Files = append(sk.Files, File{strings.TrimPrefix(fp, root+"/"), []byte(paths.Replace(string(b)))})
			return nil
		})
		if werr != nil {
			return werr
		}
		sort.Slice(sk.Files, func(i, j int) bool { return sk.Files[i].Path < sk.Files[j].Path })
		out = append(out, sk)
		return nil
	}
	for _, e := range entries {
		if e.IsDir() {
			if err := read("templates/skills/"+e.Name(), e.Name()); err != nil {
				return nil, err
			}
		}
	}
	if gitCommits {
		if err := read("templates/optional-skills/"+CommitSkill, CommitSkill); err != nil {
			return nil, err
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// Hash is the hex SHA-256 of data, as the marker records it.
func Hash(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// Marker is the content of .akm-skill.json: the tool version and the hash of
// each file the tool wrote, by path in the skill folder. It is written in the
// form Prettier keeps (two-space indent, a final newline).
func (s Skill) Marker(tool string) []byte {
	files := map[string]string{}
	for _, f := range s.Files {
		files[f.Path] = Hash(f.Data)
	}
	b, _ := json.MarshalIndent(struct {
		Tool  string            `json:"tool"`
		Files map[string]string `json:"files"`
	}{tool, files}, "", "  ")
	return append(b, '\n')
}

// ParseMarker reads a marker. ok is false when the bytes are not one.
func ParseMarker(b []byte) (tool string, files map[string]string, ok bool) {
	var m struct {
		Tool  string            `json:"tool"`
		Files map[string]string `json:"files"`
	}
	if err := json.Unmarshal(b, &m); err != nil || m.Files == nil {
		return "", nil, false
	}
	return m.Tool, m.Files, true
}

// RepoFiles returns the skill's files with repository-relative paths, and the
// marker, sorted.
func (s Skill) RepoFiles(tool string) []File {
	dir := SkillDir(s.Name) + "/"
	out := make([]File, 0, len(s.Files)+1)
	for _, f := range s.Files {
		out = append(out, File{dir + f.Path, f.Data})
	}
	out = append(out, File{dir + MarkerName, s.Marker(tool)})
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}
