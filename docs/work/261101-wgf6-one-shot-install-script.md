---
id: "261101-wgf6"
title: one shot install script
status: active
created: 2026-11-01
updated: 2026-11-01
kind: feature
summary: A one-line curl install script that downloads, verifies and installs the latest release archive
claimed_by: claude
claimed_at: 2026-11-02T00:52
---

# one shot install script

## Goal

A user installs the latest `akmctl` release with one command, on macOS first and on Linux too, without Go and without a browser download (which macOS quarantines). The `release` workflow already builds and attaches all five archives and `SHA256SUMS` on a `v*` tag (v0.3.3 shipped through it), so no workflow change is needed.

Done when:

- [x] `scripts/install.sh` picks the archive for the OS and CPU, downloads the latest release (or `AKMCTL_VERSION`), verifies it against `SHA256SUMS`, and installs `akmctl` into `AKMCTL_INSTALL_DIR` (default `~/.local/bin`).
- [x] The README install section leads with the one-line command and a PATH hint; `docs/release.md` says the script is fetched from `main`.
- [x] Tried against the real v0.3.3 release; changelog and glossary updated.

## Next

Push, run the real curl one-liner against main, then close.

## Log

- 2026-11-01: filed.
- 2026-11-01 17:52-07:00: claimed (claude).
- 2026-11-01 17:52-07:00: ticked "`scripts/install.sh` picks the archive for the OS and CPU, downloads the latest release (or `AKMCTL_VERSION`), verifies it against `SHA256SUMS`, and installs `akmctl` into `AKMCTL_INSTALL_DIR` (default `~/.local/bin`).": done; installed v0.3.3 and v0.3.2 into temp dirs, bad version refused
- 2026-11-01 17:52-07:00: ticked "The README install section leads with the one-line command and a PATH hint; `docs/release.md` says the script is fetched from `main`.": done; installed v0.3.3 and v0.3.2 into temp dirs, bad version refused
- 2026-11-01 17:52-07:00: ticked "Tried against the real v0.3.3 release; changelog and glossary updated.": done; installed v0.3.3 and v0.3.2 into temp dirs, bad version refused
