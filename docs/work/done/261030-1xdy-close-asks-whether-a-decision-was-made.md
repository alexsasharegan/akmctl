---
id: "261030-1xdy"
title: close asks whether a decision was made
status: done
created: 2026-10-30
updated: 2026-10-30
closed: 2026-10-30
kind: chore
priority: 1
summary: The akm-work-items Close step asks whether the work made a hard-to-reverse choice, and links the ADR in decisions.
---

# close asks whether a decision was made

## Goal

Skill text only, no new field and no new rule. The `akm-work-items` Close section gets a step: ask whether the work made a choice that is hard to reverse or that rejected a real alternative. If yes, draft an ADR with `akm-record-decision` (or find the existing one) and list its id in the item's `decisions`; if no, write `Decision: none` in the Outcome. `akm-record-decision` "When" gains the same case: a choice noticed at close, not only before.

Why: `akm-record-decision` says to draft before making the choice, but an item can finish having made one nobody flagged, and nothing in the close flow asks. The maintainer's vault skill ran a five-question check at close for this. The tool cannot judge it, so a skill step is the right size; the `decisions` field already links an item to its ADRs.

Done when:

- [x] The Close section of `akm-work-items` and the "When" section of `akm-record-decision` are changed in both copies (`.agents/skills/` and `internal/initrepo/templates/skills/`), and the skill drift test passes.
- [x] The step names the `decisions` field and the `Decision: none` line, and does not claim the tool checks it.
- [ ] This repository closes one item with the new step as a trial, and the Log records whether it was usable.

## Outcome

- Close in `akm-work-items` (both copies) has a new step 2: ask about a hard-to-reverse choice, list the ADR in `decisions`, or write `Decision: none`; it says the tool does not check it.
- `akm-record-decision` "When" (both copies) adds the choice noticed at close.
- Drift test: `go test ./...` passes; `akmctl check` 0 errors, 0 warnings; `npm run format:check` clean.
- Trial: this item closes with the new step. The step was usable: one question, one line in the Outcome, no friction.

Decision: none.

## Log

- 2026-10-30: filed.
- 2026-10-30: filed from the vault migration review (task-lifecycle adr-check compared with akm-record-decision).
- 2026-10-30: claimed (claude).
- 2026-10-30: skills edited in both copies; closed as the trial of the new step.
