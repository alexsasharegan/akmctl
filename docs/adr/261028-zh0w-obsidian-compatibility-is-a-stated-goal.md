---
id: "261028-zh0w"
type: adr
adr_number: 20
title: "Obsidian Compatibility Is a Stated Goal of the AKM Standard"
status: accepted
created: 2026-10-28
updated: 2026-10-28
date: 2026-10-28
summary: Obsidian compatibility is a stated goal of AKM and the tiebreaker in design conflicts, though Obsidian itself is not required.
---

# ADR 261028-zh0w: Obsidian Compatibility Is a Stated Goal of the AKM Standard

## Context

The AKM standard ([[261027-ensf-akm-an-agent-knowledge-management]]) stores a repository's records as Markdown files with frontmatter. Its author practises personal knowledge management in Obsidian, and AKM is PKM's idea applied to a repository. Compatibility with Obsidian has shaped many decisions without ever being stated: the frontmatter format was aligned to YAML 1.2 core after a test showed that Obsidian rewrites a note's properties in its own canonical form (dates unquoted, comments dropped, inline lists made block lists), which would have made quoted dates fail validation on the first edit ([[agent-ready-repo-decision-log]], section F, Y26). The constraint is critical and was uncaptured.

## Decision

Compatibility with Obsidian is a stated goal of the standard and the tiebreaker where design choices conflict.

- **What it means:** AKM records are valid Obsidian notes. Obsidian reads them, shows their properties, and can edit them (through the Properties panel, its frontmatter API and common plugins) without making them invalid under the standard.
- **What it does not mean:** Obsidian is not required. A repository need not be a vault, and the tool works without Obsidian. The standard does not depend on Obsidian features, plugins or file layout beyond what Markdown and frontmatter already provide.
- **Tiebreaker:** when a rule would break Obsidian's reading or editing of a record, Obsidian's behaviour wins and the rule changes.
- **Verification:** when the frontmatter format or the record schema changes, test it against the current Obsidian by reading the files and by editing a property through its frontmatter API (the method is in the decision log, section F, using `obsidian eval`), and record the Obsidian version tested.

Rejected: ignoring Obsidian (breaks the maintainer's own use and the vault dogfooding); requiring Obsidian (excludes repositories that are not vaults); treating compatibility as case-by-case with no stated rule (it was, and it was lost).

## Consequences

- The frontmatter format is YAML-lite, a strict subset of YAML 1.2 core, the schema Obsidian reads and writes ([[261027-ensf-akm-an-agent-knowledge-management]] errata; spec section 0, H5 and H7).
- Each Obsidian release can change what it writes. The test above is repeated when the standard changes, and Obsidian's version is recorded.
- References between records can be wikilinks. A quoted string holding a wikilink (`parent: "[[261027-m3qz]]"`) is valid YAML 1.2 and Obsidian shows it as a link. The vault already does this in 179 notes (found with `rg '^\s*[a-z_]+:\s*"\[\[[^]]+\]\]"'`). Whether the standard allows, prefers or requires that form for `parent` and `blocked_by`, and how the tool resolves it, is not decided here.
- Areas still to check, not decided here: the relative-link rule (L1) against wikilinks in a vault; Obsidian's property types (for example `types.json` assigns the date type by property name, vault-wide); dot-files such as `docs/.layout` are hidden in Obsidian's file explorer; the migration of the vault itself ([[261027-nrwf-vault-agent-ready-migration]]).
- The standard stays usable in repositories with no Obsidian.

## Evidence

- The Obsidian test: [[agent-ready-repo-decision-log]], section F (Obsidian 1.14.5, 2026-10-28), and Y26.
- Forum report of the same rewrite behaviour: Obsidian forum, "YAML & Properties & API: processFrontMatter removes / alters string quotes, comments, types, formatting".

## Relations

- Related to:
  - [[261027-ensf-akm-an-agent-knowledge-management]] (the standard this constrains)
  - [[261023-mtpt-frontmatter-field-registry]] (frontmatter field registry: the vault's own field rules)
  - [[260912-bcdr-obsidian-cli-as-the-agent-interface]] (the Obsidian CLI as the agent interface: the test method)

## Errata

- 2026-10-28: the wikilink question under Consequences is settled in the decision log as Y30: a frontmatter reference may be a quoted wikilink with an id alias, bare ids stay canonical, and the tool reads both. Obsidian keeps such links unchanged when it edits a property.
- 2026-10-28: the relative-link question is settled as Y31: relative Markdown links are the standard and are checked by L1; wikilinks in bodies are tolerated and checked as warnings by L2.
- 2026-10-28: property types and hidden dot-files are settled as Y32: the standard ships no type settings, a vault resolves its own name conflicts (this vault's `id` is typed number), and dot-files stay invisible in Obsidian by design.
