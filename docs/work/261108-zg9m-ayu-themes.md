---
id: "261108-zg9m"
title: ayu themes
status: ready
kind: feature
parent: 261108-gvzn
created: 2026-11-08
updated: 2026-11-08
summary: "Add Ayu to the workbench: Light, Mirage and Dark"
---

# ayu themes

## Goal

Add the Ayu theme to the workbench on the structure the factorization item built. Ayu has three flavours: light, mirage and dark. Ship light and one dark at least (mirage and dark both if the picker handles more than one dark variant). Map the palette onto every token of the theme unit; keep the status colors (raw, ready, active, blocked) and the warning and error pairs distinguishable and readable.

Done when:

- [ ] the theme is a self-contained unit that passes the missing-token test
- [ ] it appears in the picker with its light and dark variants
- [ ] body text, muted text, the accent and the status colors meet WCAG AA contrast in each variant (record the check)
- [ ] checked in the browser in light and dark across the board, the detail drawer, the review diff and the palette
- [ ] the upstream palette's licence is noted if it needs attribution

## Next

Read docs/theming.md, then this item's Discovery. Read the upstream palette named there, and write themes/<id>.css and the themes.json entry; run go test ./internal/server and look at it in the workbench.

## Discovery

Source: Ghostty's bundled theme files (`/Applications/Ghostty.app/Contents/Resources/ghostty/themes/<name>`), read 2026-11-08. These are terminal palettes: background, foreground, cursor, selection and 16 ANSI colors. They fix the base, text and hue families, but not the workbench's surface tiers, lines and muted text (`--surface`, `--surface-2`, `--tray`, `--line`, `--text-2`), which are derived from the upstream palette below. Mapping used: red=ANSI 1, green=2, yellow=3, blue=4, magenta=5, cyan=6, muted=8 (bright black).

**Variants.** Three flavours, one light and two dark: Light, Mirage (dark, the lighter of the two) and Dark. Ghostty names them `Ayu Light`, `Ayu Mirage` and `Ayu` (the last is the Dark flavour; its file has no suffix). Plan: ship all three; Dark is the default dark and Light the default light.

| Role                         | Light   | Mirage  | Dark    |
| ---------------------------- | ------- | ------- | ------- |
| background                   | #f8f9fa | #1f2430 | #0b0e14 |
| foreground                   | #5c6166 | #cccac2 | #bfbdb6 |
| cursor (the Ayu orange-gold) | #ffaa33 | #ffcc66 | #e6b450 |
| red                          | #ea6c6d | #ed8274 | #ea6c73 |
| green                        | #6cbf43 | #87d96c | #7fd962 |
| yellow                       | #eca944 | #facc6e | #f9af4f |
| blue                         | #3199e1 | #6dcbfa | #53bdfa |
| magenta                      | #9e75c7 | #dabafa | #cda1fa |
| cyan                         | #46ba94 | #90e1c6 | #90e1c6 |
| ANSI 0                       | #000000 | #171b24 | #11151c |
| ANSI 8 (muted)               | #686868 | #686868 | #686868 |

Ghostty's selection color for all three is a saturated blue (#035bd6 on Light, #409fff on the dark flavours), too strong for a workbench highlight: derive a soft one. Ghostty's ANSI 0 on the dark flavours (#171b24, #11151c) is a usable surface tone.

**Upstream.** The Ayu color scheme is at github.com/ayu-theme/ayu-colors (`ayu-colors`) and the VS Code port, which has the complete UI palette (`ui.bg`, `ui.panel.bg`, `ui.line`, `ui.fg`, `ui.selection.bg` and so on) for each flavour. Licence: MIT. Not yet read; the build step reads it and records the commit used.

**Second source: Zed's bundled Ayu.** Zed embeds `themes/ayu/ayu.json` in its binary (`/Applications/Zed.app/Contents/MacOS/zed` holds the path, but the contents are compressed and not readable in place). I took the same file from the public repository (`zed-industries/zed`, `assets/themes/ayu/ayu.json`, `main` as of 2026-11-08), so it may differ from the installed build. It is Zed's adaptation of Ayu with a full UI token set, not Ayu's own VS Code theme. It shows the tiers Ghostty lacks, with Zed's slightly different bases:

| Role                           | Light     | Mirage    | Dark      |
| ------------------------------ | --------- | --------- | --------- |
| page (editor.background)       | #fcfcfc   | #242835   | #0d1016   |
| surface and elements           | #ececed   | #353944   | #1f2127   |
| chrome (background)            | #dcddde   | #464a52   | #313337   |
| selected element               | #cfd0d2   | #53565d   | #3e4043   |
| border                         | #cfd1d2   | #53565d   | #3f4043   |
| border (variant, softer)       | #dfe0e1   | #43464f   | #2d2f34   |
| text                           | #5c6166   | #cccac2   | #bfbdb6   |
| muted text                     | #8b8e92   | #9a9a98   | #8a8986   |
| placeholder                    | #a9acae   | #7b7d7f   | #696a6a   |
| accent and info (blue)         | #3b9ee5   | #72cffe   | #5ac1fe   |
| error                          | #ef7271   | #f18779   | #ef7177   |
| warning                        | #f1ad49   | #fecf72   | #feb454   |
| success                        | #85b304   | #d5fe80   | #aad84c   |
| selection fill (accent at 24%) | #3b9ee53d | #72cffe3d | #5ac1fe3d |

Text colors agree with Ghostty for Light, Mirage and Dark (#5c6166, #cccac2, #bfbdb6). The page differs: Ghostty has #f8f9fa, #1f2430, #0b0e14; Zed has #fcfcfc, #242835, #0d1016, which are lighter than Ayu's published bases for the dark flavours. Zed also makes the text accent blue where Ghostty's cursor is the Ayu gold (#ffaa33, #ffcc66, #e6b450). Ayu's own VS Code theme was not installed, so the two disagree on the page color with no tie-breaker yet: the build step should read `ayu-colors` and pick the base from there (Ayu's published `ui.bg` values: check, do not trust memory), taking Zed's tier steps and Ghostty's hues as cross-checks.

## Log

- 2026-11-08: filed.
- 2026-11-08 11:00-07:00: Palette discovery from Ghostty's bundled themes: variants, hex values and upstream sources recorded in Discovery.
- 2026-11-08 11:05-07:00: Second discovery pass from the installed VS Code themes (Catppuccin, GitHub, One Dark/Light); Ayu has none installed. Tiers, borders, accents and diff fills added to Discovery.
- 2026-11-08 11:21-07:00: Found Ayu in Zed (embedded in the binary, compressed); read the same file from zed-industries/zed main. Added its tiers to Discovery; the page color disagrees with Ghostty, so the build reads ayu-colors.
