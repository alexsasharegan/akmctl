---
id: "261029-wpak"
title: "Rule S2: accept block scalars in SKILL.md frontmatter"
status: done
created: 2026-10-29
updated: 2026-11-01
closed: 2026-11-01
kind: chore
priority: 3
summary: Eight vault skills used YAML block scalars and failed S2 until they were folded by hand.
---

# Rule S2: accept block scalars in SKILL.md frontmatter

## Goal

Rule S2 rejects folded and literal block scalars in a skill's frontmatter. Accept valid block scalars with a skill-specific reader, or state the restriction in the spec.

Done when:

- [x] S2 accepts valid block scalars, or the spec says why not.
- [x] A fixture covers both.

## Next

## Outcome

- S2 accepts valid block scalars: `readSkillFrontmatter` in `internal/check/skillfm.go` swaps each top-level block scalar for a placeholder and returns its YAML value; rule S2 in `internal/check/check.go` checks that value. Proven by `TestS2BlockScalars`.
- Fixture covers both: `TestS2BlockScalars` has accepting cases (folded, literal, chomping indicators, comment after the header, block `name`) and rejecting ones (empty block description, block name that differs, anchor, duplicate key); `TestBlockValueMatchesYAML` compares the folded value with `go.yaml.in/yaml/v3`.

Decision: none. The record frontmatter reader stays strict; the reader is local to S2 and easy to change.

## Log

- 2026-10-29: filed.
- 2026-11-01 18:03-07:00: claimed (claude).
- 2026-11-01 18:05-07:00: ticked "S2 accepts valid block scalars, or the spec says why not.": internal/check/skillfm.go, rule S2 in check.go; spec S2 row and subset table
- 2026-11-01 18:05-07:00: ticked "A fixture covers both.": TestS2BlockScalars (accept and reject cases) and TestBlockValueMatchesYAML in internal/check/skillfm_test.go
- 2026-11-01 18:05-07:00: S2 reads block scalars in SKILL.md frontmatter through readSkillFrontmatter (internal/check/skillfm.go); other lines still go through frontmatter.Parse. Tests, spec, glossary and changelog updated.
