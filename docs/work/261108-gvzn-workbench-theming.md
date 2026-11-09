---
id: "261108-gvzn"
title: Workbench theming
status: ready
kind: feature
created: 2026-11-08
updated: 2026-11-08
summary: Make the workbench's look themeable so a new color theme is a small, plug-in addition, then ship the classic themes in light and dark
---

# Workbench theming

## Goal

The workbench look is fixed to one palette ("Quiet paper", light and dark) whose tokens live in `internal/server/web/app.css`. Make the styling themeable, so that adding a theme is plug-and-play, then ship the classic themes the maintainer uses. Every theme has a light and a dark variant.

Children, in dependency order:

- [Theme factorization](done/261108-hp9k-theme-factorization.md): the refactor that makes a theme one self-contained unit; the rest wait for it
- [Catppuccin themes](done/261108-m4a6-catppuccin-themes.md): Latte and the dark flavours
- [Ayu themes](261108-zg9m-ayu-themes.md): light, mirage and dark
- [Atom One themes](261108-pmkp-atom-one-themes.md): light and dark
- [GitHub themes](261108-dv5e-github-themes.md): light and dark

More themes can be added as children when wanted. A _color theme_ (this project) is not a _move theme_ (the lane-move animation in `js/anim.js`, see the glossary); the factorization item registers the distinction.

Done when:

- [ ] every child is `done` or `dropped`
- [ ] the workbench offers every shipped theme in a light and a dark variant, and the choice is kept across reloads

## Next

Work the four theme items in any order (Catppuccin 261108-m4a6 first, as the largest): the format is done, docs/theming.md is the how-to. Close this item when all four are done.

## Discovery

Palette discovery was done on 2026-11-08 from Ghostty's bundled themes, with the upstream palette sources named in each child. Variants found: Catppuccin Latte, Frappé, Macchiato, Mocha; Ayu Light, Mirage, Dark; Atom One Light and Dark; GitHub Light and Dark Default plus Dimmed, High Contrast and Colorblind. Each child's Discovery section has the hex values, the source files and the licence. Ghostty is a terminal palette, so surface tiers and lines are taken from the upstream sources (or derived) when each theme is built; the factorization item decides which tokens derive.

## Log

- 2026-11-08: filed.
- 2026-11-08 11:00-07:00: Palette discovery from Ghostty's bundled themes: variants, hex values and upstream sources recorded in Discovery.
- 2026-11-08 11:05-07:00: Second discovery pass from the installed VS Code themes (Catppuccin, GitHub, One Dark/Light); Ayu has none installed. Tiers, borders, accents and diff fills added to Discovery.
- 2026-11-08 11:21-07:00: Ayu second source found in Zed's bundled theme (see 261108-zg9m).
