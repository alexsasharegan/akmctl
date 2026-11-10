---
id: "261031-18t8"
title: spike drawer layout
summary: "Spike: the drawer as a first-class surface at half, two-thirds or centered width, without leaving the board"
status: raw
kind: spike
parent: "261031-gtw6"
created: 2026-10-31
updated: 2026-10-31
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

## Log

- 2026-10-31: filed.
