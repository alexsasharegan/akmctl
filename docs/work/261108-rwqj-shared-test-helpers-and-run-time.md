---
id: "261108-rwqj"
title: shared test helpers and run time
status: active
kind: chore
priority: 3
created: 2026-11-08
updated: 2026-11-08
claimed_by: claude-sonnet-5-5-rest
claimed_at: 2026-11-08T21:05
---

# shared test helpers and run time

## Goal

Findings F11, F12 and F13 of [the audit](../design/test-strategy-audit.md), from [261030-456q](done/261030-456q-test-strategy-and-coverage-audit.md): seven copies of a repository-copy helper that ignore their errors, three of a snapshot helper, a serial `internal/api` run bound by file copies, and a race run dominated by `internal/body`. Add one `internal/testutil`, use it everywhere, run independent tests in parallel and lower the battery iterations under `-race`.

Done when:

- [x] One `internal/testutil` package (`CopyTree` that fails the test on error, `Snapshot`, `MustOK`) replaces the copies.
- [x] `t.Parallel` is on the `internal/api` tests that copy their own repository; the run time before and after is in the Log.
- [x] The `internal/body` race run takes under 10 s, and the Log says how.

## Next

Close after review. Tell the Windows thread to merge main before it adds tests (helpers moved to internal/testutil).

## Log

- 2026-11-08: filed.
- 2026-11-08 14:05-07:00: claimed (claude-sonnet-5-5-rest).
- 2026-11-08 14:06-07:00: internal/testutil now has CopyTree, Snapshot (both fail the test on error) next to RequireTool; the copies in api, cli, server, migrate, check and initrepo use it (14 files, -109 lines); callOK folded into mustOK. MustOK stays in internal/api: it needs the api types, and testutil importing api would be a cycle with api's own tests.
- 2026-11-08 14:13-07:00: Run time, same laptop, -count=1: internal/api 17.8 s -> 11.7 s without -race (24.7 s -> 15.3 s with); 102 tests that copy their own repository call t.Parallel (those that set migrate.Registry, timeNow, migrateRand, the environment or the working directory stay serial, so they finish before the parallel ones resume). internal/body with -race 33.5 s -> 7.2 s: the three generated batteries are parallel and run an eighth of their documents under the race detector (race_test.go and norace_test.go set raceBuild; thresholds scale with the count, full-size runs keep their old floors). Whole go test -race ./... 33.8 s -> 18.2 s.
