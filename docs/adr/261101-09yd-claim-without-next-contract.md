---
id: "261101-09yd"
title: claim refuses a ready item that has no Next
summary: "Proposed: claim is refused when the item's Next is empty, so claim cannot create a WK6 error; read before changing claim or WK6"
status: accepted
date: 2026-11-01
scope: [docs/spec/akm-api-spec.md, internal/api]
---

# claim without next contract

## Context

In the first Codex trial, `akmctl claim` moved a `ready` item to `active` although its `## Next` was empty, and `akmctl check` then failed with WK6. The same happened to the item that tracks this fix. The API spec leaves WK6 to `check` for `set` because an agent sets the status first and writes the section next. `claim` differs: the item is `ready` and `section` works on it, so the agent can write `## Next` before it claims.

## Considered Options

- Refuse the claim (`refused`, `details.rule: "WK6"`) when `## Next` is empty; the message says to write it first. The item stays `ready`.
- Accept an optional `next` parameter on `claim` and write it in the same change; refuse when neither is present.
- Make WK6 apply to `ready` as well, so the state never arises. This changes the standard, and every ready item in this repository and in adopters' repositories fails until filled.
- Do nothing: keep `check` as the only guard.

## Decision

Refuse the claim. It is the smallest change, it reads as the order the `akm-track-work` skill already gives (write `## Next`, then claim), it needs no new parameter, and it leaves a failing `check` impossible to reach through `claim`. A caller that wants one step can `section` then `claim`. Confirmed by Alexander on 2026-11-01.

## Consequences

- A cold agent that claims a ready item with no `## Next` gets a refusal that says what to do, not a state `check` rejects.
- `claim` gains a `refused` case with `details.rule` `WK6`; the API spec and both skill copies document it. Existing clients that claim only items with a `## Next` are unaffected; others get an error they did not get before, so the next MINOR notes it.
- `set status=active` still leaves WK6 to `check`; only `claim` changes.

## Confirmation

`TestClaimAndRelease` and conformance cases `claim-1` and `claim-3` (fixture repository `claimsnonext`).

## Review

Accepted by Alexander, 2026-11-01, in conversation; accepted for the merge of the branch that implements it.

## Revert

Remove the check in `opClaim` (`internal/api/ops_write.go`), the `WK6` sentence in section 4.13 of the API spec and the two skill copies, and conformance case `claim-3`.
