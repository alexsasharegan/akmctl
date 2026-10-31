---
id: "261031-z50q"
title: check blocked by cycles and parent child deadlocks
status: done
created: 2026-10-31
updated: 2026-10-31
closed: 2026-10-31
summary: check rule for blocked_by cycles and for a parent blocked by its descendant
priority: 10
---

# check blocked by cycles and parent child deadlocks

## Goal

From spike 261031-078w. `check` accepted two items blocked on each other. Build the union graph (see 261031-xprq): `blocked_by` edges plus parent-waits-for-child edges, and find cycles in O(items + edges). Report one cycle per strongly connected component, not one line per member. Free-text `blocked_by` entries are ignored.

Correction to the spike: a parent that lists its own child in `blocked_by` is not a deadlock. It repeats an edge the parent already has, so it is no cycle. The real deadlock is a child that waits for its parent or grandparent. Done when:

- [x] `check` reports WK18 for a mutual block, a longer loop and a child blocked by an ancestor, once per cycle, with each a test and a fixture; a parent listing its child is accepted.
- [x] `set` refuses a `blocked_by` entry, and a `parent`, that would close a cycle.

## Next

Close this item.

## Outcome

- `internal/depgraph` (Tarjan SCC, paths): `internal/depgraph/depgraph_test.go`. `check` WK18 and WK4-once-per-loop: `TestRules` (WK18 mutual block, child blocked by its parent, grandchild chain back to its parent), `TestDependencyCycleIsReportedOncePerCycle`, `TestParentLoopIsReportedOnceAsWK4`, `TestParentListingItsChildInBlockedByIsNotACycle`; fixture `check-6` (repo `cyclic`).
- Write-time refusal: `TestSetRefusesBlockedByCycles`, `TestBlockedByCyclesThroughTheParentTree`, `TestParentEdgeCannotCompleteADependencyCycle`; fixture `set-11`. Fixtures are in both copies.
- Checks: `go vet`, `go test -race ./...`, `gofmt -l`, `npm run schemas` (75 cases) and `npm run format:check` pass.
- Corrected a mistake in the spike (261031-078w gap 4) and in the spec text of 261031-xprq: a parent blocked by its own child is redundant, not a deadlock. Spec, API spec and changelog now say so.
- Not done: a warning for the redundant entry (parent listing a descendant); `new` with `blocked_by` is not checked, since a new id cannot already be waited on.

Decision: none.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
- 2026-10-31: claimed (claude).
- 2026-10-31: depgraph package, WK18 in check, write-time refusals, fixtures check-6 and set-11, spec correction.
