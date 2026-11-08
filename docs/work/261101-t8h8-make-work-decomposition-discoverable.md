---
id: "261101-t8h8"
title: make work decomposition discoverable
status: ready
created: 2026-11-01
updated: 2026-11-01
kind: feature
summary: "Project for the outcome of spike 261101-baze: a newcomer or agent finds how to split an oversized item, and the tools make the split cheap"
---

# make work decomposition discoverable

## Goal

Project from spike 261101-baze (design: docs/design/decomposing-work.md). A maintainer knows to turn an oversized item into a parent with children; a newcomer or an agent has no surface that says so. The spike chose a section of `akm-track-work` ("Item too big", already in the tree, with its template mirror) over a new skill, because a new skill touches the skill mirror, six `init` fixtures in two copies, `make_fixtures.py` and two Go tests. The work is the tools that make the split cheap and visible:

- [261031-ph53](261031-ph53-new-parent-kind-flags-and-split-command.md): `new --parent --kind` and `split`, with priority inheritance, per-child commit messages and a parent reminder.
- [261031-72bd](done/261031-72bd-rollup-counts-for-parent-items.md): rollup counts, and `close` of a child reporting the parent and its open siblings.
- [261101-3erc](done/261101-3erc-block-commit-message-names-every-blocker.md): the block commit message names every blocker.

Done when:

- [ ] Every child is `done` or `dropped`.
- [ ] The `Item too big` section uses `split` and matches in both skill copies, and the guide (`docs/spec/akm-guide.md`, work items) points to it in one line.
- [ ] A cold agent splits an oversized item using only the skill, and the Log records that it needed no guess the trial of 261101-baze listed.

## Next

Take 261101-3erc (small bug), then 261031-ph53 (split), then 261031-72bd; finish with the guide pointer and the cold-agent retest.

## Log

- 2026-11-01: filed.
- 2026-11-01 10:25-07:00: Filed from spike 261101-baze with ph53, 72bd and 3erc as children.
