---
id: "261029-5799"
title: Use the akm-commit skill in this repository
status: done
created: 2026-10-29
updated: 2026-10-29
closed: 2026-10-29
kind: chore
priority: 1
summary: Turn on git_commits for this repository so agents use the akm-commit skill; AGENTS.md settles the push rule
---

# Use the akm-commit skill in this repository

## Goal

This repository uses its own optional skill: `git_commits` is true in the manifest and `.agents/skills/akm-commit/` is installed.

Done when:

- [x] `akmctl config git-commits on` has run here (manifest key and skill installed)
- [x] AGENTS.md says that pushing verified commits to origin is allowed here, since the skill's default is never to push
- [x] the checks pass

## Outcome

git_commits is true in docs/.layout and .agents/skills/akm-commit is installed. AGENTS.md now says verified commits are pushed to origin here (the skill's default is not to push) and that tags, releases and visibility stay with the maintainer.

## Log

- 2026-10-29: filed.
- 2026-10-29: claimed (claude).
