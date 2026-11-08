---
id: "261108-2ncf"
title: basic windows support
status: ready
kind: feature
priority: 2
created: 2026-11-08
updated: 2026-11-08
---

# basic windows support

## Goal

A first Windows user is coming. Windows was best effort ([ADR 261029-eegx](../adr/261029-eegx-platform-support-tiers.md)); this parent takes it to a basic supported level: a user can install `akmctl` the usual way, and the file behavior that differs on Windows (paths, line endings, locks, rename) is covered by tests that run in PowerShell on real Windows. The CI Windows job already counts ([261029-q8z3](done/261029-q8z3-cas-race-windows.md), [261108-ayw9](done/261108-ayw9-windows-ci-job-fails-on-every-run.md)).

Not in scope: self-update (once WinGet and Homebrew install the tool, they update it), and polish of the interactive commands beyond a check.

Done when:

- [ ] Every child is `done` or `dropped`.
- [ ] A Windows user installs with `winget install` (or the documented zip steps), and a macOS or Linux user with `brew install`.
- [ ] A PowerShell harness on the Windows runner exercises the real `akmctl.exe` on a repository.
- [ ] A superseding ADR states the support tier, and the README says the same.

## Next

Take the children in this order: the PowerShell harness, the CRLF and path audit, the installers (zip steps, Homebrew, WinGet), then the ADR and the first-run checklist.

## Log

- 2026-11-08: filed.
