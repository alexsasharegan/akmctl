// The record drawer's field table: one row per frontmatter field, typed and labelled,
// with the editable fields of an open work item edited in place through `set`.
// The typing, the edit rules and the request are pure (node tests); the table at the
// bottom wires them to the DOM. The edit rules are a fixed copy of the matrix in the
// Outcome of 261031-z0n8 until the tool reports them (261101-3che).
import { el } from "./util.js";
import { api, describeError } from "./store.js";
import { createSelect, createSummary, SUMMARY_MAX } from "./controls.js";
import { createPicker } from "./picker.js";
import { KINDS } from "./capture.js";

/** What each type is called and drawn with. The word is always shown, so colour is never the only signal. */
export const TYPE_LEGEND = [
  { type: "text", word: "text", note: "plain text" },
  { type: "date", word: "date", note: "a day, in monospace" },
  { type: "number", word: "number", note: "a whole number, in monospace" },
  { type: "flag", word: "flag", note: "yes or no" },
  { type: "list", word: "list", note: "one item per line" },
  { type: "link", word: "link", note: "a wikilink, opens its note" },
  { type: "id", word: "id", note: "a record id, opens the record when it is here" },
];

const DATE = /^\d{4}-\d{2}-\d{2}(T\d{2}:\d{2})?$/;
const ID = /^\d{6}-[0-9a-hjkmnp-tv-z]{4}$/;

/** The type of a parsed value. `wikilink(s)` tells whether a string is a wikilink. */
export function valueType(v, wikilink = (s) => /^\[\[.+\]\]$/.test(s)) {
  if (Array.isArray(v)) return "list";
  if (typeof v === "boolean") return "flag";
  if (typeof v === "number") return "number";
  const s = String(v);
  if (wikilink(s)) return "link";
  if (DATE.test(s)) return "date";
  if (ID.test(s)) return "id";
  return "text";
}

export const OPEN = ["raw", "ready", "active", "blocked"];
// the controls the drawer has; every other field is shown and locked
export const EDITORS = { summary: "text", kind: "select", priority: "number", parent: "picker", blocked_by: "picker-multi", cycle: "label", keep: "toggle" };
// what an unset editable field will hold, for its badge
const UNSET_TYPE = { toggle: "flag", picker: "id", "picker-multi": "list", number: "number" };
const MANAGED = {
  id: "the id never changes",
  created: "set by the tool",
  updated: "set by the tool on every write",
  closed: "set by close",
  status: "move the item on the board",
  claimed_by: "set by claim and release",
  claimed_at: "set by claim and release",
  title: "edit the title in the file",
};

/**
 * Whether a field can be edited here: { state: "yes" | "warn" | "no", why }.
 * `type` is the record type and `status` its status. A `warn` edit is allowed and
 * shows `why` next to the control.
 */
export function editState(name, type, status) {
  if (type !== "work") return { state: "no", why: "ADRs are read-only here" };
  if (MANAGED[name]) return { state: "no", why: MANAGED[name] };
  if (!EDITORS[name]) return { state: "no", why: "no control for this field yet" };
  if (!OPEN.includes(status)) return name === "keep" ? { state: "yes", why: "" } : { state: "no", why: "a closed item is read-only; a correction comes with 261101-3vwz" };
  const under = status === "active" || status === "blocked";
  if ((name === "kind" || name === "parent") && under) return { state: "warn", why: "changing " + name + " after work has started rewrites what the item was" };
  if (name === "priority" && status === "active") return { state: "warn", why: "changing priority while it is active rewrites what it was ranked as" };
  return { state: "yes", why: "" };
}

/**
 * The `set` request for an edit, or the first problem.
 * `raw` is what the control holds: a string for text, number, label, kind and parent;
 * a list of ids for blocked_by; a boolean for keep. An empty value removes the field.
 * ctx: { ids: the ids a parent or blocker may be, self }.
 * Returns { changes } or { error }.
 */
