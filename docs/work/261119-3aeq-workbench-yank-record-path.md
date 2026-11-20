---
id: "261119-3aeq"
title: "Workbench: yank a record's absolute path to the clipboard"
status: raw
kind: feature
created: 2026-11-19
updated: 2026-11-19
---

# Workbench: yank a record's absolute path to the clipboard

## Goal

In the workbench (`akmctl ui`) `Y` copies the absolute path of the focused or open record. Done when the API exposes the repository root so the page can build the path (a small additive change with a changelog entry and a fixture), `Y` copies it with a toast, and the keys help lists it. Needs the id yank item first. GitHub issue 97; relative paths are left out on purpose.

## Log

- 2026-11-19: filed.
