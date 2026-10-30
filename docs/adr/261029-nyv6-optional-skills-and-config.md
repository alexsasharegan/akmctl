---
id: "261029-nyv6"
title: Optional skills are switched by a manifest key and akmctl config, which may remove only akm- skills
summary: Optional skills are switched by a manifest key and akmctl config, which may remove only akm- skills
status: proposed
date: 2026-10-29
---

# Optional skills are switched by a manifest key and akmctl config, which may remove only akm- skills

## Context

The first optional skill is `akm-commit`, which tells an agent to commit after each logical step (status changes and `## Next` checkpoints). Not every maintainer wants an agent to commit. The choice must be made at `init` and be reversible later. The tool itself never stages or commits; only the skill does.

## Considered Options

- A skill file that is simply present or absent, with no recorded setting. Rejected: nothing in the repository says what was chosen, and `check` cannot see a mismatch.
- A manifest key with several modes. Rejected for the beta: each mode needs a definition and tests.
- A boolean manifest key (`git_commits`) in `docs/.layout`, with the skill following from it.

## Decision

`git_commits: true|false` in `docs/.layout` is the source of truth. `akmctl init` takes it as a parameter (a bare `init` on a terminal asks; flags and `--yes` skip the question; a non-terminal never asks). `akmctl config git-commits on|off` changes it later and installs or removes `.agents/skills/akm-commit/` to match. The skill's first step is to read the key and stop if it is off, so a stale copy does nothing. A `check` rule warns when the key and the directory disagree.

Removal is the tool's first deletion outside a migration, so it is narrow: it only removes a directory under `.agents/skills/` whose name starts with `akm-` and belongs to a known optional skill; it prints the paths first and asks (or needs `--yes`, or returns the plan under `dry_run`); if the installed files differ from the shipped copy it shows the difference and asks again; it never touches anything else. Writing refuses to overwrite a path that is not ours and reports the conflict.

## Consequences

- One new manifest key (a MINOR layout change), one new API operation `config`, one new check rule.
- The tool still never commits; the skill does, and only when the key is on.
- More optional skills can use the same mechanism later.
- The questionnaire is plain y/N prompts; a Bubble Tea wizard is deferred (work item 261029-95f2).
