---
id: "261031-6bzq"
title: spike body writer api
summary: "Spike: a Markdown writer and body write operations so no agent or UI hand-writes Markdown"
status: raw
kind: spike
parent: "261031-gtw6"
created: 2026-10-31
updated: 2026-10-31
---

# spike body writer api

## Goal

A strong priority: every write to a record body goes through the tool. An agent or the workbench says "tick this item", "append this log line", "replace the goal text", and the tool produces correct Markdown, so no one spends tokens on formatting. Depends on the model in [261031-qbat](261031-qbat-spike-work-item-body-model.md). Find out:

- The writer: a parse to the section tree and a render back, standard library only (no new dependency without an ADR). What it must preserve byte for byte (unknown sections, hand formatting, wikilinks, tables) and what it may normalize. Compare a full re-render with targeted splices of the changed span; the diff a maintainer sees in git is the test.
- The operations: `section` already exists, so start from what it does and where it falls short. Candidates: set or append text in a section, add, tick or untick a checklist item, append a dated log entry, add or remove an extension section, move a section. Addressing by section and item id rather than by text.
- The surface: CLI commands and API operations, request and response shapes, errors (unknown section, stale etag, a section the status forbids), dry run and a diff in the response, atomicity with the frontmatter write that goes with it (`updated`).
- Token cost for an agent: the size of a typical request and response against hand-editing the file, including the failure and retry cases.
- How the drawer uses the same operations, and what etag or section-level conflict handling it needs.

Done when:

- [ ] A prototype writer round-trips every item in `docs/work/` and `docs/work/done/` with a stated rule for any diff it produces.
- [ ] Three operations are run end to end on a scratch repository, with the commands and the resulting diffs.
- [ ] A recommendation (operations, addressing, CLI and API shapes, token comparison) is in `## Outcome`; a hard-to-reverse choice has a `proposed` ADR; follow-up items are filed.

## Log

- 2026-10-31: filed.
