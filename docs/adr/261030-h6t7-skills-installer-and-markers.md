---
id: "261030-h6t7"
title: A skills command installs, checks and updates the akm skills, and each skill folder carries a marker file
summary: akmctl skills install|status|update writes the embedded akm skills with a marker per skill folder, telling outdated copies from hand-edited ones.
status: proposed
date: 2026-10-30
---

# A skills command installs, checks and updates the akm skills, and each skill folder carries a marker file

## Context

`akmctl init` is the only way the `akm` skills reach a repository, and it is for new projects. A repository adopted through a one-off migration (the maintainer's Obsidian vault, 2026-10-29) has the records and none of the skills, and the only way to get them is to copy files out of a source checkout, which drifts silently from the binary. Skills are text an agent reads, so a stale copy is a stale instruction.

ADRs 261029-nyv6 and 261029-tb3d already fix the rules for anything the tool writes under `.agents/skills/`: only the router `akm` and skills named `akm-<name>`; print the paths first and ask (or `--yes`, or a plan under `dry_run`); refuse to overwrite a path that is not ours; if an installed file differs from the shipped copy, show the difference and ask again. This record decides the command, how a copy's provenance is recorded, and how `init` fits.

## Considered Options

- Keep copying by hand. Rejected: the cause of the drift.
- Extend `akmctl config` (`config skills install`). Rejected: `config` switches manifest keys; an installer with a status view and a diff is a command family of its own.
- Record one `skills_version` in `docs/.layout`. Rejected: the manifest is committed and shared, one number cannot tell an outdated copy from one edited by hand, and a skill folder can be copied or deleted without touching it.
- Record nothing and compare content with the shipped copy. Rejected: it can say "differs" but not why, so `update` could not tell a safe replacement from overwriting a local edit.
- A command family, `skills install|status|update`, with a marker file in each skill folder.
- Check on every command and update unasked. Rejected: the tool writes nothing it was not asked to.

## Decision

**The managed set** is what `init` writes: the router `akm`, `akm-work-items` and `akm-record-decision`, and `akm-commit` only when `git_commits` is true. Nothing outside that set is read for update or written, whatever it is called.

**The command** is `akmctl skills status|install|update`, one API operation `skills` with a `mode` of `status` (the default, read-only), `install` or `update`, plus `dry_run`, `confirm` and `confirm_edited` as `config` has them. A write mode without `confirm` returns the plan (paths, and a diff for each changed file) and writes nothing. The API stays 1.0: the operation is additive.

**The marker** is a file `.akm-skill.json` in each managed skill folder: `{"tool": "<version>", "files": {"<path in the folder>": "<sha256 of the bytes the tool wrote>"}}`. It is written with the skill and is committed with it. It records what the tool shipped, not what is on disk.

**Status of a skill**, from the marker, the shipped files and the files on disk:

- `missing`: no folder.
- `current`: every file equals the shipped one.
- `outdated`: the files equal what the marker says was written, and the shipped files differ. Replacing is safe.
- `edited`: a file differs from what the marker says was written. A local change; `update` shows the diff and needs `confirm_edited`.
- `unmarked`: a folder of a managed name with no marker (a copy by hand, or written by an older `init`). If every file equals the shipped one, `install` or `update` writes the marker; if not, it is treated as `edited`.
- `conflict`: a path where a file or directory of the skill should be is not ours (for example a file where a folder belongs). Reported and never overwritten.

**`install`** writes `missing` skills and marks `unmarked` ones that match; it leaves `current`, `outdated` and `edited` ones alone and says so. **`update`** replaces `outdated` skills, and `edited` ones only with `confirm_edited` after showing the diff. Neither deletes anything: a file a newer version no longer ships is reported as `stale` and left. Removal stays with `config git-commits off`, under its existing narrow rule.

**One code path.** `init` and `config git-commits on` call the same installer, so a skill they write has its marker, and a repository made by an older `init` is adopted the first time `skills install` or `update` runs.

**No new `check` rule.** The status view is the report. A rule would change the layout version for a case `status` already covers; revisit with the next layout change.

**No network.** The skills come from the binary; `update` means "to the version of this binary".

## Consequences

- One new operation and command, one new file per skill folder, the `akm-skill` marker in the specification and the guide, a schema and fixtures, a glossary entry (skill marker), and a changelog line. The marker is a JSON file inside a folder the tool already owns by name.
- The vault gets the skills with `akmctl skills install`, and later binaries can update them without losing a local edit.
- `init` output lists one more file per skill; the `init` fixtures change.
- A repository's own `akm-` skill that is not in the managed set is never touched, and `status` does not list it as ours.
- A skill folder written before this record is `unmarked` until first run; there is no migration step.

## Confirmation

A test fails if `init` writes a managed skill without a marker, if `update` replaces an `edited` skill without `confirm_edited`, or if any run writes outside the managed set. The status of each state above has a test.

## Review

Proposed 2026-10-30 after the maintainer chose the command shape, the marker file and the shared code path.

## Revert

Remove the `skills` operation and command and stop writing the markers; the skill folders and what `init` and `config` write are unchanged. Delete the `.akm-skill.json` files.
