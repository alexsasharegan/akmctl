---
id: 261028-wdvj
title: Provide set and close for records and no renumber command
summary: akmctl set changes a field and refreshes updated; close sets the status, moves the item to done/ and refuses an empty Outcome. Collisions are repaired by hand.
status: accepted
date: 2026-10-28
scope: [internal/api]
authored_by: agent
---

# Provide set and close for records and no renumber command

## Context

Agents edit fields and close items many times a day; plain-file edits are easy to get subtly wrong (a forgotten `updated`, a move without a rename). A renumber command would have to rewrite references across files. (Decision log entries Y10, Y14.)

## Considered Options

- Plain-file edits only: slow and error-prone.
- `set` and `close`, plus `renumber`.
- `set` and `close`, without `renumber`; document the manual repair.

## Decision

`set <id> field=value` and `close <id> done|dropped` exist, with the plain-file steps as their fallback. There is no `renumber`; the specification documents the manual repair of an id collision.

## Consequences

- Writes are compare-and-swap under the repository lock.
- A collision (rare) costs a manual rename and a search for references.

## Confirmation

Write tests in `internal/api`.
