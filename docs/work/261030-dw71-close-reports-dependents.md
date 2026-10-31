---
id: "261030-dw71"
title: close reports dependents
status: ready
created: 2026-10-30
updated: 2026-10-30
kind: feature
priority: 2
summary: close and drop report the open items that name the closed item in blocked_by, so a blocker is not left pointing at finished work.
---

# close reports dependents

## Goal

When an item is closed as `done` or `dropped`, `close` lists the open items whose `blocked_by` names it, in the CLI output and as a field in the API result. It reports; it never edits those items and does not refuse. A `parent` with open children is already handled (WK8), so only `blocked_by` is new.

Why: an agent that closes an item has no cue that something else was waiting on it, so the waiting item stays `blocked` until someone notices. The maintainer's vault skill did this by hand (a search for open items linking to the closed one); the structured part of that belongs in the tool, the free-text link scan does not.

Open question to settle in this item: whether `check` should also warn on a `blocked` item whose `blocked_by` entries are all `done` or `dropped` (a stale blocker). If yes, it needs its own rule id and a spec row.

Done when:

- [ ] `close` output and the API result list the open dependents (id, status, title); an empty list is omitted or empty, and the schema and fixtures cover both.
- [ ] A close with no dependents is unchanged.
- [ ] The stale-blocker `check` question is decided and recorded in the Log, and the rule is added if the answer is yes.
- [ ] The `akm-work-items` Close step tells the agent to read the list and unblock or re-point each item, and the `akm` reference for `close` shows an example.
- [ ] Specification, guide, glossary and changelog are updated.

## Log

- 2026-10-30: filed.
- 2026-10-30: filed from the vault migration review (task-lifecycle step 4.3 compared with akm-work-items).
