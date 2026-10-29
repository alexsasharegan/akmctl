---
type: design
status: draft
created_at: 2026-10-26
updated_at: 2026-10-26
authored_by: agent
---
# Agent-ready repository layout: inputs for the spike

Inputs for [261026-m74k-agent-ready-repo-layout-spike](../work/done/261026-m74k-agent-ready-repo-layout-spike.md). This note records what was said and found on 2026-10-26 so a fresh thread can weigh the options without starting cold. It decides nothing.

## What Alexander wants

A guide to what he calls an **agent-ready repo**: a repository that keeps its own history (decisions over time, roadmap items, issues, work in flight), so an agent with no memory can tell what the project is, how to work in it, why it is the way it is, and what is happening now. The guide covers repo layout, procedures and skills.

Two candidate bases, to be weighed on equal footing:

- **Option A:** his existing practice across his repos.
- **Option B:** the agent's own proposal, built from web sources and original ideas (below).
- **Option C:** a mix of both. The spike decides.

Alexander's request for the spike: treat A comprehensively first, then re-propose B at matching scope (the first B was short and would look weaker than it is), then compare. Do not carry the agent's leanings from this first pass into the new thread's choice.

## Sources (web, 2026-10-26)

The searches returned plenty on `AGENTS.md` and the long-running-agent progress file pattern, and almost nothing on ADRs, roadmaps, plans directories or skills. Parts of Option B are therefore the agent's own design. Most sources are blogs and gists of uneven authority.

