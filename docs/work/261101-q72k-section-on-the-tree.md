---
id: "261101-q72k"
title: move section writes onto the section tree
status: raw
created: 2026-11-01
updated: 2026-11-01
kind: feature
parent: 261031-gtw6
summary: AddSection on body.Tree; section and the API write path use the tree and refuse unsupported constructs; delete SetSection
blocked_by: [261101-4j86]
---

# move section writes onto the section tree

## Goal

From spike [261031-6bzq](done/261031-6bzq-spike-body-writer-api.md) and proposed ADR [261101-tr8e](../adr/261101-tr8e-body-writes-through-the-section-tree.md). Give `body.Tree` an `AddSection(name, text)` that creates a missing section before `## Log` (otherwise at the end), as `record.SetSection` does, then move `section` and the other body writes (`AppendToSection`, `Sections`) onto the tree and delete the line scan. `TestTreeMatchesSetSection` (`internal/record/section_test.go`) holds the two to one behaviour today; extend it to creation before landing, and keep the cases `SetSection` tests now as tree tests. One behaviour changes: a body with a structural construct, or a local one in the edited section, can no longer be written by the tool; the response is `refused` with the construct and line (the classification is work item 261101-4j86). Done when:

- [ ] `AddSection` has golden tests for a missing section before Log, with no Log, in an empty body and in CRLF, and matches `SetSection` on those cases before it is removed.
- [ ] `section`, `close --outcome`, `AppendToSection` and the UI write path use the tree; `SetSection` and the duplicate fence scans in `internal/record/fields.go` are gone.
- [ ] A body with a structural construct, or a local one in the edited section, is refused by `section` with the construct and line, and nothing is written (test, spec section 4.16, fixture in both copies).
- [ ] Every existing `section` test and fixture still passes.

## Next

Not started. Read the spike's Outcome and ADR 261101-tr8e; make ready once the ADR is agreed.

## Log

- 2026-11-01: filed.
- 2026-11-01: filed by spike 261031-6bzq.
- 2026-11-01: refusal rule narrowed in ADR 261101-tr8e; blocked_by 261101-4j86 (the classification).
