---
id: "261108-hwcq"
title: windows powershell harness
status: raw
kind: chore
parent: 261108-2ncf
priority: 2
created: 2026-11-08
updated: 2026-11-08
---

# windows powershell harness

## Goal

A PowerShell harness for the Windows filesystem. Part of [261108-2ncf](261108-2ncf-basic-windows-support.md). Windows file behavior (rename over an open file, lock files, path separators, CRLF) is checked only by `go test`; nothing runs the real `akmctl.exe` the way a user does. Add a script (`scripts/windows-smoke.ps1`) and a CI step on the Windows runner that builds `akmctl.exe` and, through `pwsh`, runs `init`, `new`, `set`, `claim`, `close` (which moves a file and rewrites links) and `check` on a temporary repository, with paths written Windows-style and with `core.autocrlf` on and off.

Done when:

- [ ] The script runs in CI on `windows-latest` and fails the job on any wrong exit code or output.
- [ ] It covers a file move with links, a write while a reader holds the file, and a repository path with spaces and backslashes.
- [ ] `docs/architecture.md` says what the harness covers (with [261028-jzsx](261028-jzsx-windows-path-local-test.md)).

## Log

- 2026-11-08: filed.