export function buildChanges(name, raw, ctx = { ids: [], self: "" }) {
  const one = (v) => ({ changes: { [name]: v } });
  switch (name) {
    case "summary": {
      const s = String(raw).trim();
      if (!s) return one(null);
      if (s.length > SUMMARY_MAX) return { error: "The summary is " + s.length + " characters; the limit is " + SUMMARY_MAX + "." };
      return one(s);
    }
    case "kind":
      return KINDS.includes(raw) ? one(raw) : { error: "The kind is one of " + KINDS.join(", ") + "." };
    case "priority": {
      const s = String(raw).trim();
      if (!s) return one(null);
      return /^-?(0|[1-9][0-9]*)$/.test(s) ? one(Number(s)) : { error: "Priority is a whole number." };
    }
    case "parent": {
      const s = String(raw).trim();
      if (!s) return one(null);
      if (s === ctx.self) return { error: "An item cannot be its own parent." };
      return ctx.ids.includes(s) ? one(s) : { error: "No open work item has the id " + s + "." };
    }
    case "blocked_by": {
      const list = [...raw];
      return list.length ? one(list) : one(null);
    }
    case "cycle": {
      const s = String(raw).trim();
      if (!s) return one(null);
      return /^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$/.test(s) ? one(s) : { error: "A cycle is a label of 1 to 64 characters: letters, digits, '.', '_' or '-', starting with a letter or digit." };
    }
    case "keep":
      return one(Boolean(raw));
  }
  return { error: name + " cannot be edited here." };
}

/** What a stale write says: nothing was written. */
export const STALE_TEXT = "This record changed since you opened it, so nothing was written. Reload to see the current values, then edit again.";

/**
 * The rows to show for a record: its present fields in file order (skipping the title,
 * which is the heading, and empty values), then the editable fields it does not have
 * yet when the item is open.
 */
export function rowNames(r) {
  const names = Object.entries(r.fields)
    .filter(([k, v]) => k !== "title" && v !== null && v !== undefined && !(Array.isArray(v) && !v.length))
    .map(([k]) => k);
  if (r.type === "work") {
    for (const k of Object.keys(EDITORS)) {
      if (!names.includes(k) && editState(k, r.type, String(r.fields.status || "")).state !== "no") names.push(k);
    }
  }
  return names;
}

// ---- the table ----

/**
 * The field table of one record. deps: { valueNodes(v, byId), byId, records() (the open
 * work items, for the pickers), reload() (draw the drawer again) }.
 */
export function fieldTable(r, deps) {
  const status = String(r.fields.status || "");
  const open = r.type === "work" && OPEN.includes(status);
  const legend = el(
    "details",
    { class: "fm-legend" },
    el("summary", { text: "Value types" }),
    el("ul", {}, TYPE_LEGEND.map((t) => el("li", {}, el("span", { class: "fm-type", "data-type": t.type, text: t.word }), " " + t.note))),
  );
  const rows = rowNames(r).map((k) => row(r, k, status, deps));
  const table = el(
    "table",
    { class: "fm", "aria-label": "Fields" },
    el("thead", { class: "sr-only" }, el("tr", {}, el("th", { text: "Field" }), el("th", { text: "Value" }), el("th", { text: "Edit" }))),
    el("tbody", {}, rows),
  );
  return el("div", { class: "fm-wrap" }, open ? null : r.type === "work" ? el("p", { class: "help", text: "A closed item is read-only except for keep." }) : null, legend, table);
}

function valueCell(v, deps) {
  if (v === null || v === undefined || (Array.isArray(v) && !v.length)) return el("span", { class: "fm-empty", text: "not set" });
  if (Array.isArray(v)) return el("ul", { class: "fm-list" }, v.map((x) => el("li", {}, deps.valueNodes(x, deps.byId))));
  if (typeof v === "boolean") return el("span", { class: "fm-flag", text: v ? "yes" : "no" });
  return el("span", { class: "fm-val", "data-type": valueType(v) }, deps.valueNodes(v, deps.byId));
}

