---
type: design
status: active
created_at: 2026-10-28
updated_at: 2026-10-28
---

# AKM guide: a repository that remembers and runs on rails

Version 0.2.0 (beta).

AKM (agent knowledge management) is an opinionated standard that gives a repository a memory and rails, so any agent can pick up where work stopped. This guide is for a maintainer, solo or on a team, who works with AI coding agents and wants to adopt it. It explains what AKM is, why it exists, how it works day to day, how to adopt it, and what is not proven. The rules themselves are in the specification ([agent-ready-repo-spec](agent-ready-repo-spec.md)); the tool's operations are in the API spec ([akm-api-spec](akm-api-spec.md)).

Everything here works by hand, using plain files, the id recipe and the fallbacks in section 7. Regenerating the ADR index and running the full checks need the tool, `akmctl`.

## 1. What AKM is

AKM is the same idea as personal knowledge management (PKM), moved inside a repository and addressed to agents. PKM externalises memory into notes so the mind need not hold it. AKM externalises a project's operational memory into the repo, so an agent that starts with no memory can read what it needs and carry on.

It has four parts:

| Part         | What it is                                                                                            |
| ------------ | ----------------------------------------------------------------------------------------------------- |
| The standard | A written spec: the layout, record formats, rules, ids and versions                                   |
| The starter  | `akmctl init`, which writes the layout, `AGENTS.md`, record templates and skills into a new directory |
| The tool     | `akmctl`, a Go CLI that queries and checks the records, plus a workbench UI started with `akmctl ui`  |
| The skills   | Instructions an agent loads to follow the procedures                                                  |

The tool is an accelerator, not a dependency. Everything it does can be done by reading plain files.

## 2. Why

**The problem.** Every agent session starts with no memory. A repo's history sits in chat logs, in people's heads and in docs that drift. The result is settled decisions argued again, mistakes repeated, work lost mid-task, parallel agents colliding, and tokens spent reading everything to find the one thing needed.

**The promise.** A cold agent, or you after months away, can answer four questions in a few reads: what is this, why is it like this, what is planned, and where did we stop.

**Two headline differentiators**

1. **Memory, not only orientation.** A good `AGENTS.md` tells an agent how to work in a repo. AKM does that and also keeps the repo's history: decisions and their reasons, direction, and work in flight.
2. **Rails for the work.** Work follows a simple Kanban-style flow with one active item and a stated next step, and no ceremonies. The agent keeps the records current; the maintainer steers.

**Supporting principles:** conducive to parallel work; token-efficient; less drift; works with any agent, with or without tooling; upgradable; open for extension. The full list, with the reasoning behind the design, is in section 0 of the specification.

## 3. The four questions and where they are answered

| Question                              | Answered by                                            | Cost to read                  |
| ------------------------------------- | ------------------------------------------------------ | ----------------------------- |
| What is this and how do I work in it? | `AGENTS.md` (60 to 80 lines), `docs/architecture.md`   | A few dozen lines             |
| Why is it like this?                  | `docs/adr/` (decision records), with a generated index | The index, then one record    |
| What is planned?                      | `docs/ROADMAP.md` (ordered lists)                      | One short file                |
| Where did we stop?                    | `docs/work/` (one file per item), each with `## Next`  | The listing, then one section |

A reader routes by headers. Every record has frontmatter in YAML-lite (id, title, summary, status, dates) that is cheap to read; the body is read only when needed. The ADR index is generated from those headers, and the open-work listing is derived on demand.

## 4. The layout

```text
AGENTS.md                  60 to 80 lines: identity, commands, rules, boundaries, pointers
CLAUDE.md                  one line importing AGENTS.md (optional, tool-specific)
README.md                  for humans
docs/
  ROADMAP.md               hand-ordered lists: Now, Next, Later, Dropped
  architecture.md          the map (optional)
  adr/
    README.md              generated routing index; never edit by hand
    261015-k7f2-use-sqlite.md
  work/
    261027-m3qz-add-json-output.md
    done/                  closed items
  .layout                  manifest: layout and extension versions
.agents/skills/            the shared skills
```

