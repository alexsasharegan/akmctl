---
id: 261028-mnc0
title: "Offer a claims extension: advisory claimed_by and claimed_at, never a lock"
summary: With the claims extension an active item carries claimed_by and claimed_at, set by akmctl claim. A claim is advisory and stale after 24 hours.
status: accepted
date: 2026-10-28
scope: [internal/api, internal/check]
authored_by: agent
---

# Offer a claims extension: advisory claimed_by and claimed_at, never a lock

## Context

Several agents may pick the same active item. A lock needs a server; a repository has none. (Decision log entry Y2.)

## Considered Options

- No claims: collisions found at merge.
- A hard lock: needs shared state outside git.
- Advisory fields set and cleared by the tool, stale after `claim_stale_hours`.

## Decision

`akmctl claim <id> --as <name>` moves the item to `active` and sets both fields; any move out of `active` and `release` clear them; each claim is logged. Rules CL1 and CL2 warn on an active item with no claim and on one holder with several active items. Without the tool the first Log line says `claimed (where)`.

## Consequences

- Two branches can claim the same item and meet only at merge, where the field conflicts.
- Agents have no stable identity; the name is a free string.

## Confirmation

Claims tests in `internal/api`; rules CL1 and CL2 in `internal/check`.
