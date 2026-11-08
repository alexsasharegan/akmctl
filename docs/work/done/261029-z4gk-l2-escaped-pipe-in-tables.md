---
id: "261029-z4gk"
title: "Rule L2: understand escaped pipes in wikilinks inside tables"
status: done
created: 2026-10-29
updated: 2026-11-08
closed: 2026-11-08
kind: bug
priority: 3
summary: A wikilink with an alias inside a Markdown table cell is written `[[target\|alias]]`; L2 reads the backslash as part of the target and reports it unresolved.
---

# Rule L2: understand escaped pipes in wikilinks inside tables

## Goal

Obsidian requires the pipe in a wikilink alias to be escaped inside a table cell: `[[target\|alias]]`. Rule L2 (and the mint_ids reference scan) treat the backslash as part of the target name and report it unresolved. Found while writing the migration ledger in the vault.

Done when:

- [x] L2 and the reference scans strip the escape before resolving.
- [x] Fixtures cover an alias link inside a table cell.
- [x] The tool does not report the escaped form as unresolved.

## Next

## Outcome

L2 and the mint_ids reference scan both read `[[name\|alias]]` in a table cell. The scan missed it because its pattern did not allow a backslash before the pipe. The basic fixture repository carries such a link, so check-1 covers it. To allow that, links and wikilinks in folders whose name starts with a dot are no longer scanned or counted.

## Log

- 2026-10-29: filed.
- 2026-10-29: rule L2 now strips the escape (internal/check/links.go, TestL2EscapedPipeInTableCell). Still open: the mint_ids reference scan and an API fixture.
- 2026-11-08 00:38-07:00: claimed (claude-sonnet-5-5).
- 2026-11-08 00:38-07:00: ticked "L2 and the reference scans strip the escape before resolving.": L2 strips the escape (internal/check/links.go, earlier fix); findReferences accepts \| after the name; TestFindReferencesSeesEscapedPipeInTableCell
- 2026-11-08 00:38-07:00: ticked "Fixtures cover an alias link inside a table cell.": repos/basic/docs/links.md holds an aliased link in a table cell; check-1 expects it clean; a missing target reports L2 without the backslash
- 2026-11-08 00:38-07:00: ticked "The tool does not report the escaped form as unresolved.": check on the fixture is clean (0 warnings); the reference scan lists the table-cell link
- 2026-11-08 00:38-07:00: Remaining part: the mint_ids reference scan ended its pattern at [\]|#] so name\| never matched; now accepts an optional backslash before the pipe. Other wikilink parsers checked: check/links.go trims the backslash, record/repo.go reads frontmatter values (no table cells). Added docs/links.md to the basic fixture (source, then sync).
