---
id: "261028-jzsx"
title: Cover the Windows path beyond CI
summary: Windows behaviour (paths, line endings, locks, atomic rename) is tested only by the CI runner; add a way to run or simulate it locally.
status: ready
kind: chore
priority: 1
created: 2026-10-28
updated: 2026-10-28
---

# Cover the Windows path beyond CI

## Goal

The Windows behaviour (path separators, CRLF, file locks, rename over an existing file) is exercised only by the Windows job in CI. A change cannot be tried on Windows before it is pushed.

Done when:

- [ ] A documented way to run the Windows-sensitive tests locally, or tests that simulate the differences (path separators and line endings) on any platform.
- [ ] The CI Windows job stays the final check.
- [ ] A note in `docs/architecture.md` lists what is covered where.

## Log

- 2026-10-28: filed.
