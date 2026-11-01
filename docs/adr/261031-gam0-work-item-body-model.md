---
id: "261031-gam0"
title: A work item body is a list of sections, with typed core sections and open extension
summary: A work item body is an ordered list of named sections; Goal, Next, Log and Outcome are typed, any other ## section is kept as opaque Markdown
status: accepted
date: 2026-10-31
---

# A work item body is a list of sections, with typed core sections and open extension

## Context

Spike 261031-qbat surveyed 123 work items (53 open, 70 closed). Every one has `## Goal` and `## Log`; 120 keep the order Goal, Next, Outcome, Log; 112 have no other `##` heading; 11 (mostly early, long items) have extra sections such as `Options`, `Plan` or `To settle`. 448 of 449 log lines are `- YYYY-MM-DD: text`. 106 Goals hold a checklist (359 items, none nested, none with an id or commit marker), and 102 of those end the Goal with a "Done when" line and the list. 24 closed items still have unticked boxes. The workbench drawer and agents both need to edit the body, and the writer must be the tool (spike 261031-6bzq).

## Considered Options

- **(a) Closed.** Only the four core sections; anything else is rejected. 11 of 123 items fail; the fix is to rewrite real content.
- **(b) Declared extension.** Core sections plus extension sections that a repository declares by name, kind and node type; undeclared headings are errors. Strict, and costs a declaration mechanism before any value.
- **(c) Open with a typed core.** Any `##` heading is allowed. The four core sections have node types; other sections are opaque Markdown, addressed by name. Declaration of typed extensions can be added later without changing files.

## Decision

Option (c). A body is a preamble (text before the first `##`) and an ordered list of sections, each `{name, text}`; the name is the heading compared case-insensitively. Splitting and joining are lossless (123 of 123 items round-trip byte for byte), so unknown content is never rewritten.

The core sections are typed:

| Section | Node type                                                                                     |
| ------- | --------------------------------------------------------------------------------------------- |
| Goal    | prose, plus one checklist: the first contiguous run of `- [ ]` / `- [x]` lines, one line each |
| Next    | prose                                                                                         |
| Log     | append-only entries, `- YYYY-MM-DD: text`                                                     |
| Outcome | prose (checklist-aligned evidence is a later refinement)                                      |

Typed nodes are addressed by section name and 1-based position, guarded by a section etag, not by ids stored in the file (no existing item carries one, and ids would add noise to hand-edited Markdown). A write splices only the lines it changes (ticking an item changes `[ ]` to `[x]` in place) and normalizes nothing else. Where a typed section does not parse (4 Goals with prose after the checklist, 1 Log line with a continuation), the writer treats the rest as opaque text and still allows edits to the typed part.

## Consequences

- Nothing existing has to be rewritten; 123 of 123 items are valid.
- The drawer and agents can address a checklist item or a log line without matching text.
- Typed extension sections remain possible later (261031-x50y) and need no file change.
- Closing does not tick the "done when" boxes (24 closed items are unticked); this model makes `close` ticking them possible (261031-zek2).
- A hand edit that breaks a typed shape degrades to opaque text rather than failing; `check` rules for it are for 261031-x50y.
- The parser behind the section tree is decided separately (proposed ADR 261031-cpbg, from spike 261031-g9zk); this decision holds either way, since it is about the model and the edit rule.
