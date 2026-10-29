---
type: design
status: active
created_at: 2026-10-28
updated_at: 2026-10-28
authored_by: agent
---

# AKM template repository: what the dogfood taught

Written by Claude Code on 2026-10-28 after stamping the AKM layout onto the `akmctl` repository ([261028-eqwq-akmctl-dogfood-akm-layout](../work/done/261028-eqwq-akmctl-dogfood-akm-layout.md)). The repository is the first one outside the vault to carry the layout; the branch is `dogfood-layout` in `~/github.com/alexsasharegan/akmctl`. This note is the input for building the template repository (specification, section 7, [agent-ready-repo-spec](../spec/agent-ready-repo-spec.md)). There is no template repository yet, so the layout was copied by hand and each file below is one the template must ship.

## How it was applied

- **Copy, not `akmctl migrate`.** A repository with no manifest reads as layout `0.0.0`, and 1.0.0 is the first layout, so no migration exists ([akm-guide](../spec/akm-guide.md), section 7). Confirmed with the tool: a directory with an `AGENTS.md` and no manifest reports `layout 0.0.0 -> 1.0.0: unsupported`; a directory with neither `AGENTS.md` nor `docs/.layout` makes `check` and `migrate` stop with exit 2, "no repository root". Adoption is a one-time conversion, as the guide says.
- **Order that worked:** manifest and `AGENTS.md` first (the tool then finds the root and `check` runs clean on an empty repository), then the roadmap and the ADR index (run `akmctl index` once so rule A4 has a file to compare), then records.
- **`README.md` and `CHANGELOG.md` stay at the repository root.** The layout already lists `README.md` at the root; a changelog is release notes for the product, not a work record, and the release process reads it. The records live under `docs/`. Nothing was absorbed.

## What the template needs to ship

Stamped into every new repository (copied, then edited):

| File                                                               | Notes from the dogfood                                                                                                                                                                                                                                                                                                                  |
| ------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `AGENTS.md`                                                        | Skeleton with the five parts, the required rules, the plain-file fallback table and the id recipe. The dogfood version is 42 lines; the target is 60 to 80. It names no tool version in the template: the tool's own repository builds the tool from source, others need a minimum-version line.                                        |
| `CLAUDE.md`                                                        | One line, `@AGENTS.md`.                                                                                                                                                                                                                                                                                                                 |
| `.prettierrc.json`                                                 | Must be written the way Prettier writes it: `{ "embeddedLanguageFormatting": "off" }` with the spaces. The one-line form in the specification (`{"embeddedLanguageFormatting": "off"}`) fails `prettier --check .`.                                                                                                                     |
| `.prettierignore`                                                  | `docs/adr/README.md`, plus whatever the project has that must not be reformatted (here `testdata/`, `dist/`, `internal/`). The template needs a place for project lines.                                                                                                                                                                |
| `.gitignore`                                                       | A line for `/docs/.akm.lock`.                                                                                                                                                                                                                                                                                                           |
| `docs/.layout`                                                     | `layout` and the extensions the repository keeps. The dogfood keeps `adr-index`, `work-items`, `claims` and `safety`; `safety` also wants `docs/.akm-deny` (comment lines only to start).                                                                                                                                               |
| `docs/ROADMAP.md`                                                  | The four headings, empty. `check` accepts it.                                                                                                                                                                                                                                                                                           |
| `docs/adr/README.md`                                               | Generated: run `akmctl index` when stamping. An empty repository's index is the one comment line.                                                                                                                                                                                                                                       |
| `docs/work/`                                                       | `docs/work/done/` need not ship: git keeps no empty directory, and `check` passes without it; verified that `akmctl close` creates `done/` on the first close.                                                                                                                                                                          |
| `docs/architecture.md`                                             | Optional (tier). A short skeleton: packages or components, invariants, versions.                                                                                                                                                                                                                                                        |
| `.agents/skills/work-items/` and `.agents/skills/record-decision/` | `SKILL.md` plus `assets/work-item.md` and `assets/adr.md`, the record templates for the no-tool path. No such skills existed anywhere: the dogfood wrote them from the specification (about 40 lines each). They are the first draft of the template's skills; the vault's `task-lifecycle` skill is vault-specific and was not reused. |

