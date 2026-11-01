---
id: "261031-2j14"
title: section fence tracking
status: done
created: 2026-10-31
updated: 2026-10-31
closed: 2026-10-31
kind: bug
summary: SetSection ends a code fence on any backtick or tilde fence line, so a mixed or longer fence can hide or invent a section heading
---

# section fence tracking

## Goal

Found by the fuzz prototype in spike [261031-g9zk](261031-g9zk-spike-markdown-parser-dependency.md). `SetSection` in `internal/record/fields.go` flips its in-fence flag on any line starting with three backticks or `~~~`, so a fence opened with one and "closed" by the other, or a four-backtick fence holding a three-backtick line, makes it see a `## ` heading that is inside code (or miss one that is not). In the prototype, tracking the fence character and length (a closing fence is the same character, at least as long, with nothing after it) removed every such disagreement with goldmark. Done when:

- [x] A test with a mixed fence (backticks closed by `~~~`) and a four-backtick fence holding three fails today and passes after the fix.
- [x] `section` and anything else that finds headings use the same fence rule.
- [x] The fix is made in or ahead of the body package from [261031-j4bc](../261031-j4bc-body-section-tree-package.md), so the rule lives in one place.

## Next

Done and verified; close with the Outcome.

## Outcome

- Mixed and four-backtick fences: `TestFenceRule` in `internal/record/section_test.go` (backticks closed by `~~~`, four backticks holding three, tildes holding backticks) failed for `Sections`, `SetSection`, `AppendToSection` and `FirstHeading` before the fix and passes after; `TestFenceRule` in `internal/body/fence_test.go` pins the rule itself. A check case ("G1 heading hidden in a longer fence") covers the AGENTS.md heading search.
- One rule everywhere: `Sections`, `FirstHeading`, `AppendToSection`, `SetSection` (`internal/record/fields.go`) and the G1 search (`internal/check/check.go`) use `body.FenceTracker`. Not changed on purpose: `internal/check/links.go` strips code before link checks and tolerates indented fences; it finds no headings, so it is for 261031-wjs3.
- In one place: `internal/body/fence.go` holds `fenceOpen` and `fenceClose`, used by both the scanner and `FenceTracker`; the scanner's battery, generator and saved fuzz findings still pass.

`go vet`, `go test -race ./...`, `gofmt -l .`, the records check and `npm run format:check` pass (records check run after closing). Not done: `SetSection` still scans lines itself rather than calling `Locate`, so it does not refuse the constructs the scanner refuses (HTML lines, blockquotes); that is 261031-j4bc's typed layer, and changing it would make section writes fail closed (ADR 261031-nvka), a decision for that item.

Decision: none.

## Log

- 2026-10-31: filed.
- 2026-10-31: claimed (claude).
- 2026-10-31: added body.FenceTracker (one fence rule, shared with the scanner); Sections, FirstHeading, AppendToSection, SetSection and the G1 heading search use it.
