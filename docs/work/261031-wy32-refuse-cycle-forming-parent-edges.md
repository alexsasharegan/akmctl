---
id: "261031-wy32"
title: refuse cycle forming parent edges
status: raw
created: 2026-10-31
updated: 2026-10-31
summary: set and new refuse a parent edge that closes a cycle at any depth
---

# refuse cycle forming parent edges

## Goal

From spike 261031-078w. Today only `parent` equal to the item itself is refused; `check` finds longer cycles later (WK4) and reports every member. Refuse at write time: walk up the new parent's ancestor chain (parent is single-valued, so this is O(depth)); a cycle exists exactly when the item is on that chain. The error names the edge. Done when:

- [ ] `set` (and API `set`/`new`) returns `refused` with `details.rule` WK4 and the chain, for a parent that is the item's descendant at depth 2 and 3.
- [ ] A test and an API fixture cover it; WK4 in `check` names the edge that closes the cycle.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
