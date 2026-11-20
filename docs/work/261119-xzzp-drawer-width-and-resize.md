---
id: "261119-xzzp"
title: drawer centered with a scrim
status: active
kind: feature
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
summary: Show the record drawer centered over the board with a scrim, up to 1000 px wide, a full-screen sheet when narrow
claimed_by: claude-sonnet-5-5-drawer-build
claimed_at: 2026-11-20T04:59
---

# drawer centered with a scrim

## Goal

The spike 261031-18t8 settled the drawer's shape after Alexander tried it in the workbench: centered over the board, up to 1000 px wide with 16 px of margin at the sides and 4 percent of the height above and below, rounded corners, a scrim over an inert board, and it rises in instead of sliding. Under about 700 px it is a full-screen sheet. Every width is capped with min(100vw, ...), because a fixed minimum above 390 px overflows a phone. The trial CSS is in the spike's worktree history and was not committed.

Done when:

- [x] The drawer opens centered at the size above, with the scrim and the inert board, and Esc and a click on the scrim close it.
- [x] Under about 700 px it is a full-screen sheet.
- [ ] Checked at 1440, 1024 and 390 px wide, light and dark, with screenshots on the pull request.

## Next

Attach the six screenshots in .evidence/drawer-centered to the pull request, tick the last box, then close the item. Known and not in this item: at 390 px the header squeezes the title (item 4em1).

## Log

- 2026-11-19: filed.
- 2026-11-19 21:59-07:00: claimed (claude-sonnet-5-5-drawer-build).
- 2026-11-19 22:01-07:00: ticked "The drawer opens centered at the size above, with the scrim and the inert board, and Esc and a click on the scrim close it.": browser check on the built layer at 1440x900: drawer 1000 px wide at x 220-1220, y 36-864; Esc and a scrim click close it, a click inside does not; board inert while open
- 2026-11-19 22:01-07:00: ticked "Under about 700 px it is a full-screen sheet.": .evidence/drawer-centered/centered-narrow-light.png and -dark.png: 390 px wide, the drawer fills the screen with no radius or border
- 2026-11-19 22:01-07:00: Rewrote .drawer as a centered sheet with a narrow full-screen rule and a rise animation (replacing slide); changelog and glossary updated; go test -race ./..., go vet and gofmt clean; six screenshots in .evidence/drawer-centered. Left: attach the screenshots to the pull request.
- 2026-11-19 22:18-07:00: Follow-up from Alexander inspecting the tip: nothing was left under the Log section at the bottom. The body shrank to fit the drawer (flex-shrink 1), so its 24 px bottom padding spilled out of view (measured: 0 px gap below the last section). .dlg-body is flex:none inside the drawer now.
