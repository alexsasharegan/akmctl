---
id: "261029-w3d9"
title: git_commits manifest key, akmctl config, init questionnaire
status: active
created: 2026-10-29
updated: 2026-10-29
kind: feature
priority: 1
summary: Add git_commits to the manifest, the config operation and command, the optional-skill install and removal, and the init questionnaire.
claimed_by: config-agent
claimed_at: 2026-10-30T01:36
---

# git_commits manifest key, akmctl config, init questionnaire

## Goal

Per ADR 261029-nyv6: a boolean `git_commits` key in `docs/.layout`; `init` takes it as a parameter and, on a terminal, asks plain y/N questions (capital letter is the default) for it and for `git init` when the directory is not a repository; `--yes`/flags skip the questions and a non-terminal never asks; `akmctl config` lists settings and `akmctl config git-commits on|off` changes it, installing or removing `.agents/skills/akm-commit/` under the removal limits in the ADR. Windows is best effort (ADR 261029-eegx).

Done when:

- [ ] `config` is an API operation (list/get/set, `dry_run` returns the plan) with schema, fixtures in both copies, CLI and HTTP bindings, and a section in the API spec.
- [ ] Removal only touches `.agents/skills/akm-*` directories of known optional skills, prints paths first, needs confirmation or `--yes`, shows a diff when files were edited, and reports conflicts without overwriting.
- [ ] A new check rule warns when `git_commits` and the skill directory disagree; the spec names the key.
- [ ] `init` accepts `git_commits` and `git_init`; the questionnaire is in the CLI layer only and never blocks without a terminal; tests cover both paths.

## Next

Not started.

## Log

- 2026-10-29: filed.
- 2026-10-29: claimed (config-agent).
