---
id: "261031-xprq"
title: specify recursive closable and blocked
status: ready
created: 2026-10-31
updated: 2026-10-31
summary: Define closable and blocked as recursive predicates over parent and blocked_by
priority: 10
---

# specify recursive closable and blocked

## Goal

From spike 261031-078w. Specify, in the AKM spec and guide, the dependency graph as one graph: explicit `blocked_by` edges plus an implicit edge from every parent to each child (a parent waits for its children). `closable(r)` is true when r has no open children and every child is itself closable; `blocked(r)` is true when a `blocked_by` entry is open or any child is blocked. A cycle in the union graph is an error at any depth. Done when:

- [ ] The spec defines the two predicates and the cycle rule, with fixtures in `docs/spec/akm-api/.fixtures` (and the `testdata/api` copy).
- [ ] WK8 is restated in terms of the predicate; the glossary has the new terms.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
