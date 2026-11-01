---
id: "261031-j4bc"
title: body section tree package
status: active
created: 2026-10-31
updated: 2026-10-31
kind: feature
parent: 261031-gtw6
summary: internal/body gains a lossless section tree, typed Goal/Next/Log/Outcome and splice edits, on top of Locate, per ADR 261031-gam0
claimed_by: claude
claimed_at: 2026-11-01T06:59
---

# body section tree package

## Goal

Build the in-memory body model that the body writer (spike [261031-6bzq](261031-6bzq-spike-body-writer-api.md)) and the drawer will share, in `internal/body`, on top of `Locate` ([261031-yw5v](done/261031-yw5v-locate-scanner-and-goldmark-oracle.md), done) and the shared `FenceTracker` ([261031-2j14](done/261031-2j14-section-fence-tracking.md), done). The model and edit rule are decided: ADR [261031-gam0](../adr/261031-gam0-work-item-body-model.md) (open body, typed core) and ADR [261031-cpbg](../adr/261031-cpbg-goldmark-markdown-parser.md) (own scanner, goldmark as a test-only oracle); both are accepted, so nothing here waits on a decision.

This item is the package, not the CLI or API surface; the operations, addressing syntax and error shapes are 261031-6bzq's. It provides the primitives those will call:

- `Split(body) (Tree, error)` and `Tree.Join() []byte`: a preamble (text before the first `##`) and an ordered list of sections `{name, text}`, the name compared case-insensitively, found with `Locate` so headings in fences are not headings. Lossless: unknown sections and hand formatting are never rewritten.
- Typed views of the four core sections: Goal (prose plus the first contiguous run of `- [ ]` / `- [x]` lines, one line each), Next (prose), Log (append-only `- YYYY-MM-DD: text` entries) and Outcome (prose).
- Splice edits that change only the touched lines and normalize nothing else: tick or untick a checklist item (`[ ]` to `[x]` in place), append a Log entry, replace a section's text. Items are addressed by section name and 1-based position, with a section etag to detect a stale write; no ids are stored in the file.
- A typed section that does not parse (a Goal with prose after the checklist, a Log line with a continuation) keeps the rest as opaque text, and the typed part still edits.
- Line endings: a CRLF file stays CRLF, and an edit writes the line ending of the lines around it.

Decided (Alexander, 2026-10-31): when `Locate` refuses a construct (`ErrUnsupported`), `Split` returns that error and nothing is written (fail closed, ADR 261031-nvka), instead of falling back to the line scan `SetSection` uses. Today `SetSection` and `section` keep their own line scan and are not moved onto this package by this item.

Done when:

- [ ] `Split` then `Join` is byte-identical on every item in `docs/work/` and `docs/work/done/` (a test over the real files), and a refused construct returns `ErrUnsupported` with its line.
- [ ] Ticking a checklist item, appending a Log entry and replacing a section each change only their own lines (golden tests, including CRLF and a file without a final newline).
- [ ] A section that does not parse as its type is kept as opaque text and still edits (tests for the four known deviations: prose after the checklist and a Log continuation line).
- [ ] A seeded generator and the existing oracle battery confirm that a typed edit changes nothing outside its span, and that sections found equal the oracle's `##` headings on every input the scanner accepts.
- [ ] New terms (section tree, preamble, section etag, opaque section) are in `docs/glossary.md` and the package row in `docs/architecture.md` is updated.

## Next

Not started. Read ADR 261031-gam0, internal/body (scanner.go, fence.go) and internal/record/fields.go (the line scan that stays for now). The fail-closed Split is confirmed. Claim the item and start with the lossless Split/Join over the real files.

## Log

- 2026-10-31: filed.
- 2026-10-31: blocked_by 261031-g9zk: the parser choice (dependency or our scanner) sets how this is built.
- 2026-10-31: spike 261031-g9zk closed; unblocked, but wait for ADRs 261031-gam0 and 261031-cpbg to be accepted.
- 2026-10-31: the `Locate` interface, scanner, oracle and fuzz moved to 261031-yw5v, built first; ADRs 261031-gam0 and 261031-cpbg are accepted.
- 2026-10-31: Goal rewritten after 261031-yw5v, 261031-2j14 and 261031-wjs3 closed: ADRs accepted, Locate and FenceTracker exist, no wait on 261031-6bzq for the primitives; blocked_by removed.
- 2026-10-31: Alexander confirmed Split fails closed when Locate refuses; made ready.
- 2026-10-31: Alexander confirmed Split fails closed when Locate refuses; made ready.
- 2026-10-31: claimed (claude).
