---
id: "261108-arn2"
title: measure test coverage
status: done
kind: spike
parent: 261030-456q
priority: 3
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
summary: "Part 3 of the test audit: measured coverage per package, an operation by parameter and error code matrix, and a check rule id matrix"
---

# measure test coverage

## Goal

Part 3 of the audit in [261030-456q](../261030-456q-test-strategy-and-coverage-audit.md): coverage, measured and not guessed. Run `go test -cover` per package; build a matrix of operations by parameter and error code against the fixtures; build a matrix of `check` rule ids against tests (a case that trips the rule and a case that does not); list the spec's own "not tested" admissions.

It adds a "Coverage" section to `docs/design/test-strategy-audit.md` (create the file if it is missing). Every figure comes from a command recorded in this item's Log, run on the audit's date. It is an audit: no test is written except to measure something.

Done when:

- [x] The report has a "Coverage" section with the per-package figures; each command and its date are in this item's Log.
- [x] The operation matrix and the rule-id matrix are in the report, with the uncovered cells listed.
- [x] The spec's "not tested" admissions are listed with where each is stated.
- [x] Findings are written as candidates in that section, unnumbered, each with an area, a severity (gap, risk or nuisance), the evidence and a recommended fix.

## Next

## Outcome

Measured coverage on 2026-11-08: 86.9 % of statements across all tests, 14 functions never run, 47 of 120 request parameters never named in a fixture. The Coverage section of docs/design/test-strategy-audit.md holds the per-package table, the operation and rule-id matrices, the spec's not-tested admissions and six candidate findings (V1 to V6) for 261108-kmwm to number.

## Log

- 2026-11-08: filed.
- 2026-11-08 01:28-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08: `go clean -testcache; go test -count=1 -cover ./...` (own-test figures); `go test -count=1 -coverpkg=./... -coverprofile=cover.out ./...` then `go tool cover -func=cover.out` (all-test figures, 86.9 % of 7,118 statements); `go tool cover -func` for the zero-coverage and weakest functions; read `internal/api/ops.go` lines 268-299 for the `ls` filter blocks.
- 2026-11-08: built the operation by parameter matrix from `api.Params` against `docs/spec/akm-api/.fixtures/cases/*.json` (109 cases), the error code and refusal reason tables from the spec against case expectations, and the rule id matrix from the spec table against `internal/check` tests.
- 2026-11-08: wrote the Coverage section of `docs/design/test-strategy-audit.md` with six candidate findings (V1 to V6), unnumbered until 261108-kmwm.
