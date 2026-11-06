---
id: "261101-wghm"
title: upgrade path for existing users
status: done
created: 2026-11-01
updated: 2026-11-06
closed: 2026-11-06
parent: 261101-vs5c
---

# upgrade path for existing users

## Goal

How do changes to the agent rails reach a repository that already ran `init`? Skills refresh with `akmctl skills update`; the binary arrives by release or self update (261030-55cs) and carries refusals such as `claim` on an empty Next, `close` emptying Next and CM4; the CHANGELOG is the only notice. `AGENTS.md` is the gap: `init` writes it once from `AGENTS.md.tmpl` and never overwrites it, and it is the repository's own file, so the tool must not edit it. Decision (maintainer, 2026-11-01): the tool ships the template and prints it, filled with the repository's settings, so a user or agent audits their own file and splices, replaces or appends as they judge. Layout migration is out of scope. No ADR: a read-only print. Done when:

- [x] `docs/spec/akm-guide.md` has an Upgrading section covering the binary, the skills and `AGENTS.md`: the command a user runs for each and what it leaves alone.
- [x] A `ready` item exists for `akmctl agents` (print the filled template, read-only), and one for the CHANGELOG convention that names a template change.
- [x] `akmctl skills status` is checked for whether it says skills are behind the binary; if not, a `ready` item is filed.

## Next

## Outcome

The guide's section 8 has an Upgrading table for the binary, the skills, AGENTS.md and the layout. skills status already reports 'outdated' (checked on a scratch repo), so no item was filed for it. The two gaps, akmctl agents (261101-abw2) and the CHANGELOG phrase (261101-hsh4), already existed and are now on the roadmap under Next. Until abw2 ships, the guide tells a user to compare against a fresh init.

## Log

- 2026-11-01: filed.
- 2026-11-01 20:26-07:00: Found: 'akmctl skills update' exists and refreshes skills and retires renamed ones, so the skills channel has a path. Still unanswered: AGENTS.md (never touched by init), new refusals (claim with empty Next, close empties Next, CM4) and where CHANGELOG notes reach users. Needs: the Goal narrowed to an answer written in docs/guide, plus one filed gap per channel.
- 2026-11-01 20:39-07:00: refined: Goal narrowed to a guide section plus filed gaps; facts gathered (skills update, migrate, AGENTS.md written once, CHANGELOG as only notice)
- 2026-11-01 20:45-07:00: decision: AGENTS.md channel is a read-only 'akmctl agents' print of the filled template; layout out of scope; optional --missing later
- 2026-11-06 12:05-07:00: claimed (claude-sonnet-5-5).
- 2026-11-06 12:05-07:00: Checked skills status on a scratch repo with an outdated skill: it reports 'outdated' and update --dry-run shows the diff, so no item is needed. Wrote the guide's Upgrading section; abw2 and hsh4 already exist and are now on the roadmap.
