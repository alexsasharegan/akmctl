---
type: design
status: draft
created_at: 2026-10-26
updated_at: 2026-10-26
authored_by: agent
tags:
    - pkm/rails
---
# Agent-ready repository: contract and A/B/C comparison

> **Historical (2026-10-28).** This note records the comparison and the decisions as they stood when they were made. [[agent-ready-repo-spec]] supersedes it where they differ: names such as `docs-check`, `NNNN` ids and the Python script were replaced by `akmctl`, `YYMMDD-xxxx` ids and a Go tool, and later decisions are only in the spec and [[261027-ensf-akm-an-agent-knowledge-management]].

Step 3 of [[261026-m74k-agent-ready-repo-layout-spike]]. Evidence: [[agent-ready-repo-discovery]]. This note compares; it decides nothing. Decision points are marked **D1** to **D8**; Alexander's answers (2026-10-27) are in the next section.

## Decisions (2026-10-27)

Made by interview, one point at a time, by the person who will maintain the template. The spec is for any maintainer or team, so decisions are worded generally ("maintainer", "team"). These make Option C concrete: it is B's work-item model and short `AGENTS.md` on A's names, ownership rules and boundaries. The sections below keep the comparison as it stood before the interview.

| # | Decision | Chosen | Notes |
|---|---|---|---|
| D1 | Unit of in-flight work | One file per item, `docs/work/NNNN-slug.md` with Goal, Next, Log, Outcome (B) | Add a dimension linking items into a feature or project: frontmatter fields, and the feature id is itself a work item. So a `parent` field; no separate project-folder type. |
| D2 | Resume pointer | `## Next`, overwritten; Log append-only | A parent item keeps its own Next. |
| D3 | Source of "now" | Derived listing from item frontmatter, plus a hand-ordered `ROADMAP.md` | Roadmap is only the ordered lists. Status comes from the derived listing. The listing may use a `priority` frontmatter field to sort; the roadmap supersedes it. Both "Where Things Stand" and "Waiting on Alex" prose blocks are dropped. |
| D4 | Checks | One stdlib script, shipped in the shared skill's `scripts/` | Run by the agent (AGENTS.md and the skills say so). No hook or CI requirement. |
| D5 | Skill home | Canonical in a versioned repo; each project repo carries a copy under `.agents/skills/` with a version | Repo `AGENTS.md` wins over the skill on conflict. Distribution settled in O4. |
| D6 | ADR acceptance | A maintainer's commit sets `accepted` (the merge of the PR that carries it, where a remote exists) | Statuses: proposed, accepted, rejected, deprecated, superseded. Immutability of accepted bodies is convention only, no check. |
| D7 | `AGENTS.md` | Short, pointers out, about 60 to 80 lines | Structure tree to `docs/architecture.md`; long policies to named docs linked from the Do-not list; script warns over 80, fails over 150. |
| D8 | Folder names | `docs/adr/`, `docs/work/`, `docs/ROADMAP.md` | `docs/projects/` is replaced by `docs/work/`; existing project folders need a migration note. |

### Open consequences, settled 2026-10-27

| # | Question | Chosen | Notes |
|---|---|---|---|
| O1 | How an item says it waits on someone | `status: blocked` with `blocked_by` | `blocked_by` is a list. An entry that looks like an item id must resolve to an item; any other entry is free text naming a person, role or outside dependency. The generated list can group blocked items by `blocked_by`. Dependencies and people share one field. |
| O2 | Closing a parent item | All children must be `done` or `dropped` | The check fails otherwise. Dropping a parent requires dropping or re-parenting its open children. |
| O3 | `priority` | Optional integer rank; higher is sooner; an item with none counts as 0 | Sorts the generated list only; the roadmap supersedes it. Ranks can collide or go stale; the roadmap resolves that. |
| O4 | How the boilerplate reaches a repo | One template repo holds the skill, check script and templates with a version; a project copies them in | Each copy carries a `VERSION`; the check script reports when the copy is behind the template's version (it cannot fetch, so a human gives it the template's version or path). Updating is a deliberate copy. |

### Still to write

