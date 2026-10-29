---
type: design
status: draft
created_at: 2026-10-28
updated_at: 2026-10-28
authored_by: agent
---
# Go display-width library evaluation for `akmctl` (draft for review)

Evaluation for [261028-0snt-go-width-library-evaluation](../work/done/261028-0snt-go-width-library-evaluation.md), serving [agent-ready-repo-spec](../spec/agent-ready-repo-spec.md) (section 4, "The generated index") and decision Y16 in [agent-ready-repo-decision-log](agent-ready-repo-decision-log.md). It answers one question: how should the Go tool compute the display width of table cells, so that the generated ADR index is padded the way Prettier pads it? Criteria, in Alexander's order: (1) baseline: output matches Prettier, (2) performance, (3) security and supply chain, (4) popularity and support as a tiebreaker. Run on 2026-10-28; machine and toolchain in the performance section. The method and layout follow [go-yaml-library-evaluation](go-yaml-library-evaluation.md); the Prettier test it extends is in [akm-spec-review](akm-spec-review.md), section 3.

## Summary

**Recommendation: `github.com/clipperhouse/displaywidth` v0.11.0, behind a one-function wrapper (`cellWidth(s string) int`) in `akmctl`, with a test that compares the wrapper to Prettier at every Prettier pin bump (Y17).** The tool should also reject invalid UTF-8 and normalise whitespace in a cell before it measures.

- **Baseline (1), the deciding criterion.** Ground truth is Prettier 3.9.10 itself (published 2026-10-10, 18 days old): one-column tables, one per cell, formatted by Prettier, with the width it chose read back from the padding. 208 fixture cells; 2 are not scored because Prettier rewrites the cell text (runs of spaces, a tab), so **206 are scored**. Cells where the width differs from Prettier:

    | Candidate | Wrong, of 206 | Wrong, of the 129 common-text cells |
    |---|---|---|
    | `github.com/clipperhouse/displaywidth` v0.11.0 | **32** | **3** |
    | `github.com/rivo/uniseg` v0.4.7 | 38 | 3 |
    | `golang.org/x/text/width` v0.43.0 (per-rune use) | 40 | 15 |
    | `github.com/mattn/go-runewidth` v0.0.31 | 42 (71 under a Japanese locale) | 15 |
    | Baseline: embedded table, wide ranges only | 29 | 13 |
    | Baseline: embedded table plus emoji rules, tuned to Prettier's algorithm | 5 | 0 |
    | Baseline: rune count | 84 | 51 |

    "Common text" is ASCII, accented Latin (precomposed and combining), Greek and Cyrillic, CJK, kana, hangul, fullwidth forms, ambiguous-width symbols, single emoji and emoji with a variation selector. **No library reproduces Prettier.** The reason is that Prettier's rule is not Unicode's. It counts code points: control characters, U+0300 to U+036F and variation selectors count 0; everything else counts 1, or 2 for wide and fullwidth characters; an emoji sequence found by a regular expression counts 2 (1 if it is a single text-style emoji). So Prettier counts a zero-width space, a zero-width joiner, and Thai, Arabic and Devanagari combining marks as 1 column each, while the three libraries that implement Unicode grapheme rules count them as 0. See "What Prettier does".
- **Which cells the recommended library gets wrong (32 of 206):** 10 format characters (zero-width space, joiners, soft hyphen, word joiner, BOM, direction marks), 8 combining marks of Thai, Arabic, Devanagari, Khmer and Tibetan, 2 conjoining hangul jamo, 2 keycap sequences, 2 regional-indicator (flag) edge cases, 2 skin-tone edge cases, 2 control characters, 1 odd ZWJ pair, 1 letter followed by U+FE0F, 2 invalid UTF-8 cells. All are text an ADR title is unlikely to contain; the one realistic risk is an invisible character pasted into a title (zero-width space, soft hyphen).
- **Is the no-dependency baseline good enough?** Rune count: no (CJK and every emoji wrong; 51 of 129 common-text cells). A table of wide ranges alone: no (13 of 129; it misses every emoji with U+FE0F, such as the check mark and the warning sign). **A table plus emoji rules mirroring Prettier's algorithm: yes, and it is the closest of all** (5 of 206, 0 of 129). But I wrote it after seeing the Prettier results, so its score is optimistic; it costs about 110 lines and a table regenerated from the Unicode data files at each Unicode release; and that is the maintenance Y16 chose to avoid. On criterion (1) alone it wins; I recommend the dependency anyway because the penalty for a mismatch is small (the generated file is also listed in `.prettierignore`, and rule A4 compares the generator's own bytes), and the maintainer asked for a dependency. If fidelity to Prettier outranks maintenance, take the baseline. That is Alexander's call.
- **Performance (2):** not a differentiator. 26.7 ns per cell on an ASCII title, 66 ns on mixed text for `displaywidth`; the slowest candidate, `uniseg`, takes 0.55 to 0.63 us per cell. 5,000 cells cost 0.13 to 3 ms in total.
- **Security (3):** no panics and no super-linear behaviour on any hostile input (100 MB strings, a 10-million-mark grapheme cluster, 5-million-element ZWJ chains, 100 MB of invalid bytes). `displaywidth` adds one module (`clipperhouse/uax29/v2`, same author), no advisories for either. `go-runewidth` reads the locale from the environment, so the same input gives different widths on different machines (42 wrong in a default locale, 71 under `LC_ALL=ja_JP.UTF-8`); that alone rules it out for a generator whose output is committed and compared byte for byte.
- **Popularity (4):** the weak point of the recommendation. `displaywidth` is young (repository created 2025-10-08), has one maintainer, 19 stars and 94 importing packages, against 5,958 for `go-runewidth` and 1,405 for `uniseg`. It breaks no tie here; criterion (1) decided. Its author also contributes to `go-runewidth`, which now depends on the same `uax29` module, so the tables under it are in wide use.
- **Runner-up:** `github.com/rivo/uniseg` (zero dependencies, 3 of 129 common-text cells wrong) if `displaywidth` is judged too young. It has had no release since 2024-02-08 and its tables are Unicode 15.0, so four newer emoji (Unicode 16 and 17) measure 1 column instead of 2. Not recommended: `go-runewidth` (locale dependence; 13 of 19 emoji-with-U+FE0F cells wrong; two open issues or pull requests about exactly that) and `golang.org/x/text/width` (per-rune only: it has no zero-width concept, so every combining sequence is wrong, 8 of 8 decomposed-Latin cells; and its tables are Unicode 17.0.0 on Go 1.27 and later but 15.0.0 on older Go).

This is a recommendation for Alexander to confirm. It changes nothing in the spec, [261027-ensf-akm-an-agent-knowledge-management](../adr/261027-ensf-akm-an-agent-knowledge-management.md) or the decision log. If Alexander confirms, Y16 can name the module and version.

## What Prettier does

Read from `prettier@3.9.10` (`index.mjs`, function `getStringWidth`, used for table cells); the bundled packages are `get-east-asian-width` 1.6.0 and `emoji-regex` 10.6.0 (see `THIRD-PARTY-NOTICES.md` in the package). The algorithm, as observed in that source and confirmed by the 206 results:

1. If the text is only characters U+0020 to U+007F, the width is the string length (so U+007F counts 1).
2. Otherwise every match of the emoji regular expression counts 2, or 1 when the whole match is one of a fixed list of "narrow" text-style emoji such as the check mark U+2714 or the heart U+2764 (with U+FE0F added the match is longer and counts 2).
3. In what remains, each code point counts 0 if it is a control character (U+0000 to U+001F, U+007F to U+009F), U+0300 to U+036F or U+FE00 to U+FE0F; otherwise 2 if `get-east-asian-width` says Wide or Fullwidth; otherwise 1.

Two consequences matter for the generator. Prettier does not use grapheme clusters, so a ZWJ sequence that is not in the emoji regular expression (for example two smileys joined by U+200D) counts 2 + 1 + 2 = 5. And Prettier normalises cell text before it measures: runs of spaces become one space and a tab becomes a space, so the generator must do that to its own cells first.

## Candidates

| # | Candidate | Role |
|---|---|---|
| 1 | `golang.org/x/text/width` v0.43.0 | the standard-adjacent option; used per rune: Wide and Fullwidth count 2, anything else 1 |
| 2 | `github.com/mattn/go-runewidth` v0.0.31 | the most used; `StringWidth` with the default condition |
| 3 | `github.com/rivo/uniseg` v0.4.7 | grapheme-aware; `StringWidth` |
| 4 | `github.com/clipperhouse/displaywidth` v0.11.0 | not in the brief; added because `go-runewidth` v0.0.31 now builds on the same author's `uax29` module and this is that author's width library; `String` |
| 5 | Baseline: rune count | `utf8.RuneCountInString`, with the file in `.prettierignore` |
| 6 | Baseline: embedded table, wide ranges only | East Asian Wide and Fullwidth ranges from Unicode 17.0.0 `EastAsianWidth.txt`, generated by `gen_table.py`, plus Prettier's zero-width rules |
| 7 | Baseline: embedded table plus emoji rules | as 6, plus a small recogniser for emoji sequences (variation selector, skin tone, ZWJ chain, tag sequence, flag pair, keycap) built from `emoji-data.txt` |

Candidates 6 and 7 are my own code, written to mimic Prettier's rules, because the contract is "same as Prettier", not "same as Unicode". Candidates 1 to 4 are used as their authors document them, with no wrapper, so they are not tuned. That asymmetry is the main caveat on the comparison, and why row 7's score is flattering.

## Comparison on the four criteria

Facts as of 2026-10-28.

