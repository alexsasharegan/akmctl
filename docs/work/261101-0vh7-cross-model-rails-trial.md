---
id: "261101-0vh7"
title: cross model rails trial
status: active
created: 2026-11-01
updated: 2026-11-01
kind: spike
parent: 261101-vs5c
summary: Rerun the trial's task on a clean checkout with a reworded prompt on two or more models, scored against one table
claimed_by: claude-0vh7
claimed_at: 2026-11-02T01:02
---

# cross model rails trial

## Goal

The trial of 2026-11-01 had one model, one loose prompt and no baseline. Done when:

- [ ] A fixed prompt (recommend an item and wait; then implement it on a branch) and a scoring table (the rows of [the trial write-up](../design/first-codex-trial.md)) are written.
- [ ] The same task is run on a clean checkout on Codex and on Claude Code, after the other children are done, with each run's log read as the trial's was.
- [ ] The Log records each deviation per model and which layer should have held it; each one that is a tool or wording gap files its own item.

## Next

Claude Code run: start a fresh session in a clean worktree of the commit in the Log, send the two stage prompts from docs/design/rails-trial-protocol.md, score the log. Codex run: the maintainer runs the same prompts; then read its rollout log and score it. Record each deviated row per model here and file gap items.

## Log

- 2026-11-01: filed.
- 2026-11-01 18:02-07:00: claimed (claude-0vh7).
- 2026-11-01 18:02-07:00: protocol and scoring table written in docs/design/rails-trial-protocol.md (15 rows, two-stage prompt); base commit for the runs is f6089e8
