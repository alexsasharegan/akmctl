---
id: "261031-zek2"
title: close checks goal checkboxes
summary: "close refuses while a Goal checkbox is unchecked; a command ticks one box with evidence"
status: done
kind: feature
parent: "261031-v07x"
created: 2026-10-31
updated: 2026-11-01
closed: 2026-11-01
---

# close checks goal checkboxes

## Goal

From the consultation (gap 1). Rule WK7 checks that `## Goal` and `## Outcome` exist, but `close done` never reads the `- [ ]` lines, so a one-line Outcome closes an item whose criteria were not met. Done when:

- [x] `close <id> done` is `refused` while `## Goal` has an unchecked `- [ ]` line, and the error names each such line. `dropped` is not affected.
- [x] A way to close anyway is decided (a flag, or tick the box and say why) and written down; a hard-to-reverse choice gets a `proposed` ADR.
- [x] A command ticks one Goal line by number and appends the evidence to `## Log` (shape to decide, for example `check <id> <n> --evidence "..."`); an out-of-range number is `not_found`.
- [x] The API spec, the `close` reference page and the glossary cover the new behavior.

## Next

Done; ready to close.

## Outcome

close done now checks the Goal; tick takes evidence.

- Refusal: `uncheckedGoal` in `internal/api/ops_write.go`; `reason: unchecked_goal` with each box's line, text and, when tick can address it, position. Checked after open children, before any write; `dropped` and a body the tree cannot read are not checked. `TestCloseRefusesUncheckedGoal`, `TestCloseGoalEdges`; fixture `close-9`.
- Close anyway: decided as no override flag. The ways through are to tick the box with proof, change the Goal with `section` (a visible edit), or drop the item. A flag would reopen the gap this closes; it can be added later without breaking anything, so no ADR.
- Evidence: `tick --evidence` appends `ticked "<text>": <evidence>` (or `unticked`) to the Log in the same write, creating a Log if missing, also when the box was already in that state. `TestTickEvidence` (LF and CRLF); fixture `tick-6`; schema `tick.json`. The item asked for a new command; `tick` already existed, so this is a parameter, not `check <id> <n>`.
- Docs: API spec 4.12 and 4.21, `close` and `tick` references, `akm-track-work` Close step (and embedded copies), glossary (Unchecked goal, Tick), changelog.
- Fixture items: `261020-a2b3` now has a ticked Goal so the other close cases still close; several tests tick before closing. `go test -race ./...`, `npm run schemas` (93 cases), `akmctl check` and Prettier pass.

Decision: none.

## Log

- 2026-10-31: filed from the consultation. The `section` operation already writes a whole section, so the new part is the per-line tick, not general checklist editing.
- 2026-11-01 08:51-07:00: claimed (claude).
- 2026-11-01 08:55-07:00: ticked "`close <id> done` is `refused` while `## Goal` has an unchecked `- [ ]` line, and the error names each such line. `dropped` is not affected.": see Outcome
- 2026-11-01 08:55-07:00: ticked "A way to close anyway is decided (a flag, or tick the box and say why) and written down; a hard-to-reverse choice gets a `proposed` ADR.": see Outcome
- 2026-11-01 08:55-07:00: ticked "A command ticks one Goal line by number and appends the evidence to `## Log` (shape to decide, for example `check <id> <n> --evidence "..."`); an out-of-range number is `not_found`.": see Outcome
- 2026-11-01 08:55-07:00: ticked "The API spec, the `close` reference page and the glossary cover the new behavior.": see Outcome
