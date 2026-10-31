---
type: idea
status: idea
created_at: 2026-10-31
---

# Promote a task's decision

Status: `idea`.

## The idea

A builder meets a choice while working and notes it in the item's `## Log`. Later, a person or a supervisor sees that it reaches beyond the task. One command, such as `promote-decision <work-id> <n> --slug <slug>`, would copy that note into a new `proposed` ADR, add its id to the item's `decisions`, and regenerate the ADR index. Small decisions stay light; the ones that change the shape of the repository become records.

## What AKM already gives it

- `decisions` on work items, and rule WK11 that each entry resolves to an ADR.
- The rule that an agent drafts a `proposed` ADR and a maintainer accepts it.
- `new adr` and `index`.

## What it would still need

- A convention for marking a decision inside `## Log` or a section, so a command can find the nth one.
- A reason it beats `new adr` followed by `set decisions=...`. Today that is two commands and a copy; if that is enough, this stays an idea.

## First step

Count, in a month of real work, how often a decision starts in a log and is copied into an ADR by hand. If it is rare, set this aside.
