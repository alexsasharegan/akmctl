---
id: "261031-gtw6"
title: body schema and drawer redesign
summary: "Project: give a work item's body a schema, make every body write go through the tool, and redesign the drawer as a first-class editor"
status: raw
kind: chore
created: 2026-10-31
updated: 2026-10-31
---

# body schema and drawer redesign

## Goal

The project for one redesign with two halves, from Alexander's notes of 2026-10-31. It starts as spikes; the build items come from their outcomes.

1. **The body is a schema, not a text dump.** A work item already has a de facto body shape (`## Goal` with a "Done when" checklist, `## Next`, `## Log`, `## Outcome`). Decide how much of that is schematized and how much stays open for extension, and make the tool the only writer of it, so an agent never fights Markdown formatting or spends tokens getting it right.
2. **The drawer becomes a first-class editor.** Today it is a narrow panel off to the side. It should keep the iOS-sheet feel (the user never leaves the board) at a width that suits reading and editing (half, two-thirds or centered), and edit the body and the fields in place, with edits allowed only when they make sense for the item's status.

Children (spikes, in dependency order; none is built until its spike has an outcome):

- [261031-qbat](done/261031-qbat-spike-work-item-body-model.md): what a work item body contains, and the open-versus-closed question.
- [261031-6bzq](done/261031-6bzq-spike-body-writer-api.md): the Markdown writer and the body write API (needs the model).
- [261031-x50y](261031-x50y-spike-schema-declaration-and-migration.md): where the schema is declared, how `check` enforces it, how existing items migrate (needs the model).
- [261031-g9zk](done/261031-g9zk-spike-markdown-parser-dependency.md): which Markdown parser, if any, gives exact block spans for splice edits (needs the model; blocks the section tree package 261031-j4bc).
- [261031-18t8](261031-18t8-spike-drawer-layout.md): the drawer's width, layout and navigation (independent).
- [261031-z0n8](done/261031-z0n8-spike-field-edit-rules-by-status.md): which fields and body sections are editable in which status (needs the model; informs [261028-esyt](261028-esyt-ui-field-editing.md)).

[261028-esyt](261028-esyt-ui-field-editing.md) (edit fields in the drawer) is ready and first on the roadmap. Its field table and simple controls can proceed, but its "which fields are editable when" and its closed-items rule are decided by the last spike above. Done when:

- [ ] Every child is `done` or `dropped`.
- [ ] The decisions the spikes make are in `proposed` or `accepted` ADRs, and the build items they produce are filed and on the roadmap.
- [ ] 261028-esyt is updated to match the outcomes.

## Log

- 2026-10-31: filed with five spike children, from Alexander's notes on the drawer and a schematized body.
- 2026-10-31: added spike 261031-g9zk (parser dependency) after the body-model spike; 261031-j4bc waits for it.
- 2026-10-31: accepted ADR 261031-nvka (narrow implementations checked against oracles); parser ADR 261031-cpbg rewritten to that approach and still proposed; filed bug 261031-2j14 from the fuzz prototype.
