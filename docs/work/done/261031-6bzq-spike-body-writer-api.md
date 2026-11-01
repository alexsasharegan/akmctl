---
id: "261031-6bzq"
title: spike body writer api
summary: "Spike: a Markdown writer and body write operations so no agent or UI hand-writes Markdown"
status: done
kind: spike
parent: "261031-gtw6"
created: 2026-10-31
updated: 2026-11-01
closed: 2026-11-01
decisions: [261101-tr8e]
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

- [x] A prototype writer round-trips every item in `docs/work/` and `docs/work/done/` with a stated rule for any diff it produces.
- [x] Three operations are run end to end on a scratch repository, with the commands and the resulting diffs.
- [x] A recommendation (operations, addressing, CLI and API shapes, token comparison) is in `## Outcome`; a hard-to-reverse choice has a `proposed` ADR; follow-up items are filed.

## Next

Done and verified; close with the Outcome.

## Outcome

Recommendation: build typed body operations on the section tree (option (b) of the proposed ADR [261101-tr8e](../../adr/261101-tr8e-body-writes-through-the-section-tree.md)). Follow-ups: [261101-ga4s](../261101-ga4s-tick-and-log-operations.md) and [261101-q72k](../261101-q72k-section-on-the-tree.md).

- **Writer and round-trip rule** (box 1). The writer is `body.Tree` (261031-j4bc), standard library only, built by targeted splice and never by re-render. Outside the edited span nothing changes (tested on 9,207 generated documents and a 120-second fuzz run); inside it a tick changes one byte, a log entry is one new line, and a replace keeps the blank lines around the text and writes the file's line ending. Against the current writer, `TestTreeMatchesSetSection` (kept, in `internal/record/section_test.go`) replaces every section and appends a Log line on all 129 work items with both: 401 of 401 edits are byte-identical, so the tree adds no normalization that `SetSection` did not already do. The one thing it cannot do yet is create a missing section (`AddSection`, follow-up 261101-q72k). Full re-render was not tried: the splice already reproduces every file, so a render could only add diffs.
- **Three operations end to end** (box 2) on a copy of the `basic` fixture repository, through `Split`, a splice and the frontmatter editor (a throwaway test, deleted; the CLI and API operations do not exist yet). Tick Goal item 1, append a log entry, replace Next, with `updated` set in the same write; the resulting diff:

```diff
-updated: 2026-10-27
+updated: 2026-11-01
...
-- [ ] It is done.
+- [x] It is done.
...
 ## Next
-Write the encoder test.
+Wire the encoder into ls.
-
-Saved: 2026-10-27 18:00
-
 ## Log
 - 2026-10-20: started.
+- 2026-11-01: encoder test written
```

The replace dropped the optional `Saved:` line because it is part of Next's text; `section` does the same today. Nothing else moved.

- **Operations and addressing.** `tick` (section, 1-based item, done) and `log` (one line, tool-supplied date) are new; `section` stays for prose. Items are addressed by section and position (ADR 261031-gam0); `if_match` covers the file and an optional `section_etag` lets the drawer apply when only its section is unchanged. Errors: `not_found` (id, section, item), `refused` `closed`, `refused` with the construct and line for `ErrUnsupported`, `stale` with the current etag. `dry_run` and `diff` come with the existing edit machinery; the body and `updated` are one atomic file write.
- **Token cost** (bytes; about 4 per token, an estimate; sizes measured on the 129 work items: file median 1,859, p90 5,581, max 19,777; Goal median 837, p90 2,209, max 4,784; Next median 143). Ticking one box: hand edit, read the file and send an edit with enough context to be unique, about 1.9 KB median and 5.7 KB at p90; `section` replace, read the Goal and send it back, about 1.7 KB median and 4.4 KB at p90; `tick`, read the Goal then about 35 bytes plus a response of a few hundred bytes (the text output of a write is about 0.5 KB), about 1.3 KB median. So about a third less than either today, and the saving grows with the Goal. A log line already costs about 100 bytes through `section --mode append`, so `log` saves little (the date and the bullet); its value is that the tool owns the date and the shape. Failure cases: a hand edit that matches two lines or whitespace fails and re-reads; a wrong position or stale etag is one refusal and one re-read of the section.
- **The drawer** uses the same operations over HTTP, with `If-Match` as for every write and `section_etag` from the section it displayed. No new conflict machinery is needed beyond the optional etag; the UI is for 261031-18t8 and 261031-z0n8.
- **Fail closed has a cost.** Moving `section` onto the tree makes a body with an unsupported construct unwritable by the tool. All 160 records pass, but a hand-written note with a blockquote in Goal would be refused; the ADR accepts that.
- **Not covered:** the CLI and API shapes were not prototyped in Go (they are the follow-up); the Outcome checklist-evidence refinement of ADR 261031-gam0; `close` ticking the Goal boxes (261031-zek2 can now use `SetChecked`); duplicate section names (edits address the first).

Decision: [261101-tr8e](../../adr/261101-tr8e-body-writes-through-the-section-tree.md) (proposed).

## Log

- 2026-10-31: filed.
- 2026-11-01: claimed (claude).
- 2026-11-01: compared the tree with SetSection on all work items (401 of 401 identical), ran three ops on a scratch copy, drafted ADR 261101-tr8e, filed 261101-ga4s and 261101-q72k.
