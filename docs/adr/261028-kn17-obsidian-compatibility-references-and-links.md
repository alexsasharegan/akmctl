---
id: 261028-kn17
title: "Keep records valid in Obsidian: bare ids by default, quoted wikilinks accepted, Markdown links checked"
summary: Obsidian wins any conflict. References are bare ids or quoted wikilinks; body links are Markdown links (L1) or wikilinks (L2).
status: accepted
date: 2026-10-28
scope: [internal/record, internal/check]
authored_by: agent
---

# Keep records valid in Obsidian: bare ids by default, quoted wikilinks accepted, Markdown links checked

## Context

Records are edited in Obsidian by the maintainer. Obsidian rewrites frontmatter, resolves wikilinks by whole filename and types properties by name for a whole vault. A test against Obsidian 1.14.5 showed what it keeps and what it changes. (Decision log entries Y26, Y29, Y30, Y31; the standard's ADR 020 records the compatibility goal.)

## Considered Options

- Bare ids only: no links or backlinks in Obsidian.
- Wikilinks only: needs the slug in every reference and breaks plain-text use and grep.
- Bare id canonical, a quoted wikilink `"[[<filename>|<id>]]"` accepted everywhere a reference is allowed.

## Decision

The tool reads both reference forms, resolves by the leading id, writes bare ids unless the manifest sets `reference_style: wikilink`, and never rewrites a reference it did not create. Relative Markdown links are checked as errors (L1); wikilinks are resolved by filename as warnings (L2). Dot-files are untouched by Obsidian and hold no records.

## Consequences

- A renamed slug is not an error.
- A vault that types `id` as a number differs from AKM's string `id`; that vault resolves it on its side.

## Confirmation

Rules L1 and L2; the reference tests in `internal/record`.