A repository that keeps its records inside another structure, such as a notes vault, can move them: the manifest keys `work_dir`, `adr_dir` and `roadmap` name the work folder, the ADR folder and the roadmap file, and everything above then reads from there (`done/` stays inside the work folder; the manifest stays at `docs/.layout`).

### Ids

Every record has an id of the form `YYMMDD-xxxx`, such as `261027-m3qz`: the UTC creation date plus four random characters from a 32-character alphabet (no `i`, `l`, `o`, `u`). ADRs and work items share one id space. The id alone is the reference; the slug in the filename is only a label, so renaming it breaks nothing. The tool mints ids, because agents are poor random-number generators. There is no shared counter, because counters make parallel branches claim the same number and conflict on merge.

### Work items

One Markdown file per item in `docs/work/`. A feature or project is itself a work item that other items name as `parent`.

```markdown
---
id: 261027-m3qz
title: Add JSON output
summary: Let `list` print machine-readable JSON.
status: active
kind: feature
priority: 2
created: 2026-10-27
updated: 2026-10-28
decisions: [261015-k7f2]
---
# Add JSON output

## Goal
Print the same records as JSON so other tools can read them.
- [ ] `list --json` prints a stable schema
- [ ] the schema is documented and covered by a test

## Next
Write the encoder behind the existing `list` command, then add the schema test.

## Log
- 2026-10-27: filed. Chose a flat schema (ADR 261015-k7f2).
- 2026-10-28: encoder drafted.

## Outcome
```

| Status    | Meaning                                                                                            |
| --------- | -------------------------------------------------------------------------------------------------- |
| `raw`     | Captured, not specified                                                                            |
| `ready`   | Specified: a goal and a "done when" list                                                           |
| `active`  | Being worked; has a `## Next`                                                                      |
| `blocked` | Cannot proceed; `blocked_by` names item ids or free text (a person, a role, an outside dependency) |
| `done`    | Finished; `## Outcome` gives evidence for each "done when" line                                    |
| `dropped` | Decided against; `## Outcome` gives the reason                                                     |

`priority` is an optional integer (higher is sooner, default 0) that only sorts the listing; the roadmap wins. A parent can close only when all its children are done or dropped.

### Decision records (ADRs)

One file per decision in `docs/adr/`, with frontmatter (`id`, `title`, `summary`, `status`, `date`, optional `scope`, `supersedes`, `superseded_by`) and the sections Context, Considered Options, Decision, Consequences, and optionally Confirmation, Review and Revert. Statuses: `proposed`, `accepted`, `rejected`, `deprecated`, `superseded`. A maintainer accepts. An agent drafts as `proposed`. An accepted record is not edited; a change of mind adds a new record that supersedes it.

### The roadmap

`docs/ROADMAP.md` holds only ordered lists, in the sections Now, Next, Later and Dropped, one line per item: `- [261027-m3qz] Add JSON output`. Order is priority. A Dropped line gives a reason.

### `AGENTS.md`

Short, in five parts: identity, commands, where things live, rules, and a Do-not list. Long policies live in their own docs and are linked. It always includes the rule to use the plain-file fallbacks if the tool is absent.

## 5. Who does what

| The agent                                                                      | The maintainer                                                   |
| ------------------------------------------------------------------------------ | ---------------------------------------------------------------- |
| Keeps `status`, `## Next` and `## Log` current, in the same commit as the code | Orders the roadmap                                               |
| Files new items, as `raw` when unsure                                          | Accepts or rejects ADRs                                          |
| Drafts ADRs as `proposed`                                                      | Unblocks what is `blocked` on a person                           |
| Closes items with evidence                                                     | Reviews records marked `authored_by: agent` and removes the mark |
| Runs the checks                                                                | Decides what changes the layout, and runs migrations             |

The paperwork is the agent's. The decisions are the maintainer's.

## 6. A working session

### Start

1. Read `AGENTS.md`.
2. List the open work. With the tool: `akmctl ls --status ready,active,blocked`. Without it: see section 7.
3. If an item is `active`, open it and read `## Goal`, `## Next` and the last few `## Log` lines.
4. Run the checks for a clean baseline. If they fail, record that in the log and stop before building on a broken base.
5. Work one item at a time.

### While working

- Work only under an item. If none fits, file a `raw` one and ask which to work.
- Update `## Next` at checkpoints: before any pause or question, after a commit, and at the end.
- Append a dated line to `## Log` for each decision, finding or change of plan.

