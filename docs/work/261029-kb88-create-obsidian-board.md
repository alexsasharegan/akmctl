---
id: "261029-kb88"
title: Generate an Obsidian Bases file that mirrors the workbench
summary: A command that writes an Obsidian Bases file for the work items and ADRs, mirroring the workbench views; written once and never overwritten
status: active
kind: feature
created: 2026-10-29
updated: 2026-11-07
priority: 1
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-08T04:06
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

- [x] the open questions are answered in this item's Log
- [x] the operation is in the API spec with a schema and fixtures, and the CLI command writes the file and refuses to overwrite an existing one
- [x] the generated file parses as YAML, every filter and property names an AKM field, and the folder filters use the manifest's directories
- [ ] Alexander opens it in Obsidian and each view renders (manual check; the tool cannot test this)
- [x] the README and docs/glossary.md name it as Obsidian Bases

## Next

Alexander opens `docs/akm.base` in Obsidian (this repository opened as a vault) and says whether each view renders: the open, ready and active, blocked, stale and recently closed tables, the three cards views, the Kanban view and the ADR table. Fix what does not, then close. If `file.name != "README.md"` or the `closed` date arithmetic is wrong in Obsidian, those are the likeliest.

## Log

- 2026-10-29: filed.
- 2026-10-29: refined with Alexander.
- 2026-11-07 21:06-07:00: claimed (claude-sonnet-5-5).
- 2026-11-07 21:13-07:00: Open questions answered. Kanban is a core Bases layout (Obsidian 1.14; list 1.10, table and cards 1.9), so the board is a real Kanban view with a fixed group order. The file is YAML with filters, formulas, properties and views, and the subset used is the one the Obsidian help documents, written for 1.14; akmctl cannot open Obsidian, so the render check is Alexander's. The vault root is the nearest .obsidian folder above the repository (the repository root when there is none), and folder filters are relative to it. No ADR: the file is written once, never overwritten and easy to delete, so the choice is not hard to reverse.
- 2026-11-07 21:13-07:00: Built the bases operation: API spec section 4.25, schema, four fixture cases and the obsidian fixture repository, the CLI and POST /bases, skill reference, README and glossary. Chose the name bases and the default docs/akm.base with Alexander. Wrote docs/akm.base in this repository for the manual check.
- 2026-11-07 21:13-07:00: ticked "the open questions are answered in this item's Log": answers are in the Log
- 2026-11-07 21:13-07:00: ticked "the operation is in the API spec with a schema and fixtures, and the CLI command writes the file and refuses to overwrite an existing one": spec 4.25, schemas/bases.json, cases bases-1 to 4, CLI and POST /bases; go test and npm run schemas pass
- 2026-11-07 21:13-07:00: ticked "the generated file parses as YAML, every filter and property names an AKM field, and the folder filters use the manifest's directories": TestBasesFileParsesAndNamesOnlyAKMFields parses the file and checks every property; TestBasesFileUsesTheManifestsFoldersAndStaleDays checks the folders
- 2026-11-07 21:13-07:00: ticked "the README and docs/glossary.md name it as Obsidian Bases": README section Obsidian Bases and the glossary row bases
