---
id: "261029-nf6h"
title: "mint_ids: build the slug from the title"
status: active
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

- [x] A slug_from option (name or title) and a maximum length.
- [x] Tests for numeric names and long titles.
- [x] adr_slugs.py in the vault scripts is no longer needed.

## Next

Review slug_from and max_slug (internal/migrate/mintids.go). The step list in the API spec and schemas/common.json do not name the new options yet (spec edits were out of scope); adr_slugs.py is kept as migration history. Then close with the drafted Outcome.

## Log

- 2026-10-29: filed.
- 2026-10-29: claimed (life-agent).
- 2026-10-29: mint_ids takes slug_from (name or title) and max_slug (1 to 100, default 60; a cut ends on a whole word unless the default name slug). Tests in internal/migrate/slugfrom_test.go. adr_slugs.py is no longer needed for a new run and stays as history.
- 2026-10-29: released (life-agent).
