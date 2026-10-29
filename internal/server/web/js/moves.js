// Moves between lanes: what a move still needs, and the API calls that make it.
// Pure: no DOM, no network, no clock. The wizard asks for what the standard
// requires (rules WK5 to WK7) and only for what the record lacks.
import { LANES, CLOSED } from "./logic.js";

/** Targets that always go through the review, even when nothing is missing. */
export const ALWAYS_REVIEWED = ["blocked", "done", "dropped"];
/** Targets whose rule needs a `## Goal` (WK7), and a non-empty `## Next` (WK6). */
export const NEEDS_GOAL = ["ready", "active", "blocked", "done"];
export const NEEDS_NEXT = ["active", "blocked"];

/** The lane next to `status` ("left" or "right"), or null at an end. */
export function neighbour(status, dir) {
  const i = LANES.indexOf(status);
  if (i < 0) return null;
  return LANES[i + (dir === "left" ? -1 : 1)] || null;
}

/**
 * The `##` sections of a body as { lower-case name: text }, text trimmed.
 * A heading inside a code fence is not a heading; `###` does not split.
 */
export function sectionsOf(body) {
  const out = {};
  let name = null;
  let buf = [];
  let fence = "";
  const flush = () => {
    if (name !== null) out[name] = buf.join("\n").trim();
  };
  for (const line of String(body || "").split("\n")) {
    const f = /^\s{0,3}(`{3,}|~{3,})/.exec(line);
    if (f) {
      if (!fence) fence = f[1][0];
      else if (f[1][0] === fence) fence = "";
    }
    const h = !fence && /^##[ \t]+(.+?)[ \t#]*$/.exec(line.replace(/\r$/, ""));
    if (h) {
      flush();
      name = h[1].trim().toLowerCase();
      buf = [];
    } else if (name !== null) buf.push(line);
  }
  flush();
  return out;
}

/** Lines of free text as a list: trimmed, empty lines dropped, repeats dropped. */
export function parseEntries(text) {
  const seen = new Set();
  const out = [];
  for (const l of String(text || "").split(/\r?\n/)) {
    const t = l.trim();
    if (t && !seen.has(t)) {
      seen.add(t);
      out.push(t);
    }
  }
  return out;
}

/**
 * What a move needs. `i`: { from, to, claims, sections, blockedBy }.
 * Returns { kind: "noop" | "safe" | "gated", steps } where steps are the
 * questions to ask, in order. "safe" has no steps and is written at once.
 */
export function planMove(i) {
  if (i.from === i.to || CLOSED.includes(i.from)) return { kind: "noop", steps: [] };
  const has = (n) => !!(i.sections[n] || "").trim();
  const steps = [];
  if (NEEDS_GOAL.includes(i.to) && !has("goal")) steps.push("goal");
  if (NEEDS_NEXT.includes(i.to) && !has("next")) steps.push("next");
  if (i.to === "blocked" && !(i.blockedBy || []).length) steps.push("blocked_by");
  if (CLOSED.includes(i.to) && !has("outcome")) steps.push("outcome");
  if (i.from === "blocked" && !CLOSED.includes(i.to) && i.to !== "blocked" && (i.blockedBy || []).length) steps.push("unblock");
  if (i.to === "done") steps.push("decision");
  const gated = steps.length > 0 || ALWAYS_REVIEWED.includes(i.to);
  if (gated && i.claims && i.to === "active") steps.push("identity");
  return { kind: gated ? "gated" : "safe", steps };
}

/** Every `changes` list that names the ADR made in the same wizard says "$new" for its id. */
export const NEW_ADR = "$new";

const uniq = (xs) => [...new Set(xs)];

/**
 * The writes of a move, in the order they are made. `i` as for planMove plus
 * `decisions` (the record's list now). `a`: the answers { goal, next,
 * blocked_by (text), outcome, identity, decision: { mode: none|link|new, id, slug, title } }.
 */
export function buildOps(i, a) {
  const ops = [];
  const text = (v) => String(v || "").trim();
  if (text(a.goal)) ops.push({ kind: "section", name: "Goal", text: text(a.goal), mode: "replace" });
  if (text(a.next)) ops.push({ kind: "section", name: "Next", text: text(a.next), mode: "replace" });
  const d = a.decision || { mode: "none" };
  if (i.to === "done" && d.mode === "new") ops.push({ kind: "new", what: "adr", slug: d.slug, title: text(d.title) });
  if (i.to === "done" && (d.mode === "link" || d.mode === "new")) {
    ops.push({ kind: "set", changes: { decisions: uniq([...(i.decisions || []), d.mode === "new" ? NEW_ADR : d.id]) } });
  }
  const unblock = i.from === "blocked" && a.unblock !== false && (i.blockedBy || []).length;
  switch (i.to) {
    case "raw":
    case "ready":
      // leaving active with a claim: release it first (claims on)
      if (i.claims && i.from === "active" && i.claimed) ops.push({ kind: "release" });
      ops.push({ kind: "set", changes: unblock ? { status: i.to, blocked_by: null } : { status: i.to } });
      break;
    case "active":
      if (i.claims) {
        // claim sets the status only when the record is ready
        if (i.from !== "ready") ops.push({ kind: "set", changes: unblock ? { status: "active", blocked_by: null } : { status: "active" } });
        ops.push({ kind: "claim", as: text(a.identity) });
      } else ops.push({ kind: "set", changes: unblock ? { status: "active", blocked_by: null } : { status: "active" } });
      break;
    case "blocked": {
      const changes = { status: "blocked" };
      const by = parseEntries(a.blocked_by);
      if (by.length) changes.blocked_by = by;
      ops.push({ kind: "set", changes });
      break;
    }
    case "done":
    case "dropped": {
      const op = { kind: "close", status: i.to };
      if (text(a.outcome)) op.outcome = text(a.outcome);
      ops.push(op);
      break;
    }
  }
  return ops;
}

/** The HTTP call for one op: { method, path, body }. The etag goes in If-Match. */
export function requestFor(op, id, o = {}) {
  const body = {};
  if (o.dry_run) body.dry_run = true;
  const rec = "records/" + encodeURIComponent(id);
  switch (op.kind) {
    case "set":
      return { method: "PATCH", path: rec, body: { changes: op.changes, ...body } };
    case "section":
      return { method: "PUT", path: rec + "/sections/" + encodeURIComponent(op.name), body: { text: op.text, mode: op.mode || "replace", ...body } };
    case "close": {
      const b = { status: op.status, ...body };
      if (op.outcome) b.outcome = op.outcome;
      return { method: "POST", path: rec + "/close", body: b };
    }
    case "claim":
      return { method: "POST", path: rec + "/claim", body: { as: op.as, ...(o.takeover ? { takeover: true } : {}), ...body } };
    case "release":
      return { method: "DELETE", path: rec + "/claim", body: o.dry_run ? body : undefined };
    case "new": {
      const b = { kind: op.what, slug: op.slug, title: op.title };
      if (op.fields) b.fields = op.fields;
      return { method: "POST", path: "records", body: b };
    }
  }
  throw new Error("unknown op " + op.kind);
}

/** Replace the "$new" placeholder in an op with the id of the ADR just made. */
export function resolveOp(op, newId) {
  if (op.kind === "set" && op.changes.decisions && newId) return { ...op, changes: { ...op.changes, decisions: op.changes.decisions.map((x) => (x === NEW_ADR ? newId : x)) } };
  return op;
}

/**
 * Whether a dry run can show this op against the file as it is now. Not the
 * creation of an ADR, not a set that names it, and not a claim that follows a
 * status change (claim refuses a record that is not yet ready or active).
 */
export function previewable(op, i) {
  if (op.kind === "new") return false;
  if (op.kind === "set" && (op.changes.decisions || []).includes(NEW_ADR)) return false;
  if (op.kind === "claim" && i.from !== "ready") return false;
  return true;
}

/** A line for the review: what an op does. */
export function opLabel(op, id) {
  switch (op.kind) {
    case "section":
      return "Write the ## " + op.name + " section";
    case "set": {
      const bits = Object.entries(op.changes).map(([k, v]) => k + " " + (Array.isArray(v) ? "[" + v.join(", ") + "]" : v));
      return "Change " + bits.join(", ");
    }
    case "close":
      return "Close as " + op.status + (op.outcome ? " with an outcome" : "") + " and move the file to docs/work/done/";
    case "claim":
      return "Claim for " + op.as;
    case "release":
      return "Release the claim";
    case "new":
      return "Create the ADR " + op.slug + " (status proposed)";
  }
  return op.kind + " " + id;
}

/** Text of the error shown in a step, "" when the answer is fine. `a` as for buildOps. */
export function validateStep(step, a, to) {
  const t = (v) => String(v || "").trim();
  switch (step) {
    case "goal":
      return t(a.goal) ? "" : "Say what done looks like. A ready, active, blocked or done item needs a Goal.";
    case "next":
      return t(a.next) ? "" : "Write the next step. An active or blocked item needs a Next.";
    case "blocked_by":
      return parseEntries(a.blocked_by).length ? "" : "Name what it waits on: a record id or a few words.";
    case "outcome":
      return t(a.outcome) ? "" : to === "dropped" ? "Give the reason it was dropped." : "Write a short outcome. A closed item needs one.";
    case "identity":
      return t(a.identity) ? "" : "Say who is claiming it.";
    case "decision": {
      const d = a.decision || { mode: "none" };
      if (d.mode === "link" && !d.id) return "Pick the ADR it decided, or choose No.";
      if (d.mode === "new" && !t(d.title)) return "Give the new ADR a title.";
      if (d.mode === "new" && !d.slug) return "The title has no letters or digits for a name.";
      return "";
    }
  }
  return "";
}

/**
 * What a key does in the wizard: "back" (Esc), "next", "write", "reload" or "".
 * k: { key, mod, tag, type, phase, step }. Enter advances a step, except where it
 * means something else: a textarea (a new line), a select, a button (it presses itself).
 * Mod+Enter advances from anywhere and never writes; Enter writes on the review.
 */
export function wizardKey(k) {
  if (k.key === "Escape") return "back";
  if (k.key !== "Enter") return "";
  const review = k.step === "review";
  if (k.mod) return review ? "" : "next";
  if (k.tag === "TEXTAREA" || k.tag === "SELECT" || k.tag === "BUTTON") return "";
  if (review) return k.phase === "conflict" ? "reload" : k.phase === "review" ? "write" : "";
  return "next";
}
