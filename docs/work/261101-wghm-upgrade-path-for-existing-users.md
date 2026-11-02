---
id: "261101-wghm"
title: upgrade path for existing users
status: raw
created: 2026-11-01
updated: 2026-11-01
parent: 261101-vs5c
---

# upgrade path for existing users

## Goal

Maintainer question: how do changes in the agent rails reach existing users? The skills ship as templates in `internal/initrepo/templates/`, so changed wording reaches a repository only when it re-runs init or updates (unverified: whether an update command exists). Their own AGENTS.md is theirs and is not touched. Tool changes arrive with the binary, and some are refusals that can break existing scripts (claim with an empty Next, WK6; close now empties Next; CM4). Done when:

- [ ] A written answer covers each channel: skills, templates, AGENTS.md, binary behaviour, and the CHANGELOG notes, with what a user does to pick each up.
- [ ] A gap in that path (no way to refresh skills, no upgrade note for a new refusal) files its own item; a hard-to-reverse choice gets an ADR.

## Log

- 2026-11-01: filed.
