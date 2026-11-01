---
id: "261031-yw5v"
title: locate scanner and goldmark oracle
status: active
created: 2026-10-31
updated: 2026-10-31
kind: feature
parent: 261031-gtw6
summary: "internal/body: a Locate interface, our own fail-closed scanner, and goldmark as a test-only oracle with a differential battery and a fuzz target"
claimed_by: claude
claimed_at: 2026-11-01T05:02
---

# locate scanner and goldmark oracle

## Goal

The first slice of the body package, from ADR [261031-cpbg](../adr/261031-cpbg-goldmark-markdown-parser.md) (accepted) and [261031-gam0](../adr/261031-gam0-work-item-body-model.md). Build `internal/body` with a `Locate(body []byte) (Layout, error)` interface, our own scanner that implements it, and goldmark as a test-only oracle with a differential test battery and a fuzz target, in the style of the YAML-lite parser's `oracle_test.go` and `fuzz_test.go`. The scanner fails closed (ADR [261031-nvka](../adr/261031-nvka-own-narrow-implementations-with-oracles.md)): a construct it does not claim to understand is an error that names the construct and the line, never a guess. Done when:

- [ ] `body.Locate` returns a `Layout` of blocks (kind: heading with level, list item, task item, fence, other; byte start and end; line range), splitting any frontmatter first; the scanner tracks fence character and length.
- [ ] A refusal is `body.ErrUnsupported` with the construct and the line: indented fence or heading, HTML line, blockquote, tab-indented line, table delimiter row, ordered-list task, setext or thematic underline after text.
- [ ] A goldmark adapter implements the same interface in `_test.go` files only; `go.mod` gains goldmark as a test-only dependency (ADR 261031-cpbg); the binary is built without it.
- [ ] One battery runs against both: for any input the scanner refuses or its `Layout` equals the oracle's on the blocks it claims. It covers every item in `docs/work/`, `docs/work/done/` and `docs/adr/` (none may be refused) and a table of hostile cases.
- [ ] A seeded generator over Markdown fragments runs in `go test` (fixed seed, tens of thousands of documents) and reports how many were refused; a `FuzzAgainstOracle` target seeded with the records runs on demand.
- [ ] The package is in `docs/architecture.md`; new terms are in `docs/glossary.md`.

Not in this item: typed Goal, Next and Log parsing and splice edits (261031-j4bc), and moving `SetSection` onto it (261031-2j14).

## Next

Checkpoint. The 300-second fuzz run after the last fix was clean (PASS, no new findings), so step (1) below is done; steps (2) and (3) remain. Built and passing: `internal/body` (body.go, scanner.go; tests body_test.go, oracle_test.go, battery_test.go), goldmark v1.8.7 in `go.mod` (tests only; `go list -deps ./cmd/akmctl` has no goldmark), the docs/architecture.md row and three glossary rows. All 125+ records are located and none refused; the seeded generator (30,000 documents) has 0 confident errors.

Left: (1) a clean fuzz run: `go test ./internal/body -run '^$' -fuzz=FuzzAgainstOracle -fuzztime=300s`. Each failure it finds is saved under `internal/body/testdata/fuzz/` (keep them; they run as ordinary tests); fix it by adding a refusal to scanner.go plus a row in `TestRefusals`, never by loosening the comparison. The last finds were goldmark task-box quirks in ordered and nested lists and table delimiter lookalikes. (2) Then `go vet ./... && go test ./...`, `/tmp/akmctl check`, `npm run format:check`, tick the Done-when boxes, and `close` with an Outcome (include the fuzz finding rate). (3) Then 261031-2j14 (move `SetSection` onto `body`) and 261031-j4bc (typed layer, splice edits).

## Log

- 2026-10-31: filed.
- 2026-10-31: claimed (claude).
- 2026-10-31: made ready and claimed to build it.
- 2026-10-31: built `internal/body` (Locate, scanner, UnsupportedError), goldmark v1.8.7 as a test-only dependency, the oracle adapter, the battery over the records and hostile cases, a seeded generator and `FuzzAgainstOracle`; each fuzz finding became a refusal and a unit test.
- 2026-10-31: a 300-second `FuzzAgainstOracle` run passed with no new findings after the table-delimiter refusal.
