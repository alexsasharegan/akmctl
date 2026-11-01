---
id: "261101-3erc"
title: block commit message names every blocker
status: raw
created: 2026-11-01
updated: 2026-11-01
kind: bug
summary: The block step's commit message names only the first blocked_by entry, so a three-blocker item reads as blocked on one
---

# block commit message names every blocker

## Goal

Found by the cold-agent trial of spike 261101-baze. `set status=blocked` prints `<id>(block): blocked on <first>` using `firstBlocker` (`internal/api/ops_write.go:777`), so an item with several `blocked_by` entries is described as waiting on one. Done when:

- [ ] The message names every blocker (for example `blocked on a, b, c`), with a bounded length for a long list.
- [ ] A test covers one, several and no blocker; any API fixture or spec text that pins the old message is updated in both fixture copies.

## Log

- 2026-11-01: filed.
- 2026-11-01 10:22-07:00: Filed from spike 261101-baze.
