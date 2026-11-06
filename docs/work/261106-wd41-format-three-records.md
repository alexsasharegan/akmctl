---
id: "261106-wd41"
title: Format three records so CI format check passes
status: active
created: 2026-11-06
updated: 2026-11-06
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-06T18:22
---

# Format three records so CI format check passes

## Goal

`npm run format:check` fails on `main` for three files, so the `format records` step of the ubuntu CI job is red on every pull request: `docs/design/agent-environment-management.md`, `docs/work/done/261101-39k3-add-review-md.md` and `docs/work/done/261101-nsmy-revise-review-md.md`. Done when `npm run format:check` prints no warning and `akmctl check` still passes. Formatting only; no wording changes.

## Next

Open the pull request; the maintainer merges.

## Log

- 2026-11-06: filed.
- 2026-11-06 11:22-07:00: claimed (claude-sonnet-5-5).
- 2026-11-06 11:22-07:00: Ran prettier --write on the three files; only blank lines changed.
