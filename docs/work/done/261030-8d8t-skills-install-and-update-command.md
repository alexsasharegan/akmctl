---
id: "261030-8d8t"
title: skills install and update command
status: done
created: 2026-10-30
updated: 2026-10-30
closed: 2026-10-30
kind: feature
priority: 3
decisions: [261030-h6t7]
summary: A command installs, checks and updates the akm skills packaged in the binary, in a repository that was not made by init.
---

# skills install and update command

## Goal

A command (shape open: `akmctl skills install|status|update`, or an extension of `config`) writes the skills embedded in the binary (`akm`, `akm-work-items`, `akm-record-decision`, and optionally `akm-commit`) into `.agents/skills/` of an existing repository, records which version it wrote, and reports drift between the installed copy and the one the binary ships.

Why: `init` is the only way skills arrive today, and it is for new projects. A repository adopted through a one-off migration (the maintainer's Obsidian vault, 2026-10-29) has the AKM records but none of the skills, so the only way to get them is to copy files out of a source checkout, which drifts silently from the binary.

Constraints already decided (ADR 261029-nyv6, ADR 261029-tb3d): the tool owns only skills named `akm-<name>` plus the `akm` router; it prints paths first and asks (or takes `--yes`, or returns a plan under `dry_run`); it refuses to overwrite a path that is not ours and reports the conflict; if an installed file differs from the shipped copy it shows the difference and asks again.

Done when:

- [x] An ADR (`proposed`, 261030-h6t7) records the command shape, how the installed version is recorded (a manifest key, or a marker in each skill directory), and the new API operation, since it adds to the 1.0 API.
- [x] The command installs the missing skills from the binary into a repository with a layout and no skills, and leaves a repository that already has them unchanged.
- [x] It refuses to overwrite a skill directory that is not ours, and `--dry-run` and `--yes` work as for `config`.
- [x] A status view, or a `check` rule, reports a skill whose installed version differs from the shipped one, and an optional skill whose key and directory disagree.
- [x] `update` shows the diff for a changed skill before replacing it.
- [x] The API schema, fixtures, specification, guide, glossary and changelog are updated; a test covers a repository made by an older tool version.
- [x] Run in the maintainer's vault, it installs the three skills and `akmctl check` reports no new warnings.

## Next

Only the last check remains: run `akmctl skills install` in the maintainer's vault and confirm `akmctl check` reports no new warnings. That is the maintainer's repository, so it is left to them.

## Outcome

akmctl skills status|install|update ships (ADR 261030-h6t7), with a marker per skill folder and one code path shared with init and config git-commits. Verified in the maintainer's Obsidian vault: install wrote 28 files, status was clean, and akmctl check was unchanged (0 errors, 128 warnings). That vault's task 261030-wf9z (outside this repository, in the vault under 9 pkm/tasks/done/) then retired its task-lifecycle skill for akm-track-work. Not exercised on a real folder: the akm-work-items rename (the vault never had it), covered by API tests. One follow-up fix came from the vault run: the text output now says whether a plan was applied.

## Log

- 2026-10-30: filed.
- 2026-10-30: filed from the vault migration review; the vault has the records but not the skills.
- 2026-10-30: claimed (claude).
- 2026-10-30: ADR 261030-h6t7 drafted (`proposed`) after the maintainer chose `skills install|status|update`, a marker file per skill folder and one code path with `init`.
- 2026-10-30: released (claude).
- 2026-10-30: ADR 261030-h6t7 accepted by the maintainer. Sequenced after 261030-p7r1 (the rename of akm-work-items and the naming grammar).
- 2026-10-30: claimed (claude).
- 2026-10-30: built `akmctl skills status|install|update`: `internal/initrepo/skills.go` (managed set, markers, renames), `internal/api/ops_skills.go`, the CLI and `GET`/`POST /skills`, schema `skills.json`, spec 4.20, the `akm` reference `skills.md`, glossary and changelog. `init` and `config git-commits` write markers through the same code. Deviation from ADR 261030-h6t7: an unmarked folder under a retired name is treated as edited (needs `confirm_edited`), because the files the 0.3.x init wrote are not embedded to compare with.
- 2026-10-31: first run in the maintainer's vault: 28 files installed, `clean: true`, `check` unchanged. Fix from its report: the text output now ends with "applied" or "nothing written (a plan only)".
