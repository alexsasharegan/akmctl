---
id: "261028-k0ek"
title: Give the workbench live refresh
summary: The workbench page shows a snapshot; it does not update when records change on disk.
status: raw
kind: feature
priority: 0
created: 2026-10-28
updated: 2026-10-28
---

# Give the workbench live refresh

## Goal

The workbench page loads the records once. If a record changes on disk (an agent edits it, `akmctl` writes it), the page stays stale until it is reloaded. The planned UI port will decide between server-sent events and polling; this item records the gap in the current page.

Done when:

- [ ] The page updates within a few seconds of a record changing on disk, without a reload.
- [ ] A change made by the page itself does not flicker or lose the user's place.

## Log

- 2026-10-28: filed.
