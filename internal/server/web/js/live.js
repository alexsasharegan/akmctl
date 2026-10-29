// The live board: poll the work listing, compare it with the last one and
// report what changed. Nothing here draws (the animator does); the DOM is
// touched only by mountIndicator. The diff is a pure function.
//
// Why polling: the API has no event stream. One cheap request every two
// seconds (id, path, etag and status of every work item) is enough to see an
// agent, an editor or Obsidian change a file.
import { el } from "./util.js";
import { pollWork, onWrite, writesInFlight } from "./store.js";

export const POLL_MS = 2000;
export const OFFLINE_AFTER = 3; // failed polls in a row before "offline" shows

// ---- what this page wrote ----

/** The etags that writes of this page produced, per record, so their effect is labelled "you". */
export function createOwn(keep = 8) {
  const byId = new Map();
  return {
    note(id, etag) {
      const list = byId.get(id) || [];
      list.push(etag);
      byId.set(id, list.slice(-keep));
    },
    has: (id, etag) => (byId.get(id) || []).includes(etag),
  };
}

const own = createOwn();
onWrite((r) => own.note(r.id, r.etag));

// ---- the diff ----

const inDone = (path) => /(^|\/)done\//.test(String(path || "").replace(/\\/g, "/"));

/** A listing as a Map: id to { id, path, etag, status } (work items only). */
export function snapshotOf(records) {
  const m = new Map();
  for (const r of records || []) {
    if (r.type !== "work") continue;
    m.set(r.id, { id: r.id, path: r.path || "", etag: r.etag || "", status: r.fields && r.fields.status != null ? String(r.fields.status) : "" });
  }
  return m;
}

/**
 * What changed between two snapshots. Events, in id order (removed ones first):
 *   added   { id, to }
 *   removed { id, from }
 *   moved   { id, from, to, closed }  status changed, or the file moved into done/
 *   edited  { id }                    the etag changed and nothing above did
 * Each has source "you" when the new etag is one this page's own write produced,
 * else "outside". A removal is always "outside": a write of this page renames
 * a file (close) but never deletes one.
 */
export function diffSnapshots(prev, next, mine = own) {
  const events = [];
  for (const [id, p] of prev) {
    if (!next.has(id)) events.push({ type: "removed", id, from: p.status, source: "outside" });
  }
  for (const [id, n] of next) {
    const p = prev.get(id);
    if (!p) {
      events.push({ type: "added", id, to: n.status, source: mine.has(id, n.etag) ? "you" : "outside" });
      continue;
    }
    if (p.etag === n.etag && p.path === n.path && p.status === n.status) continue;
    const source = mine.has(id, n.etag) ? "you" : "outside";
    const closed = inDone(n.path) && !inDone(p.path);
    if (p.status !== n.status || closed) events.push({ type: "moved", id, from: p.status, to: n.status, closed, source });
    else events.push({ type: "edited", id, source });
  }
  return events;
}

// ---- the words ----

const short = (s, n = 60) => (s.length > n ? s.slice(0, n - 1) + "…" : s);

/** One event as a sentence for a toast. title may be empty. */
export function describeEvent(ev, title) {
  const who = ev.source === "you" ? "You" : "Outside";
  const what = ev.id + (title ? " " + short(title) : "");
  if (ev.type === "moved") {
    const lanes = ev.closed && ev.from === ev.to ? "closed" : ev.from + " to " + ev.to;
    return ev.source === "you" ? "You moved " + what + ": " + lanes + "." : "Outside: " + what + " moved, " + lanes + ".";
  }
  if (ev.type === "added") return who + (ev.source === "you" ? " added " : ": added ") + what + (ev.to ? " (" + ev.to + ")." : ".");
  if (ev.type === "removed") return who + (ev.source === "you" ? " removed " : ": removed ") + what + ".";
  return who + (ev.source === "you" ? " edited " : ": edited ") + what + ".";
}

