---
id: "261031-j4bc"
title: body section tree package
status: done
created: 2026-10-31
updated: 2026-11-01
closed: 2026-11-01
kind: feature
parent: 261031-gtw6
summary: internal/body gains a lossless section tree, typed Goal/Next/Log/Outcome and splice edits, on top of Locate, per ADR 261031-gam0
---

# body section tree package

## Goal

Build the in-memory body model that the body writer (spike [261031-6bzq](../261031-6bzq-spike-body-writer-api.md)) and the drawer will share, in `internal/body`, on top of `Locate` ([261031-yw5v](261031-yw5v-locate-scanner-and-goldmark-oracle.md), done) and the shared `FenceTracker` ([261031-2j14](261031-2j14-section-fence-tracking.md), done). The model and edit rule are decided: ADR [261031-gam0](../../adr/261031-gam0-work-item-body-model.md) (open body, typed core) and ADR [261031-cpbg](../../adr/261031-cpbg-goldmark-markdown-parser.md) (own scanner, goldmark as a test-only oracle); both are accepted, so nothing here waits on a decision.

This item is the package, not the CLI or API surface; the operations, addressing syntax and error shapes are 261031-6bzq's. It provides the primitives those will call:

- `Split(body) (Tree, error)` and `Tree.Join() []byte`: a preamble (text before the first `##`) and an ordered list of sections `{name, text}`, the name compared case-insensitively, found with `Locate` so headings in fences are not headings. Lossless: unknown sections and hand formatting are never rewritten.
- Typed views of the four core sections: Goal (prose plus the first contiguous run of `- [ ]` / `- [x]` lines, one line each), Next (prose), Log (append-only `- YYYY-MM-DD: text` entries) and Outcome (prose).
- Splice edits that change only the touched lines and normalize nothing else: tick or untick a checklist item (`[ ]` to `[x]` in place), append a Log entry, replace a section's text. Items are addressed by section name and 1-based position, with a section etag to detect a stale write; no ids are stored in the file.
- A typed section that does not parse (a Goal with prose after the checklist, a Log line with a continuation) keeps the rest as opaque text, and the typed part still edits.
- Line endings: a CRLF file stays CRLF, and an edit writes the line ending of the lines around it.

Decided (Alexander, 2026-10-31): when `Locate` refuses a construct (`ErrUnsupported`), `Split` returns that error and nothing is written (fail closed, ADR 261031-nvka), instead of falling back to the line scan `SetSection` uses. Today `SetSection` and `section` keep their own line scan and are not moved onto this package by this item.

Done when:

- [x] `Split` then `Join` is byte-identical on every item in `docs/work/` and `docs/work/done/` (a test over the real files), and a refused construct returns `ErrUnsupported` with its line.
- [x] Ticking a checklist item, appending a Log entry and replacing a section each change only their own lines (golden tests, including CRLF and a file without a final newline).
- [x] A section that does not parse as its type is kept as opaque text and still edits (tests for the four known deviations: prose after the checklist and a Log continuation line).
- [x] A seeded generator and the existing oracle battery confirm that a typed edit changes nothing outside its span, and that sections found equal the oracle's `##` headings on every input the scanner accepts.
- [x] New terms (section tree, preamble, section etag, opaque section) are in `docs/glossary.md` and the package row in `docs/architecture.md` is updated.

## Next

Done and verified; close with the Outcome.

## Outcome

- Round trip: `TestSplitJoinRecords` (`internal/body/tree_test.go`) splits and joins every record under `docs/work`, `docs/work/done` and `docs/adr` byte for byte, checks that the preamble and sections cover the whole body, and that none is refused. `TestSplitJoin` shows a blockquote returns `ErrUnsupported` from `Split` (fail closed, as Alexander confirmed).
- Edits touch only their own lines: `TestSetChecked` (tick, untick, no-op, stale etag, item outside the run; LF and CRLF), `TestLogEntriesAndAppend` (before a following section, CRLF, no final newline, empty section, heading as the last line) and `TestReplaceSection` (blank lines kept, CRLF, empty and last-line sections, Log refused, stale, missing). Ticking changes exactly one byte.
- Deviations: `TestDeviations` covers prose and a second task line after the checklist, a Log continuation line and a non-dated Log line; the rest stays opaque and the typed part still edits.
- Oracle and generator: `TestTreeAgainstOracle` ran 9,207 accepted documents (of 20,000, fixed seed) and 5,146 edits: the sections equal the oracle's level-2 headings, every edit leaves the text outside its section unchanged, and the edited body splits again. `FuzzTree` ran 120 seconds (about 10 million executions, 554 corpus entries) with no finding.
- Docs: glossary (section tree, preamble, section etag, opaque section) and the `internal/body` row in `docs/architecture.md`.

`gofmt -l .` (empty), `go vet ./...`, `go test -race ./...`, the records check (0 errors, 0 warnings) and `npm run format:check` pass; `go list -deps ./cmd/akmctl | grep -c goldmark` prints 0.

Not done: nothing calls the tree yet. `section`, the API and the drawer still use the line scan in `internal/record/fields.go` (so they do not refuse what `Split` refuses); moving them is for the writer spike 261031-6bzq, which now has the primitives. Next and Outcome are plain text with no typed view beyond `Text` (the spec's `Saved:` line is not read). Where two sections share a name, edits address the first. Section etags are in the package but no operation passes them yet.

Decision: none (the model and the fail-closed rule are in ADRs 261031-gam0, 261031-cpbg and 261031-nvka and Alexander's confirmation recorded in this item).

## Log

- 2026-10-31: filed.
- 2026-10-31: blocked_by 261031-g9zk: the parser choice (dependency or our scanner) sets how this is built.
- 2026-10-31: spike 261031-g9zk closed; unblocked, but wait for ADRs 261031-gam0 and 261031-cpbg to be accepted.
- 2026-10-31: the `Locate` interface, scanner, oracle and fuzz moved to 261031-yw5v, built first; ADRs 261031-gam0 and 261031-cpbg are accepted.
- 2026-10-31: Goal rewritten after 261031-yw5v, 261031-2j14 and 261031-wjs3 closed: ADRs accepted, Locate and FenceTracker exist, no wait on 261031-6bzq for the primitives; blocked_by removed.
- 2026-10-31: Alexander confirmed Split fails closed when Locate refuses; made ready.
- 2026-10-31: Alexander confirmed Split fails closed when Locate refuses; made ready.
- 2026-10-31: claimed (claude).
- 2026-11-01: built Split/Join, the Goal checklist, Log entries and the splice edits in internal/body/tree.go with golden tests, an oracle generator and FuzzTree; 120-second fuzz run clean.
