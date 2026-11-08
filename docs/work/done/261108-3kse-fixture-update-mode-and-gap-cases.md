---
id: "261108-3kse"
title: fixture update mode and gap cases
status: done
kind: chore
priority: 3
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
---

# fixture update mode and gap cases

## Goal

Findings F9, F10, F14, F16 and F17 of [the audit](../../design/test-strategy-audit.md), from [261030-456q](261030-456q-test-strategy-and-coverage-audit.md): the expected responses of the conformance fixtures are edited by hand, and about half the request parameters, the `ls` filters `blocked`, `scope` and `since`, a stale case for each write and six `refused` reasons have no fixture. Add an update mode to the fixture runner first, then add the missing cases, then sync the copy (`scripts/sync-fixtures.sh`).

Done when:

- [x] `go test` has an update flag that rewrites `expected` from the tool's own output, with the etag handled, and the diff is the review.
- [x] Fixture cases cover the `ls` filters `blocked`, `scope` and `since`.
- [x] Every write has a stale (`If-Match`) case, and the six refusal reasons named in F17 each have a case.
- [x] `npm run schemas` and `TestFixtureCopiesAgree` pass.

## Next

## Outcome

TestConformance has an -update flag that rewrites expected (and after, would_be) from the tool's output in both fixture copies; verified by breaking an etag and an after text and getting the files back byte for byte. 20 cases added (ls blocked, scope, since; stale for seven writes; dry_run for claim, log, release; five refusal reasons) with fixture repositories commitson and nomanifest; scopeMatch has a unit test. if_match_required is an HTTP reason and waits for 261108-f1tv; the rest of F16's unnamed parameters are not covered.

## Log

- 2026-11-08: filed.
- 2026-11-08 13:30-07:00: claimed (claude-sonnet-5-5-rest).
- 2026-11-08 13:35-07:00: TestConformance -update rewrites expected (and after, would_be for a diff) in the source fixtures and the copy; checked by breaking an etag and an after text and getting the files back byte for byte. 20 new cases: ls blocked/scope/since (ls-8 to ls-12), stale for section, tick, remove, log, close, claim, release, dry_run for claim, log, release, and the refusals extension_off, needs_confirmation, no_manifest, not_git, ticked_box (new fixture repos commitson and nomanifest). if_match_required is HTTP-only and waits for f1tv (both bindings); scopeMatch has a unit test. Not done: F16's remaining unnamed parameters beyond these.
