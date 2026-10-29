---
id: "261029-xtqj"
title: Multi-write moves can stop half-way with no rollback
summary: A wizard move makes several writes in turn; if a later one fails the earlier ones stay and nothing is rolled back.
status: ready
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-10-29
---

# Multi-write moves can stop half-way with no rollback

## Goal

A wizard move makes several writes in turn; if a later one fails the earlier ones stay and nothing is rolled back.

Done when:

- [ ] Decide how a move that fails half-way is undone or reported, and implement it.

## Log

- 2026-10-29: filed.
