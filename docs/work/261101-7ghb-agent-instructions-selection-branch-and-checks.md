---
id: "261101-7ghb"
title: agent instructions selection branch and checks
status: ready
created: 2026-11-01
updated: 2026-11-01
kind: chore
parent: 261101-vs5c
summary: AGENTS.md says what to do on a request to find work, how big a branch is, which checks to run and how to attribute a commit
---

# agent instructions selection branch and checks

## Goal

Four gaps in `AGENTS.md` from the first Codex trial. Each is one concrete sentence; Codex proposed wording for the first three in its post-mortem. Done when:

- [ ] A request to find or recommend work is answered with a recommendation and no claim, unless implementation is asked for; once a branch has a claim commit, it does not switch items.
- [ ] A branch and a PR carry one work item; a changed item means a new branch.
- [ ] The full-check line includes `-race` on Linux and macOS, and says what to do when the sandbox blocks the Go cache or Git (ask for the escalation; do not skip the check).
- [ ] The attribution convention for an agent's commit (a trailer after the `AKM-` ones) is decided and written, with the PR-description line if one is wanted.
- [ ] `akmctl check` and the skill mirrors still pass; `AGENTS.md` stays short enough to read in one pass.
- [x] `AGENTS.md` says a visible change gets a `CHANGELOG.md` entry under Unreleased in the same commit.

## Next

Four boxes remain (find-work requests, one item per branch, -race and sandbox note, attribution). Take them when 261101-jgdy and 261101-avg8 are done.

## Log

- 2026-11-01: filed.
- 2026-11-01 17:33-07:00: ticked "`AGENTS.md` says a visible change gets a `CHANGELOG.md` entry under Unreleased in the same commit.": AGENTS.md Rules, added with 261101-p2m8 after Codex's PR 1 changed visible output with no entry.
