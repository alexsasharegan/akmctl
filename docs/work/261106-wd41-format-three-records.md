---
id: "261106-wd41"
title: Format three records so CI format check passes
status: raw
created: 2026-11-06
updated: 2026-11-06
---

# Format three records so CI format check passes

## Goal

`npm run format:check` fails on `main` for three files, so the `format records` step of the ubuntu CI job is red on every pull request: `docs/design/agent-environment-management.md`, `docs/work/done/261101-39k3-add-review-md.md` and `docs/work/done/261101-nsmy-revise-review-md.md`. Done when `npm run format:check` prints no warning and `akmctl check` still passes. Formatting only; no wording changes.

## Log

- 2026-11-06: filed.
