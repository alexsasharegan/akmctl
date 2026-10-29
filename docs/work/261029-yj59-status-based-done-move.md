---
id: "261029-yj59"
title: Move closed records to done/ by status in one step
status: ready
created: 2026-10-29
updated: 2026-10-29
kind: feature
priority: 3
summary: The vault migration generated 76 move_file steps and a 106 KB .migrate-log.
---

# Move closed records to done/ by status in one step

## Goal

Add a migrate step that moves every record whose status is done or dropped into the done folder, and keep .migrate-log to counts plus changed paths.

Done when:

- [ ] A step with a dir parameter moves closed records.
- [ ] .migrate-log stays small for 100 files.
- [ ] Tests on the vault fixture.

## Log

- 2026-10-29: filed.
