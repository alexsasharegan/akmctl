---
id: "261031-6c5n"
title: refuse dropping a parent with open children
status: done
created: 2026-10-31
updated: 2026-10-31
closed: 2026-10-31
kind: chore
priority: 10
summary: close dropped and check WK8 leave open children under a dropped parent
---

# refuse dropping a parent with open children

## Goal

Found probing the patch release. The spec says a parent may be dropped only if every open child is dropped or re-parented first, but `close <parent> dropped` succeeds with an open child and `check` reports nothing (WK8 covers `done` only). Done when:

- [ ] `close ... dropped` is refused with `open_children` when a child is open, as for `done`; CLI and API, with a test and a fixture.
- [ ] WK8 covers `dropped` as well as `done` (spec rule table, guide, changelog, a `check` test and a fixture).

## Next

Close this item.

## Outcome

- API `close` refuses `dropped` with `open_children` (WK8): `TestCloseDroppedRefusesOpenChildren`, `TestCloseParentWithOpenChildren` (updated, it expected the old behaviour), fixture `close-7`.
- `check` WK8 covers `dropped` and open descendants below a settled child: `TestRules` (WK8 dropped parent with open child). No separate `check` fixture; the rule is in the spec table and the table test.
- Spec rule table and parent bullet, API spec, guide and changelog updated. `go test -race ./...`, `gofmt`, `npm run schemas` (78 cases) pass.

Decision: none.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from the release readiness check after 261031-99b6.
- 2026-10-31: claimed (claude).
- 2026-10-31: done.
