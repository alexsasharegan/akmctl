// The board: four lanes of work items (raw, ready, active, blocked), each with
// its own scroll, driven by the keyboard. Focus is a card id that survives a
// reload and a visit to another view.
import { el, ageText, text } from "./util.js";
import { lanes, locate, firstLoc, step, field, priorityOf, blockedBy, staleness } from "./logic.js";

const FOCUS_KEY = "akm-board-focus";

const LANE_INFO = {
  raw: { glyph: "○", label: "Raw", desc: "Captured, not yet shaped" },
  ready: { glyph: "◐", label: "Ready", desc: "Shaped, can start" },
  active: { glyph: "●", label: "Active", desc: "In progress" },
  blocked: { glyph: "■", label: "Blocked", desc: "Waiting on something" },
};

let remembered = null; // the focused card id, kept across mounts
try {
  remembered = sessionStorage.getItem(FOCUS_KEY);
} catch (e) {
  remembered = null;
}

export function remember(id) {
  remembered = id;
  try {
    sessionStorage.setItem(FOCUS_KEY, id);
  } catch (e) {
    /* kept in memory only */
  }
}

/** The card for a work item. */
export function cardNode(r, byId, today) {
  const pri = priorityOf(r);
  const { days, stale } = staleness(r, today);
  const parent = field(r, "parent");
  const waits = blockedBy(r);
  const meta = [];
  if (parent) {
    const p = byId.get(parent);
    meta.push(el("span", { class: "chip", title: p ? text(field(p, "title")) : "not in this repository", text: "parent " + parent }));
  }
  if (waits.length) meta.push(el("span", { class: "chip wait", text: "waits on " + waits.join(", ") }));
  if (days !== null) meta.push(el("span", { class: "chip age" + (stale ? " stale" : ""), text: (stale ? "stale " : "") + ageText(days) }));
  const detail = field(r, "summary");
  return el(
    "a",
    { class: "card" + (stale ? " is-stale" : ""), draggable: "true", href: "#/record/" + encodeURIComponent(r.id), "data-id": r.id, tabindex: "-1" },
    el(
      "div",
      { class: "card-head" },
      el("span", { class: "card-id", text: r.id }),
      pri > 0 ? el("span", { class: "chip pri", title: "priority " + pri, text: "P" + pri }) : null,
      field(r, "kind") ? el("span", { class: "chip kind", text: field(r, "kind") }) : null,
    ),
    el("div", { class: "card-title", text: text(field(r, "title")) || "(no title)" }),
    meta.length ? el("div", { class: "card-meta" }, meta) : null,
    detail ? el("div", { class: "card-detail", text: detail }) : null,
  );
}

/**
 * Draw the board into `host`. ctx: records(), byId(), today(), setPage(p).
 * Returns { keys(e), focus(), destroy() }.
 */
export function mountBoard(host, ctx) {
  const ls = lanes(ctx.records());
  const byId = ctx.byId();
  const today = ctx.today();
  let loc = (remembered && locate(ls, remembered)) || firstLoc(ls);
  let pendingG = 0;

  const cols = ls.map((lane) => {
    const info = LANE_INFO[lane.status];
    const body = el("div", { class: "col-body", role: "list", "aria-labelledby": "lane-" + lane.status });
    if (!lane.records.length) body.append(el("div", { class: "empty-col", text: "Nothing in " + lane.status + "." }));
    for (const r of lane.records) {
      const c = cardNode(r, byId, today);
      c.setAttribute("role", "listitem");
      body.append(c);
    }
    return el(
      "section",
      { class: "col s-" + lane.status },
      el(
        "header",
        { class: "col-head" },
        el("h2", { id: "lane-" + lane.status }, el("span", { class: "g col-glyph", "aria-hidden": "true", text: info.glyph }), " ", info.label),
        el("span", { class: "count", "aria-label": lane.records.length + " items", text: "(" + lane.records.length + ")" }),
        el("p", { class: "col-desc", text: info.desc }),
      ),
      body,
    );
  });
  const boardEl = el("div", { class: "board" }, cols);
  host.replaceChildren(boardEl);

  const cardAt = (l) => l && host.querySelectorAll(".col")[l.lane].querySelectorAll(".card")[l.index];

  function describe() {
    if (!loc) return { mode: "BOARD", hints: HINTS, where: "no work items" };
    const lane = ls[loc.lane];
    return { mode: "BOARD", hints: HINTS, where: lane.status + " " + (loc.index + 1) + " of " + lane.records.length };
  }

  function focusCard(scroll) {
    for (const c of host.querySelectorAll(".card")) c.tabIndex = -1;
    const c = cardAt(loc);
    if (c) {
      c.tabIndex = 0;
      c.focus({ preventScroll: true });
      if (scroll) c.scrollIntoView({ block: "nearest", inline: "nearest" });
      remember(ls[loc.lane].records[loc.index].id);
    }
    ctx.setPage(describe());
  }

  // a click or Tab onto a card makes it the focused one
  boardEl.addEventListener("focusin", (e) => {
    const c = e.target.closest && e.target.closest(".card");
    if (!c) return;
    const at = locate(ls, c.dataset.id);
    if (at && (!loc || at.lane !== loc.lane || at.index !== loc.index)) {
      loc = at;
      for (const x of host.querySelectorAll(".card")) x.tabIndex = x === c ? 0 : -1;
      remember(c.dataset.id);
      ctx.setPage(describe());
    }
  });

  // make one card reachable by Tab even before the first key press
  const c0 = cardAt(loc);
  if (c0) c0.tabIndex = 0;
  ctx.setPage(describe());

  const KEYS = { h: "left", ArrowLeft: "left", l: "right", ArrowRight: "right", j: "down", ArrowDown: "down", k: "up", ArrowUp: "up", G: "last", Home: "first", End: "last" };

  return {
    keys(e) {
      if (e.shiftKey && e.key.length > 1 && e.key !== "Home" && e.key !== "End") return false;
      let dir = KEYS[e.key];
      if (e.key === "g") {
        const now = Date.now();
        if (now - pendingG < 900) {
          pendingG = 0;
          dir = "first";
        } else {
          pendingG = now;
          return true;
        }
      } else pendingG = 0;
      if (dir) {
        if (!loc && !(loc = firstLoc(ls))) return true;
        const next = step(ls, loc, dir);
        if (next) loc = next;
        focusCard(true);
        return true;
      }
      if (loc && (e.key === "o" || (e.key === "Enter" && !(e.target.closest && e.target.closest("a, button"))))) {
        ctx.open(ls[loc.lane].records[loc.index].id);
        return true;
      }
      return false;
    },
    focus() {
      focusCard(true);
    },
    /** The id of the focused card, or null. */
    current() {
      return loc ? ls[loc.lane].records[loc.index].id : null;
    },
    destroy() {
      host.replaceChildren();
    },
  };
}

const HINTS = [
  ["h j k l", "move"],
  ["g g", "first"],
  ["G", "last"],
  ["Enter", "open"],
  ["H L", "move card"],
  ["d x b", "done drop block"],
  ["c", "capture"],
  ["u", "undo"],
  ["/", "search"],
  ["?", "help"],
];
