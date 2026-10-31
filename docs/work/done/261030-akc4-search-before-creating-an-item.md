---
id: "261030-akc4"
title: search before creating an item
status: done
created: 2026-10-30
updated: 2026-10-30
closed: 2026-10-30
kind: feature
priority: 1
summary: Creating an item first looks for open items with similar titles, so the agent adds to one instead of filing a duplicate.
---

# search before creating an item

## Goal

Before `akmctl new work` creates an item, the tool looks for open items (`raw`, `ready`, `active`, `blocked`) whose title or summary shares the new slug's words, and returns them as `similar` in the result and the CLI output. It never refuses. The `akm-work-items` skill's "No item fits" step says to read that list, and to add to a match instead of filing.

Why: the skill goes straight from "no item fits" to `new`, and an agent's idea of "fits" is only as good as the listing it read. The maintainer's vault skill ran a keyword search first and it caught real duplicates.

Done when:

- [x] `new work` returns up to five similar open items; the matching rule (words, case, stop words) is written in the specification so it is testable.
- [x] No match gives the same output as today.
- [x] The API schema and fixtures cover a match and no match; the CLI prints the matches under the new record's path.
- [x] `akm-work-items` ("No item fits") and the `akm` reference for `new` are updated in both copies (`.agents/skills/` and `internal/initrepo/templates/skills/`), and the skill drift test passes.
- [x] Glossary and changelog are updated.

## Outcome

- `new work` returns up to five `similar` open items: `internal/api/similar.go`, tests in `internal/api/similar_test.go` (a match with its shared words, no match, one shared word, a one-word slug, closed and ADR never match, the cap and the order) and `TestNewPrintsSimilarItems` in `internal/cli`. The rule is in section 4.10 of `docs/spec/akm-api-spec.md`.
- No match: the response has no `similar`, as before.
- Schema: `docs/spec/akm-api/schemas/new.json` has the member; `npm run schemas` passes (19 schemas, 67 cases). No fixture case was added: the API spec says `new` is tested for shape only because ids are random, so the match and no-match cases are Go tests rather than fixtures.
- The CLI prints the matches under the path. `akm-work-items` ("No item fits") and `references/new.md` are changed in both copies; `TestSkillsMatchRepository` passes.
- Glossary has "Similar items"; the changelog has an Unreleased entry.
- Checks: `go vet`, `go test ./...`, `gofmt -l`, `akmctl check` (0 errors, 0 warnings) and `npm run format:check`.

Decision: none. The threshold (two shared words) and the stop list are reversible and stated in the spec.

## Log

- 2026-10-30: filed.
- 2026-10-30: filed from the vault migration review.
- 2026-10-30: claimed (claude).
- 2026-10-30: claimed.
- 2026-10-30: implemented and closed.
