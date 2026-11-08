---
id: "261108-79t5"
title: skill text checks
status: done
kind: chore
priority: 3
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
---

# skill text checks

## Goal

Finding F6 of [the audit](../../design/test-strategy-audit.md), from [261030-456q](261030-456q-test-strategy-and-coverage-audit.md): the tests check that skill copies agree, not what the text tells an agent. Two tool snags (261108-h7ds, 261108-kdz1) were found only by a cold-agent retest. Record a small scripted cold-agent check as a procedure and add a test for each snag a retest finds.

Done when:

- [x] A recorded procedure for the cold-agent check exists (where it lives is decided in the item).
- [x] A test pins each snag found so far (261108-h7ds, 261108-kdz1).

## Next

## Outcome

The cold-agent procedure is recorded in docs/cold-agent-check.md (linked from AGENTS.md, term in the glossary). Both known snags (h7ds, kdz1) were already pinned by tests; the doc lists them.

## Log

- 2026-11-08: filed.
- 2026-11-08 15:29-07:00: claimed (claude-sonnet-5-5-rest).
- 2026-11-08 15:31-07:00: ticked "A test pins each snag found so far (261108-h7ds, 261108-kdz1).": already pinned: TestListFieldValueForms, TestSetBlockedByStoresACleanList (h7ds); TestSetTitleRenamesTheHeading (kdz1); listed in the doc's table
- 2026-11-08 15:31-07:00: Procedure lives in docs/cold-agent-check.md (maintainer's choice over a new skill, which would ship to users and touch the mirror and init fixtures).
- 2026-11-08 15:31-07:00: ticked "A recorded procedure for the cold-agent check exists (where it lives is decided in the item).": docs/cold-agent-check.md, linked from AGENTS.md; glossary term Cold-agent check
