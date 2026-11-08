---
id: "261108-ct4s"
title: windows interactive commands check
status: ready
kind: chore
parent: 261108-2ncf
priority: 4
created: 2026-11-08
updated: 2026-11-08
---

# windows interactive commands check

## Goal

A check of the interactive commands on Windows. Part of [261108-2ncf](261108-2ncf-basic-windows-support.md). `init`'s questionnaire (`internal/cli/prompt.go`), colours and `akmctl ui` (opens the browser with `rundll32`) were never run on a Windows terminal. Run them once in Windows Terminal and PowerShell, record what works, and file the faults. The workbench is lower value than the filesystem items.

Done when:

- [ ] `init` with the questionnaire completes in Windows Terminal and in plain PowerShell.
- [ ] `akmctl ui` opens the page and serves it.
- [ ] Faults found are filed as items.

## Next

Waiting for the first Windows user to do step 3 of docs/windows-first-run.md (init questions in Windows Terminal and PowerShell, ui --open, Ctrl+C). Then tick Goal boxes 1 and 2, file any fault (box 3) and close. The runner part is done and runs in CI.

## Log

- 2026-11-08: filed.
- 2026-11-08 15:08-07:00: claimed (claude-sonnet-5-5-win).
- 2026-11-08 15:13-07:00: runner check (scripts/windows-interactive.ps1, green): init finishes with stdin from NUL, an empty file and an answers file (stdout redirected, so nothing is asked); ui and ui --open serve the page on 127.0.0.1 and the API answers 401 without the token. No faults found. The tool emits no ANSI colors. Open, for a person: init's questions in Windows Terminal and PowerShell, the browser opening, Ctrl+C (now steps in docs/windows-first-run.md)
- 2026-11-08 15:13-07:00: released (claude-sonnet-5-5-win).
