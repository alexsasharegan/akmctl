---
id: "261119-4em1"
title: drawer header narrow widths
status: done
kind: bug
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
closed: 2026-11-19
summary: The drawer header squeezes the title to one or two characters per line at phone widths
---

# drawer header narrow widths

## Goal

At 390 px wide the drawer header holds the title, "Open in Obsidian" and Close in one row, and the title wraps to one or two characters per line (found by the drawer layout spike, 261031-18t8). The header must stack at narrow widths so the title keeps the full width.

Done when:

- [x] At 390 px wide, light and dark, the title reads on whole words and the two buttons sit on their own row.

## Next

## Outcome

1. At 390 px wide, light and dark, the title reads on whole words and the buttons sit on their own rows: checked in the built tip and by Alexander in the running build. At desktop width the Open in Obsidian link no longer wraps (66 px tall before, 42 px after).
2. The screenshot done-when box was removed at Alexander's request (see the Log): he inspected the running build himself and decided screenshots need not go on the pull request.

Not part of this item: the field table below the header is also cramped at 390 px (the value column is squeezed); it is not filed as an item.

Decision: none.

## Log

- 2026-11-19: filed.
- 2026-11-19 22:07-07:00: claimed (claude-sonnet-5-5-drawer-header).
- 2026-11-19 22:07-07:00: ticked "At 390 px wide, light and dark, the title reads on whole words and the two buttons sit on their own row.": checked in the built tip at 390 px wide, light and dark: the title reads on whole words and the two buttons sit on their own rows; 700 and 1440 px unchanged
- 2026-11-19 22:07-07:00: removed "A screenshot at 390 px, light and dark, is attached to the pull request.": Alexander decided screenshots need not go on the pull request (2026-11-19); he inspected the build himself
- 2026-11-19 22:07-07:00: Stacked the drawer header as rows under 700 px, scoped to .drawer; changelog entry added. The field table below is also cramped at 390 px (value column squeezed); not part of this item.
- 2026-11-19 22:30-07:00: Follow-up from Alexander inspecting the tip: Open in Obsidian wrapped onto two lines at desktop width (link 66 px tall, now 42 px). The header buttons are white-space:nowrap.
