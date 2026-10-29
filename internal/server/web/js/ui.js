// The chrome shared by every view: theme, hint bar, toasts, the command palette
// (Cmd/Ctrl+K and /) and the help overlay (?). All text goes in as text nodes.
import { el, clear, debounce, text } from "./util.js";
import { score, searchRecords, field } from "./logic.js";
import { searchBody, describeError } from "./store.js";

export const isMac = /Mac|iPhone|iPad/.test(navigator.platform || navigator.userAgent || "");
export const MOD = isMac ? "⌘" : "Ctrl";

const root = document.documentElement;
const $ = (s) => document.querySelector(s);

const S = {
  overlay: null, // "palette" | "help" | null
  returnFocus: null,
  page: { mode: "", hints: [], where: "" },
  hooks: { records: () => [], actions: () => [], onKey: () => false, help: () => [], openRecord: () => {} },
  pal: { q: "", sel: 0, items: [], body: [], bodyQ: "", seq: 0 },
  modal: { mode: "", hints: [], onKey: null }, // a dialog opened by openDialog (the wizard, capture)
};

/** A key such as "g g" or "Ctrl+K" as <kbd> elements. */
export function kbd(keys) {
  return String(keys)
    .split(" ")
    .filter(Boolean)
    .map((k) => el("kbd", { text: k }));
}

export const typing = (t) => !!t && (t.tagName === "TEXTAREA" || t.tagName === "SELECT" || t.isContentEditable || (t.tagName === "INPUT" && !["checkbox", "radio", "button"].includes(t.type)));

// ---- theme ----

export function currentTheme() {
  return root.dataset.theme === "dark" ? "dark" : "light";
}

export function setTheme(theme) {
  root.dataset.theme = theme;
  try {
    localStorage.setItem("akm-theme", theme);
  } catch (e) {
    /* the choice lasts until the page is closed */
  }
  renderChrome();
}

export const toggleTheme = () => setTheme(currentTheme() === "dark" ? "light" : "dark");

function renderChrome() {
  const dark = currentTheme() === "dark";
  const tb = $("#theme-btn");
  clear(tb).append(el("span", { class: "g", "aria-hidden": "true", text: dark ? "☾" : "☀" }), el("span", { class: "theme-label", text: dark ? "Dark" : "Light" }), el("kbd", { text: "t" }));
  tb.setAttribute("aria-label", "Switch to the " + (dark ? "light" : "dark") + " theme");
  const pb = $("#palette-btn");
  clear(pb).append(el("span", { class: "g", "aria-hidden": "true", text: "⌕" }), el("span", { class: "pb-label", text: "Search" }), el("kbd", { text: MOD }), el("kbd", { text: "K" }));
}

// ---- hint bar ----

function hintState() {
  if (S.overlay === "palette") return ["PALETTE", [["↑ ↓", "select"], ["Enter", "open"], ["Esc", "close"]], ""];
  if (S.overlay === "help") return ["HELP", [["Esc", "close"], ["?", "close"]], ""];
  if (S.overlay === "modal") return [S.modal.mode, S.modal.hints, ""];
  return [S.page.mode, S.page.hints, S.page.where];
}

export function renderHints() {
  const [mode, hints, where] = hintState();
  clear($("#hintbar")).append(
    el("span", { class: "mode", text: mode }),
    el("ul", { class: "hints" }, hints.map(([k, label]) => el("li", {}, kbd(k), el("span", { text: label })))),
    el("span", { class: "where", text: where || "" }),
  );
}

/** Describe the current view for the hint bar: { mode, hints: [[keys, label]], where }. */
export function setPage(p) {
  Object.assign(S.page, p);
  renderHints();
}

// ---- toast ----

export function toast(msg, kind, ms) {
  const bad = kind === "bad";
  const box = el("div", { class: "toast" + (bad ? " bad" : "") }, bad ? el("strong", { text: "Error: " }) : null, msg);
  const host = $("#toasts");
  host.append(box);
  while (host.children.length > 3) host.firstChild.remove();
  setTimeout(() => box.remove(), ms || (bad ? 9000 : 3500));
  return box;
}

// ---- overlays ----

const setInert = (on) => {
  for (const sel of [".app", "#drawer"]) {
    const n = $(sel);
    if (n) n.inert = on;
  }
};

export const overlayOpen = () => S.overlay !== null;

function openOverlay(kind, dialog) {
  S.returnFocus = document.activeElement;
  S.overlay = kind;
  const scrim = el("div", { class: "scrim", onclick: closeOverlay });
  clear($("#overlay-root")).append(scrim, dialog);
  setInert(true);
  renderHints();
}

/**
 * Show a dialog node that the caller built (the move wizard, capture). hooks:
 * { mode, hints, onKey(e) } - every key goes to onKey while it is open.
 */
export function openDialog(dialog, hooks) {
  S.modal = { mode: "", hints: [], onKey: null, ...hooks };
  clear($("#toasts")); // a toast would sit over the dialog's buttons
  openOverlay("modal", dialog);
}

