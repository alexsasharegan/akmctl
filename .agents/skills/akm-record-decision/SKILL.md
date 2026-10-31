---
name: akm-record-decision
description: Draft an AKM decision record (ADR) in docs/adr/ for a hard-to-reverse choice or one that rejects a real alternative, and supersede one. Use before making such a choice.
---

# Record a decision

The repository's `AGENTS.md` wins over this skill where they disagree.

## When

A choice that is hard to reverse, or that rejects a real alternative, needs an ADR. A reversible detail does not. If unsure, ask.

Draft before making the choice. If you notice at close that the work already made one, draft the ADR then (`akm-track-work` Close asks), and list it in the item's `decisions`.

## Draft

1. Look first: read `docs/adr/README.md` (the generated index) and the records whose `scope` matches your area. Do not re-decide a settled question.
2. `akmctl new adr <slug>` creates a `proposed` record. Without the tool, copy `assets/adr.md` and mint the id with the recipe in `AGENTS.md`.
3. Fill `summary` (one sentence, up to 160 characters: what was decided and when to read it), `scope`, and the sections `Context`, `Considered Options` (at least two, including doing nothing), `Decision`, `Consequences`. Add `Confirmation` (a test or grep that fails if the decision is violated) when one exists, and `Review` and `Revert` for a costly or hard-to-reverse choice.
4. Run `akmctl index`, then `akmctl check`.

## Accept

An agent drafts as `proposed` and does not accept. Acceptance is the maintainer's commit that sets `status: accepted`, or the merge of the pull request that carries it. Remove `authored_by: agent` when the maintainer has reviewed the record.

## Change of mind

Never edit the body of an `accepted` record except status, `superseded_by`, an errata section and typos. Write a new ADR with `supersedes: <old id>`, then set the old one to `superseded` with `superseded_by: <new id>`, and run `akmctl index`.
