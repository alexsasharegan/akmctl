---
id: "261108-w3xf"
title: color themes are css files with a registry
summary: A color theme is one CSS file of token blocks plus a registry entry; the server joins the files, so a new theme touches neither the page nor the stylesheet
status: proposed
date: 2026-11-08
---

# color themes are css files with a registry

## Context

The workbench had one palette ("Quiet paper") as custom properties in `app.css`: a light block, a dark block, the dark block again inside a `prefers-color-scheme` fallback, and the review-diff colors repeated the same three ways. `boot.js` and `ui.js` knew "light", "dark" or the system. Work item 261108-gvzn asks for themes that are plug-and-play, and its children add Catppuccin (1 light and 3 dark flavours), Ayu (1 and 2), Atom One (1 and 1) and GitHub (several), so a theme has more than one dark variant.

Constraints: the page has a strict content security policy (`style-src 'self'`, no inline style or script); the first paint must be right without a flash, which only a classic script in the head can do; the tool is a single binary with standard-library Go and no JavaScript build step; and `app.css` holds non-color tokens (fonts, spacing, radii, durations) that no theme should repeat.

## Considered Options

1. **One CSS file per theme, a registry, the server joins the files** (chosen).
2. One CSS file per theme, each linked from `index.html`: adding a theme edits the page, and the page test lists every link.
3. All themes in `app.css`: the file grows with every theme and a theme is a hunt through it, which is what this work removes.
4. Themes as JSON data, with the server or the script turning tokens into CSS: no stylesheet to write by hand, but it needs a generator or an inline style (refused by the policy), and a theme author loses plain CSS.
5. Derive most tokens from a base color with `color-mix`: a short theme, but the discovery of the real palettes (Catppuccin, GitHub Primer, One) shows the surface and border steps are chosen per theme and a derived set does not reproduce them.

## Decision

A **color theme** is one file, `internal/server/web/themes/<id>.css`, plus an entry in `themes/themes.json`. The file holds one block of custom properties per **variant**, each marked light or dark in the registry. The first variant of a mode is the theme's default for it: its block is `:root[data-theme="<id>"][data-mode="<mode>"]`; any other adds `[data-variant="<variant>"]`. Every block defines the same set of tokens and a matching `color-scheme`, and `TestThemeRegistryAndTokens` fails when a block adds or omits one (the reference set is the default theme's default light block).

The server joins every `themes/*.css` into `/themes.css` and serves the registry at `/themes/themes.json`; the page links `/themes.css` before `app.css`. `app.css` keeps structure and non-color tokens and holds no color literal (`TestAppCSSHasNoColorOfItsOwn`).

The default theme's default light block also answers for a bare `:root`, so a page whose theme id is unknown, or whose scripts did not run, is still colored. There is no `prefers-color-scheme` copy of any block: `boot.js` resolves the system preference and sets `data-theme`, `data-mode` and, when one was chosen, `data-variant` on `<html>` before the first paint.

The choice is kept in four `localStorage` keys: `akm-theme` (the theme id), `akm-mode` (`system`, `light` or `dark`, absent meaning system), and `akm-variant-light` and `akm-variant-dark` (a variant other than the default; cleared when the theme changes). The old `akm-theme` values `light` and `dark` are moved to `akm-mode` on first load. `js/theme.js` has the pure resolution; `boot.js` repeats it, and a test runs both on the same inputs.

Each theme supplies every token. No token is derived from a base color, so a theme can reproduce the official surface and border steps. A short helper may be added later without changing this format.

## Consequences

- Adding a theme is one new CSS file and one registry entry; no page, stylesheet or script changes. The how-to is `docs/theming.md`.
- A theme is about 40 tokens per variant, written by hand from the upstream palette; the test catches a missing token, not a poor choice of color. Contrast is checked per theme item.
- A user's saved choice survives a theme being removed: the unknown id falls back to the default theme.
- Pages that do not run scripts show the default light variant, not the system preference. The workbench does not work without scripts anyway.
- `boot.js` and `theme.js` hold the same resolution twice. A classic script cannot import, and a module cannot run before the first paint; the equality test is the guard.
- The set of tokens is part of the format: adding a token means adding it to every theme in one change.
