---
id: "261108-hwcq"
title: windows powershell harness
status: done
kind: chore
parent: 261108-2ncf
priority: 2
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
---

# windows powershell harness

## Goal

A PowerShell harness for the Windows filesystem. Part of [261108-2ncf](../261108-2ncf-basic-windows-support.md). Windows file behavior (rename over an open file, lock files, path separators, CRLF) is checked only by `go test`; nothing runs the real `akmctl.exe` the way a user does. Add a script (`scripts/windows-smoke.ps1`) and a CI step on the Windows runner that builds `akmctl.exe` and, through `pwsh`, runs `init`, `new`, `set`, `claim`, `close` (which moves a file and rewrites links) and `check` on a temporary repository, with paths written Windows-style and with `core.autocrlf` on and off.

Done when:

- [x] The script runs in CI on `windows-latest` and fails the job on any wrong exit code or output.
- [x] It covers a file move with links, a write while a reader holds the file, and a repository path with spaces and backslashes.
- [x] `docs/architecture.md` says what the harness covers (with [261028-jzsx](../261028-jzsx-windows-path-local-test.md)).

## Next

## Outcome

scripts/windows-smoke.ps1 runs the built akmctl.exe through init, new, set, section, claim, close (link repair) and check on a repository with a spaced path, with core.autocrlf true (records checked out with CRLF) and false, with a mixed-separator --repo and with readers holding a record. CI runs it on windows-latest after the Go tests, all three jobs green. It found 261108-mwtr and 261108-ma45 (known faults, asserted by the harness). Not covered: a console (interactive prompts, colors), see 261108-ct4s.

## Log

- 2026-11-08: filed.
- 2026-11-08 10:39-07:00: the GitHub Windows runner with PowerShell is the chosen environment (maintainer, 2026-11-08).
- 2026-11-08 13:08-07:00: claimed (claude-sonnet-5-5-win).
- 2026-11-08 13:21-07:00: ticked "The script runs in CI on `windows-latest` and fails the job on any wrong exit code or output.": PR 51: the windows smoke step passed on windows-latest and failed there on a wrong flag, a wrong exit code and two real faults while being built
- 2026-11-08 13:21-07:00: ticked "It covers a file move with links, a write while a reader holds the file, and a repository path with spaces and backslashes.": windows-smoke.ps1: close moves and rewrites the link, two reader cases, repo path 'sub dir' with spaces, backslashes and a mixed-separator --repo
- 2026-11-08 13:21-07:00: ticked "`docs/architecture.md` says what the harness covers (with [261028-jzsx](../261028-jzsx-windows-path-local-test.md)).": docs/architecture.md, Platform coverage
- 2026-11-08 13:21-07:00: the run found two faults, filed as 261108-mwtr (write refused while a reader holds the record) and 261108-ma45 (check A4 on a CRLF ADR index); the harness asserts both as known faults and fails when either is fixed
