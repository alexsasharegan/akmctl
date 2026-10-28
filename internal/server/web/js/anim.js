// The animation queue of the live board. Events from live.js play one at a
// time, in the order they came; a "theme" draws one move. When the queue is
// empty and the page is idle the views are drawn again in place, from the
// records the page already refreshed. No page reload.
//
// Rules: more than BACKLOG events at once are not worth watching (show the
// result); prefers-reduced-motion: reduce plays nothing; the re-draw waits
// while a dialog, a text field or a drag is active.
import { el } from "./util.js";
import { LANES, lanes } from "./logic.js";
import { describeEvent, summarize } from "./live.js";

export const BACKLOG = 10;
const IDLE_MS = 1500;

const wait = (ms) => new Promise((r) => setTimeout(r, ms));
const reduced = () => matchMedia("(prefers-reduced-motion: reduce)").matches;
const cardOf = (id) => [...document.querySelectorAll(".card")].find((c) => c.dataset.id === id) || null;
const laneStatus = (col) => {
  const c = [...col.classList].find((x) => x.startsWith("s-"));
  return c ? c.slice(2) : "";
};
const laneCol = (status) => (LANES.includes(status) ? [...document.querySelectorAll(".col")].find((c) => c.classList.contains("s-" + status)) || null : null);

/** Make a lane's count and empty note true to the cards it holds. */
function recount(col) {
  if (!col) return;
  const n = col.querySelectorAll(".card").length;
  const count = col.querySelector(".count");
  if (count) {
    count.textContent = "(" + n + ")";
    count.setAttribute("aria-label", n + " items");
  }
  const body = col.querySelector(".col-body");
  const note = body.querySelector(".empty-col");
  if (n && note) note.remove();
  if (!n && !note) body.append(el("div", { class: "empty-col", text: "Nothing in " + laneStatus(col) + "." }));
}

const play = (card, frames, opts) => {
  if (!card.animate) return Promise.resolve();
  // a page that is not being painted (a hidden pane) never finishes an animation: do not wait for it
  return Promise.race([card.animate(frames, opts).finished.catch(() => {}), wait(opts.duration + 200)]);
};

// ---- themes ----
// A theme is async (event, ctx) => boolean: true if it drew something.
// ctx: records() the refreshed records. "quiet" is the only one.

export const themes = {
  /** The card lifts, slides to its place in the new lane (or fades out if it left the board) and settles. */
  async quiet(ev, ctx) {
    const card = cardOf(ev.id);
    if (!card) return false;
    const from = card.closest(".col");
    const to = ev.closed ? null : laneCol(ev.to);
    if (to && from === to) return false;
    const focused = document.activeElement === card;
    card.classList.add("anim-lift");
    await wait(260);
    if (!card.isConnected) return true;
    if (!to) {
      card.classList.remove("anim-lift");
      card.classList.add("anim-flying");
      await play(card, [{ opacity: 1 }, { opacity: 0 }], { duration: 360, easing: "ease-out" });
      card.remove();
      recount(from);
      return true;
    }
    const first = card.getBoundingClientRect();
    card.classList.remove("anim-lift");
    card.classList.add("anim-flying");
    // its place in the new lane: before the next card of the refreshed order that is on the board
    const body = to.querySelector(".col-body");
    const order = (lanes(ctx.records()).find((l) => l.status === ev.to) || { records: [] }).records.map((r) => r.id);
    const after = order.slice(order.indexOf(ev.id) + 1);
    const ref = [...body.querySelectorAll(".card")].find((c) => after.includes(c.dataset.id)) || null;
    body.querySelector(".empty-col")?.remove();
    body.insertBefore(card, ref);
    if (focused) card.focus({ preventScroll: true });
    const last = card.getBoundingClientRect();
    await play(card, [{ transform: "translate(" + (first.left - last.left) + "px," + (first.top - last.top) + "px)" }, { transform: "none" }], { duration: 480, easing: "cubic-bezier(.2,.7,.2,1)" });
    card.classList.remove("anim-flying");
    recount(from);
    recount(to);
    return true;
  },
};

