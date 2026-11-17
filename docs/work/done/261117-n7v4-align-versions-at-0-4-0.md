---
id: "261117-n7v4"
title: Align the tool, layout, spec and API versions at 0.4.0
status: done
created: 2026-11-17
updated: 2026-11-17
closed: 2026-11-17
---

# Align the tool, layout, spec and API versions at 0.4.0

## Goal

The tool is 0.3.x, the layout and spec are 0.3.0 and the API is 1.0, versioned independently (`docs/architecture.md`). The API's 1.0 reads as stable beside everything else at 0.x, and readers and LLMs flag the mismatch. The maintainer chose lockstep (2026-11-17): at the 0.4.0 release the tool, layout, specification, guide, API spec and extensions all read 0.4.0, and the API's `api_version` is 0.4 or 0.4.0 (decide the exact form in the ADR).

Done when:

- [x] A `proposed` ADR records the lockstep rule before 1.0.0, what the API version string becomes, and that the post-1.0 rule stays with the compatibility policy item (261109-vxdr).
- [x] `check.LayoutVersion`, the manifest versions, `testdata/api` and `docs/spec/akm-api/.fixtures` (via `scripts/sync-fixtures.sh`), the schemas, the spec, guide and API spec headers, the README status line and `docs/architecture.md` agree on 0.4.0; `go test ./...` and `/tmp/akmctl check` pass.
- [x] No `migrate` step is needed for a manifest bump (decided in the ADR, as for 0.3.0).
- [x] `CHANGELOG.md` has an entry under `## Unreleased` naming the API version change for clients.
- [x] `docs/release.md` says the versions move together until 1.0.0 (no glossary term is new).

## Next

## Outcome

Merged in pull request 93. Tool, layout, specification, extensions and API share one version: layout 0.4.0, API 0.4, HTTP base /api/v0. The ADR on one version number before 1.0.0 stays proposed for the maintainer. No migration added for a 0.3.0 manifest (S1 warning, as after 0.3.0). Decision: ADR 261117-6qn7.

## Log

- 2026-11-17: filed.
- 2026-11-17: drafted the proposed ADR on one version number before 1.0.0. Changed `check.LayoutVersion` to 0.4.0, `api.Version` to 0.4, the HTTP base to `/api/v0`, this repository's manifest, both fixture copies (edited the case files and schemas directly, then `scripts/sync-fixtures.sh`; `make_fixtures.py` updated as the frozen reference), the specification, guide and API spec text, the README status, architecture, release doc, the `version` skill reference and one skill example that named 0.4.0. No migration added: `migrate` reports 0.3.0 to 0.4.0 as unsupported, as it did for 0.2.0 (see the ADR). Done-when items left open: none except the glossary, which gains no new term.
- 2026-11-17 00:38-07:00: claimed (claude-sonnet-5-5).
