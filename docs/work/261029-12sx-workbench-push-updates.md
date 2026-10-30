---
id: "261029-12sx"
title: Push updates to the workbench with a watcher, keeping polling as the fallback
status: raw
created: 2026-10-29
updated: 2026-10-29
kind: feature
priority: 3
summary: Server-sent events fed by a file watcher, with the 2 second poll kept as the fallback, a manifest toggle and an akmctl ui --watch override
---

# Push updates to the workbench with a watcher, keeping polling as the fallback

## Goal

Alexander's wishes, for refinement and an ADR before any code:

- Push changes to the workbench with server-sent events fed by a file watcher (`fsnotify`, a dependency that needs an ADR).
- Keep the existing poll as a fallback: a manifest key picks the behaviour, and `akmctl ui --watch` overrides it.
- Open questions: how an `EventSource` (no `Authorization` header) is authenticated; debouncing a branch switch that rewrites many files; the inotify and macOS descriptor limits; the Windows path (on hold).
- Related: a conditional poll with a repository change token (`If-None-Match`, `304`) could be the same protocol the watcher feeds later. Not decided.

## Log

- 2026-10-29: filed.
