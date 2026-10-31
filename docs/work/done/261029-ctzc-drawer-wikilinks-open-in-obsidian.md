---
id: "261029-ctzc"
title: Make wikilinks in drawer fields clickable and open them in Obsidian
status: done
created: 2026-10-29
updated: 2026-10-31
closed: 2026-10-31
kind: feature
priority: 2
summary: project and origin hold wikilinks such as Vault Views; show them as links that open the target in Obsidian.
---

# Make wikilinks in drawer fields clickable and open them in Obsidian

## Goal

In the drawer, project: `[[Vault Views]]` and origin: `[[261025-g68j-build-pkm-workbench]]` are plain text. Parse a quoted wikilink (with or without an alias, in a field or inside a list) and render it as a link. Clicking opens the target note in Obsidian, the same action as Open in Obsidian. If the target is an AKM record, offer to open it in the workbench drawer instead (open question: which is the default; Alexander asked for the Obsidian open first).

Done when:

- [x] Wikilinks in any frontmatter field are rendered as links, aliases show the alias text.
- [x] A click opens the note in Obsidian; a link whose target does not resolve is shown as unresolved and is not clickable.
- [x] The resolution uses the API (find by name or id), not a guess from the file name.
- [x] Tests for parsing (plain, alias, list, unresolved) and a browser check on the vault.

## Outcome

Wikilinks in any frontmatter field (alone or in a list) are links in the drawer: wikilinkParts in logic.js parses plain, alias, heading and unresolved forms (tested under node), detail.js shows the alias, resolves through the API and opens the note in Obsidian, adds a "here" link for a record, and marks a missing or ambiguous name as unresolved and not clickable. The lookup is a new `name` parameter on find (internal/api/ops_find.go; spec 4.8, find.json, fixtures find-2 to find-4, the akm reference, tests including hidden folders, ambiguity and path traversal). Checked in the workbench on a fixture: a project link, an alias link to a record (with the here link), an unresolved name and an alias in a list. Not checked: the real vault, large vaults (the lookup walks the repository on each new name and is remembered until the records reload), and clicking through to Obsidian itself. Open question from the item left as built: the Obsidian link is the default and the drawer link is secondary. Decision: extending find with name rather than adding an operation (maintainer chose it in chat 2026-10-31); no ADR written.

## Log

- 2026-10-29: filed.
- 2026-10-31: claimed (claude).
