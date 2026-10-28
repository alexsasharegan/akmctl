---
id: "261028-8dd4"
title: Build the mint_ids migration step
summary: A migrate step that adds ids to existing records so a repository without ids can adopt AKM without minting each id by hand.
status: ready
kind: feature
priority: 3
created: 2026-10-28
updated: 2026-10-28
---

# Build the mint_ids migration step

## Goal

Adopting AKM on a repository whose notes have no ids means minting an id for each existing record by hand: the id recipe prints a new id for today's date, so a record that should read as old needs the date part edited, ids minted in one batch collide unless the batch is tracked (the recipe checks only files that exist), and each id must be written twice (the filename and the `id` field). Add a migration primitive `mint_ids` so the tool does this.

Done when:

- [ ] A `mint_ids` step in `internal/migrate` that, for each Markdown file in a named directory with no `id`, mints an id, renames the file to `<id>-<slug>.md`, and writes the `id` field.
- [ ] The id date part comes from a chosen field of the record (for example `created`), falling back to today, so old records read as old.
- [ ] Ids are unique across the batch and the repository, and the step is idempotent: a second run changes nothing.
- [ ] It never deletes a file, and references that named the old filename are reported, not rewritten.
- [ ] Fixtures in `testdata/` cover a batch with a date field, one without, and a collision.

## Log

- 2026-10-28: filed.
