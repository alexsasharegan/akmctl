---
id: "261101-q72k"
title: move section writes onto the section tree
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
kind: feature
parent: 261031-gtw6
summary: AddSection on body.Tree; section and the API write path use the tree and refuse unsupported constructs; delete SetSection
---

# move section writes onto the section tree

## Goal

From spike [261031-6bzq](261031-6bzq-spike-body-writer-api.md) and proposed ADR [261101-tr8e](../../adr/261101-tr8e-body-writes-through-the-section-tree.md). Give `body.Tree` an `AddSection(name, text)` that creates a missing section before `## Log` (otherwise at the end), as `record.SetSection` does, then move `section` and the other body writes (`AppendToSection`, `Sections`) onto the tree and delete the line scan. `TestTreeMatchesSetSection` (`internal/record/section_test.go`) holds the two to one behavior today; extend it to creation before landing, and keep the cases `SetSection` tests now as tree tests. One behavior changes: a body with a structural construct, or a local one in the edited section, can no longer be written by the tool; the response is `refused` with the construct and line (the classification is work item 261101-4j86). Done when:

- [x] `AddSection` has golden tests for a missing section before Log, with no Log, in an empty body and in CRLF, and matches `SetSection` on those cases before it is removed.
- [x] `section`, `close --outcome`, `AppendToSection` and the UI write path use the tree; `SetSection` and the duplicate fence scans in `internal/record/fields.go` are gone.
- [x] A body with a structural construct, or a local one in the edited section, is refused by `section` with the construct and line, and nothing is written (test, spec section 4.16, fixture in both copies).
- [x] Every existing `section` test and fixture still passes.

## Next

Done; see Outcome. Next is 261101-ga4s (tick and log operations).

## Outcome

`section`, `close --outcome` and the Log line of `claim` and `release` write through the section tree (`record.EditSection`, `record.AppendToSection`). `SetSection` and the line scan that built it are gone. The tree gained `AddSection` (before the first Log, else at the end), `AppendText`, `WithEOL` and `LineEnding`.

- **Equivalence.** Before `SetSection` was deleted, a throwaway test compared it with the tree on every work item (replace, append, create, clear, three texts each); they matched byte for byte, except that clearing an already-empty section added a blank line each time in the old code. That quirk is not kept. Clearing a section with text now removes its lines whole; the tree had left an empty line (found by that comparison, fixed in `ReplaceSection`).
- **Kept as they were**, because fixtures and callers rely on them: a missing Log made by `claim` or `release` has its first line straight under the heading; a replaced last section ends the file with a line ending.
- **Refusal.** A structural construct anywhere, or a local one in the edited section, is `refused` with `details.reason: "unsupported_markdown"`, the construct, the file's line (as rule B1 reports it) and `structural`. Tests: `TestEditSectionEdges` (structural anywhere, local in another section allowed, local in the edited section refused), spec 4.16, fixture `section-7` in both copies.
- **Reads were not moved.** `record.Sections` (used by `show`, `check` and `NextLine`) stays a heading scan on the one fence rule, because a read must work on a body the tree refuses; `check` rule B1 reports those bodies. Moving it would need a read mode of the tree that tolerates structural constructs, which is not worth a second path today.
- Not done: the equivalence test `TestTreeMatchesSetSection` was retired with `SetSection`; `TestSplitJoinRecords` and `TestTreeAgainstOracle` hold the tree to lossless edits.

## Log

- 2026-11-01: filed.
- 2026-11-01: filed by spike 261031-6bzq.
- 2026-11-01: refusal rule narrowed in ADR 261101-tr8e; blocked_by 261101-4j86 (the classification).
- 2026-11-01: 261101-4j86 closed; blocked_by removed.
- 2026-11-01: claimed (claude).
- 2026-11-01: section writes on the tree; SetSection deleted; fixture section-7; closing.
