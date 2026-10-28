---
id: 261028-tphz
title: Ship the workbench UI as the akmctl ui subcommand of the same binary
summary: "The UI is not a separate product: akmctl ui serves the workbench page from the one binary. Read before proposing a separate UI program."
status: accepted
date: 2026-10-28
scope: [internal/server, cmd/akmctl]
authored_by: agent
---

# Ship the workbench UI as the akmctl ui subcommand of the same binary

## Context

The workbench needs the same operations as the CLI and a prebuilt binary people can run without a toolchain. A separate product doubles release work. (Decision log entry Y18.)

## Considered Options

- A separate AKM Workbench program with its own release.
- A subcommand of `akmctl`, the page embedded in the binary.
- No UI.

## Decision

`akmctl ui [--port n] [--open]` starts the HTTP binding (ADR 261028-fnez) on a free loopback port and serves the page embedded in the binary.

## Consequences

- One artifact to install, one release.
- The page and its server share the API's tests.
- A larger UI port (the maintainer's existing workbench) is planned against `internal/server`.

## Confirmation

`akmctl ui` runs from the release binary.
