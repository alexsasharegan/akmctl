---
id: "261101-hsh4"
title: CHANGELOG names a change to the AGENTS.md template
status: done
created: 2026-11-01
updated: 2026-11-07
closed: 2026-11-07
blocked_by: [261101-abw2]
---

# CHANGELOG names a change to the AGENTS.md template

## Goal

A release that changes `internal/initrepo/templates/AGENTS.md.tmpl` or the skill templates reaches existing repositories only if they notice (261101-wghm). Done when:

- [x] CHANGELOG entries for such a change say so in a fixed phrase (for example "The `AGENTS.md` template changed: run `akmctl agents`"), written in `docs/release.md` as a step of Prepare.
- [x] A rule or test fails the release prepare when the template changed since the last tag and no entry carries the phrase, or the choice not to is recorded in the Outcome.

## Next

## Outcome

Prepare step 4 in docs/release.md names the fixed phrase (The AGENTS.md template changed: run akmctl agents) and the command that finds the change (git diff --stat <last-tag> -- internal/initrepo/templates/AGENTS.md.tmpl). The 261101-qr8x changelog entry, the one pending template change since v0.3.3, now carries it. Decided against a test: the template and the changelog are free text, a hash or tag comparison would break on every wording fix, and the maintainer reads the Prepare list at each release; skill templates need no phrase because skills status reports them.

## Log

- 2026-11-01: filed.
- 2026-11-07 23:09-07:00: claimed (claude-sonnet-5-5).
- 2026-11-07 23:09-07:00: ticked "CHANGELOG entries for such a change say so in a fixed phrase (for example "The `AGENTS.md` template changed: run `akmctl agents`"), written in `docs/release.md` as a step of Prepare.": Prepare step 4 in docs/release.md
- 2026-11-07 23:09-07:00: ticked "A rule or test fails the release prepare when the template changed since the last tag and no entry carries the phrase, or the choice not to is recorded in the Outcome.": no test: the Prepare step is the check, choice recorded in the Outcome
