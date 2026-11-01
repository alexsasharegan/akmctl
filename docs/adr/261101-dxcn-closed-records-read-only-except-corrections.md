---
id: "261101-dxcn"
title: Closed work items are read-only except through a logged correction
summary: "set and section refuse a closed work item unless the call carries a correction reason, written to the Log in the same write; keep is always editable"
status: proposed
date: 2026-11-01
scope: [internal/api]
---

# Closed work items are read-only except through a logged correction

## Context

Work item 261031-z0n8 asked which fields and sections may change in which status. Today the tool is inconsistent: `section` refuses a closed item (`closed`), but `set` changes any settable field of one, and nothing can append to a closed item's Log (`log` refuses it). WK9 keeps a closed item closed, so there is no reopen.

The history of `docs/work` (every commit that changed a work item's frontmatter) shows what people really changed, by the item's status before the edit: `raw` (priority 7, summary 6, title 4, `blocked_by` 10, kind 1), `ready` (summary 2, `blocked_by` 1), `active` (`decisions` 2, adding ADR links) and `done` (`origin` and `project`, 12 in all, all by the layout migration). Nothing changed `priority`, `kind`, `parent` or `cycle` on an `active` or closed item. Edits to a closed item's body were of two kinds: a relative link repaired because its target moved into `done/` (five times, by hand), and one correction of an Outcome paragraph with a Log line explaining it (2026-11-01, 261031-6bzq). Both were hand edits the tool could not do.

## Considered Options

- **(a) Do nothing.** The drawer hides or disables whatever it likes, and the API stays open. Two rules in two places, and an agent can still rewrite a closed item's `priority` through `set`.
- **(b) Closed means read-only, with no way through.** A typo in an Outcome means a hand edit or a new item. The 6bzq correction could not be made by the tool.
- **(c) Reopen.** A closed item can return to `active`. WK9 and `closed` would need a reverse move, the commit trailers a new step, and a done item's children and parent's closability would shift under it. A larger change than the problem.
- **(d) Read-only, with a correction parameter.** `set` and `section` on a closed item are refused unless the call carries `correction`, a one-line reason; the tool writes `corrected <field or section>: <reason>` to the Log in the same write. The pattern is `tick --evidence`. Chosen.
- **Tiers for open items.** The API could refuse an edit that is meaningful but late (a `priority` change on an `active` item). Chosen against: no case in the history, and a warning is a drawer matter. The API refuses only what it can never be right to do.

## Decision

Option (d).

- **Closed (`done`, `dropped`) work items.** `set` and `section` refuse with `refused`, `reason: "closed"`, naming `correction`. With `correction` they change the allowed targets and append the Log line, creating `## Log` if absent. Allowed targets: `title`, `summary` and the sections `Outcome`, `Goal`, `Context` and any other prose section. Never, even with a reason: `status`, `kind`, `priority`, `parent`, `blocked_by`, `cycle`, `decisions` (what the item was while open), and `Log` (append-only; the tool writes it). The reason must be non-empty and one line.
- **`keep`** is editable in every status with no reason: it records a later decision about compaction, not the history of the item.
- **Open items** (`raw`, `ready`, `active`, `blocked`): every settable field and section as now. The `status` moves keep their own rules (WK5, WK8, WK9, WK18); `Log` stays append-only.
- **The tool reports the matrix.** A record read carries, for each field and section, `edit: yes`, `warn` or `no`, with a short `reason` for `warn` and `no`, so the drawer draws controls from it and does not copy the table. `warn` is advice (a field changed after work began) and never makes the API refuse.
- **No `check` rule.** `check` sees a state and not an edit, and it cannot know whether `priority` changed after close.
- **Exempt.** `migrate` and `close`/`claim`/`release` write what they write. The migration is the sanctioned writer of a layout change on closed items.

## Consequences

- A closed item's history is stable by default, and a correction is a visible, dated Log line, as the 6bzq one was by hand.
- Work item for the tool: refuse in `set`, accept `correction` in `set` and `section`, write the Log line; an API 1.x addition, with fixtures and spec text. A second item adds the `edit` member to `show`.
- Link repair after a `close` stays a hand edit until a third item makes `close` rewrite the links that point at the moved file; `correction` is not meant for it.
- A drawer can show a closed item's fields as read-only text with the reason, and a "Correct…" action for the few allowed targets. No "reopen" control exists.
- Reopening an item stays out of scope; a wrongly closed item is replaced by a new one that cites it.
