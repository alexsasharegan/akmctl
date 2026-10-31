---
id: "261031-wy32"
title: refuse cycle forming parent edges
status: done
created: 2026-10-31
updated: 2026-10-31
closed: 2026-10-31
summary: set and new refuse a parent edge that closes a cycle at any depth
priority: 10
---

# refuse cycle forming parent edges

## Goal

From spike 261031-078w. Today only `parent` equal to the item itself is refused. Two more states are accepted by `set` and `new`, and `check` then fails: a parent chain that loops (WK4, reported once per member), and a parent that is already `done` or `dropped` while the child is open (WK8, found in the probe by setting `parent` of an open item to a done one, or creating a child of it). Refuse both at write time. For the loop, walk up the new parent's ancestor chain (parent is single-valued, so this is O(depth)); a cycle exists exactly when the item is on that chain. Each error names the edge. Done when:

- [ ] `set` and `new` (CLI and API) return `refused` with `details.rule` WK4 and the chain, for a parent that is the item's descendant at depth 2 and 3.
- [ ] They return `refused` with `details.rule` WK8 for an open item under a closed parent; an already-closed item may still be re-parented under a closed parent.
- [ ] Tests and API fixtures cover each; WK4 in `check` names the edge that closes the loop.

## Next

Close this item.

## Outcome

- Loop at depth 2 and 3, CLI and API `set`: `TestSetRefusesParentLoopsAtAnyDepth` (`internal/api/parent_edge_test.go`); fixture `set-9`. `new` cannot make a loop (the id is new).
- Open item under a closed parent, `set` and `new`: `TestOpenItemCannotJoinAClosedParent`; fixture `set-10`. A closed item may join a closed parent: `TestClosedItemMayJoinAClosedParent`.
- A violation already present does not block unrelated changes: `TestSetKeepsAnExistingLoopEditable`.
- The check on a real repository: `akmctl set` of a loop prints `refused: parent ... would make a loop: a -> b -> a` and exits 2.
- Not done: `check` still reports WK4 once per loop member, not once per loop (that is the WK18 reporting rule, with 261031-z50q); a hand-edited file can still hold a loop.
- API spec (section 4.15 set and 4.14 new), changelog and both fixture copies updated; `go vet`, `go test -race ./...`, `gofmt -l`, `npm run schemas` (73 cases) pass.

Decision: none (it enforces rules the spec already states).

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
- 2026-10-31: claimed (claude).
- 2026-10-31: parentEdgeRule in internal/api/ops_write.go, tests, fixtures set-9 and set-10, API spec and changelog.
