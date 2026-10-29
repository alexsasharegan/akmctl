---
id: 261026-zg0g
type: task
status: ready
project: Vault Views
origin: 261025-g68j
created: 2026-10-26
updated: 2026-10-26
authored_by: agent
title: "pkm-workbench: filters"
---

# pkm-workbench: filters

## Goal

The board and tables show every open task. Alexander wants to narrow them, with sorting as the higher priority ([261026-mrdr-workbench-sort](261026-mrdr-workbench-sort.md)). It also depends on 261026-400a-task-timestamps-datetime if a filter uses dates.

Done when:

Filters by the agreed fields work from the keyboard on the board and tables, active filters are visible and clearable, the choice survives a reload, tests cover include and exclude, and README and keys help are updated.

## Ideas (Alexander, 2026-10-26; not hard and fast)

- Filter by **project**: show only this project, or **exclude** this project.
- Filters should be able to say "only" and "not".
- Do not go crazy. A small set that covers the daily cases.

## To settle

- Which fields to filter on first: project, status, text in the name, age or staleness, presence of `adr` or `authored_by`.
- The control: a palette action, a key to open a filter box, a visible chip for each active filter that can be cleared with a key.
- Combining filters (all must match is the simple rule), and whether an empty result explains itself.
- Whether filters apply to the board, the tables, or both, and whether they are remembered across reloads.
- How a filter interacts with watch mode and the animation queue: an agent's move of a task that is filtered out should not animate in view, but a toast could say it happened.
- Whether `Tasks.base` views can be expressed as saved filters, so the two stay aligned.

## Log

- 2026-10-26: filed by Claude Code on the Mac at Alexander's direction. `inbox`, not claimed; waits on the timestamp task and, preferably, on sorting.
- 2026-10-26: status inbox to ready (vocabulary change, 261026-p5cz-task-status-vocabulary).
