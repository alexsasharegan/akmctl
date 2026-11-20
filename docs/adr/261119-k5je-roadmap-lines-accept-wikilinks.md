---
id: "261119-k5je"
title: Roadmap lines accept a wikilink in place of the bracketed id
summary: A roadmap line may be a wikilink with an id alias, so Obsidian shows a link; check and the loader read both forms and the bracketed id stays valid
status: proposed
date: 2026-11-19
---

# Roadmap lines accept a wikilink in place of the bracketed id

## Context

A roadmap line was `- [id] text` and nothing else. The specification said roadmap lines use the bare id, rule RM1 matched only `[id]`, and the loader read only that form. In a repository with `reference_style: wikilink` (frontmatter references are written as wikilinks there) the roadmap was the one place that stayed plain text, so Obsidian showed no link. An agent that wrote `- [[261027-m3qz-fix-login|261027-m3qz]]` got an RM1 error from `check`, so it wrote the bracketed form to make the check pass (GitHub issue 96). The tool never writes the roadmap, so no renderer was at fault; the rule was.

## Considered Options

- Keep the bracketed id only. Rejected: the roadmap is the page a maintainer reads most, and in Obsidian it cannot be used to open an item.
- Accept a wikilink in place of the bracketed id, with the id alias. Chosen.
- Accept a wikilink only when `reference_style` is `wikilink`. Rejected: it makes `check` depend on a setting to read a line, and a repository that moves between styles would see errors in a file nobody edited.
- Let the tool render the roadmap from the work items. Rejected here: the roadmap's order is a human decision and the tool does not write it (a larger change, and not needed for this).

## Decision

A roadmap line may be `- [id] text` or `- [[target|id]] text`, where the target is the record's filename without `.md` and the alias keeps the id visible, the form frontmatter already uses. The id of a wikilink is the leading id of its target; the slug and the alias are not read, so a renamed slug is not an error. Both forms are valid in every repository. A repository with `reference_style: wikilink` should write the wikilink form.

One reader, `record.RoadmapEntry`, serves the loader and rule RM1, so they cannot disagree. A wikilink whose target has no leading id is reported by RM1 like any unresolved entry. Nothing else about the roadmap changes: sections, order, RM2 to RM4 and the `roadmap` sort key read ids and are unaffected.

## Consequences

- A repository may write wikilinks in its roadmap and pass `check`; one that writes bracketed ids needs no change.
- A tool older than this change reports RM1 on a wikilink line, so a repository that adopts the form needs the new tool.
- The specification (sections 0 and 5, rule RM1), the guide and the glossary change; the work item and the changelog name this decision. It is an additive change to the layout.
- The slug in the target is not checked, so a wikilink can point at a stale filename that Obsidian then shows as unresolved; the existing link rule L1 reports an unresolved wikilink in the roadmap file.
- Status `proposed`: nothing here is accepted until the pull request is merged.
