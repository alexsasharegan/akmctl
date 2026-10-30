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

Options for a native macOS driver (FSEvents has no pure-Go interface; it is a CoreServices C API), to weigh in the ADR:

- cgo, building the macOS archives on the macOS CI runner: works, but changes the release pipeline and breaks the standard-library-only rule.
- `purego` (calls system libraries without cgo, so Linux-built macOS archives still work): a small dependency; the FSEvents callback and dispatch-queue setup is fiddly and needs a spike first.
- A small helper binary that prints change events, spawned by `akmctl ui --watch=native`: keeps cgo out of the main binary, adds a second artifact and another executable to run.
- `fsnotify` uses kqueue on macOS (as far as we know), which needs a descriptor per watched file; check the limits on a large repository.
- A rewrite of the core in Rust would make this easier (the `notify` crate has native drivers and a polling fallback) but costs the no-dependency property and a large port; not a reason by itself.

Whatever is chosen, the stat-based change token (261029-xeph) is the foundation: a native watcher only says when to rescan, and FSEvents can itself ask for a rescan.

## Log

- 2026-10-29: filed.
