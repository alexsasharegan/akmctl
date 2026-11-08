---
id: "261029-q8z3"
title: TestCompareAndSwapRace is flaky on Windows CI
summary: The compare-and-swap race test fails now and then on the Windows CI job.
status: active
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-11-08
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-08T16:29
---

# TestCompareAndSwapRace is flaky on Windows CI

## Goal

The compare-and-swap race test fails now and then on the Windows CI job.

Done when:

- [ ] The test passes reliably on all three CI platforms, or the cause is recorded.

## Next

Stress the test on Windows CI to see the failure, fix the cause (not the test's tolerance), then drop continue-on-error.

## Log

- 2026-10-29: filed.
- 2026-11-08 09:29-07:00: claimed (claude-sonnet-5-5).
