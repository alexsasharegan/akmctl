// Pure functions over the loaded records. No DOM, no network, no clock unless
// a date is passed in, so each can be read and checked on its own.
import { ageDays } from "./util.js";

/** The folders of a repository that keeps the defaults; `version.repository` names others. */
export const DEFAULT_FOLDERS = { work_dir: "docs/work", adr_dir: "docs/adr" };

/** The work and ADR folders from a `version` response (the manifest's work_dir and adr_dir), else the defaults. */
export function foldersFrom(version) {
  const r = (version && version.repository) || {};
  const ok = (v) => typeof v === "string" && v !== "";
  return { work_dir: ok(r.work_dir) ? r.work_dir : DEFAULT_FOLDERS.work_dir, adr_dir: ok(r.adr_dir) ? r.adr_dir : DEFAULT_FOLDERS.adr_dir };
}

export const LANES = ["raw", "ready", "active", "blocked"];
export const CLOSED = ["done", "dropped"];
export const STALE_DAYS = 30; // the default of rule WK10

export const field = (r, k) => (r && r.fields && r.fields[k] !== undefined ? r.fields[k] : null);
export const isWork = (r) => r.type === "work";
export const isOpenWork = (r) => isWork(r) && !CLOSED.includes(field(r, "status"));
export const priorityOf = (r) => (typeof field(r, "priority") === "number" ? field(r, "priority") : 0);
export const blockedBy = (r) => {
  const v = field(r, "blocked_by");
  return Array.isArray(v) ? v : [];
};

/** Index records by id. */
export function indexById(records) {
  const m = new Map();
  for (const r of records) m.set(r.id, r);
  return m;
}

/** Children of a record: work items whose parent is this id, in snapshot order. */
export function childrenOf(records, id) {
  return records.filter((r) => isWork(r) && field(r, "parent") === id);
}

// ---- board ----

/** The four lanes with their work items, in the order the records came in. */
export function lanes(records) {
  return LANES.map((status) => ({ status, records: records.filter((r) => isWork(r) && field(r, "status") === status) }));
}

/** Where a card is: { lane, index }, or null. */
export function locate(ls, id) {
  for (let lane = 0; lane < ls.length; lane++) {
    const index = ls[lane].records.findIndex((r) => r.id === id);
    if (index >= 0) return { lane, index };
  }
  return null;
}

/** The first card on the board (first non-empty lane), or null. */
export function firstLoc(ls) {
  const lane = ls.findIndex((l) => l.records.length);
  return lane < 0 ? null : { lane, index: 0 };
}

/**
 * Move focus. `dir` is left, right, up, down, first or last. Left and right go
 * to the next lane that has cards and keep the row where they can; first and
 * last go to the first and last card of the current lane. Returns the new
 * location, or the old one when there is nowhere to go.
 */
export function step(ls, loc, dir) {
  if (!loc) return firstLoc(ls);
  const n = (l) => ls[l].records.length;
  switch (dir) {
    case "up":
      return { lane: loc.lane, index: Math.max(0, loc.index - 1) };
    case "down":
      return { lane: loc.lane, index: Math.min(n(loc.lane) - 1, loc.index + 1) };
    case "first":
      return { lane: loc.lane, index: 0 };
    case "last":
      return { lane: loc.lane, index: n(loc.lane) - 1 };
    case "left":
    case "right": {
      const d = dir === "left" ? -1 : 1;
      for (let l = loc.lane + d; l >= 0 && l < ls.length; l += d) {
        if (n(l)) return { lane: l, index: Math.min(loc.index, n(l) - 1) };
      }
      return loc;
    }
  }
  return loc;
}

/** The staleness of a record: days since `updated`, and whether that is past the limit. */
export function staleness(r, today, limit = STALE_DAYS) {
  const days = ageDays(field(r, "updated") || "", today);
  return { days, stale: days !== null && days >= limit };
}

// ---- search ----

const words = (q) => q.toLowerCase().split(/\s+/).filter(Boolean);

/** 0 when any word of the query is missing; higher is a better match (earlier is better). */
export function score(query, hay) {
  const ws = words(query);
  if (!ws.length) return 1;
  const h = String(hay || "").toLowerCase();
  let sc = 0;
  for (const w of ws) {
    const i = h.indexOf(w);
    if (i < 0) return 0;
    sc += 100 - Math.min(i, 90);
  }
  return sc;
}

/** Search the loaded records (id, title, summary, kind, status, scope). Best first, at most `limit`. */
export function searchRecords(records, query, limit = 8) {
  if (!words(query).length) return [];
  const out = [];
  for (const r of records) {
    const inTitle = score(query, field(r, "title")) * 2;
    const all = score(query, [field(r, "title"), r.id, field(r, "summary"), field(r, "kind"), field(r, "status"), field(r, "scope")].join(" "));
    const best = Math.max(inTitle, all);
    if (best) out.push({ record: r, score: best });
  }
  out.sort((a, b) => b.score - a.score || (a.record.id < b.record.id ? -1 : 1));
  return out.slice(0, limit);
}

// ---- tables ----

export const TABLES = [
  { slug: "ready-active", title: "Ready and active" },
  { slug: "by-parent", title: "Open by parent" },
  { slug: "blocked", title: "Blocked" },
  { slug: "stale", title: "Stale by updated" },
  { slug: "closed", title: "Recently closed" },
];

const byPriorityDesc = (a, b) => priorityOf(b) - priorityOf(a);

