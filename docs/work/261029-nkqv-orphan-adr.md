---
id: "261029-nkqv"
title: A new ADR made by the done wizard can be orphaned
summary: The done wizard creates the ADR first and links it after; if the later write fails the ADR stays unlinked.
status: active
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-11-16
claimed_by: claude-fable-5-2
claimed_at: 2026-11-16T14:05
---

# A new ADR made by the done wizard can be orphaned

## Goal

The done wizard creates the ADR first and links it after; if the later write fails the ADR stays unlinked.

Done when:

- [ ] A failed link leaves no orphan, or the wizard says so and offers to retry the link.

## Next

In the wizard's write loop (internal/server/web/js/wizard.js), keep the new ADR's id and the etag across a failure, say when the ADR is unlinked and offer a retry that makes only the remaining writes.

## Log

- 2026-10-29: filed.
- 2026-11-16 07:05-07:00: claimed (claude-fable-5-2).