| | `x/text/width` | `go-runewidth` | `uniseg` | `displaywidth` |
|---|---|---|---|---|
| Version / date | v0.43.0, 2026-10-09 | v0.0.31, 2026-10-08 | v0.4.7, 2024-02-08 | v0.11.0, 2026-02-17 |
| Licence | BSD-3-Clause | MIT | MIT | MIT |
| Maintainers (contributors endpoint) | Go team: mpvl 517 commits, nigeltao 45, gopherbot 29 | mattn 186, next 11 | rivo 50, next 14 | clipperhouse 104, sole contributor |
| Last push | 2026-10-21 (whole `x/text` repo) | 2026-10-08 | 2024-05-31 | 2026-05-17 |
| Open issues (GitHub count, includes PRs) | 30 (whole `x/text` repo) | 6 | 12 | 3 |
| Issues about width | not examined | open: "Honor VS16 emoji presentation in grapheme widths", "Count an emoji presentation sequence as 2 cells", "Update the tables to Unicode 18.0.0" | "Improve emoji consistency with older terminals" | not examined |
| Dependencies as seen from a consumer | none in the build (module graph lists `x/tools`, `x/mod`, `x/sync`; `go.sum` gets 2 lines) | 1 module: `clipperhouse/uax29/v2` (2 packages built; minimum v2.2.0) | none | 1 module: `clipperhouse/uax29/v2` v2.7.0 |
| Unicode version | 17.0.0 on Go 1.27+, 15.0.0 on older Go (build tag) | 17.0.0 (generator URLs) | 15.0.0 | not stated in the docs I read; correct on the Unicode 16 and 17 emoji tested |
| Environment dependence | none | reads `LC_ALL`/`LC_CTYPE`/`LANG` and `RUNEWIDTH_EASTASIAN` at start-up | none | none in the built packages |
| Importing packages (pkg.go.dev) | 883 | 5,958 | 1,405 | 94 |
| Stars (repository) | 810 (`golang/text`) | 725 | 738 | 19 |
| Advisories (osv.dev) | 8 entries for `x/text`, none for `width`: affected packages are `unicode/norm` (fixed 0.39.0), `secure/precis` (fixed 0.41.0), `language`, and one from 2021 (fixed 0.3.3) | 0 | 0 | 0 (and 0 for `uax29/v2`) |
| **(1) Baseline**: wrong of 206 / of 129 common | 40 / 15 | 42 / 15 | 38 / 3 | **32 / 3** |
| **(2) Performance**, ns per cell (ASCII / mixed) | 153 / 174 | 38 / 126 | 546 / 631 | **27 / 66** |
| **(3) Security** | linear; no panic | linear; no panic; locale-dependent output | linear, slowest on 100 MB invalid bytes (4.9 s); no panic | linear; no panic |
| **(4) Popularity** | high trust, mid use | highest | mid, unmaintained since 2024 | lowest, young |

Sources:
- Versions and dates: `https://proxy.golang.org/<module>/@latest`; `go.mod` files from the same proxy.
- Repository state, licence, contributors, open issue titles: `gh api repos/...` (GitHub API).
- Importers: pkg.go.dev "Imported By" tab, count of packages including internal ones (not modules).
- Advisories: `https://api.osv.dev/v1/query` for each module path, and `https://api.osv.dev/v1/vulns/<id>` for the fixed versions.
- Dependencies: `go list -m all` and `go list -deps` in a throwaway module per candidate, with a `go.sum` generated by `go mod tidy`.
- Unicode versions: `UnicodeVersion` constants and `//go:build` tags in `x/text/width`; the generator source of `go-runewidth` (`script/generate.go`); the header of `uniseg/eastasianwidth.go`.
- Environment access: a search of the built packages for `Getenv` (only `go-runewidth`).

## Per-class results

206 scored cells; the figure is the number of cells in the class where the width differs from Prettier. Class sizes are in the second column. The first five candidates are the libraries and baselines; columns follow the order of the run.

| Class | Cells | rune count | table (wide only) | table + emoji | `displaywidth` | `go-runewidth` | `uniseg` | `x/text/width` |
|---|---|---|---|---|---|---|---|---|
| ASCII | 10 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| Latin, precomposed | 10 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| Latin, combining marks | 8 | 8 | 0 | 0 | 0 | 0 | 0 | 8 |
| Greek, Cyrillic | 7 | 1 | 0 | 0 | 0 | 0 | 0 | 1 |
| CJK ideographs | 9 | 9 | 0 | 0 | 0 | 0 | 0 | 0 |
| Kana | 6 | 5 | 0 | 0 | 0 | 0 | 0 | 0 |
| Hangul | 6 | 6 | 0 | 0 | 2 | 2 | 2 | 0 |
| Fullwidth forms | 8 | 6 | 0 | 0 | 0 | 0 | 0 | 0 |
| Ambiguous width | 20 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| Emoji, single | 26 | 14 | 0 | 0 | 0 | 0 | 1 | 0 |
| Emoji with U+FE0F | 19 | 2 | 13 | 0 | 1 | 13 | 0 | 6 |
| Emoji, keycap | 5 | 3 | 0 | 0 | 2 | 5 | 5 | 3 |
| Emoji, skin tone | 6 | 2 | 5 | 1 | 2 | 1 | 3 | 4 |
| Emoji, ZWJ | 8 | 8 | 7 | 1 | 1 | 1 | 1 | 7 |
| Emoji, flags and tags | 6 | 2 | 1 | 0 | 2 | 0 | 2 | 1 |
| New in Unicode 15.1 to 17 | 6 | 6 | 0 | 0 | 0 | 0 | 4 | 0 |
| Zero-width and format characters | 11 | 1 | 0 | 0 | 10 | 10 | 10 | 1 |
| Other scripts (Thai, Arabic, Hebrew, Devanagari, Tibetan, Khmer) | 10 | 0 | 0 | 0 | 8 | 6 | 6 | 0 |
| Control characters | 11 | 6 | 1 | 1 | 2 | 2 | 2 | 6 |
| Long strings (200 to 100,000 characters) | 6 | 3 | 0 | 0 | 0 | 0 | 0 | 1 |
| Invalid UTF-8 | 8 | 2 | 2 | 2 | 2 | 2 | 2 | 2 |
| **Total** | **206** | **84** | **29** | **5** | **32** | **42** | **38** | **40** |

What causes the differences:
- **Rune count** is wrong wherever a character is not one column: every wide character and every emoji, and in the other direction every combining mark.
- **`x/text/width`** has a width per rune and no notion of zero width or sequences. Every combining mark counts 1 (8 of 8 decomposed-Latin cells wrong), and so do U+FE0F and U+200D, so emoji with a variation selector and every ZWJ sequence come out too wide (U+2705 followed by U+FE0F gives 3, not 2). Its answers on single characters are right, including the Unicode 17 additions, but only on Go 1.27 or later; on older Go it uses Unicode 15 tables.
- **`go-runewidth`** (default condition) counts an emoji with U+FE0F as 1 when the base character is text-style (13 of 19 wrong; the open issues listed above are about this). Its keycap sequences are 1. The same code gives 71 wrong in a Japanese locale, where ambiguous-width characters (Greek and Cyrillic letters, arrows, box drawing) count 2.
- **`uniseg`** and **`displaywidth`** agree with each other on 188 of 206 cells. Both follow Unicode: zero-width for format characters and for combining marks of any script, which Prettier counts as 1. `uniseg` additionally lacks the emoji added after Unicode 15 (U+1FAE9, U+1FA89, U+1FAC8, U+1FAC6 count 1) and one text-style wide character (U+3297). `displaywidth` counts a letter followed by U+FE0F as 2.
- **Hangul:** the conjoining jamo sequences (a lead consonant followed by a vowel) are one syllable for the three Unicode-based libraries, so they count 2; Prettier counts each jamo separately (2 + 1, and 2 + 1 + 1). Precomposed syllables and compatibility jamo all agree.
- **Invalid UTF-8:** all candidates disagree with Prettier on 2 of 8 cells, though not the same two. Prettier's file reader replaces each maximal invalid subsequence with one U+FFFD, while Go sees one error per byte (`displaywidth` sees a whole invalid sequence as one). The tool should reject invalid UTF-8 earlier, since it cannot be valid Markdown for Prettier either.
- **U+007F** counts 1 in Prettier (the ASCII fast path); the table baselines and the libraries count 0.
- **Long strings:** all candidates agree with Prettier on 200, 5,000 and 100,000 ASCII characters, 1,000 ideographs and 500 emoji; the rune count is wrong on the CJK and emoji ones; `x/text/width` is wrong on a long mixed string with combining marks.

## Performance

Go `testing.Benchmark`, best of 5 runs, two deterministic corpora of 5,000 cells each: "ascii" (ADR-style titles, 44 bytes on average) and "mixed" (70% ASCII, 10% accented Latin, 10% Japanese text, 10% with emoji, 50 bytes on average). Machine: Apple M1 Max, 64 GB; `go1.27.2 darwin/arm64`. The library numbers do not include any wrapper.

| Candidate | ASCII: ns per cell | mixed: ns per cell | allocations per cell |
|---|---|---|---|
| rune count | 22 | 28 | 0 |
| `displaywidth` | **27** | **66** | 0 |
| `go-runewidth` | 38 | 126 | 0 |
| `x/text/width` | 153 | 174 | 0 |
| table (wide only) | 196 | 227 | 0.7 (the `[]rune` copy) |
| table + emoji | 534 | 585 | 0.7 |
| `uniseg` | 546 | 631 | 0 |

An index of 5,000 rows with 6 columns is 30,000 cells; even the slowest candidate measures it in under 20 ms. The criterion is met by all, and none of the differences affects the choice. My table code is unoptimised (it copies the string to runes and binary-searches the ranges).

## Security

Each case ran in a child process of the evaluation binary; the parent polled resident memory every 40 ms and would kill at 6 GB or 30 s. Neither limit was reached. Peak memory is of the whole process including the 100 MB input string. Full numbers are in the appendix.

| Input | `displaywidth` | `go-runewidth` | `uniseg` | `x/text/width` | table + emoji |
|---|---|---|---|---|---|
| 100 MB ASCII | 50 ms, 106 MB | 67 ms, 105 MB | 1.26 s, 105 MB | 340 ms, 106 MB | 1.18 s, 506 MB |
| 100 MB ideographs | 535 ms, 106 MB | 182 ms, 108 MB | 1.60 s, 106 MB | 255 ms, 105 MB | 838 ms, 239 MB |
| 100 MB invalid bytes (0xFF) | 1.12 s, 106 MB | 1.33 s, 108 MB | 4.89 s, 106 MB | 744 ms, 106 MB | 2.47 s, 506 MB |
| 100 MB emoji | 523 ms, 105 MB | 303 ms, 108 MB | 1.10 s, 106 MB | 225 ms, 105 MB | 599 ms, 206 MB |
| One base letter plus 10 million combining marks (one grapheme cluster) | 43 ms, 44 MB | 83 ms, 46 MB | 210 ms, 44 MB | 59 ms, 44 MB | 135 ms, 82 MB |
| 5 million emoji joined by ZWJ | 61 ms, 72 MB | 104 ms, 74 MB | 242 ms, 72 MB | 125 ms, 72 MB | 106 ms, 111 MB |

