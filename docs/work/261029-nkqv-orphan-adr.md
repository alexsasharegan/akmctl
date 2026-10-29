---
id: "261029-nkqv"
title: A new ADR made by the done wizard can be orphaned
summary: The done wizard creates the ADR first and links it after; if the later write fails the ADR stays unlinked.
status: ready
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-10-29
---

# A new ADR made by the done wizard can be orphaned

## Goal

The done wizard creates the ADR first and links it after; if the later write fails the ADR stays unlinked.

Done when:

- [ ] A failed link leaves no orphan, or the wizard says so and offers to retry the link.

## Log

- 2026-10-29: filed.
