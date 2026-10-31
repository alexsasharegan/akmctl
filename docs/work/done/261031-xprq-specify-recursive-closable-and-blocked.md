---
id: "261031-xprq"
title: specify recursive closable and blocked
status: done
created: 2026-10-31
updated: 2026-10-31
closed: 2026-10-31
summary: Define closable and blocked as recursive predicates over parent and blocked_by
priority: 10
---

# specify recursive closable and blocked

## Goal

From spike 261031-078w. Specify, in the AKM spec and guide, the dependency graph as one graph: explicit `blocked_by` edges plus an implicit edge from every parent to each child (a parent waits for its children). `closable(r)` is true when every child is complete; `blocked(r)` is true when a `blocked_by` entry is open or any child is blocked. A cycle in the union graph is an error at any depth. Done when:

- [x] The spec defines the predicates, the cycle rule (WK18) and the write-time refusal requirement.
- [x] WK8 is restated in terms of the predicate; the guide and the glossary have the new terms.
- [x] The fixtures are NOT part of this item (reworded): a WK18 check fixture would fail the suite until `check` implements it, so each lands with its implementing item, 261031-wy32 and 261031-z50q, in both copies (`docs/spec/akm-api/.fixtures` and `testdata/api`).

## Next

Close this item.

## Outcome

- Spec: `docs/spec/agent-ready-repo-spec.md`, new subsection "Closable, blocked and the dependency graph", WK8 restated, WK18 added, write-time refusal requirement. Guide: `docs/spec/akm-guide.md` (closing rule, rule list). Glossary: settled, complete, closable, blocked (derived), dependency graph. Changelog entry marks it specification only.
- Checks: `go vet`, `go test ./...`, `gofmt -l`, `npm run format:check` and `npm run schemas` (71 cases, 0 failures) pass; the fixtures were not touched (see Goal).
- Not done: the tool does not yet enforce WK18 or the recursive WK8 or the write-time refusal; `check` is unchanged. That is 261031-z50q and 261031-wy32.

Decision: none (the definitions follow from the existing parent and blocked_by fields).

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
- 2026-10-31: claimed (claude).
- 2026-10-31: spec, guide, glossary and changelog written; fixtures reassigned to 261031-wy32 and 261031-z50q.
