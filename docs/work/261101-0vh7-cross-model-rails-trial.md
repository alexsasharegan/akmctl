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
- 2026-11-01 18:06-07:00: Claude Code run on f6089e8 (subagent, local, no push or PR): stage 1 held R1 and R2, recommended 261029-wpak and did not claim; stage 2 held R3-R8, R10-R13 (claim, 2 step commits, CM1-CM4 ok, Co-Authored-By present, CHANGELOG and glossary, Next emptied by close, full check with -race); R9 n/a, R14 not run. Gaps: claimed as plain 'claude' (CL2 risk, AGENTS.md does not say to use a unique --as); format:check fails on main at the wgf6 item, which no run caught before merge (fixed here); the run did not file it since it was not an akmctl failure
- 2026-11-01 18:50-07:00: Codex run (GPT-6, exact maintainer prompt, clean clone at f6089e8): found work with 'akmctl ls --status ready,active,blocked' and 'board' after reading the skills, so retrieval worked. Held R1, R3, R4 (unique name codex-gpt-6), R5, R8-R10 (escalated for the Go cache), R11, R13, R15. Deviated: stage rule (claimed and implemented when asked to find; the prompt also said PR flow, so ambiguous); R7 (AKM trailers in separate paragraphs, not git trailers, commits still passed: 261101-ft79); first checkpoint hand-wrote a step commit with literal backslash-n and the branch was abandoned; used a stale installed akmctl and filed a false raw item kpjx (261101-qr8x); gh token invalid, branch pushed and PR not opened (R14 not reached)
