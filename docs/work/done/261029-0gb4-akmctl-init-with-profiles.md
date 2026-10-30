---
id: "261029-0gb4"
title: Add akmctl init with the profiles code and notes
summary: akmctl init creates the AKM layout in a new repository from a profile (code or notes); it replaces the template repository.
status: done
kind: feature
priority: 2
created: 2026-10-29
updated: 2026-10-29
closed: 2026-10-29
---

# Add akmctl init with the profiles code and notes

## Goal

Decided 2026-10-29 ([261029-j0d0](../../adr/261029-j0d0-one-tool-spec-and-fixtures-in-this-repository.md)): the template repository is dropped and `akmctl init` creates the layout in a new repository. A profile picks what it writes: `code` (a software repository: `AGENTS.md` with commands, `docs/adr/`, `docs/work/`, `docs/ROADMAP.md`, the skills, the manifest) and `notes` (a notes repository or vault: the records under a notes folder, the manifest keys `work_dir`, `adr_dir` and `roadmap`). The input is the file list in [akm-template-lessons](../../design/akm-template-lessons.md). Like every write, `init` refuses to overwrite and reports what it wrote.

Done when:

- [ ] `akmctl init --profile code` and `--profile notes` write a repository on which `akmctl check` reports 0 errors.
- [ ] Neither profile overwrites an existing file; a second run changes nothing.
- [ ] `init` is an API operation with a schema and fixtures, and it is in the CLI and HTTP bindings.
- [ ] Tests cover both profiles; the README and the standard name `init` (see [261029-dsga](261029-dsga-spec-template-section-describes-init.md)).

## Next

Coordinator review: the README and the standard (261029-dsga) must name init; the spec text for the init operation (request profile and dir; response profile, dir, layout, written, skipped; error refused/not_a_directory) and the harness field after_paths are still unwritten. Work item 261029-ap3y is absorbed (init writes .gitignore with docs/.akm.lock).

## Outcome

`akmctl init --profile code|notes [dir]` writes the layout without overwriting (`internal/initrepo`). `akmctl check` reports 0 errors and 0 warnings on both profiles (`TestProfilesCheckClean`); a second run changes nothing (`TestNoOverwriteAndIdempotent`). API operation with `schemas/init.json`, fixtures init-1 to init-5 in both copies, CLI and HTTP bindings (`TestInit`, `TestInitRoute`), and section 4.17 of the API spec. Checked by hand in a scratch directory for both profiles. The Ajv validator was not run (not installed).

## Log

- 2026-10-29: filed.
- 2026-10-30: claimed (init-agent).
- 2026-10-30: init operation (internal/initrepo, api, CLI, HTTP), schema, fixtures in both copies, tests; check reports 0 errors and 0 warnings on both profiles.
- 2026-10-29: released (init-agent).
