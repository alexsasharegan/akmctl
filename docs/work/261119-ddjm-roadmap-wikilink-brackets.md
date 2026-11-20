---
id: "261119-ddjm"
title: Roadmap wraps ids in single brackets when wikilinks is true
status: raw
kind: bug
created: 2026-11-19
updated: 2026-11-19
---

# Roadmap wraps ids in single brackets when wikilinks is true

## Goal

With `reference_style: wikilink`, a roadmap line can be a wikilink with the id as its alias and still pass `check`. Today the reader and rule RM1 accept only `- [id]`, the specification says roadmap lines use the bare id, and a `- [[target|id]]` line fails RM1, so an agent that must make `check` pass writes the single-bracket form and Obsidian shows no link. GitHub issue 96.

Provisional: the decision is drafted as a `proposed` ADR and is not accepted.

Done when:

- [ ] A roadmap line `- [[<filename>|<id>]]` (and `- [[<filename>]]`) is read as that id by the loader and by `check`; `- [id]` still works; a wikilink whose target has no leading id still fails RM1.
- [ ] The specification, the guide, the skills and their `init` templates say a roadmap line may be either form, and which one to write under `reference_style: wikilink`.
- [ ] A fixture or test covers both forms, in the loader, `check` and `ls --sort roadmap`.
- [ ] A `proposed` ADR records the choice, and the changelog has an entry.

## Log

- 2026-11-19: filed.
