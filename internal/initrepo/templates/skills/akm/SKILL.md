---
name: akm
description: Use the akmctl tool to read and change AKM records (work items, ADRs, roadmap). Use whenever a task touches the work or ADR folders or the roadmap, or the user mentions akmctl.
---

# akm

AKM keeps a repository's decisions, direction and work as Markdown records. `akmctl` reads, writes and checks them; it never stages, commits or pushes.

The repository's `AGENTS.md` wins over this skill. The procedures are in `akm-track-work` and `akm-record-decision`.

If `akmctl` is not installed (`akmctl version` fails), tell the maintainer and use the plain-file fallbacks in `AGENTS.md`. Do not invent ids or hand-edit the ADR index.

## Which command

Read the reference before the first use of a command.

| Task                                     | Command                    | Reference               |
| ---------------------------------------- | -------------------------- | ----------------------- |
| Tool and API version, extensions         | `akmctl version`           | `references/version.md` |
| List records, find open work             | `akmctl ls`                | `references/ls.md`      |
| Find records by words                    | `akmctl search <text>`     | `references/search.md`  |
| Read one record                          | `akmctl show <id>`         | `references/show.md`    |
| See active, blocked and ready work       | `akmctl board`             | `references/board.md`   |
| Get the file path of an id               | `akmctl find <id>`         | `references/find.md`    |
| Read the date inside an id               | `akmctl when <id>`         | `references/when.md`    |
| Validate the records                     | `akmctl check`             | `references/check.md`   |
| Check step commits against their diff    | `akmctl commits`           | `references/commits.md` |
| See what happened to the work, in order  | `akmctl history`           | `references/history.md` |
| Regenerate the ADR index                 | `akmctl index`             | `references/index.md`   |
| Write an Obsidian Bases file of the work | `akmctl bases`             | `references/bases.md`   |
| Create a work item or ADR draft          | `akmctl new <kind> <slug>` | `references/new.md`     |
| Change status, priority or other fields  | `akmctl set <id> k=v`      | `references/set.md`     |
| Write Goal, Next or Outcome              | `akmctl section <id> <n>`  | `references/section.md` |
| Tick a box of the Goal                   | `akmctl tick <id> <n>`     | `references/tick.md`    |
| Remove a Goal box that no longer applies | `akmctl remove <id> <n>`   | `references/remove.md`  |
| Add a stamped Log line                   | `akmctl log <id> <text>`   | `references/log.md`     |
| Take an item before working it           | `akmctl claim <id> --as n` | `references/claim.md`   |
| See a list of writes before them         | `akmctl preview`           | `references/preview.md` |
| Give up an item                          | `akmctl release <id>`      | `references/release.md` |
| Finish an item                           | `akmctl close <id> done`   | `references/close.md`   |
| Check or upgrade a repository's layout   | `akmctl migrate`           | `references/migrate.md` |
| Start a new AKM repository               | `akmctl init`              | `references/init.md`    |
| Read or change repository settings       | `akmctl config`            | `references/config.md`  |
| Check or refresh the AKM skills          | `akmctl skills`            | `references/skills.md`  |

## Naming records in chat

Ids are for files, commands, links and commit trailers. In what you say to the maintainer, name a record by its title: an id read aloud (by a person or a text-to-speech voice) costs more attention than the decision it belongs to. Stop at the first step that tells the items apart among those in your reply and the last few messages:

1. The title alone, shortened ("the compatibility policy item"), when nothing else in view could be mistaken for it.
2. A word that differs: name the part that separates two similar items ("the refactor that prepares for theming", "the theme additions").
3. The date in the id, said as a date ("the one from October twelfth"). `akmctl when <id>` reads it.
4. For items made the same day, only the end of the id and only as many characters as separate them ("the one ending in x j"). It is a handle for talking: never write it into a file or a commit.
5. The whole id, when the maintainer asks for it or must type or paste it.

Judge ambiguity against the items shown, not the whole repository: a list that is all about theming needs no more than "the refactor" and "the themes". Before you send, check that no two items shorten to the same phrase, and step down a level for those two only. When you file or retitle an item, put the words that tell it from its neighbors first, because speech reads left to right.

A repository whose `AGENTS.md` asks for ids in chat overrides this section.

## Habits

- Run `akmctl check` before a commit; fix errors, read warnings.
- If `akmctl skills` says the repository's skills are newer than your `akmctl`, read `references/skills.md` before you go on.
- A refusal (`refused`, `stale`, `ambiguous`) is information: read `details.reason`, do not retry blindly.
- `--format json` gives the full response when the text form is not enough.
