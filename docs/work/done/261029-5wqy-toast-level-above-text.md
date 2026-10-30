---
id: "261029-5wqy"
title: Toast level label sits above the message
status: done
created: 2026-10-29
updated: 2026-10-29
closed: 2026-10-29
kind: chore
priority: 1
summary: In a toast the level word is on its own line above the message, not inline before it
---

# Toast level label sits above the message

## Goal

Alexander saw a toast read as one continuous line, `Success: <text>`. The level word goes on its own line above the message.

Done when:

- [x] the level label is above the message text in every toast, with the glyph on the left and the close icon on the right
- [x] undo buttons and persistent toasts still lay out correctly
- [x] checked in the browser in light and dark

## Outcome

The toast markup wraps the level word and the message in a column (toast-main): the word is on its own line above the message, with the glyph on the left and the close icon on the right. Checked in the browser in dark and light, with the undo button inline in the message.

## Log

- 2026-10-29: filed.
- 2026-10-29: claimed (claude).
