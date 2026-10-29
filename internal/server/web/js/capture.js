// Capture: a new raw work item with a live preview of its name. The slug rules
// are those of pkm-workbench (filler words dropped, six words, 48 characters).
import { el } from "./util.js";
import { api, describeError } from "./store.js";

const MAX_WORDS = 6;
const MAX_TITLE_SLUG = 48;
export const MAX_SLUG = 60; // the API's limit
export const KINDS = ["feature", "bug", "chore", "spike"];

// Words left out of a slug made from a title, so the word budget goes to the
// words that carry meaning. They stay when nothing else is left.
const FILLER = new Set(["a", "an", "the", "to", "of", "for", "and", "or", "in", "on", "at", "with", "my"]);

// Accented Latin letters made plain; a letter not here counts as a separator.
const PLAIN = {};
for (const [plain, letters] of Object.entries({
  a: "àáâãäåāăą",
  ae: "æ",
  c: "çćč",
  d: "ďđð",
  e: "èéêëēėęě",
  i: "ìíîïīįı",
  l: "łľ",
  n: "ñńň",
  o: "òóôõöøōő",
  oe: "œ",
  r: "ř",
  s: "śšş",
  ss: "ß",
  t: "ť",
  th: "þ",
  u: "ùúûüūůű",
  y: "ýÿ",
  z: "źżž",
})) {
  for (const ch of letters) PLAIN[ch] = plain;
}

/** Lower-case ASCII words of a text, accents made plain. */
export function words(text) {
  const out = [];
  let cur = "";
  const flush = () => {
    if (cur) out.push(cur);
    cur = "";
  };
  for (const ch of String(text || "").toLowerCase()) {
    if ((ch >= "a" && ch <= "z") || (ch >= "0" && ch <= "9")) cur += ch;
    else if (PLAIN[ch]) cur += PLAIN[ch];
    else flush();
  }
  flush();
  return out;
}

function join(ws, maxWords, maxChars) {
  let s = ws.slice(0, maxWords).join("-");
  if (s.length > maxChars) s = s.slice(0, maxChars).replace(/-+$/, "");
  return s;
}

/** The slug a title makes: kebab-case, no filler, at most six words and 48 characters; "" if no letters or digits. */
export function slugify(title) {
  const all = words(String(title || "").replace(/\s+/g, " "));
  const keep = all.filter((w) => !FILLER.has(w));
  return join(keep.length ? keep : all, MAX_WORDS, MAX_TITLE_SLUG);
}

/** A slug the user typed, made valid without dropping filler: the wording is theirs. At most 8 words and 60 characters. */
export function cleanSlug(text) {
  return join(words(text), MAX_WORDS + 2, MAX_SLUG);
}

export const validSlug = (s) => /^[a-z0-9]+(-[a-z0-9]+)*$/.test(s) && s.length <= MAX_SLUG;

/** The line under the slug field: where the new file will be. */
export function fileNote(workDir, typed) {
  return typed ? "File: " + workDir + "/<id>-" + typed + ".md" : "";
}

/**
 * The `new` request for the capture form, or the first problem.
 * f: { title, slug (typed, "" if not edited), summary, kind, priority, parent };
 * openIds: ids a parent may be. Returns { body } or { error, field }.
 */
export function buildCapture(f, openIds = []) {
  const title = String(f.title || "").replace(/\s+/g, " ").trim();
  if (!title) return { error: "Write a title first.", field: "title" };
  const typed = String(f.slug || "").trim();
  const slug = typed ? cleanSlug(typed) : slugify(title);
  if (!slug || !validSlug(slug)) return { error: "The name has no letters or digits. Type a name of your own.", field: "slug" };
  const fields = {};
  const summary = String(f.summary || "").trim();
  if (summary.length > 160) return { error: "The summary is " + summary.length + " characters; the limit is 160.", field: "summary" };
  if (summary) fields.summary = summary;
  if (f.kind) {
    if (!KINDS.includes(f.kind)) return { error: "The kind is one of " + KINDS.join(", ") + ".", field: "kind" };
    fields.kind = f.kind;
  }
  const pr = String(f.priority === undefined || f.priority === null ? "" : f.priority).trim();
  if (pr) {
    if (!/^-?(0|[1-9][0-9]*)$/.test(pr)) return { error: "Priority is a whole number.", field: "priority" };
    fields.priority = Number(pr);
  }
  const parent = String(f.parent || "").trim();
  if (parent) {
    if (!openIds.includes(parent)) return { error: "No open work item has the id " + parent + ".", field: "parent" };
    fields.parent = parent;
  }
  const body = { kind: "work", slug, title };
  if (Object.keys(fields).length) body.fields = fields;
  return { body };
}

