---
type: design
status: draft
created_at: 2026-10-28
updated_at: 2026-10-28
authored_by: agent
tags:
    - pkm/rails
---
# Vault migration: drafts for Alexander's review

Drafted by Claude Code during the dry run ([[vault-migration-runbook]]). Nothing here is applied to the real vault. The scripts that apply them are in `Deep Thought/9 scripts/akm-migration/`; the dry-run copy is `/private/tmp/claude-501/-Users-deepthought-obsidian-vaults/dfafa885-dc8d-435b-9a5a-ca08fd51f7df/scratchpad/copy` (branch `migrate-0.0.0-to-1.0.0`).

## 1. ADR summaries (21)

One sentence each, at most 160 characters, from each ADR's Decision section. The ids are what the seeded run produced in the dry run; they repeat in the real run if the ADR dates are unchanged. Source: `adr_summaries.json`.

| Legacy | Id | Title | Draft summary |
| --- | --- | --- | --- |
| ADR 261019-b423 | 261019-b423 | Record Vault Decisions as ADRs | Vault decisions are recorded as short ADRs in 9 pkm/adr, drafted first, with past decisions backfilled on a best-effort basis. |
| ADR 260912-3q6j | 260912-3q6j | Periodic Review Skill Design | Periodic review skills read one level down (annual from quarterly, weekly from daily) and write terse, structured notes under 1,000 words. |
| ADR 260912-extg | 260912-extg | Periodic Review Skill Version 1.0.0 | Periodic review artifacts carry semver versions, v1.0.0 is the current state, and each major bump needs its own ADR. |
| ADR 260601-wvft | 260601-wvft | Calendar Layout: Periodic Notes Paths under 1 Calendar | Periodic notes live under 1 Calendar with year and month nesting, three-letter month folders and shared calendar keys; core Daily Notes is off. |
| ADR 261014-eq0g | 261014-eq0g | Periodic Review Skill v2.0.0: Deterministic Gather and Apply Scripts | The weekly review keeps synthesis with the agent and moves the mechanical gather and apply steps into weekly_review.py; Sleep becomes a section. |
| ADR 260531-q4ts | 260531-q4ts | Transition to Tiered Atomic Multi-Interval Architecture | History logging and state serialization are split: git-autosave runs five times a day under launchd and the git bundle runs monthly after a gc. |
| ADR 260531-eh31 | 260531-eh31 | Two-Stage Materialized View Ingestion Pipeline | Weekend analytics run as two stages: a compiler turn builds a materialized-view index of the week, then a reasoning turn works from that index. |
| ADR 260707-a27j | 260707-a27j | Incident Report: Speed Layer Failure & Cloud Replication Blackout | Incident report for 2026-07-02 to 07-07: git-autosave and the cloud syncs failed on a system git shim and a reset Full Disk Access, with the fixes. |
| ADR 260528-qabh | 260528-qabh | Git Autosave Scheduling: cron, then launchd | git-autosave moved from hourly cron to once a day, then to launchd after macOS updates kept resetting cron's Full Disk Access. |
| ADR 241103-b2mb | 241103-b2mb | Git for History Only: No Remote, Cloud Copies by rsync | Git tracks the vault for history only with no remote; copies reach Google Drive and iCloud by rsync every 15 minutes, and the work laptop is not synced. |
| ADR 260912-bcdr | 260912-bcdr | Obsidian CLI as the Agent Interface | Agents use the Obsidian CLI for anything Obsidian knows how to do (links, search, properties, moves) and plain file edits for note text. |
| ADR 261020-5be4 | 261020-5be4 | Base and Hub: Every Database Gets a Hub Note | Every database is a base plus a hub note, with records of one type in a subfolder and the hub beside it, so the dashboard is reachable without opening it. |
| ADR 261023-f7rd | 261023-f7rd | Task Queue and the task-lifecycle Skill | Tasks are notes with type: task under 9 pkm/tasks, carried through their life by the task-lifecycle skill; bookkeeping is allowed in any mode. |
| ADR 261023-xp6v | 261023-xp6v | Notebook Tracking by Tag and Chain | Notebooks are tracked by an opt-in #notebook/daily tag in the daily note and a previous property that links each notebook to the one before. |
| ADR 261023-mtpt | 261023-mtpt | Frontmatter Field Registry and Conformance Test | The frontmatter field registry is the authority for field names and types, and a conformance test checks every note against it. |
| ADR 261023-890k | 261023-890k | Review Skills after v2.0: Shavian, Sleep and Notebooks | Review skills v2.1 to v2.3: a Self group, Shavian not read in reviews, stricter sleep rules, a Phoenix-time default and notebook gathering. |
| ADR 261023-n0sv | 261023-n0sv | MOC Type, Style Guide and the Projects Index | type: moc marks hand-curated maps (style guide and moc-builder skill); for projects the folder in 2 PROJECTS is the index, with no base or status field. |
| ADR 261023-nt8w | 261023-nt8w | Backup Sync with Homebrew rsync | Backup syncs use Homebrew rsync with --whole-file --delete, so cloud placeholder files cannot break a run; the launchd plists are committed. |
| ADR 261025-r0ye | 261025-r0ye | pkm-workbench: A Local Server over the Vault's Files | pkm-workbench is a local Go server in its own repo that serves task views from the vault's files and writes changes back to the frontmatter. |
| ADR 261027-ensf | 261027-ensf | AKM: An Agent Knowledge Management Standard for Repositories | AKM is adopted as the standard for repository records: flat frontmatter, YYMMDD-xxxx ids, work items and ADRs in one id space, and the akmctl tool. |
| ADR 261028-zh0w | 261028-zh0w | Obsidian Compatibility Is a Stated Goal of the AKM Standard | Obsidian compatibility is a stated goal of AKM and the tiebreaker in design conflicts, though Obsidian itself is not required. |

