// Runs before the page paints (a classic script): sets the colour theme so
// there is no flash. A saved choice wins; otherwise the system preference.
(function () {
  var theme = "";
  try {
    theme = localStorage.getItem("akm-theme") || "";
  } catch (e) {
    theme = "";
  }
  if (theme !== "dark" && theme !== "light") {
    theme = window.matchMedia && matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }
  document.documentElement.dataset.theme = theme;
})();
