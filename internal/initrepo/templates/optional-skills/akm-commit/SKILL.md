---
name: akm-commit
description: Commit after each logical step of AKM work. Use when finishing a logical step (a status change or a Next checkpoint) in a repository where git_commits is true.
---

# AKM commits

The repository's `AGENTS.md` wins over this skill where they disagree.

## First step: is it on, and is this a git repository?

1. Read `git_commits` from `docs/.layout`: `akmctl config git-commits`, or without the tool `grep '^git_commits: true' docs/.layout`. If it is not `true`, STOP: do not commit. The maintainer commits.
2. Check `git rev-parse --is-inside-work-tree` prints `true`. If it does not, STOP: do not commit and do not run `git init` yourself. Tell the maintainer this is not a git repository, so commits are off; offer to run `git init` and then `akmctl config git-commits on` if they want that, or `akmctl config git-commits off` to turn the option off and remove this skill.

## When to commit

One commit per logical step:

- a status change: claim, block, unblock, close;
- each `## Next` checkpoint (the work is in a state a cold agent could resume from).

`## Log` lines do not get a commit of their own: batch them into the next commit.

## How

1. Run `akmctl check`; fix errors first.
2. Stage ONLY by path: the work item, the roadmap, the ADR index and the code files this step changed. Never `git add -A` or `git add .`.
3. Leave unrelated dirty files alone: do not stage, stash, revert or commit them.
4. Message: for a status change, use the `commit.message` the command returned, as given: `<id>(<step>): <detail>`, a blank line, then the trailers `AKM-Step` and `AKM-Item`. Stage `commit.paths` and the code the step changed. You may put a body between the subject and the trailers and add other trailers (attribution) after them; do not change the subject or the two `AKM-` trailers. `new`, `claim`, `release`, `close` and a `set` that changes `status` return one (`--format json`, or the block the command prints). For a `## Next` checkpoint, which returns none, write `<id>: <what happened>`.
5. Follow the repository's `AGENTS.md` on attribution lines.

## Check

Before a push, `akmctl commits --range origin/main..HEAD` checks that each step commit's subject, trailers and diff agree (rules CM1 to CM3). Fix a failing commit with a new commit; do not rewrite history.

## Never

- Push, tag, amend, rebase or rewrite history.
- Use `--no-verify`. If a hook fails, fix the cause and make a new commit; the failed one did not happen.
