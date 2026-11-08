---
id: "261108-ct4s"
title: windows interactive commands check
status: active
kind: chore
parent: 261108-2ncf
priority: 4
created: 2026-11-08
updated: 2026-11-08
claimed_by: claude-sonnet-5-5-win
claimed_at: 2026-11-08T22:08
---

# windows interactive commands check

## Goal

A check of the interactive commands on Windows. Part of [261108-2ncf](261108-2ncf-basic-windows-support.md). `init`'s questionnaire (`internal/cli/prompt.go`), colours and `akmctl ui` (opens the browser with `rundll32`) were never run on a Windows terminal. Run them once in Windows Terminal and PowerShell, record what works, and file the faults. The workbench is lower value than the filesystem items.

Done when:

- [ ] `init` with the questionnaire completes in Windows Terminal and in plain PowerShell.
- [ ] `akmctl ui` opens the page and serves it.
- [ ] Faults found are filed as items.

## Next

Write scripts/windows-interactive.ps1 (init with NUL and piped stdin, ui served and fetched over HTTP, ui --open starts), run it in CI, record what the runner cannot do (a real console, Windows Terminal, the browser opening), file faults.

## Log

- 2026-11-08: filed.
- 2026-11-08 15:08-07:00: claimed (claude-sonnet-5-5-win).
