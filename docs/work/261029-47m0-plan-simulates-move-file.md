---
id: "261029-47m0"
title: "migrate plan: simulate move_file steps"
status: active
created: 2026-10-29
updated: 2026-11-09
kind: bug
priority: 2
summary: plan counted ADR Candidates as an ADR (22 ids) because it did not apply move_file steps; apply made 21.
claimed_by: claude-sonnet-5-5-47m0
claimed_at: 2026-11-09T16:34
---

# migrate plan: simulate move_file steps

## Goal

plan applies the steps on a scratch copy but plans move_file and add_file against the real tree, so a file moved out of a directory is still counted.

Done when:

- [ ] plan and apply agree on the number of ids and renames on the vault fixture.
- [ ] A test covers a move_file before a mint_ids on the same directory.

## Next

Plan move_file and add_file on the scratch copy, mirroring a non-Markdown source or target from the real tree on demand, so later steps see the moved files; add a test (move_file before mint_ids on one directory) and a plan/apply agreement check.

## Log

- 2026-10-29: filed.
- 2026-11-09 09:34-07:00: claimed (claude-sonnet-5-5-47m0).
