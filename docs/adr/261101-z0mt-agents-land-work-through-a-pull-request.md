---
id: "261101-z0mt"
title: Agents land work through a pull request
summary: An agent pushes a branch and opens a pull request; it merges to main only when the maintainer asks; commits carry a Co-Authored-By trailer
status: superseded
superseded_by: "261106-zvpm"
date: 2026-11-01
---

# agents land work through a pull request

## Context

AGENTS.md told agents to push verified commits to origin. In the first Codex trial that meant a branch and a PR only because the maintainer asked for one; the maintainer does not trust an unreviewed agent run, and a shared rule has to hold for every model. The maintainer may still ask for a local merge, as in the work after the trial.

## Considered Options

- A pull request for every agent, with the maintainer free to override for a session.
- A pull request for every agent except a trusted one.
- Direct pushes to `main` for every agent, as before.

## Decision

Every agent works on a branch, pushes the verified branch and opens a pull request. It does not push to `main` and does not merge its own pull request unless the maintainer asks for a local merge. Commits keep the configured Git author and end with `Co-Authored-By: <agent and model> <noreply address>` after the `AKM-` trailers; a pull request description ends with one line naming the agent. A branch and its pull request carry one work item. Decided by Alexander on 2026-11-01.

## Consequences

- Review happens on the pull request for any model, so the rule does not depend on which agent is trusted.
- A session where the maintainer asks to merge locally lands work faster; the override is explicit in the conversation, not in the file.
- `commits --range` and the CM rules still check the branch's commits before the push.

## Confirmation

AGENTS.md Rules; reviewed on each pull request.

## Review

Accepted by Alexander, 2026-11-01, in conversation.

## Revert

Restore the previous Rules line in AGENTS.md (push verified commits to origin) and drop the attribution sentence.
