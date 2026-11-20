---
id: "261031-18t8"
title: spike drawer layout
summary: "Spike: the drawer as a first-class surface at half, two-thirds or centered width, without leaving the board"
status: active
kind: spike
parent: "261031-gtw6"
created: 2026-10-31
updated: 2026-11-19
claimed_by: claude-sonnet-5-5-drawer-spike
claimed_at: 2026-11-20T03:16
---

# spike drawer layout

## Goal

The drawer is a preview at a constrained width or a third of the screen, and it feels pushed to the side. It needs to be a first-class surface: big enough to read and edit a whole record, while keeping the board behind it, like an iOS sheet. A new page is not wanted. Find out:

- The width: half, two-thirds, centered, or a user-resizable panel that remembers its width; how each behaves on narrow screens and with the board still visible. Whether the board stays interactive behind it.
- The layout inside: fields in a header block, the body as sections, the log and outcome last; sticky header and actions; read mode and edit mode (one surface or two).
- Navigation: opening from a card, a link or the palette; going to a parent, child or linked record without closing it (a back stack); the URL for a drawer state so a link opens it.
- Keyboard and focus: open, close, move between records, save and cancel; unsaved changes and a stale-record warning.
- Constraints: the workbench has no dependencies and is plain JavaScript (see [261029-kb33](done/261029-kb33-workbench-form-controls.md)); light and dark; no reliance on color alone.

Done when:

- [ ] Each width option is built as a throwaway page in the workbench's own styles and checked in the browser at three viewport sizes, light and dark, with screenshots.
- [ ] A recommended layout and navigation model is in `## Outcome`, with the states the design must handle (read, edit, saving, conflict, error); follow-up items are filed, including changes to the existing drawer items.

## Next

Finish the spike: (1) probe a non-modal drawer (board not inert, scrim click-through) and check whether card click and focus behave; (2) check the back stack and the keyboard model against detail.js and app.js; (3) decide where the screenshots live in the repo; (4) write the Outcome with a recommended layout and the read, edit, saving, conflict and error states; (5) file follow-up items. Read the 2026-11-20 Log line first.

## Log

- 2026-10-31: filed.
- 2026-11-19 20:16-07:00: claimed (claude-sonnet-5-5-drawer-spike).
- 2026-11-19 20:33-07:00: Read the drawer (internal/server/web/js/detail.js, app.css .drawer). It is already modal: fixed 600px on the right, a scrim, the board inert, a record deep link at #/record/<id>, Esc closes. Built five variants as injected CSS over the live workbench (baseline, half, two-thirds, centered, resizable) and shot each at 1440, 1024 and 390 wide, light and dark. Findings so far: the 12-row field table fills the first screen and pushes the Goal, Next and Log below the fold at every width; at 390 wide the header row (title, Open in Obsidian, Close) squeezes the title to one or two characters per line; a fixed min-width over 390 overflows, so every width needs min(100vw, ...). Shots are in the session scratchpad, not yet in the repo.
