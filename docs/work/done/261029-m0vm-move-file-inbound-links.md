---
id: "261029-m0vm"
title: "move_file: report the inbound links it leaves broken"
status: done
created: 2026-10-29
updated: 2026-11-08
closed: 2026-11-08
kind: feature
priority: 3
summary: move_file does not repair links; two sweep links broke in the vault run and were fixed by hand.
---

# move_file: report the inbound links it leaves broken

## Goal

When move_file moves a note, scan for links to the old path and report each (file and line), as mint_ids does for renames.

Done when:

- [x] move_file lists inbound links per moved file in its notes.
- [x] A test with a path-qualified wikilink.

## Next

## Outcome

move_file reports, in the step's notes, each relative Markdown link and path-qualified wikilink in another file that points at the moved file (file and line), in plan and apply. Nothing is rewritten. mdlinks gained Wikilinks; move_closed is filed separately.

## Log

- 2026-10-29: filed.
- 2026-10-31: 261031-wjs3 added internal/mdlinks (link scanner and rewrite) and the close-side repair; the migrate-side report can reuse it.
- 2026-11-08 08:18-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08 08:20-07:00: new internal/migrate/inbound.go lists inbound Markdown links and path-qualified wikilinks as notes of move_file (plan and apply); internal/mdlinks gains Wikilinks and a shared line walker; test TestMoveFileReportsInboundLinks; spec and changelog updated. Move_closed does not report yet.
