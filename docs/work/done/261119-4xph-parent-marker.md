---
id: "261119-4xph"
title: parent marker
status: done
kind: feature
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
closed: 2026-11-19
summary: A chip on the board card and the drawer summary line shows how many child items are done
---

# parent marker

## Goal

A work item that has children shows a compact chip with the glyph U+21B3, a space and DONE/TOTAL (for example "↳ 3/5"), on the board card and in the drawer summary line. Children are the work items whose parent is the item's id, in any status. TOTAL counts the children that are not dropped; DONE counts those of them with status done. If there are none, or every child is dropped, nothing is drawn. The chip does not depend on kind and the kind chip is unchanged. The glyph is aria-hidden; the chip has a title such as "3 of 5 child items done" ("1 of 1 child item done") and the same words as a visually hidden name. It is a plain chip (text color, normal weight), not dashed. In the drawer the chip links to the Children list.

Done when:

- [x] One helper counts the children per work item once per render, and the board and the drawer both use it; tests cover none, all dropped, mixed done, ready and dropped.
- [x] The board card shows the chip after the priority chip and before the kind chip, with its title and accessible name; a card without children shows none; a chore parent shows it; the kind chip is unchanged.
- [x] The drawer summary line shows the chip after the priority chip, as a keyboard-operable link to the Children list.
- [x] Checked in a real browser, light and dark, on a record with many children and on the board, with no header wrap at the narrowest column.
- [x] CHANGELOG Unreleased entry and glossary term "Child count".

## Next

## Outcome

1. One helper: `childCounts` in internal/server/web/js/logic.js counts every parent's children in one pass; board.js (`mountBoard`, once per mount) and detail.js (once per render) both use it. Proof: the child counts block in internal/server/testdata/logic.test.mjs (none, all dropped, mixed done, ready, dropped, ADR ignored); it fails without the change (`L.childCounts is not a function`).
2. Board card: testdata/kids.test.mjs checks "↳ 3/5" with the title "3 of 5 child items done" and the same words in a hidden span, no chip on a card without children or with only dropped children, a chore parent shows "↳ 1/1" (singular title) between the priority and kind chips, and the kind chip class stays "chip kind". It fails with the web sources reverted.
3. Drawer: the same file checks the chip in the summary line after the priority, as a button that moves focus to the Children heading; it fails with detail.js reverted.
4. Real browser: Chromium via Playwright, light and dark, 390, 700, 900 and 1440 px, on the board and on the drawer of the project with 24 live children ("↳ 22/24"). The chip is not dashed, weight 400. The drawer chip scrolls the Children heading below the sticky header (found and fixed during the check). At 1440 px no card header wraps. In columns of 174 px and narrower, card headers already wrapped without the chip (id, priority and kind); the chip, which never splits, adds a line on some cards at 174 px and on one at 156 px.
5. CHANGELOG.md Unreleased entry and the glossary term "Child count" (and the Summary line row).

Verified: gofmt -l prints nothing, go vet clean, go test -race ./... passes, akmctl check 0 errors (5 warnings about other items), prettier check passes.

Decision: none.

## Log

- 2026-11-19: filed.
- 2026-11-19 23:30-07:00: claimed (claude-sonnet-5-5-parent-marker).
- 2026-11-19 23:45-07:00: ticked "One helper counts the children per work item once per render, and the board and the drawer both use it; tests cover none, all dropped, mixed done, ready and dropped.": childCounts in js/logic.js; logic.test.mjs child counts block; board.js and detail.js both use it
- 2026-11-19 23:45-07:00: ticked "The board card shows the chip after the priority chip and before the kind chip, with its title and accessible name; a card without children shows none; a chore parent shows it; the kind chip is unchanged.": testdata/kids.test.mjs board card and chore parent blocks
- 2026-11-19 23:45-07:00: ticked "The drawer summary line shows the chip after the priority chip, as a keyboard-operable link to the Children list.": testdata/kids.test.mjs drawer block; button, focus moves to #drawer-kids
- 2026-11-19 23:45-07:00: ticked "Checked in a real browser, light and dark, on a record with many children and on the board, with no header wrap at the narrowest column.": Playwright on 261031-gtw6 and the board, light and dark, 390 to 1440 px; screenshots in .evidence
- 2026-11-19 23:45-07:00: ticked "CHANGELOG Unreleased entry and glossary term "Child count".": CHANGELOG Unreleased entry; glossary row Child count
- 2026-11-19 23:45-07:00: Built childCounts (js/logic.js), childChip (js/util.js), the board card chip and the drawer summary chip (a button to the Children list); tests in kids.test.mjs and logic.test.mjs fail without the source change; checked in Chromium, light and dark.
