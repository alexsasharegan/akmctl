---
id: "261029-ctzc"
title: Make wikilinks in drawer fields clickable and open them in Obsidian
status: active
created: 2026-10-29
updated: 2026-10-31
kind: feature
priority: 2
summary: project and origin hold wikilinks such as Vault Views; show them as links that open the target in Obsidian.
claimed_by: claude
claimed_at: 2026-10-31T07:12
---

# Make wikilinks in drawer fields clickable and open them in Obsidian

## Goal

In the drawer, project: `[[Vault Views]]` and origin: `[[261025-g68j-build-pkm-workbench]]` are plain text. Parse a quoted wikilink (with or without an alias, in a field or inside a list) and render it as a link. Clicking opens the target note in Obsidian, the same action as Open in Obsidian. If the target is an AKM record, offer to open it in the workbench drawer instead (open question: which is the default; Alexander asked for the Obsidian open first).

Done when:

- [ ] Wikilinks in any frontmatter field are rendered as links, aliases show the alias text.
- [ ] A click opens the note in Obsidian; a link whose target does not resolve is shown as unresolved and is not clickable.
- [ ] The resolution uses the API (find by name or id), not a guess from the file name.
- [ ] Tests for parsing (plain, alias, list, unresolved) and a browser check on the vault.

## Log

- 2026-10-29: filed.
- 2026-10-31: claimed (claude).
