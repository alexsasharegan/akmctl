---
id: "261030-8d8t"
title: skills install and update command
status: active
created: 2026-10-30
updated: 2026-10-30
kind: feature
priority: 3
decisions: [261030-h6t7]
summary: A command installs, checks and updates the akm skills packaged in the binary, in a repository that was not made by init.
claimed_by: claude
claimed_at: 2026-10-31T04:43
---

# skills install and update command

## Goal

A command (shape open: `akmctl skills install|status|update`, or an extension of `config`) writes the skills embedded in the binary (`akm`, `akm-work-items`, `akm-record-decision`, and optionally `akm-commit`) into `.agents/skills/` of an existing repository, records which version it wrote, and reports drift between the installed copy and the one the binary ships.

Why: `init` is the only way skills arrive today, and it is for new projects. A repository adopted through a one-off migration (the maintainer's Obsidian vault, 2026-10-29) has the AKM records but none of the skills, so the only way to get them is to copy files out of a source checkout, which drifts silently from the binary.

Constraints already decided (ADR 261029-nyv6, ADR 261029-tb3d): the tool owns only skills named `akm-<name>` plus the `akm` router; it prints paths first and asks (or takes `--yes`, or returns a plan under `dry_run`); it refuses to overwrite a path that is not ours and reports the conflict; if an installed file differs from the shipped copy it shows the difference and asks again.

Done when:

- [x] An ADR (`proposed`, 261030-h6t7) records the command shape, how the installed version is recorded (a manifest key, or a marker in each skill directory), and the new API operation, since it adds to the 1.0 API.
- [ ] The command installs the missing skills from the binary into a repository with a layout and no skills, and leaves a repository that already has them unchanged.
- [ ] It refuses to overwrite a skill directory that is not ours, and `--dry-run` and `--yes` work as for `config`.
- [ ] A status view, or a `check` rule, reports a skill whose installed version differs from the shipped one, and an optional skill whose key and directory disagree.
- [ ] `update` shows the diff for a changed skill before replacing it.
- [ ] The API schema, fixtures, specification, guide, glossary and changelog are updated; a test covers a repository made by an older tool version.
- [ ] Run in the maintainer's vault, it installs the three skills and `akmctl check` reports no new warnings.

## Next

The maintainer reviews ADR 261030-h6t7 and accepts it. Then build `internal/skills` (status, plan, install, update, markers), route `init` and `config git-commits` through it, add the `skills` operation (params, schema, fixtures), the CLI command and `GET /skills`, then the spec, guide, glossary and changelog.

## Log

- 2026-10-30: filed.
- 2026-10-30: filed from the vault migration review; the vault has the records but not the skills.
- 2026-10-30: claimed (claude).
- 2026-10-30: ADR 261030-h6t7 drafted (`proposed`) after the maintainer chose `skills install|status|update`, a marker file per skill folder and one code path with `init`.
