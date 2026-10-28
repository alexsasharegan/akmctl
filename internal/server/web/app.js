// The workbench: reads the token from the URL fragment, loads the records once
// and draws the view that the hash route names. Everything from the repository
// is shown as text, never as HTML.
import { initToken, parseToken, hasToken, loadSnapshot, loadConfig, loadClosed, check, state, describeError, ApiError } from "./js/store.js";
import { parseRoute, routeHref, indexById, TABLES } from "./js/logic.js";
import { el, isoDate } from "./js/util.js";
import { initUI, setPage, toast, toggleTheme, openHelp } from "./js/ui.js";
import { mountBoard } from "./js/board.js";
import { mountTables } from "./js/tables.js";
import { mountAdrs } from "./js/adrs.js";
import { mountCheck } from "./js/check.js";
import { createDrawer } from "./js/detail.js";

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

function mountBase(r, focus) {
  if (view) view.destroy();
  base = r;
  viewKey = routeKey(r);
  document.body.dataset.view = r.name;
  document.title = TITLES[r.name] + " · akmctl workbench";
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

const NAV_KEYS = { 1: { name: "board" }, 2: { name: "tables" }, 3: { name: "adrs" }, 4: { name: "check" } };

function onKey(e) {
  if (drawer.isOpen() && drawer.keys(e)) return true;
  if (!drawer.isOpen() && view && view.keys(e)) return true;
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

const actions = () => [
  { label: "Go to Board", keys: ["1"], run: () => go({ name: "board" }) },
  { label: "Go to Tables", keys: ["2"], run: () => go({ name: "tables" }) },
  ...TABLES.map((t) => ({ label: "Table: " + t.title, keys: [], run: () => go({ name: "tables", table: t.slug }) })),
  { label: "Go to ADRs", keys: ["3"], run: () => go({ name: "adrs" }) },
  { label: "Go to Check", keys: ["4"], run: () => go({ name: "check" }) },
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
  [
    "Record",
    [
      ["j k", "scroll"],
      ["Esc", "close"],
    ],
  ],
];

// the skip link must not become a route: move focus into the view instead
document.querySelector(".skip").addEventListener("click", (e) => {
  e.preventDefault();
  if (view) view.focus();
  else host.focus();
});

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
  });
  await reload(false);
  await cfg;
  if (!location.hash) history.replaceState(null, "", location.pathname + location.search + routeHref(base));
}

start();
