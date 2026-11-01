---
id: "261031-2j14"
title: section fence tracking
status: active
created: 2026-10-31
updated: 2026-10-31
kind: bug
summary: SetSection ends a code fence on any backtick or tilde fence line, so a mixed or longer fence can hide or invent a section heading
claimed_by: claude
claimed_at: 2026-11-01T06:08
---

# section fence tracking

## Goal

Found by the fuzz prototype in spike [261031-g9zk](done/261031-g9zk-spike-markdown-parser-dependency.md). `SetSection` in `internal/record/fields.go` flips its in-fence flag on any line starting with three backticks or `~~~`, so a fence opened with one and "closed" by the other, or a four-backtick fence holding a three-backtick line, makes it see a `## ` heading that is inside code (or miss one that is not). In the prototype, tracking the fence character and length (a closing fence is the same character, at least as long, with nothing after it) removed every such disagreement with goldmark. Done when:

- [x] A test with a mixed fence (backticks closed by `~~~`) and a four-backtick fence holding three fails today and passes after the fix.
- [x] `section` and anything else that finds headings use the same fence rule.
- [x] The fix is made in or ahead of the body package from [261031-j4bc](261031-j4bc-body-section-tree-package.md), so the rule lives in one place.

## Next

Done and verified; close with the Outcome.

## Log

- 2026-10-31: filed.
- 2026-10-31: claimed (claude).
- 2026-10-31: added body.FenceTracker (one fence rule, shared with the scanner); Sections, FirstHeading, AppendToSection, SetSection and the G1 heading search use it.
