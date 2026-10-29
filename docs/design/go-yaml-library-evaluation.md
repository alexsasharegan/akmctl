---
type: design
status: draft
created_at: 2026-10-28
updated_at: 2026-10-28
tags:
    - pkm/rails
---
# Go YAML library evaluation for `akmctl` (draft for review)

Evaluation for [[261027-wkxv-go-yaml-library-evaluation]], serving [[agent-ready-repo-spec]] (H5, section 0) and [[261027-ensf-akm-an-agent-knowledge-management]]. It answers one question: how should the Go tool read and edit frontmatter? Criteria, in Alexander's order: (1) baseline requirements, (2) performance, (3) security, (4) popularity and support as a tiebreaker. Run on 2026-10-28; machine and toolchain in the benchmark section.

## Summary

**Recommendation: keep a hand-written parser for the H5 subset. Do not add a YAML library to `akmctl`.** Reading and editing are one component, a line-oriented parser that records byte offsets, so an edit is a splice into the original bytes.

- **Baseline (1):** this is decisive. No library can leave the rest of a block byte-for-byte unchanged on edit. All four tested (both `go.yaml.in/yaml` majors, the archived `gopkg.in/yaml.v3`, `goccy/go-yaml`) rewrote spacing, blank lines and list indentation on a no-op round trip. `goccy/go-yaml` also dropped a trailing comment on a value edit and silently failed to rename a key. The libraries also accept what H5 forbids (anchors, aliases, tags, nesting, block scalars, and in node mode duplicate keys), and they retype values (`2026-10-28` becomes a timestamp, `08540` a float), so the tool would need its own subset check on top. The hand parser rejects all of that with a file line number.
- **Performance (2):** the hand parser took 2.8 us per file against 17.4 us for the fastest library (6 times), with 7.7 times fewer allocations. 5,000 files parse in 14 ms against 89 to 115 ms. Absolute numbers are small either way; this criterion does not decide the matter, it only confirms there is no cost to going without.
- **Security (3):** the hand parser is linear, has no recursion and refuses anchors and nesting, so alias bombs and deep nesting cannot occur. Among libraries, `go.yaml.in/yaml/v3` is the best behaved. `goccy/go-yaml` used over 6 GB (killed by my harness) on 100,000 nested `[`, and over 3 GB for one 100 MB scalar. See the security section.
- **Popularity (4):** it only breaks a tie and there is no tie. For reference: `gopkg.in/yaml.v3` has 44,731 importers but is archived and unmaintained; its maintained successor `go.yaml.in/yaml/v3` has 555; `goccy/go-yaml` has 2,983.

**If a library is wanted later** (for example for richer config files that are not frontmatter): `go.yaml.in/yaml/v3` v3.0.5. Stable API, security fixes only, no runtime dependencies, bounded depth and alias checks. Not v4 (still release candidate `v4.0.0-rc.6`; its README says release is intended for October 2026), not `goccy/go-yaml` (single-maintainer, 260 open issues, several are parser-correctness bugs, unbounded nesting memory), and not `gopkg.in/yaml.v3` (archived).

**A cheap mix that costs nothing at runtime:** use `go.yaml.in/yaml/v3` in the test suite only, as a differential oracle. Every fixture that the hand parser accepts must also be valid YAML and give the same keys and strings in the library. That keeps the hand parser from drifting into a private dialect, and that is the real risk of this choice (see open risks).

This is a recommendation for Alexander to confirm. It changes nothing in [[261027-ensf-akm-an-agent-knowledge-management]] or the spec; the spec already says (H5, X8) that a tool with no YAML library can parse the subset.

## Candidates

| # | Candidate | Role |
|---|---|---|
| 1 | Hand-written subset parser (this note, appendix) | baseline and recommendation |
| 2 | `go.yaml.in/yaml/v3` v3.0.5 | maintained successor of go-yaml v3, same API |
| 3 | `go.yaml.in/yaml/v4` v4.0.0-rc.6 | next major in the same repository, not released yet |
| 4 | `github.com/goccy/go-yaml` v1.19.2 | independent pure-Go implementation with its own AST |
| 5 | `gopkg.in/yaml.v3` v3.0.1 | the original, archived; kept as a reference point only |

Left out on purpose: `sigs.k8s.io/yaml` (a JSON-oriented wrapper over a go-yaml fork; no node tree, so no ordered fields or comments) and `go.yaml.in/yaml/v2` (older, frozen, no advantage over v3).

## Comparison on the four criteria

Facts as of 2026-10-28, from the GitHub API (`gh api repos/...`), the Go module proxy (`proxy.golang.org`), pkg.go.dev, and osv.dev.

| | Hand-written | `go.yaml.in/yaml/v3` | `go.yaml.in/yaml/v4` | `goccy/go-yaml` | `gopkg.in/yaml.v3` |
|---|---|---|---|---|---|
| Version / last release | n/a | v3.0.5, 2026-07-26 | v4.0.0-rc.6, 2026-06-17 (no GA yet) | v1.19.2, 2026-01-08 | v3.0.1, 2022-05-27 |
| Licence | n/a | Apache-2.0 (repo licence; not checked for per-file notices) | Apache-2.0 | MIT | Apache-2.0 (as the original; not re-checked) |
| Maintainers | n/a | `yaml` GitHub org; top contributors ingydotnet, niemeyer, ccoVeille | same repository | one dominant author (goccy: 605 commits; next contributor 22) | archived 2025-04-01; original author states it is unmaintained |
| Activity | n/a | repo pushed 2026-10-27 (v3 branch: security fixes only) | active development on `main` | last push 2026-04-11 | none |
| Open issues | n/a | 80 (whole repo) | 80 (whole repo) | 260 | 422 (frozen) |
| Runtime dependencies | 0 | 0 | 0 | 0 | 0 |
| Stars / importers | n/a | 544 stars (repo) / 555 importers | same repo / 376 importers | 2,243 / 2,983 | 7,014 / 44,731 |
| Known vulnerabilities (osv.dev, GitHub advisories) | none | none for the module | none | none listed | CVE-2022-28948 (panic on malformed input), fixed in v3.0.1 |
| **(1) Baseline** | all met; edits byte-identical (4 of 4) | parses the subset; edits **not** byte-identical (0 of 4); accepts out-of-subset input; retypes dates | same as v3 | same, and a value edit loses the trailing comment; key rename failed | same as v3 |
| **(2) Performance** (us/file, allocs/file) | **2.8, 28** | 17.9, 216 | 17.4, 237 | 22.7, 599 | 17.8, 216 |
| **(3) Security** | linear, no recursion, rejects aliases | max depth 10,000; alias-ratio check; slow decode on 1M keys | same as v3 | **unbounded** memory on deep flow nesting; 3 GB for 100 MB scalar | same as v3 (no fix for slow decode) |
| **(4) Popularity** | n/a | modest, growing (migration path of the original) | low (pre-release) | good | highest, but archived |

Sources:
- Repository state: `gh api repos/{yaml/go-yaml, go-yaml/yaml, goccy/go-yaml}` (archived flag, pushed_at, stars, issues, licence), contributors endpoint, `yaml/go-yaml` README (version table, v4 release plan).
- Versions and dates: `https://proxy.golang.org/<module>/@latest` and `@v/list`. Note `yaml/go-yaml` publishes tags, not GitHub releases.
- Importers: pkg.go.dev "Imported by" count for each module path.
- Vulnerabilities: `https://api.osv.dev/v1/query` for each module path, and the GitHub advisory database. Only `gopkg.in/yaml.v3` returned entries.
- Dependencies: `go list -m all` and `go list -deps` in the test module (no third-party modules besides the candidates).

go-yaml status in one line: the original `go-yaml/yaml` is archived (its README: "THIS PROJECT IS UNMAINTAINED"). The YAML organisation maintains the fork at `yaml/go-yaml`, published as `go.yaml.in/yaml/v2`, `/v3` and `/v4`. v2 and v3 are frozen and receive security fixes only; new work goes into v4.

## Behaviour results

Method: one program per candidate (appendix) reads each fixture's frontmatter block and normalizes it to fields (key, kind, tag, quoting style, values, line). Libraries are used in their node/AST mode so that ordering, quoting style and comments are visible. Line numbers below are file lines (block line plus 1 for the opening fence). Raw output is in the appendix.

### Subset coverage

- **Passed by all five:** plain, single-quoted and double-quoted strings with escapes; a 160-character `summary`; integers (including negative); inline `[a, b]` and block `- item` lists (also unindented items); empty values; unknown keys preserved in order; trailing and full-line comments; CRLF; empty frontmatter.
- **Dates:** the hand parser and `goccy/go-yaml` keep `2026-10-28` as a string. The go-yaml family tags it `!!timestamp` (the node value is still the original text, so a tool using nodes can read it as a string; decoding into `any` is documented to give `time.Time`, which I did not test).
- **Typing traps** (fixture `typing`): `ver: 1.10`, `zip: 08540` and a 20-digit integer become floats in the go-yaml family; `0x1F` and `1_000` become ints; `null` and `~` become null. `goccy/go-yaml` is closer to string-preserving but still types `1.10` as float and `0x1F`, `1_000` as int. `no` and `yes` stay strings in all of them. The hand parser types only `^-?(0|[1-9][0-9]*)$` as int and `~`/empty as null, everything else is a string. For `akmctl` the schema layer decides types from the field name, so reading every scalar as its text is the safer default.
- **Outside the subset** (anchor, alias, tag, nesting, flow map, block scalar, multi-line plain scalar): all four libraries parse these without complaint. Only the hand parser rejects them, with a line number. With a library, `akmctl` would have to walk the node tree and reject these itself (H5, rule W1).
- **Duplicate keys:** the go-yaml family accepts them in node mode (the decode into a map errors). `goccy/go-yaml` rejects them at parse time. The hand parser rejects them.

### Errors and line numbers

Fixtures with a known bad line. "OK" means the reported line is the offending line (file line, after adding 1).

| Fault | Hand | `go.yaml.in/yaml/v3` / `gopkg.in/yaml.v3` | `go.yaml.in/yaml/v4` | `goccy/go-yaml` |
|---|---|---|---|---|
| Unterminated `"` (file line 3) | OK, line 3 | OK (`yaml: line 2` + 1) | `L2.C10` + 1: OK | `[2:10]` + 1: OK |
| Unterminated `[a, b` (file line 4) | OK, line 4 | **wrong**: `line 2` (+1 = 3) | `L3.C7` + 1: OK | `[4:1]` + 1 = 5: points at the next line |
| `summary: a: b` (file line 3) | OK | OK | OK | OK |
| Duplicate key (file line 4) | OK | **no error** in node mode | **no error** in node mode | OK (`[3:1]` + 1) |
| Tab indent (file line 4) | OK | OK | OK | OK |
| Stray indented line (file line 3) | OK | OK | OK | `[1:8]` + 1 = 2: points at the previous line |

