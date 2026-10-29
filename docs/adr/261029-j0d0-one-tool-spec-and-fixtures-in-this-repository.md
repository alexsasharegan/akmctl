---
id: "261029-j0d0"
title: The CLI is the only tool; the standard and its fixtures live in this repository
summary: "The template repository is dropped: akmctl init and akmctl migrate are the way in, and the standard, guide, API spec and fixtures live in this repository."
status: accepted
date: 2026-10-29
authored_by: agent
---

# The CLI is the only tool; the standard and its fixtures live in this repository

## Context

AKM began with a plan for a template repository: a clone template that carries the layout, `AGENTS.md`, the skills and the CI files, with the specification, the guide and the API spec kept in the maintainer's notes. The dogfood of the `akmctl` repository ([261028-eqwq](../work/done/261028-eqwq-akmctl-dogfood-akm-layout.md)) showed that the template would have to be kept in step with the tool, the standard and the migrations, three things that already change together. The tool can do what the template does: write the files, and move an existing repository to the layout by migration. The standard, the guide, the API spec and its fixtures were still outside the repository that implements them, and the fixtures in `testdata/api` were a copy of the notes' fixtures.

## Considered Options

- Keep the template repository next to the tool: two things to version, and a clone cannot run a migration.
- Keep the template and add `akmctl init`: two ways in that drift apart.
- Drop the template: the tool is the only way in; the standard and the fixtures move into this repository so one repository holds the contract, the implementation and its tests.

## Decision

The template repository is dropped. The CLI is the only tool. `akmctl init` creates the layout in a new repository, with profiles `code` and `notes`. `akmctl migrate` takes an existing repository to the layout; adopting AKM is the migration from layout 0.0.0 and gets a discovery phase for a repository with no manifest. The standard, the guide, the API spec and the fixtures live in this repository (`docs/spec/`), with the design history in `docs/design/`, the closed AKM work in `docs/work/done/` and the earlier ADRs in `docs/adr/`. Alexander decided this on 2026-10-29 (decision log entry Y44, [agent-ready-repo-decision-log](../design/agent-ready-repo-decision-log.md)).

## Consequences

- `init` with profiles, `migrate` discovery for a repository with no manifest, and a rewrite of the standard's template section (section 7 and the paragraph "AKM template (a clone template)") are open work items.
- `testdata/api` is still a copy of `docs/spec/akm-api/.fixtures`, kept in step by hand until a work item makes one the source of the other.
- The vault copies of the moved records stay as frozen history, marked as moved.
- The standard changes in the same repository as the code, so a change to a rule and to its check can share a commit and a review.

## Confirmation

`docs/spec/` holds the standard, the guide, the API spec and `akm-api/`; the work items named above are on `docs/ROADMAP.md`; `akmctl check` passes on the repository.
