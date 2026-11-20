---
id: "261119-nbv8"
title: drawer body first layout
status: raw
kind: feature
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
summary: Show the body sections first in the drawer and fold the 12-row field table under a one-line summary
---

# drawer body first layout

## Goal

The field table fills the first screen of the drawer at every width, so the Goal, Next and Log sit below the fold (found by the spike 261031-18t8). Put a one-line summary of the key fields (status, kind, parent, priority, claim) in the header, the body sections next, and the full field table in a collapsed section that opens for editing. The header and its actions stay in view while the body scrolls. The table also has a layout bug: .fm-name is display:flex on a th, so it stops being a table cell and the name column, the row dividers and the Edit links drift apart. The fix, found in the trial, is to leave the th a cell, stack its two spans as blocks, align the cells on their first line, and trim the Edit button's vertical padding.

Done when:

- [ ] At 1440 by 900, the Goal is visible without scrolling and the field table is one click away.
- [ ] Editing a field still works from the collapsed table, and the stale-write message still shows inline.
- [ ] The header and its actions stay in view while the body scrolls.
- [ ] In the table, the name column is narrow, the dividers run the full width, and Edit lines up with its value, in light and dark.

## Log

- 2026-11-19: filed.
