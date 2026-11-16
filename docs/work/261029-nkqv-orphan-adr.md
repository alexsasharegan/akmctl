---
id: "261029-nkqv"
title: A new ADR made by the done wizard can be orphaned
summary: The done wizard creates the ADR first and links it after; if the later write fails the ADR stays unlinked.
status: active
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-11-16
claimed_by: claude-fable-5-2
claimed_at: 2026-11-16T14:05
---

# A new ADR made by the done wizard can be orphaned

## Goal

The done wizard creates the ADR first and links it after; if the later write fails the ADR stays unlinked.

Done when:

- [x] A failed link leaves no orphan, or the wizard says so and offers to retry the link.

## Next

In the wizard's write loop (internal/server/web/js/wizard.js), keep the new ADR's id and the etag across a failure, say when the ADR is unlinked and offer a retry that makes only the remaining writes.

## Outcome

- A failed link no longer goes unsaid: the wizard names the new ADR that is not linked and offers **Retry the rest**, which writes the link and the close and does not make the ADR again (`write(resume)` and `orphanNote` in `internal/server/web/js/wizard.js`, `unlinkedAdr` in `moves.js`).
- Proof: `internal/server/testdata/orphan.test.mjs` and three assertions in `logic.test.mjs`.

Not checked in a real browser. Going Back after a partial write and writing again still starts the list over; rollback of partial moves is the item on moves that stop half-way.

Decision: none (the pick is in the Log).

## Log

- 2026-10-29: filed.
- 2026-11-16 07:05-07:00: claimed (claude-fable-5-2).
- 2026-11-16 07:05-07:00: ticked "A failed link leaves no orphan, or the wizard says so and offers to retry the link.": wizard.js write(resume) and orphanNote; orphan.test.mjs: the link fails, the wizard names the ADR, Retry links it and makes no second ADR
- 2026-11-16 07:05-07:00: pick: say so and offer a retry. Writing the link first is not possible: the id exists only after new. The retry resumes the write list where it stopped, for every move, not only this one
- 2026-11-16 07:05-07:00: not verified: in a real browser. The DOM test runs under Node against the fake document. Known gap: going Back after a partial write and writing again from the review still starts over (a second new ADR); the note names the first ADR, so it is findable
