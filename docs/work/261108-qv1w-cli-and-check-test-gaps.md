---
id: "261108-qv1w"
title: cli and check test gaps
status: active
kind: chore
priority: 3
created: 2026-11-08
updated: 2026-11-08
claimed_by: claude-sonnet-5-5-rest
claimed_at: 2026-11-08T20:43
---

# cli and check test gaps

## Goal

Findings F15, F18 and F19 of [the audit](../design/test-strategy-audit.md), from [261030-456q](done/261030-456q-test-strategy-and-coverage-audit.md): the text view of the CLI (30 % covered, `progressText` 0 %) is not tested, no `check` rule has a near-miss case that must not raise it, and `cmd/akmctl` has no test while some zero-coverage functions may be dead code. Add golden-output tests, near-miss cases and a smoke test of the built binary, and delete or test the dead code.

Done when:

- [ ] Golden-output tests cover the text view of `ls`, `show`, `close`, `split` and the commit message block.
- [ ] Each `check` rule has a case that must not raise it, next to the case that does.
- [ ] A smoke test runs the built binary with `version` and `help`; `migrate.Plan` and `record.SetMetaDir` are tested or deleted.

## Next

Start with golden-output tests for `ls` and `show` (F15).

## Log

- 2026-11-08: filed.
- 2026-11-08 13:43-07:00: claimed (claude-sonnet-5-5-rest).
