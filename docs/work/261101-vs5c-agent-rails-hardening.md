---
id: "261101-vs5c"
title: agent rails hardening
status: ready
created: 2026-11-01
updated: 2026-11-01
kind: feature
summary: "Umbrella from the first Codex trial: make the repository's rails steer any model to the expected output"
---

# agent rails hardening

## Goal

The first Codex run (PR 1) shows where the instructions held and where they did not; the evidence and method are in [the trial write-up](../design/first-codex-trial.md). The code in that PR is secondary. The work is to move each deviation to the layer that can hold it: the tool where a rule can be computed, one concrete sentence in `AGENTS.md` where it needs judgment, then measure on more than one model.

- [261101-p2m8](done/261101-p2m8-claim-refuses-item-without-next.md): `claim` stops leaving an item in a state `check` rejects.
- [261101-jgdy](done/261101-jgdy-log-and-section-guard-against-shell-split-text.md): `log` and `section` notice shell-split text.
- [261101-avg8](done/261101-avg8-check-flags-stale-next-and-empty-goal.md): `check` flags a stale `## Next` on close and an empty `## Goal`.
- [261101-7ghb](done/261101-7ghb-agent-instructions-selection-branch-and-checks.md): `AGENTS.md` says how to select, scope a branch, run checks and attribute commits.
- [261101-0vh7](done/261101-0vh7-cross-model-rails-trial.md): rerun the task on two models and score both.

Done when:

- [ ] Every child is `done` or `dropped`.
- [ ] The rerun in 261101-0vh7 shows no deviation from the trial's table that a tool or an `AGENTS.md` sentence could have prevented.

## Next

Take 261101-p2m8 (claim refuses an item with no Next), then 261101-jgdy, 261101-avg8 and 261101-7ghb; 261101-0vh7 reruns the trial last. Read docs/design/first-codex-trial.md first.

## Log

- 2026-11-01: filed.
