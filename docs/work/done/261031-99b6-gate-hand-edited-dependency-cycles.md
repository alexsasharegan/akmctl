---
id: "261031-99b6"
title: gate hand edited dependency cycles
status: done
created: 2026-10-31
updated: 2026-10-31
closed: 2026-10-31
kind: chore
summary: Say how a hand-edited cycle is stopped before it lands
priority: 10
---

# gate hand edited dependency cycles

## Goal

From 261031-z50q. The tool refuses cycle-forming `set` and `new` calls, but a person or agent can still write a loop into a file by hand. The tool cannot prevent a file edit; it can make the gate easy. Done when:

- [ ] The guide shows how to run `akmctl check` as a git pre-commit hook and as a CI step, and says that the check reports WK4 and WK18 on a hand-edited loop (this repository's CI already runs it).
- [ ] The `akm-commit` skill and its repository copy say the same in one line; if a built-in hook installer is wanted instead, that is filed as a separate item.

## Next

Close this item.

## Outcome

- Guide section "Gate hand edits" (`docs/spec/akm-guide.md`): pre-commit hook and CI step; says `check` reports WK4 and WK18 on a hand-edited loop and that this repository's CI already runs it (`.github/workflows/ci.yml`).
- The `akm-commit` skill has the one-line pointer, in the repository copy and the template (`internal/initrepo/templates/optional-skills/akm-commit/SKILL.md`); `go test ./internal/initrepo` passes.
- A built-in hook installer was not built; it is not filed.

Decision: none.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from 261031-z50q; part of the patch release.
- 2026-10-31: claimed (claude).
- 2026-10-31: done.
