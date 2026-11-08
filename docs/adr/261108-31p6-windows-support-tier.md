---
id: "261108-31p6"
title: "Windows is supported for the command line: what is tested blocks, what is not tested is not promised"
summary: "Proposed: Windows (amd64, command line) is supported; what the Windows job tests blocks merges and releases, the console and untested areas carry no promise."
status: proposed
date: 2026-11-08
supersedes: 261029-eegx
scope: [.github/workflows/ci.yml, scripts/windows-smoke.ps1, README.md, AGENTS.md, docs/release.md]
authored_by: agent
---

# Windows is supported for the command line: what is tested blocks, what is not tested is not promised

## Context

ADR 261029-eegx made Windows best effort: a Windows failure did not block merging or a release below 1.0.0. Since then the facts changed (work item 261108-2ncf):

- The Windows CI job blocks merging, as the macOS and Ubuntu jobs do (261029-q8z3). Nothing in the workflow lets a red Windows job through (`continue-on-error` appears nowhere).
- Real Windows behavior is checked, not only compiled: the PowerShell smoke harness (`scripts/windows-smoke.ps1`) runs the built `akmctl.exe` through init, new, section, set, claim, close, check, show, ls and find on the runner, and the README's install block and the first-run checklist run there too (261108-hwcq, 261108-rzr3, 261108-7avk).
- The audit (261108-kbve) found faults, filed as 261108-mwtr, 261108-ma45, 261108-2z9m and 261108-sfj3. They are known and small, and none corrupts a record.
- Release archives for Windows are built (`akmctl_<version>_windows_amd64.zip`) and WinGet manifests are generated and checked (261108-fabf). A first Windows user is expected.

So the written rule (best effort, non-blocking) no longer matches what the project does, and the README says "best effort" next to a zip install guide and a checklist. The rule needs to say what Windows support is now, including what it does not cover.

## Considered Options

1. **Keep 261029-eegx.** Windows stays best effort and non-blocking. The CI job and the docs are then out of step with the rule, unless the maintainer also turns the Windows job back to non-blocking. Rejected as it stands.
2. **Windows is a full peer of macOS and Linux.** Every feature, including the interactive ones (the `init` questionnaire, colors, `akmctl ui`), is required to work and is gated. Not possible today: nothing checks the console features on a real console (261108-ct4s), and a VM or physical machine is not part of CI (decision of 2026-11-08: the GitHub Windows runner with PowerShell is the environment).
3. **Windows is supported for the command line, scoped by what is tested** (recommended). The tested behavior is a requirement: a regression in it blocks a merge and a release. Behavior the Windows job does not exercise carries no promise, and the README says which that is.
4. **Drop Windows.** Stop building the zip and the CI job. Rejected: the code is portable, the work is done and a user is coming.

## Decision

Recommended (the maintainer decides): option 3.

**Supported platforms.** macOS and Linux (as before), and Windows on amd64 with the binary from the release zip or built from source, in PowerShell.

**What blocks.** The Windows CI job (`go vet`, `go test`, the smoke harness, the README install block, the first-run checklist, the WinGet manifest check) blocks merging and a release, as the other jobs do. A Windows failure is fixed, or the test is marked with a work item that says why; it is not merged red.

**What is promised.** What those checks run on the Windows runner: `init`, `new`, `section`, `set`, `claim`, `close`, `check`, `show`, `ls` and `find`, on a path with spaces, with `--repo` in backslash, forward-slash, mixed and relative forms, and on CRLF checkouts (every edit of an existing record keeps its line endings). The other write commands (`release`, `log`, `tick`, `remove`, `split`, the migrate steps) are covered by the Go tests on all three platforms and by the CRLF audit, but not by the harness; they are expected to work and are not separately gated.

**What is not promised yet.** Each of these has a work item, and none blocks a release until it is tested:

- The console: the `init` questionnaire, colors and `akmctl ui` (261108-ct4s).
- Windows arm64 (no archive is built), UNC paths and drive-letter forms other than the runner's.
- Known faults while open: a write refused while another program holds the record (261108-mwtr), `check` A4 on a CRLF ADR index (261108-ma45), new files written as LF (261108-2z9m), backslash links not checked (261108-sfj3). Fixing them is wanted, not gating.
- A signed binary: the zip is unsigned, so SmartScreen can warn; the README says how to proceed after checking the checksum. Signing is not planned before 1.0.0.

**Package managers.** WinGet and Homebrew are published by the maintainer. Their manifests are checked in CI; a submission is not part of a release's definition of done.

**Versions.** The Windows version is the one on the runner (`windows-latest`). Older ones are not tested.

**Documents that change on acceptance.** README line 7 and the Windows bullet near the end, the `AGENTS.md` rule that calls Windows best effort, and `docs/release.md` (the release steps name the Windows job). 261029-eegx is set to `superseded` with `superseded_by: 261108-31p6` in the same commit.

## Consequences

- A Windows regression stops a merge, so contributors keep the Windows job green; the harness and the CI steps show why it failed without a Windows machine.
- The cost of Windows work is ongoing: the PowerShell harness, the install block and the checklist must be kept in step with the CLI. They run in the same job, so a drift fails loudly.
- The promise is narrow and written down, so an unsupported area (the console, arm64) is an open item, not a broken promise.
- The `windows-latest` image changes under the project. A red job that no commit caused is possible; the response is to pin or fix, not to turn the job off.
- For Windows, 1.0.0 needs the console features checked and the known faults fixed. That is a gate already implied here, not a new decision.

## Confirmation

`grep -rn continue-on-error .github/workflows` prints nothing.

## Review

When the first Windows user has run the checklist (261108-7avk), and again before 1.0.0.

## Revert

Option 1 is a small change: mark the Windows job `continue-on-error: true` and restore the old README and `AGENTS.md` wording. Nothing else depends on this decision.