/** Work with status ready or active: active first, then by priority (the snapshot order breaks ties). */
export function readyActive(records) {
  const rank = (r) => (field(r, "status") === "active" ? 0 : 1);
  return records
    .filter((r) => isWork(r) && ["ready", "active"].includes(field(r, "status")))
    .map((r, i) => ({ r, i }))
    .sort((a, b) => rank(a.r) - rank(b.r) || byPriorityDesc(a.r, b.r) || a.i - b.i)
    .map((x) => x.r);
}

/** Open work grouped by parent: [{ parent: id|null, title, known, records }]. Parents in id order, no parent last. */
export function openByParent(records) {
  const byId = indexById(records);
  const groups = new Map();
  for (const r of records.filter(isOpenWork)) {
    const p = field(r, "parent");
    if (!groups.has(p)) groups.set(p, []);
    groups.get(p).push(r);
  }
  const keys = [...groups.keys()].sort((a, b) => (a === null) - (b === null) || (a < b ? -1 : a > b ? 1 : 0));
  return keys.map((p) => {
    const rec = p ? byId.get(p) : null;
    return { parent: p, title: p === null ? "No parent" : rec ? field(rec, "title") : "Not in this repository", known: !!rec, records: groups.get(p) };
  });
}

/** Open work that is blocked: status blocked, or a non-empty blocked_by. */
export function blockedRows(records) {
  return records.filter((r) => isOpenWork(r) && (field(r, "status") === "blocked" || blockedBy(r).length > 0));
}

/** Open work, least recently updated first, with its age. */
export function staleRows(records, today) {
  return records
    .filter(isOpenWork)
    .map((r) => ({ record: r, ...staleness(r, today) }))
    .sort((a, b) => (b.days ?? -1) - (a.days ?? -1) || (a.record.id < b.record.id ? -1 : 1));
}

/** Closed work, most recently updated first. */
export function closedRows(records) {
  return [...records].sort((a, b) => String(field(b, "updated") || "").localeCompare(String(field(a, "updated") || "")) || (a.id < b.id ? 1 : -1));
}

/** ADRs, newest id first. */
export function adrRows(records) {
  return records.filter((r) => r.type === "adr").sort((a, b) => (a.id < b.id ? 1 : -1));
}

// ---- check ----

export const SEVERITY = { E: "Error", W: "Warning" };

/** Findings with errors first; the order inside each severity is kept. */
export function orderFindings(findings) {
  const rank = (x) => (x.severity === "E" ? 0 : 1);
  return findings
    .map((x, i) => ({ x, i }))
    .sort((a, b) => rank(a.x) - rank(b.x) || a.i - b.i)
    .map((o) => o.x);
}

// ---- routes ----

const IDS = /^[0-9]{6}-[0-9a-z]{4}$/;
export const isId = (s) => IDS.test(String(s));

/** A hash such as "#/tables/blocked" as { name, ... }. Unknown hashes are the board. */
export function parseRoute(hash) {
  const parts = String(hash || "")
    .replace(/^#\/?/, "")
    .split("/")
    .filter(Boolean)
    .map((p) => {
      try {
        return decodeURIComponent(p);
      } catch (e) {
        return p;
      }
    });
  const [name, arg] = parts;
  switch (name) {
    case "tables":
      return { name: "tables", table: TABLES.some((t) => t.slug === arg) ? arg : TABLES[0].slug };
    case "adrs":
    case "check":
      return { name };
    case "record":
      return isId(arg) ? { name: "record", id: arg } : { name: "board" };
  }
  return { name: "board" };
}

export function routeHref(route) {
  switch (route.name) {
    case "tables":
      return "#/tables/" + (route.table || TABLES[0].slug);
    case "record":
      return "#/record/" + encodeURIComponent(route.id);
    default:
      return "#/" + route.name;
  }
}

/**
 * The obsidian:// URL that opens a record, or "" when the root is unknown.
 * `cfg` is GET /ui/config ({ root }); `path` the record's path relative to the
 * root. Always the absolute `path=` form: Obsidian finds the vault that holds
 * the file, so the repository root need not be the vault root.
 */
export function obsidianUrl(cfg, path) {
  const root = cfg && cfg.root ? String(cfg.root) : "";
  if (!root || !path) return "";
  const rel = String(path).replace(/^\/+/, "");
  const sep = root.includes("\\") && !root.includes("/") ? "\\" : "/";
  const abs = root.replace(/[\\/]+$/, "") + sep + rel.split("/").join(sep);
  return "obsidian://open?path=" + encodeURIComponent(abs);
}

/**
 * The repository's label in the page title: the last folder of its root path, the name the header
 * shows. Two open repositories with the same folder name are told apart by the parent folder, as
 * "parent/name", when `others` (the other roots the page knows of) holds a root with that name and
 * a different path. The page only knows its own root, so it passes [] today and the label is the
 * folder name.
 */
export function repoLabel(root, others = []) {
  const split = (p) => String(p || "").split(/[\\/]/).filter(Boolean);
  const parts = split(root);
  const name = parts[parts.length - 1] || "";
  if (!name) return "";
  const full = parts.join("/");
  const clash = others.some((o) => {
    const p = split(o);
    return p[p.length - 1] === name && p.join("/") !== full;
  });
  return clash && parts.length > 1 ? parts[parts.length - 2] + "/" + name : name;
}

/**
 * The document title: "<view> - <repository> - akmctl". `view` is a view name ("Board", "ADRs"...),
 * or for a record "Record <id>: <title>". Only the record's id and title appear; no other record text.
 */
export function pageTitle({ view, repo = "" }) {
  return [view, repo, "akmctl"].filter(Boolean).join(" - ");
}
