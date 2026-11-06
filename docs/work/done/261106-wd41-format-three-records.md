---
id: "261106-wd41"
title: Format three records so CI format check passes
status: done
created: 2026-11-06
updated: 2026-11-06
closed: 2026-11-06
---

# Format three records so CI format check passes

## Goal

`npm run format:check` fails on `main` for three files, so the `format records` step of the ubuntu CI job is red on every pull request: `docs/design/agent-environment-management.md`, `docs/work/done/261101-39k3-add-review-md.md` and `docs/work/done/261101-nsmy-revise-review-md.md`. Done when `npm run format:check` prints no warning and `akmctl check` still passes. Formatting only; no wording changes.

## Next

## Outcome

Prettier --write on the three files; only blank lines changed (git diff --ignore-blank-lines is empty). npm run format:check is clean, so the ubuntu format records step no longer fails on main's files.

## Log

- 2026-11-06: filed.
- 2026-11-06 11:22-07:00: claimed (claude-sonnet-5-5).
- 2026-11-06 11:22-07:00: Ran prettier --write on the three files; only blank lines changed.
