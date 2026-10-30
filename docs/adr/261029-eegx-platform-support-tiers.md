---
id: "261029-eegx"
title: "Platform support tiers: macOS and Linux supported, Windows best effort before 1.0.0"
summary: "macOS is supported, Linux is supported, Windows is best effort and out of scope for any release below 1.0.0."
status: accepted
date: 2026-10-29
---

# Platform support tiers: macOS and Linux supported, Windows best effort before 1.0.0

## Context

The 0.2.0 beta goes to macOS users first. AGENTS.md says the binary must work on macOS, Linux and Windows and that CI covers all three. Windows has cost time (a flaky compare-and-swap test, work items 261029-q8z3 and 261028-jzsx) and will cost more with a terminal questionnaire and file removal in `init` and `config`.

## Considered Options

- Keep all three platforms as equal requirements. Rejected: the cost goes to a platform no current user runs.
- Drop Windows. Rejected: the code is portable and CI already builds it; there is no reason to break it.
- Tier the platforms: macOS and Linux supported, Windows best effort.

## Decision

macOS is the primary platform. Linux is supported, because it is close to macOS and cheap to keep working. Windows is best effort: the tool should build and the tests should run there, but a Windows-only failure does not block a release, and Windows support is out of scope for any release below 1.0.0. Release archives for Windows may still be built; they carry no support claim. Windows support becomes a requirement again at 1.0.0 unless a later decision says otherwise.

## Consequences

- CI keeps running Windows, but a Windows failure does not block merging or a release below 1.0.0.
- The README states the tiers.
- Windows items (261029-q8z3, 261028-jzsx) stay in Later and are not required for the beta.
- New features need to be correct on macOS and Linux first; Windows fixes are welcome but not gating.
- AGENTS.md's rule that the binary must work on all three platforms is amended to say this.
