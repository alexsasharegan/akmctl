---
id: "261101-abw2"
title: akmctl agents prints the filled AGENTS.md template
status: active
created: 2026-11-01
updated: 2026-11-06
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-06T20:25
---

# akmctl agents prints the filled AGENTS.md template

## Goal

Existing repositories cannot see what the shipped `AGENTS.md` template now says (261101-wghm). `init` writes `AGENTS.md` once and never overwrites it, and the file is the repository's own. Add `akmctl agents`: print the template as `init` would write it for this repository (docs dir, wikilinks and git-commit settings read from the repo), read-only, so a user or agent can audit their file and splice lines as they judge. Done when:

- [ ] `akmctl agents` prints the filled template and writes nothing, in a repository that has `init` settings and in one that lacks them (defaults).
- [ ] A test shows its output equals the `AGENTS.md` that `init` writes for the same settings.
- [ ] Usage text, guide Upgrading section, glossary entry for the command, CHANGELOG entry, and the API spec if the command is in it.

## Next

Factor the template fill in internal/initrepo/initrepo.go so init and agents share it; add the command; add the equality test.

## Log

- 2026-11-01: filed.
- 2026-11-06 13:25-07:00: claimed (claude-sonnet-5-5).
