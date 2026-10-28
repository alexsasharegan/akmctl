---
id: 261028-6e73
title: Migrate layouts with declarative steps embedded in the binary and a safe apply
summary: akmctl migrate chains declarative steps between layout versions, applies them on a backup branch and tag, logs them, never commits. Read before adding a step.
status: accepted
date: 2026-10-28
scope: [internal/migrate, internal/api]
authored_by: agent
---

# Migrate layouts with declarative steps embedded in the binary and a safe apply

## Context

Only a MAJOR layout change needs a migration. A migration must be auditable, testable on fixtures and safe on a maintainer's repository. It must work without a network. (Decision log entries Y5, Y36.)

## Considered Options

- Arbitrary scripts per migration: powerful, not auditable.
- Declarative steps (rename_field, set_default, rename_heading, move_file, add_file, manual) embedded in the binary, with a handoff prompt for the manual remainder.
- No migration support: each adopter converts by hand.

## Decision

Steps are data, chained between layout versions. `apply_safe` needs `acknowledge_backup` and a clean git tree, makes a tag and a branch, logs to `docs/.migrate-log` and writes by temp file and rename. The manifest moves only when no manual step remains. Release 0.1.0 embeds no migrations; `check --target` classifies findings by the steps' `fixes`.

## Consequences

- A newer migration needs a newer tool.
- A repository with no manifest (layout 0.0.0) cannot be migrated to 1.0.0; adoption is a one-time conversion.
- `check --target` accepts only the tool's own layout version (a work item tracks it).

## Confirmation

`internal/migrate` and `internal/api` migrate tests.
