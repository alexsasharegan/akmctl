---
id: "261031-18t8"
title: spike drawer layout
summary: "Spike: the drawer as a first-class surface at half, two-thirds or centered width, without leaving the board"
status: done
kind: spike
parent: "261031-gtw6"
created: 2026-10-31
updated: 2026-11-19
closed: 2026-11-19
---

# spike drawer layout

## Goal

The drawer is a preview at a constrained width or a third of the screen, and it feels pushed to the side. It needs to be a first-class surface: big enough to read and edit a whole record, while keeping the board behind it, like an iOS sheet. A new page is not wanted. Find out:

- The width: half, two-thirds, centered, or a user-resizable panel that remembers its width; how each behaves on narrow screens and with the board still visible. Whether the board stays interactive behind it.
- The layout inside: fields in a header block, the body as sections, the log and outcome last; sticky header and actions; read mode and edit mode (one surface or two).
- Navigation: opening from a card, a link or the palette; going to a parent, child or linked record without closing it (a back stack); the URL for a drawer state so a link opens it.
- Keyboard and focus: open, close, move between records, save and cancel; unsaved changes and a stale-record warning.
- Constraints: the workbench has no dependencies and is plain JavaScript (see [261029-kb33](261029-kb33-workbench-form-controls.md)); light and dark; no reliance on color alone.

Done when:

- [x] Each width option is built as a throwaway page in the workbench's own styles and checked in the browser at three viewport sizes, light and dark, with screenshots.
- [x] A recommended layout and navigation model is in `## Outcome`, with the states the design must handle (read, edit, saving, conflict, error); follow-up items are filed, including changes to the existing drawer items.

## Next

## Outcome

Recommendation (as settled by Alexander after trying it in the workbench, 2026-11-19): a centered drawer over the board with a scrim, opened at up to 1000 px wide and about 92 percent of the viewport height, with the body before the field table.

Width and placement. Half and two-thirds, pinned right, both read well at 1440 px. Centered was tried against them in the live workbench and chosen: it reads as a sheet over the board, which is what the drawer is. It is up to 1000 px wide with 16 px of margin at the sides and 4 percent of the height above and below, with rounded corners, and it rises in instead of sliding. Under about 700 px it is a full-screen sheet. Every width is capped with min(100vw, ...): a fixed minimum above 390 px overflows a phone.

Layout inside. The 12-row field table fills the first screen at every width, so the Goal, Next and Log sit below the fold. Header: id, title, a one-line summary of the key fields (status, kind, parent, priority, claim), the actions. Then the body sections, with the log and outcome last. The full field table goes in a collapsed section that opens for editing. The header and actions stay in view while the body scrolls. Read and edit stay one surface, as they are now (the inline editors in fields.js).

Navigation. The record route is already a hash (#/record/id), so a link opens the drawer and browser Back is the back stack: the probe clicked a card, then a linked record, and Back returned to the previous one. Esc closes the drawer and returns focus to the card it came from. No new routing is needed.

Board behind the drawer. It stays covered by a scrim and inert: nothing behind the drawer is clickable and focus stays inside it. A non-modal drawer was probed (a click on a visible card switched the drawer to that record) and tried in the workbench with a toggle, and Alexander chose the scrim. The "keeps the board interactive" follow-up is dropped.

Keyboard and focus. Esc closes, j and k scroll, and focus moves to the title on open. These stay. Add: gg to the top and G to the bottom, with the board's keys and its 900 ms window (js/board.js), and a visible saving state (see below). The scrim already keeps focus inside the drawer.

States the design must handle: loading (a Loading line, exists); read (exists); edit (the fm-edit rows; the live refresh is suppressed while one is open, exists); saving (no visible state today: add a disabled control and a busy label); conflict (the stale-write message inline, exists, text in fields.js STALE_TEXT); error (a role=alert block, exists); removed (This record was removed, exists).

Bugs found: at 390 px wide the header squeezes the title to one or two characters per line; and the field table's name cell is display:flex on a th, which breaks the table layout (a wide gap, dividers that stop short, Edit out of line with its value).

Follow-ups, all children of the redesign project: the header at narrow widths (261119-4em1), the centered drawer and its width (261119-xzzp), the body-first layout and the table fix (261119-nbv8), and the gg and G keys. The board-interactive item (261119-7wtj) is dropped. The throwaway trial that produced the decision is not committed. The existing drawer items (261029-ctzc, 261029-paw7, 261028-esyt) are done and need no change; the body-first item reworks the field table they built.

Evidence: 30 screenshots (five variants, three widths, light and dark) and the scripts that took them are in .evidence/drawer-layout/ in the worktree, which Git ignores; they go on the pull request, not in the repository.

Decision: none (the follow-up items carry the choices; none is hard to reverse).

## Log

- 2026-10-31: filed.
- 2026-11-19 20:16-07:00: claimed (claude-sonnet-5-5-drawer-spike).
- 2026-11-19 20:33-07:00: Read the drawer (internal/server/web/js/detail.js, app.css .drawer). It is already modal: fixed 600px on the right, a scrim, the board inert, a record deep link at #/record/<id>, Esc closes. Built five variants as injected CSS over the live workbench (baseline, half, two-thirds, centered, resizable) and shot each at 1440, 1024 and 390 wide, light and dark. Findings so far: the 12-row field table fills the first screen and pushes the Goal, Next and Log below the fold at every width; at 390 wide the header row (title, Open in Obsidian, Close) squeezes the title to one or two characters per line; a fixed min-width over 390 overflows, so every width needs min(100vw, ...). Shots are in the session scratchpad, not yet in the repo.
- 2026-11-19 20:39-07:00: ticked "Each width option is built as a throwaway page in the workbench's own styles and checked in the browser at three viewport sizes, light and dark, with screenshots.": 30 screenshots and the scripts in .evidence/drawer-layout/ (gitignored; attached to the pull request)
- 2026-11-19 20:39-07:00: ticked "A recommended layout and navigation model is in `## Outcome`, with the states the design must handle (read, edit, saving, conflict, error); follow-up items are filed, including changes to the existing drawer items.": Outcome; follow-ups 261119-4em1, 261119-xzzp, 261119-7wtj, 261119-nbv8
- 2026-11-19 20:39-07:00: Probed a non-modal drawer, read the back stack and keyboard model, wrote the Outcome and filed four follow-ups.
- 2026-11-19 20:50-07:00: corrected outcome: Alexander tried centered with a scrim in the live workbench and chose it; the first recommendation (right-pinned, board kept live) is replaced.
