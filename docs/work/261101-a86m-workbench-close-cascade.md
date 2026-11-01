---
id: "261101-a86m"
title: workbench close cascade
summary: On close, list the open items that link to the closing one and offer, per item, to close it or to log that the linked item closed
status: raw
created: 2026-11-01
updated: 2026-11-01
kind: feature
priority: -1
---

# workbench close cascade

## Goal

On close, list the open work items that link to the closing one (a programmatic search of the records for its id) and offer, per item, to close it with its own outcome or to log that the linked item closed. The decision to close each one stays with Alexander. Check it against 261101-1smq (close repairs links to moved items) and `blocked_by`, which may already cover part of it.

## Log

- 2026-11-01: filed. It was pinned as "later" in `docs/design/pkm-workbench-design.md` (ADR 261025-r0ye) and never became a work item.
