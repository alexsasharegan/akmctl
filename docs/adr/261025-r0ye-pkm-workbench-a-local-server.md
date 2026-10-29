---
id: "261025-r0ye"
type: adr
adr_number: 18
title: "pkm-workbench: A Local Server over the Vault's Files"
status: accepted
created: 2026-10-25
updated: 2026-10-26
authored_by: agent
date: 2026-10-25
summary: pkm-workbench is a local Go server in its own repo that serves task views from the vault's files and writes changes back to the frontmatter.
---

# ADR 261025-r0ye: pkm-workbench: A Local Server over the Vault's Files

## Context

The task views in Tasks (`Tasks.base`) render only inside Obsidian, and Bases has no board or drag-to-change-status that is known to work. On 2026-10-25 native cards and list views were tried on `Tasks.base`: useful, possibly as a phone view, but the look did not scale and Alexander wanted a UI with more control. He also wants to move things (change a status, close a task) from the view itself, and to give raw captures a lane before they become tasks. The design spike (261024-k2rc-local-views-server-design) worked through the options; the full reasoning is in [pkm-workbench-design](../design/pkm-workbench-design.md).

## Decision

Build `pkm-workbench` (short names `pkm-wb`, `pkmwb`): a local server, in its own repo, that serves task views from the vault's files and writes changes back to the YAML frontmatter. The files stay the source of truth.

- **Repo and stack.** `~/github.com/alexsasharegan/pkm-workbench`, outside the vault. Go, standard library only (a dependency comes in only when it earns its place, with the reason logged). Templates, CSS and JS are embedded in one binary. The UI is designed from scratch and shares no code with Vault Views.
- **Scope.** Version 1 serves notes with `type: task` in `9 pkm/tasks/`, read per request, no file watcher.
- **Views.** The server defines its own views. `Tasks.base` stays for Obsidian. Both implement one authoritative view spec, which moves into the Tasks hub, and a check compares their lists on the same files.
- **Writes.** Direct file writes (not the `obsidian` CLI, which needs Obsidian running). One frontmatter line is edited in place, with no YAML re-serialising. A hash of the file, taken when the page was rendered, is checked before the write, which is refused on a mismatch. The write goes through a temp file and a rename. Closing a task is the one structured edit: it requires an outcome, puts `## Outcome` above `## Log`, deletes `## Next`, and sets `status` and `updated_at`.
- **Security.** Bind `127.0.0.1` only. A per-run token and `Origin` and `Host` checks on reads and writes. No phone access in version 1.
- **Inbox lane.** A new task status `raw` for unrefined captures, shown as its own lane. `inbox` keeps its meaning: filed, not claimed.

Rejected: parsing `.base` files and Bases formulas; the `obsidian` CLI as the write path; reusing the Vault Views kit and renderer (static pages, and Alexander dislikes its look); Node, Python or Rust; a repo shared with the Vault Views fork; a file watcher in version 1.

## Consequences

- There are two implementations of the task views to keep in step. The spec is authoritative; change it first.
- The server writes into the vault. Git is the undo and the autosave picks the edits up. A note open in Obsidian with unsaved edits can race with a server write; the hash check covers the common case and the rest is accepted for a single-user tool.
- `raw` must be added to the status vocabulary in Tasks and to the field registry (frontmatter-fields) when it is built.
- A closed task keeps its `adr` field unset until the ADR check runs. The close cascade and an ADR check dialog are pinned as later work.
- Task content, which can hold private detail, is served over localhost, so the token and Origin checks are required, not optional.

## Evidence

- Tasks: 261024-k2rc-local-views-server-design (the spike, closed), [261025-g68j-build-pkm-workbench](../work/done/261025-g68j-build-pkm-workbench.md) (the build).
- Design: [pkm-workbench-design](../design/pkm-workbench-design.md).
- Commits: `1a0de22` (Bases views, checkpoint before: `3008564`), `fcb1f3a` (design note), `698aabc` (spike closed).

## Relations

- Related to:
  - 261020-5be4-base-and-hub-every-database-gets-a-hub (base and hub: the views live over the same task files and the hub carries the spec)
  - 261023-f7rd-task-queue-and-the-task-lifecycle-skill (task queue and the task-lifecycle skill: the close rules the server must follow)
  - 260912-bcdr-obsidian-cli-as-the-agent-interface (the Obsidian CLI as the agent interface: not used for server writes)
  - 241103-b2mb-git-for-history-only-no-remote-cloud (git for history only: the undo for server writes)

## Errata

- 2026-10-29: copied from the maintainer's notes into this repository. Wikilinks to records copied with it became relative links; wikilinks to notes that stayed behind became plain text. The decision is unchanged.