## 2. `## Next` for the three active tasks

Only one of the three lacked a `## Next`: [[261025-sja2-reader-tag-to-summary-cards]] and [[261025-yx6c-daily-scripture-spike]] already have one, and `akmctl check` (WK6) passes on them. Draft for the third, [[261022-sbv7-aei-comp-dataset]]:

> Fill the open cells in the comps dataset: the contact form for every row (needs page HTML, so it waits for the archivist script, work item 261022-pd5y), L&T's URL count, and the RFQ form for L&T, Guzman and PMF. Leave Precision Advanced Manufacturing as "blocks automated fetches" unless a manual check is wanted. Then write the steps that regenerate the table and log them. Suggested mode: Plan first (decide whether to check the three RFQ forms by hand while the script is unbuilt), then Build.

## 3. Goal scaffolds to review (open tasks)

`goal_scaffold` gave all 124 task notes a `## Goal` (111 with the old introduction moved under it, 13 also with `## Done when` kept as a list), not 115. The 45 open ones (ready 42, active 3) follow. Flags: SHORT (under 80 characters), LONG (over 1500), NO-DONE-WHEN (no "done when" text), EMPTY. Six have a flag other than NO-DONE-WHEN. Ids are those of the dry run.

| id | status | title | Goal chars | flags |
| --- | --- | --- | --- | --- |
| 261020-2e3s | ready | Start a backyard compost project | 618 | - |
| 261020-4p2s | ready | Archive the general election project | 355 | - |
| 261020-7397 | ready | Spike: decide how to uplift the Cycles log | 1222 | NO-DONE-WHEN |
| 261020-eppq | ready | Convert Appendix B candidates into skills | 77 | SHORT NO-DONE-WHEN |
| 261020-nq5n | ready | Draft an implementation-free philosophy doc for the review cycle | 561 | NO-DONE-WHEN |
| 261020-xhdm | ready | Add a rails review to the monthly review | 160 | NO-DONE-WHEN |
| 261021-56bh | ready | Document the pmx-01 database backup job | 823 | - |
| 261021-753k | ready | Retrofit Project List as a MOC with commitment tiers | 761 | - |
| 261021-bkf9 | ready | Scan 2 PROJECTS for drift and staleness, to rotate stale projects to archives | 777 | - |
| 261021-m9xm | ready | Move three non-projects out of 2 PROJECTS | 681 | - |
| 261021-mvt7 | ready | Spike: decide how to manage pmx-01's server info and maintenance log | 2365 | LONG NO-DONE-WHEN |
| 261022-1js5 | ready | Confirm who pays Alexander, and the W-9 and 1099 path | 700 | - |
| 261022-9pgc | ready | Build the AEI Stage 1 proposal | 978 | - |
| 261022-kmbw | ready | Inventory AEI's hosting, domains and registrar access | 831 | - |
| 261022-pd5y | ready | Write a script that archives a site from its sitemap | 914 | - |
| 261022-sbv7 | active | Turn the peer comparison into a dataset | 1254 | - |
| 261022-xefa | ready | Review the scope draft, then agree it with Landis and Ron | 467 | - |
| 261022-zdgm | ready | Collect everyone's contact method | 402 | - |
| 261023-h8dg | ready | Find a home for every file in 0 Inbox, or delete it | 1378 | - |
| 261024-52vn | ready | Fork the kit, own the code, reskin it, add PDF and PNG output | 243 | NO-DONE-WHEN |
| 261024-7vhc | ready | Write the vault-views skill | 704 | - |
| 261024-vaq2 | ready | Home Depot run: R20 bulbs, dimmer knob, AA batteries, pendant bulb | 1884 | LONG |
| 261025-35br | ready | Index Chrysostom's homilies by scripture range | 1285 | - |
| 261025-4nb8 | ready | Spike: a sturdier way to log and plan travel | 851 | - |
| 261025-8wv1 | ready | Migrate the existing Common Place Patristics quotes to the quote format | 1180 | - |
| 261025-9g5w | ready | Decide where blog writing lives and how it flows | 1046 | - |
| 261025-9hjv | ready | Organize 5 ARCHIVES | 1749 | LONG |
| 261025-9nvs | ready | Spike: capture political values piecemeal and keep them current | 800 | - |
| 261025-b2fd | ready | Check out superfile | 1058 | - |
| 261025-eq8z | ready | Rename the 22 pilot video cards to the `(video)` suffix | 900 | - |
| 261025-fct0 | ready | Decide whether the vault needs a general template skill | 902 | - |
| 261025-hqcm | ready | Brainstorm uses for the Hunter's Moon field notes | 515 | - |
| 261025-sja2 | active | Turn Readwise Reader items tagged `vault` into summary cards | 488 | NO-DONE-WHEN |
| 261025-yx6c | active | Spike: a daily scripture project with the GOARCH readings | 2066 | LONG |
| 261026-3c62 | ready | Clean the downstairs and master bath sink drains, and start a log of when it was done | 1387 | - |
| 261026-400a | ready | Record date and time on task notes | 659 | - |
| 261026-79zq | ready | Look at the master shower drain, and log it | 715 | - |
| 261026-mb2a | ready | Move the view spec into the Tasks hub | 584 | - |
| 261026-mrdr | ready | pkm-workbench: sort options | 435 | - |
| 261026-np3c | ready | pkm-workbench: delete a task, with its backlinks listed | 394 | - |
| 261026-tesh | ready | A rare priority flag on tasks | 754 | - |
| 261026-wnx3 | ready | Track down the build pipeline and hosting for parametricstudios.io | 741 | - |
| 261026-zg0g | ready | pkm-workbench: filters | 461 | - |
| 261027-82pg | ready | Spike: revive the dotfiles project and scan for drift | 1852 | LONG |
| 261027-nrwf | ready | Migrate the vault to the agent-ready layout | 869 | - |


