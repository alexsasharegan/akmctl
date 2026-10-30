package check

import (
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alexsasharegan/akmctl/internal/record"
)

// The safety extension: SF1 (built-in secret patterns), SF2 (the repo's deny
// list), SF3 (high-entropy strings) and SF4 (inline allow markers). A finding
// names the file, line and rule and the kind of pattern, never the matched text.

type pattern struct {
	name string
	re   *regexp.Regexp
}

var secretPatterns = []pattern{
	{"a private key header", regexp.MustCompile(`-----BEGIN (?:[A-Z]+ )?PRIVATE KEY(?: BLOCK)?-----`)},
	{"a cloud access key id", regexp.MustCompile(`\b(?:AKIA|ASIA|AGPA|AIDA|AROA)[0-9A-Z]{16}\b`)},
	{"a GitHub token", regexp.MustCompile(`\b(?:gh[pousr]_[A-Za-z0-9]{36,}|github_pat_[A-Za-z0-9_]{22,})\b`)},
	{"a Slack token", regexp.MustCompile(`\bxox[abposr]-[A-Za-z0-9-]{10,}`)},
	{"a payment provider live key", regexp.MustCompile(`\b[sr]k_live_[0-9A-Za-z]{16,}`)},
	{"a Google API key", regexp.MustCompile(`\bAIza[0-9A-Za-z_-]{35}\b`)},
	{"an npm token", regexp.MustCompile(`\bnpm_[A-Za-z0-9]{36}\b`)},
	{"a JSON Web Token", regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{8,}\.eyJ[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,}`)},
}

var assignment = regexp.MustCompile(`(?i)\b(?:password|passwd|pwd|secret|api[_-]?key|access[_-]?token|auth[_-]?token)\b["']?\s*[:=]\s*["']?([^\s"'<>{}$%,;]{8,})`)

var placeholder = regexp.MustCompile(`(?i)example|changeme|placeholder|your[_-]|redacted|dummy|sample|\*{3,}|x{4,}|\.{3}|^(?:true|false|null|none|string|password|secret)$`)

var allowMarker = regexp.MustCompile(`akm-allow:\s*(SF[123])\b[ \t]*(.*)$`)
var hashPrefix = regexp.MustCompile(`(?:\bh1:|\bsha(?:1|256|384|512)[-:=]|\bintegrity\W+)$`)
var urlSpan = regexp.MustCompile(`https?://[^\s)>\]"']+`)
var entropyToken = regexp.MustCompile(`[A-Za-z0-9+/=_-]{32,}`)

// A repository path is public by design: the target of a Markdown link, or a record file name
// (YYMMDD-xxxx-slug, with an optional folder prefix and .md).
var linkTargetSpan = regexp.MustCompile(`\]\([^)\n]*\)`)
var recordName = regexp.MustCompile(`^(?:[A-Za-z0-9_.-]+/)*[0-9]{6}-[0-9a-hjkmnp-tv-z]{4}-[0-9a-z-]+(?:\.md)?$`)

type denyEntry struct {
	n   int
	lit string
	re  *regexp.Regexp
}

func (c *checker) safety() {
	deny := c.denyList()
	files := []string{"AGENTS.md", c.repo.RoadmapPath()}
	for _, r := range c.repo.Records {
		files = append(files, r.Path)
	}
	for _, f := range files {
		b, err := os.ReadFile(filepath.Join(c.repo.Root, filepath.FromSlash(f)))
		if err != nil {
			continue
		}
		for i, line := range strings.Split(string(b), "\n") {
			c.safetyLine(f, i+1, strings.TrimRight(line, "\r"), deny)
		}
	}
}

func (c *checker) denyList() []denyEntry {
	b, err := os.ReadFile(filepath.Join(c.repo.Root, filepath.FromSlash(c.repo.DenyPath())))
	if err != nil {
		return nil
	}
	var out []denyEntry
	for i, l := range strings.Split(string(b), "\n") {
		l = strings.TrimRight(l, "\r")
		if strings.TrimSpace(l) == "" || strings.HasPrefix(l, "#") {
			continue
		}
		if strings.HasPrefix(l, "re:") {
			if re, err := regexp.Compile(l[3:]); err == nil {
				out = append(out, denyEntry{n: i + 1, re: re})
			} else {
				c.add("SF2", "E", c.repo.DenyPath(), i+1, "", "the regular expression on this line does not compile")
			}
			continue
		}
		out = append(out, denyEntry{n: i + 1, lit: l})
	}
	return out
}

func (c *checker) safetyLine(file string, n int, line string, deny []denyEntry) {
	allowed := map[string]bool{}
	for _, m := range allowMarker.FindAllStringSubmatch(line, -1) {
		reason := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(m[2]), "-->"))
		if reason == "" {
			c.add("SF4", "E", file, n, "", "an allow marker for %s needs a reason", m[1])
			continue
		}
		allowed[m[1]] = true
		c.add("SF4", "W", file, n, "", "allow marker for %s: %s", m[1], reason)
	}
	hit := false
	if !allowed["SF1"] {
		for _, p := range secretPatterns {
			if p.re.MatchString(line) {
				c.add("SF1", "E", file, n, "", "matches the pattern for %s", p.name)
				hit = true
			}
		}
		for _, m := range assignment.FindAllStringSubmatch(line, -1) {
			v := m[1]
			if placeholder.MatchString(v) || classes(v) < 2 {
				continue
			}
			c.add("SF1", "E", file, n, "", "matches the pattern for a credential assignment")
			hit = true
			break
		}
	}
	if !allowed["SF2"] {
		for _, d := range deny {
			if (d.re != nil && d.re.MatchString(line)) || (d.lit != "" && strings.Contains(line, d.lit)) {
				c.add("SF2", "E", file, n, "", "matches line %d of %s", d.n, c.repo.DenyPath())
				hit = true
			}
		}
	}
	if !hit && !allowed["SF3"] {
		urls := urlSpan.FindAllStringIndex(line, -1)
		urls = append(urls, linkTargetSpan.FindAllStringIndex(line, -1)...)
		for _, loc := range entropyToken.FindAllStringIndex(line, -1) {
			if inSpan(urls, loc[0]) {
				continue // ids in URLs (playlists, document ids) are public by design; SF1 still sees known token formats
			}
			if recordName.MatchString(line[loc[0]:loc[1]]) {
				continue // a record's file name, not a secret
			}
			if hashPrefix.MatchString(line[:loc[0]]) {
				continue // go.sum and package-lock integrity values are hashes
			}
			if highEntropy(line[loc[0]:loc[1]]) {
				c.add("SF3", "W", file, n, "", "holds a high-entropy string that looks like a secret")
				break
			}
		}
	}
}

// classes counts the character classes in s: lower, upper, digit, other.
func classes(s string) int {
	var lo, up, di, ot bool
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
			lo = true
		case r >= 'A' && r <= 'Z':
			up = true
		case r >= '0' && r <= '9':
			di = true
		default:
			ot = true
		}
	}
	n := 0
	for _, b := range []bool{lo, up, di, ot} {
		if b {
			n++
		}
	}
	return n
}

func highEntropy(s string) bool {
	if classes(s) < 3 || strings.ContainsAny(s, "/") && strings.Count(s, "/") > 2 {
		return false
	}
	if isHex(s) {
		return false // a hash, not a credential
	}
	freq := map[rune]float64{}
	for _, r := range s {
		freq[r]++
	}
	var h float64
	n := float64(len([]rune(s)))
	for _, f := range freq {
		p := f / n
		h -= p * math.Log2(p)
	}
	return h >= 4.5
}

func isHex(s string) bool {
	for _, r := range s {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F') {
			return false
		}
	}
	return true
}

var _ = record.Work

func inSpan(spans [][]int, i int) bool {
	for _, sp := range spans {
		if i >= sp[0] && i < sp[1] {
			return true
		}
	}
	return false
}
