---
type: design
status: draft
created_at: 2026-10-26
updated_at: 2026-10-26
authored_by: agent
tags:
    - pkm/rails
---
# Agent-ready repository: discovery reports (Option A and Option B)

Output of the discovery pass in [[261026-m74k-agent-ready-repo-layout-spike]]. Two read-only agents worked in isolation from the contract draft in the task note: the A agent read Alexander's repos and the vault; the B agent had no local file access and used the web only. Both reports are condensed here from the agents' final messages. Web sources are as the B agent reported them; none was re-checked by the main thread. Inputs from the earlier pass: [[agent-ready-repo-layout]]. This note decides nothing.

Out of scope: `parametricstudios.com` (excluded) and any work or employer repo (confidential; agents never ask).

## Contract draft (C1 to C8)

See the task note for the wording: C1 Orientation, C2 Rationale, C3 Direction, C4 Current state, C5 Boundaries, C6 Procedures, C7 Trust, C8 Cost.

### Purposes the two agents say the draft misses

Found independently by both agents (strong signal):

- **Verification evidence:** how a reader knows a change is correct; a record says what was and was not verified; an item is never "done" because code exists. A: LOG records what was *not* verified. B: C9; the long-running-agents post names declaring victory early as a core failure.
- **Provenance:** which records an agent wrote and which a human has reviewed (`authored_by: agent`).
- **Portability:** the layout and skills work across tools and machines; a cold agent on another machine may not have `~/.claude`.
- **Sensitive content:** what a record may never contain (credentials, employer material, private writing, copyrighted text). C5 covers actions, not content.

Found by the A agent only:

- **Authority and conflict rule:** a declared source of truth ("spec wins", "repo `AGENTS.md` wins over the skill") so two documents never compete.
- **Decision ownership:** agents propose, Alexander decides; ADOPTED / PROVISIONAL / OPEN lifecycle for agent-proposed rules.
- **Handoff and concurrency:** modes, claim and sign-off, "active with a log line from another thread in the last day: ask".
- **Human review surface:** PRs, review notes, "waiting on Alex" lists.
- **Capture and triage:** raw to ready refinement.
- **Registry of generated files:** what is generated, how, what guards it.

Found by the B agent only:

- **Recoverability:** small commits, one item per branch.

## Option A: what Alexander already does

Method: file reads plus mechanical checks (link and index resolution, status comparison, skill diffs). "Verified" below means checked mechanically.

### Inventory

| Case | AGENTS.md | CLAUDE.md | Roadmap | Projects | ADRs | Repo skill | CI | Other |
|---|---|---|---|---|---|---|---|---|
| alexregan.com | 123 lines | none | yes | 20 | 14 | `.agents/skills/documentation` (115 lines, not symlinked into `.claude/skills`) | none | `pii-inspector` agent, `editorial/` disclosure policy and ledger, `scripts/editorial-check.sh`, `docs/{architecture,deployment,dns}.md` |
| poetic-overtones | 142 | 1 line `@AGENTS.md` | yes (board, stable task IDs) | 3 | 18 | none (relies on user-level skill) | none | `spec/` as authority, `docs/HISTORY.md`, `docs/reviews/<date>.md` decision ledgers; no git remote |
| shavian2tengwar | 125 | 1 line import | yes | 4 | 13 | `documentation` (70 lines; `.claude/skills` symlink) | `.github/workflows/test.yml` | `spec/` with rule IDs, `corpus/` as test suite, guard tests, generated-files table |
| mysqlbak | 14 | none | none | none (`docs/tasks/task-NN-*.md`, no status) | 3, no index | none | none | `docs/spec.md`, `build.sh` |
| pkm-workbench | none | none | `ROADMAP.md` at root | none | none | none | none | README; design and ADR 261025-r0ye live in the vault; no remote |
| Vault | `AGENTS.md` | `@AGENTS.md` | none (task queue and `Tasks.base`) | none (task notes) | 19 plus index | 8 skills in `.agents/skills/`, no `.claude/skills` | none; 4 Claude Code hooks | 6 runbooks, 12 design docs, `skills-map.md`, `Repos.md`, `me.md` |

