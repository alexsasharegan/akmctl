---
id: 261026-m74k
type: task
status: done
origin: 261025-g68j
created: 2026-10-26
updated: 2026-10-28
authored_by: agent
adr: recorded
title: "Spike: an agent-ready repository layout"
---

# Spike: an agent-ready repository layout

## Goal

Alexander wants a guide to what he calls an **agent-ready repo**: a repository that keeps its own history (decisions, direction, work in flight) so that an agent with no memory, or Alexander after months away, can work well in it. The guide covers layout, procedures and skills.

Inputs so far, with web sources and a first pass at both options: [agent-ready-repo-layout](../../design/agent-ready-repo-layout.md). Its leanings are inputs only; the options below are re-derived in isolation.

Done when:

The contract is settled; Option A, Option B and the chosen layout are each described at comparable depth and mapped to it; the decision is recorded (layout ADR); the guide exists in the vault; the layout spec is complete and reviewed (sections 6 and 7 rewritten for the Go tool); the project has a category word, a name and a value proposition, and the layout ADR carries them. Migration notes for the other repos are not part of this spike: once the spec is complete, an agent can be given the spec and a target repo and produce a migration plan (decided 2026-10-27). The Go tool, the API spec, the YAML library choice and the vault migration are separate tasks filed from this one: [261027-xkd4-repo-records-api-spec](261027-xkd4-repo-records-api-spec.md), [261027-r3hh-repo-records-go-tool](261027-r3hh-repo-records-go-tool.md), [261027-wkxv-go-yaml-library-evaluation](261027-wkxv-go-yaml-library-evaluation.md), [261027-nrwf-vault-agent-ready-migration](261027-nrwf-vault-agent-ready-migration.md). The pilot on `pkm-workbench` moves with the tool, since the workbench is expected to fold into it or be coupled to it.

## Contract

Before comparing layouts, state what a repo's memory is _for_. Each mechanism (a docs folder, ADRs, a roadmap, a skill) has to earn its place by serving a clause here. C1 to C8 are a draft from Alexander's words; C9 to C12 were added after the discovery pass (2026-10-26, grouping accepted by Alexander). The spike settles the whole.

An agent-ready repo lets a cold reader:

- **C1 Orientation:** learn what the project is, how to build, test and run it, and where things live.
- **C2 Rationale:** learn why it is the way it is (decisions, rejected alternatives), so past mistakes and settled decisions are not repeated or re-litigated.
- **C3 Direction:** learn what is planned, in what order, and what was dropped.
- **C4 Current state:** learn what is in flight, where it stands and the single next step, so work is resumable.
- **C5 Boundaries:** learn what must not be touched or done without asking.
- **C6 Procedures:** repeat recurring work the same way each time.
- **C7 Trust:** the records stay accurate over time; drift is prevented or detected, and it is clear what may be edited and by whom.
- **C8 Cost:** cheap to read cold (progressive disclosure), cheap to keep up.
- **C9 Authority and ownership:** a declared source of truth when documents conflict; who decides (agents propose, a maintainer accepts); a human review surface.
- **C10 Verification:** records say what was and was not verified; nothing is "done" because code exists.
- **C11 Content safety and provenance:** what a record may never contain; which records an agent wrote and which a maintainer reviewed.
- **C12 Portability and handoff:** works across tools and machines without relying on `~/.claude`; modes, claiming and concurrent threads are handled.

Each option, and the chosen mix, is described by how it meets each clause (C1 to C12), how strongly, and at what cost. Clauses may be added, merged or dropped; that change is itself a decision to record.

## Options

- **Option A:** Alexander's practice across his repos.
- **Option B:** an independent design from web sources and original ideas.
- **Option C:** a mix, proposed after A and B are on the table.

## Steps

1. **Discovery, in isolation.** Two read-only agents with controlled context, so B is not coloured by A.
   - Option A agent: reads `alexregan.com`, `poetic-overtones`, `shavian2tengwar`, `mysqlbak` (under `~/github.com/alexsasharegan/`), the user-level `project-docs` skill and templates, the vault as a repo case, and `pkm-workbench` as the pilot target. Describes, does not advocate. Records files, procedures, skills, hooks, CI, purpose, inconsistencies and real drift, then maps them to C1 to C8.
   - Option B agent: no access to local files. Researches ADR practice, roadmaps, plans directories, agent skills, `AGENTS.md` evidence and critiques. Proposes one full design: layout, memory layers, issues, `AGENTS.md` text, draft skills, mechanical checks, lifecycle, small and large tiers, and a C1 to C8 mapping.
   - Out of scope: `parametricstudios.com` (excluded), and any work or employer repo, which Alexander cannot share and agents never request.
