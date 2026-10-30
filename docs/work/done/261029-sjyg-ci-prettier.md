---
id: "261029-sjyg"
title: Prettier is not run by CI
summary: CI does not run prettier --check on docs and web files, so formatting drifts.
status: done
kind: chore
priority: 1
created: 2026-10-29
updated: 2026-10-29
---

# Prettier is not run by CI

## Goal

CI does not run prettier --check on docs and web files, so formatting drifts.

Done when:

- [ ] CI runs prettier --check, or the decision not to is recorded.

## Outcome

CI runs `npx prettier@3.9.10 --check AGENTS.md CHANGELOG.md README.md "docs/**/*.md"` on Linux (`.github/workflows/ci.yml`). It passes locally on the current tree after formatting three files that had drifted.

## Log

- 2026-10-29: filed.
