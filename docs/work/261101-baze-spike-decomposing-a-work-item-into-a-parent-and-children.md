---
id: "261101-baze"
title: spike decomposing a work item into a parent and children
status: active
created: 2026-11-01
updated: 2026-11-01
kind: spike
summary: how a new user or agent finds out that an oversized item should become a parent with children, and the procedure for it
claimed_by: claude
claimed_at: 2026-11-01T16:43
---

# spike decomposing a work item into a parent and children

## Goal

A spike: settle the design and prove it with a prototype, then file the implementation as its own work item. The question: when a work item turns out too large (a multi-part spike, a feature with several steps), how does a newcomer or an agent discover that the answer is a parent record with child records, and how do they do the split? The maintainer knows this by heart; the workflow has no surface that says so (`akm-track-work` covers claim, work and close only).

Related, already filed: 261031-ph53 (`new --parent --kind` and `split`), 261031-qwhn (`ls --descendants` and `--tree`). Settled mechanics: spike 261031-078w.

### Done when

- [ ] The Log lists where the nudge can live (a new skill, a section of `akm-track-work`, a hint from `new` or `check` on an oversized item, the guide) with the cost and reach of each, and names one or a mix.
- [ ] The procedure for splitting an existing item is written down and walked through once on a real item: what stays in the parent's Goal, what moves to children, how the parent's history and links are kept, and when a spike's parent closes.
- [ ] A prototype of the chosen surface exists (a draft skill or hint) and was tried by a cold agent on an oversized item; the Log records what it got wrong.
- [ ] The Log says how 261031-ph53 and 261031-qwhn change, if at all, and any term the work adds is registered in `docs/glossary.md`.
- [ ] An implementation work item is filed (`akmctl new work`) whose Goal carries the acceptance criteria, adjusted for what the prototype showed.

## Next

Read how akm-track-work, new, split (ph53) and the guide surface parent/child today, then list the candidate surfaces in the Log.

## Log

- 2026-11-01: filed.
- 2026-11-01 09:43-07:00: claimed (claude).
- 2026-11-01 09:43-07:00: Survey: the guide (docs/spec/akm-guide.md:82,122,182) mentions parents only as a field and a closing rule; no skill, command output or check rule says 'split an oversized item'. akm-track-work has 'No item fits' (file a raw item) but no 'item is too big' branch. Candidate surfaces to weigh: (a) new skill akm-decompose-work, (b) a section in akm-track-work, (c) a guide chapter, (d) a hint from new/check/ls, (e) split command (261031-ph53). Not yet weighed.
