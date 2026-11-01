---
id: "261031-z0n8"
title: spike field edit rules by status
summary: "Spike: which fields and body sections may be edited in which status, and whether the tool or only the UI enforces it"
status: raw
kind: spike
parent: "261031-gtw6"
created: 2026-10-31
updated: 2026-10-31
---

# spike field edit rules by status

## Goal

Editing `priority` makes sense while an item is `raw` or `ready`; changing it while `active` or after `done` is retroactive and probably wrong. Work out the rule for every field and every body section, so the drawer shows controls only where an edit is meaningful. Depends on the section model in [261031-qbat](done/261031-qbat-spike-work-item-body-model.md) and decides what [261028-esyt](261028-esyt-ui-field-editing.md) may edit. Find out:

- A matrix: each field (`summary`, `kind`, `priority`, `parent`, `blocked_by`, `cycle`, `keep`, the extra `project`) and each body section against each status (`raw`, `ready`, `active`, `blocked`, `done`, `dropped`): editable, editable with a warning, or read-only, and why. Include appending: a log line is fine in any open status; a closed item's `## Outcome` is not.
- Where it is enforced: only in the UI, in `set` and the body writer ([261031-6bzq](261031-6bzq-spike-body-writer-api.md)) with an error code, or as a `check` rule. Whether an agent may override with a flag, and whether the status transition itself (reopening) is the sanctioned way to edit a closed item.
- How the drawer shows a locked field: hidden, disabled with the reason, or editable behind a "reopen" action.
- Which edits are legitimate corrections rather than retroactive changes (a typo in a closed item), and how they are recorded in `## Log`.

Done when:

- [ ] The matrix is written and tried against the last 20 closed items and the open ones, listing any edit it would wrongly block.
- [ ] A recommendation (matrix, enforcement point, locked-field display, correction path) is in `## Outcome`; a hard-to-reverse choice has a `proposed` ADR; follow-up items are filed and 261028-esyt is updated.

## Log

- 2026-10-31: filed.