Kept by the template, not stamped: `VERSION` (the same number as the manifest's `layout`), `CHANGELOG.md` (the template's own history), the template's `README.md` (replaced by the new project's), an `EXTENSION` file per extension (a stamped repository records the extension in the manifest only; unused extensions are deleted).

## Where the specification and the tool disagreed with Prettier

The standard promises that records pass the pinned Prettier. The dogfood found four places where they did not (Prettier 3.8.4 was available globally; 3.9.10 is the pin, see below).

1. **`akmctl new work` wrote `## Log` directly followed by a list item.** Prettier inserts a blank line after a heading, so every new item failed `prettier --check`. Fixed in the tool (commit `1c8c1a2`, with a test). The specification's section example has the same shape (`## Log` then `- YYYY-MM-DD: ...`); the example should show the blank line.
2. **`.prettierrc.json`** (above).
3. **A `*` in a heading** is escaped by Prettier (`v*` becomes `v\*`), so the heading no longer matches the title; titles avoid it.
4. **Hand-written tables** (the architecture map) need Prettier's padding. Run Prettier on every document the template ships.

Where the Prettier pin lives is unspecified. A Go repository has no `package.json`, so the dogfood puts the version in the `AGENTS.md` Commands line (`npx prettier@3.9.10 --check ...`) and CI does not run Prettier. The template should say where a pin goes for a repository with Node and for one without.

## Minting ids for existing material (data for `mint_ids`)

Fourteen ADR ids were minted by hand for decisions made on 2026-10-27 and 10-28; eleven work items came from `akmctl new`. What the hand path cost:

- The recipe prints today's date. An id for an old decision needs the date part edited, and a decision can be older than the repository (the id then names the decision, not the file).
- The recipe checks only files that exist. Ids minted in one batch before the files are written can collide; a loop with a list of ids already issued was needed.
- Each id is written twice (filename and `id` field) and then cited in other records (ADRs cite each other by id). A rename after the fact touches all of them.
- `akmctl new` quotes the id (`id: "261028-ehv6"`); the hand-written records leave it bare. Both parse; the template's examples and the tool should agree.
- `akmctl new` always dates the id and `created` today and ignores the record templates in the skills (it builds the file itself), so the templates are for the no-tool path only. A record for an old decision cannot be created with `new` at all; the 14 ADRs were generated by a throwaway script that wrote the files.
- A `mint_ids` step needs: a date source field (or an explicit date), batch uniqueness, a rename with the `id` field written, and a report of references it did not rewrite. The work item is `261028-8dd4` in the tool repository.

## Cost and fit

- The layout is small enough to stamp in a few minutes once the skills exist. The long parts were the content: 14 ADRs and 11 work items.
- `check` ran clean (0 errors, 0 warnings) with the `claims` and `safety` extensions on across 25 records; the safety rules raised no false positive.
- The ADR route worked for decisions that were already made: they were written as `accepted` with `authored_by: agent`, on the reading that merging the branch is the acceptance (specification, section 4). The template's `record-decision` skill says an agent drafts as `proposed`; a bulk import of past decisions is a case the skill does not cover.
- Sample repositories inside the repository (`testdata/api/repos/*` each have their own `AGENTS.md` and `docs/.layout`) did not confuse `check`: the root is the nearest ancestor of the working directory.
- CI: the workflow triggers on a push to `main` and on pull requests only, so a branch push ran nothing. The dogfood branch widened the trigger to every branch and added a `check` step (Linux). A template that ships CI should decide this.

## Open points for the template

- The Prettier pin location, and whether the template's CI runs Prettier.
- Whether `akmctl new` should read the repository's record templates.
- Whether the specification's examples should be Prettier-clean everywhere (the Log example, `.prettierrc.json`).
- The skills: the two written here are a starting point for review.
