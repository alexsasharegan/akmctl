---
id: "261029-z4gk"
title: "Rule L2: understand escaped pipes in wikilinks inside tables"
status: ready
created: 2026-10-29
updated: 2026-10-29
kind: bug
priority: 3
summary: A wikilink with an alias inside a Markdown table cell is written `[[target\|alias]]`; L2 reads the backslash as part of the target and reports it unresolved.
---

# Rule L2: understand escaped pipes in wikilinks inside tables

## Goal

Obsidian requires the pipe in a wikilink alias to be escaped inside a table cell: `[[target\|alias]]`. Rule L2 (and the mint_ids reference scan) treat the backslash as part of the target name and report it unresolved. Found while writing the migration ledger in the vault.

Done when:

- [ ] L2 and the reference scans strip the escape before resolving.
- [ ] Fixtures cover an alias link inside a table cell.
- [ ] The tool does not report the escaped form as unresolved.

## Log

- 2026-10-29: filed.
- 2026-10-29: rule L2 now strips the escape (internal/check/links.go, TestL2EscapedPipeInTableCell). Still open: the mint_ids reference scan and an API fixture.
