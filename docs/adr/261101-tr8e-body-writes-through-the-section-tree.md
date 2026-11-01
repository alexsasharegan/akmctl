---
id: "261101-tr8e"
title: Body writes go through the section tree, with tick and log operations
summary: Body writes are splices by the section tree; add tick and log operations, address items by section and position, refuse only unreadable constructs that matter
status: accepted
date: 2026-11-01
scope: [internal/api, internal/body, internal/record, internal/server]
---

# Body writes go through the section tree, with tick and log operations

## Context

Spike 261031-6bzq. The model is decided (ADR 261031-gam0: an open body with typed Goal, Next, Log and Outcome) and the package exists (`internal/body`: `Split`, `Join`, the Goal checklist, Log entries, splice edits; ADRs 261031-cpbg and 261031-nvka).

Today one operation writes a body: `section`, which replaces a `## ` section or appends to it with the line scan `record.SetSection`. To tick a Goal box an agent must resend the whole Goal (median 837 bytes, up to 4.8 KB) or hand-edit the file after reading it (median file 1.9 KB, up to 20 KB), and either can get the Markdown wrong. A Log line needs the right date and bullet. The drawer will need to tick, append and replace without losing a concurrent edit. `SetSection` and the tree must not drift: on all 129 work items, replacing each section and appending a Log line with the tree gives byte-for-byte the same file as `SetSection` (401 of 401 edits, `TestTreeMatchesSetSection`).

