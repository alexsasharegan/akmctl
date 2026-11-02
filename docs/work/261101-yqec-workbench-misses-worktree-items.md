---
id: "261101-yqec"
title: workbench misses worktree items
status: raw
created: 2026-11-01
updated: 2026-11-01
parent: 261101-vs5c
---

# workbench misses worktree items

## Goal

The maintainer could not open 261101-qr8x in the workbench UI while it existed only in a git worktree on an unmerged branch. The workbench reads one checkout, so work filed or claimed on another worktree or branch is invisible to it. Done when:

- [ ] It is decided whether the workbench should show items from other worktrees (a repository picker, a list of linked worktrees, or a notice that others exist), or stay one-checkout and say so.
- [ ] The chosen behaviour is built and tested, or the limit is written in the guide and glossary.

## Log

- 2026-11-01: filed.
