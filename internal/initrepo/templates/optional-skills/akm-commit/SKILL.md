---
name: akm-commit
description: Commit after each logical step of AKM work. Use when finishing a logical step (a status change or a Next checkpoint) in a repository where git_commits is true.
---

# AKM commits

The repository's `AGENTS.md` wins over this skill where they disagree.

## First step: is it on?

Read `git_commits` from `docs/.layout`: `akmctl config git-commits`, or without the tool `grep '^git_commits: true' docs/.layout`. If it is not `true`, STOP: do not commit. The maintainer commits.

## When to commit

One commit per logical step:

- a status change: claim, block, unblock, close;
- each `## Next` checkpoint (the work is in a state a cold agent could resume from).

`## Log` lines do not get a commit of their own: batch them into the next commit.

## How

1. Run `akmctl check`; fix errors first.
2. Stage ONLY by path: the work item, the roadmap, the ADR index and the code files this step changed. Never `git add -A` or `git add .`.
3. Leave unrelated dirty files alone: do not stage, stash, revert or commit them.
4. Message: `<id>: <what happened>`, for example `261029-chq2: claim`, `261029-chq2: block on ADR review`, `261029-chq2: close`.
5. Follow the repository's `AGENTS.md` on attribution lines.

## Never

- Push, tag, amend, rebase or rewrite history.
- Use `--no-verify`. If a hook fails, fix the cause and make a new commit; the failed one did not happen.
