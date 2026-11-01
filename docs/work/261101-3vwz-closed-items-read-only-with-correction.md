---
id: "261101-3vwz"
title: closed items read only with correction
summary: "set and section refuse a closed work item unless the call carries a correction reason, which is written to the Log in the same write (ADR 261101-dxcn)"
status: ready
kind: feature
parent: "261031-gtw6"
created: 2026-11-01
updated: 2026-11-01
decisions: ["[[261101-dxcn-closed-records-read-only-except-corrections|261101-dxcn]]"]
---

# closed items read only with correction

## Goal

Implement ADR [261101-dxcn](../adr/261101-dxcn-closed-records-read-only-except-corrections.md) once it is accepted.

- [ ] `set` on a closed work item is `refused` (`reason: "closed"`) unless the call has `correction`; `keep` needs none.
- [ ] With `correction`, `set` and `section` change only `title`, `summary` and prose sections; any other field, and `Log`, stay refused.
- [ ] The Log line `corrected <target>: <reason>` is written in the same write; a missing Log is created.
- [ ] Fixtures (both copies), schemas, API spec 4.11 and 4.16, the `set` and `section` skill references (and embedded copies), glossary, changelog.

## Log

- 2026-11-01: filed by the spike 261031-z0n8; waits for ADR 261101-dxcn to be accepted.
