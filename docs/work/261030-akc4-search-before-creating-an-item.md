---
id: "261030-akc4"
title: search before creating an item
status: active
created: 2026-10-30
updated: 2026-10-30
kind: feature
priority: 1
summary: Creating an item first looks for open items with similar titles, so the agent adds to one instead of filing a duplicate.
claimed_by: claude
claimed_at: 2026-10-31T03:41
---

# search before creating an item

## Goal

Before `akmctl new work` creates an item, the tool looks for open items (`raw`, `ready`, `active`, `blocked`) whose title or summary shares the new slug's words, and returns them as `similar` in the result and the CLI output. It never refuses. The `akm-work-items` skill's "No item fits" step says to read that list, and to add to a match instead of filing.

Why: the skill goes straight from "no item fits" to `new`, and an agent's idea of "fits" is only as good as the listing it read. The maintainer's vault skill ran a keyword search first and it caught real duplicates.

Done when:

- [ ] `new work` returns up to five similar open items; the matching rule (words, case, stop words) is written in the specification so it is testable.
- [ ] No match gives the same output as today.
- [ ] The API schema and fixtures cover a match and no match; the CLI prints the matches under the new record's path.
- [ ] `akm-work-items` ("No item fits") and the `akm` reference for `new` are updated in both copies (`.agents/skills/` and `internal/initrepo/templates/skills/`), and the skill drift test passes.
- [ ] Glossary and changelog are updated.

## Next

Read how `new work` is built (internal/cli and internal/api), then write the matching rule in the specification before the code.

## Log

- 2026-10-30: filed.
- 2026-10-30: filed from the vault migration review.
- 2026-10-30: claimed (claude).
- 2026-10-30: claimed.
