---
id: "261108-rzr3"
title: windows install zip steps
status: raw
kind: chore
parent: 261108-2ncf
priority: 2
created: 2026-11-08
updated: 2026-11-08
---

# windows install zip steps

## Goal

Documented zip steps for Windows. Part of [261108-2ncf](261108-2ncf-basic-windows-support.md). `scripts/install.sh` stops on Windows with a pointer to the releases page, and the README has one generic line. Write the steps for a Windows user: download the `.zip`, check `SHA256SUMS` (`Get-FileHash`), unpack, put the folder on `PATH` (`[Environment]::SetEnvironmentVariable`), and what to do about the SmartScreen warning for an unsigned binary. Decide whether a `scripts/install.ps1` is worth adding once WinGet exists.

Done when:

- [ ] The README install section has a Windows block that a user can follow in PowerShell.
- [ ] `scripts/install.sh` points at it instead of only at the releases page.
- [ ] The steps were run on a Windows machine (the CI runner is enough).

## Log

- 2026-11-08: filed.
