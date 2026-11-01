---
id: "261101-3che"
title: show reports edit states
summary: "A record read carries an edit state (yes, warn, no, with a reason) for each field and section, so the drawer draws controls from the tool's rules"
status: ready
kind: feature
parent: "261031-gtw6"
blocked_by: ["261101-3vwz"]
created: 2026-11-01
updated: 2026-11-01
decisions: ["[[261101-dxcn-closed-records-read-only-except-corrections|261101-dxcn]]"]
---

# show reports edit states

## Goal

The matrix in the Outcome of [261031-z0n8](done/261031-z0n8-spike-field-edit-rules-by-status.md) is data the tool serves, so the drawer ([261028-esyt](261028-esyt-ui-field-editing.md)) holds no copy of it.

- [ ] `show` returns `edit`: for each field and section, `state` (`yes`, `warn`, `no`) and a short `reason` for `warn` and `no`.
- [ ] The states follow the matrix; a test walks every status and field against it.
- [ ] API spec, schema, fixtures (both copies), glossary, changelog.

## Log

- 2026-11-01: filed by the spike 261031-z0n8; waits for 261101-3vwz and ADR 261101-dxcn.
