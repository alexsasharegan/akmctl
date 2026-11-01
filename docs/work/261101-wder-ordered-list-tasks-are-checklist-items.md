---
id: "261101-wder"
title: ordered list tasks are checklist items
status: ready
created: 2026-11-01
updated: 2026-11-01
parent: 261031-gtw6
summary: Read 1. [ ] and 1) [x] as checklist items so tick and the drawer work on ordered steps; the maintainer wants ordered task lists
---

# ordered list tasks are checklist items

## Goal

The maintainer wants ordered task lists (`1. [ ] step`) as legitimate checklist items, for steps that happen in order (2026-11-01). Work item 261101-xher kept them as a local refusal only because claiming them is a model change; this item makes the change.

Design to settle first:

- Position: a checklist position counts every task in the section in order, bullet or ordered, so `tick <id> <n>` stays one number. Or ordered boxes count separately; pick one and state it in the glossary and the `tick` reference.
- Markers: `1.` and `1)`, any number of digits (CommonMark allows at most nine). Nested markers such as `0) 0) [X]` stay refused.
- The oracle (goldmark, GFM) already reports these as tasks, so the Scanner claims what the oracle already says; the battery, generator and both fuzz targets must agree.
- The writers keep working: a tick changes one byte, and the drawer and `Checklist` report the item. Whether the tool renumbers on insert is out of scope; it never edits the number.

Done when:

- [ ] A decision on position counting is in this item's Outcome (an ADR only if it rejects a real alternative).
- [ ] `Locate` claims `1. [ ] x` as a Task; `Checklist`, `SetChecked` and `tick` work on it, in LF and CRLF files; the refusal row in `TestRefusals` and the `ordered-list task` construct are removed.
- [ ] The oracle battery, generator and both fuzz targets agree; fixtures `section-7` and `check-8`, which use an ordered task as their refusal example, are regenerated with another local construct and a new `tick` fixture shows an ordered task.
- [ ] Glossary, the `tick` reference and the API spec say what counts as a checklist item.

## Next

Not started. Read internal/body/scanner.go (the orderedTask case and its comment) and tree.go Checklist/SetChecked; then settle position counting.

## Log

- 2026-11-01: filed.
