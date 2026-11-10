---
id: "261106-60xv"
title: Git is the move log
summary: The order, author and time of moves come from the step commits in git, read by a new op; a journal file is rejected for now and kept as a spike
status: proposed
date: 2026-11-06
---

# Git is the move log

## Context

The workbench should replay the moves made while the maintainer was away (work item 261101-w0mz). It polls the work listing and compares two snapshots, so it sees the net change only: ids come back in id order, a path through several statuses shows as one move, a burst of more than ten is summarized and lost, and "Outside" never says who. Nothing the API returns says when a status changed or who changed it. Replay, and later undo (261029-shz5), churn detection (261031-st7s) and counts (261031-72bd), need an ordered record of moves.

The step commits already are that record when `git_commits` is on: each has a subject `<id>(<step>): <detail>`, the `AKM-Step` and `AKM-Item` trailers, an author, a time and an `Assisted-by` trailer, and `commits` already reads and checks them.

## Considered Options

- **Git history.** A read-only `history` op derived from the step commits. No new storage and no layout change.
- **A journal file in the repository.** Every write appends an event. Exact and includes uncommitted moves, but it adds a file to the layout before the 1.0.0 freeze, conflicts when branches append in parallel, misses hand and Obsidian edits, grows (retention and compaction, 261031-dkv2 and 261101-wdbj), and holds derived data next to the records.
- **Per-item `## Log` lines.** Stays in the records, but has no global order, minute resolution and hand-edit drift.
- **A journal in the `akmctl ui` process.** Cheap, but blind while the server is down, which is the case replay is for, and unavailable to the CLI.
- **Client only.** Label by `claimed_by` and keep the snapshot diff. No true order.

## Decision

Git is the move log. A read-only `history` operation lists the AKM step commits in order, with item, step, author, time and the trailers, and says plainly when it has nothing to show (not a git work tree, `git_commits` off, no step commits). It belongs to the git-aware part of the tool, beside `commits`; the core keeps working without git, and the workbench falls back to the snapshot diff and says why. Moves not yet committed, and hand edits, are seen by the snapshot diff only. No journal file is added now; 261106-wx2a keeps the question open. Decided by Alexander on 2026-11-06, in conversation.

## Consequences

- Replay gets a true order, an author and a path across commits, and survives a long absence by remembering the last commit seen.
- Order is commit order, not edit order, and several status changes inside one commit are one event.
- A repository that does not commit per step has little history; the workbench says so instead of guessing.
- An addition to the API (a new operation, fixtures, a spec section), not a layout change.
- Revisit the journal if history proves too coarse, or if undo needs uncommitted changes.

## Confirmation

The `history` op and its fixtures (261106-sdnr); the workbench replay (261101-w0mz) reads it.

## Review

Proposed 2026-11-06 after a discussion of five options with the maintainer; accepted by the maintainer on merge.

## Revert

Supersede this record with a decision for a journal and keep `history` as a reader of git for older repositories.
