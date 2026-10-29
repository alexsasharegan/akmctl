---
id: "261027-ensf"
type: adr
adr_number: 19
title: "AKM: An Agent Knowledge Management Standard for Repositories"
status: accepted
created: 2026-10-27
updated: 2026-10-28
date: 2026-10-27
summary: "AKM is adopted as the standard for repository records: flat frontmatter, YYMMDD-xxxx ids, work items and ADRs in one id space, and the akmctl tool."
---

# ADR 261027-ensf: AKM: An Agent Knowledge Management Standard for Repositories

## Context

Alexander's repositories keep their own history: decisions, a roadmap, work in flight. The practice grew by hand across several agentic tools and drifted: the core workflow skill lives outside every repo and exists in three variants; ADR status words, numbering and fields differ; entry points are uneven; only one repo has CI. A cold agent, or Alexander after months away, cannot rely on the records. The spike [261026-m74k-agent-ready-repo-layout-spike](../work/done/261026-m74k-agent-ready-repo-layout-spike.md) compared Alexander's practice (Option A), an independent design built from the web by an isolated agent (Option B), and a mix (Option C), against a twelve-clause contract (orientation, rationale, direction, current state, boundaries, procedures, trust, cost, authority and ownership, verification, content safety and provenance, portability and handoff). The evidence is in [agent-ready-repo-discovery](../design/agent-ready-repo-discovery.md) and [agent-ready-repo-comparison](../design/agent-ready-repo-comparison.md); the specification is [agent-ready-repo-spec](../spec/agent-ready-repo-spec.md). The layout is meant for any maintainer or team, not only for Alexander.

**Lineage.** Alexander practises personal knowledge management (PKM): externalising memory into notes so the mind need not hold it, with capture, organisation and review. This standard is the same discipline moved inside a repository and addressed to agents. It is **agent knowledge management (AKM)**: PKM's idea, applied to the operational memory of a codebase. That lineage is the main reason for the name below.

## Name and category (working title)

