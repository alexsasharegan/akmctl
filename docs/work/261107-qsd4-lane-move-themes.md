---
id: "261107-qsd4"
title: More lane-move themes
summary: Add the rocket-to-car, train or Pac-Man lane-move themes as opt-in entries in the theme registry
status: raw
kind: feature
created: 2026-11-07
updated: 2026-11-07
---

# More lane-move themes

## Goal

Add one or two themes to `themes` in `js/anim.js` beside quiet and tetris (the 4n0v spike): the rocket that turns into a car, a train carrying the card, or Pac-Man eating it and spitting it out. A theme gets the move event and the refreshed records, moves the card itself and returns true, so none touches the server. Each is opt-in, plays nothing under prefers-reduced-motion and adds no dependency.

Done when:

- [ ] at least one new theme plays in the browser in light and dark
- [ ] a theme with drawn parts uses inline SVG, no image files
- [ ] the settings control lists the theme

## Log

- 2026-11-07: filed.
