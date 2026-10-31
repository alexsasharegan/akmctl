---
id: "261031-078w"
title: spike nested projects and dag
status: active
created: 2026-10-31
updated: 2026-10-31
kind: spike
summary: What AKM supports for nested projects and DAG-shaped work, and which features to add
claimed_by: claude
claimed_at: 2026-10-31T18:12
---

# spike nested projects and dag

## Goal

Find out what AKM supports today for projects nested more than one level (project > sub-projects > work items), where a project is a work item that others name as `parent`, and which gaps need scripted features. `parent` is single-valued, so it forms a tree; `blocked_by` adds dependency edges, so the two together form a DAG with the no-cycle rule on both.

Done when:

- [ ] A three-level tree (tier 1, tier 2, leaves) is built in a scratch repository and `new`, `set`, `ls --parent`, `close`, `check`, `board`, `find`, the roadmap rules and the workbench are exercised on it; findings are in `## Outcome`.
- [ ] Cycle handling is tested for `parent` (WK4) and for `blocked_by` (is a cycle or a parent/child block checked?), and the result is written down.
- [ ] Each gap has a proposed feature with a one-line cost and a recommendation (build, defer or skip); candidates: subtree listing (`ls --descendants` or `tree`), progress roll-up, split a project into children, re-parent, ready-leaf routing, blocked_by cycle check.
- [ ] Follow-up work items are filed for the features recommended; a hard-to-reverse choice gets a `proposed` ADR.

## Outcome

Probed with the built tool in a scratch repository (tier 1 > two tier 2 > three leaves, `parent` set with `set` after `new`).

**Works today**

- `parent` nests to any depth; `ls --parent <id>` lists direct children only.
- `close` refuses a parent with a child not done or dropped (WK8), at every level: tier 1 is refused while a tier 2 is open.
- `parent` cycles are found by `check` (WK4, "parent chain forms a cycle").
- `close` reports `dependents` (open items whose `blocked_by` names the closed one).

**Gaps found**

1. `parent` is single-valued: the hierarchy is a tree. The DAG comes from `blocked_by` (many-to-many) and is not a grouping. An item belonging to two projects is not expressible.
2. `set parent=<descendant>` is accepted (only self-parent is refused). `check` then reports WK4 on every item in the cycle, with no pointer to the edge that closed it. `new`/`set` should refuse the edge, as they already do for WK4 self-parent.
3. `blocked_by` cycles are not checked: leaf1 blocked by leaf2 and leaf2 blocked by leaf1 gives `0 errors, 0 warnings`.
4. A parent blocked by its own child (tier2b blocked_by leaf3, where leaf3 is its child) is also accepted; with WK8 that is a deadlock, since the parent cannot close before the child.
5. No subtree view: `ls --descendants`, a tree print and `board` grouped by project do not exist (`unknown parameter "descendants"`). Deep nesting means one `ls --parent` call per level.
6. No roll-up: no "3 of 5 leaves done" per parent, at any level.
7. The CLI `new work` takes no `parent` or `kind` (the API does: `fields`). Splitting a project into N children is N `new` plus N `set` calls.
8. No re-parent helper that moves a subtree, and no way to ask "which ready leaf under project X is next?" (`ls --next` ignores the tree).

**Proposed features** (each a follow-up item; recommendation first)

- Build: refuse a `parent` edge that creates a cycle at write time, and name the edge (gap 2). Small.
- Build: `check` rule for `blocked_by` cycles and for a parent blocked by its own descendant (gaps 3, 4). Small; an error, like WK4.
- Build: `ls --descendants <id>` plus `--tree` print, with a `depth` field in JSON (gap 5). Medium.
- Build: `--parent`, `--kind` on CLI `new`, and `akmctl split <id> <slug>...` creating children in one step (gaps 7, 8). Small to medium.
- Build: roll-up counts (children done / total, direct and transitive) in `show`, `ls --tree` and the workbench (gap 6). Medium.
- Defer: multi-parent / "belongs to" links (gap 1). A second field is a spec change and needs an ADR; `blocked_by` plus the workbench `project` filter covers most cases for now.
- Defer: ready-leaf routing under a project (`ls --next --under <id>`), after `--descendants` exists.

**Not verified:** the workbench's display of three levels; whether WK8 holds when a file is edited by hand into `done/`; the roadmap rule RM4 for nested ready items.

Decision: none (the spike changes no behavior; a multi-parent field would need an ADR in its own item).

## Next

File follow-up items for the recommended features (cycle refusal and checks first), then close this spike.

## Log

- 2026-10-31: filed.
- 2026-10-31: claimed (claude).
- 2026-10-31: probed nesting, cycles and blocked_by in a scratch repository; findings in Outcome.
