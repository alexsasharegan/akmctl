---
id: "261106-sdnr"
title: "history: a read-only op that lists AKM moves from git"
status: ready
created: 2026-11-06
updated: 2026-11-06
kind: feature
---

# history: a read-only op that lists AKM moves from git

## Goal

Add a read-only `history` operation (ADR 261106-60xv): it lists the AKM step commits of a range as moves, so the workbench and agents can see what happened, in order, by whom. Today `commits` only checks them. Done when:

- [ ] `history` returns, oldest first, `{sha, time, author, item, step, from, to, trailers}` for each step commit in a range (default: since a given commit, or the last N), where `from` and `to` are the item's status before and after, read from the commit's diff.
- [ ] A repository that is not a git work tree, has `git_commits` off or has no step commits gets a clear `available: false` with a `reason`, not an error.
- [ ] The API spec has the section, `docs/spec/akm-api/` the schema and fixtures (both copies, 261029-x4ws), `akmctl history` the CLI binding, the glossary the term, the CHANGELOG an entry.
- [ ] Tests cover order, a path through several steps, a commit that changes two items, a repo without git, and a range with no step commits.

## Next

Read how commits.go lists step commits and diffs them; design the response and the since cursor; write the spec section first, then fixtures, then the op and the CLI.

## Log

- 2026-11-06: filed.