2. **Review the contract** with Alexander against both reports, including purposes either agent found that the draft misses.
3. **Compare** A, B and a mix C against the settled contract, plus cold-start cost, enforcement, upkeep cost, drift, fit for small and large repos, and fit with the vault's task notes. Show the options; do not decide for him.
4. **Decide** with Alexander, and record it as an ADR if it changes how repos are set up.
5. **Pilot** on `pkm-workbench` (it has `README.md`, `ROADMAP.md` and `docs/style-brief.md`), then write the guide in the vault.

## Outcome

AKM (agent knowledge management), a standard for repositories that keep their own memory and run on a Kanban-style flow, was specified, decided and written up. The contract (C1 to C12) is settled, Options A (existing practice), B (independent design) and a mix C were compared, and Option C was chosen: [261027-ensf-akm-an-agent-knowledge-management](../../adr/261027-ensf-akm-an-agent-knowledge-management.md) (accepted 2026-10-28), the specification [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md), the discovery and comparison notes [agent-ready-repo-discovery](../../design/agent-ready-repo-discovery.md) and [agent-ready-repo-comparison](../../design/agent-ready-repo-comparison.md), and the guide [akm-guide](../../spec/akm-guide.md). The name AKM and `akmctl` is a working title.

Why: [261027-ensf-akm-an-agent-knowledge-management](../../adr/261027-ensf-akm-an-agent-knowledge-management.md) carries the reasoning and the rejected alternatives.

Not done, by decision: migration notes for the other repos (an agent given the spec and a target repo can plan a migration), and the pilot on `pkm-workbench` (it moves with the tool). Follow-ups filed: [261027-xkd4-repo-records-api-spec](261027-xkd4-repo-records-api-spec.md), [261027-r3hh-repo-records-go-tool](261027-r3hh-repo-records-go-tool.md), [261027-wkxv-go-yaml-library-evaluation](261027-wkxv-go-yaml-library-evaluation.md), [261027-nrwf-vault-agent-ready-migration](261027-nrwf-vault-agent-ready-migration.md), [261028-1z36-akm-spec-open-questions](261028-1z36-akm-spec-open-questions.md).

Commits: `fd2233d` (comparison), `7c6802a` (D1 to D8), `d2d6df2` (O1 to O4), `19c6ed1` (ids and requirements), `7c838d4` (Go tool direction), `982d28c` (fallback, sections 6 and 7), `f31dc0d` (version manifest and migrate), `e3593db` (working title and value proposition), `dd814ce` (ADR 261027-ensf accepted), `0436c5b` (guide).

## Log

