---
id: "261029-shz5"
title: Undo covers only the last safe move or capture
summary: Undo handles one step; gated (wizard) moves and earlier moves cannot be undone.
status: active
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-11-16
claimed_by: claude-fable-5-2
claimed_at: 2026-11-16T14:06
---

# Undo covers only the last safe move or capture

## Goal

Undo handles one step; gated (wizard) moves and earlier moves cannot be undone.

Done when:

- [x] Undo covers wizard moves, or the limit is documented in the key help.

Design note (2026-10-29): consider one local operation journal for this item and 261029-xtqj / 261029-shz5 together: a gitignored, append-only intent log (what is about to change, then done), written under the existing lock, that lets a half-finished multi-write be replayed or reversed and lets undo reach back further than the last operation. It is a change journal, not a source of truth: writers outside the tool (Obsidian, an agent's edit tool, git) bypass it. A by-product is a cheap hint for the workbench to skip scans after its own writes. Open: the file name (it sits beside `.akm.lock`), rotation, and what undo means across it. Needs an ADR.

## Next

State the limit of undo in the key help (flowHelp in internal/server/web/js/flow.js) and in the toast u shows after a wizard move; update the glossary.

## Outcome

- The limit is documented in the key help: the `u` row of `flowHelp` says one step, safe moves only, and that a wizard move is not undone.
- `u` after a wizard move now says that undo does not cover it and what to do, where it said "Nothing to undo".
- Proof: two assertions in `internal/server/testdata/keys.test.mjs`; the glossary's Move row carries the limit.

Undo itself did not change, so the 1.0.0 criterion "the last operation can be undone" holds for safe moves and captures only. Not checked in a real browser.

Decision: none (the pick is in the Log).

## Log

- 2026-10-29: filed.
- 2026-11-16 07:06-07:00: claimed (claude-fable-5-2).
- 2026-11-16 07:06-07:00: ticked "Undo covers wizard moves, or the limit is documented in the key help.": flowHelp row for u and the toast after a wizard move in flow.js; keys.test.mjs asserts both; glossary Move
- 2026-11-16 07:06-07:00: pick: document the limit, not read git history. A wizard move can move a file, add Log lines and make an ADR; undoing it from git history is a revert of a commit the tool does not make, so it is left to git. Roadmap criterion 2 says the last operation can be undone: with this pick that holds for safe moves and captures only, and the criterion's wording is the maintainer's to keep or change
