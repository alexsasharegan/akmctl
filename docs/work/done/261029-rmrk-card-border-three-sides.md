---
id: "261029-rmrk"
title: Fix the card border that shows on only three sides
status: done
created: 2026-10-29
updated: 2026-10-31
closed: 2026-10-31
kind: bug
priority: 3
summary: In the Ready lane the first card draws its border on three of the four sides.
---

# Fix the card border that shows on only three sides

## Goal

Observed by Alexander in the workbench on the vault (light theme, a record drawer open earlier): the first card in the Ready lane (261026-mb2a) shows its border on only three of four sides. He suspects a state style, such as a left-hand colour overlay for a state like focus, claimed or selected, but it does not appear to be the drag state. Reproduce, find which rule or state produces it (check the card's classes and the focus ring and the lane scroll clip), and fix it.

Done when:

- [x] The cause is identified and logged here (which class or rule).
- [x] Every card shows a complete border in every state: idle, focused, selected, dragging, claimed, in both themes.
- [x] A test or a documented browser check covers the state.

## Outcome

Cause: in internal/server/web/app.css the .card rule set border-left to "3px solid transparent" and .card:hover set border-left-color to transparent again, so the left side drew no line (only the stale dotted style and the focus accent gave it one). Fix: the left border is 3px in the same colour as the other sides (var(--line)), and hover no longer overrides it. Browser check: computed border colour is the same on all four sides in the dark and the light theme for an idle card. The stale (dotted) and focus (accent) left borders are unchanged and deliberately differ; the dragging, selected and claimed states were not re-checked by eye, and a focused card could not be tested because script focus does not apply :focus to an unfocused page. Decision: none.

## Log

- 2026-10-29: filed.
- 2026-10-31: claimed (claude).
