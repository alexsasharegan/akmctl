---
id: "261108-31p6"
title: "Windows is supported: what passes there is committed to, and a feature that cannot is switched off by name"
summary: "Proposed: what passes Windows CI is committed to; new features are tested there by default and held if they fail; one that cannot work is switched off by name."
status: proposed
date: 2026-11-08
supersedes: 261029-eegx
scope: [.github/workflows/ci.yml, scripts/windows-smoke.ps1, README.md, AGENTS.md, docs/release.md]
authored_by: agent
---

# Windows is supported: what passes there is committed to, and a feature that cannot is switched off by name

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
3. **Windows is supported for the command line, scoped by what is tested.** The tested behavior is a requirement and the rest carries no promise. This was the first draft. Rejected: it ties the promise to whichever tests happen to exist, so a gap in testing becomes a gap in the promise without anyone deciding it.
4. **Windows is the default for every feature, and an exception is a named decision** (recommended). What passes on Windows today is committed to. A new feature is tested there from the start and is held if it fails. A feature the maintainer decides cannot work there is switched off on Windows, listed by name, and the rest of the tool still works.
5. **Drop Windows.** Stop building the zip and the CI job. Rejected: the code is portable, the work is done and a user is coming.

## Decision

Recommended (the maintainer decides): option 4, in four terms.

1. **What is landed and passes on Windows is committed to.** Every feature that passes the Windows CI job today is supported. The job blocks merging and a release, as the macOS and Ubuntu jobs do; a regression is fixed, and a red Windows job is never merged. The job is the evidence: `go vet`, `go test`, the PowerShell smoke harness (`init`, `new`, `section`, `set`, `claim`, `close`, `check`, `show`, `ls` and `find`, on a path with spaces, with `--repo` in backslash, forward-slash, mixed and relative forms, on CRLF checkouts), the README install block, the first-run checklist and the WinGet manifest check.
2. **A new feature is tested on Windows by default.** Its tests run in the Windows job from the first pull request. Nobody opts Windows in; leaving it out is the exception in term 4. The happy path is a feature that passes on all three platforms and is then supported on all three.
3. **A new feature that fails on Windows is held.** It does not merge, and it does not drop out of the Windows tests. Either it is made to work, or the maintainer decides under term 4. Best effort here means an honest attempt to make it work, not an automatic pass.
4. **A feature the maintainer decides cannot work on Windows is switched off there, by name.** The tool still works on Windows and the rest is untouched; the user just does not get that feature. Three things keep the exception visible:
   - **One list, in one place.** The README keeps a table of the features Windows does not have, each with the work item that recorded the decision and its reason. A feature is on the list only by a decision recorded in a work item or an ADR, never by a test that was left out.
   - **A plain refusal.** On Windows the command says it is not available there and why, with an exit code distinct from a failure, and does not crash or hang. The code is chosen when the first such feature lands.
   - **A named skip.** The feature's tests are skipped on Windows by name, checked against the list, so a test cannot disappear from Windows without the list saying so. (The test job sets `AKM_REQUIRE_TOOLS=1` for the same reason: a skip must be declared, not silent.)

**Supported platforms.** macOS and Linux (as before), and Windows on amd64 with the binary from the release zip or built from source, in PowerShell. The Windows version is the one on the runner (`windows-latest`); older ones are not tested.

**Not committed to yet.** These are the first entries of the unsupported list until a Windows check says otherwise; none blocks a release until it is tested:

- The console: the `init` questionnaire, colors and `akmctl ui` are checked on the runner without a terminal (redirected input, the page served), but not on a real console (261108-ct4s).
- Windows arm64 (no archive is built), UNC paths and drive-letter forms other than the runner's.
- Known faults while open: a write refused while another program holds the record (261108-mwtr), `check` A4 on a CRLF ADR index (261108-ma45), new files written as LF (261108-2z9m), backslash links not checked (261108-sfj3). Fixing them is wanted, not gating.
- A signed binary: the zip is unsigned, so SmartScreen can warn; the README says how to proceed after checking the checksum. Signing is not planned before 1.0.0.
- The write commands the harness does not run (`release`, `log`, `tick`, `remove`, `split`, the migrate steps) are covered by the Go tests on all three platforms and by the CRLF audit; they are expected to work and are not separately gated.

**Package managers.** WinGet and Homebrew are published by the maintainer. Their manifests are checked in CI; a submission is not part of a release's definition of done.

**Documents that change on acceptance.** README line 7 and the Windows bullet near the end (and the new list of features Windows does not have), the `AGENTS.md` rule that calls Windows best effort, and `docs/release.md` (the release steps name the Windows job). 261029-eegx is set to `superseded` with `superseded_by: 261108-31p6` in the same commit. Work item 261108-wg98 owns these edits.

## Consequences

- A Windows regression stops a merge, so contributors keep the Windows job green; the harness and CI steps show why it failed without a Windows machine.
- Windows is part of the design of every feature from the start, which is cheaper than finding out late. The cost is ongoing: the PowerShell harness, the install block and the checklist must be kept in step with the CLI. They run in the same job, so a drift fails loudly.
- A feature can be good and still not ship on Windows. That is the maintainer's decision, recorded, and visible to users in one table.
- Switching a feature off by platform is new to the tool. It needs a small mechanism (a build constraint or a runtime check, the refusal and the named skip) and must not weaken the API's contract: the API spec says operations are available, so an operation that is unavailable on a platform has to be declared in the spec (for example in the `version` response) or the fixtures and clients will disagree about what exists. The first feature that needs this files that change; this ADR does not decide the form.
- The promise is wide for what passes and explicit about what does not, so an unsupported area is an open item or a listed exception, not a broken promise.
- The `windows-latest` image changes under the project. A red job that no commit caused is possible; the response is to pin or fix, not to turn the job off.
- For Windows, 1.0.0 needs the console features checked and the known faults fixed, or each listed as an exception by a decision. That is a gate already implied here, not a new decision.

## Confirmation

`grep -rn continue-on-error .github/workflows` prints nothing. Once the list exists, a test fails if a feature's Windows tests are skipped and the feature is not on it.

## Review

When the first Windows user has run the checklist (261108-7avk), when the first feature goes on the unsupported list, and again before 1.0.0.

## Revert

Option 1 is a small change: mark the Windows job `continue-on-error: true` and restore the old README and `AGENTS.md` wording. Nothing else depends on this decision. The first draft's scoping (option 3) is a smaller commitment if term 4 proves too heavy.
