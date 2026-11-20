---
id: "261119-3aeq"
title: "Workbench: yank a record's absolute path to the clipboard"
status: active
kind: feature
created: 2026-11-19
updated: 2026-11-19
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-20T02:32
---

# Workbench: yank a record's absolute path to the clipboard

## Goal

In the workbench (`akmctl ui`) `Y` copies the absolute path of the focused or open record. Done when the API exposes the repository root so the page can build the path (a small additive change with a changelog entry and a fixture), `Y` copies it with a toast, and the keys help lists it. Needs the id yank item first. GitHub issue 97; relative paths are left out on purpose.

## Next

Join the workbench config root and the record path, add Y beside y.

## Log

- 2026-11-19: filed.
- 2026-11-19 19:32-07:00: claimed (claude-sonnet-5-5).
