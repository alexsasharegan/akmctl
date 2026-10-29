---
id: "261029-nf6h"
title: "mint_ids: build the slug from the title"
status: ready
created: 2026-10-29
updated: 2026-10-29
kind: feature
priority: 2
summary: For notes whose file name is a number, mint_ids made adr-019 as the slug; offer a title-based slug with a maximum length.
---

# mint_ids: build the slug from the title

## Goal

The vault migration overrode mint_ids slugs with a script because the old names were numbers. Add a mint_ids option that builds the slug from the title with a length limit.

Done when:

- [ ] A slug_from option (name or title) and a maximum length.
- [ ] Tests for numeric names and long titles.
- [ ] adr_slugs.py in the vault scripts is no longer needed.

## Log

- 2026-10-29: filed.
