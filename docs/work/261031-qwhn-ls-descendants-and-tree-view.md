---
id: "261031-qwhn"
title: ls descendants and tree view
status: raw
created: 2026-10-31
updated: 2026-11-01
summary: ls --descendants and a --tree print, with depth in JSON
---

# ls descendants and tree view

## Goal

From spike 261031-078w. `ls --parent` lists direct children only. Add `--descendants <id>` (whole subtree) and `--tree` (indented print); JSON records gain `depth`. Done when:

- [ ] `ls --descendants <id>` and `ls --tree` work on a three-level tree and respect the other filters (`--status`, `--sort`).
- [ ] The API `ls` operation and the reference page document the parameters.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
- 2026-11-01 10:22-07:00: Spike 261101-baze: no change. The split procedure needs only ls --parent; a tree view helps a reader but is not on its path.
