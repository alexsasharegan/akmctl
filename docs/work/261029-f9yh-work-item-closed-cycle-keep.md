---
id: "261029-f9yh"
title: Add closed, cycle and keep to work items (layout 0.3.0)
status: active
created: 2026-10-29
updated: 2026-10-29
kind: feature
priority: 1
summary: Add closed (set by close), cycle (opaque label) and keep (pin) to work items; layout, extensions and spec to 0.3.0; ADR 261029-sg82
claimed_by: claude
claimed_at: 2026-10-30T03:54
---

# Add closed, cycle and keep to work items (layout 0.3.0)

## Goal

Work items gain `closed`, `cycle` and `keep` as specified in ADR 261029-sg82, and the layout, extensions and spec read 0.3.0.

Done when:

- [ ] `close` writes `closed` (done and dropped); no other operation does; `set` refuses it
- [ ] `set` accepts `cycle` and `keep`; `show`, `ls` and the board can show them
- [ ] check rules for the three fields; layout, extensions and spec at 0.3.0
- [ ] schemas, fixtures (both copies), spec, guide and skill references updated
- [ ] this repository's done items backfilled with `closed`

## Log

- 2026-10-29: filed.
- 2026-10-29: claimed (claude).
