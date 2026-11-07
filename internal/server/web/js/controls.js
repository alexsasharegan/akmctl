// Shared form controls: a custom select (the WAI-ARIA listbox pattern) and a
// summary text area with a character counter. Plain JavaScript, no dependencies.
// The decisions (counter level, key handling) are pure functions so node can test
// them; the builders at the bottom only wire them to the DOM.
import { el } from "./util.js";

export const SUMMARY_MAX = 160; // the spec's limit (WK11)
const NEAR = 0.9; // the counter warns from 90% of the limit

/**
 * The counter for a text of `len` characters against `max`.
 * Returns { level: "ok" | "near" | "limit" | "over", text, left }. The words change
 * with the level, so colour is never the only signal.
 */
export function counterState(len, max = SUMMARY_MAX) {
  const left = max - len;
  const head = len + " / " + max;
  if (left < 0) return { level: "over", left, text: head + " - " + -left + " over the limit" };
  if (left === 0) return { level: "limit", left, text: head + " - limit reached" };
  if (len >= Math.ceil(max * NEAR)) return { level: "near", left, text: head + " - " + left + " left" };
  return { level: "ok", left, text: head };
}

/** A pasted or typed line break is a space: a summary is one line. */
export const oneLine = (s) => String(s).replace(/\s*[\r\n]+\s*/g, " ");

export const TYPEAHEAD_MS = 700;

/**
 * One key press on a select. `st` is { open, active, buf, at } (active is an index
 * into `opts`, [{ value, label }]); `key` is the KeyboardEvent key and `now` a
 * millisecond clock. Returns the next state plus:
 *   handled - the key was used (preventDefault; Escape also stops the dialog closing),
 *   choose  - an index to select now, or null.
 * Escape closes an open list and is only handled then; Tab closes without choosing.
 */
export function selectKey(st, key, opts, now = 0) {
  const last = opts.length - 1;
  const next = { ...st, buf: "", at: now };
  const move = (i) => ({ ...next, open: true, active: Math.max(0, Math.min(last, i)) });
  if (last < 0) return { st: next, handled: false, choose: null };
  switch (key) {
    case "ArrowDown":
      return { st: move(st.open ? st.active + 1 : st.active), handled: true, choose: null };
    case "ArrowUp":
      return { st: move(st.open ? st.active - 1 : st.active), handled: true, choose: null };
    case "Home":
      return { st: move(0), handled: true, choose: null };
    case "End":
      return { st: move(last), handled: true, choose: null };
    case "Enter":
    case " ":
      if (!st.open) return { st: { ...next, open: true }, handled: true, choose: null };
      return { st: { ...next, open: false }, handled: true, choose: st.active };
    case "Escape":
      return { st: { ...next, open: false }, handled: st.open, choose: null };
    case "Tab":
      return { st: { ...next, open: false }, handled: false, choose: null };
  }
  if (key.length !== 1 || key === " ") return { st: next, handled: false, choose: null };
  // type-ahead: letters typed in quick succession make a prefix; one letter repeated cycles
  const buf = (now - st.at <= TYPEAHEAD_MS ? st.buf : "") + key.toLowerCase();
  const cycle = buf.length > 1 && [...buf].every((c) => c === buf[0]);
  const prefix = cycle ? buf[0] : buf;
  const from = buf.length === 1 || cycle ? st.active + 1 : st.active;
  for (let n = 0; n <= last; n++) {
    const i = (from + n) % opts.length;
    if (opts[i].label.toLowerCase().startsWith(prefix)) {
      return { st: { ...next, buf, open: true, active: i }, handled: true, choose: null };
    }
  }
  return { st: { ...next, buf }, handled: true, choose: null };
}

let uid = 0;

/**
 * Open `list` (inside `root`) below `anchor`, or above it when there is more room
 * there. Room is measured in the nearest scrolling box: a dialog body clips what
 * overflows it. Sets the "up" class on root and the list's max-height.
 */
export function placeList(anchor, list, root) {
  const r = anchor.getBoundingClientRect();
  const box = anchor.closest(".dlg-body, .drawer");
  const b = box ? box.getBoundingClientRect() : { top: 0, bottom: innerHeight };
  const below = b.bottom - r.bottom - 8;
  const above = r.top - b.top - 8;
  const up = below < 180 && above > below;
  root.classList.toggle("up", up);
  list.style.maxHeight = Math.max(96, Math.min(224, up ? above : below)) + "px";
}

/** Keep `li` inside the visible part of its scrolling `list`. */
export function keepInView(list, li) {
  if (li.offsetTop < list.scrollTop) list.scrollTop = li.offsetTop;
  else if (li.offsetTop + li.offsetHeight > list.scrollTop + list.clientHeight) list.scrollTop = li.offsetTop + li.offsetHeight - list.clientHeight;
}

/**
 * A custom select. props: { id, options: [{ value, label }], value, onChange(value) }.
 * Returns { root, button, value (get and set), focus() }. The button carries the id,
 * so a <label for> names it; the list is a listbox opened under it.
 */
