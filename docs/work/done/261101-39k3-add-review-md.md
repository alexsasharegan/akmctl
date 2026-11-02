---
id: "261101-39k3"
title: Add a repository REVIEW.md review policy
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
---

# Add a repository REVIEW.md review policy

## Goal

Add a root `REVIEW.md` that tells an AI reviewer what to flag in this repository. The text is the one the maintainer approved on 2026-11-01 after the collaborative review-policy exercise.

Done when:

- [x] `REVIEW.md` exists at the repository root with the approved text.
- [x] `npx prettier --check REVIEW.md`, the records check and the Go checks pass.
- [x] `AGENTS.md` and `CLAUDE.md` are unchanged.

## Next


## Outcome

REVIEW.md added at the repository root with the text the maintainer approved. It holds 8 invariants, the Records rule and the output format. Checked: go vet, go test -race, gofmt, format:check, schemas and akmctl check (0 errors). AGENTS.md and CLAUDE.md are unchanged. Open: AGENTS.md forbids any edit to an accepted ADR body while REVIEW.md allows a scrub of private material; aligning AGENTS.md needs separate approval. No changelog entry: no user-visible tool change. ADR check: none.

## Log

- 2026-11-01: filed.
- 2026-11-01 23:03-07:00: claimed (claude-39k3).
- 2026-11-01 23:03-07:00: Wrote REVIEW.md from the approved draft (47 lines). prettier --check, akmctl check (0 errors) pass. AGENTS.md says an accepted ADR body is never edited; REVIEW.md adds a scrub-of-private-material exception, so AGENTS.md needs a separate approved edit to match.