All four libraries report the line relative to the block they were given, so the tool must add the offset itself, same as it must for a hand parser fed only the block. `goccy/go-yaml` and the v3 family sometimes name a neighbouring line. v4 changed the error format (`L2.C10`, "go-yaml load error in scanner ...") and has an open issue titled "BC break: odd line/column notation in errors".

### Edit round trip

Edit fixture (in the appendix) has odd spacing after colons, trailing comments, a blank line, a loose inline list, a four-space block list and mixed quoting. Four edits: no-op, set `status`, rename `updated` to `updated_at`, add a default. The expected output is the naive text edit.

| Edit | Hand | `go.yaml.in/yaml/v3` | `go.yaml.in/yaml/v4` | `goccy/go-yaml` | `gopkg.in/yaml.v3` |
|---|---|---|---|---|---|
| no-op | identical | differs | differs | differs | differs |
| set `status` | identical | differs | differs | differs, **trailing comment lost** | differs |
| rename key | identical | differs | differs | differs, **key not renamed** | differs |
| add default | identical | differs | differs | not attempted via AST | differs |

What "differs" means for the go-yaml family: `title:   "x"   # c` becomes `title: "x" # c`; the blank line is removed; `[ a,b ,'c' ]` becomes `[a, b, 'c']`; the four-space list indent becomes two (with `SetIndent(2)`; the default is four). `goccy/go-yaml` keeps the blank line and the four-space indent but also normalizes spacing and the flow list. Comments survive in the go-yaml family, which is the best result any library gives, but "comments survive" is not "bytes survive". The task's baseline says an option that cannot do this must say so: **none of the four libraries can.** The hand parser does it by construction, because it never regenerates text.

Limits of the hand parser's edit: `set` handles scalar values only (list edits would be another splice kind); `add` appends `key: value\n` and does not copy a CRLF line ending; it does not re-indent or reflow.

## Performance

5,000 generated files, each 1,375 bytes with a 407-byte, 15-line frontmatter block (one 120-character `summary`, an inline list, a two-item block list, quoted values, a comment) and a 960-byte body that is not parsed. Machine: Apple M1 Max, 10 CPUs, 64 GB; `go1.27.2 darwin/arm64`; Go `testing.Benchmark`, three runs each. Times are for the whole batch of 5,000 (median of three), plus per file.

| Candidate | Mode | ms / 5,000 files | us / file | allocs / file | bytes / file |
|---|---|---|---|---|---|
| Hand-written | fields | 14.2 | 2.8 | 28 | 7,000 |
| `go.yaml.in/yaml/v3` | node, fields | 89.4 | 17.9 | 216 | 17,100 |
| `go.yaml.in/yaml/v3` | decode `map[string]any` | 107.7 | 21.5 | 282 | 18,100 |
| `go.yaml.in/yaml/v4` | node, fields | 86.8 | 17.4 | 237 | 17,500 |
| `go.yaml.in/yaml/v4` | decode `map[string]any` | 109.9 | 22.0 | 303 | 18,400 |
| `goccy/go-yaml` | node, fields | 113.4 | 22.7 | 599 | 32,400 |
| `goccy/go-yaml` | decode `map[string]any` | 153.2 | 30.6 | 740 | 35,100 |
| `gopkg.in/yaml.v3` | node, fields | 89.8 | 17.8 | 216 | 17,100 |

Reading the files from disk and parsing each (best of 5): hand 120 ms, go-yaml family 207 to 211 ms, goccy 245 ms. Most of the hand parser's 120 ms is file system time, which every candidate pays; parse cost is the difference.

Every candidate handles thousands of small files in well under a second. The hand parser is 6 to 8 times faster and allocates 8 to 21 times less, which matters for a check that runs in a pre-commit hook or on every `akmctl` call, but not for correctness. The hand parser has not been tuned (it copies the block to a string; 28 allocations per file could drop).

## Security

Each case ran in a child process of the benchmark binary, with the parent polling resident memory every 50 ms and killing at 6 GB or 30 s. "node" means parse to the tree without expanding aliases; "decode" means `Unmarshal` into `map[string]any`. Machine as above.

| Hostile input | Hand | `go.yaml.in/yaml/v3` (and `gopkg.in/yaml.v3`, same) | `go.yaml.in/yaml/v4` | `goccy/go-yaml` |
|---|---|---|---|---|
| Alias bomb, 9 levels of 9 (387 million leaves if expanded) | rejected at the first `&` | node: ok (not expanded); decode: error "excessive aliasing", 10 ms | same | node and decode: ok in 10 ms, 6 MB (aliases not expanded) |
| 100,000 nested `[` | rejected, 6 MB | error "exceeded max depth of 10000", 10 ms | same | **killed at 6.5 GB after 1 s** |
| Block nesting, 5,000 levels | rejected | ok, 45 to 68 MB, about 100 ms | same | node ok (144 MB); decode error "exceeded max depth" |
| 100,000 `- ` on one line | rejected | error, 6 MB | error | error after 1.6 s, 189 MB |
| One 100 MB scalar (plain or quoted) | ok, 207 to 407 MB, 140 ms | ok, 460 to 560 MB, 1.1 s | same | ok, **3.2 to 3.6 GB**, about 2 s |
| 1,000,000 keys (about 20 MB) | ok, 543 MB, 0.55 s | node ok in 1.0 s; **decode over 30 s (killed)** | same as v3 | **node and decode both over 30 s (killed)**, 3.3 GB |
| 1,000,000 list items | ok, 71 MB | ok, 250 to 300 MB | same | ok, 850 to 920 MB |

Reading it:
- The hand parser holds no state beyond the current line, so time and memory grow linearly with input, and nothing it accepts can expand. It does have no size limit of its own; the tool should cap file size (for example refuse frontmatter above 64 KB, the `summary` limit of 160 characters being far lower) before parsing. That cap is a one-line check and is also needed with any library.
- The go-yaml family limits depth (10,000) and alias expansion, which is real protection. Decoding a map of a million keys is slow (not completed in 30 s); I did not isolate why and did not check whether the duplicate-key check is the cause. Irrelevant for the subset's file sizes, but it shows that "bounded" depends on the mode.
- `goccy/go-yaml` has no depth limit in flow style and a high constant memory factor (about 30 times input size for one scalar). A hostile file in a cloned repo could exhaust memory of any tool that parses it. I did not report this upstream.

History: `gopkg.in/yaml.v3` CVE-2022-28948 (panic on malformed input) was fixed in v3.0.1; the test module uses that version. No advisories were found for the other modules on 2026-10-28.

## Open risks

- **Dialect drift.** The hand parser defines a dialect. If it accepts something YAML does not, or rejects valid YAML in the subset (an escape such as `\/`, a plain scalar with `: ` inside quotes, a key with spaces), files will behave differently in other tools (GitHub rendering, Obsidian properties). Mitigation: the differential test against `go.yaml.in/yaml/v3` described above, plus a fuzz target; neither is written yet.
- **Subset not fully pinned.** Decisions the spec should state before the parser is built: allowed key characters (mine: letters, digits, `_`, `-`), whether `[a,b]` needs a space, escapes allowed in double quotes (I used Go's `strconv.Unquote`, which is close to but not the same as YAML), BOM handling (not handled), whether `---` inside a value is an issue, line ending preservation when adding a field (not handled).
- **Obsidian interplay.** Obsidian's Properties UI rewrites frontmatter in its own style (it quotes, reorders and reindents). Byte-for-byte edits by `akmctl` do not stop Obsidian from changing a block when someone edits properties there. Not tested.
- **The evidence is small.** One machine, one generated corpus, 15-line blocks. Hostile cases are my choices, not a catalogue. The libraries were used at one version each.
- **If H5 widens** (nested maps, multi-line values), the hand parser stops being cheap and the recommendation flips to `go.yaml.in/yaml/v3` for reading, with the edit problem unsolved by any library; edits would still be splices guided by node line and column positions.

## What I could not verify

- Whether `go.yaml.in/yaml/v4` will reach GA with the same API; its README says October 2026 and the latest tag today is rc.6.
- The cause of the slow `map[string]any` decode on 1M keys in the go-yaml family.
- Maintainer counts beyond the contributors endpoint (commit counts, not people with merge rights), and open-issue lists beyond keyword scans of titles.
- Whether any private or unpublished advisories exist; only osv.dev and the public GitHub advisory database were checked.
- The hand parser against a fuzzer, or a full YAML 1.2 test suite; only the fixtures here.
- Behaviour on Windows or Linux; only macOS arm64 was run.

## Appendix A: how it was run

Module in the session scratch directory; `GOFLAGS=-mod=mod` to fetch once, then the `go.sum` below. Only my own programs were run; no candidate test suite, script or `go generate` was run. Commands: `go build -o yamleval .`, then `./yamleval behave`, `./yamleval bench`, `./yamleval hostile` (child mode is internal).

`go.mod` after `go mod tidy`:

```
module yamleval

go 1.27

require (
	github.com/goccy/go-yaml v1.19.2
	go.yaml.in/yaml/v3 v3.0.5
	go.yaml.in/yaml/v4 v4.0.0-rc.6
	gopkg.in/yaml.v3 v3.0.1
)
```

`go.sum`:

```
github.com/goccy/go-yaml v1.19.2 h1:PmFC1S6h8ljIz6gMRBopkjP1TVT7xuwrButHID66PoM=
github.com/goccy/go-yaml v1.19.2/go.mod h1:XBurs7gK8ATbW4ZPGKgcbrY1Br56PdM69F7LkFRi1kA=
go.yaml.in/yaml/v3 v3.0.5 h1:N6y/pJk8buWs9NY5ERU2HSMfm+IuD/OtfdAnq6kESPw=
go.yaml.in/yaml/v3 v3.0.5/go.mod h1:HVTZu1O7/Vkt2N+BFy8Zza+lnLsABggaTM2ZpNIGuKg=
go.yaml.in/yaml/v4 v4.0.0-rc.6 h1:1h7H1ohdUh93/FyE4YaDa1Zh64K6VVbjF4K6WUxMtH4=
go.yaml.in/yaml/v4 v4.0.0-rc.6/go.mod h1:aZqd9kCMsGL7AuUv/m/PvWLdg5sjJsZ4oHDEnfPPfY0=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405 h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
```

## Appendix B: source

The three go-yaml adapters (`lib_v3old.go`, `lib_v3new.go`, `lib_v4.go`) are identical except for the import path and type name; `lib_v3new.go` is shown, and the other two differ in the import and the type name (`yaml3old`, `yaml3new`, `yaml4`).

### `types.go`

```go
package main

import "bytes"

// Field is the normalized result every candidate must produce for one top-level key.
type Field struct {
	Key   string
	Kind  string // scalar | list | null | map | alias | other
	Tag   string // str | int | null | float | bool | timestamp ...
	Style string // plain | single | double | "" (lists: flow | block)
	Vals  []string
	Line  int // 1-based line within the block
}

type Op struct{ Kind, Key, Arg string } // set | rename | add

type Cand interface {
	Name() string
	Parse(block []byte) ([]Field, error)
	Decode(block []byte) error // typical-API decode into map[string]any (hostile-input mode)
	Edit(block []byte, ops []Op) ([]byte, error)
}

// splitFM returns the frontmatter block (between the fences), and the file line of the block's first line.
func splitFM(file []byte) (block []byte, firstLine int, ok bool) {
	if !bytes.HasPrefix(file, []byte("---\n")) && !bytes.HasPrefix(file, []byte("---\r\n")) {
		return nil, 0, false
	}
	i := bytes.IndexByte(file, '\n') + 1
	rest := file[i:]
	off := 0
	for off <= len(rest) {
		j := bytes.IndexByte(rest[off:], '\n')
		line := rest[off:]
		if j >= 0 {
			line = rest[off : off+j]
		}
		if bytes.Equal(bytes.TrimRight(line, "\r"), []byte("---")) {
			return rest[:off], 2, true
		}
		if j < 0 {
			break
		}
		off += j + 1
	}
	return nil, 0, false
}
```

### `hand.go`

```go
package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Hand is a hand-written parser for the H5 subset. It records byte offsets so edits are splices.
type Hand struct{}

func (Hand) Name() string          { return "hand-written" }
func (Hand) Decode(b []byte) error { _, err := Hand{}.Parse(b); return err }

type hf struct {
	Field
	keyS, keyE, valS, valE int // byte offsets in block; valS=-1 when no inline value
	lineEnd                int // end of field incl. block-list items
}

func isKeyByte(c byte) bool {
	return c == '_' || c == '-' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// scalar parses one scalar starting at s[i]; returns value, tag, style, end offset of the scalar text.
func scalar(s string, i int, inFlow bool) (val, tag, style string, end int, err string) {
	if i >= len(s) {
		return "", "null", "plain", i, ""
	}
	switch c := s[i]; c {
	case '\'':
		var sb strings.Builder
		j := i + 1
		for j < len(s) {
			if s[j] == '\'' {
				if j+1 < len(s) && s[j+1] == '\'' {
					sb.WriteByte('\'')
					j += 2
					continue
				}
				return sb.String(), "str", "single", j + 1, ""
			}
			sb.WriteByte(s[j])
			j++
		}
		return "", "", "", j, "unterminated single-quoted string"
	case '"':
		j := i + 1
		for j < len(s) {
			if s[j] == '\\' {
				j += 2
				continue
			}
			if s[j] == '"' {
				u, e := strconv.Unquote(s[i : j+1])
				if e != nil {
					return "", "", "", j, "unsupported escape in double-quoted string"
				}
				return u, "str", "double", j + 1, ""
			}
			j++
		}
		return "", "", "", j, "unterminated double-quoted string"
	case '&', '*', '!', '|', '>', '{', '%', '@', '`':
		return "", "", "", i, fmt.Sprintf("%q starts an anchor, alias, tag, block scalar or map: outside the subset", c)
	case '[':
		return "", "", "", i, "nested list: outside the subset"
	}
	if c0 := s[i]; (c0 == '-' || c0 == '?' || c0 == ':') && (i+1 == len(s) || s[i+1] == ' ') {
		return "", "", "", i, "indicator character at the start of a plain value: outside the subset"
	}
	j := i
	for j < len(s) {
		if s[j] == ' ' && j+1 < len(s) && s[j+1] == '#' {
			break
		}
		if inFlow && (s[j] == ',' || s[j] == ']') {
			break
		}
		if s[j] == ':' && (j+1 == len(s) || s[j+1] == ' ') {
			return "", "", "", j, "mapping values are not allowed here"
		}
		j++
	}
	v := strings.TrimRight(s[i:j], " \t")
	end = i + len(v)
	tag = "str"
	if isInt(v) {
		tag = "int"
	} else if v == "" || v == "~" {
		tag = "null"
	}
	return v, tag, "plain", end, ""
}

