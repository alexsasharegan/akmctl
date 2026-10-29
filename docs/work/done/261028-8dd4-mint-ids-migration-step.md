---
id: "261028-8dd4"
title: Build the mint_ids migration step
summary: A migrate step that adds ids to existing records so a repository without ids can adopt AKM without minting each id by hand.
status: done
kind: feature
priority: 3
created: 2026-10-28
updated: 2026-10-29
---

# Build the mint_ids migration step

## Goal

Adopting AKM on a repository whose notes have no ids means minting an id for each existing record by hand: the id recipe prints a new id for today's date, so a record that should read as old needs the date part edited, ids minted in one batch collide unless the batch is tracked (the recipe checks only files that exist), and each id must be written twice (the filename and the `id` field). Add a migration primitive `mint_ids` so the tool does this.

Done when:

- [x] A `mint_ids` step in `internal/migrate` that, for each Markdown file in a named directory with no `id`, mints an id, renames the file to `<id>-<slug>.md`, and writes the `id` field.
- [x] The id date part comes from a chosen field of the record (for example `created`), falling back to today, so old records read as old.
- [x] Ids are unique across the batch and the repository, and the step is idempotent: a second run changes nothing.
- [x] It never deletes a file, and references that named the old filename are reported, not rewritten.
- [x] Fixtures in `testdata/` cover a batch with a date field, one without, and a collision.

## Outcome
mint_ids is built: a migrate step that, for each Markdown file in a directory with no id, mints an id (date from a chosen field, today otherwise; unique in the batch and across docs/adr, docs/work and docs/work/done), renames it to <id>-<slug>.md, writes id first and quoted, and reports references to the old name without rewriting them. It runs in report, plan, apply_safe and handoff, never commits, and a second run changes nothing. No migration embeds it yet. The ids shown by plan are examples: apply draws again. CI passed on Linux, macOS and Windows.

## Log

- 2026-10-28: filed.
- 2026-10-29: claimed (claude).
- 2026-10-29: built `mint_ids` in `internal/migrate/mintids.go` (`dir`, `date_field`, `recursive`; id minted from the date field, unique across the batch and docs/adr, docs/work, docs/work/done; rename to `<id>-<slug>.md`; `id` written first and quoted; references reported, never rewritten). Added `migrate.Env` (random source and clock) and step `notes` in the API. Tests: unit and API tests over `testdata/migrate/mint_ids`.
