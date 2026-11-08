---
id: "261108-kmwm"
title: write audit report and file tickets
status: done
kind: spike
parent: 261030-456q
priority: 3
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
summary: "Last part of the test audit: number the findings, group them, file one work item per group, list the ADR drafts and order the roadmap"
---

# write audit report and file tickets

## Goal

The last part of the audit in [261030-456q](261030-456q-test-strategy-and-coverage-audit.md), after the three section children have written their parts of `docs/design/test-strategy-audit.md`. It numbers the candidate findings `F1`, `F2` and so on, groups them (for example conformance fixtures, tooling and maintainability, platform, workbench, skills), files one work item per group that names the findings it addresses, lists the ADR drafts the findings call for (for instance how to fixture non-deterministic or git-dependent operations) without deciding them, and puts the filed items on the roadmap in the order the report recommends. It also collects the parent's seed-list check into one table.

Done when:

- [x] The report ends with a numbered list of findings `F1`, `F2` and so on; each has an area, a severity (gap, risk or nuisance), the evidence and a recommended fix.
- [x] The findings are grouped, and each group has a work item filed from it that names the findings it addresses; the report links the items, and each item links the finding ids.
- [x] The seed-list check from [the strategy child](261108-dzh4-audit-test-strategy.md) is collected in the report, one row per entry.
- [x] The ADR drafts the findings call for are listed in the report, not decided here.
- [x] The roadmap lists the filed items in the order the report recommends.

## Next

## Outcome

The audit report numbers its 19 findings F1 to F19, groups them into seven filed work items, lists four ADR drafts and the seed list, and the roadmap carries the items in the recommended order.

## Log

- 2026-11-08: filed.
- 2026-11-08 01:47-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08 01:48-07:00: numbered the 19 findings F1 to F19, grouped them into seven work items (z9n8, 3kse, qv1w, rwqj, f1tv, 79t5, xajw), listed four ADR drafts, and put the items on the roadmap in the report's order; `npx prettier --write`, `akmctl check`.
