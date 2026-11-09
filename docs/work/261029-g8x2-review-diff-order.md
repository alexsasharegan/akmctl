---
id: "261029-g8x2"
title: The review diff is against the current file, not after earlier writes
summary: Each dry-run diff in the wizard review is computed against the file as it is now, so a later write's diff ignores the earlier writes of the same move.
status: active
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-11-09
claimed_by: claude-sonnet-5-5-g8x2
claimed_at: 2026-11-09T17:37
---

# The review diff is against the current file, not after earlier writes

## Goal

Each dry-run diff in the wizard review is computed against the file as it is now, so a later write's diff ignores the earlier writes of the same move.

Done when:

- [ ] The review shows each diff as it would apply after the earlier writes.

## Next

Draft a proposed ADR for a read-only preview operation that applies an ordered list of writes on a scratch copy and returns each write's dry-run response; then build it (api, HTTP, CLI, spec, fixtures) and make the move wizard's review use it.

## Log

- 2026-10-29: filed.
- 2026-11-09 10:37-07:00: claimed (claude-sonnet-5-5-g8x2).
