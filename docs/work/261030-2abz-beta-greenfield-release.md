---
id: "261030-2abz"
title: Beta release for greenfield use (v0.2.0)
status: ready
created: 2026-10-30
updated: 2026-10-29
kind: feature
priority: 1
summary: "Ship a public beta (v0.2.0) for greenfield repositories: init, a sound task lifecycle, agent-ready files, the spec at 0.2.0. Adoption/migrate is out of scope."
---

# Beta release for greenfield use (v0.2.0)

## Goal

Alexander wants a public beta tomorrow (2026-10-30) so a colleague can try AKM in a fresh directory: `akmctl init` makes a new repository (and Obsidian vault) that is agent-ready. The beta supports greenfield use only; adoption of an existing repository (`migrate`, discovery) is not claimed and is a 1.0.0 requirement. The beta is tagged `v0.2.0`; `v1.0.0` stays reserved. The spec, guide and layout are versioned `0.2.0` to match the beta and keep that number through beta updates.

Scope, in order:

1. `akmctl init` ([261029-0gb4](done/261029-0gb4-akmctl-init-with-profiles.md)), with schema, fixtures and tests.
2. The task lifecycle: claim, work, close ([261029-7zcp](done/261029-7zcp-double-claim-log-on-unblock.md), [261029-yj59](done/261029-yj59-status-based-done-move.md), [261029-f5vn](done/261029-f5vn-utc-id-date.md), [261029-nf6h](done/261029-nf6h-mint-ids-slug-from-title.md)).
3. Agent-ready output: `AGENTS.md`, `CLAUDE.md`, the `work-items` and `record-decision` skills that `init` writes.
4. The spec and guide at 0.2.0, with the template section describing `init` ([261029-dsga](done/261029-dsga-spec-template-section-describes-init.md)).
5. Cosmetics tied to the workbench: page title, header centring, favicon.
6. Public-readiness: Prettier in CI, a README that says what the beta does and does not support, an audit by Alexander, then Alexander makes the repository public and tags.

Done when:

- [x] A fresh directory plus `akmctl init --profile code` and `--profile notes` gives a repository where `akmctl check` reports 0 errors and an agent can find, claim, work and close an item using only the files written.
- [x] The lifecycle items above are closed.
- [x] The spec and guide read Version 0.2.0 and no longer describe a template repository.
- [x] The README states the beta scope and the unsupported adoption path.
- [ ] CI is green on macOS, Linux and Windows.
- [ ] Alexander has made the repository public and approved the `v0.2.0` tag.

## Next

Alexander: audit, confirm macOS CI is green on the latest main, then tag v0.2.0 (go install @latest needs a non-prerelease tag). Make the repository public before or with the tag.

## Log

- 2026-10-29: filed.
- 2026-10-29: init, lifecycle fixes, UI cosmetics, spec at 0.2.0 (layout and extensions too), Prettier in CI and the beta README are done and merged on main. Remaining: CI green on three platforms, Alexander's audit, public flip, tag.
- 2026-10-29: akm- skills, git_commits and config, init questionnaire, platform tiers ADR and the README rewrite merged on main. Remaining: Alexander's audit, macOS CI green, tag v0.2.0.
