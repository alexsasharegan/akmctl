// The parent picker: a combobox that searches the open work items and shows what
// matched. Stages 1 and 2 of work item 261029-bdt9: a results list wide enough to read
// with the id, title and summary matches marked, and rows that show id, title, status
// and the summary. The search and the highlighting are pure.
import { el } from "./util.js";
import { placeList, keepInView } from "./controls.js";

export const MAX_RESULTS = 50;

/** The query as lower-case tokens; every token must match for a record to be found. */
export const tokens = (q) => String(q || "").toLowerCase().split(/\s+/).filter(Boolean);

/** The [start, end) ranges of `text` that any token matches (case-insensitive), merged and sorted. */
export function matchRanges(text, toks) {
  const low = String(text || "").toLowerCase();
  const out = [];
  for (const t of toks) {
    for (let at = low.indexOf(t); at >= 0; at = low.indexOf(t, at + t.length)) out.push([at, at + t.length]);
  }
  out.sort((a, b) => a[0] - b[0] || a[1] - b[1]);
  const merged = [];
  for (const r of out) {
    const last = merged[merged.length - 1];
    if (last && r[0] <= last[1]) last[1] = Math.max(last[1], r[1]);
    else merged.push([r[0], r[1]]);
  }
  return merged;
}

/** `text` cut into [{ text, hit }] at the ranges, so a renderer never builds HTML from it. */
export function highlight(text, ranges) {
  const s = String(text || "");
  const parts = [];
  let at = 0;
  for (const [a, b] of ranges) {
    if (a > at) parts.push({ text: s.slice(at, a), hit: false });
    parts.push({ text: s.slice(a, b), hit: true });
    at = b;
  }
  if (at < s.length) parts.push({ text: s.slice(at), hit: false });
  return parts;
}

/**
 * The records a query finds. items: [{ id, title, summary }]. A record is found when
 * every token matches its id, title or summary. Each result is
 * { item, id, title, summary: ranges, where: ["id" | "title" | "summary", ...] }.
 * Order: an id match first (a prefix before a later one), then a title match (from the
 * start of a word first), then a summary-only match; ties keep the given order.
 * No query lists the first MAX_RESULTS items as given.
 */
export function searchItems(items, query, limit = MAX_RESULTS) {
  const toks = tokens(query);
  const found = [];
  items.forEach((item, order) => {
    const id = matchRanges(item.id, toks);
    const title = matchRanges(item.title, toks);
    const summary = matchRanges(item.summary, toks);
    if (toks.length) {
      const all = toks.every((t) => [item.id, item.title, item.summary].some((f) => String(f || "").toLowerCase().includes(t)));
      if (!all) return;
    }
    const where = [id.length && "id", title.length && "title", summary.length && "summary"].filter(Boolean);
    let rank = 4;
    if (!toks.length) rank = 0;
    else if (id.length) rank = id[0][0] === 0 ? 0 : 1;
    else if (title.length) rank = wordStart(item.title, title[0][0]) ? 2 : 3;
    found.push({ rank, order, r: { item, id: id, title: title, summary: summary, where } });
  });
  found.sort((a, b) => a.rank - b.rank || a.order - b.order);
  return found.slice(0, limit).map((f) => f.r);
}

/** True when `at` begins a word of `text`. */
const wordStart = (text, at) => at === 0 || !/[a-z0-9]/i.test(text[at - 1]);

/** What the picker says about the number of results, for a screen reader. */
export const countText = (n) => (n === 0 ? "No open work item matches" : n === 1 ? "1 result" : n + " results");

/** Text with the matched parts in <mark>, built from nodes. */
const marked = (text, ranges) => highlight(text, ranges).map((p) => (p.hit ? el("mark", { text: p.text }) : p.text));

let uid = 0;

