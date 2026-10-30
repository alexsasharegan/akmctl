---
id: "261029-5pe1"
title: Dates are the writer's local day; only timestamps that are compared across machines are UTC
summary: Dates are the writer's local day; only timestamps that are compared across machines are UTC
status: proposed
date: 2026-10-29
---

# Dates are the writer's local day; only timestamps that are compared across machines are UTC

## Context

ADR 261027-ensf set two rules that disagreed. Its id paragraph said an id carries "the UTC creation date", while its field paragraph said date fields hold the local calendar day of the writing machine and timestamps are UTC. The first implementation followed the id paragraph: `new`, `mint_ids` and the `created` date came from UTC.

That produced the Phoenix evening trap. Arizona is on MST (UTC-7) all year, so after 17:00 local the UTC date is already tomorrow. Items filed on the evening of 2026-10-29 carried the id date `261030` and `created: 2026-10-30`, while `updated` (written from the local clock) said 2026-10-29, so a record could be updated before it was created. Work item 261029-f5vn fixed the id and date writers to use the local day, and the spec was edited to match. No ADR recorded the reason, and ADR 261027-ensf still reads "UTC creation date".

## Considered Options

- Everything in UTC. Rejected: it is the trap above. A person reading `created` sees a day that was not their day, and an evening session starts on tomorrow's date.
- Everything local, timestamps included. Rejected: a claim is stale when it is older than `claim_stale_hours`, which two machines in different zones (or across a daylight-saving change) must compute the same way. A local wall-clock time cannot be compared.
- Timestamps with an explicit offset (`2026-10-29T17:30-07:00`). Rejected: Obsidian and the YAML readers in our tests show them poorly, the format is longer, and the only comparison we need is "how long ago", which UTC answers.
- Local day for dates a human reads, UTC for the one timestamp a machine compares.

## Decision

The guiding principle: **use UTC where machines must compare values; use the writer's local day where a person reads the value.**

- Date fields (`created`, `updated`, `closed`, an ADR's `date`, log-line dates) and the date in an id are the **local calendar day of the machine that writes them**, or the `as_of` day when the request gives one. They are `YYYY-MM-DD` and are not timestamps.
- `claimed_at` is the only timestamp. It is UTC, `YYYY-MM-DDTHH:MM`, because staleness is compared across machines.
- Ids do not depend on the date for uniqueness. Uniqueness comes from the four random characters and a check that the id does not already exist, so a day boundary or a zone difference cannot make two writers collide.
- The tool takes its clock from one place (the request's `as_of`, else the system clock) and tests fix both the clock and the zone, including the evening case.
- Hand-written dates use the local day too; the fallback recipe is `date +%y%m%d`, not `date -u`.

This ADR replaces the words "the UTC creation date" in the id paragraph of ADR 261027-ensf with the local day. That ADR is accepted and its body is not edited.

## Consequences

- Two writers in different zones may write different dates for the same instant, and ids sort by local creation day, not by the instant. Neither matters for uniqueness or for the 30-day staleness rule, whose error is at most one day.
- A new field that holds a moment machines compare must be a UTC timestamp; one that a person reads is a local date. A proposal that breaks the rule needs a reason in its own ADR.
- A reader of ADR 261027-ensf finds the correction only through this ADR; the spec already states the local-day rule.
