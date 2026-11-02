---
id: "261101-wghm"
title: upgrade path for existing users
status: ready
created: 2026-11-01
updated: 2026-11-01
parent: 261101-vs5c
---

# upgrade path for existing users

## Goal

How do changes to the agent rails reach a repository that already ran `init`? Known so far, from reading the code and `--help`: skills refresh with `akmctl skills update`; the layout moves with `akmctl migrate` (guide, adopting AKM); `AGENTS.md` is written once by `init` from `AGENTS.md.tmpl` and never overwritten; the binary arrives by release or self update (261030-55cs) and carries refusals such as `claim` on an empty Next, `close` emptying Next and CM4; the CHANGELOG is the only notice. Unknown: whether anything tells a user their skills or `AGENTS.md` are behind the binary. Done when:

- [ ] A short section in `docs/spec/akm-guide.md` (Upgrading) lists each channel (binary, skills, `AGENTS.md`, layout, CHANGELOG) with the command a user runs to pick it up and what it does not touch.
- [ ] `akmctl skills status` is checked for whether it says skills are behind the binary; if not, a `ready` item is filed.
- [ ] Each gap gets its own item: no way to refresh the rails wording in an existing `AGENTS.md` (for example a diff against the template), and no upgrade note for a new refusal in the CHANGELOG (a `Behaviour change` marker). A hard-to-reverse choice gets an ADR draft.

## Next

Run `akmctl skills status` on a repo with an old skill copy to answer the second box; then write the Upgrading section in the guide; file the gap items.

## Log

- 2026-11-01: filed.
- 2026-11-01 20:26-07:00: Found: 'akmctl skills update' exists and refreshes skills and retires renamed ones, so the skills channel has a path. Still unanswered: AGENTS.md (never touched by init), new refusals (claim with empty Next, close empties Next, CM4) and where CHANGELOG notes reach users. Needs: the Goal narrowed to an answer written in docs/guide, plus one filed gap per channel.
- 2026-11-01 20:39-07:00: refined: Goal narrowed to a guide section plus filed gaps; facts gathered (skills update, migrate, AGENTS.md written once, CHANGELOG as only notice)
