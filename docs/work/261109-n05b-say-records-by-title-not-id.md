---
id: "261109-n05b"
title: Agents name records by title, not id, and resolve ambiguity in steps
status: active
created: 2026-11-09
updated: 2026-11-09
claimed_by: claude-sonnet-5-5-names
claimed_at: 2026-11-10T00:57
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

- [ ] Both copies of the `akm` skill state the rule (`.agents/skills/akm` and `internal/initrepo/templates/skills/akm`), and the repository's `AGENTS.md` says it may override.
- [ ] A decision record (`proposed`, for the maintainer to accept) explains why the default changes for every user of the skill, and the glossary registers any new term.
- [ ] A cold-agent check (`docs/cold-agent-check.md`) gives an agent a mixed list that contains two similar titles and two items made the same day, and the agent's reply names each one without a full id and without two items sharing a handle.
- [ ] A `CHANGELOG.md` entry under `## Unreleased` says what changed and who is affected.

Not in scope: a listing option that prints the shortest unique tail. File it only if the cold-agent check shows agents get the arithmetic wrong.

## Next

Write the rule into both copies of the akm skill and the AGENTS.md override note, draft the proposed decision record, add the changelog entry and the cold-agent check row.

## Log

- 2026-11-09: filed.
- 2026-11-09 17:57-07:00: claimed (claude-sonnet-5-5-names).
