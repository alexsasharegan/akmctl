---
id: "261101-a86m"
title: workbench close cascade
summary: On close, list the open items that link to the closing one and offer, per item, to close it or to log that the linked item closed
status: raw
created: 2026-11-01
updated: 2026-11-01
kind: feature
priority: -1
blocked_by:
  - 261031-72bd
---

# workbench close cascade

## Goal

Most of the original idea is already covered by the tool, so this item is the workbench half only:

- 261030-dw71 (done): `close` reports the open items whose `blocked_by` names the closed one.
- 261031-72bd (ready): `close` of a child names its parent and the open siblings left, and says when the parent is now closable.

What is left, once 72bd lands: the workbench close prompt shows both lists (dependents and the parent's remaining children) and offers, per item, to open it, close it with its own outcome, or leave it. The decision to close each stays with Alexander. Open: how the lists reach the UI (the `close` API result carries them), and whether a free-text link scan is wanted (dw71 decided the structured part only).

Done when:

- [ ] the workbench close result lists the dependents and the parent's open siblings, each with an action to open it
- [ ] closing the last child shows that the parent can now be closed

## Log

- 2026-11-01: filed. It was pinned as "later" in `docs/design/pkm-workbench-design.md` (ADR 261025-r0ye) and never became a work item.
- 2026-11-01: narrowed to the workbench half and blocked by 261031-72bd; the CLI and API halves are dw71 (done) and 72bd.
