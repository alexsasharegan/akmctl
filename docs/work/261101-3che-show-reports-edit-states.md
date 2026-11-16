---
id: "261101-3che"
title: show reports edit states
summary: "A record read carries an edit state (yes, warn, no, with a reason) for each field and section, so the drawer draws controls from the tool's rules"
status: active
kind: feature
parent: "261031-gtw6"
created: 2026-11-01
updated: 2026-11-16
decisions: ["[[261101-dxcn-closed-records-read-only-except-corrections|261101-dxcn]]"]
claimed_by: claude-fable-5-2
claimed_at: 2026-11-16T13:52
---

# show reports edit states

## Goal

The matrix in the Outcome of [261031-z0n8](done/261031-z0n8-spike-field-edit-rules-by-status.md) is data the tool serves, so the drawer ([261028-esyt](done/261028-esyt-ui-field-editing.md)) holds no copy of it.

- [ ] `show` returns `edit`: for each field and section, `state` (`yes`, `warn`, `no`) and a short `reason` for `warn` and `no`.
- [ ] The states follow the matrix; a test walks every status and field against it.
- [ ] API spec, schema, fixtures (both copies), glossary, changelog.

## Next

Add the edit member to opShow (internal/api/editstate.go) from the matrix in 261031-z0n8, with a test that walks every status and field; then the spec, schema, fixtures, glossary and changelog.

## Log

- 2026-11-01: filed by the spike 261031-z0n8; waits for 261101-3vwz and ADR 261101-dxcn.
- 2026-11-01 09:45-07:00: unblocked: 261101-3vwz is closed
- 2026-11-16 06:52-07:00: claimed (claude-fable-5-2).
