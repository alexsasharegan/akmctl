---
id: "261107-49ta"
title: Animation settings in the workbench
summary: A settings control for the move theme, confetti and speed, with an off switch, instead of palette actions
status: raw
kind: feature
created: 2026-11-07
updated: 2026-11-07
---

# Animation settings in the workbench

## Goal

Replace the three palette actions from the 4n0v spike with a small settings control in the workbench chrome: a theme choice, a confetti switch, a speed control (a slow watch speed) and an off switch, remembered in localStorage (key `akm-anim`, see `js/effects.js`). Quiet stays the default and prefers-reduced-motion still plays nothing.

Done when:

- [ ] the control works by keyboard and has labels a screen reader reads
- [ ] speed scales every theme's durations and off plays nothing
- [ ] checked in the browser in light and dark

## Log

- 2026-11-07: filed.
