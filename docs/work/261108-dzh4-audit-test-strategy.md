---
id: "261108-dzh4"
title: audit test strategy
status: active
kind: spike
parent: 261030-456q
priority: 3
created: 2026-11-08
updated: 2026-11-08
summary: "Part 1 of the test audit: what each test layer is for, where layers overlap or leave gaps, and the seed list checked entry by entry"
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-08T08:01
---

# audit test strategy

## Goal

Part 1 of the audit in [261030-456q](261030-456q-test-strategy-and-coverage-audit.md): what each test layer is for, and where two layers test the same thing or none does. The layers are the unit tests per package, the conformance fixtures (`docs/spec/akm-api/.fixtures`, with the `testdata` copy that `scripts/sync-fixtures.sh` writes), the YAML-lite oracle and fuzz tests, the CLI and HTTP bindings, the workbench JavaScript tested through Node from a Go test, the skill drift tests, and CI (what runs on which platform, and what may fail). It also checks the six seed-list entries of the parent.

The report is `docs/design/test-strategy-audit.md`. This item writes its "Strategy" section and creates the file if it is missing; the other audit children add their own sections, and the report child numbers the findings. It is an audit: no test is written except to measure something.

Done when:

- [ ] The report has a "Strategy" section: each layer, what it is for, whether it is the right one, and where layers overlap or none covers a thing, with evidence (a path, a command or a number).
- [ ] Findings are written as candidates in that section, unnumbered, each with an area, a severity (gap, risk or nuisance), the evidence and a recommended fix.
- [ ] Each of the six seed-list entries in the parent's Goal is marked confirmed, corrected or dropped, with the evidence.

## Next

Survey each layer with evidence, check the six seed entries, write the Strategy section of docs/design/test-strategy-audit.md.

## Log

- 2026-11-08: filed.
- 2026-11-08 01:01-07:00: claimed (claude-sonnet-5-5).
