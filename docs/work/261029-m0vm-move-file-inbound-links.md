---
id: "261029-m0vm"
title: "move_file: report the inbound links it leaves broken"
status: active
created: 2026-10-29
updated: 2026-11-08
kind: feature
priority: 3
summary: move_file does not repair links; two sweep links broke in the vault run and were fixed by hand.
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-08T15:18
---

# move_file: report the inbound links it leaves broken

## Goal

When move_file moves a note, scan for links to the old path and report each (file and line), as mint_ids does for renames.

Done when:

- [ ] move_file lists inbound links per moved file in its notes.
- [ ] A test with a path-qualified wikilink.

## Next

Read the `move_file` step in `internal/migrate` and the scanner in `internal/mdlinks`; decide how a moved file's inbound links are reported (file and line); add the report and a test with a path-qualified wikilink.

## Log

- 2026-10-29: filed.
- 2026-10-31: 261031-wjs3 added internal/mdlinks (link scanner and rewrite) and the close-side repair; the migrate-side report can reuse it.
- 2026-11-08 08:18-07:00: claimed (claude-sonnet-5-5).
