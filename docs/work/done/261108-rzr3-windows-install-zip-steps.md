---
id: "261108-rzr3"
title: windows install zip steps
status: done
kind: chore
parent: 261108-2ncf
priority: 2
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
---

# windows install zip steps

## Goal

Documented zip steps for Windows. Part of [261108-2ncf](../261108-2ncf-basic-windows-support.md). `scripts/install.sh` stops on Windows with a pointer to the releases page, and the README has one generic line. Write the steps for a Windows user: download the `.zip`, check `SHA256SUMS` (`Get-FileHash`), unpack, put the folder on `PATH` (`[Environment]::SetEnvironmentVariable`), and what to do about the SmartScreen warning for an unsigned binary. Decide whether a `scripts/install.ps1` is worth adding once WinGet exists.

Done when:

- [x] The README install section has a Windows block that a user can follow in PowerShell.
- [x] `scripts/install.sh` points at it instead of only at the releases page.
- [x] The steps were run on a Windows machine (the CI runner is enough).

## Next

## Outcome

The README has a Windows section with PowerShell steps (download the zip, Get-FileHash against SHA256SUMS, install to %LOCALAPPDATA%\Programs\akmctl, user PATH, SmartScreen). CI extracts and runs that block on windows-latest against v0.3.3 and checks that a bad checksum stops it. install.sh points at the section. No install.ps1: the block is short and tested, WinGet is the planned one-command route. Not tested: a browser-downloaded zip with the SmartScreen prompt (needs a person), and arm64.

## Log

- 2026-11-08: filed.
- 2026-11-08 14:09-07:00: claimed (claude-sonnet-5-5-win).
- 2026-11-08 14:16-07:00: ticked "The README install section has a Windows block that a user can follow in PowerShell.": README.md, section Windows: a PowerShell block with Get-FileHash, PATH and SmartScreen notes
- 2026-11-08 14:16-07:00: ticked "`scripts/install.sh` points at it instead of only at the releases page.": scripts/install.sh now points at the README Windows section
- 2026-11-08 14:16-07:00: ticked "The steps were run on a Windows machine (the CI runner is enough).": CI step 'windows zip install (README steps)' runs the README block itself on windows-latest: installs v0.3.3, akmctl version works, and a forced bad checksum is refused
- 2026-11-08 14:16-07:00: decision: no scripts/install.ps1 for now. The README block is 17 lines, CI runs it as written, and WinGet (manifests ready, submission is the maintainer's) will be the one-command route; a second script would be a second thing to keep in step with the release. Revisit if WinGet is refused or people ask.