The user-level `project-docs` skill (`~/.claude/skills/project-docs`, 5 templates) prescribes numbered project folders, a Now/Next/Done roadmap, an append-only LOG with a per-entry session summary, Nygard ADRs accepted by merging the PR, `Review` and `Revert` fields on each ADR, squash merges, and `.wip/` and `.tmp/` scratch. It says the repo's `AGENTS.md` wins.

Notable mechanisms:

- **Vault:** Plan/Query/Build mode gates; `task-lifecycle` skill with an overwritten `## Next` block (State, Next step, Read first, In flight, Last mode, Saved); hooks `SessionStart`, `UserPromptSubmit`, `Stop` (blocks once if an active task note is older than the work), `PreCompact`; an ADR check at task close and a periodic ADR sweep.
- **shavian2tengwar:** the most mechanised. Spec-first rule IDs, corpus written from the spec, tests that enforce coverage and staleness, a generated-files table.
- **poetic-overtones:** ordered "Start Here" reading list; rule IDs with ADOPTED / PROVISIONAL / OPEN that an agent must not promote; review notes as a decision ledger.
- **alexregan.com:** content and deploy boundaries (agents never publish posts, deploy only when asked), a drafts-leak warning, a PII gate.

### Strengths

- A shared spine across repos: `AGENTS.md`, ADR index, roadmap, numbered projects, append-only LOG, ignored `.wip/`.
- "Spec is the authority; when code disagrees the code is wrong", with tests guarding it.
- The vault's `## Next` protocol is the only mechanism built for cold resume, and hooks back it.
- Decision records that name the PR, the revert path and rejected alternatives.
- Explicit permission boundaries and recorded privacy and copyright policies.

### Verified clean

- Every ADR file is in its repo's index and every index link resolves (3 repos).
- Every roadmap project link resolves and every project dir appears in its roadmap (`alexregan.com`, `shavian2tengwar`).
- Project README statuses agree with roadmap sections in those two repos.
- Vault ADR index: 19 files, 19 rows, all `accepted`.

### Inconsistencies and drift found

1. The core workflow is in no repo, and repos disagree about it. Only `shavian2tengwar` names the dependency and gives a fallback; `poetic-overtones` says "load the skill" with none; `alexregan.com` points at its own older skill and never mentions `project-docs`.
2. Three variants of one skill: the `alexregan.com` copy lacks PR-based acceptance and Review/Revert; the `shavian2tengwar` copy is a delta on top of the global one.
3. LOG heading style is mixed (`: ` versus ` — `), frozen in old append-only entries.
4. ADR status vocabulary differs: Accepted / Active / `ACCEPTED`; `alexregan.com` uses "Proposed" and "Proposed (on hold)" that its own legend lacks; the vault allows both `active` and `accepted`.
5. Review/Revert fields: 0 of 14 in `alexregan.com`, all in `poetic-overtones` (18) and `shavian2tengwar` (13). `poetic-overtones` uses them with no remote, so no PR number exists.
6. Numbering styles differ (4-digit, vault 3-digit with `id:`, `mysqlbak` `# 2.`). `alexregan.com` has two independent 4-digit spaces (project and ADR) that look like collisions.
7. Entry points are uneven. `pkm-workbench` has no `AGENTS.md`; `alexregan.com` and `mysqlbak` have no `CLAUDE.md`.
8. Skills are undiscoverable by Claude Code in `alexregan.com` (no symlink) and the vault (no `.claude/skills`).
9. Vault version drift: `skills-map.md` says `task-lifecycle` 1.1.0; the skill says 1.2.0.
10. Vault ADR sweep overdue by its own rule: 12 tasks carry `adr: candidate` against a threshold of 5; last sweep 2026-10-23.
11. State lives in different places per repo; no repo has a single next-step field.
12. Stale-by-design records: `poetic-overtones` READMEs say "not yet committed when this was written"; `alexregan.com` has a dated "Where Things Stand" block that will rot.
13. Working-tree state is not a record: uncommitted files and a leftover branch for a completed project are invisible to a cold reader of the roadmap.
14. CI exists only in `shavian2tengwar`.
15. The vault `Stop` hook is a heuristic with documented holes (`runbooks/Agent Hooks.md`).
16. Three repos have no remote (vault, `poetic-overtones`, `pkm-workbench`); PR-based acceptance assumes one.

