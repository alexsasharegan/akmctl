---
id: "261108-f1tv"
title: run fixtures through both bindings
status: active
kind: chore
priority: 3
created: 2026-11-08
updated: 2026-11-08
claimed_by: claude-sonnet-5-5-rest
claimed_at: 2026-11-08T21:26
decisions: [261108-15dw]
---

# run fixtures through both bindings

## Goal

Findings F1 and F2 of [the audit](../design/test-strategy-audit.md), from [261030-456q](done/261030-456q-test-strategy-and-coverage-audit.md): the conformance fixtures run only through `api.Run`, so nothing proves the CLI and the HTTP server return the same envelope, and `new`, `commits`, `history` and `skills` have no fixture. Run every case through the HTTP handler and the reads through the CLI, and draft the ADR on fixing the clock, ids and git history for the rest.

Done when:

- [ ] Every case file runs through the HTTP handler and compares with `expected`; the reads also run through the CLI with `--format json`.
- [x] An ADR draft (`proposed`) says how to fixture `new`, `commits`, `history` and `skills`.
- [ ] The first fixtures that ADR allows are added, or filed as items.

## Next

Wait for the maintainer to accept, change or reject ADR 261108-15dw. Then: matchers and bindings in TestConformance (F1), fixtures for new, commits, history and skills (F2), -update keeps matchers.

## Log

- 2026-11-08: filed.
- 2026-11-08 14:26-07:00: claimed (claude-sonnet-5-5-rest).
- 2026-11-08 14:27-07:00: Drafted ADR 261108-15dw (proposed): matchers ($id, $hash, $any), a replayed git history in the case, one runner for api.Run, the HTTP handler and the CLI reads. Stopped here for the maintainer's decision before building the runners or the new fixtures.
