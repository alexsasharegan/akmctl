---
id: "261029-xeph"
title: Measure the workbench poll, then add a conditional poll with a change token
status: raw
created: 2026-10-29
updated: 2026-10-29
kind: feature
priority: 3
summary: Measure what one workbench poll costs on a repository of about 10000 files; if it is slow, add a change token and If-None-Match so an unchanged poll answers 304
---

# Measure the workbench poll, then add a conditional poll with a change token

## Goal

Step one is a measurement: build a scratch repository of about 10,000 work files and time one `GET records` poll (server work per request, response size) with one tab and with ten.

Only if the number is bad, step two: the server computes a change token over the work, done and ADR folders and the roadmap (name, size and modification time of each file, from `stat`, no reads), shared for about a second across tabs, and the page sends it as `If-None-Match` so an unchanged poll answers `304`. A full refresh every 30 seconds or so covers coarse modification times. A directory's modification time alone is not enough, because editors write in place.

Done when:

- [ ] the measurement is recorded in this item's Log
- [ ] if built: an unchanged poll answers 304, and a change shows within the current two seconds
- [ ] the protocol is the one a watcher could feed later (see 261029-12sx)

## Log

- 2026-10-29: filed.
