---
id: "261119-4em1"
title: drawer header narrow widths
status: active
kind: bug
parent: 261031-gtw6
created: 2026-11-19
updated: 2026-11-19
summary: The drawer header squeezes the title to one or two characters per line at phone widths
claimed_by: claude-sonnet-5-5-drawer-header
claimed_at: 2026-11-20T05:07
---

# drawer header narrow widths

## Goal

At 390 px wide the drawer header holds the title, "Open in Obsidian" and Close in one row, and the title wraps to one or two characters per line (found by the drawer layout spike, 261031-18t8). The header must stack at narrow widths so the title keeps the full width.

Done when:

- [ ] At 390 px wide, light and dark, the title reads on whole words and the two buttons sit on their own row.
- [ ] A screenshot at 390 px, light and dark, is attached to the pull request.

## Next

Stack the drawer header (title block, Open in Obsidian, Close) as rows under 700 px in app.css, scoped to .drawer; changelog entry; run the tip locally for Alexander to inspect.

## Log

- 2026-11-19: filed.
- 2026-11-19 22:07-07:00: claimed (claude-sonnet-5-5-drawer-header).
