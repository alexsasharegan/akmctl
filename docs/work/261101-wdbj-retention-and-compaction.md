---
id: "261101-wdbj"
title: retention and compaction
status: raw
created: 2026-11-01
updated: 2026-11-01
kind: feature
summary: "Umbrella for retention and compaction: keep the active working set small, keep the knowledge, keep the way back"
---

# retention and compaction

## Goal

The owner of the retention and compaction work: a feature that keeps the active working set small as closed work piles up, without losing what it taught or the way back to the originals. Compaction is so far only a glossary entry and the `keep` field, and ADR 261029-6cp5 left it "a separate decision". This item groups the decisions and the build.

Design source: [retention-and-compaction](../design/retention-and-compaction.md). Its hypothesis is that a closed parent item can serve as the archive index for its children, with no second record hierarchy. Here that parent is the umbrella item (an "umbrella is a parent", glossary), not a new kind of record; the design note's word "project" means that parent.

Children (none is built until its spike has an outcome):

- [261101-dkv2](261101-dkv2-spike-retention-and-compaction.md): the discovery spike, workstreams A to G. Refinement may split it into one spike per workstream, as children of this item.
- [261101-hefr](261101-hefr-poc-parent-survives-compaction.md): a hand-built proof of concept in a scratch repository that a compacted parent can still find and restore its children. It runs ahead of the spike, to feed it.

What already exists and bears on it:

- Decided: [261029-sg82](../adr/261029-sg82-work-item-closed-cycle-keep.md) (`closed`, `cycle`, `keep`; a calendar of periods was suggested, not decided), [261101-dxcn](../adr/261101-dxcn-closed-records-read-only-except-corrections.md) (closed items are read-only, `keep` stays editable), [261029-6cp5](../adr/261029-6cp5-roadmap-now-next-later.md) (finished work leaves the roadmap; compaction not covered), [261031-gam0](../adr/261031-gam0-work-item-body-model.md) and [261101-tr8e](../adr/261101-tr8e-body-writes-through-the-section-tree.md) (the body model and the one writer that compaction would summarize and rewrite).
- Specified: `keep` and WK13 to WK17 in the guide and spec; glossary entries Compaction and Keep.
- Related design: [decomposing-work](../design/decomposing-work.md) and [261101-t8h8](261101-t8h8-make-work-decomposition-discoverable.md) (parents, closure and rollups, which the archive index would extend); [decision-promotion](../ideas/decision-promotion.md) (the same promotion of knowledge out of a work item that workstream B wants at completion); [ralph-loop-consultation](../design/ralph-loop-consultation.md) (unattended loops are heavy users of closed work).
- Mechanics that an archive operation would need: [261029-xtqj](261029-xtqj-move-no-rollback.md) (multi-write moves with no rollback), [261029-shz5](261029-shz5-undo-one-step.md) (undo covers one step), [261029-m0vm](261029-m0vm-move-file-inbound-links.md) and [261101-1smq](261101-1smq-close-repairs-links-to-moved-items.md) (links broken or repaired when a file moves), [261029-47m0](261029-47m0-plan-simulates-move-file.md) (plans that simulate moves).

Done when:

- [ ] Every child is `done` or `dropped`.
- [ ] The decisions the spike makes are in `proposed` or `accepted` ADRs, and the build items it produces are filed and on the roadmap.
- [ ] The glossary, the guide and the spec describe the archived state, the archive index and the compaction operation, and `Compaction` no longer says "not yet built".

## Log

- 2026-11-01: filed.
