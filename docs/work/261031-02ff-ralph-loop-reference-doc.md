---
id: "261031-02ff"
title: ralph loop reference doc
summary: "A reference loop in docs that uses only ls, claim, section, check and close"
status: raw
kind: chore
parent: "261031-v07x"
created: 2026-10-31
updated: 2026-10-31
---

# ralph loop reference doc

## Goal

From the consultation. The outer loop (task master, builder, supervisor) is a script, not part of `akmctl`. What belongs here is a short reference showing the loop with the tool's own commands, so people see how the records are meant to drive one. See [ralph-harness](../ideas/ralph-harness.md). Done when:

- [ ] A page (or skill reference) lists the loop step by step with the exact commands and exit conditions, and says which parts are the harness's and which the tool's.
- [ ] Every command in it is run once against a scratch repository, and the page says that it was.
- [ ] It states what the loop cannot guarantee (see gaps 1, 2 and 4 in the design note).

## Log

- 2026-10-31: filed from the consultation.
