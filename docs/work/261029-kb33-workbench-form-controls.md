---
id: "261029-kb33"
title: "Workbench form controls: a custom select and a summary text area with a counter"
status: ready
created: 2026-10-29
updated: 2026-10-29
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

- [ ] a custom select (listbox pattern) replaces the native select for kind, with the same height as the inputs beside it
- [ ] the select works from the keyboard (open, arrows, type-ahead, Enter, Escape) and for screen readers (role, selected state announced)
- [ ] the summary field is a text area that wraps, with a counter that changes at about 90% and at 160 characters
- [ ] both controls look right in light and dark, with visible focus
- [ ] unit tests cover the counter thresholds and the select's key handling; checked in the browser
- [ ] new terms are registered in docs/glossary.md

## Next

Ready to start; no dependencies.

## Log

- 2026-10-29: filed.
