---
id: "261101-7ghb"
title: agent instructions selection branch and checks
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
kind: chore
parent: 261101-vs5c
summary: AGENTS.md says what to do on a request to find work, how big a branch is, which checks to run and how to attribute a commit
decisions: [261101-z0mt]
---

# agent instructions selection branch and checks

## Goal

Four gaps in `AGENTS.md` from the first Codex trial. Each is one concrete sentence; Codex proposed wording for the first three in its post-mortem. Done when:

- [x] A request to find or recommend work is answered with a recommendation and no claim, unless implementation is asked for; once a branch has a claim commit, it does not switch items.
- [x] A branch and a PR carry one work item; a changed item means a new branch.
- [x] The full-check line includes `-race` on Linux and macOS, and says what to do when the sandbox blocks the Go cache or Git (ask for the escalation; do not skip the check).
- [x] The attribution convention for an agent's commit (a trailer after the `AKM-` ones) is decided and written, with the PR-description line if one is wanted.
- [x] `akmctl check` and the skill mirrors still pass; `AGENTS.md` stays short enough to read in one pass.
- [x] `AGENTS.md` says a visible change gets a `CHANGELOG.md` entry under Unreleased in the same commit.
- [x] `AGENTS.md` says to fill a new item's `## Goal` (reproduction, actual and expected, done-when) before committing it, since no rule can: `raw` may be empty (ADR 261101-8cd9).

## Next

## Outcome

AGENTS.md now says: a find or recommend request gets an answer and no claim; a branch and its PR carry one item and keep it after a claim commit; the full check includes `-race` and an escalation is asked for, never skipped; a new item's Goal is filled before it is committed; commits carry a `Co-Authored-By` trailer after the `AKM-` ones and a PR ends with a line naming the agent; work lands through a PR unless the maintainer asks for a local merge (ADR 261101-z0mt). `akmctl check` passes and AGENTS.md is 48 lines. No CHANGELOG entry: no command, API, spec or rule changed. Not verified: whether another model follows the sentences, which is 261101-0vh7. Decision: 261101-z0mt.

## Log

- 2026-11-01: filed.
- 2026-11-01 17:33-07:00: ticked "`AGENTS.md` says a visible change gets a `CHANGELOG.md` entry under Unreleased in the same commit.": AGENTS.md Rules, added with 261101-p2m8 after Codex's PR 1 changed visible output with no entry.
- 2026-11-01 17:53-07:00: claimed (claude).
- 2026-11-01 17:54-07:00: ticked "A request to find or recommend work is answered with a recommendation and no claim, unless implementation is asked for; once a branch has a claim commit, it does not switch items.": AGENTS.md Rules and Commands, ADR 261101-z0mt.
- 2026-11-01 17:54-07:00: ticked "A branch and a PR carry one work item; a changed item means a new branch.": AGENTS.md Rules and Commands, ADR 261101-z0mt.
- 2026-11-01 17:54-07:00: ticked "The full-check line includes `-race` on Linux and macOS, and says what to do when the sandbox blocks the Go cache or Git (ask for the escalation; do not skip the check).": AGENTS.md Rules and Commands, ADR 261101-z0mt.
- 2026-11-01 17:54-07:00: ticked "The attribution convention for an agent's commit (a trailer after the `AKM-` ones) is decided and written, with the PR-description line if one is wanted.": AGENTS.md Rules and Commands, ADR 261101-z0mt.
- 2026-11-01 17:54-07:00: ticked "`akmctl check` and the skill mirrors still pass; `AGENTS.md` stays short enough to read in one pass.": AGENTS.md Rules and Commands, ADR 261101-z0mt.
- 2026-11-01 17:54-07:00: ticked "`AGENTS.md` says to fill a new item's `## Goal` (reproduction, actual and expected, done-when) before committing it, since no rule can: `raw` may be empty (ADR 261101-8cd9).": AGENTS.md Rules and Commands, ADR 261101-z0mt.
