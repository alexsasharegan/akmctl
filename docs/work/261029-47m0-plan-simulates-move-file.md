---
id: "261029-47m0"
title: "migrate plan: simulate move_file steps"
status: ready
created: 2026-10-29
updated: 2026-10-29
kind: bug
priority: 2
summary: plan counted ADR Candidates as an ADR (22 ids) because it did not apply move_file steps; apply made 21.
---

# migrate plan: simulate move_file steps

## Goal

plan applies the steps on a scratch copy but plans move_file and add_file against the real tree, so a file moved out of a directory is still counted.

Done when:

- [ ] plan and apply agree on the number of ids and renames on the vault fixture.
- [ ] A test covers a move_file before a mint_ids on the same directory.

## Log

- 2026-10-29: filed.