- **No panics, no quadratic behaviour, no unbounded allocation:** the libraries allocate nothing; my table code allocates 4 bytes per input byte (the rune slice), which is why it peaks at 506 MB on 100 MB. Clusters that grow without bound (combining marks, ZWJ chains, regional indicators, tag characters, variation selector and skin-tone chains, all 5 to 10 million elements) cost milliseconds.
- **A size cap is still needed:** an ADR title or summary is at most a few hundred characters (the spec limits `summary` to 160), so `akmctl` should refuse a cell above, say, 4 KB before it measures. That is one line and applies to any candidate.
- **Supply chain:** `displaywidth` and `go-runewidth` both depend on `clipperhouse/uax29/v2`, one maintainer's module. `displaywidth` is more exposed to that maintainer (it is also his; 104 of its 104 commits). The built packages of all four use no `os`, `net`, `unsafe` or `exec` except `go-runewidth` (`os` for the environment, and platform files). `displaywidth` has files with `os` and `unsafe` in its repository (a command and test data helpers), none of them in the package that `import` builds; I did not audit the source beyond this search. All four modules were fetched through the Go module proxy and are pinned by `go.sum` (appendix).
- **History:** `golang.org/x/text` has eight advisories, none in `width`; v0.43.0 includes all the fixes. `go-runewidth`, `uniseg`, `displaywidth` and `uax29/v2` have none in osv.dev today.

## Open risks

- **The fixtures are mine.** 208 cells chosen to cover the classes in the brief, not a sample of real ADR text. Real titles are probably almost all ASCII with the occasional accented letter, dash, quotation mark or emoji, where every Unicode-based candidate is right. A corpus of real titles and summaries would give a better estimate; I had none, and I did not go looking in the daily logs or journals for one.
- **The table baseline is fitted.** I wrote and fixed it against the same cells; the libraries were not tuned. A wrapper around `displaywidth` that applies Prettier's zero-width rules could close part of the gap, but I did not build or test one.
- **Prettier can change its rule.** Everything above is for 3.9.10. The width function, `get-east-asian-width` and `emoji-regex` are separate packages that move on their own, and 4.0.0-alpha.13 is out. The recommended wrapper test must run on each pin bump (Y17).
- **`displaywidth` is young and pre-1.0.** Eleven minor versions in four months, one maintainer, 94 importing packages. The API could change; its Unicode version is not documented where I looked. Mitigation: pin the version, call it from one function, and keep the 110-line table baseline as the known replacement.
- **Cell text must be normalised first.** Prettier turns runs of spaces and tabs into one space. A title with a double space would give a different cell, not only a different width. Likewise `|` must be escaped as `\|` before measuring (the spec says so; the `\|` fixture shows Prettier counts the backslash).
- **Invalid UTF-8 and control characters** should be rejected or stripped by the tool before measuring; the candidates do not agree on them.
- **`go-runewidth` is not safe for generated files** unless the program sets `runewidth.DefaultCondition.EastAsianWidth = false` at start-up, which I did not test.
- **One machine, one Go version.** `go1.27.2 darwin/arm64`. `x/text/width` gives different tables on older Go.

## What I could not verify

- Prettier 4: not tested; only 3.9.10.
- Whether the Unicode version behind `displaywidth` is 17; I saw only that it measures Unicode 16 and 17 emoji correctly, and that its `uax29` dependency says Unicode 17.
- The content of the open issues about width in `x/text` and `displaywidth` (not read); maintainers' merge rights beyond commit counts.
- Whether osv.dev and the GitHub advisory database list every advisory; private advisories are invisible.
- Windows and Linux behaviour. `go-runewidth` has separate Windows code that I did not run.
- A fuzz test of any candidate; only fixtures and the hostile cases here.
- How `akmctl` will read the file: if it never receives invalid UTF-8, those 8 cells are moot.
- The sources of the libraries beyond the file searches described; I did not run any candidate's own tests, scripts or `go generate`.

## Appendix A: how it was run

Scratch directory outside the vault. Only my own programs were run, importing the candidates as modules; no repository test suite, script, `Makefile` target or `go generate` was run. The module was created with `GOFLAGS=-mod=mod go get <module>@<version>` once, then `go mod tidy`, and built with `GOFLAGS=-mod=readonly go build -o widtheval .`. Prettier: `npx --yes prettier@3.9.10 --no-config --no-editorconfig --write <dir>`. Unicode data for the table baseline: `https://www.unicode.org/Public/17.0.0/ucd/EastAsianWidth.txt` and `.../ucd/emoji/emoji-data.txt` (data files, not code). Commands:

```bash
python3 gen_fixtures.py fixtures.json             # 208 fixture cells (hex of the exact bytes)
python3 prettier_widths.py fixtures.json ./p      # Prettier 3.9.10 -> truth.json
python3 gen_table.py EastAsianWidth17.txt emoji-data17.txt > table.go
go build -o widtheval . && ./widtheval compare ..  # per-cell widths, per-class counts
./widtheval bench x ; ./widtheval hostile x        # performance; hostile input (child processes)
LC_ALL=ja_JP.UTF-8 ./widtheval compare ..          # locale test for go-runewidth
```

Variants for `go-runewidth` alone (built in separate modules so the dependency versions are the published ones): v0.0.31 with `uax29/v2` v2.2.0: 42 wrong, same as in the main module with v2.7.0; v0.0.16 (depends on `uniseg` v0.2.0): 53 wrong.

`go.mod` of the evaluation module:

```
module widtheval

go 1.27

require (
	github.com/clipperhouse/displaywidth v0.11.0
	github.com/mattn/go-runewidth v0.0.31
	github.com/rivo/uniseg v0.4.7
	golang.org/x/text v0.43.0
)

require github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
```

`go.sum`:

```
github.com/clipperhouse/displaywidth v0.11.0 h1:lBc6kY44VFw+TDx4I8opi/EtL9m20WSEFgwIwO+UVM8=
github.com/clipperhouse/displaywidth v0.11.0/go.mod h1:bkrFNkf81G8HyVqmKGxsPufD3JhNl3dSqnGhOoSD/o0=
github.com/clipperhouse/uax29/v2 v2.7.0 h1:+gs4oBZ2gPfVrKPthwbMzWZDaAFPGYK72F0NJv2v7Vk=
github.com/clipperhouse/uax29/v2 v2.7.0/go.mod h1:EFJ2TJMRUaplDxHKj1qAEhCtQPW2tJSwu5BF98AuoVM=
github.com/mattn/go-runewidth v0.0.31 h1:1oCV4ReTlosDBc6rYRw5cL8S/4y8iv5oRl4SJ5IBgwo=
github.com/mattn/go-runewidth v0.0.31/go.mod h1:3qAiGCV4Koz/yuveO58qUefmUTRm8r0IGEXZ9jeHp/8=
github.com/rivo/uniseg v0.4.7 h1:WUdvkW8uEhrYfLC4ZzdpI2ztxP1I582+49Oc5Mq64VQ=
github.com/rivo/uniseg v0.4.7/go.mod h1:FN3SvrM+Zdj16jyLfmOkMNblXMcoc8DfTHruCPUcx88=
golang.org/x/text v0.43.0 h1:1QivrlhwAsnMOcqgOdl9my7s9OyZl8O5/UGUTJVXoWc=
golang.org/x/text v0.43.0/go.mod h1:hbSIYA/amXcRXTFSZqh8tWHEoqRvvbKPJ2uz/KU+a+s=
```

## Appendix B: Prettier's width function (from `prettier@3.9.10/index.mjs`, lines 15728 to 15755)

The long `narrowEmojiRegexp` literal on line 15724 (about 1.3 KB, a list of text-style emoji code points) is omitted here; it is the package `narrow-emojis`.

```js
var notAsciiRegex = /[^\x20-\x7F]/;
function getStringWidth(text) {
  if (!text) {
    return 0;
  }
  if (!notAsciiRegex.test(text)) {
    return text.length;
  }
  let width = 0;
  text = text.replace(emoji_regex_default(), (character) => {
    width += isNarrowEmojiCharacter(character) ? 1 : 2;
    return "";
  });
  for (const character of text) {
    const codePoint = character.codePointAt(0);
    if (codePoint <= 31 || codePoint >= 127 && codePoint <= 159) {
      continue;
    }
    if (codePoint >= 768 && codePoint <= 879) {
      continue;
    }
    if (codePoint >= 65024 && codePoint <= 65039) {
      continue;
    }
    width += isFullWidth(codePoint) || isWide(codePoint) ? 2 : 1;
  }
  return width;
}
```

## Appendix C: fixtures and the Prettier ground truth

`gen_fixtures.py` is the fixture list: 208 cells in 21 classes, written as Python string escapes, stored as hex so invalid UTF-8 is exact. `prettier_widths.py` builds one table per cell, runs Prettier, and derives the width from the padding: the separator row has as many dashes as the column is wide, and the cell row is the cell text followed by the padding.

### `gen_fixtures.py`

