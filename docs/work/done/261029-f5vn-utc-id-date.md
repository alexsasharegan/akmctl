---
id: "261029-f5vn"
title: Date new ids with the local day, not UTC
status: done
created: 2026-10-29
updated: 2026-10-29
closed: 2026-10-29
kind: bug
priority: 2
summary: akmctl new dates ids in UTC, so an evening in a western time zone gets tomorrow's date.
---

# Date new ids with the local day, not UTC

## Goal

akmctl new and mint_ids take the id date from UTC; in Arizona (MST) after 17:00 local, ids and created dates read as the next day. Use the local day (or as_of when given) and document it.

Done when:

- [x] new and mint_ids use the local calendar day by default.
- [x] A test with a fixed clock and zone covers the evening case.
- [x] The spec says which day an id carries.

## Next

Review the local-day change (clock in internal/api/ops_write.go, Env.now in internal/migrate/mintids.go). Other UTC wording remains (spec section on dates and fields, API spec dates, the date -u shell recipe, id_date, when, WK2) and needs a follow-up decision. Then close with the drafted Outcome.

## Outcome

`new`, `mint_ids`, `updated` and log dates use the local calendar day, or `as_of` when given; `claimed_at` stays UTC. Tests: `internal/api/localday_test.go`, `internal/migrate/localday_test.go`. The spec, API spec, AGENTS.md, guide recipes (`date +%y%m%d`) and `when` now say local day.

## Log

- 2026-10-29: filed.
- 2026-10-29: claimed (life-agent).
- 2026-10-29: new, set, close, claim, release and mint_ids now write the local calendar day (as_of when given); claimed_at stays UTC. Tests: internal/api/localday_test.go, internal/migrate/localday_test.go. Spec: one sentence in the id section.
- 2026-10-29: released (life-agent).
