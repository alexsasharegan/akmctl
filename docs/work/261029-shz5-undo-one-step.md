---
id: "261029-shz5"
title: Undo covers only the last safe move or capture
summary: Undo handles one step; gated (wizard) moves and earlier moves cannot be undone.
status: ready
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-10-29
---

# Undo covers only the last safe move or capture

## Goal

Undo handles one step; gated (wizard) moves and earlier moves cannot be undone.

Done when:

- [ ] Undo covers wizard moves, or the limit is documented in the key help.

Design note (2026-10-29): consider one local operation journal for this item and 261029-xtqj / 261029-shz5 together: a gitignored, append-only intent log (what is about to change, then done), written under the existing lock, that lets a half-finished multi-write be replayed or reversed and lets undo reach back further than the last operation. It is a change journal, not a source of truth: writers outside the tool (Obsidian, an agent's edit tool, git) bypass it. A by-product is a cheap hint for the workbench to skip scans after its own writes. Open: the file name (it sits beside `.akm.lock`), rotation, and what undo means across it. Needs an ADR.

## Log

- 2026-10-29: filed.
