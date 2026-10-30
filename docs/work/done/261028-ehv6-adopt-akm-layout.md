---
id: "261028-ehv6"
title: Adopt the AKM layout in this repository
summary: Stamp the layout, ADRs and work items onto akmctl itself and make akmctl check pass.
status: done
kind: chore
created: 2026-10-28
updated: 2026-10-28
closed: 2026-10-28
---

# Adopt the AKM layout in this repository

## Goal

Keep the tool's own memory the AKM way, and use the repository as the first test case of the standard. The layout (`AGENTS.md`, `docs/`, the manifest, the Prettier config, the skills) is in place; the tool's decisions are ADRs; its known gaps are work items.

Done when:

- [x] `akmctl check`, built from this branch, passes on the repository.
- [x] The tool's decisions are ADRs in `docs/adr/` and `akmctl index` output is committed.
- [x] The known gaps are work items under `docs/work/` and on the roadmap.
- [x] `go vet ./...` and `go test ./...` pass, and CI passes on the branch.

## Next

Merge the branch (the maintainer's call). Then build the template repository from the lessons note in the maintainer's notes, and start the `mint_ids` item.

## Log

- 2026-10-28: filed.
- 2026-10-28: claimed (claude-code@mac).
- 2026-10-28: choices: README.md and CHANGELOG.md stay at the repository root (the layout lists the README at the root; a changelog is release notes, read by the release process). The layout was copied by hand, not migrated (a repository with no manifest is layout 0.0.0 and no migration reaches 1.0.0). Fourteen of the decision log's entries became ADRs (Y28 to Y37 and the earlier ones the code embodies: Y2, Y3, Y5, Y6, Y8 to Y10, Y13, Y14, Y17, Y18, Y26, Y29 to Y31, Y34); the rest stay as history (name, process and standard-text decisions).
- 2026-10-28: the ADRs are `accepted` with `authored_by: agent`: the decisions were made by the maintainer in interviews; the merge of this branch is the acceptance (specification, section 4). Removing the marker is the review.
- 2026-10-28: found and fixed: `akmctl new work` output failed `prettier --check` (no blank line after `## Log`).
- 2026-10-28: CI ran on the branch only after its trigger was widened to every branch push; it passed on Linux, macOS and Windows.

## Outcome

- `akmctl check` passes: 0 errors, 0 warnings, built from this branch (`go build -o /tmp/akmctl ./cmd/akmctl && /tmp/akmctl check`); the same command is a CI step on Linux.
- The decisions are 14 ADRs in `docs/adr/` with the generated index (`akmctl index`, rule A4 passes); each cites its decision log entries.
- The gaps are work items `261028-8dd4`, `x1b4`, `a2wd`, `e7ej`, `1k0e`, `jzsx`, `k0ek`, `esyt`, `fcc1`, `nkzy`, all on `docs/ROADMAP.md`.
- `go vet ./...` and `go test ./...` pass; CI run 40783633613 passed on Linux, macOS and Windows.
- Not verified: Prettier 3.9.10 (the pin) was not run; Prettier 3.8.4, installed globally, reports every document clean. CI does not run Prettier.