/** Change the mode and hints the hint bar shows for the open dialog. */
export function setDialogHints(mode, hints) {
  S.modal.mode = mode;
  S.modal.hints = hints;
  renderHints();
}

export function closeOverlay() {
  if (!S.overlay) return;
  S.overlay = null;
  S.modal = { mode: "", hints: [], onKey: null };
  clear($("#overlay-root"));
  setInert(false);
  const back = S.returnFocus;
  S.returnFocus = null;
  if (back && back.isConnected && back.focus) back.focus({ preventScroll: true });
  renderHints();
}

// ---- help ----

function helpSections() {
  return [
    [
      "Anywhere",
      [
        [MOD + "+K", "command palette"],
        ["/", "command palette"],
        ["1", "Board"],
        ["2", "Tables"],
        ["3", "ADRs"],
        ["4", "Check"],
        ["r", "reload the records"],
        ["t", "light or dark"],
        ["?", "this list"],
        ["Esc", "close, or go back"],
      ],
    ],
    [
      "Palette",
      [
        ["type", "filter actions and records"],
        ["↑ ↓", "select"],
        ["Enter", "open"],
        ["Esc", "close"],
      ],
    ],
  ].concat(S.hooks.help());
}

export function openHelp() {
  const close = el("button", { type: "button", class: "btn quiet", onclick: closeOverlay }, "Close ", el("kbd", { text: "Esc" }));
  const dlg = el(
    "div",
    { class: "dialog help", role: "dialog", "aria-modal": "true", "aria-labelledby": "help-t" },
    el("header", { class: "dlg-head" }, el("div", {}, el("p", { class: "eyebrow", text: "Keyboard" }), el("h2", { id: "help-t", text: "Shortcuts" })), close),
    el(
      "div",
      { class: "dlg-body help-grid" },
      helpSections().map(([title, rows]) =>
        el("section", {}, el("h3", { text: title }), el("dl", {}, rows.map(([k, label]) => el("div", {}, el("dt", {}, kbd(k)), el("dd", { text: label }))))),
      ),
    ),
  );
  openOverlay("help", dlg);
  close.focus();
}

// ---- palette ----

/** The palette rows for a query: actions, then loaded records, then body hits that are not already listed. */
export function paletteItems(q, actions, records, bodyHits) {
  const acts = actions.map((a) => ({ kind: "action", label: a.label, keys: a.keys || [], run: a.run, score: score(q, a.label) })).filter((a) => a.score);
  if (!q.trim()) return acts;
  acts.sort((a, b) => b.score - a.score);
  const local = searchRecords(records, q, 8).map((x) => ({ kind: "record", record: x.record }));
  const seen = new Set(local.map((x) => x.record.id));
  const body = (bodyHits || []).filter((r) => !seen.has(r.id)).map((r) => ({ kind: "body", record: r }));
  return acts.concat(local, body);
}

function palMarks(label, q) {
  const w = q.trim().toLowerCase().split(/\s+/)[0];
  const i = w ? label.toLowerCase().indexOf(w) : -1;
  if (i < 0) return [label];
  return [label.slice(0, i), el("mark", { text: label.slice(i, i + w.length) }), label.slice(i + w.length)];
}

const GROUPS = { action: "Actions", record: "Records", body: "Text in bodies" };

function palRender() {
  const q = S.pal.q;
  const items = paletteItems(q, S.hooks.actions(), S.hooks.records(), S.pal.bodyQ === q.trim() ? S.pal.body : []);
  S.pal.items = items;
  S.pal.sel = Math.max(0, Math.min(S.pal.sel, items.length - 1));
  const list = $("#pal-list");
  clear(list);
  if (!items.length) {
    list.append(el("li", { class: "pal-empty", role: "presentation" }, "No match for ", el("b", { text: q.trim() }), ". Try a shorter word."));
    $("#pal-q").removeAttribute("aria-activedescendant");
    return;
  }
  let last = "";
  items.forEach((it, i) => {
    if (it.kind !== last) {
      list.append(el("li", { class: "pal-group", role: "presentation", text: GROUPS[it.kind] }));
      last = it.kind;
    }
    const sel = i === S.pal.sel;
    let label;
    let right;
    if (it.kind === "action") {
      label = el("span", { class: "pal-label" }, palMarks(it.label, q));
      right = it.keys.map((k) => kbd(k));
    } else {
      const r = it.record;
      const sub = it.kind === "body" ? bodySnippet(r) : "";
      label = el(
        "span",
        { class: "pal-label" },
        el("span", { class: "pal-id", text: r.id }),
        " ",
        palMarks(text(field(r, "title")) || "(no title)", q),
        sub ? el("span", { class: "pal-sub", text: sub }) : null,
      );
      right = [el("span", { class: "stat", text: text(field(r, "status")) })];
    }
    list.append(
      el(
        "li",
        {
          id: "pal-" + i,
          class: "pal-row" + (sel ? " sel" : ""),
          role: "option",
          "aria-selected": sel ? "true" : "false",
          onclick: () => palRun(i),
        },
        el("span", { class: "pal-mark", "aria-hidden": "true", text: "▸" }),
        label,
        el("span", { class: "pal-right" }, right),
      ),
    );
  });
  const row = $("#pal-" + S.pal.sel);
  if (row) {
    row.scrollIntoView({ block: "nearest" });
    $("#pal-q").setAttribute("aria-activedescendant", row.id);
  }
}

