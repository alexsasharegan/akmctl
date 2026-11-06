---
id: "261101-wghm"
title: upgrade path for existing users
status: active
created: 2026-11-01
updated: 2026-11-06
parent: 261101-vs5c
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-06T19:05
---

# upgrade path for existing users

## Goal

How do changes to the agent rails reach a repository that already ran `init`? Skills refresh with `akmctl skills update`; the binary arrives by release or self update (261030-55cs) and carries refusals such as `claim` on an empty Next, `close` emptying Next and CM4; the CHANGELOG is the only notice. `AGENTS.md` is the gap: `init` writes it once from `AGENTS.md.tmpl` and never overwrites it, and it is the repository's own file, so the tool must not edit it. Decision (maintainer, 2026-11-01): the tool ships the template and prints it, filled with the repository's settings, so a user or agent audits their own file and splices, replaces or appends as they judge. Layout migration is out of scope. No ADR: a read-only print. Done when:

- [ ] `docs/spec/akm-guide.md` has an Upgrading section covering the binary, the skills and `AGENTS.md`: the command a user runs for each and what it leaves alone.
- [ ] A `ready` item exists for `akmctl agents` (print the filled template, read-only), and one for the CHANGELOG convention that names a template change.
- [ ] `akmctl skills status` is checked for whether it says skills are behind the binary; if not, a `ready` item is filed.

## Next

Run `akmctl skills status` on a repo with an old skill copy; write the Upgrading section in the guide; file the two items.

## Log

- 2026-11-01: filed.
- 2026-11-01 20:26-07:00: Found: 'akmctl skills update' exists and refreshes skills and retires renamed ones, so the skills channel has a path. Still unanswered: AGENTS.md (never touched by init), new refusals (claim with empty Next, close empties Next, CM4) and where CHANGELOG notes reach users. Needs: the Goal narrowed to an answer written in docs/guide, plus one filed gap per channel.
- 2026-11-01 20:39-07:00: refined: Goal narrowed to a guide section plus filed gaps; facts gathered (skills update, migrate, AGENTS.md written once, CHANGELOG as only notice)
- 2026-11-01 20:45-07:00: decision: AGENTS.md channel is a read-only 'akmctl agents' print of the filled template; layout out of scope; optional --missing later
- 2026-11-06 12:05-07:00: claimed (claude-sonnet-5-5).