function row(r, name, status, deps) {
  const v = r.fields[name];
  const es = editState(name, r.type, status);
  const type = v === null || v === undefined || (Array.isArray(v) && !v.length) ? UNSET_TYPE[EDITORS[name]] || "text" : valueType(v);
  const cell = el("td", { class: "fm-value" });
  const act = el("td", { class: "fm-act" });
  const show = () => {
    cell.replaceChildren(valueCell(v, deps));
    act.replaceChildren();
    if (es.state === "no") {
      act.append(el("span", { class: "fm-locked", title: es.why, text: "locked" }), el("span", { class: "sr-only", text: ": " + es.why }));
    } else if (EDITORS[name] === "toggle") {
      act.append(el("button", { type: "button", class: "btn quiet", "aria-pressed": String(Boolean(v)), "aria-label": name + (v ? " is on; turn off" : " is off; turn on"), onclick: () => save(name, !v) }, v ? "Turn off" : "Turn on"));
    } else {
      act.append(el("button", { type: "button", class: "btn quiet", "aria-label": "Edit " + name, onclick: edit }, "Edit"));
    }
  };
  const message = (text, cls) => {
    cell.querySelector(".fm-msg")?.remove();
    if (text) cell.append(el("p", { class: "fm-msg " + cls, role: cls === "err" ? "alert" : null, text }));
  };
  async function save(field, raw) {
    const ctx = { ids: deps.records().map((x) => x.id), self: r.id };
    const built = buildChanges(field, raw, ctx);
    if (built.error) return message(built.error, "err");
    message("", "");
    try {
      await api("PATCH", "records/" + encodeURIComponent(r.id), { body: { changes: built.changes }, etag: r.etag });
      deps.reload();
    } catch (e) {
      if (e && e.code === "stale") {
        message(STALE_TEXT, "err");
        cell.append(el("button", { type: "button", class: "btn", onclick: () => deps.reload() }, "Reload"));
      } else message(describeError(e), "err");
    }
  }
  function edit() {
    const control = makeControl(name, v, r, deps, () => save(name, control.get()));
    cell.replaceChildren(control.node);
    if (es.state === "warn") cell.append(el("p", { class: "fm-msg warn", text: "Careful: " + es.why + "." }));
    act.replaceChildren(
      el("button", { type: "button", class: "btn primary", onclick: () => save(name, control.get()) }, "Save"),
      el("button", { type: "button", class: "btn quiet", onclick: show }, "Cancel"),
    );
    control.focus();
  }
  show();
  return el("tr", { class: "fm-row", "data-field": name }, el("th", { scope: "row", class: "fm-name" }, el("span", { class: "mono", text: name }), el("span", { class: "fm-type", "data-type": type, text: type })), cell, act);
}

/** The control for a field: { node, get(), focus() }. onSubmit saves (Enter in a text box). */
function makeControl(name, v, r, deps, onSubmit) {
  const id = "fm-" + name;
  const label = el("label", { class: "sr-only", for: id, text: name });
  const wrap = (c, get) => ({ node: el("div", { class: "fm-edit" }, label, c.root || c), get, focus: () => c.focus() });
  switch (EDITORS[name]) {
    case "text": {
      const c = createSummary({ id, value: v === null || v === undefined ? "" : String(v), onSubmit });
      return wrap(c, () => c.value);
    }
    case "select": {
      const c = createSelect({ id, value: v || KINDS[2], options: KINDS.map((k) => ({ value: k, label: k })) });
      return wrap(c, () => c.value);
    }
    case "picker": {
      const items = pickItems(deps, r.id);
      const c = createPicker({ id, items, self: r.id, placeholder: "none" });
      c.value = v ? String(v) : "";
      return wrap(c, () => c.value);
    }
    case "picker-multi": {
      const c = createPicker({ id, items: pickItems(deps, r.id), multiple: true, self: r.id, values: Array.isArray(v) ? v.map(String) : [], placeholder: "add" });
      return wrap(c, () => c.values);
    }
    default: {
      const input = el("input", { id, type: "text", autocomplete: "off", spellcheck: "false", inputmode: name === "priority" ? "numeric" : null, size: name === "priority" ? "6" : null, value: v === null || v === undefined ? "" : String(v) });
      input.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
          e.preventDefault();
          onSubmit();
        }
      });
      return wrap({ root: input, focus: () => input.focus() }, () => input.value);
    }
  }
}

function pickItems(deps, self) {
  return deps
    .records()
    .filter((x) => x.id !== self)
    .map((x) => ({ id: x.id, title: x.title, status: x.status, summary: x.summary }));
}
