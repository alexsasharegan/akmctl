---
id: "261028-jzsx"
title: Cover the Windows path beyond CI
summary: Windows behaviour (paths, line endings, locks, atomic rename) is tested only by the CI runner; add a way to run or simulate it locally.
status: ready
kind: chore
priority: 1
created: 2026-10-28
updated: 2026-11-08
parent: 261108-2ncf
---

# Cover the Windows path beyond CI

## Goal

The Windows behaviour (path separators, CRLF, file locks, rename over an existing file) is exercised only by the Windows job in CI. A change cannot be tried on Windows before it is pushed.

Done when:

- [ ] A documented way to run the Windows-sensitive tests locally, or tests that simulate the differences (path separators and line endings) on any platform.
- [ ] The CI Windows job stays the final check.
- [ ] A note in `docs/architecture.md` lists what is covered where.

## Log

- 2026-10-28: filed.
- 2026-11-08 10:36-07:00: options for running Windows locally (2026-11-08): no Windows containers on macOS (they need a Windows kernel); a Windows VM (UTM, Parallels, Fusion) is the high-fidelity option; Wine is quick but emulates the filesystem and will not show lock or rename behavior; the GitHub Windows runner with a PowerShell script is real and is the planned harness (261108-hwcq). Now a child of 261108-2ncf.
- 2026-11-08 10:39-07:00: decision (maintainer, 2026-11-08): the GitHub Windows runner with PowerShell is the Windows environment; no VM or Wine. Local simulation covers only what can be simulated on any platform (path style, CRLF); real behavior is checked by the PowerShell harness 261108-hwcq in CI.
