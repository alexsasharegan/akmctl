---
name: akm-track-work
description: Find, claim, work, checkpoint and close an AKM work item in docs/work/. Use at the start of a task, when stopping, and when a task finishes.
---

# Track work

The repository's `AGENTS.md` wins over this skill where they disagree. Fields, statuses and rules are in the AKM specification; this is the procedure.

## Start

1. `akmctl ls --status ready,active,blocked`. Without the tool: `grep -lE '^status: ?"?active"?$' docs/work/*.md`.
2. Resume the `active` item you hold. If you hold none, take the top `ready` item (roadmap order, then `priority`). If several `active` items are unclaimed, ask.
3. `akmctl claim <id> --as <name>` (claims extension). If someone else holds a fresh claim, do not work the item: tell the maintainer and take another. `claimed_at` is UTC, while `updated` is a local day and Log lines carry local time with an offset; a difference of hours is not an error.
4. Read `## Goal`, then `## Next`.

## No item fits

Before filing, read what `ls --status raw,ready,active,blocked` shows. `akmctl new work <slug>` creates a `raw` item and lists `similar` open items under the new path (up to five that share words with the slug and title). It never refuses. If one of them covers the work, add to that item (its `## Goal`, `## Next` or `## Log`) and, if you already created the new one, drop it with a reason in its Outcome. Otherwise keep the new item. Ask the maintainer which item to work. Do not type an id by hand; without the tool use the recipe in `AGENTS.md`.

## While working

- Keep `## Next` as the single next concrete step. Overwrite it at each checkpoint.
- Append lines to `## Log` with `akmctl log <id> <text>` (the tool stamps date, time and offset; add `--at` if the thing happened earlier); never rewrite old lines. Tick a Goal box with `akmctl tick <id> <n>`.
- Update both in the same commit as the code.
- Commit one step at a time (file, ready, start, claim, release, close, ...). Each command that changes an item prints a commit message: use it as given, never your own subject. If the `akm-commit` skill is installed (`.agents/skills/akm-commit`), it has the rest.
- Command details are in the `akm` skill (`.agents/skills/akm/references/`).
- A `ready` item needs a `## Goal` with a "done when" list; an `active` or `blocked` one needs a non-empty `## Next`; a `blocked` one needs `blocked_by`.

## Stop or hand off

Rewrite `## Next` so a cold agent can resume: where the work stands, the single next step, what to read first. `akmctl release <id>` if you will not continue.

## Close

1. Tick each "done when" box as its proof lands: `akmctl tick <id> <n> --evidence "<test or file>"`. `close done` refuses an empty box; if a requirement is gone, `akmctl remove <id> <n> --reason "<why>"`: it deletes the box, logs it and loses history (the tool does not track git), so commit first and, as an agent, ask the user before passing `--confirm`.
2. Write `## Outcome`: for each "done when" line, the test, command or file that proves it (for `dropped`, the reason).
3. Ask whether the work made a choice that is hard to reverse or that rejected a real alternative. If yes, draft an ADR with `akm-record-decision` (or find the one that already records it) and list its id in the item's `decisions` field. If no, end the Outcome with `Decision: none`. The tool does not check this; the answer is yours to give.
4. `akmctl close <id> done` (or `dropped`). It sets the status and `closed` (today), moves the file to `docs/work/done/` and refuses an empty Outcome.
5. Read the list of open items it prints (`dependents`): they named this item in `blocked_by`. For each, unblock it (`set status=ready`, and remove the entry from `blocked_by`) or re-point `blocked_by` at what it still waits on. The tool lists them and edits none.
6. `close` also rewrote the Markdown links the move breaks, and listed those files; do not fix links by hand unless it reports `link_errors`. Run `akmctl check`. Commit with the message `close` printed, staging its `paths`, the work and any roadmap change with it.

The template for a new item is `assets/work-item.md` (for use without the tool).
