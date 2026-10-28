---
name: work-items
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

`akmctl new work <slug>` creates a `raw` item. Ask the maintainer which item to work. Do not type an id by hand; without the tool use the recipe in `AGENTS.md`.

## While working

- Keep `## Next` as the single next concrete step. Overwrite it at each checkpoint.
- Append dated lines to `## Log`; never rewrite old lines.
- Update both in the same commit as the code.
- A `ready` item needs a `## Goal` with a "done when" list; an `active` or `blocked` one needs a non-empty `## Next`; a `blocked` one needs `blocked_by`.

## Stop or hand off

Rewrite `## Next` so a cold agent can resume: where the work stands, the single next step, what to read first. `akmctl release <id>` if you will not continue.

## Close

1. Write `## Outcome`: for each "done when" line, the test, command or file that proves it (for `dropped`, the reason).
2. `akmctl close <id> done` (or `dropped`). It sets the status, moves the file to `docs/work/done/` and refuses an empty Outcome.
3. Run `akmctl check`. Commit the work, the item and any roadmap change together.

The template for a new item is `assets/work-item.md` (for use without the tool).