- The work-item frontmatter schema (fields, allowed values, defaults) from D1, D2, D3, O1 to O3.
- The check script's rules, one line each, from D3, D4, D6, D7, O1, O2, O4.
- The boilerplate: `AGENTS.md` template, ADR and work-item templates, and the skill (`start-session`, `work-item`, `record-decision`) with the check script.
- The guide, and an ADR for the layout (D6 says the maintainer's commit accepts it).
- Pilot on `pkm-workbench`; a migration note for repos that use `docs/projects/`.

## The contract (C1 to C12)

C1 to C8 are the original draft. C9 to C12 are the grouping of the clauses the two discovery agents added; Alexander accepted the grouping on 2026-10-26.

A repo's memory exists so that a cold reader (an agent with no memory, or a maintainer back after months away) can work well and not repeat past mistakes. It lets that reader:

| # | Clause | Meaning |
|---|---|---|
| C1 | Orientation | Learn what the project is, how to build, test and run it, where things live. |
| C2 | Rationale | Learn why it is the way it is, including rejected alternatives, so settled decisions are not re-litigated. |
| C3 | Direction | Learn what is planned, in what order, and what was dropped. |
| C4 | Current state | Learn what is in flight, where it stands, the single next step; work is resumable. |
| C5 | Boundaries | Learn what must not be touched or done without asking. |
| C6 | Procedures | Repeat recurring work the same way each time. |
| C7 | Trust | Records stay accurate; drift is prevented or detected; it is clear what may be edited and by whom. |
| C8 | Cost | Cheap to read cold, cheap to keep up. |
| C9 | Authority and ownership | A declared source of truth when documents conflict (code and spec win over prose; the repo's `AGENTS.md` wins over a skill); who decides (agents propose, a maintainer accepts); a human review surface (what waits on a maintainer). |
| C10 | Verification | Records say what was and was not verified; nothing is "done" because code exists; done has evidence. |
| C11 | Content safety and provenance | What a record may never contain (credentials, employer material, private writing, copyrighted text); which records an agent wrote and which a maintainer has reviewed. |
| C12 | Portability and handoff | The setup works across tools and machines without relying on `~/.claude`; modes, claiming and concurrent threads are handled. |

Folded in, not separate: capture and triage (into C3 and C4), the generated-files registry (into C5 and C7), recoverability (into C8; small commits).

## Two caveats before reading the tables

- **A is practised; B is a paper design.** B's strength ratings are untested. A's drift (16 items) is real, but Alexander points out it is **not a verdict on A's design**: A grew from his own effort across several agentic tools, without a spec or target architecture to hand them. Spec-less drift is the problem this spike exists to fix, so weigh the drift list as evidence of what a spec must prevent (C7), and weigh A's *design choices* (names, ownership rules, boundaries) on their merits.
- **A is not one thing.** The vault and `shavian2tengwar` are the best cases; `mysqlbak` and `pkm-workbench` are the weakest. The tables rate A's *best case* and its *typical case*.

## Clause by clause

Strength: strong, partial, weak, none. A is rated best case / typical.

| Clause | A: Alexander's practice | B: independent design | C: mix (proposed below) |
|---|---|---|---|
| C1 Orientation | strong / partial. Long `AGENTS.md` (123 to 142 lines) with structure trees that can drift; `pkm-workbench` has none. | strong for commands (a check verifies they exist), medium for the map. `AGENTS.md` at most 60 lines, pointers out. | B's short file and command check; A's ordered "Start Here" list. |
| C2 Rationale | strong / partial. ADRs with alternatives, Review, Revert, in two repos; none of Review/Revert in `alexregan.com`. | medium. MADR-style ADRs with a required Confirmation (a test or grep). Immutability enforced by script. | A's ADRs and Review/Revert; B's Confirmation field and immutability check. |
| C3 Direction | strong / partial. Roadmap sections; "Dropped" explicit only in vault task statuses. | medium. Roadmap with a Dropped list and reasons; each line must resolve to a work file. | A's roadmap; B's Dropped list and the link check. |
| C4 Current state | strong (vault: `## Next`, claim, hooks) / partial (repos: LOG summary, no single next-step field). | medium to strong. `## Next` in each work item; state derived by `docs-check --open`; staleness check. | Vault's `## Next` pattern in the repo's work unit; derived listing (B). |
| C5 Boundaries | strong in prose; mechanical in two places (PII gate, guard tests). | weak as text, strong with hooks. Says to back must-hold rules with hooks. | A's prose boundaries plus B's rule to move must-hold ones into hooks or permissions. |
| C6 Procedures | partial. The central skill is outside every repo; three variants. | medium. Skills carry their own templates and scripts, so they copy between repos; may not trigger. | Skill in version control (see **D5**). |
| C7 Trust | partial, uneven. CI and guard tests in one repo; convention elsewhere; 16 drift items found. | structure strong, truth weak. One `docs-check` for numbering, links, statuses, staleness, index. | B's one portable check script, run per repo; A's code-wins rule. |
| C8 Cost | partial to strong (vault); high upkeep for small fixes (3 files per project plus roadmap plus ADR index). | strong by design for reading (small always-loaded file; history on demand); upkeep discipline needed; tiers for small repos. | B's tiers; A's "one-line LOG entry for a small fix". |
| C9 Authority and ownership | strong. Spec wins, "repo `AGENTS.md` wins", agents propose and Alexander decides, ADOPTED/PROVISIONAL/OPEN, review notes, "waiting on Alex". | partial. Ownership table and "code wins"; ADR acceptance is an open question. | A's, nearly whole. |
| C10 Verification | partial. LOG records what was not verified; vault close protocol. Not uniform. | medium. Outcome must name evidence for each done-when line; "never mark done because code exists". | B's evidence rule in A's close step. |
| C11 Content safety and provenance | strong in prose (privacy policies, PII gate in one repo, `authored_by: agent` in the vault). | weak. Names `authored_by` and sensitive material as missing purposes; gives no mechanism. | A's, extended to repos (see **D7**). |
| C12 Portability and handoff | partial. `AGENTS.md` as shared file, `CLAUDE.md` as import; vault has modes, claiming, hooks; the skill depends on `~/.claude`; 3 repos lack a remote. | partial. Skills self-contained; no concurrency model. | A's modes and claiming; B's self-contained skills. |

## Other criteria

| Criterion | A | B | C |
|---|---|---|---|
| Cold-start cost | High in repos (long `AGENTS.md`, project folders); low in the vault (injected `me.md`, task list). | Low by design: one small file, then on-demand. Untested. | Aim at B's. |
| What it enforces | One repo has CI and tests; vault has hooks; the rest is convention. | One portable check script plus optional hooks. | B's script everywhere, A's hooks where they exist. |
| Upkeep per change | Heavy: project README, TASKS, LOG, roadmap, ADR index. | Lighter: one work file, roadmap line; index generated. | See **D2**. |
| Drift risk | Measured: 16 items. | Predicted: silent staleness; checks catch structure, not truth. | Fewer hand-kept lists; same truth gap. |
| Small repos | Poor: `mysqlbak` has a 14-line setup, nothing consistent; a full project folder for a one-day fix is heavy. | Minimal tier (three files). | Minimal tier from B. |
| Large repos | Good: numbered projects, spec-first, guard tests. | Nested `AGENTS.md`, per-component ADRs. | A's projects for big work. |
| Fit with vault task notes | The vault is its own case; repos and vault use different units (project folder versus task note). | Work item with `## Next`, `## Log`, `## Outcome` is nearly the vault's task note. | Same unit in both, so a repo item and a vault task read alike. |

## Where A and B differ: decision points

| # | Question | A | B | Consideration |
|---|---|---|---|---|
| **D1** | Unit of in-flight work | Numbered project folder: README, TASKS, LOG | One file per item: Goal, Next, Log, Outcome | B's unit matches the vault task note; A's suits multi-phase work. |
| **D2** | Resume pointer | LOG "context state for next thread", append-only | `## Next`, overwritten | The vault already proved `## Next` works with hooks. |
| **D3** | Where current state comes from | Roadmap "Now" plus README status, kept by hand | Derived by script | Hand-kept lists drift (A, items 11 to 13); a script cannot judge truth. |
| **D4** | Checks | CI in one repo | One portable script in every repo | Cost of writing the script once, in the skill. |
| **D5** | Where the canonical skill lives | `~/.claude`, outside version control | Inside each repo | Copies drift (A, items 1 and 2); a shared copy couples repos. A middle path: canonical in `dotfiles`, repos carry a pinned copy or a pointer with fallback text. |
| **D6** | ADR acceptance | Merge of the PR | Owner flips status | Three repos have no remote; "Proposed = open PR" cannot occur there. |
| **D7** | `AGENTS.md` | Long, with structure trees and policies | At most 60 lines, pointers | The one hard study finding favours short; policies and structure trees can live in linked docs. |
| **D8** | Folder names | `docs/adr/`, `docs/projects/`, `docs/ROADMAP.md` | `docs/decisions/`, `docs/work/`, `docs/ROADMAP.md` | A's names exist in four repos; renaming costs more than it gains. |

## Option C: a mix, proposed

Option C keeps A's names, ownership rules and boundaries, and takes B's discipline on size, derived state and mechanical checks.

**Layout (standard tier)**

```text
AGENTS.md            at most about 80 lines; commands, boundaries, pointers; no structure tree
CLAUDE.md            one line: @AGENTS.md
README.md            for humans
docs/
  ROADMAP.md         Now / Next / Backlog / On Hold / Done / Dropped (with reasons)
  adr/               A's ADRs plus B's Confirmation field; one status vocabulary
  projects/NNNN-*/   A's folders, for multi-phase work; README carries a ## Next block
  architecture.md    the map (moved out of AGENTS.md)
.agents/skills/      pinned copy of the shared skill, or a pointer with fallback text
scripts/docs-check   one stdlib script
```

**What C takes from each**

- **From A:** `docs/adr/`, `docs/projects/`, roadmap sections; Review and Revert on ADRs; "agents propose, Alexander decides" and ADOPTED / PROVISIONAL / OPEN; privacy and content rules; modes and claiming (where hooks exist); the one-line LOG entry for small fixes.
- **From B:** short `AGENTS.md`; overwritten `## Next`; a Dropped list with reasons; Confirmation field on ADRs; evidence-for-done in Outcome; a portable check script; derived `--open` listing; minimal tier; hooks for must-hold boundaries.
- **New in C, to settle:** one ADR status vocabulary across repos (D6); a rule for repos with no remote; where the canonical skill lives (D5).

**Minimal tier (small repo)**

`AGENTS.md`, `docs/adr/` with an index, and one `docs/ROADMAP.md` with a `Next:` line on the active item. No project folders until a piece of work needs phases.

**What the script checks (`docs-check`)**

ADR numbering, index and status vocabulary; supersession links; roadmap entries resolve to a project or ADR; project README `status` agrees with the roadmap section; an active project has a non-empty `## Next`; `AGENTS.md` length and that its commands exist; relative links; skill validity. It reports staleness. It cannot judge whether a record is true.

**Weak spots C does not fix**

- Silent staleness: well-formed but wrong records.
- Skills may not trigger; critical triggers are duplicated in `AGENTS.md`.
- Moving existing repos is migration work with its own cost; the task asks for a note, not a migration.
- C is a proposal on paper, like B.

## Cost of adopting each

| | A (keep and tidy) | B (adopt as designed) | C (mix) |
|---|---|---|---|
| New to write | A drift cleanup in 4 repos; a shared skill | The script, 4 skills, templates | The script, a merged skill, a tier guide |
| Change to existing repos | Edit in place | Rename folders, new work unit | Small: add the script, a `## Next`, a vocabulary pass |
| Risk | Drift continues | Untested; abandons working conventions | Both, but smaller |
| Vault | Unchanged | Work item overlaps the task note | Vault stays as is; repos adopt its pattern |

## What the maintainer decides next

1. Whether the contract (C1 to C12) is right as written. This is the base for everything else.
2. D1 to D8, or the ones that matter to him; the rest can take C's defaults.
3. Whether C is the base for the pilot on `pkm-workbench`, or whether A or B alone should be.
4. Whether any of this changes how repos are set up enough to need an ADR.