func isInt(v string) bool {
	t := strings.TrimPrefix(v, "-")
	if t == "" || (len(t) > 1 && t[0] == '0') {
		return false
	}
	for i := 0; i < len(t); i++ {
		if t[i] < '0' || t[i] > '9' {
			return false
		}
	}
	return true
}

func restOK(s string, i int) bool { // only spaces and an optional comment may follow
	rest := s[i:]
	t := strings.TrimLeft(rest, " \t")
	return t == "" || t[0] == '#' && len(rest) != len(t)
}

func (Hand) parse(block []byte) ([]hf, error) {
	text := string(block)
	var out []hf
	seen := map[string]bool{}
	off, ln := 0, 0
	for off < len(text) {
		ln++
		nl := strings.IndexByte(text[off:], '\n')
		lineEnd := len(text)
		next := len(text)
		if nl >= 0 {
			lineEnd, next = off+nl, off+nl+1
		}
		line := strings.TrimRight(text[off:lineEnd], "\r")
		t := strings.TrimLeft(line, " ")
		fail := func(m string) error { return fmt.Errorf("line %d: %s", ln, m) }
		switch {
		case t == "" || t[0] == '#':
		case strings.ContainsRune(line[:len(line)-len(t)], '\t') || t[0] == '\t':
			return nil, fail("tab in indentation")
		case t == "-" || strings.HasPrefix(t, "- "):
			if len(out) == 0 || out[len(out)-1].valS >= 0 || !(out[len(out)-1].Kind == "null" || out[len(out)-1].Style == "block") {
				return nil, fail("list item without a list key")
			}
			f := &out[len(out)-1]
			f.Kind, f.Style, f.Tag = "list", "block", ""
			i := len(line) - len(t) + 1
			for i < len(line) && line[i] == ' ' {
				i++
			}
			v, tag, _, end, e := scalar(line, i, false)
			if e != "" {
				return nil, fail(e)
			}
			if !restOK(line, end) {
				return nil, fail("unexpected text after list item")
			}
			if tag == "null" && t == "-" {
				return nil, fail("empty list item")
			}
			f.Vals = append(f.Vals, v)
			f.lineEnd = next
		case len(line) != len(t):
			return nil, fail("indented line: nesting and multi-line values are outside the subset")
		default:
			k := 0
			for k < len(line) && isKeyByte(line[k]) {
				k++
			}
			if k == 0 || k >= len(line) || line[k] != ':' || (k+1 < len(line) && line[k+1] != ' ') {
				return nil, fail("expected 'key: value'")
			}
			key := line[:k]
			if seen[key] {
				return nil, fail("duplicate key " + key)
			}
			seen[key] = true
			f := hf{Field: Field{Key: key, Line: ln}, keyS: off, keyE: off + k, valS: -1, lineEnd: next}
			i := k + 1
			for i < len(line) && line[i] == ' ' {
				i++
			}
			if i >= len(line) || line[i] == '#' {
				f.Kind, f.Tag, f.Style = "null", "null", "plain"
			} else if line[i] == '[' {
				i++
				f.Kind, f.Style = "list", "flow"
				f.valS = off + i - 1
				for {
					for i < len(line) && line[i] == ' ' {
						i++
					}
					if i >= len(line) {
						return nil, fail("unterminated inline list")
					}
					if line[i] == ']' {
						i++
						break
					}
					v, _, _, end, e := scalar(line, i, true)
					if e != "" {
						return nil, fail(e)
					}
					f.Vals = append(f.Vals, v)
					i = end
					for i < len(line) && line[i] == ' ' {
						i++
					}
					if i < len(line) && line[i] == ',' {
						i++
					} else if i >= len(line) || line[i] != ']' {
						return nil, fail("unterminated inline list")
					}
				}
				f.valE = off + i
				if !restOK(line, i) {
					return nil, fail("unexpected text after list")
				}
			} else {
				v, tag, style, end, e := scalar(line, i, false)
				if e != "" {
					return nil, fail(e)
				}
				if !restOK(line, end) {
					return nil, fail("unexpected text after value")
				}
				f.Kind, f.Tag, f.Style, f.Vals = "scalar", tag, style, []string{v}
				f.valS, f.valE = off+i, off+end
			}
			out = append(out, f)
		}
		off = next
	}
	return out, nil
}

func (h Hand) Parse(block []byte) ([]Field, error) {
	fs, err := h.parse(block)
	if err != nil {
		return nil, err
	}
	res := make([]Field, len(fs))
	for i := range fs {
		res[i] = fs[i].Field
	}
	return res, nil
}

func (h Hand) Edit(block []byte, ops []Op) ([]byte, error) {
	cur := block
	for _, op := range ops {
		fs, err := h.parse(cur)
		if err != nil {
			return nil, err
		}
		var f *hf
		for i := range fs {
			if fs[i].Key == op.Key {
				f = &fs[i]
			}
		}
		s := string(cur)
		switch op.Kind {
		case "set":
			if f == nil || f.valS < 0 {
				return nil, fmt.Errorf("no scalar %s", op.Key)
			}
			s = s[:f.valS] + op.Arg + s[f.valE:]
		case "rename":
			if f == nil {
				return nil, fmt.Errorf("no key %s", op.Key)
			}
			s = s[:f.keyS] + op.Arg + s[f.keyE:]
		case "add":
			if f != nil {
				continue
			}
			if s != "" && !strings.HasSuffix(s, "\n") {
				s += "\n"
			}
			s += op.Key + ": " + op.Arg + "\n"
		}
		cur = []byte(s)
	}
	return cur, nil
}
```

### `lib_v3new.go`

```go
package main

import (
	"bytes"
	"fmt"
	"strings"

	yaml "go.yaml.in/yaml/v3"
)

type yaml3new struct{}

func (yaml3new) Name() string { return "go.yaml.in/yaml/v3" }

func (yaml3new) Parse(block []byte) ([]Field, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(block, &doc); err != nil {
		return nil, err
	}
	if doc.Kind == 0 {
		return nil, nil
	}
	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("line %d: top level is not a mapping", root.Line)
	}
	out := make([]Field, 0, len(root.Content)/2)
	for i := 0; i+1 < len(root.Content); i += 2 {
		k, v := root.Content[i], root.Content[i+1]
		f := Field{Key: k.Value, Line: k.Line}
		switch v.Kind {
		case yaml.ScalarNode:
			f.Kind, f.Tag, f.Vals = "scalar", strings.TrimPrefix(v.Tag, "!!"), []string{v.Value}
			switch v.Style {
			case yaml.SingleQuotedStyle:
				f.Style = "single"
			case yaml.DoubleQuotedStyle:
				f.Style = "double"
			default:
				f.Style = "plain"
			}
			if f.Tag == "null" {
				f.Kind = "null"
			}
		case yaml.SequenceNode:
			f.Kind, f.Style = "list", "block"
			if v.Style&yaml.FlowStyle != 0 {
				f.Style = "flow"
			}
			for _, c := range v.Content {
				f.Vals = append(f.Vals, c.Value)
			}
		case yaml.MappingNode:
			f.Kind = "map"
		case yaml.AliasNode:
			f.Kind = "alias"
		}
		out = append(out, f)
	}
	return out, nil
}

