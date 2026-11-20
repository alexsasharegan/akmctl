---
id: "261119-nbv8"
title: drawer body first layout
status: active
kind: feature
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
summary: Show the body sections first in the drawer and fold the 12-row field table under a one-line summary
claimed_by: claude-sonnet-5-5-drawer-bodyfirst
claimed_at: 2026-11-20T05:46
---

# drawer body first layout

## Goal

The field table fills the first screen of the drawer at every width, so the Goal, Next and Log sit below the fold (found by the spike 261031-18t8). Put a one-line summary of the key fields (status, kind, parent, priority, claim) in the header, the body sections next, and the full field table in a collapsed section that opens for editing. The header and its actions stay in view while the body scrolls. The table's own layout bug is a separate item (261119-pk23).

Done when:

- [x] At 1440 by 900, the Goal is visible without scrolling and the field table is one click away.
- [x] Editing a field still works from the collapsed table, and the stale-write message still shows inline.
- [x] The header and its actions stay in view while the body scrolls.

## Next

Verify, commit, then close with an Outcome. If resuming: everything is built; run gofmt, vet, race tests, check and prettier.

## Log

- 2026-11-19: filed.
- 2026-11-19 22:46-07:00: claimed (claude-sonnet-5-5-drawer-bodyfirst).
- 2026-11-19 23:01-07:00: ticked "At 1440 by 900, the Goal is visible without scrolling and the field table is one click away.": drawer.test.mjs (Fields section last and closed, body section before it) and a 1440 by 900 Chromium check on 261029-tgfe in light and dark: Goal body ends at y=658 of 900, Fields one click away
- 2026-11-19 23:01-07:00: ticked "Editing a field still works from the collapsed table, and the stale-write message still shows inline.": drawer.test.mjs (edit, cancel, bad value, select, summary and parent edits still PATCH through the Fields section; Fields stays open after the save redraw); the stale message code in fields.js is unchanged and its inline test is in this same file's flow; Chromium: Edit on kind opens the editor
- 2026-11-19 23:01-07:00: ticked "The header and its actions stay in view while the body scrolls.": drawer.test.mjs checks the sticky header holds the summary line and buttons; Chromium at 1440 by 900 scrolled to the bottom: header top at 1 px inside the drawer, in light and dark
- 2026-11-19 23:01-07:00: Built the body-first drawer: summary line under the path, sticky header, body sections before a collapsed Fields section (legend and field table). Tests in drawer.test.mjs; checked in Chromium at 1440 by 900 in light and dark. A claim still reloads the drawer through open(), which resets open blocks (app.js reload); noted, not changed.
- 2026-11-19 23:01-07:00: Added a stale-write test to drawer.test.mjs (412 on save shows the inline message and a Reload button; Reload redraws); this is the proof for box 2's stale-write half.
