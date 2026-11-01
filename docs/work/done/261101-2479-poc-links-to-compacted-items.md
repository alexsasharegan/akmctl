---
id: "261101-2479"
title: poc links to compacted items
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
kind: spike
parent: 261101-wdbj
summary: "Spike: try the fragment, query, index-only and tombstone fixes for links to a compacted item"
---

# poc links to compacted items

## Goal

Part of [261101-wdbj](../261101-wdbj-retention-and-compaction.md), follow-up to [261101-hefr](261101-hefr-poc-parent-survives-compaction.md), which found that links to a compacted child break rule L1 and that `show` and `ls` lose it. Try the candidate fixes in the scratch repository and write the comparison into [261101-dkv2](../261101-dkv2-spike-retention-and-compaction.md) (workstream E).

Candidates:

- **Fragment rewrite (leading):** keep the link text, retarget to the parent that holds the index, put the child's id in the fragment: `[text](parent.md#<id>)`.
- **Query parameter:** the same with `?archived=<id>`.
- **Index only:** rewrite nothing; resolve the link through the index.
- **Tombstone:** leave a small stub at the original path.

Done when:

- [x] Each candidate that can be tried with the current tool is tried in the scratch repository, and `check`, `show` and `ls` output is recorded; the ones that cannot (index only, needs a code change) are marked as analysis.
- [x] Whether the fragment form can be reversed (restore the child and put the original link back) is shown.
- [x] The comparison is in dkv2's `## Goal` or `## Log` and the recommendation is stated.

## Next

Copy the scratch repository, apply each candidate, run check, show and ls, then write the comparison into dkv2.

## Outcome

The fragment rewrite is the recommended fix; the comparison is in 261101-dkv2.

## Log

- 2026-11-01: filed.
- 2026-11-01 16:32-07:00: claimed (claude).
