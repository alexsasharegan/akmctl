---
id: "261108-wg98"
title: windows support tier adr
status: raw
kind: chore
parent: 261108-2ncf
priority: 2
created: 2026-11-08
updated: 2026-11-08
---

# windows support tier adr

## Goal

A superseding ADR for the platform support tiers. Part of [261108-2ncf](261108-2ncf-basic-windows-support.md). ADR 261029-eegx says Windows is best effort and a Windows failure does not block merging; since 2026-11-08 the Windows CI job counts and a first Windows user is coming. Draft a `proposed` ADR that supersedes it and states what Windows support means at this stage (what is tested, what blocks a release, what is not promised, the unsigned binary), then update the README lines that say "best effort". The decision is the maintainer's; the draft lists the options.

Done when:

- [ ] A `proposed` ADR exists with the options and consequences, for the maintainer to accept.
- [ ] After it is accepted, the README and `docs/release.md` match it.

## Log

- 2026-11-08: filed.
