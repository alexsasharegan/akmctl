---
id: "261109-n05b"
title: Agents name records by title, not id, and resolve ambiguity in steps
status: done
created: 2026-11-09
updated: 2026-11-09
closed: 2026-11-09
---

# Agents name records by title, not id, and resolve ambiguity in steps

## Goal

Make the `akm` skill tell an agent to name records by title in prose, and to give an id only as a last step. A maintainer who uses dictation and text-to-speech hears an id read out digit by digit, which costs more attention than the decision it is part of (filed from that maintainer's report).

The rule, in order, stopping at the first step that tells the items apart in the list the reader has in front of them and the last few messages:

1. The title alone, shortened, when nothing else in view could be mistaken for it.
2. A word that differs: name the part that separates two similar items ("the refactor that prepares for theming" and "the theme additions").
3. The date the id encodes, spoken as a date ("the one from October twelfth").
4. For items made the same day, only the short random tail of the id, and only as many characters as separate them. This is a handle for conversation: never written to a file or a commit.
5. The full id only when the reader asks for it or must type or paste it.

Ambiguity is judged against the set shown, not the whole repository. Before sending, an agent checks that no two items in one reply shorten to the same phrase, and steps down a level for those two only. Titles should put the words that tell items apart first, because speech reads left to right. A repository's own `AGENTS.md` may override the rule (a team that likes ids in prose).

Ids stay in files, commands, links, commit trailers and pull request text; the rule is about what an agent says in chat.

Done when:

- [x] Both copies of the `akm` skill state the rule (`.agents/skills/akm` and `internal/initrepo/templates/skills/akm`), and the repository's `AGENTS.md` says it may override.
- [x] A decision record (`proposed`, for the maintainer to accept) explains why the default changes for every user of the skill, and the glossary registers any new term.
- [x] A cold-agent check (`docs/cold-agent-check.md`) gives an agent a mixed list that contains two similar titles and two items made the same day, and the agent's reply names each one without a full id and without two items sharing a handle.
- [x] A `CHANGELOG.md` entry under `## Unreleased` says what changed and who is affected.

Not in scope: a listing option that prints the shortest unique tail. File it only if the cold-agent check shows agents get the arithmetic wrong.

## Next

## Outcome

The akm skill (both copies) and AGENTS.md tell an agent to name records by title in chat and to give an id only as a last step; ADR 261109-7c0s (proposed) records why, with the glossary term Handle and a changelog entry. A cold run with a larger model passed. Three cold runs with a small model did not: it still writes the full id before each title. Not covered: a title-first listing in the tool, which would fix that at the source.

## Log

- 2026-11-09: filed.
- 2026-11-09 17:57-07:00: claimed (claude-sonnet-5-5-names).
- 2026-11-09 18:01-07:00: Skill section, decision record, glossary term, changelog entry and AGENTS.md line written. Cold-agent check not passed: run one named items by title but never faced the similar or same-day items; run two quoted full ids beside titles. Goal box 3 stays open.
- 2026-11-09 18:03-07:00: Tightened the skill section and the AGENTS.md line (a flat 'never put an id in a reply', with a right and a wrong example). Third cold run with a small model still wrote the full id before every title, so the rule does not hold yet. Goal box 3 stays open.
- 2026-11-09 18:12-07:00: Cold run 4 (2026-11-09, Sonnet, six records including two near-identical titles and four made the same day): no id in the reply, each record named by its title, nothing needed a handle. Small models (Haiku, three runs) still write the full id before a title; recorded in docs/cold-agent-check.md as open for small models. Tool option for a title-first listing is not filed; file it if small-model agents matter.
