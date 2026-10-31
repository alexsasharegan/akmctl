---
id: "261029-57jq"
title: The closed table flickers on reload
summary: Reloading while the closed table is shown redraws it visibly.
status: done
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-10-31
closed: 2026-10-31
---

# The closed table flickers on reload

## Goal

Reloading while the closed table is shown redraws it visibly.

Done when:

- [x] The table redraws in place without a flash.

## Outcome

The closed table redraws in place: the list stays in state across reloads (store.js no longer clears it), tables.js draws it at once and redraws only when the id, status, updated or title of some row changed, and shows the "Loading closed work…" text only the first time. Checked in the workbench: pressing r on #/tables/closed remounted the table once, the loading text never appeared, and the rows stayed. Not compared against the old build side by side, and not tried with a real network delay. Decision: none.

## Log

- 2026-10-29: filed.
- 2026-10-31: claimed (claude).
