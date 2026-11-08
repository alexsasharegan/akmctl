---
id: "261108-7avk"
title: windows first run checklist
status: raw
kind: chore
parent: 261108-2ncf
priority: 3
created: 2026-11-08
updated: 2026-11-08
---

# windows first run checklist

## Goal

A first-run checklist for the first Windows user. Part of [261108-2ncf](261108-2ncf-basic-windows-support.md). A short list the user (or an agent on their machine) runs once, in PowerShell: install, `akmctl version`, `akmctl init` in a new folder, `new` / `claim` / `close`, `check`, and what to send back if one fails (the command, the output, `$PSVersionTable`). Put it in `docs/` and link it from the README.

Done when:

- [ ] The checklist is in `docs/` and the first user has run it.
- [ ] Each failure it found is filed as an item.

## Log

- 2026-11-08: filed.