export function createSelect(props) {
  const opts = props.options;
  const n = ++uid;
  let st = { open: false, active: Math.max(0, opts.findIndex((o) => o.value === props.value)), buf: "", at: 0 };
  let sel = st.active;
  const list = el("ul", { class: "sel-list", role: "listbox", id: "sel-list-" + n, tabindex: "-1", hidden: true });
  const label = el("span", { class: "sel-text" });
  const button = el(
    "button",
    { type: "button", id: props.id, class: "sel-btn", role: "combobox", "aria-haspopup": "listbox", "aria-expanded": "false", "aria-controls": list.id },
    label,
    el("span", { class: "sel-caret", "aria-hidden": "true", text: "▾" }),
  );
  const items = opts.map((o, i) =>
    el("li", {
      class: "sel-opt",
      role: "option",
      id: "sel-opt-" + n + "-" + i,
      text: o.label,
      onmousedown: (e) => e.preventDefault(), // keep the focus on the button
      onclick: () => {
        choose(i);
        close();
      },
      onmousemove: () => {
        if (st.active !== i) show({ ...st, active: i });
      },
    }),
  );
  list.append(...items);
  const root = el("div", { class: "sel" }, button, list);

  function show(next) {
    st = next;
    button.setAttribute("aria-expanded", String(st.open));
    list.hidden = !st.open;
    root.classList.toggle("open", st.open);
    items.forEach((li, i) => {
      li.setAttribute("aria-selected", String(i === sel));
      li.classList.toggle("active", st.open && i === st.active);
    });
    if (st.open) {
      button.setAttribute("aria-activedescendant", items[st.active].id);
      place();
      keepInView(list, items[st.active]);
    } else button.removeAttribute("aria-activedescendant");
    label.textContent = opts[sel] ? opts[sel].label : "";
  }
  const place = () => placeList(button, list, root);
  function choose(i) {
    const changed = i !== sel;
    sel = i;
    st = { ...st, active: i };
    if (changed && props.onChange) props.onChange(opts[i].value);
  }
  function close() {
    show({ ...st, open: false });
  }

  button.addEventListener("keydown", (e) => {
    if (e.metaKey || e.ctrlKey || e.altKey) return;
    // opening for the first time starts on the chosen option
    const from = st.open ? st : { ...st, active: sel };
    const r = selectKey(from, e.key, opts, e.timeStamp || Date.now());
    if (r.handled) {
      e.preventDefault();
      e.stopPropagation(); // the dialog must not see the Escape that closed this list
    }
    if (r.choose !== null) choose(r.choose);
    show(r.st);
  });
  button.addEventListener("click", () => show({ ...st, open: !st.open, active: sel }));
  root.addEventListener("focusout", (e) => {
    if (st.open && !root.contains(e.relatedTarget)) close();
  });
  show(st);

  return {
    root,
    button,
    get value() {
      return opts[sel] ? opts[sel].value : "";
    },
    set value(v) {
      const i = opts.findIndex((o) => o.value === v);
      if (i >= 0) {
        sel = i;
        show({ ...st, active: i });
      }
    },
    focus: () => button.focus(),
  };
}

/**
 * A summary text area: it wraps, stays one line of text (Enter submits, pasted line
 * breaks become spaces) and shows a counter that changes at 90% and at the limit.
 * props: { id, max, value, describedby, onSubmit() }.
 * Returns { root, area, value (get and set), focus() }.
 */
export function createSummary(props = {}) {
  const max = props.max || SUMMARY_MAX;
  const countId = props.id + "-count";
  const area = el("textarea", { id: props.id, rows: "2", maxlength: String(max), autocomplete: "off", "aria-describedby": [props.describedby, countId].filter(Boolean).join(" ") });
  const count = el("span", { id: countId, class: "counter" });
  const said = el("span", { class: "sr-only", "aria-live": "polite" }); // spoken only when the level changes
  let level = "";
  function update() {
    const s = counterState(area.value.length, max);
    count.textContent = s.text;
    count.dataset.level = s.level;
    if (s.level !== level) {
      said.textContent = s.level === "ok" ? "" : s.text;
      level = s.level;
    }
  }
  area.addEventListener("input", () => {
    const flat = oneLine(area.value);
    if (flat !== area.value) area.value = flat;
    update();
  });
  area.addEventListener("keydown", (e) => {
    if (e.key === "Enter" && !e.shiftKey && !e.metaKey && !e.ctrlKey && !e.isComposing) {
      e.preventDefault();
      if (props.onSubmit) props.onSubmit();
    }
  });
  if (props.value) area.value = oneLine(props.value);
  update();
  return {
    root: el("div", { class: "summary-ctl" }, area, count, said),
    area,
    get value() {
      return area.value;
    },
    set value(v) {
      area.value = oneLine(v);
      update();
    },
    focus: () => area.focus(),
  };
}