- [How to write an effective AGENTS.md](https://agent-ready.dev/how-to-write-an-effective-agents-md)
- [Writing an Effective AGENTS.md for Coding Agents (Agentic AI Foundation)](https://aaif.io/blog/writing-an-effective-agents-md)
- [Agents.md best practices (gist)](https://gist.github.com/0xfauzi/7c8f65572930a21efa62623557d83f6e)
- [AGENTS.md Guidelines (gist)](https://gist.github.com/jerdaw/3917eab775d3e4bbcf37928101fbc3db)
- [Improve your AI code output with AGENTS.md (Builder.io)](https://www.builder.io/blog/agents-md)
- [Effective harnesses for long-running agents (Anthropic)](https://anthropic.com/engineering/effective-harnesses-for-long-running-agents)
- [Long-running Agents (Addy Osmani)](https://addyosmani.com/blog/long-running-agents/)
- [Code as Agent Harness (arXiv)](https://arxiv.org/html/2605.18747v1)

What the sources support:

- Root `AGENTS.md`, nearest-file-wins in a monorepo, explicit user prompts override repo files.
- `AGENTS.md` as an index that links out to deeper docs (architecture, ADRs), kept short (one gist says about 150 lines at most).
- A layout section so agents find files and reuse patterns; build, test and lint commands with flags; boundaries (generated files, lockfiles); how to verify work is done; an escape hatch to ask when unsure.
- Review `AGENTS.md` when the runtime, CI, layout or release process changes, and after an agent mistake.
- The long-running-agent pattern: an initializer writes a progress log, a feature list and an init script; later sessions read the progress file and the git log, do one failing feature, test, and commit. Repo as shared memory; survey language: working context, session state, long-term memory.
- Critiques: a JSON feature list reinvents an issue tracker; the last stretch of a project stays expensive; memory still breaks across long sessions.

## Option B: the agent's first proposal (2026-10-26, as given)

An agent-ready repo is one where a fresh agent with no memory can answer: what is this, how do I work in it, why is it like this, what is happening now.

### Layout

```text
AGENTS.md          index, not an encyclopedia: commands, boundaries, links (~100 lines)
CLAUDE.md          one line importing AGENTS.md
README.md          for humans: what, install, run
ROADMAP.md         now / next / later, and what was dropped
CHANGELOG.md       user-visible changes (optional for private tools)
docs/
  architecture.md  the map: packages, data flow, invariants
  adr/             one decision per file, numbered, immutable once accepted
  plans/           one file per piece of work in flight, deleted or archived when done
  runbooks/        how to do recurring things (release, add a theme, debug the stream)
  issues/          local issues as files, only if there is no tracker
.agents/skills/    procedures an agent loads on demand
```

### Memory layers

- **Long-term (why):** ADRs hold decision, options, consequences. Accepted ADRs are never edited; a new one supersedes and links back. `architecture.md` states invariants.
- **Medium-term (what next):** `ROADMAP.md`, human-edited and short; each item links to a plan or issue; dropped items stay in a "Not doing" list so they are not re-proposed.
- **Session state (where we are):** a plan file with Goal, Done when, Steps, Next, Log, Open questions. A fresh agent starts from `## Next`.
- **History:** commit messages say why, and cite an ADR or plan number when one applies.

### Issues

Use the tracker if there is one. Files in `docs/issues/` pay off only for a private tool with no tracker. Start with roadmap and plans; add issue files when the lack is felt.

### Procedures as skills

Short `SKILL.md` files: **start-session** (read `AGENTS.md`, the active plan's Next, `git log -10`, run the tests first), **end-session** (refresh Next, log what was verified and not, commit), **decide** (when something needs an ADR; template; numbering), **new-plan** and **close-plan**, **verify** (exact commands and what "done" means).

### Rules in `AGENTS.md`

Short and linking out; commands first; what not to touch (generated files, accepted ADRs); ask when unsure; update after an agent mistake; review on a schedule.

### Mechanical checks

ADR numbers unique and supersession links forward; every roadmap item resolves to a plan or issue; no plan `active` for N days without a log line. Run from CI or a hook. `pkm-wb check` is the same idea.

### Open questions the first proposal left

Where a closed plan goes (archive, not delete); one roadmap or one per project (one, for a single tool); whether a code repo keeps its own plan or links the vault task.

## First-pass analysis of Option A (Alexander's repos)

Read on 2026-10-26: `alexregan.com`, `poetic-overtones`, `shavian2tengwar`, `mysqlbak`, and the user-level skill `~/.claude/skills/project-docs`. Not read: `parametricstudios.com` (Alexander excluded it), the vault as a repo case, and the work repos he says are more fleshed out and will show.

What Option A has:

- **Project folders:** `docs/projects/NNNN-title/` with `README.md` (status, objective, scope, definition of done), `TASKS.md` (checklist in phases), `LOG.md` (append-only dated journal) and an untracked `.wip/` scratchpad. Numbers are never reused.
- **ADRs:** Nygard style in `docs/adr/NNNN-title.md` with an index and status legend in `docs/adr/README.md`. **Accept on merge** (`shavian2tengwar` ADR-0013): an ADR is written as Active in the PR that implements it; merging accepts it; Proposed is for spikes and roadmap plans; squash merges preferred; the PR closes its project.
- **Roadmap with state:** `docs/ROADMAP.md` with Now, Next, Backlog, On Hold, Done, plus "Where Things Stand" and "Waiting on Alex" blocks (`alexregan.com`).
- **Reviews:** `poetic-overtones` `docs/reviews/YYYY-MM-DD.md` for overnight runs: what was built, "New decisions of mine (PROVISIONAL, for you to adopt or change)", "Needs you", risks. A decision ledger.
- **History:** `docs/HISTORY.md` narrates the project in phases.
- **Spec first:** `spec/` with rule IDs and the corpus as the test suite (`shavian2tengwar`).
- **`AGENTS.md` content beyond commands:** mandatory-skill trigger, project identity, repo structure, execution protocol, deployment rules ("only deploy when Alex asks"), Related Repositories, Content Rules (agents never publish posts), a sensitive-information policy.
- **Skills:** a user-level `project-docs` skill (152 lines plus templates) with a repo-specific `documentation` skill in `.agents/skills/` (`alexregan.com`, `shavian2tengwar`; `shavian2tengwar` also has `.claude/skills/`).
- **A small tier:** `mysqlbak` has a 14-line `AGENTS.md` pointing at `docs/spec.md`, `docs/adr/` and `docs/tasks/`.
- **CI:** `shavian2tengwar` has `.github/workflows/test.yml`.

What Option B has that A seemed to lack (to verify):

- A resume pointer (`## Next`), since `LOG.md` is append-only.
- `runbooks/` for recurring procedures.
- A "Not doing" list.
- Mechanical checks on the docs themselves (ADR numbering, roadmap-to-project links, staleness).
- A `CHANGELOG`.

Problems noticed in A:

- The user-level skill is outside version control (ADR-0013 says so); copies already drift between repos, and `poetic-overtones` has no skills folder.
- ADR status words differ (Active in two repos, Accepted in the vault).
- Roadmap prose can drift from the project folders; nothing checks it.

The vault as a case: it uses task notes (`type: task`, `## Next`, `## Log`) instead of project folders; see Tasks and the `task-lifecycle` skill.

First-pass leanings of the agent (inputs, not decisions): use A's folder convention for `pkm-workbench` (`docs/ROADMAP.md`, `docs/adr/`, `docs/projects/`) rather than B's root-level roadmap and `plans/`; add `## Next` to the project README; move the canonical skill into version control (the `dotfiles` repo is a candidate); the repo owns the project folder and the vault task links to it. Treat these as one side of the comparison.

## What the spike still has to do

See [261026-m74k-agent-ready-repo-layout-spike](../work/done/261026-m74k-agent-ready-repo-layout-spike.md).
