---
id: "261031-yw5v"
title: locate scanner and goldmark oracle
status: done
created: 2026-10-31
updated: 2026-10-31
closed: 2026-10-31
kind: feature
parent: 261031-gtw6
summary: "internal/body: a Locate interface, our own fail-closed scanner, and goldmark as a test-only oracle with a differential battery and a fuzz target"
---

# locate scanner and goldmark oracle

## Goal

The first slice of the body package, from ADR [261031-cpbg](../../adr/261031-cpbg-goldmark-markdown-parser.md) (accepted) and [261031-gam0](../../adr/261031-gam0-work-item-body-model.md). Build `internal/body` with a `Locate(body []byte) (Layout, error)` interface, our own scanner that implements it, and goldmark as a test-only oracle with a differential test battery and a fuzz target, in the style of the YAML-lite parser's `oracle_test.go` and `fuzz_test.go`. The scanner fails closed (ADR [261031-nvka](../../adr/261031-nvka-own-narrow-implementations-with-oracles.md)): a construct it does not claim to understand is an error that names the construct and the line, never a guess. Done when:

- [x] `body.Locate` returns a `Layout` of blocks (kind: heading with level, list item, task item, fence, other; byte start and end; line range), splitting any frontmatter first; the scanner tracks fence character and length.
- [x] A refusal is `body.ErrUnsupported` with the construct and the line: indented fence or heading, HTML line, blockquote, tab-indented line, table delimiter row, ordered-list task, setext or thematic underline after text.
- [x] A goldmark adapter implements the same interface in `_test.go` files only; `go.mod` gains goldmark as a test-only dependency (ADR 261031-cpbg); the binary is built without it.
- [x] One battery runs against both: for any input the scanner refuses or its `Layout` equals the oracle's on the blocks it claims. It covers every item in `docs/work/`, `docs/work/done/` and `docs/adr/` (none may be refused) and a table of hostile cases.
- [x] A seeded generator over Markdown fragments runs in `go test` (fixed seed, tens of thousands of documents) and reports how many were refused; a `FuzzAgainstOracle` target seeded with the records runs on demand.
- [x] The package is in `docs/architecture.md`; new terms are in `docs/glossary.md`.

Not in this item: typed Goal, Next and Log parsing and splice edits (261031-j4bc), and moving `SetSection` onto it (261031-2j14).

## Next

Checkpoint. The 300-second fuzz run after the last fix was clean (PASS, no new findings), so step (1) below is done; steps (2) and (3) remain. Built and passing: `internal/body` (body.go, scanner.go; tests body_test.go, oracle_test.go, battery_test.go), goldmark v1.8.7 in `go.mod` (tests only; `go list -deps ./cmd/akmctl` has no goldmark), the docs/architecture.md row and three glossary rows. All 125+ records are located and none refused; the seeded generator (30,000 documents) has 0 confident errors.

Left: (1) a clean fuzz run: `go test ./internal/body -run '^$' -fuzz=FuzzAgainstOracle -fuzztime=300s`. Each failure it finds is saved under `internal/body/testdata/fuzz/` (keep them; they run as ordinary tests); fix it by adding a refusal to scanner.go plus a row in `TestRefusals`, never by loosening the comparison. The last finds were goldmark task-box quirks in ordered and nested lists and table delimiter lookalikes. (2) Then `go vet ./... && go test ./...`, `/tmp/akmctl check`, `npm run format:check`, tick the Done-when boxes, and `close` with an Outcome (include the fuzz finding rate). Closing: `akmctl close <id> --status done --outcome "$(cat outcome.md)"` (the status flag is required, and `--outcome -` writes a literal "-", not stdin); closing moves the file to `docs/work/done/`, so fix relative links in other items and rerun `/tmp/akmctl index`. (3) Then 261031-2j14 (move `SetSection` onto `body`) and 261031-j4bc (typed layer, splice edits).

## Outcome

Each "done when" line, with its proof:

- `Locate` and `Layout`: `internal/body/body.go` and `scanner.go`; `TestLocateBlocks`, `TestFenceTracking` (fence character and length) and `TestOffsetsAcrossLineEndings`. Frontmatter is split off by the caller with `internal/frontmatter` before the body is passed in, as ADR 261031-cpbg says ("split off before the body is read"); `Locate` takes the body and its offsets are body-relative. The Goal's words "splitting any frontmatter first" are met by that caller step, not inside `Locate`; `TestRecords` does it with `frontmatter.Parse`.
- Refusals: `ErrUnsupported` and `UnsupportedError` (construct and line); `TestRefusals` covers indented fence and heading, HTML line, blockquote, tab-indented line, table delimiter row, ordered-list task and setext or thematic underline after text, plus more found by fuzzing.
- Goldmark adapter: `oracle_test.go` only; `go.mod` has goldmark v1.8.7; `go list -deps ./cmd/akmctl | grep -c goldmark` prints 0.
- Battery: `TestRecords` locates all 160 files under `docs/work`, `docs/work/done` and `docs/adr` (README excluded) with none refused and every Layout equal to the oracle's; `TestHostileCases` is the table of hostile cases.
- Generator and fuzz: `TestGenerated` (fixed seed) ran 30,000 documents: 16,113 refused, 13,887 compared, 0 wrong. `FuzzAgainstOracle` has 23 saved findings in `internal/body/testdata/fuzz/`, each fixed by a refusal and a `TestRefusals` row; the last 300-second run was clean (no new findings).
- Docs: `docs/architecture.md` has the `internal/body` row; `docs/glossary.md` has Locator and Body layout (and the other rows added with them).

Checks run: `go vet ./...`, `go test -race ./...`, `gofmt -l .` (empty), `/tmp/akmctl check` (0 errors, 0 warnings), `npm run format:check`, `npm run schemas` (78 cases, 0 failures).

Not covered: `FuzzAgainstOracle` was run once for 300 seconds after the last fix, not for longer; the scanner is not yet used by `SetSection` (261031-2j14) or the typed layer (261031-j4bc).

Decision: none (the choices were made in ADRs 261031-cpbg and 261031-gam0).

## Log

- 2026-10-31: filed.
- 2026-10-31: claimed (claude).
- 2026-10-31: made ready and claimed to build it.
- 2026-10-31: built `internal/body` (Locate, scanner, UnsupportedError), goldmark v1.8.7 as a test-only dependency, the oracle adapter, the battery over the records and hostile cases, a seeded generator and `FuzzAgainstOracle`; each fuzz finding became a refusal and a unit test.
- 2026-10-31: a 300-second `FuzzAgainstOracle` run passed with no new findings after the table-delimiter refusal.
- 2026-10-31: verified: vet, race tests, gofmt, records check, format check and schemas pass; closing.
