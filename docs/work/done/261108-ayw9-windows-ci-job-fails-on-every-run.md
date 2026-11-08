---
id: "261108-ayw9"
title: windows ci job fails on every run
status: done
kind: bug
priority: 2
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
---

# windows ci job fails on every run

## Goal

The `test (windows-latest)` job of CI fails on every run (it is `continue-on-error`, so the run still shows success), which hides a new Windows regression among the known failures (audit finding F7, [261108-z9n8](../261108-z9n8-ci-hardening-for-skipped-tests-and-scripts.md)). The hold on Windows work was lifted on 2026-11-08. The latest `main` run (16:08 UTC) shows two failing tests:

- `TestCustomDocsDirEndToEnd` (`internal/api/docsdir_test.go`): "the lock is not in the docs dir". On Windows the lock is created exclusively and removed on release (`internal/fsx/lock_other.go`), so the file is not there after the write; the test assumes the Unix advisory lock that stays on disk.
- `TestRepoInsideARepositoryIsRefused` (`internal/cli/cli_test.go`): compares the JSON output with an unescaped Windows path; JSON writes `\` for each backslash.

Done when:

- [x] Both tests pass on Windows CI (a run of the branch shows them passing).
- [x] The whole `go test ./...` step of the `test (windows-latest)` job is green, or each remaining failure is filed.
- [x] [261029-q8z3](../261029-q8z3-cas-race-windows.md) (flaky `TestCompareAndSwapRace`) is checked on the green runs and noted.
- [x] Whether the job can stop being `continue-on-error` is recorded in the Log (the decision is the maintainer's).

## Next

## Outcome

Two Windows-only test faults fixed: TestCustomDocsDirEndToEnd no longer expects the Unix advisory lock file to remain on Windows, and TestRepoInsideARepositoryIsRefused compares the JSON-escaped path. The Windows job passed on PR 44. continue-on-error is left on; q8z3 stays open.

## Log

- 2026-11-08: filed.
- 2026-11-08 09:11-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08 09:14-07:00: branch CI (head of PR 44): test (windows-latest) passed on both the push and the pull request run; go test ./... printed no FAIL line and every package ok. TestCompareAndSwapRace (261029-q8z3) passed in both runs, which is two data points, not proof; q8z3 stays open. continue-on-error for windows-latest in .github/workflows/ci.yml is left on: dropping it is the maintainer's call, suggested after several green runs on main.
