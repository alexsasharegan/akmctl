---
id: "261101-kngg"
title: project as parent alias
status: done
created: 2026-11-01
updated: 2026-11-01
closed: 2026-11-01
---

# project as parent alias

## Goal

The extra frontmatter field `project` is a leftover from the pkm vault this repository's work items came from (three items still carry `project: Project Notes`, and the migrator already copies `project` to `parent`). It means what `parent` means: the item a piece of work belongs to. Grouping related work under one goal is a project, so the word for a parent that has children is **project**; the field stays `parent`. No ADR: this corrects a migration leftover and the record is the account of it. Done when:

- [x] "Umbrella" is struck from the glossary, both copies of the `akm-track-work` skill, `docs/design/decomposing-work.md` and the open work items' prose; the Parent row defines project.
- [x] The workbench filter narrows by project (a parent and its descendants) read from `parent`, not from the `project` field, with the logic tests updated; the glossary row "Project filter" no longer says it is not `parent`.
- [x] The three items carrying `project:` lose the field.
- [x] CHANGELOG entry under Unreleased; `akmctl check`, the Go and JS tests and `format:check` pass.

## Next

## Outcome

project names a parent with children; the field stays parent. The workbench filter now groups by parent, the three leftover project: fields are gone, umbrella is struck from glossary, skills (both copies), design note and open items. Closed items and accepted ADRs keep their old wording. No ADR: this corrected a migration leftover.

## Log

- 2026-11-01: filed.
- 2026-11-01 20:26-07:00: detached from 261101-vs5c (not a rails gap). Findings: 'umbrella' is a glossary term (Parent row) and appears in the akm-track-work skill, its template and docs/design/decomposing-work.md; the workbench 'Project filter' reads a separate frontmatter field 'project'; the maintainer wants to strike 'umbrella'. Options: (a) project = prose word for a parent with children, field stays parent; (b) project accepted as a field/flag alias; (c) rename parent to project.
- 2026-11-01 20:32-07:00: claimed (claude-sonnet-5-5-kngg).
