---
id: "261108-3t7t"
title: move closed reports inbound links
status: raw
kind: feature
priority: 4
created: 2026-11-08
updated: 2026-11-08
---

# move closed reports inbound links

## Goal

`move_file` now reports the inbound links it leaves broken (261029-m0vm); `move_closed`, which is a `move_file` per closed record, reports none. Moving closed records out of a folder breaks links to them in the same way.

Done when:

- [ ] `move_closed` lists the inbound links it leaves broken, in `notes`, in the form `move_file` uses, without one line per record growing past a bound (decide the bound, as `mint_ids` does for references).
- [ ] A test with a link to a record that moves.

## Log

- 2026-11-08: filed.
