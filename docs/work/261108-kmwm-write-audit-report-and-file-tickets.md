---
id: "261108-kmwm"
title: write audit report and file tickets
status: blocked
kind: spike
parent: 261030-456q
priority: 3
created: 2026-11-08
updated: 2026-11-08
summary: "Last part of the test audit: number the findings, group them, file one work item per group, list the ADR drafts and order the roadmap"
blocked_by: [261108-arn2]
---

# write audit report and file tickets

## Goal

The last part of the audit in [261030-456q](261030-456q-test-strategy-and-coverage-audit.md), after the three section children have written their parts of `docs/design/test-strategy-audit.md`. It numbers the candidate findings `F1`, `F2` and so on, groups them (for example conformance fixtures, tooling and maintainability, platform, workbench, skills), files one work item per group that names the findings it addresses, lists the ADR drafts the findings call for (for instance how to fixture non-deterministic or git-dependent operations) without deciding them, and puts the filed items on the roadmap in the order the report recommends. It also collects the parent's seed-list check into one table.

Done when:

- [ ] The report ends with a numbered list of findings `F1`, `F2` and so on; each has an area, a severity (gap, risk or nuisance), the evidence and a recommended fix.
- [ ] The findings are grouped, and each group has a work item filed from it that names the findings it addresses; the report links the items, and each item links the finding ids.
- [ ] The seed-list check from [the strategy child](done/261108-dzh4-audit-test-strategy.md) is collected in the report, one row per entry.
- [ ] The ADR drafts the findings call for are listed in the report, not decided here.
- [ ] The roadmap lists the filed items in the order the report recommends.

## Next

Wait for the strategy, maintainability and coverage children; then number their candidate findings and write the grouping.

## Log

- 2026-11-08: filed.