```python
#!/usr/bin/env python3
"""Writes fixtures.json: [{id, cls, hex}] where hex is the exact bytes of one table cell."""
import json, sys
F = []
def add(cls, *cells, raw=False):
    for c in cells:
        b = c if isinstance(c, bytes) else c.encode('utf-8')
        F.append({"cls": cls, "hex": b.hex()})

add("ascii", "a", "ab", "abc", "hello world", "ADR-0001 Use Go for the tool", "0123456789", "2026-10-28", "pad (x) = y; z: 1, 2.", r"a\|b", "", "x   y")
add("latin-precomposed", "café", "naïve", "Zürich", "Ångström", "señor", "Œuvre", "ß", "Dvořák", "ǅ", "Ǆ")
add("latin-combining", "café", "ñ", "Å", "à́̂", "ệ", "́", "ö̈", "abćdef")
add("greek-cyrillic", "αβγδ", "Ελληνικά", "ΑΒΓ", "ά", "Привет мир", "Ё", "ЂЉ")
add("cjk", "漢字", "日本語のテキスト", "中文", "\U00020000", "〇", "一二三四五六七八九十", "汉语拼音 pinyin", "　a", "a　b")
add("kana", "ひらがな", "カタカナ", "ｶﾀｶﾅ", "ー", "ゔ", "アイウエオ")
add("hangul", "한글", "안녕하세요", "가", "각", "ㄱㅏ", "한ᅡ")
add("fullwidth", "ＡＢＣ", "１２３", "！", "＠", "（括弧）", "￥", "｡｢｣", "￩￪")
add("ambiguous", "±", "×", "÷", "°", "§¶", "‘’“”", "…", "→←↔", "█", "■", "○", "★☆", "①②③", "Ⅳ", "─│┌┐", "ø", "Å", "℃", "™", "¡¿")
add("emoji-single", "\U0001f600", "\U0001f680", "✅", "❌", "✔", "❤", "☺", "⭐", "⚠", "©", "\U0001f525", "\U0001f44d", "⌚", "⏰", "☕", "♻", "▶", "\U0001f170", "ℹ", "Ⓜ", "\U0001f201", "㊗", "‼", "↔", "\U0001f4a1 idea", "ok ✅ done")
add("emoji-vs16", "✔️", "❤️", "☺️", "⚠️", "©️", "™️", "♻️", "▶️", "\U0001f170️", "ℹ️", "Ⓜ️", "㊗️", "‼️", "↔️", "✅️", "a️", "漢️", "漢︀", "✔︎")
add("emoji-keycap", "1️⃣", "#️⃣", "0️⃣", "1⃣", "a⃣")
add("emoji-skin", "\U0001f44d\U0001f3fd", "\U0001f44d\U0001f3fb\U0001f44d\U0001f3ff", "\U0001f3fd", "☝\U0001f3fd", "☝️\U0001f3fd", "\U0001f9d1\U0001f3fe")
add("emoji-zwj", "\U0001f9d1\U0001f3fd‍\U0001f4bb", "\U0001f468‍\U0001f469‍\U0001f467‍\U0001f466", "\U0001f3f3️‍\U0001f308", "\U0001f3f4‍☠️", "❤️‍\U0001f525", "\U0001f9d1‍\U0001f91d‍\U0001f9d1", "\U0001f469‍❤️‍\U0001f48b‍\U0001f468", "\U0001f600‍\U0001f600")
add("emoji-flag", "\U0001f1fa\U0001f1f8", "\U0001f1ef\U0001f1f5\U0001f1e9\U0001f1ea", "\U0001f1fa", "\U0001f1fa\U0001f1f8\U0001f1ef", "\U0001f3f4\U000e0067\U000e0062\U000e0065\U000e006e\U000e0067\U000e007f", "\U0001f3c1")
add("zero-width", "a‍b", "a​b", "a‌b", "a⁠b", "a﻿b", "a­b", "‍", "a͏b", "a᠎b", "a‎b", "a‪b‬")
add("other-scripts", "กำ", "ก็", "สวัสดี", "مرحبا", "مُرْ", "שלום", "नमस्ते", "؀a", "ཀཱ", "ក្មែរ")
add("control", "a\tb", "a\x01b", "a\x7fb", "a\x85b", "a\x80b", "a\x1bb", "a\x00b", "a b", "a b", "a b", "a b", "a\u0008b")
add("long", "a"*200, "a"*5000, "漢"*1000, "\U0001f600"*500, ("ab漢\U0001f600é ")*60 + "end", "x"*100000)
add("new-unicode", "\U0001fae9", "\U0001fa89", "\U0001fac8", "\U0002ebf0", "\U000323b0", "\U0001fac6")
add("invalid-utf8", b"a\xffb", b"\xc3", b"\xe6\xbc", b"a\xc0\xafb", b"\xed\xa0\x80", b"\xf4\x90\x80\x80", b"\xf0\x9f\x98", b"ok\x80\x80\x80")
for i, f in enumerate(F): f["id"] = i
json.dump(F, open(sys.argv[1] if len(sys.argv) > 1 else "fixtures.json", "w"))
print(len(F), "fixtures")
```

### `prettier_widths.py`

```python
#!/usr/bin/env python3
"""Usage: prettier_widths.py fixtures.json workdir
Writes one-column tables, runs pinned Prettier 3.9.10 over them, and derives the width Prettier chose for each cell.
Output: truth.json [{id, ok, w, out_text_hex, changed}]"""
import json, os, subprocess, sys, shutil
fx = json.load(open(sys.argv[1])); wd = sys.argv[2]
src, out = os.path.join(wd, 'src'), os.path.join(wd, 'out')
for d in (src, out): shutil.rmtree(d, ignore_errors=True); os.makedirs(d)
for f in fx:
    cell = bytes.fromhex(f['hex'])
    body = b"| h |\n| --- |\n| " + cell + b" |\n| ref |\n"
    open(f"{src}/{f['id']:04d}.md", 'wb').write(body)
    open(f"{out}/{f['id']:04d}.md", 'wb').write(body)
r = subprocess.run(["npx", "--yes", "prettier@3.9.10", "--no-config", "--no-editorconfig", "--write", out], capture_output=True, text=True)
print(r.stdout[-300:], r.stderr[-1000:], "exit", r.returncode)
res = []
for f in fx:
    cell = bytes.fromhex(f['hex'])
    o = open(f"{out}/{f['id']:04d}.md", 'rb').read().decode('utf-8', 'surrogatepass')
    lines = o.split("\n")
    try:
        assert lines[0].startswith("| ") and lines[1].startswith("| -")
        W = lines[1].count('-')
        row = lines[2]
        assert row.startswith("| ") and row.endswith(" |")
        inner = row[2:-2]
        txt = inner.rstrip(' ')
        # a cell that is empty or all spaces: pad is whole inner
        pad = len(inner) - len(txt)
        w = W - pad
        want = cell.decode('utf-8', 'replace')
        res.append({"id": f['id'], "ok": True, "w": w, "colw": W, "out_hex": txt.encode().hex(), "changed": txt != want})
    except Exception as e:
        res.append({"id": f['id'], "ok": False, "err": repr(e), "out": o[:200]})
json.dump(res, open(os.path.join(wd, '..', 'truth.json') if False else 'truth.json', 'w'))
print(sum(r['ok'] for r in res), "ok of", len(res), "; changed text:", sum(r.get('changed', False) for r in res))
```

## Appendix D: Go source

One file per candidate library (`cand_*.go`), so that the `go-runewidth` variants could be built in their own modules. `table.go` is generated (360 ranges) and not reproduced; `gen_table.py` makes it.

### `gen_table.py`

```python
#!/usr/bin/env python3
"""gen_table.py EastAsianWidth.txt emoji-data.txt > table.go : ranges for the embedded-table baseline (Unicode 17.0.0)."""
import re, sys
def ranges(path, want):
    out = {}
    for line in open(path, encoding='utf-8'):
        line = line.split('#')[0].strip()
        if not line: continue
        rng, prop = [x.strip() for x in line.split(';')]
        a, _, b = rng.partition('..'); b = int(b or a, 16); a = int(a, 16)
        if prop in want: out.setdefault(prop, []).append((a, b))
    return out
def merge(rs):
    rs.sort(); m = []
    for a, b in rs:
        if m and a <= m[-1][1] + 1: m[-1] = (m[-1][0], max(m[-1][1], b))
        else: m.append((a, b))
    return m
eaw = ranges(sys.argv[1], {'W', 'F'}); emo = ranges(sys.argv[2], {'Extended_Pictographic', 'Emoji_Presentation', 'Emoji_Modifier'})
def emit(name, rs):
    print(f"var {name} = [][2]rune{{")
    for a, b in merge(rs): print(f"\t{{0x{a:X}, 0x{b:X}}},")
    print("}\n")
print("// Code generated by gen_table.py from Unicode 17.0.0 EastAsianWidth.txt and emoji-data.txt. DO NOT EDIT.\npackage main\n")
emit("wideRanges", eaw['W'] + eaw['F'])
emit("pictoRanges", emo['Extended_Pictographic'])
emit("emojiPresRanges", emo['Emoji_Presentation'])
```

### `main.go`

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
)

type Fix struct {
	ID  int    `json:"id"`
	Cls string `json:"cls"`
	Hex string `json:"hex"`
}
type Truth struct {
	ID      int  `json:"id"`
	OK      bool `json:"ok"`
	W       int  `json:"w"`
	Changed bool `json:"changed"`
}

func load(dir string) (fx []Fix, tr map[int]Truth) {
	b, _ := os.ReadFile(dir + "/fixtures.json")
	json.Unmarshal(b, &fx)
	b, _ = os.ReadFile(dir + "/truth.json")
	var ts []Truth
	json.Unmarshal(b, &ts)
	tr = map[int]Truth{}
	for _, t := range ts {
		tr[t.ID] = t
	}
	return
}

func unhex(h string) string {
	b := make([]byte, len(h)/2)
	for i := range b {
		fmt.Sscanf(h[2*i:2*i+2], "%02x", &b[i])
	}
	return string(b)
}

func safe(f func(string) int, s string) (n int, p string) {
	defer func() {
		if r := recover(); r != nil {
			p = fmt.Sprint("PANIC: ", r)
		}
	}()
	return f(s), ""
}

func short(s string) string {
	q := fmt.Sprintf("%+q", s)
	if len(q) > 44 {
		q = q[:41] + "..."
	}
	return q
}

func compare(dir string) {
	fx, tr := load(dir)
	fmt.Printf("go %s %s/%s; fixtures %d\n", runtime.Version(), runtime.GOOS, runtime.GOARCH, len(fx))
	fmt.Printf("%-4s %-18s %-46s %4s", "id", "class", "cell", "PRET")
	for _, c := range cands {
		fmt.Printf(" %5.5s", c.Name)
	}
	fmt.Println()
	wrong := map[string]int{}
	byCls := map[string]map[string]int{}
	clsN := map[string]int{}
	scored := 0
	for _, f := range fx {
		t := tr[f.ID]
		s := unhex(f.Hex)
		skip := t.Changed
		if !skip {
			scored++
			clsN[f.Cls]++
		}
		line := fmt.Sprintf("%-4d %-18s %-46s %4d", f.ID, f.Cls, short(s), t.W)
		if skip {
			line += "  (Prettier rewrote the cell text; not scored)"
			fmt.Println(line)
			continue
		}
		bad := false
		for _, c := range cands {
			n, p := safe(c.Fn, s)
			mark := " "
			if p != "" || n != t.W {
				mark = "*"
				wrong[c.Name]++
				if byCls[c.Name] == nil {
					byCls[c.Name] = map[string]int{}
				}
				byCls[c.Name][f.Cls]++
				bad = true
			}
			if p != "" {
				line += fmt.Sprintf(" %5s", "PANIC")
			} else {
				line += fmt.Sprintf(" %4d%s", n, mark)
			}
		}
		if bad {
			fmt.Println(line)
		} else {
			fmt.Println(strings.Replace(line, "  ", "  ", 0) + "   (all agree)")
		}
	}
	fmt.Printf("\nScored %d cells. Wrong per candidate (cells where width != Prettier):\n", scored)
	for _, c := range cands {
		fmt.Printf("  %-14s %3d\n", c.Name, wrong[c.Name])
	}
	real := map[string]bool{"ascii": true, "latin-precomposed": true, "latin-combining": true, "greek-cyrillic": true, "cjk": true, "kana": true, "hangul": true, "fullwidth": true, "ambiguous": true, "emoji-single": true, "emoji-vs16": true}
	nReal := 0
	for k, n := range clsN {
		if real[k] {
			nReal += n
		}
	}
	fmt.Printf("\nCommon-text subset (ascii, latin, greek/cyrillic, cjk, kana, hangul, fullwidth, ambiguous, single emoji, emoji+VS16): %d cells. Wrong:\n", nReal)
	for _, c := range cands {
		n := 0
		for k, v := range byCls[c.Name] {
			if real[k] {
				n += v
			}
		}
		fmt.Printf("  %-14s %3d\n", c.Name, n)
	}
	var cls []string
	for k := range clsN {
		cls = append(cls, k)
	}
	sort.Strings(cls)
	fmt.Printf("\nWrong per class (of n cells):\n%-18s %3s", "class", "n")
	for _, c := range cands {
		fmt.Printf(" %5.5s", c.Name)
	}
	fmt.Println()
	for _, k := range cls {
		fmt.Printf("%-18s %3d", k, clsN[k])
		for _, c := range cands {
			fmt.Printf(" %5d", byCls[c.Name][k])
		}
		fmt.Println()
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: widtheval compare|bench|hostile|child <dir>")
		return
	}
	switch os.Args[1] {
	case "compare":
		compare(os.Args[2])
	case "bench":
		bench(os.Args[2])
	case "hostile":
		hostile()
	case "child":
		child(os.Args[2:])
	}
}
```

### `cands.go`

```go
package main

