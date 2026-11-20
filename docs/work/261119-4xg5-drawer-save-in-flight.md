---
id: "261119-4xg5"
title: drawer save in flight
status: active
kind: bug
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
summary: The drawer Save and toggle buttons stay clickable while a save is in flight, so a second press is refused as stale
claimed_by: claude-sonnet-5-5-drawer-save
claimed_at: 2026-11-20T05:48
---

# drawer save in flight

## Goal

In the record drawer's field table (internal/server/web/js/fields.js, row() and save()), pressing Save sends a PATCH and waits, and nothing changes on screen while it waits: the Save button stays clickable. A second press (or Enter in the field) sends the same change again with the same etag; the first save changes the etag, so the second is refused as stale and shows 'This record changed since you opened it' for the user's own save. While a save is in flight the Save and Cancel buttons and the Turn on or Turn off toggle must be disabled, with a visible disabled look that does not rely on color alone (the Save label becomes 'Saving' and aria-busy is set), and a second submit must send nothing. They are enabled again if the save fails (an error or a stale answer). On success the drawer redraws as now.

Done when:

- [ ] while a save is pending the buttons are disabled and a second press or Enter sends no second request
- [ ] a failed or stale save enables them again and shows its message as now
- [ ] the disabled state looks disabled in light and dark (check in a browser; there is a .btn rule to extend in app.css)

## Next

Disable the buttons while a save is pending in fields.js (row, save), add the disabled look in app.css, test in drawer.test.mjs, check in a browser.

## Log

- 2026-11-19: filed.
- 2026-11-19 22:48-07:00: claimed (claude-sonnet-5-5-drawer-save).
