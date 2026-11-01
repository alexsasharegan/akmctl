---
id: "261031-ph53"
title: new parent kind flags and split command
status: ready
created: 2026-10-31
updated: 2026-11-01
summary: CLI new takes --parent and --kind; split creates several children of an item
parent: 261101-t8h8
---

# new parent kind flags and split command

## Goal

From spike 261031-078w, amended by spike 261101-baze. The API `new` takes `fields` (parent, kind, ...) but the CLI `new work` takes none, so a split needs N `new` plus N `set` plus N `section` (the cold-agent trial of 261101-baze: three commands per child). Add `--parent` and `--kind` to `new`, and `akmctl split <id> <slug>...` creating one child per slug. Done when:

- [ ] `new work x --parent <id> --kind feature` sets both; an unknown parent is `not_found`.
- [ ] `split` creates the children as `raw` items under the parent and prints one `file` commit message per child, in order, as the `akm-commit` skill expects.
- [ ] Children take the parent's `priority` unless one is given (the trial found they inherit nothing, so they sort below unrelated work).
- [ ] After a split, `split` prints what is left for the parent, as `close` prints dependents: write a `## Next`, release the claim, set it `ready`, move its old Goal boxes. It edits none of that.
- [ ] The `Item too big` section of `akm-track-work` (step 2) is shortened to use `split`, in both skill copies, and "split" is registered in `docs/glossary.md` with the word "umbrella" for a parent whose Goal only lists its children.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
- 2026-11-01 10:22-07:00: Amended from spike 261101-baze (box 4): added the priority, per-child commit messages and parent-reminder boxes, and the skill update.