// ---- keeping the page steady across a re-draw ----

/**
 * Run `redraw` (which replaces the contents of `host`) and put back what the
 * reader had: scroll positions and the focused card or row. If that card or
 * row is gone, `fallback` moves focus somewhere sensible.
 */
export function keepSteady(host, redraw, fallback) {
  const scrollers = () => [host, ...host.querySelectorAll(".col-body, .table-wrap, .tables-body")];
  const tops = scrollers().map((n) => [n.scrollTop, n.scrollLeft]);
  const a = document.activeElement;
  const was = a && host.contains(a) && a.closest("[data-id]") ? a.closest("[data-id]") : null;
  const inHost = !!(a && a !== document.body && host.contains(a));
  const id = was && was.dataset.id;
  const kind = was && (was.matches("tr.row") ? "tr.row" : ".card");
  redraw();
  scrollers().forEach((n, i) => {
    if (tops[i]) {
      n.scrollTop = tops[i][0];
      n.scrollLeft = tops[i][1];
    }
  });
  if (!inHost) return;
  const now = id ? [...host.querySelectorAll(kind)].find((n) => n.dataset.id === id) : null;
  if (now) now.focus({ preventScroll: true });
  else if (fallback) fallback();
}

// ---- the queue ----

/**
 * ctx: toast(msg), titleOf(id), records(), rerender(touchedIds) draws the
 * views again in place, busy() true while a dialog or a field is active.
 */
export function createAnimator(ctx) {
  const queue = [];
  const touched = new Set();
  const holds = [];
  let playing = false;
  let settling = false;
  let theme = "quiet";

  const busy = () => !!(ctx.busy && ctx.busy()) || !!document.querySelector(".card.dragging") || holds.some((f) => f());
  const say = (ev) => describeEvent(ev, ctx.titleOf(ev.id));

  async function playOne(ev, skip) {
    ctx.toast(say(ev));
    touched.add(ev.id);
    if (skip || ev.type !== "moved" || reduced()) return;
    try {
      await (themes[theme] || themes.quiet)(ev, ctx);
    } catch (e) {
      /* the re-draw makes the board right anyway */
    }
    await wait(180);
  }

  async function drain() {
    if (playing) return;
    playing = true;
    while (queue.length) await playOne(queue.shift(), queue.length >= BACKLOG);
    playing = false;
    settle();
  }

  // Draw the views again once nothing is playing and nothing is in the way.
  async function settle() {
    if (settling) return;
    settling = true;
    await wait(200);
    while (playing || queue.length || busy()) await wait(playing || queue.length ? 200 : IDLE_MS);
    settling = false;
    const ids = [...touched];
    touched.clear();
    ctx.rerender(ids);
    flash(ids);
  }

  return {
    /** Events from live.js; the records are already refreshed. */
    push(events) {
      if (!events.length) return;
      if (events.length > BACKLOG || reduced()) {
        // a burst is shown as its result, with one sentence
        if (events.length > 1) ctx.toast(summarize(events));
        else ctx.toast(say(events[0]));
        for (const e of events) touched.add(e.id);
        settle();
        return;
      }
      queue.push(...events);
      drain();
    },
    use(name) {
      if (themes[name]) theme = name;
    },
    /** Another module (a drag, a form) says: not now, while fn() is true. */
    holdWhile(fn) {
      holds.push(fn);
    },
  };
}

/** A short highlight on the cards and rows that changed (the motion itself is in the stylesheet). */
function flash(ids) {
  const want = new Set(ids);
  for (const n of document.querySelectorAll(".card[data-id], tr.row[data-id]")) {
    if (!want.has(n.dataset.id)) continue;
    n.classList.add("anim-settle");
    setTimeout(() => n.classList.remove("anim-settle"), 900);
  }
}
