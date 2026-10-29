---
id: "261029-tgfe"
title: Give migrate a discovery phase for repositories with no manifest
summary: Adopting AKM is the migration from layout 0.0.0; for a repository without docs/.layout, migrate first discovers what the repository already has.
status: ready
kind: feature
priority: 2
created: 2026-10-29
updated: 2026-10-29
---

# Give migrate a discovery phase for repositories with no manifest

## Goal

Decided 2026-10-29 ([261029-j0d0](../adr/261029-j0d0-one-tool-spec-and-fixtures-in-this-repository.md)): adopting AKM is the migration from layout 0.0.0, and `akmctl migrate` gets a discovery phase for a repository with no `docs/.layout`. Discovery reads what the repository has (task and ADR folders, id and status conventions, entry-point files, an `AGENTS.md`) and proposes a migration; the scripts in `scripts/vault-migration/` are the seed (they did this by hand for one vault). Dry run by default, no deletes, idempotent, as for every migration.

Done when:

- [ ] `migrate` on a repository with no manifest reports what it found and a proposed set of steps, and changes nothing.
- [ ] The report names what it could not classify, for the agent handoff.
- [ ] Applying the proposal gives a repository that passes `akmctl check`, shown on a fixture.
- [ ] The standard and the guide describe the discovery phase (section 7 of the guide, "An existing repo").

## Log

- 2026-10-29: filed.
