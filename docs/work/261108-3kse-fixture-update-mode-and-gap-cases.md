---
id: "261108-3kse"
title: fixture update mode and gap cases
status: active
kind: chore
priority: 3
created: 2026-11-08
updated: 2026-11-08
claimed_by: claude-sonnet-5-5-rest
claimed_at: 2026-11-08T20:30
---

# fixture update mode and gap cases

## Goal

Findings F9, F10, F14, F16 and F17 of [the audit](../design/test-strategy-audit.md), from [261030-456q](done/261030-456q-test-strategy-and-coverage-audit.md): the expected responses of the conformance fixtures are edited by hand, and about half the request parameters, the `ls` filters `blocked`, `scope` and `since`, a stale case for each write and six `refused` reasons have no fixture. Add an update mode to the fixture runner first, then add the missing cases, then sync the copy (`scripts/sync-fixtures.sh`).

Done when:

- [ ] `go test` has an update flag that rewrites `expected` from the tool's own output, with the etag handled, and the diff is the review.
- [ ] Fixture cases cover the `ls` filters `blocked`, `scope` and `since`.
- [ ] Every write has a stale (`If-Match`) case, and the six refusal reasons named in F17 each have a case.
- [ ] `npm run schemas` and `TestFixtureCopiesAgree` pass.

## Next

Write the update flag on `TestConformance` (F9), then add the cases.

## Log

- 2026-11-08: filed.
- 2026-11-08 13:30-07:00: claimed (claude-sonnet-5-5-rest).
