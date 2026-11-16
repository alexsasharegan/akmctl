---
id: "261101-1smq"
title: close repairs links to moved items
summary: "Closing an item moves it into done/, which breaks relative links to it from other records; close should rewrite them instead of leaving a hand edit"
status: dropped
kind: feature
parent: "261031-gtw6"
created: 2026-11-01
updated: 2026-11-16
closed: 2026-11-16
---

# close repairs links to moved items

## Goal

Find out whether `check` reports the broken links and how often they happen (the history of `docs/work` shows five hand repairs, for example in 261031-qbat and 261031-2j14), then make `close` rewrite relative links that point at the moved file in every record, in the same write, listing the files it changed. Closed records are included: this is a mechanical repair, not a correction under ADR 261101-dxcn.

- [ ] A count of broken links after a `close` and whether `check` finds them.
- [ ] `close` rewrites the links and names the files; wikilinks need nothing.
- [ ] Spec, fixtures, glossary and changelog.

## Outcome

Dropped as a duplicate: `close` has rewritten the relative links to a moved item, in open and closed records, since the item on close rewriting Markdown links ([261031-wjs3](261031-wjs3-close-rewrites-markdown-links.md)), which was built the day before this one was filed.

- Count: nine closes in the branch stack of 2026-11-16 left no broken link; `check` (rule L1) was clean after each. One of them rewrote links in two closed records and listed the files.
- `close` rewrites and names the files: `TestCloseRewritesLinks`, fixture `close-8`, `links_rewritten` in the response.
- Spec, fixtures, glossary and changelog: done by that item.

Decision: none.

## Log

- 2026-11-01: filed by the spike 261031-z0n8.
- 2026-11-16 07:11-07:00: checked on 2026-11-16: close already does this since 261031-wjs3 (built the day before this item was filed; the five hand repairs the spike counted are older). Closing 261101-3che in this branch stack rewrote the links to it in two closed records (261028-esyt and 261031-z0n8, both in done/) and listed them; check reported no L1 after any of the nine closes in the stack
