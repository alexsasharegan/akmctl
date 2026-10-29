// Moves, undo, capture and drag and drop on the board. Safe moves write at
// once with the record's etag and offer `u`; every other move goes through the
// wizard and the review (wizard.js). Nothing here builds HTML.
import { el, ageDays, ageText, text } from "./util.js";
import { api, ApiError, state, show, describeError } from "./store.js";
import { field, blockedBy, CLOSED, isOpenWork, obsidianUrl } from "./logic.js";
import { neighbour, sectionsOf, planMove, buildOps, requestFor, validateStep } from "./moves.js";
import { openWizard } from "./wizard.js";
import { openCapture } from "./capture.js";
import { openDialog, closeOverlay, setDialogHints, toast, MOD, overlayOpen } from "./ui.js";
import { remember } from "./board.js";

const IDENTITY_KEY = "akm-identity";

let H = null; // { host, reload(), byId(), records(), focusView() }
let last = null; // what `u` undoes: { kind: "move"|"capture", ... }
let claimsP = null;
let busy = false;

const deps = { openDialog, closeDialog: closeOverlay, setDialogHints, MOD };

/** Start the flow. h: host (the view element), reload(): Promise, byId(), records(), focusView(). */
export function initFlow(h) {
  H = h;
  dragSetup(h.host);
}

// ---- identity and claims ----

export function getIdentity() {
  let id = "";
  try {
    id = localStorage.getItem(IDENTITY_KEY) || "";
  } catch (e) {
    id = "";
  }
  return id || state.config.identity || "";
}

export function rememberIdentity(name) {
  try {
    localStorage.setItem(IDENTITY_KEY, name);
  } catch (e) {
    /* kept for this page only */
  }
  state.config = { ...state.config, identity: name };
}

/** Whether the repository's manifest enables claims (from `version`), read once. */
function claimsOn() {
  if (!claimsP) {
    claimsP = api("GET", "version")
      .then((v) => ((v.repository && v.repository.extensions) || []).includes("claims"))
      .catch(() => {
        claimsP = null;
        return false;
      });
  }
  return claimsP;
}

// ---- helpers ----

const arr = (v) => (Array.isArray(v) ? v : []);
const titleOf = (r) => field(r, "title") || r.id;

async function refresh(id) {
  if (id) remember(id);
  await H.reload();
  H.focusView();
}

function workRecord(id) {
  const rec = H.byId().get(id);
  return rec && rec.type === "work" ? rec : null;
}

// ---- the move ----

/** Move a card to a lane: the one path for H, L, d, x, b and drag and drop. */
export async function moveTo(id, to, extra = {}) {
  if (busy || overlayOpen()) return;
  const rec = workRecord(id);
  if (!rec) return toast("Only a work item can be moved.", "bad");
  const from = field(rec, "status");
  if (from === to) return toast(id + " is already " + to + ".");
  if (CLOSED.includes(from)) return toast(id + " is closed; a closed record is not moved.", "bad");
  busy = true;
  try {
    let data;
    try {
      data = await show(id);
    } catch (e) {
      return toast(describeError(e), "bad");
    }
    const claims = await claimsOn();
    const move = { id, from, to, claims, sections: sectionsOf(data.body), blockedBy: blockedBy(data.record), claimed: !!field(data.record, "claimed_by"), decisions: arr(data.record.fields.decisions) };
    const plan = planMove(move);
    if (plan.kind === "noop") return;
    if (plan.kind === "safe") return await safeMove(rec, move);
    last = null;
    openWizard(deps, {
      rec,
      title: titleOf(rec),
      i: move,
      plan,
      identity: getIdentity(),
      adrs: H.records().filter((r) => r.type === "adr").map((r) => ({ id: r.id, title: titleOf(r), status: field(r, "status") })),
      outcomeDefault: extra.outcomeDefault,
      onReload: () => refresh(id),
      onDone: (res) => {
        if (res.identity && move.claims && to === "active") rememberIdentity(res.identity);
        toast(doneText(rec, to, res));
        refresh(id);
      },
    });
  } finally {
    busy = false;
  }
}

