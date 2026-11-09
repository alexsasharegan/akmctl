---
id: "261108-dv5e"
title: github themes
status: done
kind: feature
parent: 261108-gvzn
created: 2026-11-08
updated: 2026-11-08
closed: 2026-11-08
summary: "Add GitHub's Primer themes to the workbench: Light and Dark Default first, then Dimmed, High Contrast and Colorblind"
---

# github themes

## Goal

Add the GitHub theme to the workbench on the structure the factorization item built. GitHub Light and GitHub Dark (the default pair; the dimmed and high-contrast flavors are out of scope unless they cost nothing). Map the palette onto every token of the theme unit; keep the status colors (raw, ready, active, blocked) and the warning and error pairs distinguishable and readable.

Done when:

- [x] the theme is a self-contained unit that passes the missing-token test
- [x] it appears in the picker with its light and dark variants
- [x] body text, muted text, the accent and the status colors meet WCAG AA contrast in each variant (record the check)
- [x] checked in the browser in light and dark across the board, the detail drawer, the review diff and the palette
- [x] the upstream palette's license is noted if it needs attribution

## Next

## Discovery

Source: Ghostty's bundled theme files (`/Applications/Ghostty.app/Contents/Resources/ghostty/themes/<name>`), read 2026-11-08. These are terminal palettes: background, foreground, cursor, selection and 16 ANSI colors. They fix the base, text and hue families, but not the workbench's surface tiers, lines and muted text (`--surface`, `--surface-2`, `--tray`, `--line`, `--text-2`), which are derived from the upstream palette below. Mapping used: red=ANSI 1, green=2, yellow=3, blue=4, magenta=5, cyan=6, muted=8 (bright black).

