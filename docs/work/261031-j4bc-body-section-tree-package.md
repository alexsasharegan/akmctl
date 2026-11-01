---
id: "261031-j4bc"
title: body section tree package
status: raw
created: 2026-10-31
updated: 2026-10-31
kind: feature
parent: 261031-gtw6
summary: A package that splits a body into sections, parses the typed core sections and splices edits, per ADR 261031-gam0
---

# body section tree package

## Goal

From spike 261031-qbat and ADR [261031-gam0](../adr/261031-gam0-work-item-body-model.md). Build the in-memory body model the writer (261031-6bzq) and the drawer share: a lossless split into a preamble and sections, typed parsing of Goal (prose plus the first checklist run), Next, Log and Outcome, and splice edits that change only the touched lines. Built on our own scanner behind a `Locate` interface, with goldmark as a test-only oracle, if the proposed ADR [261031-cpbg](../adr/261031-cpbg-goldmark-markdown-parser.md) is accepted. Wait for ADRs 261031-gam0 and 261031-cpbg to be accepted, and for spike 261031-6bzq to settle the operations. Done when:

- [ ] A `Locate(body) (Layout, error)` interface: blocks with kind, byte span and line range; the scanner implements it and returns an error naming an unsupported construct instead of guessing (fail closed, ADR 261031-nvka).
- [ ] A goldmark adapter implements the same interface in test files only; one test battery runs against both, asserting the scanner refuses or equals the oracle on what it claims.
- [ ] The 125 records and the hostile cases run as ordinary tests; a fuzz target seeded with the records and a Markdown-fragment generator runs on demand (`go test -fuzz`). None of the real records may be refused.
- [ ] Fence tracking follows work item 261031-2j14.
- [ ] Split and join are byte-identical on every item in `docs/work/` and `docs/work/done/` (a test over the real files).
- [ ] Ticking a checklist item, appending a log entry and replacing a section each change only their own lines (golden tests, including CRLF).
- [ ] A section that does not parse as its type is kept as opaque text and still edits.
- [ ] New terms are registered in `docs/glossary.md`.

## Log

- 2026-10-31: filed.
- 2026-10-31: blocked_by 261031-g9zk: the parser choice (dependency or our scanner) sets how this is built.
- 2026-10-31: spike 261031-g9zk closed; unblocked, but wait for ADRs 261031-gam0 and 261031-cpbg to be accepted.