function doneText(rec, to, res) {
  const t = titleOf(rec);
  const adr = res.newAdr ? " ADR " + res.newAdr + " is proposed; run akmctl index to refresh the index." : "";
  if (to === "done") return "Closed " + t + " as done." + adr;
  if (to === "dropped") return "Dropped " + t + ".";
  return "Moved " + t + " to " + to + "." + adr;
}

async function safeMove(rec, move) {
  const identity = getIdentity();
  if (move.claims && move.to === "active" && !identity) {
    // no name to claim as: ask for it in the wizard instead of guessing
    const plan = { kind: "gated", steps: ["identity"] };
    openWizard(deps, { rec, title: titleOf(rec), i: move, plan, identity, adrs: [], onReload: () => refresh(rec.id), onDone: (res) => (rememberIdentity(res.identity), toast("Claimed " + titleOf(rec) + "."), refresh(rec.id)) });
    return;
  }
  const ops = buildOps(move, { identity });
  let etag = rec.etag;
  let r = null;
  try {
    for (const op of ops) {
      const q = requestFor(op, rec.id);
      r = await api(q.method, q.path, { body: q.body, etag });
      etag = r.record.etag;
    }
  } catch (e) {
    if (e instanceof ApiError && e.code === "stale") toast(titleOf(rec) + " changed since it was loaded. Nothing was written; the records were reloaded.", "bad");
    else toast(describeError(e), "bad");
    await refresh(rec.id);
    return;
  }
  const claimed = ops.some((o) => o.kind === "claim");
  last = { kind: "move", id: rec.id, from: move.from, to: move.to, etag, claimed };
  await refresh(rec.id);
  undoToast("Moved " + titleOf(rec) + " from " + move.from + " to " + move.to + (claimed ? ", claimed as " + identity : "") + ".");
}

function undoToast(msg, label = "Undo") {
  const box = toast(msg + " ", "", 12000);
  box.append(el("button", { type: "button", class: "btn quiet toast-btn", onclick: () => (box.remove(), undo()) }, label + " ", el("kbd", { text: "u" })));
}

/** Undo the last safe move or capture. A capture is dropped with a reason: there is no delete. */
export async function undo() {
  if (busy || overlayOpen()) return;
  const u = last;
  if (!u) return toast("Nothing to undo.");
  for (const b of document.querySelectorAll("#toasts .toast-btn")) b.closest(".toast").remove();
  if (u.kind === "capture") {
    last = null;
    return moveTo(u.id, "dropped", { outcomeDefault: "Captured by mistake." });
  }
  busy = true;
  last = null;
  try {
    let etag = u.etag;
    if (u.claimed) {
      const q = requestFor({ kind: "release" }, u.id);
      etag = (await api(q.method, q.path, { etag })).record.etag;
    }
    const q = requestFor({ kind: "set", changes: { status: u.from } }, u.id);
    const r = await api(q.method, q.path, { body: q.body, etag });
    await refresh(u.id);
    toast("Moved " + u.id + " back to " + u.from + ". The updated date was not restored: it now reads " + (r.record.fields.updated || "today") + (u.claimed ? ", and the claim stays in the log" : "") + ".", "", 7000);
  } catch (e) {
    if (e instanceof ApiError && e.code === "stale") toast(u.id + " changed after the move, so it was not undone.", "bad");
    else toast("Undo failed: " + describeError(e), "bad");
    await refresh(u.id);
  } finally {
    busy = false;
  }
}

// ---- capture ----

export function capture() {
  if (busy || overlayOpen()) return;
  const parents = H.records().filter(isOpenWork).map((r) => ({ id: r.id, title: titleOf(r) }));
  openCapture(deps, {
    parents,
    onDone: (rec) => {
      last = { kind: "capture", id: rec.id };
      refresh(rec.id).then(() => undoToast("Captured " + titleOf(rec) + " as raw.", "Drop it"));
    },
  });
}

// ---- keys, palette actions, help ----

/**
 * A key on the board or in the record drawer. `id` is the focused or open
 * record. Returns true when the key was used.
 */
