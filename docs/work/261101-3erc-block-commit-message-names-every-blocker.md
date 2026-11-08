---
id: "261101-3erc"
title: block commit message names every blocker
status: active
created: 2026-11-01
updated: 2026-11-07
kind: bug
summary: The block step's commit message names only the first blocked_by entry, so a three-blocker item reads as blocked on one
parent: 261101-t8h8
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-08T06:25
---

# block commit message names every blocker

## Goal

Found by the cold-agent trial of spike 261101-baze. `set status=blocked` prints `<id>(block): blocked on <first>` using `firstBlocker` (`internal/api/ops_write.go:777`), so an item with several `blocked_by` entries is described as waiting on one. Done when:

- [ ] The message names every blocker (for example `blocked on a, b, c`), with a bounded length for a long list.
- [ ] A test covers one, several and no blocker; any API fixture or spec text that pins the old message is updated in both fixture copies.

## Next

Replace firstBlocker in internal/api/ops_write.go with a function that joins every blocker, bounded in length; add a test for one, several and no blocker; update both fixture copies and any spec text.

## Log

- 2026-11-01: filed.
- 2026-11-01 10:22-07:00: Filed from spike 261101-baze.
- 2026-11-07 23:25-07:00: claimed (claude-sonnet-5-5).
