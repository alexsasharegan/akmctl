---
id: "261108-7avk"
title: windows first run checklist
status: active
kind: chore
parent: 261108-2ncf
priority: 3
created: 2026-11-08
updated: 2026-11-08
claimed_by: claude-sonnet-5-5-win
claimed_at: 2026-11-08T21:32
---

# windows first run checklist

## Goal

A first-run checklist for the first Windows user. Part of [261108-2ncf](261108-2ncf-basic-windows-support.md). A short list the user (or an agent on their machine) runs once, in PowerShell: install, `akmctl version`, `akmctl init` in a new folder, `new` / `claim` / `close`, `check`, and what to send back if one fails (the command, the output, `$PSVersionTable`). Put it in `docs/` and link it from the README.

Done when:

- [ ] The checklist is in `docs/` and the first user has run it.
- [ ] Each failure it found is filed as an item.

## Next

Waiting for the first Windows user to run docs/windows-first-run.md (maintainer arranges). Then tick Goal box 1, file each failure under 261108-2ncf (box 2) and close. CI already runs the block on the runner.

## Log

- 2026-11-08: filed.
- 2026-11-08 14:32-07:00: claimed (claude-sonnet-5-5-win).
- 2026-11-08 14:38-07:00: docs/windows-first-run.md written and linked from the README; CI runs its command block on windows-latest (green, no failures found). Open: a person runs it by hand
