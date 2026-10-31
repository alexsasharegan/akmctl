---
id: "261031-pbws"
title: warn redundant blocked by on a descendant
status: done
created: 2026-10-31
updated: 2026-10-31
closed: 2026-10-31
kind: chore
summary: check warns when an item lists its own descendant in blocked_by
priority: 10
---

# warn redundant blocked by on a descendant

## Goal

From 261031-z50q. A parent already waits for its children, so a `blocked_by` entry naming one of its descendants repeats an edge it has. It is not a cycle (WK18) and is accepted, but it hides the real dependency and goes stale. Add rule WK19 (warning): an item's `blocked_by` names one of its descendants. Done when:

- [ ] Spec (rule table and guide), glossary if a term is added, changelog and a `check` test say and prove it, including that a child named by a grandparent counts and a sibling does not.
- [ ] A fixture `check-7` is in both copies.

## Next

Close this item.

## Outcome

- Rule WK19 in `internal/check/check.go` (`redundantBlockers`), proven by `TestRules` (WK19 blocked_by a child) and `TestRedundantBlockerRule` (child, grandchild named by a grandparent: warned; a sibling: not). Spec rule table and guide, changelog; fixture `check-7` (repo `redundant`) in both copies, and a warning alone leaves `clean` true.
- The earlier `TestBlockedByCyclesThroughTheParentTree` now expects the warning.
- Checks: `go vet`, `go test -race ./...`, `gofmt -l`, `npm run schemas` (76 cases), `npm run format:check`, `akmctl check` pass.

Decision: none.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from 261031-z50q; part of the patch release.
- 2026-10-31: claimed (claude).
- 2026-10-31: done.
