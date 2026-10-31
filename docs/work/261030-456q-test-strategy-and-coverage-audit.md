---
id: "261030-456q"
title: Test strategy and coverage audit
summary: Audit how akmctl is tested (strategy, maintainability, coverage), list the gaps as numbered findings, and file linked tickets for each group.
status: ready
kind: spike
priority: 3
created: 2026-10-30
updated: 2026-10-30
---

# Test strategy and coverage audit

## Goal

We have let gaps in the tests slide, and the fixtures are the clearest case: `new` is tested for shape only, `commits` has no conformance fixtures, and the spec lists more write fixtures as "not in 1.0". Nobody has looked at the whole. This item looks at it once and produces a report with a hard list of findings, then tickets that address them in groups.

It is an audit, not a fix: no test is written under this item except to measure something.

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

Done when:

- [ ] A report in `docs/design/test-strategy-audit.md` covers the three parts and ends with a numbered list of findings `F1`, `F2`, and so on. Each finding has an area, a severity (gap, risk or nuisance), the evidence (a command, a path or a number), and a recommended fix.
- [ ] The coverage figures in the report come from commands recorded in the Log, run on the audit's date.
- [ ] The findings are grouped (for example conformance fixtures, tooling and maintainability, platform, workbench, skills), and each group has a work item filed from it that names the findings it addresses. The report links the items, and each item links the finding ids.
- [ ] The seed list above is checked entry by entry in the report: confirmed, corrected or dropped.
- [ ] The ADR drafts that the findings call for (a decision on how to fixture non-deterministic or git-dependent operations, for instance) are listed in the report, not decided here.
- [ ] The roadmap lists the filed items in an order the report recommends.

## Log

- 2026-10-30: filed, from the maintainer's request after `commits` shipped without fixtures; output is an audit report and linked tickets.
