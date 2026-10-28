---
id: "261028-fcc1"
title: Add a capture form to the workbench
summary: The workbench cannot create a record; add a form that files a raw work item (and an ADR draft) through the new operation.
status: raw
kind: feature
priority: 0
created: 2026-10-28
updated: 2026-10-28
---

# Add a capture form to the workbench

## Goal

The workbench cannot create a record. The `new` operation exists in the API; the page needs a form to file a raw work item or an ADR draft.

Done when:

- [ ] A form with a title, a kind and an optional summary creates a `raw` work item through `new`; the new item appears on the board.
- [ ] An ADR draft can be created the same way.
- [ ] Validation errors from the API are shown on the form.

## Log

- 2026-10-28: filed.
