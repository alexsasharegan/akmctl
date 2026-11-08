---
id: "261101-4n0v"
title: spike workbench move and close animations
summary: Spike a set of playful lane-move and completion animations, including a confetti cannon on done, as opt-in themes
status: done
created: 2026-11-01
updated: 2026-11-07
closed: 2026-11-07
kind: spike
priority: -1
---

# spike workbench move and close animations

## Goal

A spike on the wish-list animations for the workbench, from the brainstorm "Fast follow: animations" in `docs/design/pkm-workbench-design.md`. That section was never filed as an item; this is it. The plain slide-and-flash for a move exists (`js/anim.js`); the missed-move replay is its own, higher-priority item (261101-w0mz).

Two groups to try:

1. Moving a record from one lane to another: themed movement (a rocket that becomes a car, a Mario-style figure along the lane, Frogger hops, Pac-Man, a Tetris drop, a train or marble run).
2. Clearing a task by completing it: a **confetti cannon** on `done` (the standard name for the effect Alexander wants), the "Vegas / power mode" extras (dimmed lights, sweeping spotlights, small explosions), a quiet fizzle for `dropped`, and streaks or combos for several closes in a row.

Constraints from the brainstorm: animation is presentation only (the server emits a move event, a theme draws it); the queue plays one event at a time; a speed control and an off switch; quiet by default, `prefers-reduced-motion` honoured, power themes opt-in, sound off by default; no new dependency.

Done when:

- [x] a confetti cannon prototype on `done` runs in the browser, light and dark, with reduced motion on and off
- [x] one lane-move theme prototype shows how a theme plugs into the event queue
- [x] the Outcome picks what to build, how themes are chosen and stored, and files the feature items

## Next

## Outcome

Both prototypes work and ship as opt-in features: a confetti cannon on done (js/effects.js, 70 seeded pieces, Web Animations, no dependency) and a tetris lane-move theme beside quiet. Decisions: a theme is a function in the anim.js themes registry that gets the move event and the refreshed records, moves the card and returns true, so the server stays out of it; the choice is a per-browser preference in localStorage (akm-anim), not a repository setting; quiet and no confetti are the defaults and reduced motion plays nothing. Follow-ups: 261107-49ta (settings control, speed, off switch), 261107-qsd4 (more themes), 261107-7v9y (dropped fizzle, streaks). Not checked: reduced motion with a real system setting (the page was stubbed), the tetris path by eye (sampled by script; the pane did not paint frames), and screen readers.

## Log

- 2026-11-01: filed. The idea lived only in the design doc's brainstorm; Alexander added the confetti cannon.
- 2026-11-07 22:52-07:00: Refined with Alexander: the spike builds a confetti cannon on done and one lane-move theme, a Tetris-style drop, to show how a theme plugs into the anim.js queue; the Outcome files the feature items. w0mz (replay of missed moves) is done.
- 2026-11-07 22:52-07:00: claimed (claude-sonnet-5-5).
- 2026-11-07 22:56-07:00: ticked "a confetti cannon prototype on `done` runs in the browser, light and dark, with reduced motion on and off": confetti cannon on done: 70 pieces fired on a close and drawn in light and dark; reduced motion (matchMedia stubbed, the pane has no switch) plays nothing; palette actions switch it on and off
- 2026-11-07 22:56-07:00: ticked "one lane-move theme prototype shows how a theme plugs into the event queue": tetris theme in anim.js themes registry: lift, stepped slide and drop, squash; chosen with anim.use and kept in akm-anim
- 2026-11-07 22:56-07:00: ticked "the Outcome picks what to build, how themes are chosen and stored, and files the feature items": Outcome picks it; follow-ups 261107-49ta, qsd4, 7v9y filed
