---
id: "261101-nsmy"
title: Revise REVIEW.md after the critical review
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
---

# Revise REVIEW.md after the critical review

## Goal

Revise the root `REVIEW.md` after the maintainer's critical review of 2026-11-01, which closed work item 261101-39k3 with the first version. Done when:

- [x] `REVIEW.md` carries the ten agreed revisions: check-only rules, behavior and process labels, function-scoped debt, scanner scope, ADR acceptance recorded in the work item Log, wider Records examples, the write-path wording, no Tests section, named oracle tests, and the closed work item invariant.
- [x] `npx prettier --check REVIEW.md` and `akmctl check` pass.
- [x] `AGENTS.md`, `CLAUDE.md` and the closed item 261101-39k3 are unchanged.

## Next

## Outcome

REVIEW.md revised with the ten agreed changes: check-only rule distinction, behavior and process labels, function-scoped debt, scanner scope, ADR acceptance recorded in the work item Log, wider Records examples, the write-path wording, no Tests section, named oracle tests and the closed work item invariant. Checks: prettier --check REVIEW.md and akmctl check pass. Open: the migrate fence loops (internal/migrate/steps.go, migrate.go) are known debt with no work item yet; the branch is stacked on claude/review-md, which carries an unrelated architecture.md commit. No changelog entry: no user-visible tool change. ADR check: none.

## Log

- 2026-11-01: filed.
- 2026-11-01 23:19-07:00: claimed (claude-nsmy).
- 2026-11-01 23:19-07:00: Revised REVIEW.md (46 lines). ADR acceptance needs maintainer authorization recorded in this Log; none is recorded here because no ADR changes. prettier --check and akmctl check pass; AGENTS.md, CLAUDE.md and 261101-39k3 unchanged. Stacked on claude/review-md because main has no REVIEW.md.
