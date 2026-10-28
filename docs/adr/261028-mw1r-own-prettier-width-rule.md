---
id: 261028-mw1r
title: Implement Prettier's string-width rule in the tool instead of using a width library
summary: Table padding in the ADR index follows Prettier's getStringWidth, implemented here from Unicode 17.0.0 tables. Read before changing the index or the tables.
status: accepted
date: 2026-10-28
scope: [internal/width, internal/index]
authored_by: agent
---

# Implement Prettier's string-width rule in the tool instead of using a width library

## Context

The generated ADR index must be a table that Prettier leaves unchanged, so the generator pads cells to Prettier's display width. A width library measures grapheme clusters; Prettier counts code points with its own rules for wide characters and emoji. An evaluation of libraries found none that matches Prettier. (Decision log entries Y16, Y23, Y28.)

## Considered Options

- Pad by rune count: no code, wrong for wide characters and emoji.
- A width library (`displaywidth` or similar): one dependency, and it differs from Prettier.
- Implement Prettier 3.9.10's rule with tables generated from Unicode data by a committed generator.

## Decision

Implement the rule ourselves in `internal/width`. East Asian Wide and Fullwidth ranges and the emoji tables are generated from Unicode 17.0.0 by `go generate ./internal/width` (`gen/main.go`, data in `data/`). A 206-cell comparison with Prettier's output is the test. This replaced an earlier choice of a small dependency.

## Consequences

- The standard library only is restored (ADR 261027-ebht).
- Known gap: the code differs from Prettier on non-standard emoji ZWJ pairs (a work item tracks it).
- A Prettier bump is re-tested against the fixtures.

## Confirmation

`internal/width` tests (the Prettier comparison) and the index test.
