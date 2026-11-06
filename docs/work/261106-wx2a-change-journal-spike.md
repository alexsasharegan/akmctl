---
id: "261106-wx2a"
title: "spike: an append-only change journal for moves"
status: raw
created: 2026-11-06
updated: 2026-11-06
kind: spike
---

# spike: an append-only change journal for moves

## Goal

ADR 261106-60xv chose git as the move log and rejected an append-only journal file for now. Write down what a journal would need so the choice can be revisited on facts: where the file lives (tracked or ignored), the event shape, how parallel branches merge, how hand edits are reconciled with the snapshot diff, and what retention and compaction cost (261031-dkv2, 261101-wdbj). Done when:

- [ ] A note in `docs/design/` answers those five questions with one recommendation.
- [ ] It names the trigger that would justify building it (for example undo needing uncommitted changes, 261029-shz5).

## Log

- 2026-11-06: filed.
