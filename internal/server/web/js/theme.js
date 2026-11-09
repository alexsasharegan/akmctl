// Color themes: which theme, mode and variant the page shows. Pure functions
// (checked by testdata/logic.test.mjs); ui.js applies and stores the result and
// boot.js repeats resolvePrefs for the first paint (the test keeps them equal).
//
// A theme (themes/<id>.css, listed in themes/themes.json) has variants, each
// marked light or dark. The first variant of a mode is the theme's default for
// that mode and needs no data-variant; any other is chosen by data-variant.
// mode is what the person chose: "system" follows prefers-color-scheme.

export const DEFAULT_THEME = "quiet-paper";
export const MODES = ["system", "light", "dark"];
/** localStorage keys. akm-theme once held "light" or "dark"; those values still work. */
export const KEYS = { theme: "akm-theme", mode: "akm-mode", light: "akm-variant-light", dark: "akm-variant-dark" };

const ID = /^[a-z0-9][a-z0-9-]*$/;

/**
 * Read the stored strings into {theme, mode, resolved, variant}. stored holds the
 * raw values by name (theme, mode, light, dark), each a string or null. Unknown
 * or odd values fall back, so a hand-edited store never breaks the page.
 */
export function resolvePrefs(stored, systemDark) {
  const s = stored || {};
  let theme = DEFAULT_THEME;
  let mode = null;
  if (s.theme === "light" || s.theme === "dark") mode = s.theme; // the old two-value key
  else if (ID.test(s.theme || "")) theme = s.theme;
  if (MODES.includes(s.mode)) mode = s.mode;
  if (mode === null) mode = "system";
  const resolved = mode === "system" ? (systemDark ? "dark" : "light") : mode;
  const variant = ID.test(s[resolved] || "") ? s[resolved] : "";
  return { theme, mode, resolved, variant };
}

/** The registry's variants of a theme in one mode, default first; [] for an unknown theme. */
export function variantsOf(registry, themeId, mode) {
  const t = ((registry && registry.themes) || []).find((x) => x.id === themeId);
  return t ? t.variants.filter((v) => v.mode === mode) : [];
}

/**
 * Check prefs against the registry: an unknown theme becomes the default theme,
 * and a variant the theme lacks in this mode is dropped (the default shows).
 */
export function checkPrefs(prefs, registry) {
  const themes = (registry && registry.themes) || [];
  const fallback = (registry && registry.default) || DEFAULT_THEME;
  const theme = themes.some((t) => t.id === prefs.theme) ? prefs.theme : fallback;
  const known = variantsOf(registry, theme, prefs.resolved).slice(1); // the default needs no name
  const variant = theme === prefs.theme && known.some((v) => v.id === prefs.variant) ? prefs.variant : "";
  return { ...prefs, theme, variant };
}

/** The attribute values for <html>: data-theme, data-mode, data-variant ("" = none). */
export const attrsOf = (p) => ({ theme: p.theme, mode: p.resolved, variant: p.variant });
