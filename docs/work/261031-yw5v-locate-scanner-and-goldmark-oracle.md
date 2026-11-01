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

Create `internal/body` with the `Layout` types and the scanner; then add goldmark to `go.mod` and the oracle adapter in tests.

## Log

- 2026-10-31: filed.
- 2026-10-31: claimed (claude).
- 2026-10-31: made ready and claimed to build it.
