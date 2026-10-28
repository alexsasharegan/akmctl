---
id: "261028-ehv6"
title: Adopt the AKM layout in this repository
summary: Stamp the layout, ADRs and work items onto akmctl itself and make akmctl check pass.
status: active
kind: chore
created: 2026-10-28
updated: 2026-10-28
claimed_by: claude-code@mac
claimed_at: 2026-10-28T22:53
---

# Adopt the AKM layout in this repository

## Goal

Keep the tool's own memory the AKM way, and use the repository as the first test case of the standard. The layout (`AGENTS.md`, `docs/`, the manifest, the Prettier config, the skills) is in place; the tool's decisions are ADRs; its known gaps are work items.

Done when:

- [ ] `akmctl check`, built from this branch, passes on the repository.
- [ ] The tool's decisions are ADRs in `docs/adr/` and `akmctl index` output is committed.
- [ ] The known gaps are work items under `docs/work/` and on the roadmap.
- [ ] `go vet ./...` and `go test ./...` pass, and CI passes on the branch.

## Next

Push the branch, read the CI result, then close this item with the evidence.

## Log

- 2026-10-28: filed.
- 2026-10-28: claimed (claude-code@mac).
