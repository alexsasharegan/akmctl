---
id: "261031-72bd"
title: rollup counts for parent items
status: raw
created: 2026-10-31
updated: 2026-10-31
summary: Children done / total, direct and transitive, in show, ls --tree and the API
---

# rollup counts for parent items

## Goal

From spike 261031-078w. A parent shows no progress. Add counts (`children_done`, `children_total`, and the transitive pair) computed from `parent` links, never stored (one source of truth, H4). Done when:

- [ ] `show` and `ls --tree` print the counts; the API returns them as computed fields.
- [ ] Dropped children are counted in a separate figure, not as done.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
