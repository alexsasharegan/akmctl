---
id: "261030-dw71"
title: close reports dependents
status: done
created: 2026-10-30
updated: 2026-10-30
closed: 2026-10-30
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

- [x] `close` output and the API result list the open dependents (id, status, title); an empty list is omitted or empty, and the schema and fixtures cover both.
- [x] A close with no dependents is unchanged.
- [x] The stale-blocker `check` question is decided and recorded in the Log, and the rule is added if the answer is yes.
- [x] The `akm-work-items` Close step tells the agent to read the list and unblock or re-point each item, and the `akm` reference for `close` shows an example.
- [x] Specification, guide, glossary and changelog are updated.

## Outcome

- `close` returns `dependents` (`dependents()` in `internal/api/ops_write.go`): open work items whose `blocked_by` names the closed item, for `done`, `dropped` and `dry_run`; omitted when empty, so a close with none is unchanged. It edits none of them and never refuses. Tests: `internal/api/dependents_test.go` (a dependent, none, a closed item is not one) and `TestClosePrintsDependents` in `internal/cli`.
- Schema `close.json` has the member. Fixtures: `make_fixtures.py` adds it to close cases 2, 3, 4, 5 and 6, regenerated, with `testdata/api/cases` synced by hand; `npm run schemas` passes (19 schemas, 67 cases) and the conformance run is 67 of 67.
- Stale-blocker question, decided: no `check` rule now. A new rule is a layout change (WK13 to WK17 shipped as layout 0.3.0), so it would force a layout bump, a manifest version and fixture churn for a case that arises only when an item is closed by hand. `close` now reports the dependents on the normal path. Revisit with the next layout bump if hand-closed blockers show up.
- The `akm-work-items` Close step 4 and `references/close.md` (both copies) tell the agent to read the list and unblock or re-point each item; spec section 4.12, glossary ("Dependents") and changelog are updated.
- Checks: `go vet`, `go test ./...`, `gofmt -l`, `akmctl check`, `npm run format:check`, `npm run schemas`.

Decision: none. The stale-blocker rule is deferred, not rejected, and is reversible.

## Log

- 2026-10-30: filed.
- 2026-10-30: filed from the vault migration review (task-lifecycle step 4.3 compared with akm-work-items).
- 2026-10-30: claimed (claude).
- 2026-10-30: implemented; stale-blocker check decided (not now); closed.
