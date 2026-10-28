// Small helpers with no state. The pure ones (no DOM) are used by logic.js.

/** Build an element. Text goes in with textContent only, never as HTML. */
export function el(tag, props, ...kids) {
  const n = document.createElement(tag);
  for (const [k, v] of Object.entries(props || {})) {
    if (v === undefined || v === null || v === false) continue;
    if (k === "class") n.className = v;
    else if (k === "text") n.textContent = v;
    else if (k.startsWith("on") && typeof v === "function") n.addEventListener(k.slice(2), v);
    else if (k === "dataset") Object.assign(n.dataset, v);
    else n.setAttribute(k, v === true ? "" : String(v));
  }
  append(n, kids);
  return n;
}

export function append(node, kids) {
  for (const kid of kids.flat(Infinity)) {
    if (kid === null || kid === undefined || kid === false) continue;
    node.append(typeof kid === "string" || typeof kid === "number" ? document.createTextNode(String(kid)) : kid);
  }
  return node;
}

export function clear(node) {
  node.replaceChildren();
  return node;
}

const DAY = 86400000;

/** Today as a UTC date string, YYYY-MM-DD. `now` is a Date (tests pass one). */
export function isoDate(now = new Date()) {
  return now.toISOString().slice(0, 10);
}

/** Whole days from a YYYY-MM-DD date to `today` (also YYYY-MM-DD); null if either is not a date. */
export function ageDays(date, today = isoDate()) {
  const a = Date.parse(date + "T00:00:00Z");
  const b = Date.parse(today + "T00:00:00Z");
  if (!Number.isFinite(a) || !Number.isFinite(b)) return null;
  return Math.max(0, Math.round((b - a) / DAY));
}

/** The date `days` before `today`, YYYY-MM-DD. */
export function daysBefore(days, today = isoDate()) {
  return new Date(Date.parse(today + "T00:00:00Z") - days * DAY).toISOString().slice(0, 10);
}

/** "today", "1 day", "12 days", "3 months" ... for an age in days. */
export function ageText(days) {
  if (days === null || days === undefined) return "";
  if (days === 0) return "today";
  if (days === 1) return "1 day";
  if (days < 60) return days + " days";
  if (days < 730) return Math.round(days / 30) + " months";
  return Math.round(days / 365) + " years";
}

export function debounce(fn, ms) {
  let t;
  const d = (...a) => {
    clearTimeout(t);
    t = setTimeout(() => fn(...a), ms);
  };
  d.cancel = () => clearTimeout(t);
  return d;
}

export function clamp(n, lo, hi) {
  return Math.max(lo, Math.min(hi, n));
}

/** A field as text: null and undefined are "", lists are joined. */
export function text(v) {
  if (v === null || v === undefined) return "";
  return Array.isArray(v) ? v.join(", ") : String(v);
}
