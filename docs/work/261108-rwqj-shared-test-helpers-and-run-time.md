---
id: "261108-rwqj"
title: shared test helpers and run time
status: ready
kind: chore
priority: 3
created: 2026-11-08
updated: 2026-11-08
---

# shared test helpers and run time

## Goal

Findings F11, F12 and F13 of [the audit](../design/test-strategy-audit.md), from [261030-456q](done/261030-456q-test-strategy-and-coverage-audit.md): seven copies of a repository-copy helper that ignore their errors, three of a snapshot helper, a serial `internal/api` run bound by file copies, and a race run dominated by `internal/body`. Add one `internal/testutil`, use it everywhere, run independent tests in parallel and lower the battery iterations under `-race`.

Done when:

- [ ] One `internal/testutil` package (`CopyTree` that fails the test on error, `Snapshot`, `MustOK`) replaces the copies.
- [ ] `t.Parallel` is on the `internal/api` tests that copy their own repository; the run time before and after is in the Log.
- [ ] The `internal/body` race run takes under 10 s, and the Log says how.

## Next

Write `internal/testutil` and move one package onto it (F11).

## Log

- 2026-11-08: filed.
