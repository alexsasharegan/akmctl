---
id: "261101-kngg"
title: project as parent alias
status: raw
created: 2026-11-01
updated: 2026-11-01
parent: 261101-vs5c
---

# project as parent alias

## Goal

The maintainer says "project" in natural language for what AKM calls `parent`, and wants "project" accepted as an alias of `parent` because it matches common speech. Open questions to settle first: where the alias applies (frontmatter field, `new`/`set` flags, `ls` filters, the workbench, the spec); and a collision, since the glossary already defines a workbench "Project filter" on an extra frontmatter field `project` that "is not AKM's `parent`". Done when:

- [ ] The question is answered in an ADR or a glossary change: alias, rename, or keep the two apart.
- [ ] The tool and spec accept the chosen form, with a test and a CHANGELOG entry.

## Log

- 2026-11-01: filed.
