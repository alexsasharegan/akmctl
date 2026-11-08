---
id: "261101-3erc"
title: block commit message names every blocker
status: done
created: 2026-11-01
updated: 2026-11-07
closed: 2026-11-07
kind: bug
summary: The block step's commit message names only the first blocked_by entry, so a three-blocker item reads as blocked on one
parent: 261101-t8h8
---

# block commit message names every blocker

## Goal

Found by the cold-agent trial of spike 261101-baze. `set status=blocked` prints `<id>(block): blocked on <first>` using `firstBlocker` (`internal/api/ops_write.go:777`), so an item with several `blocked_by` entries is described as waiting on one. Done when:

- [x] The message names every blocker (for example `blocked on a, b, c`), with a bounded length for a long list.
- [x] A test covers one, several and no blocker; any API fixture or spec text that pins the old message is updated in both fixture copies.

## Next

## Outcome

The block commit message names every blocker, bounded at 60 characters with 'and N more'; the first is always named. Spec and changelog updated; no fixture changed.

## Log

- 2026-11-01: filed.
- 2026-11-01 10:22-07:00: Filed from spike 261101-baze.
- 2026-11-07 23:25-07:00: claimed (claude-sonnet-5-5).
- 2026-11-07 23:26-07:00: ticked "The message names every blocker (for example `blocked on a, b, c`), with a bounded length for a long list.": blockerList in internal/api/ops_write.go; TestBlockCommitNamesBlockers
- 2026-11-07 23:26-07:00: ticked "A test covers one, several and no blocker; any API fixture or spec text that pins the old message is updated in both fixture copies.": TestBlockCommitNamesBlockers (one, several, bounded, long first) and TestBlockerListWithoutBlocker; no fixture pins the message, spec section 1 text updated
- 2026-11-07 23:26-07:00: Replaced firstBlocker with blockerList: every blocker, joined by comma, stops before 60 characters with 'and N more'. No fixture in either copy pins the message; only the spec text did.
