---
name: akm-work-items
description: Find, claim, work, checkpoint and close an AKM work item in docs/work/. Use at the start of a task, when stopping, and when a task finishes.
---

# Work items

The repository's `AGENTS.md` wins over this skill where they disagree. Fields, statuses and rules are in the AKM specification; this is the procedure.

## Start

1. `akmctl ls --status ready,active,blocked`. Without the tool: `grep -lE '^status: ?"?active"?$' docs/work/*.md`.
2. Resume the `active` item you hold. If you hold none, take the top `ready` item (roadmap order, then `priority`). If several `active` items are unclaimed, ask.
3. `akmctl claim <id> --as <name>` (claims extension). If someone else holds a fresh claim, do not work the item: tell the maintainer and take another.
4. Read `## Goal`, then `## Next`.

## No item fits

Before filing, read what `ls --status raw,ready,active,blocked` shows. `akmctl new work <slug>` creates a `raw` item and lists `similar` open items under the new path (up to five that share words with the slug and title). It never refuses. If one of them covers the work, add to that item (its `## Goal`, `## Next` or `## Log`) and, if you already created the new one, drop it with a reason in its Outcome. Otherwise keep the new item. Ask the maintainer which item to work. Do not type an id by hand; without the tool use the recipe in `AGENTS.md`.

## While working

- Keep `## Next` as the single next concrete step. Overwrite it at each checkpoint.
- Append dated lines to `## Log`; never rewrite old lines.
- Update both in the same commit as the code.
- If the `akm-commit` skill is installed (`.agents/skills/akm-commit`), follow it for commits.
- Command details are in the `akm` skill (`.agents/skills/akm/references/`).
- A `ready` item needs a `## Goal` with a "done when" list; an `active` or `blocked` one needs a non-empty `## Next`; a `blocked` one needs `blocked_by`.

## Stop or hand off

Rewrite `## Next` so a cold agent can resume: where the work stands, the single next step, what to read first. `akmctl release <id>` if you will not continue.

## Close

1. Write `## Outcome`: for each "done when" line, the test, command or file that proves it (for `dropped`, the reason).
2. Ask whether the work made a choice that is hard to reverse or that rejected a real alternative. If yes, draft an ADR with `akm-record-decision` (or find the one that already records it) and list its id in the item's `decisions` field. If no, end the Outcome with `Decision: none`. The tool does not check this; the answer is yours to give.
3. `akmctl close <id> done` (or `dropped`). It sets the status and `closed` (today), moves the file to `docs/work/done/` and refuses an empty Outcome.
4. Read the list of open items it prints (`dependents`): they named this item in `blocked_by`. For each, unblock it (`set status=ready`, and remove the entry from `blocked_by`) or re-point `blocked_by` at what it still waits on. The tool lists them and edits none.
5. Run `akmctl check`. Commit the work, the item and any roadmap change together.

The template for a new item is `assets/work-item.md` (for use without the tool).
