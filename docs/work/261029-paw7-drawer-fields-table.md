---
id: "261029-paw7"
title: Show a record's frontmatter as a typed table in the drawer
status: ready
created: 2026-10-29
updated: 2026-10-29
kind: feature
priority: 3
summary: The drawer lists fields as plain text rows; type each value (string, date, number, link) and lay them out as a table.
---

# Show a record's frontmatter as a typed table in the drawer

## Goal

Alexander's feedback after using the workbench on the vault: the detail view is good, but the frontmatter fields could be processed into a better UI. Render the fields as a table where each value is typed and styled by its kind: string, date, number, boolean, list, and wikilink (special). Use the AST of the value, not string matching, so quoted wikilinks, aliases and lists are handled the same way everywhere.

Done when:

- [ ] The drawer shows the frontmatter as a table with one row per field.
- [ ] Each value is styled by its type, with a legend or label that does not rely on colour alone.
- [ ] Lists render one item per line, wikilinks as links (see the wikilink item).
- [ ] Unit tests cover the value typing; checked in the browser on the vault in light and dark.

## Log

- 2026-10-29: filed.
