---
id: "261029-g8x2"
title: The review diff is against the current file, not after earlier writes
summary: Each dry-run diff in the wizard review is computed against the file as it is now, so a later write's diff ignores the earlier writes of the same move.
status: ready
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-10-29
---

# The review diff is against the current file, not after earlier writes

## Goal

Each dry-run diff in the wizard review is computed against the file as it is now, so a later write's diff ignores the earlier writes of the same move.

Done when:

- [ ] The review shows each diff as it would apply after the earlier writes.

## Log

- 2026-10-29: filed.
