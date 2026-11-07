---
id: "261029-kb33"
title: "Workbench form controls: a custom select and a summary text area with a counter"
status: done
created: 2026-10-29
updated: 2026-11-06
closed: 2026-11-06
kind: feature
priority: 1
summary: A custom select and a summary text area with a character counter for the workbench forms, in plain JavaScript, used by capture and the drawer
---

# Workbench form controls: a custom select and a summary text area with a counter

## Goal

Shared form controls for the capture modal (261029-aezq) and the drawer editor (261028-esyt). The workbench has no dependencies, so these are written in plain JavaScript; Radix and similar libraries are only a reference for how a custom select should feel.

1. Summary field (from Alexander's notes on the capture modal)

160 characters does not fit nicely into this field and scrolls awkwardly. Consider a text area, or some way to see lines wrap. Consider adding helper text below the input which counts characters and changes when input hits 160 or hits about 90% of the maximum.

2. Kind (and other short choices: priority, status)

The input is not the same height as its siblings in the row. Plain HTML option/select inputs always feel like plain HTML inputs. Radix UI or other popular libraries may serve as a good example of what custom UI selects feel like.

Done when:

- [x] a custom select (listbox pattern) replaces the native select for kind, with the same height as the inputs beside it
- [x] the select works from the keyboard (open, arrows, type-ahead, Enter, Escape) and for screen readers (role, selected state announced)
- [x] the summary field is a text area that wraps, with a counter that changes at about 90% and at 160 characters
- [x] both controls look right in light and dark, with visible focus
- [x] unit tests cover the counter thresholds and the select's key handling; checked in the browser
- [x] new terms are registered in docs/glossary.md

## Next

## Outcome

The capture modal now uses two shared controls from `js/controls.js`. `createSelect` is a listbox-pattern select: it has the same height as its neighbours, opens with Enter, Space or the arrows, Home and End move, a typed letter jumps to a match, Escape closes the list without closing the dialog, and Tab closes it. It announces the selected option and opens upward when there is no room below. `createSummary` is a wrapping text area for one line of text (Enter captures, a pasted line break becomes a space) with a counter that warns from 144 of 160 characters and says "limit reached" at 160. The key handling and the counter are pure functions with node tests in `testdata/logic.test.mjs`. Not done: the controls are used only by capture; the drawer editor (261028-esyt) and the capture polish (261029-aezq) adopt them next, and priority could take the select too.

## Log

- 2026-10-29: filed.
- 2026-11-06: built `js/controls.js`: `createSelect` (a button with `role=combobox` and a listbox, keys by the pure `selectKey`) and `createSummary` (a text area with the counter from `counterState`); capture uses both. Checked in the browser, dark and light.
- 2026-11-06 21:38-07:00: claimed (claude-sonnet-5-5).
