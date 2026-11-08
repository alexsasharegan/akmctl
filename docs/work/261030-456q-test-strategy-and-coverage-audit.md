---
id: "261030-456q"
title: Test strategy and coverage audit
summary: Audit how akmctl is tested (strategy, maintainability, coverage), list the gaps as numbered findings, and file linked tickets for each group.
status: ready
kind: spike
priority: 3
created: 2026-10-30
updated: 2026-11-08
---

# Test strategy and coverage audit

## Goal

We have let gaps in the tests slide, and the fixtures are the clearest case: `new` is tested for shape only, `commits` has no conformance fixtures, and the spec lists more write fixtures as "not in 1.0". Nobody has looked at the whole. This project looks at it once and produces a report with a hard list of findings, then tickets that address them in groups.

It is an audit, not a fix: no test is written under it except to measure something.

Seed list (known from earlier work; the audit finds the rest, and corrects any entry that is wrong):

- `new` has no exact fixture because ids are random (API spec, section 10); `commits` has none because a fixture directory cannot hold a git history; `similar`, `dependents` and `commit` on `new` are covered by Go tests only.
- Two copies of the API fixtures (work item 261029-x4ws), synced by hand; `make_fixtures.py` is a Python reference that knows two sample repositories, and its output is the contract the Go tests check, so a wrong generator and a wrong implementation can agree.
- The workbench JavaScript is tested through Node from a Go test (`internal/server/ui_assets_test.go`); whether it skips silently when Node is missing, and how much of the UI it covers, is unknown.
- The skills are tested for drift between their two copies, not for what they tell an agent; the instruction text has no test.
- `TestCompareAndSwapRace` is flaky on Windows CI (261029-q8z3); `-race` runs only on Linux and macOS; the Windows path beyond CI is uncovered (261028-jzsx).
- Whether every `check` rule id in the specification has a case that trips it and a case that does not is not recorded anywhere.

Scope of the audit, in three parts:

1. **Strategy.** What each layer is for and whether it is the right one: unit tests, the conformance fixtures, the YAML-lite oracle and fuzz tests, the CLI and HTTP bindings, the workbench tests, the skill drift tests, CI (what runs on which platform, and what is allowed to fail). Where two layers test the same thing and where none does.
2. **Maintainability.** How hard each suite is to change: the generator and the hand-synced copy, helper duplication across test files, tests that depend on real time, git or Node, fixtures that break when a template changes (the `init` written-files lists), and how long a full run takes locally and in CI.
3. **Coverage.** Measured, not guessed: `go test -cover` per package; a matrix of operations by parameter and error code against the fixtures; a matrix of `check` rule ids against tests; the spec's own "not tested" admissions.

Children, in dependency order (the first three can be taken in any order; the last needs them):

- [Audit the test strategy](done/261108-dzh4-audit-test-strategy.md): the layers, their overlaps, and the seed-list check.
- [Audit test maintainability](261108-x2c4-audit-test-maintainability.md): how hard each suite is to change, and run times.
- [Measure test coverage](261108-arn2-measure-test-coverage.md): per-package figures and the operation and rule-id matrices.
- [Write the audit report and file tickets](261108-kmwm-write-audit-report-and-file-tickets.md): numbered findings, groups, filed items, ADR draft list and roadmap order; waits for the other three.

All four write to one report, `docs/design/test-strategy-audit.md`; whichever child comes first creates it, and each adds its own section.

Done when:

- [ ] Every child is `done` or `dropped`.
- [ ] The report in `docs/design/test-strategy-audit.md` covers the three parts and ends with a numbered list of findings `F1`, `F2` and so on, each with an area, a severity (gap, risk or nuisance), the evidence and a recommended fix (child 261108-kmwm).
- [ ] The coverage figures come from commands recorded in the Log, run on the audit's date (child 261108-arn2).
- [ ] The findings are grouped and each group has a work item filed from it, linked both ways (child 261108-kmwm).
- [ ] The seed list is checked entry by entry in the report (child 261108-dzh4, collected by 261108-kmwm).
- [ ] The ADR drafts the findings call for are listed in the report, not decided here (child 261108-kmwm).
- [ ] The roadmap lists the filed items in the order the report recommends (child 261108-kmwm).

## Next

Take 261108-dzh4 (strategy), 261108-x2c4 (maintainability) and 261108-arn2 (coverage) in any order; 261108-kmwm (report and tickets) comes last. Close this item after kmwm.

## Log

- 2026-10-30: filed, from the maintainer's request after `commits` shipped without fixtures; output is an audit report and linked tickets.
- 2026-11-08 00:57-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08 00:57-07:00: Split into four children. The old boxes moved: the report-with-findings box to 261108-kmwm; the coverage-commands box to 261108-arn2; the grouped-items box, the ADR-drafts box and the roadmap box to 261108-kmwm; the seed-list box to 261108-dzh4 (collected by 261108-kmwm). One new box: every child done or dropped.
- 2026-11-08 00:57-07:00: released (claude-sonnet-5-5).
