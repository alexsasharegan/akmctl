---
id: "261109-e540"
title: CI steps that download from GitHub fail on a 504 and block the PR run
status: ready
kind: bug
created: 2026-11-09
updated: 2026-11-09
---

# CI steps that download from GitHub fail on a 504 and block the PR run

## Goal

Three CI steps download from GitHub or the winget feed and fail the whole `test` job when the server answers 504: `winget manifests` (the winget client, `0x801901f8 Gateway timeout (504)`), `homebrew formula` (`curl: (22) ... 504`) and `windows zip install (README steps)` (`Invoke-WebRequest ... 504`, then a checksum mismatch on the empty `SHA256SUMS`). On 2026-11-09 they failed pull request runs of 261108-hp9k and the theming stack; a rerun of the push run passed, and two reruns of the pull request run did not.

Command: `ci` workflow on a pull request. Actual: the job fails on an outage that has nothing to do with the change, and the pull request stays `UNSTABLE`. Expected: a transient download error is retried a few times, and if the service stays down the step is reported as skipped (not failed) while the rest of the job's result counts.

Done when:

- [ ] Each of the three steps retries its download (at least three tries with a pause) before it fails.
- [ ] A download that still fails after the retries fails the step only if the maintainer has said these steps must gate a merge; otherwise it is a warning annotation. Decide and record which in the Log.
- [ ] The step that checks a damaged download still fails on a real checksum mismatch.

## Next

Add the retry to the three steps in .github/workflows/ci.yml.

## Log

- 2026-11-09: filed.
- 2026-11-09 01:23-07:00: Filed after the 504s on 2026-11-09 (winget manifests, homebrew formula, windows zip install). Same commit passed on the push run.
