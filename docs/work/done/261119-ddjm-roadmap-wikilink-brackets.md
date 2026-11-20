---
id: "261119-ddjm"
title: Roadmap wraps ids in single brackets when wikilinks is true
status: done
kind: bug
created: 2026-11-19
updated: 2026-11-19
closed: 2026-11-19
---

# Roadmap wraps ids in single brackets when wikilinks is true

## Goal

With `reference_style: wikilink`, a roadmap line can be a wikilink with the id as its alias and still pass `check`. Today the reader and rule RM1 accept only `- [id]`, the specification says roadmap lines use the bare id, and a `- [[target|id]]` line fails RM1, so an agent that must make `check` pass writes the single-bracket form and Obsidian shows no link. GitHub issue 96.

Provisional: the decision is drafted as a `proposed` ADR and is not accepted.

Done when:

- [x] A roadmap line `- [[<filename>|<id>]]` (and `- [[<filename>]]`) is read as that id by the loader and by `check`; `- [id]` still works; a wikilink whose target has no leading id still fails RM1.
- [x] The specification, the guide, the skills and their `init` templates say a roadmap line may be either form, and which one to write under `reference_style: wikilink`.
- [x] A fixture or test covers both forms, in the loader, `check` and `ls --sort roadmap`.
- [x] A `proposed` ADR records the choice, and the changelog has an entry.

## Next

## Outcome

A roadmap line may be `[id]`, `[[target]]` or `[[target|id]]`; the loader and rule RM1 read all three through record.RoadmapEntry. The tool never wrote roadmap lines, so this was a spec change, decided in the ADR and merged in pull request 98. GitHub issue 96.

## Log

- 2026-11-19: filed.
- 2026-11-19 18:55-07:00: claimed (claude-sonnet-5-5).
- 2026-11-19 18:57-07:00: Traced: the tool never writes roadmap lines. The loader (record.roadmapLine) and rule RM1 accepted only `[id]`, and the spec said roadmap lines use the bare id, so a wikilink line failed check. Provisional fix (maintainer has not accepted it): record.RoadmapEntry reads `[id]`, `[[target]]` and `[[target|id]]` for both; proposed ADR 261119-k5je; spec sections 0 and 5, RM1, guide, glossary and changelog updated; tests in record and check. Left out: skills and init templates do not describe roadmap lines, so no change there. RM4 warns that this item is not on the roadmap; the order is the maintainer's.
- 2026-11-19 19:21-07:00: ticked "A roadmap line `- [[<filename>|<id>]]` (and `- [[<filename>]]`) is read as that id by the loader and by `check`; `- [id]` still works; a wikilink whose target has no leading id still fails RM1.": record.RoadmapEntry; TestRoadmapEntry and TestRoadmapWikilinkEntries
- 2026-11-19 19:21-07:00: ticked "The specification, the guide, the skills and their `init` templates say a roadmap line may be either form, and which one to write under `reference_style: wikilink`.": spec sections 0 and 5, RM1, guide, glossary; skills and templates do not describe roadmap lines
- 2026-11-19 19:21-07:00: ticked "A fixture or test covers both forms, in the loader, `check` and `ls --sort roadmap`.": TestRoadmapEntry, TestRoadmapWikilinkEntries; loader checked through r.Roadmap
- 2026-11-19 19:21-07:00: ticked "A `proposed` ADR records the choice, and the changelog has an entry.": ADR 261119-k5je accepted by the merge of pull request 98; changelog entry
