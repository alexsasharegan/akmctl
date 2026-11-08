---
id: "261108-f1tv"
title: run fixtures through both bindings
status: done
kind: chore
priority: 3
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
decisions: [261108-15dw]
---

# run fixtures through both bindings

## Goal

Findings F1 and F2 of [the audit](../../design/test-strategy-audit.md), from [261030-456q](261030-456q-test-strategy-and-coverage-audit.md): the conformance fixtures run only through `api.Run`, so nothing proves the CLI and the HTTP server return the same envelope, and `new`, `commits`, `history` and `skills` have no fixture. Run every case through the HTTP handler and the reads through the CLI, and draft the ADR on fixing the clock, ids and git history for the rest.

Done when:

- [x] Every case file runs through the HTTP handler and compares with `expected`; the reads also run through the CLI with `--format json`.
- [x] An ADR draft (`proposed`) says how to fixture `new`, `commits`, `history` and `skills`.
- [x] The first fixtures that ADR allows are added, or filed as items.

## Next

## Outcome

The conformance fixtures now run through every binding: `internal/conformance` replaces the `internal/api` test and runs each of the 140 cases through the in-process API, the HTTP handler and, for reads, the CLI with `--format json`. Matchers, a replayed git history and a `bindings` member (ADR 261108-15dw) let `new`, `commits`, `history`, `skills` and the HTTP-only `if_match_required` refusal have fixtures; `-update` keeps matchers; the validator understands them. Two binding differences surfaced and are filed raw: 261108-jbd5 (`remove` has no HTTP route) and 261108-a91d (the `if_match_required` envelope names operation `http`). The CLI is run for reads only, as the ADR says.

## Log

- 2026-11-08: filed.
- 2026-11-08 14:26-07:00: claimed (claude-sonnet-5-5-rest).
- 2026-11-08 14:27-07:00: Drafted ADR 261108-15dw (proposed): matchers ($id, $hash, $any), a replayed git history in the case, one runner for api.Run, the HTTP handler and the CLI reads. Stopped here for the maintainer's decision before building the runners or the new fixtures.
- 2026-11-08 14:41-07:00: ticked "Every case file runs through the HTTP handler and compares with `expected`; the reads also run through the CLI with `--format json`.": internal/conformance runs every case through api, http and (reads) cli; 103 of 418 subtests skip by design (CLI writes, remove over HTTP), 315 pass
- 2026-11-08 14:41-07:00: ticked "The first fixtures that ADR allows are added, or filed as items.": 11 cases: new-2, new-3, commits-2, commits-3, history-1..3, skills-2..4, set-15
- 2026-11-08 14:41-07:00: Built internal/conformance: matchers, git replay, api/http/cli runners, matcher-keeping -update; removed internal/api/conformance_test.go. All 129 existing cases pass on every binding they apply to. Found: remove has no HTTP route (261108-jbd5); the if_match_required refusal names operation http with an empty request (261108-a91d). Validator now understands matchers; history.json had a strict-mode error (never compiled before a history case existed), fixed.
