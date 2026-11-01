---
id: "261101-3vwz"
title: closed items read only with correction
summary: "set and section refuse a closed work item unless the call carries a correction reason, which is written to the Log in the same write (ADR 261101-dxcn)"
status: done
kind: feature
parent: "261031-gtw6"
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
decisions: ["[[261101-dxcn-closed-records-read-only-except-corrections|261101-dxcn]]"]
---

# closed items read only with correction

## Goal

Implement ADR [261101-dxcn](../../adr/261101-dxcn-closed-records-read-only-except-corrections.md) once it is accepted.

- [x] `set` on a closed work item is `refused` (`reason: "closed"`) unless the call has `correction`; `keep` needs none.
- [x] With `correction`, `set` and `section` change only `title`, `summary` and prose sections; any other field, and `Log`, stay refused.
- [x] The Log line `corrected <target>: <reason>` is written in the same write; a missing Log is created.
- [x] Fixtures (both copies), schemas, API spec 4.11 and 4.16, the `set` and `section` skill references (and embedded copies), glossary, changelog.

## Next

Tick, close.

## Outcome

Done: ADR 261101-dxcn is implemented.

- **Refusal.** `set` and `section` on a `done` or `dropped` work item are `refused` (`reason: "closed"`) unless the call has `correction`; `keep` needs none. `status` stays WK9. `tick`, `log` and `remove` stay refused on a closed item.
- **Correction.** With `--correction "<why>"` only `title`, `summary` and prose sections change (`closed_field` for any other field and for the Log), and `corrected <what>: <why>` is appended to the Log in the same write. A call that changes nothing writes nothing.
- **Checked:** `go test -race ./...`, schemas (102 cases), records check, Prettier. New fixtures `set-12` to `set-14`, `section-8`, `section-9`; `section-6` message changed.
- **Found:** one existing test (`TestClosedItemMayJoinAClosedParent`) set a closed item's parent; the rule now refuses it, so it became `TestClosedItemKeepsItsParent`. A corrected `title` no longer matches the first heading (WK12, a warning, as for an open item); the correction does not rewrite the heading.

Decision: ADR 261101-dxcn, accepted.

## Log

- 2026-11-01: filed by the spike 261031-z0n8; waits for ADR 261101-dxcn to be accepted.
- 2026-11-01 09:40-07:00: claimed (claude).
- 2026-11-01 09:45-07:00: ticked "`set` on a closed work item is `refused` (`reason: "closed"`) unless the call has `correction`; `keep` needs none.": set-12, TestClosedItemCorrection: set on a closed item is refused (closed) without correction; keep needs none
- 2026-11-01 09:45-07:00: ticked "With `correction`, `set` and `section` change only `title`, `summary` and prose sections; any other field, and `Log`, stay refused.": set-13, set-14, section-8, section-9: title, summary and prose sections only; closed_field otherwise and for the Log
- 2026-11-01 09:45-07:00: ticked "The Log line `corrected <target>: <reason>` is written in the same write; a missing Log is created.": set-13 and section-8 expect the corrected Log line in the same write; no-op writes nothing (TestClosedItemCorrection)
- 2026-11-01 09:45-07:00: ticked "Fixtures (both copies), schemas, API spec 4.11 and 4.16, the `set` and `section` skill references (and embedded copies), glossary, changelog.": fixtures in both copies, schemas set/section, spec 4.11/4.16, set.md/section.md and embedded copies, guide, glossary Correction, changelog
