---
id: "261028-k0ek"
title: Give the workbench live refresh
summary: The workbench page shows a snapshot; it does not update when records change on disk.
status: done
kind: feature
priority: 0
created: 2026-10-28
updated: 2026-10-29
closed: 2026-10-29
---

# Give the workbench live refresh

## Goal

The workbench page loads the records once. If a record changes on disk (an agent edits it, `akmctl` writes it), the page stays stale until it is reloaded. The planned UI port will decide between server-sent events and polling; this item records the gap in the current page.

Done when:

- [ ] The page updates within a few seconds of a record changing on disk, without a reload.
- [ ] A change made by the page itself does not flicker or lose the user's place.

## Outcome

Already built: internal/server/web/js/live.js polls the work listing every 2 seconds, diffs it, labels each change you or outside, animates it, and shows an offline indicator; an open drawer refreshes in place and the animator waits while you type or a modal is open. Read in code, not driven in the browser. Server-sent events and a conditional poll are separate, later items.

## Log

- 2026-10-28: filed.
