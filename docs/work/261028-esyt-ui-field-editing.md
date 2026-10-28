---
id: "261028-esyt"
title: Edit more than status in the workbench
summary: The workbench edits a record's status only; other fields (priority, parent, blocked_by, summary, kind) cannot be changed from the page.
status: raw
kind: feature
priority: 0
created: 2026-10-28
updated: 2026-10-28
---

# Edit more than status in the workbench

## Goal

The workbench can change a record's status only. The `set` operation can change `priority`, `parent`, `blocked_by`, `summary` and `kind`, and the page should offer them.

Done when:

- [ ] Each field `set` accepts can be edited from the record's detail view, with the etag compare-and-swap handled on conflict.
- [ ] A refused edit shows the API's error message.

## Log

- 2026-10-28: filed.
