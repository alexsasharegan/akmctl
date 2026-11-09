---
id: "261108-gvzn"
title: Workbench theming
status: done
kind: feature
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
summary: Make the workbench's look themeable so a new color theme is a small, plug-in addition, then ship the classic themes in light and dark
---

# Workbench theming

## Goal

The workbench look is fixed to one palette ("Quiet paper", light and dark) whose tokens live in `internal/server/web/app.css`. Make the styling themeable, so that adding a theme is plug-and-play, then ship the classic themes the maintainer uses. Every theme has a light and a dark variant.

Children, in dependency order:

- [Theme factorization](261108-hp9k-theme-factorization.md): the refactor that makes a theme one self-contained unit; the rest wait for it
- [Catppuccin themes](261108-m4a6-catppuccin-themes.md): Latte and the dark flavors
- [Ayu themes](261108-zg9m-ayu-themes.md): light, mirage and dark
- [Atom One themes](261108-pmkp-atom-one-themes.md): light and dark
- [GitHub themes](261108-dv5e-github-themes.md): light and dark

More themes can be added as children when wanted. A _color theme_ (this project) is not a _move theme_ (the lane-move animation in `js/anim.js`, see the glossary); the factorization item registers the distinction.

Done when:

- [x] every child is `done` or `dropped`
- [x] the workbench offers every shipped theme in a light and a dark variant, and the choice is kept across reloads

## Next

## Discovery

Palette discovery was done on 2026-11-08 from Ghostty's bundled themes, with the upstream palette sources named in each child. Variants found: Catppuccin Latte, Frappé, Macchiato, Mocha; Ayu Light, Mirage, Dark; Atom One Light and Dark; GitHub Light and Dark Default plus Dimmed, High Contrast and Colorblind. Each child's Discovery section has the hex values, the source files and the license. Ghostty is a terminal palette, so surface tiers and lines are taken from the upstream sources (or derived) when each theme is built; the factorization item decides which tokens derive.

## Outcome

- The workbench's colors are a color theme: one CSS file of token blocks plus an entry in `themes.json` (ADR 261108-w3xf, work item 261108-hp9k). Four themes ship besides the default Quiet paper: Catppuccin (4 variants), Ayu (3), Atom One (2) and GitHub (7), each with at least one light and one dark variant.
- Each child recorded its contrast check (lowest pair 4.5 or better in every variant) and what was and was not looked at in the browser; afterwards the maintainer looked at every variant of every theme in the workbench and reported that all of them look fine (2026-11-08).
- The choice is kept across reloads in `localStorage` and resolved in `boot.js` before first paint; the flash was not measured.
- Left open: PR 67 holds the format and the other four PRs are stacked on it; nothing here is merged yet.

Decision: none beyond ADR 261108-w3xf, which stays `proposed` until PR 67 merges.

## Log

- 2026-11-08: filed.
- 2026-11-08 11:00-07:00: Palette discovery from Ghostty's bundled themes: variants, hex values and upstream sources recorded in Discovery.
- 2026-11-08 11:05-07:00: Second discovery pass from the installed VS Code themes (Catppuccin, GitHub, One Dark/Light); Ayu has none installed. Tiers, borders, accents and diff fills added to Discovery.
- 2026-11-08 11:21-07:00: Ayu second source found in Zed's bundled theme (see 261108-zg9m).
- 2026-11-08 19:42-07:00: claimed (claude-sonnet-5-5-theming).
- 2026-11-08 19:42-07:00: ticked "every child is `done` or `dropped`": children 261108-hp9k, m4a6, zg9m, pmkp and dv5e are done (PRs 67, 70, 71, 72, 73)
- 2026-11-08 19:42-07:00: ticked "the workbench offers every shipped theme in a light and a dark variant, and the choice is kept across reloads": themes.json offers quiet-paper, catppuccin, ayu, atom-one and github, each with a light and a dark variant; the choice is kept in localStorage (akm-theme, akm-mode, akm-variant-*)
- 2026-11-08 20:25-07:00: the maintainer looked at every variant of every theme and reported that all look fine; recorded in the Outcome and in each child (edited by hand, since `log` refuses closed items)
