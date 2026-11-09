---
id: "261101-w0mz"
title: workbench replay missed moves
summary: When you return to the workbench, play the moves that happened while you were away (agents, the shell) as animations so none pass unseen
status: done
created: 2026-11-01
updated: 2026-11-06
closed: 2026-11-06
kind: feature
priority: 1
blocked_by: [261106-sdnr]
---

# workbench replay missed moves

## Goal

Alexander wants moves made while he is away, or made by an agent while he watches, to animate so that coming back he can see what happened. This is higher priority than the animation themes (261101-4n0v), which are wish-list polish.

The design (`docs/design/pkm-workbench-design.md`, "Watch mode and the animation queue") already says the moves play one at a time, in order, on return. The built live board (261028-a149 step 5, `js/live.js`, `js/anim.js`) does part of this: polling pauses while the tab is hidden, the first poll on return diffs against the last snapshot, and each `moved` event animates. Check what is missing against the wish:

- the diff shows only the net change (raw to done), not the path an item took through several moves while away
- a batch over 10 events, or reduced motion, skips animation and gives one summary toast; decide whether a return should still replay a few, or a summary that can be played
- a hidden tab never finishes animations, so waits are capped; check a return after a long absence
- the order and kind (agent or outside) of the events, and whether a move out of a filtered view animates (261026-zg0g)
- whether a log of moves is needed (the `log` records) rather than a snapshot diff

Scope since 2026-11-06 (ADR 261106-60xv): the order, author and path of moves come from git through a new `history` op (261106-sdnr), so this item is the workbench replay on top of it and waits for that op. Without git the board keeps the snapshot diff and says why there is no replay.

Done when:

- [x] a gap list says what a returning user sees now, driven in the browser: one move, three, twelve, a multi-step path, after a hidden tab
- [x] the moves made while away play one at a time in order on return, each labeled by who made it
- [x] reduced motion and a very long batch degrade to a clear summary
- [x] unit tests cover the event ordering and the batching rule

## Gap list

Driven on 2026-11-06 in the built-in browser against a scratch repository (14 `ready` items, moves made from the shell with `akmctl set` and `close`; the hidden tab was simulated by overriding `document.visibilityState`).

| Scenario                                                          | What the user sees now                                                                                             | Gap                                                                                                                                            |
| ----------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| One move                                                          | A toast ("Outside: 261106-08ek Item k moved, ready to active.") and the card lifts, slides to its lane and settles | None. The label says "Outside", never who.                                                                                                     |
| Three moves between two polls                                     | Three toasts and three slides, one at a time, about 0.9 s apart                                                    | Order is id order, not the order the moves happened in (`diffSnapshots` sorts by id; the listing carries no time of change).                   |
| Twelve moves                                                      | One toast, "12 changes outside (12 moved).", no animation; the board is simply redrawn                             | The 12 moves are gone: the toast lists none and cannot be played. Over 10 events nothing replays, so the longer the absence the less is shown. |
| A four-step path on one item (ready, active, ready, active, done) | One move, "ready to done"                                                                                          | The path is invisible: the poll compares two snapshots. `set status` writes no `## Log` line, so no record of the path exists to replay.       |
| Hidden tab, three moves                                           | No requests while hidden (0 fetches in 9 s); on return one poll, then the three moves play one at a time           | None for a short absence. Over 10 moves the result is the summary above.                                                                       |
| Reduced motion                                                    | By code, not driven: `push` skips animation and shows the summary, or the single sentence                          | A reduced-motion user gets no way to review the moves either.                                                                                  |
| A move out of a filtered view                                     | Not driven: the filters (261026-zg0g) do not exist yet                                                             | `cardOf` returns null when the card is not drawn, so the toast shows and nothing animates; settle once filters land.                           |

What exists to build on: `claimed_by` and `claimed_at` (UTC, to the minute) are in the polled listing, so a move of a claimed item can be labeled with the agent that holds it; the etag is the only per-item change marker; nothing in the API says when a status changed.

## Next

## Outcome

Moves made while away are read from history, replayed one at a time in commit order with who made each, including an item's path; a burst over ten or reduced motion shows a summary toast and lists every move in the away panel; a page load lists what was committed since the last visit. Verified in the browser on a scratch repo and by unit tests of the merge, order and batching. Not done: a round trip with no net change is not shown; no history without git or with git_commits off, where the old snapshot diff stays; reduced motion was reasoned from the code, not driven; the filtered-view case waits for the filters.

## Log

- 2026-11-01: filed from Alexander's note: the away-and-agent moves matter more than the wish-list polish.
- 2026-11-06 12:08-07:00: claimed (claude-sonnet-5-5).
- 2026-11-06 12:11-07:00: Drove the board in the browser and wrote the gap list: id order not move order, bursts over 10 are summarized and lost, a path shows as its net move, 'Outside' never names an agent; claimed_by is available to label agents.
- 2026-11-06 12:22-07:00: released (claude-sonnet-5-5).
- 2026-11-06 12:37-07:00: claimed (claude-sonnet-5-5).
- 2026-11-06 12:42-07:00: Built the replay: replay.js (merge history moves into the poll's events, who, the away panel), history fetch, who-aware sentences, panel CSS; verified in the browser on a scratch repo: five committed moves with a three-step path played in commit order, labeled Claude Code and Codex; twelve moves gave a summary toast and a listed panel; a reload listed what was committed since the last visit; a hidden tab returned and replayed. Limits noted in the CHANGELOG.