### Record a decision

A choice that is hard to reverse, or that rejects a real alternative, gets an ADR. Search `docs/adr/` first. If an accepted record already covers it, do not reopen it; tell the maintainer what it says. Otherwise draft one as `proposed` with at least two considered options (including doing nothing) and the reasons. Link it from the item through `decisions`.

### Close an item

For each "done when" line, name the test, command or file that proves it, in `## Outcome`. Set `done`, move the file to `docs/work/done/`, and update the roadmap. A parent closes after its children.

### Stop or hand off

Rewrite `## Next` so a cold agent can resume: where the work stands, the single next step, what to read first. Commit it with the work.

## 7. Adopting AKM

### A new repo

1. Run `akmctl init --profile code` (a software repository) or `--profile notes` (a notes repository or an Obsidian vault) in an empty directory. It writes the layout, `AGENTS.md`, the skills and the `.layout` manifest, and never overwrites a file.
2. Fill in `AGENTS.md` (identity, commands, boundaries).
3. File the first work items. Write the first ADR for the first real decision.
4. A small repo can stay minimal: `AGENTS.md`, `docs/adr/`, a few work items, `docs/ROADMAP.md`.

### An existing repo

The 0.2.0 beta supports new repositories only. This section describes the planned path and is not yet supported.

A repo with no `docs/.layout` manifest is read as layout `0.0.0`, and no migration can take it to 0.2.0, because 0.2.0 is the first layout: adopting AKM is a one-time conversion that an agent does with you. Give the agent the specification ([agent-ready-repo-spec](agent-ready-repo-spec.md)) and the repo, and ask it for a plan. Before anything is changed:

- Make sure the working tree is clean and tag a backup.
- Work on a branch and review every diff.
- Never delete a file; renaming and moving keep history.

The tool can do the mechanical part of that conversion from a migration file you keep in the working tree: `akmctl migrate --migration <file>` runs the declarative steps (rename or copy fields, set a title from the first heading, scaffold a `## Goal`, mint ids with or without renaming the files) under the same safeguards. The steps are listed in the API spec, section 4.15; the text a step cannot judge, such as the content of a new `## Goal`, is left for a person and listed in the step's notes.

Once a repo has a manifest, `akmctl migrate` handles later layout versions: it reports what differs (`--mode report`), lists the files each safe step would change (`--mode plan`), applies the safe steps after making a backup tag and a branch (`--mode apply-safe --acknowledge-backup`), and writes a prompt for an agent to do the remainder (`--mode handoff`). It never commits and never deletes a file. The design is in the specification, section 7.

### Without the tool

The files are the source of truth, so a bare agent can work:

| Need | Without the tool |
| ------------------------------------- | ----------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| Which decisions exist | Read `docs/adr/README.md` |
| What work is active | `grep -l '^status: active' docs/work/*.md` |
| What is blocked, and on whom | `grep -l '^status: blocked' docs/work/*.md`, then read `blocked_by` |
| Where to resume an item | Read its `## Next` |
| Mint an id | `id="$(date +%y%m%d)-$(LC_ALL=C tr -dc '0-9a-hjkmnp-tv-z' </dev/urandom                                  | head -c 4)"`, then check `ls docs/work docs/work/done docs/adr \| grep "^$id-"` prints nothing |
| Create a record | Copy the template, name the file `<id>-<slug>.md`, set `id` to the same id |
| Close an item | Set `status`, write `## Outcome`, set `updated` to today's UTC date, `git mv` the file to `docs/work/done/` |
| Find a record by id | `ls docs/work docs/work/done docs/adr \| grep '^<id>-'` |
| The ADR index is stale or in conflict | Do not hand-edit it; ask the maintainer to regenerate it with the tool |

This is slower and uses more tokens, but it must stay correct.

### With the tool

`akmctl` is one Go binary with subcommands: `check`, `ls`, `search`, `show`, `board`, `index`, `new`, `set`, `close`, `find`, `when`, `migrate` and `ui`, plus `claim` and `release` for the optional claims extension. (There is no `renumber` command; the specification documents the manual repair.) `akmctl ui` starts a workbench page, a Kanban over the records, on the loopback interface. The tool has a stable JSON output and an API spec with a CLI binding and an HTTP binding, so other implementations can be written against the same contract. It ships as prebuilt binaries for macOS, Linux and Windows; building from source is possible.

