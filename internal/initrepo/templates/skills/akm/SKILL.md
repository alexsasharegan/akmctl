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
| Regenerate the ADR index                 | `akmctl index`             | `references/index.md`   |
| Create a work item or ADR draft          | `akmctl new <kind> <slug>` | `references/new.md`     |
| Change status, priority or other fields  | `akmctl set <id> k=v`      | `references/set.md`     |
| Write Goal, Next, Outcome or a Log line  | `akmctl section <id> <n>`  | `references/section.md` |
| Take an item before working it           | `akmctl claim <id> --as n` | `references/claim.md`   |
| Give up an item                          | `akmctl release <id>`      | `references/release.md` |
| Finish an item                           | `akmctl close <id> done`   | `references/close.md`   |
| Check or upgrade a repository's layout   | `akmctl migrate`           | `references/migrate.md` |
| Start a new AKM repository               | `akmctl init`              | `references/init.md`    |
| Read or change repository settings       | `akmctl config`            | `references/config.md`  |

## Habits

- Run `akmctl check` before a commit; fix errors, read warnings.
- A refusal (`refused`, `stale`, `ambiguous`) is information: read `details.reason`, do not retry blindly.
- `--format json` gives the full response when the text form is not enough.
