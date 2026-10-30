---
id: "261029-9aa2"
title: "Vault migration: manifest paths and migration steps"
status: done
created: 2026-10-29
updated: 2026-10-29
closed: 2026-10-29
---

# Vault migration: manifest paths and migration steps

## Goal

The vault keeps its AKM records in folders other than `docs/work` and `docs/adr`, and its tasks and ADRs need converting without renaming through the tool. Let the manifest name the folders, and add the migration steps that conversion needs.

Done when:

- [x] `docs/.layout` accepts `work_dir`, `adr_dir` and `roadmap`; every place that hard-coded the paths uses them; bad values are rejected with a clear error.
- [x] `version.repository` reports `work_dir` and `adr_dir`, and the workbench takes the paths from it.
- [x] Steps `title_from_h1`, `goal_scaffold`, `remove_field`, `copy_field` exist, and `rename_field` takes `dir`.
- [x] `mint_ids` takes `rename: false` and `seed`, and reports `renames`.
- [x] `migrate --migration <file>` and the API `migrations` parameter run a one-off migration.
- [x] Spec, API spec, schemas, fixtures and tests cover all of it.

## Outcome

Manifest keys work_dir, adr_dir and roadmap are read everywhere the paths were hard-coded (internal/record/repo.go accessors; tests record/paths_test.go, check TestCustomFolders, api/paths_test.go; fixture repo vault with cases ls-6, check-5, index-2, version-3, close-5); a bad value is rule S3 and refuses writes; version.repository reports work_dir and adr_dir and the workbench reads them (logic.test.mjs). The steps title_from_h1, goal_scaffold, remove_field, copy_field, dir on the field steps, mint_ids rename false and seed with renames, and the order with rename_field for numeric ADR ids are in internal/migrate (steps_test.go, api/migrate_oneoff_test.go, cases migrate-2 to 4). migrate --migration and the migrations parameter run a one-off migration under the same safeguards (cli_test.go, migrate_oneoff_test.go). Spec, guide, API spec, schemas and fixtures are updated in the vault; Ajv 8 validates 17 schemas and 44 cases; akmctl check reports 0 errors and 0 warnings; go vet and go test -race pass. Not run against the real vault.

## Log

- 2026-10-29: filed.
- 2026-10-29: claimed (claude).
- 2026-10-29: built the manifest keys `work_dir`, `adr_dir`, `roadmap` (`record.Repo` accessors, rule S3 for bad values, writes refused while S3 holds), `version.repository.work_dir`/`adr_dir`, the workbench taking the folders from `version`, fixture repo `vault` with cases, and the blank line after a created `## Outcome` (`SetSection`).
- 2026-10-29: built the steps `title_from_h1`, `goal_scaffold`, `remove_field`, `copy_field`, `dir` on the field steps, `mint_ids` `rename: false` and `seed`, `migrate --migration` and the `migrations` parameter; `plan` runs the steps on a scratch copy so dependent steps plan as they will run. Fixtures: repos `vault` and `legacy`, cases ls-6, check-5, index-2, version-3, close-5, migrate-2 to 4; schemas updated.
- 2026-10-29: spec (new rule S3, Settings, section 7), guide, API spec (4.15 and the CLI binding), schemas and fixtures updated in the vault; schema and case validation with Ajv 8 passes; rule ids and section references resolve. Branch pushed.