**Variants.** Ghostty carries two generations of GitHub themes. The Primer-based ones are the ones to use: `GitHub Light Default` and `GitHub Dark Default` (the default pair), plus `GitHub Dark Dimmed`, `GitHub Light High Contrast`, `GitHub Dark High Contrast`, `GitHub Light Colorblind` and `GitHub Dark Colorblind`. The older iTerm-era `GitHub` (light, bg #f4f4f4) and `GitHub Dark` (bg #101216) are not Primer palettes: leave them out. Plan: ship Light Default and Dark Default first; the dimmed, high-contrast and colorblind variants are a good fit for the workbench's accessibility, so ship them if the picker handles several variants per theme (the factorization item decides), otherwise file a follow-up.

| Role           | Light Default | Dark Default | Dark Dimmed | Light High Contrast | Dark High Contrast |
| -------------- | ------------- | ------------ | ----------- | ------------------- | ------------------ |
| background     | #ffffff       | #0d1117      | #22272e     | #ffffff             | #0a0c10            |
| foreground     | #1f2328       | #e6edf3      | #adbac7     | #0e1116             | #f0f3f6            |
| red            | #cf222e       | #ff7b72      | #f47067     | #a0111f             | #ff9492            |
| green          | #116329       | #3fb950      | #57ab5a     | #024c1a             | #26cd4d            |
| yellow         | #4d2d00       | #d29922      | #c69026     | #3f2200             | #f0b72f            |
| blue           | #0969da       | #58a6ff      | #539bf5     | #0349b4             | #71b7ff            |
| magenta        | #8250df       | #bc8cff      | #b083f0     | #622cbc             | #cb9eff            |
| cyan           | #1b7c83       | #39c5cf      | #39c5cf     | #1b7c83             | #39c5cf            |
| ANSI 0         | #24292f       | #484f58      | #545d68     | #0e1116             | #7a828e            |
| ANSI 8 (muted) | #57606a       | #6e7681      | #636e7b     | #4b535d             | #9ea7b3            |

Colorblind variants (`GitHub Light Colorblind`, `GitHub Dark Colorblind`) swap red and green for orange (#b35900 light, #ec8e2c dark) and blue; they matter for the status colors (raw, ready, active, blocked), which must stay distinguishable. In Ghostty the Light Default yellow (#4d2d00) is a dark brown, not a yellow: the warning color needs Primer's own attention scale.

**Upstream.** GitHub Primer primitives (github.com/primer/primitives, `primer/primitives`) and the VS Code theme (github.com/primer/github-vscode-theme) hold the full tokens (`canvas.default`, `canvas.subtle`, `border.default`, `fg.muted`, `accent.fg`, `attention.*`, `danger.*`) per variant. License: MIT. The VS Code extension above is the second source already read; the build step still records the upstream commit it uses.

**Second source: the installed VS Code extension** `github.github-vscode-theme` 6.3.5 (`~/.vscode/extensions/github.github-vscode-theme-6.3.5/themes/`, MIT; the nine files are `light`, `dark`, `light-default`, `dark-default`, `dark-dimmed`, `light-high-contrast`, `dark-high-contrast`, `light-colorblind`, `dark-colorblind`), read 2026-11-08. It is the Primer source, and it confirms that Ghostty's Default, Dimmed, High Contrast and Colorblind themes match Primer, while plain `light` and `dark` are the older themes. It adds the surfaces, border, muted text and soft diff fills:

| Role                     | Light Default | Dark Default | Dark Dimmed | Light HC  | Dark HC   | Light CB  | Dark CB   |
| ------------------------ | ------------- | ------------ | ----------- | --------- | --------- | --------- | --------- |
| page (editor.background) | #ffffff       | #0d1117      | #22272e     | #ffffff   | #0a0c10   | #ffffff   | #0d1117   |
| sidebar and panel        | #f6f8fa       | #010409      | #1c2128     | #ffffff   | #010409   | #f6f8fa   | #010409   |
| widget                   | #ffffff       | #161b22      | #2d333b     | #ffffff   | #272b33   | #ffffff   | #161b22   |
| border                   | #d0d7de       | #30363d      | #444c56     | #20252c   | #7a828e   | #d0d7de   | #30363d   |
| text                     | #1f2328       | #e6edf3      | #adbac7     | #0e1116   | #f0f3f6   | #24292f   | #c9d1d9   |
| muted text               | #656d76       | #7d8590      | #768390     | #0e1116   | #f0f3f6   | #57606a   | #8b949e   |
| focus                    | #0969da       | #1f6feb      | #316dca     | #0349b4   | #409eff   | #0969da   | #1f6feb   |
| link                     | #0969da       | #2f81f7      | #539bf5     | #0349b4   | #71b7ff   | #0969da   | #58a6ff   |
| error                    | #cf222e       | #f85149      | #e5534b     | #a0111f   | #ff6a69   | #b35900   | #d47616   |
| diff added fill          | #6fdd8b80     | #3fb9504d    | #57ab5a4d   | #43c66380 | #26cd4d4d | #80ccff80 | #58a6ff4d |
| diff removed fill        | #ff818266     | #ff7b724d    | #f470674d   | #ee5a5d66 | #ff94924d | #f7993966 | #ec8e2c4d |

The colorblind variants show the status trade-off directly: error becomes orange (#b35900 light, #d47616 dark) and "added" becomes blue, so the status colors (raw, ready, active, blocked) and the review-diff pairs must read from per-variant tokens, not a fixed red and green. High-contrast variants use a strong border (#20252c, #7a828e) and identical text and muted colors; check that nothing relies on muted text being fainter. Color values are the Primer functional tokens flattened for VS Code; the Primer primitives repo remains the source for names (`canvas.default`, `canvas.subtle`, `border.default`, `fg.muted`).

## Outcome

- `internal/server/web/themes/github.css` has seven blocks and passes `TestThemeRegistryAndTokens`: Light, Light high contrast, Light colorblind, Dark, Dark dimmed, Dark high contrast and Dark colorblind. Values: the Primer VS Code theme `github.github-vscode-theme` 6.3.5 (MIT), the second source recorded in Discovery; Ghostty's Primer themes agree on the hues.
- `themes.json` lists them with Light as the default light and Dark as the default dark; the rest are chosen from the command palette (search `variant`).
- Contrast (WCAG, computed by the generator over the resolved values): the lowest pair is 4.56 for Light, 4.55 for Light colorblind, 4.61 for Dark, 4.53 for Dark dimmed and 5.08 or more for the two high contrast variants, over text and muted text on the page, surface and hover tones, the accent and its text, warning, error, added and removed diff text on their tints, key caps, toasts, a search match and the four status colors. Colors that are too light or dark to read as text are mixed toward the text color with `color-mix`; those mixes are the only colors outside the sources, and the file header says so.
- Colorblind variants follow Primer: green becomes blue and red becomes orange, so active is blue and blocked and removed are orange. Ready would also be blue there, so it is purple (#8250df light, #bc8cff dark) to keep ready and active apart. In the high contrast variants muted text equals text, as in Primer, and Light high contrast has the same page and surface, so cards are told apart by the strong border only.
- Looked at in the browser: the board in Light and Dark colorblind and the record drawer in Dark. Not looked at: the palette, the review diff, the Check view, and Dimmed and both high contrast variants rendered.
- Unverified: whether the mixes look right to the maintainer's eye; warning and error are both brown-orange in Light colorblind (Primer's own choice), so check that pair; the contrast numbers are arithmetic, not a visual test.

Decision: none (the format is ADR 261108-w3xf).

## Log

- 2026-11-08: filed.
- 2026-11-08 11:00-07:00: Palette discovery from Ghostty's bundled themes: variants, hex values and upstream sources recorded in Discovery.
- 2026-11-08 11:05-07:00: Second discovery pass from the installed VS Code themes (Catppuccin, GitHub, One Dark/Light); Ayu has none installed. Tiers, borders, accents and diff fills added to Discovery.
- 2026-11-08 19:13-07:00: claimed (claude-sonnet-5-5-github).
- 2026-11-08 19:14-07:00: ticked "the theme is a self-contained unit that passes the missing-token test": themes/github.css has seven blocks; TestThemeRegistryAndTokens passes
- 2026-11-08 19:14-07:00: ticked "it appears in the picker with its light and dark variants": themes.json lists Light, Light high contrast, Light colorblind, Dark, Dark dimmed, Dark high contrast, Dark colorblind; switched in the browser through the stored prefs
- 2026-11-08 19:14-07:00: ticked "body text, muted text, the accent and the status colors meet WCAG AA contrast in each variant (record the check)": generator contrast report: lowest pair 4.55 light-cb, 4.56 light, 4.53 dimmed, 4.61 dark over text, muted text, accent, warning, error, diff text, status colors, toast and key caps on their surfaces; high contrast variants 5.08 and up
- 2026-11-08 19:14-07:00: ticked "checked in the browser in light and dark across the board, the detail drawer, the review diff and the palette": browser pane: board in Light and Dark colorblind, record drawer in Dark. Not looked at: palette, review diff, Check view, the other variants rendered
- 2026-11-08 19:14-07:00: ticked "the upstream palette's license is noted if it needs attribution": MIT (primer/github-vscode-theme); stated in the header of themes/github.css
