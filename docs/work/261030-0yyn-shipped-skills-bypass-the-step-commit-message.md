---
id: "261030-0yyn"
title: The shipped skills let an agent skip the step commit message
summary: A vault agent ignored the printed commit message because the skills defer to the optional akm-commit and tell it to commit work and item together.
status: ready
kind: bug
priority: 2
created: 2026-10-30
updated: 2026-10-30
---

# The shipped skills let an agent skip the step commit message

## Goal

An agent that has the three default skills (no `akm-commit`) uses the `commit` message the tool returns, one commit per step, without having to find another skill first.

Found by auditing the maintainer's vault (2026-10-31): its last five commits use `<id>: <text>`, no step word, no trailers, and bundle `file and claim`. The agent saw the printed message and wrote its own. Its account:

- `akm/references/new.md`, `claim.md` and `close.md` say the optional `akm-commit` skill "says how to use" the `commit`. Without that skill nothing says to use it as given.
- `akm-track-work/SKILL.md` says to "Commit the work, the item and any roadmap change together" and "in the same commit as the code", which reads as one commit per task, not one per step.
- The agent reported `references/commits.md` missing; the vault has it (checked from its file tree), so that was a misreading.
- `git_commits` was off, so `akm-commit` was not installed; the tool did not say what turning it on would give.

Fix: state the rule in the always-installed skills (use `commit.message` as given for the step, stage `commit.paths` with the code, one commit per step), reword the two "together" lines, and have `skills status` mention `akm-commit` when `git_commits` is off. Whether `akm-commit` should stop being optional is a decision (an ADR); ADR 261030-tbkh made it optional, and this item leaves it so.

Done when:

- [ ] The `akm` references that return a `commit` say to use `commit.message` as given and stage `commit.paths` with the code, with no need for `akm-commit`.
- [ ] `akm-track-work` no longer says to commit the whole task together; it says one commit per step, built from the returned message.
- [ ] `skills status` says that `akm-commit` is not installed while `git_commits` is off, and how to turn it on.
- [ ] Both skill copies are equal, the specification, schema and fixtures follow, and the changelog says so.

## Next

Wait for the maintainer to say which fix to take; then reword the skills in `internal/initrepo/templates/skills` and `.agents/skills` together (a drift test requires equal copies).

## Log

- 2026-10-30: filed.
