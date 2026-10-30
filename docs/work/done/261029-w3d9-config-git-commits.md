---
id: "261029-w3d9"
title: git_commits manifest key, akmctl config, init questionnaire
status: done
created: 2026-10-29
updated: 2026-10-29
kind: feature
priority: 1
summary: Add git_commits to the manifest, the config operation and command, the optional-skill install and removal, and the init questionnaire.
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

Built and tested; not closed. Left for the maintainer: review the draft Outcome, replace the placeholder `internal/initrepo/templates/optional-skills/akm-commit/SKILL.md` with the real text (if the skill gains files, update the `write` lists of fixtures config-3 and config-4 and init-6 through `make_fixtures.py`), then close.

## Outcome

`git_commits` is a manifest key; API `config` lists, reads and sets it with a plan under `dry_run`, conflict reporting, and double confirmation before removing an edited skill (spec 4.18, schema `config.json`, fixtures config-1 to config-10 in both copies). `init` takes `git_commits` and `git_init` and forces the key off outside a work tree (`git_commits_forced_off`). Rule S4 warns when the key, the skill directory and the work tree disagree. CLI `config` and the y/N `init` questionnaire (terminal only) and HTTP `/config`. Checked by hand in a scratch directory: init outside git forces off; `config on` after `git init` installs the skill; `off --yes` removes it; `check` clean throughout. `npm run schemas`: 19 schemas, 61 cases, 0 failures.

## Log

- 2026-10-29: filed.
- 2026-10-29: claimed (config-agent).
- 2026-10-29: manifest key git_commits, operation config (plan, conflict, confirm, edits), init params git_commits and git_init, rule S4, CLI questionnaire and `config` command, HTTP `/config`; git_commits stays off outside a git work tree. Tests and fixtures config-1..10, init-6..7 in both copies.
- 2026-10-29: released (config-agent).
