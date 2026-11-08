---
id: "261031-ph53"
title: new parent kind flags and split command
status: done
created: 2026-10-31
updated: 2026-11-07
closed: 2026-11-07
summary: CLI new takes --parent and --kind; split creates several children of an item
parent: 261101-t8h8
---

# new parent kind flags and split command

## Goal

From spike 261031-078w, amended by spike 261101-baze. The API `new` takes `fields` (parent, kind, ...) but the CLI `new work` takes none, so a split needs N `new` plus N `set` plus N `section` (the cold-agent trial of 261101-baze: three commands per child). Add `--parent` and `--kind` to `new`, and `akmctl split <id> <slug>...` creating one child per slug. Done when:

- [x] `new work x --parent <id> --kind feature` sets both; an unknown parent is `not_found`.
- [x] `split` creates the children as `raw` items under the parent and prints one `file` commit message per child, in order, as the `akm-commit` skill expects.
- [x] Children take the parent's `priority` unless one is given (the trial found they inherit nothing, so they sort below unrelated work).
- [x] After a split, `split` prints what is left for the parent, as `close` prints dependents: write a `## Next`, release the claim, set it `ready`, move its old Goal boxes. It edits none of that.
- [x] The `Item too big` section of `akm-track-work` (step 2) is shortened to use `split`, in both skill copies, and "split" is registered in `docs/glossary.md` with the word "umbrella" for a parent whose Goal only lists its children.

## Next

## Outcome

new takes --parent, --kind and --priority; akmctl split <id> <slug>... files raw children with the parent's priority, prints a file commit message each and what the parent still needs. CLI only, not an API operation. The skill's Item too big step 2 and the glossary are updated.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
- 2026-11-01 10:22-07:00: Amended from spike 261101-baze (box 4): added the priority, per-child commit messages and parent-reminder boxes, and the skill update.
- 2026-11-07 23:52-07:00: claimed (claude-sonnet-5-5).
- 2026-11-07 23:52-07:00: ticked "`new work x --parent <id> --kind feature` sets both; an unknown parent is `not_found`.": TestNewTakesParentKindAndPriorityFlags: --parent/--kind/--priority set the fields; an unknown parent is not_found
- 2026-11-07 23:52-07:00: Commit 1: new takes --parent, --kind and --priority (the API already took them as fields). --kind is also the record type's parameter name, so the CLI reads it as a field for new and split.
- 2026-11-07 23:54-07:00: ticked "`split` creates the children as `raw` items under the parent and prints one `file` commit message per child, in order, as the `akm-commit` skill expects.": TestSplitFilesChildrenUnderTheParent: one file commit message per child, in order
- 2026-11-07 23:54-07:00: ticked "Children take the parent's `priority` unless one is given (the trial found they inherit nothing, so they sort below unrelated work).": children take the parent's priority; --priority wins (same test)
- 2026-11-07 23:54-07:00: ticked "After a split, `split` prints what is left for the parent, as `close` prints dependents: write a `## Next`, release the claim, set it `ready`, move its old Goal boxes. It edits none of that.": split prints the parent's remaining steps and edits nothing (same test)
- 2026-11-07 23:54-07:00: ticked "The `Item too big` section of `akm-track-work` (step 2) is shortened to use `split`, in both skill copies, and "split" is registered in `docs/glossary.md` with the word "umbrella" for a parent whose Goal only lists its children.": step 2 of Item too big uses split in both skill copies (identical); glossary term Split with umbrella
- 2026-11-07 23:54-07:00: Commit 2: akmctl split, CLI only (one new per slug, validated first so a bad slug creates nothing). Skill step 2 shortened in both copies; glossary Split registers umbrella, Parent row no longer strikes it.
