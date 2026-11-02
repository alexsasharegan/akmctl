---
id: "261101-hsh4"
title: CHANGELOG names a change to the AGENTS.md template
status: ready
created: 2026-11-01
updated: 2026-11-01
blocked_by: [261101-abw2]
---

# CHANGELOG names a change to the AGENTS.md template

## Goal

A release that changes `internal/initrepo/templates/AGENTS.md.tmpl` or the skill templates reaches existing repositories only if they notice (261101-wghm). Done when:

- [ ] CHANGELOG entries for such a change say so in a fixed phrase (for example "The `AGENTS.md` template changed: run `akmctl agents`"), written in `docs/release.md` as a step of Prepare.
- [ ] A rule or test fails the release prepare when the template changed since the last tag and no entry carries the phrase, or the choice not to is recorded in the Outcome.

## Next

Depends on 261101-abw2 for the command name. Add the phrase and the Prepare step to docs/release.md, then decide whether a check is worth its cost.

## Log

- 2026-11-01: filed.