### Contract mapping

| Clause | Mechanism | Strength |
|---|---|---|
| C1 Orientation | `AGENTS.md` plus README; ordered reading list; Commands blocks; `Repos.md` | Strong in `shavian2tengwar`, `poetic-overtones`; partial in `alexregan.com`, vault; weak in `mysqlbak`, `pkm-workbench` |
| C2 Rationale | Nygard ADRs with alternatives, Review, Revert; ADR check and sweep; review notes | Strong in two repos; partial elsewhere |
| C3 Direction | Roadmap sections, project statuses, `Tasks.base`, ADR Rejected/Superseded | Strong in `poetic-overtones` and vault; "Dropped" explicit only in vault task statuses |
| C4 Current state | Vault `## Next` plus hooks; repos use LOG session summary, roadmap "Now", README status | Strong in vault; partial in repos (convention, no single next-step field) |
| C5 Boundaries | `AGENTS.md` rules, mode gates | Strong in prose; mechanical only in `alexregan.com` (PII gate) and `shavian2tengwar` (guard tests) |
| C6 Procedures | `project-docs`, `docs/deployment.md`, vault runbooks, Commands blocks | Partial; the skill is external and copies diverge |
| C7 Trust | CI and guard tests in one repo; `pkm-wb check`; append-only LOG; Stop hook | Partial and uneven |
| C8 Cost | Thin `CLAUDE.md`; hub-and-skill progressive disclosure; capped `SessionStart` injection | Partial to strong in vault; `AGENTS.md` files of 123 to 142 lines carry structure trees that can drift; upkeep is high for small fixes |

## Option B: independent design

Method: web research, then one complete design. The B agent read no local file. It reported that it found no empirical study of what history layers (ADRs, roadmaps, logs) do for agents, and no citable critique of ADR staleness; the layout, checks and tiers are its own reasoning.

### Sources and authority (as reported)

- agents.md: vendor-neutral convention; nearest file wins; chat prompt overrides.
- arXiv 2602.11988 (Gloaguen et al., "Evaluating AGENTS.md"): 138 Python tasks, 4 models. Context files gave no gain or a slight loss and raised cost over 20%; LLM-generated files were worst; agents obey files, so needless requirements hurt. Small study. The one hard finding points to restraint: always-loaded text costs, so history must load on demand.
- code.claude.com memory docs: instruction files are context, not enforcement; use hooks for must-happen rules.
- agentskills.io spec: lowercase-hyphen names; description up to 1024 characters; body under 500 lines; `references/`, `scripts/`, `assets/` on demand; `skills-ref validate`.
- Anthropic, effective harnesses for long-running agents: progress file plus git; JSON feature list where agents flip only `passes`; one feature per session; failure modes are one-shotting and declaring done early. One case study.
- Nygard (2011) and MADR: sequential never-reused numbers; supersede, don't delete; MADR adds required considered options and an optional Confirmation section.
- Fowler, scaling architecture conversationally; Spec Kit and Kiro (per-feature spec/plan/tasks, heavy for small repos); beads (dependency tracker needing a binary and database); git-bug (issues in git refs, invisible to file browsing); keepachangelog.
- Blogs on `AGENTS.md` length (about 60 to 200 lines) are anecdotal.

### Layout

