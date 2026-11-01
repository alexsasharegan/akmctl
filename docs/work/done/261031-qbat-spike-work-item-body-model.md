---
id: "261031-qbat"
title: spike work item body model
summary: "Spike: what a work item body contains, how much is schema versus open text, and how extension works"
status: done
kind: spike
parent: "261031-gtw6"
created: 2026-10-31
updated: 2026-10-31
closed: 2026-10-31
decisions: [261031-gam0]
---

# spike work item body model

## Goal

Work item bodies already follow a convention: `## Goal` (prose plus a "Done when" checklist), `## Next`, `## Log` (dated lines), `## Outcome`. Treat that as a schema and find out where it should be strict and where open. Questions:

- Which sections exist per status and kind (a spike has an `## Outcome`; a bug may want repro steps; a feature a checklist). Which are required, optional or repeatable.
- The node types inside a section: prose, checklist item (with a done state and an optional commit), dated log entry, link list, table. Whether each is a typed node or opaque text.
- Open versus closed: the options to compare are (a) a fixed set of sections plus a free-form `notes` area, (b) fixed core sections plus extension sections declared by name, kind and node type, (c) any `##` heading allowed with only the core ones typed. Include what an Obsidian user typing by hand can still do, and what survives a round trip unchanged.
- The model in the data: a body as an ordered list of typed sections (the in-memory form the writer and the drawer both use), with a stable id per section and per checklist item so edits address them without matching text.
- How it relates to frontmatter: which facts belong in a field (priority, project) and which in the body (goal text); whether a "done when" item is data other rules read (see [261031-zek2](../261031-zek2-close-checks-goal-checkboxes.md)).

Done when:

- [x] The current bodies in `docs/work/` and `docs/work/done/` are surveyed: which headings, node types and deviations occur, with counts.
- [x] Options (a) to (c) are tried on a sample of real items, with what each would reject, keep or rewrite.
- [x] A recommended model is in `## Outcome`, as a draft schema (section names, node types, cardinality, extension rule) and a round-trip rule; a hard-to-reverse choice has a `proposed` ADR; follow-up items are filed.

## Next

Survey the headings, node types and deviations in `docs/work/` and `docs/work/done/` (counts), then try options (a) to (c) on a sample.

## Outcome

Recommendation: build option (c), an open body with a typed core, and write the writer on a lossless section tree. Decision drafted as the proposed ADR [261031-gam0](../../adr/261031-gam0-work-item-body-model.md).

- **Survey** (123 items: 53 open, 70 closed; throwaway scripts, not kept). `## Goal` and `## Log` are in all 123; the order Goal, Next, Outcome, Log holds in 120; 112 have no other `##` heading and 11 have extras (`Options`, `Plan`, `To settle`, mostly early long items). 448 of 449 log lines are `- YYYY-MM-DD: text`. 106 Goals have a checklist (359 items, none nested, none with an id or commit marker); 102 end the Goal with "Done when:" and the list, 4 have prose after it, 17 have no checklist (12 of them closed). `Next` is prose in 36 of 39 cases. 28 items link to other records with Markdown links, 5 use wikilinks.
- **Options on the real items.** (a) closed core: 11 items fail. (b) declared extensions: same 11 are errors until declared, and it needs a declaration mechanism first. (c) open with a typed core: 123 of 123 valid, and the lossless split and join reproduces every file byte for byte.
- **Round trip rule.** Unknown content is never rewritten; a typed edit splices only the touched lines (ticking changes `[ ]` to `[x]` in place); nothing else is normalized.
- **Addressing.** Section name plus 1-based position, guarded by a section etag. No ids in the file: none exist today and they would add noise to hand edits.
- **A finding outside the model.** 24 of 70 closed items still have unticked "done when" boxes, so `close` does not tick them; this is [261031-zek2](../261031-zek2-close-checks-goal-checkboxes.md).
- **Not covered:** the spec's rules for `Next` (`Saved:` line), the ADR body, and a prototype writer in Go. The survey ran in Python over the files; the Go split is the follow-up.
- Follow-ups: [261031-j4bc](../261031-j4bc-body-section-tree-package.md) (the section tree package), and the open spikes 261031-6bzq (operations) and 261031-x50y (declaration and `check`), which build on this model.

## Log

- 2026-10-31: filed.
- 2026-10-31: claimed (claude).
- 2026-10-31: made ready and claimed to start the spike.
- 2026-10-31: closed done; drafted ADR 261031-gam0 and filed 261031-j4bc.
