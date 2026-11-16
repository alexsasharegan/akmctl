---
id: "261109-e540"
title: CI steps that download from GitHub fail on a 504 and block the PR run
status: done
kind: bug
created: 2026-11-09
updated: 2026-11-16
closed: 2026-11-16
---

# CI steps that download from GitHub fail on a 504 and block the PR run

## Goal

Three CI steps download from GitHub or the winget feed and fail the whole `test` job when the server answers 504: `winget manifests` (the winget client, `0x801901f8 Gateway timeout (504)`), `homebrew formula` (`curl: (22) ... 504`) and `windows zip install (README steps)` (`Invoke-WebRequest ... 504`, then a checksum mismatch on the empty `SHA256SUMS`). On 2026-11-09 they failed pull request runs of 261108-hp9k and the theming stack; a rerun of the push run passed, and two reruns of the pull request run did not.

Command: `ci` workflow on a pull request. Actual: the job fails on an outage that has nothing to do with the change, and the pull request stays `UNSTABLE`. Expected: a transient download error is retried a few times, and if the service stays down the step is reported as skipped (not failed) while the rest of the job's result counts.

Done when:

- [x] Each of the three steps retries its download (at least three tries with a pause) before it fails.
- [x] A download that still fails after the retries fails the step only if the maintainer has said these steps must gate a merge; otherwise it is a warning annotation. Decide and record which in the Log.
- [x] The step that checks a damaged download still fails on a real checksum mismatch.

## Next

## Outcome

- Each of the three steps retries: `scripts/ci-retry.sh` (`ci_retry`) for the Homebrew step and `scripts/ci-retry.ps1` (`Invoke-CiRetry`) for the two Windows steps, three tries with a growing pause.
- A download that still fails is a warning annotation and the step exits 0; any failure that is not a transient download error fails the step on the first try.
- The damaged-download case now has to fail with `checksum mismatch` in its output, so a 504 cannot pass for a refusal.

Not verified locally: the PowerShell helper and the workflow edits (no `pwsh` here). `TestCIRetry` covers the shell helper. The pull request's own CI run is the first real run.

Decision: none (warn or gate is a pick, in the Log).

## Log

- 2026-11-09: filed.
- 2026-11-09 01:23-07:00: Filed after the 504s on 2026-11-09 (winget manifests, homebrew formula, windows zip install). Same commit passed on the push run.
- 2026-11-16 07:01-07:00: claimed (claude-fable-5-2).
- 2026-11-16 07:02-07:00: ticked "Each of the three steps retries its download (at least three tries with a pause) before it fails.": ci.yml: Invoke-CiRetry and ci_retry around the README block, gh release download, winget install and brew install; three tries, pauses of 15 and 30 seconds
- 2026-11-16 07:02-07:00: ticked "A download that still fails after the retries fails the step only if the maintainer has said these steps must gate a merge; otherwise it is a warning annotation. Decide and record which in the Log.": pick: a warning annotation, not a gate (Exit-CiSkipOrFail, ci_skip_or_fail); recorded in the Log
- 2026-11-16 07:02-07:00: ticked "The step that checks a damaged download still fails on a real checksum mismatch.": ci.yml: the damaged-download case must fail with 'checksum mismatch' in its output; TestCIRetry: a mismatch is not retried and fails
- 2026-11-16 07:02-07:00: pick: after the retries a transient download error is a warning annotation and the step exits 0, it does not gate the merge; the maintainer has not said these steps must gate. Transient is decided from the output (HTTP 500 to 504, timeouts, resets, the winget HRESULT 0x8019xxxx)
- 2026-11-16 07:02-07:00: not verified: the PowerShell helper and the edited workflow steps have not run. There is no pwsh or actionlint on this machine; the shell helper is tested (TestCIRetry) and the YAML parses. The first CI run of the pull request is the check
