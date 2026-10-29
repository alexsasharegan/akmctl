---
id: "261029-9aa2"
title: "Vault migration: manifest paths and migration steps"
status: active
created: 2026-10-29
updated: 2026-10-29
claimed_by: claude
claimed_at: 2026-10-29T03:06
---

# Vault migration: manifest paths and migration steps

## Goal

The vault keeps its AKM records in folders other than `docs/work` and `docs/adr`, and its tasks and ADRs need converting without renaming through the tool. Let the manifest name the folders, and add the migration steps that conversion needs.

Done when:

- [ ] `docs/.layout` accepts `work_dir`, `adr_dir` and `roadmap`; every place that hard-coded the paths uses them; bad values are rejected with a clear error.
- [ ] `version.repository` reports `work_dir` and `adr_dir`, and the workbench takes the paths from it.
- [ ] Steps `title_from_h1`, `goal_scaffold`, `remove_field`, `copy_field` exist, and `rename_field` takes `dir`.
- [ ] `mint_ids` takes `rename: false` and `seed`, and reports `renames`.
- [ ] `migrate --migration <file>` and the API `migrations` parameter run a one-off migration.
- [ ] Spec, API spec, schemas, fixtures and tests cover all of it.

## Next

Code and fixtures are done. Update the spec, guide and API spec in the vault, run the checks, close this item, push the branch.

## Log

- 2026-10-29: filed.
- 2026-10-29: claimed (claude).
- 2026-10-29: built the manifest keys `work_dir`, `adr_dir`, `roadmap` (`record.Repo` accessors, rule S3 for bad values, writes refused while S3 holds), `version.repository.work_dir`/`adr_dir`, the workbench taking the folders from `version`, fixture repo `vault` with cases, and the blank line after a created `## Outcome` (`SetSection`).
- 2026-10-29: built the steps `title_from_h1`, `goal_scaffold`, `remove_field`, `copy_field`, `dir` on the field steps, `mint_ids` `rename: false` and `seed`, `migrate --migration` and the `migrations` parameter; `plan` runs the steps on a scratch copy so dependent steps plan as they will run. Fixtures: repos `vault` and `legacy`, cases ls-6, check-5, index-2, version-3, close-5, migrate-2 to 4; schemas updated.
