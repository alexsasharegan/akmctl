---
id: "261101-vs5c"
title: agent rails hardening
status: done
created: 2026-11-01
updated: 2026-11-06
closed: 2026-11-06
kind: feature
summary: "Project from the first Codex trial: make the repository's rails steer any model to the expected output"
---

# agent rails hardening

## Goal

The first Codex run (PR 1) shows where the instructions held and where they did not; the evidence and method are in [the trial write-up](../../design/first-codex-trial.md). The code in that PR is secondary. The work is to move each deviation to the layer that can hold it: the tool where a rule can be computed, one concrete sentence in `AGENTS.md` where it needs judgment, then measure on more than one model.

- [261101-p2m8](261101-p2m8-claim-refuses-item-without-next.md): `claim` stops leaving an item in a state `check` rejects.
- [261101-jgdy](261101-jgdy-log-and-section-guard-against-shell-split-text.md): `log` and `section` notice shell-split text.
- [261101-avg8](261101-avg8-check-flags-stale-next-and-empty-goal.md): `check` flags a stale `## Next` on close and an empty `## Goal`.
- [261101-7ghb](261101-7ghb-agent-instructions-selection-branch-and-checks.md): `AGENTS.md` says how to select, scope a branch, run checks and attribute commits.
- [261101-0vh7](261101-0vh7-cross-model-rails-trial.md): rerun the task on two models and score both.
- [261101-qr8x](261101-qr8x-rails-name-and-binary.md): claim with a unique name; use the branch build; check the version before filing.
- [261101-ft79](261101-ft79-step-trailers-not-git-trailers.md): `commits` rule CM5 and the whole-message commit.
- [261101-w82x](261101-w82x-close-leaves-blank-line-after-next.md): `close` leaves a Prettier-clean record.
- [261101-y3ch](261101-y3ch-log-at-error-names-the-expected-form.md): a pattern error names the expected form.
- [261101-wghm](261101-wghm-upgrade-path-for-existing-users.md): the Upgrading section of the guide.
- [261106-jkj9](261106-jkj9-assisted-by-trailer.md): "assisted by" attribution.

Done when:

- [x] Every child is `done` or `dropped`.

## Next

## Outcome

The rails hardening is done for Claude: every rails child is done (p2m8, jgdy, avg8, 7ghb, 0vh7, qr8x, ft79, w82x, y3ch, wghm, jkj9). Deferred on 2026-11-06 by the maintainer, whose focus is Claude: the rerun on Codex and other harnesses and the PR row, now 261106-xrpr; f9fb (a gentler commits check) is its own item, not a rails gap.

## Log

- 2026-11-01: filed.
- 2026-11-01 20:05-07:00: triage: promoted ft79, qr8x, w82x, y3ch to ready with a Next; round 2 order set; 4gk3, kngg, yqec flagged as off-theme for this project; wghm and the rerun of the PR row stay open
- 2026-11-06 13:19-07:00: claimed (claude-sonnet-5-5).
- 2026-11-06 13:19-07:00: removed "The rerun in 261101-0vh7 shows no deviation from the trial's table that a tool or an `AGENTS.md` sentence could have prevented.": Narrowed on 2026-11-06: the rerun covered Claude, and the maintainer deferred Codex and other harnesses and the PR row; they move to 261106-xrpr.
- 2026-11-06 13:19-07:00: ticked "Every child is `done` or `dropped`.": f9fb detached on 2026-11-06 (own item, not a rails gap); p2m8, jgdy, avg8, 7ghb, 0vh7, qr8x, ft79, w82x, y3ch, wghm and jkj9 are done
- 2026-11-06 13:19-07:00: Closing: f9fb detached, the Codex rerun and the PR row deferred to 261106-xrpr (the maintainer's focus is Claude).
