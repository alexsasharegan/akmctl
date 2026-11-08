---
id: "261108-fabf"
title: winget package
status: active
kind: chore
parent: 261108-2ncf
priority: 2
created: 2026-11-08
updated: 2026-11-08
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-08T17:46
---

# winget package

## Goal

A WinGet package. Part of [261108-2ncf](261108-2ncf-basic-windows-support.md). Write the manifest (version, installer and locale files; `InstallerType: zip` with a nested portable `akmctl.exe`) for the Windows archive, pick a package id (for example `alexsasharegan.akmctl`), and generate the three files from a release's `SHA256SUMS` with a script (`scripts/winget-manifests.sh`), so a release can be packaged and checked without publishing anything. Submitting waits until the Windows PowerShell harness ([261108-hwcq](261108-hwcq-windows-powershell-harness.md)), the CRLF and path audit ([261108-kbve](261108-kbve-windows-crlf-and-path-audit.md)) and the first-run checklist ([261108-7avk](261108-7avk-windows-first-run-checklist.md)) give enough confidence in the build; the maintainer decides when, and submits. The first review at `microsoft/winget-pkgs` takes days.

Done when:

- [ ] The manifest passes `winget validate` and installs from a local manifest (`winget install --manifest`) on the Windows runner or a Windows machine.
- [ ] CI generates the manifests for a dummy release on every run and checks them.
- [ ] The submission steps are written down (not done) in `docs/release.md`.

## Next

Write the three manifests for the Windows zip, validate them on the Windows runner, then write the submission steps; do not submit.

## Log

- 2026-11-08: filed.
- 2026-11-08 10:46-07:00: claimed (claude-sonnet-5-5).
