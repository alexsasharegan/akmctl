---
id: "261101-4n0v"
title: spike workbench move and close animations
summary: Spike a set of playful lane-move and completion animations, including a confetti cannon on done, as opt-in themes
status: active
created: 2026-11-01
updated: 2026-11-07
kind: spike
priority: -1
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-08T05:52
---

# spike workbench move and close animations

## Goal

A spike on the wish-list animations for the workbench, from the brainstorm "Fast follow: animations" in `docs/design/pkm-workbench-design.md`. That section was never filed as an item; this is it. The plain slide-and-flash for a move exists (`js/anim.js`); the missed-move replay is its own, higher-priority item (261101-w0mz).

Two groups to try:

1. Moving a record from one lane to another: themed movement (a rocket that becomes a car, a Mario-style figure along the lane, Frogger hops, Pac-Man, a Tetris drop, a train or marble run).
2. Clearing a task by completing it: a **confetti cannon** on `done` (the standard name for the effect Alexander wants), the "Vegas / power mode" extras (dimmed lights, sweeping spotlights, small explosions), a quiet fizzle for `dropped`, and streaks or combos for several closes in a row.

Constraints from the brainstorm: animation is presentation only (the server emits a move event, a theme draws it); the queue plays one event at a time; a speed control and an off switch; quiet by default, `prefers-reduced-motion` honoured, power themes opt-in, sound off by default; no new dependency.

Done when:

- [ ] a confetti cannon prototype on `done` runs in the browser, light and dark, with reduced motion on and off
- [ ] one lane-move theme prototype shows how a theme plugs into the event queue
- [ ] the Outcome picks what to build, how themes are chosen and stored, and files the feature items

## Next

Ready. Scope agreed: a confetti cannon on `done` and one lane-move theme (a Tetris-style drop) as prototypes behind an opt-in switch, then the Outcome.

## Log

- 2026-11-01: filed. The idea lived only in the design doc's brainstorm; Alexander added the confetti cannon.
- 2026-11-07 22:52-07:00: Refined with Alexander: the spike builds a confetti cannon on done and one lane-move theme, a Tetris-style drop, to show how a theme plugs into the anim.js queue; the Outcome files the feature items. w0mz (replay of missed moves) is done.
- 2026-11-07 22:52-07:00: claimed (claude-sonnet-5-5).
