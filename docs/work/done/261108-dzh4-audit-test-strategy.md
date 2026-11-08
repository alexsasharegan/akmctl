---
id: "261108-dzh4"
title: audit test strategy
status: done
kind: spike
parent: 261030-456q
priority: 3
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
summary: "Part 1 of the test audit: what each test layer is for, where layers overlap or leave gaps, and the seed list checked entry by entry"
---

# audit test strategy

## Goal

Part 1 of the audit in [261030-456q](261030-456q-test-strategy-and-coverage-audit.md): what each test layer is for, and where two layers test the same thing or none does. The layers are the unit tests per package, the conformance fixtures (`docs/spec/akm-api/.fixtures`, with the `testdata` copy that `scripts/sync-fixtures.sh` writes), the YAML-lite oracle and fuzz tests, the CLI and HTTP bindings, the workbench JavaScript tested through Node from a Go test, the skill drift tests, and CI (what runs on which platform, and what may fail). It also checks the six seed-list entries of the parent.

The report is `docs/design/test-strategy-audit.md`. This item writes its "Strategy" section and creates the file if it is missing; the other audit children add their own sections, and the report child numbers the findings. It is an audit: no test is written except to measure something.

Done when:

- [x] The report has a "Strategy" section: each layer, what it is for, whether it is the right one, and where layers overlap or none covers a thing, with evidence (a path, a command or a number).
- [x] Findings are written as candidates in that section, unnumbered, each with an area, a severity (gap, risk or nuisance), the evidence and a recommended fix.
- [x] Each of the six seed-list entries in the parent's Goal is marked confirmed, corrected or dropped, with the evidence.

## Next

## Outcome

The Strategy section of docs/design/test-strategy-audit.md covers nine layers and eight candidate findings (C1 to C8), and checks the seed list: two entries corrected, two partly corrected and the rest confirmed. The main gaps: the bindings do not run the fixtures, three operations have no fixture, tests that need node, git or a corpus skip silently, 11 of 25 workbench modules have no test, and the install and release scripts have none. The report child numbers them.

## Log

- 2026-11-08: filed.
- 2026-11-08 01:01-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08 01:03-07:00: ticked "The report has a "Strategy" section: each layer, what it is for, whether it is the right one, and where layers overlap or none covers a thing, with evidence (a path, a command or a number).": docs/design/test-strategy-audit.md, section Strategy: nine layers with a verdict each
- 2026-11-08 01:03-07:00: ticked "Findings are written as candidates in that section, unnumbered, each with an area, a severity (gap, risk or nuisance), the evidence and a recommended fix.": eight candidates C1 to C8, each with area, severity, evidence and fix
- 2026-11-08 01:03-07:00: ticked "Each of the six seed-list entries in the parent's Goal is marked confirmed, corrected or dropped, with the evidence.": seed table in the report: seven entries, each confirmed, corrected or partly corrected, with evidence
- 2026-11-08 01:03-07:00: Evidence commands (2026-11-08, main at ee8a99f): ls internal/_/__test.go | wc -l (71); grep -c '^func Test' per package; ls docs/spec/akm-api/.fixtures/cases (109 files, none for commits, history or skills); grep -rn 'testdata/api/cases' internal (one reader, api.Run only); grep -rn 'Skip' internal --include='*_test.go'; ls internal/server/web/js (25 modules) against the 14 in TestPureModulesUnderNode; grep -rln 'install.sh' internal scripts .github (none in tests); .github/workflows/ci.yml read for -race and continue-on-error; a script over the spec rules table and internal/check tests (44 ids, all named). The seed list had seven entries once the lines are split; the table in the report keeps them in order. Not rechecked: the Windows flaky test (Windows work is on hold).
