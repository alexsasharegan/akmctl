---
id: "261119-za8w"
title: drawer top and bottom keys
status: active
kind: feature
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
summary: Add gg and G to the record drawer to scroll to the top and the bottom
claimed_by: claude-sonnet-5-5-drawer-keys
claimed_at: 2026-11-20T04:59
---

# drawer top and bottom keys

## Goal

The drawer scrolls with j and k. Add gg to the top and G (Shift+G) to the bottom, using the board's keys and its 900 ms window for the second g (js/board.js, pendingG), so the muscle memory carries over. The spike 261031-18t8 trial built it in js/detail.js keys(), and it worked on a long record. List the keys in the help overlay and the hint bar while the drawer is open.

Done when:

- [x] G scrolls the open drawer to the bottom and gg, pressed twice within 900 ms, to the top; one slow g does nothing.
- [x] The help overlay and the hint bar list both keys for the drawer.
- [x] A workbench page test covers both.

## Next

Verified; close the item.

## Log

- 2026-11-19: filed.
- 2026-11-19 21:59-07:00: claimed (claude-sonnet-5-5-drawer-keys).
- 2026-11-19 22:00-07:00: ticked "G scrolls the open drawer to the bottom and gg, pressed twice within 900 ms, to the top; one slow g does nothing.": internal/server/testdata/drawer.test.mjs: G, gg within 900 ms, one g, a slow second g, a key between, and modifier keys
- 2026-11-19 22:00-07:00: ticked "The help overlay and the hint bar list both keys for the drawer.": internal/server/testdata/drawer.test.mjs checks the hint bar and the Record help group; app.js
- 2026-11-19 22:00-07:00: ticked "A workbench page test covers both.": internal/server/testdata/drawer.test.mjs (run by TestDOMModulesUnderNode)
- 2026-11-19 22:00-07:00: Added gg and G to the drawer keys (js/detail.js), listed them in the hint bar and the Record help group (app.js), tests in drawer.test.mjs, fakedom scrollTo and scrollBy
