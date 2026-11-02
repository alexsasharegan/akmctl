---
id: "261029-wpak"
title: "Rule S2: accept block scalars in SKILL.md frontmatter"
status: active
created: 2026-10-29
updated: 2026-11-01
kind: chore
priority: 3
summary: Eight vault skills used YAML block scalars and failed S2 until they were folded by hand.
claimed_by: claude
claimed_at: 2026-11-02T01:03
---

# Rule S2: accept block scalars in SKILL.md frontmatter

## Goal

Rule S2 rejects folded and literal block scalars in a skill's frontmatter. Accept valid block scalars with a skill-specific reader, or state the restriction in the spec.

Done when:

- [ ] S2 accepts valid block scalars, or the spec says why not.
- [ ] A fixture covers both.

## Next

Add a skill-specific frontmatter reader in internal/check that folds block scalars to a plain string before frontmatter.Parse, use it in rule S2, add tests, update spec table row and CHANGELOG.

## Log

- 2026-10-29: filed.
- 2026-11-01 18:03-07:00: claimed (claude).
