---
id: "261101-gw7b"
title: Log entries carry a local time with its offset
summary: A Log line carries the writer's local time and its offset; claim, release and log write it, and an optional at records when the thing happened
status: proposed
date: 2026-11-01
scope: [internal/api, internal/body]
---

# Log entries carry a local time with its offset

## Context

A Log line is `- YYYY-MM-DD: text`, dated with the writer's local day (ADR 261029-5pe1). The tool stamps the date, and that stamp is much of why the tool writes the line and not the agent. A day is too coarse: an item is claimed, ticked, blocked and closed within one afternoon, and the Log cannot say in what order or how far apart. The maintainer asked for the time (2026-11-01).

What reads a Log line today: nothing in production. `Tree.LogEntries` is the only parser and only tests call it; `check` does not read the Log (the spec says it is append-only by convention and is not diffed); `show`, `search` and the workbench treat it as text. The 506 existing lines are date-only and must stay valid. The AKM spec template shows `- YYYY-MM-DD: dated lines`, and `claim` and `release` write their own Log lines.

ADR 261029-5pe1 gives the rule: UTC where machines compare values, the writer's local time where a person reads them. A Log line is read by people and is never compared, so it is local. That ADR rejected an offset on timestamps, for three reasons: Obsidian and YAML readers show offsets poorly, the format is longer, and the only comparison it needed was "how long ago". The first applies to frontmatter values and not to a Markdown body line, and the third does not apply to a Log. The second, the length, is the real cost.

`as_of` today takes `YYYY-MM-DD` or `YYYY-MM-DDTHH:MM` and reads the time as UTC (`time.Parse`); `claim` writes `claimed_at` from it. Fixtures depend on that to be deterministic.

## Considered Options

- **(a) Date only, as now.** No cost, no order within a day, and the tool's stamp carries little.
- **(b) Local time with no zone** (`2026-11-01 14:30`). Short and readable. Ambiguous when a file is read in another zone or after a daylight-saving change: the same trap 5pe1 fixed for dates, one level down.
- **(c) Local time with its offset** (`2026-11-01 14:30-07:00`). Local, so a person reads their own clock, and unambiguous. Twelve characters longer than (a). Chosen.
- **(d) UTC with a `Z`** (`2026-11-01T21:30Z`). Unambiguous and short, and not the writer's clock: the Phoenix evening trap again, in the Log.
- **Separator.** `T` (as `claimed_at`) or a space. A space reads better in prose and matches the `Saved: YYYY-MM-DD HH:MM` line in `## Next`. Space chosen.
- **Backdating.** `--as-of` sets the stamp, so a line can claim to be written when it was not. Or a separate `at` that records when the thing happened, in the text, and leaves the stamp as the write time. The second chosen.

## Decision

Option (c), with a space.

- **Form.** `- YYYY-MM-DD HH:MM+HH:MM: text` (or `-HH:MM`; `+00:00` for UTC). The time is the writer's local clock to the minute; the offset is the machine's offset at that moment, so a line across a daylight-saving change is right. A date-only line stays valid forever: a hand-written line, or one of the 506 existing ones, is not rewritten.
- **Who writes it.** `log`, `claim` and `release` (their Log lines are `claimed (who)`, `took over from ...`, `released (who)`). No other operation writes a Log line.
- **Reading.** `LogEntries` accepts both forms: `^- (\d{4}-\d{2}-\d{2})(?: (\d{2}:\d{2})([+-]\d{2}:\d{2}))?: ?(.*)$`. `Entry` gains `Time` and `Offset`, empty for a date-only line. The Log is in write order. The tool does not sort it by time, because lines written in different zones would reorder.
- **`as_of`.** A date-only `as_of` writes a date-only line (the tool does not invent a time of day). An `as_of` with a time writes that time with offset `+00:00`, because `as_of` is read as UTC as it is today; this keeps the fixtures deterministic. A real write uses the machine's zone.
- **`at`: when the thing happened.** `log` takes an optional `at` (`--at`): a timestamp, `YYYY-MM-DD HH:MM` or `YYYY-MM-DDTHH:MM`, with an optional offset (local zone when absent). The line is stamped with the write time as usual and ends with ` (at <timestamp>)`, the date left out when it is the day of the stamp: `- 2026-11-01 16:02-07:00: deployed (at 11:05-07:00)`. An `at` after the write time is `invalid_request`. `as_of` keeps its job of fixing the clock for tests and replays.
- **Three conventions, one table**, which the spec, the glossary and the `akm` skills state together so no reader takes one for a clock error (work item 261031-mr2q):

| Value                                                | Zone                             | Because                     |
| ---------------------------------------------------- | -------------------------------- | --------------------------- |
| `created`, `updated`, `closed`, an ADR's `date`, ids | the writer's local day           | a person reads it           |
| a Log line stamp                                     | the writer's local time + offset | a person reads it, any zone |
| `claimed_at`                                         | UTC, `YYYY-MM-DDTHH:MM`          | machines compare it         |

## Consequences

- The Log says in what order and how far apart things happened. The tool's stamp is worth more, and the agent has no reason to type a time.
- A line is 12 characters longer. `claim` and `release` lines get longer too.
- One more line shape to read (`LogEntries`); the 506 old lines need nothing. A reader of the standard that matches `^- \d{4}-\d{2}-\d{2}: ` strictly will not match the new lines. Nothing in this repository does; the template and spec say the time is optional so a strict reader has notice.
- The spec template, the `akm` and `akm-track-work` skills, the glossary, the fixtures that show a Log line (`claim-1` and others, in both copies) and the changelog change. An API 1.x addition (the `at` parameter) and a changed output of `claim`, `release` and `log`.
- Writes from machines in different zones interleave in write order, each line in its own writer's clock. That is the purpose of the offset: the line is correct wherever it is read.
- Option (d) and a `Z` for `claimed_at` (work item 261031-mr2q, option 2) stay separate decisions. This ADR does not change `claimed_at`.
- Work item 261101-vy77 implements it once this is accepted.
