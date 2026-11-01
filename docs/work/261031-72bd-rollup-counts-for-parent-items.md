---
id: "261031-72bd"
title: rollup counts for parent items
status: ready
created: 2026-10-31
updated: 2026-11-01
summary: Children done / total, direct and transitive, in show, ls --tree and the API
parent: 261101-t8h8
---

# rollup counts for parent items

## Goal

From spike 261031-078w. A parent shows no progress. Add counts (`children_done`, `children_total`, and the transitive pair) computed from `parent` links, never stored (one source of truth, H4). Done when:

- [ ] `show` and `ls --tree` print the counts; the API returns them as computed fields.
- [ ] Dropped children are counted in a separate figure, not as done.
- [ ] `close` of a child names its parent and the open siblings left (`parent`, `open_siblings` in JSON), and says when the parent is now closable (from spike 261101-baze, which found `close` is silent about the parent, so a closer does not learn that a split is finished).

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
- 2026-11-01 10:14-07:00: Added a box for close reporting the parent and open siblings, from spike 261101-baze.
