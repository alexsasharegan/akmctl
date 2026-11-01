---
id: "261101-vy77"
title: log entries carry time
status: active
kind: feature
created: 2026-11-01
updated: 2026-11-01
decisions: [261101-gw7b]
summary: Log lines carry a local time and offset (claim, release, log) and log takes at, per proposed ADR 261101-gw7b
claimed_by: claude
claimed_at: 2026-11-01T08:39
---

# log entries carry time

## Goal

Implement [ADR 261101-gw7b](../adr/261101-gw7b-log-entries-carry-a-local-time.md) once the maintainer accepts it: a Log line carries the writer's local time and its offset (`- 2026-11-01 14:30-07:00: text`), and `log` can record when the thing happened. Read ADR 261029-5pe1 (dates, UTC and local) and the ADR before starting. Done when:

- [ ] `log`, `claim` and `release` write the new form; a date-only `as_of` still writes a date-only line, and an `as_of` with a time writes offset `+00:00`.
- [ ] `log` takes `at` (`--at`): the line ends with `(at <timestamp>)`, the date dropped on the stamp's day; an `at` after the write time is `invalid_request`.
- [ ] `Tree.LogEntries` accepts both forms and `Entry` carries `Time` and `Offset` (tests with old lines).
- [ ] The spec template, the `akm` and `akm-track-work` skills (and the embedded copy), the glossary (with the three-conventions table), the changelog, the schemas and the fixtures that show a Log line (both copies) are updated.

## Next

Implementing per ADR 261101-gw7b: shared log-stamp helper in internal/api, then log/claim/release, at, LogEntries, then spec/skills/glossary/changelog/schemas/fixtures.

## Log

- 2026-11-01: filed.
- filed with the proposed ADR 261101-gw7b; split from 261031-mr2q (claimed_at UTC docs), which stays a quick fix.
- 2026-11-01: claimed (claude).
