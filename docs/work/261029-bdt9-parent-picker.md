---
id: "261029-bdt9"
title: "Workbench parent picker: a searchable list with match highlighting"
status: ready
created: 2026-10-29
updated: 2026-10-29
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

- [ ] stage 1: replaces the datalist in the capture modal, with highlighting of id and title matches, keyboard use (arrows, Enter, Escape) and screen reader support
- [ ] stage 2: result rows show id, title, status and summary
- [ ] stage 3: a multi-select variant for blocked_by
- [ ] unit tests for the matching and highlighting; checked in the browser in light and dark
- [ ] new terms are registered in docs/glossary.md

## Next

Ready to start after 261029-kb33 (shared control styles); stage 1 can begin earlier.

## Log

- 2026-10-29: filed.
