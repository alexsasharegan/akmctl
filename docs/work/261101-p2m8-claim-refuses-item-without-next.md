---
id: "261101-p2m8"
title: claim refuses item without next
status: active
created: 2026-11-01
updated: 2026-11-01
kind: bug
parent: 261101-vs5c
summary: claim leaves an item active with no Next, a state check then rejects (WK6)
claimed_by: claude
claimed_at: 2026-11-02T00:17
decisions: [261101-09yd]
---

# claim refuses item without next

## Goal

Seen in the first Codex trial: `akmctl claim 261101-3erc --as codex` succeeded and printed `active`, then `akmctl check` failed with `WK6 status active requires a non-empty ## Next`. The tool created a state its own checker forbids. Codex filed this as 261101-2n55 on the PR 1 branch (empty Goal); this item replaces it. Done when:

- [ ] `claim` of an item whose `## Next` is empty is refused with `refused`, naming WK6 and saying to write `## Next` first (the item stays `ready`).
- [ ] The API spec, the `claim` reference in both skill copies and a fixture or test cover the refusal; a `claim` of an item with a `## Next` is unchanged.
- [ ] A `ready` item rule is checked too: say in the spec whether `ready` requires `## Next` (a cold agent claims from `ready`).

## Next

Decide the contract in the proposed ADR (refuse the claim, or accept a next parameter, or let ready require Next); then implement in opClaim (internal/api/ops_write.go), add the test beside TestSetStatusChangeReturnsCommit style cases in internal/api, update docs/spec/akm-api-spec.md (claim section and the set paragraph that leaves WK6 to check) and the claim reference in both skill copies.

## Log

- 2026-11-01: filed.
- 2026-11-01 17:17-07:00: claimed (claude).
- 2026-11-01 17:17-07:00: claimed itself without a Next: the bug reproduced on this item; wrote Next to clear WK6
