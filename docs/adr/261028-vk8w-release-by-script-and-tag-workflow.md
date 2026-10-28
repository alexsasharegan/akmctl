---
id: 261028-vk8w
title: Release five platform archives with a script, driven by a version tag in CI, without GoReleaser
summary: scripts/release.sh builds five platforms with the standard toolchain; a workflow runs it on a v* tag and uploads archives with checksums. Read before releasing.
status: accepted
date: 2026-10-28
scope: [scripts/release.sh, .github/workflows]
authored_by: agent
---

# Release five platform archives with a script, driven by a version tag in CI, without GoReleaser

## Context

Users install a prebuilt binary. The release must be reproducible and need nothing beyond the Go toolchain. The repository is private and the name is final; going public is the maintainer's call and is not tied to a version. (Decision log entry Y37.)

## Considered Options

- GoReleaser: more features, another tool to pin.
- A shell script with the standard toolchain (CGO off, trimpath, version by linker flag), archives and `SHA256SUMS`.
- Manual builds.

## Decision

`scripts/release.sh <version>` builds macOS arm64 and amd64, Linux amd64 and arm64 and Windows amd64 into `dist/<version>/` with `SHA256SUMS`. `.github/workflows/release.yml` runs it on a `v*` tag and runs `gh release create`. Actions are pinned by commit SHA. The first version is 0.1.0; before 1.0 the CLI and API may change.

## Consequences

- A tag, a release or a visibility change is the maintainer's decision, never an agent's.
- The Windows archive is built on Linux; the Windows path is tested only by CI.

## Confirmation

`scripts/release.sh` and the CI release workflow.
