---
id: 261026-mrdr
type: task
status: ready
project: "[[Vault Views]]"
origin: "[[261025-g68j-build-pkm-workbench]]"
created: 2026-10-26
updated: 2026-10-26
authored_by: agent
title: "pkm-workbench: sort options"
---
# pkm-workbench: sort options

## Goal

The board and tables have a fixed order. Alexander wants to choose how they sort. This comes after [[261026-400a-task-timestamps-datetime]], so that date sorts have a stable order within a day.

Done when:

The board and each table can be sorted by the agreed keys in both directions from the keyboard; the current sort is visible; the choice survives a reload; tests cover the comparators including ties; README and keys help updated.

## Ideas (Alexander, 2026-10-26; not hard and fast)

- Every sort has **ascending and descending**.
- Sort **chronologically**: by created, by updated.
- Sort **lexicographically**: by name.
- **Group by project**, which may be a grouping and not a sort.
- Keep it modest. He is more interested in sorting than in filtering, but the two go together ([[261026-zg0g-workbench-filter]]).

## To settle

- Which keys and groupings ship first (created, updated, name, project, status, age, priority if it exists).
- Where the control lives: a key (cycle the sort, flip direction), a palette action, a visible indicator of the current sort.
- Per view or global; whether the choice is remembered (the board's theme choice is stored in the browser; the sort could be too).
- Within a lane or table group: sort inside each column, and what a drag means when the order is not manual.
- Missing values sort last, as the existing views do.
- A check that the server's order matches `Tasks.base` when it sorts the same way (`pkm-wb check` compares membership only today).

## Log

- 2026-10-26: filed by Claude Code on the Mac at Alexander's direction. `inbox`, not claimed; waits on the timestamp task.
- 2026-10-26: status inbox to ready (vocabulary change, [[261026-p5cz-task-status-vocabulary]]).