func (yaml3new) Decode(block []byte) error {
	var m map[string]any
	return yaml.Unmarshal(block, &m)
}

func (yaml3new) Edit(block []byte, ops []Op) ([]byte, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(block, &doc); err != nil {
		return nil, err
	}
	root := doc.Content[0]
	for _, op := range ops {
		idx := -1
		for i := 0; i+1 < len(root.Content); i += 2 {
			if root.Content[i].Value == op.Key {
				idx = i
			}
		}
		switch op.Kind {
		case "set":
			root.Content[idx+1].Value = op.Arg
		case "rename":
			root.Content[idx].Value = op.Arg
		case "add":
			if idx < 0 {
				root.Content = append(root.Content,
					&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: op.Key},
					&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: op.Arg})
			}
		}
	}
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&doc); err != nil {
		return nil, err
	}
	enc.Close()
	return buf.Bytes(), nil
}
```

### `lib_goccy.go`

```go
package main

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
)

type goccy struct{}

func (goccy) Name() string { return "github.com/goccy/go-yaml" }

func (goccy) Parse(block []byte) ([]Field, error) {
	f, err := parser.ParseBytes(block, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if len(f.Docs) == 0 || f.Docs[0].Body == nil {
		return nil, nil
	}
	var mvs []*ast.MappingValueNode
	switch b := f.Docs[0].Body.(type) {
	case *ast.MappingNode:
		mvs = b.Values
	case *ast.MappingValueNode:
		mvs = []*ast.MappingValueNode{b}
	default:
		return nil, fmt.Errorf("top level is %T, not a mapping", b)
	}
	var out []Field
	for _, mv := range mvs {
		kt := mv.Key.GetToken()
		fl := Field{Key: kt.Value, Line: kt.Position.Line}
		switch v := mv.Value.(type) {
		case *ast.StringNode:
			fl.Kind, fl.Tag, fl.Vals, fl.Style = "scalar", "str", []string{v.Value}, "plain"
			switch v.Token.Type {
			case token.SingleQuoteType:
				fl.Style = "single"
			case token.DoubleQuoteType:
				fl.Style = "double"
			}
		case *ast.IntegerNode:
			fl.Kind, fl.Tag, fl.Vals, fl.Style = "scalar", "int", []string{v.GetToken().Value}, "plain"
		case *ast.FloatNode:
			fl.Kind, fl.Tag, fl.Vals, fl.Style = "scalar", "float", []string{v.GetToken().Value}, "plain"
		case *ast.BoolNode:
			fl.Kind, fl.Tag, fl.Vals, fl.Style = "scalar", "bool", []string{v.GetToken().Value}, "plain"
		case *ast.NullNode:
			fl.Kind, fl.Tag, fl.Style = "null", "null", "plain"
		case *ast.SequenceNode:
			fl.Kind, fl.Style = "list", "block"
			if v.IsFlowStyle {
				fl.Style = "flow"
			}
			for _, e := range v.Values {
				fl.Vals = append(fl.Vals, strings.TrimSpace(e.GetToken().Value))
			}
		case *ast.MappingNode, *ast.MappingValueNode:
			fl.Kind = "map"
		case *ast.AliasNode:
			fl.Kind = "alias"
		default:
			fl.Kind = fmt.Sprintf("other:%T", v)
		}
		out = append(out, fl)
	}
	return out, nil
}

func (goccy) Decode(block []byte) error {
	var m map[string]any
	return yaml.Unmarshal(block, &m)
}

func (g goccy) Edit(block []byte, ops []Op) ([]byte, error) {
	f, err := parser.ParseBytes(block, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	for _, op := range ops {
		switch op.Kind {
		case "set":
			p, err := yaml.PathString("$." + op.Key)
			if err != nil {
				return nil, err
			}
			if err := p.ReplaceWithReader(f, strings.NewReader(op.Arg)); err != nil {
				return nil, err
			}
		case "rename":
			mn, ok := f.Docs[0].Body.(*ast.MappingNode)
			if !ok {
				return nil, fmt.Errorf("body is %T", f.Docs[0].Body)
			}
			for _, mv := range mn.Values {
				if mv.Key.GetToken().Value == op.Key {
					t := mv.Key.GetToken()
					t.Value = op.Arg
					t.Origin = strings.Replace(t.Origin, op.Key, op.Arg, 1)
				}
			}
		case "add":
			return nil, fmt.Errorf("add not attempted via AST (would be a text splice)")
		}
	}
	return []byte(f.String()), nil
}
```

### `fixtures.go`

```go
package main

import "strings"

var long160 = strings.Repeat("0123456789", 15) + "abcdefghij" // 160 chars

type Fixture struct{ Name, Doc, Note string }

var fixtures = []Fixture{
	{"basic", `---
title: Plain string with spaces
quoted_s: 'it''s single'
quoted_d: "say \"hi\"\ttab"
count: 42
neg: -7
date: 2026-10-28
status: active
tags: [a, b, "c d", 'e']
blocked_by:
  - one
  - 'two'
  - 3
flush:
- x
- y
nothing:
unknown_key: kept # trailing comment
# full-line comment
summary: ` + long160 + `
---
body
`, "subset in one file"},
	{"typing", `---
country: no
flag: yes
ver: 1.10
zip: 08540
big: 99999999999999999999
stamp: 2026-10-28T10:00:00Z
tilde: ~
nulltxt: null
hex: 0x1F
under: 1_000
id: 261027-k7f2
colon_in_value: "a: b"
hash_in_value: a#b
url: http://example.com/x
---
`, "values that YAML libraries retype; the tool wants strings kept as written"},
	{"anchor", `---
base: &b value
copy: *b
tagged: !!str 12
---
`, "anchor, alias and tag: outside the subset; tool must reject or detect"},
	{"nested", `---
a: 1
b:
  c: 2
  d: [1, 2]
e: {x: 1}
---
`, "nesting and flow map: outside the subset"},
	{"multiline", `---
a: |
  literal
  block
b: >
  folded
c: plain
  continued
---
`, "block scalars and multi-line plain: outside the subset"},
	{"empty", "---\n---\nbody\n", "empty frontmatter"},
	{"crlf", "---\r\ntitle: x\r\nstatus: active\r\ntags: [a]\r\n---\r\n", "CRLF line endings"},
	{"bad-unterminated-quote", "---\ntitle: ok\nsummary: \"never closed\nstatus: active\n---\n", "error expected near file line 3"},
	{"bad-unterminated-list", "---\ntitle: ok\nstatus: active\ntags: [a, b\nowner: me\n---\n", "error expected near file line 4"},
	{"bad-colon", "---\ntitle: ok\nsummary: a: b\nstatus: x\n---\n", "error expected at file line 3"},
	{"bad-dup", "---\ntitle: ok\nstatus: a\nstatus: b\n---\n", "error expected at file line 4"},
	{"bad-tab", "---\ntitle: ok\nblocked_by:\n\t- x\nstatus: a\n---\n", "error expected at file line 4"},
	{"bad-indent", "---\ntitle: ok\n  status: a\nowner: me\n---\n", "error expected at file line 3"},
}

// Edit fixture: odd spacing, comments, mixed quoting, blank lines, key order.
const editDoc = `# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: active   # state
updated: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
`

var editCases = []struct {
	Name string
	Ops  []Op
	Want string
}{
	{"no-op", nil, editDoc},
	{"set-status", []Op{{"set", "status", "done"}}, strings.Replace(editDoc, "status: active   # state", "status: done   # state", 1)},
	{"rename-updated", []Op{{"rename", "updated", "updated_at"}}, strings.Replace(editDoc, "updated:", "updated_at:", 1)},
	{"add-default", []Op{{"add", "priority", "3"}}, editDoc + "priority: 3\n"},
}
```

### `main.go`

```go
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var cands = []Cand{Hand{}, yaml3old{}, yaml3new{}, yaml4{}, goccy{}}

func byName(n string) Cand {
	for _, c := range cands {
		if c.Name() == n {
			return c
		}
	}
	return nil
}

func show(fs []Field) string {
	var sb strings.Builder
	for _, f := range fs {
		fmt.Fprintf(&sb, "    L%-2d %-14s %-6s %-5s %-6s %q\n", f.Line, f.Key, f.Kind, f.Tag, f.Style, f.Vals)
	}
	return sb.String()
}

func safe(f func()) (msg string) {
	defer func() {
		if r := recover(); r != nil {
			msg = fmt.Sprintf("PANIC: %v", r)
		}
	}()
	f()
	return ""
}

var lineRe = regexp.MustCompile(`(?:line |\[)(\d+)`)

func behave() {
	for _, fx := range fixtures {
		fmt.Printf("=== fixture %s (%s)\n", fx.Name, fx.Note)
		block, off, ok := splitFM([]byte(fx.Doc))
		if !ok {
			fmt.Println("  split failed")
			continue
		}
		for _, c := range cands {
			var fs []Field
			var err error
			p := safe(func() { fs, err = c.Parse(block) })
			fmt.Printf("  [%s]\n", c.Name())
			if p != "" {
				fmt.Println("    " + p)
			} else if err != nil {
				e := strings.ReplaceAll(err.Error(), "\n", " / ")
				if m := lineRe.FindStringSubmatch(e); m != nil {
					n := 0
					fmt.Sscan(m[1], &n)
					e += fmt.Sprintf("   => file line %d", n+off-1)
				}
				fmt.Println("    ERROR: " + e)
			} else {
				fmt.Print(show(fs))
			}
		}
	}
	fmt.Println("=== edit round-trip (block only; want = naive text edit)")
	for _, ec := range editCases {
		for _, c := range cands {
			var got []byte
			var err error
			p := safe(func() { got, err = c.Edit([]byte(editDoc), ec.Ops) })
			switch {
			case p != "":
				fmt.Printf("  %-12s %-28s %s\n", ec.Name, c.Name(), p)
			case err != nil:
				fmt.Printf("  %-12s %-28s ERROR %v\n", ec.Name, c.Name(), err)
			case string(got) == ec.Want:
				fmt.Printf("  %-12s %-28s BYTE-IDENTICAL\n", ec.Name, c.Name())
			default:
				fmt.Printf("  %-12s %-28s DIFFERS\n--- got\n%s--- want\n%s", ec.Name, c.Name(), got, ec.Want)
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: behave | bench | hostile")
		return
	}
	switch os.Args[1] {
	case "behave":
		behave()
	case "bench":
		bench()
	case "hostile":
		hostile()
	case "child":
		child(os.Args[2:])
	}
}
```

### `bench.go`

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func genFile(i int) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "---\nid: 2610%02d-%04x\ntitle: Work item number %d about something\nsummary: %s\nstatus: active\ncreated: 2026-10-%02d\nupdated: 2026-10-28\nparent: 261001-aaaa\nowner: alex\npriority: %d\ntags: [akm, yaml, \"spec %d\"]\nblocked_by:\n  - 261002-bbbb\n  - free text reason\nnotes: 'quoted ''value'' here'\n# a comment\n---\n", i%28+1, i, i, strings.Repeat("s", 120), i%28+1, i%5, i)
	sb.WriteString(strings.Repeat("Body text line.\n", 60))
	return []byte(sb.String())
}

func bench() {
	const N = 5000
	files := make([][]byte, N)
	dir, _ := os.MkdirTemp("", "yamleval")
	defer os.RemoveAll(dir)
	for i := range files {
		files[i] = genFile(i)
		os.WriteFile(filepath.Join(dir, fmt.Sprintf("%05d.md", i)), files[i], 0o644)
	}
	fmt.Printf("%d files, %d bytes each (frontmatter ~%d), go %s %s/%s, %d CPU\n", N, len(files[0]), len(strings.SplitN(string(files[0]), "---\n", 3)[1]), runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.NumCPU())
	for _, c := range cands {
		for _, mode := range []string{"node/fields", "decode map[string]any"} {
			if c.Name() == "hand-written" && mode != "node/fields" {
				continue
			}
			for run := 0; run < 3; run++ {
				r := testing.Benchmark(func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						for _, f := range files {
							blk, _, _ := splitFM(f)
							if mode == "node/fields" {
								if _, err := c.Parse(blk); err != nil {
									b.Fatal(err)
								}
							} else if err := c.Decode(blk); err != nil {
								b.Fatal(err)
							}
						}
					}
				})
				fmt.Printf("%-28s %-22s run%d  %8.2f ms/5000 files  %6.2f us/file  %9d B/op  %8d allocs/op (per 5000)\n", c.Name(), mode, run+1,
					float64(r.NsPerOp())/1e6, float64(r.NsPerOp())/1e3/N, r.AllocedBytesPerOp(), r.AllocsPerOp())
			}
		}
	}
	// disk: read + parse every file once, best of 5
	for _, c := range cands {
		best := time.Hour
		for run := 0; run < 5; run++ {
			t := time.Now()
			ents, _ := os.ReadDir(dir)
			for _, e := range ents {
				d, _ := os.ReadFile(filepath.Join(dir, e.Name()))
				blk, _, _ := splitFM(d)
				if _, err := c.Parse(blk); err != nil {
					panic(err)
				}
			}
			if el := time.Since(t); el < best {
				best = el
			}
		}
		fmt.Printf("disk read+parse 5000 files  %-28s best of 5: %v\n", c.Name(), best)
	}
}
```

### `hostile.go`

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func hcase(name string) []byte {
	var sb strings.Builder
	switch name {
	case "alias-bomb-9x9":
		sb.WriteString("a0: &a0 [x,x,x,x,x,x,x,x,x]\n")
		for i := 1; i <= 9; i++ {
			fmt.Fprintf(&sb, "a%d: &a%d [%s]\n", i, i, strings.TrimSuffix(strings.Repeat(fmt.Sprintf("*a%d,", i-1), 9), ","))
		}
	case "deep-flow-100k":
		sb.WriteString("k: " + strings.Repeat("[", 100000) + strings.Repeat("]", 100000) + "\n")
	case "deep-block-5k":
		for i := 0; i < 5000; i++ {
			sb.WriteString(strings.Repeat(" ", i) + "k:\n")
		}
	case "deep-seq-100k":
		sb.WriteString("k: " + strings.Repeat("- ", 100000) + "x\n")
	case "scalar-100MB-plain":
		sb.WriteString("k: " + strings.Repeat("a", 100<<20) + "\n")
	case "scalar-100MB-dquote":
		sb.WriteString("k: \"" + strings.Repeat("a", 100<<20) + "\"\n")
	case "1M-keys":
		for i := 0; i < 1000000; i++ {
			fmt.Fprintf(&sb, "key%07d: value %d\n", i, i)
		}
	case "list-1M-items":
		sb.WriteString("k:\n")
		for i := 0; i < 1000000; i++ {
			sb.WriteString("  - item\n")
		}
	}
	return []byte(sb.String())
}

var hnames = []string{"alias-bomb-9x9", "deep-flow-100k", "deep-block-5k", "deep-seq-100k", "scalar-100MB-plain", "scalar-100MB-dquote", "1M-keys", "list-1M-items"}

func child(a []string) { // child <cand> <case> <mode>
	c := byName(a[0])
	in := hcase(a[1])
	t := time.Now()
	var err error
	if a[2] == "node" {
		_, err = c.Parse(in)
	} else {
		err = c.Decode(in)
	}
	msg := "ok"
	if err != nil {
		msg = "error: " + strings.ReplaceAll(err.Error(), "\n", " / ")
		if len(msg) > 110 {
			msg = msg[:110] + "..."
		}
	}
	fmt.Printf("%s|%v\n", msg, time.Since(t).Round(time.Millisecond))
}

func rssKB(pid int) int {
	out, _ := exec.Command("ps", "-o", "rss=", "-p", strconv.Itoa(pid)).Output()
	n, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	return n
}

func hostile() {
	exe, _ := os.Executable()
	const limitKB = 6 << 20 // 6 GB
	fmt.Printf("%-28s %-26s %-6s %-14s %-9s %s\n", "candidate", "case", "mode", "peak RSS", "wall", "outcome")
	for _, c := range cands {
		if o := os.Getenv("ONLY"); o != "" && o != c.Name() {
			continue
		}
		for _, hn := range hnames {
			for _, mode := range []string{"node", "decode"} {
				if c.Name() == "hand-written" && mode == "decode" {
					continue
				}
				cmd := exec.Command(exe, "child", c.Name(), hn, mode)
				var out strings.Builder
				cmd.Stdout, cmd.Stderr = &out, &out
				start := time.Now()
				cmd.Start()
				done := make(chan error, 1)
				go func() { done <- cmd.Wait() }()
				outcome := ""
				peak := 0
			loop:
				for {
					select {
					case <-done:
						break loop
					case <-time.After(50 * time.Millisecond):
						if r := rssKB(cmd.Process.Pid); r > peak {
							peak = r
						}
						if r := peak; r > limitKB {
							cmd.Process.Kill()
							<-done
							outcome = "KILLED: RSS over 6 GB"
							break loop
						}
						if time.Since(start) > 30*time.Second {
							cmd.Process.Kill()
							<-done
							outcome = "KILLED: over 30 s"
							break loop
						}
					}
				}
				if ru, ok := cmd.ProcessState.SysUsage().(*syscall.Rusage); ok && int(ru.Maxrss/1024) > peak {
					peak = int(ru.Maxrss / 1024)
				}
				wall := time.Since(start).Round(10 * time.Millisecond)
				if outcome == "" {
					o := strings.TrimSpace(out.String())
					if cmd.ProcessState.ExitCode() != 0 {
						first := strings.SplitN(o, "\n", 2)[0]
						outcome = fmt.Sprintf("CRASH exit %d: %.90s", cmd.ProcessState.ExitCode(), first)
					} else {
						outcome = strings.Split(o, "|")[0]
					}
				}
				fmt.Printf("%-28s %-26s %-6s %-14s %-9v %s\n", c.Name(), hn, mode, fmt.Sprintf("%d MB", peak/1024), wall, outcome)
			}
		}
	}
}
```

## Appendix C: raw results

### Benchmark (`bench_final.txt`)

```
5000 files, 1375 bytes each (frontmatter ~407), go go1.27.2 darwin/arm64, 10 CPU
gopkg.in/yaml.v3             node/fields            run1     89.76 ms/5000 files   17.95 us/file   85738588 B/op   1080037 allocs/op (per 5000)
gopkg.in/yaml.v3             node/fields            run2     91.41 ms/5000 files   18.28 us/file   85738775 B/op   1080037 allocs/op (per 5000)
gopkg.in/yaml.v3             node/fields            run3     89.00 ms/5000 files   17.80 us/file   85739246 B/op   1080039 allocs/op (per 5000)
gopkg.in/yaml.v3             decode map[string]any  run1    109.80 ms/5000 files   21.96 us/file   90459815 B/op   1410040 allocs/op (per 5000)
gopkg.in/yaml.v3             decode map[string]any  run2    107.68 ms/5000 files   21.54 us/file   90459772 B/op   1410039 allocs/op (per 5000)
gopkg.in/yaml.v3             decode map[string]any  run3    109.76 ms/5000 files   21.95 us/file   90459879 B/op   1410040 allocs/op (per 5000)
go.yaml.in/yaml/v3           node/fields            run1     89.38 ms/5000 files   17.88 us/file   85738709 B/op   1080037 allocs/op (per 5000)
go.yaml.in/yaml/v3           node/fields            run2     89.76 ms/5000 files   17.95 us/file   85738534 B/op   1080037 allocs/op (per 5000)
go.yaml.in/yaml/v3           node/fields            run3     88.66 ms/5000 files   17.73 us/file   85738540 B/op   1080037 allocs/op (per 5000)
go.yaml.in/yaml/v3           decode map[string]any  run1    107.39 ms/5000 files   21.48 us/file   90459662 B/op   1410039 allocs/op (per 5000)
go.yaml.in/yaml/v3           decode map[string]any  run2    107.65 ms/5000 files   21.53 us/file   90459829 B/op   1410040 allocs/op (per 5000)
go.yaml.in/yaml/v3           decode map[string]any  run3    107.87 ms/5000 files   21.57 us/file   90460096 B/op   1410041 allocs/op (per 5000)
go.yaml.in/yaml/v4           node/fields            run1     89.83 ms/5000 files   17.97 us/file   87698926 B/op   1185038 allocs/op (per 5000)
go.yaml.in/yaml/v4           node/fields            run2     86.77 ms/5000 files   17.35 us/file   87698977 B/op   1185038 allocs/op (per 5000)
go.yaml.in/yaml/v4           node/fields            run3     86.71 ms/5000 files   17.34 us/file   87698892 B/op   1185037 allocs/op (per 5000)
go.yaml.in/yaml/v4           decode map[string]any  run1    109.87 ms/5000 files   21.97 us/file   92420001 B/op   1515040 allocs/op (per 5000)
go.yaml.in/yaml/v4           decode map[string]any  run2    109.95 ms/5000 files   21.99 us/file   92420244 B/op   1515041 allocs/op (per 5000)
go.yaml.in/yaml/v4           decode map[string]any  run3    112.42 ms/5000 files   22.48 us/file   92420220 B/op   1515041 allocs/op (per 5000)
github.com/goccy/go-yaml     node/fields            run1    115.02 ms/5000 files   23.00 us/file  162174648 B/op   2994475 allocs/op (per 5000)
github.com/goccy/go-yaml     node/fields            run2    113.09 ms/5000 files   22.62 us/file  162174065 B/op   2994467 allocs/op (per 5000)
github.com/goccy/go-yaml     node/fields            run3    113.42 ms/5000 files   22.68 us/file  162173322 B/op   2994463 allocs/op (per 5000)
github.com/goccy/go-yaml     decode map[string]any  run1    153.18 ms/5000 files   30.64 us/file  175744669 B/op   3699509 allocs/op (per 5000)
github.com/goccy/go-yaml     decode map[string]any  run2    151.75 ms/5000 files   30.35 us/file  175734112 B/op   3699439 allocs/op (per 5000)
github.com/goccy/go-yaml     decode map[string]any  run3    153.58 ms/5000 files   30.72 us/file  175745509 B/op   3699517 allocs/op (per 5000)
hand-written                 node/fields            run1     14.17 ms/5000 files    2.83 us/file   34960220 B/op    140000 allocs/op (per 5000)
hand-written                 node/fields            run2     13.97 ms/5000 files    2.79 us/file   34960154 B/op    140000 allocs/op (per 5000)
hand-written                 node/fields            run3     14.24 ms/5000 files    2.85 us/file   34960090 B/op    140000 allocs/op (per 5000)
disk read+parse 5000 files  gopkg.in/yaml.v3             best of 5: 210.648667ms
disk read+parse 5000 files  go.yaml.in/yaml/v3           best of 5: 206.7045ms
disk read+parse 5000 files  go.yaml.in/yaml/v4           best of 5: 206.878792ms
disk read+parse 5000 files  github.com/goccy/go-yaml     best of 5: 244.993875ms
disk read+parse 5000 files  hand-written                 best of 5: 120.266875ms
```

The hand-written rows are from a rerun after the last parser fix; the others are from the first run (the library code was not changed between runs).

### Hostile input (`hostile_final.txt`)

```
candidate                    case                       mode   peak RSS       wall      outcome
hand-written                 alias-bomb-9x9             node   5 MB           10ms      error: line 1: '&' starts an anchor, alias, tag, block scalar or map: outside the subset
hand-written                 deep-flow-100k             node   6 MB           10ms      error: line 1: nested list: outside the subset
hand-written                 deep-block-5k              node   54 MB          20ms      error: line 2: indented line: nesting and multi-line values are outside the subset
hand-written                 deep-seq-100k              node   6 MB           10ms      error: line 1: indicator character at the start of a plain value: outside the subset
hand-written                 scalar-100MB-plain         node   207 MB         140ms     ok
hand-written                 scalar-100MB-dquote        node   407 MB         130ms     ok
hand-written                 1M-keys                    node   543 MB         550ms     ok
hand-written                 list-1M-items              node   71 MB          80ms      ok
gopkg.in/yaml.v3             alias-bomb-9x9             node   5 MB           10ms      ok
gopkg.in/yaml.v3             alias-bomb-9x9             decode 6 MB           10ms      error: yaml: document contains excessive aliasing
gopkg.in/yaml.v3             deep-flow-100k             node   14 MB          10ms      error: yaml: exceeded max depth of 10000
gopkg.in/yaml.v3             deep-flow-100k             decode 14 MB          10ms      error: yaml: exceeded max depth of 10000
gopkg.in/yaml.v3             deep-block-5k              node   45 MB          90ms      ok
gopkg.in/yaml.v3             deep-block-5k              decode 68 MB          100ms     ok
gopkg.in/yaml.v3             deep-seq-100k              node   6 MB           10ms      error: yaml: block sequence entries are not allowed in this context
gopkg.in/yaml.v3             deep-seq-100k              decode 6 MB           0s        error: yaml: block sequence entries are not allowed in this context
gopkg.in/yaml.v3             scalar-100MB-plain         node   523 MB         1.15s     ok
gopkg.in/yaml.v3             scalar-100MB-plain         decode 463 MB         1.2s      ok
gopkg.in/yaml.v3             scalar-100MB-dquote        node   463 MB         1.07s     ok
gopkg.in/yaml.v3             scalar-100MB-dquote        decode 464 MB         1.07s     ok
gopkg.in/yaml.v3             1M-keys                    node   554 MB         1.04s     ok
gopkg.in/yaml.v3             1M-keys                    decode 480 MB         30.03s    KILLED: over 30 s
gopkg.in/yaml.v3             list-1M-items              node   285 MB         460ms     ok
gopkg.in/yaml.v3             list-1M-items              decode 256 MB         510ms     ok
go.yaml.in/yaml/v3           alias-bomb-9x9             node   5 MB           10ms      ok
go.yaml.in/yaml/v3           alias-bomb-9x9             decode 6 MB           10ms      error: yaml: document contains excessive aliasing
go.yaml.in/yaml/v3           deep-flow-100k             node   14 MB          10ms      error: yaml: exceeded max depth of 10000
go.yaml.in/yaml/v3           deep-flow-100k             decode 13 MB          10ms      error: yaml: exceeded max depth of 10000
go.yaml.in/yaml/v3           deep-block-5k              node   53 MB          90ms      ok
go.yaml.in/yaml/v3           deep-block-5k              decode 61 MB          100ms     ok
go.yaml.in/yaml/v3           deep-seq-100k              node   6 MB           10ms      error: yaml: block sequence entries are not allowed in this context
go.yaml.in/yaml/v3           deep-seq-100k              decode 6 MB           10ms      error: yaml: block sequence entries are not allowed in this context
go.yaml.in/yaml/v3           scalar-100MB-plain         node   563 MB         1.15s     ok
go.yaml.in/yaml/v3           scalar-100MB-plain         decode 502 MB         1.14s     ok
go.yaml.in/yaml/v3           scalar-100MB-dquote        node   463 MB         1.06s     ok
go.yaml.in/yaml/v3           scalar-100MB-dquote        decode 464 MB         1.05s     ok
go.yaml.in/yaml/v3           1M-keys                    node   550 MB         1.03s     ok
go.yaml.in/yaml/v3           1M-keys                    decode 458 MB         30.02s    KILLED: over 30 s
go.yaml.in/yaml/v3           list-1M-items              node   305 MB         430ms     ok
go.yaml.in/yaml/v3           list-1M-items              decode 254 MB         490ms     ok
go.yaml.in/yaml/v4           alias-bomb-9x9             node   5 MB           0s        ok
go.yaml.in/yaml/v4           alias-bomb-9x9             decode 6 MB           0s        error: go-yaml load error in constructor at L2.C22: document contains excessive aliasing
go.yaml.in/yaml/v4           deep-flow-100k             node   12 MB          10ms      error: go-yaml load error in scanner (while increasing flow level) at L1.C10004: exceeded max depth of 10000
go.yaml.in/yaml/v4           deep-flow-100k             decode 13 MB          10ms      error: go-yaml load error in scanner (while increasing flow level) at L1.C10004: exceeded max depth of 10000
go.yaml.in/yaml/v4           deep-block-5k              node   62 MB          90ms      ok
go.yaml.in/yaml/v4           deep-block-5k              decode 67 MB          90ms      ok
go.yaml.in/yaml/v4           deep-seq-100k              node   6 MB           0s        error: go-yaml load error in scanner at L1.C4: block sequence entries are not allowed in this context
go.yaml.in/yaml/v4           deep-seq-100k              decode 6 MB           0s        error: go-yaml load error in scanner at L1.C4: block sequence entries are not allowed in this context
go.yaml.in/yaml/v4           scalar-100MB-plain         node   464 MB         1.12s     ok
go.yaml.in/yaml/v4           scalar-100MB-plain         decode 564 MB         1.13s     ok
go.yaml.in/yaml/v4           scalar-100MB-dquote        node   564 MB         1.06s     ok
go.yaml.in/yaml/v4           scalar-100MB-dquote        decode 463 MB         1.05s     ok
go.yaml.in/yaml/v4           1M-keys                    node   574 MB         1.04s     ok
go.yaml.in/yaml/v4           1M-keys                    decode 466 MB         30.03s    KILLED: over 30 s
go.yaml.in/yaml/v4           list-1M-items              node   304 MB         430ms     ok
go.yaml.in/yaml/v4           list-1M-items              decode 283 MB         500ms     ok
github.com/goccy/go-yaml     alias-bomb-9x9             node   6 MB           10ms      ok
github.com/goccy/go-yaml     alias-bomb-9x9             decode 6 MB           10ms      ok
github.com/goccy/go-yaml     deep-flow-100k             node   6505 MB        1.1s      KILLED: RSS over 6 GB
github.com/goccy/go-yaml     deep-flow-100k             decode 6559 MB        910ms     KILLED: RSS over 6 GB
github.com/goccy/go-yaml     deep-block-5k              node   144 MB         240ms     ok
github.com/goccy/go-yaml     deep-block-5k              decode 156 MB         240ms     error: exceeded max depth
github.com/goccy/go-yaml     deep-seq-100k              node   189 MB         1.63s     error: [1:4] block sequence entries are not allowed in this context / >  1 
github.com/goccy/go-yaml     deep-seq-100k              decode 104 MB         1.64s     error: [1:4] block sequence entries are not allowed in this context / >  1 
github.com/goccy/go-yaml     scalar-100MB-plain         node   3288 MB        2.02s     ok
github.com/goccy/go-yaml     scalar-100MB-plain         decode 3332 MB        2s        ok
github.com/goccy/go-yaml     scalar-100MB-dquote        node   3209 MB        1.82s     ok
github.com/goccy/go-yaml     scalar-100MB-dquote        decode 3580 MB        1.88s     ok
github.com/goccy/go-yaml     1M-keys                    node   3275 MB        30.11s    KILLED: over 30 s
github.com/goccy/go-yaml     1M-keys                    decode 3343 MB        30.09s    KILLED: over 30 s
github.com/goccy/go-yaml     list-1M-items              node   846 MB         790ms     ok
github.com/goccy/go-yaml     list-1M-items              decode 923 MB         890ms     ok
```

The first hostile run accepted `k: - - - x` as a plain scalar in the hand parser (a bug, fixed by rejecting an indicator at the start of a plain value); the hand rows above are from the rerun.

### Behaviour (`behave.txt`)

```
=== fixture basic (subset in one file)
  [hand-written]
    L1  title          scalar str   plain  ["Plain string with spaces"]
    L2  quoted_s       scalar str   single ["it's single"]
    L3  quoted_d       scalar str   double ["say \"hi\"\ttab"]
    L4  count          scalar int   plain  ["42"]
    L5  neg            scalar int   plain  ["-7"]
    L6  date           scalar str   plain  ["2026-10-28"]
    L7  status         scalar str   plain  ["active"]
    L8  tags           list         flow   ["a" "b" "c d" "e"]
    L9  blocked_by     list         block  ["one" "two" "3"]
    L13 flush          list         block  ["x" "y"]
    L16 nothing        null   null  plain  []
    L17 unknown_key    scalar str   plain  ["kept"]
    L19 summary        scalar str   plain  ["012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789abcdefghij"]
  [gopkg.in/yaml.v3]
    L1  title          scalar str   plain  ["Plain string with spaces"]
    L2  quoted_s       scalar str   single ["it's single"]
    L3  quoted_d       scalar str   double ["say \"hi\"\ttab"]
    L4  count          scalar int   plain  ["42"]
    L5  neg            scalar int   plain  ["-7"]
    L6  date           scalar timestamp plain  ["2026-10-28"]
    L7  status         scalar str   plain  ["active"]
    L8  tags           list         flow   ["a" "b" "c d" "e"]
    L9  blocked_by     list         block  ["one" "two" "3"]
    L13 flush          list         block  ["x" "y"]
    L16 nothing        null   null  plain  [""]
    L17 unknown_key    scalar str   plain  ["kept"]
    L19 summary        scalar str   plain  ["012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789abcdefghij"]
  [go.yaml.in/yaml/v3]
    L1  title          scalar str   plain  ["Plain string with spaces"]
    L2  quoted_s       scalar str   single ["it's single"]
    L3  quoted_d       scalar str   double ["say \"hi\"\ttab"]
    L4  count          scalar int   plain  ["42"]
    L5  neg            scalar int   plain  ["-7"]
    L6  date           scalar timestamp plain  ["2026-10-28"]
    L7  status         scalar str   plain  ["active"]
    L8  tags           list         flow   ["a" "b" "c d" "e"]
    L9  blocked_by     list         block  ["one" "two" "3"]
    L13 flush          list         block  ["x" "y"]
    L16 nothing        null   null  plain  [""]
    L17 unknown_key    scalar str   plain  ["kept"]
    L19 summary        scalar str   plain  ["012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789abcdefghij"]
  [go.yaml.in/yaml/v4]
    L1  title          scalar str   plain  ["Plain string with spaces"]
    L2  quoted_s       scalar str   single ["it's single"]
    L3  quoted_d       scalar str   double ["say \"hi\"\ttab"]
    L4  count          scalar int   plain  ["42"]
    L5  neg            scalar int   plain  ["-7"]
    L6  date           scalar timestamp plain  ["2026-10-28"]
    L7  status         scalar str   plain  ["active"]
    L8  tags           list         flow   ["a" "b" "c d" "e"]
    L9  blocked_by     list         block  ["one" "two" "3"]
    L13 flush          list         block  ["x" "y"]
    L16 nothing        null   null  plain  [""]
    L17 unknown_key    scalar str   plain  ["kept"]
    L19 summary        scalar str   plain  ["012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789abcdefghij"]
  [github.com/goccy/go-yaml]
    L1  title          scalar str   plain  ["Plain string with spaces"]
    L2  quoted_s       scalar str   single ["it's single"]
    L3  quoted_d       scalar str   double ["say \"hi\"\ttab"]
    L4  count          scalar int   plain  ["42"]
    L5  neg            scalar int   plain  ["-7"]
    L6  date           scalar str   plain  ["2026-10-28"]
    L7  status         scalar str   plain  ["active"]
    L8  tags           list         flow   ["a" "b" "c d" "e"]
    L9  blocked_by     list         block  ["one" "two" "3"]
    L13 flush          list         block  ["x" "y"]
    L16 nothing        null   null  plain  []
    L17 unknown_key    scalar str   plain  ["kept"]
    L19 summary        scalar str   plain  ["012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789abcdefghij"]
=== fixture typing (values that YAML libraries retype; the tool wants strings kept as written)
  [hand-written]
    L1  country        scalar str   plain  ["no"]
    L2  flag           scalar str   plain  ["yes"]
    L3  ver            scalar str   plain  ["1.10"]
    L4  zip            scalar str   plain  ["08540"]
    L5  big            scalar int   plain  ["99999999999999999999"]
    L6  stamp          scalar str   plain  ["2026-10-28T10:00:00Z"]
    L7  tilde          scalar null  plain  ["~"]
    L8  nulltxt        scalar str   plain  ["null"]
    L9  hex            scalar str   plain  ["0x1F"]
    L10 under          scalar str   plain  ["1_000"]
    L11 id             scalar str   plain  ["261027-k7f2"]
    L12 colon_in_value scalar str   double ["a: b"]
    L13 hash_in_value  scalar str   plain  ["a#b"]
    L14 url            scalar str   plain  ["http://example.com/x"]
  [gopkg.in/yaml.v3]
    L1  country        scalar str   plain  ["no"]
    L2  flag           scalar str   plain  ["yes"]
    L3  ver            scalar float plain  ["1.10"]
    L4  zip            scalar float plain  ["08540"]
    L5  big            scalar float plain  ["99999999999999999999"]
    L6  stamp          scalar timestamp plain  ["2026-10-28T10:00:00Z"]
    L7  tilde          null   null  plain  ["~"]
    L8  nulltxt        null   null  plain  ["null"]
    L9  hex            scalar int   plain  ["0x1F"]
    L10 under          scalar int   plain  ["1_000"]
    L11 id             scalar str   plain  ["261027-k7f2"]
    L12 colon_in_value scalar str   double ["a: b"]
    L13 hash_in_value  scalar str   plain  ["a#b"]
    L14 url            scalar str   plain  ["http://example.com/x"]
  [go.yaml.in/yaml/v3]
    L1  country        scalar str   plain  ["no"]
    L2  flag           scalar str   plain  ["yes"]
    L3  ver            scalar float plain  ["1.10"]
    L4  zip            scalar float plain  ["08540"]
    L5  big            scalar float plain  ["99999999999999999999"]
    L6  stamp          scalar timestamp plain  ["2026-10-28T10:00:00Z"]
    L7  tilde          null   null  plain  ["~"]
    L8  nulltxt        null   null  plain  ["null"]
    L9  hex            scalar int   plain  ["0x1F"]
    L10 under          scalar int   plain  ["1_000"]
    L11 id             scalar str   plain  ["261027-k7f2"]
    L12 colon_in_value scalar str   double ["a: b"]
    L13 hash_in_value  scalar str   plain  ["a#b"]
    L14 url            scalar str   plain  ["http://example.com/x"]
  [go.yaml.in/yaml/v4]
    L1  country        scalar str   plain  ["no"]
    L2  flag           scalar str   plain  ["yes"]
    L3  ver            scalar float plain  ["1.10"]
    L4  zip            scalar float plain  ["08540"]
    L5  big            scalar float plain  ["99999999999999999999"]
    L6  stamp          scalar timestamp plain  ["2026-10-28T10:00:00Z"]
    L7  tilde          null   null  plain  ["~"]
    L8  nulltxt        null   null  plain  ["null"]
    L9  hex            scalar int   plain  ["0x1F"]
    L10 under          scalar int   plain  ["1_000"]
    L11 id             scalar str   plain  ["261027-k7f2"]
    L12 colon_in_value scalar str   double ["a: b"]
    L13 hash_in_value  scalar str   plain  ["a#b"]
    L14 url            scalar str   plain  ["http://example.com/x"]
  [github.com/goccy/go-yaml]
    L1  country        scalar str   plain  ["no"]
    L2  flag           scalar str   plain  ["yes"]
    L3  ver            scalar float plain  ["1.10"]
    L4  zip            scalar str   plain  ["08540"]
    L5  big            scalar str   plain  ["99999999999999999999"]
    L6  stamp          scalar str   plain  ["2026-10-28T10:00:00Z"]
    L7  tilde          null   null  plain  []
    L8  nulltxt        null   null  plain  []
    L9  hex            scalar int   plain  ["0x1F"]
    L10 under          scalar int   plain  ["1_000"]
    L11 id             scalar str   plain  ["261027-k7f2"]
    L12 colon_in_value scalar str   double ["a: b"]
    L13 hash_in_value  scalar str   plain  ["a#b"]
    L14 url            scalar str   plain  ["http://example.com/x"]
=== fixture anchor (anchor, alias and tag: outside the subset; tool must reject or detect)
  [hand-written]
    ERROR: line 1: '&' starts an anchor, alias, tag, block scalar or map: outside the subset   => file line 2
  [gopkg.in/yaml.v3]
    L1  base           scalar str   plain  ["value"]
    L2  copy           alias               []
    L3  tagged         scalar str   plain  ["12"]
  [go.yaml.in/yaml/v3]
    L1  base           scalar str   plain  ["value"]
    L2  copy           alias               []
    L3  tagged         scalar str   plain  ["12"]
  [go.yaml.in/yaml/v4]
    L1  base           scalar str   plain  ["value"]
    L2  copy           alias               []
    L3  tagged         scalar str   plain  ["12"]
  [github.com/goccy/go-yaml]
    L1  base           other:*ast.AnchorNode              []
    L2  copy           alias               []
    L3  tagged         other:*ast.TagNode              []
=== fixture nested (nesting and flow map: outside the subset)
  [hand-written]
    ERROR: line 3: indented line: nesting and multi-line values are outside the subset   => file line 4
  [gopkg.in/yaml.v3]
    L1  a              scalar int   plain  ["1"]
    L2  b              map                 []
    L5  e              map                 []
  [go.yaml.in/yaml/v3]
    L1  a              scalar int   plain  ["1"]
    L2  b              map                 []
    L5  e              map                 []
  [go.yaml.in/yaml/v4]
    L1  a              scalar int   plain  ["1"]
    L2  b              map                 []
    L5  e              map                 []
  [github.com/goccy/go-yaml]
    L1  a              scalar int   plain  ["1"]
    L2  b              map                 []
    L5  e              map                 []
=== fixture multiline (block scalars and multi-line plain: outside the subset)
  [hand-written]
    ERROR: line 1: '|' starts an anchor, alias, tag, block scalar or map: outside the subset   => file line 2
  [gopkg.in/yaml.v3]
    L1  a              scalar str   plain  ["literal\nblock\n"]
    L4  b              scalar str   plain  ["folded\n"]
    L6  c              scalar str   plain  ["plain continued"]
  [go.yaml.in/yaml/v3]
    L1  a              scalar str   plain  ["literal\nblock\n"]
    L4  b              scalar str   plain  ["folded\n"]
    L6  c              scalar str   plain  ["plain continued"]
  [go.yaml.in/yaml/v4]
    L1  a              scalar str   plain  ["literal\nblock\n"]
    L4  b              scalar str   plain  ["folded\n"]
    L6  c              scalar str   plain  ["plain continued"]
  [github.com/goccy/go-yaml]
    L1  a              other:*ast.LiteralNode              []
    L4  b              other:*ast.LiteralNode              []
    L6  c              scalar str   plain  ["plain continued"]
=== fixture empty (empty frontmatter)
  [hand-written]
  [gopkg.in/yaml.v3]
  [go.yaml.in/yaml/v3]
  [go.yaml.in/yaml/v4]
  [github.com/goccy/go-yaml]
=== fixture crlf (CRLF line endings)
  [hand-written]
    L1  title          scalar str   plain  ["x"]
    L2  status         scalar str   plain  ["active"]
    L3  tags           list         flow   ["a"]
  [gopkg.in/yaml.v3]
    L1  title          scalar str   plain  ["x"]
    L2  status         scalar str   plain  ["active"]
    L3  tags           list         flow   ["a"]
  [go.yaml.in/yaml/v3]
    L1  title          scalar str   plain  ["x"]
    L2  status         scalar str   plain  ["active"]
    L3  tags           list         flow   ["a"]
  [go.yaml.in/yaml/v4]
    L1  title          scalar str   plain  ["x"]
    L2  status         scalar str   plain  ["active"]
    L3  tags           list         flow   ["a"]
  [github.com/goccy/go-yaml]
    L1  title          scalar str   plain  ["x"]
    L2  status         scalar str   plain  ["active"]
    L3  tags           list         flow   ["a"]
=== fixture bad-unterminated-quote (error expected near file line 3)
  [hand-written]
    ERROR: line 2: unterminated double-quoted string   => file line 3
  [gopkg.in/yaml.v3]
    ERROR: yaml: line 2: found unexpected end of stream   => file line 3
  [go.yaml.in/yaml/v3]
    ERROR: yaml: line 2: found unexpected end of stream   => file line 3
  [go.yaml.in/yaml/v4]
    ERROR: go-yaml load error in scanner (while scanning a quoted scalar) at L2.C10-L4.C1: found unexpected end of stream
  [github.com/goccy/go-yaml]
    ERROR: [2:10] could not find end character of double-quoted text /    1 | title: ok / >  2 | summary: "never closed /    3 | status: active /                 ^ /    => file line 3
=== fixture bad-unterminated-list (error expected near file line 4)
  [hand-written]
    ERROR: line 3: unterminated inline list   => file line 4
  [gopkg.in/yaml.v3]
    ERROR: yaml: line 2: did not find expected ',' or ']'   => file line 3
  [go.yaml.in/yaml/v3]
    ERROR: yaml: line 2: did not find expected ',' or ']'   => file line 3
  [go.yaml.in/yaml/v4]
    ERROR: go-yaml load error in parser (while parsing a flow sequence) at L3.C7-L4.C6: did not find expected ',' or ']'
  [github.com/goccy/go-yaml]
    ERROR: [4:1] ',' or ']' must be specified /    1 | title: ok /    2 | status: active /    3 | tags: [a, b / >  4 | owner: me /        ^ /    => file line 5
=== fixture bad-colon (error expected at file line 3)
  [hand-written]
    ERROR: line 2: mapping values are not allowed here   => file line 3
  [gopkg.in/yaml.v3]
    ERROR: yaml: line 2: mapping values are not allowed in this context   => file line 3
  [go.yaml.in/yaml/v3]
    ERROR: yaml: line 2: mapping values are not allowed in this context   => file line 3
  [go.yaml.in/yaml/v4]
    ERROR: go-yaml load error in scanner at L2.C11: mapping values are not allowed in this context
  [github.com/goccy/go-yaml]
    ERROR: [2:10] mapping value is not allowed in this context /    1 | title: ok / >  2 | summary: a: b /                 ^ /    3 | status: x   => file line 3
=== fixture bad-dup (error expected at file line 4)
  [hand-written]
    ERROR: line 3: duplicate key status   => file line 4
  [gopkg.in/yaml.v3]
    L1  title          scalar str   plain  ["ok"]
    L2  status         scalar str   plain  ["a"]
    L3  status         scalar str   plain  ["b"]
  [go.yaml.in/yaml/v3]
    L1  title          scalar str   plain  ["ok"]
    L2  status         scalar str   plain  ["a"]
    L3  status         scalar str   plain  ["b"]
  [go.yaml.in/yaml/v4]
    L1  title          scalar str   plain  ["ok"]
    L2  status         scalar str   plain  ["a"]
    L3  status         scalar str   plain  ["b"]
  [github.com/goccy/go-yaml]
    ERROR: [3:1] mapping key "status" already defined at [2:1] /    1 | title: ok /    2 | status: a / >  3 | status: b /        ^ /    => file line 4
=== fixture bad-tab (error expected at file line 4)
  [hand-written]
    ERROR: line 3: tab in indentation   => file line 4
  [gopkg.in/yaml.v3]
    ERROR: yaml: line 3: found character that cannot start any token   => file line 4
  [go.yaml.in/yaml/v3]
    ERROR: yaml: line 3: found character that cannot start any token   => file line 4
  [go.yaml.in/yaml/v4]
    ERROR: go-yaml load error in scanner (while scanning for the next token) at L3.C1: found character that cannot start any token
  [github.com/goccy/go-yaml]
    ERROR: [3:1] found character '	' that cannot start any token /    1 | title: ok /    2 | blocked_by: / >  3 | 	- x /        ^ /    4 | status: a   => file line 4
=== fixture bad-indent (error expected at file line 3)
  [hand-written]
    ERROR: line 2: indented line: nesting and multi-line values are outside the subset   => file line 3
  [gopkg.in/yaml.v3]
    ERROR: yaml: line 2: mapping values are not allowed in this context   => file line 3
  [go.yaml.in/yaml/v3]
    ERROR: yaml: line 2: mapping values are not allowed in this context   => file line 3
  [go.yaml.in/yaml/v4]
    ERROR: go-yaml load error in scanner at L2.C9: mapping values are not allowed in this context
  [github.com/goccy/go-yaml]
    ERROR: [1:8] mapping value is not allowed in this context / >  1 | title: ok /    2 |   status: a /               ^ /    3 | owner: me   => file line 2
=== edit round-trip (block only; want = naive text edit)
  no-op        hand-written                 BYTE-IDENTICAL
  no-op        gopkg.in/yaml.v3             DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: active # state
updated: 2026-10-28
tags: [a, b, 'c']
blocked_by:
  - x # four-space indent
  - "y"
summary: 'it''s kept'
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: active   # state
updated: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
  no-op        go.yaml.in/yaml/v3           DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: active # state
updated: 2026-10-28
tags: [a, b, 'c']
blocked_by:
  - x # four-space indent
  - "y"
summary: 'it''s kept'
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: active   # state
updated: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
  no-op        go.yaml.in/yaml/v4           DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: active # state
updated: 2026-10-28
tags: [a, b, 'c']
blocked_by:
  - x # four-space indent
  - "y"
summary: 'it''s kept'
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: active   # state
updated: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
  no-op        github.com/goccy/go-yaml     DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: active # state
updated: 2026-10-28

tags: [a, b, 'c']
blocked_by:
    - x # four-space indent
    - "y"
summary: 'it''s kept'
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: active   # state
updated: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
  set-status   hand-written                 BYTE-IDENTICAL
  set-status   gopkg.in/yaml.v3             DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: done # state
updated: 2026-10-28
tags: [a, b, 'c']
blocked_by:
  - x # four-space indent
  - "y"
summary: 'it''s kept'
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: done   # state
updated: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
  set-status   go.yaml.in/yaml/v3           DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: done # state
updated: 2026-10-28
tags: [a, b, 'c']
blocked_by:
  - x # four-space indent
  - "y"
summary: 'it''s kept'
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: done   # state
updated: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
  set-status   go.yaml.in/yaml/v4           DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: done # state
updated: 2026-10-28
tags: [a, b, 'c']
blocked_by:
  - x # four-space indent
  - "y"
summary: 'it''s kept'
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: done   # state
updated: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
  set-status   github.com/goccy/go-yaml     DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: done
updated: 2026-10-28

tags: [a, b, 'c']
blocked_by:
    - x # four-space indent
    - "y"
summary: 'it''s kept'
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: done   # state
updated: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
  rename-updated hand-written                 BYTE-IDENTICAL
  rename-updated gopkg.in/yaml.v3             DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: active # state
updated_at: 2026-10-28
tags: [a, b, 'c']
blocked_by:
  - x # four-space indent
  - "y"
summary: 'it''s kept'
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: active   # state
updated_at: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
  rename-updated go.yaml.in/yaml/v3           DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: active # state
updated_at: 2026-10-28
tags: [a, b, 'c']
blocked_by:
  - x # four-space indent
  - "y"
summary: 'it''s kept'
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: active   # state
updated_at: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
  rename-updated go.yaml.in/yaml/v4           DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: active # state
updated_at: 2026-10-28
tags: [a, b, 'c']
blocked_by:
  - x # four-space indent
  - "y"
summary: 'it''s kept'
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: active   # state
updated_at: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
  rename-updated github.com/goccy/go-yaml     DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: active # state
updated: 2026-10-28

tags: [a, b, 'c']
blocked_by:
    - x # four-space indent
    - "y"
summary: 'it''s kept'
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: active   # state
updated_at: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
  add-default  hand-written                 BYTE-IDENTICAL
  add-default  gopkg.in/yaml.v3             DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: active # state
updated: 2026-10-28
tags: [a, b, 'c']
blocked_by:
  - x # four-space indent
  - "y"
summary: 'it''s kept'
priority: 3
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: active   # state
updated: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
priority: 3
  add-default  go.yaml.in/yaml/v3           DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: active # state
updated: 2026-10-28
tags: [a, b, 'c']
blocked_by:
  - x # four-space indent
  - "y"
summary: 'it''s kept'
priority: 3
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: active   # state
updated: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
priority: 3
  add-default  go.yaml.in/yaml/v4           DIFFERS
--- got
# leading comment
id: 261027-k7f2
title: "Spaced   and quoted" # keep this
status: active # state
updated: 2026-10-28
tags: [a, b, 'c']
blocked_by:
  - x # four-space indent
  - "y"
summary: 'it''s kept'
priority: 3
--- want
# leading comment
id: 261027-k7f2
title:   "Spaced   and quoted"   # keep this
status: active   # state
updated: 2026-10-28

tags: [ a,b ,'c' ]
blocked_by:
    - x   # four-space indent
    - "y"
summary: 'it''s kept'
priority: 3
  add-default  github.com/goccy/go-yaml     ERROR add not attempted via AST (would be a text splice)
```