```text
AGENTS.md            at most 60 lines, always loaded: orientation, boundaries, pointers
CLAUDE.md            one line: @AGENTS.md
docs/
  decisions/         ADRs, MADR-minimal; README.md index generated by the check script
  ROADMAP.md         Now / Next / Later / Dropped, one line per item linking a work item
  work/              one file per issue or plan; done/ holds closed items, never deleted
  ARCHITECTURE.md    optional one-page code map, owner-edited
.agents/skills/      start-session, work-item, record-decision, release (optional)
scripts/docs-check   one stdlib script, no dependencies
```

No hand-kept state file: current state is derived by `scripts/docs-check --open`, which lists non-closed work items with status, Next and age.

### Memory layers

| Layer | Files | Writer | Rule |
|---|---|---|---|
| L0 Orientation and boundaries | `AGENTS.md`, `ARCHITECTURE.md` | Owner; agent proposes | Mutable; each line must prevent a mistake |
| L1 Rationale | `decisions/*.md` | Agent drafts as proposed; owner accepts | Accepted body immutable except status, superseded-by, typos; numbers never reused |
| L2 Direction | `ROADMAP.md` | Owner orders; agent proposes | Items leave only to Dropped with reason and date |
| L3 In flight | `work/NNNN-*.md` | Agent and owner | Frontmatter and `## Next` overwritten; Log append-only |
| L4 Procedure | `.agents/skills/` | Owner-approved | Agent suggests, does not self-edit mid-task |
| L5 Reality | code, tests, git log | n/a | Wins any conflict with prose |

### Work items

`docs/work/NNNN-slug.md`, frontmatter `id`, `title`, `status` (raw, ready, active, blocked, done, dropped), `kind`, `roadmap` (now, next, later, none; must match `ROADMAP.md`), `decisions`, `blocked_by`, dates. Sections: Goal with a checkable done-when list, Next (required when active), Log (append-only), Outcome (evidence for each done-when line). Roadmap lines are `- [0007] Title`; a roadmap id with no work file fails the check. B rejects beads and git-bug for a solo repo because a cold reader cannot read non-file data; GitHub issue numbers go in frontmatter for outside reports.

### `AGENTS.md` rules (draft, about 30 lines)

- Sections: one-sentence identity; Commands (build, test, lint, run, full check before commit); Where things live (pointers to `decisions/`, `ROADMAP.md`, `docs-check --open`); Rules; Do not, without asking.
- Rules: run `docs-check --open` before starting and resume an active item; work only under a work item (create a `raw` one and ask if none fits); a hard-to-reverse choice or one that rejects a real alternative needs an ADR draft; never edit an accepted ADR's body, supersede it; update Next and Log in the same commit as the code; code and tests win over prose.
- Do not, without asking: change public flags, output formats or schema; add dependencies or leave the item's scope; reorder the roadmap, accept an ADR or edit skills; push, publish, tag or release.
- Must-hold boundaries belong in hooks or permission rules too, since instruction files are advisory.

### Skills (drafted by the B agent)

- **start-session:** run `docs-check --open`; resume the one active item (ask if several); else take the top `ready` item under Now and set it active with a Next; run the full check for a green baseline and stop if red; read only the ADRs the item lists; one item at a time; never mark done because code exists.
- **record-decision:** search first and do not reopen an accepted ADR; next number zero-padded to 4; template with Context, Considered Options (at least two, including "do nothing"), Decision, Consequences, Confirmation (a test or grep that fails if violated); under a page; supersede by changing only the old record's status; stays `proposed` until the owner accepts.
- **work-item:** create (`raw`; `ready` only with a Goal and done-when), update (overwrite Next, append Log), block, close (name the test, command or file proving each done-when line, move to `done/`), drop (reason, `Dropped` on the roadmap); roadmap order is the owner's call.
- Templates and scripts live in each skill's `assets/` and `scripts/`, so skills copy cleanly between repos.

### Mechanical checks (`scripts/docs-check`, local, pre-commit, CI)

