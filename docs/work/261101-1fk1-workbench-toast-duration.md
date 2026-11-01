---
id: "261101-1fk1"
title: workbench toast duration
summary: Raise the default time a workbench toast stays up from 8 to 15 seconds, then judge whether it is too long
status: active
created: 2026-11-01
updated: 2026-11-01
kind: feature
claimed_by: claude
claimed_at: 2026-11-01T21:38
---

# workbench toast duration

## Goal

Alexander finds 8 seconds too short for a toast (`EPHEMERAL_MS` in `internal/server/web/js/toasts.js`). He is inclined to 10 to 15 seconds: set 15 and see if it is too much. Persistent toasts and an explicit `ms` are unchanged.

Done when:

- [x] the default is 15000 ms in `toasts.js`, with the CSS fallback in `app.css` and the unit tests updated
- [ ] Alexander has tried it and the Outcome records the value kept

## Next

Alexander tries 15 s in the workbench and says keep, lower (10 to 12) or go back.

## Log

- 2026-11-01: filed.
- 2026-11-01: default raised from 8 s to 15 s (`EPHEMERAL_MS`, the `--toast-ms` fallback, the tests) as a trial.
- 2026-11-01 14:38-07:00: claimed (claude).
