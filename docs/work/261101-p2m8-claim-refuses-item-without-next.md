---
id: "261101-p2m8"
title: claim refuses item without next
status: ready
created: 2026-11-01
updated: 2026-11-01
kind: bug
parent: 261101-vs5c
summary: claim leaves an item active with no Next, a state check then rejects (WK6)
---

# claim refuses item without next

## Goal

Seen in the first Codex trial: `akmctl claim 261101-3erc --as codex` succeeded and printed `active`, then `akmctl check` failed with `WK6 status active requires a non-empty ## Next`. The tool created a state its own checker forbids. Codex filed this as 261101-2n55 on the PR 1 branch (empty Goal); this item replaces it. Done when:

- [ ] `claim` of an item whose `## Next` is empty is refused with `refused`, naming WK6 and saying to write `## Next` first (the item stays `ready`).
- [ ] The API spec, the `claim` reference in both skill copies and a fixture or test cover the refusal; a `claim` of an item with a `## Next` is unchanged.
- [ ] A `ready` item rule is checked too: say in the spec whether `ready` requires `## Next` (a cold agent claims from `ready`).

## Log

- 2026-11-01: filed.