import (
	"sort"
	"unicode/utf8"
)

type Cand struct {
	Name string
	Fn   func(string) int
}

var cands = []Cand{
	{"runecount", utf8.RuneCountInString},
	{"table-eaw", func(s string) int { return tableWidth(s, false) }},
	{"table-emoji", func(s string) int { return tableWidth(s, true) }},
}

// ---- embedded-table baseline (ranges in table.go, generated from Unicode 17.0.0) ----

func inRanges(t [][2]rune, r rune) bool {
	i := sort.Search(len(t), func(i int) bool { return t[i][1] >= r })
	return i < len(t) && t[i][0] <= r
}

func isRI(r rune) bool { return r >= 0x1F1E6 && r <= 0x1F1FF }

// runeW follows the rules Prettier applies to one code point: controls, U+0300-036F and U+FE00-FE0F are 0, wide/fullwidth are 2.
func runeW(r rune) int {
	switch {
	case r <= 0x1F || r >= 0x7F && r <= 0x9F, r >= 0x300 && r <= 0x36F, r >= 0xFE00 && r <= 0xFE0F:
		return 0
	case r < 0x1100:
		return 1
	case inRanges(wideRanges, r):
		return 2
	}
	return 1
}

// emojiCluster recognises a simplified emoji sequence at rs[i]: n runes consumed (0 = none) and its width.
func emojiCluster(rs []rune, i int) (n, w int) {
	r := rs[i]
	if (r >= '0' && r <= '9' || r == '#' || r == '*') && i+1 < len(rs) {
		j := i + 1
		if rs[j] == 0xFE0F {
			j++
		}
		if j < len(rs) && rs[j] == 0x20E3 {
			return j + 1 - i, 2
		}
		return 0, 0
	}
	if isRI(r) {
		if i+1 < len(rs) && isRI(rs[i+1]) {
			return 2, 2
		}
		return 0, 0
	}
	if !inRanges(pictoRanges, r) {
		return 0, 0
	}
	j := i + 1
	for {
		if j < len(rs) && rs[j] == 0xFE0F {
			j++
		}
		if j < len(rs) && rs[j] >= 0x1F3FB && rs[j] <= 0x1F3FF {
			j++
		}
		if j < len(rs) && rs[j] >= 0xE0020 && rs[j] <= 0xE007F { // tag sequence
			for j < len(rs) && rs[j] >= 0xE0020 && rs[j] <= 0xE007F {
				j++
			}
		}
		if j+1 < len(rs) && rs[j] == 0x200D && inRanges(pictoRanges, rs[j+1]) {
			j += 2
			continue
		}
		break
	}
	if j-i == 1 && !inRanges(emojiPresRanges, r) && !inRanges(wideRanges, r) {
		return 1, 1 // a bare text-presentation emoji is narrow
	}
	return j - i, 2
}

func tableWidth(s string, emoji bool) int {
	rs := []rune(s)
	w := 0
	for i := 0; i < len(rs); {
		if emoji {
			if n, cw := emojiCluster(rs, i); n > 0 {
				w += cw
				i += n
				continue
			}
		}
		w += runeW(rs[i])
		i++
	}
	return w
}
```

### `cand_xtext.go`

```go
package main

import "golang.org/x/text/width"

func init() { cands = append(cands, Cand{"x/text/width", xtextWidth}) }

// xtextWidth: per rune, Wide and Fullwidth are 2, everything else 1 (the library has no zero-width notion).
func xtextWidth(s string) int {
	n := 0
	for _, r := range s {
		switch width.LookupRune(r).Kind() {
		case width.EastAsianWide, width.EastAsianFullwidth:
			n += 2
		default:
			n++
		}
	}
	return n
}
```

### `cand_runewidth.go`

```go
package main

import "github.com/mattn/go-runewidth"

func init() { cands = append(cands, Cand{"go-runewidth", runewidth.StringWidth}) }
```

### `cand_uniseg.go`

```go
package main

import "github.com/rivo/uniseg"

func init() { cands = append(cands, Cand{"uniseg", uniseg.StringWidth}) }
```

### `cand_displaywidth.go`

```go
package main

import "github.com/clipperhouse/displaywidth"

func init() { cands = append(cands, Cand{"displaywidth", displaywidth.String}) }
```

### `bench.go`

```go
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

// corpus builds n deterministic table cells. kind "ascii": ADR-like titles; "mixed": 70% ASCII, 10% accented Latin, 10% CJK, 10% emoji.
func corpus(kind string, n int) []string {
	rng := rand.New(rand.NewSource(42))
	words := strings.Fields("use go for the tool choose a database adopt prettier generate index from frontmatter accept records as markdown keep supersede ADR when decision changes validate schema cache build pipeline")
	title := func() string {
		k := 3 + rng.Intn(9)
		var p []string
		for i := 0; i < k; i++ {
			p = append(p, words[rng.Intn(len(words))])
		}
		return strings.Join(p, " ")
	}
	out := make([]string, n)
	for i := range out {
		t := title()
		if kind == "mixed" {
			switch r := rng.Intn(10); {
			case r == 0:
				t = "Café résumé naïve " + t
			case r == 1:
				t = "日本語の決定 " + t
			case r == 2:
				t = t + " ✅ \U0001f680 ❤️"
			}
		}
		out[i] = t
	}
	return out
}

func bench(_ string) {
	for _, kind := range []string{"ascii", "mixed"} {
		cells := corpus(kind, 5000)
		bytes := 0
		for _, c := range cells {
			bytes += len(c)
		}
		fmt.Printf("corpus %s: %d cells, %d bytes (avg %d)\n", kind, len(cells), bytes, bytes/len(cells))
		for _, c := range cands {
			var best testing.BenchmarkResult
			for run := 0; run < 5; run++ {
				r := testing.Benchmark(func(b *testing.B) {
					b.ReportAllocs()
					sink := 0
					for i := 0; i < b.N; i++ {
						for _, s := range cells {
							sink += c.Fn(s)
						}
					}
					_ = sink
				})
				if run == 0 || r.NsPerOp() < best.NsPerOp() {
					best = r
				}
			}
			fmt.Printf("  %-14s %9.0f ns per 5000 cells  %7.1f ns/cell  %6.2f allocs/cell  %7.1f B/cell\n", c.Name,
				float64(best.NsPerOp()), float64(best.NsPerOp())/5000, float64(best.AllocsPerOp())/5000, float64(best.AllocedBytesPerOp())/5000)
		}
	}
}
```

### `hostile.go`

```go
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var hnames = []string{
	"ascii-100MB", "cjk-100MB", "invalid-100MB", "combining-10M", "zwj-chain-5M", "regional-5M", "skin-vs-5M", "tags-5M", "crlf-ascii-1MBx", "emoji-100MB",
}

func hcase(name string) string {
	switch name {
	case "ascii-100MB":
		return strings.Repeat("a", 100<<20)
	case "cjk-100MB":
		return strings.Repeat("漢", (100<<20)/3)
	case "invalid-100MB":
		return strings.Repeat("\xff", 100<<20)
	case "combining-10M": // one base then 10M combining marks: a single huge grapheme cluster
		return "a" + strings.Repeat("́", 10_000_000)
	case "zwj-chain-5M": // one emoji ZWJ chain
		return strings.Repeat("\U0001f600‍", 5_000_000) + "\U0001f600"
	case "regional-5M":
		return strings.Repeat("\U0001f1fa", 5_000_000)
	case "skin-vs-5M":
		return "\U0001f44d" + strings.Repeat("️\U0001f3fd", 2_500_000)
	case "tags-5M":
		return "\U0001f3f4" + strings.Repeat("\U000e0067", 5_000_000)
	case "emoji-100MB":
		return strings.Repeat("\U0001f600", (100<<20)/4)
	case "crlf-ascii-1MBx":
		return strings.Repeat("ab\r\n", 250_000)
	}
	return ""
}

func child(a []string) { // child <cand> <case>
	var fn func(string) int
	for _, c := range cands {
		if c.Name == a[0] {
			fn = c.Fn
		}
	}
	s := hcase(a[1])
	var m0, m1 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m0)
	t := time.Now()
	n, p := safe(fn, s)
	el := time.Since(t)
	runtime.ReadMemStats(&m1)
	fmt.Printf("%d|%d|%d|%s", n, el.Milliseconds(), (m1.TotalAlloc-m0.TotalAlloc)>>20, p)
}

func rss(pid int) int64 {
	out, err := exec.Command("ps", "-o", "rss=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return 0
	}
	v, _ := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
	return v / 1024
}

