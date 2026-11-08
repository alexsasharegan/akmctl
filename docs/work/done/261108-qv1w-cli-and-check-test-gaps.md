---
id: "261108-qv1w"
title: cli and check test gaps
status: done
kind: chore
priority: 3
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
---

# cli and check test gaps

## Goal

Findings F15, F18 and F19 of [the audit](../../design/test-strategy-audit.md), from [261030-456q](261030-456q-test-strategy-and-coverage-audit.md): the text view of the CLI (30 % covered, `progressText` 0 %) is not tested, no `check` rule has a near-miss case that must not raise it, and `cmd/akmctl` has no test while some zero-coverage functions may be dead code. Add golden-output tests, near-miss cases and a smoke test of the built binary, and delete or test the dead code.

Done when:

- [x] Golden-output tests cover the text view of `ls`, `show`, `close`, `split` and the commit message block.
- [x] Each `check` rule has a case that must not raise it, next to the case that does.
- [x] A smoke test runs the built binary with `version` and `help`; `migrate.Plan` and `record.SetMetaDir` are tested or deleted.

## Next

## Outcome

Golden files pin the text view of ls, ls --tree, show, board, search, check, close, set, claim and split (-update rewrites them). Every rule of the spec table has a near-miss case (internal/check/near_miss_test.go; SF1, SF3, SF4 point at tests that hold theirs) and a test fails when a rule has none. A smoke test builds cmd/akmctl and runs version, help and an unknown command. record.SetMetaDir was deleted and migrate.Plan folded into PlanEnv. runUI and openBrowser stay untested (no browser in a test).

## Log

- 2026-11-08: filed.
- 2026-11-08 13:43-07:00: claimed (claude-sonnet-5-5-rest).
- 2026-11-08 13:46-07:00: Golden files for the text view (internal/cli/testdata/golden, go test ./internal/cli -run TestTextView -update): ls, ls --tree (progress counts), show, board, search, check, close and set and claim dry runs with the commit block, split. Near misses: internal/check/near_miss_test.go has a boundary case for 38 rules, three (SF1, SF3, SF4) point at tests that already hold theirs, and a test fails when a rule of the spec table has none. cmd/akmctl smoke test builds the binary and runs version, help and an unknown command. Dead code deleted: record.SetMetaDir; migrate.Plan folded into PlanEnv. Not tested: runUI and openBrowser (a browser cannot open in a test).