- **The standard:** AKM, short for agent knowledge management. Said as letters, "A-K-M", the way PKM is. It is a methodology and a standard: a written spec with conformance, extensions and a migration path.
- **The tool:** `akmctl`, the control interface (the `ctl` convention of `kubectl` and `systemctl`, here over a repository's records). `akmctl ui` starts the workbench UI. It is one Go codebase in a separate repository.
- **Why this name:** the PKM lineage above; it names the discipline, not a feature; Alexander stands behind it as the most coherent option he found.
- **Status:** working title, not final. The name is a find-and-replace across the spec and tasks if it changes.
- **Overlap, accepted knowingly:** `itlackey/akm` ("Agent Knowledge Manager", a CLI for agent skills, scripts and memories; npm `akm-cli`); the AgenticAKM paper (2026), which uses AKM for Architecture Knowledge Management and generates ADRs from repositories; and AKM as an established architecture-knowledge research term. Different job from the first, but the same acronym, so search results and conversation can mix them up.
- **Alternatives considered:** Memenda, Jidoka, Andon, Memrail. Rejected: Cairn, Waymark, Waybill, Wayfinder, Trailmark, Switchyard, Roundhouse, Fingerpost (awkward), Agemo (a Yoruba deity and festival).

## Value proposition

**One line.** AKM gives a repository a memory and rails, so any agent can pick up where work stopped. As a descriptor: an opinionated standard for repositories with memory and rails.

**For whom.** A maintainer, solo or on a team, who works with AI coding agents.

**The problem.** Every agent session starts with no memory. Project history sits in chat logs, in people's heads and in docs that drift. The result is settled decisions argued again, mistakes repeated, work lost mid-task, parallel agents colliding, and tokens spent reading everything to find the one thing needed.

**The promise.** A repo that carries its own memory in plain files. A cold agent, or the maintainer after months away, can answer four questions in a few reads: what is this, why is it like this, what is planned, and where did we stop.

**Two headline differentiators.**

1. **Memory, not only orientation.** A good `AGENTS.md` tells an agent how to work in a repo. AKM does that and also keeps the repo's history: decisions and their reasons, direction, and work in flight.
2. **Rails for the work.** Work follows a simple Kanban-style flow (raw, ready, active, blocked, done, dropped) with one active item and a stated next step, and no ceremonies. Agents have a path to follow and people get a board that keeps itself current. The paperwork is the agent's; the decisions stay the maintainer's.

**Supporting principles.** What else the standard commits to. The mechanisms are in [agent-ready-repo-spec](../spec/agent-ready-repo-spec.md), not here.

3. **Conducive to parallel work.** Many agents and people can start work and record decisions at once with little to reconcile. Not a guarantee: shared files such as the roadmap can still conflict.
4. **Token-efficient by design.** A reader finds the right record by reading a little and loading detail on demand, instead of reading everything.
5. **Less drift.** Records are structured and checked mechanically, and each fact lives in one place, so they are more likely to stay accurate as work goes on. Checks catch structure, not truth.
6. **Works with any agent, with or without tooling.** Plain files are the source of truth. A bare agent can use the repo; the tool makes it faster without making it possible.
7. **Upgradable.** The standard is versioned, and a repo can move to a newer version along a defined, reviewable path.
8. **Open for extension.** Plain files are the source of truth and the standard is a written spec, so others can add extensions or alternative implementations without changing the core. Versions and the manifest keep each extension reconciled with the layout it needs.

**What is not shown.** The research found no study of what history layers do for agents, and one study found that context files can hurt task success. The Kanban-style flow echoes the pattern in Anthropic's long-running-agents post (one case study), but the claim that agents work better on it is a hypothesis. The value is unproven until pilots show it, and the guide should say so.

## Decision

Adopt Option C as the standard **AKM (agent knowledge management) repository layout**, specified in [agent-ready-repo-spec](../spec/agent-ready-repo-spec.md). In summary:

- **Records.** Work items and ADRs are Markdown files with simple flat frontmatter (the header layer), in one id space, in `docs/work/` (closed items in `docs/work/done/`) and `docs/adr/`.
- **Ids.** `YYMMDD-xxxx`: the UTC creation date plus four random characters. No shared counter, because counters cause merge conflicts. The id alone is the reference; the slug lives only in the filename. The tool mints ids.
- **Work items.** One file each, with `status` (`raw`, `ready`, `active`, `blocked`, `done`, `dropped`), an optional `parent` (a feature is itself an item), `blocked_by` (ids or free text), integer `priority` (higher is sooner, default 0), an overwritten `## Next` and an append-only `## Log`.
- **ADRs.** Statuses `proposed`, `accepted`, `rejected`, `deprecated`, `superseded`. A maintainer's commit accepts; agents only draft. The ADR index is generated from frontmatter, committed, and never edited by hand. There is no committed work index; the work listing is derived on demand.
- **Roadmap.** `docs/ROADMAP.md` holds hand-ordered lists only (`Now`, `Next`, `Later`, `Dropped`).
- **`AGENTS.md`.** Short (60 to 80 lines), five fixed parts, pointers out. Must-hold boundaries are also enforced by tool-level controls where available.
- **Project and tool.** The project is a **clone template** with instructions and opt-in default extensions. The tool and a workbench UI are one Go codebase in a separate repository, released as prebuilt binaries, with one subcommand set, a stable JSON output and a transport-neutral API spec with CLI and HTTP bindings. The tool is an accelerator, not a dependency: every operation has a plain-file fallback.
- **Versions.** A manifest `docs/.layout` records the layout and extension versions. Only MAJOR changes need a migration. `migrate` reconciles a repo: a conformance report, declarative safe steps (dry run by default, no deletes), and an agent handoff prompt with safeguards.

Rejected: keeping Option A as it is (the drift, and the skill outside version control); adopting Option B unchanged (untested, abandons working conventions); a shared counter or sequential numbers for ids; a hand-kept status board or roadmap prose; per-project folders as the unit of work; a committed work index; a Python script copied into each repo.

## Consequences

- Existing repos (`alexregan.com`, `poetic-overtones`, `shavian2tengwar`, `mysqlbak`) and the vault each need a migration note; none is forced. The vault's task notes differ in field names, naming and ids ([261027-nrwf-vault-agent-ready-migration](../work/done/261027-nrwf-vault-agent-ready-migration.md)).
- The Go tool, the API spec and the YAML library choice are separate tasks: [261027-r3hh-repo-records-go-tool](../work/done/261027-r3hh-repo-records-go-tool.md), [261027-xkd4-repo-records-api-spec](../work/done/261027-xkd4-repo-records-api-spec.md), [261027-wkxv-go-yaml-library-evaluation](../work/done/261027-wkxv-go-yaml-library-evaluation.md). Until the tool exists, the plain-file fallbacks are the only interface, and the layout is usable by hand.
- Trust still rests on agent discipline, because no hook or CI forces the check to run. A well-formed but wrong record passes every check.
- `pkm-workbench` ([261025-r0ye-pkm-workbench-a-local-server](261025-r0ye-pkm-workbench-a-local-server.md)) is expected to fold into, or be coupled to, the new tool. That will need a follow-up ADR; this one does not supersede [261025-r0ye-pkm-workbench-a-local-server](261025-r0ye-pkm-workbench-a-local-server.md).
- The formatter question (leaning to Prettier), whether `renumber` ships in the first version, and several naming choices are still open in the spec.

## Evidence

- Spike: [261026-m74k-agent-ready-repo-layout-spike](../work/done/261026-m74k-agent-ready-repo-layout-spike.md). Discovery: [agent-ready-repo-discovery](../design/agent-ready-repo-discovery.md). Comparison and decisions D1 to D8, O1 to O4: [agent-ready-repo-comparison](../design/agent-ready-repo-comparison.md). Specification: [agent-ready-repo-spec](../spec/agent-ready-repo-spec.md).
- Accepted by Alexander on 2026-10-28 (acceptance commit `dd814ce`, the rename from the drafts folder). The name is a working title.
- Commits: `fd2233d` (comparison), `7c6802a` (D1 to D8), `d2d6df2` (O1 to O4), `19c6ed1` (ids and requirements), `7c838d4` (Go tool direction), `982d28c` (fallback and sections 6, 7), `f31dc0d` (version manifest and migrate), `e3593db` (working title and value proposition).

## Relations

- Related to:
  - 261023-f7rd-task-queue-and-the-task-lifecycle-skill (task queue and the task-lifecycle skill: the work-item model follows its `## Next` and close rules)
  - 261023-mtpt-frontmatter-field-registry (frontmatter field registry: the layout adds fields that will need registering)
  - [261025-r0ye-pkm-workbench-a-local-server](261025-r0ye-pkm-workbench-a-local-server.md) (pkm-workbench: expected to merge with the new tool)

## Errata

- 2026-10-28: three corrections, from [akm-spec-review](../design/akm-spec-review.md): `dropped` added to the flow of statuses; the acceptance commit hash added to Evidence; an unsourced star count for `itlackey/akm` removed. The decision is unchanged.
- 2026-10-28: decisions made since acceptance, in [agent-ready-repo-decision-log](../design/agent-ready-repo-decision-log.md) (Y1 to Y20): Prettier is required (pinned, embedded formatting off); `renumber` is not in v1; the UI is `akmctl ui`; a small dependency is allowed for display width (an exception to "standard library only"); optional `claims` and `safety` extensions exist. They refine the spec, not this decision. The `authored_by` marker was removed.
- 2026-10-28: the display-width dependency mentioned above was dropped; the tool implements Prettier's width rule itself (Y28). Obsidian compatibility, which the decisions above rely on, is recorded in [261028-zh0w-obsidian-compatibility-is-a-stated-goal](261028-zh0w-obsidian-compatibility-is-a-stated-goal.md).
- 2026-10-29: copied from the maintainer's notes into this repository. Wikilinks to records copied with it became relative links; wikilinks to notes that stayed behind became plain text. The decision is unchanged.
