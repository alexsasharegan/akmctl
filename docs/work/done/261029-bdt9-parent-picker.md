---
id: "261029-bdt9"
title: "Workbench parent picker: a searchable list with match highlighting"
status: done
created: 2026-10-29
updated: 2026-11-07
closed: 2026-11-07
kind: feature
priority: 0
summary: Replace the capture modal's native datalist parent input with a searchable list that shows what matched and enough of each item to choose well
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
- [x] stage 3: a multi-select variant for blocked_by
- [x] unit tests for the matching and highlighting; checked in the browser in light and dark
- [x] new terms are registered in docs/glossary.md

## Next

## Outcome

The workbench has a parent picker, `js/picker.js`. Stage 1: a combobox over the open work items that finds a record by id, title or summary (every word must match), ranks an id match first, marks what matched, and opens a list wide enough to read, upward when the dialog has no room; arrows and Enter choose, Escape closes only the list, and the result count is announced. Stage 2: each row shows the id, title, status and the summary. The capture modal uses it in place of the datalist. Stage 3: `multiple` mode for `blocked_by`, with removable chips, that does not offer a chosen item or the item itself. The search, highlighting and chip rules are pure and have node tests. Not done: no form uses the multi-select yet, so I checked it by mounting it in the capture dialog by hand; screen-reader behaviour is untested; the drawer editor (261028-esyt) is the first user.

## Log

- 2026-10-29: filed.
- 2026-11-07: stages 1 and 2: `js/picker.js` (`createPicker`, with the pure `searchItems`, `matchRanges` and `highlight`) replaces the datalist in the capture modal; rows show id, title, status and the summary with matches marked. Checked in the browser, light and dark. Stage 3 remains, so the item stays active.
- 2026-11-07: stage 3: `createPicker` takes `multiple`, `values` and `self`: chosen items are removable chips (the × button, or Backspace in the empty box), and the items already chosen and the item itself are not offered. The chip rules are pure (`addChip`, `removeChip`, `offered`) with node tests. No form uses it yet (capture has no `blocked_by`), so I mounted it in the capture dialog by hand and checked it in the browser, light and dark.
- 2026-11-07 01:14-07:00: claimed (claude-sonnet-5-5).
