---
id: "261031-ph53"
title: new parent kind flags and split command
status: raw
created: 2026-10-31
updated: 2026-10-31
summary: CLI new takes --parent and --kind; split creates several children of an item
---

# new parent kind flags and split command

## Goal

From spike 261031-078w. The API `new` takes `fields` (parent, kind, ...) but the CLI `new work` takes none, so a project split needs N `new` plus N `set`. Add `--parent` and `--kind` to `new`, and `akmctl split <id> <slug>...` creating one child per slug. Done when:

- [ ] `new work x --parent <id> --kind feature` sets both; an unknown parent is `not_found`.
- [ ] `split` creates the children as `raw` items and prints one commit message per file or one for the batch, as the `akm-commit` skill expects.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
