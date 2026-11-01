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

- [x] The Log lists where the nudge can live (a new skill, a section of `akm-track-work`, a hint from `new` or `check` on an oversized item, the guide) with the cost and reach of each, and names one or a mix.
- [x] The procedure for splitting an existing item is written down and walked through once on a real item: what stays in the parent's Goal, what moves to children, how the parent's history and links are kept, and when a spike's parent closes.
- [ ] A prototype of the chosen surface exists (a draft skill or hint) and was tried by a cold agent on an oversized item; the Log records what it got wrong.
- [ ] The Log says how 261031-ph53 and 261031-qwhn change, if at all, and any term the work adds is registered in `docs/glossary.md`.
- [ ] An implementation work item is filed (`akmctl new work`) whose Goal carries the acceptance criteria, adjusted for what the prototype showed.

## Next

Box 3: put the procedure where the cold agent will find it (a draft section in akm-track-work, no new skill file), then have a fresh agent split an oversized item in a scratch copy of the repo and log what it got wrong.

## Log

- 2026-11-01: filed.
- 2026-11-01 09:43-07:00: claimed (claude).
- 2026-11-01 09:43-07:00: Survey: the guide (docs/spec/akm-guide.md:82,122,182) mentions parents only as a field and a closing rule; no skill, command output or check rule says 'split an oversized item'. akm-track-work has 'No item fits' (file a raw item) but no 'item is too big' branch. Candidate surfaces to weigh: (a) new skill akm-decompose-work, (b) a section in akm-track-work, (c) a guide chapter, (d) a hint from new/check/ls, (e) split command (261031-ph53). Not yet weighed.
- 2026-11-01 09:46-07:00: ticked "The Log lists where the nudge can live (a new skill, a section of `akm-track-work`, a hint from `new` or `check` on an oversized item, the guide) with the cost and reach of each, and names one or a mix.": Log entry weighing surfaces a to e
- 2026-11-01 09:46-07:00: Weighed the surfaces. (a) New skill akm-decompose-work: found by description when an agent says 'too big'; a user who never loads skills never sees it. (b) Section in akm-track-work: cheapest and rides the skill every agent loads, but buries a procedure in a lifecycle skill. (c) Guide chapter: reaches humans, not agents at the moment of need. (d) Hint from new/check/ls: reaches everyone at the moment of need, but 'oversized' is a judgment the tool cannot make; the only safe trigger is kind: spike, and a nag on every spike is noise. (e) split (ph53): the mechanism the others point at, not a surface. Pick: (a) a small skill for the procedure, with a one-line pointer in akm-track-work 'No item fits' and in the guide, and (e) as the tool it calls. Skip (d) unless the cold-agent trial shows the pointers are missed.
- 2026-11-01 09:56-07:00: Drafted the procedure and walked it through on 261030-456q as a dry run (nothing changed): docs/design/decomposing-work.md. Findings: a split is 3 commands per child (new, set, section) until split (ph53) exists; moving a done-when box between items has no command and means retyping the Goal; children do not inherit priority; releasing a claim on the new parent is not mentioned anywhere. Cost of a new skill: it is mirrored in internal/initrepo/templates/skills and listed in 6 init fixtures (x2 copies), make_fixtures.py, initrepo_test.go and skills_test.go; a section in akm-track-work avoids all of that, so the pick leans to (b) unless the maintainer wants a standalone skill.
- 2026-11-01 09:56-07:00: ticked "The procedure for splitting an existing item is written down and walked through once on a real item: what stays in the parent's Goal, what moves to children, how the parent's history and links are kept, and when a spike's parent closes.": docs/design/decomposing-work.md, walk-through on 261030-456q