- 2026-10-26: filed by Claude Code on the Mac at Alexander's direction, with the design note [agent-ready-repo-layout](../../design/agent-ready-repo-layout.md). `inbox`, not claimed.
- 2026-10-26: claimed (Claude Code on the Mac).
- 2026-10-26: Alexander reshaped the task. (1) Work repos are not available: employer confidentiality, never asked of him; the task's "ask for them" was a mistake. (2) Add a contract: a general statement of why a repo keeps its own history, against which A, B and C are each measured. (3) Option A and Option B are investigated by separate isolated agents so B is not coloured by A. Task body rewritten in Build the same day.
- 2026-10-26: Alexander granted Build to rewrite the task body. Body now carries the contract, the three options and the isolated-discovery step; work repos removed.
- 2026-10-26: Step 1 done. Option A and Option B agents reported; condensed in [agent-ready-repo-discovery](../../design/agent-ready-repo-discovery.md). Both independently flagged verification, provenance, portability and sensitive content as missing from the contract.
- 2026-10-26: Alexander accepted the grouping of added clauses; contract is now C1 to C12 in this note. Comparison of A, B and a proposed C drafted in [agent-ready-repo-comparison](../../design/agent-ready-repo-comparison.md); decision points D1 to D8 left for him.
- 2026-10-27: Alexander corrected the comparison's caveat: A's drift comes from his unaided effort across tools and is not a verdict on A's design; it shows what a spec must prevent. Note edited. Interviewing him on D1 to D8 one at a time.
- 2026-10-27: Alexander decided D1 to D8 by interview; recorded in [agent-ready-repo-comparison](../../design/agent-ready-repo-comparison.md). Key: one work-item file with a `parent` field (feature = a work item), overwritten Next, derived listing plus hand roadmap, shared skill copied into repos with its check script, owner-commit ADR acceptance, short AGENTS.md, `docs/adr` and `docs/work`.
- 2026-10-27: Alexander asked for generalised wording (any maintainer or team) and settled O1 to O4: `blocked` + `blocked_by`; a parent closes only when all children are done or dropped; `priority` is an integer, higher sooner, default 0; template repo with versioned copies. Recorded in [agent-ready-repo-comparison](../../design/agent-ready-repo-comparison.md).
- 2026-10-27: drafted [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md) (schema and check rules) from D1 to D8 and O1 to O4. Assumptions X1 to X8 flagged for review. Templates, skill and script not yet built.
- 2026-10-27: Alexander settled X1 to X8. X2 changed to date-prefixed string ids for work items (ADRs stay sequential). Spec updated in [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md).
- 2026-10-27: Alexander set requirements for ids (one id space for ADRs and work items, no counter because of merge conflicts, fixed width, date visible, slug out of the id) and chose P1 (`YYMMDD-xxxx`, UTC). He asked that the requirements be in the spec, and that the ADR index be a generated routing table from frontmatter. Spec rewritten accordingly; X9 and X10 are my additions, not confirmed.
- 2026-10-27: Alexander accepted X9 and X10, and asked about the formatter, the work listing and a Kanban view over the spec. Open items recorded in [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md).
- 2026-10-27: Alexander: no committed work index; JSON listing is a must and should grow into a query/Kanban CLI over ADRs and work items (filters, sort, search) to cut agent token use; leaning to require Prettier. pkm-workbench to be summarised by a subagent. Recorded in [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md).
- 2026-10-27: pkm-workbench summarised by a subagent (Go, stdlib only; statuses match; parser cannot read lists; no JSON listing). Findings and the parser-location question recorded in [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md).
- 2026-10-27: Alexander chose Go over Python (one shared codebase for a CLI and the workbench UI, folding in or extending pkm-workbench), one tool with subcommands, and a transport-neutral API spec with a stable JSON response and CLI and HTTP bindings. Names deferred. Recorded in [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md); X8 changed.
- 2026-10-27: Alexander agreed the split. Filed [261027-xkd4-repo-records-api-spec](261027-xkd4-repo-records-api-spec.md), [261027-r3hh-repo-records-go-tool](261027-r3hh-repo-records-go-tool.md), [261027-wkxv-go-yaml-library-evaluation](261027-wkxv-go-yaml-library-evaluation.md) and [261027-nrwf-vault-agent-ready-migration](261027-nrwf-vault-agent-ready-migration.md). Decided: the project is a clone template with opt-in extensions; the tool and UI are a separate repo with prebuilt binaries; Go stdlib preferred, YAML library by due diligence; vault migration is a known open item.
- 2026-10-27: Alexander confirmed the fallback: the tool is an accelerator, not a dependency (H6). Spec sections 2, 6 and 7 rewritten for the Go tool and the clone template; X11 added for confirmation.
- 2026-10-27: Alexander accepted widening X11 into a version manifest and asked for a migrate proposal (codemod-style safe steps plus a structured report and an agent handoff with safeguards). Recorded as a proposal in [agent-ready-repo-spec](../../spec/agent-ready-repo-spec.md); the API spec and tool tasks note the scope.
- 2026-10-27: Alexander accepted the migrate design. Spec updated. Drafted the layout ADR as `proposed` in `9 pkm/adr/drafts/DRAFT - Agent-Ready Repository Layout.md`; acceptance and numbering are Alexander's.
- 2026-10-27: Alexander dropped migration notes from the spike (an agent given the spec and a target repo can plan a migration). Asked to settle what the thing is, what it is called and its value proposition, to update the ADR, then to write the guide.
- 2026-10-27: Alexander reviewed the value proposition: principles must be philosophical, mechanisms kept out. Rewrote the six principles and added them to the draft ADR. Name and category pending his thinking.
- 2026-10-27: Alexander added principle 7 (open for extension) and asked for trail-metaphor names and an elevator pitch. Searched GitHub and the web (one search each; absence is weak evidence): Cairn collides in the same domain (agentcairn, project-cairn, oritera/Cairn, cairn-dev/cairn); Waymark has several agent-tool collisions; Trailmark (Trail of Bits and another); Wayfinder (Laravel); Fingerpost has no real project on GitHub. Spec Kit (GitHub) and OpenSpec (Fission AI) are spec-driven development workflows, read from search summaries only.
- 2026-10-28: Alexander: pitch 1 is the base; second differentiator is rails (a Kanban-style flow agents follow, with the bookkeeping done by agents). Rail-themed names checked on GitHub: Waybill (wardmos/waybill, agent handoff), Switchyard (NVIDIA), Roundhouse (migration tool) all collide; Railhead is mostly free but reads as Ruby on Rails. Fingerpost still free.
- 2026-10-28: Alexander rejected Fingerpost as awkward and asked for NeMo-style blended names from our vocabulary. NeMo = Neural Modules (NVIDIA docs). Checked blends on GitHub names only: Memrail (tiny same-space repo, 0 stars), Lograil (tiny), Merail, Recrail, Kanmem free; Rememo, Repomem, Agrail crowded. Lean: Memrail. Go module, npm, domain and trademark checks not done.
- 2026-10-28: second naming round from Alexander's roots (history, memory, agent, model, kanban, board, Greek or Latin). GitHub name checks: free or near-free: Memenda (memoranda + agenda), Memkan, Histell, Hismem, Mnemod, Kanmo; taken or crowded: Kanbo (Kanboard), Mello, Pinax (Django), Pinakes (two small tools), Stele (a few small), Agemem, Agendum (small). Dropped Agemo: it is a Yoruba deity and festival. Idea worth keeping: agenda (things to be done) and memoranda (things to be remembered) map to rails and memory.
- 2026-10-28: Alexander brought `akmctl` / "AKM Standard" (Agent Knowledge Management) from another model, which claimed it was uncontested. Checked: the exact string `akmctl` is free, but the acronym AKM is taken in the same space by itlackey/akm ("Agent Knowledge Manager", a CLI for agent knowledge and memories, npm akm-cli, ~55 to 67 stars, active) and by the AgenticAKM paper (Architecture Knowledge Management, generating ADRs). Recommended against. Name still open; leads: Memenda, Jidoka, Andon, Memrail.
- 2026-10-28: Alexander argued for AKM / `akmctl` as a working title (PKM lineage; collisions are acceptable if known; `ctl` fits a control plane; stars are not a veto). Conceded my overstated points; kept the overlap as a recorded risk. Awaiting his word to record it as the working title.
- 2026-10-28: Alexander confirmed AKM and `akmctl` as the working title, with the PKM-to-AKM lineage in the ADR (the main reason for the name). Recorded in the draft ADR (context, name section, restructured value proposition), the spec (`wb` placeholders replaced by `akmctl`) and [261027-r3hh-repo-records-go-tool](261027-r3hh-repo-records-go-tool.md).
- 2026-10-28: Alexander accepted the layout ADR. Moved from drafts to `9 pkm/adr/ADR-019.md`, status `accepted`, index row added. `authored_by: agent` left in place until he reviews it. The ADR check at close will find this recorded.
- 2026-10-28: wrote [akm-guide](../../spec/akm-guide.md) (what AKM is, why, the four questions, the layout, roles, a working session, adoption with and without the tool, versions, rules of thumb, what is unproven). Example ids checked against the id pattern.
- 2026-10-28: Alexander registered that he may prefer AKM Workbench as the UI's own name with `akmctl` as the CLI; not decided, convention has momentum. Logged in the spec, the Go tool and API spec tasks and the guide. [261027-ensf-akm-an-agent-knowledge-management](../../adr/261027-ensf-akm-an-agent-knowledge-management.md) is accepted and left unchanged (the name there is a working title). Linked tasks checked: none resolved by this one; a log line added to each open one.
- 2026-10-28: ADR check: candidate (questions 1 to 5), already recorded by [261027-ensf-akm-an-agent-knowledge-management](../../adr/261027-ensf-akm-an-agent-knowledge-management.md); `adr: recorded`.
- 2026-10-28: closed. Follow-up filed: [261028-1z36-akm-spec-open-questions](261028-1z36-akm-spec-open-questions.md).
