---
id: "261101-xher"
title: "scanner: ordered-list tasks are refused"
status: raw
created: 2026-11-01
updated: 2026-11-01
blocked_by: [261101-4j86]
kind: feature
parent: 261031-gtw6
summary: "Locate refuses an ordered-list task (1. [ ] x); 2 of the vault's work items have one. Decide: read it as a list item, or keep refusing (local)"
---

# scanner: ordered-list tasks are refused

## Goal

Found by the refusal census in ADR [261101-tr8e](../adr/261101-tr8e-body-writes-through-the-section-tree.md) (principle 4: every kind of refusal met on real records gets a backlog item). The scanner refuses an ordered-list task such as `1. [ ] step` because goldmark reads the box differently in ordered lists. Two of the maintainer's vault work items have one, both in a section other than the one a tick would edit, so a write only refuses when the edited section holds one (local refusal). Done when:

- [ ] A decision, recorded in this item's Outcome: claim `1. [ ]` as a list item that is not a checklist item (and say so in the glossary), or keep refusing it as a local construct with a reason in the code.
- [ ] If claimed: the oracle battery, the generator and the fuzz target agree, and the refusal row in `TestRefusals` is removed or retagged.

## Next

Not started; waits on 261101-4j86 (the structural/local classification).

## Log

- 2026-11-01: filed.
- 2026-11-01: filed from the refusal census (ADR 261101-tr8e, principle 4).
