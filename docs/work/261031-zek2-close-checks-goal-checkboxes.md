---
id: "261031-zek2"
title: close checks goal checkboxes
summary: "close refuses while a Goal checkbox is unchecked; a command ticks one box with evidence"
status: raw
kind: feature
parent: "261031-v07x"
created: 2026-10-31
updated: 2026-10-31
---

# close checks goal checkboxes

## Goal

From the consultation (gap 1). Rule WK7 checks that `## Goal` and `## Outcome` exist, but `close done` never reads the `- [ ]` lines, so a one-line Outcome closes an item whose criteria were not met. Done when:

- [ ] `close <id> done` is `refused` while `## Goal` has an unchecked `- [ ]` line, and the error names each such line. `dropped` is not affected.
- [ ] A way to close anyway is decided (a flag, or tick the box and say why) and written down; a hard-to-reverse choice gets a `proposed` ADR.
- [ ] A command ticks one Goal line by number and appends the evidence to `## Log` (shape to decide, for example `check <id> <n> --evidence "..."`); an out-of-range number is `not_found`.
- [ ] The API spec, the `close` reference page and the glossary cover the new behavior.

## Log

- 2026-10-31: filed from the consultation. The `section` operation already writes a whole section, so the new part is the per-line tick, not general checklist editing.