/**
 * The picker. props: { id, items: [{ id, title, status, summary }], placeholder }.
 * Returns { root, input, value (the text in the box: an id once one is chosen) }.
 * Only the items given are offered; typing a free id that is not among them is the
 * caller's validation to refuse.
 */
export function createPicker(props) {
  const n = ++uid;
  const items = props.items;
  let results = [];
  let active = -1;
  let open = false;
  const list = el("ul", { class: "pick-list", role: "listbox", id: "pick-list-" + n, "aria-label": "Open work items", hidden: true });
  const input = el("input", {
    id: props.id,
    type: "text",
    class: "pick-input",
    role: "combobox",
    autocomplete: "off",
    spellcheck: "false",
    placeholder: props.placeholder || "none",
    "aria-autocomplete": "list",
    "aria-expanded": "false",
    "aria-controls": list.id,
  });
  const said = el("span", { class: "sr-only", "aria-live": "polite" });
  const root = el("div", { class: "pick" }, input, list, said);

  function render() {
    list.replaceChildren(
      ...results.map((r, i) =>
        el(
          "li",
          {
            class: "pick-opt" + (i === active ? " active" : ""),
            role: "option",
            id: "pick-opt-" + n + "-" + i,
            "aria-selected": String(i === active),
            onmousedown: (e) => e.preventDefault(), // keep the focus in the box
            onclick: () => choose(i),
            onmousemove: () => {
              if (active !== i) {
                active = i;
                mark();
              }
            },
          },
          el("span", { class: "pick-id" }, ...marked(r.item.id, r.id)),
          el("span", { class: "pick-title" }, ...marked(r.item.title, r.title)),
          r.item.status ? el("span", { class: "pick-status", "data-status": r.item.status, text: r.item.status }) : null,
          // the first line of the summary, with what matched marked, so a choice needs no opening
          r.item.summary ? el("span", { class: "pick-sum" }, ...marked(r.item.summary, r.summary)) : null,
        ),
      ),
    );
    if (!results.length) list.append(el("li", { class: "pick-none", role: "presentation", text: countText(0) }));
  }
  function mark() {
    [...list.querySelectorAll(".pick-opt")].forEach((li, i) => {
      li.classList.toggle("active", i === active);
      li.setAttribute("aria-selected", String(i === active));
    });
    const li = list.children[active];
    if (open && li && active >= 0) {
      input.setAttribute("aria-activedescendant", li.id);
      keepInView(list, li);
    } else input.removeAttribute("aria-activedescendant");
  }
  function show(next) {
    open = next;
    input.setAttribute("aria-expanded", String(open));
    list.hidden = !open;
    root.classList.toggle("open", open);
    if (open) placeList(input, list, root);
    mark();
  }
  function search(keepActive = false) {
    results = searchItems(items, input.value);
    if (!keepActive) active = results.length ? 0 : -1;
    render();
    said.textContent = countText(results.length);
    show(true);
  }
  function choose(i) {
    input.value = results[i].item.id;
    show(false);
    input.focus();
  }

  input.addEventListener("input", () => search());
  input.addEventListener("keydown", (e) => {
    if (e.metaKey || e.ctrlKey || e.altKey) return;
    const used = () => {
      e.preventDefault();
      e.stopPropagation();
    };
    if (e.key === "ArrowDown" || e.key === "ArrowUp") {
      used();
      if (!open) return search();
      active = results.length ? (active + (e.key === "ArrowDown" ? 1 : -1) + results.length) % results.length : -1;
      mark();
    } else if (e.key === "Enter" && open && active >= 0) {
      used();
      choose(active);
    } else if (e.key === "Escape" && open) {
      used();
      show(false);
    } else if (e.key === "Tab") show(false);
  });
  root.addEventListener("focusout", (e) => {
    if (open && !root.contains(e.relatedTarget)) show(false);
  });

  return {
    root,
    input,
    get value() {
      return input.value;
    },
    set value(v) {
      input.value = v;
    },
    focus: () => input.focus(),
  };
}
