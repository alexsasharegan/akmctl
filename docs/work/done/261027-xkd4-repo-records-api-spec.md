---
id: 261027-xkd4
type: task
status: done
origin: "[[261026-m74k-agent-ready-repo-layout-spike]]"
created: 2026-10-27
updated: 2026-10-28
authored_by: agent
adr: candidate
title: Write the API spec for repo records
---
# Write the API spec for repo records

## Goal

A transport-neutral spec for querying and checking the records in an agent-ready repository (ADRs and work items), so that several implementations can be written against one contract. Source: [[agent-ready-repo-spec]], "Decided in review" and the query-layer design in the spike's conversation.

Scope as decided so far:

- Operations (names as in the spec's section 6): `ls`, `show`, `search`, `check`, `index`, `new`, `find`, `when`, `renumber`, `board`, `set`, `close`, `migrate`, plus `claim` and `release` for the optional claims extension. `renumber` is not in v1. Each gets a request schema and a response schema.
- A stable JSON response. At least two bindings: the CLI (arguments in, JSON out, exit codes) and HTTP (routes and query parameters in, the same JSON out). Another transport such as gRPC can follow.
- Filters, sort, field selection, limits and section extraction, to keep agent token use low.
- Conformance fixtures: sample repos with the expected outputs, as the executable form of the spec.
- The spec is versioned; a schema change is a major version.

Depends on the layout spec being settled (sections 0 to 5 define the records; section 6 and 7 define the tool and versions, see [[agent-ready-repo-spec]]). The working title is AKM and `akmctl`; the UI is `akmctl ui`.

Done when: the spec exists in the vault (or the new repo) with schemas for each operation, the CLI and HTTP bindings, and a first set of conformance fixtures; the open questions it raises are logged.

- Includes the `migrate` operations: report (`check --target`), plan and apply safe steps, and handoff; see "Versions, manifest and reconciliation" in [[agent-ready-repo-spec]].

## Owned here (Y24)

The detail of the `claims` and `safety` extensions: request and response schemas for `claim` and `release`, the safety findings (file, line, rule, never the matched text), and the conformance fixtures, including the YAML-lite fixtures in section 0 of [[agent-ready-repo-spec]]. Dates and timestamps are strings, in JSON output and in frontmatter (YAML 1.2 core, Y26).

## Outcome

- The API spec is [[akm-api-spec]] (API 1.0, draft): 15 operations, envelope and error codes, CLI and HTTP bindings, write and compare-and-swap rules, versioning, conformance method. Schemas: `design/akm-api/schemas/` (16 files). Fixtures: `design/akm-api/.fixtures/` (two sample repos, 15 read-only cases). Validated with Ajv 8; no implementation has run against it yet. Commits 8fdaec6, a95a29f, 12d01e2, fc68bc2.
- Left for the tool build: write fixtures, the `docs/.akm.lock` name, search ranking, `renumber`, and the YAML-lite fixtures run on the parser.
- Why: HTTP carries writes, loopback only, with a per-run bearer token, a Host and Origin check, and compare-and-swap through etags (required on HTTP, optional on the CLI), so the workbench UI can edit records while an agent works the same files. Rejected: read-only HTTP (the UI would shell out for every change); OpenAPI (more tooling, and the CLI still needs its own description); required etags on the CLI (an extra read per write for little gain).

## Log

- 2026-10-27: filed by Claude Code on the Mac at Alexander's direction, split out of [[261026-m74k-agent-ready-repo-layout-spike]]. Raw until the layout spec is final and the project has a name.
- 2026-10-27: scope widened with `migrate` (proposed in [[agent-ready-repo-spec]]).
- 2026-10-28: [[261026-m74k-agent-ready-repo-layout-spike]] closed. [[261027-ensf-akm-an-agent-knowledge-management]] is accepted; the working title is AKM and `akmctl`. The UI may be named AKM Workbench; undecided.
- 2026-10-28: operation names aligned with the spec after the review ([[akm-spec-review]]).
- 2026-10-28: spec review interview ([[agent-ready-repo-decision-log]], Y1 to Y14): added `set`, `close`, `claim`, `release`; the frontmatter grammar in the spec is the base for the conformance fixtures; the index table format is fixed (Prettier-stable).
- 2026-10-28: edge rules of the frontmatter grammar confirmed (Y15); UI is `akmctl ui`. See [[agent-ready-repo-decision-log]].
- 2026-10-28: took over the extension details (Y24); YAML-lite aligned to YAML 1.2 core (Y26).
- 2026-10-28: claimed (Claude Code on the Mac). Reference forms (Y30), body links (Y31) and property types (Y32) settled first; they affect the schemas.
- 2026-10-28: decisions with Alexander: prose plus JSON Schema; HTTP allows writes, local only; spec and fixtures in the vault design folder (Y34).
- 2026-10-28: draft 1 written and committed: [[akm-api-spec]], `design/akm-api/schemas/` (16 files), `design/akm-api/.fixtures/` (two sample repos, 15 cases, `make_fixtures.py`), `design/akm-api/validate.js`. Ajv 8 (strict, draft 2020-12, installed in the session scratchpad with Alexander's OK) compiles every schema and validates every fixture response. No write fixtures yet. Dot-folder keeps the sample records out of Obsidian's index (Y32).
- 2026-10-28: review round 1 with Alexander: exit codes stay as drafted (every `ok: false` exits 2); writes use compare-and-swap: `if_match` optional on the CLI, required on HTTP (428, `if_match_required`), with a short advisory lock for the check-and-rename; `scope` filter is path prefix for entries ending in `/`, equality otherwise.
- 2026-10-28: review round 2 with Alexander: the response `request` echoes the request as given (no defaults); `ls`, `search` and `board` return a short default field set (`["*"]` for all); the CLI default format is the human table whatever stdout is. Fixtures regenerated and revalidated.
- 2026-10-28: review round 3 with Alexander: references stay bare ids in responses; every HTTP route needs the token; the migrate handoff prompt is returned in the response (CLI `--out` writes a file) and only `docs/.migrate-log` is a file, committed with the migration.
- 2026-10-28: closed. ADR check: candidate (3 new tool and service surface; 5 a future agent could weaken the HTTP security rules without knowing why). Not large enough to draft now; leave to the sweep, or fold into the consolidated v1.0 ADR ([[261028-j3m9-akm-spec-v1-cleanup]]).
