---
id: "261101-tr8e"
title: Body writes go through the section tree, with tick and log operations
summary: Body writes are splices by the section tree; add tick and log operations, address items by section and position, and refuse what the scanner cannot read
status: proposed
date: 2026-11-01
scope: [internal/api, internal/body, internal/record, internal/server]
authored_by: agent
---

# Body writes go through the section tree, with tick and log operations

## Context

Spike 261031-6bzq. The model is decided (ADR 261031-gam0: an open body with typed Goal, Next, Log and Outcome) and the package exists (`internal/body`: `Split`, `Join`, the Goal checklist, Log entries, splice edits; ADRs 261031-cpbg and 261031-nvka).

Today one operation writes a body: `section`, which replaces a `## ` section or appends to it with the line scan `record.SetSection`. To tick a Goal box an agent must resend the whole Goal (median 837 bytes, up to 4.8 KB) or hand-edit the file after reading it (median file 1.9 KB, up to 20 KB), and either can get the Markdown wrong. A Log line needs the right date and bullet. The drawer will need to tick, append and replace without losing a concurrent edit. `SetSection` and the tree must not drift: on all 129 work items, replacing each section and appending a Log line with the tree gives byte-for-byte the same file as `SetSection` (401 of 401 edits, `TestTreeMatchesSetSection`).

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
- **Fail closed.** The write operations use `Split`. A body with a construct the scanner refuses (`ErrUnsupported`) is `refused`, with the construct and line in `details`, and nothing is written; the maintainer edits that file by hand. All 160 records pass today.
- **Missing sections.** Creating a missing section (before `## Log`, as `SetSection` does) becomes a tree operation (`AddSection`) before `section` moves onto the tree.
- **`section` moves onto the tree**, and `SetSection` is deleted, once `AddSection` exists and the equivalence test passes for the new behaviour.

## Consequences

- An agent ticks a box with about 35 bytes of request after reading the Goal, instead of resending it or re-reading the file; a stale or wrong address is an error and not a silent edit.
- The drawer can tick, append and replace with section-level conflict handling and no Markdown of its own.
- The API grows by two operations and one optional parameter, additive for clients that ignore them (API 1.x). Two fixtures per operation in both copies, the spec, the skills and the HTTP table change.
- `section` becomes stricter: a body with an unsupported construct can no longer be written by the tool. This is the price of one scanner and one fence rule; the line scan stops being a second opinion on where sections are.
- Position addressing can pick the wrong item if the Goal changed between a read and a write. `section_etag` and `if_match` exist for that; an agent that writes without either gets last-writer-wins on positions, as it does for text today.
- Duplicate section names are addressed by the first. `Next`'s optional `Saved:` line is part of the section text and is replaced with it.
