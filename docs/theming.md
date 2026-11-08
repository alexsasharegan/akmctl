# Adding a colour theme to the workbench

A colour theme is one CSS file and one line of data. No other file changes. The format is decided in ADR 261108-w3xf; the words are in `docs/glossary.md` (colour theme, theme variant, colour mode).

## The files

- `internal/server/web/themes/<id>.css`: the tokens, one block per variant.
- `internal/server/web/themes/themes.json`: the list the picker reads. The server also joins every `themes/*.css` into `/themes.css`, which the page loads before `app.css`.

`app.css` holds structure and the tokens that are not colours (fonts, spacing, radii, durations). It never holds a colour.

## Steps

1. Pick an id: lower-case letters, digits and hyphens (`catppuccin`, `atom-one`).
2. Copy `themes/quiet-paper.css` to `themes/<id>.css` and change the selectors and the values.
3. Add the theme to `themes/themes.json`: its `id`, a `name` for the picker, and its `variants`. Each variant has an `id` (unique within the theme), a `name` and a `mode` (`light` or `dark`). A theme needs at least one of each mode.
4. Run `go test ./internal/server`. It tells you which token a block misses or adds.
5. Run `akmctl ui`, open the command palette (`Cmd+K` or `Ctrl+K`), type `colour`, and pick the theme, then each mode and variant. Check the board, a record's drawer, a review diff and the palette itself in light and dark.

## The selectors

The first variant of a mode in `themes.json` is the theme's default for that mode. Its block has no `data-variant`; every other variant of the mode adds one:

```css
/* default light */
:root[data-theme="demo"][data-mode="light"] { color-scheme: light; /* tokens */ }
/* default dark */
:root[data-theme="demo"][data-mode="dark"] { color-scheme: dark; /* tokens */ }
/* a second dark variant */
:root[data-theme="demo"][data-mode="dark"][data-variant="night"] { color-scheme: dark; /* tokens */ }
```

`color-scheme` must match the variant's `mode`. The default theme's default light block also carries a bare `:root` in its selector list (see `quiet-paper.css`); copy that only into the default theme.

## The tokens

Every block defines every token, with no extras; the test compares each block with the default theme's default light block. The groups, as `quiet-paper.css` writes them:

- Surfaces: `--bg` (the page), `--surface` (cards and the drawer), `--surface-2`, `--input-bg`, `--tray`.
- Text and lines: `--text`, `--text-2` (muted text), `--line`, `--line-strong`.
- Accent: `--accent`, `--accent-fg` (text on the accent), `--accent-soft`, `--focus`, `--scrim` (behind a dialog).
- Warning and error: `--warn-fg`, `--warn-bg`, `--err-fg`, `--err-bg`.
- Small parts: `--kbd-bg`, `--kbd-fg`, `--pal-sel` (the selected palette row), `--mark-bg` (a search match), `--toast-bg`, `--toast-fg`.
- Traffic lights (the check summary): `--tl-success`, `--tl-info`, `--tl-warning`, `--tl-error`.
- Status colours of the board lanes: `--st-raw`, `--st-ready`, `--st-active`, `--st-blocked`. They must stay distinguishable from one another, including for colour-blind readers.
- The review diff: `--add-fg`, `--add-bg`, `--rem-fg`, `--rem-bg`, `--hunk-bg`.
- Shadows: `--sh-ink` (an `r,g,b` triple), `--sh-sm`, `--sh-lg`, `--sh-2`, `--card-hi` (a highlight inside a card; `none` if the theme has none).

A new token is added to every theme in the same change.

## Where the colours come from

Take them from the theme's own palette, not from a terminal theme: a terminal palette has no surface steps or borders. Each theme's work item records the sources it used, with their licences. Check the contrast of text on every surface (WCAG AA, 4.5:1 for text) in each variant, and note it in the work item.

## What the picker keeps

The choice lives in the browser (`localStorage`): `akm-theme` (theme id), `akm-mode` (`system`, `light` or `dark`) and `akm-variant-light` and `akm-variant-dark` (a variant other than the default). A saved id that no longer exists falls back to the default theme. `t` still switches light and dark.
