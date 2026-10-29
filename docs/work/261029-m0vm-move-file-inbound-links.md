---
id: "261029-m0vm"
title: "move_file: report the inbound links it leaves broken"
status: ready
created: 2026-10-29
updated: 2026-10-29
kind: feature
priority: 3
summary: move_file does not repair links; two sweep links broke in the vault run and were fixed by hand.
---

# move_file: report the inbound links it leaves broken

## Goal

When move_file moves a note, scan for links to the old path and report each (file and line), as mint_ids does for renames.

Done when:

- [ ] move_file lists inbound links per moved file in its notes.
- [ ] A test with a path-qualified wikilink.

## Log

- 2026-10-29: filed.
