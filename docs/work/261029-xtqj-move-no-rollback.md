---
id: "261029-xtqj"
title: Multi-write moves can stop half-way with no rollback
summary: A wizard move makes several writes in turn; if a later one fails the earlier ones stay and nothing is rolled back.
status: active
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-11-16
claimed_by: claude-fable-5-2-xtqj
claimed_at: 2026-11-17T02:02
---

# Multi-write moves can stop half-way with no rollback

## Goal

A wizard move makes several writes in turn; if a later one fails the earlier ones stay and nothing is rolled back.

Done when:

- [ ] Decide how a move that fails half-way is undone or reported, and implement it.

Design note (2026-10-29): consider one local operation journal for this item and 261029-xtqj / 261029-shz5 together: a gitignored, append-only intent log (what is about to change, then done), written under the existing lock, that lets a half-finished multi-write be replayed or reversed and lets undo reach back further than the last operation. It is a change journal, not a source of truth: writers outside the tool (Obsidian, an agent's edit tool, git) bypass it. A by-product is a cheap hint for the workbench to skip scans after its own writes. Open: the file name (it sits beside `.akm.lock`), rotation, and what undo means across it. Needs an ADR.

## Next

The apply operation, its CLI and HTTP bindings, fixtures, spec section 4.27, schema, skill reference, ADR 261116-2qyf (proposed) and the workbench switch are written and uncommitted-then-committed as a checkpoint. Left: run the full check (go vet, go test -race ./..., gofmt -l, /tmp/akmctl check, npm run format:check, npm run schemas), tick the Goal box, write the Outcome, list the ADR in decisions, close, akmctl commits --range origin/main..HEAD, push and open the pull request.

## Log

- 2026-10-29: filed.
- 2026-11-16 19:02-07:00: claimed (claude-fable-5-2-xtqj).
- 2026-11-16 19:13-07:00: Built apply (one lock, restore from memory, no journal) and moved the wizard, safe moves and undo to it. Checked by hand in the workbench: a close with a new ADR whose file move was refused left git status clean and the wizard said nothing was written. go test ./... passed before the docs pass; the race run and the final full check are still to do.
