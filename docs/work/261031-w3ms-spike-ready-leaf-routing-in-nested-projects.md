---
id: "261031-w3ms"
title: spike ready leaf routing in nested projects
status: raw
created: 2026-10-31
updated: 2026-10-31
kind: spike
summary: How should next-work lookup behave when the id it is given is a project with nested projects?
---

# spike ready leaf routing in nested projects

## Goal

From spike 261031-078w. Idea: `ls --next --under <id>` returns the work to take next inside a project. The open question is nesting. Compare: (a) flatten: the ready leaves of the whole subtree, ordered by roadmap, `priority`, id, with each row's path of ancestors; (b) one level: the direct children only, with a project marked as a project and its own roll-up count, so the caller descends; (c) a mix, flattened by default with `--depth <n>` or `--level`. Also decide: whether parents are ever returned (a parent with all children done but itself open); how blocked children, and children blocked through a parent, are shown or hidden; how claimed items count; what an agent receives as its single next item; the API and JSON shape. Depends on `ls --descendants` (261031-qwhn) and the recursive predicates (261031-xprq). Done when:

- [ ] Each option is run on a three-level scratch tree, with sample output, for an agent picking work and a maintainer reviewing a project.
- [ ] A recommendation and a command shape are in `## Outcome`; follow-up items are filed for it.

## Log

- 2026-10-31: filed.
- 2026-10-31: filed from spike 261031-078w.