## 4. Draft `AGENTS.md` (five parts, 42 lines)

On the copy, `akmctl check` reports no G1, G2 or G3 finding. G3 only resolves `scripts/...` paths and make or npm targets, so it did not test the conformance path in Commands. The text it replaces (from "User & Interaction" on) moves unchanged to a new note `Deep Thought/9 pkm/vault-conventions.md` (type reference), which the draft links.

```markdown
# AGENTS.md

Alexander's personal knowledge vault (an Obsidian vault in `Deep Thought/`, tracked by git at this root), worked by him and by agents under the rules below.

## Commands

- Check the vault: `python3 "Deep Thought/9 scripts/frontmatter/conformance.py"` (frontmatter against the field registry) and `akmctl check` (task and ADR records; run it before committing records).
- Regenerate the ADR index after any ADR change: `akmctl index`.
- Obsidian CLI (needs the app running): `obsidian --vault="Deep Thought" <command>`; see `.agents/skills/obsidian/SKILL.md`.
- Commit only the files you changed, staged by path: `git add -- <path>`, then `git commit`. No remote: nothing to pull or push.

## Where things live

- Vault: `Deep Thought/` (PARA: `0 Inbox`, `1 Calendar`, `2 PROJECTS`, `3 AREAS`, `4 RESOURCES`, `5 ARCHIVES`, `8 assets`, `9 pkm`, `9 scripts`, `9 templates`). Full layout, Obsidian settings and periodic-note paths: [[vault-conventions]].
- Who Alexander is and how to talk to him: `Deep Thought/9 pkm/me.md` (read it at the start of a session; the `SessionStart` hook injects it in Claude Code).
- Open work: `akmctl ls --status ready,active,blocked`. Records: `Deep Thought/9 pkm/tasks/` (closed in `done/`). Procedure: `Deep Thought/9 pkm/Tasks.md` and the `task-lifecycle` skill.
- Decisions: `Deep Thought/9 pkm/adr/README.md` (generated index); conventions in `Deep Thought/9 pkm/design/adr-conventions.md`.
- Roadmap: `Deep Thought/9 pkm/design/akm-roadmap.md`. Design and runbooks: `Deep Thought/9 pkm/design/`, `Deep Thought/9 pkm/runbooks/`. Hooks: `.claude/hooks/` (see the runbook "Agent Hooks").
- Skills: `.agents/skills/<name>/SKILL.md`; index in `Deep Thought/9 pkm/skills-map.md`. Read a skill when its description fits the work.
- Code lives in clones outside the vault; `Deep Thought/9 pkm/Repos.md` maps the names Alexander uses to repos.

## Rules

- Default mode is Query (read only). Plan / Design discusses and proposes; Build edits only within the approved scope. Task bookkeeping (status, `updated`, `## Next`, `## Log`, `## Outcome` of the task you work) is allowed in any mode.
- List the open work before starting; resume the `active` task you hold, else take the top `ready` one; if several are unclaimed, ask. Work only under a task note.
- A hard-to-reverse choice needs an ADR draft with status `proposed`. Never edit the body of an `accepted` ADR; supersede it.
- Update the task's `## Next` and `## Log` in the same commit as the work. Run `akmctl check` before committing records.
- Links are Obsidian wikilinks (`[[File name]]`); rename notes with the Obsidian CLI so links are repaired. Escape a bracket that opens prose: `\[Photo: ...]`.
- A note an agent creates carries `authored_by: agent` in its frontmatter; frontmatter fields come from `Deep Thought/9 pkm/design/frontmatter-fields.md`.
- New unclassified notes go to `0 Inbox/`. Check `Deep Thought/3 AREAS/LifeOS/Names & Spellings.md` before guessing a name (dictation).
- Modify only the files the prompt names or the operation needs. No vault-wide formatting, link or metadata passes.
- Code and tests win over prose; when a note is wrong, fix it or say so.
- Without `akmctl`, use `rg`/`grep` on the task frontmatter (`rg -l '^status: active' "Deep Thought/9 pkm/tasks"`) and tell Alexander how to install it.
- Mint an id: `id="$(date -u +%y%m%d)-$(LC_ALL=C tr -dc '0-9a-hjkmnp-tv-z' </dev/urandom | head -c 4)"`, then check no file in the task, done and ADR folders starts with `$id-`.

## Do not, without asking

- Read, translate, quote or gloss Shavian in Alexander's private writing (daily logs, field notes); treat it as opaque. Exceptions: [[vault-conventions]].
- Edit past periodic notes (daily, weekly, monthly, quarterly, yearly) before the current period, or anything in `4 RESOURCES/` or `5 ARCHIVES/`.
- Push, create a remote, or change the autosave job (`com.deepthought.git-autosave`); the autosave commits about five times a day.
- Delete a note, change `.obsidian/` settings, or change a skill, hook or template outside a task that says so.
- Put credentials, tokens or personal identifiers into notes, commits or scripts.
```
