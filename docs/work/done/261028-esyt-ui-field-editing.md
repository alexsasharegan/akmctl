---
id: "261028-esyt"
title: Show and edit a record's fields in the drawer
summary: The drawer shows frontmatter as a typed table and edits each editable field in place through set, with etag conflict handling; merges the typed-table item
status: done
kind: feature
priority: 1
created: 2026-10-28
updated: 2026-11-07
closed: 2026-11-07
---

# Show and edit a record's fields in the drawer

## Goal

The drawer shows a record's frontmatter as a typed table, and each editable field can be edited in place. This merges 261029-paw7 (the typed table) with the original request to edit more than status.

Table (from 261029-paw7): one row per field, each value typed and styled by its kind (string, date, number, boolean, list, wikilink), working from the parsed value rather than string matching so quoted wikilinks, aliases and lists behave the same everywhere. Styling does not rely on color alone. Lists show one item per line; wikilinks are links (see 261029-ctzc).

Editing: on an open work item, edit summary (text area, 261029-kb33), kind (select, 261029-kb33), priority (whole number), parent (picker, 261029-bdt9), blocked_by (multi-picker, 261029-bdt9), cycle (a label, checked as `set` checks it) and keep (a toggle). Read-only: id, created, updated, closed, claimed_by and claimed_at, which the tool manages. Which controls appear, and their locked reasons, come from the tool's `edit` states ([261101-3che](../261101-3che-show-reports-edit-states.md); the matrix is in the Outcome of 261031-z0n8): `kind`, `priority` and `parent` are editable with a warning once an item is active, and a closed item is read-only with a "Correct…" action for its title, summary and prose sections (ADR 261101-dxcn). ADRs are read-only in this first cut. Every edit goes through `set` with the record's etag.

Done when:

- [x] the drawer shows one row per field, typed and styled, with a label or legend that does not rely on color
- [x] lists render one item per line and wikilinks as links
- [x] each editable field above can be edited in place on an open work item
- [x] a stale write shows the conflict and writes nothing more; a refused edit shows the API's error next to its field
- [x] unit tests cover the value typing and the edit requests; checked in the browser in light and dark
- [x] new terms are registered in docs/glossary.md

## Next

## Outcome

The record drawer shows a record's fields as a table (`js/fields.js`): one row per field with its value type as a word (text, date, number, flag, list, link, id) and a legend; lists show one item per line and wikilinks and ids are links. On an open work item, summary, kind, priority, parent, blocked_by, cycle and keep are edited in place through `set` with the record's etag, using the controls from kb33 and bdt9. A locked field says why; changing kind or parent on an active item (and priority while active) warns; a closed item is read-only except keep. A write the API refuses shows its message beside the field and keeps the edit open; a stale write shows the conflict and writes nothing. A live refresh no longer redraws an edit in progress. Value typing, the edit rules and the edit requests are pure functions with node tests. Not done: the edit rules are a fixed copy of the 261031-z0n8 matrix until the tool reports them (261101-3che); the "Correct…" action for a closed item waits on 261101-3vwz; a blocked_by entry that is free text can be kept or removed but not added; screen-reader behavior is untested.

## Log

- 2026-10-28: filed.
- 2026-10-29: refined with Alexander; merged 261029-paw7 into this item.
- 2026-10-31: the edit rules (which fields when, closed items) are decided by spike 261031-z0n8 under 261031-gtw6; the table and simple fields can start first.
- 2026-10-31: blocked on 261031-z0n8; the edit rules and the drawer redesign (261031-gtw6) may change this item's scope.
- 2026-11-01 09:09-07:00: unblocked: 261031-z0n8 decided the edit rules; follow-ups 261101-3vwz and 261101-3che
- 2026-11-07 20:59-07:00: claimed (claude-sonnet-5-5).
- 2026-11-07: built `js/fields.js`: the typed table (`valueType`, a legend, lists one item per line, wikilinks as links) and in-place editing of summary, kind, priority, parent, blocked_by, cycle and keep through `set` with the etag, using the controls from 261029-kb33 and 261029-bdt9. The edit rules are the 261031-z0n8 matrix as a fixed table (`editState`). A live refresh no longer redraws an edit in progress. Checked in the browser against a copy of the repository, dark and light.
