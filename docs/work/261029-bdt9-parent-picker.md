---
id: "261029-bdt9"
title: "Workbench parent picker: a searchable list with match highlighting"
status: active
created: 2026-10-29
updated: 2026-11-07
kind: feature
priority: 0
summary: Replace the capture modal's native datalist parent input with a searchable list that shows what matched and enough of each item to choose well
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-07T08:14
---

# Workbench parent picker: a searchable list with match highlighting

## Goal

A rich parent picker, used by the capture modal (261029-aezq) and the drawer editor (261028-esyt, which also needs it for blocked_by). Today the parent is a native datalist input (`cap-parent` in capture.js).

From Alexander's notes on the capture modal:

This is a massive opportunity for rich UI interactions, but also for big builds, and needs a lot of runway. The current search is squeezed into a small space, does not offer any kind of emphasis on where matches are being made (in description or id), and the result options are too small to inspect to inform a good choice.

Build it in stages, each shippable:

1. A combobox with a results list that is wide enough to read, and match highlighting that shows whether the id, the title or the summary matched.
2. Result rows that show the id, title, status and the first line of the summary, so a choice can be made without opening anything.
3. Multi-select for blocked_by, with the chosen items shown as removable chips.

Only open work items are offered for parent (the existing rule); closed items are not.

Done when:

- [x] stage 1: replaces the datalist in the capture modal, with highlighting of id and title matches, keyboard use (arrows, Enter, Escape) and screen reader support
- [x] stage 2: result rows show id, title, status and summary
- [ ] stage 3: a multi-select variant for blocked_by
- [x] unit tests for the matching and highlighting; checked in the browser in light and dark
- [x] new terms are registered in docs/glossary.md

## Next

Stages 1 and 2 are done (`js/picker.js`, used by the capture modal). Stage 3 is next: a multi-select variant for `blocked_by`, with removable chips; the drawer editor (261028-esyt) needs it.

## Log

- 2026-10-29: filed.
- 2026-11-07: stages 1 and 2: `js/picker.js` (`createPicker`, with the pure `searchItems`, `matchRanges` and `highlight`) replaces the datalist in the capture modal; rows show id, title, status and the summary with matches marked. Checked in the browser, light and dark. Stage 3 remains, so the item stays active.
- 2026-11-07 01:14-07:00: claimed (claude-sonnet-5-5).
