---
id: "261119-3aeq"
title: "Workbench: yank a record's absolute path to the clipboard"
status: active
kind: feature
created: 2026-11-19
updated: 2026-11-19
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-20T02:32
---

# Workbench: yank a record's absolute path to the clipboard

## Goal

In the workbench (`akmctl ui`) `Y` copies the absolute path of the focused or open record. Done when `Y` copies the repository root joined with the record path (the page already gets the root from `/ui/config`, so no API change), shows a toast, says so when the root is unknown or the clipboard refuses, and appears in the keys help, the hints and the palette. GitHub issue 97; relative paths are left out on purpose. Builds on the id yank.

## Next

Try Y in a real browser; close the item after the pull request merges.

## Log

- 2026-11-19: filed.
- 2026-11-19 19:32-07:00: claimed (claude-sonnet-5-5).
- 2026-11-19 19:33-07:00: Corrected the Goal: the root is already on /ui/config (used by the Obsidian URL), so no API change is needed. Added logic.absolutePath (obsidianUrl uses it), flow.yankPath, the Y key, palette action, help and hints. Tests: keys.test.mjs, logic.test.mjs. Fake DOM only, not a real browser. Stacked on the id yank branch.
