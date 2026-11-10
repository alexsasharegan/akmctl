---
id: "261109-7c0s"
title: Agents name records by title in chat, and give an id only as a last step
summary: The akm skill tells agents to name records by title in chat and to give an id only when nothing shorter tells the items apart.
status: proposed
date: 2026-11-09
---

# Agents name records by title in chat, and give an id only as a last step

## Context

An id such as `261109-n05b` is a stable reference: it survives a rename, sorts by date and never collides. It is a poor thing to put in a sentence. A reader who takes in text by dictation and text-to-speech hears it digit by digit, and a reply that lists a dozen items becomes a stream of numbers the reader must ignore or decode before deciding anything. A reader who does not use speech still has to match each id to an item to follow the reply. The maintainer who reported this cannot decide from such replies (work item 261109-n05b).

The skill says nothing about it today, so agents quote ids as the commands print them. Two places could carry the fix: the maintainer's own agent file, which helps one person and no one else, or the `akm` skill, which every user of the tool gets.

## Considered Options

1. **The `akm` skill names records by title, with ids as a last step** (chosen).
2. Only the maintainer's global agent file says it. No work for the project, but every other user keeps the problem, and so does the maintainer in a repository whose agent never read that file.
3. A tool option that prints the shortest unique tail of each id in a listing, with the skill left as it is. It fixes the source of the readback, but it is code and a spec change, and it leaves the agent free to keep quoting full ids from `show` and `board`.
4. Drop ids from chat altogether, always. The simplest rule, and it breaks down when two items share a title or the reader has to type the id.

## Decision

The `akm` skill gets a section, "Naming records in chat". In what an agent says to the maintainer it names a record by its title, and moves down this list only as far as it must to tell the items in its reply and the last few messages apart: the title alone; a word that differs; the date in the id, said as a date; for items made the same day, only the end of the id and only as many characters as separate them; the whole id when the reader asks for it or must type or paste it.

Ambiguity is judged against the items shown, not the whole repository. Before it sends, the agent checks that no two items shorten to the same phrase and steps down a level for those two only. The short end of an id is a handle for conversation: it is never written to a file or a commit. Ids stay in files, commands, links, commit trailers and pull request text; the rule concerns chat. A repository's `AGENTS.md` may ask for ids in chat and then wins, as it does over every skill.

Both copies of the skill (`.agents/skills/akm` and `internal/initrepo/templates/skills/akm`) carry the text, so `init` writes it and `akmctl skills` reports an older copy as outdated.

## Consequences

- Replies read the same aloud and on screen: a title says what the item is, which an id never did.
- Every user of the skill gets the default, including people who like ids in chat. They can restore ids with one line in `AGENTS.md`; that is the cost of making it a default instead of an option.
- An agent has to judge ambiguity. A model can get it wrong: two items can shorten to the same phrase. The cold-agent check has a case for it (`docs/cold-agent-check.md`). If agents keep getting it wrong, option 3 can be added later without changing this rule.
- Titles matter more. A title that begins with the same words as its neighbors forces the agent down the list. The skill asks for the distinguishing words first when an item is filed or retitled; nothing checks it.
- The repositories that already have the skill see it as outdated until they run `akmctl skills`. The text is guidance, so an old copy changes nothing in how the tool works.
- No new operation, field or rule. The tool is unchanged.
