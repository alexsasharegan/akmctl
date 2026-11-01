---
id: "261101-dkv2"
title: spike retention and compaction
status: raw
created: 2026-11-01
updated: 2026-11-01
kind: spike
summary: "Spike: a lifecycle for archiving and compacting completed work and keeping its knowledge, with a Git recovery path"
parent: 261101-wdbj
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

## Links to compacted items

Findings of [261101-2479](done/261101-2479-poc-links-to-compacted-items.md), which tried each fix in a scratch repository (`.tmp/` copies of the POC). The problem: after a child is compacted, a link to it breaks rule L1, `show` returns `not_found` and `ls --parent` omits it. Bare ids in frontmatter (`parent`, `blocked_by`) cannot be rewritten by any link syntax, so every option still needs id lookup to fall back to the index.

| Option                                   | L1 (`check`)                                        | Reversible                                | Obsidian / GitHub click                                      | `show`, `ls`                  | Cost                                                                                                                               |
| ---------------------------------------- | --------------------------------------------------- | ----------------------------------------- | ------------------------------------------------------------ | ----------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| Fragment: `[text](parent.md#<id>)`       | passes (tried: 0 errors)                            | yes, byte-identical (tried)               | opens the parent; fragment ignored (analysis, untried)       | no change; needs index lookup | rewrite of links in closed files, by the tool                                                                                      |
| Query: `[text](parent.md?archived=<id>)` | passes (tried: 0 errors)                            | yes in principle (not tried)              | untried; a query on a local file may not resolve in Obsidian | no change; needs index lookup | same rewrite; no benefit over the fragment                                                                                         |
| Index only                               | needs a code change (`links()` only stats the file) | n/a: nothing rewritten                    | dead link when clicked (analysis)                            | needs index lookup            | no edits to closed files; every check reads every index                                                                            |
| Tombstone at the original path           | passes (tried: 0 errors)                            | trivially: replace the stub with the file | works unchanged (analysis)                                   | works, shows the stub (tried) | files stay (353 bytes each); needs a body that satisfies WK7 or a new `archived` status (tried `status: archived`: refused by WK3) |

Notes:

- **Query:** L1 drops `?…` and `#…` before the file test ([mdlinks.go:121](../../internal/mdlinks/mdlinks.go)), so it passes. I said earlier that Obsidian fails on a query; that was not tested. No benefit over the fragment either way, so drop it.
- **Fragment:** the id in the fragment is what makes the rewrite reversible: restoring a child finds links to the pointer with that id and puts the original path back, with the path from the index. Link text is untouched.
- **Tombstone:** as `status: done` with `## Goal` and `## Outcome` it passes every rule, and `show` and `ls` keep working with no code change. A real `archived` status would need a new value in WK3 and in `ls`.

Recommendation: the fragment rewrite for body Markdown links, plus lookup through the index for bare ids, `show`, `ls` and `check`. Keep the tombstone as the fallback if the lookup proves costly. Drop the query form.

## Log

- 2026-11-01: filed.
- 2026-11-01 16:15-07:00: POC 261101-hefr: a Compaction Index section on a closed parent restores a child from a Git recovery commit. Open for D, E, F: closed parents need a sanctioned write path (not --correction), L1 and show/ls must resolve through the index, and a recovery SHA must be checked and kept reachable (a squash plus gc loses it silently).
- 2026-11-01 16:32-07:00: Compared four fixes for links to compacted items (261101-2479); recommended the fragment rewrite plus index lookup; query dropped; tombstone kept as fallback.
