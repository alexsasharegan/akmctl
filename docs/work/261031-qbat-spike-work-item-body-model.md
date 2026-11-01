---
id: "261031-qbat"
title: spike work item body model
summary: "Spike: what a work item body contains, how much is schema versus open text, and how extension works"
status: active
kind: spike
parent: "261031-gtw6"
created: 2026-10-31
updated: 2026-10-31
claimed_by: claude
claimed_at: 2026-11-01T03:39
---

# spike work item body model

## Goal

Work item bodies already follow a convention: `## Goal` (prose plus a "Done when" checklist), `## Next`, `## Log` (dated lines), `## Outcome`. Treat that as a schema and find out where it should be strict and where open. Questions:

- Which sections exist per status and kind (a spike has an `## Outcome`; a bug may want repro steps; a feature a checklist). Which are required, optional or repeatable.
- The node types inside a section: prose, checklist item (with a done state and an optional commit), dated log entry, link list, table. Whether each is a typed node or opaque text.
- Open versus closed: the options to compare are (a) a fixed set of sections plus a free-form `notes` area, (b) fixed core sections plus extension sections declared by name, kind and node type, (c) any `##` heading allowed with only the core ones typed. Include what an Obsidian user typing by hand can still do, and what survives a round trip unchanged.
- The model in the data: a body as an ordered list of typed sections (the in-memory form the writer and the drawer both use), with a stable id per section and per checklist item so edits address them without matching text.
- How it relates to frontmatter: which facts belong in a field (priority, project) and which in the body (goal text); whether a "done when" item is data other rules read (see [261031-zek2](261031-zek2-close-checks-goal-checkboxes.md)).

Done when:

- [ ] The current bodies in `docs/work/` and `docs/work/done/` are surveyed: which headings, node types and deviations occur, with counts.
- [ ] Options (a) to (c) are tried on a sample of real items, with what each would reject, keep or rewrite.
- [ ] A recommended model is in `## Outcome`, as a draft schema (section names, node types, cardinality, extension rule) and a round-trip rule; a hard-to-reverse choice has a `proposed` ADR; follow-up items are filed.

## Next

Survey the headings, node types and deviations in `docs/work/` and `docs/work/done/` (counts), then try options (a) to (c) on a sample.

## Log

- 2026-10-31: filed.
- 2026-10-31: claimed (claude).
- 2026-10-31: made ready and claimed to start the spike.
