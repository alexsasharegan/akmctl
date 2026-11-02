---
id: "261101-avg8"
title: check flags stale next and empty goal
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
kind: feature
parent: 261101-vs5c
summary: check flags a closed item whose Next is the pre-work plan, a raw item with an empty Goal and an unlabelled work-file commit
decisions: [261101-8cd9]
---

# check flags stale next and empty goal

## Goal

Three deviations in the first Codex trial are computable: 3erc closed with its pre-work `## Next`; 2n55 filed `raw` with an empty `## Goal`; a checkpoint commit (`261101-w0mz: documented replay history gap`) swept in a `release` change and, with no trailers, `commits` skipped it. Done when:

- [x] Decide which of the three the checker can state as a rule without false positives (a closed item's `## Next` should be empty or say it is finished; a `raw` item needs a non-empty Goal before it is committed; a commit that changes `status` or `claimed_*` must be a step commit). Record each as a rule id, or say why not.
- [x] Each chosen rule is implemented with tests, listed in the spec, and does not fire on this repository's existing records.

## Next

Done.

## Outcome

Two of the three deviations moved into the tool: `close` empties `## Next`, and `commits` reports CM4 for an unlabelled status or claim change (ADR 261101-8cd9, proposed). The empty `raw` Goal is not a rule, since raw may be empty; it is a box on 261101-7ghb. Not covered: a stale `## Next` on an open item, and a step commit that carries unrelated work. Decision: 261101-8cd9.

## Log

- 2026-11-01: filed.
- 2026-11-01 17:47-07:00: claimed (claude).
- 2026-11-01 17:48-07:00: ticked "Decide which of the three the checker can state as a rule without false positives (a closed item's `## Next` should be empty or say it is finished; a `raw` item needs a non-empty Goal before it is committed; a commit that changes `status` or `claimed_*` must be a step commit). Record each as a rule id, or say why not.": ADR 261101-8cd9 decides all three: close empties Next, CM4 added, no rule for an empty raw Goal or a stale Next (59 of 98 closed items would warn).
- 2026-11-01 17:48-07:00: ticked "Each chosen rule is implemented with tests, listed in the spec, and does not fire on this repository's existing records.": opClose, CM4 in ops_commits.go, conformance close-2..6 and close-8, TestCommitsFlagsAnUnlabelledStatusChange, spec section 4.12 and 4.19; no existing record fails check.