/** A burst as one sentence: "12 outside changes (9 moved, 2 added, 1 edited)." */
export function summarize(events) {
  const outside = events.filter((e) => e.source !== "you").length;
  const count = {};
  for (const e of events) count[e.type] = (count[e.type] || 0) + 1;
  const parts = ["moved", "added", "removed", "edited"].filter((t) => count[t]).map((t) => count[t] + " " + t);
  const who = outside === events.length ? "outside" : outside === 0 ? "of your own" : "(" + outside + " outside)";
  return events.length + " changes " + who + " (" + parts.join(", ") + ").";
}

// ---- the loop ----

/**
 * hooks:
 *   onEvents(events)  may return a promise; if it rejects the same events come again next poll
 *   onStatus(s)       "ok", "offline" (after repeated failures) or "signed-out" (the token no longer works)
 *   fetchWork()       the listing (default: the API)
 * Returns { start(records), rebase(records), stop(), poll() }. The loop runs only while the tab is visible.
 */
export function createLive({ onEvents, onStatus = () => {}, fetchWork = pollWork, interval = POLL_MS } = {}) {
  let prev = null;
  let timer = null;
  let inflight = false;
  let fails = 0;
  let status = "ok";
  let running = false;

  const visible = () => typeof document === "undefined" || document.visibilityState !== "hidden";
  const setStatus = (s) => {
    if (s !== status) {
      status = s;
      onStatus(s);
    }
  };
  const later = (ms) => {
    clearTimeout(timer);
    timer = running && visible() ? setTimeout(poll, ms) : null;
  };

  async function poll() {
    clearTimeout(timer);
    timer = null;
    if (!running || !visible()) return;
    if (inflight || writesInFlight() > 0) return later(300);
    inflight = true;
    const before = prev;
    try {
      const data = await fetchWork();
      // a write of this page finished or started meanwhile: its etag may not be noted yet, look again
      if (writesInFlight() > 0) {
        inflight = false;
        return later(300);
      }
      const next = snapshotOf(data.records);
      if (prev) {
        const events = diffSnapshots(prev, next);
        prev = next;
        if (events.length) await onEvents(events);
      } else prev = next;
      fails = 0;
      setStatus("ok");
    } catch (e) {
      prev = before;
      if (e && e.code === "unauthorized") {
        // a restarted akmctl ui has a new token; nothing here can fix that
        running = false;
        inflight = false;
        setStatus("signed-out");
        return;
      }
      fails++;
      if (fails >= OFFLINE_AFTER) setStatus("offline");
    }
    inflight = false;
    later(fails ? Math.min(interval * (1 + fails), 10000) : interval);
  }

  const onVisibility = () => {
    if (!running) return;
    if (visible()) poll();
    else clearTimeout(timer);
  };

  return {
    /** Begin polling; `records` is the listing the page already holds (the first comparison is with it). */
    start(records) {
      prev = records ? snapshotOf(records) : null;
      fails = 0;
      status = "ok";
      onStatus("ok");
      if (!running && typeof document !== "undefined") document.addEventListener("visibilitychange", onVisibility);
      running = true;
      later(interval);
    },
    /** The page reloaded everything: compare future polls with that. */
    rebase(records) {
      prev = snapshotOf(records);
    },
    stop() {
      running = false;
      clearTimeout(timer);
      if (typeof document !== "undefined") document.removeEventListener("visibilitychange", onVisibility);
    },
    poll,
  };
}

// ---- the offline mark ----

/** A quiet mark in the top bar, hidden until the server cannot be reached. Returns set(status). */
export function mountIndicator(parent) {
  const mark = el("span", { class: "live-status", role: "status", hidden: true });
  parent.prepend(mark);
  return (s) => {
    mark.hidden = s === "ok";
    mark.textContent = s === "offline" ? "◌ offline" : s === "signed-out" ? "◌ signed out" : "";
    mark.title = s === "offline" ? "The server does not answer; trying again." : s === "signed-out" ? "The token no longer works (akmctl ui was restarted?). Open the URL it printed." : "";
  };
}
