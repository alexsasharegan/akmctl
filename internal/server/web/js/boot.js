// Runs before the page paints (a classic script): sets the colour theme so
// there is no flash. It repeats resolvePrefs from theme.js, which a classic
// script cannot import; testdata/logic.test.mjs runs both and compares them.
(function () {
  var KEYS = { theme: "akm-theme", mode: "akm-mode", light: "akm-variant-light", dark: "akm-variant-dark" };
  var ID = /^[a-z0-9][a-z0-9-]*$/;
  var s = {};
  try {
    for (var k in KEYS) s[k] = localStorage.getItem(KEYS[k]);
  } catch (e) {
    s = {};
  }
  var theme = "quiet-paper";
  var mode = null;
  if (s.theme === "light" || s.theme === "dark") mode = s.theme;
  else if (ID.test(s.theme || "")) theme = s.theme;
  if (s.mode === "system" || s.mode === "light" || s.mode === "dark") mode = s.mode;
  if (mode === null) mode = "system";
  var dark = window.matchMedia && matchMedia("(prefers-color-scheme: dark)").matches;
  var resolved = mode === "system" ? (dark ? "dark" : "light") : mode;
  var root = document.documentElement;
  root.dataset.theme = theme;
  root.dataset.mode = resolved;
  if (ID.test(s[resolved] || "")) root.dataset.variant = s[resolved];
  else delete root.dataset.variant;
})();
