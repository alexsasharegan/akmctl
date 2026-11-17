---
id: "261117-n7v4"
title: Align the tool, layout, spec and API versions at 0.4.0
status: active
created: 2026-11-17
updated: 2026-11-17
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-17T07:38
---

# Align the tool, layout, spec and API versions at 0.4.0

## Goal

The tool is 0.3.x, the layout and spec are 0.3.0 and the API is 1.0, versioned independently (`docs/architecture.md`). The API's 1.0 reads as stable beside everything else at 0.x, and readers and LLMs flag the mismatch. The maintainer chose lockstep (2026-11-17): at the 0.4.0 release the tool, layout, specification, guide, API spec and extensions all read 0.4.0, and the API's `api_version` is 0.4 or 0.4.0 (decide the exact form in the ADR).

Done when:

- [ ] A `proposed` ADR records the lockstep rule before 1.0.0, what the API version string becomes, and that the post-1.0 rule stays with the compatibility policy item (261109-vxdr).
- [ ] `check.LayoutVersion`, the manifest versions, `testdata/api` and `docs/spec/akm-api/.fixtures` (via `scripts/sync-fixtures.sh`), the schemas, the spec, guide and API spec headers, the README status line and `docs/architecture.md` agree on 0.4.0; `go test ./...` and `/tmp/akmctl check` pass.
- [ ] A `migrate` step carries a repository from layout 0.3.0 to 0.4.0, if the tool requires one for a manifest bump.
- [ ] `CHANGELOG.md` has an entry under `## Unreleased` naming the API version change for clients.
- [ ] The glossary and `docs/release.md` say the versions move together until 1.0.0.

## Next

Draft the proposed ADR, then list every place the three version strings live and change them in one commit.

## Log

- 2026-11-17: filed.
- 2026-11-17 00:38-07:00: claimed (claude-sonnet-5-5).
