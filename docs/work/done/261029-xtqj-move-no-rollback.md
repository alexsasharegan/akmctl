---
id: "261029-xtqj"
title: Multi-write moves can stop half-way with no rollback
summary: A wizard move makes several writes in turn; if a later one fails the earlier ones stay and nothing is rolled back.
status: done
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-11-16
closed: 2026-11-16
decisions: [261116-2qyf]
---

# Multi-write moves can stop half-way with no rollback

## Goal

A wizard move makes several writes in turn; if a later one fails the earlier ones stay and nothing is rolled back.

Done when:

- [x] Decide how a move that fails half-way is undone or reported, and implement it.

Design note (2026-10-29): consider one local operation journal for this item and 261029-xtqj / 261029-shz5 together: a gitignored, append-only intent log (what is about to change, then done), written under the existing lock, that lets a half-finished multi-write be replayed or reversed and lets undo reach back further than the last operation. It is a change journal, not a source of truth: writers outside the tool (Obsidian, an agent's edit tool, git) bypass it. A by-product is a cheap hint for the workbench to skip scans after its own writes. Open: the file name (it sits beside `.akm.lock`), rotation, and what undo means across it. Needs an ADR.

## Next

## Outcome

Decided: a move that fails half-way is undone, not only reported. The new API operation apply (spec section 4.27, POST /apply, akmctl apply --writes) makes a list of writes under one repository lock and, when one fails, restores every file from the bytes read before the first write; a new file is removed and a moved file is put back. No journal, so a crash between two writes is not covered (the journal spike 261106-wx2a keeps that). The workbench wizard, a safe move that releases or claims, and undo send their writes through it; Retry the rest is gone.

Proof for the done-when line: TestApplyRollsBackWhenAWriteFails and TestApplyMakesTheWritesInOrder (internal/api), TestApplyOverHTTP (internal/server), the fixtures apply-1 to apply-3 through the API and HTTP bindings, rollback.test.mjs and wizard.test.mjs for the page, and a run by hand in the workbench where a close with a new ADR was refused at the file move and git status stayed clean.

Decision: ADR 261116-2qyf (proposed), for the maintainer to accept at merge.

## Log

- 2026-10-29: filed.
- 2026-11-16 19:02-07:00: claimed (claude-fable-5-2-xtqj).
- 2026-11-16 19:13-07:00: Built apply (one lock, restore from memory, no journal) and moved the wizard, safe moves and undo to it. Checked by hand in the workbench: a close with a new ADR whose file move was refused left git status clean and the wizard said nothing was written. go test ./... passed before the docs pass; the race run and the final full check are still to do.
- 2026-11-16 23:11-07:00: ticked "Decide how a move that fails half-way is undone or reported, and implement it.": apply operation (ADR 261116-2qyf, proposed); TestApplyRollsBackWhenAWriteFails, fixtures apply-1 to apply-3, rollback.test.mjs
- 2026-11-16 23:11-07:00: Full check passed: go vet, go test -race ./..., gofmt, GOOS=windows go vet, akmctl check (0 errors), format:check, schemas (147 cases).
