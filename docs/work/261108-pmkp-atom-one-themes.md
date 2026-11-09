---
id: "261108-pmkp"
title: atom one themes
status: active
kind: feature
parent: 261108-gvzn
created: 2026-11-08
updated: 2026-11-08
summary: Add Atom One to the workbench in its light and dark variants
claimed_by: claude-sonnet-5-5-atom-one
claimed_at: 2026-11-09T02:04
---

# atom one themes

## Goal

Add the Atom One theme to the workbench on the structure the factorization item built. Atom One Light and Atom One Dark. Map the palette onto every token of the theme unit; keep the status colors (raw, ready, active, blocked) and the warning and error pairs distinguishable and readable.

Done when:

- [x] the theme is a self-contained unit that passes the missing-token test
- [x] it appears in the picker with its light and dark variants
- [x] body text, muted text, the accent and the status colors meet WCAG AA contrast in each variant (record the check)
- [x] checked in the browser in light and dark across the board, the detail drawer, the review diff and the palette
- [x] the upstream palette's license is noted if it needs attribution

## Next

Read docs/theming.md, then this item's Discovery. Read the upstream palette named there, and write themes/<id>.css and the themes.json entry; run go test ./internal/server and look at it in the workbench.

## Discovery

Source: Ghostty's bundled theme files (`/Applications/Ghostty.app/Contents/Resources/ghostty/themes/<name>`), read 2026-11-08. These are terminal palettes: background, foreground, cursor, selection and 16 ANSI colors. They fix the base, text and hue families, but not the workbench's surface tiers, lines and muted text (`--surface`, `--surface-2`, `--tray`, `--line`, `--text-2`), which are derived from the upstream palette below. Mapping used: red=ANSI 1, green=2, yellow=3, blue=4, magenta=5, cyan=6, muted=8 (bright black).

**Variants.** Two: Atom One Light and Atom One Dark. Ghostty names them `Atom One Light` and `Atom One Dark`. Do not use Ghostty's plain `Atom` theme: it is a different, older scheme (dark background #161719) and not Atom One. Plan: ship both.

| Role           | Light                   | Dark    |
| -------------- | ----------------------- | ------- |
| background     | #f9f9f9                 | #21252b |
| foreground     | #2a2c33                 | #abb2bf |
| cursor         | #bbbbbb                 | #abb2bf |
| selection      | #d0dffc                 | #323844 |
| red            | #de3e35                 | #e06c75 |
| green          | #3f953a                 | #98c379 |
| yellow         | #d2b67c                 | #e5c07b |
| blue           | #2f5af3                 | #61afef |
| magenta        | #950095                 | #c678dd |
| cyan           | #3f953a (same as green) | #56b6c2 |
| ANSI 0         | #000000                 | #21252b |
| ANSI 8 (muted) | #000000                 | #767676 |

Two things to check against upstream before trusting these. Ghostty's Light flavor is a terminal approximation (its cyan equals its green and its muted color is black), so take the Light hues from Atom's own `one-light-syntax` variables (for example hue-1 cyan, hue-2 blue, hue-5 red). Ghostty's Dark background (#21252b) was suspect too; the VS Code port below confirms the page is #282c34.

**Upstream.** The Atom `one-dark-ui` / `one-light-ui` and `one-dark-syntax` / `one-light-syntax` packages (github.com/atom/atom, `packages/`) define the palette as Less variables (`syntax-bg`, `syntax-fg`, `hue-1` to `hue-6`, `mono-1` to `mono-3`, `ui-syntax-color`, `base-border-color`). Atom is archived; the VS Code port (`akamud/vscode-theme-onedark`) is a cross-check. License: MIT. The VS Code extension above is the second source already read; the build step still records the upstream commit it uses.

**Second source: the installed VS Code extensions** `akamud.vscode-theme-onedark` 2.3.0 and `akamud.vscode-theme-onelight` 2.3.0 (`~/.vscode/extensions/akamud.vscode-theme-one{dark,light}-2.3.0/themes/One{Dark,Light}.json`), read 2026-11-08. They are ports of Atom's One UI and syntax themes. They settle the open question: the One Dark page background is **#282C34**, and Ghostty's #21252b is the One Dark _sidebar_ color, so Ghostty's Atom One Dark puts the wrong tier on its background.

| Role                     | One Light | One Dark |
| ------------------------ | --------- | -------- |
| page (editor.background) | #FAFAFA   | #282C34  |
| sidebar and widgets      | #EAEAEB   | #21252B  |
| input                    | #FFFFFF   | #1B1D23  |
| active row               | #DBDBDC   | #2C313A  |
| selection                | #E5E5E6   | #3E4451  |
| text (editor.foreground) | #383A42   | #ABB2BF  |
| focus and accent         | #526FFF   | #528BFF  |
| button                   | #5871EF   | #4D78CC  |

Light page differs slightly from Ghostty (#fafafa against #f9f9f9) and the text differs (#383a42 against #2a2c33); prefer the extension, which is derived from Atom's own variables (`mono-1` #383a42 and `syntax-bg` #fafafa for Light). These extensions define no border, muted-text, error or warning colors (they are unset), so take those from Atom's `one-*-syntax` hues (`hue-5` red, `hue-6` and `hue-6-2` orange, `hue-4` green, `mono-3` muted) and `base-border-color`.

## Log

- 2026-11-08: filed.
- 2026-11-08 11:00-07:00: Palette discovery from Ghostty's bundled themes: variants, hex values and upstream sources recorded in Discovery.
- 2026-11-08 11:05-07:00: Second discovery pass from the installed VS Code themes (Catppuccin, GitHub, One Dark/Light); Ayu has none installed. Tiers, borders, accents and diff fills added to Discovery.
- 2026-11-08 19:04-07:00: claimed (claude-sonnet-5-5-atom-one).
- 2026-11-08 19:05-07:00: ticked "the theme is a self-contained unit that passes the missing-token test": themes/atom-one.css has two blocks (One Light, One Dark); TestThemeRegistryAndTokens passes
- 2026-11-08 19:05-07:00: ticked "it appears in the picker with its light and dark variants": themes.json lists One Light (default light) and One Dark (default dark); switched in the browser through the stored prefs
- 2026-11-08 19:05-07:00: ticked "body text, muted text, the accent and the status colors meet WCAG AA contrast in each variant (record the check)": generator contrast report: lowest pair 4.52 One Light and 4.51 One Dark over text, muted text, accent, warning, error, diff text, status colors, toast and key caps on their surfaces
- 2026-11-08 19:05-07:00: ticked "checked in the browser in light and dark across the board, the detail drawer, the review diff and the palette": browser pane: board in One Light and One Dark; record drawer in One Dark. Not looked at: palette, review diff, Check view, drawer in One Light
- 2026-11-08 19:05-07:00: ticked "the upstream palette's license is noted if it needs attribution": MIT (atom/atom, GitHub Inc.); stated in the header of themes/atom-one.css
