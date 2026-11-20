---
id: "261119-xzzp"
title: drawer width and resize
status: raw
kind: feature
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
summary: Open the drawer at about 58 percent of the viewport, resizable, remembering its width, a full-width sheet when narrow
---

# drawer width and resize

## Goal

The spike 261031-18t8 recommends a drawer that opens at about 58 percent of the viewport (never under 480 px), can be resized by dragging its left edge, remembers the width, and becomes a full-width sheet under about 700 px. Every width is capped with min(100vw, ...), because a fixed minimum above 390 px overflows a phone.

Done when:

- [ ] The drawer opens at the remembered width, or about 58 percent the first time, capped at the viewport.
- [ ] Dragging the left edge resizes it, with the keyboard as an alternative, and the width survives a reload.
- [ ] Under about 700 px the drawer is a full-width sheet.
- [ ] Checked at 1440, 1024 and 390 px wide, light and dark, with screenshots on the pull request.

## Log

- 2026-11-19: filed.