function bodySnippet(r) {
  const m = (r.matches || []).find((x) => x.field === "body") || (r.matches || [])[0];
  return m ? "line " + m.line + ": " + m.snippet : "";
}

const bodyLookup = debounce(async (q, seq) => {
  try {
    const data = await searchBody(q);
    if (seq !== S.pal.seq || !S.overlay) return;
    S.pal.body = data.records;
    S.pal.bodyQ = q;
    palRender();
  } catch (e) {
    if (seq === S.pal.seq && S.overlay === "palette") toast("Body search failed: " + describeError(e), "bad");
  }
}, 250);

function palRun(i) {
  const it = S.pal.items[i];
  if (!it) return;
  S.returnFocus = null;
  closeOverlay();
  if (it.kind === "action") it.run();
  else S.hooks.openRecord(it.record.id);
}

export function openPalette() {
  S.pal = { q: "", sel: 0, items: [], body: [], bodyQ: "", seq: S.pal.seq + 1 };
  const input = el("input", {
    id: "pal-q",
    type: "text",
    autocomplete: "off",
    spellcheck: "false",
    placeholder: "Type a command, a title, an id or a word from a record",
    role: "combobox",
    "aria-expanded": "true",
    "aria-controls": "pal-list",
    "aria-autocomplete": "list",
    "aria-label": "Command palette",
  });
  input.addEventListener("input", () => {
    S.pal.q = input.value;
    S.pal.sel = 0;
    S.pal.seq++;
    bodyLookup.cancel();
    const q = input.value.trim();
    if (q.length >= 3) bodyLookup(q, S.pal.seq);
    palRender();
  });
  const dlg = el(
    "div",
    { class: "dialog palette", role: "dialog", "aria-modal": "true", "aria-label": "Command palette" },
    el("div", { class: "pal-input" }, el("span", { class: "g", "aria-hidden": "true", text: "⌕" }), input, el("kbd", { text: "Esc" })),
    el("ul", { id: "pal-list", class: "pal-list", role: "listbox", "aria-label": "Results" }),
    el("div", { class: "pal-foot" }, el("span", {}, kbd("↑ ↓"), " select"), el("span", {}, kbd("Enter"), " open"), el("span", {}, kbd("Esc"), " close")),
  );
  openOverlay("palette", dlg);
  palRender();
  input.focus();
}

// ---- keys ----

function onKeyDown(e) {
  const mod = e.metaKey || e.ctrlKey;
  if (mod && !e.altKey && e.key.toLowerCase() === "k") {
    e.preventDefault();
    if (S.overlay === "palette") closeOverlay();
    else if (!S.overlay) openPalette();
    return;
  }
  if (S.overlay === "modal") {
    if (S.modal.onKey) S.modal.onKey(e);
    return;
  }
  if (S.overlay === "palette") {
    const n = S.pal.items.length;
    if (e.key === "Escape") {
      e.preventDefault();
      closeOverlay();
    } else if (e.key === "ArrowDown" || (e.ctrlKey && e.key === "n")) {
      e.preventDefault();
      S.pal.sel = Math.min(n - 1, S.pal.sel + 1);
      palRender();
    } else if (e.key === "ArrowUp" || (e.ctrlKey && e.key === "p")) {
      e.preventDefault();
      S.pal.sel = Math.max(0, S.pal.sel - 1);
      palRender();
    } else if (e.key === "Enter") {
      e.preventDefault();
      palRun(S.pal.sel);
    } else if (e.key === "Tab") e.preventDefault();
    return;
  }
  if (S.overlay === "help") {
    if (e.key === "Escape" || e.key === "?") {
      e.preventDefault();
      closeOverlay();
    }
    return;
  }
  if (typing(e.target) || mod || e.altKey) return;
  if (e.key === "?") openHelp();
  else if (e.key === "/") openPalette();
  else if (e.key === "t") toggleTheme();
  else if (S.hooks.onKey(e)) {
    /* handled by the view or the router */
  } else return;
  e.preventDefault();
}

/**
 * Start the chrome. hooks: records() the loaded records, actions() the palette
 * actions [{label, keys, run}], onKey(e) a key the views handle (true if used),
 * help() extra help sections, openRecord(id).
 */
export function initUI(hooks) {
  Object.assign(S.hooks, hooks);
  document.addEventListener("keydown", onKeyDown);
  $("#theme-btn").addEventListener("click", toggleTheme);
  $("#palette-btn").addEventListener("click", () => {
    if (!S.overlay) openPalette();
  });
  matchMedia("(prefers-color-scheme: dark)").addEventListener("change", (e) => {
    let saved = "";
    try {
      saved = localStorage.getItem("akm-theme") || "";
    } catch (err) {
      saved = "";
    }
    if (!saved) {
      root.dataset.theme = e.matches ? "dark" : "light";
      renderChrome();
    }
  });
  renderChrome();
  renderHints();
}
