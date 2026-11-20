---
id: "261119-3f87"
title: "Workbench: yank a record's id to the clipboard"
status: raw
kind: feature
created: 2026-11-19
updated: 2026-11-19
---

# Workbench: yank a record's id to the clipboard

## Goal

In the workbench (`akmctl ui`) the `y` key copies the id of the focused or open record to the clipboard and shows a brief toast. Done when `y` is bound in the views that show a record (board, tables, the record view), is unbound while a text field has focus, falls back when the Clipboard API is blocked, appears in the keys help, and the front-end tests cover it. GitHub issue 97. Do the path (`Y`) in its own item.

## Log

- 2026-11-19: filed.
