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

## Next

Verify and merge the branch, then backfill `closed` on this repository's done items (from the history, once) and raise `docs/.layout` to 0.3.0 in the same step; Alexander approves the tag. Open: whether the README status line should read 0.3.0, and whether to embed a 0.2.0-to-0.3.0 migration (today a 0.2.0 manifest gets the S1 "no migration path" warning).

## Log

- 2026-10-29: filed.
- 2026-10-29: claimed (claude).
- 2026-10-29: implemented on a worktree branch: closed/cycle/keep in close, set, ls (--cycle), show and the board; check rules WK13 to WK17; layout and extensions 0.3.0; schemas, fixtures (both copies), spec, guide, API spec, skill references, changelog. Backfill of done items not done.
