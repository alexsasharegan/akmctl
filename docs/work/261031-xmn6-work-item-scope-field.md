---
id: "261031-xmn6"
title: work item scope field
summary: "Spike: an optional scope list on work items and a rule that flags changes outside it"
status: raw
kind: spike
parent: "261031-v07x"
created: 2026-10-31
updated: 2026-10-31
---

# work item scope field

## Goal

From the consultation (gap 4). ADRs carry `scope` (paths and topic words). An unattended agent can edit files outside the task it was given. Find out whether a work item `scope` plus a check on the changed files is worth having. Done when:

- [ ] Options compared on a scratch repository: a `scope` field with a rule over `git status --porcelain` or a commit range (like `commits` for CM1 to CM3), an extension, or leave it to the harness. Include what counts as outside (the item's own file, the roadmap, generated files).
- [ ] A recommendation (build, defer or skip) is in `## Outcome`; a hard-to-reverse choice has a `proposed` ADR; follow-up items are filed.

## Log

- 2026-10-31: filed from the consultation.
