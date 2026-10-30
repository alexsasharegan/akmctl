---
id: 261027-nrwf
type: task
status: done
origin: 261026-m74k
created: 2026-10-27
updated: 2026-10-29
closed: 2026-10-29
authored_by: agent
title: Migrate the vault to the agent-ready layout
---

# Migrate the vault to the agent-ready layout

## Goal

Alexander expects to dogfood the agent-ready layout in the vault, which needs a migration and matching changes in the workbench. Known open item, not a decision yet.

Known differences, from [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md):

- The vault's task notes use `created_at`, `updated_at`, `project`, `origin` and `adr`, are named `YYYY-MM-DD slug`, and carry no `id`. The layout uses `created`, `updated`, `parent`, `decisions`, `YYMMDD-xxxx` ids and `blocked_by`.
- The vault is a notes repository, not a code repository, so `docs/work/` and `docs/adr/` would map to `Deep Thought/9 pkm/tasks/` and `adr/`.
- The workbench loads only `9 pkm/tasks/` with `type: task` and cannot read lists or integers today.

Done when: a mapping from the vault's fields and folders to the layout is written, a migration plan is agreed with Alexander, and the workbench changes it needs are listed.

## YAML-lite and the vault (Y26)

YAML-lite is aligned to YAML 1.2 core, the schema Obsidian reads and writes, so unquoted dates (`created: 2026-10-28`) are valid and no date quoting is needed. Known conflict (Y32, Y33: the standard keeps `id`, so the vault changes): `.obsidian/types.json` types `id` as number, and the layout's `id` is a string with letters; decide whether the vault keeps numbered series under another field name or retypes `id`. Still to check in the migration: checkbox and number properties (booleans and integers are accepted; floats are errors), Obsidian's rewrites (comments dropped, inline lists made block lists), and keys or values the grammar rejects. Also record the dogfooding facts: the vault's task notes use `created_at`/`updated_at`, `project`, `origin`, `adr`, a `YYYY-MM-DD slug` name and no `id`; the layout uses `created`/`updated`, `parent`, `decisions` and `YYMMDD-xxxx` ids; the vault is a notes repo, not a code repo, so its folders would differ from `docs/work/`.

## Next

The real run is done and merged into main (Y43). Left: Alexander looks at the Tasks board and the ADR index in Obsidian and opens one renamed task to see the links; then close this task. Open follow-ups: the roadmap is empty (45 RM4 warnings), akmctl work items from the dry run, the UTC id date, the Templater trigger on new files, pkm-workbench stays frozen, the two Shavian-containing files with their own ADR numbering (LifeNotes_ADR, Tagging Manifest) were not rewritten. Suggested mode: Query.

## Outcome

Done on 2026-10-28 (decision Y43). The vault's 124 tasks and 21 ADRs are AKM records: ids, titles, created/updated, a Goal scaffold on every task, closed tasks in tasks/done/, ADR summaries and a generated index. 145 renames ran through the Obsidian CLI and Obsidian repaired the links; verification found nothing lost. akmctl check: 0 errors, 47 warnings. Alexander checked the Tasks board and the ADR index in Obsidian. Left as follow-ups: fill the roadmap (45 RM4 warnings), akmctl work items from the dry run, the UTC id date, Templater's trigger on moved notes, two Shavian-containing files with their own ADR numbering, pkm-workbench stays frozen. ADR check: none.

## Log