**A census of real notes** (2026-11-01, the maintainer's Obsidian vault; counts only, nothing from the notes is kept). Whole-file refusals under the scanner as it is: this repository's 160 records, 0 refused. The vault's work items (137 files, 568 `##` sections): 18 files refused, all local (16 blockquotes, 2 ordered-list tasks), which is 13% of files but 18 of 568 sections (3%). Its ADRs: 1 of 26 refused (an HTML line, structural). Every note in the vault (2,544 files, outside the assets folder): 31% refused, led by a thematic break or setext underline straight after text (272 files), blockquotes (227), tab-indented lines (132) and empty list items (46). The last group is not what the tool edits, but it shows how a stricter scanner would fare on prose notes.

## Considered Options

- **(a) Keep `section` only.** No new surface. Ticking stays a full Goal rewrite, and nothing protects a section from a stale write except the whole-file etag.
- **(b) Typed operations on the tree.** Keep `section` and add `tick` and `log`, all built on `body.Tree`; address checklist items by section and 1-based position; a section etag for fine-grained conflicts; `section` moves onto the tree and refuses what the scanner refuses.
- **(c) A general edit language** (a patch or JSON-pointer style API over the tree). Flexible, and a second language for agents to get wrong; nothing in the 129 items needs it.
- **(d) Ids stored in the file** for checklist items and log lines, so edits address them without a position. Rejected by ADR 261031-gam0: no item carries one, and ids add noise to hand-edited Markdown.

## Decision

Option (b).

- **Operations.** `tick` sets or clears one checklist box: `id`, `item` (1-based position in the Goal's checklist; `section` defaults to `goal`), `done` (default true), plus the usual `if_match`, `section_etag`, `dry_run` and `as_of`. `log` appends one dated entry: `id`, `text` (one line), with the date taken from `as_of` or today's local day. `section` stays for prose and for multi-line text. CLI: `akmctl tick <id> <n> [--undo]`, `akmctl log <id> <text>`. HTTP: `PUT /records/{id}/sections/{name}/items/{n}` and `POST /records/{id}/log`.
- **One write.** Each operation reads the file under the repository lock, splits the body, splices, sets `updated` and writes the file once by the existing compare-and-swap, so a body edit and `updated` are atomic. `dry_run` and its `diff` come with it.
- **Addressing and conflicts.** Items are addressed by section name and position, as in ADR 261031-gam0. `if_match` is the whole-file etag as for `set`. `section_etag` is optional: a client that holds one section for a long time (the drawer) sends the etag it read, and the operation applies if that section is unchanged even when other parts of the file moved on; otherwise it answers `stale` and returns the current section etag. Responses carry the section and its etag.
- **Round-trip rule.** Bytes outside the edited span are never changed. Inside it: a tick changes one byte; a log entry is one line `- YYYY-MM-DD: text` added after the last non-blank line of `## Log`, with the line ending of the line before it; a replace keeps the blank lines before and after the old text and writes the new text with the file's line ending, trimmed of leading and trailing blank lines.
- **Fail closed, narrowly.** The scanner's refusals (`ErrUnsupported`) come in two kinds, and `Locate` reports which. A _structural_ refusal is a construct that could change where a heading, fence or checklist is: an HTML line, an indented fence, heading or code line, a tab-indented line, a setext underline or thematic break after text, a block inside a list item, a control character or a bare carriage return. A _local_ refusal is confined to its own lines: a blockquote, a table-delimiter lookalike, an ordered-list task, a task box with unusual contents, an empty list item. A write is `refused`, with the construct and line in `details`, and nothing is written, when the body has a structural refusal anywhere, or a local one inside the section being edited. A local construct in another section is not a reason to refuse: a Goal box ticks in a file whose Notes section has a callout. All 160 records pass today.
- **Fewer refusals over time.** Local constructs the scanner can skip as opaque blocks (a blockquote region, an indented code block) stop being refusals one at a time, each checked against the goldmark oracle (work item 261101-4j86). Refusing is the default; claiming a construct needs the battery to agree.
- **Missing sections.** Creating a missing section (before `## Log`, as `SetSection` does) becomes a tree operation (`AddSection`) before `section` moves onto the tree.
- **`section` moves onto the tree**, and `SetSection` is deleted, once `AddSection` exists and the equivalence test passes for the new behaviour.

## Principles

What the questions around this decision settled; they bind the code that implements it.

1. **The document is the source of truth and the tool is a guest.** The tool never requires a file to conform to it, never changes the format to suit its parser (no ids or markers in the Markdown), and never writes a byte it was not asked to change.
2. **Understand as little Markdown as possible.** The tree reads headings, fences, task lines and list lines; everything else is opaque bytes carried through unchanged. A construct is claimed only when an edit needs it, and then only with the oracle's agreement.
3. **A refusal is cautious, narrow and named.** When the tool cannot be sure, it declines rather than guesses. It declines only where the doubt can affect the edit (a structural construct anywhere, a local one in the edited section), says which construct and line, and loses nothing.
4. **A refusal is a defect to burn down, and it is measured.** The refusal rate on real notes is a tracked number, not an assumption. Every kind of refusal met on real work items or ADRs gets its own backlog item, however rare: one exception is still a note the tool cannot edit. Separately, when the rate passes a threshold the maintainer sets, the answer is to claim more constructs or to promote goldmark, not to tell people to change their notes. The threshold, set by the maintainer on 2026-11-01: more than 5% of the work items in a real repository have a refusal in a section an operation edits (today, in the vault: 3% of sections, 13% of files whole).
5. **The parser is swappable, so the choice is cheap to reverse.** Nothing written into a file depends on how it was read. The `Locator` interface returns blocks with byte offsets; our scanner, a larger scanner or goldmark (ADR 261031-cpbg keeps that option, which needs an ADR for the dependency) can sit behind it without a file migration.
6. **Reach for another encoding only when the humans stop reading the files.** A structured encoding would make writing trivial and give up the reason for Markdown here: records that a person and an agent both read and edit in an editor and in git.

## Confirmation

`TestTreeMatchesSetSection` (`internal/record/section_test.go`) holds the section tree and `SetSection` to one behaviour on every work item until `section` moves onto the tree; `TestSplitJoinRecords`, `TestTreeAgainstOracle` and `FuzzTree` (`internal/body`) hold the tree to lossless edits. The operations add their own tests (work items 261101-ga4s and 261101-q72k).

## Consequences

- An agent ticks a box with about 35 bytes of request after reading the Goal, instead of resending it or re-reading the file; a stale or wrong address is an error and not a silent edit.
- The drawer can tick, append and replace with section-level conflict handling and no Markdown of its own.
- The API grows by two operations and one optional parameter, additive for clients that ignore them (API 1.x). Two fixtures per operation in both copies, the spec, the skills and the HTTP table change.
- `section` becomes stricter only where it must: a body with a structural construct, or a local one in the edited section, can no longer be written by the tool. This is the price of one scanner and one fence rule; the line scan stops being a second opinion on where sections are. `check` should report such a body before an agent meets it (a rule for 261031-x50y).
- Position addressing can pick the wrong item if the Goal changed between a read and a write. `section_etag` and `if_match` exist for that; an agent that writes without either gets last-writer-wins on positions, as it does for text today.
- Duplicate section names are addressed by the first. `Next`'s optional `Saved:` line is part of the section text and is replaced with it.
