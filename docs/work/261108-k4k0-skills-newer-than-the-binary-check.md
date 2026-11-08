---
id: "261108-k4k0"
title: skills newer than the binary check
status: raw
kind: feature
priority: 2
created: 2026-11-08
updated: 2026-11-08
---

# skills newer than the binary check

## Goal

A project that carries the skills of a newer `akmctl` while the global binary is older finds commands missing (2026-11-08). Make the mismatch visible instead of a surprise: `akmctl skills status`, or `check`, should say when a repository's skill files name a command or version the running binary does not have, and which `akmctl version` to install. Options to weigh: a `min_version` in each skill's frontmatter checked against `Version`; a list of operations each skill references checked against the binary's operations (the skill tests already map operations to references); or both. Related: [261108-8m4a](done/261108-8m4a-global-akmctl-follows-green-main.md).

Done when:

- [ ] Running an older binary on a repository whose skills need a newer one prints one clear line naming the needed version.
- [ ] A test covers it with a skill that needs a version above the binary's.
- [ ] The `akm` skill tells an agent what to do on that line.

## Log

- 2026-11-08: filed.
