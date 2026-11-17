---
id: "261117-bja4"
title: Release v0.4.0
status: done
kind: chore
priority: 10
created: 2026-11-17
updated: 2026-11-17
closed: 2026-11-17
---

# Release v0.4.0

## Goal

Cut v0.4.0 following `docs/release.md`. It is a MINOR bump: the tool, layout, specification and API now share one version (ADR 261117-6qn7), the API path is `/api/v0`, and `api_version` is `0.4`.

Done when:

- [x] `CHANGELOG.md` has a `## 0.4.0 - <date>` heading with the Unreleased entries beneath it, and the `AGENTS.md` template phrase is present (it is, from work item 261101-qr8x).
- [x] The full check passes and CI on the merged commit passes.
- [x] The tag `v0.4.0` is pushed by the maintainer and the release workflow passes.
- [x] `gh release view v0.4.0` lists the five archives and `SHA256SUMS`.
- [x] `go install github.com/alexsasharegan/akmctl/cmd/akmctl@v0.4.0 && akmctl version` prints `0.4.0 (API 0.4)`.

## Next

## Outcome

Changelog heading 0.4.0 - 2026-11-17 merged in pull request 94. Tag v0.4.0 pushed by the maintainer; the release workflow run 44042799991 passed in 1m35s. gh release view v0.4.0 lists the five archives (darwin amd64 and arm64, linux amd64 and arm64, windows amd64) and SHA256SUMS. install.sh with AKMCTL_VERSION=0.4.0 and go install at v0.4.0 both print akmctl 0.4.0 (API 0.4). Decision: none.

## Log

- 2026-11-17: filed.
- 2026-11-17 00:52-07:00: claimed (claude-sonnet-5-5).
- 2026-11-17 00:53-07:00: Prepared: changelog heading 0.4.0 - 2026-11-17 (AGENTS.md template phrase already present from the claim-name item); full check passed (vet, race tests, gofmt, records check, format, schemas); dry run SKIP_TESTS=1 scripts/release.sh 0.4.0 built five archives and SHA256SUMS, darwin arm64 binary prints akmctl 0.4.0 (API 0.4). Closed the alignment item that this branch carries as a rider. CI on the merged commit passed on all three platforms. Tag is the maintainer's step.
- 2026-11-17 01:03-07:00: ticked "`CHANGELOG.md` has a `## 0.4.0 - <date>` heading with the Unreleased entries beneath it, and the `AGENTS.md` template phrase is present (it is, from work item 261101-qr8x).": Merged in pull request 94; the template phrase is in the claim-name entry
- 2026-11-17 01:03-07:00: ticked "The full check passes and CI on the merged commit passes.": Full check passed on the branch; CI on the merged commit passed
- 2026-11-17 01:03-07:00: ticked "The tag `v0.4.0` is pushed by the maintainer and the release workflow passes.": Release workflow run 44042799991 passed
- 2026-11-17 01:03-07:00: ticked "`gh release view v0.4.0` lists the five archives and `SHA256SUMS`.": gh release view v0.4.0 lists five archives and SHA256SUMS
- 2026-11-17 01:03-07:00: ticked "`go install github.com/alexsasharegan/akmctl/cmd/akmctl@v0.4.0 && akmctl version` prints `0.4.0 (API 0.4)`.": install.sh and go install both print akmctl 0.4.0 (API 0.4)
