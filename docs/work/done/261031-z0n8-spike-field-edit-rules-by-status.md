---
id: "261031-z0n8"
title: spike field edit rules by status
summary: "Spike: which fields and body sections may be edited in which status, and whether the tool or only the UI enforces it"
status: done
kind: spike
parent: "261031-gtw6"
created: 2026-10-31
updated: 2026-11-01
closed: 2026-11-01
---

# spike field edit rules by status

## Goal

Editing `priority` makes sense while an item is `raw` or `ready`; changing it while `active` or after `done` is retroactive and probably wrong. Work out the rule for every field and every body section, so the drawer shows controls only where an edit is meaningful. Depends on the section model in [261031-qbat](261031-qbat-spike-work-item-body-model.md) and decides what [261028-esyt](../261028-esyt-ui-field-editing.md) may edit. Find out:

- A matrix: each field (`summary`, `kind`, `priority`, `parent`, `blocked_by`, `cycle`, `keep`, the extra `project`) and each body section against each status (`raw`, `ready`, `active`, `blocked`, `done`, `dropped`): editable, editable with a warning, or read-only, and why. Include appending: a log line is fine in any open status; a closed item's `## Outcome` is not.
- Where it is enforced: only in the UI, in `set` and the body writer ([261031-6bzq](261031-6bzq-spike-body-writer-api.md)) with an error code, or as a `check` rule. Whether an agent may override with a flag, and whether the status transition itself (reopening) is the sanctioned way to edit a closed item.
- How the drawer shows a locked field: hidden, disabled with the reason, or editable behind a "reopen" action.
- Which edits are legitimate corrections rather than retroactive changes (a typo in a closed item), and how they are recorded in `## Log`.

Done when:

- [x] The matrix is written and tried against the last 20 closed items and the open ones, listing any edit it would wrongly block.
- [x] A recommendation (matrix, enforcement point, locked-field display, correction path) is in `## Outcome`; a hard-to-reverse choice has a `proposed` ADR; follow-up items are filed and 261028-esyt is updated.

## Next

Write the Outcome, close, update 261028-esyt.

## Outcome

Recommendation: the tool, not the drawer, owns the rules. A closed work item is read-only except through a logged correction (proposed ADR [261101-dxcn](../../adr/261101-dxcn-closed-records-read-only-except-corrections.md)); open items stay as editable as now.

**Matrix.** `yes` editable, `warn` editable with a warning in the drawer (the API never refuses on a warning), `no` read-only. Closed means `done` or `dropped`.

| Field or section                                       | raw                                                              | ready  | active      | blocked | closed                       |
| ------------------------------------------------------ | ---------------------------------------------------------------- | ------ | ----------- | ------- | ---------------------------- |
| `title`, `summary`                                     | yes                                                              | yes    | yes         | yes     | correction                   |
| `kind`                                                 | yes                                                              | yes    | warn        | warn    | no                           |
| `priority`                                             | yes                                                              | yes    | warn        | yes     | no                           |
| `parent`                                               | yes                                                              | yes    | warn        | warn    | no                           |
| `blocked_by`                                           | yes                                                              | yes    | yes         | yes     | no                           |
| `cycle`                                                | yes                                                              | yes    | yes         | yes     | no                           |
| `keep`                                                 | yes                                                              | yes    | yes         | yes     | yes                          |
| `decisions`                                            | yes                                                              | yes    | yes         | yes     | no                           |
| `project` and others                                   | no                                                               | no     | no          | no      | no (not settable; migration) |
| `status`, claim fields, `created`, `updated`, `closed` | tool-managed (`set` moves among open statuses, `claim`, `close`) |        |             |         |                              |
| Goal, Next, Context                                    | yes                                                              | yes    | warn (Goal) | yes     | correction (not Next)        |
| Outcome                                                | yes                                                              | yes    | yes         | yes     | correction                   |
| Log                                                    | append (`log`)                                                   | append | append      | append  | correction line only         |

Why the `warn` and `no` cells: `kind`, `priority` and `parent` describe the item before work and ranking it, so changing them once it is under way, or after it closed, rewrites what it was; `blocked_by` and `cycle` change as plans slip, so they stay free while open; `blocked_by` cannot empty a `blocked` item (WK5 already refuses). `keep` is the maintainer's later decision about compaction, so it is the one field a closed item keeps open. Goal edits on an active item move the goalposts and carry a warning that suggests a `log` line.

**Tried against the history.** Every commit that changed a work item's frontmatter, by status before the edit: `raw` priority 7, summary 6, title 4, `blocked_by` 10, kind 1; `ready` summary 2, `blocked_by` 1; `active` `decisions` 2; `done` `origin` and `project` 12, all by the layout migration. The matrix allows all of these, `migrate` being exempt. Nothing in the history changed `priority`, `kind`, `parent` or `cycle` on an active or closed item, so the `warn` cells rest on reasoning, not data. Body edits to closed items in the last 20 closes: five hand repairs of a relative link whose target moved into `done/`, and one correction of the Outcome with a Log line (261031-6bzq, 2026-11-01). The matrix would block both; the correction goes through the correction path, and the link repair becomes `close`'s job ([261101-1smq](../261101-1smq-close-repairs-links-to-moved-items.md)).

**Enforcement.** Today `section` refuses a closed item but `set` does not, and the Log of a closed item cannot be appended to. Recommend enforcing the one hard line, closed, in `set` and `section` with `refused` and `reason: "closed"`, because an agent can otherwise rewrite `priority` on a done item. Not as a `check` rule: `check` reads a state, not an edit. The warnings are a drawer matter. No override flag: the sanctioned ways to change a closed item are `keep`, and a call carrying `correction` (a one-line reason), which writes `corrected <target>: <reason>` to the Log in the same write, the pattern of `tick --evidence`. Reopening is not the way: WK9 keeps a closed item closed, a reverse move would shift its parent's closability and the commit steps, and a wrongly closed item is better replaced by a new one that cites it.

**Drawer.** A locked field shows as read-only text with its reason in the cell (not hidden: the table shows every field), and a closed item has a "Correct…" action for the allowed targets, which asks for the reason. No "reopen". The controls come from an `edit` member of the record read ([261101-3che](../261101-3che-show-reports-edit-states.md)), so the matrix lives only in the tool.

**Follow-ups filed (raw, waiting on ADR 261101-dxcn):** [261101-3vwz](261101-3vwz-closed-items-read-only-with-correction.md) (the refusal and `correction`), [261101-3che](../261101-3che-show-reports-edit-states.md) (the `edit` member), [261101-1smq](../261101-1smq-close-repairs-links-to-moved-items.md) (`close` repairs links). 261028-esyt is unblocked and updated.

## Log

- 2026-10-31: filed.
- 2026-11-01 09:06-07:00: claimed (claude).
- 2026-11-01 09:09-07:00: ticked "The matrix is written and tried against the last 20 closed items and the open ones, listing any edit it would wrongly block.": matrix in the Outcome, tried against the git history of docs/work
- 2026-11-01 09:09-07:00: ticked "A recommendation (matrix, enforcement point, locked-field display, correction path) is in `## Outcome`; a hard-to-reverse choice has a `proposed` ADR; follow-up items are filed and 261028-esyt is updated.": recommendation in the Outcome; ADR 261101-dxcn proposed; 3vwz, 3che and 1smq filed
