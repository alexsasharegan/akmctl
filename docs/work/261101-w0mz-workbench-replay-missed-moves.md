---
id: "261101-w0mz"
title: workbench replay missed moves
summary: When you return to the workbench, play the moves that happened while you were away (agents, the shell) as animations so none pass unseen
status: active
created: 2026-11-01
updated: 2026-11-01
kind: feature
priority: 1
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

Done when:

- [ ] a gap list says what a returning user sees now, driven in the browser: one move, three, twelve, a multi-step path, after a hidden tab
- [ ] the moves made while away play one at a time in order on return, each labelled by who made it
- [ ] reduced motion and a very long batch degrade to a clear summary
- [ ] unit tests cover the event ordering and the batching rule

## Next

The current `pollWork` listing exposes only each work item’s latest status, path and etag, so it cannot reconstruct several moves while the tab is hidden. Decide whether the API needs to expose move history; read `internal/server/web/js/live.js`, `internal/server/web/js/anim.js` and `internal/server/web/js/store.js`, then drive the board in a browser to document current behavior.

## Log

- 2026-11-01: filed from Alexander's note: the away-and-agent moves matter more than the wish-list polish.
- 2026-11-01 16:35-07:00: claimed (codex).
- 2026-11-01 16:36-07:00: released (codex).
- 2026-11-01 16:36-07:00: confirmed pollWork carries only the latest status, path and etag
- 2026-11-01 16:36-07:00: confirmed pollWork carries only the latest status, path and etag; the API has no replay history
