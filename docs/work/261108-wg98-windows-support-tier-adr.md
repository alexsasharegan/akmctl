---
id: "261108-wg98"
title: windows support tier adr
status: active
kind: chore
parent: 261108-2ncf
priority: 2
created: 2026-11-08
updated: 2026-11-08
claimed_by: claude-sonnet-5-5-win
claimed_at: 2026-11-08T21:49
decisions: [261108-31p6]
---

# windows support tier adr

## Goal

A superseding ADR for the platform support tiers. Part of [261108-2ncf](261108-2ncf-basic-windows-support.md). ADR 261029-eegx says Windows is best effort and a Windows failure does not block merging; since 2026-11-08 the Windows CI job counts and a first Windows user is coming. Draft a `proposed` ADR that supersedes it and states what Windows support means at this stage (what is tested, what blocks a release, what is not promised, the unsigned binary), then update the README lines that say "best effort". The decision is the maintainer's; the draft lists the options.

Done when:

- [x] A `proposed` ADR exists with the options and consequences, for the maintainer to accept.
- [ ] After it is accepted, the README and `docs/release.md` match it.

## Next

Waiting for the maintainer to accept or change ADR 261108-31p6. On acceptance: set 261029-eegx to superseded, edit README (line 7, last Windows bullet), AGENTS.md rule and docs/release.md, tick box 2, close.

## Log

- 2026-11-08: filed.
- 2026-11-08 14:49-07:00: claimed (claude-sonnet-5-5-win).
- 2026-11-08 14:50-07:00: ticked "A `proposed` ADR exists with the options and consequences, for the maintainer to accept.": docs/adr/261108-31p6-windows-support-tier.md: proposed, four options, consequences, recommendation (option 3)
- 2026-11-08 14:50-07:00: drafted ADR 261108-31p6 (proposed, supersedes 261029-eegx); waiting for the maintainer. Box 2 and the README, AGENTS.md and release.md edits follow acceptance
