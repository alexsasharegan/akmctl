---
id: "261031-z50q"
title: check blocked by cycles and parent child deadlocks
status: ready
created: 2026-10-31
updated: 2026-10-31
summary: check rule for blocked_by cycles and for a parent blocked by its descendant
priority: 10
---

# check blocked by cycles and parent child deadlocks

## Goal

From spike 261031-078w. `check` accepts two items blocked on each other and a parent blocked on its own child. Build the union graph (see specify-recursive-closable-and-blocked): `blocked_by` edges plus parent-waits-for-child edges, and find cycles with a depth-first search using three colours (or Kahn's algorithm), O(items + edges). Report one cycle path per strongly connected component, not one line per member. Free-text `blocked_by` entries are ignored. Done when:

- [ ] A new error rule fires for a mutual block, a longer loop, a parent blocked by a descendant and a child blocked by an ancestor; each has a test and a fixture.
- [ ] `set status=blocked blocked_by=...` is refused when the entry would close a cycle.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
