---
id: "261108-fabf"
title: winget package
status: raw
kind: chore
parent: 261108-2ncf
priority: 3
created: 2026-11-08
updated: 2026-11-08
---

# winget package

## Goal

A WinGet package. Part of [261108-2ncf](261108-2ncf-basic-windows-support.md). Write the manifest (version, installer and locale files; `InstallerType: zip` with a nested portable `akmctl.exe`) for the Windows archive, pick a package id (for example `alexsasharegan.akmctl`), and submit it to `microsoft/winget-pkgs` (the first review takes days; later versions can go by `wingetcreate update` or a release action). Submitting is the maintainer's call and needs their GitHub identity.

Done when:

- [ ] The manifest passes `winget validate` and installs from a local manifest (`winget install --manifest`) on the Windows runner or a Windows machine.
- [ ] The pull request to `winget-pkgs` is open or merged.
- [ ] `docs/release.md` says how each release updates the package, and the README names the command.

## Log

- 2026-11-08: filed.
