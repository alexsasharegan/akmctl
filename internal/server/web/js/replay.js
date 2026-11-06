// Replay of the moves made while the page was away (work item 261101-w0mz,
// ADR 261106-60xv). The poll in live.js only sees the net change between two
// snapshots, in id order, with no author. The `history` operation lists the
// step commits in order with who made them; this module turns those moves into
// the events live.js and anim.js already understand, and lists a long batch so
// it is not lost. Everything here is pure except mountAway.
import { el, clear } from "./util.js";
import { describeEvent } from "./live.js";

/** More moves than this are listed, not played one by one (the animator's backlog). */
export const AWAY_MAX = 10;

/** Who made a move: the Assisted-by trailer's value, else the commit's author. */
export function whoOf(move) {
  for (const t of move.trailers || []) {
    const m = /^Assisted-by:\s*(.+)$/i.exec(t);
    if (m && m[1].trim()) return m[1].trim();
  }
  return move.author || "";
}

/**
 * One history move as an event: a lane change is "moved", a new file "added",
 * and a step that changes no status (claim, release) "step". Every event is
 * "outside" with `who`, `at` (the commit time) and `sha`.
 */
export function moveToEvent(m) {
  const base = { id: m.item, source: "outside", who: whoOf(m), at: m.time || "", sha: m.sha || "" };
  if (m.step === "file") return { ...base, type: "added", to: m.to };
  const closed = m.step === "close" || m.step === "drop";
  if (m.from !== m.to || closed) return { ...base, type: "moved", from: m.from, to: m.to, closed };
  return { ...base, type: "step", step: m.step };
}

/**
 * The poll's net events with the history's moves put in their place: for an
 * item that moved outside and whose last committed move ends where the poll
 * says it ended, the moves replace the net event, in the order they were
 * committed. Everything else (added, removed, edited, an uncommitted move, a
 * move the history does not explain) stays, after the replayed moves.
 */
export function mergeWithHistory(events, moves) {
  const byId = new Map();
  for (const m of moves || []) {
    if (!byId.has(m.item)) byId.set(m.item, []);
    byId.get(m.item).push(m);
  }
  const covered = new Set();
  for (const e of events) {
    if (e.type !== "moved" || e.source === "you") continue;
    const ms = (byId.get(e.id) || []).filter((m) => m.from !== m.to || m.step === "close" || m.step === "drop");
    if (ms.length && ms[ms.length - 1].to === e.to) covered.add(e.id);
  }
  const replay = (moves || []).filter((m) => covered.has(m.item)).map(moveToEvent);
  const rest = events.filter((e) => !(e.type === "moved" && e.source !== "you" && covered.has(e.id)));
  return [...replay, ...rest];
}

/** A commit time as a short label: HH:MM today, else the date and HH:MM; "" if it is not a time. */
export function timeLabel(iso, now = new Date()) {
  const m = /^(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2})/.exec(String(iso || ""));
  if (!m) return "";
  const today = new Date(now.getTime() - now.getTimezoneOffset() * 60000).toISOString().slice(0, 10);
  return m[1] === today ? m[2] : m[1] + " " + m[2];
}

/** The lines of the "while you were away" list: [{ time, text }], in the order given. */
export function awayLines(events, titleOf = () => "", now = new Date()) {
  return events.map((e) => ({ time: timeLabel(e.at, now), text: describeEvent(e, titleOf(e.id)) }));
}

/** The heading: "12 moves while you were away" (and "since the last 100" when older ones were left out). */
export function awayHeading(count, truncated) {
  return count + (count === 1 ? " move" : " moves") + " while you were away" + (truncated ? " (older ones are not listed)" : "");
}

// ---- the panel ----

/**
 * A quiet panel at the foot of the page. show(events, { truncated, titleOf })
 * replaces what it lists; it is dismissed with its button or Escape.
 */
export function mountAway(parent) {
  const title = el("h2", { class: "away-title" });
  const list = el("ol", { class: "away-list" });
  const close = el("button", { type: "button", class: "away-close", text: "Dismiss", onclick: () => hide() });
  const panel = el("aside", { class: "away", role: "region", "aria-label": "Moves while you were away", hidden: true }, el("div", { class: "away-head" }, title, close), list);
  panel.addEventListener("keydown", (e) => {
    if (e.key === "Escape") {
      e.stopPropagation();
      hide();
    }
  });
  parent.append(panel);
  function hide() {
    panel.hidden = true;
  }
  return {
    show(events, { truncated = false, titleOf = () => "" } = {}) {
      if (!events.length) return;
      title.textContent = awayHeading(events.length, truncated);
      clear(list);
      for (const l of awayLines(events, titleOf)) list.append(el("li", {}, l.time ? el("time", { text: l.time }) : null, " ", l.text));
      panel.hidden = false;
    },
    hide,
    isOpen: () => !panel.hidden,
  };
}
