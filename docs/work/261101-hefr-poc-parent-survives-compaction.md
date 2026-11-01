---
id: "261101-hefr"
title: poc a parent survives compaction
status: active
created: 2026-11-01
updated: 2026-11-01
kind: spike
parent: 261101-wdbj
summary: "Spike: hand-build a compacted parent in a scratch repository and see whether it can still find and restore its children"
claimed_by: claude
claimed_at: 2026-11-01T23:13
---

# poc a parent survives compaction

## Goal

Part of [261101-wdbj](261101-wdbj-retention-and-compaction.md), and the first test of the hypothesis in [retention-and-compaction](../design/retention-and-compaction.md): a closed parent can serve as the archive index for its children, with no second record hierarchy.

Build the smallest throwaway version by hand in a scratch repository (not this one, and not a feature of the tool), then try to break it. The scratch repository is `akmctl init` plus a few closed children under one parent, one nested parent with its own children, and one link from retained work to a child.

Questions:

- What does the parent need to hold so that a child can be found and restored (recovery commit, original path, id)? What shape reads and parses well as a Markdown section?
- What does `akmctl check` say after the children are deleted: the parent's rollup, the child links, the closure rules (WK rules), `ls`, `show`? Which of those need to change?
- Can one child be restored from the recovery commit alone, without the rest?
- Does a nested archive stay recoverable when its parent is compacted in turn?
- What breaks the recovery: an uncommitted child, a rewritten branch, a moved file?

Done when:

- [ ] A scratch repository exists with the layout above and a script or command log that reproduces the compaction and the restore.
- [ ] Each question above is answered in `## Outcome` with the command output that shows it.
- [ ] The tool defects found are filed as items; the findings feed 261101-dkv2 (they are not a design decision, so no ADR here).

## Log

- 2026-11-01: filed.
- 2026-11-01 16:13-07:00: claimed (claude).
