---
id: "261101-p2m8"
title: claim refuses item without next
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
kind: bug
parent: 261101-vs5c
summary: claim leaves an item active with no Next, a state check then rejects (WK6)
decisions: [261101-09yd]
---

# claim refuses item without next

## Goal

Seen in the first Codex trial: `akmctl claim 261101-3erc --as codex` succeeded and printed `active`, then `akmctl check` failed with `WK6 status active requires a non-empty ## Next`. The tool created a state its own checker forbids. Codex filed this as 261101-2n55 on the PR 1 branch (empty Goal); this item replaces it. Done when:

- [x] `claim` of an item whose `## Next` is empty is refused with `refused`, naming WK6 and saying to write `## Next` first (the item stays `ready`).
- [x] The API spec, the `claim` reference in both skill copies and a fixture or test cover the refusal; a `claim` of an item with a `## Next` is unchanged.
- [x] A `ready` item rule is checked too: say in the spec whether `ready` requires `## Next` (a cold agent claims from `ready`).

## Next

Done.

## Outcome

`claim` is refused with `reason: rule`, `rule: WK6` when `## Next` is empty (ADR 261101-09yd, accepted). Tests that claimed a Next-less item now write it first; conformance `claim-1` uses an item with a Next and `claim-3` pins the refusal. `go vet`, `go test ./...`, `gofmt`, `npm run schemas`, format and records checks pass. Decision: 261101-09yd.

## Log

- 2026-11-01: filed.
- 2026-11-01 17:17-07:00: claimed (claude).
- 2026-11-01 17:17-07:00: claimed itself without a Next: the bug reproduced on this item; wrote Next to clear WK6
- 2026-11-01 17:26-07:00: ticked "`claim` of an item whose `## Next` is empty is refused with `refused`, naming WK6 and saying to write `## Next` first (the item stays `ready`).": TestClaimAndRelease refuses a claim with an empty Next and accepts it after section Next; conformance case claim-3.
- 2026-11-01 17:26-07:00: ticked "The API spec, the `claim` reference in both skill copies and a fixture or test cover the refusal; a `claim` of an item with a `## Next` is unchanged.": docs/spec/akm-api-spec.md 4.13, claim.md and akm-track-work in both skill copies; fixture claim-3 and the claimsnonext repository in both copies.
- 2026-11-01 17:26-07:00: ticked "A `ready` item rule is checked too: say in the spec whether `ready` requires `## Next` (a cold agent claims from `ready`).": The spec now says a ready item needs no Next; only claim does (4.13).
