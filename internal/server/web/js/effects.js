// Opt-in animation effects (the spike of work item 261101-4n0v): the preferences
// that choose a lane-move theme and switch the confetti cannon on, and the cannon
// itself. Presentation only: nothing here touches a record. The preferences and the
// pieces are pure, so node can test them; the burst at the bottom only draws.
//
// Quiet by default: the theme is "quiet" and the confetti is off until the reader
// turns it on, and the animator plays nothing at all under prefers-reduced-motion.

export const THEMES = ["quiet", "tetris"];
export const PREFS_KEY = "akm-anim";

/** The stored preferences, made valid: an unknown theme is "quiet", confetti is on only when true. */
export function parsePrefs(raw) {
  let o = {};
  try {
    const v = JSON.parse(raw || "{}");
    if (v && typeof v === "object") o = v;
  } catch (e) {
    /* damaged storage reads as the defaults */
  }
  return { theme: THEMES.includes(o.theme) ? o.theme : "quiet", confetti: o.confetti === true };
}

/** The theme after `theme` in the list, wrapping round. */
export const nextTheme = (theme) => THEMES[(THEMES.indexOf(theme) + 1) % THEMES.length];

export function loadPrefs() {
  try {
    return parsePrefs(localStorage.getItem(PREFS_KEY));
  } catch (e) {
    return parsePrefs("");
  }
}

export function savePrefs(p) {
  try {
    localStorage.setItem(PREFS_KEY, JSON.stringify(p));
  } catch (e) {
    /* storage may be blocked: the choice then lasts until the page closes */
  }
}

// Mid-tone colours that read on both the light and the dark surface.
export const CONFETTI_COLORS = ["#E4572E", "#F3A712", "#29A19C", "#5C7AEA", "#C04CFD", "#7BC043"];

/** A small repeatable random generator (mulberry32): the same seed gives the same burst. */
export function rng(seed) {
  let a = seed >>> 0;
  return () => {
    a = (a + 0x6d2b79f5) >>> 0;
    let t = a;
    t = Math.imul(t ^ (t >>> 15), t | 1);
    t ^= t + Math.imul(t ^ (t >>> 7), t | 61);
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
  };
}

/**
 * The pieces of one burst. Each is { dx, rise, fall, spin, w, h, color, delay, dur }:
 * it flies from the origin to (dx, -rise), then falls to (dx * 1.15, fall) while it
 * spins and fades. The fan opens upward, like a cannon.
 */
export function confettiPieces(n = 70, seed = 1) {
  const rand = rng(seed);
  const out = [];
  for (let i = 0; i < n; i++) {
    const angle = (-90 + (rand() - 0.5) * 150) * (Math.PI / 180); // up, spread 75 degrees each way
    const power = 140 + rand() * 260;
    out.push({
      dx: Math.cos(angle) * power,
      rise: -Math.sin(angle) * power * 0.8,
      fall: 120 + rand() * 260,
      spin: (rand() - 0.5) * 1080,
      w: 5 + Math.round(rand() * 5),
      h: 8 + Math.round(rand() * 6),
      color: CONFETTI_COLORS[Math.floor(rand() * CONFETTI_COLORS.length)],
      delay: Math.round(rand() * 120),
      dur: 1100 + Math.round(rand() * 700),
    });
  }
  return out;
}

/** Fire one burst from the middle of `rect`. Resolves when the last piece is gone. */
export function burst(rect, seed = Date.now()) {
  if (typeof document === "undefined" || !document.body.animate) return Promise.resolve();
  const layer = document.createElement("div");
  layer.className = "fx-layer";
  layer.setAttribute("aria-hidden", "true");
  const ox = rect.left + rect.width / 2;
  const oy = rect.top + rect.height / 2;
  const done = confettiPieces(70, seed).map((p) => {
    const bit = document.createElement("i");
    bit.className = "fx-bit";
    for (const [k, v] of [["left", ox + "px"], ["top", oy + "px"], ["width", p.w + "px"], ["height", p.h + "px"], ["background", p.color]]) bit.style.setProperty(k, v);
    layer.append(bit);
    const frames = [
      { transform: "translate(0,0) rotate(0deg)", opacity: 1, offset: 0, easing: "cubic-bezier(.1,.7,.4,1)" },
      { transform: "translate(" + p.dx + "px," + -p.rise + "px) rotate(" + p.spin / 2 + "deg)", opacity: 1, offset: 0.4, easing: "cubic-bezier(.5,0,.9,.6)" },
      { transform: "translate(" + p.dx * 1.15 + "px," + p.fall + "px) rotate(" + p.spin + "deg)", opacity: 0, offset: 1 },
    ];
    return bit.animate(frames, { duration: p.dur, delay: p.delay, fill: "forwards" }).finished.catch(() => {});
  });
  document.body.append(layer);
  // a page that is not being painted never finishes an animation: do not wait for it
  return Promise.race([Promise.all(done), new Promise((r) => setTimeout(r, 2400))]).then(() => layer.remove());
}
