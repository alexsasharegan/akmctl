---
id: "261031-72bd"
title: rollup counts for parent items
status: done
created: 2026-10-31
updated: 2026-11-07
closed: 2026-11-07
summary: Children done / total, direct and transitive, in show, ls --tree and the API
parent: 261101-t8h8
---

# rollup counts for parent items

## Goal

From spike 261031-078w. A parent shows no progress. Add counts (`children_done`, `children_total`, and the transitive pair) computed from `parent` links, never stored (one source of truth, H4). Done when:

- [x] `show` and `ls --tree` print the counts; the API returns them as computed fields.
- [x] Dropped children are counted in a separate figure, not as done.
- [x] `close` of a child names its parent and the open siblings left (`parent`, `open_siblings` in JSON), and says when the parent is now closable (from spike 261101-baze, which found `close` is silent about the parent, so a closer does not learn that a split is finished).

## Next

## Outcome

show and ls --tree carry a computed progress member (children and descendants: done, dropped, total); close of a child reports parent, open_siblings and parent_closable. ls gained --tree. Nothing is stored. 261101-a86m is unblocked by this item.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
- 2026-11-01 10:14-07:00: Added a box for close reporting the parent and open siblings, from spike 261101-baze.
- 2026-11-01: 261101-a86m (workbench close prompt) is blocked by this item's close reporting.
- 2026-11-07 23:46-07:00: claimed (claude-sonnet-5-5).
- 2026-11-07 23:46-07:00: ticked "`show` and `ls --tree` print the counts; the API returns them as computed fields.": progress member in show and ls --tree (TestShowProgressCountsChildren, TestLsTreeOrdersChildrenUnderParent); computed in internal/api/progress.go, never stored
- 2026-11-07 23:46-07:00: ticked "Dropped children are counted in a separate figure, not as done.": children_dropped and descendants_dropped are separate from *_done; TestShowProgressCountsChildren
- 2026-11-07 23:46-07:00: ticked "`close` of a child names its parent and the open siblings left (`parent`, `open_siblings` in JSON), and says when the parent is now closable (from spike 261101-baze, which found `close` is silent about the parent, so a closer does not learn that a split is finished).": close returns parent, open_siblings, parent_closable; TestCloseNamesParentAndOpenSiblings; close-2..8 fixtures
- 2026-11-07 23:46-07:00: Named the member progress, not rollup: the glossary already lists rollup as an alias of Compaction. Added ls --tree (it did not exist) with depth. 261101-a86m can now read parent_closable from close.
