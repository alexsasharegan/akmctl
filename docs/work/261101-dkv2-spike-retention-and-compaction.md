---
id: "261101-dkv2"
title: spike retention and compaction
status: raw
created: 2026-11-01
updated: 2026-11-01
kind: spike
summary: "Spike: a lifecycle for archiving and compacting completed work and keeping its knowledge, with a Git recovery path"
---

# spike retention and compaction

## Goal

Alexander's discovery spike on how AKM keeps the active working set useful as work piles up: how finished work is evaluated, archived, compacted and recovered. The full text, with the principles, the four questions, the open questions and the glossary, is [retention-and-compaction](../design/retention-and-compaction.md). Compaction today is only a glossary entry and the `keep` field ([261029-sg82](../adr/261029-sg82-work-item-closed-cycle-keep.md)), a design in progress.

The leading hypothesis to test: the project record becomes the archive index for its children (recovery commit, original paths, stable ids), so no second hierarchy is needed. Four questions stay apart: what to compact (retention), when (lifecycle boundary), how (operation) and how recovery is guaranteed (integrity).

Filed as one `raw` item. Whether to turn it into a parent with children is decided at refinement. The note proposes seven workstreams, in this dependency order:

- A. Lifecycle and archival semantics: `raw → ready → active → done → archived`, done versus archived.
- B. Context-local knowledge capture at completion (needs A).
- E. Record links and graph integrity: containment versus semantic links, links to archived records (needs A).
- C. Exit refinement and the retention policy (needs B and E).
- D. Project-as-archive-index architecture, nested archives (needs C).
- F. Compaction mechanics and Git recovery: resumable, rollback, recovery commit (needs D).
- G. Integration and conformance scenarios (needs F).

Done when:

- [ ] Refinement has decided: one spike with findings per workstream, or a parent with a child spike per workstream (use `split`, [261101-t8h8](261101-t8h8-make-work-decomposition-discoverable.md)).
- [ ] Each workstream has findings, a recommendation and its open questions answered or filed.
- [ ] A `proposed` ADR covers each hard-to-reverse choice (archive representation, the `archived` status, record-link vocabulary).
- [ ] The build items are filed and on the roadmap.

## Log

- 2026-11-01: filed.
