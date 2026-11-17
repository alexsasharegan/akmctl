---
id: "261117-6qn7"
title: One version number for the tool, layout, specification and API before 1.0.0
summary: Until 1.0.0 the tool, the layout, the specification, the extensions and the API share one version number; the API 1.0 becomes 0.4 and its HTTP path /api/v0
status: proposed
date: 2026-11-17
---

# One version number for the tool, layout, specification and API before 1.0.0

## Context

The tool was 0.3.x, the layout, the specification and the extensions 0.3.0, and the API 1.0, each versioned on its own (`docs/architecture.md`). The API's 1.0 sat beside everything else at 0.x. Readers, and language models reading the specification, took it as a stable surface next to an unstable one, or as a typo, and asked which number to trust. The README says that nothing below 1.0.0 is stable, so the API was the one component whose number said otherwise.

Independent numbers are a good rule once there is something to keep stable, because a client then needs to know only the API's MAJOR. Before 1.0.0 nothing is promised, and the cost of tracking four numbers is paid by every reader.

## Considered Options

- Keep the numbers independent and add a compatibility table. Rejected: it documents the confusion instead of removing it, and the table needs upkeep at every release.
- Move only the API down to 0.x and leave the layout on its own number. Rejected: the layout and the specification already follow the tool's beta number by hand, and they drift on the patch number.
- One number for all of them until 1.0.0.

## Decision

Until 1.0.0 the tool, the layout, the specification and guide, the extensions and the API carry one MAJOR.MINOR. A release is `0.N.P`: the tool is `0.N.P`, and the layout, the extensions and the specification are `0.N.0`, because a PATCH never changes them. The API's `api_version` is `0.N`.

- The first release under this rule is 0.4.0. The layout and the extensions go from 0.3.0 to 0.4.0 and the API from 1.0 to 0.4. A MINOR bump of the tool bumps all of them, even when a component has no change, so that the number says which release a reader holds.
- `api_version` keeps its `MAJOR.MINOR` form, so a client that reads the field needs no change but the value. The HTTP path carries the MAJOR, so it becomes `/api/v0`, and `--api-version 0` pins it. The schema `$id` URIs read `https://akm.invalid/api/0.4/`.
- No migration is added for the layout. Nothing in a 0.3.0 repository changes, only the number, and the beta supports new repositories only (the README). `check` warns (S1) that a 0.3.0 manifest has no migration path and `migrate` reports it as unsupported, as it did for 0.2.0 after the 0.3.0 bump; a repository raises the versions in its manifest by hand. A migration with a single informational step would remove the warning; add one if a beta user asks.
- At 1.0.0 all of them become 1.0.0, and the numbers separate only by the compatibility policy that the work item "State what a MINOR and a MAJOR change mean" is drafting. This decision does not say what happens after 1.0.0.

## Consequences

- A change a client can see: `api_version` is `"0.4"` in every response, the HTTP routes move from `/api/v1` to `/api/v0`, and `--api-version 1` is refused. The workbench uses the new path in the same release. A client written against 1.0 must change the value it checks and the path it calls.
- The API's MINOR and MAJOR rules (specification, section 8) apply with the MAJOR at 0: before 1.0.0 any release may break a client, as the README already says.
- The fixtures, the schemas, the specification headers, the README and the architecture note change in one commit; every release item from now on raises all of them together (`docs/release.md`).
- The cost is that the layout and API numbers move when their content did not. That is accepted for the clarity of one number.
- Supersede this decision, rather than edit it, if the compatibility policy settles on independent numbers after 1.0.0.
