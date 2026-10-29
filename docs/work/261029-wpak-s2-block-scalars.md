---
id: "261029-wpak"
title: "Rule S2: accept block scalars in SKILL.md frontmatter"
status: ready
created: 2026-10-29
updated: 2026-10-29
kind: chore
priority: 3
summary: Eight vault skills used YAML block scalars and failed S2 until they were folded by hand.
---

# Rule S2: accept block scalars in SKILL.md frontmatter

## Goal

Rule S2 rejects folded and literal block scalars in a skill's frontmatter. Accept valid block scalars with a skill-specific reader, or state the restriction in the spec.

Done when:

- [ ] S2 accepts valid block scalars, or the spec says why not.
- [ ] A fixture covers both.

## Log

- 2026-10-29: filed.
