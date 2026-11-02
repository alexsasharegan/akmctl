---
id: "261101-0vh7"
title: cross model rails trial
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
kind: spike
parent: 261101-vs5c
summary: Rerun the trial's task on a clean checkout with a reworded prompt on two or more models, scored against one table
---

# cross model rails trial

## Goal

The trial of 2026-11-01 had one model, one loose prompt and no baseline. Done when:

- [x] A fixed prompt (recommend an item and wait; then implement it on a branch) and a scoring table (the rows of [the trial write-up](../../design/first-codex-trial.md)) are written.
- [x] The same task is run on a clean checkout on Codex and on Claude Code, after the other children are done, with each run's log read as the trial's was.
- [x] The Log records each deviation per model and which layer should have held it; each one that is a tool or wording gap files its own item.

## Next


## Outcome

Both models held most rows. Codex found work with ls and board unaided and claimed under a unique name; Claude claimed as plain claude. Gaps: step trailers outside git's trailer block (261101-ft79), a unique --as and which binary to use (261101-qr8x). Caveats: the Claude run used a two-stage prompt, not the maintainer's, and was not repeated; the PR row was not reached on either. Scores are in the Log and docs/design/rails-trial-protocol.md.

## Log

- 2026-11-01: filed.
- 2026-11-01 18:02-07:00: claimed (claude-0vh7).
- 2026-11-01 18:02-07:00: protocol and scoring table written in docs/design/rails-trial-protocol.md (15 rows, two-stage prompt); base commit for the runs is f6089e8
- 2026-11-01 18:06-07:00: Claude Code run on f6089e8 (subagent, local, no push or PR): stage 1 held R1 and R2, recommended 261029-wpak and did not claim; stage 2 held R3-R8, R10-R13 (claim, 2 step commits, CM1-CM4 ok, Co-Authored-By present, CHANGELOG and glossary, Next emptied by close, full check with -race); R9 n/a, R14 not run. Gaps: claimed as plain 'claude' (CL2 risk, AGENTS.md does not say to use a unique --as); format:check fails on main at the wgf6 item, which no run caught before merge (fixed here); the run did not file it since it was not an akmctl failure
- 2026-11-01 18:50-07:00: Codex run (GPT-6, exact maintainer prompt, clean clone at f6089e8): found work with 'akmctl ls --status ready,active,blocked' and 'board' after reading the skills, so retrieval worked. Held R1, R3, R4 (unique name codex-gpt-6), R5, R8-R10 (escalated for the Go cache), R11, R13, R15. Deviated: stage rule (claimed and implemented when asked to find; the prompt also said PR flow, so ambiguous); R7 (AKM trailers in separate paragraphs, not git trailers, commits still passed: 261101-ft79); first checkpoint hand-wrote a step commit with literal backslash-n and the branch was abandoned; used a stale installed akmctl and filed a false raw item kpjx (261101-qr8x); gh token invalid, branch pushed and PR not opened (R14 not reached)
- 2026-11-01 18:57-07:00: ticked "A fixed prompt (recommend an item and wait; then implement it on a branch) and a scoring table (the rows of [the trial write-up](../../design/first-codex-trial.md)) are written.": docs/design/rails-trial-protocol.md
- 2026-11-01 18:57-07:00: ticked "The same task is run on a clean checkout on Codex and on Claude Code, after the other children are done, with each run's log read as the trial's was.": Codex run on the maintainer's exact prompt; Claude run on the two-stage prompt, repeat skipped for usage limits
- 2026-11-01 18:57-07:00: ticked "The Log records each deviation per model and which layer should have held it; each one that is a tool or wording gap files its own item.": Log entries per model; items ft79, qr8x filed under vs5c
