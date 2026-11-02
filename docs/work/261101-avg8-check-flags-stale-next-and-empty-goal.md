---
id: "261101-avg8"
title: check flags stale next and empty goal
status: ready
created: 2026-11-01
updated: 2026-11-01
kind: feature
parent: 261101-vs5c
summary: check flags a closed item whose Next is the pre-work plan, a raw item with an empty Goal and an unlabelled work-file commit
---

# check flags stale next and empty goal

## Goal

Three deviations in the first Codex trial are computable: 3erc closed with its pre-work `## Next`; 2n55 filed `raw` with an empty `## Goal`; a checkpoint commit (`261101-w0mz: documented replay history gap`) swept in a `release` change and, with no trailers, `commits` skipped it. Done when:

- [ ] Decide which of the three the checker can state as a rule without false positives (a closed item's `## Next` should be empty or say it is finished; a `raw` item needs a non-empty Goal before it is committed; a commit that changes `status` or `claimed_*` must be a step commit). Record each as a rule id, or say why not.
- [ ] Each chosen rule is implemented with tests, listed in the spec, and does not fire on this repository's existing records.

## Log

- 2026-11-01: filed.
