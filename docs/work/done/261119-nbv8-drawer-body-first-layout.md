---
id: "261119-nbv8"
title: drawer body first layout
status: done
kind: feature
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
closed: 2026-11-19
summary: Show the body sections first in the drawer and fold the 12-row field table under a one-line summary
---

# drawer body first layout

## Goal

The field table fills the first screen of the drawer at every width, so the Goal, Next and Log sit below the fold (found by the spike 261031-18t8). Put a one-line summary of the key fields (status, kind, parent, priority, claim) in the header, the body sections next, and the full field table in a collapsed section that opens for editing. The header and its actions stay in view while the body scrolls. The table's own layout bug is a separate item (261119-pk23).

Done when:

- [x] At 1440 by 900, the Goal is visible without scrolling and the field table is one click away.
- [x] Editing a field still works from the collapsed table, and the stale-write message still shows inline.
- [x] The header and its actions stay in view while the body scrolls.

## Next

## Outcome

The drawer is body first: a summary line under the path, a sticky header, the body sections, and last a collapsed Fields section with the value types legend and the field table (js/detail.js, app.css).\n\n- Box 1 (Goal visible without scrolling, table one click away): drawer.test.mjs asserts the summary line, that the first body section comes before the Fields section and that Fields is last and closed on first draw; in Chromium at 1440 by 900 on 261029-tgfe, in light and dark, the Goal body ends at y=658 of 900 and one click on Fields opens the table.\n- Box 2 (editing and the stale write): drawer.test.mjs edits priority, kind, summary and parent through the opened Fields section, and sends a stale 412 to see the inline message and the Reload button; Fields stays open through the redraw after a save. In Chromium, Edit opens the editor.\n- Box 3 (header stays in view): drawer.test.mjs checks the summary line and buttons live in the header; in Chromium, scrolled to the bottom, the header sits at the top of the drawer (1 px inside the border) in light and dark. The fake DOM has no layout, so position:sticky is proven in the browser only.\n\nNot changed: a claim reloads the drawer through open() (app.js reload), which resets the open state of every block, Fields included; a redraw through refresh() keeps it.\n\nDecision: none.

## Log

- 2026-11-19: filed.
- 2026-11-19 22:46-07:00: claimed (claude-sonnet-5-5-drawer-bodyfirst).
- 2026-11-19 23:01-07:00: ticked "At 1440 by 900, the Goal is visible without scrolling and the field table is one click away.": drawer.test.mjs (Fields section last and closed, body section before it) and a 1440 by 900 Chromium check on 261029-tgfe in light and dark: Goal body ends at y=658 of 900, Fields one click away
- 2026-11-19 23:01-07:00: ticked "Editing a field still works from the collapsed table, and the stale-write message still shows inline.": drawer.test.mjs (edit, cancel, bad value, select, summary and parent edits still PATCH through the Fields section; Fields stays open after the save redraw); the stale message code in fields.js is unchanged and its inline test is in this same file's flow; Chromium: Edit on kind opens the editor
- 2026-11-19 23:01-07:00: ticked "The header and its actions stay in view while the body scrolls.": drawer.test.mjs checks the sticky header holds the summary line and buttons; Chromium at 1440 by 900 scrolled to the bottom: header top at 1 px inside the drawer, in light and dark
- 2026-11-19 23:01-07:00: Built the body-first drawer: summary line under the path, sticky header, body sections before a collapsed Fields section (legend and field table). Tests in drawer.test.mjs; checked in Chromium at 1440 by 900 in light and dark. A claim still reloads the drawer through open(), which resets open blocks (app.js reload); noted, not changed.
- 2026-11-19 23:01-07:00: Added a stale-write test to drawer.test.mjs (412 on save shows the inline message and a Reload button; Reload redraws); this is the proof for box 2's stale-write half.
