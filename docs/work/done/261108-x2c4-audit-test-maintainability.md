---
id: "261108-x2c4"
title: audit test maintainability
status: done
kind: spike
parent: 261030-456q
priority: 3
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
summary: "Part 2 of the test audit: how hard each suite is to change, which tests depend on time, git or Node, and how long a full run takes"
---

# audit test maintainability

## Goal

Part 2 of the audit in [261030-456q](../261030-456q-test-strategy-and-coverage-audit.md): how hard each suite is to change. Look at the fixtures and their sync (one source since 261029-x4ws, and `make_fixtures.py` as a frozen reference), helper duplication across test files, tests that depend on real time, git or Node (and whether they skip silently when Node or git is missing), fixtures that break when a template changes (the `init` written-files lists), and how long a full run takes locally and in CI on each platform.

It adds a "Maintainability" section to `docs/design/test-strategy-audit.md` (create the file if it is missing; [the strategy child](261108-dzh4-audit-test-strategy.md) normally does). It is an audit: no test is written except to measure something.

Done when:

- [x] The report has a "Maintainability" section that covers each concern named above, with evidence (a path, a command or a number).
- [x] Run times are measured with the commands recorded in this item's Log and the date.
- [x] Findings are written as candidates in that section, unnumbered, each with an area, a severity (gap, risk or nuisance), the evidence and a recommended fix.

## Next

## Outcome

The Maintainability section of docs/design/test-strategy-audit.md has measured run times, the helper and clock/git/Node survey, and a measured table of what breaks when a template or fixture changes. Five candidate findings (M1 to M5): no update mode for the expected responses, six init cases that list written files by hand, seven repository-copy helpers that ignore errors, a serial file-copy-bound internal/api, and a race run dominated by internal/body.

## Log

- 2026-11-08: filed.
- 2026-11-08 01:14-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08 01:20-07:00: ticked "The report has a "Maintainability" section that covers each concern named above, with evidence (a path, a command or a number).": docs/design/test-strategy-audit.md, section Maintainability: run times, helpers, time/git/Node, template and fixture coupling
- 2026-11-08 01:20-07:00: ticked "Run times are measured with the commands recorded in this item's Log and the date.": commands and dates recorded in this Log; table of run times in the report
- 2026-11-08 01:20-07:00: ticked "Findings are written as candidates in that section, unnumbered, each with an area, a severity (gap, risk or nuisance), the evidence and a recommended fix.": five candidates M1 to M5, each with area, severity, evidence and fix
- 2026-11-08 01:20-07:00: Run times, 2026-11-08, main at 75c8b7d: go clean -testcache; go test ./... (22.0 s real; internal/api 19.4 s); go clean -testcache; go test -race ./... (33.8 s; internal/body 31.8 s); go test -count=1 -v ./internal/api for the slowest tests; gh run view <id> --json jobs for CI (Linux 73 to 108 s, macOS 27 to 54 s, Windows 101 to 115 s). Template coupling measured in a throwaway git worktree, one edit at a time, go test -count=1 ./...: reworded AGENTS.md sentence (nothing fails); renamed a required heading (3 tests); edited one skill copy (TestSkillsMatchRepository); added a reference file to both skill copies (6 init cases). My first try used a sed range macOS does not support and edited nothing; it was redone with confirmed edits.