func hostile() {
	self, _ := os.Executable()
	fmt.Printf("%-18s %-14s %10s %9s %12s %9s  %s\n", "case", "candidate", "width", "ms", "alloc MB", "peak MB", "note")
	for _, h := range hnames {
		for _, c := range cands {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			cmd := exec.CommandContext(ctx, self, "child", c.Name, h)
			var buf strings.Builder
			cmd.Stdout = &buf
			cmd.Stderr = &buf
			cmd.Start()
			done := make(chan error, 1)
			go func() { done <- cmd.Wait() }()
			var peak int64
			note := ""
		loop:
			for {
				select {
				case <-done:
					break loop
				case <-time.After(40 * time.Millisecond):
					if r := rss(cmd.Process.Pid); r > peak {
						peak = r
					}
					if peak > 6000 {
						cmd.Process.Kill()
						note = "KILLED at 6 GB"
					}
				}
			}
			if ctx.Err() != nil {
				note = "KILLED after 30 s"
			}
			cancel()
			f := strings.SplitN(buf.String(), "|", 4)
			if len(f) == 4 && note == "" {
				fmt.Printf("%-18s %-14s %10s %9s %12s %9d  %s\n", h, c.Name, f[0], f[1], f[2], peak, strings.TrimSpace(f[3]))
			} else {
				fmt.Printf("%-18s %-14s %10s %9s %12s %9d  %s %s\n", h, c.Name, "-", "-", "-", peak, note, strings.TrimSpace(buf.String()))
			}
		}
	}
}
```

## Appendix E: raw results

### Per-cell widths (`./widtheval compare`)

Columns: Prettier 3.9.10, then the candidates in the order rune count, table (wide only), table + emoji, `displaywidth`, `go-runewidth`, `uniseg`, `x/text/width`. An asterisk marks a width that differs from Prettier. Cells are shown as Go `%+q` escapes, cut at 44 characters. "all agree" marks rows where every candidate matches Prettier.

```text
go go1.27.2 darwin/arm64; fixtures 208
id   class              cell                                           PRET runec table table displ go-ru unise x/tex
0    ascii              "a"                                               1    1     1     1     1     1     1     1    (all agree)
1    ascii              "ab"                                              2    2     2     2     2     2     2     2    (all agree)
2    ascii              "abc"                                             3    3     3     3     3     3     3     3    (all agree)
3    ascii              "hello world"                                    11   11    11    11    11    11    11    11    (all agree)
4    ascii              "ADR-0001 Use Go for the tool"                   28   28    28    28    28    28    28    28    (all agree)
5    ascii              "0123456789"                                     10   10    10    10    10    10    10    10    (all agree)
6    ascii              "2026-10-28"                                     10   10    10    10    10    10    10    10    (all agree)
7    ascii              "pad (x) = y; z: 1, 2."                          21   21    21    21    21    21    21    21    (all agree)
8    ascii              "a\\|b"                                           4    4     4     4     4     4     4     4    (all agree)
9    ascii              ""                                                0    0     0     0     0     0     0     0    (all agree)
10   ascii              "x   y"                                           3  (Prettier rewrote the cell text; not scored)
11   latin-precomposed  "caf\u00e9"                                       4    4     4     4     4     4     4     4    (all agree)
12   latin-precomposed  "na\u00efve"                                      5    5     5     5     5     5     5     5    (all agree)
13   latin-precomposed  "Z\u00fcrich"                                     6    6     6     6     6     6     6     6    (all agree)
14   latin-precomposed  "\u00c5ngstr\u00f6m"                              8    8     8     8     8     8     8     8    (all agree)
15   latin-precomposed  "se\u00f1or"                                      5    5     5     5     5     5     5     5    (all agree)
16   latin-precomposed  "\u0152uvre"                                      5    5     5     5     5     5     5     5    (all agree)
17   latin-precomposed  "\u00df"                                          1    1     1     1     1     1     1     1    (all agree)
18   latin-precomposed  "Dvo\u0159\u00e1k"                                6    6     6     6     6     6     6     6    (all agree)
19   latin-precomposed  "\u01c5"                                          1    1     1     1     1     1     1     1    (all agree)
20   latin-precomposed  "\u01c4"                                          1    1     1     1     1     1     1     1    (all agree)
21   latin-combining    "cafe\u0301"                                      4    5*    4     4     4     4     4     5*
22   latin-combining    "n\u0303"                                         1    2*    1     1     1     1     1     2*
23   latin-combining    "A\u030a"                                         1    2*    1     1     1     1     1     2*
24   latin-combining    "a\u0300\u0301\u0302"                             1    4*    1     1     1     1     1     4*
25   latin-combining    "e\u0323\u0302"                                   1    3*    1     1     1     1     1     3*
26   latin-combining    "\u0301"                                          0    1*    0     0     0     0     0     1*
27   latin-combining    "o\u0308\u0308"                                   1    3*    1     1     1     1     1     3*
28   latin-combining    "abc\u0301def"                                    6    7*    6     6     6     6     6     7*
29   greek-cyrillic     "\u03b1\u03b2\u03b3\u03b4"                        4    4     4     4     4     4     4     4    (all agree)
30   greek-cyrillic     "\u0395\u03bb\u03bb\u03b7\u03bd\u03b9\u03...      8    8     8     8     8     8     8     8    (all agree)
31   greek-cyrillic     "\u0391\u0392\u0393"                              3    3     3     3     3     3     3     3    (all agree)
32   greek-cyrillic     "\u03b1\u0301"                                    1    2*    1     1     1     1     1     2*
33   greek-cyrillic     "\u041f\u0440\u0438\u0432\u0435\u0442 \u0...     10   10    10    10    10    10    10    10    (all agree)
34   greek-cyrillic     "\u0401"                                          1    1     1     1     1     1     1     1    (all agree)
35   greek-cyrillic     "\u0402\u0409"                                    2    2     2     2     2     2     2     2    (all agree)
36   cjk                "\u6f22\u5b57"                                    4    2*    4     4     4     4     4     4 
37   cjk                "\u65e5\u672c\u8a9e\u306e\u30c6\u30ad\u30...     16    8*   16    16    16    16    16    16 
38   cjk                "\u4e2d\u6587"                                    4    2*    4     4     4     4     4     4 
39   cjk                "\U00020000"                                      2    1*    2     2     2     2     2     2 
40   cjk                "\u3007"                                          2    1*    2     2     2     2     2     2 
41   cjk                "\u4e00\u4e8c\u4e09\u56db\u4e94\u516d\u4e...     20   10*   20    20    20    20    20    20 
42   cjk                "\u6c49\u8bed\u62fc\u97f3 pinyin"                15   11*   15    15    15    15    15    15 
43   cjk                "\u3000a"                                         3    2*    3     3     3     3     3     3 
44   cjk                "a\u3000b"                                        4    3*    4     4     4     4     4     4 
45   kana               "\u3072\u3089\u304c\u306a"                        8    4*    8     8     8     8     8     8 
46   kana               "\u30ab\u30bf\u30ab\u30ca"                        8    4*    8     8     8     8     8     8 
47   kana               "\uff76\uff80\uff76\uff85"                        4    4     4     4     4     4     4     4    (all agree)
48   kana               "\u30fc"                                          2    1*    2     2     2     2     2     2 
49   kana               "\u3094"                                          2    1*    2     2     2     2     2     2 
50   kana               "\u30a2\u30a4\u30a6\u30a8\u30aa"                 10    5*   10    10    10    10    10    10 
51   hangul             "\ud55c\uae00"                                    4    2*    4     4     4     4     4     4 
52   hangul             "\uc548\ub155\ud558\uc138\uc694"                 10    5*   10    10    10    10    10    10 
53   hangul             "\u1100\u1161"                                    3    2*    3     3     2*    2*    2*    3 
54   hangul             "\u1100\u1161\u11a8"                              4    3*    4     4     2*    2*    2*    4 
55   hangul             "\u3131\u314f"                                    4    2*    4     4     4     4     4     4 
56   hangul             "\ud55c\u1161"                                    3    2*    3     3     3     3     3     3 
57   fullwidth          "\uff21\uff22\uff23"                              6    3*    6     6     6     6     6     6 
58   fullwidth          "\uff11\uff12\uff13"                              6    3*    6     6     6     6     6     6 
59   fullwidth          "\uff01"                                          2    1*    2     2     2     2     2     2 
60   fullwidth          "\uff20"                                          2    1*    2     2     2     2     2     2 
61   fullwidth          "\uff08\u62ec\u5f27\uff09"                        8    4*    8     8     8     8     8     8 
62   fullwidth          "\uffe5"                                          2    1*    2     2     2     2     2     2 
63   fullwidth          "\uff61\uff62\uff63"                              3    3     3     3     3     3     3     3    (all agree)
64   fullwidth          "\uffe9\uffea"                                    2    2     2     2     2     2     2     2    (all agree)
65   ambiguous          "\u00b1"                                          1    1     1     1     1     1     1     1    (all agree)
66   ambiguous          "\u00d7"                                          1    1     1     1     1     1     1     1    (all agree)
67   ambiguous          "\u00f7"                                          1    1     1     1     1     1     1     1    (all agree)
68   ambiguous          "\u00b0"                                          1    1     1     1     1     1     1     1    (all agree)
69   ambiguous          "\u00a7\u00b6"                                    2    2     2     2     2     2     2     2    (all agree)
70   ambiguous          "\u2018\u2019\u201c\u201d"                        4    4     4     4     4     4     4     4    (all agree)
71   ambiguous          "\u2026"                                          1    1     1     1     1     1     1     1    (all agree)
72   ambiguous          "\u2192\u2190\u2194"                              3    3     3     3     3     3     3     3    (all agree)
73   ambiguous          "\u2588"                                          1    1     1     1     1     1     1     1    (all agree)
74   ambiguous          "\u25a0"                                          1    1     1     1     1     1     1     1    (all agree)
75   ambiguous          "\u25cb"                                          1    1     1     1     1     1     1     1    (all agree)
76   ambiguous          "\u2605\u2606"                                    2    2     2     2     2     2     2     2    (all agree)
77   ambiguous          "\u2460\u2461\u2462"                              3    3     3     3     3     3     3     3    (all agree)
78   ambiguous          "\u2163"                                          1    1     1     1     1     1     1     1    (all agree)
79   ambiguous          "\u2500\u2502\u250c\u2510"                        4    4     4     4     4     4     4     4    (all agree)
80   ambiguous          "\u00f8"                                          1    1     1     1     1     1     1     1    (all agree)
81   ambiguous          "\u00c5"                                          1    1     1     1     1     1     1     1    (all agree)
82   ambiguous          "\u2103"                                          1    1     1     1     1     1     1     1    (all agree)
83   ambiguous          "\u2122"                                          1    1     1     1     1     1     1     1    (all agree)
84   ambiguous          "\u00a1\u00bf"                                    2    2     2     2     2     2     2     2    (all agree)
85   emoji-single       "\U0001f600"                                      2    1*    2     2     2     2     2     2 
86   emoji-single       "\U0001f680"                                      2    1*    2     2     2     2     2     2 
87   emoji-single       "\u2705"                                          2    1*    2     2     2     2     2     2 
88   emoji-single       "\u274c"                                          2    1*    2     2     2     2     2     2 
89   emoji-single       "\u2714"                                          1    1     1     1     1     1     1     1    (all agree)
90   emoji-single       "\u2764"                                          1    1     1     1     1     1     1     1    (all agree)
91   emoji-single       "\u263a"                                          1    1     1     1     1     1     1     1    (all agree)
92   emoji-single       "\u2b50"                                          2    1*    2     2     2     2     2     2 
93   emoji-single       "\u26a0"                                          1    1     1     1     1     1     1     1    (all agree)
94   emoji-single       "\u00a9"                                          1    1     1     1     1     1     1     1    (all agree)
95   emoji-single       "\U0001f525"                                      2    1*    2     2     2     2     2     2 
96   emoji-single       "\U0001f44d"                                      2    1*    2     2     2     2     2     2 
97   emoji-single       "\u231a"                                          2    1*    2     2     2     2     2     2 
98   emoji-single       "\u23f0"                                          2    1*    2     2     2     2     2     2 
99   emoji-single       "\u2615"                                          2    1*    2     2     2     2     2     2 
100  emoji-single       "\u267b"                                          1    1     1     1     1     1     1     1    (all agree)
101  emoji-single       "\u25b6"                                          1    1     1     1     1     1     1     1    (all agree)
102  emoji-single       "\U0001f170"                                      1    1     1     1     1     1     1     1    (all agree)
103  emoji-single       "\u2139"                                          1    1     1     1     1     1     1     1    (all agree)
104  emoji-single       "\u24c2"                                          1    1     1     1     1     1     1     1    (all agree)
105  emoji-single       "\U0001f201"                                      2    1*    2     2     2     2     2     2 
106  emoji-single       "\u3297"                                          2    1*    2     2     2     2     1*    2 
107  emoji-single       "\u203c"                                          1    1     1     1     1     1     1     1    (all agree)
108  emoji-single       "\u2194"                                          1    1     1     1     1     1     1     1    (all agree)
109  emoji-single       "\U0001f4a1 idea"                                 7    6*    7     7     7     7     7     7 
110  emoji-single       "ok \u2705 done"                                 10    9*   10    10    10    10    10    10 
111  emoji-vs16         "\u2714\ufe0f"                                    2    2     1*    2     2     1*    2     2 
112  emoji-vs16         "\u2764\ufe0f"                                    2    2     1*    2     2     1*    2     2 
113  emoji-vs16         "\u263a\ufe0f"                                    2    2     1*    2     2     1*    2     2 
114  emoji-vs16         "\u26a0\ufe0f"                                    2    2     1*    2     2     1*    2     2 
115  emoji-vs16         "\u00a9\ufe0f"                                    2    2     1*    2     2     1*    2     2 
116  emoji-vs16         "\u2122\ufe0f"                                    2    2     1*    2     2     1*    2     2 
117  emoji-vs16         "\u267b\ufe0f"                                    2    2     1*    2     2     1*    2     2 
118  emoji-vs16         "\u25b6\ufe0f"                                    2    2     1*    2     2     1*    2     2 
119  emoji-vs16         "\U0001f170\ufe0f"                                2    2     1*    2     2     1*    2     2 
120  emoji-vs16         "\u2139\ufe0f"                                    2    2     1*    2     2     1*    2     2 
121  emoji-vs16         "\u24c2\ufe0f"                                    2    2     1*    2     2     1*    2     2 
122  emoji-vs16         "\u3297\ufe0f"                                    2    2     2     2     2     2     2     3*
123  emoji-vs16         "\u203c\ufe0f"                                    2    2     1*    2     2     1*    2     2 
124  emoji-vs16         "\u2194\ufe0f"                                    2    2     1*    2     2     1*    2     2 
125  emoji-vs16         "\u2705\ufe0f"                                    2    2     2     2     2     2     2     3*
126  emoji-vs16         "a\ufe0f"                                         1    2*    1     1     2*    1     1     2*
127  emoji-vs16         "\u6f22\ufe0f"                                    2    2     2     2     2     2     2     3*
128  emoji-vs16         "\u6f22\ufe00"                                    2    2     2     2     2     2     2     3*
129  emoji-vs16         "\u2714\ufe0e"                                    1    2*    1     1     1     1     1     2*
130  emoji-keycap       "1\ufe0f\u20e3"                                   2    3*    2     2     2     1*    1*    3*
131  emoji-keycap       "#\ufe0f\u20e3"                                   2    3*    2     2     2     1*    1*    3*
132  emoji-keycap       "0\ufe0f\u20e3"                                   2    3*    2     2     2     1*    1*    3*
133  emoji-keycap       "1\u20e3"                                         2    2     2     2     1*    1*    1*    2 
134  emoji-keycap       "a\u20e3"                                         2    2     2     2     1*    1*    1*    2 
135  emoji-skin         "\U0001f44d\U0001f3fd"                            2    2     4*    2     2     2     2     4*
136  emoji-skin         "\U0001f44d\U0001f3fb\U0001f44d\U0001f3ff"        4    4     8*    4     4     4     4     8*
137  emoji-skin         "\U0001f3fd"                                      2    1*    2     2     2     2     0*    2 
138  emoji-skin         "\u261d\U0001f3fd"                                2    2     3*    2     1*    2     1*    3*
139  emoji-skin         "\u261d\ufe0f\U0001f3fd"                          4    3*    3*    2*    2*    2*    2*    4 
140  emoji-skin         "\U0001f9d1\U0001f3fe"                            2    2     4*    2     2     2     2     4*
141  emoji-zwj          "\U0001f9d1\U0001f3fd\u200d\U0001f4bb"            2    4*    7*    2     2     2     2     7*
142  emoji-zwj          "\U0001f468\u200d\U0001f469\u200d\U0001f4...      2    7*   11*    2     2     2     2    11*
143  emoji-zwj          "\U0001f3f3\ufe0f\u200d\U0001f308"                2    4*    4*    2     2     2     2     5*
144  emoji-zwj          "\U0001f3f4\u200d\u2620\ufe0f"                    2    4*    4*    2     2     2     2     5*
145  emoji-zwj          "\u2764\ufe0f\u200d\U0001f525"                    2    4*    4*    2     2     2     2     5*
146  emoji-zwj          "\U0001f9d1\u200d\U0001f91d\u200d\U0001f9d1"      2    5*    8*    2     2     2     2     8*
147  emoji-zwj          "\U0001f469\u200d\u2764\ufe0f\u200d\U0001...      2    8*   10*    2     2     2     2    11*
148  emoji-zwj          "\U0001f600\u200d\U0001f600"                      5    3*    5     2*    2*    2*    2*    5 
149  emoji-flag         "\U0001f1fa\U0001f1f8"                            2    2     2     2     2     2     2     2    (all agree)
150  emoji-flag         "\U0001f1ef\U0001f1f5\U0001f1e9\U0001f1ea"        4    4     4     4     4     4     4     4    (all agree)
151  emoji-flag         "\U0001f1fa"                                      1    1     1     1     2*    1     2*    1 
152  emoji-flag         "\U0001f1fa\U0001f1f8\U0001f1ef"                  3    3     3     3     4*    3     4*    3 
153  emoji-flag         "\U0001f3f4\U000e0067\U000e0062\U000e0065...      2    7*    8*    2     2     2     2     8*
154  emoji-flag         "\U0001f3c1"                                      2    1*    2     2     2     2     2     2 
155  zero-width         "a\u200db"                                        3    3     3     3     2*    2*    2*    3 
156  zero-width         "a\u200bb"                                        3    3     3     3     2*    2*    2*    3 
157  zero-width         "a\u200cb"                                        3    3     3     3     2*    2*    2*    3 
158  zero-width         "a\u2060b"                                        3    3     3     3     2*    2*    2*    3 
159  zero-width         "a\ufeffb"                                        3    3     3     3     2*    2*    2*    3 
160  zero-width         "a\u00adb"                                        3    3     3     3     2*    2*    2*    3 
161  zero-width         "\u200d"                                          1    1     1     1     0*    0*    0*    1 
162  zero-width         "a\u034fb"                                        2    3*    2     2     2     2     2     3*
163  zero-width         "a\u180eb"                                        3    3     3     3     2*    2*    2*    3 
164  zero-width         "a\u200eb"                                        3    3     3     3     2*    2*    2*    3 
165  zero-width         "a\u202ab\u202c"                                  4    4     4     4     2*    2*    2*    4 
166  other-scripts      "\u0e01\u0e33"                                    2    2     2     2     1*    2     2     2 
167  other-scripts      "\u0e01\u0e47"                                    2    2     2     2     1*    1*    1*    2 
168  other-scripts      "\u0e2a\u0e27\u0e31\u0e2a\u0e14\u0e35"            6    6     6     6     4*    4*    4*    6 
169  other-scripts      "\u0645\u0631\u062d\u0628\u0627"                  5    5     5     5     5     5     5     5    (all agree)
170  other-scripts      "\u0645\u064f\u0631\u0652"                        4    4     4     4     2*    2*    2*    4 
171  other-scripts      "\u05e9\u05dc\u05d5\u05dd"                        4    4     4     4     4     4     4     4    (all agree)
172  other-scripts      "\u0928\u092e\u0938\u094d\u0924\u0947"            6    6     6     6     3*    4*    4*    6 
173  other-scripts      "\u0600a"                                         2    2     2     2     0*    2     2     2 
174  other-scripts      "\u0f40\u0f71"                                    2    2     2     2     1*    1*    1*    2 
175  other-scripts      "\u1780\u17d2\u1798\u17c2\u179a"                  5    5     5     5     2*    4*    4*    5 
176  control            "a\tb"                                            3  (Prettier rewrote the cell text; not scored)
177  control            "a\x01b"                                          2    3*    2     2     2     2     2     3*
178  control            "a\x7fb"                                          3    3     2*    2*    2*    2*    2*    3 
179  control            "a\u0085b"                                        2    3*    2     2     2     2     2     3*
180  control            "a\u0080b"                                        2    3*    2     2     2     2     2     3*
181  control            "a\x1bb"                                          2    3*    2     2     2     2     2     3*
182  control            "a\x00b"                                          2    3*    2     2     2     2     2     3*
183  control            "a\u00a0b"                                        3    3     3     3     3     3     3     3    (all agree)
184  control            "a\u2003b"                                        3    3     3     3     3     3     3     3    (all agree)
185  control            "a\u2007b"                                        3    3     3     3     3     3     3     3    (all agree)
186  control            "a\u2028b"                                        3    3     3     3     2*    2*    2*    3 
187  control            "a\bb"                                            2    3*    2     2     2     2     2     3*
188  long               "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa...    200  200   200   200   200   200   200   200    (all agree)
189  long               "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa...   5000 5000  5000  5000  5000  5000  5000  5000    (all agree)
190  long               "\u6f22\u6f22\u6f22\u6f22\u6f22\u6f22\u6f...   2000 1000* 2000  2000  2000  2000  2000  2000 
191  long               "\U0001f600\U0001f600\U0001f600\U0001f600...   1000  500* 1000  1000  1000  1000  1000  1000 
192  long               "ab\u6f22\U0001f600e\u0301 ab\u6f22\U0001...    483  423*  483   483   483   483   483   543*
193  long               "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx...   100000 100000  100000  100000  100000  100000  100000  100000    (all agree)
194  new-unicode        "\U0001fae9"                                      2    1*    2     2     2     2     1*    2 
195  new-unicode        "\U0001fa89"                                      2    1*    2     2     2     2     1*    2 
196  new-unicode        "\U0001fac8"                                      2    1*    2     2     2     2     1*    2 
197  new-unicode        "\U0002ebf0"                                      2    1*    2     2     2     2     2     2 
198  new-unicode        "\U000323b0"                                      2    1*    2     2     2     2     2     2 
199  new-unicode        "\U0001fac6"                                      2    1*    2     2     2     2     1*    2 
200  invalid-utf8       "a\xffb"                                          3    3     3     3     3     3     3     3    (all agree)
201  invalid-utf8       "\xc3"                                            1    1     1     1     1     1     1     1    (all agree)
202  invalid-utf8       "\xe6\xbc"                                        1    2*    2*    2*    1     2*    2*    2*
203  invalid-utf8       "a\xc0\xafb"                                      4    4     4     4     4     4     4     4    (all agree)
204  invalid-utf8       "\xed\xa0\x80"                                    3    3     3     3     1*    3     3     3 
205  invalid-utf8       "\xf4\x90\x80\x80"                                4    4     4     4     1*    4     4     4 
206  invalid-utf8       "\xf0\x9f\x98"                                    1    3*    3*    3*    1     3*    3*    3*
207  invalid-utf8       "ok\x80\x80\x80"                                  5    5     5     5     5     5     5     5    (all agree)

