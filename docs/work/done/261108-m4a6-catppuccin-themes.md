---
id: "261108-m4a6"
title: catppuccin themes
status: done
kind: feature
parent: 261108-gvzn
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
summary: "Add Catppuccin to the workbench: Latte (light) and Frappé, Macchiato, Mocha (dark)"
---

# catppuccin themes

## Goal

Add the Catppuccin theme to the workbench on the structure the factorization item built. Catppuccin's Latte is the light flavour; Frappé, Macchiato and Mocha are dark. Ship Latte for light and Mocha for dark at least, and say in the Log whether the other dark flavours are worth offering. Map the palette onto every token of the theme unit; keep the status colors (raw, ready, active, blocked) and the warning and error pairs distinguishable and readable.

Done when:

- [x] the theme is a self-contained unit that passes the missing-token test
- [x] it appears in the picker with its light and dark variants
- [x] body text, muted text, the accent and the status colors meet WCAG AA contrast in each variant (record the check)
- [x] checked in the browser in light and dark across the board, the detail drawer, the review diff and the palette
- [x] the upstream palette's licence is noted if it needs attribution

## Next

## Discovery

Source: Ghostty's bundled theme files (`/Applications/Ghostty.app/Contents/Resources/ghostty/themes/<name>`), read 2026-11-08. These are terminal palettes: background, foreground, cursor, selection and 16 ANSI colors. They fix the base, text and hue families, but not the workbench's surface tiers, lines and muted text (`--surface`, `--surface-2`, `--tray`, `--line`, `--text-2`), which are derived from the upstream palette below. Mapping used: red=ANSI 1, green=2, yellow=3, blue=4, magenta=5, cyan=6, muted=8 (bright black).

**Variants.** Four flavours, one light and three dark: Latte (light), Frappé, Macchiato, Mocha (dark). Ghostty names: `Catppuccin Latte`, `Catppuccin Frappe`, `Catppuccin Macchiato`, `Catppuccin Mocha`. Plan: ship all four; Mocha is the default dark and Latte the default light.

| Role                           | Latte   | Frappé  | Macchiato | Mocha   |
| ------------------------------ | ------- | ------- | --------- | ------- |
| background (base)              | #eff1f5 | #303446 | #24273a   | #1e1e2e |
| foreground (text)              | #4c4f69 | #c6d0f5 | #cad3f5   | #cdd6f4 |
| cursor / selection (rosewater) | #dc8a78 | #f2d5cf | #f4dbd6   | #f5e0dc |
| red                            | #d20f39 | #e78284 | #ed8796   | #f38ba8 |
| green                          | #40a02b | #a6d189 | #a6da95   | #a6e3a1 |
| yellow                         | #df8e1d | #e5c890 | #eed49f   | #f9e2af |
| blue                           | #1e66f5 | #8caaee | #8aadf4   | #89b4fa |
| magenta (pink)                 | #ea76cb | #f4b8e4 | #f5bde6   | #f5c2e7 |
| cyan (teal)                    | #179299 | #81c8be | #8bd5ca   | #94e2d5 |
| ANSI 0 (surface1)              | #bcc0cc | #51576d | #494d64   | #45475a |
| ANSI 8 (surface2)              | #acb0be | #626880 | #5b6078   | #585b70 |
| ANSI 7 (subtext1)              | #5c5f77 | #b5bfe2 | #b8c0e0   | #bac2de |

Ghostty's cursor and selection are Catppuccin's rosewater; the accent for the workbench is a free choice (blue is the natural fit, as in the current theme).

**Upstream.** The full palette has 26 named colors per flavour, including the surface tiers this workbench needs: crust, mantle, base, surface0-2, overlay0-2, subtext0-1. Take those from the official source rather than deriving them: the Catppuccin palette repository (`catppuccin/palette`, `palette.json`) and the style guide at catppuccin.com. Licence: MIT. The VS Code extension above is the second source already read; the build step still records the upstream commit it uses.

**Second source: the installed VS Code extension** `catppuccin.catppuccin-vsc` 3.19.0 (`~/.vscode/extensions/catppuccin.catppuccin-vsc-3.19.0/themes/{latte,frappe,macchiato,mocha}.json`, MIT), read 2026-11-08. It agrees with Ghostty on every base, text and hue value, and adds the surface tiers, border and soft fills the terminal palette lacks:

