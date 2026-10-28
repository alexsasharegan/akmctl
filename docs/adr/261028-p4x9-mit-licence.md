---
id: 261028-p4x9
title: License the tool under the MIT licence
summary: akmctl is MIT-licensed (copyright Alexander Regan); the Unicode data under internal/width/data keeps the Unicode licence. Read before accepting a contribution.
status: accepted
date: 2026-10-28
scope: [LICENSE, internal/width/data]
authored_by: agent
---

# License the tool under the MIT licence

## Context

The tool is meant to be adopted widely if it is ever published. A licence is hard to change once others hold copies. (Decision log entry Y35.)

## Considered Options

- Apache-2.0: patent grant, NOTICE attribution.
- MPL-2.0: modified files stay open.
- MIT: the widest adoption; does not stop a closed derivative.
- No licence: nobody may use it.

## Decision

MIT. The maintainer accepts that a closed derivative is possible. The Unicode data files keep the Unicode licence.

## Consequences

- Simple terms for adopters.
- The licence applies to the tool; the AKM standard text has its own terms.

## Confirmation

`LICENSE` at the repository root.