/** Open a record in Obsidian; false when the root is unknown (nothing to do). */
export function openInObsidian(id) {
  const r = H.byId().get(id);
  const url = r ? obsidianUrl(state.config, r.path) : "";
  if (!url) return false;
  window.location.href = url;
  return true;
}

export function flowKeys(e, id) {
  if (e.metaKey || e.ctrlKey || e.altKey) return false;
  switch (e.key) {
    case "c":
      capture();
      return true;
    case "u":
      undo();
      return true;
  }
  if (!id) return false;
  const rec = workRecord(id);
  switch (e.key) {
    case "o":
      return openInObsidian(id);
    case "H":
    case "L": {
      const to = rec && neighbour(field(rec, "status"), e.key === "H" ? "left" : "right");
      if (!rec) return false;
      if (!to) toast("Already in the " + (e.key === "H" ? "first" : "last") + " lane.");
      else moveTo(id, to);
      return true;
    }
    case "d":
      if (rec) moveTo(id, "done");
      return !!rec;
    case "x":
      if (rec) moveTo(id, "dropped");
      return !!rec;
    case "b":
      if (rec) moveTo(id, "blocked");
      return !!rec;
  }
  return false;
}

export function flowActions(currentId) {
  const cur = () => currentId();
  const on = (fn) => () => {
    const id = cur();
    if (id) fn(id);
    else toast("Focus a card first.");
  };
  return [
    { label: "Move the focused card one lane left", keys: ["H"], run: on((id) => flowKeys({ key: "H" }, id)) },
    { label: "Move the focused card one lane right", keys: ["L"], run: on((id) => flowKeys({ key: "L" }, id)) },
    { label: "Close the focused item as done", keys: ["d"], run: on((id) => moveTo(id, "done")) },
    { label: "Drop the focused item", keys: ["x"], run: on((id) => moveTo(id, "dropped")) },
    { label: "Block the focused item", keys: ["b"], run: on((id) => moveTo(id, "blocked")) },
    { label: "Capture a new raw item", keys: ["c"], run: capture },
    { label: "Undo the last move or capture", keys: ["u"], run: undo },
    { label: "Set the name used for claims", keys: [], run: setIdentity },
  ];
}

export const flowHelp = () => [
  "Moves",
  [
    ["H L", "move the card one lane left, right (writes at once, u undoes)"],
    ["d", "close as done (asks what is missing, then reviews)"],
    ["x", "drop (asks for the reason, then reviews)"],
    ["b", "block (asks what it waits on, then reviews)"],
    ["u", "undo the last move; after a capture, drop it with a reason"],
    ["c", "capture a new raw item"],
    [MOD + "+Enter", "next step in the wizard"],
    ["Enter", "write, on the review"],
    ["Esc", "back, or cancel on the first step"],
    ["drag", "drop a card on a lane: same as H and L"],
  ],
];

/** A small dialog to change the name claims are made under. */
export function setIdentity() {
  if (overlayOpen()) return;
  const input = el("input", { id: "id-name", type: "text", autocomplete: "off", value: getIdentity() });
  const err = el("p", { class: "err", role: "alert", hidden: true });
  const save = () => {
    const name = input.value.trim();
    const e = validateStep("identity", { identity: name });
    if (e) {
      err.hidden = false;
      err.textContent = e;
      return;
    }
    rememberIdentity(name);
    closeOverlay();
    toast("Claims will be made as " + name + ".");
  };
  const dlg = el(
    "div",
    { class: "dialog capture", role: "dialog", "aria-modal": "true", "aria-labelledby": "id-t" },
    el("header", { class: "dlg-head" }, el("div", {}, el("p", { class: "eyebrow", text: "Claims" }), el("h2", { id: "id-t", text: "Claim as" }))),
    el("div", { class: "dlg-body" }, el("label", { class: "field", for: "id-name" }, el("span", { class: "label", text: "Your name" }), el("span", { class: "help", text: "Used when a move to active claims the item. Remembered in this browser; the default is the operating system user." }), input), err),
    el("footer", { class: "dlg-foot" }, el("button", { type: "button", class: "btn quiet", onclick: () => closeOverlay() }, "Cancel ", el("kbd", { text: "Esc" })), el("button", { type: "button", class: "btn primary", onclick: save }, "Save ", el("kbd", { text: "Enter" }))),
  );
  openDialog(dlg, {
    mode: "IDENTITY",
    hints: [["Enter", "save"], ["Esc", "cancel"]],
    onKey(e) {
      if (e.key === "Escape") {
        e.preventDefault();
        closeOverlay();
      } else if (e.key === "Enter") {
        e.preventDefault();
        save();
      }
    },
  });
  input.focus();
  input.select();
}