## 8. Versions, upgrades and extensions

- The `.layout` manifest records the layout version and each extension's version. Versions are semantic: only a MAJOR change can require a migration.
- Each extension states which layout versions it works with. The tool reports whether a repo is current, behind, ahead, unsupported or in conflict.
- Because files are the source of truth and the standard is a written spec, anyone can add an extension or write an alternative implementation without changing the core.

## 9. Rules of thumb

- Keep `AGENTS.md` short. Put anything long in a doc and link it.
- One source of truth per fact. If a list can be derived from headers, derive it.
- Never hand-edit a generated file.
- Code and tests win over prose. If a doc is wrong, fix it or say so.
- A well-formed record can still be wrong. Review, do not trust.
- Records are ordinary Obsidian notes: Obsidian reads and edits them without making them invalid (it drops comments and turns inline lists into block lists when it saves a property, which is cosmetic). A repository does not have to be a vault.
  - To get clickable references in Obsidian, write a quoted wikilink with an id alias (`parent: "[[261027-m3qz-fix-login|261027-m3qz]]"`) or set `reference_style: wikilink`; bare ids stay valid.
  - In the body of a note, write relative Markdown links so the note reads the same on a Git host. Obsidian's default wikilinks are tolerated, and the tool warns when one does not resolve. In a vault, turn off Obsidian's "Use Wikilinks" setting if you want it to write Markdown links.
  - Obsidian types each property by name for the whole vault, so if you open a repo inside a vault that already uses `id` as a number, Obsidian shows a type mismatch (the value is not changed).
  - Files whose names start with a dot (`docs/.layout`, `.agents/`) do not appear in Obsidian.
- Format with the pinned Prettier (config: `embeddedLanguageFormatting: off`), so frontmatter stays readable by the tool.
- Do not turn the flow into ceremony. The rails are a path for the agent, not a report for a meeting.

## 10. What is not proven

- No study found shows what history layers (ADRs, roadmaps, work logs) do for agents. One study found that context files can hurt task success. AKM is built to keep always-loaded text small, but its benefit is a hypothesis until pilots show it.
- The Kanban-style flow echoes a pattern described in public writing on long-running agents (one case study). That agents work better on it is a belief, not a result.
- The checks catch structure, not truth. A well-formed, wrong record passes every rule.
- Nothing forces an agent to run the checks; the standard asks, and the maintainer reviews.
- Parallel work is conducive, not guaranteed: shared files such as the roadmap can still conflict.
- Claiming an `active` item is an optional, advisory extension (never a lock), and the optional safety extension checks only known secret patterns and a deny list. Private prose and unusual secrets are not caught.
- The migration automation has so far been exercised only on a made-up older layout, because no layout before 0.2.0 exists.

## 11. Where things are

| What                                               | Where                                                                |
| -------------------------------------------------- | -------------------------------------------------------------------- |
| The rules: layout, formats, checks, versions       | [agent-ready-repo-spec](agent-ready-repo-spec.md)                    |
| The tool's operations, JSON, bindings and fixtures | [akm-api-spec](akm-api-spec.md)                                      |
| The tool                                           | the `akmctl` repository: prebuilt binaries, its README and changelog |

## Glossary

- **Cold reader:** an agent with no memory, or a maintainer back after months away.
- **Header layer:** the frontmatter of a record, read without its body.
- **Rails:** the simple flow of statuses and the stated next step that agents follow.
- **Routing index:** the generated table of ADRs a reader uses to find the right record.
- **Manifest:** the `.layout` file recording versions.
- **YAML-lite:** the frontmatter format: a strict subset of YAML 1.2 core, the schema Obsidian reads and writes, that the tool parses itself. It means exactly the same as YAML 1.2 where accepted, so Obsidian can edit the same files, and it errors on anything unsupported (nesting, multi-line values, numbers other than integers, tabs). Dates are plain `2026-10-28`; quote a value that looks like a number you mean as text.
- **Plain-file fallback:** doing a tool's job by reading the files directly.
