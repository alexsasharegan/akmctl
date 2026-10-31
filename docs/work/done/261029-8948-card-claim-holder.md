---
id: "261029-8948"
title: Show the claim holder on workbench cards
summary: Cards do not show who holds a claim, the snapshot leaves out claimed_by and claimed_at.
status: done
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-10-31
closed: 2026-10-31
---

# Show the claim holder on workbench cards

## Goal

Cards do not show who holds a claim: the snapshot leaves out claimed_by and claimed_at.

Done when:

- [x] Cards show holder and age where claims are on.

## Outcome

Cards show holder and age: board.js adds a chip "held by <name>, <age>" from claimed_by and claimed_at (the snapshot fields in store.js now include both); a record with no claim shows no chip, so repositories without claims are unchanged. Checked in the workbench against a fixture with two claims: the board showed "held by agent-9, 3 days" and "held by claude, 2 days". The node logic tests and Go tests pass. Decision: none.

## Log

- 2026-10-29: filed.
- 2026-10-30: claimed (claude).
