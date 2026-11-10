---
id: "261101-t8h8"
title: make work decomposition discoverable
status: done
created: 2026-11-01
updated: 2026-11-08
closed: 2026-11-08
kind: feature
summary: "Project for the outcome of spike 261101-baze: a newcomer or agent finds how to split an oversized item, and the tools make the split cheap"
---

# make work decomposition discoverable

## Goal

Project from spike 261101-baze (design: docs/design/decomposing-work.md). A maintainer knows to turn an oversized item into a parent with children; a newcomer or an agent has no surface that says so. The spike chose a section of `akm-track-work` ("Item too big", already in the tree, with its template mirror) over a new skill, because a new skill touches the skill mirror, six `init` fixtures in two copies, `make_fixtures.py` and two Go tests. The work is the tools that make the split cheap and visible:

- [261031-ph53](261031-ph53-new-parent-kind-flags-and-split-command.md): `new --parent --kind` and `split`, with priority inheritance, per-child commit messages and a parent reminder.
- [261031-72bd](261031-72bd-rollup-counts-for-parent-items.md): rollup counts, and `close` of a child reporting the parent and its open siblings.
- [261101-3erc](261101-3erc-block-commit-message-names-every-blocker.md): the block commit message names every blocker.

Done when:

- [x] Every child is `done` or `dropped`.
- [x] The `Item too big` section uses `split` and matches in both skill copies, and the guide (`docs/spec/akm-guide.md`, work items) points to it in one line.
- [x] A cold agent splits an oversized item using only the skill, and the Log records that it needed no guess the trial of 261101-baze listed.

## Next

## Outcome

Splitting is now cheap and visible. new takes --parent, --kind and --priority; akmctl split files the children and lists what the parent still needs; show and ls --tree give a parent progress counts; close names the parent and its open siblings; the block commit message names every blocker; the Item too big section uses split and the guide points to it. A cold subagent split 261030-456q using only the skill and hit none of the baze guesses except the judgment call on moving Goal boxes. Found on the way and filed: 261108-h7ds, 261108-kdz1, and 261107-1gha extended.

## Log

- 2026-11-01: filed.
- 2026-11-01 10:25-07:00: Filed from spike 261101-baze with ph53, 72bd and 3erc as children.
- 2026-11-08 00:01-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08 00:03-07:00: Cold-agent retest (a sonnet subagent in a scratch copy of main, guided only by AGENTS.md and the akm-track-work skill, no docs/design, on 261030-456q): split into a parent and four spike children, claimed and released, parent ready, check 0 errors. None of the guesses the 261101-baze trial listed came back except (5): moving a Goal box with no clean owner is still a judgment call, which the skill leaves to the agent by rule. New snags, unrelated to the split steps: a bracketed blocked_by list is stored mangled and check is silent (filed 261108-h7ds); set title= leaves the heading out of sync (filed 261108-kdz1); --changes refuses an object (added to 261107-1gha); split cannot give a child a title, so a child title comes from its slug. The agent also could not file those raw items itself (no git, a scratch copy).
- 2026-11-08 00:03-07:00: ticked "Every child is `done` or `dropped`.": 3erc (PR 28), 72bd (PR 29) and ph53 (PR 30) are done
- 2026-11-08 00:03-07:00: ticked "The `Item too big` section uses `split` and matches in both skill copies, and the guide (`docs/spec/akm-guide.md`, work items) points to it in one line.": step 2 of Item too big uses split in both skill copies (identical); the guide's Work items section points to it in one line
- 2026-11-08 00:03-07:00: ticked "A cold agent splits an oversized item using only the skill, and the Log records that it needed no guess the trial of 261101-baze listed.": retest in the Log: none of the baze guesses recurred except box moves (by rule); new snags filed as h7ds, kdz1, 1gha
