---
id: "261119-4xg5"
title: drawer save in flight
status: done
kind: bug
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
closed: 2026-11-19
summary: The drawer Save and toggle buttons stay clickable while a save is in flight, so a second press is refused as stale
---

# drawer save in flight

## Goal

In the record drawer's field table (internal/server/web/js/fields.js, row() and save()), pressing Save sends a PATCH and waits, and nothing changes on screen while it waits: the Save button stays clickable. A second press (or Enter in the field) sends the same change again with the same etag; the first save changes the etag, so the second is refused as stale and shows 'This record changed since you opened it' for the user's own save. While a save is in flight the Save and Cancel buttons and the Turn on or Turn off toggle must be disabled, with a visible disabled look that does not rely on color alone (the Save label becomes 'Saving' and aria-busy is set), and a second submit must send nothing. They are enabled again if the save fails (an error or a stale answer). On success the drawer redraws as now.

Done when:

- [x] while a save is pending the buttons are disabled and a second press or Enter sends no second request
- [x] a failed or stale save enables them again and shows its message as now
- [x] the disabled state looks disabled in light and dark (check in a browser; there is a .btn rule to extend in app.css)

## Next

## Outcome

Done: the drawer's Save, Cancel and Turn on or off buttons are disabled while a field save is pending (fields.js, row() and save()), and a second submit sends nothing.

- Box 1 (pending save disables the buttons, a second press or Enter sends nothing): internal/server/testdata/drawer.test.mjs, the new block at the end. It holds the PATCH with a promise, asserts Save and Cancel are disabled, Save reads Saving with aria-busy, and a second click and an Enter keydown add no call to app.calls; the toggle is covered the same way. It fails with fields.js reverted (go test ./internal/server -run TestDOMModulesUnderNode) and passes with it. In Chrome one PATCH was sent for a click, a forced click and Enter.
- Box 2 (failed or stale save enables them again, message as now): the same test, for an error answer and a stale answer: buttons enabled, label Save again, aria-busy gone, the message shown (the stale one is the usual text).
- Box 3 (looks disabled in light and dark): checked in Chrome with Playwright in both color schemes; the .evidence/light-saving.png and dark-saving.png screenshots show the faded, dashed-border, Saving button and a faded dashed Cancel. The rule is in app.css next to .btn[disabled].

Verified: gofmt -l prints nothing, go vet and go test -race ./... pass, akmctl check shows 0 errors, prettier passes on CHANGELOG.md and docs.

Decision: none.

## Log

- 2026-11-19: filed.
- 2026-11-19 22:46-07:00: claimed (claude-sonnet-5-5-drawer-save).
- 2026-11-19 22:47-07:00: ticked "while a save is pending the buttons are disabled and a second press or Enter sends no second request": internal/server/testdata/drawer.test.mjs: a pending save disables Save and Cancel and a second click and Enter add no call to app.calls; fails without the change
- 2026-11-19 22:47-07:00: ticked "a failed or stale save enables them again and shows its message as now": drawer.test.mjs: after an error answer and a stale answer the buttons are enabled, labeled Save again, and the message shows
- 2026-11-19 22:47-07:00: ticked "the disabled state looks disabled in light and dark (check in a browser; there is a .btn rule to extend in app.css)": browser check in Chrome light and dark: .evidence/{light,dark}-saving.png show faded, dashed, labeled Saving; one PATCH sent for click, forced click and Enter
- 2026-11-19 22:47-07:00: Disabled Save, Cancel and the toggle while a save is pending (fields.js row and save), added the dashed disabled look in app.css, tests in drawer.test.mjs, checked in Chrome in light and dark.
