---
id: "261031-btdz"
title: refuse cycles when creating an item
status: done
created: 2026-10-31
updated: 2026-10-31
closed: 2026-10-31
kind: chore
priority: 10
summary: new with parent and blocked_by can create a dependency cycle
---

# refuse cycles when creating an item

## Goal

Found probing the patch release. API `new` (and so the workbench capture form) accepts `fields` with both `parent` and `blocked_by`: a child of P with `blocked_by: [P]`, or blocked by an item that already waits on P, makes a cycle that `check` then reports (WK18). Done when:

- [ ] `new` is refused with `rule` WK18, naming the path, for both cases; free-text entries and an unrelated blocker still work; a test and a fixture.
- [ ] The API spec says so.

## Next

Close this item.

## Outcome

- `newItemCycle` in `internal/api/ops_write.go`: `TestNewRefusesACycleFromParentAndBlockedBy` covers the parent as blocker, a blocker that waits on the parent, free text, an unrelated blocker and a blocker without a parent; fixture `new-1`; API spec updated.
- The same function is used for the workbench capture form, which calls the API `new`.

Decision: none.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from the release readiness check after 261031-99b6.
- 2026-10-31: claimed (claude).
- 2026-10-31: done.
