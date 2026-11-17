---
id: "261117-bja4"
title: Release v0.4.0
status: active
kind: chore
priority: 10
created: 2026-11-17
updated: 2026-11-17
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-17T07:52
---

# Release v0.4.0

## Goal

Cut v0.4.0 following `docs/release.md`. It is a MINOR bump: the tool, layout, specification and API now share one version (ADR 261117-6qn7), the API path is `/api/v0`, and `api_version` is `0.4`.

Done when:

- [ ] `CHANGELOG.md` has a `## 0.4.0 - <date>` heading with the Unreleased entries beneath it, and the `AGENTS.md` template phrase is present (it is, from work item 261101-qr8x).
- [ ] The full check passes and CI on the merged commit passes.
- [ ] The tag `v0.4.0` is pushed by the maintainer and the release workflow passes.
- [ ] `gh release view v0.4.0` lists the five archives and `SHA256SUMS`.
- [ ] `go install github.com/alexsasharegan/akmctl/cmd/akmctl@v0.4.0 && akmctl version` prints `0.4.0 (API 0.4)`.

## Next

Move the changelog entries under the 0.4.0 heading, run the full check, open the pull request; the tag is the maintainer's step.

## Log

- 2026-11-17: filed.
- 2026-11-17 00:52-07:00: claimed (claude-sonnet-5-5).
- 2026-11-17 00:53-07:00: Prepared: changelog heading 0.4.0 - 2026-11-17 (AGENTS.md template phrase already present from the claim-name item); full check passed (vet, race tests, gofmt, records check, format, schemas); dry run SKIP_TESTS=1 scripts/release.sh 0.4.0 built five archives and SHA256SUMS, darwin arm64 binary prints akmctl 0.4.0 (API 0.4). Closed the alignment item that this branch carries as a rider. CI on the merged commit passed on all three platforms. Tag is the maintainer's step.
