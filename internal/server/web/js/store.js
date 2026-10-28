// The token, the API client and the snapshot of the repository's records.
// Nothing here builds HTML.
import { isoDate, daysBefore } from "./util.js";

const TOKEN_KEY = "akm-token";

let token = "";

/** The value of `token` in a URL fragment such as "#token=abc", or "". */
export function parseToken(hash) {
  return new URLSearchParams(String(hash || "").replace(/^#/, "")).get("token") || "";
}

/**
 * Read the token once: from the fragment if there (store it in sessionStorage
 * and strip it from the address bar), else from sessionStorage.
 */
export function initToken(loc = location, hist = history) {
  const fromUrl = parseToken(loc.hash);
  if (fromUrl) {
    token = fromUrl;
    try {
      sessionStorage.setItem(TOKEN_KEY, token);
    } catch (e) {
      /* the token stays in memory for this page */
    }
    hist.replaceState(null, "", loc.pathname + loc.search);
  } else {
    try {
      token = sessionStorage.getItem(TOKEN_KEY) || "";
    } catch (e) {
      token = "";
    }
  }
  return token;
}

export const hasToken = () => token !== "";

export class ApiError extends Error {
  constructor(message, { code = "internal", status = 0, details = null } = {}) {
    super(message);
    this.name = "ApiError";
    this.code = code;
    this.status = status;
    this.details = details;
  }
}

/** The text shown for an API error: the API's own message, with the reason when it gives one. */
export function describeError(e) {
  if (!(e instanceof ApiError)) return e && e.message ? e.message : String(e);
  const d = e.details || {};
  const why = [d.reason, d.rule].filter(Boolean).join(", ");
  return e.code + ": " + e.message + (why ? " (" + why + ")" : "");
}

/** Build "path?a=1&b=x,y": lists are joined with commas, null and undefined are left out. */
export function withQuery(path, query) {
  const q = new URLSearchParams();
  for (const [k, v] of Object.entries(query || {})) {
    if (v === null || v === undefined || v === "") continue;
    q.set(k, Array.isArray(v) ? v.join(",") : String(v));
  }
  const s = q.toString();
  return s ? path + "?" + s : path;
}

// Live board hooks: the writes this page makes are reported, so a change it
// made itself is told apart from one made outside (see live.js).
const writeHooks = [];
let writesPending = 0;

/** Call fn(record) with the record of every successful write this page makes (a dry run is not one). */
export const onWrite = (fn) => writeHooks.push(fn);

/** How many writes of this page have not finished. */
export const writesInFlight = () => writesPending;

/** One call to /api/v1 (or another path starting with "/"). Resolves to the response object; rejects with ApiError. */
export async function api(method, path, opts = {}) {
  const writing = method !== "GET";
  if (writing) writesPending++;
  try {
    const data = await call(method, path, opts);
    if (writing && data.record && data.record.id && data.record.etag && !(opts.body && opts.body.dry_run)) {
      for (const fn of writeHooks) fn(data.record);
    }
    return data;
  } finally {
    if (writing) writesPending--;
  }
}

async function call(method, path, { query, body, etag } = {}) {
  const url = path.startsWith("/") ? path : "/api/v1/" + path;
  const headers = { Authorization: "Bearer " + token };
  if (body !== undefined) headers["Content-Type"] = "application/json";
  if (etag) headers["If-Match"] = '"' + etag + '"';
  let res;
  try {
    res = await fetch(withQuery(url, query), { method, headers, body: body === undefined ? undefined : JSON.stringify(body) });
  } catch (e) {
    throw new ApiError("cannot reach the server (is akmctl ui still running?)", { code: "unreachable" });
  }
  let data = null;
  try {
    data = await res.json();
  } catch (e) {
    /* not JSON */
  }
  if (!res.ok || !data || data.ok === false) {
    const err = (data && data.error) || {};
    throw new ApiError(err.message || res.status + " " + res.statusText, { code: err.code || (res.status === 401 ? "unauthorized" : "internal"), status: res.status, details: err.details });
  }
  return data;
}

// ---- the snapshot ----

// What the board, tables, ADR list and palette need from `ls`.
export const SNAPSHOT_FIELDS = ["id", "title", "summary", "status", "kind", "parent", "priority", "created", "updated", "blocked_by", "date", "scope"];

export const state = {
  records: [], // work and ADRs, work in roadmap order
  loadedAt: null,
  config: { root: "", identity: "", obsidian_vault: "" },
  closed: null, // recently closed work, loaded when the table is first shown
};

/** Load every record's header once. Work comes in the order the board shows (roadmap, priority, id). */
export async function loadSnapshot() {
  const data = await api("GET", "records", { query: { fields: SNAPSHOT_FIELDS, sort: ["roadmap", "-priority", "id"] } });
  state.records = data.records;
  state.loadedAt = new Date();
  state.closed = null;
  return state.records;
}

export async function loadConfig() {
  try {
    state.config = { ...state.config, ...(await api("GET", "/ui/config")) };
  } catch (e) {
    /* the page works without it */
  }
  return state.config;
}

/** Work closed in the last `days` days (done or dropped), newest first. */
export async function loadClosed(days = 30) {
  const data = await api("GET", "records", { query: { type: "work", status: ["done", "dropped"], since: daysBefore(days, isoDate()), fields: SNAPSHOT_FIELDS, sort: ["-updated", "-id"] } });
  state.closed = data.records;
  return state.closed;
}

export const show = (id, section) => api("GET", "records/" + encodeURIComponent(id), { query: { section } });

/** Body search: the records whose title, summary, scope or body hold every word. */
export const searchBody = (text, limit = 10) => api("GET", "search", { query: { text, body: true, limit, fields: ["id", "title", "status", "kind"] } });

export const check = () => api("GET", "check");

/** The cheap listing the live board polls: id, path and etag of every work item, with its status. */
export const pollWork = () => api("GET", "records", { query: { type: "work", fields: ["status", "updated"], sort: ["id"] } });