// ---- drag and drop ----
// A card dropped on a lane takes the same path as H and L. Delegated on the
// view element, so it survives the board being drawn again.

function dragSetup(host) {
  let dragId = null;
  const cardOf = (t) => t.closest && t.closest(".card[data-id]");
  const laneOf = (t) => t.closest && t.closest(".col");
  const laneStatus = (col) => ["raw", "ready", "active", "blocked"].find((s) => col.classList.contains("s-" + s));
  const clear = () => {
    for (const x of host.querySelectorAll(".dragging, .drop-target")) x.classList.remove("dragging", "drop-target");
  };
  host.addEventListener("dragstart", (e) => {
    const c = cardOf(e.target);
    if (!c) return;
    dragId = c.dataset.id;
    e.dataTransfer.setData("text/plain", dragId);
    e.dataTransfer.effectAllowed = "move";
    c.classList.add("dragging");
  });
  host.addEventListener("dragend", () => {
    dragId = null;
    clear();
  });
  host.addEventListener("dragover", (e) => {
    const col = laneOf(e.target);
    if (!dragId || !col) return;
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
    if (!col.classList.contains("drop-target")) {
      for (const x of host.querySelectorAll(".drop-target")) x.classList.remove("drop-target");
      col.classList.add("drop-target");
    }
  });
  host.addEventListener("drop", (e) => {
    const col = laneOf(e.target);
    if (!dragId || !col) return;
    e.preventDefault();
    const id = dragId;
    dragId = null;
    clear();
    const to = laneStatus(col);
    if (to) moveTo(id, to);
  });
}

// ---- claim controls in the drawer ----

/**
 * The claim line and buttons for a record shown in the drawer, or null when the
 * repository does not enable claims. `r` is the full record from `show`;
 * `reload()` redraws the drawer after a write.
 */
export async function claimControls(r, reload) {
  if (r.type !== "work" || !(await claimsOn())) return null;
  const f = r.fields || {};
  const holder = f.claimed_by ? text(f.claimed_by) : "";
  const me = getIdentity();
  const status = text(f.status);
  const run = async (what, op, o) => {
    try {
      const q = requestFor(op, r.id, o);
      await api(q.method, q.path, { body: q.body, etag: r.etag });
      toast(what + " " + r.id + ".");
    } catch (e) {
      toast(describeError(e), "bad");
    }
    await refresh(r.id);
    reload();
  };
  const btn = (label, fn) => el("button", { type: "button", class: "btn quiet", onclick: fn }, label);
  const kids = [];
  if (holder) {
    const age = f.claimed_at ? ageText(ageDays(String(f.claimed_at).slice(0, 10))) : "";
    kids.push(el("span", { class: "claim-holder", text: "Claimed by " + holder + (age ? ", " + age : "") }));
    if (holder === me) kids.push(btn("Release", () => run("Released", { kind: "release" })));
    else if (me) kids.push(btn("Take over", () => (window.confirm("Take over " + r.id + " from " + holder + "?") ? run("Took over", { kind: "claim", as: me }, { takeover: true }) : null)));
  } else if ((status === "ready" || status === "active") && me) {
    kids.push(el("span", { class: "claim-holder", text: "Not claimed" }), btn("Claim as " + me, () => run("Claimed", { kind: "claim", as: me })));
  } else if (!me) kids.push(el("span", { class: "claim-holder", text: "Set a name for claims from the palette." }));
  return el("div", { class: "claim-bar" }, kids);
}