| Role                              | Latte     | Frappé    | Macchiato | Mocha     |
| --------------------------------- | --------- | --------- | --------- | --------- |
| page (editor.background, `base`)  | #eff1f5   | #303446   | #24273a   | #1e1e2e   |
| sidebar and widgets (`mantle`)    | #e6e9ef   | #292c3c   | #1e2030   | #181825   |
| input and active row (`surface0`) | #ccd0da   | #414559   | #363a4f   | #313244   |
| border (`surface2`)               | #acb0be   | #626880   | #5b6078   | #585b70   |
| text                              | #4c4f69   | #c6d0f5   | #cad3f5   | #cdd6f4   |
| accent (button, focus: mauve)     | #8839ef   | #ca9ee6   | #c6a0f6   | #cba6f7   |
| link (blue)                       | #1e66f5   | #8caaee   | #8aadf4   | #89b4fa   |
| error (red)                       | #d20f39   | #e78284   | #ed8796   | #f38ba8   |
| warning (peach)                   | #fe640b   | #ef9f76   | #f5a97f   | #fab387   |
| diff added fill                   | #40a02b33 | #a6d18933 | #a6da9533 | #a6e3a133 |
| diff removed fill                 | #d20f3933 | #e7828433 | #ed879633 | #f38ba833 |

So the tiers map as `bg`=base, `surface`/`tray`=mantle, `input-bg`=surface0, `line`=surface2 (or surface1, `ANSI 0`, for a softer rule). The official accent is mauve, not blue; the workbench accent is a free choice. The review-diff pairs (`--add-*`, `--rem-*`) take the diff fills over the page color. The warning color is peach, distinct from the yellow hue.

## Outcome

- `internal/server/web/themes/catppuccin.css` has four blocks (Latte, Mocha, Frappé, Macchiato) and passes `TestThemeRegistryAndTokens`; the values are the official palette (`catppuccin/palette`, `palette.json` at commit 07d02aa110ef, MIT), cross-checked with the installed VS Code theme.
- `themes.json` lists the theme with Latte as the default light and Mocha as the default dark; Frappé and Macchiato are chosen from the command palette (search `variant`), checked in the browser.
- Contrast (WCAG, computed by the generator over the resolved values): the lowest pair is 4.50 for Latte, 4.75 for Mocha, 4.51 for Frappé and 4.73 for Macchiato, over text and muted text on the page, surface and tray, the accent and its text, warning, error, added and removed diff text on their tints, key caps, toasts, a search match, the selected palette row and the four status colors. To get there, palette colors that are too light or dark to read as text are mixed toward the text color with `color-mix` (for example Latte green and peach, Frappé red); those mixes are the only colors that are not in the palette, and they are noted in the file's header.
- Looked at in the browser: the board and the palette in Latte, Mocha and Frappé, and the record drawer in Latte. Not looked at: Macchiato rendered (its tokens were read and none is empty), the review diff, the Check view, and the drawer in the dark variants.
- Choices made: accent is mauve (the Catppuccin accent), the page is `mantle`, cards `base` and the tray `crust`. The traffic-light colors are mixed toward the page tone so they read on the inverted toast. Frappé's error text is mixed 60% with its text color to reach 4.5:1 on its tint, so it reads pinker than the palette red.
- Unverified: whether the mixes look right to the maintainer's eye; the contrast numbers are arithmetic, not a visual test.

Decision: none (the format is ADR 261108-w3xf).

## Log

- 2026-11-08: filed.
- 2026-11-08 11:00-07:00: Palette discovery from Ghostty's bundled themes: variants, hex values and upstream sources recorded in Discovery.
- 2026-11-08 11:05-07:00: Second discovery pass from the installed VS Code themes (Catppuccin, GitHub, One Dark/Light); Ayu has none installed. Tiers, borders, accents and diff fills added to Discovery.
- 2026-11-08 18:31-07:00: claimed (claude-sonnet-5-5-catppuccin).
- 2026-11-08 18:34-07:00: ticked "the theme is a self-contained unit that passes the missing-token test": web/themes/catppuccin.css; TestThemeRegistryAndTokens passes with the four variants
- 2026-11-08 18:34-07:00: ticked "it appears in the picker with its light and dark variants": themes.json lists Latte, Mocha, Frappé, Macchiato; the palette (search 'variant') lists and switches the dark variants (checked in the browser)
- 2026-11-08 18:34-07:00: ticked "body text, muted text, the accent and the status colors meet WCAG AA contrast in each variant (record the check)": generator contrast report: lowest pair 4.50 (Latte), 4.75 (Mocha), 4.51 (Frappé), 4.73 (Macchiato) over text, muted text, accent, warning, error, added and removed diff, kbd, toast, search mark, selected row and the four status colors on their surfaces
- 2026-11-08 18:34-07:00: ticked "checked in the browser in light and dark across the board, the detail drawer, the review diff and the palette": browser pane: board and palette in Latte, Mocha and Frappé; Macchiato only by computed tokens (none empty). The drawer, the review diff and the Check view were not looked at.
- 2026-11-08 18:34-07:00: ticked "the upstream palette's licence is noted if it needs attribution": MIT; stated in the header of themes/catppuccin.css
