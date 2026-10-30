---
id: "261029-kb88"
title: create obsidian board
summary: a cmd to generate an obsidian base which mirrors the workbench ui
status: raw
kind: feature
created: 2026-10-29
updated: 2026-10-29
---

# create obsidian board

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

## Log

- 2026-10-29: filed.
