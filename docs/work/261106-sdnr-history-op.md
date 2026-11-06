---
id: "261106-sdnr"
title: "history: a read-only op that lists AKM moves from git"
status: active
created: 2026-11-06
updated: 2026-11-06
kind: feature
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-06T19:24
---

# history: a read-only op that lists AKM moves from git

## Goal

Add a read-only `history` operation (ADR 261106-60xv): it lists the AKM step commits of a range as moves, so the workbench and agents can see what happened, in order, by whom. Today `commits` only checks them. Done when:

- [x] `history` returns, oldest first, `{sha, time, author, item, step, from, to, trailers}` for each step commit in a range (default: since a given commit, or the last N), where `from` and `to` are the item's status before and after, read from the commit's diff.
- [x] A repository that is not a git work tree, or has `git_commits` off, gets a clear `available: false` with a `reason`, not an error. A range with no step commits is an empty list (a poll with nothing new), not an error.
- [x] The API spec has the section, `docs/spec/akm-api/` the schema and fixtures (both copies, 261029-x4ws), `akmctl history` the CLI binding, the glossary the term, the CHANGELOG an entry.
- [x] Tests cover order, a path through several steps, a commit that changes two items, a repo without git, and a range with no step commits.

## Next

Open the pull request; the maintainer merges. Then unblock w0mz.

## Log

- 2026-11-06: filed.
- 2026-11-06 12:24-07:00: claimed (claude-sonnet-5-5).
- 2026-11-06 12:30-07:00: Built history: spec 4.24, schema, op, CLI table, HTTP route, skill reference (both copies), 12 init fixtures in both copies, glossary and CHANGELOG. 'No step commits' is an empty list rather than available false, since a poll with nothing new is normal; the box says so.
