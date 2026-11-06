---
id: "261106-jkj9"
title: Commits and pull requests say assisted by
status: done
parent: 261101-vs5c
created: 2026-11-06
updated: 2026-11-06
closed: 2026-11-06
---

# Commits and pull requests say assisted by

## Goal

Agent attribution in this repository says "assisted by", not "co-authored" or "generated". Done when `AGENTS.md` names `Assisted-by: <agent and model>` for commits and `Assisted by <agent>` for pull requests, ADR 261101-z0mt is superseded by 261106-zvpm, the trial protocol row R13 matches, and `check` passes.

## Next


## Outcome

AGENTS.md, the trial protocol and a new ADR (261106-zvpm, superseding 261101-z0mt) now use Assisted-by for commits and Assisted by for pull requests. The user-level setting and the dotfiles instructions say the same. Not changed: older commits keep Co-Authored-By.

## Log

- 2026-11-06: filed.
- 2026-11-06 10:54-07:00: claimed (claude-assisted).
- 2026-11-06 10:54-07:00: Added AGENTS.md rule, ADR 261106-zvpm superseding z0mt, protocol R13 and a CHANGELOG entry.
