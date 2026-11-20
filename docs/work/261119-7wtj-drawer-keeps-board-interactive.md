---
id: "261119-7wtj"
title: drawer keeps board interactive
status: raw
kind: feature
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
summary: "Do not make the board inert behind the drawer: a click on a visible card opens that record in the drawer"
---

# drawer keeps board interactive

## Goal

Today the drawer makes the board inert and covers it with a scrim, so nothing behind it can be clicked. The spike 261031-18t8 probed the alternative: with the board not inert and a transparent, click-through scrim, a click on a visible card switches the drawer to that record, browser Back returns to the previous record, and Esc closes the drawer and puts focus back on the card. This item makes that the behavior on wide screens, and keeps the modal behavior in the full-width sheet.

Done when:

- [ ] On a wide screen a click on a visible card opens that record in the drawer, Back returns to the previous record, and Esc closes the drawer and restores focus to the card.
- [ ] In the full-width sheet the board stays inert and the focus stays inside the drawer.
- [ ] The workbench page tests cover both modes.

## Log

- 2026-11-19: filed.
