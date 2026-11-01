---
id: "261101-1smq"
title: close repairs links to moved items
summary: "Closing an item moves it into done/, which breaks relative links to it from other records; close should rewrite them instead of leaving a hand edit"
status: raw
kind: feature
parent: "261031-gtw6"
created: 2026-11-01
updated: 2026-11-01
---

# close repairs links to moved items

## Goal

Find out whether `check` reports the broken links and how often they happen (the history of `docs/work` shows five hand repairs, for example in 261031-qbat and 261031-2j14), then make `close` rewrite relative links that point at the moved file in every record, in the same write, listing the files it changed. Closed records are included: this is a mechanical repair, not a correction under ADR 261101-dxcn.

- [ ] A count of broken links after a `close` and whether `check` finds them.
- [ ] `close` rewrites the links and names the files; wikilinks need nothing.
- [ ] Spec, fixtures, glossary and changelog.

## Log

- 2026-11-01: filed by the spike 261031-z0n8.
