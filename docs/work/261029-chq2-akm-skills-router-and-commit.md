---
id: "261029-chq2"
title: "akm- skills: rename, router with references, optional akm-commit"
status: ready
created: 2026-10-29
updated: 2026-10-29
kind: feature
priority: 1
summary: Rename skills to akm-*, add the akm router skill with one reference per command, add the optional akm-commit skill.
---

# akm- skills: rename, router with references, optional akm-commit

## Goal

Per ADR 261029-tb3d (skill namespace) and ADR 261029-nyv6 (optional skills): rename the skills `init` writes and this repository's `.agents/skills/` to `akm-work-items` and `akm-record-decision`; add a router skill `akm` (short SKILL.md table plus `references/<command>.md` per operation); add the optional skill `akm-commit`, whose first step reads `git_commits` from `docs/.layout` and stops if it is not true.

Done when:

- [ ] Skill directories are renamed in `.agents/skills/` and in `internal/initrepo/templates/skills/`; the manifest extension `work-items` keeps its name.
- [ ] `akm` has a reference for every API operation, with when to use it, a real example, the files it changes and common errors; a test fails when an operation has none.
- [ ] `akm-commit` lives in `internal/initrepo/templates/optional-skills/akm-commit/`, checks the manifest first, commits per logical step (status changes and `## Next` checkpoints), stages by path, never pushes.
- [ ] The spec, guide, `AGENTS.md` and templates use the new names; the drift test covers the always-installed skills.

## Next

Not started.

## Log

- 2026-10-29: filed.
