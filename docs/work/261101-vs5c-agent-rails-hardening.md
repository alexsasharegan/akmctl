---
id: "261101-vs5c"
title: agent rails hardening
status: active
created: 2026-11-01
updated: 2026-11-06
kind: feature
summary: "Project from the first Codex trial: make the repository's rails steer any model to the expected output"
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-06T20:19
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

Round 1 is done. Round 2 from the rerun: take 261101-qr8x and 261101-ft79 (the two gaps that decide whether Codex and Claude behave the same), then 261101-w82x and 261101-y3ch. Then repeat the 0vh7 protocol on the PR row, which neither model reached. 261101-wghm (upgrade path) needs an answer before the next release. 261101-4gk3, kngg and yqec are not rails gaps: consider moving them to their own parent.

## Log

- 2026-11-01: filed.
- 2026-11-01 20:05-07:00: triage: promoted ft79, qr8x, w82x, y3ch to ready with a Next; round 2 order set; 4gk3, kngg, yqec flagged as off-theme for this project; wghm and the rerun of the PR row stay open
- 2026-11-06 13:19-07:00: claimed (claude-sonnet-5-5).