ADR numbers unique and contiguous, status in the allowed set; supersession symmetric and links resolve; accepted ADR body unchanged since acceptance (stored hash or CI diff against main); work frontmatter valid, active implies non-empty Next, blocked implies `blocked_by`, closed items in `done/`; every roadmap id resolves and matches `roadmap:`; staleness (active item older than the last commit touching its paths, or over 30 days); `AGENTS.md` length (warn over 80, fail over 150) and every named command exists; relative-link check; `skills-ref validate`; regenerated ADR index, CI fails on diff. No script can judge whether a record is true.

### Lifecycle and tiers

- **First commit:** `AGENTS.md` and the minimal tier; ADR 0001 for the first real choice.
- **Growth:** split the single work file at about 150 lines; add `ARCHITECTURE.md` past a handful of modules; add `docs-check` past about five ADRs or items; add a skill once a procedure has been done correctly twice by hand.
- **Closing:** mark `AGENTS.md` archived and superseded; close or drop all items; keep history.
- **Minimal tier (small CLI or library):** `AGENTS.md`, `docs/decisions/`, one `docs/WORK.md` holding Now/Next/Later/Dropped and active items with a `Next:` line; link check only.
- **Standard:** full layout, `docs-check`, four skills.
- **Large:** nested `AGENTS.md` (nearest wins), per-component ADR folders, optional spec/plan/tasks files for big items.
- **Data or research repos:** add `docs/data.md` for sources, licences and generated versus curated data; mark generated data read-only.

### Contract mapping and failure modes

| Clause | Mechanism | Strength |
|---|---|---|
| C1 | `AGENTS.md`, `ARCHITECTURE.md`, command-exists check | Strong for commands; medium for the map |
| C2 | ADRs with rejected options, immutability check | Medium; depends on writing discipline |
| C3 | Roadmap, work items, Dropped with reasons | Medium |
| C4 | Next/Log, `docs-check --open`, `start-session`, staleness check | Medium to strong; staleness detected, not prevented |
| C5 | `AGENTS.md` Do-not plus hooks or permissions | Weak as text alone; strong with hooks |
| C6 | Validated skills | Medium; skills may not trigger |
| C7 | Schema, link and immutability checks; code wins | Structure strong; truth weak |
| C8 | Tiny always-loaded file; on-demand layers; per-item files | Strong by design |

Failure modes the agent names: silent staleness (well-formed but wrong); the agent forgets the log; the skill does not load (duplicate critical triggers in `AGENTS.md`); over-documentation hurting task success; ADR inflation by agents; a fiddly immutability hash; context loss mid-session (write Next at checkpoints); over-engineering tiny repos.

### Open questions B left for Alexander

1. Must records live under `docs/`, or do some repos (a published static site) need another place?
2. Who accepts ADRs: only Alexander, or the agent for low-stakes ones? A commit, a PR merge, or a frontmatter flip?
3. Issues as files only, or GitHub issues authoritative in some repos? Which wins if both?
4. Work-item numbering: per-repo counters or globally unique?
5. Immutability check: stored hash or CI diff against main?
6. A periodic review that lists stale-looking records?
7. `authored_by: agent` on every agent-written file; who removes it?
8. Skills copied into each repo or kept in one personal repo and linked?
9. Pre-commit strictness: block or warn?
10. Tier rule, such as "move up past five ADRs, or when relearning the repo after a month away".

## First observations (main thread, not a comparison)

- The A report is mostly an audit of one spine that already works in places (`shavian2tengwar`, the vault) and drifts in others. The B report is a single design of similar shape with more mechanical enforcement.
- The two agents did not see each other's work, yet both flagged verification, provenance, portability and sensitive content as missing from the contract.
- B's open question 8 (copy skills into each repo versus one shared copy) is the same problem as A's findings 1 and 2.
- Whether B's mechanisms (single check script, derived state, per-item work files) are better than A's (numbered project folders, LOG, roadmap) is for the comparison step. This note does not rank them.