Scored 206 cells. Wrong per candidate (cells where width != Prettier):
  runecount       84
  table-eaw       29
  table-emoji      5
  displaywidth    32
  go-runewidth    42
  uniseg          38
  x/text/width    40

Common-text subset (ascii, latin, greek/cyrillic, cjk, kana, hangul, fullwidth, ambiguous, single emoji, emoji+VS16): 129 cells. Wrong:
  runecount       51
  table-eaw       13
  table-emoji      0
  displaywidth     3
  go-runewidth    15
  uniseg           3
  x/text/width    15

Wrong per class (of n cells):
class                n runec table table displ go-ru unise x/tex
ambiguous           20     0     0     0     0     0     0     0
ascii               10     0     0     0     0     0     0     0
cjk                  9     9     0     0     0     0     0     0
control             11     6     1     1     2     2     2     6
emoji-flag           6     2     1     0     2     0     2     1
emoji-keycap         5     3     0     0     2     5     5     3
emoji-single        26    14     0     0     0     0     1     0
emoji-skin           6     2     5     1     2     1     3     4
emoji-vs16          19     2    13     0     1    13     0     6
emoji-zwj            8     8     7     1     1     1     1     7
fullwidth            8     6     0     0     0     0     0     0
greek-cyrillic       7     1     0     0     0     0     0     1
hangul               6     6     0     0     2     2     2     0
invalid-utf8         8     2     2     2     2     2     2     2
kana                 6     5     0     0     0     0     0     0
latin-combining      8     8     0     0     0     0     0     8
latin-precomposed   10     0     0     0     0     0     0     0
long                 6     3     0     0     0     0     0     1
new-unicode          6     6     0     0     0     0     4     0
other-scripts       10     0     0     0     8     6     6     0
zero-width          11     1     0     0    10    10    10     1
```

### Benchmark (`./widtheval bench x`)

```text
corpus ascii: 5000 cells, 224946 bytes (avg 44)
  runecount         108964 ns per 5000 cells     21.8 ns/cell    0.00 allocs/cell      0.0 B/cell
  table-eaw         979271 ns per 5000 cells    195.9 ns/cell    0.71 allocs/cell    158.2 B/cell
  table-emoji      2667638 ns per 5000 cells    533.5 ns/cell    0.71 allocs/cell    158.2 B/cell
  displaywidth      133286 ns per 5000 cells     26.7 ns/cell    0.00 allocs/cell      0.0 B/cell
  go-runewidth      188464 ns per 5000 cells     37.7 ns/cell    0.00 allocs/cell      0.0 B/cell
  uniseg           2728669 ns per 5000 cells    545.7 ns/cell    0.00 allocs/cell      0.0 B/cell
  x/text/width      763877 ns per 5000 cells    152.8 ns/cell    0.00 allocs/cell      0.0 B/cell
