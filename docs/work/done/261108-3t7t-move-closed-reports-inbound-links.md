---
id: "261108-3t7t"
title: move closed reports inbound links
status: done
kind: feature
priority: 4
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
---

# move closed reports inbound links

## Goal

`move_file` now reports the inbound links it leaves broken (261029-m0vm); `move_closed`, which is a `move_file` per closed record, reports none. Moving closed records out of a folder breaks links to them in the same way.

Done when:

- [x] `move_closed` lists the inbound links it leaves broken, in `notes`, in the form `move_file` uses, without one line per record growing past a bound (decide the bound, as `mint_ids` does for references).
- [x] A test with a link to a record that moves.

## Next

## Outcome

move_closed lists the relative links and path-qualified wikilinks it leaves broken in notes, in move_file's form, read in one pass before any file moves; links between records that move together are not listed; the list ends at 50 lines with a count. Tested with a link to a moved record, a link between two moved records and the bound. Like move_file it reads docs/ and AGENTS.md only.

## Log

- 2026-11-08: filed.
- 2026-11-08 14:21-07:00: claimed (claude-sonnet-5-5-rest).
- 2026-11-08 14:21-07:00: move_closed reads the links of all its moves in one pass (inboundLinksOfMoves), before any file moves, and lists them in notes in move_file's form; a link between two records that move together is not listed; past 50 lines the notes end in '... and N more inbound links left broken'. move_file is unchanged and unbounded. Like move_file it scans docs/ and AGENTS.md only (what L1 reads). Tests TestMoveClosedReportsInboundLinks and TestMoveClosedBoundsInboundNotes.
