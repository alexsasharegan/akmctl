---
id: "261119-xzzp"
title: drawer centered with a scrim
status: raw
kind: feature
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
summary: Show the record drawer centered over the board with a scrim, up to 1000 px wide, a full-screen sheet when narrow
---

# drawer centered with a scrim

## Goal

The spike 261031-18t8 settled the drawer's shape after Alexander tried it in the workbench: centered over the board, up to 1000 px wide with 16 px of margin at the sides and 4 percent of the height above and below, rounded corners, a scrim over an inert board, and it rises in instead of sliding. Under about 700 px it is a full-screen sheet. Every width is capped with min(100vw, ...), because a fixed minimum above 390 px overflows a phone. The trial CSS is in the spike's worktree history and was not committed.

Done when:

- [ ] The drawer opens centered at the size above, with the scrim and the inert board, and Esc and a click on the scrim close it.
- [ ] Under about 700 px it is a full-screen sheet.
- [ ] Checked at 1440, 1024 and 390 px wide, light and dark, with screenshots on the pull request.

## Log

- 2026-11-19: filed.
