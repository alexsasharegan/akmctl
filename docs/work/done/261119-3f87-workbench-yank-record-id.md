---
id: "261119-3f87"
title: "Workbench: yank a record's id to the clipboard"
status: done
kind: feature
created: 2026-11-19
updated: 2026-11-19
closed: 2026-11-19
---

# Workbench: yank a record's id to the clipboard

## Goal

In the workbench (`akmctl ui`) the `y` key copies the id of the focused or open record to the clipboard and shows a brief toast. Done when `y` is bound in the views that show a record (board, tables, the record view), is unbound while a text field has focus, falls back when the Clipboard API is blocked, appears in the keys help, and the front-end tests cover it. GitHub issue 97. Do the path (`Y`) in its own item.

## Next

## Outcome

In the workbench, y copies the id of the focused card or table row, or of the open record, with a toast; a refused clipboard says so. The Clipboard API is used with a hidden-textarea fallback. Tested under the fake DOM and by the maintainer in a real browser. Merged in pull request 100. GitHub issue 97.

## Log

- 2026-11-19: filed.
- 2026-11-19 19:28-07:00: claimed (claude-sonnet-5-5).
- 2026-11-19 19:30-07:00: Built y: web/js/clipboard.js (Clipboard API, hidden-textarea fallback), flow.yank, focused() on the row views, app.js key, palette action, help section Copy and hints. y was unbound everywhere. Tests in keys.test.mjs: card, table row, refused clipboard. Checked under the fake DOM only, not in a real browser; the clipboard in a real page over localhost is untested here. Open point: the check view copies a finding row's id only if the row has a data-id.
