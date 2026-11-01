---
id: "261031-x50y"
title: spike schema declaration and migration
summary: "Spike: where the body schema is declared, how check enforces it, and how existing items move to it"
status: raw
kind: spike
parent: "261031-gtw6"
created: 2026-10-31
updated: 2026-10-31
---

# spike schema declaration and migration

## Goal

Once the body has a model ([261031-qbat](261031-qbat-spike-work-item-body-model.md)), someone has to declare it, check it and carry old records to it. Find out:

- Where it lives: in the spec and the tool only, in `docs/.layout`, or in an extension a repository can add to (custom sections and fields for its own workflow). Whether a repository may add sections and fields, and how that is validated.
- Enforcement: new `check` rules (a rule per violation class, with a rule id), their severity, and whether the rules apply to closed items. How a hand-edited file that breaks the schema is reported and repaired.
- A machine-readable schema for the API (the `docs/spec/akm-api/` schemas), so an agent can fetch the section list for a kind instead of reading prose.
- Migration: the existing items, the `migrate` command, and a version bump of the layout; what happens to a body that does not fit.
- Scope: work items first, or ADRs, ideas and design notes as well. Recommend an order.

Done when:

- [ ] `check` is prototyped with one or two schema rules and run over this repository's items, with the findings counted.
- [ ] A migration of this repository's items is planned with `migrate --mode plan`, or a manual equivalent, and its output is shown.
- [ ] A recommendation (home of the schema, rules, API shape, migration, record types in scope) is in `## Outcome`; a hard-to-reverse choice has a `proposed` ADR; follow-up items are filed.

## Log

- 2026-10-31: filed.
