// The workbench: reads the token from the URL fragment, loads the records once
// and draws the view that the hash route names. Everything from the repository
// is shown as text, never as HTML.
import { initToken, parseToken, hasToken, loadSnapshot, loadConfig, loadClosed, check, state, describeError, ApiError } from "./js/store.js";
import { createLive, mountIndicator } from "./js/live.js";
import { createAnimator, keepSteady } from "./js/anim.js";
import { parseRoute, routeHref, indexById, TABLES, pageTitle, repoLabel } from "./js/logic.js";
import { el, isoDate } from "./js/util.js";
import { initUI, setPage, toast, toggleTheme, openHelp, overlayOpen, typing } from "./js/ui.js";
import { mountBoard } from "./js/board.js";
import { mountTables } from "./js/tables.js";
import { mountAdrs } from "./js/adrs.js";
import { mountCheck } from "./js/check.js";
import { createDrawer } from "./js/detail.js";
import { initFlow, flowKeys, flowActions, flowHelp } from "./js/flow.js";

initToken();

const host = document.getElementById("view");
const drawerEl = document.getElementById("drawer");

let byIdCache = new Map();
let loaded = false;
let view = null; // the mounted view: { keys, focus, destroy }
let viewKey = "";
let base = { name: "board" }; // the last route that was not a record

const go = (route) => {
  location.hash = routeHref(route);
};

const drawer = createDrawer(drawerEl, host, {
  byId: () => byIdCache,
  records: () => state.records,
  close: () => go(base),
});

const viewCtx = {
  records: () => state.records,
  byId: () => byIdCache,
  today: () => isoDate(),
  setPage,
  open: (id) => go({ name: "record", id }),
  go: (slug) => go({ name: "tables", table: slug }),
  loadClosed: () => loadClosed(),
  check,
};

function routeKey(r) {
  return r.name === "tables" ? "tables:" + r.table : r.name;
}

const TITLES = { board: "Board", tables: "Tables", adrs: "ADRs", check: "Check" };
let repoName = ""; // the repository's label in the title, from GET /ui/config

/** The tab title for a route: the view, and for a record its id and title only. */
function setTitle(r) {
  let view = TITLES[r.name] || "";
  if (r.name === "tables") {
    const t = TABLES.find((x) => x.slug === r.table);
    view = t ? t.title : "Tables";
  } else if (r.name === "record") {
    const t = titleOfNow(r.id);
    view = "Record " + r.id + (t ? ": " + t : "");
  }
  document.title = pageTitle({ view, repo: repoName });
}

function mountBase(r, focus) {
  if (view) view.destroy();
  base = r;
  viewKey = routeKey(r);
  document.body.dataset.view = r.name;
  setTitle(r);
  for (const a of document.querySelectorAll("#nav [data-nav]")) {
    if (a.dataset.nav === r.name) a.setAttribute("aria-current", "page");
    else a.removeAttribute("aria-current");
  }
  switch (r.name) {
    case "tables":
      view = mountTables(host, viewCtx, r.table);
      break;
    case "adrs":
      view = mountAdrs(host, viewCtx);
      break;
    case "check":
      view = mountCheck(host, viewCtx);
      break;
    default:
      view = mountBoard(host, viewCtx);
  }
  if (focus) view.focus();
}

function showProblem(title, message) {
  if (view) {
    view.destroy();
    view = null;
    viewKey = "";
  }
  setPage({ mode: "ERROR", hints: [["r", "try again"], ["?", "help"]], where: "" });
  host.replaceChildren(el("div", { class: "state-card err", role: "alert" }, el("p", { class: "eyebrow", text: "Problem" }), el("h2", { text: title }), el("p", { text: message })));
}

/** Draw what the hash says. `focus` is true after a key or a click, false for the first draw. */
function route(focus) {
  if (!loaded) return;
  const r = parseRoute(location.hash);
  setTitle(r);
  if (r.name === "record") {
    if (!view) mountBase(base, false);
    setPage({ mode: "RECORD", hints: [["Esc", "close"], ["j k", "scroll"], ["1 2 3 4", "views"], ["?", "help"]], where: r.id });
    drawer.open(r.id);
    return;
  }
  const back = drawer.isOpen() ? drawer.close() : null;
  if (!view || routeKey(r) !== viewKey) mountBase(r, focus);
  else if (back) {
    // closing the drawer: put focus back on the card or row it came from
    view.focus();
  }
}

async function reload(announce) {
  try {
    await loadSnapshot();
    byIdCache = indexById(state.records);
    loaded = true;
    live.rebase(state.records);
    if (view) {
      viewKey = "";
      mountBase(base, false);
    } else route(false);
    if (drawer.isOpen()) drawer.open(drawer.id());
    if (announce) toast("Reloaded " + state.records.length + " records.");
  } catch (e) {
    if (e instanceof ApiError && e.code === "unauthorized") showProblem("No valid token", "Open the URL that akmctl ui printed; it holds the token in its fragment (#token=...). The token lasts as long as that akmctl ui process.");
    else if (!loaded) showProblem("The records did not load", describeError(e));
    else toast(describeError(e), "bad");
  }
}

