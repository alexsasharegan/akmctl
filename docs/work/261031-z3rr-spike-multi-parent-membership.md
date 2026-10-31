---
id: "261031-z3rr"
title: spike multi parent membership
status: raw
created: 2026-10-31
updated: 2026-10-31
kind: spike
summary: Can an item belong to more than one project, and what does it cost the rules?
---

# spike multi parent membership

## Goal

From spike 261031-078w. `parent` holds one value, so an item belongs to one project. Find out whether to allow more, and at what cost. Questions: a list-valued `parent` versus a second field (for example `part_of`) that is membership without ownership; what "closable" means for a shared child (does it count toward every parent, and may one parent close while the other is open); how WK4 cycle checks and the union dependency graph change; how roll-up counts avoid counting a shared item twice; the effect on the YAML-lite field grammar, the API shapes and the workbench; whether `blocked_by` plus the `project` field already covers the need. Done when:

- [ ] Each option has a worked example on a scratch tree, with the rule changes it needs.
- [ ] A recommendation (build, defer or skip) is in `## Outcome`; a hard-to-reverse choice has a `proposed` ADR.
- [ ] Follow-up items are filed for what is recommended.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
