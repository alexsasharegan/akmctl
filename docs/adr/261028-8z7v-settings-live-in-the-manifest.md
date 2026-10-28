---
id: 261028-8z7v
title: Keep settings as flat keys in docs/.layout; the status enum and roadmap sections are fixed
summary: Settings (stale_days, claim_stale_hours, AGENTS.md limits, reference_style) are flat keys in docs/.layout. No second config file. Read before adding one.
status: accepted
date: 2026-10-28
scope: [docs/.layout, internal/record]
authored_by: agent
---

# Keep settings as flat keys in docs/.layout; the status enum and roadmap sections are fixed

## Context

Tunable thresholds need a home. A repository already has a manifest recording versions and extensions. (Decision log entry Y6.)

## Considered Options

- A separate `.akm.toml`: a second file and a second parser.
- Flat keys in the existing manifest, in the frontmatter subset.
- Hard-coded values.

## Decision

Settings are flat keys in the manifest. Absent keys mean the defaults. The status enum and the roadmap sections are fixed in layout 1.

## Consequences

- One file and one parser for versions, extensions and settings.
- A new status or section needs a new layout version.

## Confirmation

`internal/record` loads the manifest; `check` honours the keys.
