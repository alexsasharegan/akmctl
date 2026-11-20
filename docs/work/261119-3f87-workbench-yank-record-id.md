---
id: "261119-3f87"
title: "Workbench: yank a record's id to the clipboard"
status: active
kind: feature
created: 2026-11-19
updated: 2026-11-19
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-20T02:28
---

# Workbench: yank a record's id to the clipboard

## Goal

In the workbench (`akmctl ui`) the `y` key copies the id of the focused or open record to the clipboard and shows a brief toast. Done when `y` is bound in the views that show a record (board, tables, the record view), is unbound while a text field has focus, falls back when the Clipboard API is blocked, appears in the keys help, and the front-end tests cover it. GitHub issue 97. Do the path (`Y`) in its own item.

## Next

Open the workbench in a real browser and press y on the board, a table row and an open record; then close the item after the pull request merges. The path key (Y) is its own item.

## Log

- 2026-11-19: filed.
- 2026-11-19 19:28-07:00: claimed (claude-sonnet-5-5).
- 2026-11-19 19:30-07:00: Built y: web/js/clipboard.js (Clipboard API, hidden-textarea fallback), flow.yank, focused() on the row views, app.js key, palette action, help section Copy and hints. y was unbound everywhere. Tests in keys.test.mjs: card, table row, refused clipboard. Checked under the fake DOM only, not in a real browser; the clipboard in a real page over localhost is untested here. Open point: the check view copies a finding row's id only if the row has a data-id.
