---
id: "261108-f1tv"
title: run fixtures through both bindings
status: ready
kind: chore
priority: 3
created: 2026-11-08
updated: 2026-11-08
---

# run fixtures through both bindings

## Goal

Findings F1 and F2 of [the audit](../design/test-strategy-audit.md), from [261030-456q](done/261030-456q-test-strategy-and-coverage-audit.md): the conformance fixtures run only through `api.Run`, so nothing proves the CLI and the HTTP server return the same envelope, and `new`, `commits`, `history` and `skills` have no fixture. Run every case through the HTTP handler and the reads through the CLI, and draft the ADR on fixing the clock, ids and git history for the rest.

Done when:

- [ ] Every case file runs through the HTTP handler and compares with `expected`; the reads also run through the CLI with `--format json`.
- [ ] An ADR draft (`proposed`) says how to fixture `new`, `commits`, `history` and `skills`.
- [ ] The first fixtures that ADR allows are added, or filed as items.

## Next

Run the existing cases through the HTTP handler (F1); draft the ADR.

## Log

- 2026-11-08: filed.