// ---- live board: poll, animate, draw again in place ----

const titleOfNow = (id) => {
  const r = byIdCache.get(id);
  return r && r.fields && r.fields.title ? String(r.fields.title) : "";
};
let titleMemory = new Map(); // titles from before the last refresh, for records that are gone

function redrawInPlace() {
  if (!loaded || !view) return;
  byIdCache = indexById(state.records);
  keepSteady(
    host,
    () => {
      viewKey = "";
      mountBase(base, false);
    },
    () => view && view.focus(),
  );
  if (drawer.isOpen()) drawer.refresh();
}

const anim = createAnimator({
  toast,
  titleOf: (id) => titleOfNow(id) || titleMemory.get(id) || "",
  records: () => state.records,
  rerender: redrawInPlace,
  busy: () => overlayOpen() || typing(document.activeElement),
});

const setLiveStatus = mountIndicator(document.querySelector(".top-actions"));
const live = createLive({
  onStatus: setLiveStatus,
  async onEvents(events) {
    titleMemory = new Map(state.records.map((r) => [r.id, r.fields && r.fields.title ? String(r.fields.title) : ""]));
    await loadSnapshot(); // throws on a hiccup: live.js then offers the same events again
    byIdCache = indexById(state.records);
    anim.push(events);
  },
});

const NAV_KEYS = { 1: { name: "board" }, 2: { name: "tables" }, 3: { name: "adrs" }, 4: { name: "check" } };

function onKey(e) {
  if (drawer.isOpen() && drawer.keys(e)) return true;
  if (!drawer.isOpen() && view && view.keys(e)) return true;
  if (flowKeys(e, currentId())) return true;
  if (e.key === "r") {
    reload(true);
    return true;
  }
  const n = NAV_KEYS[e.key];
  if (n && !e.shiftKey) {
    go(n.name === "tables" ? { name: "tables", table: base.name === "tables" ? base.table : TABLES[0].slug } : n);
    return true;
  }
  return false;
}

/** The record a move key acts on: the open record, else the focused card of the board. */
const currentId = () => (drawer.isOpen() ? drawer.id() : view && view.current ? view.current() : null);

const actions = () => [
  { label: "Go to Board", keys: ["1"], run: () => go({ name: "board" }) },
  { label: "Go to Tables", keys: ["2"], run: () => go({ name: "tables" }) },
  ...TABLES.map((t) => ({ label: "Table: " + t.title, keys: [], run: () => go({ name: "tables", table: t.slug }) })),
  { label: "Go to ADRs", keys: ["3"], run: () => go({ name: "adrs" }) },
  { label: "Go to Check", keys: ["4"], run: () => go({ name: "check" }) },
  ...flowActions(currentId),
  { label: "Reload the records", keys: ["r"], run: () => reload(true) },
  { label: "Switch light and dark", keys: ["t"], run: toggleTheme },
  { label: "Show keyboard shortcuts", keys: ["?"], run: openHelp },
];

const help = () => [
  [
    "Board",
    [
      ["h l", "lane left, right"],
      ["j k", "card down, up"],
      ["g g", "first card in the lane"],
      ["G", "last card in the lane"],
      ["Enter", "open the record"],
    ],
  ],
  [
    "Tables, ADRs, Check",
    [
      ["j k", "row down, up"],
      ["g g", "first row"],
      ["G", "last row"],
      ["h l", "previous, next table"],
      ["Enter", "open the record"],
    ],
  ],
  flowHelp(),
  [
    "Record",
    [
      ["j k", "scroll"],
      ["Esc", "close"],
      ["d x b H L", "move it, as on the board"],
    ],
  ],
];

// the skip link must not become a route: move focus into the view instead
document.querySelector(".skip").addEventListener("click", (e) => {
  e.preventDefault();
  if (view) view.focus();
  else host.focus();
});

initFlow({ host, reload: () => reload(false), byId: () => byIdCache, records: () => state.records, focusView: () => view && !drawer.isOpen() && view.focus() });
initUI({ records: () => state.records, actions, onKey, help, openRecord: (id) => go({ name: "record", id }) });
window.addEventListener("hashchange", () => {
  // a new #token= on an open page (a second akmctl ui, a pasted URL): take it and start again
  if (parseToken(location.hash)) {
    initToken();
    loaded = false;
    start();
  } else route(true);
});

async function start() {
  if (!hasToken()) {
    showProblem("No token", "Open the URL that akmctl ui printed; it holds the token in its fragment (#token=...).");
    return;
  }
  setPage({ mode: "LOADING", hints: [], where: "" });
  const cfg = loadConfig().then((c) => {
    const name = (c.root || "").split(/[\\/]/).filter(Boolean).pop();
    document.getElementById("repo").textContent = name || "";
    document.getElementById("repo").title = c.root || "";
    repoName = repoLabel(c.root || "");
    setTitle(parseRoute(location.hash));
  });
  await reload(false);
  if (loaded) live.start(state.records);
  await cfg;
  if (!location.hash) history.replaceState(null, "", location.pathname + location.search + routeHref(base));
}

start();
