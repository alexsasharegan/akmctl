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

- [ ] At 1440 by 900, the Goal is visible without scrolling and the field table is one click away.
- [ ] Editing a field still works from the collapsed table, and the stale-write message still shows inline.
- [ ] The header and its actions stay in view while the body scrolls.

## Next

Build the body-first drawer in js/detail.js: summary line under the path, sticky header, body sections before a collapsed Fields section holding the Value types legend and the field table. Then tests in drawer.test.mjs, CSS, a browser check at 1440 by 900 in light and dark, the changelog and glossary.

## Log

- 2026-11-19: filed.
- 2026-11-19 22:46-07:00: claimed (claude-sonnet-5-5-drawer-bodyfirst).
