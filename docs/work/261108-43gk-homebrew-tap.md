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

A Homebrew tap. Part of [261108-2ncf](261108-2ncf-basic-windows-support.md). Generate `Formula/akmctl.rb` for a future `homebrew-tap` repository (a separate repository under the maintainer's account, which the maintainer creates later) with a script (`scripts/homebrew-formula.sh`) that installs the release archives for macOS and Linux (arm64 and amd64) with their SHA256. Publishing the tap waits for the same confidence items as the WinGet package ([261108-fabf](done/261108-fabf-winget-package.md)); the maintainer decides when. The update on each tag needs a token with write access to the tap: the maintainer sets that secret; no credential goes in a file.

Done when:

- [ ] The generated formula installs the current release on macOS from a local file (`brew install --formula ./akmctl.rb`) and runs `akmctl version`.
- [ ] The same formula works on Linux (Homebrew on Linux).
- [ ] CI generates the formula for a dummy release on every run and checks it.
- [ ] The publishing steps are written down (not done) in `docs/release.md`.

## Next

Write Formula/akmctl.rb for the four release archives, test it on macOS, then write the publishing steps; do not publish.

## Log

- 2026-11-08: filed.
