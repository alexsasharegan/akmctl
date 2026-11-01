---
id: "261101-hefr"
title: poc a parent survives compaction
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
kind: spike
parent: 261101-wdbj
summary: "Spike: hand-build a compacted parent in a scratch repository and see whether it can still find and restore its children"
---

# poc a parent survives compaction

## Goal

Part of [261101-wdbj](../261101-wdbj-retention-and-compaction.md), and the first test of the hypothesis in [retention-and-compaction](../../design/retention-and-compaction.md): a closed parent can serve as the archive index for its children, with no second record hierarchy.

Build the smallest throwaway version by hand in a scratch repository (not this one, and not a feature of the tool), then try to break it. The scratch repository is `akmctl init` plus a few closed children under one parent, one nested parent with its own children, and one link from retained work to a child.

Questions:

- What does the parent need to hold so that a child can be found and restored (recovery commit, original path, id)? What shape reads and parses well as a Markdown section?
- What does `akmctl check` say after the children are deleted: the parent's rollup, the child links, the closure rules (WK rules), `ls`, `show`? Which of those need to change?
- Can one child be restored from the recovery commit alone, without the rest?
- Does a nested archive stay recoverable when its parent is compacted in turn?
- What breaks the recovery: an uncommitted child, a rewritten branch, a moved file?

Done when:

- [x] A scratch repository exists with the layout above and a script or command log that reproduces the compaction and the restore.
- [x] Each question above is answered in `## Outcome` with the command output that shows it.
- [x] The findings feed 261101-dkv2; no tool defect was found, so nothing was filed (they are not a design decision, so no ADR here).

## Next

Build the scratch repository in the session scratchpad: init, a parent with closed children, a nested parent, a link from retained work. Then compact by hand and run check, ls, show and a restore.

## Outcome

Yes, on the happy path: a closed parent can hold a `## Compaction Index` (schema, recovery commit, id and original path per child), its children can be deleted, and one child comes back with `git show <recovery>:<path>`. The scratch repository is `.tmp/poc-compaction-scratch/` (untracked); the sequence was: close four children and two parents with `akmctl close`, write an index into the parent, `git rm` the children, commit, check, restore.

- **Shape.** A plain `## Compaction Index` section (`Schema: 1`, `Recovery: <sha>`, a list of `id: path`) is accepted by `akmctl section` and reads fine. Each index names its own recovery commit, taken as HEAD before the removal.
- **Writing it.** The parent was already `done`, so the section write was refused until `--correction "<why>"` was passed (ADR 261101-dxcn); the correction then wrote a `corrected compaction index` line to the Log. Compaction needs a sanctioned write path to a closed parent (a `compact` operation, or an exemption like `keep`'s) so the Log does not say "corrected" for every archive.
- **`check` after removal.** Rules L1 (a link resolves to a file) fails for every link to a removed child: from a retained item (`audit-log` links to the rotation item) and from the parent's own Goal. Three errors in this run. L1 has to resolve such a link through an archive index, or the links must be rewritten when the child is compacted (compare 261101-1smq, which repairs links on a move).
- **`ls --parent`, `show`.** `ls --parent` just omits the removed children; there is no sign they ever existed. `show <id>` of a removed child says `not_found`. Both need to fall back to the index.
- **Restore.** One child restored from the index of its parent, alone, with no other record: yes. After restoring, its L1 error went away, so the links did resolve again.
- **Nesting.** The outer parent kept the inner parent as a retained child ("carries its own index"), so it listed only its direct compacted children. The inner index still restored independently. Compacting the inner parent itself away was not tried.
- **What breaks recovery.** After squashing to one commit, expiring the reflog and `git gc --prune=now`, both recovery commits were gone (`could not get object info`), and `akmctl check` still reported nothing about it. A recovery SHA is only as durable as the commit's reachability, so the index needs a check (`cat-file -e`) and a way to keep the commit reachable (a tag or ref per recovery point). Not run: an uncommitted child, a moved file.

Feeds 261101-dkv2: workstreams D (index shape holds, with the write path and L1 caveats), E (links to archived records) and F (recovery commit durability). No tool defect was found, so nothing was filed; the three tool changes above are design points for the spike.

## Log

- 2026-11-01: filed.
- 2026-11-01 16:13-07:00: claimed (claude).