corpus mixed: 5000 cells, 251776 bytes (avg 50)
  runecount         140677 ns per 5000 cells     28.1 ns/cell    0.00 allocs/cell      0.0 B/cell
  table-eaw        1133399 ns per 5000 cells    226.7 ns/cell    0.75 allocs/cell    173.2 B/cell
  table-emoji      2925389 ns per 5000 cells    585.1 ns/cell    0.75 allocs/cell    173.2 B/cell
  displaywidth      329788 ns per 5000 cells     66.0 ns/cell    0.00 allocs/cell      0.0 B/cell
  go-runewidth      631275 ns per 5000 cells    126.3 ns/cell    0.00 allocs/cell      0.0 B/cell
  uniseg           3156373 ns per 5000 cells    631.3 ns/cell    0.00 allocs/cell      0.0 B/cell
  x/text/width      872197 ns per 5000 cells    174.4 ns/cell    0.00 allocs/cell      0.0 B/cell
```

### Hostile input (`./widtheval hostile x`)

Columns: case, candidate, width returned, milliseconds, MB allocated during the call, peak resident MB of the child process (sampled every 40 ms, includes the input string).

```text
case               candidate           width        ms     alloc MB   peak MB  note
ascii-100MB        runecount       104857600        39            0       105  
ascii-100MB        table-eaw       104857600       425          400       506  
ascii-100MB        table-emoji     104857600      1177          400       506  
ascii-100MB        x/text/width    104857600       340            0       106  
ascii-100MB        go-runewidth    104857600        67            0       105  
ascii-100MB        uniseg          104857600      1256            0       105  
ascii-100MB        displaywidth    104857600        50            0       106  
cjk-100MB          runecount        34952533        85            0       106  
cjk-100MB          table-eaw        69905066       532          133       239  
cjk-100MB          table-emoji      69905066       838          133       239  
cjk-100MB          x/text/width     69905066       255            0       105  
cjk-100MB          go-runewidth     69905066       182            0       108  
cjk-100MB          uniseg           69905066      1595            0       106  
cjk-100MB          displaywidth     69905066       535            0       106  
invalid-100MB      runecount       104857600       238            0       105  
invalid-100MB      table-eaw       104857600      1385          400       506  
invalid-100MB      table-emoji     104857600      2474          400       506  
invalid-100MB      x/text/width    104857600       744            0       106  
invalid-100MB      go-runewidth    104857600      1333            0       108  
invalid-100MB      uniseg          104857600      4890            0       106  
invalid-100MB      displaywidth    104857600      1123            0       106  
combining-10M      runecount        10000001        20            0         0  
combining-10M      table-eaw               1        70           38        68  
combining-10M      table-emoji             1       135           38        82  
combining-10M      x/text/width     10000001        59            0        44  
combining-10M      go-runewidth            1        83            0        46  
combining-10M      uniseg                  1       210            0        44  
combining-10M      displaywidth            1        43            0        44  
zwj-chain-5M       runecount        10000001        26            0         0  
zwj-chain-5M       table-eaw        15000002       146           38       111  
zwj-chain-5M       table-emoji             2       106           38       111  
zwj-chain-5M       x/text/width     15000002       125            0        72  
zwj-chain-5M       go-runewidth            2       104            0        74  
zwj-chain-5M       uniseg                  2       242            0        72  
zwj-chain-5M       displaywidth            2        61            0        72  
regional-5M        runecount         5000000        13            0         0  
regional-5M        table-eaw         5000000        70           19        44  
regional-5M        table-emoji       5000000        35           19         0  
regional-5M        x/text/width      5000000        42            0        24  
regional-5M        go-runewidth      5000000        91            0        27  
regional-5M        uniseg            5000000       127            0        24  
regional-5M        displaywidth      5000000        63            0        25  
skin-vs-5M         runecount         5000001        12            0         0  
skin-vs-5M         table-eaw         5000002        53           19        58  
skin-vs-5M         table-emoji       5000000       104           19        58  
skin-vs-5M         x/text/width      7500002        40            0        39  
skin-vs-5M         go-runewidth            2        50            0        41  
skin-vs-5M         uniseg                  2       102            0        39  
skin-vs-5M         displaywidth            2        26            0         0  
tags-5M            runecount         5000001        13            0         0  
tags-5M            table-eaw         5000002        64           19        44  
tags-5M            table-emoji             2        32           19         0  
tags-5M            x/text/width      5000002        43            0        44  
tags-5M            go-runewidth            2        52            0        46  
tags-5M            uniseg                  2       126            0        44  
tags-5M            displaywidth            2        29            0         0  
crlf-ascii-1MBx    runecount         1000000         0            0         0  
crlf-ascii-1MBx    table-eaw          500000         3            3         0  
crlf-ascii-1MBx    table-emoji        500000        11            3         0  
crlf-ascii-1MBx    x/text/width      1000000         3            0         0  
crlf-ascii-1MBx    go-runewidth       500000         0            0         0  
crlf-ascii-1MBx    uniseg             500000        10            0         0  
crlf-ascii-1MBx    displaywidth       500000         4            0         0  
emoji-100MB        runecount        26214400        76            0       105  
emoji-100MB        table-eaw        52428800       378          100       206  
emoji-100MB        table-emoji      52428800       599          100       206  
emoji-100MB        x/text/width     52428800       225            0       105  
emoji-100MB        go-runewidth     52428800       303            0       108  
emoji-100MB        uniseg           52428800      1097            0       106  
emoji-100MB        displaywidth     52428800       523            0       105  
```
