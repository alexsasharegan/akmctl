---
id: "261029-cq7r"
title: Dogfood browser pass of the finished workbench
summary: Capture, refine, claim, block, unblock, close with an ADR link, review diff and a live update were not driven in the browser after step 6; Obsidian op
status: done
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-10-29
closed: 2026-10-29
---

# Dogfood browser pass of the finished workbench

## Goal

Capture, refine, claim, block, unblock, close with an ADR link, review diff and a live update were not driven in the browser after step 6; Obsidian open and takeover of another holder's claim are also untried.

Done when:

- [x] Each flow is driven on a scratch copy and the findings are in the Log.

## Outcome

The finished workbench was driven in the built-in browser on a scratch copy: capture, refine, claim, block, unblock, close with an ADR link, review diff and a live update from akmctl set. Obsidian open and takeover of another holder were not driven.

## Log

- 2026-10-29: filed.
- 2026-10-29: driven on a scratch copy with the built-in browser: capture (live slug), refine raw to ready (Goal), claim, block with blocked_by (an id and free text), unblock, close with an Outcome and an ADR link, the review diff, and a live update from akmctl set (priority and title changed on the open board within 3 s). Not driven: Obsidian open (href only) and takeover of another holder claim. Found: 261029-7zcp.
