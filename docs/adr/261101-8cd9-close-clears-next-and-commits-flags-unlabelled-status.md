---
id: "261101-8cd9"
title: close empties Next and commits flags a status change with no step
summary: close clears a non-empty Next; commits reports an unlabelled status or claim change (CM4); no check rule for a stale Next or an empty raw Goal
status: proposed
date: 2026-11-01
---

# close clears next and commits flags unlabelled status

## Context

The first Codex trial left a closed item whose `## Next` was still the pre-work plan, filed a `raw` item with an empty `## Goal`, and committed a `release` inside an unlabelled checkpoint commit that `commits` skipped (ADR 261030-tbkh). The record was fixed by prose rules each time; these three can be tested against the repository first. 59 of 98 closed items hold a `## Next`, so a `check` rule on them would warn on legacy files and cannot say what is stale. 23 items are `raw`, and `raw` is allowed to be empty.

## Considered Options

- A `check` warning for a closed item with a non-empty `## Next`.
- `close` empties `## Next`; no rule needed.
- A `check` or `new` rule that a raw item needs a `## Goal`.
- `commits` reports a commit with no step that changes status or claim (CM4).
- Do nothing and rely on the skill text.

## Decision

`close` empties a non-empty `## Next` (the heading stays; an item with none is not given one), because a closed item has no next step and the tool, not a reminder, should ensure it. `commits` adds CM4 for a commit with no step claim that changes the `status` or `claimed_by` of an open work item in place. No rule for a `raw` item with an empty `## Goal`: raw is a capture state and is allowed to be empty, so the fix is wording in `AGENTS.md` (work item 261101-7ghb). No `check` rule for a stale `## Next`: it would fire on 59 legacy items and staleness is not computable.

## Consequences

- A closed item loses its plan text. Git history keeps it; the Outcome is the record. A client that read `## Next` of a just-closed item gets it empty.
- `commits` can now report a commit it used to skip. A run over old history (no `--range`) may report earlier status edits made by hand; `--range origin/main..HEAD` reports only the branch's own.
- Fixtures `close-2` to `close-6` and `close-8` change in both copies; the generator clears Next in each.

## Confirmation

Conformance cases `close-2` to `close-6`, `close-8`; `TestCommitsFlagsAnUnlabelledStatusChange`.

## Review

Proposed 2026-11-01 under the agent-rails umbrella 261101-vs5c. Accepted by merging the branch that implements it, if Alexander agrees.

## Revert

Remove the Next clearing in `opClose` (`internal/api/ops_write.go`), CM4 and `unlabelledChange` in `ops_commits.go`, regenerate the fixtures, and revert the spec and skill sentences.