- 2026-10-27: filed by Claude Code on the Mac at Alexander's direction, split out of [261026-m74k-agent-ready-repo-layout-spike](261026-m74k-agent-ready-repo-layout-spike.md). Depends on the layout spec and on the Go tool; do not start before them.
- 2026-10-28: [261026-m74k-agent-ready-repo-layout-spike](261026-m74k-agent-ready-repo-layout-spike.md) closed. [261027-ensf-akm-an-agent-knowledge-management](../../adr/261027-ensf-akm-an-agent-knowledge-management.md) is accepted and [akm-guide](../../spec/akm-guide.md) describes adoption; the migration method is to give an agent the spec and the target repo and ask for a plan.
- 2026-10-28: took over the dogfooding detail from the decision log (Y25) and added the quoted-dates consequence of YAML-lite (Y21).
- 2026-10-28: tested YAML-lite against Obsidian 1.14.5 ([agent-ready-repo-decision-log](../../design/agent-ready-repo-decision-log.md), section F). Obsidian rewrites dates unquoted on any property edit, so requiring quoted dates (Y21) breaks it; recommended change: align YAML-lite to YAML 1.2 core (unquoted dates are strings), pending the maintainer's decision. This task's quoted-dates problem disappears if that is accepted.
- 2026-10-28: the quoted-dates problem is gone: YAML-lite was aligned to YAML 1.2 core after the Obsidian test (Y26).
- 2026-10-28: [261027-r3hh-repo-records-go-tool](261027-r3hh-repo-records-go-tool.md) closed: `akmctl` v0.1.0 exists (private repo, release assets), with `check`, `index` and `migrate`. This migration can use `akmctl check` against the vault once it has a layout manifest; the id clash (`id` typed number in `types.json`) is still open.
- 2026-10-28: Alexander's order for the next work: [261028-j3m9-akm-spec-v1-cleanup](261028-j3m9-akm-spec-v1-cleanup.md) first, then this migration.
- 2026-10-28: order changed: [261028-j3m9-akm-spec-v1-cleanup](261028-j3m9-akm-spec-v1-cleanup.md), then [261028-eqwq-akmctl-dogfood-akm-layout](261028-eqwq-akmctl-dogfood-akm-layout.md), then this migration.
- 2026-10-28: queue order now: [261028-eqwq-akmctl-dogfood-akm-layout](261028-eqwq-akmctl-dogfood-akm-layout.md) and [261028-a149-akmctl-ui-port-pkm-workbench](261028-a149-akmctl-ui-port-pkm-workbench.md) first, then this migration.
- 2026-10-28: `mint_ids` is built in `akmctl` (work item 261028-8dd4, branch `mint-ids`, not merged): the migration can mint ids for existing notes with the date from `created`; plan reports old name, new name and id, and lists references it does not rewrite.
- 2026-10-28: shaped with Alexander (decision log Y41). Scope: tasks and ADRs. Every vault-specific task field is converted to the AKM form. Task files are renamed to `<id>-<slug>.md` and no link may be lost: either rename through Obsidian (it rewrites links) or a scripted rewrite. Alexander grants scoped approval to edit past periodic notes (daily, weekly, monthly, quarterly, yearly) ONLY to repair links to renamed notes: no other change to them. Next: a read-only planning agent writes the mapping and the plan; the first run is on a copy of the vault.
- 2026-10-28: Shavian scope for this migration (Alexander): editing Shavian text is allowed where needed to repair links (to renamed notes, or links inside a note being migrated). Do not read, translate or quote Shavian prose beyond the link text. Any Shavian file touched, and any Shavian task or ADR title swept up, is listed by path for his review. He expects Shavian notes are not migration targets.
- 2026-10-28: planning report received (read-only). Alexander decided (Y42): manifest keys map the folders in the vault; `project`/`origin`/`adr` stay as extra fields; links repaired by Obsidian renames after a pilot (47 links in 46 files, link-only changes, readable names lost, accepted). Pilot reverted; baseline tag `pre-pilot-rename`. Next: tool work in `akmctl` (manifest keys `work_dir`/`adr_dir`/`roadmap`; `title_from_h1`, `goal_scaffold`, `remove_field`, `copy_field` steps; `mint_ids` option to emit the rename map and write ids without renaming), then a full dry run on a copy of the vault, then the real run in steps. Open for Alexander before the real run: ADR `date` source, the ADR summaries, Goal text for 115 tasks, the three active tasks' `## Next`, AGENTS.md restructure.
- 2026-10-28: Alexander's answers for the real run (Y43 candidates, to log when the run is planned): ADR `date` comes from `created_at`; the 115 tasks without `## Goal` get the mechanical scaffold, and he reviews only the 45 open ones; the vault's AGENTS.md is restructured to the five-part contract under 80 lines. Still open: the 21 ADR summaries (agent drafts, he reviews) and `## Next` for the three active tasks without one.
- 2026-10-28: more answers from Alexander: an agent drafts `## Next` for the three active tasks and he checks them; bare prose mentions of ADR numbers (about 94, for example "ADR 261027-ensf") are rewritten to the new id form ("ADR <id>" per the spec); the ADR conventions and sweep procedure move to `design/adr-conventions.md`, `adr/README.md` becomes the generated index, and `ADR Candidates.md` moves beside the conventions note. Boundary I am applying until he says otherwise: the prose rewrite does not touch past periodic notes (6 mentions in 3 calendar files), `4 RESOURCES` or `5 ARCHIVES`, because his grant covers link repair only; those mentions are listed for his decision.
- 2026-10-28: tool work done in `akmctl` (work item 261029-9aa2, branch `vault-paths`, not merged): manifest keys `work_dir`, `adr_dir`, `roadmap` (rule S3 for bad values); steps `title_from_h1`, `goal_scaffold`, `remove_field`, `copy_field`, `dir` on the field steps; `mint_ids` `rename: false` with `renames` and `seed`; `migrate --migration <file>` and the `migrations` parameter; `plan` now simulates the steps in order. Spec, guide, API spec, schemas and fixtures updated (rule S3, section 6 Settings, API 4.15). ADRs with a numeric `id` run `rename_field` `id` to `adr_number` before `mint_ids`. Next: the full dry run on a copy of the vault.
- 2026-10-28: dry run done on a copy (scratchpad, branch `migrate-0.0.0-to-1.0.0`): `akmctl check` 0 errors, 47 warnings; link verification passed; 145 renames. Findings, timings and the real-run runbook: [vault-migration-runbook](../../design/vault-migration/vault-migration-runbook.md); drafts to review (ADR summaries, Next, Goal scaffolds, AGENTS.md): [vault-migration-drafts](../../design/vault-migration/vault-migration-drafts.md); scripts in `9 scripts/akm-migration/`. Decisions needed before the real run are listed in the runbook.
- 2026-10-28: dry run done (runbook and drafts in `design/vault-migration-runbook.md` and `vault-migration-drafts.md`; scripts in `9 scripts/akm-migration/`). Alexander's answers: `id` is retyped as text vault-wide (Field Notes keep their numeric ids as text); ADR prose mentions are also rewritten in the 14 Shavian-containing files (each file is listed first and checked for its own numbering scheme before the rewrite; the AREAS/AMEX series stays untouched); the 21 ADR file names are reviewed by Alexander before the real run.
- 2026-10-28: Alexander approved the 21 ADR names as generated. Field Notes check before retyping `id` as text: the Field Notes base sorts by `id` DESC; Obsidian reads all 121 Field Notes ids as numbers from the YAML (checked read-only with the metadata cache), so the existing sort is unaffected by the `types.json` change, which only changes how the Properties panel types the field. Risk: a new Field Note created or edited through the Properties panel may store `id` as a string and then sort among strings; hardening if it shows up: sort the base by a formula `number(id)`.
