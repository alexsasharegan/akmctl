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

- [ ] G scrolls the open drawer to the bottom and gg, pressed twice within 900 ms, to the top; one slow g does nothing.
- [ ] The help overlay and the hint bar list both keys for the drawer.
- [ ] A workbench page test covers both.

## Next

Add gg and G to createDrawer().keys in js/detail.js (reuse the trial in the scratchpad), list them in the hint bar and the Record help group in app.js, add a drawer keys test in internal/server/testdata, a CHANGELOG entry, then verify and close.

## Log

- 2026-11-19: filed.
- 2026-11-19 21:59-07:00: claimed (claude-sonnet-5-5-drawer-keys).
