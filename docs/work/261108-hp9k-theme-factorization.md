---
id: "261108-hp9k"
title: theme factorization
status: active
kind: feature
parent: 261108-gvzn
created: 2026-11-08
updated: 2026-11-08
summary: Refactor the workbench styling so a theme is one self-contained token set and adding one touches no other file
claimed_by: claude-sonnet-5-5
claimed_at: 2026-11-08T22:48
---

# theme factorization

## Goal

Today `app.css` holds the palette as CSS custom properties in `:root`, then repeats the whole dark block twice (`:root[data-theme=dark]` and the `prefers-color-scheme` fallback), and a few later rules (the review diff colours) add their own light, dark and fallback copies. `boot.js` and `ui.js` know only "light", "dark" or the system (key `akm-theme`). Refactor so that a theme is one self-contained unit and adding one is a data change, not a hunt through the stylesheet.

Questions to settle in the work, and record in an ADR if the answer is hard to reverse: how a theme is declared and registered (a CSS file per theme, a token block per theme, or data the server serves); how a theme pairs its light and dark variants; the variant model (the discovery below shows a theme has one light variant and one to three dark ones, so a theme is a set of variants, each marked light or dark, and "system" picks the default light and the default dark variant of the chosen theme), the picker and how `akm-theme` migrates; the one token set every theme must supply, with a check or test that a theme leaves none out; how the system-preference fallback stays without repeating each block.

Done when:

- [ ] a theme is a single self-contained unit, and the current palette is converted to it with no visible change
- [ ] no colour literal is left in `app.css` outside a theme (the review diff colours move into the tokens)
- [ ] the workbench has a picker for the theme and its variant, with a system mode that follows light or dark, remembered in the browser; an existing `akm-theme` value still works
- [ ] a test fails when a theme omits a token, and a short how-to for adding a theme is written
- [ ] `docs/glossary.md` has the new terms (colour theme, theme variant) and says how they differ from the move theme
- [ ] checked in the browser in light and dark, with no flash of the wrong theme on load

## Next

Read internal/server/web/app.css (tokens at the top, the review-diff block near line 278), js/boot.js and js/ui.js (akm-theme), then draft the theme-unit design and an ADR if the choice is hard to reverse.

## Discovery

Findings from the palette discovery (details in each theme item's Discovery section):

- Ghostty's bundled themes (`/Applications/Ghostty.app/Contents/Resources/ghostty/themes/`) give each variant's background, foreground, cursor, selection and 16 ANSI colours. They cover every requested theme.
- Variants differ in number: Catppuccin 1 light + 3 dark, Ayu 1 light + 2 dark, Atom One 1 + 1, GitHub 1 light + 1 dark default plus dimmed, high-contrast and colorblind variants. The theme unit must therefore hold several variants, each marked light or dark.
- A terminal palette lacks the workbench's tiers (surface, surface-2, tray, line, text-2, input, scrim, tl-* and the review-diff pairs). The theme unit should say which tokens a theme must supply and which derive from the base (for example with `color-mix`), so a new theme needs only its base, text, accent, status and semantic hues. Weigh that against a theme being able to override any token.
- Status and warning colours must stay distinguishable in every variant (the GitHub colorblind variants swap red and green for orange and blue).

- A second source, the VS Code themes installed on this machine (Catppuccin, GitHub Primer, Atom One), supplies the tiers a terminal palette lacks. Across all of them the pattern is the same: a page colour, a slightly different sidebar or widget colour, an input colour, a border, a muted text colour, an accent, and soft diff fills. That is the token set to design the theme unit around, and it shows the tiers are chosen per theme (Catppuccin steps through mantle and surface0; GitHub Dark puts the sidebar darker than the page), so deriving them from the base would not reproduce the official look. Lean toward each theme supplying every tier, with a derivation helper only as a convenience.

## Log

- 2026-11-08: filed.
- 2026-11-08 11:00-07:00: Palette discovery from Ghostty's bundled themes: variants, hex values and upstream sources recorded in Discovery.
- 2026-11-08 11:05-07:00: Second discovery pass from the installed VS Code themes (Catppuccin, GitHub, One Dark/Light); Ayu has none installed. Tiers, borders, accents and diff fills added to Discovery.
- 2026-11-08 15:48-07:00: claimed (claude-sonnet-5-5).
