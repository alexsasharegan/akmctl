---
id: "261108-x2c4"
title: audit test maintainability
status: ready
kind: spike
parent: 261030-456q
priority: 3
created: 2026-11-08
updated: 2026-11-08
summary: "Part 2 of the test audit: how hard each suite is to change, which tests depend on time, git or Node, and how long a full run takes"
---

# audit test maintainability

## Goal

Part 2 of the audit in [261030-456q](261030-456q-test-strategy-and-coverage-audit.md): how hard each suite is to change. Look at the fixtures and their sync (one source since 261029-x4ws, and `make_fixtures.py` as a frozen reference), helper duplication across test files, tests that depend on real time, git or Node (and whether they skip silently when Node or git is missing), fixtures that break when a template changes (the `init` written-files lists), and how long a full run takes locally and in CI on each platform.

It adds a "Maintainability" section to `docs/design/test-strategy-audit.md` (create the file if it is missing; [the strategy child](261108-dzh4-audit-test-strategy.md) normally does). It is an audit: no test is written except to measure something.

Done when:

- [ ] The report has a "Maintainability" section that covers each concern named above, with evidence (a path, a command or a number).
- [ ] Run times are measured with the commands recorded in this item's Log and the date.
- [ ] Findings are written as candidates in that section, unnumbered, each with an area, a severity (gap, risk or nuisance), the evidence and a recommended fix.

## Log

- 2026-11-08: filed.
