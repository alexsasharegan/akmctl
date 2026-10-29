---
id: "261029-rmrk"
title: Fix the card border that shows on only three sides
status: ready
created: 2026-10-29
updated: 2026-10-29
kind: bug
priority: 3
summary: In the Ready lane the first card draws its border on three of the four sides.
---

# Fix the card border that shows on only three sides

## Goal

Observed by Alexander in the workbench on the vault (light theme, a record drawer open earlier): the first card in the Ready lane (261026-mb2a) shows its border on only three of four sides. He suspects a state style, such as a left-hand colour overlay for a state like focus, claimed or selected, but it does not appear to be the drag state. Reproduce, find which rule or state produces it (check the card's classes and the focus ring and the lane scroll clip), and fix it.

Done when:

- [ ] The cause is identified and logged here (which class or rule).
- [ ] Every card shows a complete border in every state: idle, focused, selected, dragging, claimed, in both themes.
- [ ] A test or a documented browser check covers the state.

## Log

- 2026-10-29: filed.
