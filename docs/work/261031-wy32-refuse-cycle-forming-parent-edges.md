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

From spike 261031-078w. Today only `parent` equal to the item itself is refused. Two more states are accepted by `set` and `new`, and `check` then fails: a parent chain that loops (WK4, reported once per member), and a parent that is already `done` or `dropped` while the child is open (WK8, found in the probe by setting `parent` of an open item to a done one, or creating a child of it). Refuse both at write time. For the loop, walk up the new parent's ancestor chain (parent is single-valued, so this is O(depth)); a cycle exists exactly when the item is on that chain. Each error names the edge. Done when:

- [ ] `set` and `new` (CLI and API) return `refused` with `details.rule` WK4 and the chain, for a parent that is the item's descendant at depth 2 and 3.
- [ ] They return `refused` with `details.rule` WK8 for an open item under a closed parent; an already-closed item may still be re-parented under a closed parent.
- [ ] Tests and API fixtures cover each; WK4 in `check` names the edge that closes the loop.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
