---
id: "261029-f5vn"
title: Date new ids with the local day, not UTC
status: ready
created: 2026-10-29
updated: 2026-10-29
kind: bug
priority: 2
summary: akmctl new dates ids in UTC, so an evening in a western time zone gets tomorrow's date.
---

# Date new ids with the local day, not UTC

## Goal

akmctl new and mint_ids take the id date from UTC; in Arizona (MST) after 17:00 local, ids and created dates read as the next day. Use the local day (or as_of when given) and document it.

Done when:

- [ ] new and mint_ids use the local calendar day by default.
- [ ] A test with a fixed clock and zone covers the evening case.
- [ ] The spec says which day an id carries.

## Log

- 2026-10-29: filed.
