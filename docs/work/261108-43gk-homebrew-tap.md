---
id: "261108-43gk"
title: homebrew tap
status: ready
kind: chore
parent: 261108-2ncf
priority: 2
created: 2026-11-08
updated: 2026-11-08
---

# homebrew tap

## Goal

A Homebrew tap. Part of [261108-2ncf](261108-2ncf-basic-windows-support.md). Create `homebrew-tap` (a separate repository under the maintainer's account, the maintainer creates it) with `Formula/akmctl.rb` that installs the release archives for macOS and Linux (arm64 and amd64) with their SHA256. Add a step to `release.yml` or a documented command that updates the formula on each tag. The update needs a token with write access to the tap: the maintainer sets that secret; no credential goes in a file.

Done when:

- [ ] `brew install alexsasharegan/tap/akmctl` installs the current release on macOS and runs `akmctl version`.
- [ ] The same formula works on Linux (Homebrew on Linux).
- [ ] The release process in `docs/release.md` says how the formula is updated, and the README names the command.

## Next

Write Formula/akmctl.rb for the four release archives, test it on macOS, then ask the maintainer to create the homebrew-tap repository and the update token.

## Log

- 2026-11-08: filed.
