---
id: "261108-arn2"
title: measure test coverage
status: ready
kind: spike
parent: 261030-456q
priority: 3
created: 2026-11-08
updated: 2026-11-08
summary: "Part 3 of the test audit: measured coverage per package, an operation by parameter and error code matrix, and a check rule id matrix"
---

# measure test coverage

## Goal

Part 3 of the audit in [261030-456q](261030-456q-test-strategy-and-coverage-audit.md): coverage, measured and not guessed. Run `go test -cover` per package; build a matrix of operations by parameter and error code against the fixtures; build a matrix of `check` rule ids against tests (a case that trips the rule and a case that does not); list the spec's own "not tested" admissions.

It adds a "Coverage" section to `docs/design/test-strategy-audit.md` (create the file if it is missing). Every figure comes from a command recorded in this item's Log, run on the audit's date. It is an audit: no test is written except to measure something.

Done when:

- [ ] The report has a "Coverage" section with the per-package figures; each command and its date are in this item's Log.
- [ ] The operation matrix and the rule-id matrix are in the report, with the uncovered cells listed.
- [ ] The spec's "not tested" admissions are listed with where each is stated.
- [ ] Findings are written as candidates in that section, unnumbered, each with an area, a severity (gap, risk or nuisance), the evidence and a recommended fix.

## Log

- 2026-11-08: filed.
