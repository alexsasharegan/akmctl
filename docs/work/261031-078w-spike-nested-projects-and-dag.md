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

## Log

- 2026-10-31: filed.
- 2026-10-31: claimed (claude).
