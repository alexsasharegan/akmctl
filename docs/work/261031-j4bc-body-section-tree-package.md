---
id: "261031-j4bc"
title: body section tree package
status: raw
created: 2026-10-31
updated: 2026-10-31
kind: feature
parent: 261031-gtw6
summary: A package that splits a body into sections, parses the typed core sections and splices edits, per ADR 261031-gam0
blocked_by: [261031-yw5v]
---

# body section tree package

## Goal

From spike 261031-qbat and ADR [261031-gam0](../adr/261031-gam0-work-item-body-model.md). On top of `Locate` ([261031-yw5v](261031-yw5v-locate-scanner-and-goldmark-oracle.md)), build the in-memory body model the writer (261031-6bzq) and the drawer share: a lossless split into a preamble and sections, typed parsing of Goal (prose plus the first checklist run), Next, Log and Outcome, and splice edits that change only the touched lines. Built on our own scanner behind a `Locate` interface, with goldmark as a test-only oracle, if the proposed ADR [261031-cpbg](../adr/261031-cpbg-goldmark-markdown-parser.md) is accepted. Wait for ADRs 261031-gam0 and 261031-cpbg to be accepted, and for spike 261031-6bzq to settle the operations. Done when:

- [ ] Split and join are byte-identical on every item in `docs/work/` and `docs/work/done/` (a test over the real files).
- [ ] Ticking a checklist item, appending a log entry and replacing a section each change only their own lines (golden tests, including CRLF).
- [ ] A section that does not parse as its type is kept as opaque text and still edits.
- [ ] New terms are registered in `docs/glossary.md`.

## Log

- 2026-10-31: filed.
- 2026-10-31: blocked_by 261031-g9zk: the parser choice (dependency or our scanner) sets how this is built.
- 2026-10-31: spike 261031-g9zk closed; unblocked, but wait for ADRs 261031-gam0 and 261031-cpbg to be accepted.
- 2026-10-31: the `Locate` interface, scanner, oracle and fuzz moved to 261031-yw5v, built first; ADRs 261031-gam0 and 261031-cpbg are accepted.
