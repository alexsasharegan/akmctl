---
id: "261029-kb88"
title: Generate an Obsidian Bases file that mirrors the workbench
summary: A command that writes an Obsidian Bases file for the work items and ADRs, mirroring the workbench views; written once and never overwritten
status: ready
kind: feature
created: 2026-10-29
updated: 2026-10-29
priority: 1
---

# Generate an Obsidian Bases file that mirrors the workbench

## Goal

Build an `akmctl` subcommand to generate an Obsidian Base which presents a mirror view from what the `akmctl ui` shows.

An example Obsidian Base:

```yaml
filters:
  and:
    - file.inFolder("docs/work")
formulas:
  age_days: if(updated, (today() - updated).days, "")
  stale: if(updated, (today() - updated).days > 30, true)
properties:
  formula.age_days:
    displayName: days since update
views:
  - type: table
    name: Open by parent
    filters:
      and:
        - status != "done"
        - status != "dropped"
    order:
      - parent
      - file.name
      - status
      - formula.age_days
    sort:
      - property: parent
        direction: ASC
      - property: status
        direction: ASC
    columnSize:
      note.parent: 282
      file.name: 547
      formula.age_days: 181
  - type: table
    name: Ready and active
    filters:
      or:
        - status == "ready"
        - status == "active"
    order:
      - file.name
      - status
      - parent
      - updated
  - type: table
    name: Blocked
    filters:
      and:
        - status == "blocked"
    order:
      - file.name
      - parent
      - updated
  - type: table
    name: Stale
    filters:
      and:
        - status != "done"
        - status != "dropped"
        - formula.stale
    order:
      - file.name
      - status
      - parent
      - formula.age_days
  - type: table
    name: Recently closed
    filters:
      and:
        - status == "done" || status == "dropped"
        - (today() - updated).days <= 14
    order:
      - file.name
      - status
      - parent
      - decisions
      - updated
  - type: cards
    name: Cards by status
    filters:
      and:
        - status != "done"
        - status != "dropped"
    groupBy:
      property: status
      direction: ASC
    order:
      - file.name
      - parent
      - formula.age_days
    sort:
      - property: updated
        direction: DESC
  - type: cards
    name: Cards by parent
    filters:
      and:
        - status != "done"
        - status != "dropped"
    groupBy:
      property: parent
      direction: ASC
    order:
      - file.name
      - status
      - formula.age_days
    sort:
      - property: status
        direction: ASC
  - type: list
    name: List ready and active
    filters:
      or:
        - status == "ready"
        - status == "active"
    groupBy:
      property: parent
      direction: ASC
    order:
      - file.name
      - status
  - type: kanban
    name: Kanban by status
    filters:
      and:
        - status != "done"
        - status != "dropped"
    groupBy:
      property: status
      direction: ASC
    groupOrder:
      - raw
      - ready
      - active
      - blocked
    order:
      - file.name
      - parent
      - formula.age_days
```

Decisions (Alexander, 2026-10-29): the file is written once and never overwritten, like `init`; the user then edits it freely in Obsidian.

The Obsidian Bases file mirrors the workbench: table views for open work, ready and active, blocked, stale (by `updated` against the manifest's `stale_days`) and recently closed (by `closed`); cards grouped by status and by parent or cycle; and an ADR table. Folder filters come from the manifest's `work_dir` and `adr_dir`, and the properties are the AKM fields (id, title, status, kind, priority, parent, cycle, keep, created, updated, closed, claimed_by). The example above is the starting shape, with the vault-specific names replaced.

It is an API operation like the others (schema, fixtures, CLI and HTTP), with a name to settle (`akmctl obsidian` is the working name).

Open questions to answer first:

- Does Obsidian Bases support a `kanban` view type in core, or only through a plugin? If not, the board is cards grouped by status with a fixed group order.
- Obsidian Bases' file format may change between Obsidian releases; generate a documented subset and say which Obsidian version it was checked against.
- A vault can be a parent of the repository, so folder filters are relative to the vault root, not the repository root. The command needs to find the vault root (the nearest `.obsidian` folder upward) or take it as an option.
- Whether this needs an ADR (a generated file that tracks another product's format).

Done when:

- [ ] the open questions are answered in this item's Log
- [ ] the operation is in the API spec with a schema and fixtures, and the CLI command writes the file and refuses to overwrite an existing one
- [ ] the generated file parses as YAML, every filter and property names an AKM field, and the folder filters use the manifest's directories
- [ ] Alexander opens it in Obsidian and each view renders (manual check; the tool cannot test this)
- [ ] the README and docs/glossary.md name it as Obsidian Bases

## Next

Ready. Start by answering the open questions, because they decide the file's shape.

## Log

- 2026-10-29: filed.
- 2026-10-29: refined with Alexander.