/**
 * The capture dialog. deps: { openDialog, closeDialog, setDialogHints, MOD }.
 * spec: { parents: [{ id, title }], onDone(record) }.
 */
export function openCapture(deps, spec) {
  const input = (id, props) => el("input", { id, type: "text", autocomplete: "off", ...props });
  const title = input("cap-title", { maxlength: "200" });
  const slug = input("cap-slug", { spellcheck: "false", "aria-describedby": "cap-note" });
  const note = el("p", { id: "cap-note", class: "help", "aria-live": "polite" });
  const summary = input("cap-summary", { maxlength: "160" });
  const kind = el("select", { id: "cap-kind" }, el("option", { value: "", text: "(default: chore)" }), KINDS.map((k) => el("option", { value: k, text: k })));
  const priority = input("cap-priority", { inputmode: "numeric", placeholder: "0", size: "6" });
  const parent = input("cap-parent", { list: "cap-parents", placeholder: "none" });
  const parents = el("datalist", { id: "cap-parents" }, spec.parents.map((p) => el("option", { value: p.id, label: p.title })));
  const err = el("p", { id: "cap-err", class: "err", role: "alert", hidden: true });
  const submit = el("button", { type: "submit", class: "btn primary" }, "Capture ", el("kbd", { text: "Enter" }));
  let dirty = false;
  let busy = false;

  const field = (label, help, control) => el("div", { class: "field" }, el("label", { class: "label", for: control.id, text: label }), help ? el("span", { class: "help", text: help }) : null, control);

  function preview() {
    const made = slugify(title.value);
    if (!dirty) slug.value = made;
    const typed = dirty ? cleanSlug(slug.value) : made;
    note.textContent = fileNote(spec.workDir || "docs/work", typed) + (!typed ? "" : dirty && typed !== slug.value.trim() ? " (cleaned from what you typed)" : "");
  }
  title.addEventListener("input", preview);
  slug.addEventListener("input", () => {
    dirty = slug.value.trim() !== "";
    preview();
  });

  const form = el(
    "form",
    { id: "cap-form", novalidate: true },
    el("header", { class: "dlg-head" }, el("div", {}, el("p", { class: "eyebrow", text: "New raw work item" }), el("h2", { id: "cap-t", text: "Capture" }))),
    el(
      "div",
      { class: "dlg-body" },
      field("What is it?", "A short title. You shape it later, when you move it out of raw.", title),
      el("div", { class: "field" }, el("label", { class: "label", for: "cap-slug", text: "Name" }), el("span", { class: "help", text: "Made from the title. Edit it to word it yourself." }), slug, note),
      field("Summary (optional)", "One line, at most 160 characters.", summary),
      el("div", { class: "row3" }, field("Kind", "", kind), field("Priority", "", priority), field("Parent", "", parent), parents),
      err,
    ),
    el("footer", { class: "dlg-foot" }, el("button", { type: "button", class: "btn quiet", onclick: () => deps.closeDialog() }, "Cancel ", el("kbd", { text: "Esc" })), submit),
  );
  const dlg = el("div", { class: "dialog capture", role: "dialog", "aria-modal": "true", "aria-labelledby": "cap-t" }, form);

  function show(e, field) {
    err.hidden = !e;
    err.textContent = e || "";
    const f = { title, slug, summary, kind, priority, parent }[field];
    if (f) f.focus();
  }

  async function send() {
    if (busy) return;
    const r = buildCapture({ title: title.value, slug: dirty ? slug.value : "", summary: summary.value, kind: kind.value, priority: priority.value, parent: parent.value }, spec.parents.map((p) => p.id));
    if (r.error) return show(r.error, r.field);
    busy = true;
    submit.disabled = true;
    try {
      const data = await api("POST", "records", { body: r.body });
      deps.closeDialog();
      spec.onDone(data.record);
    } catch (e) {
      busy = false;
      submit.disabled = false;
      show(describeError(e), "title");
    }
  }
  form.addEventListener("submit", (e) => {
    e.preventDefault();
    send();
  });

  deps.openDialog(dlg, {
    mode: "CAPTURE",
    hints: [["Enter", "capture"], ["Esc", "cancel"]],
    onKey(e) {
      if (e.key === "Escape") {
        e.preventDefault();
        deps.closeDialog();
      } else if (e.key === "Enter" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        send();
      }
    },
  });
  title.focus();
}
