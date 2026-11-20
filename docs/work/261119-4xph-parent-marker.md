---
id: "261119-4xph"
title: parent marker
status: ready
kind: feature
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
summary: A chip on the board card and the drawer summary line shows how many child items are done
---

# parent marker

## Goal

A work item that has children shows a compact chip with the glyph U+21B3, a space and DONE/TOTAL (for example "↳ 3/5"), on the board card and in the drawer summary line. Children are the work items whose parent is the item's id, in any status. TOTAL counts the children that are not dropped; DONE counts those of them with status done. If there are none, or every child is dropped, nothing is drawn. The chip does not depend on kind and the kind chip is unchanged. The glyph is aria-hidden; the chip has a title such as "3 of 5 child items done" ("1 of 1 child item done") and the same words as a visually hidden name. It is a plain chip (text color, normal weight), not dashed. In the drawer the chip links to the Children list.

Done when:

- [ ] One helper counts the children per work item once per render, and the board and the drawer both use it; tests cover none, all dropped, mixed done, ready and dropped.
- [ ] The board card shows the chip after the priority chip and before the kind chip, with its title and accessible name; a card without children shows none; a chore parent shows it; the kind chip is unchanged.
- [ ] The drawer summary line shows the chip after the priority chip, as a keyboard-operable link to the Children list.
- [ ] Checked in a real browser, light and dark, on a record with many children and on the board, with no header wrap at the narrowest column.
- [ ] CHANGELOG Unreleased entry and glossary term "Child count".

## Next

Add childCounts to logic.js, use it in board.js and detail.js, add tests, CSS, changelog, glossary.

## Log

- 2026-11-19: filed.
