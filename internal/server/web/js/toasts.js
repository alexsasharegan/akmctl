// The toast queue: no DOM, no timers, so it is unit tested. The caller passes
// the clock in (`now`, milliseconds) and owns the setTimeout; the queue knows
// the limit, which toasts are persistent, how long each has left, and what a
// hover or keyboard focus (a "hold") pauses.

export const LEVELS = ["success", "info", "warning", "error"];
export const EPHEMERAL_MS = 8000;
export const LIMIT = 3;

/**
 * Turn toast(msg, opts) into a spec. opts: { level, persistent, ms }. The level
 * defaults to info; errors are persistent unless persistent is false; the other
 * levels are ephemeral unless persistent is true. A bad level or a non-object
 * opts throws, so a caller left on the old string kinds fails loudly.
 */
export function toastSpec(msg, opts) {
  if (opts === undefined || opts === null) opts = {};
  if (typeof opts !== "object") throw new TypeError("toast options must be an object { level, persistent, ms }");
  const level = opts.level === undefined ? "info" : opts.level;
  if (!LEVELS.includes(level)) throw new RangeError("unknown toast level: " + level);
  const persistent = opts.persistent === undefined ? level === "error" : !!opts.persistent;
  const ms = Number.isFinite(opts.ms) && opts.ms > 0 ? opts.ms : EPHEMERAL_MS;
  return { msg: String(msg), level, persistent, ms };
}

/** warning and error interrupt a screen reader (alert); the rest wait (status). */
export const roleFor = (level) => (level === "warning" || level === "error" ? "alert" : "status");

export function createToastQueue({ limit = LIMIT } = {}) {
  let seq = 0;
  let items = []; // oldest first: { id, spec, left, since, holds:Set }

  const find = (id) => items.find((t) => t.id === id);
  const leftAt = (t, now) => (t.since === null ? t.left : Math.max(0, t.left - (now - t.since)));

  return {
    /**
     * Add a toast (a spec from toastSpec). Returns { id, evicted: [ids] }; id is
     * null when the queue is full of persistent toasts and the new one is
     * ephemeral (a persistent toast is never evicted by an ephemeral one).
     * When full, the oldest ephemeral toast goes first, else the oldest
     * persistent one (for a persistent newcomer).
     */
    add(spec, now) {
      const evicted = [];
      while (items.length >= limit) {
        let victim = items.find((t) => !t.spec.persistent);
        if (!victim && spec.persistent) victim = items[0];
        if (!victim) return { id: null, evicted };
        items = items.filter((t) => t !== victim);
        evicted.push(victim.id);
      }
      const id = ++seq;
      items.push({ id, spec, left: spec.ms, since: spec.persistent ? null : now, holds: new Set() });
      return { id, evicted };
    },
    /** Hold a toast's countdown for a reason ("hover", "focus"); the first hold stops the clock. */
    pause(id, reason, now) {
      const t = find(id);
      if (!t) return;
      if (!t.holds.size && t.since !== null) {
        t.left = leftAt(t, now);
        t.since = null;
      }
      t.holds.add(reason);
    },
    /** Release a hold; the clock restarts with what was left once none remain. */
    resume(id, reason, now) {
      const t = find(id);
      if (!t || !t.holds.delete(reason)) return;
      if (!t.holds.size && !t.spec.persistent) t.since = now;
    },
    paused: (id) => !!find(id)?.holds.size,
    /** Milliseconds left on an ephemeral toast; null for a persistent one or an unknown id. */
    remaining(id, now) {
      const t = find(id);
      return !t || t.spec.persistent ? null : leftAt(t, now);
    },
    /** Ephemeral toasts whose time is up and that no hold keeps. */
    expired(now) {
      return items.filter((t) => !t.spec.persistent && !t.holds.size && leftAt(t, now) <= 0).map((t) => t.id);
    },
    dismiss(id) {
      const n = items.length;
      items = items.filter((t) => t.id !== id);
      return items.length < n;
    },
    /** Id of the newest persistent toast, or null. */
    newestPersistent() {
      for (let i = items.length - 1; i >= 0; i--) if (items[i].spec.persistent) return items[i].id;
      return null;
    },
    clear() {
      const ids = items.map((t) => t.id);
      items = [];
      return ids;
    },
    ids: () => items.map((t) => t.id),
    get: (id) => find(id)?.spec ?? null,
  };
}
