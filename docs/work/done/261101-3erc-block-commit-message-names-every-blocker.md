---
id: "261101-3erc"
title: block commit message names every blocker
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
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

Inspect `firstBlocker` in `internal/api/ops_write.go` and the existing API write tests; update the block commit subject to list all blockers with a safe length limit, add the one/many/none test cases, then check for fixtures or spec text that pins the old output.

## Outcome

The block commit detail lists blockers in stored order and caps its length at 120 characters, reporting an omitted-entry count for long lists. `TestSetStatusChangeReturnsCommit` covers one and several entries; `TestBlockedByDetailBoundsLongListsAndHandlesNoBlockers` covers no entries and a long list. `docs/spec/akm-api-spec.md` documents the detail format. `go vet ./...`, `go test ./...`, `gofmt -l .`, `npm run format:check`, `npm run schemas`, and `/tmp/akmctl check` pass. Decision: none.

## Log

- 2026-11-01: filed.
- 2026-11-01 10:22-07:00: Filed from spike 261101-baze.
- 2026-11-01 16:37-07:00: claimed (codex).
- 2026-11-01 16:37-07:00: claim activated this item without a Next section
- 2026-11-01 16:37-07:00: claim activated this item without a Next section; recorded WK6 failure as 261101-2n55
- 2026-11-01 16:39-07:00: ticked "The message names every blocker (for example `blocked on a, b, c`), with a bounded length for a long list.": TestSetStatusChangeReturnsCommit covers one and several blockers; TestBlockedByDetailBoundsLongListsAndHandlesNoBlockers covers empty and long lists.
- 2026-11-01 16:39-07:00: ticked "A test covers one, several and no blocker; any API fixture or spec text that pins the old message is updated in both fixture copies.": go test ./... passes; no API fixture pins the old detail, and docs/spec/akm-api-spec.md now defines the bounded list.
