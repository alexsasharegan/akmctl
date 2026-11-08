---
id: "261029-q8z3"
title: TestCompareAndSwapRace is flaky on Windows CI
summary: The compare-and-swap race test fails now and then on the Windows CI job.
status: done
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-11-08
closed: 2026-11-08
---

# TestCompareAndSwapRace is flaky on Windows CI

## Goal

The compare-and-swap race test fails now and then on the Windows CI job.

Done when:

- [x] The test passes reliably on all three CI platforms, or the cause is recorded.

## Next

## Outcome

The Windows lock was the flaky part: a busy lock gave up after 5 s and a failed release left the file behind. LockRepo now waits 30 s with a random pause and the release retries; a stress run of all Race tests passed on Windows CI. The Windows job no longer has continue-on-error.

## Log

- 2026-10-29: filed.
- 2026-11-08 09:29-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08 09:57-07:00: cause: TestCompareAndSwapRace itself passed 300 of 300 runs on Windows CI, but a stress run of every Race test (go test ./internal/api ./internal/server -run Race -count=100 -cpu 1,2,8) failed in TestWritesAndDryRunsRace with 'open .akm.lock: The file exists'. Windows LockRepo waited 5 s with a fixed 25 ms pause and ignored a failed remove of the lock file. Fix in internal/fsx/lock_other.go: 30 s wait with a 2-15 ms random pause, remove retried for up to 10 s, plus TestLockRepoManyWriters. After the fix the same stress at -cpu 2,8 passed (357 s, run 42501704686). One stress run, not a proof. The temporary stress step is removed and continue-on-error for windows-latest is dropped from ci.yml, as the maintainer asked.
