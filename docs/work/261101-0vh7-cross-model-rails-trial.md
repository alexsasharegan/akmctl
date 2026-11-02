---
id: "261101-0vh7"
title: cross model rails trial
status: blocked
created: 2026-11-01
updated: 2026-11-01
kind: spike
parent: 261101-vs5c
summary: Rerun the trial's task on a clean checkout with a reworded prompt on two or more models, scored against one table
blocked_by: [261101-p2m8, 261101-jgdy, 261101-avg8, 261101-7ghb]
---

# cross model rails trial

## Goal

The trial of 2026-11-01 had one model, one loose prompt and no baseline. Done when:

- [ ] A fixed prompt (recommend an item and wait; then implement it on a branch) and a scoring table (the rows of [the trial write-up](../design/first-codex-trial.md)) are written.
- [ ] The same task is run on a clean checkout on Codex and on Claude Code, after the other children are done, with each run's log read as the trial's was.
- [ ] The Log records each deviation per model and which layer should have held it; each one that is a tool or wording gap files its own item.

## Next

Wait for 261101-p2m8, 261101-jgdy, 261101-avg8 and 261101-7ghb. Then write the fixed prompt and the scoring table from docs/design/first-codex-trial.md, and run it on Codex and Claude Code.

## Log

- 2026-11-01: filed.
