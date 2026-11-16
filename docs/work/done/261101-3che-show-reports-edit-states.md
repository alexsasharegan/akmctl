---
id: "261101-3che"
title: show reports edit states
summary: "A record read carries an edit state (yes, warn, no, with a reason) for each field and section, so the drawer draws controls from the tool's rules"
status: done
kind: feature
parent: "261031-gtw6"
created: 2026-11-01
updated: 2026-11-16
closed: 2026-11-16
decisions: ["[[261101-dxcn-closed-records-read-only-except-corrections|261101-dxcn]]"]
---

# show reports edit states

## Goal

The matrix in the Outcome of [261031-z0n8](261031-z0n8-spike-field-edit-rules-by-status.md) is data the tool serves, so the drawer ([261028-esyt](261028-esyt-ui-field-editing.md)) holds no copy of it.

- [x] `show` returns `edit`: for each field and section, `state` (`yes`, `warn`, `no`) and a short `reason` for `warn` and `no`.
- [x] The states follow the matrix; a test walks every status and field against it.
- [x] API spec, schema, fixtures (both copies), glossary, changelog.

## Next

## Outcome

- `show` returns `edit`: `internal/api/editstate.go`, wired in `opShow`; conformance cases `show-1`, `show-3` and the new `show-4` (a closed item) carry it.
- The states follow the matrix: `TestEditStatesFollowTheMatrix` walks every status against every field and section; `TestEditStatesAgreeWithTheWritesOnAClosedItem` checks that what `edit` says matches what `set` does.
- API spec section 4.4, `schemas/show.json`, the fixtures in both copies, the glossary term Edit state and the changelog are in the same change.

Not done here: the drawer still has to read `edit` (its parent item), and an ADR has no `edit`.

Decision: none. The shape of the member is a set of small picks, listed in the Log.

## Log

- 2026-11-01: filed by the spike 261031-z0n8; waits for 261101-3vwz and ADR 261101-dxcn.
- 2026-11-01 09:45-07:00: unblocked: 261101-3vwz is closed
- 2026-11-16 06:52-07:00: claimed (claude-fable-5-2).
- 2026-11-16 06:54-07:00: ticked "`show` returns `edit`: for each field and section, `state` (`yes`, `warn`, `no`) and a short `reason` for `warn` and `no`.": internal/api/editstate.go, the edit member of opShow; fixture show-4
- 2026-11-16 06:54-07:00: ticked "The states follow the matrix; a test walks every status and field against it.": TestEditStatesFollowTheMatrix walks every status and field; TestEditStatesAgreeWithTheWritesOnAClosedItem
- 2026-11-16 06:54-07:00: ticked "API spec, schema, fixtures (both copies), glossary, changelog.": API spec 4.4, schemas/show.json, cases show-1, show-3, show-4 in both copies, glossary Edit state, CHANGELOG
- 2026-11-16 06:54-07:00: picks, for the maintainer to review at merge: (1) edit is a top-level member of show with fields and sections keyed by name; (2) a closed item's title, summary and sections are state no with correction true, so the three states stay as the item asks; (3) the Log is yes with append true while open; (4) status is yes while open because set moves it, no when closed; (5) next on a closed item is no, as the matrix says, although section still accepts it with a correction; (6) an ADR has no edit yet: the matrix covers work items only
