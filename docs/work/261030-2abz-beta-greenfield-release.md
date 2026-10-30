---
id: "261030-2abz"
title: Beta release for greenfield use (v0.2.0)
status: ready
created: 2026-10-30
updated: 2026-10-30
kind: feature
priority: 1
summary: "Ship a public beta (v0.2.0) for greenfield repositories: init, a sound task lifecycle, agent-ready files, the spec at 0.2.0. Adoption/migrate is out of scope."
---

# Beta release for greenfield use (v0.2.0)

## Goal

Alexander wants a public beta tomorrow (2026-10-30) so a colleague can try AKM in a fresh directory: `akmctl init` makes a new repository (and Obsidian vault) that is agent-ready. The beta supports greenfield use only; adoption of an existing repository (`migrate`, discovery) is not claimed and is a 1.0.0 requirement. The beta is tagged `v0.2.0-beta.1`; `v1.0.0` stays reserved. The spec, guide and layout are versioned `0.2.0` to match the beta and keep that number through beta updates.

Scope, in order:

1. `akmctl init` ([261029-0gb4](261029-0gb4-akmctl-init-with-profiles.md)), with schema, fixtures and tests.
2. The task lifecycle: claim, work, close ([261029-7zcp](261029-7zcp-double-claim-log-on-unblock.md), [261029-yj59](261029-yj59-status-based-done-move.md), [261029-f5vn](261029-f5vn-utc-id-date.md), [261029-nf6h](261029-nf6h-mint-ids-slug-from-title.md)).
3. Agent-ready output: `AGENTS.md`, `CLAUDE.md`, the `work-items` and `record-decision` skills that `init` writes.
4. The spec and guide at 0.2.0, with the template section describing `init` ([261029-dsga](261029-dsga-spec-template-section-describes-init.md)).
5. Cosmetics tied to the workbench: page title, header centring, favicon.
6. Public-readiness: Prettier in CI, a README that says what the beta does and does not support, an audit by Alexander, then Alexander makes the repository public and tags.

Done when:

- [ ] A fresh directory plus `akmctl init --profile code` and `--profile notes` gives a repository where `akmctl check` reports 0 errors and an agent can find, claim, work and close an item using only the files written.
- [ ] The lifecycle items above are closed.
- [ ] The spec and guide read Version 0.2.0 and no longer describe a template repository.
- [ ] The README states the beta scope and the unsupported adoption path.
- [ ] CI is green on macOS, Linux and Windows.
- [ ] Alexander has made the repository public and approved the `v0.2.0-beta.1` tag.

## Next

Create the roadmap tiers, then run the init and lifecycle work in parallel worktrees.

## Log

- 2026-10-29: filed.
