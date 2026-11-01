---
id: "261101-s3x5"
title: workbench phone access
summary: Reach the workbench from a phone over a LAN or Tailscale bind, with real authentication
status: raw
created: 2026-11-01
updated: 2026-11-01
kind: feature
priority: -1
---

# workbench phone access

## Goal

Today `akmctl ui` binds 127.0.0.1 only, with a per-run token (ADR 261025-r0ye: no phone access in version 1). Phone access needs a LAN or Tailscale bind and real auth rather than a URL token. This is a decision first: an ADR draft (`proposed`) on the bind, the auth and the threat model.

## Log

- 2026-11-01: filed. It was pinned as "later" in `docs/design/pkm-workbench-design.md` (ADR 261025-r0ye) and never became a work item.
